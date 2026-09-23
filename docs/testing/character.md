# Character Tests

## 1. Amaç
Karakter lifecycle (create/enter/delete), stat calculation, XP/level, freeze/unfreeze, prestige, presence.

## 2. Kritik Invariant'lar

- **Slot cap 2** karakter/account v1 (DD-OPEN-8).
- **Unique name** karakter isim.
- **Level up event single-fire** — XP boost concurrent → tek level up.
- **Freeze deterministic** — mutation reject in single check.
- **Delete 48h grace** — cancel window.

## 3. Unit — Create

- ✓ Slot cap 2; 3rd create → `E_SLOT_LIMIT`.
- ✓ Duplicate name → `E_NAME_TAKEN`.
- ✓ Invalid name format → `E_NAME_INVALID`.
- ✓ Successful create → row + starter inventory + starter location.

## 4. Unit — XP → Level

- XP curve deterministic (content).
- ✓ Level 41, XP need to 42: 12,500.
- ✓ Gain 15,000 → level 42, remaining 2,500 XP toward 43.
- ✓ Multi-level up in single grant → all level_up events fired sequentially; stat gain accumulated.

## 5. Unit — Stat Calculation

Modifier stack: flat → mult_add → mult_prod → clamp.
- Base STR 20, +5 from race, +8 from class, +3 gear flat, +10% gear mult_add, +5% buff mult_prod.
- ✓ Result: `((20+5+8+3) × (1+0.1)) × (1+0.05) = 41.58`.

## 6. Unit — Freeze / Unfreeze

- ✓ Freeze → all mutating endpoints reject `E_CHARACTER_FROZEN`.
- ✓ Read endpoints OK.
- ✓ Unfreeze → mutations resume.
- ✓ Pending reward queue drains post-unfreeze.

## 7. Unit — Delete Grace

- ✓ Delete request → state='delete_grace', delete_at=now+48h.
- ✓ Cancel within window → state='active'.
- ✓ Grace expire → worker hard-delete (soft-delete first + audit).

## 8. Unit — Presence

- Login → presence:online (Redis TTL 5m).
- Heartbeat every 25s → TTL refresh.
- No heartbeat 5m → auto offline.

## 9. Integration — Ascension

- Character lvl 60 + playtime 60h + questline complete + cost paid.
- POST ascend →
  - Level → 1, XP → 0.
  - Class skill points refunded.
  - Mastery XP UNCHANGED.
  - Prestige_points_available +1.
  - Ledger cost debit.
  - Audit entry.

## 10. Integration — Prestige Node Choose

- Prestige_points_available 3, choose 3 nodes → available=0.
- Modifier stack invalidate → next stat fetch reflects bonuses.

## 11. Property-Based — Stat Modifier Order

- Random modifier stack input → apply order idempotent + deterministic.
- No permutation changes result (except intentional order).

## 12. E2E — Full Character Journey

- Register → create character → mining 1h → reveal → level up → ascend.
- Verify chain of events + audit.

## 13. Coverage
- `domain/character`: ≥ 90%.

## 14. Cross-Reference
- [character](../character/README.md), [progression/ascension.md](../progression/ascension.md), [skill-engine/modifiers.md](../skill-engine/modifiers.md)
