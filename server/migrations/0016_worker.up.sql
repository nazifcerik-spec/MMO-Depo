-- 0016_worker.up.sql
-- Worker jobs (Sprint 19).
--
-- ready_at: `activity_expire` job tags active AFK contracts whose duration has
-- elapsed ("scan + tag", docs/performance/load-afk-calc.md §2.3). The
-- notification sprint consumes these rows for "your hunt is ready" pushes.

ALTER TABLE activity.contract ADD COLUMN ready_at TIMESTAMPTZ NULL;
ALTER TABLE hunt.contract     ADD COLUMN ready_at TIMESTAMPTZ NULL;
ALTER TABLE dungeon.contract  ADD COLUMN ready_at TIMESTAMPTZ NULL;
