# Calculate XP

## 1. Amaç
Aktivite yield'inde XP hesabı: character, profession, mastery. Multiplier stack. XP debt payoff. Level up trigger cascading.

Genel XP sistem: [character/xp.md](../character/xp.md). Bu doküman **activity-reveal bağlamında** XP hesap detayı.

## 2. Üç Kategori

| Kategori | Ledger | Level triggers |
|---|---|---|
| **Character XP** | `character.xp_ledger` | `character.character.xp` / `level` |
| **Profession XP** | (kayıt DB tarafında `character.profession.xp/level`) | Profession level up |
| **Mastery XP** | (kayıt DB tarafında `character.mastery.xp/level`) | Mastery milestone |

Her kategori ayrı ledger; conversion yok.

## 3. Character XP Hesap

```go
func computeCharacterXP(rate Rate, elapsedHours float64, snap Snapshot) int64 {
    baseXP := rate.XPPerHour * elapsedHours

    multipliers := 1.0
    // Race (human_versatility)
    multipliers *= 1.0 + racialXPBonusPct(snap.RaceCode)/100.0
    // Rested (Phase 02)
    // multipliers *= (1.0 + rested_bonus)
    // Event server-wide
    multipliers *= 1.0 + activeEventXPBonusPct(snap.ActiveBuffs)/100.0
    // Party bonus
    multipliers *= 1.0 + partyXPBonusPct(snap.InPartyOfSize)/100.0
    // Item XP potion
    multipliers *= 1.0 + itemXPBoostPct(snap.ActiveBuffCodes)/100.0
    // Ascension level effect (küçük)
    // multipliers *= 1.0 + ascensionXPFactor(snap.AscensionLevel)

    // Total cap 3.0×
    if multipliers > 3.0 { multipliers = 3.0 }

    return int64(baseXP * multipliers)
}
```

## 4. Profession XP Hesap

```go
func computeProfessionXP(rate Rate, elapsedHours float64, snap Snapshot) map[uuid.UUID]int64 {
    xpMap := map[uuid.UUID]int64{}
    for profID, xpPerHour := range rate.ProfessionXPPerHour {
        base := xpPerHour * elapsedHours
        // Profession affinity race bonus
        // Guild research profession xp bonus (Phase 02)
        multipliers := 1.0
        multipliers *= 1.0 + raceProfessionAffinityBonusPct(snap.RaceCode, profID)/100.0
        xpMap[profID] = int64(base * multipliers)
    }
    return xpMap
}
```

Rate per activity computeRate içinde her aktivite ilgili profession'a xp'sini attach eder (bkz. [mining.md §5](mining.md), [fishing.md §5](fishing.md)).

## 5. Mastery XP Hesap

Multi-track (weapon, combat discipline, profession, region, bestiary):
```go
func computeMasteryXP(rate Rate, elapsedHours float64, snap Snapshot, activity Activity) map[uuid.UUID]int64 {
    xpMap := map[uuid.UUID]int64{}
    // Weapon mastery — combat aktivite ise weapon based
    if activity.Category() == "hunting" || activity.Category() == "combat" || activity.Category() == "dungeon" {
        weaponKind := snap.EquippedWeaponKind
        weaponTrackID := masteryTrackByWeapon(weaponKind)
        xpMap[weaponTrackID] = int64(rate.MasteryXPPerHour * elapsedHours)
    }
    // Bestiary — hunting only
    if activity.Category() == "hunting" {
        bestiaryTrackID := masteryTrackByFamily(activity.MobFamilies())
        xpMap[bestiaryTrackID] = int64(rate.MasteryXPPerHour * elapsedHours * 0.5)
    }
    // Region mastery — always some XP for time spent
    regionTrackID := masteryTrackForRegion(snap.CurrentRegionCode)
    xpMap[regionTrackID] = int64(elapsedHours * 1.0)  // 1 xp/hour region time

    // Profession mastery covered by profession XP already
    return xpMap
}
```

## 6. XP Debt Payoff

