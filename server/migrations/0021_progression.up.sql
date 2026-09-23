-- 0021_progression.up.sql
-- Progression (Sprint 24): mastery, achievements, titles, reputation.
-- docs/progression/{mastery,achievements,titles,reputation}.md.
--
-- One event stream drives all three counters: game systems report events
-- (kind + amount, e.g. mob_killed ×12) and content decides what each event
-- feeds. Content lives here as a seed; ascension/collections/milestones are
-- later work.

CREATE SCHEMA IF NOT EXISTS progression;

-- --- Content ---------------------------------------------------------------

CREATE TABLE content.mastery_track (
  code              VARCHAR(48) PRIMARY KEY,
  display_name_key  VARCHAR(64) NOT NULL,
  category          VARCHAR(24) NOT NULL,
  max_tier          SMALLINT NOT NULL DEFAULT 10,
  tier_xp_curve     JSONB NOT NULL,   -- [xp for tier 1, xp for tier 2, ...] (increments)
  tier_bonuses      JSONB NOT NULL DEFAULT '[]', -- [{tier, bonus_effects:[{kind,target,amount}]}]
  sources           JSONB NOT NULL    -- [{kind, xp_per_unit}]
);

CREATE TABLE content.title_template (
  code         VARCHAR(48) PRIMARY KEY,
  display_key  VARCHAR(64) NOT NULL,
  placement    VARCHAR(16) NOT NULL DEFAULT 'suffix' CHECK (placement IN ('prefix','suffix')),
  rarity       VARCHAR(16) NOT NULL DEFAULT 'common',
  source_kind  VARCHAR(24) NOT NULL DEFAULT 'achievement'
);

CREATE TABLE content.achievement_template (
  code              VARCHAR(64) PRIMARY KEY,
  display_name_key  VARCHAR(64) NOT NULL,
  description_key   VARCHAR(128) NOT NULL,
  category          VARCHAR(24) NOT NULL,
  point_value       SMALLINT NOT NULL,
  is_hidden         BOOLEAN NOT NULL DEFAULT false,
  trigger_kind      VARCHAR(64) NOT NULL,
  progress_target   INT NOT NULL DEFAULT 1 CHECK (progress_target >= 1),
  reward_title_code VARCHAR(48) NULL REFERENCES content.title_template(code),
  reward_gold_copper BIGINT NOT NULL DEFAULT 0
);
CREATE INDEX idx_achievement_trigger ON content.achievement_template (trigger_kind);

CREATE TABLE content.faction (
  code              VARCHAR(48) PRIMARY KEY,
  display_name_key  VARCHAR(64) NOT NULL,
  daily_cap         INT NOT NULL DEFAULT 1500,
  sources           JSONB NOT NULL    -- [{kind, delta}]
);

-- --- Per-character state ---------------------------------------------------

CREATE TABLE progression.character_mastery (
  character_id  UUID NOT NULL REFERENCES character.character(id) ON DELETE CASCADE,
  track_code    VARCHAR(48) NOT NULL REFERENCES content.mastery_track(code),
  xp            BIGINT NOT NULL DEFAULT 0,
  current_tier  SMALLINT NOT NULL DEFAULT 0,
  last_gain_at  TIMESTAMPTZ NULL,
  PRIMARY KEY (character_id, track_code)
);

CREATE TABLE progression.achievement_progress (
  character_id      UUID NOT NULL REFERENCES character.character(id) ON DELETE CASCADE,
  achievement_code  VARCHAR(64) NOT NULL REFERENCES content.achievement_template(code),
  progress_current  INT NOT NULL DEFAULT 0,
  unlocked_at       TIMESTAMPTZ NULL,
  PRIMARY KEY (character_id, achievement_code)
);

CREATE TABLE progression.character_title (
  character_id     UUID NOT NULL REFERENCES character.character(id) ON DELETE CASCADE,
  title_code       VARCHAR(48) NOT NULL REFERENCES content.title_template(code),
  unlocked_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
  source_ref_kind  VARCHAR(24) NOT NULL,
  PRIMARY KEY (character_id, title_code)
);

CREATE TABLE progression.character_reputation (
  character_id     UUID NOT NULL REFERENCES character.character(id) ON DELETE CASCADE,
  faction_code     VARCHAR(48) NOT NULL REFERENCES content.faction(code),
  rep_xp           INT NOT NULL DEFAULT 0,
  current_tier     SMALLINT NOT NULL DEFAULT 3,
  daily_gain_used  INT NOT NULL DEFAULT 0,
  daily_reset_on   DATE NOT NULL DEFAULT (now() AT TIME ZONE 'UTC')::date,
  last_gain_at     TIMESTAMPTZ NULL,
  PRIMARY KEY (character_id, faction_code)
);

ALTER TABLE character.character
  ADD COLUMN achievement_points INT NOT NULL DEFAULT 0,
  ADD COLUMN active_title_code  VARCHAR(48) NULL REFERENCES content.title_template(code);

-- --- Seed ------------------------------------------------------------------

