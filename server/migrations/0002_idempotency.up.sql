-- 0002_idempotency.up.sql
-- DD-188: her mutating endpoint idempotency_key zorunlu; cache 24h.

CREATE TABLE api.idempotency_record (
  key               VARCHAR(64)  NOT NULL,
  endpoint          VARCHAR(96)  NOT NULL,
  account_id        UUID         NOT NULL,
  request_hash      BYTEA        NOT NULL,
  response_status   SMALLINT     NOT NULL,
  response_body     JSONB        NOT NULL,
  created_at        TIMESTAMPTZ  NOT NULL DEFAULT now(),
  expires_at        TIMESTAMPTZ  NOT NULL DEFAULT now() + INTERVAL '24 hours',
  PRIMARY KEY (key, endpoint, account_id)
);

CREATE INDEX idx_idem_expire ON api.idempotency_record (expires_at);
