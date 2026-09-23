ALTER TABLE content.item_template DROP COLUMN IF EXISTS use_effect;
ALTER TABLE character.character
  DROP COLUMN IF EXISTS hp_updated_at,
  DROP COLUMN IF EXISTS hp_current;
