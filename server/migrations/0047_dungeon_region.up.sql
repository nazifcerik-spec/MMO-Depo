-- 0047_dungeon_region.up.sql
-- Dungeons and raids sit in a region (DD-302): seen and entered from there,
-- from a region linked to it, or from one of its child regions.
ALTER TABLE content.dungeon_template ADD COLUMN region_code VARCHAR(48) NULL;
