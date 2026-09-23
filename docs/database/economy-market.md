# Database — `economy` Schema

## Kapsam
Currency ledger (kaynak-of-truth), marketplace, mail, market watchlist. Ekonomi audit'in temeli.

---

## 1. `economy.ledger_entry`

Currency ve item transfer'lerinin append-only kaydı. Wallet balance bu kayıtların projeksiyonu (verify job ile).

| Kolon | Tip | Constraint |
|---|---|---|
| `id` | `BIGSERIAL` | PK |
| `character_id` | `UUID` | NOT NULL |
| `currency_code` | `VARCHAR(24)` | NOT NULL |
| `delta` | `BIGINT` | NOT NULL — positive credit, negative debit |
| `balance_after` | `BIGINT` | NOT NULL — reconcile check |
| `reason` | `economy.ledger_reason` | NOT NULL |
| `source_kind` | `VARCHAR(32)` | NOT NULL — activity_reveal, quest, boss, mob, market_sale, market_purchase, repair, teleport, listing_fee, sale_fee, respec, guild_deposit, admin |
| `source_ref_id` | `UUID` | NULL — polymorphic reference |
| `correlation_id` | `UUID` | NULL — request/event correlation |
| `idempotency_key` | `TEXT` | NOT NULL UNIQUE |
| `applied_at` | `TIMESTAMPTZ` | NOT NULL DEFAULT now() |
| `meta_json` | `JSONB` | NOT NULL DEFAULT `'{}'` |

Enum:
```sql
CREATE TYPE economy.ledger_reason AS ENUM (
  'credit_activity','credit_quest','credit_boss','credit_mob','credit_event',
  'credit_market_sale','credit_mail','credit_admin','credit_ascension_refund',
  'debit_vendor_buy','debit_repair','debit_teleport','debit_respec',
  'debit_listing_fee','debit_sale_fee','debit_mail_send','debit_guild_deposit',
  'debit_purchase_service','debit_admin','debit_marketplace_purchase',
  'transfer_mail','transfer_market'
);
```

Indexler:
- `uq_ledger_entry__idempotency_key`.
- `ix_ledger_entry__character_id__applied_at DESC` — history read.
- `ix_ledger_entry__correlation_id`.
- `ix_ledger_entry__applied_at` (BRIN, büyük tablo).

Partitioning (Phase 02): month-range on `applied_at`.

Notlar:
- **Append-only**: UPDATE/DELETE yasak (permission).
- Reconcile job: `SUM(delta) per (character_id, currency_code) == character.wallet.balance`; fark alarm.
- Rollback: yeni compensating entry ile (never delete).

---

## 2. `economy.item_transfer_log`

Item ownership transferlerinin audit log'u.

| Kolon | Tip |
|---|---|
| `id` | `BIGSERIAL` PK |
| `item_instance_id` | `UUID` NOT NULL |
| `from_owner_kind` | `content.item_owner_kind` NOT NULL |
| `from_owner_id` | `UUID` NOT NULL |
| `to_owner_kind` | `content.item_owner_kind` NOT NULL |
| `to_owner_id` | `UUID` NOT NULL |
| `reason` | `VARCHAR(32)` NOT NULL — market_sale, mail, loot, craft_output, gm_grant, salvage_consume |
| `correlation_id` | `UUID` NULL |
| `applied_at` | `TIMESTAMPTZ` NOT NULL DEFAULT now() |
| `meta_json` | `JSONB` NOT NULL DEFAULT `'{}'` |

Indexler: `ix_item_transfer__item_instance_id__applied_at DESC`.

---

## 3. `economy.market_listing`

Aktif ve tamamlanmış listing'ler.

