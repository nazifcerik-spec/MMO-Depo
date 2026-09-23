# Price (Listing Price Rules)

## 1. Amaç
Satıcının belirlediği fiyatın kabul kriterleri, suggested price, min/max clamp.

## 2. Prensipler

- **Free market** — satıcı fiyatı kendi belirler.
- **Min clamp**: item vendor sell fiyatının %110'undan aşağı olamaz (vendor arbitraj engeli).
- **Max clamp**: item vendor buy fiyatının 10.000× katından yüksek olamaz (fat-finger + inflate koruma).
- **Suggested price** UI'da gösterilir (son 7 gün median), zorlayıcı değil.

## 3. Min / Max Formülü

```
min_unit_price = max(1_c, round(item.base_price_copper × vendor.sell_multiplier × 1.10))
max_unit_price = round(item.base_price_copper × vendor.buy_multiplier × 10_000)
```

Örnek: base_price = 100_c, vendor.sell = 0.25, vendor.buy = 1.0:
- min = 27_c
- max = 1_000_000_c (100 g)

Client fiyat aralık dışı → server `E_PRICE_OUT_OF_RANGE`.

Legendary / mastery-crafted item için ayrı **premium clamp** (Phase 02); v1'de aynı kurallar.

## 4. Suggested Price

Server hesaplar (bkz. [market-history.md](market-history.md)):
```
suggested = median(sold_unit_price WHERE item_template_id=X AND sold_at > now()-7d)
```

Yeterli data yoksa (< 5 satış): `item.base_price × 1.5` fallback.

**Endpoint:** `GET /v1/market/suggest?item_template_code=...` → `{ suggested_copper, sample_size, p25, p50, p75 }`.

## 5. Price Race Guard

Buy request'i `expected_unit_price_copper` zorunlu (bkz. [buy.md](buy.md)). Server listing.unit_price ile mismatch → `E_PRICE_MISMATCH`.
V1'de listing update yok (fiyat değişmez); race sadece cache staleness kaynaklı.

## 6. Dynamic Pricing (Yok — v1)

V1'de otomatik fiyat düşürme (undercut) veya dinamik pricing **yok**. Satıcı manuel cancel + repost yapmalı (fee tekrar). Phase 03'te "auto-relist at % lower" opt.

## 7. Currency

- Tüm fiyat `gold` (copper). Diğer currency listing v1'de yok (DD-137).
- Fiyat integer copper; float yok.

## 8. Hata Kodları

| Kod | Açıklama |
|---|---|
| `E_PRICE_OUT_OF_RANGE` | Min/max clamp dışı |
| `E_PRICE_MISMATCH` | Buy zamanında listing fiyatı farklı (cache stale) |
| `E_PRICE_ZERO_OR_NEGATIVE` | ≤ 0 |

## 9. Metrikler
- `market_price_reject_out_of_range_total`
- `market_suggested_price_lookup_total`
- `market_suggested_price_no_data_total`

## 10. Cross-Reference
- [listing.md](listing.md), [sell.md](sell.md), [buy.md](buy.md), [market-history.md](market-history.md)
- [economy/vendor.md](../economy/vendor.md) — base price + multiplier

## 11. Bekleyen Kararlar
- Premium clamp legendary/mastery item Phase 02.
- Undercut / relist Phase 03.
- Auction reserve min/max Phase 02 (bkz. [auction.md](auction.md)).
