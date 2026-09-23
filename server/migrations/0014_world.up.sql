-- 0014_world.up.sql
-- World Engine — region hierarchy + travel (Sprint 16).
--
-- Design:
--   world.region                  region ağacı; parent_code ile self-reference.
--   world.region_link             non-parent adjacency (fast-travel path).
--   world.travel                  aktif in-transit state; slot cap 1 per char.
--   character.character.location_region_code  karakterin son varış bölgesi.
--
-- Region kind: 'continent' | 'zone' | 'city' | 'wilderness' | 'dungeon'.
-- Travel süresi:
--   - parent↔child ise: child.travel_seconds_from_parent
--   - explicit region_link varsa: region_link.travel_seconds
--   - hiçbiri değilse ErrNoRoute
-- MVP'de sadece "aynı parent'ın kardeşleri" ve "parent↔child" yolları
-- doğrudan (kardeşler için parent üzerinden hop). Multi-hop path finder
-- (A*) Phase 02.

-- --- Region tree -----------------------------------------------------------

CREATE TABLE world.region (
  code                       VARCHAR(48) PRIMARY KEY,
  display_name_key           VARCHAR(64) NOT NULL,
  kind                       VARCHAR(16) NOT NULL
                             CHECK (kind IN ('continent','zone','city','wilderness','dungeon')),
  parent_code                VARCHAR(48) NULL REFERENCES world.region(code) ON DELETE RESTRICT,
  min_level                  SMALLINT NOT NULL DEFAULT 1,
  travel_seconds_from_parent INT NOT NULL DEFAULT 60
                             CHECK (travel_seconds_from_parent >= 0
                                    AND travel_seconds_from_parent <= 3600),
  safe                       BOOLEAN NOT NULL DEFAULT true,
  is_starting                BOOLEAN NOT NULL DEFAULT false,   -- yeni karakter buraya spawn
  version                    INT NOT NULL DEFAULT 1,
  created_at                 TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_region_parent ON world.region (parent_code);
-- One starting region cap enforced by partial UNIQUE.
CREATE UNIQUE INDEX ux_region_starting ON world.region ((true)) WHERE is_starting;

-- --- Explicit adjacency (non-parent shortcuts) -----------------------------
-- (a, b) stored with a < b so lookups are order-independent.

CREATE TABLE world.region_link (
  region_a_code   VARCHAR(48) NOT NULL REFERENCES world.region(code) ON DELETE CASCADE,
  region_b_code   VARCHAR(48) NOT NULL REFERENCES world.region(code) ON DELETE CASCADE,
  travel_seconds  INT NOT NULL CHECK (travel_seconds >= 0 AND travel_seconds <= 3600),
  PRIMARY KEY (region_a_code, region_b_code),
  CHECK (region_a_code < region_b_code)     -- canonical ordering
);

CREATE INDEX idx_region_link_a ON world.region_link (region_a_code);
CREATE INDEX idx_region_link_b ON world.region_link (region_b_code);

-- --- Travel state ----------------------------------------------------------

CREATE TABLE world.travel (
  id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  character_id        UUID NOT NULL REFERENCES character.character(id) ON DELETE CASCADE,
  from_region_code    VARCHAR(48) NOT NULL REFERENCES world.region(code) ON DELETE RESTRICT,
  to_region_code      VARCHAR(48) NOT NULL REFERENCES world.region(code) ON DELETE RESTRICT,
  started_at          TIMESTAMPTZ NOT NULL DEFAULT now(),
  arrives_at          TIMESTAMPTZ NOT NULL,
  state               VARCHAR(16) NOT NULL DEFAULT 'in_progress'
                      CHECK (state IN ('in_progress','completed','cancelled')),
  completed_at        TIMESTAMPTZ NULL,
  cancelled_at        TIMESTAMPTZ NULL,
  created_at          TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at          TIMESTAMPTZ NOT NULL DEFAULT now(),
  version             INT NOT NULL DEFAULT 1
);

-- One in-progress travel per character.
CREATE UNIQUE INDEX ux_travel_in_progress
  ON world.travel (character_id)
  WHERE state = 'in_progress';

CREATE INDEX idx_travel_history
  ON world.travel (character_id, started_at DESC);

-- Worker-friendly (auto-complete or timeout scan; Sprint 19).
CREATE INDEX idx_travel_arrives
  ON world.travel (arrives_at)
  WHERE state = 'in_progress';

-- --- Character binding -----------------------------------------------------

ALTER TABLE character.character
  ADD COLUMN location_region_code VARCHAR(48) NULL
    REFERENCES world.region(code) ON DELETE SET NULL;

CREATE INDEX idx_character_location ON character.character (location_region_code)
  WHERE location_region_code IS NOT NULL;

-- --- Seed baseline (content pack YAML in Sprint 12+ is authoritative; this
-- guarantees a fresh DB is playable if the loader has not run yet) ---------

INSERT INTO world.region (code, display_name_key, kind, parent_code, min_level, travel_seconds_from_parent, safe, is_starting) VALUES
  ('sundermere',      'region.sundermere',      'continent',  NULL,           1,   0, true,  false),
  ('ironhold_city',   'region.ironhold_city',   'city',       'sundermere',   1,  60, true,  true),   -- yeni karakter spawn
  ('greenwood',       'region.greenwood',       'wilderness', 'sundermere',   1, 120, false, false),
  ('whispering_reach','region.whispering_reach','wilderness', 'sundermere',   8, 240, false, false),
  ('crypt_gates',     'region.crypt_gates',     'dungeon',    'greenwood',    3, 180, false, false),
  ('emberpeak',       'region.emberpeak',       'zone',       'sundermere',  18, 360, false, false),
  ('ashen_hollow',    'region.ashen_hollow',    'dungeon',    'emberpeak',   24, 300, false, false);

-- Kardeş bölgeler arası shortcut: greenwood ↔ whispering_reach.
INSERT INTO world.region_link (region_a_code, region_b_code, travel_seconds) VALUES
  ('greenwood',  'whispering_reach', 300);
