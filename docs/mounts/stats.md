# Stats (Mount Stats)

## 1. Amaç
Mount'un mekanik stat'ları — speed, cargo, HP, stamina.

## 2. Ana Stats

| Stat | Baz Değer | Level Scale |
|---|---|---|
| **Speed Multiplier** | 1.5-2.5× | +2% per level |
| **Cargo Bonus Slots** | 0-10 | +0 (fixed template) |
| **Max HP** | 300-800 | +2% per level |
| **Max Stamina** | 100 | +1 per level (max 130 at level 30) |
| **Stamina Regen** | 5/min | +0.1 per level |

## 3. Effective Compute

```go
mount.EffectiveSpeed = species.base_speed_multiplier × (1 + mount.Level × 0.02)
mount.EffectiveHP = species.base_hp × (1 + mount.Level × 0.02)
mount.EffectiveStaminaMax = species.stamina_max + mount.Level × 1
mount.EffectiveCargoSlots = species.cargo_bonus_slots + barding.cargo_bonus (Phase 02)
```

## 4. Stat Modifiers

Modifier stack (mount barding Phase 02, buff, food):
- **Mount barding** (Phase 02 item): +speed, +cargo, +HP.
- **Speed potion** (character buff): +N% for X min.
- **Race passive** (opsiyonel, ör. Merfolk aquatic bonus).

Compute:
```
final_speed = mount.EffectiveSpeed × (1 + barding_bonus + buff_bonus + race_bonus)
```

Cap: total × 3.5 (bkz. [travel/base-time.md §4](../travel/base-time.md)).

## 5. Rarity Stat Multiplier

Aynı pattern as pet:
```
common: base × 0.8
uncommon: base × 1.0
rare: base × 1.15
epic: base × 1.3
legendary: base × 1.5
mythic: base × 1.75
```

## 6. Speed Multiplier Tiers

Content-driven tier examples:
```
Tier 1 (basic): 1.5×
Tier 2 (mid): 2.0×
Tier 3 (fast): 2.5×
Tier 4 (legendary): 3.0×
Tier 5 (mythic): 3.5× (cap)
```

Cap: 3.5× — bkz. [travel/base-time.md §4](../travel/base-time.md).

## 7. Cargo Bonus Slots

Bkz. [carry-capacity.md](carry-capacity.md).

Base:
- Ground mount: 0.
- Cargo mount: 5-10.
- Aquatic: 0-3 (deep sea limited).

Barding (Phase 02) +extra.

## 8. HP + Stamina

- **HP**: mount stays alive; combat interrupt (dismount) doesn't kill mount.
  - HP tracking Phase 02 mounted combat için.
- **Stamina**: travel drain rate.

Stamina drain (travel):
```
walking mount = no stamina drain
mount active travel = 1 stamina / min baseline
sprint mode (Phase 02) = 5 stamina / min + 50% speed
```

Stamina recovery: idle passive; feed treats.

## 9. Content Örneği

```yaml
# content/mount_species/dire_wolf.yaml (excerpt)
species_family: wolf
type: ground
rarity: rare
max_level: 30

speed_multiplier: 2.2
cargo_bonus_slots: 3

base_stats:
  hp: 700
  stamina_max: 110
  stamina_regen_per_min: 6

skill_codes: [mount_leap, mount_scent_track]

terrain_compatibility:
  forest: 0.85       # forest specialist
  grass: 1.0
  mountain: 1.4
  water: 6.0

acquisition_source: quest
```

## 10. Stats API

**GET `/v1/characters/{id}/mounts/{instance_id}/stats`**
```json
{
  "effective_speed_multiplier": 2.4,
  "cargo_bonus_slots": 3,
  "max_hp": 840,
  "stamina": { "current": 85, "max": 110, "regen_per_min": 6 },
  "terrain_penalties": {
    "forest": 0.85,
    "mountain": 1.4,
    ...
  }
}
```

## 11. Metrics
- `mount_speed_avg{species}` gauge.
- `mount_cargo_used_ratio` gauge.

## 12. Cross-Reference
- [Mount](mount.md), [Level](level.md), [Skills](skills.md), [Terrain](terrain.md), [Carry Capacity](carry-capacity.md), [Travel](travel.md)
- [travel/mount.md](../travel/mount.md), [travel/base-time.md](../travel/base-time.md)

## 13. Bekleyen Kararlar
- Mount barding Phase 02 slot ekleme.
- HP concept mounted combat için Phase 03.
- Sprint mode Phase 02 (temp speed burst).
- Multi-passenger mount Phase 03.
