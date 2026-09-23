# Fishing Activity

## 1. Amaç
Karakterin belirli bir fishing spot'unda balık ve deniz mahsulü toplaması. Gathering profession.

## 2. Params

```json
{
  "spot_id": "loc-emerald-coast",
  "autoloot": true,
  "bait_item_code": "worm_bait",       // opsiyonel; bait bonus
  "target_fish_types": ["all"]         // opsiyonel filter
}
```

## 3. Prerequisites (ValidateStart)

- Character `active_professions` içinde `fishing`.
- `character.profession[fishing].level >= spot.min_profession_level`.
- `character.current_location_id == spot_id`.
- Inventory capacity > 0.
- Bait item envanterde varsa (opsiyonel — bonus yield).

## 4. Rate Compute

Mining paterninin aynısı, `fishing` katsayılarıyla:
```go
baseYieldPerHour := config.rates.fishing.baseYieldPerHourByTier[spot.Tier]
yieldMult := 1.0 + mods.SumAdditive("gather_yield_pct") / 100.0
if params.BaitCode != "" && char.HasItem(params.BaitCode) {
    yieldMult += config.rates.fishing.baitBonusPct / 100.0
}
```

Fishing için race passive: `merfolk` → `tide_born` +%5 fishing yield.

## 5. Yield Compute (Analytical)

```go
elapsedHours := float64(elapsedMs) / 3600000.0
totalCatches := samplePoisson(rate.YieldPerHour, elapsedHours, seed.Sub("catches"))

// Bait consumption
baitUsed := 0
if params.BaitCode != "" {
    baitPerCatch := config.rates.fishing.baitPerCatch  // 0.1 (10% catches use bait)
    baitUsed = int(float64(totalCatches) * baitPerCatch)
    if baitUsed > char.BaitCount {
        // bait tükendi; ancak per-tick simüle etmiyoruz.
        // t_bait_exhaust = char.BaitCount / (rate * baitPerCatch)
        // Threshold ve elapsed_effective ile bait-out senaryosu yönetilir (bkz. offline-progress §Threshold).
    }
}

// Fish distribution — spot'a bağlı çoklu balık türü
catches := []ItemDrop{}
for fishCode, share := range spot.CatchDistribution {
    count := int32(float64(totalCatches) * share)
    catches = append(catches, ItemDrop{TemplateCode: fishCode, Quantity: count})
}

// Rare catch (nadir balık, sea trinket)
rareCount := samplePoisson(rate.LootTableSampleRate, elapsedHours, seed.Sub("rare"))
rareDrops := lootSystem.MultiRoll(spot.RareLootTableID, rareCount, seed.Sub("rare_pool"))

return Yield{
    XPCharacter:  int64(rate.XPPerHour * elapsedHours),
    XPProfession: map[uuid.UUID]int64{ fishingProfessionID: int64(rate.ProfessionXPPerHour * elapsedHours) },
    ItemDrops:    append(catches, rareDrops...),
    // Bait consume (item delta negative):
    // ... (rewards.md içinde item consumption pattern)
}
```

## 6. Threshold Detection

Inventory dolma:
```
inventory_fill_ms = (char.InventoryCap - char.InventoryUsed) / rate.YieldPerHour * 3600000
```

Bait tükenme:
```
if params.BaitCode != "" {
    catches_until_bait_out = char.BaitCount / config.rates.fishing.baitPerCatch
    ms_until_bait_out = catches_until_bait_out / rate.YieldPerHour * 3600000
}
```

Return smallest hit threshold. Bait tükenirse `pause_reason='bait_exhausted'` (opsiyonel; v1: bait auto-off, continue without bonus).

## 7. Content: Fishing Spot

```yaml
# content/locations/emerald_coast_pier.yaml
code: emerald_coast_pier
kind: fishing_spot
region_code: emerald_coast
tier: 1
min_profession_level: 1

catch_distribution:
  common_trout: 0.5
  common_bream: 0.3
  emerald_pike: 0.15
  spiny_crab: 0.05

rare_loot_table_code: emerald_coast_rare_fishing
rare_proc_per_hour: 0.25
```

## 8. Bait Sistemi

Bait item consumable:
- Envanterde bulunmalı.
- Rate bonus (config +%20 yield).
- Consume rate: catches × baitPerCatch.
- Rare bait (Phase 02) → farklı bonus (rare drop chance, specific fish bias).

## 9. Merfolk Bonus

`tide_born` passive:
- `+5% gather_yield_pct scope:profession:fishing`.
- Modifier stack automatic.

Aquatic mount:
- `merfolk` characters + aquatic mount → additional fishing rate (Phase 02).

## 10. Lifecycle

Aynı mining pattern (bkz. [mining.md §10](mining.md)).

## 11. Reveal Örneği

```
Fishing @ Emerald Coast, 6h elapsed
Merfolk char, fishing lvl 25, bait (worm x100)
Rate: 40 base × 1.05 (racial) × 1.20 (bait) = 50.4/hour

Poisson sample: 305 catches
Bait use: 30.5 (rounded 31) → within 100 stock (OK)

Distribution:
  common_trout: 152
  common_bream: 91
  emerald_pike: 45
  spiny_crab: 15

Rare loot (0.25/h × 6h = 1.5 → poisson sample = 2 rare drops):
  - "Pearl of the Deep" x1
  - "Ancient Anchor Fragment" x1

XP: 40 × 6 × 4 = ~960 character XP + 4800 fishing profession XP
Inventory used: 305 items into slots (many stacks; auto-merge)
```

## 12. Metrics
- `activity_fishing_yield_avg{spot_tier}`.
- `activity_fishing_bait_consumed_total`.
- `activity_fishing_rare_catch_total`.

## 13. Cross-Reference
- [Mining](mining.md), [Activity](activity.md), [Offline Progress](offline-progress.md)
- [game-design/professions.md](../game-design/professions.md)
- [race/racial-bonuses.md](../race/racial-bonuses.md) — merfolk

## 14. Bekleyen Kararlar
- Bait mekaniği v1 aktif mi (öneri: evet; basit; item consume + rate bonus).
- Weather bonus (rainy = +yield) Phase 02.
- Fishing tournament event (Phase 02).
