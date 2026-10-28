-- Player trading (DD-315): two characters put items up; when both are
-- ready the items change hands and each side pays its tax.
CREATE SCHEMA IF NOT EXISTS trade;
CREATE TABLE trade.session (
  id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  a_character_id UUID NOT NULL REFERENCES character.character(id) ON DELETE CASCADE,
  b_character_id UUID NOT NULL REFERENCES character.character(id) ON DELETE CASCADE,
  state          VARCHAR(12) NOT NULL DEFAULT 'open' CHECK (state IN ('open', 'done', 'cancelled')),
  a_ready        BOOLEAN NOT NULL DEFAULT false,
  b_ready        BOOLEAN NOT NULL DEFAULT false,
  a_tax_copper   BIGINT NOT NULL DEFAULT 0,
  b_tax_copper   BIGINT NOT NULL DEFAULT 0,
  expires_at     TIMESTAMPTZ NOT NULL,
  created_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
  closed_at      TIMESTAMPTZ NULL,
  CHECK (a_character_id <> b_character_id)
);
CREATE INDEX idx_trade_open_a ON trade.session (a_character_id) WHERE state = 'open';
CREATE INDEX idx_trade_open_b ON trade.session (b_character_id) WHERE state = 'open';
CREATE TABLE trade.offer (
  session_id       UUID NOT NULL REFERENCES trade.session(id) ON DELETE CASCADE,
  character_id     UUID NOT NULL,
  item_instance_id UUID NOT NULL REFERENCES inventory.item_instance(id) ON DELETE CASCADE,
  quantity         INT NOT NULL CHECK (quantity > 0),
  PRIMARY KEY (session_id, item_instance_id)
);
