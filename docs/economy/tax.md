# Tax (Ekonomi Vergisi)

## 1. Amaç
Sistemin gold sink'lerinden biri. Tax, ledger'da `system.tax_sink` hesabına gider — sink metriği ölçülebilir.

## 2. Tax Türleri

| Tür | Aktif v1 | Oran | Uygulanır |
|---|---|---|---|
| **Vendor spread** | ✅ (implicit) | ~%75 | vendor.sell_multiplier ile (0.25) |
| **Trade tax** | ❌ | 0 | P2P trade (Phase 02: opt %2) |
| **Auction house tax** | Faz 26 | %5 listing + %5 cut | AH işlemleri |
| **Mail postage** | Faz 26 | flat + item value % | Mail gönderim |
| **Repair tax** | Phase 02 | %10 upcharge | Repair fee üzerinden |

**Terminoloji:** "tax" gold-out; item transfer üzerinden hiç tax alınmaz (item conservation).

## 3. Vendor Spread (v1 default sink)

Vendor sell_multiplier (0.25) → oyuncunun aldığı gold, item base fiyatının %25'i. Kalan %75 sistemden **kayıt olarak** silinir. Ledger:

```
kind='vendor_sell'
entries:
  +sold_price → character.wallet.gold
  -sold_price → system.vendor_sink
```

Ek explicit tax entry gerekmez; spread implicit sink.

## 4. Trade Tax (Phase 02 opt)

V1'de tax yok. Phase 02'de yüksek değerli trade için opsiyonel formül:

```
if gold_leg > 10_000_g:
  tax = round(gold_leg × 0.02)  -- %2
```

Ledger (Phase 02):
```
kind='trade_tax'
entries:
  -tax → sender.wallet.gold      (net alınan gold - tax)
  +tax → system.tax_sink
```

Hedef: whale swap'ı yaklaşık %2 sink; küçük trade tax-free (yeni oyuncu dostu).

## 5. Auction House Tax (Faz 26 taslak)

```
listing_fee   = round(reserve_price × 0.05, min=10_c)
sale_cut      = round(final_price × 0.05)
```

Listing fee **geri iade edilmez** (iptalde bile) → spam koruma sink.

## 6. Repair Tax (Phase 02 taslak)

Repair fee (bkz. [fee.md](fee.md)) üzerine %10 tax üstüne binerek uygulanır. Vendor faction/reputation indirimi fee'ye uygulanır, tax'a uygulanmaz.

## 7. Hesaplama Kuralı

- Her tax formülü **server-authoritative**.
- **Round HALF_UP**; float kullanma, integer copper üzerinden.
- Minimum tax: `1_c` (var ise). `0_c` tax entry yazılmaz (noise).
- Tax entry her zaman **aynı TX** içinde ana operasyonla birlikte (atomic).

## 8. Konfig

```yaml
# config/economy/tax.yml
trade_tax:
  enabled: false
  threshold_copper: 100000000      # 10.000 g
  rate: 0.02
auction_house:
  enabled: false                    # Faz 26
  listing_rate: 0.05
  sale_cut_rate: 0.05
  listing_min_copper: 10
repair:
  tax_rate: 0.10
```

Hot reload — bkz. [architecture/config.md](../architecture/config.md).

## 9. Metrikler
- `economy_tax_collected_copper_total{tax_kind}`
- `economy_tax_events_total{tax_kind}`

## 10. Cross-Reference
- [transaction.md](transaction.md), [vendor.md](vendor.md), [trade.md](trade.md), [fee.md](fee.md)
- Auction House — Faz 26

## 11. Bekleyen Kararlar
- Trade tax v1 kapalı — Phase 02 balance pass sonrası aç/kapa.
- Faction-based tax indirimi Phase 02.
- Auction house tax detayı Faz 26.
- Whale tier progressive tax mı — Phase 03 tartışma.
