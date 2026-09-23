# Filter (Search Dimensions)

## 1. Amaç
Search endpoint'inin desteklediği filtre kümelerinin tam listesi + validasyon.

## 2. Filter Boyutları (v1)

| Filter | Tip | Kaynak |
|---|---|---|
| `item_template_code` | string | content.item_template.code |
| `item_class` | enum | weapon / armor / consumable / material / misc / recipe |
| `item_subclass` | enum | class'a bağlı (sword/axe/bow, cloth/leather/mail/plate, potion/food/scroll, ore/wood/herb...) |
| `rarity` | enum comma-list | common/uncommon/rare/epic/legendary |
| `min_level` / `max_level` | int | item.level_req |
| `min_price` / `max_price` | bigint (copper) | listing.unit_price_copper |
| `socket_count_min` | int | snapshot.sockets |
| `gem_present` | bool | snapshot.gems.length > 0 |
| `enchant_present` | bool | snapshot.enchant IS NOT NULL |
| `stat_priority` | enum comma-list | snapshot.stats key match (str/agi/int/vit/…) |
| `affix_prefix` / `affix_suffix` | string | snapshot.affix_prefix/suffix |
| `set_code` | string | snapshot.set_code |
| `bind_kind` | enum | boe / boa (bop zaten hariç) |
| `crafted_by_character_id` | uuid | snapshot.crafter_id (mastery signature) |

## 3. Validation

- Enum values whitelist (aksi → `E_FILTER_INVALID`).
- Range kontrol (`min_level ≤ max_level` gibi).
- Price copper > 0.
- Filter count max 12 (over-complex query koruma).

## 4. Filter Kombinasyonları

Hepsi AND (or opsiyonu yok v1'de).
Comma-list olan filter'lar (rarity, stat_priority) internal olarak OR:
```
rarity=rare,epic  →  rarity IN ('rare','epic')
```

## 5. Snapshot Alanları

`market.listing_snapshot.item_snapshot` JSONB şekli:
```json
{
  "item_template_code": "sword_iron",
  "level_req": 10,
  "class": "weapon",
  "subclass": "sword",
  "rarity": "uncommon",
  "sockets": 1,
  "gems": [],
  "enchant": null,
  "affix_prefix": "sharp",
  "affix_suffix": "of_the_bear",
  "stats": [{"key":"str","value":3},{"key":"agi","value":1}],
  "set_code": null,
  "bind_kind": "boe",
  "crafter_id": null,
  "quality": "normal"
}
```

Snapshot listing post'ta oluşur; life-of-listing sabittir.

## 6. Index Stratejisi (Postgres v1)

```sql
-- Aktif listing partial index
CREATE INDEX idx_market_listing_active_price
  ON market.listing (item_template_id, unit_price_copper)
  WHERE state = 'active';

-- Snapshot JSON GIN (facet + attribute filter)
CREATE INDEX idx_market_listing_snapshot_gin
  ON market.listing_snapshot USING GIN (item_snapshot jsonb_path_ops);

-- Rarity + level partial (yaygın filter)
CREATE INDEX idx_market_listing_active_rarity_level
  ON market.listing_snapshot ((item_snapshot->>'rarity'), ((item_snapshot->>'level_req')::int));
```

Phase 02'de Elastic index'e mirror + facet.

## 7. Cross-Reference
- [search.md](search.md), [listing.md](listing.md)
- [item-engine/README.md](../item-engine/README.md) — item attribute kaynağı

## 8. Bekleyen Kararlar
- Facet aggregation Phase 02.
- Range slider UI için histogram endpoint Phase 02.
- Saved search (favori filter) Phase 02.
