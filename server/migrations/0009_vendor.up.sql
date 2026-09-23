-- 0009_vendor.up.sql
-- Vendor NPC — player-vendor economic loop.
--
-- v1 kapsamı:
--   ✓ vendor_template: name + type + buy_multiplier + sell_multiplier + sells[]
--   ✓ vendor_buyback: 2h TTL snapshot-based item restore
--   ✓ Buy: infinite stock (v1); sold price = base × buy_mult
--   ✓ Sell: whole item (stackables sold as full stack)
-- v1 dışı (Phase 02+):
--   ✗ Limited stock (daily/weekly refresh)
--   ✗ Reputation discount
--   ✗ Faction-gated vendors
--   ✗ Location-based access (vendor'la aynı lokasyonda olma zorunluluğu)

CREATE TABLE content.vendor_template (
  code             VARCHAR(48) PRIMARY KEY,
  display_name_key VARCHAR(64) NOT NULL,
  vendor_type      VARCHAR(24) NOT NULL,        -- 'general'|'weapon'|'armor'|'alchemy'|'craft_supply'
  buy_multiplier   NUMERIC(5,3) NOT NULL DEFAULT 1.000,   -- vendor→player (default: base price)
  sell_multiplier  NUMERIC(5,3) NOT NULL DEFAULT 0.250,   -- player→vendor (DD-139)
  -- Sells: JSONB array of {template_code, price_override_copper?}
  --   Empty [] = vendor doesn't sell anything, only buys.
  sells            JSONB NOT NULL DEFAULT '[]',
  version          INT NOT NULL DEFAULT 1
);

-- Buyback: player sold an item; can restore within 2h at same price.
-- Snapshot is JSONB copy of item attributes so restore recreates identical instance.
CREATE TABLE economy.vendor_buyback (
  id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  character_id      UUID NOT NULL REFERENCES character.character(id) ON DELETE CASCADE,
  vendor_code       VARCHAR(48) NOT NULL,
  item_snapshot     JSONB NOT NULL,       -- {template_code, rarity, prefix, suffix, stack_quantity}
  sold_price_copper BIGINT NOT NULL,
  sold_at           TIMESTAMPTZ NOT NULL DEFAULT now(),
  expires_at        TIMESTAMPTZ NOT NULL DEFAULT now() + INTERVAL '2 hours',
  claimed_at        TIMESTAMPTZ NULL       -- when buyback happens
);

CREATE INDEX idx_buyback_character_active
  ON economy.vendor_buyback (character_id, sold_at DESC)
  WHERE claimed_at IS NULL;
CREATE INDEX idx_buyback_expire
  ON economy.vendor_buyback (expires_at)
  WHERE claimed_at IS NULL;

-- Seed one general-purpose vendor for Ironhold city (world/location integration Sprint 16).
INSERT INTO content.vendor_template (code, display_name_key, vendor_type, buy_multiplier, sell_multiplier, sells) VALUES
  ('general_ironhold', 'vendor.general_ironhold', 'general', 1.000, 0.250, '[
    {"template_code": "potion_hp_small"},
    {"template_code": "potion_mp_small"},
    {"template_code": "sword_iron"},
    {"template_code": "bow_yew"},
    {"template_code": "staff_apprentice"},
    {"template_code": "chest_leather"},
    {"template_code": "helm_iron"},
    {"template_code": "boots_iron"}
  ]'::jsonb);
