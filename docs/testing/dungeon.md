# Dungeon Tests

## 1. Amaç
Floor progression, encounter engine, AFK run reveal, weekly lockout, personal loot.

## 2. Kritik Invariant'lar

- **Weekly lockout** karakter başına (reset Salı 06:00 UTC — DD-016).
- **AFK run analytical** — floor/room ilerleme O(1) hesap.
- **Personal loot** her party member ayrı roll.
- **Fail state atomic** — wipe = state='failed', loot ödül yarım tier.
- **Boss first-clear** first successful defeat için event + achievement.

## 3. Unit — Floor Progression

- Dungeon 5 floor, 4 room each; enemy pack + boss last.
- ✓ Room clear → `current_room++`.
- ✓ Floor clear → `current_floor++, current_room=0`.
- ✓ Last floor last boss → state='succeeded'.

## 4. Unit — Weekly Lockout

- Character run heroic → lockout marker set for template + difficulty until Salı 06:00 UTC.
- ✓ 2nd run same char + template + difficulty → `E_LOCKOUT_ACTIVE`.
- ✓ Reset job at Salı 06:00 clears markers.
- ✓ Different difficulty (mythic) not blocked.

## 5. Unit — AFK Run Analytical

- 45dk AFK run, character gear_score 1200.
- ✓ Expected floor progression via server analytical (per-floor DPS proxy).
- ✓ Reveal yields loot per floor cleared.
- ✓ Boss chance: if final floor reached & DPS threshold met → boss defeated.

## 6. Unit — Personal Loot

- 3-member party wipes boss.
- Loot roll 3 times (one per member).
- ✓ Each member gets own item instance.
- ✓ No shared loot pool contention.

## 7. Integration — Full Live Run

- Testcontainer.
- Party of 3 enter iron_mine (heroic).
- Advance through 5 floors via mock combat.
- Boss defeat → loot rolls per member + XP + guild XP + achievement first_clear if applicable.

## 8. Integration — AFK Run Cancel

- Start AFK run; before completion, character cancels.
- Assert: partial loot NOT granted (yield forfeited); state='cancelled'.

## 9. Integration — Wipe Scenario

- Party wipes on floor 3.
- Assert: state='failed'; ödül tier 1 (attend) minimum; XP some; no boss loot.

## 10. Property-Based — Loot Rate

- 10k dungeon completions → drop rate ± 3σ of expected per rarity.

## 11. E2E — First Kill Bonus

- Server-wide first defeat of new boss → server broadcast + achievement + title.
- Subsequent defeats normal reward.

## 12. Regression

- `testdata/dungeon/iron_mine_heroic_afk_run_seed_101.json`.

## 13. Coverage
- `domain/dungeon`: ≥ 85%.

## 14. Cross-Reference
- [dungeon](../dungeon/README.md), [dungeon/afk-run.md](../dungeon/afk-run.md), [testing/loot.md](loot.md)
