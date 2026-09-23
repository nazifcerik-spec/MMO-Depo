# Random Event

## 1. Amaç
Region-wide veya global scale event injection — geçici bonus/challenge. Encounter'dan farklı: **duration-based** (event window), **broadcast** (all characters aware).

## 2. Event Türleri

| Kind | Örnek | Duration | Scope |
|---|---|---|---|
| `region_event` | "Bandit Uprising in Ashen Hills" | 1-4 saat | Region |
| `world_event` | "Blood Moon" | 24 saat | Global |
| `mini_event` | "Rare mob wave" | 30 min | Location |
| `emergent_event` | "Meteor Strike" | 1 saat | Region + notify |
| `seasonal` | "Harvest Festival" | 2 hafta | Global |

Bkz. [game-design/events.md](../game-design/events.md).

## 3. Model

```sql
CREATE TABLE content.random_event_template (
  id UUID PRIMARY KEY,
  code VARCHAR(48) UNIQUE NOT NULL,
  kind VARCHAR(24) NOT NULL,
  scope VARCHAR(16) NOT NULL,          -- region/world/location
  region_ids UUID[] NOT NULL DEFAULT '{}',
  location_ids UUID[] NOT NULL DEFAULT '{}',
  duration_ms BIGINT NOT NULL,
  buff_effects_json JSONB NOT NULL,     -- character-side buff during event
  spawn_modifier_json JSONB NULL,        -- mob spawn rate/loot mod during event
  reward_pool_id UUID NULL,
  broadcast_notification BOOLEAN NOT NULL DEFAULT true,
  version INT NOT NULL DEFAULT 1
);

CREATE TABLE world.random_event_instance (
  id UUID PRIMARY KEY,
  template_id UUID NOT NULL,
  region_id UUID NULL,
  location_id UUID NULL,
  started_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  ends_at TIMESTAMPTZ NOT NULL,
  status VARCHAR(16) NOT NULL DEFAULT 'active',   -- active/completed/expired
  contribution_meta_json JSONB NOT NULL DEFAULT '{}'
);
```

## 4. Content Pack

```yaml
# content/random_events/bandit_uprising.yaml
code: bandit_uprising
kind: region_event
scope: region
region_ids: [ashen_hills]
duration_ms: 7200000                # 2 hours

spawn_modifier:
  elite_bandit_captain: 3.0         # elite bandit spawn ×3 during event
  loot_rare_chance_mult: 1.5        # loot rare +50%

buff_effects: []                    # no character buff

reward_pool_code: bandit_uprising_reward
broadcast_notification: true
```

Örnek world event:
```yaml
# content/random_events/blood_moon.yaml
code: blood_moon
kind: world_event
scope: world
duration_ms: 86400000               # 24 hours

buff_effects:
  - kind: mult_add
    target: crit_dmg_mult
    amount: 0.10                    # +10% crit dmg for everyone
  - kind: mult_add
    target: xp_gain_pct
    amount: 15                       # +15% XP

spawn_modifier:
  undead_family: 2.0                # undead mobs ×2 spawn

reward_pool_code: blood_moon_reward
broadcast_notification: true
```

## 5. Trigger

Event scheduler (bkz. [backend/worker.md §8.4](../backend/worker.md)):
- Random trigger daily (per region template'ye göre).
- Cron-based: seasonal event fixed dates.
- Server-wide random: worker rolls chance every 30 min.

Trigger workflow:
```
1. Worker picks eligible template (region_events, world_events).
2. Roll spawn chance.
3. If trigger:
   a. Insert world.random_event_instance.
   b. Apply buff to eligible characters (region or world scope).
   c. Modify spawn rates for duration.
   d. Broadcast notification (bkz. [systems/notification-system.md](../systems/notification-system.md)).
4. Duration timer → auto-expire.
```

## 6. Character Interaction

Event active while character AFK:
- Snapshot content_version includes event.
- Buff modifier applied.
- Spawn modifier reflects in rate compute.
- Reveal payload highlights event participation.

Event participation credit tracked (`progression.event_participation` bkz. [database/progression-quests-achievements.md §5](../database/progression-quests-achievements.md)).

## 7. Contribution & Reward

Event participation: kill count, boss damage, quest completion during event.

Reward:
- Event token → vendor exclusive cosmetic/mount.
- Achievement.
- Rare drop.

Reward claim: event end sonrası mailbox veya event vendor.

## 8. Broadcast

Event start notification:
- All characters in scope: in-game toast + push (opt-in).
- Chat channel announce.
- Web dashboard event indicator.

## 9. Scheduled vs Emergent

- **Scheduled**: fixed cron (blood moon every 30 days, harvest festival Q4).
- **Emergent**: random trigger (bandit uprising chance 5% per day).
- **Player-triggered** (Phase 02): quest completion → server event.

## 10. Reveal Aggregate

AFK reveal içinde:
- Active event during activity → contribution counted.
- Buff modifier applied to yield.
- Payload:
  ```json
  {
    "active_events": [
      { "event_code": "blood_moon", "hours_participated": 6 }
    ],
    "event_rewards_earned": [ ... ]
  }
  ```

## 11. Endpoint

**GET `/v1/events/active`** — server-wide active events (public).

**GET `/v1/characters/{id}/events/participation`** — character event history.

## 12. Metrics
- `event_triggered_total{kind,scope}`.
- `event_participants_count{event}` gauge.
- `event_reward_claimed_total`.

## 13. Audit
- `event.triggered`.
- `event.character_participated`.
- `event.reward_claimed`.

## 14. Cross-Reference
- [game-design/events.md](../game-design/events.md)
- [Encounter](encounter.md), [Adventure](adventure.md)
- [systems/notification-system.md](../systems/notification-system.md)
- [backend/worker.md §8.4](../backend/worker.md)
- [database/progression-quests-achievements.md §5](../database/progression-quests-achievements.md)

## 15. Bekleyen Kararlar
- Event scheduler cron pattern final.
- Player-triggered events Phase 02.
- Event opt-out (kaotik event immunity) Phase 02.
- Cross-region cascading events Phase 03.
