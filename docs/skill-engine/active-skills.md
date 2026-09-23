# Active Skills (Cast Pipeline)

## 1. Amaç
Kullanıcı/rotation tetikli skill'lerin cast pipeline'ı: validation → cost → cast time → effect resolve → cooldown set → event emit.

## 2. Sözleşme
```go
package skillengine

type Target struct {
    Kind      TargetKind      // self, enemy, ally, aoe_point
    ID        *uuid.UUID       // spesifik entity
    Position  *Vec3            // AoE için (Phase 02 positional)
}

type CastRequest struct {
    CharacterID uuid.UUID
    SkillCode   string
    Target      Target
    Rank        int8            // opsiyonel; default learned rank
    Source      CastSource      // manual, rotation, item_use
    Seed        Seed
}

type CanCastResult struct {
    OK             bool
    Reason         string       // "cooldown", "resource", "range", "target_invalid"
    CooldownLeftMs int32
    ResourceHave   int32
    ResourceNeed   int32
}

type CastResult struct {
    SkillCode      string
    ResourceSpent  int32
    Cooldowns      []CooldownSet    // main, gcd, category
    Effects        []AppliedEffect
    Events         []EngineEvent
    RNGDraws       []RNGDraw
    ExecutedAt     time.Time
}
```

## 3. Pipeline

```
Cast(req):
  1. Load character (stats + rotation + learn).
  2. Load skill definition (SkillRegistry.Get).
  3. Verify character has learned skill.
  4. Verify unlock condition.
  5. Global Cooldown (GCD) check unless gcd_bypass.
  6. Skill cooldown check (CooldownManager).
  7. Category cooldown check (opsiyonel).
  8. Resource check (rage/mana/energy/focus/nature).
  9. Target validate (range, alive, faction).
 10. Cast time check (channel etc.).
 11. Roll RNG (crit, proc).
 12. Resolve effects (EffectResolver).
 13. Debit resource.
 14. Set cooldowns.
 15. Grant mastery XP (MasteryEngine).
 16. Emit events (skill_used, damage_dealt, buff_applied).
 17. Return CastResult.
```

Adımlar atomic içinde (1-11 read-only; 12-17 transaction).

## 4. Cast Time (Channel/Instant/Cast)

- **Instant** (`cast_time_ms=0`): pipeline hemen tamamlanır.
- **Cast** (`cast_time_ms>0`): interrupt-able.
- **Channel** (Phase 02): tick-based; her tick effect fire.

Cast timing:
```
t0: Cast() start
    validate, resource check, cooldown check
    reserve cooldown (partial? — HAYIR, cast fail iade)
    schedule effect resolve at t0 + cast_time_ms
    schedule interrupt window
    if interrupted before completion:
        release reservation
    else:
        resolve effects, finalize
```

v1'de tercih: **cast start'ta hemen resource debit + full cooldown set**; interrupt olursa iade (fair).

## 5. Interrupt

Trigger:
- Attacker uses interrupt skill.
- Damage taken beyond threshold (opsiyonel).
- Silence debuff.

Effect:
- Cast cancel.
- Refund resource (%50-100 config).
- Skill goes on **partial cooldown** (5 sn) — locked cast can't be retried immediately.

## 6. Rotation Integration

AFK rotation:
- Rotation editor JSON `priority_list`.
- Combat tick evaluator:
  ```
  for each entry in priority_list:
    if canCast(entry.skill, entry.cond, currentState):
      cast(entry.skill)
      break
  ```
- Default rotation: spec'in `default_rotation_json`.

Manual play: client sends `POST /v1/characters/{id}/cast`.

## 7. Endpoint

**POST `/v1/characters/{id}/cast`**
Auth + CharacterScope.

Request:
```json
{
  "skill_code": "berserker_rampage",
  "target": { "kind": "self" },
  "source": "manual"
}
```

Response 200:
```json
{
  "result": {
    "cooldowns": [{ "skill_code": "berserker_rampage", "remaining_ms": 60000 }],
    "resource_spent": 40,
    "effects": [
      { "kind": "buff", "buff_code": "rampage_active", "duration_ms": 20000 }
    ]
  },
  "vitals": { "hp": 425, "resource": 20 }
}
```

