# Weight

> ⚠️ **v1'de DEVRE DIŞI** (DD-021).
> Bu doküman resmi tasarımdır; ancak v1 inventory sadece **slot cap** ile çalışır. Weight mekaniği Phase 02+ opsiyonu.

## 1. Amaç
Envanter kapasitesine ikinci bir kısıtlayıcı: item ağırlığı. Slot dolu değil ama "aşırı yüklü" olabilir.

## 2. Neden v1'de Kapalı?

DD-021 gerekçesi:
- **Sadelik**: slot cap tek başına anlaşılır.
- **Idle MMO psikolojisi**: weight loop "envanter yönetme yükü" hissettirir; casual friendly değil.
- **Slot cap yeter**: 40..120 slot mekanizması engagement için yeterli.

Aktive edilirse (Phase 02):
- Slot cap + weight cap kombine kısıtlayıcı.
- Item template `weight INT` zaten schema'da; kullanıma hazır.

## 3. Model (Phase 02 Aktive)

`content.item_template.weight INT` (mevcut).

`character.character.max_weight INT` (Phase 02 schema ekleme):
```sql
ALTER TABLE character.character
  ADD COLUMN max_weight INT NOT NULL DEFAULT 100,
  ADD COLUMN current_weight INT NOT NULL DEFAULT 0;
```

Current weight compute:
```
current_weight = Σ(inventory_slot.stack_count × template.weight)
             + Σ(equipped_item.template.weight)
```

Max weight kaynakları (Phase 02):
- Base: `100` (config).
- STR bonus: `+2 per STR` (bkz. [character/stats.md](../character/stats.md) tertiary).
- Mount cargo bonus: `+50` (kargo mount).
- Race passive (bkz. [race/passive.md](../race/passive.md)): örn. dwarven +%20.
- Buff: consumable "Ox Strength" +%50 30dk.

## 4. Weight Threshold Etkileri (Phase 02)

Weight ratio (`current / max`):

| Ratio | Durum | Etki |
|---|---|---|
| 0 - 60% | Normal | Yok |
| 60 - 90% | Encumbered | -10% movement speed |
| 90 - 100% | Heavily Encumbered | -25% movement + -5% haste |
| 100 - 130% | Overloaded | Movement -50%, can't run, can't mount |
| > 130% | Immobile | Can't move (edge case; prevent'able) |

Server-side check:
- Item pickup → weight overflow → drop refused veya mailbox overflow.
- Manual add: 400 `weight.overload_prevented`.

## 5. Weight UI (Phase 02)

- Weight bar (green/yellow/red).
- Item tooltip'te weight.
- Character sheet total weight display.

## 6. Content Weight Rules

- **Materyaller**: 1 weight per unit (iron ore 1, iron ingot 2, coal 1).
- **Consumables**: 0.5 weight (potion, food).
- **Weapon**: 5-15 (sword 8, staff 6, 2H sword 15).
- **Armor**: 3-25 (cloth 3, plate chest 20).
- **Currency-item**: 0 (weightless).
- **Bag** (Phase 02): 5-10.

Config-driven; balance simülasyon.

## 7. Aktivasyon Kararı

Karar tetikleyicileri:
- Ekonomi ihtiyacı (aşırı material stack).
- Class differentiation (STR class'ları daha çok taşır).
- Class player feedback (Phase 02 playtest).

**Şu an için**: aktivasyon **planı yok**; schema hazır kalıcı.

## 8. Fail Modları (Phase 02)
- `weight.overload_prevented` — item add reddedilir.
- `weight.mount_forbidden` — mount summon overload'da bloklu.

## 9. Cross-Reference
- [Inventory](inventory.md), [Slots](slots.md), [Equipment](equipment.md)
- [character/stats.md §4.3 tertiary carry_capacity](../character/stats.md)
- [database/character.md](../database/character.md)
- [race/racial-bonuses.md](../race/racial-bonuses.md) — race weight bonus
- DD-021 (v1'de kapalı kararı)

## 10. Bekleyen Kararlar
- Aktivasyon eşiği (Phase 02 playtest metric).
- Weight rate curve (STR per weight; movement penalty formulü).
- UI presence (weight bar).
- Alternate model: "encumbrance" tek bir cap (Phase 02'de dahi opsiyonel).
