-- 0045_enchanting.up.sql
-- Enchanting (DD-293): a scroll template carries what it puts on gear
-- (slot groups + modifiers); a piece of gear holds one enchant, the code of
-- the scroll that made it.
ALTER TABLE content.item_template ADD COLUMN enchant JSONB NULL;
ALTER TABLE inventory.item_instance
  ADD COLUMN enchant_code VARCHAR(64) NULL REFERENCES content.item_template(code);
