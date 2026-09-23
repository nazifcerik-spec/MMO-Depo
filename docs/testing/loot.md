# Loot Tests

## 1. Amaç
Loot roll deterministic RNG, personal loot (DD-012), rarity distribution, drop table integrity.

## 2. Kritik Invariant'lar

- **Deterministic RNG** (seed + counter) — DD-OPEN-7 (v1: xoshiro256++ öneri).
- **Personal loot default** — party her üye ayrı roll.
- **Drop rate content-driven** — code'da hardcode yok.
- **Rarity threshold linear** — [common..mythic] contiguous.
- **Pity system opsiyonel** (Phase 02) — v1 saf shuffle.

## 3. Unit — RNG Determinism

- ✓ `RNG(seed=42, counter=0).Roll()` → sabit değer.
- ✓ Aynı seed + counter → aynı sonuç 10.000 iteration.
- ✓ Farklı counter → farklı sonuç (bit-diff > 0).
- ✓ Cross-platform (Windows/Linux) aynı sonuç (Go xoshiro implementation portable).

## 4. Unit — Loot Table Resolve

Content:
```yaml
drop_table_ashen_boss:
  guaranteed:
    - { item: gold_pile, min: 100, max: 500 }
  weighted:
    - { item: iron_sword,     rarity: uncommon, weight: 40 }
    - { item: iron_plate,     rarity: uncommon, weight: 30 }
    - { item: fire_ring,      rarity: rare,     weight: 20 }
    - { item: legendary_axe,  rarity: legendary,weight:  1 }
```

- ✓ Weighted sum = 91 (guaranteed hariç).
- ✓ 100.000 roll → legendary_axe count ≈ 100.000 × 1/91 ≈ 1099 ± 3σ.
- ✓ Guaranteed her roll = 1 gold_pile in range.
- ✓ Empty table → empty loot (no crash).

## 5. Unit — Rarity Distribution

- ✓ Common 60%, Uncommon 25%, Rare 10%, Epic 4%, Legendary 0.9%, Mythic 0.1% content'te bir default table için — 100k roll → count ± 3σ.
- ✓ Rarity boost multiplier (world event `+2% rare`) → uncommon → rare cross rate ölçüle.

## 6. Unit — Personal Loot (DD-012)

- 3-member party, defeats boss.
- Server rolls 3 × loot table (once per member).
- ✓ Each member's `character_id` recorded on loot instance.
- ✓ No shared "greed/need" contention (v1).
- ✓ Master loot party option (raid Phase 02) — flag toggles.

## 7. Property-Based — Chi-Square Distribution

- 100.000 roll aynı drop table → chi-square test p > 0.01 (dist beklenene uyar).
- Bias tespit: consecutive seed'ler farklı sonuç (RNG bias yok).

## 8. Integration — End-to-End Loot Grant

- Testcontainer.
- Combat session end → loot roll → inventory insert + ledger (gold + XP source event).
- Assert: item instance owner_character_id = self; source_ref = combat_session_id.

## 9. Integration — AFK Reveal Loot

- 1h mining contract; reveal.
- ✓ Yield hesabı analytical; loot instance insert; ledger + audit.
- ✓ Personal loot cache — reveal 2x → aynı loot.

## 10. Fuzz — Drop Table Config

- Random malformed drop table config yükleme → validation error, panic yok.

## 11. Regression

- `testdata/loot/dragon_boss_seed_101.json` — fixed seed loot output.

## 12. Coverage
- `domain/loot`: ≥ 90%.
- `systems/rng`: ≥ 95% (kritik).

## 13. Cross-Reference
- [systems/loot.md](../systems/loot.md), [systems/rng.md](../systems/rng.md), [systems/reward-engine.md](../systems/reward-engine.md)
