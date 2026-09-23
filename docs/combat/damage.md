# Damage (Pipeline Detay)

## 1. Amaç
Damage compute pipeline — skill effect + auto-attack + DoT tick. Formula, mitigation, crit, variance, cap.

Foundation formula: [skill-engine/scaling.md](../skill-engine/scaling.md). Bu doküman combat-özel entegrasyon.

## 2. Damage Request

```go
type DamageRequest struct {
    ActorID      uuid.UUID       // saldıran
    TargetID     uuid.UUID       // hedef
    Skill        SkillDef        // veya nil (auto-attack)
    BaseAmount   int32           // effect.base veya weapon damage roll
    Element      string          // physical/fire/ice/nature/holy/shadow/arcane
    Kind         DamageKind      // direct/dot_tick/aoe/proc/reflection
    Scaling      Scaling         // effect.scaling
    Seed         Seed            // RNG stream
    IsCritForced bool             // opsiyonel; test/preview override
}

type DamageResult struct {
    RawAmount    int64
    ArmorMitigated int64
    ResistMitigated int64
    Blocked      int64
    Absorbed     int64            // shield
    FinalAmount  int64
    Crit         bool
    Evaded       bool
    Effects      []SecondaryEffect  // lifesteal, threat, on-hit proc
}
```

## 3. Pipeline (Full Detail)

```
1. raw = base + Σ(primary_stat × scaling.coef)
2. raw *= level_scale(actor.Level)
3. raw += gear_power_contribution (weapon_dmg × coef)
4. raw *= (1 + damage_pct.additive_stack)   // race/class/spec/set modifiers
5. raw *= (1 + damage_pct.element_bonus)     // damage_pct.fire, etc.
6. raw *= rank_multiplier (skill rank; v1: 1.0)
7. Roll evasion: if RNG < target.evasion_pct: evaded=true, final=0, exit.
8. Roll crit: if RNG < actor.crit_chance_pct: crit=true, raw *= actor.crit_dmg_mult.
9. Apply variance: raw *= 1 + (RNG - 0.5) × variance_pct/100 (±5%)
10. Element mitigation:
    if element = physical: armor_dr calc.
    else: resist_dr calc.
11. Block: if target has shield && RNG < block_pct: raw *= (1 - block_reduction)
12. Absorb: if target has shield buff, damage against shield first; residual → HP.
13. On-hit procs: actor items/passive triggers on_hit_dealt.
14. Threat generation: threat_amount = final × threat_modifier.
15. Return.
```

## 4. Armor DR Formula

```
armor_dr = target.armor / (target.armor + K_ARMOR × actor.level)
armor_dr = clamp(armor_dr, 0, 0.75)
mitigated = raw × armor_dr
final -= mitigated
```

`K_ARMOR = 60` (v1; bkz. [skill-engine/scaling.md §7](../skill-engine/scaling.md)).

## 5. Element Resist

```
element_resist = target.magic_resist_for_element
resist_dr = element_resist / (element_resist + K_RESIST × actor.level)
resist_dr = clamp(0, 0.75)
mitigated = raw × resist_dr
final -= mitigated
```

Element-specific resist (fire_resist, ice_resist etc.) — v1'de tek `magic_resist`; Phase 02 element split.

## 6. Crit

```
crit_chance = clamp(actor.crit_chance_pct, 0, 75)
crit_dmg_mult = actor.crit_dmg_mult   // 1.5..2.5

if crit_chance > 50:
    excess = crit_chance - 50
    effective = 50 + excess × 0.5   // diminishing
    crit_chance = clamp(effective, 50, 75)

if seed.NextFloat() < crit_chance/100:
    crit = true
    raw *= crit_dmg_mult
```

## 7. Evasion

```
evasion_pct = clamp(target.evasion_pct, 0, 30)  // soft cap
if seed.NextFloat() < evasion_pct/100:
    evaded = true
    final = 0
    return
```

Evade → 0 damage; threat 0.

## 8. Block

Shield equipped + not 2H weapon:
```
block_chance = clamp(target.block_pct, 0, 40)
if seed.NextFloat() < block_chance/100:
    blocked = true
    block_reduction = 0.5 (v1)  // %50 azalt
    raw *= (1 - block_reduction)
```

## 9. Absorb (Shield Buff)

