-- 0063_character_body.up.sql
-- Full-body portrait (DD-344): the figure a character is drawn with.
ALTER TABLE character.character
  ADD COLUMN body VARCHAR(8) NOT NULL DEFAULT 'male' CHECK (body IN ('male', 'female'));
