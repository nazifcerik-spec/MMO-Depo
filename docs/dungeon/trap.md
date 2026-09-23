# Trap

## 1. Amaç
Dungeon room içinde tuzak. Skill/luck check → başarısızsa damage veya CC.

## 2. Trap Türleri

| Kind | Etki | Check |
|---|---|---|
| `floor_spikes` | Physical damage | DEX check |
| `poison_dart` | Poison DoT | DEX check |
| `arcane_ward` | Silence + mana burn | INT check |
| `pit` | Fall damage | DEX or luck check |
| `ambush_trigger` | Spawn mobs | Luck check |
| `curse_glyph` | Debuff -20% dmg 30s | WIS check |

## 3. Trap Content

```yaml
# content/traps/floor_spikes.yaml
code: floor_spikes
kind: floor_spikes
damage_element: physical
damage_base: 200
damage_scaling_pct_max_hp: 15         # 15% max_hp damage
check_kind: dex
check_threshold: 30
check_chance_pct_bonus: 0             # additional %
bypass_class_codes: [rogue]            # Rogue skips trap
```

## 4. Trap Check

Character check:
```
check_pass_chance = (character.stat.dex - trap.check_threshold) / (trap.check_threshold × 2) × 100
check_pass_chance = clamp(check_pass_chance, 5, 95)   // %5 min, %95 max

if character.class in trap.bypass_class_codes:
    check_pass = true

// Roll
if seed.NextFloat() < check_pass_chance/100:
    passed = true
else:
    passed = false
```

Party: each character rolls independently.

## 5. On Fail

- Damage: `damage_base × (1 + damage_scaling_pct_max_hp/100 × character.max_hp / 1000)`.
- Or CC: apply debuff (silence, slow).
- Both: mixed effect.

## 6. Rogue Trap Detection (Phase 02)

Rogue class skill:
- "Detect Traps" passive → +30% pass chance.
- "Disarm Trap" active → guaranteed pass + %5 chest yield bonus.

v1'de bypass class flag simple.

## 7. Trap Advance

Trap encounter always advances (pass or fail):
- Pass: no damage, next room.
- Fail: damage/CC applied; advance.

Trap doesn't wipe run (usually); safety mechanic.

## 8. AFK Trap Compute

Reveal:
```
for each trap room:
    for each character:
        pass_chance = compute(...)
        passed = seed.NextFloat() < pass_chance/100
        if !passed:
            damage = ...
            character.current_hp -= damage   // AFK context
            if character.current_hp <= 0:
                death event
```

Bounded (per trap room per character).

## 9. Metrics
- `dungeon_trap_triggered_total{trap}`.
- `dungeon_trap_failed_total{trap}`.
- `dungeon_trap_death_total`.

## 10. Cross-Reference
- [Room](room.md), [Encounter](encounter.md), [AFK Run](afk-run.md)
- [combat/status-effects.md](../combat/status-effects.md) — CC apply

## 11. Bekleyen Kararlar
- Rogue disarm skill Phase 02.
- Trap variety content depth Phase 02.
- Party trap synchronization (all fail → wipe risk) — v1 individual check.
- Trap loot bonus (skip trap → chest bonus) Phase 02.