Status effect `shield_amount`:
```
if target has shield buff (kind=shield):
    absorbed = min(shield.remaining, raw)
    shield.remaining -= absorbed
    raw -= absorbed
    if shield.remaining <= 0: expire shield
```

## 10. Variance

```
variance_pct = skill.scaling.variance_pct   // default 5
variance_multiplier = 1 + (seed.NextFloat() - 0.5) × 2 × variance_pct/100
raw *= variance_multiplier   // ±5%
```

Deterministic seeded.

## 11. Damage Clamp

- **Min 1**: her damage minimum 1 (evade hariç).
- **Max int32**: overflow koruma; extreme buff kombo clamp.

## 12. Damage Kinds

| Kind | Notlar |
|---|---|
| `direct` | Skill cast tetiklenmiş |
| `auto_attack` | Weapon-based |
| `dot_tick` | DoT periodic |
| `aoe` | AoE ability; multi-target |
| `proc` | Passive/gear proc trigger |
| `reflection` | Reflection damage (nadir; special) |
| `percent_max_hp` | Boss execute; % based (cap 30%/sec) |
| `true_damage` | Mitigation bypass (rare) |

Farklı kind farklı mitigation path.

## 13. On-Hit Procs

Damage dealt sonrası passive trigger:
- `on_hit_dealt` → passives (bkz. [skill-engine/passives.md](../skill-engine/passives.md)).
- Gear proc (Flame Weapon → fire dmg).
- Set bonus signature trigger.

Actor's passives evaluated.

## 14. Lifesteal / Vamp

Actor `lifesteal_pct` (item bonus):
```
heal_return = final × lifesteal_pct/100
actor.CurrentHP += heal_return   // clamp to max_hp
```

Spell vamp benzeri (spell damage için).

## 15. Threat Generation

```
threat_amount = final × actor.threat_modifier
if actor.SpecTagsIncludes("tank_main"):
    threat_amount *= 1.5
target.ThreatTable[actor.ID] += threat_amount
```

Detay: [threat.md](threat.md).

## 16. Event Emission

Post-damage event `combat.damage_dealt`:
```json
{
  "actor_id": "...",
  "target_id": "...",
  "amount": 152,
  "kind": "direct",
  "element": "physical",
  "crit": true,
  "evaded": false,
  "blocked": false,
  "absorbed": 20,
  "skill_code": "berserker_rampage"
}
```

Subscribers: passive triggers, notification, achievement, session log.

## 17. AFK Aggregate

AFK reveal içinde:
- Per-cast damage compute analytical (skillEngine.EstimateDPS).
- Kill count = f(mob HP, DPS, spawn).
- Per-damage event yok; aggregate summary only.

## 18. Config

```yaml
combat:
  k_armor: 60
  k_resist: 60
  crit_cap_pct: 75
  crit_diminishing_start: 50
  evasion_cap_pct: 30
  block_cap_pct: 40
  block_reduction: 0.5
  armor_dr_cap: 0.75
  variance_default_pct: 5
  min_damage: 1
  percent_max_hp_dps_cap: 30
```

## 19. Test Stratejisi

- Table-driven: (actor.stats, target.stats, skill) → expected damage.
- Crit boundary.
- Evade zero damage.
- Absorb → HP residual.
- Golden replay: seed → same damage.

## 20. Metrics
- `combat_damage_dealt_total{element,kind}`.
- `combat_damage_crit_ratio` gauge.
- `combat_damage_evade_ratio` gauge.
- `combat_damage_max_amount` histogram.

## 21. Cross-Reference
- [Combat Loop](combat-loop.md), [Auto-Attack](auto-attack.md), [Healing](healing.md), [Status Effects](status-effects.md), [Threat](threat.md)
- [skill-engine/scaling.md](../skill-engine/scaling.md), [skill-engine/passives.md](../skill-engine/passives.md)
- [character/stats.md](../character/stats.md), [character/health.md](../character/health.md)
- [item-engine/stats.md](../item-engine/stats.md)

## 22. Bekleyen Kararlar
- Elemental resist split (v1: unified magic_resist; Phase 02 split).
- PvP damage split formula final Phase 02.
- Reflection damage cap.
- True damage yaygınlık — nadir kalmalı (nerf risk).
