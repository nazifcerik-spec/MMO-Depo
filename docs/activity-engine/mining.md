# Mining Activity

## 1. Amaç
Karakterin belirli bir mining spot'unda ore/stone/gem topluması. Profession-based gathering aktivitesi.

## 2. Params

```json
{
  "spot_id": "loc-iron-valley",       // required
  "autoloot": true,                    // default true; false ise sadece XP yield
  "target_node_tier_min": 1,           // opsiyonel filter
  "target_node_tier_max": 3
}
```

## 3. Prerequisites (ValidateStart)

- Character `active_professions` içinde `mining` var.
- `character.profession[mining].level >= spot.min_profession_level`.
- `character.current_location_id == params.spot_id` (spot'a fiziksel gelmiş olması).
- Inventory capacity > 0.
- Aktif primary channel contract yok (veya bu aktiviteyle değiştirilecek).

Fail modları:
- `mining.profession_not_active`
- `mining.level_insufficient`
- `mining.not_at_spot`
- `mining.inventory_full` (soft; auto-pause zaten var)

## 4. Rate Compute

```go
func (a *MiningActivity) ComputeRate(ctx, snap, params) Rate {
    spot := content.MiningSpots.ByID(params.SpotID)
    baseYieldPerHour := config.rates.mining.baseYieldPerHourByTier[spot.Tier]

    // Modifier stack apply (gather_yield_pct, profession affinity)
    mods := modifierEngine.Gather(ctx, snap, StatContext{Purpose:"activity", Activity:"mining"})
    yieldMult := 1.0 + mods.SumAdditive("gather_yield_pct") / 100.0

    // Profession mastery bonus
    mLevel := snap.MasteryLevels[miningTrackID]
    yieldMult *= (1.0 + float64(mLevel) * 0.001)  // 100 mastery = +%10

    // Rare gem sample rate — proc rate (per hour)
    rareProcRate := spot.RareProcPerHour * (1.0 + mods.SumAdditive("rarity_bonus_pct")/100.0)

    return Rate{
        YieldPerHour:        baseYieldPerHour * yieldMult,
        XPPerHour:           config.rates.mining.xpPerYield * baseYieldPerHour * yieldMult,
        ProfessionXPPerHour: 10 * baseYieldPerHour * yieldMult,  // profession xp katsayısı
        MasteryXPPerHour:    2 * baseYieldPerHour * yieldMult,
        LootTableSampleRate: rareProcRate,
    }
}
```

## 5. Yield Compute (Analytical, No Per-Tick)

```go
func (a *MiningActivity) ComputeYield(ctx, contract, snap, elapsedMs) Yield {
    rate := a.ComputeRate(ctx, snap, contract.Params)
    elapsedHours := float64(elapsedMs) / 3600000.0

    // Base yield — expected count
    expectedYield := rate.YieldPerHour * elapsedHours

    // Variance: Poisson distribution around expected (sample once, not per-tick)
    actualYield := samplePoisson(rate.YieldPerHour, elapsedHours, contract.Seed.SubStream("yield"))

    // Loot table roll — rare drops (gems, ancient ore)
    rareDropCount := samplePoisson(rate.LootTableSampleRate, elapsedHours, contract.Seed.SubStream("rare"))
    rareDrops := []ItemDrop{}
    for i := 0; i < rareDropCount; i++ {
        drop := lootSystem.RollTable(spot.RareLootTableID, contract.Seed.SubStream(fmt.Sprintf("rare-%d", i)))
        rareDrops = append(rareDrops, drop...)
    }

    // Primary ore + secondary stone (weighted distribution)
    primaryOreCount := actualYield * spot.PrimaryOreShare
    secondaryStoneCount := actualYield * spot.SecondaryStoneShare

    return Yield{
        XPCharacter: int64(rate.XPPerHour * elapsedHours),
        XPProfession: map[uuid.UUID]int64{
            miningProfessionID: int64(rate.ProfessionXPPerHour * elapsedHours),
        },
        XPMastery: map[uuid.UUID]int64{
            miningMasteryTrackID: int64(rate.MasteryXPPerHour * elapsedHours),
        },
        ItemDrops: append(
            []ItemDrop{
                {TemplateCode: spot.PrimaryOreCode, Quantity: int32(primaryOreCount)},
                {TemplateCode: spot.SecondaryStoneCode, Quantity: int32(secondaryStoneCount)},
            },
            rareDrops...,
        ),
    }
}
```

**Not**: `samplePoisson` deterministik (seed'li); tek call, per-tick loop yok. Detay: [offline-progress.md §Variance Sampling](offline-progress.md).

## 6. Threshold Detection

Inventory dolma zamanı:
```
inventory_capacity_free = char.InventoryCap - char.InventoryUsed
yield_per_hour = rate.YieldPerHour
inventory_fill_hours = inventory_capacity_free / yield_per_hour
inventory_fill_ms = inventory_fill_hours * 3600000

if inventory_fill_ms < elapsedMs:
    return Threshold{
        Kind: "inventory_full",
        HitAtMs: inventory_fill_ms,
    }
```

Node dry (opsiyonel; v1'de personal instance sonsuz):
- Personal instance → dry yok.
- Shared world node (rare) → depletion timer.

Return `nil` → threshold hit yok, full elapsed.

## 7. Snapshot Alanları

Mining için gerekli snapshot alanları:
- `MasteryLevels[mining_track_id]`
- `ProfessionLevels[mining_profession_id]`
- Race passive (`dwarven_endurance` etc.) modifier
- Ascension node (gather bonuses)
- Gear (mining tool item's yield bonus)
- Active buffs (mining potion)

## 8. Content: Mining Spot

```yaml
# content/locations/iron_valley.yaml
code: iron_valley
kind: mining_spot
region_code: ashen_hills
tier: 2
min_profession_level: 20

primary_ore_code: iron_ore
secondary_stone_code: gray_stone
primary_ore_share: 0.75
secondary_stone_share: 0.25

rare_loot_table_code: iron_valley_rare
rare_proc_per_hour: 0.5

art_ref: locations/iron_valley.png
```

## 9. Pet Assist

Aktif combat pet varsa mining aktivitesinde gather pet slot kontribüsyonu:
- Snapshot'ta `ActivePetInstanceID` var.
- Rate compute'da pet_yield_bonus_pct eklenir (bkz. [game-design/pets.md](../game-design/pets.md)).

## 10. Lifecycle Events

- **Start**: contract create, snapshot al.
- **Pause**: yield freeze at pause moment; resume'da devam.
- **Resume**: new snapshot alınabilir (opsiyonel; v1'de resume snapshot refresh yok — mevcut kullanılır).
- **Cancel**: partial yield settle.
- **Complete**: full yield settle.

Detay: [lifecycle.md](lifecycle.md).

## 11. Reveal Örnek

```
Contract: mining @ iron_valley, started 8h ago, 40 inv used out of 40 cap.
Snapshot: dwarf, mining lvl 42, mastery 60.
Rate: 90 base × 1.05 (racial) × 1.06 (mastery) × 1.05 (gear) = 105/hour
Inventory free = 0 slots → immediate threshold.
Yield: 0 ore. Pause reason: inventory_full.
```

Sonrası oyuncu envanteri boşaltıp resume:
```
elapsed since resume = 6h
Rate: 105/hour
Actual (poisson sample): 632 ore
Threshold: inv full at ~600 → yield 600, pause again.
```

## 12. Metrics
- `activity_mining_yield_avg{spot_tier}`.
- `activity_mining_paused_inventory_full_total`.
- `activity_mining_rare_proc_total`.

## 13. Cross-Reference
- [Activity](activity.md), [Fishing](fishing.md), [Offline Progress](offline-progress.md), [Rewards](rewards.md)
- [game-design/professions.md](../game-design/professions.md)
- [game-design/hunting.md](../game-design/hunting.md)
- [item-engine/README.md](../item-engine/README.md)

## 14. Bekleyen Kararlar
- Rare loot table entries content (v1 örnek).
- Multi-node aggregation (bir spot'ta N node arasında rotation) — v1'de single spot yield.
- Weather effect (Phase 02) — mining rate modifier.
- Tool durability drain per yield (v1'de tool durability yok).
