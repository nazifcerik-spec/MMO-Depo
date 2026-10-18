-- 0037_all_classes.down.sql
DELETE FROM content.class_skill
 WHERE effect IN ('barrier','dot','hot','weaken','stun')
    OR class_code NOT IN ('warrior','mage','hunter','priest','rogue');
ALTER TABLE content.class_skill DROP CONSTRAINT IF EXISTS class_skill_effect_check;
ALTER TABLE content.class_skill ADD CONSTRAINT class_skill_effect_check CHECK (effect IN ('strike','heal','guard','empower'));

-- Mastery achievements/titles registered by the content pack (DD-252).
DELETE FROM progression.achievement_progress
 WHERE achievement_code IN (SELECT code FROM content.achievement_template WHERE category = 'mastery');
DELETE FROM content.achievement_template WHERE category = 'mastery';
UPDATE character.character SET active_title_code = NULL
 WHERE active_title_code IN (SELECT code FROM content.title_template WHERE source_kind = 'class_mastery'
                               AND code NOT IN ('the_unbreakable','lord_of_battle','bloodbound','grandmaster_of_arms'));
DELETE FROM progression.character_title
 WHERE title_code IN (SELECT code FROM content.title_template WHERE source_kind = 'class_mastery'
                        AND code NOT IN ('the_unbreakable','lord_of_battle','bloodbound','grandmaster_of_arms'));
DELETE FROM content.title_template WHERE source_kind = 'class_mastery'
   AND code NOT IN ('the_unbreakable','lord_of_battle','bloodbound','grandmaster_of_arms');

UPDATE content.item_template SET class_restriction = 'hunter' WHERE class_restriction = 'ranger';
UPDATE content.item_template SET class_restriction = 'priest' WHERE class_restriction = 'cleric';
UPDATE character.character SET spec_code = 'hunter' WHERE spec_code = 'ranger';
UPDATE character.character SET spec_code = 'priest' WHERE spec_code = 'cleric';
UPDATE character.character SET class_code = 'hunter' WHERE class_code = 'ranger';
UPDATE character.character SET class_code = 'priest' WHERE class_code = 'cleric';
