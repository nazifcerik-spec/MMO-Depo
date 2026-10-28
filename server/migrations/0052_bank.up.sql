-- Account bank (DD-314): a banked item belongs to the account, not to a
-- character, in one of 200 slots every character of the account shares.
ALTER TABLE inventory.item_instance
  ADD COLUMN bank_account_id UUID NULL REFERENCES identity.account(id) ON DELETE CASCADE,
  ADD COLUMN bank_slot SMALLINT NULL CHECK (bank_slot BETWEEN 0 AND 199),
  ADD CONSTRAINT item_bank_both CHECK ((bank_account_id IS NULL) = (bank_slot IS NULL)),
  ADD CONSTRAINT item_bank_or_owner CHECK (bank_account_id IS NULL OR owner_character_id IS NULL);
CREATE UNIQUE INDEX uq_item_bank_slot ON inventory.item_instance (bank_account_id, bank_slot)
  WHERE bank_account_id IS NOT NULL;
