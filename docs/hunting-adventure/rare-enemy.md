# Rare Enemy

## 1. Amaç
Nadir mob spawn — ~%1-0.5 baz hunt sırasında; ekstra HP + loot + kill credit hakkı. Server-wide anons opsiyonu (Phase 02).

## 2. Spawn Model

Rare mob template:
```sql
CREATE TABLE content.rare_mob_spawn (
  id UUID PRIMARY KEY,
  code VARCHAR(48) UNIQUE NOT NULL,
  npc_template_id UUID NOT NULL,           -- content.npc_template reference
  monster_area_spawn_id UUID NULL,           -- hangi spot'ta
  region_id UUID NULL,                        -- alternate: region-wide
  spawn_chance_per_kill NUMERIC(6,4) NOT NULL,   -- 0.0001..0.01 (0.01%..1%)
  spawn_windows_json JSONB NULL,               -- cron-based schedule optional
  server_wide_announce BOOLEAN NOT NULL DEFAULT false,
  rare_loot_table_id UUID NOT NULL,
  first_kill_reward_json JSONB NULL,          -- server-first bonus
  version INT NOT NULL DEFAULT 1
);
```

## 3. Content Pack

```yaml
# content/rare_mob_spawns/shadowfang_alpha.yaml
code: shadowfang_alpha
npc_template_code: shadowfang_alpha
monster_area_spawn_code: whispering_forest_wolves

spawn_chance_per_kill: 0.005        # her kill'de %0.5 spawn
spawn_windows: null                  # random spawn (no schedule)
server_wide_announce: false          # v1 kapalı; Phase 02 opt

rare_loot_table_code: shadowfang_alpha_loot
first_kill_reward:
  title_code: shadowfang_slayer
  bonus_gold_copper: 5000
```

## 4. Spawn Trigger

Reveal içinde (hunting aktivite):
```go
kill_count = actual_kills_from_hunting
for _, rareTemplate := range rareSpawnsFor(spot) {
    expected_rare_spawns = kill_count × rareTemplate.SpawnChancePerKill
    rare_count = samplePoisson(expected_rare_spawns, 1.0, seed.Sub("rare-"+rareTemplate.Code))
    for i := 0; i < rare_count; i++ {
        // Character auto-defeats rare mob (if DPS >= mob HP × 3 factor)
        if canDefeatRareMob(snapshot, rareTemplate.NPCTemplate):
            grantRareLoot(rareTemplate, seed.Sub(fmt.Sprintf("rare-loot-%d", i)))
            if firstServerKill(rareTemplate):
                grantFirstKillReward(rareTemplate)
                if rareTemplate.ServerWideAnnounce:
                    broadcastAnnouncement(...)
```

## 5. Defeat Condition

AFK rare mob defeat:
- Character DPS `dps >= rareMob.HP / max_engagement_time` (60 sec default).
- Aksi: mob despawn (escape).

Manuel oyun: normal combat.

## 6. Rare Loot Table

Nadir drop (Phase 02 tercih):
- Unique cosmetic (title, dye).
- Rare crafting material.
- Legendary item chance +%0.5.
- Rare pet unlock (Phase 02).

Rare loot tables `content.loot_table` extension.

## 7. Server-First Kill

`first_kill_reward`:
- Server-wide flag (bir kez claim).
- Title (permanent).
- Bonus gold + special drop.
- Anons "First to defeat Shadowfang Alpha: <character_name>".

Race condition: SERIALIZABLE tx (bkz. [database/progression-quests-achievements.md §6](../database/progression-quests-achievements.md)).

## 8. Personal vs Global Rare

- **Personal rare** (v1 default): karakter kendi rare mob'unu görür + öldürür.
- **Global rare** (Phase 02): sunucu-wide; ilk vuran alır.

v1 personal — idle MMO fairness; Phase 02 global rare event opt-in.

## 9. Cooldown

Character per rare mob:
- Kill credit sonrası 7 gün cooldown (rare drop spam engelleme).
- `character.rare_mob_kill_history` table (Phase 17 önerisi).

Alternate: global cooldown template-wide (30 days) — content karar.

## 10. Reveal Payload

```json
{
  "notable_encounters": [
    {
      "kind": "rare_mob",
      "mob_code": "shadowfang_alpha",
      "defeated": true,
      "rare_drops": [
        { "item_code": "shadowfang_tooth", "quantity": 1 },
        { "item_code": "shadow_essence", "quantity": 3 }
      ],
      "first_server_kill": false
    }
  ]
}
```

Client UI highlights.

## 11. Notification

- Personal: reveal panel highlight.
- Server-wide (Phase 02): global chat + WS push.

## 12. Metrics
- `rare_mob_spawned_total{code}`.
- `rare_mob_defeated_total{code}`.
- `rare_mob_first_kill_total{code}`.
- `rare_drop_awarded_total{item}`.

## 13. Audit
- `rare_mob.defeated` (severity=info).
- `rare_mob.first_server_kill` (severity=notable).

## 14. Cross-Reference
- [Hunting](hunting.md), [Encounter](encounter.md), [Elite Enemy](elite-enemy.md)
- [activity-engine/hunting.md](../activity-engine/hunting.md)
- [combat/pve-encounters.md](../combat/pve-encounters.md)
- [systems/loot-system.md](../systems/loot-system.md)

## 15. Bekleyen Kararlar
- Server-wide anons Phase 02 aktivasyon eşiği.
- Global rare event framework Phase 02.
- Rare drop rate tune final.
- Cooldown pattern (personal vs global) — v1 personal 7d.
