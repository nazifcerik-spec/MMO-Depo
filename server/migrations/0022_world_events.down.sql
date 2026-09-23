-- 0022_world_events.down.sql
DROP SCHEMA IF EXISTS worldevent CASCADE;
DROP TABLE IF EXISTS content.seasonal_event;
DROP TABLE IF EXISTS content.world_event_rotation;
DROP TABLE IF EXISTS content.world_event_template;
DELETE FROM economy.currency WHERE code = 'event_token';
