# Scaling (Stat Coefficients + Level Scaling)

## 1. Amaç
Skill effect'inin base amount + stat scaling + level scaling + gear power formulünü tanımlamak. Combat damage/heal/proc formullerinin merkezi.

## 2. Effect Amount Formula

```
final = base
      + Σ (primary_stat × coefficient)
      × level_scale
      + gear_power_contribution
      × global_modifiers
      × context_modifiers
      × rank_multiplier
      × crit_multiplier   (if crit)
      × variance          (RNG jitter)
      - target_defense
```

Detaylı katman:
```
1. base (skill.effects[i].base)
2. + Σ (primary[stat] × scaling.coef[stat])
3. * level_scale(caster_level)
4. + gear_power (weapon_ap veya spell_power)
5. * (1 + damage_pct.additive_stack)
6. * damage_pct.multiplicative_products
7. * rank_multiplier
8. * (crit ? crit_dmg_mult : 1)
9. * (1 + variance in ±%5)
10. - defense mitigation (armor/resist/DR)
11. clamp min 1 (damage)
```

## 3. Scaling JSON Şeması

Skill definition:
```json
{
  "scaling": {
    "_v": 1,
    "primary": {
      "str": 1.0,
      "dex": 0.0,
      "int_": 0.0
    },
    "attack_power_coef": 1.5,
    "spell_power_coef": 0.0,
    "heal_power_coef": 0.0,
    "weapon_damage_coef": 1.0,
    "level_scale_factor": 1.0,
    "variance_pct": 5
  }
}
```

Coefficient budget:
- `sum(primary.coef) + attack_power_coef + spell_power_coef + heal_power_coef + weapon_damage_coef ≤ 3.0`
- Aşan skill'ler cheat sinyal; loader warn.

## 4. Level Scale

```
level_scale(caster_level) = 1 + (caster_level - 1) * 0.02
```
Level 1 → ×1.0, Level 50 → ×1.98.

Skill effect level ilerledikçe otomatik güçlenir (skill rank olmasa da).

Formul config `config.scaling.yaml`:
```yaml
level_scale:
  per_level_multiplier: 0.02
```

## 5. PvE vs PvP Ayrımı (Phase 02)

PvP damage formulünde ayrı katman:
```
if target.is_player:
    final *= pvp_damage_multiplier  # config 0.5 varsayılan
```

v1'de player vs player yok; formul aktivasyona hazır.

## 6. Damage Type / Element

Element'ler:
- physical, fire, ice, nature, holy, shadow, arcane.

Modifier target'ları:
- `damage_pct.fire` — fire damage'a %bonus.
- `resist.fire` — hedef fire resist.

Hesap:
```
damage_type_bonus = mods.gather("damage_pct." + element)
final *= (1 + damage_type_bonus)

target_resist = target.derived.resistances[element]
mitigated = final * (target_resist / (target_resist + K_RESIST * caster_level))
final -= mitigated
```

## 7. Armor Reduction (Physical)

```
armor_dr = target.armor / (target.armor + K_ARMOR * caster.level)
armor_dr = clamp(armor_dr, 0, 0.75)  # max %75 DR
mitigated = final * armor_dr
final -= mitigated
```

`K_ARMOR = 60` (v1 tune); config'te.

## 8. Magic Resist (Magical)

Aynı pattern, `K_MAGIC = 60`.

## 9. Crit Damage

```
crit_chance = clamp(caster.crit_chance_pct, 0, 75)
if RNG < crit_chance:
    final *= caster.crit_dmg_mult  # base 1.5, gear bonus
```

Crit dmg cap: 2.5×.

## 10. Variance

RNG jitter:
```
variance = 1 + (RNG_stream - 0.5) * (variance_pct / 100)  # ±%5 default
final *= variance
```

Determinism: RNG stream contract seed'inden türetilir.

## 11. Heal Formula

Damage'a benzer ancak defense yok:
```
final_heal = base
           + Σ (primary × coef)
           + caster.heal_power * heal_power_coef
           × level_scale
           × (1 + heal_pct.additive)
           × (crit ? crit_dmg_mult : 1)
```

Cap: `min(final_heal, target.max_hp - target.current_hp)` — overheal metrik olarak kayıt.

## 12. DoT / HoT

Tick-based:
```
tick_amount = base_tick
            + Σ (primary × coef_tick)
            × level_scale
            × (1 + relevant_pct)
```

Toplam damage/heal `tick_amount × (duration_ms / tick_ms)`.

Snapshot vs re-calc:
- **Snapshot** (v1): buff apply anında caster stat'ları kilitleniyor.
- **Re-calc** (Phase 02 opsiyonel): her tick fresh calc.
- Tercih: snapshot (basit, öngörülebilir).

## 13. Proc Formula

Proc chance:
```
if RNG < proc_chance_pct: fire proc
```
Bad luck protection (Phase 02): `effective_chance = chance × (1 + missed_count × factor)`.

## 14. Gear Power Contribution

Weapon damage:
```
weapon_dmg = item.rolled_stats.min_damage + rng * (item.max_damage - item.min_damage)
```
Skill scaling `weapon_damage_coef` çarpımı → base'e eklenir.

Spell power benzer.

## 15. Content Loader Validation

- Scaling.primary sum ≤ 3.0 (damage skills).
- Coefficient non-negative.
- Level_scale_factor 0..2 aralığı.
- Variance 0..15.

## 16. Test Stratejisi

- Table-driven damage compute (base + stat + gear + level).
- Crit boundary.
- Armor DR clamp.
- Level scaling curve.
- Element resist matrix.
- Variance seeded determinism.

## 17. Formulas Config

`config.combat.formulas.yaml`:
```yaml
armor_k: 60
resist_k: 60
crit_cap_pct: 75
crit_dmg_cap: 2.5
armor_dr_cap: 0.75
level_scale_per_level: 0.02
pvp_damage_mult: 0.5   # Phase 02
variance_default_pct: 5
```

## 18. Cross-References
- [Stats](stats.md), [Skills](skills.md), [Active Skills](active-skills.md), [Modifiers](modifiers.md)
- [character/stats.md](../character/stats.md), [character/health.md](../character/health.md)
- [game-design/combat.md](../game-design/combat.md)

## 19. Bekleyen Kararlar
- K_ARMOR / K_RESIST final değerleri (simülasyon).
- PvE-only cap (uzun boss encounter'da over-scaling koruma).
- Snapshot vs re-calc DoT (v2 experiment).
- Bad luck protection formulü final.
- Elemental interaction (mesela fire mob'a ice bonus) — Phase 02.
