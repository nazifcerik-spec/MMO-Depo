# Terrain (Mount Terrain Compatibility)

## 1. Amaç
Mount her terrain'de aynı hızda değil. Compatibility multiplier — species'a göre penalty/bonus.

Foundation: [travel/terrain.md](../travel/terrain.md).

## 2. Terrain Multiplier per Mount Type

Content field `terrain_compatibility_json`:
```yaml
terrain_compatibility:
  grass: 1.0
  grass_road: 0.9        # road hızlı
  forest: 1.2            # normal ground penalty
  hills: 1.1
  mountain: 1.6          # ground mount slow
  desert: 1.3
  tundra: 1.3
  water: 5.0             # walking through water very slow
  coast: 1.0
  swamp: 1.6
  cave: 1.4
```

Multiplier > 1.0 → süre artar (penalty).
Multiplier < 1.0 → süre azalır (bonus).

## 3. Mount Type Bias

| Type | Terrain Strengths | Terrain Weaknesses |
|---|---|---|
| **Ground horse** | grass, road, hills | mountain, water |
| **Ground wolf** | forest, hills | mountain, water |
| **Ground bear** | mountain, hills, tundra | swamp, water |
| **Aquatic** | water, coast | grass (slow) |
| **Winged** (Phase 02) | all + air | none |
| **Cosmetic** | universal same as ground horse | – |

## 4. Aquatic Mount Water Bonus

Aquatic mount + water terrain:
```
effective_multiplier = 0.5 (very fast)
```

Non-aquatic mount + water:
```
effective_multiplier = 5.0-10.0 (very slow or forbidden)
```

Route compute may forbid path (bkz. [travel/distance.md](../travel/distance.md)).

## 5. Ground Mount + Water Fallback

Ground mount encounters water:
- Auto-swim mode: character dismount + swimming penalty.
- Alternative: route reroute around water.

## 6. Winged Mount (Phase 02)

Winged mount + mountain/water:
```
effective_multiplier = 1.0 (bypass)
```

Cross-terrain freedom — Phase 02 dedicated design.

## 7. Compute Segment-Based

Bkz. [travel/base-time.md](../travel/base-time.md):
```go
for _, edge := range route {
    baseMultiplier := terrainMultipliers[edge.TerrainKind]
    // Override with mount compatibility
    if mount != nil && mount.TerrainCompatibility[edge.TerrainKind] != nil {
        mountMultiplier := mount.TerrainCompatibility[edge.TerrainKind]
    } else {
        mountMultiplier := baseMultiplier
    }
    segmentTime := edge.DistanceUnits × mountMultiplier × 1000 / ...
    totalMs += segmentTime
}
```

## 8. Race + Mount Terrain Interaction

- Merfolk race + aquatic mount + water = super fast (%50 stack).
- Dwarf race + ground mount + mountain = %20 mount penalty reduced.

## 9. Endpoint

Estimate response includes terrain breakdown (bkz. [travel/travel.md §3.1](../travel/travel.md)):
```json
{
  "segments": [
    { "terrain": "grass_road", "distance_units": 300, "duration_ms": 90000 },
    { "terrain": "forest", "distance_units": 200, "duration_ms": 96000 }
  ],
  "mount_terrain_penalty_applied": true
}
```

## 10. Metrics
- `mount_terrain_encounter{terrain,species}` gauge.
- `mount_terrain_penalty_applied_total`.

## 11. Cross-Reference
- [Mount](mount.md), [Travel](travel.md), [Stats](stats.md)
- [travel/terrain.md](../travel/terrain.md), [travel/base-time.md](../travel/base-time.md)

## 12. Bekleyen Kararlar
- Winged mount Phase 02 tam terrain bypass.
- Multi-terrain hybrid species (bear=ground/mountain specialist).
- Terrain compatibility per level up (mountain wolf +%X per level Phase 02).
- Emergency swim mode (ground → water penalty formula final).
