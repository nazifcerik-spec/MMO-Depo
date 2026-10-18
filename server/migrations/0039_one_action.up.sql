-- 0039_one_action.up.sql
-- One action at a time (DD-263) and timed crafting (DD-264).
--
-- A character runs one of gathering, hunt, dungeon, raid, travel or a
-- crafting job at a time. Each action already lives in its own table with
-- an "active" state; character.active_action is their union, and a trigger
-- on every action insert refuses a character that is busy with another
-- kind. pg_advisory_xact_lock(character) serialises concurrent starts, so
-- two different actions can never both slip in. Nothing has to be released:
-- an action stops counting when its own row leaves the active state.

-- --- Timed crafting ----------------------------------------------------------

ALTER TABLE content.recipe
  ADD COLUMN craft_seconds INT NOT NULL DEFAULT 60 CHECK (craft_seconds > 0);

CREATE TABLE crafting.job (
  id            UUID PRIMARY KEY,
  character_id  UUID NOT NULL REFERENCES character.character(id) ON DELETE CASCADE,
  recipe_code   VARCHAR(64) NOT NULL,
  recipe        JSONB NOT NULL,          -- frozen recipe (content reloads cannot change a running job)
  count         SMALLINT NOT NULL CHECK (count BETWEEN 1 AND 100),
  unit_seconds  INT NOT NULL CHECK (unit_seconds > 0),
  started_at    TIMESTAMPTZ NOT NULL,
  ready_at      TIMESTAMPTZ NOT NULL,
  state         VARCHAR(16) NOT NULL DEFAULT 'active' CHECK (state IN ('active', 'collected')),
  made          SMALLINT NULL,           -- units delivered at collection; the rest was refunded
  result        JSONB NULL,
  collected_at  TIMESTAMPTZ NULL
);

CREATE UNIQUE INDEX ux_crafting_job_active ON crafting.job (character_id) WHERE state = 'active';
CREATE INDEX idx_crafting_job_character ON crafting.job (character_id, started_at DESC);

-- --- One action at a time ----------------------------------------------------

CREATE VIEW character.active_action AS
  SELECT character_id, 'gathering'::text AS kind FROM activity.contract WHERE state = 'active'
  UNION ALL SELECT character_id, 'hunt' FROM hunt.contract WHERE state = 'active'
  UNION ALL SELECT character_id, 'dungeon' FROM dungeon.contract WHERE state = 'active'
  UNION ALL SELECT character_id, 'raid' FROM raid.member WHERE active
  UNION ALL SELECT character_id, 'travel' FROM world.travel WHERE state = 'in_progress'
  UNION ALL SELECT character_id, 'crafting' FROM crafting.job WHERE state = 'active';

-- claim_action(kind): the inserted row's character must not be busy with
-- another kind (the same kind is each table's own one-active index). Raises
-- SQLSTATE OT001 with the busy kind as HINT.
CREATE FUNCTION character.claim_action() RETURNS trigger LANGUAGE plpgsql AS $$
DECLARE
  busy text;
BEGIN
  PERFORM pg_advisory_xact_lock(hashtextextended('otter.action:' || NEW.character_id::text, 0));
  SELECT a.kind INTO busy
    FROM character.active_action a
   WHERE a.character_id = NEW.character_id AND a.kind <> TG_ARGV[0]
   LIMIT 1;
  IF busy IS NOT NULL THEN
    RAISE EXCEPTION 'character is busy: %', busy USING ERRCODE = 'OT001', HINT = busy;
  END IF;
  RETURN NEW;
END $$;

CREATE TRIGGER claim_action BEFORE INSERT ON activity.contract
  FOR EACH ROW WHEN (NEW.state = 'active') EXECUTE FUNCTION character.claim_action('gathering');
CREATE TRIGGER claim_action BEFORE INSERT ON hunt.contract
  FOR EACH ROW WHEN (NEW.state = 'active') EXECUTE FUNCTION character.claim_action('hunt');
CREATE TRIGGER claim_action BEFORE INSERT ON dungeon.contract
  FOR EACH ROW WHEN (NEW.state = 'active') EXECUTE FUNCTION character.claim_action('dungeon');
CREATE TRIGGER claim_action BEFORE INSERT ON raid.member
  FOR EACH ROW WHEN (NEW.active) EXECUTE FUNCTION character.claim_action('raid');
CREATE TRIGGER claim_action BEFORE INSERT ON world.travel
  FOR EACH ROW WHEN (NEW.state = 'in_progress') EXECUTE FUNCTION character.claim_action('travel');
CREATE TRIGGER claim_action BEFORE INSERT ON crafting.job
  FOR EACH ROW WHEN (NEW.state = 'active') EXECUTE FUNCTION character.claim_action('crafting');
