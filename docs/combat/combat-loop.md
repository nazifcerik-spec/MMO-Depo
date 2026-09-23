# Combat Loop (Tick + Scheduler)

## 1. Amaç
Live combat için server-side tick loop; rotation eval; state update; event emit. AFK path bu loop'u kullanmaz (analytical yield); ancak manuel oyun + dungeon manuel run buna ihtiyaç duyar.

## 2. Tick Rate

- **Server tick**: 250 ms (default; config).
- **GCD**: 1000 ms baseline (haste ile min 500 ms).
- **Combat state timeout**: 5 sn no-hostile → exit combat.
- **Regen tick**: 2 sn (HP/mana/resource).

## 3. Session Model

Session = ephemeral combat context (manuel oyun veya dungeon):
```
combat_session:
  id, started_at
  participants: [character_ids, npc_ids]
  location_id (or instance_id for dungeon)
```

Redis-based (`combat:session:<id>`); persist opsiyonel (`character.combat_session` Phase 02).

## 4. Tick Loop (High-Level)

```go
func (e *CombatEngine) Tick(ctx context.Context, session Session, deltaMs int64) TickReport {
    report := TickReport{}
    for _, actor := range session.Participants {
        // 1. Status effect ticks (DoT/HoT expire check)
        expiredStatuses := statusEngine.CheckExpirations(ctx, actor.ID)
        for _, st := range expiredStatuses {
            statusEngine.Expire(ctx, st.ID)
            report.StatusExpired = append(report.StatusExpired, st)
        }
        // 2. DoT/HoT tick apply
        activeDotHot := statusEngine.ActiveByKind(ctx, actor.ID, "dot", "hot")
        for _, st := range activeDotHot {
            if shouldTick(st, deltaMs) {
                effectResolver.ApplyPeriodic(ctx, st, actor)
                report.DotHotTicks = append(report.DotHotTicks, st)
            }
        }
        // 3. HP / Mana / Resource regen (every 2 sec)
        if regenDue(actor, deltaMs) {
            regenEngine.Apply(ctx, actor)
        }
        // 4. Auto-attack (weapon speed based)
        if autoAttackDue(actor, deltaMs) {
            result := autoAttack.Execute(ctx, actor, actor.Target)
            report.AutoAttacks = append(report.AutoAttacks, result)
        }
        // 5. AI rotation (NPC) or Player rotation (AFK party dungeon)
        if actor.IsNPC {
            aiEngine.Decide(ctx, actor, session)
        } else if actor.RotationActive {
            skill := rotationEngine.PickSkill(ctx, actor, session.State)
            if skill != nil {
                skillEngine.Cast(ctx, actor, skill, actor.Target, session.Seed.Sub(...))
            }
        }
        // 6. Threat table decay + update (mob AI)
        threatEngine.Decay(ctx, actor)
    }
    // 7. Emit tick events (WS push)
    eventBus.Publish(ctx, "combat.tick", session, report)
    return report
}
```

## 5. Scheduler

Session tick scheduler:
- **Time-based**: 250 ms interval; goroutine per session.
- **Event-driven** (alternate): mesaj gelince tick (client action, skill cooldown ready).

v1 tercih: **hybrid**:
- Continuous 250 ms tick session aktif iken.
- İnactive session (no participants active) → sleep + wakeup on event.

Scheduler resource:
- Concurrent session cap per node: ~1000 (v1 hedef).
- Session hibernate: 5 dk activity yok → hibernate; player action ile resume.

## 6. AFK Path (Aggregate, No Live Tick)

AFK hunting/dungeon idle contract:
- **Live tick loop çalıştırılmaz**.
- Reveal içinde analytical yield compute (bkz. [activity-engine/hunting.md §5](../activity-engine/hunting.md)).
- Damage/kill/loot aggregate; per-tick simulate YOK (DD-084).

Session yalnızca **manuel oyun** için (aktif dungeon manual run, world boss, PvP Phase 02).

## 7. Auto-Attack Interval

Weapon-based; bkz. [auto-attack.md](auto-attack.md):
- `weapon.attack_speed_ms` (2400ms default).
- Haste'a bağlı reduction (max ~ %40).

