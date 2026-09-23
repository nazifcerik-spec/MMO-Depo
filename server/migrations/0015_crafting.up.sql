-- 0015_crafting.up.sql
-- Professions + Crafting (Sprint 17).
--
-- Katmanlar:
--   content.profession        meslek katalog (blacksmithing, alchemy, ...)
--   content.recipe            reçete: inputs + output + skill req + xp
--   character.profession      karakterin öğrendiği meslek + level/xp
--   character.craft_log       audit: kim ne zaman ne yaptı (opsiyonel diagnostic)
--
-- Level = 1..MaxLevel; XP CHECK constraint yerine app-side clamp (LevelCap).

CREATE SCHEMA IF NOT EXISTS crafting;

-- --- Profession catalog -----------------------------------------------------

CREATE TABLE content.profession (
  code             VARCHAR(48) PRIMARY KEY,
  display_name_key VARCHAR(64) NOT NULL,
  max_level        SMALLINT NOT NULL DEFAULT 60
                   CHECK (max_level > 0 AND max_level <= 100),
  version          INT NOT NULL DEFAULT 1,
  created_at       TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- --- Recipe catalog ---------------------------------------------------------

CREATE TABLE content.recipe (
  code                 VARCHAR(64) PRIMARY KEY,
  display_name_key     VARCHAR(64) NOT NULL,
  profession_code      VARCHAR(48) NOT NULL REFERENCES content.profession(code) ON DELETE RESTRICT,
  min_skill            SMALLINT NOT NULL DEFAULT 1,

  -- Inputs JSONB: [{"template_code": "iron_ore", "quantity": 3}, ...]
  inputs               JSONB NOT NULL,

  -- Optional gold cost paid via TxnService to the crafting sink.
  gold_cost_copper     BIGINT NOT NULL DEFAULT 0 CHECK (gold_cost_copper >= 0),

  -- Output: one mint of `output_template_code × output_stack`.
  output_template_code VARCHAR(64) NOT NULL,
  output_stack         INT NOT NULL DEFAULT 1 CHECK (output_stack >= 1),

  -- Profession XP awarded on successful craft (independent from character XP).
  profession_xp_reward INT NOT NULL DEFAULT 0 CHECK (profession_xp_reward >= 0),

  version              INT NOT NULL DEFAULT 1,
  created_at           TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_recipe_profession ON content.recipe (profession_code, min_skill);

-- --- Per-character profession state -----------------------------------------

CREATE TABLE crafting.character_profession (
  character_id     UUID NOT NULL REFERENCES character.character(id) ON DELETE CASCADE,
  profession_code  VARCHAR(48) NOT NULL REFERENCES content.profession(code) ON DELETE RESTRICT,
  level            SMALLINT NOT NULL DEFAULT 1 CHECK (level >= 1),
  xp               INT NOT NULL DEFAULT 0 CHECK (xp >= 0),
  learned_at       TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at       TIMESTAMPTZ NOT NULL DEFAULT now(),
  version          INT NOT NULL DEFAULT 1,
  PRIMARY KEY (character_id, profession_code)
);

CREATE INDEX idx_character_profession_char ON crafting.character_profession (character_id);

-- --- Audit log (append-only) ------------------------------------------------

CREATE TABLE crafting.craft_log (
  id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  character_id      UUID NOT NULL REFERENCES character.character(id) ON DELETE CASCADE,
  recipe_code       VARCHAR(64) NOT NULL,
  output_instance_id UUID NULL,
  gold_txn_id       UUID NULL,
  crafted_at        TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_craft_log_char ON crafting.craft_log (character_id, crafted_at DESC);

-- --- Seed baseline (YAML in content/professions.yaml + content/recipes.yaml
-- authoritative; this guarantees a fresh DB is playable without loader) ----

INSERT INTO content.profession (code, display_name_key, max_level) VALUES
  ('blacksmithing', 'profession.blacksmithing', 60),
  ('alchemy',       'profession.alchemy',       60),
  ('cooking',       'profession.cooking',       60);

INSERT INTO content.recipe (code, display_name_key, profession_code, min_skill,
                            inputs, gold_cost_copper,
                            output_template_code, output_stack, profession_xp_reward)
VALUES
  -- Alchemy — starter potions
  ('recipe_potion_hp_small', 'recipe.potion_hp_small', 'alchemy', 1,
    '[{"template_code":"linseed","quantity":2}]'::jsonb, 5,
    'potion_hp_small', 1, 15),
  ('recipe_potion_mp_small', 'recipe.potion_mp_small', 'alchemy', 3,
    '[{"template_code":"mandrake","quantity":2}]'::jsonb, 5,
    'potion_mp_small', 1, 20),

  -- Blacksmithing — starter weapon/armor
  ('recipe_sword_iron', 'recipe.sword_iron', 'blacksmithing', 5,
    '[{"template_code":"iron_ore","quantity":5},{"template_code":"oak_log","quantity":1}]'::jsonb, 25,
    'sword_iron', 1, 40),
  ('recipe_helm_iron', 'recipe.helm_iron', 'blacksmithing', 5,
    '[{"template_code":"iron_ore","quantity":4}]'::jsonb, 20,
    'helm_iron', 1, 35),
  ('recipe_boots_iron', 'recipe.boots_iron', 'blacksmithing', 5,
    '[{"template_code":"iron_ore","quantity":3}]'::jsonb, 15,
    'boots_iron', 1, 30);
