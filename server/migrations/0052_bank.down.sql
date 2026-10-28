DELETE FROM inventory.item_instance WHERE bank_account_id IS NOT NULL;
DROP INDEX IF EXISTS inventory.uq_item_bank_slot;
ALTER TABLE inventory.item_instance
  DROP CONSTRAINT IF EXISTS item_bank_or_owner,
  DROP CONSTRAINT IF EXISTS item_bank_both,
  DROP COLUMN IF EXISTS bank_slot,
  DROP COLUMN IF EXISTS bank_account_id;
