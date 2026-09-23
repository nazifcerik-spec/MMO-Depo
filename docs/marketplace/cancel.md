# Cancel (Iptal Listing)

## 1. Amaç
Satıcı aktif listing'i iptal edip itemi geri alır. Fee iade politikası kısmi.

## 2. Endpoint

**DELETE `/v1/market/listing/{id}`**
Body: `{ "idempotency_key": "..." }`

## 3. Akış (Server)

Tek TX:

1. Idempotency check.
2. Listing SELECT FOR UPDATE — `state='active'`, `seller_character_id == caller`.
3. Envanter kapasite kontrol (item geri geliyor).
4. Item escrow release:
   - `market.escrow_item` sil.
   - `item_instance.is_locked=false`.
   - `owner` zaten değişmemişti → satıcının envanterine yerleştir.
5. Stackable ise stack merge.
6. Fee iade (bkz. §4).
7. `state='cancelled', closed_at=now(), close_reason='cancelled'`.
8. Outbox: search index invalidation + notification.
9. COMMIT.

## 4. Fee İade Politikası (DD-143)

**Kural:** cancel fee iadesi süreye göre kademeli. Amaç: spam post/cancel arbitraj engeli.

| Post'tan Sonra Geçen Süre | Fee İade |
|---|---|
| < 1 saat | **%50** iade |
| 1–6 saat | %25 iade |
| ≥ 6 saat | **%0** iade |

Iade yolu:
```
kind='market_listing_fee_refund'
+refund → seller.wallet.gold
-refund → system.tax_sink   (negatif tax = iade)
```

Fee `system.tax_sink`'ten geri alındığı için sink metriği net doğru kalır.

## 5. Admin Cancel (Ops)

Admin, kural ihlali (RMT, exploit) tespitinde listing'i iptal edebilir:
- `close_reason='admin'`.
- Fee iadesi **yok** (default) veya ops kararı ile tam iade.
- audit.market_admin_action zorunlu (operator_id + reason).

## 6. Auto Expire (Worker)

Cron `market_listing_expire` her 1 dk:
```sql
SELECT id FROM market.listing
 WHERE state='active' AND expires_at <= now()
 FOR UPDATE SKIP LOCKED
 LIMIT 500;
```

Her batch için:
- Item escrow release → satıcı envanterine iade (**envanter doluysa** overflow: item Redis "market_overflow" queue → 7 gün içinde claim; sonra vendor sell auto).
  - v1'de basit: envanter doluysa `market.overflow_item` tablosuna yaz; UI'da "Marketplace overflow" claim.
- `state='expired', closed_at=now(), close_reason='expired'`.
- Fee iade **yok** (expire penaltısı).

```sql
CREATE TABLE market.overflow_item (
  id UUID PRIMARY KEY,
  character_id UUID NOT NULL,
  item_instance_id UUID NOT NULL UNIQUE,
  source VARCHAR(24) NOT NULL DEFAULT 'listing_expired',
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  claim_deadline TIMESTAMPTZ NOT NULL DEFAULT now() + INTERVAL '7 days'
);
```

Overflow claim endpoint: `POST /v1/market/overflow/{id}/claim`.

## 7. Hata Kodları

| Kod | Açıklama |
|---|---|
| `E_LISTING_NOT_FOUND` | Yok / farklı satıcı |
| `E_LISTING_STATE` | Zaten sold/expired/cancelled |
| `E_INVENTORY_FULL` | Iade sırasında kapasite yok → overflow'a yönlendirilir; hata değil |
| `E_IDEMPOTENT_REPLAY` | Cached response |

## 8. Metrikler
- `market_listing_cancelled_total{reason}`
- `market_listing_fee_refunded_copper_total`
- `market_listing_expired_total`
- `market_overflow_item_gauge`

## 9. Cross-Reference
- [listing.md](listing.md), [sell.md](sell.md), [tax.md](tax.md)
- [economy/transaction.md](../economy/transaction.md)
- [inventory/item-transfer.md](../inventory/item-transfer.md)

## 10. Bekleyen Kararlar
- Iade eşiği (1h/6h) v1; Phase 02 balance pass'te data ile ayarlanır.
- Overflow TTL 7 gün v1; Phase 02'de mail sistemi ile birleşecek (Faz 26+).
- Admin cancel iade politikası ops runbook Phase 02.
