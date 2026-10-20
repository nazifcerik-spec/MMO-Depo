-- 0043_region_map.up.sql
-- World map (DD-279): each region's position on the travel map, in map
-- units (content/map.yaml). The client draws the map from these.
ALTER TABLE world.region ADD COLUMN map_x SMALLINT, ADD COLUMN map_y SMALLINT;
