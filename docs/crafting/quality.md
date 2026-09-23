# Quality (Output Quality Tier)

## 1. Amaç
Crafting output'unun kalite tier'ı — mastery + luck bazlı. Poor/Normal/Exceptional (v1); Phase 02'de daha zengin.

## 2. Quality Tier'ları

| Tier | Rarity Bump | Stat Bonus | Chance Distribution (baseline) |
|---|---|---|---|
| **poor** | -1 rarity | ×0.7 | %10 |
| **normal** | template.rarity | ×1.0 | %80 |
| **exceptional** | +1 rarity | ×1.2 | %10 |

Rarity bump: template.rarity uncommon ise `poor=common`, `exceptional=rare`.

## 3. Chance Distribution

Baseline:
```
poor: 10%
normal: 80%
exceptional: 10%
```

Modifier'lar:
- Character mastery level (bkz. [mastery.md](mastery.md)).
- Recipe rank (Phase 02).
- Tool bonus.
- Elixir buff.

Örnek high-mastery blacksmith:
```
Mastery 80 (Master):
  poor: 3%
  normal: 60%
  exceptional: 37%
```

## 4. Quality Formula

```go
func rollQuality(char, recipe, seed) Quality {
    poorBase := 10.0
    exceptionalBase := 10.0
    normalBase := 80.0

    // Mastery bias
    masteryLevel := char.MasteryLevels[recipe.ProfessionID]
    exceptionalBase += masteryLevel * 0.3       // 0.3% per mastery level
    poorBase -= masteryLevel * 0.1                // 0.1% reduction per level

    // Recipe quality_bias
    switch recipe.QualityBias {
    case "high":
        exceptionalBase += 5
        poorBase -= 3
    case "low":
        exceptionalBase -= 3
        poorBase += 3
    }

    // Tool bonus
    exceptionalBase += tool.QualityBonusPct

    // Clamp
    exceptionalBase = clamp(exceptionalBase, 0, 60)
    poorBase = clamp(poorBase, 0, 20)
    normalBase = 100 - exceptionalBase - poorBase

    // Roll
    r := seed.NextFloat() * 100
    if r < poorBase { return Poor }
    if r < poorBase + normalBase { return Normal }
    return Exceptional
}
```

## 5. Quality Etkisi

Equipment output için:
- **Poor**: rarity -1 tier; rolled stats × 0.7; UI görsel dull.
- **Normal**: template.rarity aynen; rolled stats normal.
- **Exceptional**: rarity +1 tier; rolled stats × 1.2; UI görsel sparkle.

Consumable / material output için:
- **Poor**: potion effect × 0.8, duration × 0.8.
- **Normal**: standart.
- **Exceptional**: effect × 1.2, duration × 1.1.

## 6. Item Instance ile Entegrasyon

Crafted equipment output item_engine.Create:
- `quality` alanı item_instance'a (yeni field Phase 16):
  ```sql
  ALTER TABLE content.item_instance
    ADD COLUMN craft_quality VARCHAR(16) NULL;  -- 'poor'|'normal'|'exceptional'|NULL
  ```
- Stat compute'da quality multiplier uygulanır.
- Client tooltip'te "Exceptional" tag + rarity color glow.

Non-crafted item (loot drop) `craft_quality=NULL` — default normal.

## 7. Batch Craft Quality Distribution

Batch iteration → analytical sample:
```
successCount = 19
qualityDist = sampleMultinomial(successCount, {poor: 3, normal: 60, exceptional: 37}, seed)
=> { poor: 1, normal: 11, exceptional: 7 }
```

`sampleMultinomial` = N independent binomial reduced; single-call determinism seed'li.

Output:
- 1 poor steel_longsword.
- 11 normal steel_longsword.
- 7 exceptional steel_longsword.

## 8. Exceptional as Rare Drop

Exceptional item ekonomide değerlidir:
- Vendor price ×2.
- Marketplace listing için "Exceptional" filter.
- Guild bank görsel işaret.

## 9. Fail (No Output)

Craft fail edilirse quality irrelevant (no output).

Bkz. [crafting-xp.md](crafting-xp.md) fail XP handling.

## 10. Quality Bias per Recipe

`recipe.quality_bias`:
- `low`: basic recipes (poor bias); learn material için.
- `balanced` (default): normal distribution.
- `high`: masterwork recipes (exceptional bias); rare recipe.

## 11. Config

```yaml
# config.crafting.quality.yaml
quality_distribution:
  base:
    poor: 10
    normal: 80
    exceptional: 10
  mastery_bonus_exceptional_per_level: 0.3
  mastery_reduction_poor_per_level: 0.1
  tool_quality_bonus_pct_max: 15
  clamp:
    exceptional_max: 60
    poor_max: 20

quality_multipliers:
  poor: 0.7
  normal: 1.0
  exceptional: 1.2
  poor_duration_mult: 0.8
  exceptional_duration_mult: 1.1
```

## 12. Metrics
- `craft_output_quality_total{quality,recipe}`.
- `craft_exceptional_ratio{recipe}` gauge.

## 13. Cross-Reference
- [Recipes](recipes.md), [Materials](materials.md), [Critical Craft](critical-craft.md), [Mastery](mastery.md)
- [item-engine/rarity.md](../item-engine/rarity.md), [item-engine/item-instance.md](../item-engine/item-instance.md)
- [database/content-templates.md §9](../database/content-templates.md)

## 14. Bekleyen Kararlar
- Quality tier expansion (masterwork, legendary quality) Phase 02.
- Guild bonuses (guild research → quality bump) Phase 02.
- Cross-quality material use (exceptional ingot → higher chance exceptional weapon) Phase 02.
