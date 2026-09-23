# Skills (Registry + Template + Effect DSL)

## 1. Amaç
Skill'in kimliği + davranış tanımı. Registry + in-memory cache + effect resolver arayüzü.

Data model: [database/content-templates.md §4-5](../database/content-templates.md).

## 2. Skill Kayıt Modeli

Skill = `content.skill_template` satırı + opsiyonel rank tablosu.

Ana alanlar:
- `code`, `type` (active/passive/utility/racial/ultimate).
- `class_id`, `spec_id`, `race_id` (opsiyonel bağlantılar).
- `resource_type`, `resource_cost`.
- `cooldown_ms`, `gcd_bypass`, `cast_time_ms`.
- `range_kind`, `aoe_shape`.
- `effects_json` — effect DSL list.
- `scaling_json` — stat coefficient (bkz. [scaling.md](scaling.md)).
- `unlock_condition_json` — level, spec, race.
- `max_rank`.

## 3. Effect DSL

JSON schema:
```json
{
  "_v": 1,
  "effects": [
    {
      "kind": "damage",
      "target": "enemy",
      "element": "physical",
      "base": 100,
      "scaling": { "str_coef": 1.0 }
    },
    {
      "kind": "buff",
      "target": "self",
      "buff": {
        "code": "berserker_rampage_buff",
        "modifiers": [
          { "kind": "mult_add", "target": "damage_pct", "amount": 0.30 }
        ],
        "duration_ms": 20000,
        "stacks": 1
      }
    }
  ]
}
```

## 4. Effect Kinds (Whitelisted)

| Kind | Amaç | Payload |
|---|---|---|
| `damage` | Hasar | element, base, scaling, crit_bonus |
| `heal` | İyileş | base, scaling, heal_type |
| `shield` | Absorb barrier | amount, duration |
| `buff` | Positive status effect | modifiers, duration, stacks |
| `debuff` | Negative status effect | modifiers, duration |
| `dot` | Damage over time | tick_amount, tick_ms, duration, element |
| `hot` | Heal over time | tick_amount, tick_ms, duration |
| `resource_grant` | Kaynak ver (rage gain) | resource_type, amount |
| `resource_drain` | Kaynak al (mana burn) | resource_type, amount |
| `teleport` | Position (v1 opsiyonel) | destination |
| `interrupt` | Cast break | duration |
| `cc_root/stun/silence` | Crowd control | duration, dr_category |
| `dispel` | Buff remove | filter |
| `summon` | Pet/ally (druid, cleric holy weapon) | entity_code, duration |
| `taunt` | Threat force | duration |
| `execute` | HP threshold damage boost | threshold_pct, bonus_mult |
| `trigger` | Passive proc marker | trigger_id (bkz. passives.md) |
| `custom` | Handler-registered | handler_code, params |

## 5. Effect Target Types
```
self
enemy (single)
enemy_aoe
enemy_cone
enemy_line
enemy_chain
ally (single)
ally_aoe
party
raid (Phase 02)
random_enemy
random_ally
```

## 6. Registry API

```go
package skill

type Registry interface {
    ByID(ctx context.Context, id uuid.UUID) (SkillDef, error)
    ByCode(ctx context.Context, code string) (SkillDef, error)
    ForClass(ctx context.Context, classID uuid.UUID) ([]SkillDef, error)
    ForSpec(ctx context.Context, specID uuid.UUID) ([]SkillDef, error)
    ForRace(ctx context.Context, raceID uuid.UUID) ([]SkillDef, error)
    ReloadOnVersionBump(ctx context.Context, newVersion string) error
}
```

In-memory cache; content_version bump ile invalidate.

## 7. Skill Definition (Go struct)

```go
type SkillDef struct {
    ID              uuid.UUID
    Code            string
    Type            SkillType       // active/passive/utility/racial/ultimate
    ClassID         *uuid.UUID
    SpecID          *uuid.UUID
    RaceID          *uuid.UUID
    ResourceType    ResourceType
    ResourceCost    int32
    CooldownMs      int32
    GCDBypass       bool
    CastTimeMs      int32
    RangeKind       RangeKind
    AoEShape        *AoEShape
    Effects         []Effect
    Scaling         Scaling
    UnlockCondition UnlockCondition
    MaxRank         int8
    ClientHints     ClientHints    // sound, animation
}
```

## 8. Rank Sistemi
`content.skill_rank` (opsiyonel; per-skill farklıysa):
```
skill_id, rank, effects_json_override, scaling_json_override
```

Read:
```go
func (r *Registry) EffectiveEffects(ctx context.Context, skillID uuid.UUID, rank int8) []Effect
```

Skill rank up mekaniği: [mastery.md](mastery.md).

## 9. Learn / Forget

Character `character.skill_learn`:
```
character_id, skill_id, rank, xp, learned_at
```

Learn kaynakları:
- Class starter (creation'da).
- Level up unlock (`unlock_condition.min_level`).
- Spec selection.
- Race (racial).
- Recipe / trainer NPC (bazı utility skill'ler).
- Quest reward.

Forget:
- Respec ile spec skill'leri unlearn.
- Class değişimi yok (DD-063).

## 10. Skill Bar

`character.skill_bar` slot 0..11 (12 slot).
- 0..9: aktif slot (kullanıcı düzenler).
- 10: Tier II signature (varsa).
- 11: Racial active.

Reset triggers:
- Respec: 0..9 kısmen korunur (öğrenilen skill hala varsa); slot 10 signature refresh.
- Manual: `POST /v1/characters/{id}/skill-bar`.

## 11. Unlock Kontrol

```go
func (r *Registry) CanLearn(ctx context.Context, char Character, skill SkillDef) error {
    uc := skill.UnlockCondition
    if uc.MinLevel > char.Level { return errs.ErrForbidden }
    if uc.RaceRequired != nil && *uc.RaceRequired != char.RaceID { return errs.ErrForbidden }
    if uc.ClassRequired != nil && *uc.ClassRequired != char.ClassID { return errs.ErrForbidden }
    if uc.SpecRequired != nil && !specIn(char, *uc.SpecRequired) { return errs.ErrForbidden }
    return nil
}
```

## 12. Skill Compile

Content deploy sırasında **compile step**:
- Effect DSL parse + validate.
- Scaling coefficient normalize.
- Cache-friendly binary encoding (opsiyonel Phase 02 optimization).

## 13. Content Loader Validation

- Effect DSL `_v` version.
- Each effect `kind` whitelisted.
- Target type valid.
- Scaling coefficient bütçesi (sum <= 3.0 for damage skills).
- Cooldown reasonable (0..600000 ms).
- Class/race/spec ref var.

## 14. Test Stratejisi

- Registry load: valid content set → 100% skill compiled.
- Invalid effect → build fail.
- Rank override apply.
- Learn permission matrix.
- Content_version bump reload.

## 15. Metrics
- `skill_registry_reload_total`.
- `skill_definitions_active` gauge.
- `skill_effect_compile_errors_total`.

## 16. Cross-References
- [Passives](passives.md), [Active Skills](active-skills.md), [Cooldown](cooldown.md), [Scaling](scaling.md), [Mastery](mastery.md)
- [class/specialization.md](../class/specialization.md), [race/racial-skills.md](../race/racial-skills.md)
- [database/content-templates.md](../database/content-templates.md), [game-design/skills.md](../game-design/skills.md)

## 17. Bekleyen Kararlar
- Rank sistemi full aktif mi (v1'de max_rank=1 sabit; Phase 02'de rank progression).
- Trigger-based passive effect DSL (bkz. passives.md) — Phase 02 detaylı iterasyon.
- Skill trees UI (Phase 02).
