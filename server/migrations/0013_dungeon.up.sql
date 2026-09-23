-- 0013_dungeon.up.sql
-- Dungeon Engine — template + contract + weekly lockout (Sprint 15).
--
-- Three tables:
--   content.dungeon_template   yaşam boyu değişmez (content pack)
--   dungeon.contract           run state; slot cap 1 per character
--   dungeon.weekly_lockout     bir dungeon'ı bir haftada bir kez temizleyebilirsin
--
-- Lockout key = ISO year-week (UTC). Sprint 15 v1 sabit hafta; Phase 02'de
-- dungeon-per-template rollover ayarı (Perşembe reset, sezon takvimi) DD ile
-- gelir.

CREATE SCHEMA IF NOT EXISTS dungeon;

-- --- Template ---------------------------------------------------------------

CREATE TABLE content.dungeon_template (
  code               VARCHAR(48) PRIMARY KEY,
  display_name_key   VARCHAR(64) NOT NULL,
  min_level          SMALLINT NOT NULL DEFAULT 1,
  duration_seconds   INT NOT NULL CHECK (duration_seconds > 0 AND duration_seconds <= 3600),

  -- Waves shape: [{"mob_code":"grey_wolf","count":3}, ...]
  waves              JSONB NOT NULL,
  boss_mob_code      VARCHAR(48) NOT NULL,

  -- Bonus rewards on boss kill (on top of accumulated encounter loot).
  bonus_gold_copper  BIGINT NOT NULL DEFAULT 0,
  bonus_xp           INT NOT NULL DEFAULT 0,

  version            INT NOT NULL DEFAULT 1,
  created_at         TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_dungeon_template_level ON content.dungeon_template (min_level);

-- --- Contract ---------------------------------------------------------------

CREATE TABLE dungeon.contract (
  id                 UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  character_id       UUID NOT NULL REFERENCES character.character(id) ON DELETE CASCADE,
  dungeon_code       VARCHAR(48) NOT NULL,
  duration_seconds   INT NOT NULL CHECK (duration_seconds > 0 AND duration_seconds <= 3600),
  started_at         TIMESTAMPTZ NOT NULL DEFAULT now(),
  expires_at         TIMESTAMPTZ NOT NULL,
  seed               BIGINT NOT NULL,
  state              VARCHAR(16) NOT NULL DEFAULT 'active'
                     CHECK (state IN ('active','revealed','cancelled')),
  player_snapshot    JSONB NOT NULL,
  revealed_at        TIMESTAMPTZ NULL,
  reveal_result      JSONB NULL,
  reveal_txn_id      UUID NULL,
  cancelled_at       TIMESTAMPTZ NULL,
  created_at         TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at         TIMESTAMPTZ NOT NULL DEFAULT now(),
  version            INT NOT NULL DEFAULT 1
);

-- Slot cap 1 active per character (matches hunt/activity).
CREATE UNIQUE INDEX ux_dungeon_contract_active
  ON dungeon.contract (character_id)
  WHERE state = 'active';

CREATE INDEX idx_dungeon_contract_history
  ON dungeon.contract (character_id, created_at DESC);

CREATE INDEX idx_dungeon_contract_expire
  ON dungeon.contract (expires_at)
  WHERE state = 'active';

-- --- Weekly lockout ---------------------------------------------------------

CREATE TABLE dungeon.weekly_lockout (
  character_id       UUID NOT NULL REFERENCES character.character(id) ON DELETE CASCADE,
  dungeon_code       VARCHAR(48) NOT NULL,
  week_key           VARCHAR(8)  NOT NULL,       -- "YYYY-WW" ISO week
  cleared_at         TIMESTAMPTZ NOT NULL DEFAULT now(),
  contract_id        UUID NOT NULL,               -- which run produced the lock
  PRIMARY KEY (character_id, dungeon_code, week_key)
);

CREATE INDEX idx_dungeon_weekly_lockout_char
  ON dungeon.weekly_lockout (character_id, week_key);
