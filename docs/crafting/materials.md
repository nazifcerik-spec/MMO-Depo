# Materials

## 1. Amaç
Crafting input'larının modellenmesi: tier'lar, kategori, stack, convert, sink/source.

## 2. Material Kategorileri

| Kategori | Örnek | Kaynak |
|---|---|---|
| **Raw ore** | copper_ore, iron_ore | mining |
| **Ingot / bar** | copper_ingot, iron_ingot | smelting recipe (blacksmith) |
| **Wood log** | pine_wood, oak_wood | woodcutting |
| **Processed wood** | pine_plank, oak_grip | carpentry recipe (blacksmith/leather) |
| **Leather** | leather_scrap, thick_hide | skinning |
| **Cloth** | linen_scrap, silk_thread | tailoring input (drop + herbalism) |
| **Herb** | red_thistle, moonleaf | herbalism |
| **Alchemy essence** | fire_essence, life_essence | intermediate |
| **Gem raw** | rough_ruby, rough_sapphire | mining rare drop |
| **Gem cut** | cut_ruby (tier 1..5) | jewelcrafting |
| **Cooking ingredient** | wheat, meat_raw, spice | farming/hunting/vendor |
| **Enchant dust** | arcane_dust, magic_shard | disenchanting |
| **Coal / fuel** | coal, magic_ember | mining/vendor |
| **Currency-material** | dungeon_token, event_token | dungeon/event |

## 3. Tier'lar

Her material kategori 5-6 tier:
```
Tier 1: copper_ore → copper_ingot
Tier 2: iron_ore → iron_ingot
Tier 3: silver_ore → silver_ingot
Tier 4: mithril_ore → mithril_ingot
Tier 5: adamantite_ore → adamantite_ingot
Tier 6: orichalcum_ore → orichalcum_ingot
```

Recipe input material tier'ı önemli — daha yüksek tier armor daha yüksek tier ingot ister.

## 4. Item Template

Material = normal item; `content.item_template.category='material'`, `subcategory=<kind>`.
- `stack_max`: material için 200 default; nadir material için 50.
- `weight`: 1 (v1 kapalı).
- `binding_type='no_bind'` — trade edilebilir.
- `vendor_sell_price_copper`: base value.

## 5. Convert Recipe (Örnek Smelting)

Raw → processed:
```yaml
# content/recipes/smelt_iron.yaml
code: smelt_iron
profession_code: blacksmithing
tier_required: 15
inputs:
  - item_code: iron_ore
    qty: 2
  - item_code: coal
    qty: 1
outputs:
  - item_code: iron_ingot
    qty: 1
    chance: 1.0
time_ms: 3000
xp_reward: 10
station_required_code: forge
learn_source: trainer   # tier 15 blacksmith trainer'ında öğrenilir
```

Convert recipe genellikle:
- Time short (3-5 sec).
- XP low.
- %100 chance (basic operation).

## 6. Cross-Profession Material Flow

```
Mining → ore
  ↓
Blacksmithing smelt → ingot
  ↓
Blacksmithing craft → weapon/armor

Woodcutting → wood
  ↓
Carpentry (blacksmith recipe) → grip
  ↓
Blacksmithing craft → weapon (with grip)

Skinning → hide
  ↓
Leatherworking → leather scrap
  ↓
Tailoring input? no, leather → leatherworker armor
```

Cross-profession sinerji ekonomiye derinlik verir.

## 7. Material Slot Assignment

Recipe input her material için envanter'de aranır:
```go
func gatherRecipeInputs(char, recipe) []InventorySlot {
    for _, in := range recipe.Inputs {
        slots := inventory.FindMaterialSlots(char.ID, in.item_code)
        needed := in.qty * iterations
        // Iterate; take from stacks.
    }
}
```

Multi-stack support: aynı material farklı slot'larda toplu tüket.

## 8. Bank / Guild Bank Withdraw

- Craft aktivite başlatırken envanterde yeterli material olmalı.
- Bank/guild bank'tan withdraw manuel (v1); Phase 02 auto-withdraw opsiyonu.

## 9. Endpoint

**GET `/v1/content/materials`**
Public; filter by category.
```json
{
  "materials": [
    { "code": "iron_ore", "tier": 2, "category": "raw_ore", "vendor_price_copper": 250 },
    { "code": "iron_ingot", "tier": 2, "category": "ingot", "vendor_price_copper": 800 }
  ]
}
```

## 10. Material Sink & Source

Faucet:
- Gathering profession yield.
- Mob drop (leather, cloth).
- Vendor buy (coal, basic ingredients).
- Marketplace.

Sink:
- Crafting consume (primary sink).
- Salvage input (reverse — Phase 02 disenchanting).
- Vendor sell (junk material).

Balance monitor: [database/economy-market.md §6](../database/economy-market.md).

## 11. Material Rarity

- **Common**: iron_ore, wood, herb.
- **Uncommon**: silver_ore, hide, cotton.
- **Rare**: mithril_ore, dragon_hide, moonleaf.
- **Epic**: adamantite_ore, ancient_essence.
- **Legendary**: orichalcum_ore, precursor_shard.

Rarity template.rarity ile eşleşir; vendor price + drop rate.

## 12. Metrics
- `material_produced_total{code,source}`.
- `material_consumed_total{code,recipe}`.
- `material_market_price_avg{code}` — [database/economy-market.md §6](../database/economy-market.md).

## 13. Cross-Reference
- [Recipes](recipes.md), [Crafting Queue](crafting-queue.md)
- [item-engine/item-template.md](../item-engine/item-template.md), [item-engine/rarity.md](../item-engine/rarity.md)
- [inventory/stacking.md](../inventory/stacking.md)
- [game-design/professions.md](../game-design/professions.md), [game-design/economy.md](../game-design/economy.md)
- [professions/README.md](../professions/README.md)

## 14. Bekleyen Kararlar
- Auto-withdraw from bank Phase 02.
- Material salvage (disenchant Phase 02).
- Cross-profession material bundle (crafter combo discount) Phase 02.
- Rare material seasonal event drops Phase 02.
