-- 0036_leveling.up.sql
-- Character leveling 1..1000 + class progression (DD-250). Definitions (XP
-- curve, rewards, promotions, talents) are content (content/progression.yaml);
-- this holds per-character state and the audit trail.

ALTER TABLE character.character DROP CONSTRAINT IF EXISTS character_level_check;
ALTER TABLE character.character ADD CONSTRAINT character_level_check CHECK (level >= 1 AND level <= 1000);

-- Allocated attribute points (raw; available = earned by level - spent) and
-- XP held while a promotion choice is pending / past the level cap.
ALTER TABLE character.character
  ADD COLUMN attr_str  SMALLINT NOT NULL DEFAULT 0 CHECK (attr_str  >= 0),
  ADD COLUMN attr_vit  SMALLINT NOT NULL DEFAULT 0 CHECK (attr_vit  >= 0),
  ADD COLUMN attr_dex  SMALLINT NOT NULL DEFAULT 0 CHECK (attr_dex  >= 0),
  ADD COLUMN attr_int  SMALLINT NOT NULL DEFAULT 0 CHECK (attr_int  >= 0),
  ADD COLUMN attr_wis  SMALLINT NOT NULL DEFAULT 0 CHECK (attr_wis  >= 0),
  ADD COLUMN attr_luck SMALLINT NOT NULL DEFAULT 0 CHECK (attr_luck >= 0),
  ADD COLUMN promotion_overflow_xp BIGINT NOT NULL DEFAULT 0 CHECK (promotion_overflow_xp >= 0),
  ADD COLUMN postcap_xp BIGINT NOT NULL DEFAULT 0 CHECK (postcap_xp >= 0);

-- One row per XP-granting source: a retried reveal / turn-in never pays twice.
CREATE TABLE character.xp_grant (
  character_id UUID NOT NULL REFERENCES character.character(id) ON DELETE CASCADE,
  source       VARCHAR(32) NOT NULL,
  source_id    VARCHAR(96) NOT NULL,
  amount       BIGINT NOT NULL CHECK (amount >= 0),
  created_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
  PRIMARY KEY (character_id, source, source_id)
);

-- Auditable progression events; (kind, ref) unique per character so a level
-- reward, unlock or title is recorded exactly once.
CREATE TABLE character.progression_event (
  id           BIGSERIAL PRIMARY KEY,
  character_id UUID NOT NULL REFERENCES character.character(id) ON DELETE CASCADE,
  kind         VARCHAR(32) NOT NULL,
  ref          VARCHAR(96) NOT NULL,
  level        SMALLINT NOT NULL,
  metadata     JSONB NOT NULL DEFAULT '{}',
  created_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
  UNIQUE (character_id, kind, ref)
);
CREATE INDEX idx_progression_event_char ON character.progression_event (character_id, id DESC);

CREATE TABLE character.class_history (
  id           BIGSERIAL PRIMARY KEY,
  character_id UUID NOT NULL REFERENCES character.character(id) ON DELETE CASCADE,
  from_spec    VARCHAR(48) NOT NULL,
  to_spec      VARCHAR(48) NOT NULL,
  level        SMALLINT NOT NULL,
  reason       VARCHAR(16) NOT NULL CHECK (reason IN ('promotion','respec')),
  created_at   TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX idx_class_history_char ON character.class_history (character_id, id);

CREATE TABLE character.talent (
  character_id UUID NOT NULL REFERENCES character.character(id) ON DELETE CASCADE,
  talent_code  VARCHAR(64) NOT NULL,
  rank         SMALLINT NOT NULL CHECK (rank >= 1),
  updated_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
  PRIMARY KEY (character_id, talent_code)
);

-- Class skills: levels to 1000, a specialization owner and data-driven
-- activation (auto-trigger rules, resource, duration) in config.
ALTER TABLE content.class_skill DROP CONSTRAINT IF EXISTS class_skill_unlock_level_check;
ALTER TABLE content.class_skill ADD CONSTRAINT class_skill_unlock_level_check CHECK (unlock_level BETWEEN 1 AND 1000);
ALTER TABLE content.class_skill DROP CONSTRAINT IF EXISTS class_skill_effect_check;
ALTER TABLE content.class_skill ADD CONSTRAINT class_skill_effect_check CHECK (effect IN ('strike','heal','guard','empower'));
ALTER TABLE content.class_skill
  ADD COLUMN spec_code VARCHAR(48) NULL,
  ADD COLUMN config    JSONB NOT NULL DEFAULT '{}';

-- Equipment requirements beyond level/class (DD-250 §34): a class branch
-- and raw attribute minimums, e.g. {"spec":"guardian","strength":220}.
ALTER TABLE content.item_template ADD COLUMN requirements JSONB NULL;

-- Level titles (prefix: "Veteran Warrior") and class mastery titles (suffix).
INSERT INTO content.title_template (code, display_key, placement, rarity, source_kind) VALUES
  ('veteran',             'title.veteran',             'prefix', 'common',    'level'),
  ('elite',               'title.elite',               'prefix', 'uncommon',  'level'),
  ('champion',            'title.champion',            'prefix', 'uncommon',  'level'),
  ('hero',                'title.hero',                'prefix', 'rare',      'level'),
  ('legendary',           'title.legendary',           'prefix', 'rare',      'level'),
  ('mythic',              'title.mythic',              'prefix', 'epic',      'level'),
  ('ancient',             'title.ancient',             'prefix', 'epic',      'level'),
  ('ascendant',           'title.ascendant',           'prefix', 'legendary', 'level'),
  ('immortal',            'title.immortal',            'prefix', 'legendary', 'level'),
  ('eternal',             'title.eternal',             'prefix', 'mythic',    'level'),
  ('the_unbreakable',     'title.the_unbreakable',     'suffix', 'mythic',    'class_mastery'),
  ('lord_of_battle',      'title.lord_of_battle',      'suffix', 'mythic',    'class_mastery'),
  ('bloodbound',          'title.bloodbound',          'suffix', 'mythic',    'class_mastery'),
  ('grandmaster_of_arms', 'title.grandmaster_of_arms', 'suffix', 'mythic',    'class_mastery')
ON CONFLICT (code) DO NOTHING;
