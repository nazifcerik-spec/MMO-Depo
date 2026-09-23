# Gems

## 1. Amaç
Socket'e takılan stat/effect kaynağı. Jewelcrafter meslek çıktısı (Phase 02); v1'de drop / vendor.

## 2. Model

Gem = özel bir item tipi. **Gem template** + **gem instance** ayrımı (item pattern).

Kaynak: [database/content-templates.md §13](../database/content-templates.md).

```
content.gem_template
  id, code, rarity, tier (1..5)
  color (red|blue|yellow|any)
  base_effect JSONB (modifiers)
  cut_state (raw|cut|flawless|perfect|imperfect)
  art_ref, display_name_key

content.gem_instance
  id, template_id, owner_id, seed
  cut_state, custom_state_json
  bound_flag (BoP after cut? — opsiyonel)
```

## 3. Cut States

Gem lifecycle:
```
RAW (drop)
  ↓ Jewelcrafter cut
CUT (usable, base effect)
  ↓ (opsiyonel)
FLAWLESS (better effect)
  ↓
PERFECT (best)
  ↓ (fail risk)
IMPERFECT (worse; downgrade)
```

Cut'ların stat multiplier'ı:
- raw: kullanılamaz (socket'e girmez).
- cut: 1.0×.
- flawless: 1.5×.
- perfect: 2.0×.
- imperfect: 0.5×.

## 4. Gem Colors + Effects

| Color | Primary stat | Secondary stat |
|---|---|---|
| **Red** | STR | attack_power, physical_dmg_pct |
| **Blue** | INT | spell_power, mana_regen |
| **Yellow** | DEX | crit_chance_pct, evasion_pct |
| **Any** (rare) | 1 stat of choice (Phase 02) | – |

## 5. Örnek Gem Kataloğu

| Kod | Cut | Color | Tier | Effect |
|---|---|---|---|---|
| `ruby_cut_t1` | cut | red | 1 | +3 STR |
| `ruby_flawless_t1` | flawless | red | 1 | +5 STR |
| `sapphire_cut_t1` | cut | blue | 1 | +3 INT |
| `topaz_cut_t1` | cut | yellow | 1 | +1% crit_chance |
| `ruby_cut_t3` | cut | red | 3 | +10 STR |
| `ruby_perfect_t5` | perfect | red | 5 | +40 STR |

Tier drop rarity'ye bağlı (tier 1 common; tier 5 legendary+).

## 6. Cutting Flow (Phase 02 Jewelcraft)

Raw gem + tool → cut:
```
POST /v1/characters/{id}/craft
  body: { recipe_code: "cut_ruby_cut", inputs: {raw_ruby: 1, jewelcraft_tool: 1} }

Server:
  - Validate profession level.
  - Roll cut quality (based on skill + luck):
     %70 cut, %25 flawless, %4 perfect, %1 imperfect
  - Consume raw gem.
  - Create cut gem instance.
```

v1'de cutting yok; drops directly as "cut" tier.

## 7. Gem Install

Bkz. [sockets.md §5](sockets.md).

## 8. Effect Application

Gem effect item stat contribution katmanına eklenir (bkz. [stats.md §5](stats.md)):
```go
for _, socket := range instance.Sockets {
    if socket.GemInstanceID == nil { continue }
    gem := load(socket.GemInstanceID)
    multiplier := cutStateMultiplier[gem.CutState]
    for _, m := range gem.Effect.Modifiers {
        copy := m
        copy.Amount *= multiplier
        copy.SourceKind = "item_gem"
        mods = append(mods, copy)
    }
}
```

## 9. Ownership

Gem instance:
- Character inventory: available for install.
- Socketed: owner_kind='socketed_in_item', owner_id=item_instance_id.

Alternative schema simplification:
- Gem satırı character inventory'de kalır; item sockets sadece reference.
- Trade-off: dual reference vs. clearer ownership. Karar: **socketed_in_item** kind (net ownership).

## 10. Trade Rules

Raw gem: BoE (tradeable).
Cut gem: BoE (tradeable; still valuable).
Socketed gem: item ile birlikte trade (item bound değilse).

## 11. Endpoint

**GET `/v1/characters/{id}/gems`** — envanterdeki gems.

**POST `/v1/gems/{instance_id}/cut`** (Phase 02) — cutting flow.

## 12. Fail Modları
- Cut gem for non-jewelcrafter → 400.
- Insufficient inputs → 400.
- Item slot color mismatch (Phase 02) → 400.

## 13. Metrics
- `gem_installed_total{color,tier}`.
- `gem_cut_total{result}` (cut/flawless/perfect/imperfect).
- `gem_market_avg_price{color,tier}`.

## 14. Cross-References
- [Sockets](sockets.md), [Item Instance](item-instance.md), [Item Template](item-template.md), [Stats](stats.md), [Prefix](prefix.md)
- [database/content-templates.md §13](../database/content-templates.md)
- [game-design/professions.md](../game-design/professions.md)

## 15. Bekleyen Kararlar
- Color coding v1 skip vs enable.
- Cutting profession v1'de aktif mi (öneri: hayır; Phase 02).
- Gem stacking (aynı gem type stackable? — cut gem tercihen instance).
- Meta gem (Phase 02+): 2+ farklı color gems ile aktive olan set-benzeri effect.
