# Inventory

## 1. Amaç
Karakterin taşıdığı itemleri görüntüleme, taşıma, sort, filter, kullan, sil, satış (vendor/market), transfer (bank/mail).

## 2. Layout (Mobile)

```
┌────────────────────────────┐
│ ← Envanter          🔍 ⚙  │
│ Slot 42/60 · Ağırlık 18/25│
├────────────────────────────┤
│ [Tümü][Silah][Zırh][Mat][Con]│
├────────────────────────────┤
│ ┌──┐┌──┐┌──┐┌──┐┌──┐┌──┐   │
│ │⚔ ││🛡││🧪││⛰ ││🌿││📜│... │  ← 6-col grid
│ │un││ra││x3││x8││x5││x2│    │
│ └──┘└──┘└──┘└──┘└──┘└──┘   │
│  ...                       │
├────────────────────────────┤
│ [Sırala ▼] [Filtre ▼]      │
│ Toplu Sat  |  Toplu Sil    │
└────────────────────────────┘
```

Slot: 64×64 icon + rarity border + count badge + tap → detail sheet.

## 3. Bileşenler

- **Header**: slot count / cap, weight, search.
- **Category tabs**: All / Weapon / Armor / Material / Consumable / Misc / Recipe.
- **Grid**: 6-col mobile (adaptive). Empty slots grey.
- **Sort dropdown**: recent / name / rarity / value / weight.
- **Filter dropdown**: rarity, level range, bind, tradable, favorite.
- **Batch mode**: multi-select checkbox.
- **Item detail sheet** (bottom sheet): icon + name + rarity + stats + affix + socket/gem + bind + actions (equip, use, sell, list, mail, bank, split stack, delete).

## 4. Item Detail Sheet

```
┌────────────────────────────┐
│  [icon]  Iron Sword         │
│  Uncommon · Level 10 · BoE  │
├────────────────────────────┤
│  ⚔ Dmg 12-18  ★ Crit +2%   │
│  Prefix: Sharp              │
│  Suffix: of the Bear (+3 STR)│
│  Socket: 1 empty            │
├────────────────────────────┤
│  Ağırlık 1.5 kg             │
│  Değer ~450c (vendor)       │
├────────────────────────────┤
│ [Kuşan] [Sat] [Listele]     │
│ [Mail] [Bank] [Sil]         │
└────────────────────────────┘
```

## 5. State

- Inventory provider — `GET /v1/character/{id}/inventory` + WS `inventory:{id}` events (item_added, item_removed, item_moved).
- Optimistic UI on drag/drop with rollback on error.
- Cached snapshot for offline first-paint.

## 6. Drag & Drop (Tablet/Desktop)

- Long-press mobile → drag start.
- Drop target: another slot (swap), equipment slot (equip), bank icon, mail icon, delete zone.
- Drop feedback: highlight + vibration (mobile).

## 7. Actions & API

| Action | API |
|---|---|
| Move slot | `PUT /v1/inventory/move` |
| Split stack | `POST /v1/inventory/split` |
| Equip | `POST /v1/character/{id}/equip` |
| Sell vendor | `POST /v1/vendor/{code}/sell` |
| List market | `POST /v1/market/listing` |
| Attach mail | `POST /v1/mail` (draft) |
| Deposit bank | `POST /v1/bank/deposit` |
| Delete | `DELETE /v1/inventory/{item_instance_id}` (confirm modal) |

## 8. Weight ve Slot Cap

- Weight overflow → hareket hızı %20 azalır (Phase 02 debuff); UI warning bar üstte.
- Slot cap dolu → yeni item pickup queue (mail overflow fallback Phase 02).

## 9. Metrikler
- `ui_inventory_view_total`
- `ui_inventory_move_total`
- `ui_inventory_batch_action_total{action}`

## 10. Cross-Reference
- [inventory/README.md](../inventory/README.md), [inventory/slots.md](../inventory/slots.md), [inventory/weight.md](../inventory/weight.md)
- [item-engine/README.md](../item-engine/README.md)
