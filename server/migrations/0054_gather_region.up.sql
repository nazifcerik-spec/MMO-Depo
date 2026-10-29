-- Gathering resources grow on a map (DD-317): gathered only in reach of it.
ALTER TABLE content.item_template ADD COLUMN gather_region VARCHAR(48) NULL;
