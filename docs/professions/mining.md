# Mining

## 1. Amaç
Cevher/taş/gem çıkarma profession'u. Ana material üreticisi (blacksmithing, jewelcrafting için).

Aktivite detay: [activity-engine/mining.md](../activity-engine/mining.md).

## 2. Profession Data

```yaml
# content/professions/mining.yaml
code: mining
category: gathering
display_name_key: profession.mining.name
max_level: 100
xp_curve:
  base: 50
  factor: 1.10
tool_slot: main_hand              # pickaxe
tool_alternate_slots: [off_hand]  # (opsiyonel dual pickaxe)
tool_template_codes:
  - basic_pickaxe
  - steel_pickaxe
  - mithril_pickaxe
race_affinity_bonus_pct:
  dwarf: 8
```

## 3. Node Türleri

| Node Code | Tier | Primary Output | Secondary | Rare Drop Table |
|---|---|---|---|---|
| `copper_ore_node` | 1 | copper_ore | gray_stone | copper_rare |
| `iron_ore_node` | 2 | iron_ore | gray_stone | iron_rare |
| `silver_ore_node` | 3 | silver_ore | silver_stone | silver_rare |
| `mithril_ore_node` | 4 | mithril_ore | mithril_shard | mithril_rare |
| `adamantite_node` | 5 | adamantite_ore | adamantite_shard | adamantite_rare |
| `orichalcum_node` | 6 | orichalcum_ore | ancient_stone | orichalcum_rare |

Tier belirler: profession min_level (Novice/Apprentice/... eşleşmesi).

Node content: [world/resource-node.md](../world/resource-node.md).

## 4. Rare Drop Kataloğu (Örnek)

```yaml
# content/loot_tables/iron_rare.yaml
code: iron_rare
source_type: node
entries:
  - item_code: rough_ruby
    weight: 40
    min_qty: 1
    max_qty: 1
  - item_code: rough_sapphire
    weight: 35
  - item_code: rough_emerald
    weight: 20
  - item_code: ancient_iron_relic
    weight: 5
```

Rare drop rate node'a bağlı (`rare_proc_per_hour`).

## 5. Rank Tier'ları

| Level | Tier | Yield Bonus |
|---|---|---|
| 1-19 | Novice | 0% |
| 20-39 | Apprentice | +5% |
| 40-59 | Journeyman | +10% |
| 60-79 | Expert | +15% |
| 80-99 | Master | +20% |
| 100 | Grandmaster | +25% + unique recipe |

## 6. Grandmaster Unlock

Level 100 → unlock unique recipe or ability:
- **Grandmaster Miner**: "Master of Stone" title + %2 shot at gem drop from any ore node.
- Unique content patch item.

## 7. Race Affinity

**Dwarf**: `dwarven_endurance` passive → +5% mining yield always; +additional +8% from profession affinity content field.

Combined: +13% mining yield baseline dwarf miner.

## 8. Tool Bonus

- `basic_pickaxe`: no bonus.
- `steel_pickaxe`: +5% yield.
- `mithril_pickaxe`: +10% yield.
- `orichalcum_pickaxe`: +15% + %1 rare gem chance.

Equip in main_hand (weapon slot alternate — combat weapon replace edilir; character player choice).

## 9. Endpoint

**GET `/v1/content/professions/mining`**
```json
{
  "code": "mining",
  "display_name_key": "profession.mining.name",
  "max_level": 100,
  "nodes": [ { "code": "copper_ore_node", "tier": 1, "min_level": 1 }, ... ],
  "tools": [ { "code": "basic_pickaxe", "yield_bonus_pct": 0 }, ... ],
  "race_affinity_bonus_pct": { "dwarf": 8 }
}
```

## 10. AFK Aktivite

Bkz. [activity-engine/mining.md](../activity-engine/mining.md).

Karakter mining spot'ta AFK → analitik yield.

## 11. Cross-Reference

- [Woodcutting](woodcutting.md), [Fishing](fishing.md), [Herbalism](herbalism.md)
- [activity-engine/mining.md](../activity-engine/mining.md)
- [world/resource-node.md](../world/resource-node.md)
- [race/racial-bonuses.md](../race/racial-bonuses.md) — dwarven

## 12. Metrics
- `profession_mining_yield_total`.
- `profession_mining_active_characters` gauge.
- `mining_rare_gem_drops_total`.

## 13. Bekleyen Kararlar
- Grandmaster unique reward final.
- Tool durability drain v1 vs Phase 02.
- Ore tier balance simulation.
- Multi-node cluster gather (Phase 02).
