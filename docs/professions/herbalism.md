# Herbalism (Foraging)

## 1. Amaç
Ot/mantar/mineral toplama. Alchemy (potion, elixir), cooking (spice), enchanting (dust) besler.

Not: `foraging` bazı öncekilerde "generic gathering" olarak geçti; bu klasörde **herbalism** ismini kullanır (spesifik: herbs). Diğer non-herb foraging (mushroom, moss) aynı profession.

## 2. Profession Data

```yaml
# content/professions/herbalism.yaml
code: herbalism
category: gathering
display_name_key: profession.herbalism.name
max_level: 100
xp_curve: { base: 50, factor: 1.10 }
tool_slot: null                    # tool opsiyonel (sickle Phase 02)
tool_template_codes: []
race_affinity_bonus_pct:
  elf: 5                            # secondary affinity (elf primary woodcutting)
  arcanborn: 5                      # magical herbs bias
```

## 3. Node Türleri

| Node | Tier | Primary | Secondary | Rare |
|---|---|---|---|---|
| `common_grass_herb` | 1 | common_herb | wild_flower | – |
| `red_thistle_node` | 1 | red_thistle | thistle_seed | – |
| `moonleaf_node` | 2 | moonleaf | moon_essence | – |
| `sunroot_node` | 2 | sunroot | sun_shard | – |
| `whisper_mushroom_node` | 3 | whisper_mushroom | fungal_dust | mushroom_rare |
| `frost_bloom_node` | 3 | frost_bloom | ice_crystal | – |
| `arcane_lotus_node` | 4 | arcane_lotus | arcane_petal | lotus_rare (legendary elixir) |
| `heartstem_node` | 5 | heartstem | vital_pollen | heartstem_rare |

## 4. Node Distribution by Biome

- **Forest**: red_thistle, moonleaf, whisper_mushroom.
- **Plains**: common_herb, sunroot.
- **Tundra**: frost_bloom.
- **Swamp**: whisper_mushroom, unique swamp herbs.
- **Magical zones** (Phase 02): arcane_lotus.

Region biome content-driven.

## 5. Rare Drops

Rare herb → legendary alchemy recipe ingredients (Phase 02).

## 6. Race Affinity

- **Elf**: secondary affinity (+5%; primary woodcutting).
- **Arcanborn**: magical herb bias (+5%).

Diğer race normal (+0%).

## 7. Tool (Sickle Phase 02)

v1 tool yok. Phase 02:
- **Basic sickle**: +5% yield.
- **Silver sickle**: +10%.
- Sickle enchant → herb drop rate bonus.

## 8. AFK Aktivite

Aktivite tipi: `herbalism` (mining paterni analog; bkz. [activity-engine/mining.md](../activity-engine/mining.md)).

## 9. Cross-Reference

- [Mining](mining.md), [Farming](farming.md)
- [activity-engine/mining.md](../activity-engine/mining.md) — pattern
- [world/resource-node.md](../world/resource-node.md)
- [race/racial-bonuses.md](../race/racial-bonuses.md)

## 10. Metrics
- `profession_herbalism_yield_total`.
- `herbalism_rare_drop_total`.

## 11. Bekleyen Kararlar
- Sickle tool Phase 02.
- Alchemy ingredient content depth (Phase 02).
- Seasonal herb (Phase 02 seasonal drop).
- Cross-node cluster gather Phase 02.
