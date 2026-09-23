# Race (Travel Speed Contribution)

## 1. Amaç
Karakterin ırkının travel hızına passive katkısı.

## 2. Race Speed Bonusları (Placeholder — Bkz. Race Docs)

Race passive'lerinden movement_speed_pct modifier'ları:

| Race | Passive | Movement Bonus | Scope |
|---|---|---|---|
| `human` | human_versatility | 0 | – |
| `dwarf` | dwarven_endurance | +2% | always |
| `elf` | elven_grace | +5% | always |
| `orc` | orcish_bloodlust | 0 | – |
| `merfolk` | tide_born | +10% | scope: aquatic |
| `arcanborn` | arcane_attunement | +2% | scope: magical_terrain |

Race passive modifier stack'e katılır (bkz. [skill-engine/modifiers.md](../skill-engine/modifiers.md)).

## 3. Race Multiplier Formula

```
race_multiplier = 1.0 + racial_movement_bonus_pct / 100.0
```

Örnek elf:
```
race_multiplier = 1.0 + 5/100 = 1.05
```

## 4. Scope-Aware

Bazı race bonusu koşula bağlı:
- Merfolk +10% aquatic terrain (water crossing).
- Elf +5% always (forest bonus opsiyonel Phase 02).

Compute:
```go
for _, mod := range char.PassiveModifiers {
    if mod.Target == "movement_speed_pct" {
        if scopeMatches(mod.Scope, currentContext) {
            racialBonus += mod.Amount
        }
    }
}
race_multiplier = 1.0 + racialBonus/100.0
```

Segment-based: her edge terrain'e bakılır; race bonus scope match kontrol.

## 5. Merfolk Aquatic Bonus

Merfolk water crossing:
```
Terrain=water:
  race_multiplier = 1.10   (aquatic bonus)
  ancak aquatic mount kombine: mount_multiplier × 1.10
```

## 6. Race + Mount Kombine

Multiplikatif:
```
combined = mount_multiplier × race_multiplier

örnek: elf + mount tier 2 = 2.0 × 1.05 = 2.10
```

Cap toplam bkz. [base-time.md §4](base-time.md).

## 7. Race Change (Phase 02)

Race change (Phase 02+) travel bonusu değişir.
- v1'de race değişmez.
- Contract snapshot'ta race locked.

## 8. Cross-Reference
- [Travel](travel.md), [Base Time](base-time.md), [Terrain](terrain.md), [Mount](mount.md)
- [race/passive.md](../race/passive.md), [race/racial-bonuses.md](../race/racial-bonuses.md)
- [game-design/races.md](../game-design/races.md)

## 9. Metrics
- `travel_race_bonus_applied_avg{race}` gauge.

## 10. Bekleyen Kararlar
- Race bonus magnitude final playtest.
- Winged race (Phase 03+) — cross-terrain travel.
- Racial travel unique mechanic (dwarf dark_vision → mountain tunnel Phase 02).
