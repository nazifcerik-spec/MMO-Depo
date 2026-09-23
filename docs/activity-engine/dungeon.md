# Dungeon Activity (Manuel + Idle-Run)

## 1. Amaç
Dungeon aktivitesi. İki tür: **manuel run** (party + oyuncu aktif) ve **idle-run** (contract-based AFK, düşük yield).

Yüksek seviye: [game-design/dungeons.md](../game-design/dungeons.md).

## 2. Params

```json
{
  "dungeon_template_id": "...",
  "difficulty": "normal",              // normal|heroic|mythic
  "mode": "manual",                     // manual|idle
  "party_id": null                       // manual mode party ise
}
```

## 3. Modlar

### 3.1 Manuel Mode
- `activity.instance` oluşturulur (contract değil).
- Instance-scope: dungeon expire ~2 saat.
- Party ile birlikte aktif oyun.
- Boss encounter'lar oyuncular tarafından karar-tabanlı geçilir.
- Yield: dungeon complete ile boss loot table.

Instance model: [database/activity.md §3-4](../database/activity.md).

### 3.2 Idle Mode
- Contract-based; primary channel işgal eder.
- Party opsiyonel.
- Analytical yield compute (per-tick simulasyon yok).
- Yield ~%50 manuel run'un (bkz. [game-design/dungeons.md §4](../game-design/dungeons.md)).

## 4. Prerequisites

- Character level >= dungeon.min_level.
- Party size dungeon.party_size_min..party_size_max (v1: 1-5).
- Weekly lockout not active (bkz. [database/progression-quests-achievements.md §8](../database/progression-quests-achievements.md)).
- Character at dungeon entrance location (v1 opsiyonel; Phase 02 zorunlu).

Fail:
- `dungeon.level_low`
- `dungeon.party_size_invalid`
- `dungeon.lockout_active`
- `dungeon.entrance_required`

## 5. Manuel Run — Lifecycle

```
1. Party leader → POST /v1/dungeons/enter
2. Server create activity.instance (seed, expires_at = now+2h).
3. Party members' activity.instance_participant satırları.
4. Encounter progression:
   - Each boss encounter → skill cast, damage exchange.
   - Loot roll on boss defeat.
5. Complete: instance.state='completed'.
6. Weekly lockout write (if applicable).
7. Rewards granted (per participant).
```

Manuel run **contract-based değil** — instance-based. Reveal etkilenmez.

## 6. Idle Run — Contract Based

```
1. Character (or party) start dungeon_idle contract:
   POST /v1/characters/{id}/activities/start
     body: { activity_type: "dungeon_idle", params: { dungeon_template_id: "...", difficulty: "normal" }}
2. Server:
   - Validate.
   - Create contract (channel=primary; contract.snapshot).
3. Reveal (analytical):
   - Rate compute: clears_per_hour = f(character/party DPS, dungeon total HP)
   - Yield: rewards × clears × idle_yield_multiplier (~0.5)
```

## 7. Idle-Run Rate Compute

```go
func (a *DungeonIdleActivity) ComputeRate(ctx, snap, params) Rate {
    dungeon := content.Dungeons.ByID(params.DungeonTemplateID)

    // Estimate dungeon clear time based on DPS + defense
    totalMobHP := dungeon.EstimatedTotalHP(params.Difficulty)
    partyDPS := estimatePartyDPS(snap, party)
    // fallback solo: snap DPS
    clearTimeSec := totalMobHP / partyDPS + dungeon.OverheadSec  // trash + boss transitions

    clearsPerHour := 3600.0 / clearTimeSec

    // Idle multiplier
    clearsPerHour *= config.rates.dungeon.idle_yield_multiplier  // 0.5

    // XP + loot per clear
    xpPerClear := dungeon.EstimatedXPReward(params.Difficulty)

    return Rate{
        YieldPerHour: clearsPerHour,   // full clears
        XPPerHour:    clearsPerHour * xpPerClear,
        LootTableSampleRate: clearsPerHour * dungeon.LootRollsPerClear,
    }
}
```

## 8. Idle Yield Compute

```go
elapsedHours := float64(elapsedMs) / 3600000.0
clears := int32(rate.YieldPerHour * elapsedHours)  // rounded down

// Weekly lockout: first clear this week → full reward; subsequent → normal
firstClearRewardApplied := checkWeeklyLockout(char, dungeon)

// Loot roll per clear
lootDrops := []ItemDrop{}
for i := int32(0); i < clears; i++ {
    for _, boss := range dungeon.BossEncounters {
        clearDrops := lootSystem.RollTable(boss.LootTableID, seed.Sub(fmt.Sprintf("clear-%d-boss-%s", i, boss.Code)))
        lootDrops = append(lootDrops, clearDrops...)
    }
}

// Death probability (idle-run — party ile daha az)
deathCount := samplePoisson(0.02 * elapsedHours, 1.0, seed.Sub("deaths"))  // low rate

return Yield{
    XPCharacter: int64(rate.XPPerHour * elapsedHours),
    ItemDrops:   lootDrops,
    DeathCount:  deathCount,
}
```

