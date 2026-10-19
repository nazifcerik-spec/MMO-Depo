-- 0042_quest_levels.down.sql
-- NOT VALID: rows above 60 written meanwhile stay (and keep their progress).
ALTER TABLE content.quest_template DROP CONSTRAINT IF EXISTS quest_template_min_level_check;
ALTER TABLE content.quest_template ADD CONSTRAINT quest_template_min_level_check CHECK (min_level BETWEEN 1 AND 60) NOT VALID;
