-- 0012_hunt.up.sql
-- Hunt — AFK combat encounter contracts (Sprint 14).
--
-- Ayrı bir tablo kullanıyoruz (activity.contract'a kolon eklemek yerine)
-- çünkü hunt lifecycle farklı: mid-hunt player death, per-encounter cadence,
-- player_snapshot ile frozen sheet. Ortak yanı sadece "başla → reveal → grant"
-- akışı.

CREATE SCHEMA IF NOT EXISTS hunt;

CREATE TABLE hunt.contract (
  id                 UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  character_id       UUID NOT NULL REFERENCES character.character(id) ON DELETE CASCADE,
  mob_code           VARCHAR(48) NOT NULL,
  encounter_seconds  INT NOT NULL CHECK (encounter_seconds >= 30),
  duration_seconds   INT NOT NULL CHECK (duration_seconds > 0 AND duration_seconds <= 43200),
  started_at         TIMESTAMPTZ NOT NULL DEFAULT now(),
  expires_at         TIMESTAMPTZ NOT NULL,
  seed               BIGINT NOT NULL,
  state              VARCHAR(16) NOT NULL DEFAULT 'active'
                     CHECK (state IN ('active','revealed','cancelled')),
  player_snapshot    JSONB NOT NULL,               -- combat.Combatant frozen at Start
  revealed_at        TIMESTAMPTZ NULL,
  reveal_result      JSONB NULL,                    -- HuntResult, cached for replay
  reveal_txn_id      UUID NULL,                     -- ledger txn from gold grant
  cancelled_at       TIMESTAMPTZ NULL,
  created_at         TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at         TIMESTAMPTZ NOT NULL DEFAULT now(),
  version            INT NOT NULL DEFAULT 1
);

-- One active hunt per character (matches activity slot cap = 1).
CREATE UNIQUE INDEX ux_hunt_contract_active
  ON hunt.contract (character_id)
  WHERE state = 'active';

-- History listing (character-scoped, newest first).
CREATE INDEX idx_hunt_contract_history
  ON hunt.contract (character_id, created_at DESC);

-- Worker-friendly (Sprint 19 expire job).
CREATE INDEX idx_hunt_contract_expire
  ON hunt.contract (expires_at)
  WHERE state = 'active';
