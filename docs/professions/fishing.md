# Fishing

## 1. Amaç
Su kütlelerinde balık ve deniz mahsulü toplama. Cooking, alchemy, quest reagent besler.

Aktivite detay: [activity-engine/fishing.md](../activity-engine/fishing.md).

## 2. Profession Data

```yaml
# content/professions/fishing.yaml
code: fishing
category: gathering
display_name_key: profession.fishing.name
max_level: 100
xp_curve: { base: 50, factor: 1.10 }
tool_slot: main_hand
tool_template_codes: [basic_rod, silk_rod, ancient_rod]
bait_slot: off_hand           # bait item; opsiyonel
race_affinity_bonus_pct:
  merfolk: 8
```

## 3. Spot Türleri

| Spot Code | Water Type | Tier | Catch Distribution |
|---|---|---|---|
| `emerald_coast_pier` | coast | 1 | common trout/bream/pike/crab |
| `river_ashfield` | river | 1 | freshwater fish |
| `deep_ocean` | ocean | 3 | tuna/marlin/sea_serpent (rare) |
| `frozen_lake` | ice | 3 | arctic char, ice pike |
| `swamp_pond` | swamp | 2 | swamp fish + rare alchemy reagent |
| `underground_lake` | cave | 4 | blind fish (rare enchant material) |

Detay content pack: [world/resource-node.md](../world/resource-node.md).

## 4. Bait System

Optional bait item (envanterde bulunmalı):
- **Worm Bait**: +10% yield.
- **Rare Lure**: +5% rare drop chance.
- **Magic Fly**: +15% yield + tier +1 hedefleme.

Consume rate: catches × baitPerCatch (config 0.1 default).

## 5. Merfolk Affinity

**Merfolk**: `tide_born` +5% fishing yield + affinity +8% = **+13% total**.

Underwater fishing (Phase 02) merfolk-exclusive; aquatic mount ile diğer race unlock.

## 6. Rank Tier'lar

Aynı pattern. Grandmaster unlock: **"Master Angler"** — %2 shot at legendary sea trinket.

## 7. Tool

- `basic_rod`: no bonus.
- `silk_rod`: +5% yield + +2% rare.
- `ancient_rod`: +10% yield + +5% rare + tier scaling.

## 8. Fishing Event (Phase 02)

- Server-wide "Fishing Tournament" weekend event.
- Leaderboard for biggest catch.
- Unique cosmetic reward.

## 9. Weather Interaction (Phase 02)

Rain → +10% fishing yield (opt-in bonus). Bkz. [travel/weather.md](../travel/weather.md).

## 10. AFK Aktivite

Bkz. [activity-engine/fishing.md](../activity-engine/fishing.md).

## 11. Cross-Reference

- [Mining](mining.md), [Herbalism](herbalism.md)
- [activity-engine/fishing.md](../activity-engine/fishing.md)
- [world/resource-node.md](../world/resource-node.md)
- [race/racial-bonuses.md](../race/racial-bonuses.md) — merfolk

## 12. Metrics
- `profession_fishing_yield_total`.
- `fishing_bait_consumed_total`.
- `fishing_rare_catch_total`.

## 13. Bekleyen Kararlar
- Bait v1 aktif (öneri: evet basit).
- Underwater fishing Phase 02.
- Weather bonus Phase 02.
- Fishing tournament event Phase 02.
