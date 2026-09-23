# Terrain

## 1. Amaç
Edge (adjacency) terrain türünün travel süresini etkilemesi. Grass hızlı, mountain yavaş.

## 2. Terrain Kinds

Whitelist:
```
grass         (open plains)
grass_road    (paved road on grass)
forest
mountain
hills
desert
tundra        (snowy plains)
water         (river / lake / sea)
coast
swamp
cave          (underground)
lava_zone     (Phase 02 rare)
magical_terrain  (Phase 02 arcane)
```

## 3. Terrain Multiplier (Süre Penalty)

Multiplier > 1.0 → süre artar (yavaş):
```yaml
terrain_multipliers:
  grass: 1.0
  grass_road: 0.90     # road faster than grass
  forest: 1.3
  mountain: 1.8
  hills: 1.2
  desert: 1.4
  tundra: 1.3
  water: 3.0           # walking through water very slow (or forbidden)
  coast: 1.0
  swamp: 1.6
  cave: 1.2
  lava_zone: 2.5
  magical_terrain: 1.1
```

Aksi taraf multiplier `<1.0` (road) speed boost.

## 4. Mount-Terrain Interaction

Bazı terrain mount tipe göre farklı:
- Aquatic mount + water = 0.5 (fast).
- Ground mount + water = 5.0 (impossible; auto-swim, çok yavaş) veya `travel.route_forbidden`.
- Winged mount (Phase 02) + mountain = 1.0 (bypasses).

Compute:
```go
effective_terrain_mult := terrainMultipliers[edge.TerrainKind]

if mount.Type == "aquatic" && edge.TerrainKind == "water" {
    effective_terrain_mult = 0.5   // aquatic bonus
} else if mount.Type == "ground" && edge.TerrainKind == "water" {
    effective_terrain_mult = 5.0   // penalty
}
```

## 5. Race-Terrain Interaction

Bazı race terrain'e bias:
- Elf + forest = 1.1 (elven_grace scope forest → -20% penalty).
- Dwarf + mountain/cave = 1.4 instead of 1.8 (dwarven passive).
- Merfolk + water = 1.5 without mount; 0.3 with aquatic mount.
- Orc + hills/plains = 1.0 (no penalty).

Config-driven race_terrain_resist table.

## 6. Road Bonus

Edge `road_bonus_pct SMALLINT`:
- Grass road: +15% speed.
- Highroad (mid): +20% speed.
- Ancient road (rare): +25% speed.

Compute:
```
road_multiplier = 1.0 + edge.road_bonus_pct / 100.0
```

Road multiplier `>1.0` → boost.

## 7. Segment-Based

Route her edge için terrain + road ayrı hesap; total sum:
```go
for _, edge := range route {
    terrainMult := terrainMultipliers[edge.TerrainKind]
    // Mount-terrain adjust
    // Race-terrain adjust
    roadMult := 1.0 + edge.RoadBonusPct/100.0

    segTime := (edge.DistanceUnits * terrainMult * 1000) /
               (baseSpeed * mountMult * raceMult * equipmentMult * roadMult)
    totalMs += segTime
}
```

## 8. Content Pack (Adjacency)

Bkz. [distance.md §2.2](distance.md):
```yaml
from_location: ashfield
to_location: wheatcross
distance_units: 100
terrain_kind: grass_road
road_bonus_pct: 15
```

Content validator: `terrain_kind` whitelist içinde olmalı.

## 9. Terrain Overlay (Phase 02)

- Weather + terrain kombine (fresh snow + forest = deeper penalty).
- Time of day (night + swamp = extra penalty; Phase 02).

## 10. Cross-Reference
- [Travel](travel.md), [Base Time](base-time.md), [Mount](mount.md), [Race](race.md), [Distance](distance.md)
- [world/region.md](../world/region.md) — region biome hint
- [world/location.md](../world/location.md) — adjacency edges

## 11. Bekleyen Kararlar
- Terrain multiplier final playtest.
- Cave/lava/magical terrain content (Phase 02).
- Elevation system (climbing penalty) Phase 03.
- Terrain-triggered event (ambush chance in forest) Phase 02.
