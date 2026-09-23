# Passive (Racial Passive Ability)

## 1. Amaç
Race'in **her zaman aktif** mekanik bonusu. Combat, gathering, XP veya konfor katmanında küçük ama sürekli.

## 2. Model
Race'in `passive_ability_id` FK → `content.skill_template`.
`skill_template.type = 'racial'`, `skill_template.effects_json` içinde passive effect list.

Skill template pattern (bkz. [database/content-templates.md §4](../database/content-templates.md)):
- `type='racial'`.
- `race_id` FK back-reference.
- `resource_cost=0`.
- `cooldown_ms=0`.
- `effects_json` = passive modifier list.

## 3. Effect DSL

Effect JSON şeması (versionlu):
```json
{
  "_v": 1,
  "effects": [
    {
      "kind": "flat_bonus",
      "target": "xp_gain_pct",
      "amount": 5,
      "scope": "always"
    },
    {
      "kind": "flat_bonus",
      "target": "gather_yield_pct",
      "amount": 3,
      "scope": "profession:mining"
    }
  ]
}
```

### 3.1 Effect Kinds
- `flat_bonus` — sabit değer ekle (percent veya integer).
- `multiplier` — çarpımsal (%110 = ×1.1).
- `resource_regen` — HP/mana/resource regen boost.
- `damage_modifier` — belirli damage type'ına bonus (yalnız combat).
- `resistance` — belirli hasar türüne direnç.
- `movement` — hız modifier.
- `vision` — night vision, tile radius etc.
- `custom` — special mechanic (elf `forest_stride`, dwarf `dark_vision`).

### 3.2 Effect Target'ları (Whitelisted)
Whitelist enforced by loader; kod dışında bir target eklenemez:
```
xp_gain_pct
gather_yield_pct
gather_yield_pct.{profession_code}
craft_success_pct
crit_chance_pct
crit_dmg_mult
hp_regen_flat
mana_regen_flat
armor_flat
magic_resist_flat
movement_speed_pct
inventory_slot_bonus
carry_capacity_flat
```

### 3.3 Scope
- `always` — sürekli aktif.
- `combat` — sadece combat sırasında.
- `out_of_combat` — combat dışında.
- `profession:<code>` — belirli meslek aktivitesinde.
- `activity:<kind>` — belirli aktivite türünde (hunting, adventure).

## 4. Passive Uygulama

Effect combiner (bkz. [game-design/skills.md](../game-design/skills.md)) recompute sırasında:
```
for each learned passive effect:
  if scope matches current context:
    apply to derived cache
```

Karakter `character.stats.derived_cache_json` içinde bu değerlerin toplu sonucu.

Passive her zaman aktif olduğundan `derived_cache_invalidated_at` set edilir; recompute bir sonraki okumada.

## 5. v1 Race Passive'leri (Placeholder)

| Race | Passive Kod | Effect |
|---|---|---|
| `human` | `human_versatility` | `xp_gain_pct +5, always` |
| `dwarf` | `dwarven_endurance` | `gather_yield_pct +5, profession:mining` + `hp_regen_flat +1, always` |
| `elf` | `elven_grace` | `gather_yield_pct +5, profession:woodcutting` + `movement_speed_pct +2, always` |
| `orc` | `orcish_bloodlust` | `crit_dmg_mult +0.05, combat` (5% ekstra crit dmg) |
| `merfolk` | `tide_born` | `gather_yield_pct +5, profession:fishing` + `resistance.aquatic +10, always` |
| `arcanborn` | `arcane_attunement` | `mana_regen_flat +2, always` + `spell_power_pct +3, combat` |

Detay YAML örneği (human):
```yaml
# content/skills/human_versatility.yaml
code: human_versatility
type: racial
race_code: human
display_name_key: race.human.passive.name
description_key: race.human.passive.desc
effects:
  _v: 1
  effects:
    - kind: flat_bonus
      target: xp_gain_pct
      amount: 5
      scope: always
```

## 6. Server-Side Compute

Character creation:
```go
// grant racial passive
racialPassive := content.Races.ByID(character.RaceID).PassiveAbilityID
_ = characterSkill.Learn(ctx, character.ID, racialPassive, rank=1)
// invalidate cache
_ = characterStats.InvalidateCache(ctx, character.ID)
```

Effective stat compute:
```go
func computeDerivedStats(char Character, learned []Skill, ...) Derived {
    // ... base compute ...
    for _, sk := range learned {
        if sk.Type != "racial" && sk.Type != "passive" { continue }
        for _, eff := range sk.Effects {
            if !scopeMatches(eff.Scope, currentContext) { continue }
            applyEffect(&derived, eff)
        }
    }
    return derived
}
```

## 7. Content Loader Validation

- `effects._v` version present.
- Each effect `kind` in whitelist.
- Each effect `target` in whitelist.
- Scope well-formed.
- Amount tip check (int vs float).
- Toplam magnitude soft cap: bir race'in tüm passive effect'lerinin toplam etki bütçesi:
  ```
  |amount_pct_sum| <= 20  (net passive % impact)
  ```

## 8. Test Stratejisi
- Effect DSL parser table-driven.
- Scope matcher (mining aktivitesi vs woodcutting).
- Passive apply order (deterministic: sorted by effect id).
- Nerf simülasyonu: passive değişimi → derived cache invalidate → new compute.

## 9. Balance
- Passive'lerin gerçek etkisi %2-%8 aralığı; kimlik verir ama meta-defining değil.
- Class synergy: bazı race passive'leri belirli class'larla iyi çalışır ancak "sadece bu race + class" tuzağı olmamalı.

## 10. Metrics
- `race_passive_active_count` — kaç karakter bu passive'i taşıyor.
- `race_passive_effect_impact_pct` — measured meta-impact per race (weekly report).

## 11. Cross-Platform
Server compute; client sadece display. UI'da passive tooltip'i i18n resolver ile.

## 12. Fail Modları
- Passive effect data corrupt → loader fail; deploy blok.
- Runtime effect unknown → skip + warn log.

## 13. Content Update Flow
1. YAML patch.
2. CI validate.
3. Staging deploy → race pack apply.
4. Impact simulation report (weekly).
5. Production deploy → content_version bump.
6. Existing characters: cache invalidate broadcast.

## 14. İlişkiler
- [Race](race.md), [Racial Skills](racial-skills.md), [Racial Bonuses](racial-bonuses.md), [character/stats.md](../character/stats.md), [game-design/skills.md](../game-design/skills.md), [database/content-templates.md](../database/content-templates.md)

## 15. Bekleyen Kararlar
- Passive stack rules: iki passive aynı target → additive vs multiplicative (öneri: additive).
- Level scaling opsiyonu (passive amount level ile ↑) — game-design/races.md §9 belirtilmiş; final Phase 02.
- Passive "toggle" opsiyonu (bazı out-of-combat passive'leri manuel kapatma) — Phase 03.
