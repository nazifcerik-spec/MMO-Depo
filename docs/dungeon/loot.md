# Loot (Dungeon Distribution)

## 1. Amaç
Dungeon reward payload — mob kills, elite drops, treasure, boss loot, event rewards. Party rule + weekly lockout.

Foundation: [systems/loot-system.md](../systems/loot-system.md), [activity-engine/dungeon.md](../activity-engine/dungeon.md).

## 2. Loot Kaynağı Katmanları

Per dungeon clear:
1. **Mob loot** (per enemy_pack room).
2. **Elite loot** (per elite room; rare bonus).
3. **Treasure** (per chest room).
4. **Boss loot** (per boss defeated).
5. **Random event loot** (opsiyonel bonus).
6. **First clear bonus** (server-first title + cosmetic).

## 3. Party Loot Rule

**Personal loot** (v1 default, DD-012):
- Each character rolls independently.
- No competition; no share.
- Applies to mob/elite/boss/treasure.

**Round-robin** (Phase 02 opt).
**Need Before Greed** (Phase 02 opt).
**Master Loot** (raid; Phase 02+).

## 4. Boss Loot Composition

Boss loot per character (personal):
```
1. Guaranteed items (weekly first kill only):
   - unique_drop_item_codes: dedup enforce
2. Set item pool:
   - set_item_drop_pool_id: %30 chance one set piece
3. Standard boss loot table:
   - weighted roll (rare +50% vs mob)
4. Currency (dungeon_token): guaranteed 20-50 per boss
```

## 5. Weekly Lockout Effect

Character × dungeon × difficulty (weekly counter):
- **First clear this week**: full loot + bonus.
- **Second clear+**: rare loot skip; only currency + basic material; XP normal.
- **Daily lockout**: subset (bazı dungeon daily reset).

Detay: [database/progression-quests-achievements.md §8](../database/progression-quests-achievements.md).

Config:
```yaml
lockout_effect:
  first_clear:
    rare_drop_multiplier: 1.0
    set_item_chance: 1.0
    unique_drop_enabled: true
  subsequent_clears:
    rare_drop_multiplier: 0.3
    set_item_chance: 0.0
    unique_drop_enabled: false
    currency_multiplier: 0.5
```

## 6. First Server Clear

Dungeon full clear + server-first:
- Title reward.
- Cosmetic drop (crown, cape).
- Server-wide anons.
- Achievement.

Race condition SERIALIZABLE tx (bkz. [database/progression-quests-achievements.md §6](../database/progression-quests-achievements.md)).

## 7. Dungeon Token

Currency item `dungeon_token`:
- Guaranteed drop per boss.
- Vendor exclusive item (dungeon-specific cosmetic, mount blueprint).
- BoE (tradeable).

Token vendor city hub'da.

## 8. Set Item Distribution

Dungeon-özel set (bkz. [item-engine/set.md](../item-engine/set.md)):
- Ashen Keep Set → Ashen Keep boss drop.
- Set items pool weight; per-slot balance (%15 helm, %15 chest, ..., %10 ring).

Set completion 2/4/6 piece bonuses (bkz. set.md).

## 9. AFK Dungeon Loot

Idle mode reveal (bkz. [afk-run.md](afk-run.md)):
```
clears_completed = min(clears_per_hour × elapsed_hours, weekly_lockout_remaining × character_slot)

for each clear:
    for each boss encounter:
        loot = boss.loot_table roll (weekly_lockout_effect applied)
    for each treasure room:
        loot = treasure.loot_table roll
    aggregate
```

## 10. Reveal Payload

```json
{
  "dungeon_completions": [
    { "dungeon_code": "ashen_keep", "difficulty": "normal", "clears": 3 }
  ],
  "boss_defeated": [
    { "boss_code": "keep_guardian", "loot": [...], "first_clear_this_week": true },
    { "boss_code": "fire_wraith", "loot": [...] },
    { "boss_code": "ashen_lich_king", "loot": [...], "first_clear_this_week": true }
  ],
  "chests_opened": 5,
  "tokens_gained": { "ashen_keep_token": 120 }
}
```

## 11. Anti-Cheat

Boss loot idempotency:
- Weekly lockout enforcement.
- Duplicate BoP boss item detection.
- Rare drop rate audit (per-day cap).

## 12. Metrics
- `dungeon_loot_awarded_total{dungeon,item}`.
- `dungeon_set_piece_dropped_total{set}`.
- `dungeon_first_clear_reward_total{dungeon}`.
- `dungeon_lockout_penalty_applied_total`.

## 13. Cross-Reference
- [Dungeon](dungeon.md), [Boss](boss.md), [Treasure](treasure.md), [Random Event](random-event.md), [AFK Run](afk-run.md)
- [systems/loot-system.md](../systems/loot-system.md)
- [item-engine/set.md](../item-engine/set.md)
- [database/progression-quests-achievements.md §8](../database/progression-quests-achievements.md)

## 14. Bekleyen Kararlar
- Weekly lockout penalty final tune (rare×0.3 vs 0.5).
- Round-robin loot opt Phase 02.
- Master loot raid Phase 02+.
- Bonus roll consumable (extra roll per boss) Phase 02.
- Cross-realm loot dedup Phase 03.
