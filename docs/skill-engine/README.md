# Skill + Stat Engine (Faz 9)

Combat, AFK, aktiviteler ve UI'nin altında çalışan **hesap motoru**. Statları toplar, modifier'ları uygular, skill'lerin cost/CD/effect'ini yürütür, mastery'yi ilerletir, scaling coefficient'lerini uygular.

**Bu motor combat'ın parçası değildir; combat bu motoru kullanır.** Aynı motor manuel oyun, AFK simülasyonu, test/replay için tek sözleşme sunar.

## 1. Doküman Haritası

| Konu | Doküman |
|---|---|
| Genel bakış + engine flow (bu doküman) | [README.md](README.md) |
| [Stats](stats.md) | Stat compute pipeline + cache + invalidation |
| [Modifiers](modifiers.md) | Modifier stack + apply order + budget |
| [Skills](skills.md) | Skill registry + template + effect DSL |
| [Passives](passives.md) | Passive skill mechanics + trigger types |
| [Active Skills](active-skills.md) | Cast pipeline + validation + effect resolve |
| [Cooldown](cooldown.md) | Redis-backed CD manager + GCD + reduction |
| [Mastery](mastery.md) | Use-based skill mastery + rank up |
| [Scaling](scaling.md) | Stat coefficients + level scaling + PvE/PvP split |

## 2. Engine Katmanları

```
┌───────────────────────────────────────────────────────┐
│              Domain Entry Points                      │
│  (combat.tick, activity.reveal, action.cast, buff.apply)
└───────────────────┬───────────────────────────────────┘
                    │
                    ▼
┌───────────────────────────────────────────────────────┐
│          SkillStatEngine (internal/engine)            │
│                                                        │
│  Compute effective stats (StatEngine)                  │
│  Resolve modifiers (ModifierStack)                     │
│  Evaluate skill (SkillEngine)                          │
│  Check cost + CD (ResourceManager, CooldownManager)    │
│  Apply effects (EffectResolver)                        │
│  Roll RNG (deterministic seed)                         │
│  Emit events (EventEmitter)                            │
└───────────────────┬───────────────────────────────────┘
                    │
                    ▼
┌───────────────────────────────────────────────────────┐
│  Data Layer: content (skills, race, class, spec)      │
│              character (stats, learn, buffs)          │
│              config (formulas)                        │
└───────────────────────────────────────────────────────┘
```

## 3. Ortak İlkeler

1. **Server-authoritative** — client sadece intent gönderir; engine karar verir.
2. **Deterministic** — aynı input + seed → aynı output (test + replay + AFK).
3. **Pure functions** — effect resolver stateless; side effect'ler explicit (event emit + repo write).
4. **Cache-friendly** — derived stat cache; invalidation event-driven.
5. **Data-driven** — kod skill/class/race code'u bilmez; effect DSL çözer.
6. **Combat-agnostic** — engine combat için özelleşmez; buff apply, gathering yield modifier de aynı pattern.

## 4. Ana Interface (Yüksek Seviye)

```go
package engine

// Engine — tüm skill/stat operasyonlarının giriş noktası.
type Engine interface {
    // Stat compute
    ComputeDerived(ctx context.Context, charID uuid.UUID, ctx StatContext) (Derived, error)

    // Skill eval
    CanCast(ctx context.Context, charID uuid.UUID, skillCode string, target Target) (CastCheck, error)
    Cast(ctx context.Context, charID uuid.UUID, skillCode string, target Target, seed Seed) (CastResult, error)

    // Modifier lifecycle
    ApplyBuff(ctx context.Context, target Target, buff BuffDef) (BuffInstance, error)
    ExpireBuff(ctx context.Context, buffID uuid.UUID) error

    // Cooldown
    RemainingCooldown(ctx context.Context, charID uuid.UUID, skillCode string) (time.Duration, error)

    // Mastery
    GrantMasteryXP(ctx context.Context, charID uuid.UUID, skillCode string, xp int32) error
}
```

## 5. Bileşenler

### 5.1 StatEngine
Bkz. [stats.md](stats.md). Primary + derived + tertiary hesabı.

### 5.2 ModifierStack
Bkz. [modifiers.md](modifiers.md). Additive/multiplicative sıralı zincir.

### 5.3 SkillEngine
Bkz. [skills.md](skills.md). Skill registry + template resolve.

### 5.4 PassiveEngine
Bkz. [passives.md](passives.md). Passive trigger + on-hit + on-cast.

### 5.5 ActiveSkillPipeline
Bkz. [active-skills.md](active-skills.md). Cast validation → cost → effect.

### 5.6 CooldownManager
Bkz. [cooldown.md](cooldown.md). Redis-backed CD + GCD.

### 5.7 MasteryEngine
Bkz. [mastery.md](mastery.md). Use-based rank up.

