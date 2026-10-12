-- Sprint 37 (DD-219): password reset, e-mail verification, e-mail outbox.

-- Single-use account tokens. Only SHA-256(token) is stored; the raw token
-- exists in the e-mail alone. kind: password_reset (1h) | email_verify (24h).
CREATE TABLE identity.account_token (
  id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  account_id  UUID NOT NULL REFERENCES identity.account(id) ON DELETE CASCADE,
  kind        VARCHAR(16) NOT NULL CHECK (kind IN ('password_reset','email_verify')),
  token_hash  BYTEA NOT NULL UNIQUE,
  expires_at  TIMESTAMPTZ NOT NULL,
  used_at     TIMESTAMPTZ NULL,
  created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX idx_account_token_recent ON identity.account_token (account_id, kind, created_at DESC);

ALTER TABLE identity.account ADD COLUMN password_updated_at TIMESTAMPTZ NULL;

-- Transactional e-mail outbox (worker email_dispatch). body is cleared once
-- sent or given up on, so reset links do not linger at rest.
CREATE SCHEMA IF NOT EXISTS ops;
CREATE TABLE ops.email_outbox (
  id          BIGSERIAL PRIMARY KEY,
  to_address  TEXT NOT NULL,
  subject     TEXT NOT NULL,
  body        TEXT NULL,
  attempts    INT NOT NULL DEFAULT 0,
  last_error  TEXT NULL,
  created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
  sent_at     TIMESTAMPTZ NULL,
  failed_at   TIMESTAMPTZ NULL
);
CREATE INDEX idx_email_outbox_pending ON ops.email_outbox (id) WHERE sent_at IS NULL AND failed_at IS NULL;
