# Mount System (Faz 21)

Karakterin binek hayvanı. Travel hızlandırıcı, cargo taşıyıcı, koleksiyon katmanı, terrain compatibility. Pet ile paralel karakter companion.

Foundation:
- [game-design/mounts.md](../game-design/mounts.md) — vizyon.
- [database/companions.md §3-4](../database/companions.md) — schema.
- [travel/mount.md](../travel/mount.md) — travel katkısı.

## 1. Doküman Haritası

| Konu | Doküman |
|---|---|
| Genel bakış (bu doküman) | [README.md](README.md) |
| [Mount](mount.md) | Mount template + instance |
| [Level](level.md) | Mount leveling + progression |
| [Stats](stats.md) | Speed, cargo, HP, stamina |
| [Skills](skills.md) | Mount skills + abilities |
| [Travel](travel.md) | Travel integration + method |
| [Terrain](terrain.md) | Terrain compatibility + penalty |
| [Carry Capacity](carry-capacity.md) | Cargo bonus + inventory extend |

## 2. Mount Türleri

Bkz. [game-design/mounts.md §4](../game-design/mounts.md):
- **Ground** — universal; base tier.
- **Cargo** — inventory bonus; slow.
- **Aquatic** — su travel bonus.
- **Winged** (Phase 02) — cross-terrain.
- **Cosmetic** — hız aynı, kozmetik farklı.

## 3. Karakter Etkileşimi

- Karakter **1 aktif mount** (v1).
- Mount collection = koleksiyon (BoP — trade edilemez, DD-024).
- Mount summon/dismount action.
- Combat → dismount + kısa cooldown.

## 4. Pet vs Mount

- **Pet**: level + XP + skills + happiness + stamina; combat/gathering pet.
- **Mount**: level + stats + skills + terrain; travel/cargo mount.
- Farklı slot; independent.
- Aynı content pattern (template + instance).

## 5. Ortak İlkeler

1. **BoP** — bound to character (DD-024).
2. **Ölmez** — combat dismount; ancak mount HP concept (Phase 02 combat mount).
3. **Data-driven** — species content pack.
4. **AFK travel çalışır** — snapshot mount speed reveal içinde.

## 6. Endpoint Genel

```
GET  /v1/characters/{id}/mounts                   # owned mounts
POST /v1/characters/{id}/mounts/{instance_id}/activate
POST /v1/characters/{id}/mounts/{instance_id}/summon
POST /v1/characters/{id}/mounts/{instance_id}/dismount
POST /v1/characters/{id}/mounts/{instance_id}/rename
POST /v1/characters/{id}/mounts/{instance_id}/level-up    # (implicit or explicit)
GET  /v1/content/mounts
```

## 7. Cross-Reference
- [game-design/mounts.md](../game-design/mounts.md)
- [database/companions.md §3-4](../database/companions.md)
- [database/content-templates.md §22](../database/content-templates.md)
- [travel/mount.md](../travel/mount.md), [travel/terrain.md](../travel/terrain.md)
- [inventory/inventory.md](../inventory/inventory.md) — capacity extend

## 8. Ortak Bekleyen Kararlar
- Mount collection UI + compendium Phase 02.
- Mount racing mini-game Phase 03.
- Winged mount Phase 02.
- Mount breeding/hybrid Phase 03.
- Mount combat (mounted combat) Phase 03.
