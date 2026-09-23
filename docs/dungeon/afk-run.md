# AFK Run (Dungeon Idle Mode)

## 1. Amaç
Dungeon AFK contract — offline devam. Analitik hesap; per-tick simulation YOK (DD-084).

Foundation: [activity-engine/dungeon.md](../activity-engine/dungeon.md).

## 2. Idle Mode Kural Seti

- `activity.contract.activity_type = 'dungeon_idle'`.
- `channel = 'primary'`.
- Yield multiplier `afk_yield_multiplier = 0.5` (default; content override).
- Weekly lockout paylaşımlı manuel + idle.

## 3. Rate Compute

```go
func (a *DungeonIdleActivity) ComputeRate(ctx, snap, params) Rate {
    dungeon := content.Dungeons.ByID(params.DungeonTemplateID)

    // Total dungeon clear time
    party_dps := estimatePartyDPS(snap, party)
    total_hp := dungeon.EstimatedTotalHP(params.Difficulty)   // aggregate mob + boss HP

    per_clear_time_sec := (total_hp / party_dps) + dungeon.OverheadSec  // trash + boss + transitions
    per_clear_time_sec /= dungeon.AFKYieldMultiplier   // 0.5 → yavaşlar (yield düşer)

    clears_per_hour := 3600.0 / per_clear_time_sec

    reward_per_clear := dungeon.EstimatedRewardPerClear(difficulty)

    return Rate{
        YieldPerHour:        clears_per_hour,
        XPPerHour:           clears_per_hour * dungeon.XPRewardPerClear,
        LootTableSampleRate: clears_per_hour * dungeon.LootRollsPerClear,
    }
}
```

## 4. Yield Compute (Analitik)

```go
elapsedHours := elapsedMs / 3600000.0
raw_clears := rate.YieldPerHour * elapsedHours

// Weekly lockout cap
weekly_lockout_remaining := loadWeeklyLockoutRemaining(char, dungeon, difficulty)
clears_capped := min(int(raw_clears), weekly_lockout_remaining × dungeon.LockoutMultiplier)

// Bounded loop per clear (küçük N; ~1-10 clears in 8h)
loot_drops := []ItemDrop{}
xp := int64(0)
for i := 0; i < clears_capped; i++ {
    // Per clear: bounded loop per boss + per treasure
    for _, boss := range dungeon.BossEncounters {
        drops := lootSystem.RollTable(boss.LootTableID, seed.Sub(fmt.Sprintf("clear-%d-boss-%s", i, boss.Code)))
        // Apply weekly lockout effect per drop
        if i > 0 || weekly_lockout_hit:
            drops = applyLockoutPenalty(drops)
        loot_drops = append(loot_drops, drops...)
    }
    for _, chest := range dungeon.TreasureRooms:
        drops := lootSystem.RollTable(chest.LootTableCode, seed.Sub(...))
        loot_drops = append(loot_drops, drops...)
    for _, event := range dungeon.RandomEvents:
        // Poisson sample event trigger
        if triggered := sampleBernoulli(event.chance, seed.Sub(...)); triggered:
            apply_event(...)
    xp += dungeon.XPRewardPerClear
}

// Death probability (idle-run — low with party)
death_rate_per_hour := 0.02 * (dungeon.DangerScore)   // düşük
death_count := samplePoisson(death_rate_per_hour * elapsedHours, 1.0, seed.Sub("deaths"))

xp_debt := int64(0)
if snap.Level >= 20:
    xp_debt = int64(death_count) * levelXPRequired(snap.Level) / 20

return Yield{
    XPCharacter:  int64(rate.XPPerHour * elapsedHours),
    ItemDrops:    loot_drops,
    DeathCount:   death_count,
    XPDebt:       xp_debt,
    Currencies: map[string]int64{
        "dungeon_token": int64(clears_capped * 100),   // 100 token per clear
    },
}
```

Bounded loop per clear × per boss (küçük N, ~5-15 iteration total) OK — per-tick simulation değil.

## 5. Threshold Detection

- **Weekly lockout dolar**: threshold — kısmi clear + pause.
- **Inventory dolar**: threshold — pause.
- **Death cascade** (rare): pause.

```
clears_until_inv_full := inv_free_slots / avg_slots_per_clear
ms_until_inv_full := clears_until_inv_full / rate.YieldPerHour × 3600000
```

## 6. Party Idle Run

Party ile idle:
- Party leader'ın contract'ı reference.
- Her karakter kendi contract'ı (v1); yield per member.
- DPS aggregate (party total = Σ member DPS).
- Clear time party ile daha kısa.

## 7. Idle vs Manuel Loot

Content field:
```yaml
afk_yield_multiplier: 0.5
```

- Yield rate × 0.5 (clear time × 2).
- Rare drop rate aynı per clear (proportionate).
- Currency token aynı per clear.

Sadece rate düşer; drop kalitesi düşmez.

## 8. Weekly Lockout Priority

- Manuel run + idle run **paylaşımlı lockout counter**.
- Character full week clears manuel'de yaparsa idle run rare loot skip.
- Weekly reset Salı 06:00 UTC.

## 9. Contract Snapshot

Snapshot içinde:
- Party DPS estimate.
- Weekly lockout state.
- Character equipment/stats.

Contract süresince kararlı hesap (DD-088).

## 10. Reveal Örnek

```
Contract: dungeon_idle @ Ashen Keep normal, 6h elapsed
Snapshot: Warrior lvl 25, DPS 200, solo
Total HP: 25000 (all mobs + bosses)
Per clear time = 25000/200 + 60 sec overhead = 185 sec → idle × 2 = 370 sec ≈ 6.2 min
Clears per hour = 60/6.2 ≈ 9.7 → 6h × 9.7 = 58 raw clears

Weekly lockout: 3 clears remaining × 2 idle multiplier = 6 lockout units → clears_capped = 6
XP: 3000 × 6 = 18000

Loot per clear: 3 boss loot rolls × 6 = 18 boss rolls
  ~4 rare item, ~4 epic, ~2 legendary (weekly lockout: %30 penalty after first)
  
Death rate: 0.02 × 6 = 0.12 → sample 0 death
```

## 11. Cross-Reference
- [Dungeon](dungeon.md), [Boss](boss.md), [Loot](loot.md), [Random Event](random-event.md)
- [activity-engine/dungeon.md](../activity-engine/dungeon.md)
- [activity-engine/offline-progress.md](../activity-engine/offline-progress.md) — analitik pipeline
- [database/progression-quests-achievements.md §8](../database/progression-quests-achievements.md) — lockout

## 12. Metrics
- `dungeon_idle_clears_avg{dungeon}`.
- `dungeon_idle_lockout_hit_total`.
- `dungeon_idle_active_contracts` gauge.

## 13. Bekleyen Kararlar
- AFK yield multiplier final (v1: 0.5).
- Party idle run coordination (shared contract Phase 02).
- Idle run difficulty selection (v1 normal only; heroic Phase 02).
- Auto-restart contract post-lockout Phase 02.
