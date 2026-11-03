ALTER TABLE character.character
  DROP COLUMN IF EXISTS face_id,
  DROP COLUMN IF EXISTS hair_style_id,
  DROP COLUMN IF EXISTS hair_color_id,
  DROP COLUMN IF EXISTS eye_color_id;
