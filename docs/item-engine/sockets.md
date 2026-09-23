# Sockets

## 1. Amaç
Item'ın gem yuvası — gem takarak stat/effect eklemek. Rarity'ye bağlı slot sayısı; slot type opsiyonel filter.

## 2. Model

Instance `sockets` JSONB (bkz. [database/content-templates.md §9](../database/content-templates.md)):
```json
{
  "sockets": [
    { "slot": 0, "type": "red", "gem_instance_id": "gi-abc" },
    { "slot": 1, "type": "any", "gem_instance_id": null }
  ]
}
```

- `slot` index.
- `type` opsiyonel color filter (red/blue/yellow/any) — hangi gem takılabilir.
- `gem_instance_id` — null = boş.

## 3. Socket Slot Sayısı

Rarity default:
- common/uncommon/rare: 0.
- epic: 0-1 (roll).
- legendary: 1-2 (roll).
- mythic: 2 (guaranteed).

Prefix "Runic" +1 socket (bkz. [prefix.md §4](prefix.md)).

Cap per item: **3** slot.

## 4. Socket Type (Phase 02 detay)

v1: tek tip "any" (herhangi gem takılabilir).

Phase 02 color coding:
- Red slot → red gem (STR, physical dmg).
- Blue slot → blue gem (INT, spell power).
- Yellow slot → yellow gem (DEX, crit).
- Any slot → herhangi.

Off-color penalty (Phase 02): color mismatch → gem effect %50.

v1 tercih: **type=any tümü**. Renk sistemi Phase 02 opsiyonel.

## 5. Install Gem

**POST `/v1/items/{instance_id}/sockets/{slot}/install-gem`**
Auth + ownership.

Request:
```json
{ "gem_instance_id": "gi-..." }
```

Server:
1. Load item instance FOR UPDATE.
2. Validate:
   - Slot exists.
   - Slot boş.
   - Gem owner = current character.
   - Type match (v1: skip).
   - Item bound_to == current char (opsiyonel; boş item bound olmayabilir).
3. Transaction:
   - Set `sockets[slot].gem_instance_id`.
   - Transfer gem instance owner → `void` (consume; gem burnt into item).
   - Alternatif: gem persist ownership `socketed_in_item`; sonradan çıkarılabilir. Karar: **socketed pattern** (removable).
4. Recompute item stat contribution.
5. Invalidate character stat cache (item equipped ise).
6. Emit event `item.gem_installed`.
7. Response 200: updated instance.

## 6. Remove Gem

**POST `/v1/items/{instance_id}/sockets/{slot}/remove-gem`**
Auth + ownership.

Server:
1. Load item.
2. Validate: slot has gem.
3. Cost: **removal cost** — small gold sink (Phase 02'de removal token opsiyonel).
   - v1: 500 copper flat.
   - Alternative: rarity-based tier.
4. Transaction:
   - Set `sockets[slot].gem_instance_id = null`.
   - Gem instance owner → character.inventory (return).
   - Debit removal cost.
5. Invalidate cache.
6. Emit event.
7. Response.

## 7. Socket Add (Post-Craft Consumable, Phase 02)

Rare consumable "Socket Prism":
- POST `/v1/items/{instance_id}/add-socket`.
- Cost: consume prism.
- Result: sockets array +1 slot (up to cap).
- Fail risk: consume prism, +1 socket %90; %10 fail (prism yanır).

## 8. Character Sheet UI

Equipped item tooltip:
```
Steel Longsword (rare)
  Weapon Damage: 52
  +5 Strength [prefix: Steel]
  +12 Stamina [suffix: of the Bear]
  Sockets:
   [Ruby +3 STR]
   [empty]
  Enchant: Sharpness II (+12 AP)
  Durability: 87/100
```

## 9. Fail Modları
- Slot invalid → 400.
- Gem not owned → 403.
- Slot occupied → 409 (remove first).
- Type mismatch → 400 (v1: skip).
- Cache stale after install → next read refresh.

## 10. Concurrency
- Optimistic lock on item_instance.version.
- Concurrent install → biri success, biri 409 retry.

## 11. Metrics
- `item_gem_installed_total`.
- `item_gem_removed_total`.
- `item_socket_avg_filled_ratio` gauge.

## 12. Audit
- `item.gem_installed` (source=install).
- `item.gem_removed` (source=remove).

## 13. Test Stratejisi
- Install → stat contribution.
- Remove → gem back to inventory.
- Cost debit.
- Ownership validate.

## 14. Cross-References
- [Item Instance](item-instance.md), [Gems](gems.md), [Rarity](rarity.md), [Stats](stats.md), [Prefix](prefix.md)
- [database/content-templates.md §9](../database/content-templates.md)

## 15. Bekleyen Kararlar
- Removal cost final (500 copper vs rarity-tier).
- Color coding v1 vs Phase 02.
- Socket cap 3 vs higher.
- Off-color penalty (Phase 02) formul.
