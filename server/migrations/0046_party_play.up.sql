-- 0046_party_play.up.sql
-- Party play (DD-296..299): XP shares among party members doing the same
-- kind of work, and finding a raid party (queue + listing board).

-- One row per (job, partner): what a member received from another's job.
-- The primary key pays each share once.
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

-- Raid queue: characters waiting to be put into a random party for a raid.
CREATE TABLE raid.queue (
  character_id UUID PRIMARY KEY REFERENCES character.character(id) ON DELETE CASCADE,
  account_id   UUID NOT NULL,
  raid_code    VARCHAR(48) NOT NULL,
  queued_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX idx_raid_queue_raid ON raid.queue (raid_code, queued_at);

-- Raid listings: a party advertising for members; anyone may join.
CREATE TABLE raid.listing (
  id         UUID PRIMARY KEY,
  party_id   UUID NOT NULL UNIQUE REFERENCES social.party(id) ON DELETE CASCADE,
  raid_code  VARCHAR(48) NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX idx_raid_listing_raid ON raid.listing (raid_code, created_at);
