-- Credit shop + cosmetics (DD-321): credits are an account currency; the
-- account's bought cosmetics and its chosen avatar and background.
INSERT INTO economy.currency (code, display_name_key, scope, tradable, cap, base_unit_ratio) VALUES
  ('credit', 'currency.credit', 'account', false, 100000000, 1)
ON CONFLICT (code) DO NOTHING;
CREATE SCHEMA IF NOT EXISTS cosmetic;
CREATE TABLE cosmetic.owned (
  account_id UUID NOT NULL REFERENCES identity.account(id) ON DELETE CASCADE,
  code       VARCHAR(16) NOT NULL,
  bought_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
  PRIMARY KEY (account_id, code)
);
CREATE TABLE cosmetic.profile (
  account_id UUID PRIMARY KEY REFERENCES identity.account(id) ON DELETE CASCADE,
  avatar     VARCHAR(16) NOT NULL DEFAULT 'av01',
  background VARCHAR(16) NOT NULL DEFAULT 'bg01'
);
