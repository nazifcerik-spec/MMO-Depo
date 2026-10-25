-- Return any held supplies to the bag before the column goes.
UPDATE inventory.item_instance
   SET is_locked = false, lock_owner = NULL, locked_ref_id = NULL, slot = NULL
 WHERE lock_owner = 'hunt';
ALTER TABLE hunt.contract DROP COLUMN IF EXISTS loadout;
