-- 0017_party_guild.up.sql
-- Party + Guild (Sprint 20). docs/party/*, docs/guild/*,
-- docs/database/social-guilds-parties.md §1-10.
--
-- v1 scope:
--   party: ephemeral, size cap 5, leader-only invite, 60s invites, loot 'personal'.
--   guild: 100k copper creation fee (ledger), fixed 5-rank matrix in code,
--          48h scheduled disband (worker `guild_disband_expire` finalises).
-- Out of scope: treasury / bank / research / war / LFG (later sprints).

-- --- Party -----------------------------------------------------------------

CREATE TABLE social.party (
  id                   UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  leader_character_id  UUID NOT NULL REFERENCES character.character(id) ON DELETE CASCADE,
  size_cap             SMALLINT NOT NULL DEFAULT 5 CHECK (size_cap BETWEEN 1 AND 5),
  loot_rule            VARCHAR(24) NOT NULL DEFAULT 'personal',
  created_at           TIMESTAMPTZ NOT NULL DEFAULT now(),
  disbanded_at         TIMESTAMPTZ NULL,
  version              INT NOT NULL DEFAULT 1
);

CREATE TABLE social.party_member (
  party_id      UUID NOT NULL REFERENCES social.party(id) ON DELETE CASCADE,
  character_id  UUID NOT NULL REFERENCES character.character(id) ON DELETE CASCADE,
  role_hint     VARCHAR(16) NULL,
  joined_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
  left_at       TIMESTAMPTZ NULL,
  PRIMARY KEY (party_id, character_id)
);

-- A character is in at most one live party.
CREATE UNIQUE INDEX ux_party_member_active
  ON social.party_member (character_id) WHERE left_at IS NULL;
CREATE INDEX idx_party_member_party
  ON social.party_member (party_id) WHERE left_at IS NULL;

CREATE TABLE social.party_invite (
  id                      UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  party_id                UUID NOT NULL REFERENCES social.party(id) ON DELETE CASCADE,
  target_character_id     UUID NOT NULL REFERENCES character.character(id) ON DELETE CASCADE,
  invited_by_character_id UUID NOT NULL REFERENCES character.character(id) ON DELETE CASCADE,
  status                  VARCHAR(16) NOT NULL DEFAULT 'pending'
                          CHECK (status IN ('pending','accepted','declined','expired')),
  expires_at              TIMESTAMPTZ NOT NULL,
  created_at              TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX ux_party_invite_pending
  ON social.party_invite (party_id, target_character_id) WHERE status = 'pending';
CREATE INDEX idx_party_invite_target
  ON social.party_invite (target_character_id) WHERE status = 'pending';

-- --- Guild -----------------------------------------------------------------

CREATE TABLE social.guild (
  id                    UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name                  VARCHAR(24) NOT NULL CHECK (char_length(name) BETWEEN 3 AND 24),
  name_normalized       VARCHAR(24) NOT NULL,
  tag                   VARCHAR(5)  NOT NULL CHECK (char_length(tag) BETWEEN 2 AND 5),
  tag_normalized        VARCHAR(5)  NOT NULL,
  leader_character_id   UUID NOT NULL REFERENCES character.character(id) ON DELETE RESTRICT,
  motd                  VARCHAR(256) NOT NULL DEFAULT '',
  description           VARCHAR(512) NOT NULL DEFAULT '',
  level                 SMALLINT NOT NULL DEFAULT 1,
  xp                    BIGINT NOT NULL DEFAULT 0,
  member_cap            SMALLINT NOT NULL DEFAULT 50,
  creation_txn_id       UUID NULL,                    -- ledger txn of the creation fee
  created_at            TIMESTAMPTZ NOT NULL DEFAULT now(),
  disband_scheduled_at  TIMESTAMPTZ NULL,
  disbanded_at          TIMESTAMPTZ NULL,
  version               INT NOT NULL DEFAULT 1
);

-- Names/tags are free again once a guild is disbanded.
CREATE UNIQUE INDEX ux_guild_name ON social.guild (name_normalized) WHERE disbanded_at IS NULL;
CREATE UNIQUE INDEX ux_guild_tag  ON social.guild (tag_normalized)  WHERE disbanded_at IS NULL;
CREATE INDEX idx_guild_disband_due
  ON social.guild (disband_scheduled_at)
  WHERE disband_scheduled_at IS NOT NULL AND disbanded_at IS NULL;

CREATE TABLE social.guild_member (
  guild_id       UUID NOT NULL REFERENCES social.guild(id) ON DELETE CASCADE,
  character_id   UUID NOT NULL REFERENCES character.character(id) ON DELETE CASCADE,
  rank_index     SMALLINT NOT NULL DEFAULT 4 CHECK (rank_index BETWEEN 0 AND 4),
  joined_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
  PRIMARY KEY (guild_id, character_id),
  UNIQUE (character_id)
);

CREATE INDEX idx_guild_member_guild ON social.guild_member (guild_id, rank_index);

CREATE TABLE social.guild_invite (
  id                      UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  guild_id                UUID NOT NULL REFERENCES social.guild(id) ON DELETE CASCADE,
  target_character_id     UUID NOT NULL REFERENCES character.character(id) ON DELETE CASCADE,
  invited_by_character_id UUID NOT NULL REFERENCES character.character(id) ON DELETE CASCADE,
  status                  VARCHAR(16) NOT NULL DEFAULT 'pending'
                          CHECK (status IN ('pending','accepted','declined','expired')),
  expires_at              TIMESTAMPTZ NOT NULL,
  created_at              TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX ux_guild_invite_pending
  ON social.guild_invite (guild_id, target_character_id) WHERE status = 'pending';
CREATE INDEX idx_guild_invite_target
  ON social.guild_invite (target_character_id) WHERE status = 'pending';
