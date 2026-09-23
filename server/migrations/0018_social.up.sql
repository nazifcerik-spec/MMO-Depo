-- 0018_social.up.sql
-- Social (Sprint 21): friends, block list, mail, chat.
-- docs/social/*, docs/database/social-guilds-parties.md §12-17.

-- --- Friends ---------------------------------------------------------------

CREATE TABLE social.friendship (
  character_a  UUID NOT NULL REFERENCES character.character(id) ON DELETE CASCADE,
  character_b  UUID NOT NULL REFERENCES character.character(id) ON DELETE CASCADE,
  since        TIMESTAMPTZ NOT NULL DEFAULT now(),
  PRIMARY KEY (character_a, character_b),
  CHECK (character_a < character_b)
);
CREATE INDEX idx_friendship_b ON social.friendship (character_b);

CREATE TABLE social.friend_request (
  id                     UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  requester_character_id UUID NOT NULL REFERENCES character.character(id) ON DELETE CASCADE,
  target_character_id    UUID NOT NULL REFERENCES character.character(id) ON DELETE CASCADE,
  status                 VARCHAR(16) NOT NULL DEFAULT 'pending'
                         CHECK (status IN ('pending','accepted','declined','cancelled')),
  created_at             TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE UNIQUE INDEX ux_friend_request_pending
  ON social.friend_request (requester_character_id, target_character_id) WHERE status = 'pending';
CREATE INDEX idx_friend_request_target
  ON social.friend_request (target_character_id) WHERE status = 'pending';

-- --- Block list ------------------------------------------------------------

CREATE TABLE social.ignore_list (
  owner_character_id  UUID NOT NULL REFERENCES character.character(id) ON DELETE CASCADE,
  target_character_id UUID NOT NULL REFERENCES character.character(id) ON DELETE CASCADE,
  added_at            TIMESTAMPTZ NOT NULL DEFAULT now(),
  PRIMARY KEY (owner_character_id, target_character_id)
);
CREATE INDEX idx_ignore_target ON social.ignore_list (target_character_id);

-- --- Mail ------------------------------------------------------------------

CREATE TABLE social.mail (
  id                      UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  sender_character_id     UUID NULL REFERENCES character.character(id) ON DELETE SET NULL,
  recipient_character_id  UUID NOT NULL REFERENCES character.character(id) ON DELETE CASCADE,
  subject                 VARCHAR(64) NOT NULL,
  body                    VARCHAR(2000) NOT NULL DEFAULT '',
  gold_attachment_copper  BIGINT NOT NULL DEFAULT 0 CHECK (gold_attachment_copper >= 0),
  item_attachment_count   SMALLINT NOT NULL DEFAULT 0,
  postage_copper          BIGINT NOT NULL DEFAULT 0,
  send_txn_id             UUID NULL UNIQUE,    -- ledger txn of postage+gold escrow (replay guard)
  state                   VARCHAR(16) NOT NULL DEFAULT 'sent'
                          CHECK (state IN ('sent','read','claimed','returned','expired')),
  returned_mail_id        UUID NULL,           -- set on the return mail: the original
  created_at              TIMESTAMPTZ NOT NULL DEFAULT now(),
  read_at                 TIMESTAMPTZ NULL,
  attachments_claimed_at  TIMESTAMPTZ NULL,
  expires_at              TIMESTAMPTZ NOT NULL DEFAULT now() + INTERVAL '30 days',
  deleted_at              TIMESTAMPTZ NULL
);
CREATE INDEX idx_mail_inbox
  ON social.mail (recipient_character_id, created_at DESC) WHERE deleted_at IS NULL;
CREATE INDEX idx_mail_expires
  ON social.mail (expires_at) WHERE state IN ('sent','read');

CREATE TABLE social.mail_attachment (
  mail_id           UUID NOT NULL REFERENCES social.mail(id) ON DELETE CASCADE,
  slot_index        SMALLINT NOT NULL CHECK (slot_index BETWEEN 0 AND 11),
  item_instance_id  UUID NOT NULL REFERENCES inventory.item_instance(id) ON DELETE CASCADE,
  PRIMARY KEY (mail_id, slot_index),
  UNIQUE (item_instance_id)
);

-- --- Chat ------------------------------------------------------------------

CREATE TABLE social.chat_message (
  id                   BIGSERIAL PRIMARY KEY,
  channel_kind         VARCHAR(16) NOT NULL
                       CHECK (channel_kind IN ('global','guild','party','whisper')),
  channel_ref_id       UUID NULL,          -- guild_id / party_id
  sender_character_id  UUID NOT NULL REFERENCES character.character(id) ON DELETE CASCADE,
  target_character_id  UUID NULL REFERENCES character.character(id) ON DELETE CASCADE,
  body                 VARCHAR(512) NOT NULL,
  sent_at              TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX idx_chat_channel ON social.chat_message (channel_kind, channel_ref_id, sent_at DESC);
CREATE INDEX idx_chat_whisper ON social.chat_message (sender_character_id, target_character_id, sent_at DESC)
  WHERE channel_kind = 'whisper';
CREATE INDEX idx_chat_sent ON social.chat_message USING BRIN (sent_at);
