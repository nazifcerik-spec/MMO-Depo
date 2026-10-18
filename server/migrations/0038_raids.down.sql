-- 0038_raids.down.sql
DROP SCHEMA IF EXISTS raid CASCADE;
DELETE FROM dungeon.weekly_lockout
 WHERE dungeon_code IN (SELECT code FROM content.dungeon_template WHERE group_max > 1);
DELETE FROM content.dungeon_template WHERE group_max > 1;
DELETE FROM content.mob_template WHERE raid_boss;
ALTER TABLE content.mob_template DROP COLUMN IF EXISTS raid_boss;
ALTER TABLE content.dungeon_template
  DROP CONSTRAINT IF EXISTS dungeon_template_group_check,
  DROP COLUMN IF EXISTS group_max,
  DROP COLUMN IF EXISTS group_min;