Reveal başlangıcında karakter mevcut xp_debt varsa:
```go
func applyXPDebt(charID uuid.UUID, grantedXP int64) (appliedXP int64, debtPaid int64) {
    char := loadCharacter(charID)
    if char.XPDebt == 0 { return grantedXP, 0 }
    if char.XPDebt >= grantedXP {
        char.XPDebt -= grantedXP
        return 0, grantedXP  // hepsi debt'e gitti
    }
    // Kısmi
    debtPaid = char.XPDebt
    appliedXP = grantedXP - debtPaid
    char.XPDebt = 0
    return
}
```

Debt öder → `applied_delta` = kalan.

## 7. Level Up Cascading

XP grant sonrası level up check:
```go
for {
    need := levelXPRequired(char.Level, char.AscensionLevel)
    if char.XP < need { break }
    char.XP -= need
    char.Level++
    grantLevelRewards(char, char.Level)  // stat points, unlocks
    emit event 'character.leveled_up'
    if char.Level >= config.LevelCap {
        // cap — kalan XP birikimde
        break
    }
}
```

Multi-level jump destekli.

Detay: [character/level.md §4-5](../character/level.md).

## 8. Profession Level Up

Same pattern per profession:
```go
for {
    need := professionXPRequired(prof.Level)
    if prof.XP < need { break }
    prof.XP -= need
    prof.Level++
    unlockProfessionMilestone(char, prof, prof.Level)
    emit event 'profession.leveled_up'
    if prof.Level >= profession.MaxLevel { break }
}
```

## 9. Mastery Level Up + Milestone

Detay: [skill-engine/mastery.md](../skill-engine/mastery.md).

Level up → milestone check → passive apply + notification.

## 10. Reveal Sırasında XP Order

```
1. Grant character XP (debt payoff → applied)
   - Level up cascading
2. Grant profession XP per profession
   - Profession level up cascading
3. Grant mastery XP per track
   - Milestone claim
4. Emit summary event (character.reveal_xp_granted)
```

## 11. XP Cap During Reveal

Reveal büyük XP granting mümkün (uzun offline, high-level activity).
- Multi-level jump güvenli (level up loop bounded by XP cap).
- Level cap'te fazla XP birikir (ascension pending).

Anti-cheat:
- Reveal per-hour XP anomaly detect (bkz. [systems/afk-engine.md §12](../systems/afk-engine.md)).

## 12. Idempotency

Her XP grant idempotency key:
```
reveal:<contract_id>:<reveal_at>:char_xp
reveal:<contract_id>:<reveal_at>:prof_<profession_id>
reveal:<contract_id>:<reveal_at>:mastery_<track_id>
```

xp_ledger unique index enforce.

## 13. Test Stratejisi

- Multiplier stack cap (3.0×).
- Debt payoff scenarios (kısmi/tam).
- Level up cascading multi-jump.
- Profession level up.
- Mastery milestone claim.
- Idempotent double-reveal.

## 14. Metrics

- `xp_granted_total{kind,source}` counter.
- `xp_debt_paid_total`.
- `character_leveled_up_total{level}`.
- `profession_leveled_up_total{profession}`.
- `mastery_leveled_up_total{track}`.

## 15. Audit

- `character.xp_granted` (severity=debug; high volume).
- `character.leveled_up` (severity=info).
- `profession.leveled_up`.
- `mastery.leveled_up`.
- `mastery.milestone_claimed`.

## 16. Cross-Reference

- [character/xp.md](../character/xp.md), [character/level.md](../character/level.md)
- [skill-engine/mastery.md](../skill-engine/mastery.md)
- [game-design/progression.md](../game-design/progression.md), [game-design/professions.md](../game-design/professions.md)
- [Offline Progress](offline-progress.md), [Rewards](rewards.md)

## 17. Bekleyen Kararlar

- Multiplier cap final (3.0× v1; Phase 02 test).
- Rested XP buffer Phase 02.
- Ascension XP damping formulü final.
- Cross-mastery XP conversion (Phase 03).
