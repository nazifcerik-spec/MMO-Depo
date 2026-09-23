# Item Stats

## 1. Amaç
Item'ın karaktere kattığı **base stats** (template'ten roll edilen) + affix + socket + enchant kontribüsyonlarını modifier stack'e nasıl yolladığını tanımlar.

## 2. Katmanlar

Item'ın stat kaynakları:
```
1. Rolled Base Stats     (item_template.base_stats min-max range → rolled)
2. Prefix Modifier       (magnitude × prefix effect)
3. Suffix Modifier       (magnitude × suffix effect)
4. Socket + Gem          (gems installed → gem effects)
5. Enchant               (enchant.effect)
6. Set Bonus             (aktif set piece sayısı gereksinimi)
7. Unique Effect         (legendary+ special)
8. Durability Penalty    (durability 0'da passive/proc devre dışı)
```

## 3. Rolled Base Stats

Template'te range:
```yaml
base_stats:
  weapon_damage:
    min: 45
    max: 60
  attack_speed_ms: 2400  # sabit; roll edilmez
  armor:
    min: 30
    max: 45
```

Roll:
```go
rolled_stats := map[string]int32{}
for stat, spec := range template.BaseStats {
    if spec.Min == spec.Max {
        rolled_stats[stat] = spec.Min
    } else {
        rolled_stats[stat] = seed.RangeInt(spec.Min, spec.Max)
    }
}
rolled_stats[stat] *= rarityMultiplier[rarity]
```

Rarity multiplier: [rarity.md §6](rarity.md).

## 4. Persistence

`content.item_instance.rolled_stats` JSONB:
```json
{
  "weapon_damage": 52,
  "attack_speed_ms": 2400,
  "armor": null    // silah için armor yok
}
```

## 5. Modifier Stack Contribution

Item Engine, item'ı Modifier'lara çevirir (bkz. [skill-engine/modifiers.md](../skill-engine/modifiers.md)):

```go
func (e *ItemEngine) ModifiersFor(instance Instance) []engine.Modifier {
    mods := []engine.Modifier{}

    // 1. Rolled base stats → flat modifier per stat
    for stat, val := range instance.RolledStats {
        mods = append(mods, engine.Modifier{
            Kind:       engine.KindFlat,
            Target:     stat,   // "weapon_damage", "armor", "str", ...
            Amount:     float64(val),
            SourceKind: "item_base",
            SourceID:   instance.ID,
        })
    }

    // 2. Prefix
    if instance.Affixes.Prefix != nil {
        affix := e.affixReg.ByCode(instance.Affixes.Prefix.Code)
        for _, m := range affix.Effect.Modifiers {
            copy := m
            copy.Amount = float64(instance.Affixes.Prefix.Magnitude)
            copy.SourceKind = "item_prefix"
            copy.SourceID = instance.ID
            mods = append(mods, copy)
        }
    }

    // 3. Suffix — same pattern

    // 4. Sockets + gems
    for _, socket := range instance.Sockets {
        if socket.GemInstanceID == nil { continue }
        gem := e.gemReg.Load(*socket.GemInstanceID)
        for _, m := range gem.Effect.Modifiers {
            copy := m
            copy.SourceKind = "item_gem"
            copy.SourceID = *socket.GemInstanceID
            mods = append(mods, copy)
        }
    }

    // 5. Enchant
    if instance.EnchantID != nil {
        enc := e.enchantReg.Load(*instance.EnchantID)
        for _, m := range enc.Effect.Modifiers {
            copy := m
            copy.SourceKind = "item_enchant"
            copy.SourceID = instance.ID
            mods = append(mods, copy)
        }
    }

    // 6. Unique effect trigger'ları burada değil; PassiveEngine dispatch (bkz. unique-effect.md)

    // 7. Durability check
    if instance.DurabilityCurrent != nil && *instance.DurabilityCurrent <= 0 {
        // filter out proc/passive modifiers; base stat kalır (v1 tercih)
        mods = filterOutProcOnly(mods)
    }

    return mods
}
```

## 6. Set Bonus Ayrı

Set bonus item'ın kendi ModifiersFor'undan gelmez; ekipmanın **toplu** evaluate'inden gelir:
```go
setBonusMods := e.SetBonusesFor(equippedInstances)
```

Detay: [set.md](set.md).

## 7. Cache Discipline

Item değişimi (equip/unequip, gem install, enchant, durability) → stat cache invalidate (bkz. [skill-engine/stats.md §6](../skill-engine/stats.md)).

## 8. Endpoint

Item detail response (bkz. [item-instance.md §9](item-instance.md)) `rolled_stats` + `effective_stats_preview` (opsiyonel).

Effective stats preview:
```json
{
  "rolled_stats": { "weapon_damage": 52 },
  "affix_contribution": {
    "prefix": { "str": 5 },
    "suffix": { "sta": 12 }
  },
  "socket_contribution": { "str": 3 },
  "enchant_contribution": { "attack_power": 12 },
  "total_stat_contribution": {
    "weapon_damage": 52,
    "str": 8,
    "sta": 12,
    "attack_power": 12
  }
}
```

## 9. Stat Compare (Client)

Tooltip'te "equipped vs candidate" karşılaştırma:
- Client mevcut equipment stat contribution ile yeni item contribution farkı → delta rendering.
- Server-side helper (opsiyonel Phase 02): `/v1/characters/{id}/item/{instance_id}/compare`.

## 10. Fail Modları
- Rolled stat key template'te yok → skip + warn.
- Enchant expired/removed → filter out.
- Gem instance void → filter out; socket boş render.

## 11. Test Stratejisi
- Rolled stat deterministic.
- Full modifier list from complex item (prefix+suffix+2 gems+enchant).
- Durability 0 → proc filter.
- Cache invalidation on equip.

## 12. Metrics
- `item_stat_contribution_count` gauge.
- `item_effective_stat_recompute_total`.

## 13. Cross-References
- [Affix](affix.md), [Prefix](prefix.md), [Suffix](suffix.md), [Sockets](sockets.md), [Gems](gems.md), [Enchantments](enchantments.md), [Set](set.md), [Unique Effect](unique-effect.md)
- [skill-engine/modifiers.md](../skill-engine/modifiers.md), [skill-engine/stats.md](../skill-engine/stats.md)
- [character/stats.md](../character/stats.md), [character/attributes.md](../character/attributes.md)

## 14. Bekleyen Kararlar
- Rolled stat min/max range balance final.
- Durability 0'da hangi modifier'lar filtered (proc-only vs all) — v1: proc-only.
- Stat range display client'ta (min-max göstermek mi rolled göstermek mi).
