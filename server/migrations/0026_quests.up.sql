-- 0026_quests.up.sql
-- Quests v1 (Sprint 41, DD-221): main / side / daily templates from the
-- content pack; objectives are progression event kinds (mob_killed,
-- mob_killed:<mob>, gather_<kind>, dungeon_cleared[:<code>], boss_killed,
-- item_crafted[:<template>]). docs/quests/*.

CREATE TABLE content.quest_template (
  code                VARCHAR(64) PRIMARY KEY,
  display_name_key    VARCHAR(96) NOT NULL,
  category            VARCHAR(16) NOT NULL CHECK (category IN ('main','side','daily')),
  min_level           SMALLINT NOT NULL DEFAULT 1 CHECK (min_level BETWEEN 1 AND 60),
  requires_quest_code VARCHAR(64) NULL,        -- must be turned in first (main chain)
  objectives          JSONB NOT NULL,          -- [{event, count}]
  reward_xp           INT NOT NULL DEFAULT 0 CHECK (reward_xp >= 0),
  reward_gold_copper  BIGINT NOT NULL DEFAULT 0 CHECK (reward_gold_copper >= 0),
  sort_order          INT NOT NULL DEFAULT 0,
  version             INT NOT NULL DEFAULT 1
);

CREATE TABLE progression.character_quest (
  id             UUID PRIMARY KEY,
  character_id   UUID NOT NULL REFERENCES character.character(id) ON DELETE CASCADE,
  quest_code     VARCHAR(64) NOT NULL REFERENCES content.quest_template(code),
  state          VARCHAR(16) NOT NULL CHECK (state IN ('active','completed','turned_in','abandoned','expired')),
  progress       JSONB NOT NULL,               -- [int] parallel to objectives
  day_key        DATE NULL,                    -- UTC day for dailies; NULL otherwise
  accepted_at    TIMESTAMPTZ NOT NULL,
  completed_at   TIMESTAMPTZ NULL,
  closed_at      TIMESTAMPTZ NULL,             -- turned in, abandoned or expired (past-day daily)
  reward_txn_id  UUID NULL
);

-- One open instance per quest.
CREATE UNIQUE INDEX uq_character_quest_open ON progression.character_quest (character_id, quest_code)
  WHERE state IN ('active','completed');
-- Main/side quests are turned in once.
CREATE UNIQUE INDEX uq_character_quest_done ON progression.character_quest (character_id, quest_code)
  WHERE state = 'turned_in' AND day_key IS NULL;
-- A daily can be taken once per UTC day (abandoning does not re-roll it).
CREATE UNIQUE INDEX uq_character_quest_daily ON progression.character_quest (character_id, quest_code, day_key)
  WHERE day_key IS NOT NULL;
CREATE INDEX idx_character_quest_open ON progression.character_quest (character_id)
  WHERE state IN ('active','completed');
