# Happiness (Pet Happiness Stat)

## 1. Amaç
Pet care meta-stat. Karakter pet ile ilgilenirse (feed, activate, quest) happiness artar; unattended düşer.

Happiness effect: XP share multiplier, skill unlock gate, evolution eligibility.

## 2. Model

`companion.pet_owned.happiness SMALLINT 0..100`.

Baseline: 50 (new pet).

## 3. Happiness Kaynağı

Artırıcılar:
- **Feed**: +5-15 per feed (item consume; food-specific).
- **Active use** (per day active in combat): +2/day.
- **Quest completion**: +10.
- **Level up**: +5.
- **Bond milestone** (Phase 02): +10.

Azaltıcılar:
- **Inactive** (24h+ no activation): -1/day.
- **Downed in combat**: -3.
- **Retired (stable)**: -0.5/day.

Cap: 0..100.

## 4. Happiness Effect

Multiplier alanlar:
- **Pet XP share**: `0.5 + happiness/100 × 0.5` (bkz. [xp.md §3](xp.md)).
- **Combat damage**: happy pet +10%.
- **Skill effect magnitude**: +N% per happiness (Phase 02 detail).
- **Evolution eligibility**: happiness_min gate.

Formula örneği:
```
happiness_multiplier = 0.5 + happiness/100 × 0.5   // 0.5..1.0
pet_dps_effective = pet.base_dps × happiness_multiplier
```

## 5. Feed Endpoint

**POST `/v1/characters/{id}/pets/{instance_id}/feed`**
```json
{ "food_item_code": "pet_treat_basic" }
```

Server:
1. Validate food item owned.
2. Consume item.
3. Happiness += food.happiness_bonus (5-15).
4. Cooldown: 6 hours per pet (over-feed protection).
5. Emit event `pet.fed`.

Cooldown Redis: `pet:feed:<instance_id>` TTL 6h.

## 6. Food Items

Content:
```yaml
# content/items/pet_treats/pet_treat_basic.yaml
code: pet_treat_basic
category: consumable
subcategory: pet_food
happiness_bonus: 5
species_bonus:
  beast: 2      # +2 extra for beast pets (dietary preference)
```

Species-specific treats bonus.

## 7. Happiness Decay Worker

Worker `pet_happiness_decay` (daily):
```
UPDATE companion.pet_owned
SET happiness = GREATEST(0, happiness - 1)
WHERE stable_since IS NULL
  AND (last_feed_at IS NULL OR last_feed_at < now() - interval '24 hours')
  AND (last_activation_at IS NULL OR last_activation_at < now() - interval '24 hours');
```

## 8. Happiness Display

Client tooltip:
- 100 → "Ecstatic" (green glow).
- 80-99 → "Happy" (green).
- 50-79 → "Content" (yellow).
- 25-49 → "Unhappy" (orange).
- 0-24 → "Miserable" (red).

## 9. Metrics
- `pet_happiness_avg` gauge.
- `pet_fed_total{food}`.
- `pet_happiness_decay_applied_total`.

## 10. Cross-Reference
- [Instance](instance.md), [XP](xp.md), [Evolution](evolution.md), [Combat](combat.md)
- [inventory/inventory.md](../inventory/inventory.md) — consumable food

## 11. Bekleyen Kararlar
- Feed cooldown 6h vs 4h.
- Happiness cap 100 vs 150 (with bond).
- Species-specific food full catalog Phase 02.
- Auto-feed subscription Phase 03 (idle-friendly).
