-- 0019_market.up.sql
-- Marketplace (Sprint 22) — fixed-price listings. docs/marketplace/*.
-- Auctions are Phase 02 (docs/marketplace/README.md §2).
--
-- Escrow (listing.md §Karar v1): the item keeps its owner but is locked
-- (lock_owner='market', locked_ref_id=listing_id). Partial stacks are split
-- into their own locked instance at post time.
--
-- Buy is two-phase so the ledger write and the item hand-over cannot drift:
--   active → pending_sale (buyer reserved) → ledger 'market-buy-<id>' → sold.
-- The worker finishes or reverts pending_sale rows left behind by a crash.

CREATE SCHEMA IF NOT EXISTS market;

CREATE TABLE market.listing (
  id                   UUID PRIMARY KEY,
  seller_character_id  UUID NOT NULL REFERENCES character.character(id) ON DELETE CASCADE,
  item_instance_id     UUID NOT NULL REFERENCES inventory.item_instance(id) ON DELETE CASCADE,
  template_code        VARCHAR(64) NOT NULL,
  item_class           VARCHAR(24) NOT NULL,
  rarity               VARCHAR(16) NOT NULL,
  prefix_code          VARCHAR(48) NULL,
  suffix_code          VARCHAR(48) NULL,
  stack_quantity       INT NOT NULL CHECK (stack_quantity >= 1),
  unit_price_copper    BIGINT NOT NULL CHECK (unit_price_copper >= 1),
  total_price_copper   BIGINT GENERATED ALWAYS AS (unit_price_copper * stack_quantity) STORED,
  listing_fee_copper   BIGINT NOT NULL,
  fee_txn_id           UUID NULL UNIQUE,     -- replay guard for the post
  duration_hours       SMALLINT NOT NULL CHECK (duration_hours IN (12, 24, 48, 168)),
  state                VARCHAR(16) NOT NULL DEFAULT 'active'
                       CHECK (state IN ('active','pending_sale','sold','expired','cancelled')),
  buyer_character_id   UUID NULL REFERENCES character.character(id) ON DELETE SET NULL,
  pending_at           TIMESTAMPTZ NULL,
  created_at           TIMESTAMPTZ NOT NULL DEFAULT now(),
  expires_at           TIMESTAMPTZ NOT NULL,
  closed_at            TIMESTAMPTZ NULL,
  close_reason         VARCHAR(24) NULL
);

CREATE INDEX idx_market_search  ON market.listing (template_code, unit_price_copper) WHERE state = 'active';
CREATE INDEX idx_market_class   ON market.listing (item_class, unit_price_copper)    WHERE state = 'active';
CREATE INDEX idx_market_seller  ON market.listing (seller_character_id, state);
CREATE INDEX idx_market_expires ON market.listing (expires_at) WHERE state = 'active';
CREATE INDEX idx_market_pending ON market.listing (pending_at) WHERE state = 'pending_sale';

CREATE TABLE market.price_history (
  id                BIGSERIAL PRIMARY KEY,
  template_code     VARCHAR(64) NOT NULL,
  unit_price_copper BIGINT NOT NULL,
  quantity          INT NOT NULL,
  sold_at           TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX idx_price_history ON market.price_history (template_code, sold_at DESC);
