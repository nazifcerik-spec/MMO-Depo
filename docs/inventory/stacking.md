# Stacking

## 1. Amaç
Aynı template item'ların tek slot'ta birleştirilme kuralları. Materials, consumables, stackable currency-items için hayati.

## 2. Stack Kuralı

Aynı slot'ta stack olabilir mi? 6 kriter tam eşleşmeli:
1. **Template ID** aynı.
2. **Rarity** aynı (aynı template farklı rarity ⇒ ayrı — item_instance edge case).
3. **Bind status** aynı (BoP vs no_bind farklı slot).
4. **Affixes/sockets/enchant** — instance-özel değişkenler → **stack yok** (equipment tarzı unique roll).
5. **Durability** ≠ 100 (ekipman tarzı) → stack yok.
6. **Meta json** (custom flags) eşleşme.

Kısaca:
- **Materials / consumables**: instance-özel state yoksa → stack.
- **Equipment / instance-özel roll**: stack yok, her instance ayrı row.

## 3. stack_max

Template alanı:
```
content.item_template.stack_max INT NOT NULL DEFAULT 1
```

Örnek:
- `iron_ore`: `stack_max = 200`.
- `minor_health_potion`: `stack_max = 20`.
- `steel_longsword`: `stack_max = 1` (equipment; stack imkansız zaten).

## 4. Stack Merge

Yeni item eklenirken slot picker:
1. Uygun slot mevcut mu (same template + stack < max)?
2. Evet → merge (`stack_count += amount`, cap `stack_max`).
3. Overflow varsa: kalan yeni slot'a insert.

Örnek: `iron_ore x180` inventory'de + yeni drop `iron_ore x50` → `iron_ore x200` + yeni slot `iron_ore x30`.

## 5. Split Endpoint

**POST `/v1/characters/{id}/inventory/split`**
```json
{ "from_slot": 3, "amount": 20, "to_slot": 8 }
```

Server:
- Validate: `from.stack_count > amount`.
- Validate: `to` empty or same template + stack room.
- Transaction:
  - `from.stack_count -= amount`.
  - `to` empty → insert; same template → merge.

## 6. Merge Endpoint

**POST `/v1/characters/{id}/inventory/merge`**
```json
{ "from_slot": 3, "into_slot": 8 }
```

Server:
- Validate: same template + stack.
- Combined ≤ stack_max → move; > stack_max → 400 `stacking.max_exceeded`.

Alternatif: auto-partial merge (max fill, kalan from'da kalır).

## 7. Instance Kimliği

**Kritik**: Stackable item için tek `item_instance` mi, her stack ayrı instance mi?

**Karar (v1)**:
- **Aynı template + hepsi identical** → **tek item_instance**; slot bazında `stack_count` alanı.
- Farklı slot'larda aynı template stack varsa: her slot **kendi item_instance** taşır. Yani `iron_ore` inventory'nde 3 farklı slot'ta ise 3 instance vardır (`stack_count` her birinde farklı).

Alternatif (Phase 02 optimization): tek instance + reference count. Karmaşıklığa değmez v1.

## 8. Currency-Item vs Currency

- **Currency-item** (event token, dungeon token): stackable item; envanter slot işgal eder.
- **Currency** (gold): `character.wallet` — envanter işgal etmez.

Karar item template alanında `is_currency_item BOOL` (Phase 02) veya category='currency' subcategory=...

## 9. Cross-Storage Stack

Bank/inventory arası move sırasında:
- Destination slot'ta aynı template stack yer varsa → merge.
- Yoksa → yeni slot (varsa).

Bulk deposit:
```
inventory'de 3 slot iron_ore (100/200/50)
bank tab 0 slot 5: iron_ore (150)
deposit all iron_ore → 100 + 150 = 250 → slot 5 = 200 (cap), overflow 50 → slot 6 = 50 + eski inventory kalan (50 + 100 = 150 → slot 7 vs merge)
```

Complex — server auto-optimize; client "consolidate stacks" ayrı endpoint (Phase 02).

## 10. Consolidate (Compact) — Phase 02

**POST `/v1/characters/{id}/inventory/consolidate`**
Server:
- Envanter tara → aynı template stack'leri birleştir.
- Slot rearrangement.
- Rate limit: 5/dk.

## 11. Fail Modları
- `stacking.template_mismatch`
- `stacking.max_exceeded`
- `stacking.slot_occupied_different_template`
- `stacking.count_invalid` (negatif veya 0)

## 12. Concurrency

- PK slot unique; concurrent insert same slot → 409.
- Merge two slots → both FOR UPDATE.

## 13. Test Stratejisi

- Split → merge round-trip.
- Overflow allocation.
- Consolidate compression.
- Cross-storage stack picker.
- Instance identity preservation (unique roll'lu item stack olmaz).

## 14. Metrics
- `stacking_split_total`, `stacking_merge_total`.
- `stacking_max_exceeded_total`.
- `stacking_consolidate_total`.

## 15. Cross-References
- [Inventory](inventory.md), [Bank](bank.md), [Slots](slots.md), [Item Transfer](item-transfer.md)
- [database/character.md §3-5](../database/character.md)
- [item-engine/item-instance.md §7](../item-engine/item-instance.md)

## 16. Bekleyen Kararlar
- Consolidate endpoint v1 vs Phase 02 (öneri: Phase 02).
- Currency-item flag (subcategory yeter mi).
- Reference-counted single instance (Phase 02 optimization).
