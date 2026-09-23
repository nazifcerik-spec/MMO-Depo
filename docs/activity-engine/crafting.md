# Crafting Activity

## 1. Amaç
Recipe-based üretim aktivitesi (blacksmithing, alchemy, cooking, enchanting, jewelcrafting, leatherworking, tailoring). Input consume, output produce.

## 2. Params

```json
{
  "recipe_id": "...",
  "iterations": 20,               // kaç kez craft
  "station_location_id": "..."    // opsiyonel; recipe requires station ise
}
```

## 3. Prerequisites

- Character `active_professions` içinde recipe.profession.
- `character.profession[X].level >= recipe.tier_required`.
- Character `character.recipe_known` içinde recipe var.
- Station required ise character at station location.
- Envanterde recipe.inputs yeterli miktarda.

## 4. Rate Compute

Crafting rate = `1 iteration / recipe.time_ms`.

```go
iterationsPerHour := 3600000.0 / float64(recipe.TimeMs)
craftingSuccessPct := char.Derived.CraftSuccessPct  // 85 base + WIS bonus + tool
```

Snapshot içinde `craft_success_pct` derived stat.

## 5. Yield Compute (Analytical)

```go
maxIterations := min(params.Iterations, availableInputs())
elapsedIterations := min(maxIterations, int(iterationsPerHour * elapsedHours))

successCount := sampleBinomial(elapsedIterations, craftingSuccessPct/100.0, seed.Sub("success"))
failCount := elapsedIterations - successCount

outputs := []ItemDrop{}
for _, out := range recipe.Outputs {
    baseCount := successCount * out.Qty
    if out.Chance < 1.0 {
        // proc-based output (rare byproduct)
        actualCount := sampleBinomial(baseCount, out.Chance, seed.Sub("proc-"+out.Code))
        outputs = append(outputs, ItemDrop{TemplateCode: out.Code, Quantity: int32(actualCount)})
    } else {
        outputs = append(outputs, ItemDrop{TemplateCode: out.Code, Quantity: int32(baseCount)})
    }
}

// Input consumption
inputs := []ItemConsumption{}
for _, in := range recipe.Inputs {
    // Fail'da %50 material refund
    fullConsume := successCount * in.Qty
    failConsume := failCount * in.Qty / 2
    totalConsume := fullConsume + failConsume
    inputs = append(inputs, ItemConsumption{TemplateCode: in.Code, Quantity: totalConsume})
}

xpPerIter := recipe.XPReward
totalXP := successCount * xpPerIter + failCount * xpPerIter/2

return Yield{
    XPCharacter: int64(rate.XPPerHour * elapsedHours),
    XPProfession: map[uuid.UUID]int64{ recipe.ProfessionID: totalXP },
    ItemDrops: outputs,
    ItemConsumptions: inputs,
    // ... rest
}
```

**Not**: `sampleBinomial` binomial distribution single sample; per-iteration loop yok.

## 6. Threshold Detection

Input tükenme:
```
iterations_until_input_exhaust = min(availableInputs[X] / recipe.Inputs[X].Qty  for each X)
ms_until_input_exhaust = iterations_until_input_exhaust / iterationsPerHour * 3600000
```

Iterations tamamlanma:
```
ms_to_complete_all = params.Iterations / iterationsPerHour * 3600000
```

Output inventory dolma:
```
each successful output ⇒ slots used
threshold when inv_free < remaining outputs required
```

Return smallest threshold.

## 7. Content: Recipe

Bkz. [database/content-templates.md §7](../database/content-templates.md).

```yaml
# content/recipes/steel_longsword.yaml
code: steel_longsword
profession_code: blacksmithing
tier_required: 15
inputs:
  - item_code: steel_ingot
    qty: 3
  - item_code: leather_grip
    qty: 1
outputs:
  - item_code: steel_longsword_template
    qty: 1
    chance: 1.0
  - item_code: metal_scrap
    qty: 1
    chance: 0.2
time_ms: 8000
xp_reward: 40
station_required_code: forge
learn_source: trainer
```

## 8. Item Instance for Crafted Weapon

Crafted equipment (weapon, armor) → item_engine.Create ile instance yaratılır (bkz. [item-engine/item-instance.md §4](../item-engine/item-instance.md)):
- Seed = hash(contract.Seed, iteration_index).
- Rolled stats + affix roll normal loot pattern gibi.
- Owner = character_inventory.

Crafted item genelde BoE (recipe binding_type override).

## 9. Craft Success Chance

Base: 85%.
Bonus:
- `+0.5 per WIS`
- `+tool bonus`
- `+rank bonus (recipe learn rank Phase 02)`
- `+event bonus`

Cap: 100%.

## 10. Lifecycle

- **Start**: params.iterations set; snapshot al.
- **Reveal**: analytical yield based on time × craft rate; check input exhaust threshold.
- **Complete**: iterations tamamlanınca contract auto-complete.
- **Cancel**: partial iteration settle.

## 11. Reveal Örnek

```
Crafting: 20× steel_longsword, 8s each = 160s total
Elapsed: 5 min = 300s = ~37 iteration; capped at 20.

Success chance: 92% → binomial sample = 18 success, 2 fail
Outputs:
  - steel_longsword instance × 18 (each yeni seed'li roll)
  - metal_scrap × 3 (0.2 chance × 18 success = 3.6 → binomial sample = 3)

Inputs consumed:
  - steel_ingot: 18*3 + 2*3/2 = 54 + 3 = 57
  - leather_grip: 18 + 1 = 19

XP:
  - character: 40 × 18 + 20 × 2 = 760
  - blacksmithing: 400 × 18 + 200 × 2 = 7600

Contract auto-complete after 20 iterations.
```

## 12. Metrics
- `activity_crafting_iterations_total{recipe}`.
- `activity_crafting_success_rate{recipe}` histogram.
- `activity_crafting_input_shortage_total`.

## 13. Cross-Reference
- [Activity](activity.md), [Rewards](rewards.md), [XP](xp.md)
- [game-design/professions.md](../game-design/professions.md)
- [item-engine/item-instance.md](../item-engine/item-instance.md)
- [database/content-templates.md §7](../database/content-templates.md)

## 14. Bekleyen Kararlar
- Craft quality tier (rare item output based on success streak) Phase 02.
- Batch cook (multiple recipes queued) Phase 02.
- Failure downgrade (crit-fail → destroyed input) Phase 02 opsiyonel.
