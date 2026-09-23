-- 0010_skill_equipment.up.sql
-- Skill Engine + Equipment slots.
--
-- item_template gains:
--   equip_slot         — 'mainhand'|'helm'|'chest'|'boots'|'legs'|'gloves'|'necklace'|'ring'|'trinket' (NULL for materials/consumables)
--   class_restriction  — 'warrior'|'mage'|... (NULL for universal)
--   stat_bonuses       — JSONB array of {kind, target, amount}
--
-- item_instance gains:
--   equipment_slot     — when equipped, which slot occupies (mirrors template.equip_slot)
--                        UNIQUE (owner, equipment_slot) partial index enforces one item per slot.

ALTER TABLE content.item_template
  ADD COLUMN equip_slot        VARCHAR(24) NULL,
  ADD COLUMN class_restriction VARCHAR(48) NULL,
  ADD COLUMN stat_bonuses      JSONB NOT NULL DEFAULT '[]';

ALTER TABLE inventory.item_instance
  ADD COLUMN equipment_slot VARCHAR(24) NULL;

CREATE UNIQUE INDEX ux_item_instance_equip_slot
  ON inventory.item_instance (owner_character_id, equipment_slot)
  WHERE equipment_slot IS NOT NULL;

CREATE INDEX idx_item_instance_equipped
  ON inventory.item_instance (owner_character_id)
  WHERE equipment_slot IS NOT NULL;

-- Weapon templates — universal (no class restriction) for MVP; Phase 02
-- separates by class-specific weapon types.
UPDATE content.item_template SET
  equip_slot='mainhand',
  stat_bonuses='[
    {"kind":"flat","target":"attack","amount":10},
    {"kind":"flat","target":"str","amount":3}
  ]'::jsonb
  WHERE code='sword_iron';

UPDATE content.item_template SET
  equip_slot='mainhand',
  stat_bonuses='[
    {"kind":"flat","target":"attack","amount":18},
    {"kind":"flat","target":"str","amount":6},
    {"kind":"flat","target":"crit_pct","amount":2}
  ]'::jsonb
  WHERE code='sword_steel';

UPDATE content.item_template SET
  equip_slot='mainhand',
  class_restriction='hunter',
  stat_bonuses='[
    {"kind":"flat","target":"attack","amount":8},
    {"kind":"flat","target":"agi","amount":3},
    {"kind":"flat","target":"hit_pct","amount":3}
  ]'::jsonb
  WHERE code='bow_yew';

UPDATE content.item_template SET
  equip_slot='mainhand',
  class_restriction='mage',
  stat_bonuses='[
    {"kind":"flat","target":"attack","amount":5},
    {"kind":"flat","target":"int","amount":5},
    {"kind":"flat","target":"mp_max","amount":20}
  ]'::jsonb
  WHERE code='staff_apprentice';

-- Armor templates
UPDATE content.item_template SET
  equip_slot='helm',
  stat_bonuses='[
    {"kind":"flat","target":"defense","amount":5},
    {"kind":"flat","target":"vit","amount":2}
  ]'::jsonb
  WHERE code='helm_iron';

UPDATE content.item_template SET
  equip_slot='chest',
  stat_bonuses='[
    {"kind":"flat","target":"defense","amount":8},
    {"kind":"flat","target":"vit","amount":1},
    {"kind":"flat","target":"hp_max","amount":25}
  ]'::jsonb
  WHERE code='chest_leather';

UPDATE content.item_template SET
  equip_slot='boots',
  stat_bonuses='[
    {"kind":"flat","target":"defense","amount":3},
    {"kind":"flat","target":"agi","amount":2}
  ]'::jsonb
  WHERE code='boots_iron';
