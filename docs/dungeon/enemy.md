# Enemy (Dungeon Standard Mob)

## 1. Amaç
Dungeon içindeki standart mob'lar. `content.npc_template` normal tier (elite_flag=false).

## 2. Mob Pack Composition

Room `enemy_pack`:
```yaml
kind: enemy_pack
mob_pool: [skeleton_warrior, skeleton_archer]
mob_count: 4
```

Server spawn:
- Pool weight-based selection per mob slot.
- Aynı pool'dan farklı mob mix.

## 3. Mob Stat Scaling (Dungeon)

Dungeon mob stat baz `content.npc_template` + dungeon multiplier:
```
dungeon_hp_mult = 1.0 (normal) + 0.5 (heroic) + 1.5 (mythic)
dungeon_dmg_mult = 1.0 (normal) + 0.3 (heroic) + 1.0 (mythic)
floor_deep_mult = 1 + (floor_index - 1) × 0.20    // deeper floor +%20
```

Effective:
```
mob.HP = base.hp × dungeon_hp_mult × floor_deep_mult
mob.Damage = base.damage × dungeon_dmg_mult × floor_deep_mult
```

## 4. Mob AI

`content.ai_profile` reference (bkz. [combat/pve-encounters.md §4](../combat/pve-encounters.md)):
- `melee_bruiser`, `caster`, `healer_support`, `skirmisher`.

Dungeon-özel AI variant (Phase 02):
- `phalanx` — 3+ melee'nin defensive formation.
- `flanker` — sneak behind highest DPS.

## 5. Combat Behavior

Standart combat loop (bkz. [combat/combat-loop.md](../combat/combat-loop.md)):
- Auto-attack per attack_speed.
- Ability rotation per AI decision.
- Threat table (per mob).

## 6. Loot

Kill per mob:
- `npc_template.loot_table_id` roll.
- Standart material + occasional rare drop.
- Party: personal loot each.

## 7. XP Reward

`npc_template.xp_reward` per kill.

Dungeon multiplier:
- Normal: ×1.0.
- Heroic: ×1.5.

## 8. AFK Yield

Reveal aggregate per room:
```
mob_count_in_room = room.mob_count
kills = mob_count_in_room   (all mobs cleared analytical)
loot = mob.loot_table × kills (aggregated via loot system)
xp = mob.xp_reward × kills × dungeon_xp_mult
```

Per-mob simulation yok; bounded aggregate.

## 9. Cross-Reference
- [Elite](elite.md), [Boss](boss.md), [Enemy](enemy.md)
- [combat/pve-encounters.md](../combat/pve-encounters.md)
- [database/content-templates.md §17](../database/content-templates.md)

## 10. Bekleyen Kararlar
- Dungeon-özel AI variant Phase 02.
- Elite spawn chance in normal room (%5) Phase 02.
- Mob level scaling (dynamic vs static) Phase 02.
