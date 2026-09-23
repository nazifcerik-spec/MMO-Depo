# Carry Capacity (Mount Cargo Bonus)

## 1. Amaç
Mount'un envanter kapasitesine katkısı — cargo mount slot bonus + barding (Phase 02).

## 2. Cargo Bonus Slot

Bkz. [mount.md §3](mount.md), [stats.md §7](stats.md).

`content.mount_template.cargo_bonus_slots`:
- Ground horse: 0.
- Cargo mount (steel bear): 5-10.
- Aquatic turtle: 0-3.

Effective:
```
inventory_capacity_from_mount = mount.cargo_bonus_slots + barding_bonus (Phase 02)
```

## 3. Inventory Extension

Karakter aktif cargo mount ile:
```
total_inventory_capacity = base(40) + bag_bonus + account.bag_slot_bonus + mount.cargo_bonus_slots
```

Mount summon → slots ekle.
Mount dismount → slots subtract.

Slot loss durumu:
- Envanter dolu iken mount dismount → overflow mailbox (bkz. [inventory/inventory.md §7](../inventory/inventory.md)).

Client uyarı:
- "Cargo mount dismount → 3 items will overflow to mailbox. Confirm?"

## 4. AFK Travel Cargo

Travel active with cargo mount:
- Increased inv capacity during travel.
- If travel completes and mount stays active: capacity retained.
- If mount dismounted post-travel: overflow check.

## 5. Cargo Only, No Combat Contribution

Cargo mount type:
- No combat pet contribution (mount ≠ pet).
- No DPS boost.
- Focus: transport + storage.

Trade-off design: cargo mounts slow but big carrying.

## 6. Barding (Phase 02)

Mount barding = mount equipment:
- Slot: `mount_saddle` (Phase 02 new slot for mount).
- Stat: `cargo_bonus_pct` (opsiyonel `+5 slots`).
- Content pack Phase 02 detay.

v1'de barding yok; base cargo slots only.

## 7. Cargo Overflow

Envanter full + new drop:
- Cargo bonus slots avail → auto-fill.
- Full → mailbox overflow.

## 8. Weight Interaction (Phase 02)

Weight sistemi aktive edilirse (bkz. [inventory/weight.md](../inventory/weight.md)):
- Cargo mount max_weight bonus (Phase 02).

v1'de weight kapalı.

## 9. Metrics
- `mount_cargo_slots_active_total`.
- `mount_cargo_overflow_total`.

## 10. Cross-Reference
- [Mount](mount.md), [Stats](stats.md), [Travel](travel.md)
- [inventory/inventory.md](../inventory/inventory.md), [inventory/slots.md](../inventory/slots.md)
- [database/companions.md §4](../database/companions.md) — mount_barding schema Phase 02

## 11. Bekleyen Kararlar
- Barding Phase 02 tam schema.
- Weight system Phase 02 aktivasyon impact.
- Cargo slots stack cap (v1: 10).
- Multi-mount cargo pool Phase 03 (2 mount aktif).