| Kolon | Tip | Constraint |
|---|---|---|
| `id` | `UUID` | PK |
| `seller_character_id` | `UUID` | NOT NULL |
| `item_instance_id` | `UUID` | NOT NULL UNIQUE — aktif iken |
| `item_template_id` | `UUID` | NOT NULL — search index için denormalize |
| `qty` | `INT` | NOT NULL, `CHECK qty > 0` |
| `unit_price_copper` | `BIGINT` | NOT NULL, `CHECK unit_price_copper > 0` |
| `listing_fee_paid_copper` | `BIGINT` | NOT NULL |
| `sale_fee_pct` | `SMALLINT` | NOT NULL — snapshot değeri (config değişse de listing devam) |
| `duration_ms` | `INT` | NOT NULL — 12h/24h/48h |
| `created_at` | `TIMESTAMPTZ` | NOT NULL DEFAULT now() |
| `expires_at` | `TIMESTAMPTZ` | NOT NULL |
| `status` | `economy.listing_status` | NOT NULL DEFAULT `'active'` |
| `sold_at` | `TIMESTAMPTZ` | NULL |
| `cancelled_at` | `TIMESTAMPTZ` | NULL |
| `version` | `INT` | NOT NULL DEFAULT 1 |

Enum:
```sql
CREATE TYPE economy.listing_status AS ENUM ('active','sold','expired','cancelled');
```

Indexler:
- `ix_market_listing__item_template_id__status__unit_price_copper` — search hot path.
- `ix_market_listing__seller_character_id__status`.
- `ix_market_listing__expires_at` (worker sweep).
- `ix_market_listing__item_instance_id` (UNIQUE partial `WHERE status = 'active'`).

---

## 4. `economy.market_transaction`

Tamamlanmış alım-satım.

| Kolon | Tip |
|---|---|
| `id` | `UUID` PK |
| `listing_id` | `UUID` FK NOT NULL |
| `buyer_character_id` | `UUID` NOT NULL |
| `seller_character_id` | `UUID` NOT NULL |
| `item_instance_id` | `UUID` NOT NULL |
| `qty` | `INT` NOT NULL |
| `unit_price_copper` | `BIGINT` NOT NULL |
| `total_gross_copper` | `BIGINT` NOT NULL |
| `sale_fee_copper` | `BIGINT` NOT NULL |
| `net_to_seller_copper` | `BIGINT` NOT NULL |
| `executed_at` | `TIMESTAMPTZ` NOT NULL DEFAULT now() |
| `correlation_id` | `UUID` NOT NULL |

Indexler:
- `ix_market_transaction__buyer_character_id__executed_at DESC`.
- `ix_market_transaction__seller_character_id__executed_at DESC`.
- `ix_market_transaction__item_template_id` (denormalize — pricing analytics).

---

## 5. `economy.market_watchlist` (Phase 02)

```
id UUID PK,
character_id UUID FK NOT NULL,
item_template_id UUID FK NOT NULL,
max_unit_price_copper BIGINT NULL,
created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
UNIQUE (character_id, item_template_id)
```

---

## 6. `economy.commodity_price_index`

Ekonomi metrikleri (worker günlük hesaplama).

| Kolon | Tip |
|---|---|
| `day` | `DATE` PK-part |
| `currency_code` | `VARCHAR(24)` PK-part |
| `item_template_id` | `UUID` PK-part |
| `avg_unit_price_copper` | `BIGINT` NOT NULL |
| `median_unit_price_copper` | `BIGINT` NOT NULL |
| `min_unit_price_copper` | `BIGINT` NOT NULL |
| `max_unit_price_copper` | `BIGINT` NOT NULL |
| `volume` | `INT` NOT NULL |
| `active_listings_end_of_day` | `INT` NOT NULL |

---

## 7. `economy.currency_supply_snapshot`

Toplam para arzı takibi — enflasyon monitoring.

| Kolon | Tip |
|---|---|
| `snapshot_at` | `TIMESTAMPTZ` PK-part |
| `currency_code` | `VARCHAR(24)` PK-part |
| `total_supply` | `BIGINT` NOT NULL |
| `active_wallets_count` | `INT` NOT NULL |
| `median_wallet_balance` | `BIGINT` NOT NULL |

---

## 8. `economy.mail`

Karakter mail sistemi (item + gold attach).

