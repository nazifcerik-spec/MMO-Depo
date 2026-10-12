-- 0005_economy.up.sql
-- DD-135 uygulaması: ledger-first ekonomi.
--
-- Şema:
--   economy.currency                       — currency catalog + seed
--   economy.ledger_entry                   — append-only double-entry
--   economy.character_wallet               — hot path gold projection
--   economy.character_currency_balance     — non-gold character currency
--   economy.account_currency_balance       — account-scope currency (gem, vs)
--
-- Kritik constraint: UNIQUE(idempotency_key, entry_seq) — DD-188 uyumlu; aynı
-- key + aynı entry_seq iki kez INSERT edilemez → retry-safe.

CREATE TABLE economy.currency (
  code             VARCHAR(48) PRIMARY KEY,
  display_name_key VARCHAR(64) NOT NULL,
  scope            VARCHAR(16) NOT NULL,           -- 'character' | 'account' | 'guild'
  tradable         BOOLEAN NOT NULL,
  cap              BIGINT  NOT NULL,               -- max balance (per scope entity)
  base_unit_ratio  INT     NOT NULL DEFAULT 1,     -- gold: 10000 (copper), others: 1
  version          INT     NOT NULL DEFAULT 1
);

-- Seed (docs/economy/currency.md tablosu, v1 subset).
-- Yalnızca 'gold' v1'de trade edilebilir (DD-137).
INSERT INTO economy.currency (code, display_name_key, scope, tradable, cap, base_unit_ratio) VALUES
  ('gold',               'currency.gold',            'character', true,  1000000000000, 10000),
  ('gem',                'currency.gem',             'account',   false,       10000000,     1),
  ('honor',              'currency.honor',           'character', false,         100000,     1),
  ('guild_contribution', 'currency.gcontribution',   'guild',     false,        1000000,     1);


CREATE TABLE economy.ledger_entry (
  id               BIGSERIAL PRIMARY KEY,
  txn_id           UUID NOT NULL,
  entry_seq        SMALLINT NOT NULL,               -- 0..N-1 within a txn
  kind             VARCHAR(32) NOT NULL,            -- 'quest_reward'|'vendor_buy'|'mint'|...
  currency_code    VARCHAR(48) NOT NULL REFERENCES economy.currency(code),
  account_kind     VARCHAR(16) NOT NULL,            -- character|account|guild|system
  account_ref      UUID NOT NULL,
  delta            BIGINT NOT NULL,                 -- + credit, - debit
  balance_after    BIGINT NOT NULL,                 -- post-entry projection sanity
  ref_kind         VARCHAR(32) NULL,                -- 'vendor'|'listing'|'character'|...
  ref_id           UUID NULL,
  idempotency_key  VARCHAR(96) NOT NULL,            -- required per TxnService
  created_at       TIMESTAMPTZ NOT NULL DEFAULT now(),
  UNIQUE (idempotency_key, entry_seq)
);

CREATE INDEX idx_ledger_account
  ON economy.ledger_entry (account_kind, account_ref, created_at DESC);
CREATE INDEX idx_ledger_txn ON economy.ledger_entry (txn_id);
CREATE INDEX idx_ledger_ref ON economy.ledger_entry (ref_kind, ref_id);
CREATE INDEX idx_ledger_idem ON economy.ledger_entry (idempotency_key);


-- --- Projection tables ------------------------------------------------------
--
-- Ledger'dan hesaplanabilir; hot path okuma için materialised. Reconciler
-- (Faz 01 sprint 8+) drift ≥ 1c'de character freeze tetikler.

CREATE TABLE economy.character_wallet (
  character_id  UUID PRIMARY KEY REFERENCES character.character(id) ON DELETE CASCADE,
  gold_copper   BIGINT NOT NULL DEFAULT 0 CHECK (gold_copper >= 0),
  updated_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
  version       INT NOT NULL DEFAULT 1
);

CREATE TABLE economy.character_currency_balance (
  character_id  UUID NOT NULL REFERENCES character.character(id) ON DELETE CASCADE,
  currency_code VARCHAR(48) NOT NULL REFERENCES economy.currency(code),
  amount        BIGINT NOT NULL DEFAULT 0 CHECK (amount >= 0),
  updated_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
  version       INT NOT NULL DEFAULT 1,
  PRIMARY KEY (character_id, currency_code)
);

CREATE TABLE economy.account_currency_balance (
  account_id    UUID NOT NULL REFERENCES identity.account(id) ON DELETE CASCADE,
  currency_code VARCHAR(48) NOT NULL REFERENCES economy.currency(code),
  amount        BIGINT NOT NULL DEFAULT 0 CHECK (amount >= 0),
  updated_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
  version       INT NOT NULL DEFAULT 1,
  PRIMARY KEY (account_id, currency_code)
);
