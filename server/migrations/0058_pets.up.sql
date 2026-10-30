-- Pets (DD-328): hatched from raid eggs, levelled on hunts of their own; one
-- active pet adds its bonus to the character. A hunt lives on the pet row.
CREATE SCHEMA IF NOT EXISTS pet;
CREATE TABLE pet.pet (
  id              UUID PRIMARY KEY,
  character_id    UUID NOT NULL REFERENCES character.character(id) ON DELETE CASCADE,
  species_code    VARCHAR(32) NOT NULL,
  level           SMALLINT NOT NULL DEFAULT 1 CHECK (level BETWEEN 1 AND 100),
  xp              BIGINT NOT NULL DEFAULT 0 CHECK (xp >= 0),
  hatch_ready_at  TIMESTAMPTZ NOT NULL,
  active          BOOLEAN NOT NULL DEFAULT false,
  hunt_mob        VARCHAR(64) NULL,
  hunt_started_at TIMESTAMPTZ NULL,
  hunt_ready_at   TIMESTAMPTZ NULL,
  hunt_seed       BIGINT NULL,
  created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
  CHECK ((hunt_mob IS NULL) = (hunt_started_at IS NULL))
);
CREATE INDEX ix_pet_character ON pet.pet (character_id);
CREATE UNIQUE INDEX uq_pet_active ON pet.pet (character_id) WHERE active;
