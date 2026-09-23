-- 0006_activity.up.sql
-- DD-084 uygulaması: analytical AFK reveal.
--
-- Contract state machine:
--   active     — çalışıyor; reveal öncesi, iptal edilebilir.
--   revealed   — reveal başarıyla tamamlandı; reveal_result cached (idempotency).
--   cancelled  — kullanıcı iptal etti; ödül yok.
--
-- Snapshot: contract start'ta character state dondurulur (DD-088).
-- Seed: (contract_id, character_id) hash — deterministic RNG.
-- Concurrent cap: aktif contract sayısı karakter başına 1 (v1, DD-OPEN-3).

CREATE TABLE activity.contract (
  id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  character_id     UUID NOT NULL REFERENCES character.character(id) ON DELETE CASCADE,
  kind             VARCHAR(32) NOT NULL,       -- 'mining'|'fishing'|'hunting'|'woodcutting'|'herbalism'
  duration_seconds INT NOT NULL CHECK (duration_seconds > 0),
  started_at       TIMESTAMPTZ NOT NULL DEFAULT now(),
  expires_at       TIMESTAMPTZ NOT NULL,       -- started_at + duration_seconds
  seed             BIGINT NOT NULL,            -- deterministic RNG seed for reveal
  state            VARCHAR(16) NOT NULL DEFAULT 'active',
  snapshot         JSONB NOT NULL,             -- frozen character state (class, level, bonus)
  params           JSONB NOT NULL DEFAULT '{}', -- kind-specific (location, node_type, ...)

  -- Reveal cache (idempotency): once revealed, response is frozen.
  revealed_at      TIMESTAMPTZ NULL,
  reveal_result    JSONB NULL,
  reveal_txn_id    UUID NULL,

  -- Cancel bookkeeping.
  cancelled_at     TIMESTAMPTZ NULL,

  created_at       TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at       TIMESTAMPTZ NOT NULL DEFAULT now(),
  version          INT NOT NULL DEFAULT 1
);

-- Concurrent slot cap enforcement: at most one active contract per character
-- (v1 — DD-OPEN-3 — Phase 02 opt raises to 1 main + 1 pet + 1 travel).
CREATE UNIQUE INDEX ux_activity_active_per_character
  ON activity.contract (character_id)
  WHERE state = 'active';

-- Worker lookup for expire scanner (Sprint 19).
CREATE INDEX idx_activity_expires
  ON activity.contract (expires_at)
  WHERE state = 'active';

-- Player history query.
CREATE INDEX idx_activity_character_recent
  ON activity.contract (character_id, started_at DESC);
