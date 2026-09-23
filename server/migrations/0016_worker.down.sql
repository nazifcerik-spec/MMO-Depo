-- 0016_worker.down.sql
ALTER TABLE dungeon.contract  DROP COLUMN IF EXISTS ready_at;
ALTER TABLE hunt.contract     DROP COLUMN IF EXISTS ready_at;
ALTER TABLE activity.contract DROP COLUMN IF EXISTS ready_at;
