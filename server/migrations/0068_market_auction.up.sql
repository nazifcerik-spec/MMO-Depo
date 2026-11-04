-- Auctions on the market (DD-359): same listing and item escrow as fixed-price
-- sales; bids hold gold in system.auction_escrow.
ALTER TABLE market.listing
  ADD COLUMN kind VARCHAR(8) NOT NULL DEFAULT 'fixed' CHECK (kind IN ('fixed', 'auction')),
  ADD COLUMN start_bid_copper BIGINT NULL CHECK (start_bid_copper >= 1),
  ADD COLUMN buyout_copper BIGINT NULL CHECK (buyout_copper >= 1),
  ADD COLUMN current_bid_copper BIGINT NULL,
  ADD COLUMN bidder_character_id UUID NULL REFERENCES character.character(id) ON DELETE SET NULL,
  ADD COLUMN bid_count INT NOT NULL DEFAULT 0,
  ADD COLUMN extensions SMALLINT NOT NULL DEFAULT 0;

CREATE INDEX idx_market_auction_due ON market.listing (expires_at)
  WHERE state = 'active' AND kind = 'auction' AND bidder_character_id IS NOT NULL;

CREATE TABLE market.auction_bid (
  id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  listing_id          UUID NOT NULL REFERENCES market.listing(id) ON DELETE CASCADE,
  bidder_character_id UUID NULL REFERENCES character.character(id) ON DELETE SET NULL,
  amount_copper       BIGINT NOT NULL CHECK (amount_copper >= 1),
  -- pending: gold being charged. active: the top bid. outbid: gold due back.
  -- won: paid to the seller. void: never counted, gold due back if charged.
  state               VARCHAR(8) NOT NULL DEFAULT 'pending' CHECK (state IN ('pending', 'active', 'outbid', 'won', 'void')),
  refunded_at         TIMESTAMPTZ NULL,
  created_at          TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_auction_bid_listing ON market.auction_bid (listing_id, created_at DESC);
CREATE INDEX idx_auction_bid_refund ON market.auction_bid (created_at)
  WHERE state IN ('pending', 'outbid', 'void') AND refunded_at IS NULL;
