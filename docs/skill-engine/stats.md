# Stats (Compute Pipeline)

## 1. Amaç
Karakterin efektif stat setini deterministik olarak hesaplamak; cache; invalidation.

Yüksek seviye formul: [character/stats.md](../character/stats.md), [character/attributes.md](../character/attributes.md). Bu doküman **engine implementation** detayı.

## 2. Katmanlar (Layer Model)

Hesap her stat için katman zinciri:
```
Layer 0: base_row               (character.stats.str/dex/...)
Layer 1: race_modifier          (content.race.stat_modifiers)
Layer 2: class_bias             (content.class.base_stat_bias)
Layer 3: spec_bias              (content.specialization.stat_bias, tier1 + tier2)
Layer 4: ascension_bonus        (sum of ascension_choice effects)
Layer 5: gear_stat_sum          (equipped item_instance.rolled_stats + affix + enchant)
Layer 6: set_bonus              (item_set active bonuses)
Layer 7: mastery_bonus          (mastery milestone effects)
Layer 8: buff_add / buff_mult   (active status_effect)
Layer 9: debuff_add / debuff_mult (negative modifiers)
Layer 10: cap / clamp
```

## 3. Primary Compute
```go
func computePrimary(inputs Inputs) PrimaryStats {
    result := PrimaryStats{}
    for _, stat := range primaryStatKeys {
        v := inputs.BaseRow[stat]                        // layer 0
        v += inputs.Race.StatModifiers[stat]             // 1
        v += inputs.Class.BaseStatBias[stat]             // 2
        v += inputs.SpecTier1.StatBias[stat]             // 3a
        v += inputs.SpecTier2.StatBias[stat]             // 3b
        v += inputs.AscensionBonus[stat]                 // 4
        v += inputs.GearStatSum[stat]                    // 5
        v += inputs.SetBonusSum[stat]                    // 6
        v += inputs.MasteryBonus[stat]                   // 7

        for _, b := range inputs.BuffAdditives[stat] { v += b }
        for _, d := range inputs.DebuffAdditives[stat] { v -= d }
        for _, m := range inputs.BuffMultipliers[stat] { v = int32(float64(v) * m) }

        if v < 0 { v = 0 }                                // floor
        result[stat] = v
    }
    return result
}
```

## 4. Derived Compute
Primary hesaplandıktan sonra:
```go
derived := computeDerived(primary, gear, class, spec, race, ascension, buffs)
```

Formüller [character/stats.md §3](../character/stats.md) tanımlı. Örnek:
```go
// Attack Power
ap := 2*primary.STR + 1*primary.DEX + gear.WeaponAP + ap_flat_buffs
ap = int32(float64(ap) * ap_mult_stack.Total())
result.AttackPower = ap

// Max HP
maxHP := 50 + 10*primary.STA + 5*level + gear.HPBonus
result.MaxHP = maxHP
```

## 5. Cap & Clamp

Hard caps derived'de uygulanır:
```
crit_chance_pct        max 75, diminishing after 50
crit_dmg_mult          max 2.5
haste_pct              max 40
cdr_pct                max 30
armor_dr_pct           max 75
xp_gain_pct_total_mult max 3.0
```

Diminishing returns:
```go
// crit chance eğrisi (%50 sonrası azalan)
func softCapCrit(raw float64) float64 {
    if raw <= 50 { return raw }
    excess := raw - 50
    diminished := 50 + excess * 0.5
    if diminished > 75 { diminished = 75 }
    return diminished
}
```

## 6. Cache

Kaynak: `character.stats.derived_cache_json`.

Cache invalidation triggers (bkz. [character/attributes.md §8](../character/attributes.md)):
- Level up
- Allocate stats
- Respec (spec veya stat)
- Equip / unequip
- Buff apply / expire
- Ascension spend
- Mastery milestone
- Content_version bump

Signal:
```sql
UPDATE character.stats
SET derived_cache_invalidated_at = now(), version = version + 1
WHERE character_id = $1 AND version = $expected;
```

Read path:
```go
func (e *StatEngine) ComputeDerived(ctx context.Context, charID uuid.UUID, opts ...) (Derived, error) {
    stats, _ := e.repo.LoadCharacterStats(ctx, charID)
    if stats.CacheInvalidatedAt == nil && stats.DerivedCache != nil {
        return stats.DerivedCache, nil  // hit
    }
    // recompute
    inputs := e.gatherInputs(ctx, charID)
    derived := computePrimary(inputs), computeDerived(...)
    _ = e.repo.PersistDerivedCache(ctx, charID, derived)
    return derived, nil
}
```

## 7. Context

Bazı stat'lar context'e bağlı (scope):
- `gather_yield_pct` — sadece profession aktivitesinde.
- `damage_pct.fire` — sadece fire element damage'a.

`StatContext` parametresi:
```go
type StatContext struct {
    Purpose ContextPurpose  // combat, activity, vendor, ...
    Activity string         // "mining", "hunting", ...
    Element string          // "fire", "holy", ...
    TargetKind string       // "mob", "player", "boss"
}
```

ComputeDerived context-aware; scope-limited modifier'lar sadece uygun context'te uygulanır.

## 8. In-Combat vs Out-of-Combat

Regen combat state'e göre farklı:
```go
if inCombat {
    hpRegen *= 0.2
    manaRegen *= 0.5
}
```

Combat state Redis flag (`char:combat:<char>` TTL 5s; her damage taken/verilme yenilenir).

## 9. Cross-Character Compute
Bir defada tek character. Party bonuslar (rally, aura) her karakter için ayrı compute'da apply.

Party bonus ayrı katman — buff olarak modelleniyor (Layer 8).

## 10. Fail Modları
- Missing race/class → 500; audit critical.
- Gear referans invalid → skip + warn.
- Buff expired ancak DB'de → filter out.
- Overflow (int32 taşma) → clamp max int32 + warn.

## 11. Test Stratejisi
- Table-driven: (inputs) → expected derived.
- Cache hit/miss.
- Invalidation cascade.
- Context-scoped modifier isolation.
- Layer 8 buff multi-stack.

## 12. Metrics
- `stat_engine_compute_total{source}` (hit/miss).
- `stat_engine_compute_duration_seconds`.
- `stat_cache_invalidations_total{reason}`.
- `stat_engine_cache_age_seconds`.

## 13. Bekleyen Kararlar
- Cache TTL (opsiyonel; şu an invalidation-based).
- Precomputed variants per scope (combat vs profession) — optimization Phase 02.
- Cross-character sim engine (party avg) — Phase 02.
