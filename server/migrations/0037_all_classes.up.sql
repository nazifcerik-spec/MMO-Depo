-- 0037_all_classes.up.sql
-- All ten base classes (DD-252). Class definitions, promotions, skills,
-- talents, titles and mastery achievements are content (progression.yaml,
-- skills.yaml; the content pack registers titles/achievements). Here: the
-- v1 class codes move to their progression names and class skills may use
-- the shared effect definitions.
UPDATE character.character SET class_code = 'ranger' WHERE class_code = 'hunter';
UPDATE character.character SET class_code = 'cleric' WHERE class_code = 'priest';
UPDATE character.character SET spec_code = 'ranger' WHERE spec_code = 'hunter';
UPDATE character.character SET spec_code = 'cleric' WHERE spec_code = 'priest';
UPDATE content.item_template SET class_restriction = 'ranger' WHERE class_restriction = 'hunter';
UPDATE content.item_template SET class_restriction = 'cleric' WHERE class_restriction = 'priest';
-- v1 skills of the renamed classes; the content pack writes their successors.
DELETE FROM content.class_skill WHERE class_code IN ('hunter', 'priest');

ALTER TABLE content.class_skill DROP CONSTRAINT IF EXISTS class_skill_effect_check;
ALTER TABLE content.class_skill ADD CONSTRAINT class_skill_effect_check
  CHECK (effect IN ('strike','heal','guard','empower','barrier','dot','hot','weaken','stun'));
