# Crafting Queue (Batch Iteration)

## 1. Amaç
Karakter N iteration crafting başlatır; server AFK contract olarak yürütür. Analytical yield (per-iteration loop yok).

Foundation: [activity-engine/crafting.md](../activity-engine/crafting.md).

## 2. Contract Model

`activity.contract` `activity_type='crafting_generic'`:
- `channel='primary'`.
- `params_json`:
  ```json
  {
    "recipe_id": "...",
    "iterations": 20,
    "station_location_id": "loc-ashfield-forge",
    "auto_pause_on_output_full": true
  }
  ```

Snapshot: WIS, craft_success_pct, mastery level, tool bonus.

## 3. Batch Iteration Model

Kullanıcı 20 iteration ister:
```
Total time = time_ms × iterations = 8000 × 20 = 160000 ms (2.67 min)
```

Aktivite süresi hesaplanır; contract'a persist.

Reveal içinde:
```
elapsed_ms
iterations_processed = min(params.iterations, floor(elapsed_ms / recipe.time_ms))
```

## 4. Analytical Yield

Bkz. [activity-engine/crafting.md §5](../activity-engine/crafting.md):
```
successCount = sampleBinomial(iterations_processed, craft_success_pct/100.0, seed.sub("success"))
failCount = iterations_processed - successCount

outputs = per-recipe.Outputs; success ile qty × chance sample
inputs consumed = success × input_qty + fail × input_qty × consume_on_fail_pct
```

Per-iteration loop **YOK**. Tek binomial sample; determinism seed'li.

## 5. Threshold Detection

Threshold:
- **Input exhaustion**: envanter material tükendi → paused.
- **Output inventory full**: yeni item slot yok → paused.
- **Iterations complete**: batch tamamlandı → auto-complete.
- **Character death** (nadir crafting sırasında; safe zone genellikle) → paused.

```go
iterations_until_input_exhaust = min(available_input[X] / recipe.Inputs[X].Qty for X)
iterations_until_output_inv_full = character.InventoryFree / avg_output_slots_per_iteration
threshold_iter = min(iterations, iterations_until_exhaust, iterations_until_full)
threshold_ms = threshold_iter × recipe.time_ms
```

En küçük threshold kazanır; `elapsed_effective = threshold_ms`.

## 6. Start Endpoint

**POST `/v1/characters/{id}/craft`**
```json
{
  "recipe_code": "steel_longsword",
  "iterations": 20,
  "station_location_id": "loc-ashfield-forge"
}
```

Server:
1. Load character.
2. Validate:
   - Recipe known.
   - Profession level >= recipe.tier_required.
   - Station present at location (if required).
   - Envanterde material (iterations × input.qty).
3. `Activity.ValidateStart` + snapshot.
4. Create contract; return 201.

## 7. Preview Endpoint

**POST `/v1/characters/{id}/craft/preview`**
```json
{ "recipe_code": "steel_longsword", "iterations": 20 }
```
Response:
```json
{
  "expected_success": 18,
  "expected_fail": 2,
  "expected_outputs": [
    { "item_code": "steel_longsword_template", "expected_qty": 18 },
    { "item_code": "metal_scrap", "expected_qty": 4 }
  ],
  "expected_inputs_consumed": [
    { "item_code": "steel_ingot", "expected_qty": 57 }
  ],
  "expected_duration_ms": 160000,
  "expected_xp": 720
}
```

Kullanıcı UI'da preview görür, iterations ayarlar.

## 8. Cancel Endpoint

**POST `/v1/characters/{id}/craft/{contract_id}/cancel`**
- Partial iteration settle (`activity-engine/crafting.md §10`).
- Mid-iteration'da olan iteration abort; material iade %50.

## 9. Iteration Chunking (Phase 02 Optimization)

Reveal içinde bulk item_instance insert perf:
- 100 iteration = 100 item_instance (equipment ise).
- Bulk INSERT ile 10× hızlanma.

v1'de tek-tek insert loop (küçük N için OK).

## 10. Cross-Recipe Queue (Phase 02)

Aynı anda birden fazla recipe queue:
- Recipe A 20 iter → Recipe B 10 iter → Recipe C 5 iter.
- Sequential execution.
- v1'de tek recipe / contract.
- Phase 02: multi-recipe queue in single contract.

## 11. Fail Modları

- `craft.recipe_unknown`.
- `craft.station_missing`.
- `craft.insufficient_materials`.
- `craft.inventory_full` (start'ta; runtime threshold pause).
- `craft.profession_level_insufficient`.

## 12. Concurrency

- Primary channel contract unique.
- Auto-close previous.

## 13. Reveal Örnek

```
Recipe: steel_longsword, 20 iter, 8s each
Elapsed: 3 min = 180 sec = ~22 iter → capped at 20.

Craft success 90%:
  successCount binomial sample = 19, fail = 1
Outputs:
  steel_longsword_template × 19 (each item_instance yaratılır, seed'li roll)
  metal_scrap 0.2 × 19 = 3.8 → binomial = 4
Inputs consumed:
  steel_ingot: 19*3 + 1*3/2 = 57 + 1.5 = 58 (rounded)
  leather_grip: 19 + 0.5 = 19
  coal: 19*2 + 1 = 39
XP:
  character: 40 × 19 + 20 × 1 = 780
  blacksmithing profession: 200 × 19 + 100 × 1 = 3900
Iterations complete → contract auto-complete.
```

## 14. Metrics
- `crafting_queue_iterations_total{recipe}`.
- `crafting_queue_success_rate{recipe}` histogram.
- `crafting_queue_paused_reason_total{reason}`.

## 15. Cross-Reference
- [Recipes](recipes.md), [Materials](materials.md), [Quality](quality.md), [Critical Craft](critical-craft.md), [Crafting XP](crafting-xp.md)
- [activity-engine/crafting.md](../activity-engine/crafting.md), [activity-engine/lifecycle.md](../activity-engine/lifecycle.md)
- [item-engine/item-instance.md](../item-engine/item-instance.md)

## 16. Bekleyen Kararlar
- Iteration cap final (200 v1).
- Multi-recipe queue Phase 02.
- Auto-withdraw materials from bank Phase 02.
- Bulk INSERT optimization Phase 02.
