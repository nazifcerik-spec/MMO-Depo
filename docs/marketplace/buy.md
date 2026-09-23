# Buy (Instant Buy)

## 1. Amaç
Aktif listing'i **anında** satın alma. Gold hareketi, item transferi ve sale cut tek TX'te atomik.

## 2. Endpoint

**POST `/v1/market/listing/{id}/buy`**
```json
{
  "buyer_character_id": "uuid",
  "expected_unit_price_copper": 4500,   // price race guard
  "expected_stack_quantity": 1,
  "idempotency_key": "..."
}
```

## 3. Akış (Server)

Tek SERIALIZABLE TX:

1. **Idempotency check**.
2. **Listing SELECT FOR UPDATE** — `state='active'` (aksi → `E_LISTING_UNAVAILABLE`).
3. **Price / quantity race guard** — client'ın gördüğü değerlerle eşleşme (`E_PRICE_MISMATCH`).
4. **Expiry check** — `expires_at > now()` (aksi → `E_LISTING_EXPIRED`; worker temizleyecek zaten).
5. **Same-account guard** — `listing.seller_account_id != buyer.account_id` (v1 zorunlu).
6. **Self-buy guard** — `seller_character_id != buyer_character_id`.
7. **Bind rule** — item hâlâ trade edilebilir (double-check).
8. **Envanter kapasite** — alıcı için slot / weight / stack limit yeterli.
9. **Sale cut hesapla** — bkz. [tax.md](tax.md).
    ```
    net_to_seller = total_price - sale_cut
    ```
10. **Wallet check** — `buyer.wallet.gold_copper ≥ total_price`.
11. **Ledger (atomic)**:
    ```
    kind='market_buy'
    entries:
      -total_price       → buyer.wallet.gold
      +net_to_seller     → seller.wallet.gold
      +sale_cut          → system.tax_sink
    ref: listing_id
    ```
12. **Item transfer**:
    - Escrow unlock (`market.escrow_item` sil, `item_instance.is_locked=false`).
    - `item_instance.owner_character_id = buyer.id`.
    - Alıcı envanterine yerleştir (auto-slot).
    - Stackable ise buyer inventory stack güncelle.
13. **Listing kapat**: `state='sold', closed_at=now(), close_reason='sold'`.
14. **`market.sale_history` insert** (bkz. [market-history.md](market-history.md)).
15. **Outbox**:
    - WS "listing sold" → satıcıya (mail-if-offline).
    - WS "listing bought" → alıcıya.
    - Search index invalidation.
16. **COMMIT.**

## 4. Price Race Handling

Aynı listing için concurrent iki buy:
- `SELECT FOR UPDATE` ile ilk alan lock alır.
- İkinci `state='sold'` görür → `E_LISTING_UNAVAILABLE`.
- Idempotency key ile duplicate buy safe (retry same key → cached success/error).

## 5. Envanter Yetersiz Durumu

Alıcı envanteri doluysa buy başarısız olur (rollback). Öneri: **auto-bank fallback** (Phase 02) — bank slot varsa oraya git; v1'de yok.

## 6. Offline Satıcı

Satıcı offline ise:
- Gold satıcının wallet'ına ledger yazılır (offline hesaba yansıyor).
- Notification queue → satıcı online olduğunda WS + inbox message.
- Mail sistemi (Phase 02) yerine wallet direkt kredi; itemi zaten escrow'daydı.

## 7. Hata Kodları

| Kod | Açıklama |
|---|---|
| `E_LISTING_UNAVAILABLE` | Sold / cancelled / expired |
| `E_PRICE_MISMATCH` | Listing güncellendi (v1'de update yok ama race) |
| `E_LISTING_EXPIRED` | Süre dolmuş, worker henüz temizlemedi |
| `E_SAME_ACCOUNT_BLOCK` | Kendi account'un satıcısı |
| `E_SELF_BUY_BLOCK` | Kendi karakter listing'i |
| `E_INSUFFICIENT_BALANCE` | Gold yetmiyor |
| `E_INVENTORY_FULL` | Buyer inventory yetersiz |
| `E_ITEM_NOT_TRADABLE` | Bind state değişmiş (paranoia check) |

## 8. Metrikler
- `market_buy_total{item_template}`
- `market_buy_latency_ms`
- `market_price_mismatch_total`
- `market_same_account_block_total`
- `market_inventory_full_reject_total`

## 9. Cross-Reference
- [listing.md](listing.md), [sell.md](sell.md), [tax.md](tax.md), [market-history.md](market-history.md)
- [economy/transaction.md](../economy/transaction.md)
- [inventory/item-transfer.md](../inventory/item-transfer.md)
- [systems/notification.md](../systems/notification.md)

## 10. Bekleyen Kararlar
- Partial buy (stackable listing'in bir kısmı) Phase 02.
- Auto-bank fallback (envanter doluysa bank'a git) Phase 02.
- Buy limit (spam bot koruma) — v1 rate limit 20/dk yeterli.
