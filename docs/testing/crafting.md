# Crafting Tests

## 1. Amaç
Recipe validation, material consume atomicity, queue lifecycle, quality/critical craft, mastery XP.

## 2. Kritik Invariant'lar

- **Material atomic consume** — TX içinde inventory decrement + queue insert veya rollback.
- **Queue slot cap** — 5 aktif (Phase 02'de opsiyonel expand).
- **Quality deterministic** — seed + mastery tier + skill; snapshot at enqueue.
- **Critical craft chance formula** — content-driven.
- **Cancel iade** — bekleyen slot mat iade; başlanan slot iade yok.

## 3. Unit — Recipe Validation

- ✓ Recipe requires blacksmith level ≥ 15 → char lvl 14 → `E_INSUFFICIENT_LEVEL`.
- ✓ Recipe requires 3 iron + 1 coal → env'de sadece 2 iron → `E_INSUFFICIENT_MATERIAL`.
- ✓ Recipe active only if unlocked (mastery / quest gate).

## 4. Unit — Quality / Critical Craft

- Formula: `crit_chance = base_crit + mastery_tier × 0.02 + skill × 0.005`.
- ✓ Tier 3 + skill 50 → crit chance ≈ 0.16.
- ✓ Seed roll < 0.16 → critical craft; item quality "masterwork" + affix bonus.
- ✓ Critical → mastery XP × 1.5.

## 5. Unit — Queue Lifecycle

State machine:
- `queued` → `in_progress` → `completed` | `cancelled`.
- ✓ Enqueue → state='queued', material lock.
- ✓ First slot start → state='in_progress', timer set.
- ✓ Timer expire → state='completed', item instance insert, ledger, XP.
- ✓ Cancel `queued` → mat iade %100.
- ✓ Cancel `in_progress` → mat kayıp; time refund yok.

## 6. Integration — Full Craft Flow

- Testcontainer.
- Character with materials, enqueue 3 recipes.
- Clock advance → first complete → queue moves to slot 2 auto.
- Verify: inventory item added, mat consumed, XP gained, mastery XP gained.

## 7. Integration — Batch Craft

- Enqueue 10× "iron ingot".
- Advance time → each slot completes sequentially.
- Verify: 10 iron ingot instance, 10× mat consumed, single queue exit event on last.

## 8. Concurrent Craft Slot Test

- Aynı character 5 slot dolu → 6. enqueue → `E_QUEUE_FULL`.
- Different characters concurrent → each own queue.

## 9. Property-Based — Critical Craft Distribution

- 100k roll, mastery tier 5, skill 100 → observed crit rate ± 1% of expected.

## 10. Regression

- `testdata/crafting/iron_sword_seed_42.json` — deterministic quality output.

## 11. Coverage
- `domain/crafting`: ≥ 85%.

## 12. Cross-Reference
- [crafting](../crafting/README.md), [crafting/critical-craft.md](../crafting/critical-craft.md), [crafting/queue.md](../crafting/queue.md), [progression/mastery.md](../progression/mastery.md)
