# Trade Tests

## 1. Amaç
P2P trade session güvenliği — state machine, auto-unlock, confirm hash, race, block, bind bypass.

## 2. Kritik Invariant'lar

- **Auto-unlock on offer change** (DD-138) — bait-and-switch guard.
- **Confirm hash match** — `seen_partner_offer_hash` mismatch → reject.
- **Finalize SERIALIZABLE** — recheck herşey.
- **Block enforce** — silent drop.
- **Concurrent session cap 1** — karakter başına.

## 3. Unit — State Machine

Table-driven:
- `open + lock(a)` → `initiator_locked`.
- `initiator_locked + lock(b)` → `both_locked`.
- `both_locked + confirm(a)` → `initiator_confirmed`.
- `both_confirmed` → auto finalize.
- ✓ Illegal transition (`confirm from open`) → reject.

## 4. Unit — Auto-Unlock on Update

- Both locked → offer update by a → both unlocked (version++).
- ✓ Confirm request after change with old hash → `E_TRADE_OFFER_CHANGED`.

## 5. Unit — Confirm Hash

- Server offer hash H1; client `seen_partner_offer_hash=H1` → OK.
- Offer changed to H2; client sends H1 → reject.
- Client sends H2 (yeni) → OK.

## 6. Unit — Bind Rule

- BoP item added to offer → `E_ITEM_NOT_TRADABLE`.
- BoE item OK.

## 7. Unit — Wallet Sufficiency

- Offer gold 5000c; wallet 3000c → offer update reject.

## 8. Integration — Finalize Atomic

- Two chars, offer set, both confirm.
- Finalize:
  - ✓ Items swapped (owner change).
  - ✓ Gold ledger (a: -X, b: +X).
  - ✓ Item locks released.
  - ✓ Session state='finalized'.
  - ✓ Audit trade_log entry.

## 9. Integration — Bait-and-Switch Scenario

- Setup:
  - a offers iron_sword.
  - Both lock. Both confirm.
- Race: just before finalize, a offer changed to nothing (simulated timing bug).
- ✓ Finalize recheck fails → session cancelled; item locks released.

## 10. Integration — Concurrent Same-Item

- a offers item X in session_1.
- a tries offer same item X in session_2 → session_1 lock prevents.

## 11. Integration — Block Enforcement

- b blocks a.
- a → trade_request(b) → silent drop (a sees "sent"; b sees nothing).
- ✓ No trade session created.

## 12. Integration — Rate Limit

- a sends 5 trade_request(b) in 30s → 4th+ reject (spam guard cooldown).

## 13. Property-Based — Item Conservation

- Random 100 successful trades between random chars.
- Assert: `SUM(items across all chars) == constant`.
- Assert: `SUM(gold) == constant` (no minting).

## 14. Fuzz — Session Message Sequence

- Random sequence of update/lock/confirm/cancel → state machine never enters invalid state; server never panic.

## 15. Coverage
- `domain/marketplace/trade`: ≥ 90%.

## 16. Cross-Reference
- [marketplace/trade.md](../marketplace/trade.md), [security/trade-validation.md](../security/trade-validation.md), [social/block.md](../social/block.md)
