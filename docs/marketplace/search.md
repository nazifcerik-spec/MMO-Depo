# Search (Marketplace Query)

## 1. Amaç
Aktif listing'leri hızlı, sayfalanabilir, cache dostu şekilde sorgulama.

## 2. Prensipler

- **Server-side query only**; client filter yok (10k+ listing).
- **Postgres primary + Redis cache** (query key → result set) v1.
- **Phase 02**: OpenSearch/Elastic sekonder index (facet + fuzzy).
- **Rate limit**: karakter başına 20 req/dk.
- **Cache TTL** 30 sn (staleness kabul edilebilir; sold list'ten hızla çıkar).
- **Sorted by** kullanıcı seçimi; default `unit_price_copper ASC`.

## 3. Endpoint

**GET `/v1/market/search`**

Query params:
```
item_template_code=       (opsiyonel; belirli item search)
item_class=weapon|armor|consumable|material|misc
item_subclass=sword|axe|bow|... 
rarity=common,uncommon,rare,epic,legendary   (comma-list)
min_level=1
max_level=60
min_price=0
max_price=1000000
socket_count_min=0
gem_present=true
enchant_present=true
stat_priority=str,agi,int,vit   (item stat'ı içeriyorsa)
sort=price_asc|price_desc|expires_soon|recent
cursor=<opaque>
limit=50   (max 100)
```

## 4. Response

```json
{
  "results": [
    {
      "listing_id": "uuid",
      "seller_character_name": "Grim",
      "item_template_code": "sword_iron",
      "display_name_key": "item.sword_iron",
      "rarity": "uncommon",
      "level_req": 10,
      "unit_price_copper": 4500,
      "stack_quantity": 1,
      "total_price_copper": 4500,
      "expires_at": "2026-09-10T12:34:56Z",
      "snapshot_summary": {
        "stats": [{"key":"str","value":3},{"key":"agi","value":1}],
        "sockets": 1,
        "gems": [],
        "enchant": null,
        "affix_prefix": "sharp",
        "affix_suffix": "of_the_bear"
      }
    }
  ],
  "next_cursor": "opaque",
  "total_est": 1240,             // yaklaşık; kesin count yok (pahalı)
  "cache_hit": true
}
```

## 5. Cache Anahtarı

```
market:search:<hash(normalized_params)>:<page_cursor>
TTL: 30 sn
```

`normalized_params`: params key-sorted + trim + lower + default-fill sonucu.

## 6. Sıralama

| Sort | SQL |
|---|---|
| `price_asc` | `ORDER BY unit_price_copper ASC, listing_id ASC` |
| `price_desc` | `ORDER BY unit_price_copper DESC, listing_id DESC` |
| `expires_soon` | `ORDER BY expires_at ASC` |
| `recent` | `ORDER BY listed_at DESC` |

Cursor: `(sort_field_value, listing_id)` çifti base64.

## 7. Facet (Phase 02)

- **Rarity dağılımı** her rarity için count.
- **Price histogram** min/max buckets.
- **Level distribution**.

V1'de facet **yok** (Elastic Phase 02).

## 8. Anti-Abuse

- Rate limit 20 req/dk karakter, 100 req/dk account.
- Aynı query key'e (hash) 5 sn dropped → cached response.
- Search log'lanır (`analytics.market_search_event`) trend analiz için.

## 9. Metrikler
- `market_search_qps`
- `market_search_latency_ms{p50, p95, p99}`
- `market_search_cache_hit_ratio`
- `market_search_result_count_histogram`

## 10. Cross-Reference
- [filter.md](filter.md), [listing.md](listing.md), [buy.md](buy.md)
- [database/economy.md](../database/economy.md) — index
- [backend/redis.md](../backend/redis.md) — cache

## 11. Bekleyen Kararlar
- Facet Phase 02 (Elastic gerekli).
- Fuzzy name search Phase 02.
- Personalized suggestion (kullanıcı geçmiş satın alma) Phase 03.
- Global vs regional market (v1 global).
