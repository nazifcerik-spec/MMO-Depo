# Dungeon (Entity + Template + Instance)

## 1. Amaç
Dungeon'un kalıcı tanımı (`content.dungeon_template`) + runtime instance (`activity.instance`).

## 2. Template Model

Bkz. [database/content-templates.md §14](../database/content-templates.md).

Phase 19 eklemeleri:
```sql
ALTER TABLE content.dungeon_template
  ADD COLUMN floor_count SMALLINT NOT NULL DEFAULT 3,
  ADD COLUMN estimated_clear_time_ms JSONB NOT NULL DEFAULT '{}',  -- per difficulty
  ADD COLUMN layout_json JSONB NOT NULL DEFAULT '{}',
  ADD COLUMN entry_flavor_key VARCHAR(96) NULL,
  ADD COLUMN completion_lockout_kind VARCHAR(16) NOT NULL DEFAULT 'weekly',  -- weekly/daily/none
  ADD COLUMN daily_lockout_max SMALLINT NOT NULL DEFAULT 1,
  ADD COLUMN afk_yield_multiplier NUMERIC(4,2) NOT NULL DEFAULT 0.5,
  ADD COLUMN deprecated_at TIMESTAMPTZ NULL,
  ADD COLUMN tags VARCHAR(24)[] NOT NULL DEFAULT '{}';
```

## 3. Content Pack

```yaml
# content/dungeons/ashen_keep.yaml
code: ashen_keep
display_name_key: dungeon.ashen_keep.name
lore_key: dungeon.ashen_keep.lore
region_code: ashen_hills

min_level: 20
max_level: 30
party_size_min: 1
party_size_max: 5
weekly_lockout: true
daily_lockout_max: 1
idle_run_supported: true
afk_yield_multiplier: 0.5

floor_count: 3
estimated_clear_time_ms:
  normal: 1200000    # 20 min
  heroic: 2100000    # 35 min

layout:
  floors:
    - floor: 1
      rooms:
        - kind: enemy_pack
          mob_pool: [skeleton_warrior, skeleton_archer]
          room_count: 6
        - kind: elite
          elite_code: skeleton_captain
        - kind: treasure
          treasure_code: ashen_keep_chest_small
        - kind: boss
          boss_code: keep_guardian
    - floor: 2
      rooms:
        - kind: enemy_pack
          mob_pool: [shadow_wraith, shadow_witch]
          room_count: 8
        - kind: trap
          trap_code: shadow_darts
        - kind: enemy_pack
          mob_pool: [shadow_wraith]
          room_count: 4
        - kind: treasure
          treasure_code: ashen_keep_chest_mid
        - kind: boss
          boss_code: fire_wraith
    - floor: 3
      rooms:
        - kind: enemy_pack
          mob_pool: [ashen_lich_minion]
          room_count: 5
        - kind: elite
          elite_code: lich_lieutenant
        - kind: random_event
          event_pool: [rare_summon, treasure_boost]
        - kind: boss
          boss_code: ashen_lich_king

reward_pool_code: ashen_keep_reward

entry_flavor_key: dungeon.ashen_keep.entry_flavor
tags: [undead, tomb, story]
```

## 4. Instance Model

`activity.instance` (bkz. [database/activity.md §3](../database/activity.md)):
```
id, template_id, template_kind='dungeon',
party_id (opsiyonel),
owner_character_id (solo veya leader),
difficulty,
seed,
current_encounter_index,
state,
created_at, expires_at, completed_at
```

Ek Phase 19:
```sql
ALTER TABLE activity.instance
  ADD COLUMN current_floor SMALLINT NOT NULL DEFAULT 1,
  ADD COLUMN current_room_index SMALLINT NOT NULL DEFAULT 0,
  ADD COLUMN mode VARCHAR(16) NOT NULL DEFAULT 'manual';  -- manual/idle
```

## 5. Lifecycle

```
1. Create instance (POST /v1/dungeons/enter).
2. Room progression:
   - Manuel: player triggers room combat/interact.
   - Idle: analytical clear per elapsed_ms.
3. Boss encounter per floor end.
4. Complete on last boss defeat.
5. Weekly lockout record (bkz. loot.md, lockout.md analog).
6. Reward distribution.
```

Manuel expire 2h; idle mode → contract lifetime.

## 6. Difficulty Modes

- `normal` — baseline.
- `heroic` — HP × 1.8, damage × 1.3, loot better; weekly lockout separate.
- `mythic` (Phase 02) — affix system; +variety modifier.

Difficulty content-driven scaling.

## 7. Endpoint

**POST `/v1/dungeons/enter`**
Auth + CharacterScope.
```json
{
  "dungeon_template_code": "ashen_keep",
  "difficulty": "normal",
  "mode": "manual",
  "party_id": null
}
```

Server:
1. Validate:
   - Character level in range.
   - Party size ok.
   - Weekly lockout not exceeded.
   - Character at dungeon area (Phase 02 zorunlu).
2. Mode dispatch:
   - **manual**: activity.instance create.
   - **idle**: activity.contract create (channel=primary).
3. Emit event `dungeon.entered`.
4. Response 201.

**GET `/v1/dungeons/instances/{id}`** — progress.
**POST `/v1/dungeons/instances/{id}/room/{room_id}/enter`** — manuel room.
**POST `/v1/dungeons/instances/{id}/leave`** — leave.

## 8. AFK Compute (Idle Mode)

Bkz. [afk-run.md](afk-run.md).

Contract yield analytical:
```
Total clear time = dungeon.estimated_clear_time_ms[difficulty]
clears_per_hour = 3600000 / clear_time × afk_yield_multiplier
elapsed_hours = contract.elapsed_ms / 3600000
clears_completed = min(int(clears_per_hour × elapsed_hours), weekly_lockout_remaining)
```

## 9. Weekly Lockout

Bkz. [database/progression-quests-achievements.md §8](../database/progression-quests-achievements.md).

- Weekly reset Salı 06:00 UTC.
- Character × dungeon × kind (weekly/daily) counter.
- Lockout dolarsa loot rare penalty (XP normal); ancak clear counted.

## 10. Metrics
- `dungeon_instances_created_total{dungeon,mode}`.
- `dungeon_completions_total{dungeon,difficulty}`.
- `dungeon_weekly_lockout_hit_total`.

## 11. Cross-Reference
- [Floor](floor.md), [Room](room.md), [Boss](boss.md), [Loot](loot.md), [AFK Run](afk-run.md)
- [activity-engine/dungeon.md](../activity-engine/dungeon.md)
- [game-design/dungeons.md](../game-design/dungeons.md)
- [database/content-templates.md §14](../database/content-templates.md)

## 12. Bekleyen Kararlar
- Mega dungeon Phase 02+.
- Mythic+ affix system Phase 02.
- Dungeon reset item Phase 02 monetization.
- Positional room layout Phase 02.