INSERT INTO content.mastery_track (code, display_name_key, category, tier_xp_curve, tier_bonuses, sources) VALUES
  ('mining',        'mastery.mining',        'gathering',   '[100,500,2000,8000,30000,100000,300000,800000,2000000,5000000]',
   '[{"tier":1,"bonus_effects":[{"kind":"flat","target":"gather_yield_mining","amount":1}]},{"tier":5,"bonus_effects":[{"kind":"mult_add","target":"gather_yield_mining","amount":0.05}]},{"tier":10,"bonus_effects":[{"kind":"mult_add","target":"gather_yield_mining","amount":0.10}]}]',
   '[{"kind":"gather_mining","xp_per_unit":10}]'),
  ('fishing',       'mastery.fishing',       'gathering',   '[100,500,2000,8000,30000,100000,300000,800000,2000000,5000000]', '[]',
   '[{"kind":"gather_fishing","xp_per_unit":10}]'),
  ('woodcutting',   'mastery.woodcutting',   'gathering',   '[100,500,2000,8000,30000,100000,300000,800000,2000000,5000000]', '[]',
   '[{"kind":"gather_woodcutting","xp_per_unit":10}]'),
  ('herbalism',     'mastery.herbalism',     'gathering',   '[100,500,2000,8000,30000,100000,300000,800000,2000000,5000000]', '[]',
   '[{"kind":"gather_herbalism","xp_per_unit":10}]'),
  ('combat',        'mastery.combat',        'combat',      '[100,500,2000,8000,30000,100000,300000,800000,2000000,5000000]',
   '[{"tier":3,"bonus_effects":[{"kind":"flat","target":"attack","amount":2}]},{"tier":6,"bonus_effects":[{"kind":"flat","target":"attack","amount":5}]}]',
   '[{"kind":"mob_killed","xp_per_unit":15},{"kind":"gather_hunting","xp_per_unit":5}]'),
  ('dungeoneering', 'mastery.dungeoneering', 'combat',      '[100,500,2000,8000,30000,100000,300000,800000,2000000,5000000]', '[]',
   '[{"kind":"dungeon_cleared","xp_per_unit":250},{"kind":"boss_killed","xp_per_unit":100}]'),
  ('crafting',      'mastery.crafting',      'crafting',    '[100,500,2000,8000,30000,100000,300000,800000,2000000,5000000]', '[]',
   '[{"kind":"item_crafted","xp_per_unit":40}]');

INSERT INTO content.title_template (code, display_key, placement, rarity) VALUES
  ('the_slayer',   'title.the_slayer',   'suffix', 'rare'),
  ('the_delver',   'title.the_delver',   'suffix', 'rare'),
  ('bossbane',     'title.bossbane',     'prefix', 'epic'),
  ('the_artisan',  'title.the_artisan',  'suffix', 'rare'),
  ('prospector',   'title.prospector',   'prefix', 'uncommon');

INSERT INTO content.achievement_template
  (code, display_name_key, description_key, category, point_value, trigger_kind, progress_target, reward_title_code, reward_gold_copper)
VALUES
  ('first_blood',     'ach.first_blood',     'ach.first_blood.desc',     'combat',      5, 'mob_killed',       1,   NULL,          500),
  ('slayer_100',      'ach.slayer_100',      'ach.slayer_100.desc',      'combat',     10, 'mob_killed',     100,   NULL,         5000),
  ('slayer_1000',     'ach.slayer_1000',     'ach.slayer_1000.desc',     'combat',     25, 'mob_killed',    1000,   'the_slayer', 25000),
  ('dungeon_first',   'ach.dungeon_first',   'ach.dungeon_first.desc',   'dungeon',    10, 'dungeon_cleared',  1,   NULL,         2500),
  ('dungeon_25',      'ach.dungeon_25',      'ach.dungeon_25.desc',      'dungeon',    25, 'dungeon_cleared', 25,   'the_delver', 20000),
  ('boss_hunter',     'ach.boss_hunter',     'ach.boss_hunter.desc',     'dungeon',    50, 'boss_killed',     10,   'bossbane',   50000),
  ('miner_500',       'ach.miner_500',       'ach.miner_500.desc',       'gathering',  10, 'gather_mining',  500,   'prospector',  5000),
  ('crafter_first',   'ach.crafter_first',   'ach.crafter_first.desc',   'crafting',    5, 'item_crafted',     1,   NULL,          500),
  ('crafter_100',     'ach.crafter_100',     'ach.crafter_100.desc',     'crafting',   25, 'item_crafted',   100,   'the_artisan',20000);

INSERT INTO content.faction (code, display_name_key, sources) VALUES
  ('ironhold_guard',    'faction.ironhold_guard',
   '[{"kind":"dungeon_cleared:crypt_of_ashes","delta":150},{"kind":"mob_killed","delta":1}]'),
  ('greenwood_wardens', 'faction.greenwood_wardens',
   '[{"kind":"dungeon_cleared:forest_hollow","delta":200},{"kind":"gather_herbalism","delta":1},{"kind":"gather_woodcutting","delta":1}]'),
  ('ember_covenant',    'faction.ember_covenant',
   '[{"kind":"dungeon_cleared:shadowfen_ruins","delta":300},{"kind":"dungeon_cleared:obsidian_spire","delta":400}]');
