DELETE FROM market.listing WHERE item_instance_id IS NULL;
ALTER TABLE market.listing DROP CONSTRAINT listing_item_instance_id_fkey;
ALTER TABLE market.listing ADD CONSTRAINT listing_item_instance_id_fkey
    FOREIGN KEY (item_instance_id) REFERENCES inventory.item_instance(id) ON DELETE CASCADE;
ALTER TABLE market.listing ALTER COLUMN item_instance_id SET NOT NULL;
