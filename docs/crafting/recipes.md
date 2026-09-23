# Recipes

## 1. Amaç
Crafting'in birim tanımı. Her recipe: inputs + output(s) + time + XP + station + learn source. Content-driven.

Kaynak: [database/content-templates.md §7](../database/content-templates.md).

## 2. Model

```sql
content.recipe
  id, code (unique)
  profession_id FK
  tier_required SMALLINT
  inputs_json      -- [{item_code, qty}]
  outputs_json     -- [{item_code, qty, chance}]
  time_ms INT
  xp_reward INT
  station_required_code VARCHAR(32) NULL
  learn_source VARCHAR(32) NOT NULL  -- trainer, drop, quest, event, discovery
  version INT
```

Ek alanlar (Phase 16):
```sql
ALTER TABLE content.recipe
  ADD COLUMN description_key VARCHAR(96) NULL,
  ADD COLUMN critical_craft_chance_pct SMALLINT NOT NULL DEFAULT 5,
  ADD COLUMN quality_bias VARCHAR(16) NOT NULL DEFAULT 'balanced',  -- low/balanced/high
  ADD COLUMN max_rank SMALLINT NOT NULL DEFAULT 1,      -- Phase 02 recipe rank
  ADD COLUMN deprecated_at TIMESTAMPTZ NULL,
  ADD COLUMN tags VARCHAR(24)[] NOT NULL DEFAULT '{}';
```

## 3. Content Pack (YAML)

```yaml
# content/recipes/steel_longsword.yaml
code: steel_longsword
profession_code: blacksmithing
tier_required: 15

inputs:
  - item_code: steel_ingot
    qty: 3
  - item_code: leather_grip
    qty: 1
  - item_code: coal
    qty: 2

outputs:
  - item_code: steel_longsword_template   # item template code
    qty: 1
    chance: 1.0
  - item_code: metal_scrap                 # byproduct
    qty: 1
    chance: 0.2

time_ms: 8000
xp_reward: 40
station_required_code: forge

learn_source: trainer
critical_craft_chance_pct: 5
quality_bias: balanced

description_key: recipe.steel_longsword.description
tags: [weapon, one_hand, sword, mid_tier]
```

## 4. Recipe Kategori

Recipe içerik tipleri:
- **Equipment**: weapon/armor/accessory (output = item instance, unique roll).
- **Consumable**: potion/food/scroll (output = stackable material).
- **Material intermediate**: ore → ingot dönüşümü.
- **Enchant scroll**: enchanter output.
- **Gem cut**: raw → cut gem (jewelcrafter).
- **Recipe scroll** (Phase 02): öğretici parşömen üretimi.

## 5. Inputs

`inputs_json`:
```json
[
  { "item_code": "steel_ingot", "qty": 3, "consume_on_fail_pct": 50 },
  { "item_code": "leather_grip", "qty": 1 }
]
```

Fail sırasında material iade oranı:
- **Consume on success**: full qty consumed.
- **Consume on fail**: default %50 (config).

Override per recipe (`consume_on_fail_pct`).

## 6. Outputs

```json
[
  { "item_code": "steel_longsword_template", "qty": 1, "chance": 1.0 },
  { "item_code": "metal_scrap", "qty": 1, "chance": 0.2 }
]
```

- **Primary output**: `chance=1.0` genelde.
- **Byproduct**: proc-based (`chance<1.0`).
- **Critical output**: bkz. [critical-craft.md](critical-craft.md).

Equipment output → item_engine.Create ile yeni instance (bkz. [item-engine/item-instance.md](../item-engine/item-instance.md)).

## 7. Time & XP

- `time_ms`: bir iteration'ın süresi (8s = 8000 ms).
- `xp_reward`: iteration başına character + profession XP.
- Iteration failure: %50 XP (bkz. [crafting-xp.md](crafting-xp.md)).

## 8. Station Required

