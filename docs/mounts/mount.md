# Mount (Entity + Species + Instance)

## 1. Amaç
Mount'un kalıcı tanımı — species template + character owned instance.

## 2. Template Model

Kaynak: [database/content-templates.md §22](../database/content-templates.md).

Phase 21 eklemeleri:
```sql
ALTER TABLE content.mount_template
  ADD COLUMN species_family VARCHAR(24) NOT NULL DEFAULT 'horse',
  ADD COLUMN max_level SMALLINT NOT NULL DEFAULT 30,
  ADD COLUMN xp_curve_json JSONB NOT NULL DEFAULT '{}',
  ADD COLUMN base_stats_json JSONB NOT NULL DEFAULT '{}',
  ADD COLUMN skill_codes VARCHAR(48)[] NOT NULL DEFAULT '{}',
  ADD COLUMN terrain_compatibility_json JSONB NOT NULL DEFAULT '{}',
  ADD COLUMN acquisition_source VARCHAR(32) NOT NULL DEFAULT 'quest',
  ADD COLUMN deprecated_at TIMESTAMPTZ NULL;
```

## 3. Content Pack

```yaml
# content/mount_species/steel_warhorse.yaml
code: steel_warhorse
species_family: horse
type: ground
rarity: uncommon
display_name_key: mount.steel_warhorse.name

speed_multiplier: 2.0
cargo_bonus_slots: 5
max_level: 30

xp_curve:
  base: 100
  factor: 1.10

base_stats:
  hp: 500
  stamina_max: 100

skill_codes: [mount_gallop, mount_endurance]

terrain_compatibility:
  grass: 1.0
  grass_road: 0.9         # penalty lower (fast)
  forest: 1.2
  mountain: 1.6           # penalty
  water: 5.0              # very slow
  cave: 1.4

unlock_condition:
  quest_completed: horse_taming
  min_level: 10

acquisition_source: quest
art_ref: mounts/steel_warhorse.png
tags: [ground, universal, starter]
```

## 4. Species Family

| Family | Örnek | Type | Terrain Bias |
|---|---|---|---|
| `horse` | steel warhorse, wild stallion | ground | grass/road excellent |
| `wolf` | dire wolf | ground | forest/hills good |
| `bear` | war bear | ground/cargo | mountain/tundra good |
| `dolphin` | sea dolphin | aquatic | water only |
| `turtle` | sea turtle | aquatic/cargo | water + coast |
| `griffon` (Phase 02) | griffon | winged | all terrain |
| `dragon` (Phase 03) | wyvern | winged/combat | all + air |
| `mystical` | unicorn, phoenix | ground/cosmetic | unique |

## 5. Instance Model

Bkz. [database/companions.md §3](../database/companions.md).

Phase 21 eklemeleri:
```sql
ALTER TABLE companion.mount_owned
  ADD COLUMN level SMALLINT NOT NULL DEFAULT 1,
  ADD COLUMN xp BIGINT NOT NULL DEFAULT 0,
  ADD COLUMN stamina_current SMALLINT NOT NULL DEFAULT 100,
  ADD COLUMN stamina_max SMALLINT NOT NULL DEFAULT 100,
  ADD COLUMN stamina_last_recovery_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  ADD COLUMN skill_ranks_json JSONB NOT NULL DEFAULT '{}',
  ADD COLUMN last_active_at TIMESTAMPTZ NULL;
```

## 6. Acquisition

Mount unlock kaynakları:
- **quest** — story/side quest reward.
- **drop** — nadir rare mob drop.
- **event** — seasonal.
- **purchase** — Phase 02 cosmetic monetization.
- **achievement** — collection reward.
- **crafted** — Phase 02 special (rare recipe).

## 7. Summon / Dismount

Aktif mount summon:
- Karakter safe zone dışında summon.
- Combat combat başlarsa auto-dismount + cooldown 30s.
- Manuel dismount → cooldown 5s.

State (Redis):
- `char:mount:<char_id>` TTL infinite while active — hangi mount summoned.
- Combat state → dismount trigger.

## 8. Endpoint

**GET `/v1/characters/{id}/mounts`**
```json
{
  "mounts": [
    {
      "instance_id": "...",
      "species_code": "steel_warhorse",
      "custom_name": "Thunder",
      "level": 12,
      "xp": 4200,
      "stamina": 85,
      "is_active": true,
      "speed_multiplier": 2.0,
      "cargo_bonus_slots": 5
    }
  ]
}
```

**POST `/v1/characters/{id}/mounts/{instance_id}/summon`**
Server:
1. Validate safe zone bypass (out of combat).
2. Set active mount state.
3. Apply speed + cargo modifier.
4. Emit event.

## 9. BoP Discipline

DD-024: mount bound to character. Trade edilemez, mailable değil.

## 10. Metrics
- `mount_owned_total{species}` gauge.
- `mount_active_total{type}` gauge.
- `mount_summoned_total`.

## 11. Cross-Reference
- [Level](level.md), [Stats](stats.md), [Skills](skills.md), [Travel](travel.md), [Terrain](terrain.md), [Carry Capacity](carry-capacity.md)
- [database/companions.md §3-4](../database/companions.md)
- [database/content-templates.md §22](../database/content-templates.md)
- [travel/mount.md](../travel/mount.md)

## 12. Bekleyen Kararlar
- Species count v1 (~10-15 önerisi).
- Mount rename limit (v1: 1 per 30 gün).
- Mount trade Phase 03 (BoA opsiyonu).
- Winged mount schedule Phase 02.
