# Hunting Activity

## 1. Amaç
AFK combat aktivitesi. Belirli spot'ta mob spawn ile dövüş, kill yield ile XP + loot.

## 2. Params

```json
{
  "spot_id": "loc-whispering-forest",
  "mob_type_codes": ["ashen_wolf", "shadow_wolf"],   // filter
  "autopotion": true,
  "stop_on_hp_below_pct": 30,                         // safety threshold
  "use_pet": true
}
```

## 3. Prerequisites

- Character alive.
- Character at hunt spot location.
- Character rotation defined (`character.rotation.priority_list`).
- Envanter en az 1 slot boş.

Fail:
- `hunting.not_at_spot`
- `hunting.no_rotation`
- `hunting.dead`

## 4. Rate Compute (Analytical DPS Model)

Hunting rate = kill count per hour. Bunu **DPS ile mob HP** oranından çözeriz:

```go
func (a *HuntingActivity) ComputeRate(ctx, snap, params) Rate {
    spot := content.HuntSpots.ByID(params.SpotID)

    // Effective DPS from snapshot rotation + gear
    dps := skillEngine.EstimateDPS(ctx, snap, params.MobTypeCodes)
    // Approx: (Σ ability_damage × cast_rate) with cooldown + rotation

    // Average mob HP in filter
    avgMobHP := averageMobHP(spot, params.MobTypeCodes)

    // Time per kill (seconds) = mob HP / DPS
    timePerKillSec := avgMobHP / dps

    // Add mob spawn interval + gap
    effectivePerKillSec := timePerKillSec + spot.SpawnIntervalSec

    killsPerHour := 3600.0 / effectivePerKillSec

    // XP per kill from mob templates
    avgXPPerKill := averageMobXP(spot, params.MobTypeCodes)

    // Loot table roll — per kill drops (weighted)
    lootRollsPerHour := killsPerHour

    return Rate{
        YieldPerHour:        killsPerHour,
        XPPerHour:           killsPerHour * avgXPPerKill,
        MasteryXPPerHour:    killsPerHour * 3,  // bestiary mastery
        LootTableSampleRate: lootRollsPerHour,
    }
}
```

`skillEngine.EstimateDPS` — snapshot içindeki rotation, skill damage, cooldown, resource'u analyzing DPS estimate döner (formula: Σ damage × castRate ÷ rotation duration).

## 5. Yield Compute (Analytical)

```go
elapsedHours := float64(elapsedMs) / 3600000.0
expectedKills := rate.YieldPerHour * elapsedHours
actualKills := samplePoisson(rate.YieldPerHour, elapsedHours, seed.Sub("kills"))

// Loot roll for each kill — aggregated
lootDrops := []ItemDrop{}
for _, mobCode := range params.MobTypeCodes {
    mob := content.NPCs.ByCode(mobCode)
    killsForMob := actualKills / int32(len(params.MobTypeCodes))  // eşit dağıtım approx

    // Loot table: her drop entry weight'e göre expected count
    for _, entry := range mob.LootTable.Entries {
        expectedDropCount := float64(killsForMob) * entry.DropChance * float64(entry.MinQty+entry.MaxQty)/2
        if entry.Guaranteed {
            lootDrops = append(lootDrops, ItemDrop{TemplateCode: entry.ItemCode, Quantity: int32(expectedDropCount)})
        } else {
            actualCount := sampleBinomial(killsForMob, entry.DropChance, seed.Sub(fmt.Sprintf("loot-%s-%s", mobCode, entry.ItemCode)))
            avgStack := (entry.MinQty + entry.MaxQty) / 2
            lootDrops = append(lootDrops, ItemDrop{TemplateCode: entry.ItemCode, Quantity: int32(actualCount * avgStack)})
        }
    }
}

// Death probability per hour
deathRate := estimateDeathRate(snap, spot)  // 0-1 per hour, based on mob damage vs char HP + regen
deathCount := samplePoisson(deathRate, elapsedHours, seed.Sub("death"))

xpDebt := int64(0)
if snap.Level >= 20 {
    xpDebt = int64(deathCount) * levelXPRequired(snap.Level) / 20  // ~5%
}

return Yield{
    XPCharacter:  int64(rate.XPPerHour * elapsedHours),
    ItemDrops:    lootDrops,
    DeathCount:   deathCount,
    XPDebt:       xpDebt,
    MasteryXP: map[uuid.UUID]int64{
        bestiaryTrackID: int64(actualKills * 3),
    },
}
```

