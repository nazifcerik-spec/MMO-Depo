# Farming

> ⚠️ **v1'de kapalı**; **Phase 02** aktivasyon planlı. Schema + design hazır.

## 1. Amaç
Karakter kendi arazisinde bitki/hayvan yetiştirir. Sürekli passive yield (uzun cycle); ekonomiye özel material besler.

## 2. Neden Phase 02?

- **Housing** sistemi ile bağlı (farmland = personal plot).
- Uzun cycle (saatler-günler) balance testi gerekir.
- Content pack (seed, crop, animal) büyük yatırım.
- v1 core loop öncelik.

## 3. Model (Phase 02 Schema)

```sql
CREATE TABLE character.farm_plot (
  id UUID PRIMARY KEY,
  character_id UUID NOT NULL,
  location_id UUID NOT NULL,        -- housing location
  slot_index SMALLINT NOT NULL,      -- plot içi slot 0..N
  crop_template_id UUID NULL,        -- ekili crop
  planted_at TIMESTAMPTZ NULL,
  next_harvest_at TIMESTAMPTZ NULL,
  water_last_at TIMESTAMPTZ NULL,
  fertilize_last_at TIMESTAMPTZ NULL,
  version INT NOT NULL DEFAULT 1
);

CREATE TABLE content.farm_crop_template (
  id UUID PRIMARY KEY,
  code VARCHAR(48) UNIQUE NOT NULL,
  display_name_key VARCHAR(64),
  seed_item_code VARCHAR(48) NOT NULL,
  grow_time_ms BIGINT NOT NULL,       -- örn 6h
  water_cycle_ms BIGINT NOT NULL,     -- her N saat water
  yield_item_code VARCHAR(48) NOT NULL,
  yield_qty_min SMALLINT NOT NULL,
  yield_qty_max SMALLINT NOT NULL,
  rare_drop_item_code VARCHAR(48) NULL,
  rare_drop_chance_pct SMALLINT DEFAULT 0
);
```

## 4. Farming Flow (Phase 02)

```
1. Character has housing plot (unlocked).
2. Plant seed (POST /v1/characters/{id}/farm/plot/{slot}/plant).
3. Water periodically (opsiyonel; skip = %50 yield penalty).
4. Fertilize (opsiyonel; +%20 yield).
5. Wait grow_time.
6. Harvest (POST .../harvest) → item yield.
```

## 5. Crops (Örnek)

| Crop | Seed | Grow Time | Yield |
|---|---|---|---|
| `wheat_crop` | wheat_seed | 6h | wheat ×3-5 |
| `carrot_crop` | carrot_seed | 4h | carrot ×2-4 |
| `magic_herb_crop` | magic_herb_seed | 24h | magic_herb ×1-2 + rare essence |
| `dragon_fruit_crop` | dragon_fruit_seed | 48h | dragon_fruit ×1 (rare) |

## 6. Animal Husbandry (Phase 02+)

- Chicken (eggs, hourly).
- Cow (milk, 4h cycle).
- Sheep (wool, 24h cycle).

Ayrı schema (`character.animal_pen`).

## 7. Race Affinity

Elf / Merfolk / Arcanborn için nature affinity (Phase 02 balance).

## 8. AFK Model

Farming doğuştan AFK: plant → wait → harvest. Aktivite contract'ı **yok** (background timer). Reveal analog: login'de "crops ready" notification.

Alternatif model: farming aktivite contract'ı (channel=farming, primary'yi bloke etmez) — Phase 02 karar.

## 9. Endpoint (Phase 02)

```
GET   /v1/characters/{id}/farm/plot
POST  /v1/characters/{id}/farm/plot/{slot}/plant
POST  /v1/characters/{id}/farm/plot/{slot}/water
POST  /v1/characters/{id}/farm/plot/{slot}/fertilize
POST  /v1/characters/{id}/farm/plot/{slot}/harvest
DELETE /v1/characters/{id}/farm/plot/{slot}       # uproot (yield lose)
```

## 10. Cross-Profession

Farming output:
- Wheat → cooking (bread).
- Magic herb → alchemy.
- Cotton → tailoring.

## 11. Cross-Reference
- [Herbalism](herbalism.md), [Cooking](../game-design/professions.md) (crafting)
- Housing sistem (Phase 02).

## 12. Bekleyen Kararlar
- Aktivasyon Phase 02 milestone.
- Animal husbandry Phase 02 vs Phase 03.
- Multi-plot expansion (housing tier ile).
- Community farming (guild farm) Phase 03.
- Watering can consumable / cooldown.