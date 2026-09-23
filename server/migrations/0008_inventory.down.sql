-- 0008_inventory.down.sql
DROP INDEX IF EXISTS inventory.idx_item_instance_bag;
DROP INDEX IF EXISTS inventory.ux_item_instance_owner_slot;
ALTER TABLE inventory.item_instance DROP COLUMN IF EXISTS slot;