Response 400:
```json
{ "error": { "code": "skill.cooldown_active", "details": { "remaining_ms": 45000 } } }
```

Fail codes:
- `skill.not_learned`
- `skill.cooldown_active`
- `skill.gcd_active`
- `skill.insufficient_resource`
- `skill.target_invalid`
- `skill.out_of_range`
- `skill.silenced`
- `skill.dead`

## 8. Target Validation

```go
func validateTarget(caster Character, target Target, skill SkillDef) error {
    if skill.RangeKind == Melee && target.Kind == Enemy {
        // v1: mesafe kontrolü minimal (positional yok)
        // ancak hunt spot instance context'inde geçerlilik
    }
    if skill.EffectTarget == Ally && target.Kind != Ally && target.Kind != Self {
        return errs.ErrValidation
    }
    ...
}
```

## 9. Effect Resolver

Effect list her element için:
```go
func resolveEffect(caster, target, skill, seed) AppliedEffect {
    switch effect.Kind {
    case "damage":
        return computeDamage(caster.Derived, target.Derived, effect, seed)
    case "heal":
        return computeHeal(caster.Derived, target.Derived, effect)
    case "buff":
        return applyBuff(target, effect.Buff)
    case "dot":
        return applyDoT(target, effect)
    ...
    }
}
```

Damage formul detay: [scaling.md](scaling.md).

## 10. GCD (Global Cooldown)

- Baz 1000 ms.
- Haste ile azalır: `gcd = 1000 * (1 - haste_pct/100)`; min 500 ms.
- `gcd_bypass=true` skill'ler için triggerlanmaz.
- Redis key: `char:gcd:<char_id>` TTL gcd_ms.

## 11. Category Cooldown

Bazı skill grupları paylaşımlı cooldown:
- "Consumable potion" grubu: 60 sn hepsi.
- "Trinket use" grubu: 90 sn hepsi.

Skill definition `category_cooldown_key` (opsiyonel):
```yaml
category_cooldown_key: potion_group
category_cooldown_ms: 60000
```

Redis key `char:catcd:<char>:<key>`.

## 12. Cast Result Persistence

Cast eventi audit için:
```
audit.event_log:
  event_kind='combat.skill_cast'
  subject=character_id
  payload={skill, effects, target, seed}
```

Volume yüksek — bulk insert; retention 14 gün (debug).

## 13. RNG Determinism

Cast per RNG stream:
- seed = `hash(character.seed, correlation_id, skill_code, tick_index)`.
- crit roll, proc roll, damage variance — hepsi bu stream'den.

Test: same seed → same result.

## 14. AFK Sim Uyumu

Reveal içinde:
- Each combat tick, engine.Cast simulate.
- Determinism preserved.
- Effects aggregate ⟶ reveal summary.

## 15. Test Stratejisi
- Table-driven CanCast matrix (all fail reasons).
- Cast pipeline full trace assertion.
- Cast time interrupt behavior.
- Category cooldown shared.
- Rotation eval simulation.

## 16. Metrics
- `skill_cast_total{skill_code,result}`.
- `skill_cast_duration_seconds{skill_code}`.
- `skill_cast_reject_reason_total{reason}`.
- `skill_interrupt_total`.

## 17. Cross-References
- [Skills](skills.md), [Cooldown](cooldown.md), [Scaling](scaling.md), [Modifiers](modifiers.md), [Passives](passives.md)
- [game-design/combat.md](../game-design/combat.md), [game-design/skills.md](../game-design/skills.md)
- [character/health.md](../character/health.md), [character/mana.md](../character/mana.md)

## 18. Bekleyen Kararlar
- Positional targeting (line/cone/circle) v1 basitlik: soft targeting.
- Cast interrupt threshold configurable per skill.
- Channel skill'ler v2.
- Auto-attack modeli (bkz. combat) — pipeline'a nasıl entegre.
