# Equipment

## 1. Amaç
Karakter üzerinde takılı item'ları modellemek + stat/effect contribution + equip/unequip flow.

Yüksek seviye: [game-design/equipment.md](../game-design/equipment.md).

## 2. Model

Kaynak: [database/character.md §6](../database/character.md).

```
character.equipment
  character_id UUID (PK-part)
  slot character.equipment_slot (PK-part)
  item_instance_id UUID NOT NULL
  equipped_at TIMESTAMPTZ
```

`equipment_slot` enum (bkz. [database/character.md §6](../database/character.md)):
```
head, neck, shoulder, back, chest, wrist, hands, waist, legs, feet,
ring1, ring2, trinket1, trinket2,
main_hand, off_hand, ranged
```

## 3. Slot Compatibility

Item equip için `content.item_template.equipment_slot` slot'la eşleşmeli.

Special cases:
- **Ring**: template `equipment_slot='ring'`; UI ring1 veya ring2 slot'una takar.
- **Trinket**: benzer.
- **2H weapon**: template `weapon_kind` 2H; main_hand'e takılır; off_hand'i auto-unequip.
- **Dual wield**: class + spec bağlı; main_hand + off_hand ayrı 1H weapon.

## 4. Requirements (Equip Check)

Bkz. [item-engine/item-template.md §5](../item-engine/item-template.md):
- `min_level`, `class_codes`, `spec_codes`, `race_codes`, `stat_mins`.
- Legendary+ cap (2 max) — bkz. [item-engine/unique-effect.md §6](../item-engine/unique-effect.md).

## 5. Equip Endpoint

**POST `/v1/characters/{id}/equipment`**
Auth + CharacterScope.

Request:
```json
{
  "item_instance_id": "ii-...",
  "slot": "main_hand"   // opsiyonel; template'ten inferred
}
```

Server:
1. Load character FOR UPDATE.
2. Load item_instance FOR UPDATE.
3. Validate ownership (character_inventory'de).
4. Resolve target slot (client hint vs template default).
5. Validate:
   - Slot compatibility.
   - Requirements.
   - Legendary cap.
   - Item not locked.
   - Item not bound-to-other.
6. Transaction:
   - Current item (if any) in target slot → move to source inventory slot (from where new item came) or first empty inventory slot.
   - 2H equip → off_hand auto-unequip (mirror move).
   - BoE trigger: `bound_to_character_id = character.id`.
   - `character.equipment` upsert (slot → item_instance_id).
   - `character.inventory_slot` remove new item entry.
   - `character.stats` cache invalidate.
   - `character.character` version bump.
7. Emit event `item.equipped`.
8. Response 200: updated equipment map + stats delta.

## 6. Unequip Endpoint

**DELETE `/v1/characters/{id}/equipment/{slot}`**
Auth + CharacterScope.

Server:
1. Load character + current equipment slot.
2. Validate: slot has item.
3. Validate: has inventory space (1 slot minimum).
4. Transaction:
   - Move item → first empty inventory_slot.
   - Delete character.equipment row.
   - Cache invalidate.
5. Emit event `item.unequipped`.
6. Response 200.

## 7. Swap Loadout (Phase 02)

**POST `/v1/characters/{id}/equipment/loadout-apply`**
```json
{ "loadout_id": "..." }
```

Multi-item swap tek transaction. Combat sırasında yasak (out-of-combat 5s).

## 8. Auto-Equip (Loot Convenience) — Phase 02

Bazı oyuncular yeni loot drop için auto-equip flag açar:
- Yeni item drop → daha iyi (item_score) mi kıyasla → auto-equip.
- v1'de yok; UI-driven manual.

## 9. Off-hand Interaction

Main hand equip sırasında:
- Yeni item **2H** ise: off_hand auto-move to inventory.
  - Confirmation UI: "Off-hand shield will be moved to inventory. Continue?"
- Yeni item **1H** ise + off_hand'de shield varsa: shield tutulur.

## 10. Slot Locking (Item Lock vs Slot Lock)

Bkz. [item-lock.md](item-lock.md).

Item lock:
- Equipped item lock → unequip endpoint reddedilir; önce unlock.
- Bir slotun içeriği "locked" ise equip başka item o slot'a olduğunda swap yapılamaz (koruma).

## 11. Stat Recomputation

Equip/unequip her seferinde:
- `character.stats.derived_cache_invalidated_at = now()`.
- Next stat read → recompute (bkz. [skill-engine/stats.md](../skill-engine/stats.md)).
- Set bonus evaluation da re-run.

## 12. Durability & Combat Loss

Combat sırasında equipped item durability azalır:
- Weapon: her hit -0.1 durability.
- Armor: her hit taken -0.05.
- Death: her equipped item -5% max durability.
- Repair: NPC city veya self-repair kit.

Server-side durability_current tracking; batch updates (per combat end).

## 13. Concurrency

- Optimistic lock per equipment slot (composite unique).
- Concurrent equip → biri success, biri 409.
- Character version her equip/unequip'te artar → skill engine cache invalidate.

## 14. Fail Modları
- `equipment.slot_incompatible`
- `equipment.requirement_not_met`
- `equipment.legendary_cap`
- `equipment.item_locked`
- `equipment.item_bound_to_other`
- `equipment.inventory_full` (unequip için hedef slot yok)
- `equipment.in_combat` (bazı equip aksiyonları combat'ta yasak — Phase 02)

## 15. Rate Limit
- Equip/unequip: 30/dk (gerçek use case küçük).
- Loadout swap: 5/dk.

## 16. Metrics
- `equipment_equip_total{slot,rarity}`.
- `equipment_unequip_total{slot}`.
- `equipment_binding_triggered_total`.
- `equipment_2h_auto_unequip_offhand_total`.

## 17. Audit
- `equipment.equipped`, `equipment.unequipped`.

## 18. Event Payload
```json
{
  "_v": 1,
  "character_id": "c-...",
  "slot": "main_hand",
  "item_instance_id": "ii-...",
  "action": "equipped",
  "stat_delta": { "str": 5, "attack_power": 22 }
}
```

Subscribers:
- `set_bonus.reevaluate` (bkz. [item-engine/set.md](../item-engine/set.md)).
- `passive_engine.reregister` (unique effect + set trigger).
- `skill_engine.stat_cache_invalidate`.

## 19. Test Stratejisi
- Equip happy path + stat delta.
- 2H unequip off_hand.
- Requirement fail matrix.
- Legendary cap enforce.
- BoE binding trigger.
- Loadout swap Phase 02.

## 20. Cross-References
- [Inventory](inventory.md), [Item Lock](item-lock.md), [Slots](slots.md), [Item Transfer](item-transfer.md)
- [database/character.md §6](../database/character.md)
- [item-engine/item-template.md](../item-engine/item-template.md), [item-engine/set.md](../item-engine/set.md), [item-engine/unique-effect.md](../item-engine/unique-effect.md)
- [skill-engine/stats.md](../skill-engine/stats.md)
- [game-design/equipment.md](../game-design/equipment.md)

## 21. Bekleyen Kararlar
- Loadout Phase 02 tam schema.
- Auto-equip smart loot Phase 02.
- Combat-locked equip (Phase 02 opsiyonu).
- Transmog Phase 02 (görsel override slot).