**Not**: `clears × boss_count × loot_roll` bounded loop OK (küçük N ~1-10 clears × 3-5 boss = 5-50 iteration). Per-tick simulation yok (per-clear döngü, deterministic + küçük).

## 9. Threshold Detection

Idle run threshold'lar:
- **Inventory dolma**: kill başına avg loot slot;
- **Weekly lockout**: bir sonraki clear'da lockout hit → contract'ı auto-pause + notification.
- **Party disband**: mid-contract party leader ayrılırsa contract solo devam eder.

## 10. Weekly Lockout

`progression.dungeon_lockout` tablosu (bkz. [database/progression-quests-achievements.md §8](../database/progression-quests-achievements.md)):
- Weekly reset Salı 06:00 UTC (DD-016).
- Reveal sırasında lockout check:
  - Manuel/idle her ikisi de lockout artırır.
  - Lockout dolar → rare loot drop rate düşer (bkz. dungeon config); XP normal devam.

## 11. Party Idle Run

Party ile idle run:
- Party leader'ın contract'ı; her üye kendi contract'ını başlatır (paylaşımlı seed opsiyonel).
- Yield per member kişisel (bkz. party loot rule).

v1'de basit: her karakter kendi contract'ı, DPS'i solo hesapla + party bonus modifier (+5% xp, +5% loot chance).

## 12. Content: Dungeon Template

Bkz. [database/content-templates.md §14-15](../database/content-templates.md).

```yaml
# content/dungeons/ashen_keep.yaml
code: ashen_keep
display_name_key: dungeon.ashen_keep.name
region_code: ashen_hills
min_level: 20
max_level: 30
party_size_min: 1
party_size_max: 5
weekly_lockout: true
idle_run_supported: true

boss_encounters:
  - code: keep_guardian
    order: 1
    loot_table_code: ashen_keep_guardian_loot
  - code: fire_wraith
    order: 2
    loot_table_code: ashen_keep_wraith_loot
  - code: ashen_lich_king
    order: 3
    loot_table_code: ashen_keep_boss_loot

estimated_total_hp:
  normal: 250000
  heroic: 500000
  mythic: 1000000

overhead_sec: 60
loot_rolls_per_clear: 3
```

## 13. Endpoint

**POST `/v1/dungeons/enter`** — manuel run start.
**POST `/v1/characters/{id}/activities/start`** — idle-run contract.
**POST `/v1/dungeons/instances/{id}/leave`** — manuel exit.

## 14. Reveal Örnek

```
Dungeon Idle: Ashen Keep normal, 6h elapsed
Snapshot: Warrior lvl 25, DPS 200, solo
Total HP: 250000 → clear time = 1250s + 60s overhead = ~22 min
Clears per hour = 60/22 = ~2.7 → idle × 0.5 = 1.35

6h × 1.35 = 8.1 → 8 clears
XP per clear: 3000
Total XP: 24000

Loot:
  each clear rolls 3 boss loot tables
  boss_loot chance rare 40%, epic 40%, legendary 15%
  8 clears × 3 bosses = 24 rolls
  ~10 rare items, ~10 epic, ~4 legendary

Weekly lockout: first clear full reward; 7 clears reduced (rare loot penalty).
```

## 15. Metrics
- `activity_dungeon_manual_completed_total{dungeon}`.
- `activity_dungeon_idle_clears_avg{dungeon}`.
- `activity_dungeon_party_size_distribution`.
- `activity_dungeon_weekly_lockout_hit_total`.

## 16. Cross-Reference
- [Activity](activity.md), [Combat](combat.md), [Hunting](hunting.md), [Rewards](rewards.md)
- [game-design/dungeons.md](../game-design/dungeons.md), [game-design/party.md](../game-design/party.md)
- [database/activity.md](../database/activity.md), [database/progression-quests-achievements.md §8](../database/progression-quests-achievements.md)

## 17. Bekleyen Kararlar
- Mythic+ affix (Phase 02+).
- Party DPS estimate accuracy — v1 solo DPS × 1.5 approx.
- Raid content (10+) — Phase 02.
- Cross-party dungeon queue (Phase 02).
