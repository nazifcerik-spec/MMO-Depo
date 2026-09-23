# City

## 1. Amaç
Başkent / regional büyük yerleşim. Tam servis kümesi (bank, market, trainer, vendor, guild hall, quest hub). Karakter'in "hub" noktası.

Yüksek seviye: [game-design/cities.md](../game-design/cities.md).

## 2. Model

Ortak `world.settlement` tablosu; `kind='city'` (bkz. [database/world.md §4](../database/world.md)):
```sql
world.settlement
  id, code, territory_id
  kind (city/village/outpost)
  display_name_key, size (capital/regional/village/outpost)
  spawn_point_json
  services_json
  ambient_theme
  home_faction_id
  art_ref
```

City için `kind='city'`, `size IN ('capital', 'regional')`.

## 3. Content Pack

```yaml
# content/world/settlements/ashfield.yaml
code: ashfield
territory_code: ashen_flats
kind: city
size: capital
display_name_key: city.ashfield.name
lore_key: city.ashfield.lore

spawn_point:
  x: 0
  y: 0
  heading: 0

services:
  bank: true
  marketplace: true
  guild_hall: true
  inn: true
  vendors:
    - weapon
    - armor
    - potion
    - food
    - materials
    - mount
  trainers:
    - class
    - spec
    - profession
    - recipe
  quest_hub: true
  waygate: true
  home_city_available: true

specialty: blacksmithing         # opsiyonel — hangi profession bu şehirde öne çıkar
ambient_theme: town_daytime
home_faction_code: emerald_wardens
art_ref: cities/ashfield/portrait.png
```

## 4. City Size Katmanları

| Size | Örnek | Servis Kapsamı |
|---|---|---|
| `capital` | Ashfield | Tam servis + guild hall + waygate + specialty |
| `regional` | Wheatcross-City | Tam servis - guild hall (v1); bazı specialty |

Sadece capital'de guild hall v1'de. Phase 02'de regional guild hall açılabilir.

## 5. Services Kataloğu

Services JSONB flag-based:
```json
{
  "bank": true,
  "marketplace": true,
  "guild_hall": true,
  "inn": true,
  "waygate": true,
  "vendors": ["weapon", "armor", "potion", "food", "materials", "mount"],
  "trainers": ["class", "spec", "profession", "recipe"],
  "quest_hub": true,
  "home_city_available": true,
  "faction_hub_id": "emerald_wardens"
}
```

Server-side flag check ile location-gated aksiyonlar.

## 6. Spawn Point

Karakter respawn/teleport/home city:
- `home_city_id` set edilmişse respawn buraya.
- Aksi halde default starter city.
- `spawn_point_json` x/y/heading coord (v1 abstract; Phase 02 positional).

## 7. Home City Setting

**POST `/v1/characters/{id}/home-city`**
```json
{ "city_code": "ashfield" }
```

Server:
- Validate: `services.home_city_available == true`.
- Cooldown: 7 gün (change frequency).
- Update `character.home_city_id`.
- Audit event.

Cost: gold (Phase 02'de free ilk, sonra ücretli).

## 8. City Specialty

Bazı şehir profession odaklı:
- Ashfield: blacksmithing.
- Silverbough: alchemy.
- Hearthforge: mining/jewelcrafting.
- Tidehaven: fishing.
- ...

Specialty etkisi:
- Trainer master-tier recipe access.
- Vendor exclusive material.
- Craft success bonus (+%2 within city).

Config-driven.

## 9. Endpoint

**GET `/v1/content/world/settlements/{code}`**
```json
{
  "code": "ashfield",
  "territory_code": "ashen_flats",
  "kind": "city",
  "size": "capital",
  "display_name_key": "city.ashfield.name",
  "services": { ... },
  "specialty": "blacksmithing",
  "art_ref": "..."
}
```

## 10. Safe Zone

- City = safe zone (v1 tüm settlement safe).
- Combat yasak; PvP (Phase 02) opt-in.
- Consumable use OK.

## 11. Player Instance / Layering

Peak time yoğunluğu için soft-sharding (opsiyonel):
- Aynı city'yi ziyaret eden N oyuncu → layer'lara ayrılır.
- Chat, market, guild hall paylaşımlı; visual crowd density kontrol.
- v1'de single layer; Phase 02 optional.

## 12. Cross-Platform

City list content endpoint her platformda aynı; client cache-friendly.

## 13. Metrics
- `city_active_visitors{city}` gauge.
- `city_home_city_set_total{city}`.
- `city_service_usage_total{city,service}`.

## 14. Audit
- `content.city_created/updated`.
- `character.home_city_changed`.

## 15. Cross-Reference
- [Territory](territory.md), [Village](village.md), [Outpost](outpost.md)
- [game-design/cities.md](../game-design/cities.md)
- [database/world.md §4](../database/world.md)
- [inventory/bank.md](../inventory/bank.md) — bank service dependency
- [character/character.md](../character/character.md) — home_city_id

## 16. Bekleyen Kararlar
- Layer/sharding threshold v1 vs Phase 02.
- Home city change cost final (v1 free first + 100g sonra, cooldown 7 gün).
- Guild hall regional expansion Phase 02.
- City events (parade, market day) Phase 02.
