# Healing

## 1. Amaç
Heal + shield + regen flow. Damage sisteminin pozitif ikizi. Cleric/Druid ana odak.

## 2. Heal Request

```go
type HealRequest struct {
    ActorID      uuid.UUID
    TargetID     uuid.UUID
    Skill        SkillDef
    BaseAmount   int32
    Kind         HealKind        // direct/hot_tick/percent_max_hp
    Scaling      Scaling
    Seed         Seed
}

type HealResult struct {
    RawAmount    int64
    AppliedAmount int64
    OverhealAmount int64
    Crit         bool
    FinalHP      int64
    MaxHP        int64
}
```

## 3. Pipeline

```
1. raw = base + Σ(primary × scaling.coef)
2. raw += heal_power × heal_power_coef
3. raw *= level_scale(actor.Level)
4. raw *= (1 + heal_pct.additive_stack)
5. Roll crit: if RNG < actor.crit_chance_pct: crit=true, raw *= crit_dmg_mult (heal genellikle aynı crit_dmg)
6. Variance: raw *= 1 ± variance_pct
7. Target dead? → reject (revive skill'i ayrı).
8. Apply:
   applied = min(raw, target.MaxHP - target.CurrentHP)
   overheal = raw - applied
   target.CurrentHP += applied
9. Overheal metric log.
10. Threat generation: healed threat_amount = applied × 0.5 (heal aggro half).
11. Emit event 'combat.heal_dealt'.
```

## 4. Heal Formula (Örnek Cleric Basic Heal)

```
Cleric.wis = 50, heal_power = 120, level = 25
BasicHeal.base = 30, scaling.wis_coef = 1.5, heal_power_coef = 1.0

raw = 30 + 50×1.5 + 120×1.0 = 225
raw *= level_scale(25) = 225 × 1.48 = 333
raw *= (1 + heal_pct.additive)  // e.g. +10% from spec = 366
crit_chance = 5%; miss crit; variance 1.02 → 373
target.CurrentHP was 400/900, so:
  applied = min(373, 500) = 373
  overheal = 0
  target.CurrentHP = 773
```

## 5. Heal Kind

| Kind | Notlar |
|---|---|
| `direct` | Skill cast tetiklenmiş |
| `hot_tick` | HoT periodic tick |
| `percent_max_hp` | % based heal (rare; e.g. Cleric ultimate) |
| `on_hit_heal` | Attack heal (vampiric) |

## 6. Overheal

Overheal = heal beyond MaxHP. Metric için kaydedilir; oyuncu-görünmez veya UI'da azaltılmış görsel.

## 7. Shield Absorb (Positive Shield)

Shield status effect apply:
```go
statusEngine.Apply(target, StatusInstance{
    Kind: shield,
    ShieldAmount: shield_amount,
    Duration: duration,
})
```

Damage geldiğinde shield önce absorb eder (bkz. [damage.md §9](damage.md)).

Shield stacking: multiple shield → first-in-first-out veya highest first (v1: multi shield paralel, sum absorb).

## 8. Regen (Passive)

Combat regen tick (2 sn):
```
hp_regen_per_2s = 1 + 0.1 × effective_sta
combat_multiplier = 0.2   (in-combat)
                    1.0   (out-of-combat)

target.CurrentHP += hp_regen_per_2s × combat_multiplier
target.CurrentHP = min(target.CurrentHP, target.MaxHP)
```

Bkz. [character/health.md §3.2](../character/health.md).

## 9. Rest Bonus (Phase 02)

City inn'de rest → +100% HP regen 15 sn sonra. Phase 02.

## 10. Resurrection (Special)

Dead target'a heal reject; ancak "Resurrect" skill'i özel:
```go
if req.Skill.Code == "resurrect" && target.IsDead:
    revive target with HP = raw
    target.IsDead = false
    emit event 'combat.resurrected'
```

Party revive flow: [character/health.md §5.4](../character/health.md).

## 11. AoE Heal

AoE heal effect:
```
targets = findAlliesInRadius(actor.Position, aoe.radius)
for _, t := range targets:
    computeHeal(t, ...)
```

v1'de "party heal" — party members only (no positional).

## 12. Threat from Heal

Heal → half threat (bkz. §3.10):
```
threat_amount = applied × 0.5
```

Heal-through-damage spec (Confessor Cleric): threat_amount = applied × 0.8 (more threat).

Bkz. [threat.md](threat.md).

## 13. Config

```yaml
heal:
  crit_cap_pct: 75     # heal crit chance same as damage
  variance_pct: 5
  overheal_track_enabled: true
  heal_threat_ratio: 0.5
```

## 14. Metrics
- `combat_heal_dealt_total`.
- `combat_heal_overheal_ratio` gauge.
- `combat_heal_crit_ratio` gauge.

## 15. Event

`combat.heal_dealt`:
```json
{
  "actor_id": "...",
  "target_id": "...",
  "applied": 373,
  "overheal": 0,
  "crit": false,
  "skill_code": "cleric_basic_heal"
}
```

## 16. Cross-Reference
- [Combat Loop](combat-loop.md), [Damage](damage.md), [Status Effects](status-effects.md), [Threat](threat.md), [Death & Revive](death-revive.md)
- [character/health.md](../character/health.md), [character/mana.md](../character/mana.md)
- [skill-engine/scaling.md](../skill-engine/scaling.md)

## 17. Bekleyen Kararlar
- Overheal metric visibility (client UI).
- Heal crit damage multiplier vs healing crit multiplier (v1 aynı; Phase 02 farklı olabilir).
- Cross-heal from foes (charm effect) Phase 02.
- Rez skill mekanikleri final.
