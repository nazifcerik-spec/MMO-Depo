# Outpost

## 1. Amaç
Uzmanlaşmış küçük hub — spesifik servis (crafting station, quest giver, teleport waypoint). Genelde tek amaç.

## 2. Model

`world.settlement` `kind='outpost'`, `size='outpost'`.

## 3. Content Pack

```yaml
# content/world/settlements/ironforge_outpost.yaml
code: ironforge_outpost
territory_code: hearthhills
kind: outpost
size: outpost
display_name_key: outpost.ironforge.name
lore_key: outpost.ironforge.lore

spawn_point:
  x: 250
  y: -100
  heading: 90

services:
  bank: false
  marketplace: false
  guild_hall: false
  inn: false
  waygate: false
  vendors:
    - materials
  trainers: []
  quest_hub: true                # single quest giver typically
  home_city_available: false

specialty: mining                # mining specialist
purpose: crafting_station        # ana amaç
crafting_stations:
  - forge
  - anvil
ambient_theme: industrial

tags: [mining, forge, dwarf_settlement]
```

## 4. Outpost Türleri (v1 örnekler)

| Purpose | Örnek | Servis |
|---|---|---|
| `crafting_station` | Ironforge Outpost | Forge + anvil workbench |
| `quest_hub` | Wandering Camp | Time-limited quest giver |
| `waypoint` | Highroad Station | Waygate + rest inn |
| `resource_gate` | Deep Mine Camp | Access to premium mining nodes |
| `faction_camp` | Warden Watchtower | Faction rep grinding |
| `research_station` | Arcane Spire | Enchanter workbench + rare recipe |

## 5. Purpose Alanı

`purpose` string tag; birden fazla purpose olabilir (Phase 02) → array.

v1'de single purpose.

## 6. Crafting Station

Bazı recipe belirli station gerektirir (bkz. [database/content-templates.md §7](../database/content-templates.md), [activity-engine/crafting.md](../activity-engine/crafting.md)):
- Recipe `station_required_code` → outpost `crafting_stations` array match.
- Karakter outpost location'unda → station available.

## 7. Endpoint

City / village ile aynı endpoint (`GET /v1/content/world/settlements/{code}`); `kind` alanı `outpost`.

## 8. Small Footprint

Outpost genelde:
- 1-5 NPC.
- 0-2 vendor.
- 1 hedef aktivite (crafting station, quest, waygate).

## 9. Special Outposts

### 9.1 Faction Camp
- Faction quest chain hub.
- Bazı camp faction locked (rep threshold).

### 9.2 Waypoint Outpost
- Cross-region travel node.
- Ship dock (Phase 02).
- Waygate cluster.

### 9.3 Resource Gate
- Premium resource nodes (tier 5+).
- Level or rep gated.

## 10. Cross-Platform

Content endpoint uniform.

## 11. Metrics
- `outpost_active_visitors{outpost}` gauge.
- `outpost_purpose_usage_total{purpose}`.

## 12. Cross-Reference
- [City](city.md), [Village](village.md), [Territory](territory.md), [Waygate](location.md) (waygate location'da modellenir)
- [game-design/cities.md](../game-design/cities.md)
- [database/world.md §4](../database/world.md)

## 13. Bekleyen Kararlar
- Multi-purpose outpost Phase 02.
- Time-limited outpost (event spawn) Phase 02.
- Outpost upgrade (Phase 02+ guild investment).
