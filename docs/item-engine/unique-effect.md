# Unique Effect (Legendary+ Special Mechanics)

## 1. Amaç
Legendary / Mythic item'ların stat + affix'in ötesinde **karakter davranışını değiştiren** özgün mekaniği. Handler-registered (data + code hybrid).

## 2. Model

Item template:
- `unique_effect_id UUID NULL FK → content.unique_effect`.

`content.unique_effect` tablosu (Phase 10 yeni):
```sql
CREATE TABLE content.unique_effect (
  id UUID PRIMARY KEY,
  code VARCHAR(48) UNIQUE NOT NULL,
  handler_code VARCHAR(48) NOT NULL,       -- Go tarafında handler map key
  display_name_key VARCHAR(96) NOT NULL,
  description_key VARCHAR(160) NOT NULL,
  params_json JSONB NOT NULL DEFAULT '{}',
  version INT NOT NULL DEFAULT 1
);
```

## 3. Handler Registry (Go)

```go
package uniqueeffect

type Handler interface {
    // Item equipped → aktive
    OnEquip(ctx context.Context, char Character, item Instance, params map[string]any) error
    // Item unequipped → deaktif
    OnUnequip(ctx context.Context, char Character, item Instance) error
    // Combat / activity event
    OnEvent(ctx context.Context, char Character, item Instance, event EngineEvent, params map[string]any)
    // Modifier contribution (aura-benzeri)
    Modifiers(char Character, item Instance, params map[string]any) []engine.Modifier
}

var handlers = map[string]Handler{
    "phoenix_rebirth":       phoenixRebirthHandler{},
    "vampiric_touch":        vampiricTouchHandler{},
    "chronomancer_echo":     chronomancerEchoHandler{},
    "titan_grip":            titanGripHandler{},
    "wanderers_step":        wanderersStepHandler{},
    "beast_speaker":         beastSpeakerHandler{},
}
```

Loader validation: unique_effect.handler_code must exist in handlers map; else build fail.

## 4. Örnek Unique Effect

### 4.1 Phoenix Rebirth
- **Item**: legendary chest "Phoenix Ashes Robes".
- **Effect**: karakter HP 0'a düşerse (öl olurken), 60% max HP ile 3 sn immune ile revive. Cooldown 10 dk.
- **Handler**: `phoenix_rebirth`.
- **Params**: `{ "revive_hp_pct": 60, "immune_sec": 3, "cooldown_sec": 600 }`.

Implementation:
```go
type phoenixRebirthHandler struct{}

func (h phoenixRebirthHandler) OnEvent(ctx context.Context, char Character, item Instance, ev EngineEvent, params map[string]any) {
    if ev.Kind != "on_lethal_damage" { return }
    icdKey := "phoenix:" + char.ID.String()
    if !e.cache.Available(ctx, icdKey) { return }  // cooldown
    // Prevent death
    e.combat.Revive(ctx, char.ID, params["revive_hp_pct"].(float64))
    e.combat.ApplyBuff(ctx, char.ID, "phoenix_immune", params["immune_sec"].(int))
    e.cache.SetCooldown(ctx, icdKey, params["cooldown_sec"].(int))
    e.bus.Publish(ctx, bus.Event{Kind: "unique_effect.triggered", Payload: {char, item, effect: "phoenix_rebirth"}})
}
```

### 4.2 Vampiric Touch
- **Item**: legendary ring "Vampire's Signet".
- **Effect**: her damage dealt %10 lifesteal + kills give +100 HP for 15 sn.
- **Handler**: `vampiric_touch`.

Implementation: OnEvent(on_hit_dealt) → damage_returned = damage * 0.10 → heal caster.

### 4.3 Titan Grip
- **Item**: mythic weapon "Titan's Fist".
- **Effect**: 2-hand weapon'ları tek elde tutabilir; ancak attack speed -%20.
- **Handler**: `titan_grip`.
- **Modifiers**: static aura modifier `{kind: mult_add, target: attack_speed_pct, amount: -20}`. Ayrıca equip constraint bypass.

## 5. Legendary Progression (Phase 02+)

Bazı legendary'ler zaman içinde güçlenir:
- Unique effect params item stats leveled up ile scale.
- "Grow with your character" model.
- Legendary progression sistem: [game-design/ascension.md](../game-design/ascension.md) parallel'i.

v1'de: params static; Phase 02'de progression alanı `character.legendary_progress` (schema hazır).

## 6. Sadece Bir Unique Effect Item

Cap: karakter aynı anda **max 2 legendary+ item** kuşanabilir (soft cap). Nedeni: unique effect combo aşırı OP olmasın.

Enforcement:
```go
func canEquip(char Character, tmpl TemplateDef) error {
    if tmpl.Rarity >= "legendary" {
        equippedLegendariesCount := countEquippedLegendaries(char)
        if equippedLegendariesCount >= 2 { return errs.ErrForbidden }
    }
    ...
}
```

Cap konfigüre.

## 7. Modifier Contribution

Statik aura modifier'lar:
```go
mods := handler.Modifiers(char, item, unique_effect.params)
// bunlar ItemEngine.ModifiersFor sonucuna eklenir (bkz. stats.md §5).
```

## 8. Trigger Effects

Reactive unique effect'ler PassiveEngine ile beraber çalışır (bkz. [skill-engine/passives.md](../skill-engine/passives.md)):
- OnEquip: register passive trigger.
- OnUnequip: deregister.

## 9. Endpoint

**GET `/v1/content/unique-effects/{code}`**
Public.
```json
{
  "code": "phoenix_rebirth",
  "handler_code": "phoenix_rebirth",
  "display_name_key": "unique.phoenix_rebirth.name",
  "description_key": "unique.phoenix_rebirth.desc",
  "params": { "revive_hp_pct": 60, "immune_sec": 3, "cooldown_sec": 600 }
}
```

## 10. Cross-Character Interaction
- Party içinde 2 kişide aynı unique effect var → bağımsız cooldown.
- Handler decides interaction ("aura shares", "no overlap") — content ve handler tasarım.

## 11. Fail Modları
- Handler yok → build fail (loader).
- Handler panic → recover + audit critical.
- Params tipi yanlış → runtime warn + noop.

## 12. Test Stratejisi
- Handler unit test (mock combat/cache).
- OnEquip → OnUnequip cleanup.
- Cooldown enforce.
- Cap enforcement (2 legendary max).

## 13. Metrics
- `unique_effect_active_total{code}`.
- `unique_effect_triggered_total{code}`.
- `unique_effect_cooldown_blocked_total{code}`.

## 14. Audit
- `unique_effect.triggered` (severity=info; rare events).

## 15. Content Update
- Yeni unique effect eklenirken:
  1. Handler PR (Go kod).
  2. YAML content pack.
  3. Handler registry test.
  4. Deploy.

## 16. Cross-References
- [Item Template](item-template.md), [Item Instance](item-instance.md), [Rarity](rarity.md), [Set](set.md)
- [skill-engine/passives.md](../skill-engine/passives.md)
- [race/racial-bonuses.md](../race/racial-bonuses.md) — benzer handler pattern (DD-059)
- [game-design/items.md](../game-design/items.md) — legendary progression vision

## 17. Bekleyen Kararlar
- 2 legendary cap final değer.
- Legendary progression Phase 02 detay.
- Mythic-only handler'lar (progression gated).
- Handler-less pure DSL alternatif (effect DSL yeter mi) — çoğu unique için handler gerekli (custom mekanik).
