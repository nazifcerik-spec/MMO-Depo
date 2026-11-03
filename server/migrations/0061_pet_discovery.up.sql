-- Pet index (DD-342): the species a character has ever held as an egg or a
-- pet. Recorded when the pet screen lists them, so an egg sold or traded
-- later stays discovered.
CREATE TABLE pet.discovery (
  character_id  UUID NOT NULL REFERENCES character.character(id) ON DELETE CASCADE,
  species_code  VARCHAR(32) NOT NULL,
  discovered_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  PRIMARY KEY (character_id, species_code)
);
INSERT INTO pet.discovery (character_id, species_code)
  SELECT DISTINCT character_id, species_code FROM pet.pet
  ON CONFLICT DO NOTHING;
