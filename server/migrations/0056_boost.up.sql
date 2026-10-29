-- Timed XP boosts (DD-320): +100% XP for two hours in one category, from a
-- verified rewarded ad or the credit shop. ref keys a grant once.
CREATE SCHEMA IF NOT EXISTS boost;
CREATE TABLE boost.grant (
  id           BIGSERIAL PRIMARY KEY,
  character_id UUID NOT NULL REFERENCES character.character(id) ON DELETE CASCADE,
  category     VARCHAR(16) NOT NULL CHECK (category IN ('gather', 'adventure', 'craft')),
  source       VARCHAR(8) NOT NULL CHECK (source IN ('ad', 'credit')),
  ref          VARCHAR(128) NOT NULL UNIQUE,
  starts_at    TIMESTAMPTZ NOT NULL,
  ends_at      TIMESTAMPTZ NOT NULL
);
CREATE INDEX idx_boost_char ON boost.grant (character_id, category, ends_at);
