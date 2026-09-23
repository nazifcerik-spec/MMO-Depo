# Combat Events (Event Bus Emissions)

## 1. Amaç
Combat sırasında oluşan tüm önemli olayların event bus'ta yayınlanması. Passive trigger, notification, achievement, session log — hepsi bu event'lere abone.

Foundation: [backend/event-bus.md](../backend/event-bus.md), [skill-engine/passives.md](../skill-engine/passives.md).

## 2. Event Kataloğu (Combat)

### 2.1 Cast & Skill
- `combat.skill_used` — skill cast complete.
- `combat.skill_interrupted` — cast broken.
- `combat.skill_reject_cooldown` — reject.

### 2.2 Damage & Heal
- `combat.damage_dealt` — direct damage.
- `combat.damage_taken` — inverse of dealt.
- `combat.dot_tick` — DoT periodic.
- `combat.heal_dealt` — heal cast.
- `combat.hot_tick` — HoT periodic.
- `combat.shield_absorbed` — shield tick.

### 2.3 Crit / Miss
- `combat.crit_dealt`.
- `combat.crit_taken`.
- `combat.evaded`.
- `combat.blocked`.

### 2.4 Auto-Attack
- `combat.auto_attack` (high volume — sample).

### 2.5 Status Effects
- `combat.status_applied`.
- `combat.status_expired`.
- `combat.status_dispelled`.
- `combat.status_refreshed`.

### 2.6 CC / Interrupt
- `combat.cc_applied`.
- `combat.cc_dr_hit`.
- `combat.cc_immune`.

### 2.7 Threat
- `combat.threat_generated` (aggregate; not per event).
- `combat.taunt_applied`.
- `combat.target_changed`.

### 2.8 Death & Revive
- `combat.death_recorded`.
- `combat.revived`.
- `combat.mob_killed`.
- `combat.boss_defeated`.
- `combat.world_boss_defeated`.

### 2.9 Session
- `combat.session_started`.
- `combat.session_ended`.
- `combat.wipe`.
- `combat.entered_combat`.
- `combat.exited_combat`.

### 2.10 Encounter
- `combat.encounter_started`.
- `combat.encounter_phase_transition`.
- `combat.boss_enraged`.
- `combat.adds_spawned`.

## 3. Event Envelope

Standart shape (bkz. [backend/event-bus.md §3](../backend/event-bus.md)):
```json
{
  "id": "01J8T...",
  "kind": "combat.damage_dealt",
  "occurred_at": "2026-09-09T15:00:00.000Z",
  "actor_id": "char-777",
  "subject_kind": "character",
  "subject_id": "npc-...",
  "payload": {
    "_v": 1,
    "damage": 152,
    "kind": "direct",
    "element": "physical",
    "crit": true,
    "skill_code": "berserker_rampage"
  },
  "correlation_id": "session-...",
  "trace_id": "..."
}
```

## 4. Subscriber Katalog

### 4.1 Passive Trigger Engine

Bkz. [skill-engine/passives.md §5](../skill-engine/passives.md):
- `on_hit_dealt` → passive triggers (frenzy, execute, ripcast).
- `on_crit_dealt` → crit boost triggers.
- `on_kill` → kill triggers (bloodthirst).
- `on_damage_taken` → defensive triggers.
- `on_hp_below_pct` → survival triggers.

### 4.2 Achievement Engine

- `combat.boss_defeated` → boss achievement.
- `combat.mob_killed` → kill milestone.
- `combat.world_boss_defeated` → first-server check.
- `combat.session_ended` → clean run (no death).

### 4.3 Notification System

- `combat.death_recorded` → death notification.
- `combat.boss_defeated` → celebration toast.
- `combat.world_boss_defeated` → global announcement.

### 4.4 Mastery Engine

- `combat.mob_killed{family}` → bestiary XP.
- `combat.skill_used` → weapon mastery XP.
- `combat.heal_dealt` → healing discipline mastery XP.
- `combat.damage_taken` → tanking mastery XP.

### 4.5 Session Log / Recap

- All events log to session summary (post-combat).
- Client "combat log" UI (Phase 02 opsiyonel).

### 4.6 WebSocket Push

- Tüm major event WS push:
  - `damage_dealt`, `heal_dealt`, `death_recorded`, `boss_defeated`, ...
- Throttle: max 10 Hz per client; sample if higher.

### 4.7 Audit Log

- Kritik event `audit.event_log`:
  - `combat.boss_defeated` (severity=info).
  - `combat.world_boss_defeated` (severity=info + broadcast).
  - `combat.wipe` (severity=info).
  - Suspicious pattern (rapid kill anomaly) → security_signal.

## 5. Publish Discipline

Domain modül event publish (bkz. [backend/event-bus.md §7](../backend/event-bus.md)):
- **Fire-and-forget** for analytics/notification.
- **Outbox pattern** for kalıcı side-effect (achievement, ledger).
- Publish in-process (fast); Redis pub/sub cross-node (WS fanout).

## 6. Event Ordering

Combat event ordering within tick:
```
1. damage_taken (per hit)
2. status_applied (if any triggered)
3. death_recorded (if HP=0)
4. mob_killed (if target died)
5. kill_credit (attributor)
6. loot roll
7. session_ended (if wipe)
```

Ordered emit; subscribers receive in order per topic.

## 7. High-Volume Sampling

Auto-attack + DoT tick high volume:
- Sample 10% for analytics.
- Full-fidelity per player session (session log).
- Aggregate per-minute for metrics.

## 8. Event Test

- Unit: publish → subscriber receive assert.
- Ordering: multi-event sequence assert.
- Idempotent subscriber (achievement duplicate check).

## 9. Metrics
- `combat_event_published_total{kind}`.
- `combat_event_dispatch_lag_ms{kind}` histogram.
- `combat_event_subscriber_error_total{subscriber}`.

## 10. Cross-Reference
- [Combat Loop](combat-loop.md), [Damage](damage.md), [Healing](healing.md), [Status Effects](status-effects.md), [Death & Revive](death-revive.md), [PvE Encounters](pve-encounters.md)
- [backend/event-bus.md](../backend/event-bus.md)
- [skill-engine/passives.md](../skill-engine/passives.md) — trigger engine
- [systems/notification-system.md](../systems/notification-system.md)
- [systems/reward-system.md](../systems/reward-system.md)

## 11. Bekleyen Kararlar
- Event schema versioning (`_v` bump policy) Phase 02.
- Session log persistence (session summary DB) Phase 02.
- Combat log client UI Phase 02.
- Cross-region event fanout Phase 03.
