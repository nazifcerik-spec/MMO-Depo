-- Hired NPC helpers (DD-319): a character hires helpers and sends each on
-- its own gathering or crafting job beside its own work.
CREATE SCHEMA IF NOT EXISTS helper;
CREATE TABLE helper.hired (
  character_id UUID NOT NULL REFERENCES character.character(id) ON DELETE CASCADE,
  type_code    VARCHAR(32) NOT NULL,
  hired_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
  PRIMARY KEY (character_id, type_code)
);
CREATE TABLE helper.job (
  id           UUID PRIMARY KEY,
  character_id UUID NOT NULL REFERENCES character.character(id) ON DELETE CASCADE,
  type_code    VARCHAR(32) NOT NULL,
  kind         VARCHAR(8) NOT NULL CHECK (kind IN ('gather', 'craft')),
  spec         JSONB NOT NULL,
  seed         BIGINT NOT NULL,
  state        VARCHAR(12) NOT NULL DEFAULT 'running' CHECK (state IN ('running', 'done', 'cancelled')),
  started_at   TIMESTAMPTZ NOT NULL,
  ready_at     TIMESTAMPTZ NOT NULL,
  collected_at TIMESTAMPTZ NULL,
  result       JSONB NULL
);
CREATE UNIQUE INDEX uq_helper_job_running ON helper.job (character_id, type_code) WHERE state = 'running';
