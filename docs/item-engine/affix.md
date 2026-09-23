# Affix (System Architecture)

## 1. Amaç
Prefix + Suffix'in ortak sistemi: pool tanımı, roll mekaniği, magnitude, weight, mutually exclusive tag'ler.

Prefix detayı: [prefix.md](prefix.md). Suffix detayı: [suffix.md](suffix.md).

## 2. Model

Kaynak: [database/content-templates.md §11](../database/content-templates.md).

```
content.affix_pool
  id, code (unique), kind (prefix|suffix)
  entries JSONB: [{affix_id, weight, magnitude_min, magnitude_max, tag[]}]

content.affix
  id, code (unique), kind (prefix|suffix)
  display_name_key
  effect_json   (modifier list)
  tags (mutually exclusive için)
```

Affix pool → item template'e attach edilir (`content.item_template.affix_pool_ids`).

## 3. Content Pack Format

Affix:
```yaml
# content/affixes/steel_prefix.yaml
code: steel
kind: prefix
display_name_key: affix.prefix.steel
effect:
  _v: 1
  modifiers:
    - kind: flat
      target: str
      scope: always
tags: [physical, warrior_friendly]
```

Pool:
```yaml
# content/affix_pools/uncommon_weapon_prefix_pool_v1.yaml
code: uncommon_weapon_prefix_pool_v1
kind: prefix
entries:
  - affix_code: steel
    weight: 30
    magnitude_min: 3
    magnitude_max: 8
    tags: [physical]
  - affix_code: swift
    weight: 25
    magnitude_min: 5
    magnitude_max: 10
    tags: [physical, ranged_bias]
  - affix_code: cunning
    weight: 20
    magnitude_min: 2
    magnitude_max: 6
    tags: [precision]
  # ... daha fazla
```

## 4. Roll Algoritması

Prefix + Suffix ayrı ancak aynı algoritma:
```go
func rollAffix(pool AffixPool, item TemplateDef, rarity Rarity, seed Seed) (AffixInstance, error) {
    // 1. Filter eligible entries by rarity, tags (mutual exclusion)
    eligible := filterEligible(pool.Entries, item, rarity)
    // 2. Weighted pick
    picked := seed.WeightedPick(eligible, func(e Entry) int { return e.Weight })
    // 3. Roll magnitude
    mag := seed.RangeInt(picked.MagnitudeMin, picked.MagnitudeMax)
    // 4. Return
    return AffixInstance{
        AffixID:   picked.AffixID,
        Code:      picked.AffixCode,
        Magnitude: mag,
    }, nil
}
```

## 5. Magnitude & Rarity Multiplier

Magnitude rarity multiplier ile scale:
```
final_magnitude = base_magnitude × rarity_multiplier
```
(bkz. [rarity.md §6](rarity.md))

Örnek: `steel` (base 3-8) rare item → 3.3-8.8 → integer round.

## 6. Mutually Exclusive Tags

Bazı affix'ler birlikte olamaz (aynı stat çakışması):
- Prefix "Fierce" (STR+X) + Suffix "of Fury" (STR+Y) — çakışma; birinden vazgeç.

Tag-based exclusion:
- Pool entry `tags: ["stat:str"]`.
- Roll sırasında: eğer prefix "stat:str" içeriyorsa, suffix'te "stat:str" tag'li affix'ler filtered out.

Loader validation: aynı item'a hem prefix hem suffix ile aynı tag çift roll yapılamamalı — pool tasarımı gereği.

## 7. Affix Instance Serialization

`content.item_instance.affixes` JSONB:
```json
{
  "prefix": { "code": "steel", "magnitude": 5 },
  "suffix": { "code": "of_the_bear", "magnitude": 12 }
}
```

Renderer client bunu name composition için kullanır (bkz. [prefix.md](prefix.md), [suffix.md](suffix.md)).

## 8. Effect Application

Affix'in `effect_json.modifiers` listesi:
```json
{
  "modifiers": [
    { "kind": "flat", "target": "str", "amount": null }
  ]
}
```

Amount `null` çünkü magnitude roll'dan geliyor. Item instance modifier gather'ında amount = magnitude.

```go
func (e *ItemEngine) ModifiersFor(instance Instance) []engine.Modifier {
    mods := []engine.Modifier{}
    // Base stats from rolled_stats
    for stat, value := range instance.RolledStats {
        mods = append(mods, flatModifier(stat, value, sourceItem(instance)))
    }
    // Prefix
    if instance.Affixes.Prefix != nil {
        affix := e.affixRegistry.ByCode(instance.Affixes.Prefix.Code)
        for _, m := range affix.Effect.Modifiers {
            m.Amount = float64(instance.Affixes.Prefix.Magnitude)
            m.SourceKind = "item_prefix"
            m.SourceID = instance.ID
            mods = append(mods, m)
        }
    }
    // Suffix — same
    // Gems, enchant — ayrı bölümler
    return mods
}
```

## 9. Reroll Consumable (Phase 02)

Reroll scroll / vendor:
```
POST /v1/items/{instance_id}/reroll-affixes
  body: { "reroll_scroll_instance_id": "..." }
```

Server:
- Consume scroll.
- New seed = hash(instance.seed, reroll_count).
- Re-run affix roll on same pool.
- Update item_instance.affixes.
- Cache invalidate.

Cost/frequency Phase 02 balance.

## 10. Affix Pool per Rarity + Slot

Aynı item template farklı rarity + farklı slot için farklı pool alabilir:
```yaml
# item_template
affix_pool_refs:
  common:
    prefix: common_weapon_prefix
    suffix: common_weapon_suffix
  uncommon:
    prefix: uncommon_weapon_prefix
    suffix: uncommon_weapon_suffix
  rare:
    prefix: rare_weapon_prefix
    suffix: rare_weapon_suffix
  # ...
```

Loader validation zorlar; opsiyonel simplification: rarity-agnostic single pool.

## 11. Content Loader Validation

- Pool entry weight > 0.
- magnitude_min <= magnitude_max.
- affix_code var (affix satırı var).
- Kind consistency (pool.kind == affix.kind).
- Total weight > 0.

## 12. Testing
- Table-driven roll: seed → expected (affix_code, magnitude).
- Distribution over 10k rolls matches weights within ±%2.
- Tag exclusion works.
- Reroll preserves other fields (rolled_stats, sockets).

## 13. Metrics
- `affix_rolled_total{pool,affix_code}`.
- `affix_magnitude_avg{affix_code}`.
- `affix_reroll_total`.

## 14. Cross-References
- [Prefix](prefix.md), [Suffix](suffix.md), [Item Instance](item-instance.md), [Rarity](rarity.md), [Stats](stats.md)
- [database/content-templates.md §11](../database/content-templates.md)
- [skill-engine/modifiers.md](../skill-engine/modifiers.md)

## 15. Bekleyen Kararlar
- Rarity-per-slot pool separation zorunlu mı, opsiyonel mı.
- Reroll consumable v1 vs Phase 02.
- Advanced tier affix (crafted-only, exalted) Phase 02.
- Affix crafting (choose specific affix with material) Phase 02.
