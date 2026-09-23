# Boss

## 1. Amaç
Floor sonu boss encounter — phase-based mechanic, enrage timer, unique loot.

Foundation: [combat/pve-encounters.md §7](../combat/pve-encounters.md).

## 2. Boss Model

`content.boss_encounter` (bkz. [database/content-templates.md §15](../database/content-templates.md)):
```
id, dungeon_id, order_index
phase_definitions JSONB
loot_table_id
ai_profile_id
```

Ek Phase 19:
```sql
ALTER TABLE content.boss_encounter
  ADD COLUMN estimated_hp_by_difficulty JSONB NOT NULL DEFAULT '{}',
  ADD COLUMN enrage_timer_ms INT NOT NULL DEFAULT 300000,   -- 5 min default
  ADD COLUMN first_clear_reward_json JSONB NULL,
  ADD COLUMN unique_drop_item_codes VARCHAR(48)[] NOT NULL DEFAULT '{}',
  ADD COLUMN set_item_drop_pool_id UUID NULL;
```

## 3. Phase System

Boss encounter phase-based:
```yaml
# content/boss_encounters/ashen_lich_king.yaml
code: ashen_lich_king
dungeon_code: ashen_keep
order_index: 3       # last boss

phases:
  - phase: 1
    trigger: encounter_start
    hp_pct_range: [100, 66]
    abilities: [dark_bolt, summon_skeletons]
    ability_rotation:
      - skill_code: dark_bolt
        cooldown_ms: 4000
      - skill_code: summon_skeletons
        cooldown_ms: 20000
  - phase: 2
    trigger: hp_below_pct
    hp_pct: 66
    on_trigger:
      - kind: cast
        skill_code: enrage_cast
    abilities: [dark_bolt, shadow_nova, summon_skeletons]
    ability_rotation:
      - skill_code: shadow_nova
        cooldown_ms: 8000
      - skill_code: dark_bolt
        cooldown_ms: 3000
  - phase: 3
    trigger: hp_below_pct
    hp_pct: 33
    on_trigger:
      - kind: cast
        skill_code: final_desperation
    abilities: [dark_bolt, shadow_nova, meteor]
    ability_rotation:
      - skill_code: meteor
        cooldown_ms: 10000
      - skill_code: shadow_nova
        cooldown_ms: 6000

estimated_hp_by_difficulty:
  normal: 25000
  heroic: 45000

enrage_timer_ms: 300000    # 5 min

loot_table_code: ashen_lich_boss_loot
unique_drop_item_codes: [lich_crown, phylactery_shard]
set_item_drop_pool_code: ashen_keep_set
first_clear_reward:
  title_code: lich_slayer
  cosmetic_code: shadow_crown
```

## 4. Phase Transition

Trigger types:
- `encounter_start` — phase 1 opening.
- `hp_below_pct` — HP threshold.
- `time_elapsed` — after N sec.
- `add_defeated` — boss adds killed.
- `player_action` — special (Phase 02).

Transition:
- `on_trigger` list — casts + spawns.
- New ability set replace old.
- New rotation.

## 5. Enrage Timer

Boss enrage:
- Belirli süre içinde kill edilmezse enrage.
- Enrage: damage × 3; unstoppable; wipe.
- Party DPS check.

DPS threshold formula:
```
required_dps = boss.HP / (enrage_timer_ms / 1000)
```

Örnek: 25000 HP / 300 sec = 83 DPS required.

## 6. Adds Mechanic

Boss summon adds:
- `summon_skeletons` cast → spawn N mob.
- Adds threat separate.
- Adds killed → count toward achievement.
- Boss immune while adds alive (opsiyonel — kill order strategy).

## 7. Fixation / Threat Override

Boss-özel threat behavior:
- **Fixation**: hedef sabit (random party member N sec).
- **Reset threat**: N sec threat table clear.
- **Fear/CC** aoe: multi target CC.

Content-driven per phase.

## 8. AFK Boss Compute

Reveal (idle mode):
```
boss_hp = boss.estimated_hp[difficulty]
party_dps = estimatePartyDPS(snapshot)
if party_dps >= required_dps:
    clear_time = boss_hp / party_dps
    loot = boss.loot_table roll
    xp = boss.xp_reward
else:
    wipe (enrage)
    no loot for this run
    partial progress kept
```

Analytical; per-tick simulate yok.

## 9. First Clear Reward

Server-first kill:
- Title reward.
- Cosmetic drop.
- Server chat anons.
- Achievement.

Race condition SERIALIZABLE tx.

## 10. Boss Loot Priority

Loot roll:
1. Guaranteed items (`unique_drop_item_codes` weekly first kill).
2. Set item pool (`set_item_drop_pool_id`).
3. Standard loot table (weighted).

Party: personal loot (each character rolls).

## 11. Weekly Lockout

Boss defeated → dungeon lockout incremented.
Subsequent kills within week: rare loot penalty (bkz. [loot.md](loot.md)).

## 12. Metrics
- `dungeon_boss_defeated_total{dungeon,boss,difficulty}`.
- `dungeon_boss_enrage_wipe_total{boss}`.
- `dungeon_boss_first_clear_total{dungeon}`.

## 13. Audit
- `dungeon.boss_defeated` (severity=info).
- `dungeon.boss_first_server_clear` (severity=notable).
- `dungeon.boss_enraged` (severity=warn).

## 14. Cross-Reference
- [Dungeon](dungeon.md), [Floor](floor.md), [Room](room.md), [Loot](loot.md), [AFK Run](afk-run.md)
- [combat/pve-encounters.md](../combat/pve-encounters.md)
- [game-design/dungeons.md](../game-design/dungeons.md)
- [database/content-templates.md §15](../database/content-templates.md)

## 15. Bekleyen Kararlar
- Positional boss mechanics Phase 02 (avoid stack point).
- Mythic affix random boss modifier Phase 02.
- Raid boss content (10+ oyuncu) Phase 02+.
- Boss unique drop dedup (BoP prevent duplicate) v1 policy.
