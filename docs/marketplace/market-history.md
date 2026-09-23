# Market History (Historical Price Data)

## 1. Amaç
- Suggested price hesaplama kaynağı.
- Oyuncu için "son 7 gün fiyat trendi" grafik.
- Ekonomi analitiği (inflation, item scarcity, meta shift).
- Anti-cheat sinyali (unusual price pattern).

## 2. Model

```sql
CREATE TABLE market.sale_history (
  id BIGSERIAL PRIMARY KEY,
  item_template_id UUID NOT NULL,
  listing_id UUID NULL,           -- fixed-price sale
  auction_id UUID NULL,           -- auction sale
  stack_quantity INT NOT NULL,
  unit_price_copper BIGINT NOT NULL,
  total_price_copper BIGINT NOT NULL,
  sale_cut_copper BIGINT NOT NULL,
  seller_character_id UUID NOT NULL,
  buyer_character_id UUID NOT NULL,
  sold_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  snapshot JSONB NOT NULL          -- filter/analiz için item snapshot
);

CREATE INDEX idx_sale_hist_item_time
  ON market.sale_history (item_template_id, sold_at DESC);

CREATE INDEX idx_sale_hist_time
  ON market.sale_history (sold_at DESC);
```

Retention:
- **Hot** 90 gün (aktif query).
- **Warm** 12 ay (analytics query, daha yavaş).
- **Cold** 12+ ay parquet arşiv.

## 3. Agregasyon (Rolling)

Materialized view (nightly refresh + on-demand):

```sql
CREATE MATERIALIZED VIEW market.item_price_stats_7d AS
SELECT
  item_template_id,
  COUNT(*) AS sample_size,
  MIN(unit_price_copper) AS p_min,
  percentile_cont(0.25) WITHIN GROUP (ORDER BY unit_price_copper) AS p25,
  percentile_cont(0.50) WITHIN GROUP (ORDER BY unit_price_copper) AS p50,
  percentile_cont(0.75) WITHIN GROUP (ORDER BY unit_price_copper) AS p75,
  MAX(unit_price_copper) AS p_max,
  MAX(sold_at) AS last_sold_at
FROM market.sale_history
WHERE sold_at > now() - INTERVAL '7 days'
GROUP BY item_template_id;
```

Refresh: cron `market_price_stats_refresh` her 1 saat.

## 4. Endpoint

**GET `/v1/market/history?item_template_code=&limit=100&cursor=`**
```json
{
  "item_template_code": "sword_iron",
  "recent_sales": [
    { "sold_at": "...", "unit_price_copper": 4500, "stack": 1, "rarity":"uncommon" },
    ...
  ],
  "stats_7d": {
    "sample_size": 42,
    "min": 3800,
    "p25": 4200,
    "median": 4500,
    "p75": 4800,
    "max": 5500
  }
}
```

**GET `/v1/market/suggest?item_template_code=...`** — kısa suggest (bkz. [price.md](price.md)).

## 5. Chart Görselleştirme (Client)

- Son 7 gün fiyat scatter/line.
- Median + P25/P75 band.
- Volume (satış adedi/gün) bar.

Client tarafından chart; data server'dan aggregated (raw sale list hafif).

## 6. Analytics Kullanımı

- Inflation tracker: `median_price_over_time` per item class.
- Meta shift: sudden volume spike / price drop.
- RMT signal: single buyer'a çoklu high-price alım (anti-fraud model input).

Analytics erişimi ops-only; player'a export edilmez.

## 7. Privacy

- Public history: `sold_at, price, rarity, snapshot summary`.
- Gizli: `seller_id, buyer_id` (public API'da açılmaz; oyuncu spy önlemi).
- Kendi history: `GET /v1/character/{id}/market/sales` (auth zorunlu).

## 8. Metrikler
- `market_history_query_total`
- `market_history_query_latency_ms`
- `market_price_stats_refresh_duration_ms`
- `market_price_median_gauge{item_template}`

## 9. Cross-Reference
- [price.md](price.md), [buy.md](buy.md), [auction.md](auction.md)
- [database/economy.md](../database/economy.md)

## 10. Bekleyen Kararlar
- Materialized view refresh sıklığı (1h v1; Phase 02 15dk).
- Public history granularity (buyer_id hep gizli — Phase 02 review).
- Cold arşiv format (parquet — DD-open economy).
- Chart data endpoint pre-aggregated bucket (10dk/1h/1d) Phase 02.
