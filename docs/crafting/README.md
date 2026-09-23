# Crafting (Faz 16)

Recipe-based üretim aktivitesi. Gathering profession'larının çıktıları burada işlenir; ekonomiye equipment/consumable/enchant besler. AFK-friendly.

Foundation: [activity-engine/crafting.md](../activity-engine/crafting.md) (aktivite katmanı). Yüksek seviye: [game-design/professions.md](../game-design/professions.md).

## 1. Crafting Profession Kataloğu

7 crafting profession (v1):
- `blacksmithing` — silah, ağır zırh (plate/mail).
- `leatherworking` — orta/hafif zırh (leather), aksesuar.
- `tailoring` — kumaş zırh, bag.
- `alchemy` — potion, elixir, tonic.
- `cooking` — food buff'ları.
- `enchanting` — enchant scroll + disenchant.
- `jewelcrafting` — gem cut + ring/trinket.

Gathering profession'lar farklı klasör: [/docs/professions/](../professions/README.md).

## 2. Doküman Haritası

| Konu | Doküman |
|---|---|
| Genel bakış (bu doküman) | [README.md](README.md) |
| [Recipes](recipes.md) | Recipe entity + template + öğrenme + trainer |
| [Materials](materials.md) | Material tier + slot + convert + stack |
| [Crafting Queue](crafting-queue.md) | Batch iteration + activity contract |
| [Quality](quality.md) | Output quality (poor/normal/exceptional) + rarity |
| [Critical Craft](critical-craft.md) | Crit chance + rare output + big miss |
| [Mastery](mastery.md) | Craft mastery track + rank up + unique unlock |
| [Crafting XP](crafting-xp.md) | XP formula + fail penalty + level up |

## 3. Ortak Sözleşme

Her crafting profession:
- **Recipe-driven**: content.recipe row.
- **Iteration-based**: batch (1..N iteration).
- **AFK**: activity contract; per-tick simulation YOK (DD-084).
- **Success chance**: base + WIS + tool + rank.
- **Station requirement** (opsiyonel): forge, workbench, alchemy lab.
- **Output**: item instance (equipment) veya stack (material/consumable).

## 4. Ana Interface (Yüksek Seviye)

```go
package craftingengine

type Engine interface {
    // Recipe query
    CanCraft(ctx context.Context, char Character, recipeID uuid.UUID, iterations int32) (CanCraftResult, error)
    Craft(ctx context.Context, req CraftRequest) (Contract, error)     // AFK-start
    // Reveal integration analitik (bkz. activity-engine/crafting.md)
    // Mastery grant
    GrantMasteryXP(ctx context.Context, char, trackID, xp int32) error
    // Recipe learn
    LearnRecipe(ctx context.Context, char, recipeCode string, source LearnSource) error
}
```

## 5. Cross-Reference

- [activity-engine/crafting.md](../activity-engine/crafting.md) — activity kontratı + reveal
- [database/content-templates.md §7](../database/content-templates.md) — recipe schema
- [item-engine/README.md](../item-engine/README.md) — output item instance
- [inventory/inventory.md](../inventory/inventory.md) — input consume + output place
- [character/attributes.md](../character/attributes.md) — WIS impact craft success
- [game-design/professions.md](../game-design/professions.md) — vizyon
- [skill-engine/mastery.md](../skill-engine/mastery.md) — mastery track pattern

## 6. Ortak İlkeler

1. **AFK-first** — crafting AFK aktivite; oyuncu logout devam.
2. **Analitik yield** — batch binomial sample (per-iteration loop yok).
3. **Deterministic** — seed'li; replay + audit.
4. **Server-authoritative** — output item instance server tarafında yaratılır.
5. **Idempotent iteration** — aynı iteration_id iki kez fire etmez.
6. **Data-driven** — recipe YAML; hot reload.

## 7. Endpoint Genel

```
GET  /v1/content/recipes                       # tümü discovery (public)
GET  /v1/content/recipes/{code}                # detay
GET  /v1/characters/{id}/recipes-known         # karakterin öğrendikleri
POST /v1/characters/{id}/recipes/learn         # trainer/scroll ile öğren
POST /v1/characters/{id}/craft                 # activity başlat (bkz. activity-engine/crafting.md)
POST /v1/characters/{id}/craft/{cid}/cancel
GET  /v1/characters/{id}/craft/preview         # yield tahmin
```

## 8. Bekleyen (Ortak) Kararlar

- Recipe rank Phase 02 (recipe rank up ile output kalitesi ↑).
- Batch craft cap iteration sayısı (v1 max 200).
- Guild recipe (guild-only crafts) Phase 02.
- Crafting order (player-to-player custom order) Phase 02+.
- Auto-restock feature Phase 03.
