-- A listing outlives its item instance: once the item leaves escrow (sold or
-- cancelled) it may merge into another stack and be deleted. CASCADE used to
-- delete the listing row with it (lost sale history, cancel answered 404).
ALTER TABLE market.listing ALTER COLUMN item_instance_id DROP NOT NULL;
ALTER TABLE market.listing DROP CONSTRAINT listing_item_instance_id_fkey;
ALTER TABLE market.listing ADD CONSTRAINT listing_item_instance_id_fkey
    FOREIGN KEY (item_instance_id) REFERENCES inventory.item_instance(id) ON DELETE SET NULL;
