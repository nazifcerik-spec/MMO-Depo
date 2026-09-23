# Skinning

## 1. Amaç
Öldürülen mob'lardan deri/kürk/kemik/pençe çıkarma. Leatherworking + jewelcrafting + alchemy'ye besler. **Yan aktivite** — hunting kill'lerinden pasif olarak tetiklenir.

## 2. Profession Data

```yaml
# content/professions/skinning.yaml
code: skinning
category: gathering
display_name_key: profession.skinning.name
max_level: 100
xp_curve: { base: 50, factor: 1.10 }
tool_slot: null                # skinning knife opsiyonel (Phase 02)
tool_template_codes: []
race_affinity_bonus_pct:
  orc: 5                        # secondary affinity (orc primary hunting)
```

## 3. Mekanik

Skinning **standalone aktivite değil** — hunting kill'inden tetiklenir:
1. Hunting aktivitesi çalışıyor (bkz. [activity-engine/hunting.md](../activity-engine/hunting.md)).
2. Mob kill (skinnable family — beast/reptile/dragon vs).
3. Skinning aktif ise: mob loot table'a **skinning drops** eklenir.

## 4. Mob Family Skinning Table

| Family | Skinning Yield |
|---|---|
| Beast (wolf, bear) | leather_scrap, thick_hide, pelt_alpha |
| Reptile (lizard, dragon) | scale_shard, dragon_hide (rare) |
| Insect | chitin_plate, silk_thread |
| Elemental | essence_shard (rare) |
| Humanoid | – (not skinnable v1) |

Detay content pack: `content/skinning_tables/beast.yaml`, vs.

## 5. Yield Compute

Hunting aktivite yield içinde:
```go
if skinningActive && mob.IsSkinnable {
    skinning_table := content.SkinningTables[mob.Family]
    for _, entry := range skinning_table.Entries {
        // Same aggregate binomial pattern (bkz. activity-engine/hunting.md §5)
        expected := kills × entry.dropChance × avgQty
        actualCount := sampleBinomial(kills, entry.dropChance, seed.Sub("skin-"+entry.itemCode))
        yield.ItemDrops = append(...)
    }
    yield.XPProfession[skinningID] += kills × 3  // baseline
}
```

Kill başına ~3 skinning XP.

## 6. Tool (Phase 02)

- **Basic skinning knife**: +5% skin yield.
- **Steel knife**: +10% yield + rare drop bonus.

v1 tool yok; skinning otomatik.

## 7. Race Affinity

**Orc**: secondary affinity (+5%; primary hunting).

## 8. Rank Tier'lar

Aynı pattern. Grandmaster: **"Master Skinner"** — rare "pristine hide" drop chance +2%.

## 9. Skinning Item Kataloğu (Örnek)

| Item | Kaynak Mob Family | Kullanım |
|---|---|---|
| leather_scrap | beast (any tier) | leatherworking basic |
| thick_hide | beast tier 2+ | leatherworking mid |
| alpha_pelt | alpha mob (rare) | leatherworking rare |
| dragon_hide | dragon | legendary armor material |
| chitin_plate | insect | jewelcrafting + leatherworking |
| scale_shard | reptile | alchemy reagent |

## 10. Cross-Reference

- [Hunting](hunting.md), [Herbalism](herbalism.md)
- [activity-engine/hunting.md](../activity-engine/hunting.md)
- [game-design/professions.md](../game-design/professions.md)
- [race/racial-bonuses.md](../race/racial-bonuses.md) — orc

## 11. Metrics
- `profession_skinning_yield_total{item}`.
- `profession_skinning_rare_drops_total`.

## 12. Bekleyen Kararlar
- Skinning knife Phase 02.
- Standalone skinning aktivite (recycle carcass mid-hunt) — hayır; passive.
- Family-specific unique skin materials Phase 02.
- Cross-profession leather quality tier balance.
