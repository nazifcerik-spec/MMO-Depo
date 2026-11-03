-- 0064_character_look.up.sql
-- Character customization (DD-346): face, hair style, hair colour, eye colour.
-- Ids are validated against the option counts in the service, so new art only
-- raises a count there.
ALTER TABLE character.character
  ADD COLUMN face_id SMALLINT NOT NULL DEFAULT 1 CHECK (face_id >= 1),
  ADD COLUMN hair_style_id SMALLINT NOT NULL DEFAULT 1 CHECK (hair_style_id >= 1),
  ADD COLUMN hair_color_id SMALLINT NOT NULL DEFAULT 1 CHECK (hair_color_id >= 1),
  ADD COLUMN eye_color_id SMALLINT NOT NULL DEFAULT 1 CHECK (eye_color_id >= 1);
