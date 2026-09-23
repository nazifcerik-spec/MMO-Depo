# Distance

## 1. Amaç
İki lokasyon arasındaki mesafenin belirlenmesi. Positional coord yok (v1 abstract); pre-computed route table + graph traversal.

## 2. Model

### 2.1 Route Table (Static)

`world.travel_route` (bkz. [database/world.md §9](../database/world.md)):
```sql
world.travel_route
  id, code, origin_location_id, destination_location_id,
  method, duration_ms, cost_gold_copper, schedule_json
```

Bu tablo **fixed method (wagon/ship)** routes için; walking/mount/road için dinamik hesap.

### 2.2 Distance Graph (v1 Yaklaşım)

Location-to-location adjacency graph:
```
CREATE TABLE world.location_adjacency (
  from_location_id UUID NOT NULL,
  to_location_id UUID NOT NULL,
  distance_units INT NOT NULL,          -- abstract unit
  terrain_kind VARCHAR(24) NOT NULL,    -- grass, forest, mountain, water, road, cave, desert
  road_bonus_pct SMALLINT NOT NULL DEFAULT 0,
  bidirectional BOOLEAN NOT NULL DEFAULT true,
  PRIMARY KEY (from_location_id, to_location_id)
);
```

Adjacency edge = "birbirine komşu" iki location.

## 3. Pathfinding

Karakter travel isterse:
```
origin: loc-ashfield
destination: loc-emerald-coast

Server pathfind:
  Dijkstra veya A* over world.location_adjacency
  weight = distance_units × terrain_multiplier(edge.terrain_kind)
  path: [ashfield → wheatcross → highroad → emerald_coast]
  total distance = Σ edges.distance
```

Cache: her (origin, destination) çifti kısa TTL Redis'te — sık kullanılan rotalar hızlı.

## 4. Distance Unit

Abstract; "unit" = normalize edilmiş coğrafi ölçek.

Örnek content:
```yaml
# content/world/adjacencies/ashfield_to_wheatcross.yaml
from_location: ashfield
to_location: wheatcross
distance_units: 100
terrain_kind: grass_road
road_bonus_pct: 15
bidirectional: true
```

Content pack tarafından toplu YAML.

## 5. Base Walking Speed

Config:
```yaml
base_walking_speed_units_per_sec: 1.0
```

Yani karakter walking ile 100 unit'i **100 sn** kat eder.

Speed multipliers (mount, race, etc.) bu baz üstünden çarpımsal (bkz. [base-time.md](base-time.md)).

## 6. Terrain-Aware Distance

Terrain kind edge'e bağlı; distance ile birlikte multiplier:
```
effective_distance = distance_units × terrain_multiplier[edge.terrain_kind]
```

Terrain multipliers: bkz. [terrain.md](terrain.md).

## 7. Cross-Region / Cross-Continent

- Region içi walking/mount.
- Cross-region: **road** (highroad) edge veya **waygate**.
- Cross-continent: **ship** route (fixed schedule).

Fixed routes `world.travel_route` içinde ayrı; walking pathfind bunları skip eder (fixed duration).

## 8. Endpoint

**POST `/v1/characters/{id}/travel/route`**
```json
{
  "origin_location_id": "loc-ashfield",
  "destination_location_id": "loc-emerald-coast",
  "method": "mount"
}
```

Response:
```json
{
  "path": [
    { "from": "loc-ashfield", "to": "loc-wheatcross", "distance_units": 100, "terrain": "grass_road", "road_bonus_pct": 15 },
    { "from": "loc-wheatcross", "to": "loc-highroad", "distance_units": 200, "terrain": "grass_road", "road_bonus_pct": 15 },
    { "from": "loc-highroad", "to": "loc-emerald-coast", "distance_units": 200, "terrain": "forest", "road_bonus_pct": 0 }
  ],
  "total_distance_units": 500,
  "unreachable": false
}
```

Route hesap client tarafından "estimate" endpoint'in içinde de yapılır.

## 9. Pathfind Cache

Redis key: `route:<origin>:<destination>:<method>` TTL 5 dakika.

Content patch (adjacency değişse) content_version bump → cache invalidate prefix.

## 10. Unreachable Paths

- Bazı location'lar sadece belirli method ile erişilebilir (waygate-only, ship-only).
- Pathfind fail → `travel.route_not_found` error.
- Alternatif method önerisi (client hint).

## 11. Multi-Method Route

v1'de tek method; yolun tümü aynı method.
Phase 02'de multi-method:
- Örn: walk to nearest waygate → waygate teleport → walk to destination.
- Cost + duration compose.

## 12. Fail Modları
- `route.origin_unknown` (v1: character location invalid — edge case).
- `route.destination_unknown`.
- `route.unreachable`.
- `route.method_incompatible` (walking'a "ship" istemek).

## 13. Metrics
- `route_computed_total`.
- `route_cache_hit_ratio` gauge.
- `route_unreachable_total`.
- `route_avg_distance_units` gauge.

## 14. Cross-Reference
- [Travel](travel.md), [Base Time](base-time.md), [Terrain](terrain.md)
- [world/location.md](../world/location.md)
- [database/world.md §9](../database/world.md)

## 15. Bekleyen Kararlar
- Pathfind algorithm final (Dijkstra vs A*).
- Multi-method route Phase 02.
- Adjacency graph precompute vs live query.
- Rich terrain features (elevation, river crossings) Phase 03.
