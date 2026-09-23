# Status Effects (Buff / Debuff / DoT / HoT / CC / Shield)

## 1. Amaç
Karakter/mob üzerinde geçici etki (magnitude + duration) modeli. Buff/debuff/DoT/HoT/CC/shield hepsi tek soyutlama.

Foundation: [skill-engine/modifiers.md](../skill-engine/modifiers.md) — modifier stack. Bu doküman combat-özel entegrasyon.

## 2. Model

Kaynak: [database/character.md §... status_effect (yeni Phase 17)](../database/character.md).

`character.stats_buff` şeması (bkz. [skill-engine/modifiers.md §9](../skill-engine/modifiers.md)):
```sql
character.stats_buff
  id, character_id, buff_code
  modifiers JSONB
  source_kind, source_id
  applied_at, expires_at, stacks
```

Ek Phase 17:
```sql
ALTER TABLE character.stats_buff
  ADD COLUMN kind VARCHAR(16) NOT NULL DEFAULT 'buff',   -- buff|debuff|dot|hot|shield|cc
  ADD COLUMN tick_ms INT NULL,                            -- DoT/HoT tick interval
  ADD COLUMN last_tick_at TIMESTAMPTZ NULL,
  ADD COLUMN tick_amount INT NULL,                        -- DoT tick damage
  ADD COLUMN shield_amount INT NULL,                      -- shield absorb remaining
  ADD COLUMN dispel_category VARCHAR(24) NULL;            -- physical|magic|curse|poison
```

## 3. Status Kinds

| Kind | Amaç | Örnek |
|---|---|---|
| `buff` | Positive modifier | Rampage +30% dmg 20s |
| `debuff` | Negative modifier | Weakness -10% dmg |
| `dot` | Damage over time | Bleeding 50 dmg/2s |
| `hot` | Heal over time | Rejuvenation 30 heal/2s |
| `shield` | Absorb barrier | Stone Shield 500 absorb |
| `cc` | Crowd control | Stun/Root/Silence |

## 4. Apply Request

```go
type ApplyStatusRequest struct {
    ActorID       uuid.UUID       // uygulayıcı
    TargetID      uuid.UUID       // hedef
    Kind          StatusKind
    Code          string           // "berserker_rampage_buff"
    DurationMs    int32
    TickMs        int32            // DoT/HoT için
    TickAmount    int32            // DoT/HoT için
    ShieldAmount  int32            // shield için
    Modifiers     []Modifier       // buff/debuff için
    Stacks        int32
    DispelCategory string
    SnapshotStats bool             // stat scaling snapshot al mı
    Source        SourceRef
    Seed          Seed
}
```

Snapshot (DD-070): buff apply anında caster stat'lar kilitlenir; sonraki tick'ler fresh recompute değil.

## 5. Apply Pipeline

```
1. Load target.
2. Check existing status same code:
   - Refresh (same source): reset duration, refresh stacks (cap).
   - Different source: additive stack.
   - Same source different stacks: apply based on stack rules.
3. Insert status_effect row.
4. Modifier stack invalidate (cache).
5. Emit event 'combat.status_applied'.
6. If kind=cc: DR check (bkz. §CC Diminishing Returns).
```

## 6. Stack Rules

Aynı buff_code multiple source:
- **Additive stacks**: sum stacks; cap N (config per buff).
- **Refresh**: duration reset to max.
- **Independent**: different source_id → farklı row (paralel).

Örnek Bleed:
- Warrior Bleed + Rogue Bleed → aynı buff_code, farklı source → paralel; her biri ayrı DoT.
- Warrior Bleed re-cast → refresh (same source).

## 7. DoT / HoT Tick

Live combat loop tick (bkz. [combat-loop.md §4](combat-loop.md)):
```go
for _, st := range activeDotHotStatuses {
    if now - st.LastTickAt >= st.TickMs {
        if st.Kind == "dot":
            // snapshot amount; caster stat'lar dondurulmuş
            damageEngine.Apply(target, DamageRequest{
                BaseAmount: st.TickAmount,
                Kind: dot_tick,
                Element: st.Meta.Element,
            })
        else if st.Kind == "hot":
            healEngine.Apply(target, HealRequest{...})
        st.LastTickAt = now
    }
}
```

AFK aggregate:
- Reveal'de DoT/HoT total tick count = duration / tick_ms.
- Total damage = tick_count × tick_amount.
- Analytical; per-tick loop yok.

