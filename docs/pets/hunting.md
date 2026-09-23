# Hunting (Pet Hunting Contribution)

## 1. Amaç
Pet aktif iken hunting AFK aktivitesine katkı — DPS boost, loot bonus, gathering assist.

## 2. Combat Pet Hunt Contribution

Bkz. [combat.md §8](combat.md):
- Pet DPS contribution — kills faster.
- Threat: auxiliary (owner primary tank).

## 3. Utility Pet Hunt Contribution

Utility slot pet:
- `pet_forage_boost`, `pet_carry_capacity`, `pet_lucky_paw` passives.
- Gathering yield or inventory bonus.

Hunt sırasında:
- Loot rare drop chance +N%.
- Inventory slot +N (delays inv full threshold).

## 4. Species Affinity for Hunt

- **Beast pets** (wolf, bear): combat DPS bias.
- **Avian pets** (falcon, owl): scouting → rare mob detection bonus.
- **Reptile pets** (viper, gecko): DoT bonus.
- **Insect pets** (spider, wasp): poison bonus.

Species-specific abilities per pet skill.

## 5. Pet Skinning Bonus

Hunting profession + skinning:
- Pet aktif kill sonrası pet-assisted skinning bonus:
  - `pet_hunter_instinct` trait: +5% skinning yield.
  - Beast species: +3% skin drop chance.

## 6. AFK Reveal Integration

Reveal içinde hunting yield:
```go
hunt_yield := computeBaseHuntYield(rate, elapsed)
if pet.IsActive:
    pet_boost := computePetBoost(pet, hunting_context)
    hunt_yield.XPCharacter *= (1 + pet_boost.xp_pct/100)
    hunt_yield.Kills *= (1 + pet_boost.dps_pct/100)
    hunt_yield.Loot.RareDropChance += pet_boost.rare_drop_bonus
```

## 7. Pet XP Return

Hunting sırasında pet XP earn (bkz. [xp.md §2](xp.md)):
- 20% char XP → pet.
- Happiness multiplier applied.

## 8. Stamina Consume

Pet AFK hunting:
- -2 stamina per hour.
- Threshold: pet exhausts → contribution stops (bkz. [stamina.md §6](stamina.md)).

## 9. Rare Mob Detection (Phase 02)

Avian pet special:
- `pet_scout_ability`: rare mob spawn chance +5%.
- Reveal payload highlight.

v1'de yok; species trait bias only.

## 10. Cross-Reference
- [Combat](combat.md), [Skills](skills.md), [Traits](traits.md), [Stamina](stamina.md), [XP](xp.md)
- [activity-engine/hunting.md](../activity-engine/hunting.md)
- [hunting-adventure/hunting.md](../hunting-adventure/hunting.md)

## 11. Metrics
- `pet_hunt_contribution_avg_dps_pct` gauge.
- `pet_hunt_active_total`.

## 12. Bekleyen Kararlar
- Rare mob detection Phase 02.
- Multi-pet hunt (combat + utility eşzamanlı) — v1 desteklenir; contribute limits.
- Species affinity per hunt spot (Phase 02 zone-specific bonus).
