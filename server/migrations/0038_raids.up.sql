-- 0038_raids.up.sql
-- Raids (DD-261): group dungeons for a party. A raid template is a dungeon
-- template whose group_max is above one; its clears share
-- dungeon.weekly_lockout (contract_id = the member row).
--
--   raid.run      one raid of a party: forming (lobby) → running → done
--   raid.member   a joined character, its frozen combat snapshot and its own
--                 reward collection

ALTER TABLE content.dungeon_template
  ADD COLUMN group_min SMALLINT NOT NULL DEFAULT 1,
  ADD COLUMN group_max SMALLINT NOT NULL DEFAULT 1,
  ADD CONSTRAINT dungeon_template_group_check
    CHECK (group_min >= 1 AND group_max >= group_min AND group_max <= 5);

-- Raid bosses are fought only in raids: hunts refuse them.
ALTER TABLE content.mob_template ADD COLUMN raid_boss BOOLEAN NOT NULL DEFAULT false;

CREATE SCHEMA IF NOT EXISTS raid;

CREATE TABLE raid.run (
  id                   UUID PRIMARY KEY,
  party_id             UUID NOT NULL,         -- the party it was formed in (no FK: parties are disbanded, not deleted)
  raid_code            VARCHAR(48) NOT NULL,
  leader_character_id  UUID NOT NULL,         -- who may start or call it off
  state                VARCHAR(16) NOT NULL CHECK (state IN ('forming', 'running', 'cancelled')),
  duration_seconds     INT NOT NULL CHECK (duration_seconds > 0),
  seed                 BIGINT NOT NULL DEFAULT 0,
  started_at           TIMESTAMPTZ NULL,
  expires_at           TIMESTAMPTZ NULL,
  created_at           TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at           TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- One open lobby per party.
CREATE UNIQUE INDEX ux_raid_run_forming ON raid.run (party_id) WHERE state = 'forming';

CREATE TABLE raid.member (
  id             UUID PRIMARY KEY,            -- item grant + ledger source of this member's rewards
  run_id         UUID NOT NULL REFERENCES raid.run(id) ON DELETE CASCADE,
  character_id   UUID NOT NULL REFERENCES character.character(id) ON DELETE CASCADE,
  slot           SMALLINT NOT NULL,           -- fight order
  snapshot       JSONB NOT NULL,
  joined_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
  active         BOOLEAN NOT NULL DEFAULT true, -- false once collected or called off
  revealed_at    TIMESTAMPTZ NULL,
  reveal_result  JSONB NULL,
  reveal_txn_id  UUID NULL,
  UNIQUE (run_id, character_id),
  UNIQUE (run_id, slot)
);

-- A character is in at most one raid until it collects the rewards.
CREATE UNIQUE INDEX ux_raid_member_active ON raid.member (character_id) WHERE active;
CREATE INDEX idx_raid_member_character ON raid.member (character_id, joined_at DESC);
