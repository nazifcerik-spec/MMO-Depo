-- 0022_world_events.up.sql
-- World Events (Sprint 25): world boss + invasion + seasonal event.
-- docs/world-events/{world-boss,invasion,seasonal-event}.md.
--
-- World boss and invasion share one model: a scheduled spawn with a shared
-- timer and a personal contribution per attendee (DD-150). Invasions are
-- region-gated (the character must stand in the region to enlist). v1
-- contribution is AFK-only: power snapshot at enlist × minutes attended,
-- capped at the tier-3 threshold (world-boss.md §5). Regional and dynamic
-- events are Phase 02.

CREATE SCHEMA IF NOT EXISTS worldevent;

INSERT INTO economy.currency (code, display_name_key, scope, tradable, cap, base_unit_ratio) VALUES
  ('event_token', 'currency.event_token', 'account', false, 1000000, 1)
ON CONFLICT (code) DO NOTHING;

CREATE TABLE content.world_event_template (
  code              VARCHAR(48) PRIMARY KEY,
  kind              VARCHAR(16) NOT NULL CHECK (kind IN ('world_boss','invasion')),
  display_name_key  VARCHAR(64) NOT NULL,
  min_level         SMALLINT NOT NULL DEFAULT 1,
  region_code       VARCHAR(48) NULL REFERENCES world.region(code),  -- required for invasions
  hp_base           BIGINT NOT NULL CHECK (hp_base > 0),               -- total contribution needed to win
  duration_seconds  INT NOT NULL DEFAULT 3600 CHECK (duration_seconds BETWEEN 600 AND 21600),
  tiers             JSONB NOT NULL,   -- [{min_score, gold_copper, event_tokens}] ascending, 4 entries
  afk_cap_tier      SMALLINT NOT NULL DEFAULT 3,
  CHECK (kind <> 'invasion' OR region_code IS NOT NULL)
);

CREATE TABLE content.world_event_rotation (
  template_code  VARCHAR(48) NOT NULL REFERENCES content.world_event_template(code) ON DELETE CASCADE,
  day_of_week    SMALLINT NOT NULL CHECK (day_of_week BETWEEN 0 AND 6),  -- 0 = Sunday
  hour_utc       SMALLINT NOT NULL CHECK (hour_utc BETWEEN 0 AND 23),
  PRIMARY KEY (template_code, day_of_week, hour_utc)
);

CREATE TABLE worldevent.spawn (
  id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  template_code  VARCHAR(48) NOT NULL REFERENCES content.world_event_template(code),
  starts_at      TIMESTAMPTZ NOT NULL,
  ends_at        TIMESTAMPTZ NOT NULL,
  state          VARCHAR(16) NOT NULL DEFAULT 'scheduled'
                 CHECK (state IN ('scheduled','active','succeeded','failed')),
  total_score    BIGINT NOT NULL DEFAULT 0,
  closed_at      TIMESTAMPTZ NULL,
  UNIQUE (template_code, starts_at)
);
CREATE INDEX idx_spawn_open ON worldevent.spawn (starts_at) WHERE state IN ('scheduled','active');

CREATE TABLE worldevent.attendance (
  spawn_id            UUID NOT NULL REFERENCES worldevent.spawn(id) ON DELETE CASCADE,
  character_id        UUID NOT NULL REFERENCES character.character(id) ON DELETE CASCADE,
  account_id          UUID NOT NULL,
  registered_at       TIMESTAMPTZ NOT NULL DEFAULT now(),
  power               INT NOT NULL,
  contribution_score  BIGINT NOT NULL DEFAULT 0,
  reward_tier         SMALLINT NULL,
  rewarded_at         TIMESTAMPTZ NULL,
  PRIMARY KEY (spawn_id, character_id)
);

CREATE TABLE content.seasonal_event (
  code                VARCHAR(48) PRIMARY KEY,
  display_name_key    VARCHAR(64) NOT NULL,
  starts_at           TIMESTAMPTZ NOT NULL,
  ends_at             TIMESTAMPTZ NOT NULL,
  daily_login_tokens  INT NOT NULL DEFAULT 15,
  CHECK (ends_at > starts_at)
);

CREATE TABLE worldevent.seasonal_claim (
  event_code  VARCHAR(48) NOT NULL REFERENCES content.seasonal_event(code) ON DELETE CASCADE,
  account_id  UUID NOT NULL,
  claim_date  DATE NOT NULL,
  tokens      INT NOT NULL,
  claimed_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
  PRIMARY KEY (event_code, account_id, claim_date)
);

-- --- Seed ------------------------------------------------------------------

INSERT INTO content.world_event_template (code, kind, display_name_key, min_level, region_code, hp_base, duration_seconds, tiers) VALUES
  ('ashen_colossus', 'world_boss', 'worldevent.ashen_colossus', 10, NULL, 500000, 3600,
   '[{"min_score":1,"gold_copper":100,"event_tokens":1},{"min_score":1000,"gold_copper":500,"event_tokens":5},{"min_score":10000,"gold_copper":2000,"event_tokens":15},{"min_score":50000,"gold_copper":10000,"event_tokens":40}]'),
  ('greenwood_incursion', 'invasion', 'worldevent.greenwood_incursion', 5, 'greenwood', 150000, 7200,
   '[{"min_score":1,"gold_copper":80,"event_tokens":1},{"min_score":800,"gold_copper":400,"event_tokens":4},{"min_score":8000,"gold_copper":1500,"event_tokens":12},{"min_score":40000,"gold_copper":8000,"event_tokens":30}]');

INSERT INTO content.world_event_rotation (template_code, day_of_week, hour_utc) VALUES
  ('ashen_colossus', 2, 19), ('ashen_colossus', 5, 19), ('ashen_colossus', 0, 15),
  ('greenwood_incursion', 3, 18), ('greenwood_incursion', 6, 14);

INSERT INTO content.seasonal_event (code, display_name_key, starts_at, ends_at, daily_login_tokens) VALUES
  ('harvest_festival_2026', 'seasonal.harvest_festival', '2026-09-20T00:00:00Z', '2026-10-18T00:00:00Z', 15);
