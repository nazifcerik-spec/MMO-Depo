-- Client error reports (DD-235): uncaught app errors grouped by fingerprint.
-- One row per distinct error; no account, no IP.
CREATE TABLE ops.client_error (
    fingerprint CHAR(64)     PRIMARY KEY,
    kind        VARCHAR(16)  NOT NULL,
    message     TEXT         NOT NULL CHECK (length(message) <= 1000),
    stack       TEXT         NOT NULL CHECK (length(stack) <= 8000),
    platform    VARCHAR(16)  NOT NULL,
    app_version VARCHAR(32)  NOT NULL DEFAULT '',
    flavor      VARCHAR(16)  NOT NULL DEFAULT '',
    count       BIGINT       NOT NULL DEFAULT 1,
    first_seen  TIMESTAMPTZ  NOT NULL DEFAULT now(),
    last_seen   TIMESTAMPTZ  NOT NULL DEFAULT now()
);

CREATE INDEX ix_client_error_last_seen ON ops.client_error (last_seen DESC);
