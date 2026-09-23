# Species (Pet Template)

## 1. Amaç
Pet'in kalıcı tanımı — species, base stats, ability set, evolution path.

Kaynak: [database/content-templates.md §21](../database/content-templates.md).

## 2. Model

Phase 20 eklemeleri:
```sql
ALTER TABLE content.pet_template
  ADD COLUMN species_family VARCHAR(24) NOT NULL DEFAULT 'beast',
  ADD COLUMN evolution_stage SMALLINT NOT NULL DEFAULT 1,
  ADD COLUMN next_evolution_template_id UUID NULL,
  ADD COLUMN evolution_condition_json JSONB NULL,
  ADD COLUMN base_traits_pool_json JSONB NOT NULL DEFAULT '[]',
  ADD COLUMN base_skill_codes VARCHAR(48)[] NOT NULL DEFAULT '{}',
  ADD COLUMN max_level SMALLINT NOT NULL DEFAULT 50,
  ADD COLUMN xp_curve_json JSONB NOT NULL DEFAULT '{}',
  ADD COLUMN acquisition_source VARCHAR(32) NOT NULL DEFAULT 'drop',
  ADD COLUMN deprecated_at TIMESTAMPTZ NULL;
```

## 3. Content Pack

```yaml
# content/pet_species/wolf_cub.yaml
code: wolf_cub
species_family: beast
type: combat
rarity: uncommon
display_name_key: pet.wolf_cub.name
lore_key: pet.wolf_cub.lore
art_ref: pets/wolf_cub.png

max_level: 50
xp_curve:
  base: 80
  factor: 1.12

base_stats:
  attack_power: 15
  hp: 200
  crit_pct: 5

active_ability_code: pet_bite
passive_ability_code: pet_pack_hunter    # gathering-hint

base_traits_pool:                        # trait roll pool on hatch
  - trait_code: keen_senses
    weight: 40
  - trait_code: swift_paw
    weight: 30
  - trait_code: iron_hide
    weight: 20
  - trait_code: pack_leader
    weight: 10

base_skill_codes: [pet_bite, pet_growl]

evolution_stage: 1
next_evolution_template_id: alpha_wolf   # evolves to alpha wolf
evolution_condition:
  level_min: 30
  bond_level_min: 50
  quest_completed: pack_leader_quest     # Phase 02

acquisition_source: drop                 # drop/quest/event/starter/purchase

tags: [beast, ranger_friendly, dps]
```

## 4. Species Family

| Family | Örnek | Combat Bias |
|---|---|---|
| `beast` | wolf, bear, hawk | DPS/utility |
| `insect` | spider, scorpion | DoT/poison |
| `reptile` | lizard, snake | armor/tank |
| `avian` | falcon, owl | scout/speed |
| `aquatic` | crab, seahorse | fishing bonus |
| `elemental` | fire wisp, ice sprite | magic |
| `undead` | ghost, skeleton | dark magic |
| `mystical` | phoenix, unicorn | rare/unique |

Family = bestiary track ilişki.

## 5. Rarity

Pet template rarity:
- `common` — starter, easy drops.
- `uncommon` — common quest reward.
- `rare` — dungeon drop.
- `epic` — world boss / event.
- `legendary` — server-first, unique.
- `mythic` — endgame prestige.

Rarity → stat multiplier + trait pool weight.

## 6. Acquisition Source

- `drop` — random mob drop (baby pet as loot).
- `quest` — quest reward.
- `event` — seasonal event.
- `starter` — race/class starter.
- `purchase` — Phase 02 monetization (cosmetic only).
- `evolution` — evolved from previous stage.
- `breeding` (Phase 02+).

## 7. Endpoint

**GET `/v1/content/pet-species`**
```json
{
  "species": [
    {
      "code": "wolf_cub",
      "family": "beast",
      "type": "combat",
      "rarity": "uncommon",
      "max_level": 50,
      "evolution_next": "alpha_wolf",
      "acquisition_source": "drop"
    }
  ]
}
```

## 8. Deprecation

Species deprecated:
- Mevcut instance korunur.
- Yeni acquisition blocked.
- Notification "This species is no longer available" (nadir).

## 9. Metrics
- `pet_species_active_total` gauge.
- `pet_species_distribution{code}` gauge.

## 10. Cross-Reference
- [Instance](instance.md), [Traits](traits.md), [Skills](skills.md), [Evolution](evolution.md)
- [database/content-templates.md §21](../database/content-templates.md)
- [game-design/pets.md](../game-design/pets.md)

## 11. Bekleyen Kararlar
- Species count v1 (~20-30 species önerisi).
- Breeding Phase 02+ species combine.
- Rare species event scheduling Phase 02.
