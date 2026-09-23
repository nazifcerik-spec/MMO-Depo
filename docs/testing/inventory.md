# Inventory Tests

## 1. Amaç
Item ownership, lock semantik, transfer atomicity, weight/slot cap, stack merge/split, dupe prevention.

## 2. Kritik Invariant'lar

- **Item instance UUID unique** (PK) — dupe imkansız.
- **Ownership tek** — bir instance tek karakterde.
- **Lock exclusive** — lock varken başka mutation yok.
- **Transfer atomic** — TX: source dec + dest inc; rollback either fail.
- **Weight/slot enforce** — cap aşımı reject.
- **Bind rule enforce** — BoP/BoA guard.

## 3. Unit — Slot / Weight Cap

- ✓ 60/60 slot; pickup → reject `E_INVENTORY_FULL`.
- ✓ Weight cap 25 kg; 24 kg env + 2 kg item pickup → reject.
- ✓ Stackable existing 90/99, pickup 15 same → 99 stack + overflow 6 new slot.

## 4. Unit — Lock Semantik

- ✓ Item `is_locked=true, lock_owner='market'` → sell endpoint → `E_ITEM_LOCKED{market}`.
- ✓ Aynı item 2 farklı endpoint concurrent → SERIALIZABLE + `SELECT FOR UPDATE` → biri success, biri reject.
- ✓ Lock release atomic (listing cancel) → unlock + follow-up sell OK.

## 5. Unit — Bind Rule

- ✓ BoP item trade → `E_ITEM_NOT_TRADABLE`.
- ✓ BoA item market listing → reject.
- ✓ BoE item first equip → binds; sonrasında trade reject.

## 6. Unit — Stack Split / Merge

- ✓ Split 99 stack → 40 + 59 iki instance (aynı template).
- ✓ Merge 40 + 30 same template → 70 single (yeni instance UUID veya kalan UUID kalır — v1 kararı documented).
- ✓ Merge overflow (99 cap): 99+30 → 99 + 30 (yeni).

## 7. Integration — Transfer Atomicity

- Testcontainer.
- Case: envanter → bank transfer, 100 iron.
- Fault inject: bank insert fail → env decrement rollback.
- Assert: env stack unchanged, bank empty.

## 8. Integration — Concurrent Transfer

- 2 goroutine aynı 50-stack'i farklı hedefe transfer.
- Assert: biri success; diğeri `E_INSUFFICIENT_STACK`.

## 9. Fuzz — Instance UUID

- Random UUID insert attempts → PK conflict prevents dupe.
- Reconciler sim: iki row (fake dupe) → drift alert.

## 10. Property-Based

- Property: `SUM(instances.count per template) == inventory.total(template)`.
- Property: `SUM(equipped) + SUM(inventory) + SUM(bank) + SUM(market_escrow) + SUM(mail_escrow) == character.total_owned` (append-only history + ledger reconcile).

## 11. Regression

- `testdata/inventory/complex_transfer_case_1.json` — golden output for edge cases.

## 12. Coverage
- `domain/inventory`: ≥ 90%.
- `systems/item-lock`: ≥ 95%.

## 13. Cross-Reference
- [inventory](../inventory/README.md), [inventory/item-lock.md](../inventory/item-lock.md), [inventory/item-transfer.md](../inventory/item-transfer.md), [inventory/stacking.md](../inventory/stacking.md)
