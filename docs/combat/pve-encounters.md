# PvE Encounters (Mob + Boss AI)

## 1. Amaç
NPC saldırgan davranış modelleri. Mob standart + boss encounter phase + world boss + AI profile.

## 2. NPC Model

Kaynak: [database/content-templates.md §17](../database/content-templates.md).

```sql
content.npc_template
  id, code, family, level
  hp_base, damage_base, armor, magic_resist
  abilities_json      -- ability_ids listesi
  ai_profile_id      -- AI davranış türü
  loot_table_id
  xp_reward
  elite_flag
```

## 3. AI Profile

`content.ai_profile` (yeni Phase 17):
```sql
CREATE TABLE content.ai_profile (
  id UUID PRIMARY KEY,
  code VARCHAR(32) UNIQUE NOT NULL,
  behavior_kind VARCHAR(24) NOT NULL,  -- melee_bruiser/caster/healer/skirmisher/tank_taunt
  decision_tree_json JSONB NOT NULL,
  version INT NOT NULL DEFAULT 1
);
```

Behavior kinds (v1):
- `melee_bruiser` — kısa range attack; charge target.
- `caster` — cast spells; kite if melee approaches.
- `healer_support` — heal allies, buff.
- `skirmisher` — hit-and-run; low HP retreats.
- `tank_taunt` — draws aggro; taunts.
- `boss_phase` — multi-phase encounter (bkz. §Boss).

## 4. AI Decision Tree

```yaml
# content/ai_profiles/melee_bruiser.yaml
code: melee_bruiser
behavior_kind: melee_bruiser
decision_tree:
  root:
    kind: select
    children:
      - kind: condition
        cond: target_in_melee_range
        then:
          kind: cast
          skill_code: bash
      - kind: condition
        cond: hp_below_pct
        params: { threshold: 30 }
        then:
          kind: cast
          skill_code: retreat_and_heal
      - kind: default
        action: move_toward_target   # v1 positional yok; abstract "chase"
```

v1 simple; Phase 02 rich behavior tree.

## 5. Ability Set

NPC ability_ids listesi (`content.skill_template` reference):
- Basic attack (auto).
- Special ability (cooldown).
- Boss unique mechanic.

Ability cast tick loop:
```
if ability.CooldownReady && ai.Decide(cast):
    skillEngine.Cast(mob, ability, target, seed)
```

## 6. Elite Mobs

`elite_flag=true`:
- HP × 3 average.
- Damage × 1.5.
- Loot table richer.
- Special abilities (more).
- Group hunt recommended.

## 7. Boss Encounters

Bkz. [database/content-templates.md §15](../database/content-templates.md):
```sql
content.boss_encounter
  id, dungeon_id, order_index
  phase_definitions JSONB
  loot_table_id
  ai_profile_id
```

Phase definition örneği:
```yaml
# content/boss_encounters/ashen_lich_king.yaml
phases:
  - phase: 1
    trigger: encounter_start
    hp_pct_range: [100, 66]
    abilities: [dark_bolt, summon_skeletons]
    ability_rotation:
      - skill: dark_bolt
        cd: 4000
      - skill: summon_skeletons
        cd: 20000
  - phase: 2
    trigger: hp_below
    hp_pct: 66
    on_trigger: cast_enrage
    abilities: [dark_bolt, shadow_nova, summon_skeletons]
    ability_rotation:
      - skill: shadow_nova
        cd: 8000
      - skill: dark_bolt
        cd: 3000
  - phase: 3
    trigger: hp_below
    hp_pct: 33
    on_trigger: enrage_final
    abilities: [dark_bolt, shadow_nova, meteor]
    ability_rotation:
      - skill: meteor
        cd: 10000
      - skill: shadow_nova
        cd: 6000
    enrage_timer_ms: 60000   # 60 sn sonra kill party (safety)
```

## 8. Enrage Timer

Boss enrage (safety):
- Belirli süre içinde boss defeat edilmezse enrage.
- Enrage: damage × 3; unstoppable; wipe.
- DPS check: party yeterli DPS'e sahip olmalı.

## 9. World Boss

Bkz. [game-design/world-bosses.md](../game-design/world-bosses.md):
- Spawn window; content template.
- HP scale: participants_count × baseline.
- Contribution tracked (`progression.world_boss_participation`).
- Manuel katılım gerekir (AFK yok).
- Loot: contribution threshold; personal roll.

Phase mechanics:
- Same pattern as dungeon boss + participant scaling.

## 10. Adds Mechanic

Boss summon adds:
- `summon_skeletons` ability → spawn 3-5 add mob.
- Adds threat kendi tablosu.
- Encounter until adds killed or boss killed (config).

## 11. Fixation / Threat Override

Boss-specific threat behaviors:
- **Fixation**: hedef sabit karakter (highest DPS or random).
- **Reset threat**: her N sec threat table clear.
- **Kite mechanic**: boss chases specific target.

Phase-based; content-driven.

## 12. Mob Family

`content.npc_template.family`:
- `beast`, `humanoid`, `undead`, `elemental`, `demon`, `dragon`, `insect`.

Bestiary mastery per family (bkz. [character/mastery.md](../character/mastery.md), [skill-engine/mastery.md](../skill-engine/mastery.md)).

Damage bonuses:
- `damage_pct.family.beast` — beast damage boost.
- Mastery milestones ile aktive.

## 13. Loot Roll

Boss/mob kill → loot table roll (bkz. [systems/loot-system.md](../systems/loot-system.md)).

- Contribution-based (world boss).
- Personal loot (dungeon default DD-012).
- Party loot rule (party level).

## 14. Boss First Kill

Server-first tracking:
```
if progression.world_boss_spawn or dungeon_first_clear:
    if !first_server_flag_set:
        mark first_server_flag=true
        grant unique title/cosmetic
        server-wide announcement
```

## 15. AFK Encounter (Idle Dungeon)

Idle-run dungeon:
- Per-encounter analytical clear:
  - Total HP / party DPS = clear time per boss.
  - Loot per boss table analytical roll.
- No live tick.

Detay: [activity-engine/dungeon.md](../activity-engine/dungeon.md).

## 16. Metrics
- `npc_killed_total{npc_family,elite}`.
- `boss_encounter_completed_total{dungeon,boss}`.
- `boss_enrage_wipe_total`.
- `world_boss_kill_total{boss}`.

## 17. Event

- `combat.mob_killed`.
- `combat.boss_defeated`.
- `combat.boss_enraged`.
- `combat.encounter_phase_transition`.
- `combat.adds_spawned`.
- `combat.world_boss_defeated`.

## 18. Cross-Reference
- [Combat Loop](combat-loop.md), [Damage](damage.md), [Threat](threat.md), [Death & Revive](death-revive.md), [Combat Events](combat-events.md)
- [database/content-templates.md §17-18](../database/content-templates.md)
- [game-design/dungeons.md](../game-design/dungeons.md), [game-design/world-bosses.md](../game-design/world-bosses.md)
- [activity-engine/dungeon.md](../activity-engine/dungeon.md), [activity-engine/hunting.md](../activity-engine/hunting.md)

## 19. Bekleyen Kararlar
- Behavior tree DSL Phase 02 detay.
- Positional mechanics (spread, stack, avoid AoE) Phase 02.
- Boss uniqueness (mekanik variance) Phase 02.
- Raid content (10-20 kişi) Phase 02+.
- AI difficulty scaling (mythic affix) Phase 02+.
