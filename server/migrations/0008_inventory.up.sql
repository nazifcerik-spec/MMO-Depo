-- 0008_inventory.up.sql
-- Inventory slot assignment layer.
--
-- Approach: single-table model. inventory.item_instance gains a `slot` column
-- that positions the item in the character's bag (0..N-1). Slot NULL means
-- "detached" — owned but not visible in inventory (overflow, bank Phase 02,
-- mail escrow, market escrow — the escrows use is_locked + lock_owner).
--
-- Sprint 9 v1:
--   slot cap 60
--   stackable items merge on pickup up to max_stack
--   Move / Split / Delete via service (weight cap enforced)
-- Phase 02:
--   Bank tabs (separate table)
--   Equipment slots (separate table or is_equipped flag)

ALTER TABLE inventory.item_instance
  ADD COLUMN slot SMALLINT NULL CHECK (slot IS NULL OR (slot >= 0 AND slot < 240));
  -- Cap of 240 allows Phase 02 bank slots + expanded bag without further migration.
  -- v1 service enforces effective cap = 60 for bag.

-- One item per (owner, slot) — race-safe slot allocation.
CREATE UNIQUE INDEX ux_item_instance_owner_slot
  ON inventory.item_instance (owner_character_id, slot)
  WHERE slot IS NOT NULL;

-- Fast bag view.
CREATE INDEX idx_item_instance_bag
  ON inventory.item_instance (owner_character_id, slot)
  WHERE slot IS NOT NULL AND owner_character_id IS NOT NULL;
