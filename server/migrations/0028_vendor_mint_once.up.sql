-- 0028_vendor_mint_once.up.sql
-- A vendor purchase or buyback mints exactly one item per source reference
-- (the ledger txn of the purchase / the buyback entry). Retries and
-- concurrent twins lose on this index instead of minting a free copy.
CREATE UNIQUE INDEX uq_item_instance_vendor_source ON inventory.item_instance (source_ref_id)
  WHERE source_kind IN ('vendor_buy', 'vendor_buyback') AND source_ref_id IS NOT NULL;
