# Items (Template + Instance Browser)

## 1. Amaç
Item template katalog gözden geçirme (content-driven), item instance forensic (kim yarattı, kim taşıdı, şimdi kim), suspicious pattern tespiti.

## 2. Template Browser

- Content pack'ten yüklenen tüm item template.
- Filter: class, subclass, rarity, level, bind_kind, is_no_sell.
- Detail: full JSON (base stats, affix pool, socket, drop table refs).
- Read-only; content edit **content pipeline** (Phase 02 tool) üzerinden yapılır.

## 3. Instance Browser

- Search: item_instance_id, item_template_code (list of instances).
- Filter: owner_character_id, rarity, obtained_at range.
- Columns: instance_id · template · rarity · current owner · obtained_at · state (locked/listed/mailed).

### Instance Detail (Forensic)
```
Instance: item_ab12…
Template: sword_iron
Rarity: uncommon
Sockets: 1 (empty)
Enchant: none
Bind: BoE

Current owner: Grim (char_…)
Obtained: 2026-08-15 by loot roll (dungeon_run_id: …)

History:
  2026-08-15 loot → Grim (dungeon iron_mine_hardmode)
  2026-08-20 listed on market (listing_id: …)
  2026-08-21 cancelled
  2026-08-22 mailed → Aria
  2026-08-23 claimed by Aria
  2026-08-25 sold to Grim via market

Ledger entries:
  ...

Related audit:
  ...
```

Instance history event stream (`event.item.*`) 12 ay hot retention.

## 4. Suspicious Pattern Detection (Phase 02)

- Bir instance kısa süre içinde çok owner değiştirmişse (`>5 owner in 24h`).
- Dupe candidate: aynı signature (crafter+timestamp) iki instance.
- RMT flag: yeni account 24h içinde high-value item transfer.

Admin UI'da suspicious label + drill-down.

## 5. Aksiyonlar

- **Instance remove** (dupe/exploit): item soft delete + ledger reversal opsiyonel.
- **Instance transfer** (support restoration): owner değişimi + audit.

Her aksiyon `admin` role step-up.

## 6. Endpoint

- `GET /internal/admin/item/template?…`
- `GET /internal/admin/item/instance?…`
- `GET /internal/admin/item/instance/{id}` (history + ledger + audit compound)
- `POST /internal/admin/item/instance/{id}/remove` `{ reason }`
- `POST /internal/admin/item/instance/{id}/transfer` `{ new_owner_character_id, reason }`

## 7. Cross-Reference
- [item-engine](../item-engine/README.md)
- [systems/loot.md](../systems/loot.md)
- [economy/transaction.md](../economy/transaction.md)

## 8. Bekleyen Kararlar
- Suspicious pattern detection Phase 02.
- Item retention full history 12 ay v1; Phase 02 cold arşiv extend.
- Content edit tool ayrı Phase 02 pipeline.
