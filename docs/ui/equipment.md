# Equipment (Paper Doll)

## 1. Amaç
Karakter üzerine kuşanılmış itemleri gösterme, tap-to-equip/unequip, set/gear score özet, karşılaştırma (compare tooltip).

## 2. Layout (Mobile)

```
┌────────────────────────────┐
│ ← Kuşam · Gear Score 1250 │
├────────────────────────────┤
│         [helmet]           │
│    [neck]     [shoulder]   │
│  [chest]        [back]     │
│    [wrist]     [gloves]    │
│  [belt]         [legs]     │
│    [boots]      [ring1]    │
│  [ring2]      [trinket1]   │
│    [trinket2]              │
│  [mainhand]   [offhand]    │
├────────────────────────────┤
│ Set: Iron Battlegear 3/5   │
│ 2-set: +5% dmg (aktif)     │
├────────────────────────────┤
│ [Envanterden Kuşan] btn    │
└────────────────────────────┘
```

Silhouette illustration üzerinde slot ikonları; boş slot grey.

## 3. Bileşenler

- **Header**: Gear Score (item level ortalama × slot ağırlık).
- **Slot grid** on silhouette (12+ slot, class'a göre bazıları disabled).
- **Set indicator**: aktif set bonuslar (yeşil check).
- **Bottom CTA**: envanterden hızlı kuşan modal.

## 4. Actions

- **Tap slot** (dolu) → item detail sheet (bkz. [inventory.md](inventory.md)) + "Çıkar" action.
- **Tap slot** (boş) → envanterden uygun item picker (filtered by slot + class + level).
- **Long-press** → compare tooltip (envanterde başka uygun item ile side-by-side stat delta).

## 5. Compare Tooltip

```
Aktif: Iron Sword           Alt:    Steel Sword
Dmg 12-18                   Dmg 15-22  (+3)
Str +3                      Str +5  (+2)
Crit +2%                    Crit +0%  (-2%)
```

Renk kodları: yeşil (+), kırmızı (-).

## 6. State

- Equipment provider — `GET /v1/character/{id}/equipment`.
- WS `equipment:{id}` events (equipped, unequipped, set_bonus_change).
- Gear score cache invalidate on equipment change.

## 7. API

- `GET /v1/character/{id}/equipment`
- `POST /v1/character/{id}/equip` body: `{ slot, item_instance_id }`
- `POST /v1/character/{id}/unequip` body: `{ slot }`

Bkz. [inventory/equipment.md](../inventory/equipment.md).

## 8. Metrikler
- `ui_equipment_view_total`
- `ui_equipment_equip_total{slot}`
- `ui_equipment_compare_open_total`

## 9. Cross-Reference
- [inventory/equipment.md](../inventory/equipment.md), [inventory/slots.md](../inventory/slots.md)
- [item-engine/set.md](../item-engine/set.md)
