-- 0021_progression.down.sql
ALTER TABLE character.character DROP COLUMN IF EXISTS active_title_code, DROP COLUMN IF EXISTS achievement_points;
DROP SCHEMA IF EXISTS progression CASCADE;
DROP TABLE IF EXISTS content.faction;
DROP TABLE IF EXISTS content.achievement_template;
DROP TABLE IF EXISTS content.title_template;
DROP TABLE IF EXISTS content.mastery_track;
