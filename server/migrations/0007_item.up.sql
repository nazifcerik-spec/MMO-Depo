-- 0007_item.up.sql
-- Item Engine — template + instance schema.
--
-- Creates the inventory schema (missing from baseline before Sprint 8).

CREATE SCHEMA IF NOT EXISTS inventory;

--
--
-- Templates yaşam boyu değişmez (content pack ile deploy; Sprint 12 YAML loader).
-- Instances mutable owner + affix snapshot (roll edildiği zaman dondurulur).
--
-- v1 kapsamı:
--   - template: code + class/subclass + level + base_value
--   - instance: template ref + rarity + prefix + suffix + owner
-- v1 dışı (Phase 02):
--   - socket + gem + enchantment
--   - set item bonus
--   - durability
--   - unique effect

CREATE TABLE content.item_template (
  code             VARCHAR(64) PRIMARY KEY,
  display_name_key VARCHAR(64) NOT NULL,
  item_class       VARCHAR(24) NOT NULL,       -- weapon|armor|consumable|material|misc
  item_subclass    VARCHAR(24) NOT NULL,       -- sword|axe|helm|potion|ore|...
  min_level        SMALLINT NOT NULL DEFAULT 1,
  base_value_copper BIGINT   NOT NULL DEFAULT 0,   -- vendor sell reference
  base_weight_g    INT       NOT NULL DEFAULT 100, -- weight in grams (Sprint 9 inventory cap)
  stackable        BOOLEAN   NOT NULL DEFAULT false,
  max_stack        INT       NOT NULL DEFAULT 1,
  bind_kind        VARCHAR(8) NOT NULL DEFAULT 'boe',  -- 'boe'|'bop'|'boa'
  no_trade         BOOLEAN   NOT NULL DEFAULT false,
  no_sell          BOOLEAN   NOT NULL DEFAULT false,
  affixable        BOOLEAN   NOT NULL DEFAULT false,  -- true for weapons/armor
  version          INT       NOT NULL DEFAULT 1,
  created_at       TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_item_template_class ON content.item_template (item_class, item_subclass);

CREATE TABLE inventory.item_instance (
  id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  template_code     VARCHAR(64) NOT NULL REFERENCES content.item_template(code),
  owner_character_id UUID NULL REFERENCES character.character(id) ON DELETE SET NULL,
  stack_quantity    INT NOT NULL DEFAULT 1 CHECK (stack_quantity >= 1),
  rarity            VARCHAR(16) NOT NULL DEFAULT 'common',
  prefix_code       VARCHAR(48) NULL,           -- affix prefix (rare+ items)
  suffix_code       VARCHAR(48) NULL,           -- affix suffix (rare+ items)
  is_bound          BOOLEAN NOT NULL DEFAULT false, -- BoP items or equipped BoE
  is_locked         BOOLEAN NOT NULL DEFAULT false, -- market/trade/mail/craft escrow
  lock_owner        VARCHAR(16) NULL,           -- 'market'|'trade'|'mail'|'craft'|'admin'
  locked_ref_id     UUID NULL,                  -- listing_id / trade_session_id / ...
  source_kind       VARCHAR(24) NULL,           -- 'activity_reveal'|'loot'|'craft'|'vendor'|'mint'
  source_ref_id     UUID NULL,                  -- contract_id / dungeon_run_id / ...
  created_at        TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at        TIMESTAMPTZ NOT NULL DEFAULT now(),
  version           INT NOT NULL DEFAULT 1
);

CREATE INDEX idx_item_instance_owner ON inventory.item_instance (owner_character_id) WHERE owner_character_id IS NOT NULL;
CREATE INDEX idx_item_instance_template ON inventory.item_instance (template_code);
CREATE INDEX idx_item_instance_source ON inventory.item_instance (source_kind, source_ref_id) WHERE source_ref_id IS NOT NULL;
CREATE INDEX idx_item_instance_lock ON inventory.item_instance (lock_owner, locked_ref_id) WHERE is_locked = true;

-- --- Seed: v1 item template catalog -----------------------------------------
-- 30 template covering mining/fishing/woodcutting/herbalism/hunting drops +
-- basic weapons/armor for future combat. Balance placeholder.

INSERT INTO content.item_template (code, display_name_key, item_class, item_subclass, min_level, base_value_copper, base_weight_g, stackable, max_stack, bind_kind, affixable) VALUES
  -- Mining
  ('iron_ore',       'item.iron_ore',       'material', 'ore',   1,   5, 200, true,  99, 'boe', false),
  ('copper_ore',     'item.copper_ore',     'material', 'ore',   5,   8, 200, true,  99, 'boe', false),
  ('tin_ore',        'item.tin_ore',        'material', 'ore',   5,   7, 200, true,  99, 'boe', false),
  ('silver_ore',     'item.silver_ore',     'material', 'ore',  15,  25, 200, true,  99, 'boe', false),
  ('gold_ore',       'item.gold_ore',       'material', 'ore',  20,  40, 250, true,  99, 'boe', false),
  ('mythril_ore',    'item.mythril_ore',    'material', 'ore',  40, 250, 300, true,  99, 'boe', false),

  -- Fishing
  ('salmon',         'item.salmon',         'material', 'fish',  1,   8, 500, true,  99, 'boe', false),
  ('bass',           'item.bass',           'material', 'fish',  5,  12, 400, true,  99, 'boe', false),
  ('golden_carp',    'item.golden_carp',    'material', 'fish', 15,  50, 500, true,  99, 'boe', false),
  ('ancient_leviathan_scale', 'item.leviathan_scale', 'material', 'fish', 35, 300, 100, true, 99, 'boe', false),

  -- Woodcutting
  ('oak_log',        'item.oak_log',        'material', 'wood',  1,   4, 1000, true, 99, 'boe', false),
  ('elm_log',        'item.elm_log',        'material', 'wood',  8,   9, 1000, true, 99, 'boe', false),
  ('yew_log',        'item.yew_log',        'material', 'wood', 20,  30, 1000, true, 99, 'boe', false),
  ('ironwood_log',   'item.ironwood_log',   'material', 'wood', 40, 200, 1200, true, 99, 'boe', false),

  -- Herbalism
  ('linseed',        'item.linseed',        'material', 'herb',  1,   7, 50,  true, 99, 'boe', false),
  ('mandrake',       'item.mandrake',       'material', 'herb',  5,  12, 50,  true, 99, 'boe', false),
  ('king_bloom',     'item.king_bloom',     'material', 'herb', 20,  60, 60,  true, 99, 'boe', false),
  ('phoenix_petal',  'item.phoenix_petal',  'material', 'herb', 40, 500, 30,  true, 99, 'boe', false),

  -- Hunting
  ('boar_hide',      'item.boar_hide',      'material', 'hide',  1,  15, 800, true, 99, 'boe', false),
  ('wolf_pelt',      'item.wolf_pelt',      'material', 'hide',  8,  25, 700, true, 99, 'boe', false),
  ('dire_bear_hide', 'item.dire_bear_hide', 'material', 'hide', 20,  90, 1500,true, 99, 'boe', false),
  ('dragon_scale',   'item.dragon_scale',   'material', 'hide', 45, 800, 300, true, 99, 'boe', false),

  -- Weapons (rare+ get affix; used post-item drop in dungeons/hunts)
  ('sword_iron',     'item.sword_iron',     'weapon',   'sword', 5,  100, 2500, false, 1, 'boe', true),
  ('sword_steel',    'item.sword_steel',    'weapon',   'sword', 15, 400, 2700, false, 1, 'boe', true),
  ('bow_yew',        'item.bow_yew',        'weapon',   'bow',   10, 250, 1500, false, 1, 'boe', true),
  ('staff_apprentice', 'item.staff_apprentice', 'weapon','staff', 3,  80, 1800, false, 1, 'boe', true),

  -- Armor
  ('helm_iron',      'item.helm_iron',      'armor',    'helm',  5,   80, 2000, false, 1, 'boe', true),
  ('chest_leather',  'item.chest_leather',  'armor',    'chest', 1,   60, 3000, false, 1, 'boe', true),
  ('boots_iron',     'item.boots_iron',     'armor',    'boots', 5,   50, 1500, false, 1, 'boe', true),

  -- Consumables
  ('potion_hp_small','item.potion_hp_small','consumable','potion',1,  20, 200, true, 20, 'boe', false),
  ('potion_mp_small','item.potion_mp_small','consumable','potion',1,  20, 200, true, 20, 'boe', false);
