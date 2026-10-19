-- 0040_food_buff.up.sql
-- Food buffs (DD-269): a dish can add attack and defense percent for a while.
-- One food buff per character: eating again replaces it. The character
-- sheet reads the live row (expires_at in the future) as a modifier, so a
-- hunt or dungeon started while fed keeps the bonus for the whole run.
CREATE TABLE character.buff (
  character_id     UUID PRIMARY KEY REFERENCES character.character(id) ON DELETE CASCADE,
  source_template  VARCHAR(64) NOT NULL,
  attack_pct       SMALLINT NOT NULL DEFAULT 0 CHECK (attack_pct BETWEEN 0 AND 100),
  defense_pct      SMALLINT NOT NULL DEFAULT 0 CHECK (defense_pct BETWEEN 0 AND 100),
  expires_at       TIMESTAMPTZ NOT NULL
);