**Kritik**: her kill için per-tick simulation YOK. Toplu aggregate sample.

## 6. Threshold Detection

Inventory dolma:
```
avg_loot_slots_per_kill := estimateAvgLootSlotUsage(mob.LootTable)
kills_until_inv_full := char.InventoryFree / avg_loot_slots_per_kill
ms_until_inv_full := kills_until_inv_full / rate.YieldPerHour * 3600000
```

Death cascade (safety):
```
if params.StopOnHPBelowPct is set and death rate too high:
    // simplified: if deathRate > threshold → threshold at first death expected time
    if deathRate > 0.5:  // >0.5 death/hour → risky
        ms_until_first_death = expected time until first death (poisson inverse)
        pause with reason 'death_cascade_prevention'
```

Potion out:
```
if params.Autopotion and char.HasPotion:
    potion_consume_per_kill := estimatePotionUsage(snap, mob)
    potions_until_out := char.PotionCount / potion_consume_per_kill
    ms_until_potion_out := potions_until_out / rate.YieldPerHour * 3600000
```

Return smallest.

## 7. Rotation Estimate DPS

`skillEngine.EstimateDPS`:
```
weight = 0
totalTime = 0
for entry in rotation.priority_list:
    skill = load(entry.skill_id)
    cast_rate = 1 / (skill.cooldown_ms + skill.cast_time_ms + gcd_ms)
    damage_per_cast = computeSkillDamage(snap.DerivedStats, skill, target=mob)
    weight += damage_per_cast * cast_rate * conditionProbability(entry.cond)
    totalTime += fractional time this skill hits

dps = weight / rotation_duration
```

Approx; heuristic. Deterministic for same snapshot.

## 8. Pet Contribution

Snapshot içinde `ActivePetInstanceID` var:
- Pet stats + pet skill contribution DPS'e eklenir (~%5-15).
- Rate compute'da entegre.

## 9. Content: Hunt Spot

```yaml
# content/locations/whispering_forest.yaml
code: whispering_forest
kind: hunt_spot
region_code: emerald_valley
recommended_level_min: 15
recommended_level_max: 25
mob_pool_codes: [ashen_wolf, shadow_wolf, forest_bear]
spawn_interval_sec: 3
danger_tier: normal
```

Mob pool template'te (bkz. [database/content-templates.md §17](../database/content-templates.md)).

## 10. Death Handling in Reveal

Death expected in AFK:
- Death event insert (`character.death_log`).
- XP debt accumulate.
- Durability loss.
- Contract paused during death respawn window (30s each).
- Auto-revive after 5min (bkz. [character/health.md §5.4](../character/health.md)).

Death cascade prevention: `stop_on_hp_below_pct` safety → threshold kill switch.

## 11. Reveal Örnek

```
Hunting @ Whispering Forest, ashen_wolf targets, 6h elapsed
Snapshot: Ranger lvl 20, DPS ~180
Mob HP: avg 850 → time per kill ~5s + spawn 3s = 8s per kill
Kills per hour = 3600/8 = 450

Poisson sample: 2703 kills
Loot:
  common leather_scrap × ~2200 (~85% drop, x1-3 avg)
  wolf_fang × ~600 (25% chance)
  alpha_pelt × ~5 (0.5% rare)
  ancient_wolf_fang × 0 (0.01% legendary, sample = 0)

XP: 15 per kill × 2703 = 40545 character XP
Bestiary mastery XP: 3 × 2703 = 8109

Death: rate 0.05/h × 6h = 0.3 → sample 0 deaths (lucky)
Inventory used: many stacks; auto-merge → 12 slots used
```

## 12. Metrics
- `activity_hunting_kills_avg{spot,mob}`.
- `activity_hunting_dps_estimated_avg`.
- `activity_hunting_death_total`.
- `activity_hunting_paused_reason_total{reason}`.

## 13. Cross-Reference
- [Activity](activity.md), [Combat](combat.md), [Offline Progress](offline-progress.md), [Rewards](rewards.md)
- [skill-engine/README.md](../skill-engine/README.md)
- [game-design/hunting.md](../game-design/hunting.md)
- [character/health.md](../character/health.md)

## 14. Bekleyen Kararlar
- DPS estimate accuracy (Phase 02 refinement — cond probability model).
- Party AFK hunt (bkz. [combat.md](combat.md) — party channel).
- Rare mob spawn (server-wide anons) Phase 02.
- Weather / time-of-day effect (Phase 02).
