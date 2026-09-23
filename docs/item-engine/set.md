# Set (Set Bonus System)

## 1. Amaç
Aynı "set"e ait item'lardan N parça giymek tier bonuslar açar. Endgame hedef; class kimliği + build şekillendirici.

## 2. Model

Kaynak: [database/content-templates.md §10](../database/content-templates.md).

```
content.item_set
  id, code (unique)
  display_name_key
  description_key
  art_ref
  tier (v1: 1; Phase 02 set tier progression)
  member_item_template_ids []
  version

content.item_set_bonus
  set_id, pieces_required (2 / 4 / 6)
  effect_json (modifiers + optional trigger)
  display_name_key
  description_key
```

## 3. Örnek Set

```yaml
# content/item_sets/emerald_wardens.yaml
code: emerald_wardens
display_name_key: set.emerald_wardens.name
member_item_template_codes:
  - emerald_wardens_helm
  - emerald_wardens_chest
  - emerald_wardens_gloves
  - emerald_wardens_legs
  - emerald_wardens_boots
  - emerald_wardens_ring

bonuses:
  - pieces_required: 2
    display_name_key: set.emerald_wardens.2p
    modifiers:
      - kind: flat
        target: nature_resist
        amount: 20
  - pieces_required: 4
    display_name_key: set.emerald_wardens.4p
    modifiers:
      - kind: mult_add
        target: damage_pct.nature
        amount: 8
  - pieces_required: 6
    display_name_key: set.emerald_wardens.6p
    triggers:
      - trigger: on_hit_dealt
        params: { chance_pct: 15 }
        cooldown_ms: 6000
        then:
          - kind: buff
            target: self
            buff:
              code: emerald_bloom
              modifiers:
                - { kind: mult_add, target: heal_pct, amount: 30 }
              duration_ms: 8000
```

## 4. Set Bonus Evaluation

Character sheet load / equipment change:
```go
func (e *ItemEngine) SetBonusesFor(equipped []Instance) []engine.Modifier {
    counts := map[uuid.UUID]int{}  // set_id → count
    for _, inst := range equipped {
        tmpl := e.templateReg.ByID(inst.TemplateID)
        if tmpl.SetID != nil {
            counts[*tmpl.SetID]++
        }
    }

    mods := []engine.Modifier{}
    for setID, count := range counts {
        set := e.setReg.ByID(setID)
        for _, bonus := range set.Bonuses {
            if count < bonus.PiecesRequired { continue }
            // apply modifiers
            for _, m := range bonus.Modifiers {
                copy := m
                copy.SourceKind = "item_set_bonus"
                copy.SourceID = setID
                mods = append(mods, copy)
            }
            // register triggers separately (PassiveEngine)
            for _, t := range bonus.Triggers {
                e.passives.RegisterTemporary(character, t, setID)
            }
        }
    }
    return mods
}
```

## 5. Piece Counting Rules

- Aynı slot iki item aynı set → 2 sayılmaz; sadece equipped iken 1 sayılır.
- Ring 1 + Ring 2 aynı setten → 2 sayılır (farklı slot).
- Bound / transmog etkilemez (set membership template.set_id ile).

## 6. Tiered Bonuses

Bonus tier'ları kesin (v1):
- 2 piece: konfor / minor bonus.
- 4 piece: sinergi / meta bonus.
- 6 piece: signature / trigger'lı big effect.

Set üyesi < 2 → hiç bonus.

## 7. Signature Trigger

6-piece bonuses genelde trigger-based (proc). PassiveEngine ile entegre (bkz. [skill-engine/passives.md](../skill-engine/passives.md)):
- Set bonus register → temporary passive.
- Karakter equipped iken active; unequip → deregister.

## 8. Endpoint

**GET `/v1/content/sets/{code}`**
Public.

Response:
```json
{
  "code": "emerald_wardens",
  "display_name_key": "set.emerald_wardens.name",
  "members": [
    { "code": "emerald_wardens_helm", "slot": "head" },
    { "code": "emerald_wardens_chest", "slot": "chest" },
    ...
  ],
  "bonuses": [
    { "pieces_required": 2, "display_name_key": "set.emerald_wardens.2p", "summary": "+20 nature resist" },
    ...
  ]
}
```

## 9. Character Sheet Integration

Character sheet response:
```json
{
  "equipment": {...},
  "active_sets": [
    { "set_code": "emerald_wardens", "pieces_equipped": 4, "active_bonuses": ["2p", "4p"] }
  ]
}
```

Tooltip: item detail'da "Set: Emerald Wardens (4/6)" line.

## 10. Content Loader Validation
- `member_item_template_codes` all resolve.
- `pieces_required` in [2, 4, 6] (v1 kesin; Phase 02 esneklik).
- `pieces_required` ascending.
- Modifier + trigger DSL valid.

## 11. Set Deprecation
- `deprecated_at` ile set gizlenir (loot yok).
- Mevcut item instance'lar hala bonusları alır (backward compat).

## 12. Cross-Set Interaction
- İki set aynı anda active olabilir (ör. 4/6 emerald + 2/6 crimson).
- Modifier stack additive.

## 13. Transmog + Set (Phase 02)
- Transmog görsel override; set membership stat item'a bağlı.
- "Görsel farklı, set aktif" mümkün.

## 14. Fail Modları
- Set item template reference invalid → build fail.
- Duplicate bonus tier → build fail.
- Missing bonus for defined pieces_required → warn.

## 15. Metrics
- `set_bonus_active_total{set,pieces}`.
- `set_bonus_trigger_fired_total{set,tier}`.

## 16. Test Stratejisi
- 2p / 4p / 6p threshold enforcement.
- Ring1+Ring2 same set counts.
- Set + trigger register/deregister on equip change.

## 17. Cross-References
- [Item Template](item-template.md), [Item Instance](item-instance.md), [Stats](stats.md), [Unique Effect](unique-effect.md)
- [skill-engine/modifiers.md](../skill-engine/modifiers.md), [skill-engine/passives.md](../skill-engine/passives.md)
- [database/content-templates.md §10](../database/content-templates.md)
- [game-design/equipment.md](../game-design/equipment.md)

## 18. Bekleyen Kararlar
- Set tier progression (Phase 02) — same set 2 tier level upgrade.
- Cross-set anti-synergy tags (Phase 02).
- Set trigger scope (single target vs aoe).
- Half-set replacement (kısmi replace, transition period bonus).
