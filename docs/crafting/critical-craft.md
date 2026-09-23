# Critical Craft

## 1. Amaç
Nadir "critical" iteration — normal output'un ötesinde bonus: ekstra qty, rare byproduct, quality bump. Diablo-esque "crit craft" motivasyonu.

## 2. Critical Chance

Recipe field: `critical_craft_chance_pct` (default 5%).

Modifier'lar:
- Character mastery level: +0.2% per level.
- Tool bonus (crit tool +5%).
- Elixir "Focused Craft" +10% (10 min).
- Guild research (Phase 02).
- Recipe quality_bias `high` → +5%.

Cap: %30 max.

## 3. Critical Effect

Iteration critical'sa:
- **Output qty ×2** (steel_longsword × 2 yerine × 1).
- **Byproduct guaranteed** (chance byproduct'lar hepsi drop).
- **Quality bump**: current quality → +1 tier.
  - poor → normal
  - normal → exceptional
  - exceptional → **masterwork** (Phase 02 tier)

## 4. Big Miss (Rare)

Cap fail state: rare "big miss" (%1 baseline):
- Tüm input consumed (fail refund yok).
- Zero XP.
- Notification "Disaster! All materials lost."

Cap fail chance: `1% - mastery * 0.01%` (mastery koruma).

Recipe rare-fail flag (`quality_bias='high'` recipe risk yüksek).

## 5. Batch Critical Sampling

Analytical:
```go
critChance := recipe.CriticalCraftChancePct + masteryBonus + toolBonus  // %5..30
critCount := sampleBinomial(successCount, critChance/100.0, seed.sub("crit"))
// critCount kadar iteration critical output alır
```

Determinism seed'li.

## 6. Combine with Quality

Iteration order:
```
1. Roll success (binomial with craft_success_pct)
2. If success:
   a. Roll critical (binomial with critical_craft_chance_pct)
   b. Roll quality (multinomial with quality distribution)
   c. If critical: quality bump +1 tier
   d. Output qty × (critical ? 2 : 1)
3. If fail:
   a. Roll big miss (%1)
   b. Input consumed with refund policy
```

## 7. Item Instance Marking

Critical craft output item_instance:
```sql
ALTER TABLE content.item_instance
  ADD COLUMN craft_critical BOOLEAN NOT NULL DEFAULT false;
```

Tooltip'te "Crit Craft" badge (Phase 02 UI).

## 8. Event

`crafting.critical_craft` event → notification:
```json
{ "recipe_code": "steel_longsword", "iterations_critical": 2, "extra_output": [...] }
```

Client UI feedback.

## 9. Big Miss Event

`crafting.big_miss` event → dramatic notification:
```json
{ "recipe_code": "steel_longsword", "materials_lost": [...] }
```

Rare + player-visible; craft narrative.

## 10. Config

```yaml
critical_craft:
  base_chance_pct: 5
  mastery_bonus_per_level: 0.2
  tool_bonus_max: 10
  elixir_bonus: 10
  cap_pct: 30
  effect:
    qty_multiplier: 2
    quality_bump: 1        # +1 tier
    byproduct_guaranteed: true

big_miss:
  base_chance_pct: 1
  mastery_reduction_per_level: 0.01
  cap_pct: 0.5             # min chance after mastery
  input_refund_pct: 0      # zero refund
  quality_bias_high_recipe_extra_pct: 2
```

## 11. Metrics
- `craft_critical_total{recipe}` counter.
- `craft_big_miss_total{recipe}` counter.
- `craft_critical_ratio{recipe}` gauge.

## 12. Audit
- `crafting.critical_craft` (severity=info).
- `crafting.big_miss` (severity=warn — nadir).

## 13. Cross-Reference
- [Recipes](recipes.md), [Quality](quality.md), [Mastery](mastery.md), [Crafting XP](crafting-xp.md)
- [activity-engine/crafting.md](../activity-engine/crafting.md)
- [systems/random-system.md](../systems/random-system.md)

## 14. Bekleyen Kararlar
- Critical craft ×2 qty vs quality bump alternatif (v1 kombine; Phase 02 opt).
- Big miss content patch tolerance (community feedback).
- Recipe-specific critical outputs (masterwork variant) Phase 02.
- Cross-craft chain critical (5 crit in row → bonus) Phase 03.
