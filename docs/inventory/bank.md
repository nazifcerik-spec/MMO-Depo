# Character Bank

## 1. Amaç
Şehir bankasında karakter özel uzun-süreli depolama. Tab-based; her tab satın alınır (gold sink).

## 2. Model

Kaynak: [database/character.md §4-5](../database/character.md).

```
character.bank_slot
  character_id, tab_index, slot_index, item_instance_id, stack_count, updated_at

character.bank_tab
  character_id, tab_index, name, purchased_at, purchase_cost_gold, capacity SMALLINT DEFAULT 40
```

Composite PK: `(character_id, tab_index, slot_index)`.

## 3. Tab Model

- **İlk tab**: karakter oluştururken ücretsiz.
- **Sonraki tabs**: gold sink pattern:
  - Tab 2: 5000 gold (500000 copper)
  - Tab 3: 20000 gold
  - Tab 4: 100000 gold
  - Tab 5: 500000 gold
  - Tab 6: 2000000 gold (cap v1)
- **Max tab**: 6 (v1); Phase 02'de 10.
- **Tab capacity**: 40 slot each (uniform).

Tab isimleri kullanıcı özelleştirir; default "Tab 1", "Tab 2", ...

## 4. Erişim Kısıtı

Bank erişimi sadece şehirlerdeki **bank service** POI'sinde:
```
GET /v1/characters/{id}/bank
Auth + CharacterScope + character.current_location has bank service
```

Location check:
- `character.current_location_id` → `world.location` → parent settlement.
- `settlement.services_json.bank == true`.

Aksi halde 400 `bank.not_available_here`.

## 5. Endpoint

### 5.1 Bank Listing

**GET `/v1/characters/{id}/bank`**
Response:
```json
{
  "tabs": [
    {
      "tab_index": 0,
      "name": "General",
      "capacity": 40,
      "used": 22,
      "slots": [
        { "slot_index": 0, "item_instance_id": "...", "stack_count": 1, "template_code": "..." },
        ...
      ]
    },
    { "tab_index": 1, "name": "Mats", "capacity": 40, "used": 5, "slots": [...] }
  ],
  "next_tab_cost_copper": 2000000000,
  "next_tab_index": 2,
  "max_tabs": 6
}
```

### 5.2 Deposit

**POST `/v1/characters/{id}/bank/deposit`**
```json
{
  "moves": [
    { "from_inv_slot": 5, "to_tab": 0, "to_slot": 10, "count": 1 }
  ]
}
```

Server (transaction per move):
- Validate:
  - Source inventory slot has item.
  - Item not locked, not bound-to-other, not equipped.
  - Destination tab exists (purchased).
  - Destination slot empty or same template + stackable.
- Transfer: item_instance owner_kind='character_bank', owner_id=character_id.
- Update `character.inventory_slot` remove; `character.bank_slot` insert.
- Log transfer.

### 5.3 Withdraw

**POST `/v1/characters/{id}/bank/withdraw`**
```json
{
  "moves": [
    { "from_tab": 0, "from_slot": 10, "to_inv_slot": 5, "count": 1 }
  ]
}
```

Server: reverse of deposit.

### 5.4 Purchase Tab

**POST `/v1/characters/{id}/bank/purchase-tab`**
Server:
- Validate: next tab available.
- Debit gold via ledger (`debit_bank_tab_purchase`).
- Insert `character.bank_tab`.
- Response 200: new tab.

### 5.5 Rename Tab

**POST `/v1/characters/{id}/bank/rename-tab`**
```json
{ "tab_index": 0, "name": "Materials" }
```

Server:
- Validate name length (1-32).
- Update `bank_tab.name`.

### 5.6 Move Within Bank

**POST `/v1/characters/{id}/bank/move`**
```json
{ "from_tab": 0, "from_slot": 5, "to_tab": 1, "to_slot": 10, "count": 1 }
```

Cross-tab move ok.

## 6. Bulk Operations

**POST `/v1/characters/{id}/bank/deposit-all-matching`**
```json
{
  "to_tab": 1,
  "filter": { "template_codes": ["iron_ore", "copper_ore"] }
}
```

Server:
- Envanter tara → filter match → deposit iterate.
- Any fail → rollback (transaction).

Rate limit: 10/dk (bulk operations).

## 7. Bank Item Rules

- Any owner_kind item deposit edilebilir (BoP dahil — bank karakter özel).
- Guild bank ayrı sistem (bkz. [storage.md](storage.md)).
- Bank item satılamaz doğrudan (önce inventory'e withdraw).
- Bank item marketplace'e list edilemez (önce inventory).

## 8. Item Loss Risk

- Bank içinde item asla kaybolmaz.
- Karakter silme → bank içeriği:
  - Soft delete (7 gün grace): bank item korunur.
  - Hard delete: bank item → void.
  - Guild bank ayrı; karakter guild üyeliği çıkarken bank etkisiz.

## 9. Concurrency

Optimistic lock per bank slot (`updated_at` bump).

## 10. Fail Modları
- Not in bank location → 400 `bank.not_available_here`.
- Tab not purchased → 404 `bank.tab_not_owned`.
- Slot occupied → 409 `bank.slot_occupied`.
- Insufficient funds (tab purchase) → 402.

## 11. Metrics
- `bank_tabs_purchased_total{tab_index}`.
- `bank_deposit_total`, `bank_withdraw_total`.
- `bank_avg_used_slots_per_char` gauge.

## 12. Audit
- `bank.deposit`, `bank.withdraw`, `bank.tab_purchased`, `bank.tab_renamed`.

## 13. Test Stratejisi
- Location check enforce.
- Deposit → withdraw round-trip.
- Bulk deposit atomic (partial fail rollback).
- Tab purchase gold debit + insert.
- Concurrent deposit → optimistic lock.

## 14. Cross-References
- [Inventory](inventory.md), [Storage](storage.md), [Item Transfer](item-transfer.md), [Stacking](stacking.md)
- [database/character.md §4-5](../database/character.md)
- [game-design/inventory.md §4](../game-design/inventory.md)
- [game-design/cities.md](../game-design/cities.md) — bank service

## 15. Bekleyen Kararlar
- Tab cost curve final (5k → 2M gold).
- Max tab 6 vs 10.
- Bank search (Phase 02).
- Cross-character bank (BoA Phase 02).
- Sorted view save (Phase 02).
