# Suffix

## 1. Amaç
Item isminin **sonuna** eklenen affix. Genelde utility, secondary stat, situational bonus. Item kimliğinin arka yüzü.

Genel sistem: [affix.md](affix.md).

## 2. Semantic Rol
- **Attribute suffix** (of the Bear, of the Falcon, of the Wolf) — STA/DEX/AP tarzı.
- **Utility suffix** (of Endurance, of Speed, of Wisdom) — regen, movement, XP gain.
- **Resistance suffix** (of Warding, of Shadow-Break) — elemental resist.
- **Meta suffix** (of Discovery, of Fortune) — gold/exp/rarity.

## 3. Naming Convention

Format: **`<template.display_name> <suffix>`**.

Suffix grammatik:
- **Prepositional phrase**: "of the Bear", "of the Falcon", "of Fury".
- **Complex**: "of Ancient Wisdom", "of the Shattered Blade".
- Language-specific reversal: bazı dillerde suffix aslında prefix olabilir; i18n format ayrı.

## 4. Örnek Suffix Kataloğu (Placeholder)

| Kod | Ad (en) | Effect | Tag | Rarity band |
|---|---|---|---|---|
| `of_the_bear` | of the Bear | +N STA | stat:sta | common+ |
| `of_the_falcon` | of the Falcon | +N DEX | stat:dex | common+ |
| `of_the_wolf` | of the Wolf | +N% attack speed | speed | uncommon+ |
| `of_wisdom` | of Wisdom | +N WIS + N mana regen | stat:wis, caster | uncommon+ |
| `of_fury` | of Fury | +N% dmg -N armor | offense, glass_cannon | rare+ |
| `of_warding` | of Warding | +N all resist | defense | uncommon+ |
| `of_endurance` | of Endurance | +N HP regen | vital | common+ |
| `of_the_scholar` | of the Scholar | +N INT + N% XP gain | caster, utility | rare+ |
| `of_fortune` | of Fortune | +N% gold find | utility, magic_find | uncommon+ |
| `of_the_shadow` | of the Shadow | +N% evasion | evasion, dark | rare+ |
| `of_the_dragon` | of the Dragon | +N% fire dmg + N HP | dual, dragon | epic+ |
| `of_ancient_might` | of Ancient Might | +N STR + N% dmg | offense, strong | legendary |

## 5. Roll & Magnitude

Roll akışı: [affix.md §4](affix.md).

Magnitude range örnekleri:
```
of_the_bear:      common 3-6, uncommon 6-12, rare 10-18, epic 15-25, legendary 22-35
of_the_wolf:      uncommon 2-5%, rare 4-8%, epic 6-12%, legendary 10-16%
of_the_dragon:    epic 10-15, legendary 15-22
```

## 6. Tag Structure

Prefix ile aynı taxonomy (bkz. [prefix.md §6](prefix.md)). Mutually exclusive check aynı sistem.

## 7. Class Bias

Suffix'ler genelde universal — utility odaklı. Ancak bazı suffix caster-only:
- of Wisdom → caster classes bias.
- of the Falcon → ranger/rogue bias (soft).

## 8. Display Composition

Client render (bkz. [prefix.md §8](prefix.md)):
```
"Steel Longsword of the Bear"
"Sharp Bow of the Wolf"
"Wizard's Robe of Wisdom"
```

Grammar note: suffix "of the X" default İngilizce; Türkçe için i18n format:
```yaml
tr: "{name} — {suffix}"
```
Sample TR: "Uzun Kılıç — Ayı'nın Gücü".

Format YAML'da:
```yaml
# content/i18n/format.tr.yaml
item_name_format: "{prefix_optional}{name}{suffix_optional_dashed}"
```

## 9. Test Stratejisi
- Roll deterministic.
- Distribution.
- Tag exclusion with prefix.
- Grammar/format correctness (spot check).

## 10. Metrics
- `affix_suffix_rolled_total{code,rarity}`.
- `affix_suffix_magnitude_mean{code}`.

## 11. Cross-References
- [Affix](affix.md), [Prefix](prefix.md), [Item Instance](item-instance.md), [Rarity](rarity.md)
- [database/content-templates.md §11](../database/content-templates.md)

## 12. Bekleyen Kararlar
- Compound suffix (2 parça: "of the Bear's Fury") Phase 02 opsiyonel.
- Suffix-only rare item'lar (no prefix) — content design.
- Rework worker Phase 02.