Bazı recipe workbench gerektirir:
- `forge` — blacksmith.
- `anvil` — armor.
- `alchemy_lab` — potion.
- `enchanter_table` — scroll.
- `jeweler_bench` — gem.
- `loom` — tailoring.
- `tanning_rack` — leatherworking.
- `cooking_fire` — cooking.

Station location'da bulunmalı (bkz. [world/city.md](../world/city.md), [world/outpost.md](../world/outpost.md)).

Basic recipe (starter tier) station required değil (elde işlenir); ileri tier zorunlu.

## 9. Learn Source

- **trainer**: city/village trainer NPC'de öğren; cost gold.
- **drop**: nadir mob drop scroll.
- **quest**: quest reward.
- **event**: seasonal event reward.
- **discovery**: crafting sırasında rare unlock (Phase 02).

Öğrenilmiş recipe'ler `character.recipe_known` (bkz. [database/character.md §11](../database/character.md)).

## 10. Learn Endpoint

**POST `/v1/characters/{id}/recipes/learn`**
Request:
```json
{
  "recipe_code": "steel_longsword",
  "source": "trainer",
  "trainer_npc_code": "ashfield_smith_trainer"
}
```

Server:
1. Validate:
   - Character profession level >= recipe.tier_required.
   - Recipe available at source (trainer NPC teaches; scroll owned).
2. Transaction:
   - If trainer: debit gold via ledger.
   - If scroll: consume scroll item.
   - Insert `character.recipe_known`.
   - Emit event `crafting.recipe_learned`.
3. Response 200.

## 11. Recipe Rank (Phase 02)

Recipe rank up mekaniği:
- Recipe rank 1..5.
- Higher rank → better output quality bias + reduced fail chance + rare byproduct chance.
- Rank up: recipe use N times + rare token consumption.

v1'de `max_rank=1` sabit.

## 12. Endpoint

**GET `/v1/content/recipes`**
Public; profession + tier filter.

```json
{
  "content_version": "...",
  "recipes": [
    {
      "code": "steel_longsword",
      "profession_code": "blacksmithing",
      "tier_required": 15,
      "display_name_key": "recipe.steel_longsword",
      "inputs_summary": "steel_ingot ×3, leather_grip, coal ×2",
      "outputs_summary": "steel_longsword (100%), metal_scrap (20%)",
      "time_ms": 8000,
      "xp_reward": 40,
      "learn_source": "trainer"
    }
  ]
}
```

**GET `/v1/characters/{id}/recipes-known`**
```json
{
  "recipes": [
    { "code": "steel_longsword", "rank": 1, "learned_at": "..." }
  ]
}
```

## 13. Deprecation

- Recipe deprecated → creation'da drop; mevcut karakter için known olarak korunur (backward compat).
- Aktif contract snapshot content → deprecated recipe kararlı sonuç.

## 14. Metrics
- `recipe_learned_total{recipe}`.
- `recipe_crafted_total{recipe}`.
- `recipe_content_active_total` gauge.

## 15. Audit
- `content.recipe_created/updated/deprecated`.
- `character.recipe_learned`.

## 16. Cross-Reference
- [Materials](materials.md), [Crafting Queue](crafting-queue.md), [Quality](quality.md), [Critical Craft](critical-craft.md), [Mastery](mastery.md), [Crafting XP](crafting-xp.md)
- [activity-engine/crafting.md](../activity-engine/crafting.md)
- [database/content-templates.md §7](../database/content-templates.md), [database/character.md §11](../database/character.md)
- [item-engine/item-instance.md](../item-engine/item-instance.md)
- [game-design/professions.md](../game-design/professions.md)

## 17. Bekleyen Kararlar
- Recipe rank aktivasyon Phase 02.
- Recipe discovery mekaniği (crafting sırasında rare unlock) Phase 02.
- Recipe scroll trade (BoE) Phase 02.
- Guild-exclusive recipe (guild research reward) Phase 02.
