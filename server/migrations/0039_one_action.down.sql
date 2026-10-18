-- 0039_one_action.down.sql
DROP TRIGGER IF EXISTS claim_action ON crafting.job;
DROP TRIGGER IF EXISTS claim_action ON world.travel;
DROP TRIGGER IF EXISTS claim_action ON raid.member;
DROP TRIGGER IF EXISTS claim_action ON dungeon.contract;
DROP TRIGGER IF EXISTS claim_action ON hunt.contract;
DROP TRIGGER IF EXISTS claim_action ON activity.contract;
DROP FUNCTION IF EXISTS character.claim_action();
DROP VIEW IF EXISTS character.active_action;
DROP TABLE IF EXISTS crafting.job;
ALTER TABLE content.recipe DROP COLUMN IF EXISTS craft_seconds;
