# Death & Revive

## 1. Amaç
Karakter HP 0 → death sequence. Respawn/revive modları + XP debt + durability loss + auto-revive.

Foundation: [character/health.md §5](../character/health.md).

## 2. Death Trigger

Server-side: `current_hp <= 0` any source.

Race condition:
- Multi-hit last blow: sıralı damage → HP dips negative → death.
- Kill credit: highest damage attributor (v1); Phase 02 assist system.

## 3. Death Sequence

Bkz. [character/health.md §5.3](../character/health.md):
```
BEGIN;
1. UPDATE character.character
   SET is_dead=true, died_at=now(),
       respawn_available_at=now() + interval '30 sec',
       version = version + 1
   WHERE id=$1;
2. Insert character.death_log:
   character_id, cause_kind, cause_ref_id, died_at,
   xp_debt_incurred, durability_lost
3. Compute XP debt (level >= 20):
   debt = level_xp_required(level) × 0.05
   char.xp_debt += debt
4. Durability loss:
   for equipped item: durability -= 5% × max_durability
5. Interrupt active activity contract (bkz. activity-engine/lifecycle.md).
6. Party/dungeon: wipe check.
7. Emit event 'combat.death_recorded'.
8. WS push 'death' frame.
COMMIT;
```

## 4. Kill Credit

Multi-attacker:
```
Kill credit rules:
  1. Last-hit attributor (v1 default).
  2. Highest damage attributor (Phase 02 alternative).
  3. Party assist: if any party member did >10% total damage → credit.
```

XP reward per kill goes to credit-holder. Party split (Phase 02 alt) — v1 winner-take-all with party bonus.

## 5. Respawn Modes

### 5.1 Graveyard Respawn (Self)

**POST `/v1/characters/{id}/revive`**
```json
{ "mode": "graveyard" }
```

Server:
1. Validate `now >= respawn_available_at`.
2. Set `current_location_id = home_city_id.spawn_point`.
3. `current_hp = max_hp × 0.5`.
4. `current_mana/resource = 0`.
5. `is_dead=false`, timestamps clear.
6. Apply "Weakness" debuff 60 sn (-10% dmg out & healing in).
7. Emit event `combat.revived`.

### 5.2 Party Revive Skill

Cleric "Resurrect" skill:
- Target dead within 60 sn.
- Cast time 5 sec (channel).
- On success:
  - Prompt dead player accept.
  - On accept: revive at death location.
  - HP = max_hp × 0.5.

Prompt window: 30 sec; ignore → auto-decline.

### 5.3 Resurrect Scroll (Consumable)

- Item template: `resurrect_scroll` (rare drop).
- Self use on death.
- Revive at death location; HP = max_hp × 0.5.
- Consumed.

### 5.4 Auto-Revive (5 min)

Passive: 5 min post-death auto-revive at graveyard.
- Config `auto_revive_sec: 300`.
- Idle-friendly.
- Weakness applied.

## 6. Weakness Debuff

Post-revive:
- 60 sn.
- -10% damage dealt.
- -10% healing received.
- Stacking: each death within 60s → +1 stack (max 3); -30% penalty max.

## 7. XP Debt (Level 20+)

Bkz. [character/xp.md §7](../character/xp.md):
```
if char.level >= 20:
    debt_amount = level_xp_required(char.level) × 0.05
    char.xp_debt = min(char.xp_debt + debt_amount, level_xp_required(char.level))
```

Debt cap: 1 level worth.

Next XP grant → pay debt first (bkz. [character/xp.md §7](../character/xp.md)).

## 8. Durability Loss

Bkz. [inventory/equipment.md §12](../inventory/equipment.md):
- Each equipped item: `durability_current -= 5% × durability_max`.
- Item passive/proc disabled at durability 0 (base stats kalır).
- Repair NPC city.

## 9. Item Loss

**Yok** — DD-019: karakter item düşürmez ölümde.

Hardcore mode (Phase 03) opsiyonu.

## 10. AFK Death Handling

Hunting AFK reveal içinde:
- Death count from Poisson sample (bkz. [activity-engine/hunting.md §5](../activity-engine/hunting.md)).
- Each death:
  - Insert death_log row.
  - XP debt accumulate.
  - Durability loss.
- Contract paused during death respawn window (30 sn each).
- Auto-revive after 5 min.
- Continue if params.stop_on_death_cascade false.

Death cascade prevention:
- `params.stop_on_hp_below_pct: 30` safety.
- Or `stop_after_n_deaths: 3` (Phase 02).

## 11. Wipe Check (Party/Dungeon)

Party context:
- All members dead → wipe.
- Wipe triggers:
  - Dungeon: encounter reset (mobs full HP).
  - World boss: contribution counted; boss remains.
  - Adventure: adventure fail.

Wipe recovery:
- Party members respawn (graveyard or party revive).
- Instance not reset (retryable).

## 12. Boss Kill Reward

Boss killed:
- Loot roll (bkz. [systems/loot-system.md](../systems/loot-system.md)).
- XP grant.
- Achievement trigger.
- First kill (server-first): unique reward + title.

## 13. Event Payloads

`combat.death_recorded`:
```json
{
  "character_id": "...",
  "cause_kind": "mob",
  "cause_ref_id": "npc-...",
  "location_id": "...",
  "xp_debt_incurred": 2500,
  "durability_lost_avg_pct": 5
}
```

`combat.revived`:
```json
{
  "character_id": "...",
  "mode": "graveyard",
  "new_location_id": "...",
  "weakness_stacks": 1
}
```

## 14. Metrics
- `combat_deaths_total{cause_kind}`.
- `combat_revive_total{mode}`.
- `combat_xp_debt_accumulated_total`.
- `combat_deaths_per_hour_avg`.

## 15. Audit
- `character.died` (audit + character.death_log).
- `character.revived`.
- Suspicious rapid death cycle → security_signal.

## 16. Cross-Reference
- [Combat Loop](combat-loop.md), [Damage](damage.md), [PvE Encounters](pve-encounters.md)
- [character/health.md](../character/health.md), [character/xp.md](../character/xp.md)
- [inventory/equipment.md](../inventory/equipment.md) — durability loss
- [systems/loot-system.md](../systems/loot-system.md)
- [activity-engine/hunting.md](../activity-engine/hunting.md)

## 17. Bekleyen Kararlar
- Auto-revive 5 min v1 vs 3 min Phase 02.
- Weakness stack cap 3 vs 5.
- XP debt %5 vs %3 (playtest).
- Hardcore permadeath Phase 03.
- Corpse retrieval mechanic (WoW classic) — v1 skip.
