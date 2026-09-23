# Pet System (Faz 20)

Karakter yanında sürekli çalışan companion — combat, gathering, travel, koleksiyon. Uzun süreli progression + kimlik.

Foundation:
- [game-design/pets.md](../game-design/pets.md) — vizyon.
- [database/companions.md §1-2](../database/companions.md) — schema.

## 1. Doküman Haritası

| Konu | Doküman |
|---|---|
| Genel bakış (bu doküman) | [README.md](README.md) |
| [Species](species.md) | Pet species template + rarity |
| [Instance](instance.md) | Owned pet instance + ownership |
| [Level](level.md) | Pet level progression + cap |
| [XP](xp.md) | Pet XP sources + curve |
| [Traits](traits.md) | Trait system + unique modifiers |
| [Skills](skills.md) | Pet abilities + slots |
| [Evolution](evolution.md) | Species evolution + conditions |
| [Happiness](happiness.md) | Happiness stat + feed |
| [Stamina](stamina.md) | Pet stamina + rest |
| [Combat](combat.md) | Pet combat behavior + AI |
| [Hunting](hunting.md) | Hunting contribution + species affinity |

## 2. Pet Türleri

Bkz. [game-design/pets.md §4](../game-design/pets.md):
- **Combat Pets** — DPS/tank/heal minor contribution.
- **Gathering Pets** — profession yield bonus.
- **Utility Pets** — bag/speed/luck.
- **Cosmetic Pets** — sadece görsel.

## 3. Ana Yapı

```
Species (template — content.pet_template)
  └── Instance (owned — companion.pet_owned)
       ├── Level + XP
       ├── Traits (random on hatch/tame)
       ├── Skills (species base + rank up)
       ├── Happiness (0..100)
       ├── Stamina (0..100)
       └── Evolution stage (Phase 02)
```

## 4. Ortak İlkeler

1. **Bound to character** (DD-024): trade edilemez.
2. **Aktif slot 2** (v1): 1 combat + 1 utility.
3. **Ölmez** — combat'ta downed → cooldown revive.
4. **AFK contribute** — karakter aktivite yield hesabına dahil.
5. **Data-driven** — species content pack.

## 5. Endpoint Genel

```
GET  /v1/characters/{id}/pets
POST /v1/characters/{id}/pets/{instance_id}/activate     # slot atama
POST /v1/characters/{id}/pets/{instance_id}/deactivate
POST /v1/characters/{id}/pets/{instance_id}/rename
POST /v1/characters/{id}/pets/{instance_id}/feed          # happiness
POST /v1/characters/{id}/pets/{instance_id}/rest          # stamina
POST /v1/characters/{id}/pets/{instance_id}/evolve        # Phase 02
POST /v1/characters/{id}/pets/{instance_id}/retire        # stable
GET  /v1/content/pet-species
GET  /v1/content/pet-species/{code}
```

## 6. Cross-Reference
- [game-design/pets.md](../game-design/pets.md)
- [database/companions.md](../database/companions.md)
- [database/content-templates.md §21](../database/content-templates.md)
- [combat/README.md](../combat/README.md)
- [activity-engine/hunting.md](../activity-engine/hunting.md)

## 7. Ortak Bekleyen Kararlar
- Pet compendium Phase 02 (koleksiyon).
- Pet battle mini-game Phase 03.
- Cross-account BoA pets Phase 03.
- Pet breeding (Phase 02+) — new species birth.
