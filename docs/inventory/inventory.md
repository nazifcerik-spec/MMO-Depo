# Character Inventory (Bag)

## 1. Amaç
Karakter üstündeki ana taşıma alanı. Slot-bazlı; kapasite bag item ile genişleyebilir.

## 2. Model

Kaynak: [database/character.md §3](../database/character.md).

```
character.inventory_slot
  character_id UUID (PK)
  slot_index SMALLINT (PK) CHECK 0..119
  item_instance_id UUID FK
  stack_count INT
  updated_at TIMESTAMPTZ
```

Sparse: boş slot = row yok.

## 3. Kapasite (Slot Count)

- **Baz**: 40 slot.
- **Bag item ile artış**: her equip edilmiş bag item slot ekler.
- **Cap**: 120 (bag'lerle).
- **Bag slot bonus (account)**: `auth.account.bag_slot_bonus` (monetization Phase 02 için hazır).

Effective:
```
capacity = 40 + Σ(equipped_bag.slot_bonus) + account.bag_slot_bonus
```

## 4. Bag Item Modeli (Phase 02)

v1'de bag item şablonu opsiyonel. Phase 02 için hazır schema:
- Bag item `category='utility', subcategory='bag'`.
- Base stat: `slot_bonus INT` (örn. Small Bag +6, Medium Bag +12, Large Bag +18).
- Karakter'ın "bag slots"ı (Phase 02): 4 slot; her birine bag item takılır → cumulative.

v1 tercih: bag mekaniği devre dışı; direct slot expansion Phase 02 açık.

## 5. Endpoint

### 5.1 Inventory Detail

**GET `/v1/characters/{id}/inventory`**
Auth + CharacterScope.

Response:
```json
{
  "capacity": 40,
  "used": 32,
  "slots": [
    {
      "slot_index": 0,
      "item_instance_id": "ii-...",
      "stack_count": 1,
      "template_code": "steel_longsword",
      "rarity": "rare",
      "locked": false
    },
    {
      "slot_index": 1,
      "item_instance_id": "ii-...",
      "stack_count": 87,
      "template_code": "iron_ore",
      "rarity": "common"
    }
  ]
}
```

Response format compact; ancak client tarafında template detay lookup için ayrı `GET /v1/items/{id}` çağırılabilir.

### 5.2 Move (Slot içi)

**POST `/v1/characters/{id}/inventory/move`**
Request:
```json
{ "from_slot": 5, "to_slot": 12 }
```

Server:
- Load both slots FOR UPDATE.
- Swap veya put (destination boşsa).
- Version bump.
- Response 200.

### 5.3 Split

**POST `/v1/characters/{id}/inventory/split`**
```json
{ "from_slot": 3, "amount": 20, "to_slot": 8 }
```

Server:
- Validate: `from_slot.stack_count > amount`.
- Validate: `to_slot` boş veya aynı template + stackable.
- Transaction:
  - `from_slot.stack_count -= amount`.
  - `to_slot` insert veya merge (aynı template ise stack_count += amount).
- Response.

### 5.4 Merge

**POST `/v1/characters/{id}/inventory/merge`**
```json
{ "from_slot": 3, "into_slot": 8 }
```

Server:
- Same template check.
- Combined <= stack_max.
- Transaction.

### 5.5 Salvage

**POST `/v1/characters/{id}/inventory/salvage`**
Body:
```json
{
  "slots": [3, 4, 8, 12]
}
```

Server:
- Validate ownership + not equipped + not locked + not bound (BoP salvage v1'de allowed? — karar: **allowed**, ancak confirm).
- For each: consume item; roll salvage materials from `content.item_template.salvage_output` (Phase 02 disenchanting).
- Update inventory.
- Audit event.

v1'de basit salvage: item → material tier'a bağlı iade.

### 5.6 Destroy

**POST `/v1/characters/{id}/inventory/destroy`**
```json
{
  "slots": [7],
  "confirmations": [{ "slot": 7, "confirm_name": "Steel Longsword" }]
}
```

Server:
- Validate name match (accidental delete koruma).
- Rare+ item için ek confirmation gerekir (client 3sn hold).
- Item → void.
- Audit.

### 5.7 Vendor Sell

**POST `/v1/characters/{id}/inventory/vendor-sell`**
Auth + CharacterScope + city vendor location required.
```json
{ "slots": [1, 5, 6, 9] }
```

Server:
- Validate: her item not equipped + not locked + BoE veya no_bind (BoP genelde sell edilemez; template config).
- For each: compute vendor_sell_price × rarity_multiplier.
- Transaction:
  - Credit gold via ledger (`credit_vendor_sale`).
  - Item → void.
- Response: total gold + item list.

## 6. Auto-Pause on Full

AFK aktivite envanter dolduğunda:
- Activity contract `status='paused'`, `pause_reason='inventory_full'`.
- Notification.
- Yeni yield discard değil; sadece durur (bkz. [systems/afk-engine.md](../systems/afk-engine.md)).

Slot açıldığında (salvage/sell/withdraw) → activity resume opsiyonu (auto vs manual):
- **Manual** tercih (v1) — oyuncu explicit `POST /v1/characters/{id}/activities/{cid}/resume`.
- Phase 02 opsiyonu: auto-resume config bool.

## 7. Overflow Mailbox

Envanter dolu ve mob loot geldi (rare drop):
- Öncelik: envanter aktif slot varsa insert.
- Yoksa: **mailbox'a overflow** (system mail, 7 gün TTL).
- Notification "New item awaits in mail".

Konfig:
```yaml
inventory:
  overflow_target: mailbox   # mailbox | drop_refuse | ground (Phase 02)
```

## 8. Sıralama & Filtreleme

Client-side (server slot_index sıralı döner):
- Sort by rarity, name, level, value.
- Filter: category, rarity, usable, tradable, locked.
- Bulk actions (checkbox multi-select).

Server-side arama Phase 02 opsiyonu.

## 9. Concurrency

- Optimistic lock per slot (updated_at + slot_index unique).
- Move/split/merge atomic; concurrent modification → 409 retry.

## 10. Fail Modları
- Slot not found → 404.
- Slot boş → 400 `inventory.slot_empty`.
- Stack cap exceed → 400 `inventory.stack_max_exceeded`.
- Locked item → 400 `inventory.item_locked`.
- Bound + salvage → 400 (confirm gerekiyorsa 202 confirmation flow).

## 11. Rate Limit
- Move / split / merge: 60/dk.
- Salvage bulk: 30/dk.
- Destroy: 20/dk (rare event).
- Vendor sell: 30/dk.

## 12. Metrics
- `inventory_slot_used_ratio` gauge.
- `inventory_action_total{action}` counter.
- `inventory_overflow_to_mail_total` counter.
- `inventory_auto_pause_total` counter.

## 13. Audit
- `inventory.item_moved`.
- `inventory.item_salvaged`.
- `inventory.item_destroyed` (severity=warn; rare+ severity=info).
- `inventory.item_sold_vendor`.
- `inventory.overflow_to_mail`.

## 14. Test Stratejisi
- Move happy path + swap.
- Split → merge round-trip.
- Salvage transaction rollback.
- Destroy name confirmation.
- Overflow mailbox routing.
- Concurrent move race.

## 15. Cross-References
- [Slots](slots.md), [Stacking](stacking.md), [Item Transfer](item-transfer.md), [Item Lock](item-lock.md), [Equipment](equipment.md), [Bank](bank.md), [Storage](storage.md)
- [database/character.md §3](../database/character.md)
- [systems/afk-engine.md](../systems/afk-engine.md)
- [item-engine/item-instance.md](../item-engine/item-instance.md)
- [game-design/inventory.md](../game-design/inventory.md)

## 16. Bekleyen Kararlar
- Auto-resume aktivite Phase 02.
- Bag item mekaniği Phase 02 vs v1 skip.
- Auto-salvage rule engine (rarity/level threshold based) Phase 02.
- Search endpoint server-side (Phase 02).
