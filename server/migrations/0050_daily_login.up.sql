-- Daily login reward (DD-311): one claim per account and UTC day on a
-- seven-day cycle (streak = the cycle day last paid).
CREATE TABLE progression.daily_login (
  account_id UUID PRIMARY KEY,
  last_day   DATE NOT NULL,
  streak     SMALLINT NOT NULL CHECK (streak BETWEEN 0 AND 7),
  claims     INT NOT NULL DEFAULT 0,
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
