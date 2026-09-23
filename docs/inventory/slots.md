# Slots (Slot Mechanics)

## 1. Amaç
Slot-based storage'ların ortak mekaniği: index, capacity, sparse representation, expansion.

## 2. Slot Modeli

Sparse: boş slot = row yok.
```
character.inventory_slot
  (character_id, slot_index) PK
  item_instance_id NOT NULL   -- boşsa row yok
  stack_count NOT NULL
```

Slot doldurma: `INSERT ON CONFLICT DO NOTHING` veya `UPDATE`.

## 3. Slot Türleri

| Kind | Slot Count | Slot Aralığı | Expansion |
|---|---|---|---|
| Inventory | 40..120 | 0..119 | Bag item (Phase 02) |
| Bank Tab | 40 (uniform) | 0..39 per tab | Tab satın al |
| Equipment | 17 (fixed) | slot enum | – |
| Guild Bank Tab | 100 | 0..99 per tab | Tab satın al (guild) |
| Mail attachments | 12 per mail | 0..11 | Mail sayısı ile |

## 4. Slot Index vs Position

Slot index sadece **kalıcı adres**; UI layout ile ayrı:
- Client sort/filter free.
- Server sadece slot_index index'i takip eder.

## 5. Capacity Compute

**Inventory**:
```
capacity = 40 + Σ(equipped_bags.slot_bonus) + account.bag_slot_bonus
```

**Bank total**:
```
capacity = num_purchased_tabs * 40
```

**Guild Bank**:
```
capacity = num_purchased_tabs * 100
```

Used:
```sql
SELECT COUNT(*) FROM character.inventory_slot WHERE character_id = $1;
```

## 6. Slot Assignment Strategy

Yeni item eklenirken slot seçimi:

### 6.1 Explicit Slot
Client belirli slot ister:
- Validate boş / same-template merge.
- Fail: 409.

### 6.2 Auto-Assign (Overflow senaryosu)
Server ilk uygun slot seçer:
```
1. Iterate slot_index 0..capacity-1
2. Empty slot bul → yerleştir
3. Aynı template stack yer varsa → merge tercih (Phase 02 opsiyonel)
4. Yer yoksa → overflow (mailbox)
```

## 7. Slot Expansion

**Bag Item** (Phase 02):
- Bag slot: 4 (character can equip 4 bags).
- Bag item stat: `slot_bonus`.
- Effective capacity: bag equip/unequip ile artar/azalır.
- Capacity düşerse item'lar taşar → mailbox overflow (with warning).

**Account Bag Bonus**:
- `auth.account.bag_slot_bonus` (monetization; Phase 02).

**Bank Tab Purchase**:
- Gold cost curve (bkz. [bank.md §3](bank.md)).
- Insert `character.bank_tab`.

## 8. Slot Reordering

Client sort için:
- Move endpoint (bkz. [inventory.md §5.2](inventory.md)).
- Server: swap or shift; slot_index preserve.
- Alternatif: client-side sort visualization only (server slot_index sabit); Phase 02 opsiyonu.

v1 tercih: **explicit move endpoint** (persist).

## 9. Slot Preview (Phase 02)

Server-side slot preview UI için:
- `GET /v1/characters/{id}/inventory/preview?filter=...`
- Compact response (id + template_code + count) — full item detail lazy load.

## 10. Concurrency

- PK `(character_id, slot_index)` unique.
- Insert conflict → 409 slot occupied.
- Cross-slot swap → 2 rows in transaction.

## 11. Slot Migration

Content patch bag mekaniği eklerse (Phase 02):
- Migration: mevcut karakter default 40 slot.
- Ek slot expansion opt-in.

## 12. Slot Corruption

Data integrity check (worker):
- Item instance owner_kind='character_inventory' ancak inventory_slot row yok → warn + auto-place first empty slot.
- Inventory_slot row var ancak item_instance void → cleanup.

Worker frequency: günlük.

## 13. Test Stratejisi
- Move + swap + split + merge kombinasyonu.
- Overflow senaryosu.
- Slot conflict race.
- Sort persistence.

## 14. Metrics
- `slot_used_ratio{kind}` histogram.
- `slot_expansion_purchase_total`.

## 15. Cross-References
- [Inventory](inventory.md), [Bank](bank.md), [Equipment](equipment.md), [Stacking](stacking.md), [Storage](storage.md)
- [database/character.md](../database/character.md)
- [game-design/inventory.md](../game-design/inventory.md)

## 16. Bekleyen Kararlar
- Bag slot Phase 02 tam schema.
- Client-side sort vs persist server-side.
- Guild bank tab count expansion (guild level ile).
- Auto-fill sort ("compact inventory") Phase 02 endpoint.
