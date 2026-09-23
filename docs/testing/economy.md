# Economy Tests

## 1. Amaç
Ledger integrity (double-entry sum=0), idempotency, wallet cap, reconcile drift, reversal correctness.

## 2. Kritik Invariant'lar

- **Sum=0** her currency per TX.
- **Idempotency** — retry safe (`idempotency_key`).
- **Wallet cap** — 10^12 copper aşımı reject.
- **Reversal exact reverse** — sign flip, aynı account.
- **Reconciler drift** — < 1c warn, ≥ 1c freeze + P0.

## 3. Unit — Double-Entry Sum Invariant

- Table-driven vendor buy:
  - Entries: buyer -1000, seller +950, tax_sink +50.
  - ✓ `SUM(delta per currency) == 0`.
  - ✓ Balance_after correct her account.
- Malformed entry (sum=42) → `E_LEDGER_INVARIANT` (bug guard).

## 4. Unit — Idempotency Cache

- ✓ Same idempotency_key + same body → cached response.
- ✓ Same key + different body → `E_IDEMPOTENT_MISMATCH` (409).
- ✓ Different key → new TX.
- ✓ Retention 24h; expire sonrası aynı key yeni TX olabilir.

## 5. Unit — Wallet Cap

- ✓ Wallet 10^12 - 500; grant 1000 → `E_CURRENCY_CAP_EXCEEDED`; ledger NOT written.
- ✓ Cap aşımı inbox notification queued.

## 6. Unit — Currency Scope Validation

- ✓ Character endpoint + account-scope currency (gem) → `E_CURRENCY_SCOPE_MISMATCH`.
- ✓ Guild endpoint + character scope → reject.

## 7. Unit — Reversal

- Original TX: buyer -1000, seller +950, tax +50.
- Reversal: buyer +1000, seller -950, tax -50.
- ✓ Sum invariant korunur.
- ✓ Aynı account, sign flip.
- ✓ Reversal audit entry: operator, reason, original_txn_id ref.

## 8. Integration — Full Vendor Buy

- Testcontainer.
- Character wallet 5000c, vendor iron_sword 4500c.
- POST buy → wallet 500c, inventory +1 iron_sword, ledger 3 entries, audit entry.
- Assert: character_wallet.gold_copper = 500, ledger sum=0, audit row.

## 9. Integration — Concurrent Buy Same Wallet

- 5 goroutine aynı wallet'ten 1000c buy attempt (total 5000c ama wallet 3000c).
- Assert: 3 success, 2 reject; wallet 0c; no drift.

## 10. Integration — Reconciler Drift

- Manually simulate wallet drift (test hook: insert ledger without wallet update).
- Run reconciler.
- Assert: drift detected, character.state='frozen', P0 event emitted.

## 11. Property-Based

- Property: after 1000 random valid txns, `SUM(ledger.delta per account) == wallet.balance` her account için.
- Property: SERIALIZABLE isolation + concurrent txns → no double-write.

## 12. Fuzz — Ledger Entry

- Random malformed entry (sum≠0, negative amount on unbounded=false) → validation reject before DB.

## 13. E2E — Market Flow

- Two characters: sell + buy P2P via market.
- Verify: seller net = price - cut; buyer -price; sink += cut.
- Ledger 3 entries per side; audit trace complete.

## 14. Coverage
- `domain/economy`: ≥ 92%.
- `systems/txn-service`: ≥ 95%.

## 15. Cross-Reference
- [economy/transaction.md](../economy/transaction.md), [economy/gold.md](../economy/gold.md), [security/currency-validation.md](../security/currency-validation.md), [security/transaction-validation.md](../security/transaction-validation.md)
