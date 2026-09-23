# Level (Pet Level Progression)

## 1. Amaç
Pet leveling: cap, XP curve, level up rewards, trait/skill unlock.

## 2. Level Cap

- v1: **50** (matches character cap).
- Species-özel override: `content.pet_template.max_level`.
- Evolved pet: cap yeniden (evolution resets level; bkz. [evolution.md](evolution.md)).

## 3. XP Curve

Formula (species field):
```
xp_to_next(L) = base * pow(factor, L-1)
base = 80   (default)
factor = 1.12
```

Örnek:
| Level | XP Need | Cumulative |
|---|---|---|
| 1→2 | 80 | 80 |
| 5→6 | 126 | ~600 |
| 20→21 | 691 | ~7500 |
| 50 total | – | ~90k |

Pet cap 50 → 90k XP total; ~15-20 hour AFK oyun.

## 4. Level Up Rewards

Her level:
- Base stats increase (%2 per level).
- Every 5 level: unlock trait slot (Phase 02 additional trait).
- Every 10 level: skill rank up eligible.
- Level cap (50): evolution eligible (species-özel).

## 5. Level Up Flow

```go
func (s *PetService) GrantXP(ctx, petInstanceID, xp) error {
    pet := load(petInstanceID)
    pet.XP += xp
    for {
        need := xpToNext(pet.Level, pet.Species.XPCurve)
        if pet.XP < need: break
        pet.XP -= need
        pet.Level++
        if pet.Level >= pet.Species.MaxLevel: break
        applyLevelUpRewards(pet)
    }
    save(pet)
    emit event 'pet.leveled_up'
}
```

## 6. Base Stat Scale

```
pet.EffectiveAP = species.base_stats.attack_power × (1 + pet.Level × 0.02) × rarity_mult × happiness_mult
pet.EffectiveHP = species.base_stats.hp × (1 + pet.Level × 0.02) × rarity_mult
```

Rarity multiplier:
- common: 0.8
- uncommon: 1.0
- rare: 1.15
- epic: 1.3
- legendary: 1.5
- mythic: 1.75

## 7. Trait Slot Unlock

Bkz. [traits.md](traits.md).

Level 5, 15, 30 → extra trait slot açılır (Phase 02):
- Level 5: 2nd trait slot.
- Level 15: 3rd trait slot.
- Level 30: 4th trait slot.

v1'de default 2 trait sabit.

## 8. Skill Rank Up

Bkz. [skills.md](skills.md).

Pet skill ranks (1..5):
- Level 10: rank 2 unlock eligible.
- Level 20: rank 3.
- Level 30: rank 4.
- Level 40: rank 5.

Ranking için happiness ve bond level kontrolü Phase 02.

## 9. Cap Behavior

Level cap'te:
- XP overflow biriktir; evolution unlock.
- Evolution yapmadan XP birikimi kaybolur (level cap post).

## 10. Metrics
- `pet_leveled_up_total{species,level}`.
- `pet_reached_cap_total{species}`.

## 11. Cross-Reference
- [Instance](instance.md), [XP](xp.md), [Traits](traits.md), [Skills](skills.md), [Evolution](evolution.md)
- [character/level.md](../character/level.md) — pattern reference

## 12. Bekleyen Kararlar
- Level cap 50 vs character-cap sync Phase 02.
- Trait slot progression tune.
- Rarity mult final playtest.
