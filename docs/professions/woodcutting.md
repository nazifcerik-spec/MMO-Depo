# Woodcutting

## 1. Amaç
Ağaç kesip odun/reçine/kabuk toplama. Blacksmithing (grip), leatherworking (kap), tailoring (fiber), alchemy (ingredient) ihtiyaçlarına besler.

## 2. Profession Data

```yaml
# content/professions/woodcutting.yaml
code: woodcutting
category: gathering
display_name_key: profession.woodcutting.name
max_level: 100
xp_curve: { base: 50, factor: 1.10 }
tool_slot: main_hand
tool_template_codes: [basic_axe, steel_axe, mithril_axe]
race_affinity_bonus_pct:
  elf: 8
```

## 3. Node Türleri

| Node | Tier | Primary | Secondary | Rare |
|---|---|---|---|---|
| `pine_tree_node` | 1 | pine_wood | pine_resin | pine_rare |
| `oak_tree_node` | 2 | oak_wood | oak_bark | oak_rare |
| `birch_tree_node` | 2 | birch_wood | birch_bark | birch_rare |
| `mahogany_tree_node` | 3 | mahogany_wood | mahogany_bark | mahogany_rare |
| `ironwood_tree_node` | 4 | ironwood_log | ironwood_shard | ironwood_rare |
| `elder_ancient_tree_node` | 5 | elder_wood | ancient_bark | elder_rare (unique gem) |
| `heartwood_tree_node` | 6 | heartwood_log | vital_essence | heartwood_rare (legendary craft) |

## 4. Rare Drops

Örnek `oak_rare`:
```yaml
entries:
  - item_code: rare_acorn
    weight: 40
  - item_code: ancient_leaf
    weight: 30
  - item_code: forest_essence
    weight: 20
  - item_code: dryad_shard
    weight: 10
```

## 5. Race Affinity

**Elf**: `elven_grace` +5% woodcutting yield + affinity +8% = **+13% total**.

Kültürel bağlam: elf forest'lı sec.

## 6. Rank Tier'lar

Aynı pattern (Novice → Grandmaster).

Grandmaster unlock: **"Master Woodwright"** — %2 shot at rare forest essence from any tree.

## 7. Tool

- `basic_axe`: no bonus.
- `steel_axe`: +5% yield.
- `mithril_axe`: +10%.
- `orichalcum_axe`: +15% + %1 rare essence.

## 8. Woodcutting + Herbalism Cluster

Bazı forest location'lar hem tree hem herb node barındırır. Karakter tek aktivite ile bir profession'a odaklanır; ancak Phase 02'de "dual gather" (aynı anda 2 profession) opsiyonu.

## 9. AFK Aktivite

Aktivite tipi: `woodcutting` (mining paterni; bkz. [activity-engine/mining.md](../activity-engine/mining.md) analog).

## 10. Cross-Reference

- [Mining](mining.md), [Herbalism](herbalism.md)
- [world/resource-node.md](../world/resource-node.md)
- [race/racial-bonuses.md](../race/racial-bonuses.md) — elf

## 11. Metrics
- `profession_woodcutting_yield_total`.
- `profession_woodcutting_active_characters` gauge.

## 12. Bekleyen Kararlar
- Elven affinity magnitude tune.
- Forest event content (rare tree spawn) Phase 02.
- Dual gather (Phase 02).
