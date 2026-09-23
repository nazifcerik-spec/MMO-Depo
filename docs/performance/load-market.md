# Load Test: Market

## 1. Amaç
Market search + list + buy + history endpoint'lerin yüksek arama trafiğinde SLO korunması. Cache hit ratio + serialization retry.

## 2. Senaryolar

### 2.1 Steady Market Browsing
- Yük: T3'te 3000 active × 3 search/saat = **2.5 search/sn**.
- Filter dağılımı: rarity(30%) + level(25%) + item_class(30%) + price_range(15%).
- Cache hit target > %70 (30sn TTL).
- Beklenti: search p95 < 200ms.

### 2.2 Buy Rush (Rare Legendary Drop)
- Yük: legendary listing → 60 sn içinde 500 buy attempt.
- Beklenti: 1 success + 499 `E_LISTING_UNAVAILABLE`; SERIALIZABLE retry < 3.
- Latency spike controlled; no cascade fail.

### 2.3 List Storm (Content Patch)
- Yeni content release → oyuncular yeni item listeler.
- Yük: 100 list/sn 10 dk.
- Beklenti: list p95 < 400ms; escrow lock atomicity korunur.

### 2.4 Long-Running Search Load
- 4 saat soak; 2.5 search/sn + %20 filter varyasyon.
- Beklenti: cache hit stabil > %70; DB CPU stabil.

## 3. Kritik Path

### Search
1. HTTP receive.
2. Rate limit + query parse.
3. Normalized cache key hash.
4. Redis GET cache → hit ise return.
5. Miss → Postgres JSONB GIN query.
6. Response format.
7. Cache SET (TTL 30sn).

**Cache miss cost:** ~50-150ms DB (index optimize).
**Cache hit cost:** ~5-10ms Redis.

### Buy
1. Auth + rate limit + validate.
2. SERIALIZABLE TX begin.
3. Listing SELECT FOR UPDATE.
4. Wallet + inventory check.
5. Ledger 3 entries.
6. Item ownership transfer.
7. Listing state → sold.
8. sale_history insert.
9. Outbox WS "sold" + cache invalidate.
10. COMMIT.

**Buy cost:** ~80-200ms (SERIALIZABLE serialization; retry on conflict).

## 4. Metrics

- `market_search_p95_ms{cache_hit}`.
- `market_search_cache_hit_ratio`.
- `market_buy_p95_ms`.
- `market_buy_serialization_retry_total`.
- `market_listing_active_gauge`.

## 5. Postgres Index Kritikleri

```sql
CREATE INDEX idx_market_listing_active_price
  ON market.listing (item_template_id, unit_price_copper)
  WHERE state = 'active';

CREATE INDEX idx_market_listing_snapshot_gin
  ON market.listing_snapshot USING GIN (item_snapshot jsonb_path_ops);

CREATE INDEX idx_market_listing_active_rarity_level
  ON market.listing_snapshot ((item_snapshot->>'rarity'), ((item_snapshot->>'level_req')::int));
```

Query planner check EXPLAIN periodic (nightly perf test).

## 6. Cache Strategy

Redis:
- Key: `market:search:{sha256(normalized_params)}:{page_cursor}`.
- Value: JSON result set (~50 KB max page).
- TTL: 30 sn (staleness kabul).

Invalidation:
- Listing sold/cancelled/expired → publish invalidation pattern.
- Basit v1: TTL-only invalidation; Phase 02 targeted invalidation (item_template_code bazlı).

## 7. Concurrent Buy Race

Legendary listing 500 concurrent buy:
- Postgres `SELECT FOR UPDATE` serializes.
- İlk success; sonrasında `E_LISTING_UNAVAILABLE`.
- Retry storm client'ta beklenmez; error state respect edilir.
- Server-side lock hold < 50ms.

## 8. Beklenen Sonuç

T3 steady:
- Search p95: ~150ms (cache hit > 70%).
- Search cache miss p95: ~200ms.
- Buy p95: ~200ms.
- List p95: ~300ms.
- Serialization retry rate < 2%.

T4 buy rush:
- Latency spike 500ms - 1s peak; recovery 60 sn.
- 429 rate limit engage (buy endpoint 30/dk/character).

## 9. Cross-Reference
- [marketplace/search.md](../marketplace/search.md), [marketplace/buy.md](../marketplace/buy.md)
- [backend/redis.md](../backend/redis.md), [backend/database.md](../backend/database.md)

## 10. Bekleyen Kararlar
- OpenSearch/Elastic sekonder index Phase 02 (search scale).
- Cache invalidation targeted Phase 02.
- Read replica Postgres Phase 02 (search read scale).
