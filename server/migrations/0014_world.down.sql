-- 0014_world.down.sql
ALTER TABLE character.character DROP COLUMN IF EXISTS location_region_code;
DROP TABLE IF EXISTS world.travel;
DROP TABLE IF EXISTS world.region_link;
DROP TABLE IF EXISTS world.region;
