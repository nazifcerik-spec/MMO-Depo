# Prefix

## 1. Amaç
Item isminin **başına** eklenen affix. Genelde primary/derived stat bonusu; item kimliğinin ön yüzünü verir.

Genel sistem: [affix.md](affix.md).

## 2. Semantic Rol
- **Combat prefix** (Steel, Sharp, Cruel) — physical damage / crit / attack power.
- **Magic prefix** (Wizard's, Enchanted, Blazing) — spell power / element damage.
- **Defensive prefix** (Warrior's, Bulwark, Steadfast) — armor / HP / block.
- **Utility prefix** (Swift, Nimble, Lucky) — movement / evasion / gold find.

## 3. Naming Convention

Item name composition: **`<prefix> <template.display_name> <suffix>`**.

Prefix türü grammatik:
- **Adjective**: "Steel Longsword", "Blazing Staff".
- **Possessive noun**: "Wizard's Robe", "Hunter's Bow".
- **Rare adjective**: "Cruel Longsword" (rare).

Karışım pool'da; loader lokalizasyon disciplini için i18n key ile format:
```
en: "{prefix} {name}"
tr: "{prefix} {name}"  # bazı diller reversed olabilir; format YAML'da
```

## 4. Örnek Prefix Kataloğu (Placeholder)

| Kod | Ad (en) | Effect | Tag | Rarity band |
|---|---|---|---|---|
| `steel` | Steel | +N STR | physical, stat:str | common+ |
| `swift` | Swift | +N% movement + N DEX | movement, stat:dex | common+ |
| `wizards` | Wizard's | +N INT | magic, stat:int | uncommon+ |
| `cruel` | Cruel | +N% crit dmg | precision | rare+ |
| `blazing` | Blazing | +N fire dmg | fire, elemental | uncommon+ |
| `bulwark` | Bulwark | +N armor + N% block | defense | uncommon+ |
| `hunters` | Hunter's | +N% dmg vs beasts | family_bonus | rare+ |
| `lucky` | Lucky | +N% gold find | utility, magic_find | uncommon+ |
| `runic` | Runic | +1 socket | socket | epic+ |
| `ancient` | Ancient | +N% XP gain | utility | rare+ |
| `demonic` | Demonic | +N shadow dmg -N holy resist | shadow, dual | epic+ |
| `radiant` | Radiant | +N holy dmg | holy | rare+ |

## 5. Roll & Magnitude

Roll akışı [affix.md §4](affix.md).

Magnitude range examples:
```
steel:    common 1-3, uncommon 3-6, rare 5-10, epic 8-15, legendary 12-20
swift:    common 2-4, uncommon 4-7, rare 6-10, epic 8-14, legendary 12-18
cruel:    rare only  10-15%, epic 15-22%, legendary 20-30%
```

## 6. Tag Structure

Prefix tags (loader validated):
- **Stat**: `stat:str`, `stat:dex`, `stat:int`, `stat:wis`, `stat:sta`, `stat:cha`.
- **Damage type**: `physical`, `fire`, `ice`, `nature`, `holy`, `shadow`, `arcane`.
- **Class hint**: `warrior_friendly`, `caster_friendly`, `rogue_friendly`.
- **Category**: `defense`, `offense`, `utility`, `precision`, `movement`.
- **Family bonus**: `family_bonus:beast`, `family_bonus:humanoid`, ...

Mutual exclusion (bkz. [affix.md §6](affix.md)) tag'ler üzerinden.

## 7. Class Bias

Bazı prefix'ler class-friendly:
- Wizard's → Mage/Cleric/Druid.
- Steel → Warrior/Ranger/Rogue.

Bias enforcement:
- **Soft**: pool weight class-friendly item template'te daha yüksek.
- **Hard**: pool entry `class_codes` filter (bazı prefix'ler sadece belirli class'a takılan item'a).

v1 tercih: **soft bias** (pool tasarımı ile).

## 8. Display Composition

Client render (i18n resolver):
```
prefix.name_key = "affix.prefix.steel"
suffix.name_key = "affix.suffix.of_the_bear"
template.name_key = "item.steel_longsword.name"

format = "{prefix} {name} {suffix}"

output_en = "Steel Steel Longsword of the Bear"  // duplicate word — content team dikkatli
output_tr = "Steel Steel Uzun Kılıç Ayı'nın"
```

Content team tarafından "Steel Longsword" template ismi zaten sabit "Steel" içerdiği için prefix "Steel" seçmek redundant görünür — content pool tasarımında item name karşılığı prefix pool'undan `steel` çıkarılır (opsiyonel filter).

Loader warn: duplicate word detection.

## 9. Endpoint

Public content:
**GET `/v1/content/affixes/prefix/{code}`**

Response:
```json
{
  "code": "steel",
  "display_name_key": "affix.prefix.steel",
  "effect_summary_key": "affix.prefix.steel.summary",
  "magnitude_range": { "min": 3, "max": 8 }  // rarity bağlı; endpoint filter
}
```

## 10. Test Stratejisi
- Roll seed → expected prefix code + magnitude.
- Class bias distribution.
- Tag exclusion — suffix çakışma.
- Duplicate word content lint.

## 11. Metrics
- `affix_prefix_rolled_total{code,rarity}`.
- `affix_prefix_magnitude_mean{code}`.

## 12. Cross-References
- [Affix](affix.md), [Suffix](suffix.md), [Item Instance](item-instance.md), [Rarity](rarity.md)
- [database/content-templates.md §11](../database/content-templates.md)

## 13. Bekleyen Kararlar
- Prefix rework worker Phase 02 (patch sonrası mevcut item reroll).
- Custom prefix (crafted-choose) Phase 02.
- i18n format grammar (bazı diller prefix noun formu farklı).
