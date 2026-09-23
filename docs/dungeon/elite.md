# Elite (Dungeon Elite Mob)

## 1. Amaç
Dungeon içindeki elite mob'lar. Normal'den güçlü, boss'tan hafif. Loot bonus.

Foundation: [hunting-adventure/elite-enemy.md](../hunting-adventure/elite-enemy.md).

## 2. Elite Kind (Dungeon)

- **Solo elite**: tek elite mob (harder combat).
- **Elite group**: 2-3 elite co-spawn.
- **Elite pack** (Phase 02): 1 elite + 3-4 mob.

Room `kind: elite`:
```yaml
kind: elite
elite_code: skeleton_captain
```

Group:
```yaml
kind: elite
elite_group_code: bandit_trio
```

## 3. Elite Stat

Dungeon elite → base elite mult × dungeon mult:
```
elite_HP = base.hp × 3 × dungeon_hp_mult × floor_deep_mult
elite_damage = base.damage × 1.5 × dungeon_dmg_mult × floor_deep_mult
```

## 4. Special Abilities

Elite mob 2-3 unique ability:
- AoE cleave.
- CC (stun, root).
- Enrage (low HP).
- Summon adds (Phase 02).

Content per elite template.

## 5. Elite Loot

Elite defeated:
- Standard loot table.
- +Rare drop chance (baseline +%10 vs normal mob).
- Occasional set item piece.

Rare drop kaynağı — set collection için önemli.

## 6. Champion vs Mini-Boss (Dungeon)

- **Elite** (tier 1): standart elite; per room.
- **Champion** (tier 2): rare elite; ~%5 pool weight; guaranteed rare drop.
- **Mini-boss** (tier 3): dungeon mid-floor optional; phase mechanics-lite.

Content-driven per dungeon.

## 7. Combat

- Longer TTK vs normal mob (~2-4x).
- Party recommended.
- Solo doable but slow.

## 8. AFK Yield

Reveal:
```
elite_kill_count = 1 (per elite room)
elite_time_ms = base.hp × 3 / party_dps × dungeon_mult
elite_loot = elite.loot_table + rare_drop roll (guaranteed for champion)
elite_xp = base.xp × 3 × dungeon_xp_mult
```

## 9. Cross-Reference
- [Enemy](enemy.md), [Boss](boss.md), [Room](room.md)
- [hunting-adventure/elite-enemy.md](../hunting-adventure/elite-enemy.md)
- [combat/pve-encounters.md](../combat/pve-encounters.md)

## 10. Bekleyen Kararlar
- Champion vs elite balance (v1 birleşik; Phase 02 net ayrım).
- Elite ability variety per content pack (v1: ~2 per elite; Phase 02 daha zengin).
- Elite pack (elite + minions) Phase 02.
