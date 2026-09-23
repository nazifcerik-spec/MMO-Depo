# Market (Marketplace)

## 1. Amaç
Marketplace browse, search, filter, buy, list, history chart. Sekmeler: Browse / My Listings / History.

## 2. Layout — Browse

```
┌────────────────────────────┐
│ ← Market                   │
│ [Browse][My List][History] │
├────────────────────────────┤
│ 🔍 [ara]        Filtre ⇧  │
│ Rarity: [Un][Ra][Ep] Level:│
│ [1-20][21-40][41+]   Sort ▼│
├────────────────────────────┤
│ ┌───────────────────────┐  │
│ │ Iron Sword Uncommon L15│  │
│ │ 450c · 24dk kaldı     │  │
│ │ Grim satıyor          │  │
│ │ [Detay]               │  │
│ └───────────────────────┘  │
│ ...                        │
└────────────────────────────┘
```

## 3. Listing Detail Sheet

```
┌────────────────────────────┐
│ [icon] Iron Sword          │
│ Uncommon · L15 · BoE       │
├────────────────────────────┤
│ ⚔ Dmg 12-18  ★ Crit +2%   │
│ (+3 STR, Sharp / of Bear)  │
├────────────────────────────┤
│ Fiyat: 450c                │
│ Adet: 1                    │
│ Suggested (7d median): 480c │
├────────────────────────────┤
│ [Şimdi Al]  [Fiyat Grafik] │
└────────────────────────────┘
```

## 4. Buy Confirm

Modal: "450c ödeyecek, cüzdanında 12,340c var. Onayla?" → `POST /v1/market/listing/{id}/buy` with `expected_unit_price_copper`.

## 5. My Listings Tab

- Aktif listings list.
- Progress bar (time remaining).
- Sold / expired / cancelled history.
- **Yeni Listeleme** CTA.

## 6. List Flow

1. Envanterden item seç.
2. Fiyat gir (suggested price shown; min/max clamp validation client-side).
3. Süre seç (12h/24h/48h/168h).
4. Preview: listing fee + estimated net if sold.
5. Confirm → `POST /v1/market/listing`.

## 7. History Chart

- Recent sales chart (last 7 days).
- Price median band.
- Volume bar.
- Data source: `GET /v1/market/history?item_template_code=`.

## 8. State

- Search provider — cached 30s query hash.
- My listings — `GET /v1/character/{id}/market/listings`.

## 9. API

- `GET /v1/market/search?…`
- `GET /v1/market/listing/{id}`
- `POST /v1/market/listing`
- `POST /v1/market/listing/{id}/buy`
- `DELETE /v1/market/listing/{id}`
- `GET /v1/market/history?item_template_code=`
- `GET /v1/market/suggest?item_template_code=`

## 10. Metrikler
- `ui_market_search_total`
- `ui_market_buy_total`
- `ui_market_list_total`

## 11. Cross-Reference
- [marketplace](../marketplace/README.md)