### 5.8 ScalingEngine
Bkz. [scaling.md](scaling.md). Coefficient apply + level scaling.

## 6. Ortak Kullanılan Tipler

```go
type Derived struct {
    Primary   PrimaryStats
    Combat    CombatStats
    Vital     VitalStats
    Defense   DefenseStats
    Tertiary  TertiaryStats
    Version   int32   // cache invalidation
    ComputedAt time.Time
}

type Modifier struct {
    Kind      ModifierKind  // flat_bonus, multiplier, ...
    Target    string        // "crit_chance_pct", "str", "damage_pct.fire", ...
    Amount    float64
    Scope     Scope         // always, combat, profession:mining, ...
    SourceKind string       // race, class, spec, gear, buff
    SourceID  uuid.UUID
    Priority  int           // apply order
}

type BuffInstance struct {
    ID        uuid.UUID
    Def       BuffDef
    Target    Target
    AppliedAt time.Time
    ExpiresAt time.Time
    Stacks    int32
    SourceID  uuid.UUID  // caster / item
}

type CastResult struct {
    Skill        SkillDef
    Effects      []AppliedEffect
    ResourceCost int32
    Cooldowns    []CooldownSet
    Events       []Event
}
```

## 7. Data Flow — Örnek: Skill Cast

```
1. Client → POST /v1/characters/{id}/cast {skill: "berserker_rampage", target: self}
2. Middleware auth + character scope.
3. Handler → engine.Cast(ctx, char, "berserker_rampage", target, deriveSeed(...))
4. engine:
   a. resolve skill (SkillEngine.Get)
   b. compute derived stat (StatEngine.ComputeDerived)
   c. check cost (ResourceManager.HasResource)
   d. check cooldown (CooldownManager.RemainingCooldown)
   e. check range/target validity
   f. compute effect list (EffectResolver)
   g. apply RNG (Random System stream)
   h. write buffs (character.stats_buff)
   i. set cooldowns (CooldownManager.Set)
   j. debit resource (ResourceManager.Consume)
   k. grant mastery XP (MasteryEngine.Grant)
   l. emit events (combat.skill_used, buff.applied, ...)
5. Handler → WS push 'skill_effect' + HTTP 200.
```

Kritik: adımlar (h-l) atomic — tümü Postgres transaction içinde; Redis (cooldown) kısmı outbox pattern ile.

## 8. Deterministic RNG Kullanımı

Skill result rastgele (crit, proc, RNG effect):
- `seed = hash(character.seed, skill_code, tick_or_correlation_id)`.
- RNG stream deterministik.
- Replay support: same seed + same input → same result.
- AFK Engine reveal aynı motor kullanır; determinism zorunluluğu (bkz. [systems/afk-engine.md](../systems/afk-engine.md)).

## 9. Cache Stratejisi

- **Derived stat cache**: `character.stats.derived_cache_json` (Postgres); invalidation flag ile lazy recompute.
- **Content cache**: skill_template, race, class, spec — in-memory (content_version bump ile reload).
- **Cooldown cache**: Redis (fast expiry).
- **Buff snapshot**: Postgres (`character.stats_buff`) + Redis flag (aktif buff sayısı hint).

## 10. Test İskeleti

- Golden replay: seed + input → expected result JSON.
- Table-driven modifier apply order.
- Cooldown expiration timing.
- Buff stack + refresh semantics.
- Cache invalidation triggers.

## 11. Cross-Cutting

- **Combat**: engine'in en büyük tüketicisi; her tick engine.CanCast + Cast.
- **AFK Reveal**: activity.simulate içinde engine kullanır (aynı sözleşme).
- **Crafting**: engine.CanCast benzeri "can craft" check (Phase 02).
- **Vendor**: modifier stack vendor price'a uygulanır (charisma bonus, race affinity).

## 12. İlişkiler
- [character/stats.md](../character/stats.md), [character/attributes.md](../character/attributes.md)
- [race/passive.md](../race/passive.md), [race/racial-skills.md](../race/racial-skills.md)
- [class/specialization.md](../class/specialization.md)
- [game-design/combat.md](../game-design/combat.md), [game-design/skills.md](../game-design/skills.md)
- [systems/random-system.md](../systems/random-system.md), [systems/afk-engine.md](../systems/afk-engine.md)
- [database/content-templates.md](../database/content-templates.md), [database/character.md](../database/character.md)

## 13. Ortak Bekleyen Kararlar
- Combat damage/heal formula final tuning (bkz. [scaling.md](scaling.md)).
- PvP damage split (Phase 02) — ayrı katman.
- Client-side simulation preview (Phase 02) — read-only, offline verify için.
- Advanced trigger sistem (combo point, resource threshold) — Phase 02 detay.
