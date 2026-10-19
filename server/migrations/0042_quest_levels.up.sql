-- 0042_quest_levels.up.sql
-- Quest chains run to level 1000 (DD-278): the 1..60 cap dates from the
-- level-60 era (DD-250 raised the character cap).
ALTER TABLE content.quest_template DROP CONSTRAINT IF EXISTS quest_template_min_level_check;
ALTER TABLE content.quest_template ADD CONSTRAINT quest_template_min_level_check CHECK (min_level BETWEEN 1 AND 1000);
