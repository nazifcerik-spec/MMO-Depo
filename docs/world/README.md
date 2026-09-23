# World / Map (Faz 13)

Dünyanın hiyerarşik yapısı ve karakter konumlanma sistemi. Aktivite, seyahat, keşif, ekonomi hepsi bu iskelete oturur.

Yüksek seviye vizyon: [game-design/world.md](../game-design/world.md), [game-design/map.md](../game-design/map.md), [game-design/cities.md](../game-design/cities.md), [game-design/travel.md](../game-design/travel.md). Bu klasör Faz 13 implementation-hazır detay.

## 1. Hiyerarşi

```
World
 └── Continent
     └── Region
         └── Territory
             ├── City
             ├── Village
             ├── Outpost
             └── Location (POI anchor)
                 ├── Resource Node       (mining/fishing/foraging)
                 ├── Monster Area        (hunt spot)
                 ├── Dungeon Area        (dungeon entrance)
                 └── Waygate / Travel POI
```

Her seviye kendi entity'si; parent-child FK ile bağlı.

## 2. Doküman Haritası

| Konu | Doküman |
|---|---|
| Genel bakış (bu doküman) | [README.md](README.md) |
| [World](world.md) | Top-level singleton container |
| [Continent](continent.md) | Kıta seviyesi + unlock |
| [Region](region.md) | Bölge + biyom + kimlik |
| [Territory](territory.md) | Politik alt bölünüm + faction + territory war hazır |
| [City](city.md) | Başkent + servis + spawn point |
| [Village](village.md) | Küçük yerleşim + kısıtlı servis |
| [Outpost](outpost.md) | Uzmanlaşmış küçük hub |
| [Location](location.md) | Genel POI anchor + activity hook |
| [Resource Node](resource-node.md) | Mining/fishing/foraging node instance |
| [Monster Area](monster-area.md) | Hunt spot + mob spawn |
| [Dungeon Area](dungeon-area.md) | Dungeon entrance + instance yönetimi |

## 3. Content-Driven

Tüm world içeriği YAML content pack:
```
content/world/
  continents/*.yaml
  regions/*.yaml
  territories/*.yaml
  settlements/*.yaml       (city + village + outpost — settlement_kind ile)
  locations/*.yaml
  resource_nodes/*.yaml    (template; instance runtime)
  monster_areas/*.yaml
  dungeon_areas/*.yaml
  waygates/*.yaml
  travel_routes/*.yaml
```

Loader validation + content_version bump ile hot reload. Karakter mevcut lokasyonu snapshot'ta değil (dinamik değişebilir); ancak content patch mid-session temkinli:
- Deprecated location → karakter safe city'ye teleport + notification.

## 4. Database Şeması

Bkz. [database/world.md](../database/world.md).

Ana tablolar (özet):
- `world.continent`, `world.region`, `world.territory`, `world.settlement` (kind ayrımı ile city/village/outpost tek tablo), `world.location`, `world.poi`, `world.resource_node` (instance), `world.waygate`, `world.travel_route`, `world.character_waypoint`, `world.character_discovery`.

Ekleme (Phase 13):
- Monster area ve dungeon area için ayrı tablolar veya `world.location.tags[]` üzerinden filter (öneri: tags + `activity_hooks_json` yeter; ayrı tablo gerekmez). Detay ilgili dokümanlarda.

## 5. Character Location

`character.character.current_location_id UUID FK NULL → world.location(id)`.

Aktivite başlatırken: karakter belirli location'da olmalı (spot-based aktiviteler).

Travel aktivitesi → location değişimi.

## 6. Endpoint Genel

```
GET /v1/content/world                              # continent + region özet
GET /v1/content/world/continents                   # tüm kıta
GET /v1/content/world/continents/{code}            # continent detail
GET /v1/content/world/regions/{code}               # region + territories + settlements + locations
GET /v1/content/world/territories/{code}           # territory + settlements
GET /v1/content/world/settlements/{code}           # city/village/outpost detay + services
GET /v1/content/world/locations/{code}             # location + POI + activity hooks
GET /v1/content/world/locations/{code}/nodes        # resource_node instance state
GET /v1/content/world/locations/{code}/monsters    # mob spawn info
GET /v1/content/world/locations/{code}/dungeon     # dungeon entrance info

POST /v1/characters/{id}/travel                    # travel start (bkz. activity-engine/travel.md)
GET  /v1/characters/{id}/discovery                 # explored POI listesi
POST /v1/characters/{id}/waypoints                 # kişisel waypoint
```

Public content endpoint'ler cache-friendly (content_version ETag).

## 7. Discovery (Fog of War) — Phase 02

v1'de tüm POI görünür (`visited` flag hariç).

Phase 02:
- `hidden_by_default = true` POI'ler oyuncu yakınlaşana kadar map'te görünmez.
- Discovery event → `world.character_discovery` insert + discovery XP grant.

## 8. Runtime State

**Static content** (change with content_version):
- Continent, region, territory, settlement, location definitions.

**Dynamic runtime state**:
- `world.resource_node` instance (respawn timer).
- Dungeon instance (activity.instance).
- World boss spawn (activity based, [game-design/world-bosses.md](../game-design/world-bosses.md)).
- Territory ownership (Phase 02 territory war).

## 9. Server-Authoritative Location

Karakter konum:
- `current_location_id` server tutar.
- Client display; server otoritedir.
- Location değişimi sadece travel activity üzerinden (veya spawn/teleport/death).

Location-gated aksiyonlar:
- Bank açma (city bank service).
- Vendor.
- Trainer.
- Aktivite başlatma (hunt spot, resource node).

## 10. Cross-Region Party

Farklı kıta/region karakterler party olabilir (chat, guild) ancak instance aktivitesi (dungeon, adventure) coğrafi olarak tek anchor.

## 11. Config

```yaml
# config.world.yaml
default_starter_city_by_race:
  human: ashfield
  dwarf: hearthforge
  elf: silverbough
  orc: stormcrag
  merfolk: tidehaven
  arcanborn: azuremind
travel:
  base_walking_speed_units_per_sec: 1.0
  ambush_enabled: false           # Phase 02
discovery:
  fog_of_war_enabled: false       # Phase 02
  discovery_xp_per_poi: 100
```

## 12. Cross-Reference

- [game-design/world.md](../game-design/world.md), [game-design/map.md](../game-design/map.md), [game-design/cities.md](../game-design/cities.md), [game-design/travel.md](../game-design/travel.md)
- [database/world.md](../database/world.md)
- [activity-engine/travel.md](../activity-engine/travel.md), [activity-engine/mining.md](../activity-engine/mining.md), [activity-engine/hunting.md](../activity-engine/hunting.md), [activity-engine/dungeon.md](../activity-engine/dungeon.md)
- [character/character.md](../character/character.md) — home_city + current_location

## 13. Ortak Bekleyen Kararlar
- Fog of war Phase 02 aktivasyon.
- Weather/season sistem Phase 02.
- Territory war (Phase 02+) → territory ownership dynamic.
- Real coord system vs abstract (v1 abstract; Phase 02 opsiyonel positional).
- Cross-region market (Phase 03).
