# Professions (Faz 15) — Gathering Emphasized

Karakterin ekonomik/uzmanlaşma katmanı. AFK aktivite motoruna oturur, ekonomiye kaynak besler. Bu klasör **gathering + primary profession** deep-dive.

Yüksek seviye vizyon: [game-design/professions.md](../game-design/professions.md). Bu klasör Faz 15 implementation-hazır detay.

## 1. Profession Kataloğu (v1 + Phase 02)

Bu faz odak — **gathering + primary specialist**:

| Kod | Kategori | v1 | Doküman |
|---|---|---|---|
| `mining` | Gathering | ✓ | [mining.md](mining.md) |
| `woodcutting` | Gathering | ✓ | [woodcutting.md](woodcutting.md) |
| `fishing` | Gathering | ✓ | [fishing.md](fishing.md) |
| `hunting` | Gathering + Combat | ✓ | [hunting.md](hunting.md) |
| `herbalism` | Gathering | ✓ | [herbalism.md](herbalism.md) |
| `farming` | Gathering + Production | Phase 02 | [farming.md](farming.md) |
| `skinning` | Gathering (yan aktivite) | ✓ | [skinning.md](skinning.md) |
| `archaeology` | Gathering (unique) | Phase 02 | [archaeology.md](archaeology.md) |

Crafting professions (blacksmithing/leatherworking/tailoring/alchemy/cooking/enchanting/jewelcrafting) ayrı bir faz — bu klasörde değil.

## 2. Ortak Profession Modeli

Kaynak: [database/content-templates.md §6](../database/content-templates.md), [database/character.md §10](../database/character.md).

- **Template**: `content.profession` — code, category, max_level, xp_curve.
- **Character State**: `character.profession` — level, xp, is_active.

Karakter aynı anda **3 profession slot** (v1). Yeni öğrenme trainer'da; unlearn 24h cooldown (Phase 02).

## 3. Ortak Sözleşme

Her profession:
- **AFK-friendly** — Activity engine tarafında aktivite tipi.
- **Data-driven** — content pack YAML.
- **Deterministic yield** — analitik hesap (bkz. [activity-engine/offline-progress.md](../activity-engine/offline-progress.md)).
- **Modifier stack integration** — race passive, mastery, gear bonus.
- **Node/spawn model** — location-based content.

## 4. Level & XP Curve

Level range: 1..100 (Novice/Apprentice/Journeyman/Expert/Master/Grandmaster tier'lar).

XP curve (per profession):
```
xp_to_next(L) = base * pow(factor, L-1)
base = 50
factor = 1.10
```

Grandmaster (100) ≈ 700k cumulative XP.

## 5. Yield Multiplier Katmanları

Ortak:
```
yield = base_rate
      × (1 + gather_yield_pct_modifiers/100)
      × (1 + mastery_bonus_fraction)
      × race_affinity_multiplier
      × tool_bonus
      × event_bonus (etkin ise)
```

Detay her profession'da.

## 6. Tool Requirements

Bazı profession gathering tool gerektirir:
- **Mining**: pickaxe (weapon slot alternate).
- **Woodcutting**: axe.
- **Fishing**: fishing rod.
- **Herbalism**: sickle (opsiyonel).
- **Skinning**: skinning knife.
- **Farming**: hoe + seeds.
- **Archaeology**: brush + shovel.

Tool durability drain (v1 opsiyonel; genelde skip).

## 7. Endpoint Genel

```
GET  /v1/characters/{id}/professions              # list active + levels
POST /v1/characters/{id}/professions/learn        # trainer'da öğren
POST /v1/characters/{id}/professions/unlearn      # (Phase 02)
GET  /v1/content/professions                      # tümü discovery
GET  /v1/content/professions/{code}               # detay + recipe/node listesi
```

## 8. Race Affinity

Her ırk için 1 profession affinity (bkz. [race/racial-bonuses.md](../race/racial-bonuses.md)):
- Human: trading (special).
- Dwarf: mining (+8%).
- Elf: woodcutting (+8%).
- Orc: hunting (+8%).
- Merfolk: fishing (+8%).
- Arcanborn: enchanting (+8% success; crafting affinity).

## 9. Cross-Reference

- [game-design/professions.md](../game-design/professions.md)
- [database/character.md §10-11](../database/character.md), [database/content-templates.md §6-7](../database/content-templates.md)
- [activity-engine/mining.md](../activity-engine/mining.md), [activity-engine/fishing.md](../activity-engine/fishing.md), [activity-engine/hunting.md](../activity-engine/hunting.md), [activity-engine/crafting.md](../activity-engine/crafting.md)
- [world/resource-node.md](../world/resource-node.md), [world/monster-area.md](../world/monster-area.md)
- [race/racial-bonuses.md](../race/racial-bonuses.md)
- [skill-engine/mastery.md](../skill-engine/mastery.md) — profession mastery track

## 10. Ortak Bekleyen Kararlar

- Tool durability drain v1 vs Phase 02.
- Trading profession Phase 02 tam tanım.
- Profession specialization sub-tree (Phase 02).
- Cross-profession sinerji (mining + jewelcrafting bundle discount) Phase 02.
- Profession lockout (aynı iki oyuncu farklı profession — pazar sağlığı) Phase 03 değerlendirme.
