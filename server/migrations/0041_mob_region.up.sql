-- 0041_mob_region.up.sql
-- Hunting is regional (DD-273): a mob lives in one region and is hunted from
-- that region or one linked to it by road. Raid bosses have none.
ALTER TABLE content.mob_template ADD COLUMN region_code VARCHAR(64);
