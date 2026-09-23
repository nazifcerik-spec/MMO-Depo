# Marketplace Tax (Listing Fee + Sale Cut)

> Ekonomi tax altyapısı için [economy/tax.md](../economy/tax.md) ve [economy/transaction.md](../economy/transaction.md).

## 1. Amaç
Marketplace iki tür tax uygular:
- **Listing fee** — post sırasında satıcıdan alınır (non-refundable expire; kısmi cancel iade).
- **Sale cut** — buy sırasında alıcı payment'ından kesilir; satıcıya net gider.

Her ikisi `system.tax_sink` hesabına gider.

## 2. Listing Fee

**Formül (v1):**
```
listing_fee = round(unit_price × stack_quantity × base_rate × duration_multiplier)
base_rate = 0.05   -- %5
```

**Duration multiplier:**

| Duration | Multiplier |
|---|---|
| 12h | 1.0× |
| 24h | 1.5× |
| 48h | 2.5× |
| 168h (7d) | 5.0× |

Minimum listing fee: **10_c** (0.1 silver) — spam koruma.

**Örnek:** unit_price = 10.000_c, qty=1, duration=24h:
`fee = round(10000 × 1 × 0.05 × 1.5) = 750_c`

Fee ödendiği anda satıcının wallet'ından düşer; listing sold olsun/olmasın (expire → yanmış; cancel → §cancel.md kademeli iade).

## 3. Sale Cut

**Formül (v1):**
```
sale_cut = round(total_price × sale_cut_rate)
sale_cut_rate = 0.05   -- %5
net_to_seller = total_price - sale_cut
```

Buyer `total_price` öder; seller `net_to_seller` alır; delta `system.tax_sink`'e.

**Örnek:** total_price = 10.000_c:
- buyer -10.000
- seller +9.500
- tax_sink +500

## 4. Kombinasyon Örneği

Satıcı bir itemi 10.000_c (100g) 24h listeliyor, satılıyor:
- Listing fee: 750_c (satıştan bağımsız).
- Sale cut: 500_c (satıştan sonra).
- **Toplam sink**: 1.250_c.
- Satıcının net gold gelirinden gideri: 750 + 500 = 1.250_c → net kazanç 8.750_c.
- Alıcı öder: 10.000_c.

## 5. Ledger Şeması

### Post
```
kind='market_listing_fee'
-fee → seller.wallet.gold
+fee → system.tax_sink
ref_kind='listing' ref_id=listing_id
```

### Buy
```
kind='market_buy'
-total_price   → buyer.wallet.gold
+net_to_seller → seller.wallet.gold
+sale_cut      → system.tax_sink
ref_kind='listing' ref_id=listing_id
```

### Cancel iade (v1 kademeli)
```
kind='market_listing_fee_refund'
+refund → seller.wallet.gold
-refund → system.tax_sink
ref_kind='listing' ref_id=listing_id
```

## 6. Konfig

```yaml
# config/marketplace/tax.yml
listing_fee:
  base_rate: 0.05
  min_copper: 10
  duration_multiplier:
    12: 1.0
    24: 1.5
    48: 2.5
    168: 5.0
sale_cut:
  rate: 0.05
cancel_refund:
  under_1h_ratio: 0.50
  under_6h_ratio: 0.25
  over_6h_ratio: 0.00
```

Hot reload — bkz. [architecture/config.md](../architecture/config.md).

## 7. Reputation / Guild İndirimi (Phase 02)

V1'de indirim yok. Phase 02'de faction reputation ile listing fee %10'a kadar indirim; guild perk ile sale cut %5→%4 (bkz. [guild/research.md](../guild/research.md)).

## 8. Metrikler
- `market_listing_fee_collected_copper_total`
- `market_sale_cut_collected_copper_total`
- `market_fee_refund_copper_total`
- `market_tax_total_by_kind{kind}`

## 9. Cross-Reference
- [economy/tax.md](../economy/tax.md), [economy/transaction.md](../economy/transaction.md)
- [listing.md](listing.md), [sell.md](sell.md), [buy.md](buy.md), [cancel.md](cancel.md)

## 10. Bekleyen Kararlar
- Rate final ayar Phase 02 balance pass (%5 fee + %5 cut yüksek mi?).
- Auction listing fee ayrı formül Phase 02 (bkz. [auction.md](auction.md)).
- Reputation/guild indirimi Phase 02.
- Legendary/high-value tier progressive tax Phase 03.
