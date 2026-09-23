# Base Travel Time

## 1. Amaç
Distance + speed katmanlarını nihai `duration_ms`'e çevirme formulü.

## 2. Ana Formül

```
base_time_sec = distance_units / base_walking_speed_units_per_sec

effective_multiplier = mount_multiplier
                    × race_multiplier
                    × equipment_multiplier
                    × road_bonus_multiplier
                    × weather_multiplier         (Phase 02; v1=1.0)
                    / terrain_penalty_multiplier

duration_ms = int64((base_time_sec / effective_multiplier) × 1000)
```

**Not**: terrain **penalty** (arttırıcı) olduğu için bölünür; diğerleri "boost" (arttırıcı hız) çarpılır.

Genel formulasyon (bölme normalize):
```
duration_ms = (distance_units × 1000) / (base_speed × effective_multiplier)
```

## 3. Katmanların Rol Özeti

| Katman | Çarpan Türü | Kaynak | Cap |
|---|---|---|---|
| **Base Speed** | Referans | Config (1.0 units/sec) | – |
| **Mount** | Hız çarpımı | Active mount speed_multiplier | Cap 3.0× |
| **Race** | Hız çarpımı | Racial passive | Cap 1.15× |
| **Equipment** | Hız çarpımı | Boots + accessory + barding | Cap 1.20× |
| **Road Bonus** | Hız çarpımı | Edge road_bonus_pct | Cap 1.20× |
| **Terrain** | Süre penalty | Edge terrain_kind | 1.0..2.0 |
| **Weather** (Phase 02) | Hız çarpımı ↓ | Weather system | 0.5..1.0 |

Detaylar ilgili dokümanlarda.

## 4. Cap Discipline

Toplam movement speed effective cap:
- Combat: **cap +25%** (out-of-combat).
- Out-of-combat travel: **cap +50%** (v1 default).

Cap check:
```go
totalBoost := mount * race * equipment * road_bonus
if totalBoost > 3.5 {   // soft cap 3.5×
    totalBoost = 3.5
}
```

Terrain penalty **cap dışı** — yani mount+race maksimize edilse bile mountain terrain ×2 yavaşlık uygulanır.

## 5. Method Bypass

Yukarıdaki formül **walking + mount + road** için. Diğerler fixed:

- **Wagon**: `duration_ms = route.duration_ms` (content'ten sabit).
- **Ship**: `duration_ms = route.duration_ms` + schedule wait.
- **Waygate**: `duration_ms = config.waygate_cast_time_ms` (5000 ms default).

Fixed method içinde race/equipment/mount irrelevant; cost sabit gold sink.

## 6. Örnek Hesap

Karakter: Elf Ranger, mount tier 2 (2× speed), boots +5%.
Route: ashfield → emerald_coast, 500 unit, mix terrain (300 grass_road, 200 forest).

```
distance_units = 500
base_time_sec = 500 / 1.0 = 500 sec

Segment-based:
  grass_road (300u, terrain 1.0, road_bonus 1.15):
    time = 300 / (1 × 2.0 × 1.05 × 1.05 × 1.15) = 118 sec
  forest (200u, terrain 1.3, road_bonus 1.0):
    time = 200 × 1.3 / (1 × 2.0 × 1.05 × 1.05 × 1.0) = 118 sec

Total: ~236 sec ≈ 4 min
```

Basitleştirme (v1): tek toplam multiplier + terrain weighted average.

## 7. Segment-Based Compute

Route segment-by-segment hesabı doğru; ancak client karmaşıklığı ↑.

v1 tercih: **segment-based** — her edge'in kendi terrain'i uygulanır.

Formul:
```go
totalMs := int64(0)
for _, seg := range route {
    segTime := (float64(seg.DistanceUnits) * terrainPenalty(seg.TerrainKind) * 1000) /
               (baseSpeed * mountMult * raceMult * equipmentMult * roadBonus(seg.RoadBonusPct) * weatherMult)
    totalMs += int64(segTime)
}
duration_ms = totalMs
```

## 8. Preview Response

Estimate endpoint (bkz. [travel.md §3.1](travel.md)):
```json
{
  "expected_duration_ms": 236000,
  "segments": [
    { "distance_units": 300, "terrain": "grass_road", "segment_duration_ms": 118000 },
    { "distance_units": 200, "terrain": "forest", "segment_duration_ms": 118000 }
  ],
  "modifiers_applied": {
    "base_speed": 1.0,
    "mount_multiplier": 2.0,
    "race_multiplier": 1.05,
    "equipment_multiplier": 1.05,
    "road_bonus_avg": 1.09,
    "terrain_penalty_avg": 1.12,
    "total_effective_multiplier": 2.32
  }
}
```

UI display: total + segment breakdown.

## 9. Deterministic

- Aynı character snapshot + route → aynı duration.
- Server-side compute; client override yasak.

## 10. Reveal Complete Check

```go
if elapsed_ms >= contract.ExpectedDurationMs {
    complete()
}
```

Simple check; per-segment simulation yok.

## 11. Metrics
- `travel_computed_duration_ms{method}` histogram.
- `travel_multiplier_avg{layer}` gauge (mount/race/eq/terrain).

## 12. Cross-Reference
- [Travel](travel.md), [Distance](distance.md), [Mount](mount.md), [Race](race.md), [Terrain](terrain.md), [Equipment](equipment.md), [Weather](weather.md)

## 13. Bekleyen Kararlar
- Segment-based vs aggregate v1 (öneri: segment).
- Cap final (3.5× v1).
- Multi-method compose Phase 02.
- Dynamic speed change (buff mid-travel) — v1 snapshot; Phase 02 dinamik.
