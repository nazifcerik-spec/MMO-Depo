# Passives

## 1. Amaç
Sürekli aktif etki + belirli olaylara tepki veren skill'lerin engine tarafında modeli. Race passive, class trait, spec passive, gear proc — hepsi tek soyutlama.

## 2. Passive Türleri

### 2.1 Aura (Sürekli)
Her zaman aktif modifier stack'ine katkı.

Örnek: race `human_versatility` — `xp_gain_pct +5, always`.

Engine kullanımı:
- Character load'ta modifier gather'a dahil.
- Stat compute'da yansır.
- Combat/aktivite yield modifier zincirine katılır.

### 2.2 Trigger-based (Reactive)
Belirli olay olduğunda tetiklenir; kısa süreli buff/damage/heal üretir.

Örnek: warrior `berserker_frenzy` — düşük HP'de +%15 attack speed.

Engine kullanımı:
- Trigger event olduğunda evaluate.
- Effect apply (buff insert).
- Duration expiry.

## 3. Trigger Event Kataloğu

```
on_hit_dealt         — melee/ranged hit landed
on_spell_cast        — any spell cast complete
on_damage_taken      — any damage taken
on_crit_dealt        — crit landed
on_crit_taken        — crit received
on_hp_below_pct      — HP dropped below threshold
on_hp_above_pct      — HP rose above threshold
on_resource_full     — resource pool full
on_resource_empty    — resource pool drained
on_kill              — target killed
on_ally_killed       — party member died
on_buff_applied      — specific buff applied
on_debuff_applied    — specific debuff applied
on_dodge             — evasion success
on_block             — block success
on_activity_tick     — AFK aktivite tick (gathering yield vs)
on_gather_success    — resource node gather
on_craft_success     — crafting complete
```

## 4. Passive Definition

```json
{
  "code": "berserker_frenzy",
  "type": "passive",
  "class_code": "warrior",
  "spec_code": "berserker",
  "effects": [
    {
      "kind": "trigger",
      "trigger": "on_hp_below_pct",
      "params": { "threshold_pct": 30 },
      "cooldown_ms": 15000,
      "then": [
        {
          "kind": "buff",
          "target": "self",
          "buff": {
            "code": "frenzy_active",
            "modifiers": [
              { "kind": "mult_add", "target": "haste_pct", "amount": 15 }
            ],
            "duration_ms": 15000
          }
        }
      ]
    }
  ]
}
```

## 5. Trigger Engine

```go
type TriggerEngine interface {
    // Domain event geldiğinde çağırılır.
    OnEvent(ctx context.Context, char Character, event EngineEvent)
}
```

Implementation:
```go
func (t *triggerEngine) OnEvent(ctx context.Context, char Character, ev EngineEvent) {
    passives := t.registry.PassivesFor(char)
    for _, p := range passives {
        for _, eff := range p.Effects {
            if eff.Kind != "trigger" { continue }
            if !matchTrigger(eff.Trigger, eff.Params, ev) { continue }
            // internal cooldown check
            if !t.icd.Ready(ctx, char.ID, p.ID) { continue }
            t.icd.Set(ctx, char.ID, p.ID, eff.CooldownMs)
            for _, then := range eff.Then {
                _ = t.effects.Apply(ctx, char, then)
            }
        }
    }
}
```

Internal Cooldown (ICD): Redis short-lived key `icd:<char>:<passive_id>`; TTL cooldown_ms.

## 6. Priority / Order
- Passive'ler content pack'te tanımlanan sırayla evaluate.
- Multiple passives aynı event'e tepki verebilir; her biri kendi ICD ile.

## 7. Passive Stack
Aynı passive tekrar öğrenilirse (edge case) rank artışı (Phase 02); v1'de duplicate learn no-op.

## 8. Aura Passive Modifier'ları

Aura effects direct modifier oluşturur; modifier gather'da kaynak `SourceKind=passive`, `SourceID=passive_skill_id`.

Cache invalidate: passive learn/forget → stat cache invalidate.

## 9. Cross-Character Aura (Party Aura)

Bazı passive'ler party-wide etki verir (Cleric aura, Warden aura Phase 02):
- Passive effect `target: party_members`.
- Party join → apply modifier to each member.
- Party leave → remove modifier.
- Range check (opsiyonel; MMO ölçek zayıflatır).

v1'de basit: party aura karakterin passive'i olarak duruyor; join sırasında ally'ye buff insert.

## 10. Gear Proc

Gear affix'i "on hit chance to X" pattern — passive olarak modelleniyor:
```json
{
  "kind": "trigger",
  "trigger": "on_hit_dealt",
  "params": { "chance_pct": 15 },
  "cooldown_ms": 6000,
  "then": [
    { "kind": "damage", "target": "enemy_target", "element": "fire", "base": 200 }
  ]
}
```

Item template `unique_effect_id` bu formatı taşır; equip sırasında passive registry'e ekle.

## 11. Combat Event Dispatch

Combat tick (Phase 10) → engine.Emit(event) → TriggerEngine dispatch.

Örnek olay:
```go
combat.Emit(EngineEvent{
    Kind: "on_hit_dealt",
    Actor: attackerCharID,
    Target: defenderID,
    Data: map[string]any{
        "damage": 152,
        "crit": true,
    },
})
```

## 12. AFK Sim Uyumu

AFK reveal içinde:
- Simulate tick engine.Emit çağrıları da yapar.
- Deterministic seed altında proc RNG.
- Reveal log'a passive proc event'leri özet olarak.

## 13. Test Stratejisi
- Table-driven trigger matcher.
- ICD respect.
- Aura apply / remove.
- Proc chance distribution (seed'li 10k iteration).

## 14. Metrics
- `passive_trigger_fired_total{code}`.
- `passive_trigger_icd_blocked_total{code}`.
- `passive_aura_active_count` gauge.

## 15. Fail Modları
- Trigger param invalid → skip + warn.
- ICD Redis unreachable → fallback fail-closed (proc miss); rare.
- Circular trigger (passive triggers itself infinitely) → cycle detector, max depth 3.

## 16. Cross-References
- [Skills](skills.md), [Active Skills](active-skills.md), [Modifiers](modifiers.md), [Cooldown](cooldown.md), [Stats](stats.md)
- [race/passive.md](../race/passive.md), [class/specialization.md](../class/specialization.md)
- [game-design/combat.md](../game-design/combat.md)

## 17. Bekleyen Kararlar
- Trigger condition compound expression ("hp<30 AND target_hp<50") — Phase 02 DSL.
- Aura range enforcement (Phase 02).
- Passive rank scaling (Phase 02).
- Cross-passive interaction (2 passive aynı event → sıralama) — CI validation.
