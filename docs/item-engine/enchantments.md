# Enchantments

## 1. Amaç
Item'a Enchanter meslek çıktısı ile eklenen kalıcı bonus. Enchant slot tekil per item; gem'lerden ayrı katman.

## 2. Model

Item instance:
- `enchant_id UUID NULL FK → content.enchant`.

Enchant template `content.enchant` (bkz. [database/content-templates.md §12](../database/content-templates.md)):
```
id, code, display_name_key,
applicable_slots [] (equipment_slot enum list — hangi item'a takılabilir),
applicable_weapon_kinds [] (opsiyonel filter),
applicable_armor_kinds [] (opsiyonel),
tier (1..5),
effect_json (modifiers)
```

## 3. Örnek Enchant Kataloğu

| Kod | Tier | Slot | Effect |
|---|---|---|---|
| `sharpness_i` | 1 | main_hand, off_hand (weapon) | +8 attack power |
| `sharpness_ii` | 2 | weapon | +12 AP |
| `sharpness_iii` | 3 | weapon | +18 AP |
| `elemental_flame` | 3 | weapon | +N fire dmg per hit |
| `spell_focus_i` | 1 | main_hand (staff) | +5% spell power |
| `warding_i` | 1 | chest | +20 all resist |
| `endurance_i` | 1 | legs, feet | +100 HP |
| `swiftness_ii` | 2 | feet | +8% movement speed |
| `wisdom_iii` | 3 | head | +15 INT + N mana regen |

## 4. Apply Enchant

**POST `/v1/items/{instance_id}/enchant`**
Auth + ownership.

Request:
```json
{
  "enchant_code": "sharpness_ii",
  "consume_scroll_instance_id": "ii-scroll-...",   // enchant scroll consumed
  "at_workbench_location_id": "loc-enchanter-shop"  // opsiyonel; bazı enchant workbench required
}
```

Server:
1. Load item instance FOR UPDATE.
2. Validate:
   - Item ownership.
   - Enchant slot uygulanabilir (`enchant.applicable_slots`).
   - Existing enchant varsa: replace (üstüne yaz) veya remove first (config).
   - Scroll owned + not expired.
   - Workbench location (opsiyonel).
3. Transaction:
   - Consume scroll (owner → void).
   - Set `item_instance.enchant_id = enchant_id`.
   - Increment version.
   - Recompute stat contribution.
   - Invalidate cache.
   - Emit event `item.enchanted`.
4. Response 200.

## 5. Remove Enchant

**POST `/v1/items/{instance_id}/enchant/remove`**
Server:
1. Load item.
2. Cost: gold or "disenchant kit" consumable.
3. Set `enchant_id = null`.
4. Optionally: return residual "essence" material (Phase 02 disenchanting).
5. Recompute.

## 6. Enchant Upgrade Path (Phase 02)

Tier progression: `sharpness_i` → `sharpness_ii` → `sharpness_iii`.

Upgrade endpoint (Phase 02):
- Consume current enchant + materials → upgrade tier.
- Success/fail risk (Phase 02 opsiyonel):
  - Tier 1→2: %100 (safe).
  - Tier 2→3: %70; fail = tier lost, back to 1.
  - Tier 3→4: %40; fail = destroyed (severe).
- Protection scroll (rare): fail'de tier düşer ama destroyed olmaz.

**v1'de**: no upgrade path; enchant tier drop or crafted directly.

## 7. Enchant Scroll (Consumable Item)

Enchanter production:
- Scroll = special item template, category='consumable', subcategory='enchant_scroll'.
- Scroll içinde `enchant_code` alanı; consume ile o enchant apply olur.

Alternatif model: Enchanter workbench'te doğrudan apply (recipe olarak).

v1 tercih: **scroll model** (basit; envanter'de tutulabilir + trade edilebilir).

## 8. Trade Rules

- Enchant scroll BoE.
- Enchanted item: enchant'ı ile birlikte trade (kaybolmaz).
- BoP item enchanted olabilir; trade zaten yok.

## 9. Effect Application

Enchant effect item stat contribution katmanına (bkz. [stats.md §5](stats.md)):
```go
if instance.EnchantID != nil {
    enc := load(instance.EnchantID)
    for _, m := range enc.Effect.Modifiers {
        copy := m
        copy.SourceKind = "item_enchant"
        copy.SourceID = instance.ID
        mods = append(mods, copy)
    }
}
```

## 10. Enchant + Elemental Damage

Bazı enchant'lar element bazlı damage ekler ("Flame Enchant" +N fire dmg per hit):
- Passive-benzeri: on_hit_dealt trigger + damage effect.
- Bkz. [skill-engine/passives.md](../skill-engine/passives.md) — enchant "gear proc" gibi.

Model:
```json
{
  "modifiers": [],
  "triggers": [
    {
      "trigger": "on_hit_dealt",
      "params": { "chance_pct": 100 },
      "cooldown_ms": 0,
      "then": [
        { "kind": "damage", "target": "enemy", "element": "fire", "base": 15 }
      ]
    }
  ]
}
```

## 11. Config

```yaml
enchant:
  removal_cost_copper: 5000
  workbench_required_for_tiers: [3, 4, 5]
  upgrade_enabled: false      # v2 aktive
```

## 12. Fail Modları
- Slot mismatch → 400 `enchant.slot_invalid`.
- Scroll not owned → 403.
- Workbench required but absent → 400.
- Concurrent apply → version conflict retry.

## 13. Metrics
- `enchant_applied_total{code,tier}`.
- `enchant_removed_total`.
- `enchant_upgrade_total{result}` (Phase 02).

## 14. Audit
- `item.enchanted` (source_kind=scroll).
- `item.enchant_removed`.

## 15. Test Stratejisi
- Apply → stat contribution.
- Remove → cache invalidate.
- Scroll consume atomic.
- Workbench requirement enforce.

## 16. Cross-References
- [Item Instance](item-instance.md), [Stats](stats.md), [Sockets](sockets.md), [Gems](gems.md)
- [database/content-templates.md §12](../database/content-templates.md)
- [game-design/professions.md](../game-design/professions.md) — Enchanter
- [skill-engine/passives.md](../skill-engine/passives.md) — trigger enchants

## 17. Bekleyen Kararlar
- Upgrade tier path Phase 02 aktivasyon + fail formula.
- Disenchant/salvage material output (Phase 02 disenchanting profession).
- Scroll model vs direct apply — v1 scroll.
- Enchant slot conflict rules (replace vs manual remove first).
