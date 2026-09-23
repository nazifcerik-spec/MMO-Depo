# Crafting XP

## 1. Amaç
Crafting iteration'dan gelen üç XP kanalı (character, profession, mastery) hesabı; fail penalty; level up cascading.

Foundation: [character/xp.md](../character/xp.md), [activity-engine/xp.md](../activity-engine/xp.md), [skill-engine/mastery.md](../skill-engine/mastery.md).

## 2. Üç XP Kanalı

| Kanal | Kaynak | Ledger |
|---|---|---|
| **Character XP** | Recipe xp_reward × multipliers | `character.xp_ledger` |
| **Profession XP** | Recipe xp_reward × 5 (approx) | `character.profession.xp` |
| **Mastery XP** | Recipe xp_reward × 0.3 | `character.mastery.xp` |

## 3. Success vs Fail XP

Success iteration:
- Full character XP.
- Full profession XP.
- Full mastery XP.

Fail iteration:
- %50 character XP.
- %50 profession XP.
- %50 mastery XP.
- Material fail refund policy separate (bkz. [recipes.md §5](recipes.md)).

Big miss (rare fail):
- **Zero XP** (all channels).
- Full material loss.

## 4. Reveal Sırasında XP Hesabı

Craft aktivite yield (bkz. [activity-engine/crafting.md §5](../activity-engine/crafting.md)):
```go
successCount = sampleBinomial(iterations, craftSuccessPct/100)
failCount = iterations - successCount
bigMissCount = sampleBinomial(failCount, bigMissChance/100)
regularFailCount = failCount - bigMissCount

characterXP := recipe.XPReward × (successCount + failCount×0.5)   // big miss dahil zero contribution
professionXP := recipe.XPReward × 5 × (successCount + regularFailCount×0.5)
masteryXP := recipe.XPReward × 0.3 × (successCount + regularFailCount×0.5)

// Multiplier stack (character XP için)
characterXP *= raceMult × eventMult × partyMult
```

## 5. Multiplier Stack

Character XP multiplier'lar:
- Race passive (`human_versatility +5%`).
- Rested XP (Phase 02).
- Event XP bonus.
- Party XP bonus.
- Item XP boost potion.

Cap: 3.0× (bkz. [activity-engine/xp.md §3](../activity-engine/xp.md)).

Profession + mastery XP genelde flat (multiplier daha az).

## 6. Level Up Cascading

Craft aktivite reveal içinde:
1. Character XP grant → level up cascading (bkz. [character/level.md](../character/level.md)).
2. Profession XP grant → profession level up cascading.
3. Mastery XP grant → milestone check.

## 7. Cross-Recipe XP

Bazı recipe'ler crafting profession'a değil, farklı domain'e XP verir:
- Cooking recipe → cooking profession + character XP (küçük).
- Smelt raw ore → blacksmithing profession XP küçük (basic operation).

Advanced recipe (masterwork tier) → +200% profession XP + character XP.

## 8. Anti-Cheat

Per-hour crafting XP cap (macro/bot koruma):
- Character XP: crafting kanalından max ~%20 total XP/gün.
- Profession XP: 5000/hour cap.
- Mastery XP: bkz. [skill-engine/mastery.md §11](../skill-engine/mastery.md).

Cap aşımı log warn + throttle (silent).

## 9. XP Debt Interaction

XP debt yalnız combat death'ten gelir; crafting death yok (safe zone).

Ancak character XP grant sırasında xp_debt varsa payoff (bkz. [character/xp.md §7](../character/xp.md)):
```go
if char.XPDebt > 0 {
    debtPaid := min(char.XPDebt, characterXP)
    char.XPDebt -= debtPaid
    characterXP -= debtPaid
}
```

## 10. Config

```yaml
crafting_xp:
  fail_xp_pct: 50
  big_miss_xp: 0
  profession_xp_multiplier: 5
  mastery_xp_multiplier: 0.3
  per_hour_cap:
    character: 20000
    profession: 5000
    mastery: 1500
```

## 11. Metrics
- `craft_xp_granted_total{kind,recipe}`.
- `craft_xp_per_hour_avg{profession}` gauge.
- `craft_xp_cap_hit_total{kind}`.

## 12. Audit
- `crafting.xp_granted` (severity=debug; high volume).
- `crafting.profession_leveled_up`.
- `crafting.mastery_leveled_up`.

## 13. Test Stratejisi

- Multiplier stack cap.
- Fail XP proportional.
- Big miss zero XP.
- Level up cascading (multi-level jump).
- Per-hour cap enforcement.

## 14. Cross-Reference
- [Recipes](recipes.md), [Crafting Queue](crafting-queue.md), [Quality](quality.md), [Critical Craft](critical-craft.md), [Mastery](mastery.md)
- [activity-engine/xp.md](../activity-engine/xp.md), [activity-engine/crafting.md](../activity-engine/crafting.md)
- [character/xp.md](../character/xp.md), [character/level.md](../character/level.md)
- [skill-engine/mastery.md](../skill-engine/mastery.md)

## 15. Bekleyen Kararlar
- Fail XP oran final (v1: 50%).
- Big miss compensation (yumuşat vs harsh) v1 harsh.
- Profession XP multiplier tune (v1: 5×).
- Rested XP applied to crafting (Phase 02).
