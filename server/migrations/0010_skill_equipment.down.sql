-- 0010_skill_equipment.down.sql
DROP INDEX IF EXISTS inventory.idx_item_instance_equipped;
DROP INDEX IF EXISTS inventory.ux_item_instance_equip_slot;
ALTER TABLE inventory.item_instance DROP COLUMN IF EXISTS equipment_slot;

ALTER TABLE content.item_template
  DROP COLUMN IF EXISTS stat_bonuses,
  DROP COLUMN IF EXISTS class_restriction,
  DROP COLUMN IF EXISTS equip_slot;
