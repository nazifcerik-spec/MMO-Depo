-- 0027_class_skills.up.sql
-- Class skills v1 (Sprint 42, DD-222): per-class passives (sheet modifiers)
-- and actives (auto-cast in combat), unlocked by level. No per-character
-- state: what a character has follows from class + level.

CREATE TABLE content.class_skill (
  code             VARCHAR(64) PRIMARY KEY,
  display_name_key VARCHAR(96) NOT NULL,
  class_code       VARCHAR(24) NOT NULL,
  unlock_level     SMALLINT NOT NULL CHECK (unlock_level BETWEEN 1 AND 60),
  kind             VARCHAR(8) NOT NULL CHECK (kind IN ('active','passive')),
  sort_order       INT NOT NULL DEFAULT 0,            -- actives: cast priority (low first)
  effect           VARCHAR(16) NULL CHECK (effect IN ('strike','heal')),
  power            INT NOT NULL DEFAULT 0,
  mp_cost          INT NOT NULL DEFAULT 0 CHECK (mp_cost >= 0),
  cooldown_turns   INT NOT NULL DEFAULT 1 CHECK (cooldown_turns >= 1),
  modifiers        JSONB NOT NULL DEFAULT '[]',       -- passives: [{kind,target,amount}]
  version          INT NOT NULL DEFAULT 1,
  CHECK ((kind = 'active') = (effect IS NOT NULL))
);
CREATE INDEX idx_class_skill_class ON content.class_skill (class_code, unlock_level);
