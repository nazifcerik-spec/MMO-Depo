# Class Tests

## 1. Amaç
Class + specialization + skill tree, respec, skill unlock, class restriction, spec change gating.

## 2. Kritik Invariant'lar

- **Class immutable** — race/class create sonrası değişmez v1 (Phase 03 opt).
- **Spec picked at lvl 10** — 3 spec/class.
- **Skill points level başına 1** — skill tree spend.
- **Respec cost formula** — DD (bkz. [economy/fee.md](../economy/fee.md)).
- **Class restriction on gear** — enforce equip.

## 3. Unit — Class Assignment

- ✓ Create character with class 'warrior' → `character.class_code='warrior'`.
- ✓ Class code whitelist enum.
- ✓ Post-create class change → `E_CLASS_IMMUTABLE`.

## 4. Unit — Specialization Pick

- ✓ Character lvl < 10 → `POST /spec/choose` → `E_LEVEL_TOO_LOW`.
- ✓ Character lvl 10 + first pick → spec set + skill tree unlock.
- ✓ Aynı spec 2. pick → no-op (idempotent).
- ✓ Different spec after respec allowed.

## 5. Unit — Skill Tree Unlock

- Warrior/Berserker spec tree 20 node total.
- ✓ Level 10 → 1 skill point available.
- ✓ Level up → +1 skill point.
- ✓ Spend on node → node unlocked, point --.
- ✓ Node prerequisite check (parent required).
- ✓ Insufficient points → `E_INSUFFICIENT_SKILL_POINTS`.

## 6. Unit — Respec

- Character lvl 42, skill nodes 32 spent.
- Cost: `level² × 100 = 176_400c` (cap 500g = 5_000_000c).
- ✓ Wallet insufficient → `E_INSUFFICIENT_BALANCE`.
- ✓ Respec success → all skill points refunded (32/32), all nodes cleared, ledger debit.

## 7. Unit — Class Restriction on Equipment

- Priest class trying to equip 'iron_axe' (warrior restriction) → `E_CLASS_RESTRICTION`.
- Universal item (no restriction) OK for any class.

## 8. Integration — Full Class Journey

- Register warrior → level 10 → pick berserker spec → spend 5 nodes → respec → pick guardian spec → spend 5 nodes.
- Verify: modifier stack reflects nodes correctly at each stage.

## 9. Integration — Modifier Stack from Skill Tree

- Spec node "Fury Strike +10% attack" active.
- Base attack 100 → stat compute 110.
- Deactivate node (respec) → 100.

## 10. Property-Based — Skill Point Conservation

- Property: `spent_points + available_points == level_earned_points` her karakter için.

## 11. Regression

- `testdata/class/warrior_berserker_full_tree.json`.

## 12. Coverage
- `domain/class`: ≥ 88%.
- `domain/skill-engine` (class spec integration): ≥ 85%.

## 13. Cross-Reference
- [class](../class/README.md), [class/respec.md](../class/respec.md), [skill-engine](../skill-engine/README.md), [economy/fee.md](../economy/fee.md)
