# Equipment (Travel Speed Contribution)

## 1. Amaç
Boots + accessory + mount barding gibi ekipmanların travel hızına katkısı.

## 2. Etki Kaynakları

| Slot | Örnek | Bonus |
|---|---|---|
| **Feet** (boots) | "Boots of Swiftness" | +5% movement |
| **Cape** (back) | "Traveler's Cloak" | +3% movement |
| **Ring/Trinket** | "Nomad's Ring" | +2% movement |
| **Mount Barding** (Phase 02) | "Steel Barding" | +5% mount speed |
| **Consumable** | "Speed Potion" | +15% movement for 5 min |
| **Enchant** | "Swiftness Enchant" | +5% on boots |
| **Set Bonus** | "Wanderer's Set 4p" | +10% movement |

## 3. Equipment Multiplier

Item stat contribution (bkz. [item-engine/stats.md](../item-engine/stats.md)):
- Item `movement_speed_pct` modifier stack'e katılır.
- Compute:
  ```
  equipment_speed_bonus = Σ movement_speed_pct modifiers from equipped items
  equipment_multiplier = 1.0 + equipment_speed_bonus / 100.0
  ```

## 4. Buff / Consumable

Speed potion:
- Consume → 5 dk buff (+%15).
- Modifier `movement_speed_pct +15 scope:out_of_combat` (opsiyonel).
- Duration mid-travel biterse mekanik snapshot'a bağlı:
  - **v1**: snapshot; buff aktifken travel start ise tam buff süre boyunca (contract snapshot content).
  - Alternatif: buff bitince yavaşla (dynamic). Karar: **snapshot** v1.

## 5. Cap Discipline

- Toplam equipment speed bonus cap **+20%** (v1).
- Aşağıda: mount + race + equipment kombine cap 3.5× (bkz. [base-time.md §4](base-time.md)).

## 6. Mount Barding (Phase 02)

Mount barding item:
- Slot: `mount_saddle` (Phase 02 yeni slot).
- Stat: `mount_speed_bonus_pct`.
- Mount speed_multiplier ile kombine:
  ```
  effective_mount_mult = mount.speed_multiplier × (1 + barding_pct/100)
  ```

v1 skip; schema hazır (bkz. [database/companions.md §4](../database/companions.md)).

## 7. Encumbrance (Phase 02 Weight)

Weight sistemi aktive ederse (bkz. [inventory/weight.md](../inventory/weight.md)):
- Overloaded karakter movement -50%.
- Encumbered -10%.
- Cap: negative equipment_multiplier durumu.

v1 kapalı.

## 8. Snapshot Dahil

Travel contract snapshot'ında:
```json
{
  "equipment_speed_bonus_pct": 12,
  "active_speed_buff_codes": ["speed_potion_minor"]
}
```

Snapshot al; travel süresince equipment değişse bile duration korunur.

## 9. Compute Örneği

Elf ranger, mount tier 2 (2×), boots +5%, cape +3%, speed potion +15%:
```
race_mult = 1.05
mount_mult = 2.0
equipment_bonus_pct = 5 + 3 + 15 = 23
equipment_mult = 1.23

combined_boost = 1.05 × 2.0 × 1.23 = 2.583

500 unit grass_road (terrain 1.0, road_bonus 1.15):
duration = (500 × 1.0 × 1000) / (1.0 × 2.583 × 1.15)
         = 500000 / 2.970 = 168358 ms ≈ 168 sec (~2.8 min)
```

## 10. Fail Modları
- Equipment removed mid-travel: contract snapshot → OK (v1); Phase 02 dynamic uyarı.
- Cap exceed: silent clamp.

## 11. Cross-Reference
- [Travel](travel.md), [Base Time](base-time.md), [Mount](mount.md)
- [item-engine/stats.md](../item-engine/stats.md), [item-engine/enchantments.md](../item-engine/enchantments.md), [item-engine/set.md](../item-engine/set.md)
- [inventory/equipment.md](../inventory/equipment.md)
- [character/stats.md](../character/stats.md) — movement_speed_pct

## 12. Metrics
- `travel_equipment_bonus_avg` gauge.
- `travel_speed_potion_active_total` counter.

## 13. Bekleyen Kararlar
- Mount barding Phase 02 tam.
- Weight/encumbrance Phase 02 aktivasyon.
- Speed potion cooldown / stack rules.
- Legendary "Wanderer's Set" content Phase 02.
