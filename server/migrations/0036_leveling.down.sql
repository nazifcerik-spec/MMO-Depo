-- 0036_leveling.down.sql
DELETE FROM progression.character_title WHERE title_code IN
  ('veteran','elite','champion','hero','legendary','mythic','ancient','ascendant','immortal','eternal',
   'the_unbreakable','lord_of_battle','bloodbound','grandmaster_of_arms');
UPDATE character.character SET active_title_code = NULL WHERE active_title_code IN
  ('veteran','elite','champion','hero','legendary','mythic','ancient','ascendant','immortal','eternal',
   'the_unbreakable','lord_of_battle','bloodbound','grandmaster_of_arms');
DELETE FROM content.title_template WHERE source_kind IN ('level','class_mastery');

DELETE FROM content.class_skill WHERE spec_code IS NOT NULL OR unlock_level > 60 OR effect IN ('guard','empower');
ALTER TABLE content.class_skill DROP COLUMN config, DROP COLUMN spec_code;
ALTER TABLE content.class_skill DROP CONSTRAINT IF EXISTS class_skill_effect_check;
ALTER TABLE content.class_skill ADD CONSTRAINT class_skill_effect_check CHECK (effect IN ('strike','heal'));
ALTER TABLE content.class_skill DROP CONSTRAINT IF EXISTS class_skill_unlock_level_check;
ALTER TABLE content.class_skill ADD CONSTRAINT class_skill_unlock_level_check CHECK (unlock_level BETWEEN 1 AND 60);

ALTER TABLE content.item_template DROP COLUMN requirements;
DROP TABLE character.talent;
DROP TABLE character.class_history;
DROP TABLE character.progression_event;
DROP TABLE character.xp_grant;

ALTER TABLE character.character
  DROP COLUMN postcap_xp, DROP COLUMN promotion_overflow_xp,
  DROP COLUMN attr_luck, DROP COLUMN attr_wis, DROP COLUMN attr_int,
  DROP COLUMN attr_dex, DROP COLUMN attr_vit, DROP COLUMN attr_str;
UPDATE character.character SET level = 60 WHERE level > 60;
ALTER TABLE character.character DROP CONSTRAINT IF EXISTS character_level_check;
ALTER TABLE character.character ADD CONSTRAINT character_level_check CHECK (level >= 1 AND level <= 60);
