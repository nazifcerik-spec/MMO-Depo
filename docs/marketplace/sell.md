# Sell (Post to Market)

## 1. Amaç
Bir itemi belirli fiyattan satışa çıkarma. Escrow lock + listing insert + listing fee tahsili tek TX'te.

## 2. Endpoint

**POST `/v1/market/listing`**
```json
{
  "item_instance_id": "uuid",         // unique/non-stackable
  "stack_quantity": 1,                // stackable için >1
  "unit_price_copper": 4500,
  "duration_hours": 24,
  "idempotency_key": "..."
}
```

## 3. Akış (Server)

Tek SERIALIZABLE TX:

1. **Idempotency check** — `idempotency_key` varsa cached response dön.
2. **Item ownership** — item_instance owner == seller.
3. **Bind rule** — `is_no_trade=false`, `bind_kind != 'bop'` (bkz. [item-engine](../item-engine/README.md)).
4. **Not already locked** — `is_locked=false`.
5. **Not quest / soulbound** — `is_no_sell` benzer flag'ler.
6. **Price sınırları** — `unit_price_copper ∈ [min, max]` (bkz. [price.md](price.md)).
7. **Duration valid** — {12, 24, 48, 168} h.
8. **Listing cap** — seller aktif listing sayısı < 50.
9. **Listing fee hesapla** — bkz. [tax.md](tax.md).
10. **Wallet yeterli fee için** — `wallet.gold_copper ≥ listing_fee`.
11. **Item escrow lock** — `is_locked=true, lock_owner='market'`.
12. **Insert `market.listing`** state='active'.
13. **Insert `market.escrow_item`**.
14. **Insert `market.listing_snapshot`** (item instance JSONB'sini snapshot al).
15. **Ledger (fee)**:
    ```
    kind='market_listing_fee'
    -listing_fee → seller.wallet.gold
    +listing_fee → system.tax_sink
    ref: listing_id
    ```
16. **Outbox** → search index invalidation event + notification (opsiyonel "listing active" WS satıcıya).
17. **COMMIT.**

## 4. Response

```json
{
  "listing_id": "uuid",
  "expires_at": "2026-09-10T12:34:56Z",
  "listing_fee_paid_copper": 45,
  "estimated_net_if_sold": 42750,   // buyer payment - sale cut
  "state": "active"
}
```

## 5. Same-Account Guard (v1)

Post sırasında **kontrol yok** (satıcı → alıcı kontrolü buy zamanında).
Buy sırasında `seller_account_id == buyer_account_id` → `E_SAME_ACCOUNT_BLOCK` (bkz. [buy.md](buy.md)).

Gerekçe: Post zamanında bilinemez; buy zamanında zorunlu.

## 6. Stackable Item

- `item_template.stackable=true` ise `stack_quantity ≥ 1`.
- Envanterdeki toplam stack ≥ `stack_quantity` olmalı.
- Post sonrası envanterden `stack_quantity` düşülür.
- Alıcı **tümünü** alır; kısmi buy v1'de yok (Phase 02: partial buy opt).

## 7. Unique Item

- Craft mastery signature veya legendary unique instance:
  - Item snapshot signature (crafter_id, mastery, crit) UI'da gösterilir.
  - Fiyat lower bound genellikle yüksek (mastery premium — bkz. [price.md](price.md)).

## 8. Hata Kodları

| Kod | Açıklama |
|---|---|
| `E_ITEM_NOT_FOUND` | Instance yok / farklı owner |
| `E_ITEM_LOCKED` | Zaten locked |
| `E_ITEM_NOT_TRADABLE` | BoP / soulbound |
| `E_PRICE_OUT_OF_RANGE` | Min/max clamp dışı |
| `E_INVALID_DURATION` | Süre desteklenmeyen |
| `E_MARKET_LISTING_CAP` | 50 aktif cap doldu |
| `E_INSUFFICIENT_BALANCE` | Fee ödenemez |
| `E_IDEMPOTENT_REPLAY` | Aynı key yeniden (cached response) |

## 9. Metrikler
- `market_listing_posted_total{item_template, rarity}`
- `market_listing_fee_collected_copper_total`
- `market_post_latency_ms`
- `market_post_error_total{code}`

## 10. Cross-Reference
- [listing.md](listing.md), [buy.md](buy.md), [cancel.md](cancel.md), [tax.md](tax.md), [price.md](price.md)
- [inventory/item-lock.md](../inventory/item-lock.md)
- [economy/transaction.md](../economy/transaction.md)

## 11. Bekleyen Kararlar
- Partial buy (stackable) Phase 02.
- Bulk post (birden fazla listing tek istek) Phase 02.
- Craft signature premium formula (Phase 02 mastery pass).
