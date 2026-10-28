CREATE TABLE social.party_xp_share (
  source_ref             UUID NOT NULL,
  recipient_character_id UUID NOT NULL REFERENCES character.character(id) ON DELETE CASCADE,
  giver_character_id     UUID NOT NULL REFERENCES character.character(id) ON DELETE CASCADE,
  work_kind              VARCHAR(48) NOT NULL,
  xp                     BIGINT NOT NULL CHECK (xp > 0),
  created_at             TIMESTAMPTZ NOT NULL DEFAULT now(),
  PRIMARY KEY (source_ref, recipient_character_id)
);
CREATE INDEX idx_party_xp_share_recipient ON social.party_xp_share (recipient_character_id, created_at DESC);