## 8. Regen Tick (2 sn)

Her aktör için `last_regen_at`:
```go
if now - actor.LastRegenAt >= 2000 {
    hp_regen = derived.HPRegenPer2s × combatMultiplier
    actor.CurrentHP = min(actor.MaxHP, actor.CurrentHP + hp_regen)
    // same for mana/resource
    actor.LastRegenAt = now
}
```

Combat state → `combatMultiplier = 0.2` (HP), `0.5` (mana), `1.0` (energy/rage).

## 9. Rotation Evaluation

Character rotation (bkz. [skill-engine/active-skills.md §6](../skill-engine/active-skills.md)):
```go
for _, entry := range rotation.PriorityList {
    if !conditionMet(entry.Cond, actor, target) { continue }
    if !canCast(entry.SkillID, actor) { continue }
    return entry.SkillID
}
return nil  // fallback auto-attack
```

Auto-attack fires each interval regardless of rotation.

## 10. AI Decision (NPC)

NPC AI decide (per-tick):
```go
switch npc.AIProfile {
case "melee_bruiser":
    if target.InRange(melee) {
        cast(bash_skill)
    } else {
        move_toward(target)  // Phase 02 positional
    }
case "caster":
    if target.InRange(caster) {
        cast(fireball)
    }
case "healer_support":
    lowest_hp_ally := findLowestHPAlly()
    if lowest_hp_ally.HPPct < 50 {
        cast(heal)
    }
}
```

v1'de simple decision tree; Phase 02 behavior tree / FSM.

## 11. WebSocket Push

Tick sonrası:
- Aggregate events per participant.
- Push `combat.tick` frame WS'e (bkz. [backend/websocket.md](../backend/websocket.md)).
- Throttle: max 10 Hz per client.

Frame payload:
```json
{
  "type": "combat_tick",
  "session_id": "...",
  "events": [
    { "kind": "damage_dealt", "actor": "...", "target": "...", "amount": 152, "crit": true },
    { "kind": "auto_attack", "actor": "...", "target": "...", "amount": 45 },
    { "kind": "status_expired", "actor": "...", "status": "berserker_rampage" }
  ]
}
```

## 12. Session Termination

Session sona erer:
- Tüm hostile NPC dead → victory.
- Tüm party wipe → defeat (respawn).
- Tüm participants leave.
- Timeout (5 dk no activity).

Termination:
- Save summary (kills, XP, loot).
- Emit event `combat.session_ended`.
- Cleanup Redis state.

## 13. Concurrency

Session state Redis-based; her tick tek node (goroutine) execute.
Multi-node coordination: session bir node'a "sticky" (Redis leader election).
Node down → session migrate (Phase 02).

## 14. Performance Target

- **Tick latency**: p50 < 50 ms per session.
- **Sessions per node**: 1000 concurrent.
- **Total sessions**: horizontal scale.

## 15. Metrics
- `combat_session_active_count` gauge.
- `combat_tick_duration_seconds{session_kind}` histogram.
- `combat_tick_events_per_tick_avg` gauge.

## 16. Cross-Reference
- [Auto-Attack](auto-attack.md), [Damage](damage.md), [Healing](healing.md), [Status Effects](status-effects.md), [Threat](threat.md), [Death & Revive](death-revive.md), [PvE Encounters](pve-encounters.md), [Combat Events](combat-events.md)
- [skill-engine/active-skills.md](../skill-engine/active-skills.md), [skill-engine/cooldown.md](../skill-engine/cooldown.md)
- [character/health.md](../character/health.md)
- [activity-engine/hunting.md](../activity-engine/hunting.md), [activity-engine/dungeon.md](../activity-engine/dungeon.md)
- [backend/websocket.md](../backend/websocket.md)

## 17. Bekleyen Kararlar
- Session persistence Postgres vs Redis-only (v1 Redis).
- Tick rate 250ms vs 500ms (perf trade-off; v1 250).
- Positional gameplay Phase 02.
- Cross-node session migrate Phase 02.
- Behavior tree AI Phase 02.
