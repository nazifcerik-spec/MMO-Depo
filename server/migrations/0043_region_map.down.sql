-- 0043_region_map.down.sql
ALTER TABLE world.region DROP COLUMN IF EXISTS map_x, DROP COLUMN IF EXISTS map_y;
