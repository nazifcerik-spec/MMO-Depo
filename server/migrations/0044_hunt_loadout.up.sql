-- 0044_hunt_loadout.up.sql
-- A hunt takes potions and dishes from the bag (DD-288). The supplies are
-- held as item instances locked to the hunt (lock_owner = 'hunt',
-- locked_ref_id = contract id); loadout freezes their effects and a food
-- buff already running at the start, so a replay never changes.
ALTER TABLE hunt.contract ADD COLUMN loadout JSONB NULL;
