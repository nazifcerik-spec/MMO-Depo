# Hunting + Adventure (Faz 18)

Karakterin dünyayla combat + explore ilişkisinin somut mekanikleri. Combat engine üstünde encounter, rare/elite, treasure, random event, adventure story flow.

Foundation:
- [combat/README.md](../combat/README.md) — combat engine.
- [activity-engine/hunting.md](../activity-engine/hunting.md), [activity-engine/adventure.md](../activity-engine/adventure.md), [activity-engine/hunting.md](../activity-engine/hunting.md).
- [world/monster-area.md](../world/monster-area.md), [world/location.md](../world/location.md).

## 1. Doküman Haritası

| Konu | Doküman |
|---|---|
| Genel bakış (bu doküman) | [README.md](README.md) |
| [Hunting](hunting.md) | Hunt aktivite deep-dive + filter + party |
| [Encounter](encounter.md) | Random encounter framework |
| [Rare Enemy](rare-enemy.md) | Rare mob spawn + kill credit + reward |
| [Elite Enemy](elite-enemy.md) | Elite mob + party recommendation |
| [Treasure](treasure.md) | Treasure chest + world drop + hidden loot |
| [Random Event](random-event.md) | Random event injection + rewards |
| [Adventure](adventure.md) | Story adventure + stage-based |
| [Exploration](exploration.md) | Discovery + fog of war + region mastery |

## 2. Katman Diyagramı

```
Location (world)
   │
   ├── Monster Area  ─┬── Hunting activity → combat → kill loot
   │                  ├── Rare enemy spawn
   │                  └── Elite enemy
   │
   ├── Encounter Pool ─┬── Random encounter (mid-activity)
   │                   └── Treasure spawn
   │
   ├── Adventure Anchor ─── Adventure activity → stage checks → narration
   │
   └── Random Event ─── injected event (server-wide veya per-region)
```

## 3. Content-Driven

Tüm mekanik YAML content:
```
content/monster_area_spawns/*.yaml       (bkz. world/monster-area.md)
content/adventures/*.yaml                 (bkz. adventure.md)
content/adventure_stages/*.yaml
content/encounters/*.yaml
content/treasures/*.yaml
content/rare_mob_spawns/*.yaml
content/random_events/*.yaml
```

## 4. Ortak İlkeler

1. **AFK-first** — hunting + adventure AFK contract; encounter mid-reveal.
2. **Analytical yield** — DD-084; per-encounter simulation değil, Poisson sample.
3. **Deterministic RNG** — seed-based; replay + audit.
4. **Content-driven** — YAML pack + hot reload.
5. **Cross-layer** — encounter hunting/adventure/travel içine inject edilebilir.

## 5. Cross-Reference

- [combat/README.md](../combat/README.md)
- [activity-engine/hunting.md](../activity-engine/hunting.md), [activity-engine/adventure.md](../activity-engine/adventure.md)
- [world/monster-area.md](../world/monster-area.md), [world/location.md](../world/location.md)
- [game-design/hunting.md](../game-design/hunting.md), [game-design/adventure.md](../game-design/adventure.md)
- [systems/loot-system.md](../systems/loot-system.md)

## 6. Ortak Bekleyen Kararlar

- Encounter pool granularity (per location vs per region).
- Rare mob global anons (Phase 02).
- Adventure branching narrative (Phase 02).
- Treasure competition (personal vs shared).
- Random event pacing rate.