## 8. Expire

Expire check (per tick veya per read):
```
if now >= expires_at:
    delete status_effect row
    modifier cache invalidate
    emit event 'combat.status_expired'
```

Cleanup worker: her dakika expired batch delete (bkz. [skill-engine/modifiers.md §10](../skill-engine/modifiers.md)).

## 9. Dispel

Dispel skill kaldırır belirli status:
```go
func dispel(target, categoryFilter):
    candidates := activeStatusesByCategory(target, categoryFilter)
    if len(candidates) == 0: return
    highest := pickHighestPriority(candidates)
    expire(highest)
    emit 'combat.status_dispelled'
```

Dispel category:
- `magic` — mage debuff.
- `curse` — dark magic.
- `physical` — bleed, weapon debuff.
- `poison` — alchemy debuff.

Bazı boss debuff `undispellable=true` flag.

## 10. CC (Crowd Control)

CC status'lar:
- **Stun**: hedef aksiyon yapamaz (skill cast + movement blocked).
- **Root**: movement blocked, cast OK.
- **Silence**: skill cast blocked, movement OK.
- **Fear**: hedef random hareket (Phase 02 positional).
- **Charm**: hedef ally saldırır (Phase 02).

Duration ms; genelde kısa (1-5 sn).

## 11. CC Diminishing Returns (DR)

Aynı CC kategori repeatedly apply → duration azalır (immunity build):
```
cc_category = "stun"
if target.CCHistory[stun].applied_within_last_20_sec:
    diminishing_level = target.CCHistory[stun].level
    diminishing_multiplier = [1.0, 0.5, 0.25, 0.0]  // level 0/1/2/3 immune
    effective_duration = raw_duration × multiplier
    target.CCHistory[stun].level = min(3, level + 1)
```

DR reset: 20 sn no CC in category → reset to 0.

DR categories: stun, root, silence, disorient, disable.

## 12. Shield (Absorb Buff)

Shield status:
```
apply: shield_amount = 500, duration_ms = 8000
damage hit:
    absorbed = min(shield.shield_amount, damage_raw)
    shield.shield_amount -= absorbed
    damage_raw -= absorbed
    if shield.shield_amount <= 0:
        expire shield
```

Shield stacking:
- Multi-shield paralel; damage önce highest priority shield'ı tüketir.
- Priority: kaynağa göre (buff manifest'inde `shield_priority`).

## 13. Immunity Buffs

Bazı buff hedefi CC-immune yapar:
- "Divine Shield": %5 sn CC immunity.
- Bkz. [item-engine/unique-effect.md](../item-engine/unique-effect.md) — Phoenix Rebirth 3s immune.

Immunity flag status effect'in kendisinde.

## 14. Snapshot vs Re-Calc (DD-070)

- **Snapshot** (v1): buff apply anında caster stat'lar tick_amount hesabına gider; sabit.
- **Re-calc** (Phase 02 opt): her tick fresh calc.

## 15. Event

```json
{
  "kind": "status_applied",
  "actor_id": "...",
  "target_id": "...",
  "status_code": "berserker_rampage_buff",
  "duration_ms": 20000,
  "stacks": 1
}
```

`combat.status_applied`, `combat.status_expired`, `combat.status_dispelled`, `combat.dot_tick`, `combat.hot_tick`.

## 16. Config

```yaml
status_effects:
  cc_diminishing_reset_sec: 20
  cc_diminishing_multipliers: [1.0, 0.5, 0.25, 0.0]
  buff_stack_cap_default: 3
  cleanup_worker_interval_sec: 60
```

## 17. Metrics
- `status_applied_total{kind,code}`.
- `status_expired_total`.
- `status_dispelled_total{category}`.
- `cc_dr_hit_total{category,level}`.
- `active_status_count` gauge.

## 18. Cross-Reference
- [Combat Loop](combat-loop.md), [Damage](damage.md), [Healing](healing.md)
- [skill-engine/modifiers.md](../skill-engine/modifiers.md), [skill-engine/passives.md](../skill-engine/passives.md)
- [character/health.md](../character/health.md)

## 19. Bekleyen Kararlar
- Snapshot vs re-calc DoT/HoT (v1 snapshot; Phase 02 opsiyonu).
- CC DR levels final playtest.
- Shield stacking priority pattern.
- Persistent DoT (mid-session save) Phase 02.
- Buff visibility client filter (spam counter Phase 02).
