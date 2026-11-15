DROP TABLE progression.story_scene;
DELETE FROM progression.character_quest
 WHERE quest_code IN (SELECT code FROM content.quest_template WHERE category = 'story');
DELETE FROM content.quest_template WHERE category = 'story';
ALTER TABLE content.quest_template DROP CONSTRAINT quest_template_category_check;
ALTER TABLE content.quest_template ADD CONSTRAINT quest_template_category_check
    CHECK (category IN ('main','side','daily'));