| Kolon | Tip | Constraint |
|---|---|---|
| `id` | `UUID` | PK |
| `sender_kind` | `economy.mail_sender_kind` | NOT NULL |
| `sender_character_id` | `UUID` | NULL — player mail için |
| `sender_system_key` | `VARCHAR(48)` | NULL — 'market_sale', 'overflow_inventory', ... |
| `receiver_character_id` | `UUID` | NOT NULL |
| `subject` | `VARCHAR(96)` | NOT NULL |
| `body` | `TEXT` | NOT NULL DEFAULT '' |
| `attachments_json` | `JSONB` | NOT NULL DEFAULT `'{}'` — {items:[item_instance_id], gold_copper: N} |
| `cost_gold_copper` | `BIGINT` | NOT NULL DEFAULT 0 — sender ödedi |
| `read` | `BOOLEAN` | NOT NULL DEFAULT false |
| `claimed` | `BOOLEAN` | NOT NULL DEFAULT false — attachments alındı mı |
| `expires_at` | `TIMESTAMPTZ` | NOT NULL |
| `sent_at` | `TIMESTAMPTZ` | NOT NULL DEFAULT now() |
| `read_at` | `TIMESTAMPTZ` | NULL |
| `claimed_at` | `TIMESTAMPTZ` | NULL |
| `returned_to_sender_at` | `TIMESTAMPTZ` | NULL — expire olunca sender'a döner |

Enum:
```sql
CREATE TYPE economy.mail_sender_kind AS ENUM ('player','system','guild','event');
```

Indexler:
- `ix_mail__receiver_character_id__sent_at DESC`.
- `ix_mail__expires_at` (worker).

Notlar:
- Mail atomicity: gönderme sırasında sender's ledger debit + item ownership transfer to mailbox owner (`content.item_instance.owner_kind='mailbox'`).
- Expire flow: 30 gün sonra unread mail sender'a iade (item + gold).

---

## 9. Kullanım Örnekleri

**Marketplace search:**
```sql
SELECT id, seller_character_id, item_instance_id, qty, unit_price_copper, expires_at
FROM economy.market_listing
WHERE item_template_id = $1
  AND status = 'active'
ORDER BY unit_price_copper ASC, expires_at ASC
LIMIT 50;
```

**Buy transaction:**
```sql
BEGIN ISOLATION LEVEL SERIALIZABLE;

-- 1. Load listing FOR UPDATE (optimistic version check)
SELECT * FROM economy.market_listing WHERE id = $1 AND status = 'active' FOR UPDATE;

-- 2. Debit buyer wallet
INSERT INTO economy.ledger_entry (character_id, currency_code, delta, balance_after, reason, source_kind, source_ref_id, idempotency_key)
VALUES ($buyer, 'gold', -$total_gross, $new_balance, 'debit_marketplace_purchase', 'market_listing', $listing_id, $idem);

-- 3. Update listing status
UPDATE economy.market_listing
SET status='sold', sold_at=now(), version=version+1
WHERE id=$1 AND version=$v;

-- 4. Insert transaction
INSERT INTO economy.market_transaction (...) VALUES (...);

-- 5. Deliver item to buyer via mail (BoE safe path)
INSERT INTO economy.mail (sender_kind='system', receiver=$buyer, attachments={items:[item_instance_id], gold_copper:0}, ...);
UPDATE content.item_instance SET owner_kind='mailbox', owner_id=$mail_id WHERE id=$item_instance_id;

-- 6. Credit seller (net)
INSERT INTO economy.ledger_entry (character_id=$seller, delta=+$net, reason='credit_market_sale', ...);

COMMIT;
```

---

## 10. İlişkiler
- Ledger → character (character_id), source events (activity.reveal_log, quest, mail, market).
- Market listing → item_instance, character.
- Mail → item_instance (attached), character.wallet (via ledger entry).
- Item transfer log → item_instance, mixed owner references.

---

## 11. Notlar
- **Serializable isolation** market_transaction ve ledger_entry insert path'inde kritik.
- **Currency code registry**: `content.currency_definition`.
- **Cent/copper strategy**: bkz. DD-open. Öneri: 1 gold = 10000 copper, tüm ledger `BIGINT copper`.
- **Anti-cheat**: reconcile mismatch alarm; sale fee override attempt log; rapid listing/cancel pattern flag.

---

## 12. Bekleyen Kararlar
- Sale fee tam oran (%5 veya %7) — Phase 01 balance.
- Auction / bid mode (Phase 02).
- Cross-region market (Phase 03).
- Escrow direct trade (Phase 02 opsiyonel).
