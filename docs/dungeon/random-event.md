# Random Event (Dungeon-Inject)

## 1. Amaç
Dungeon içinde random event injection — spice + variety. Encounter'a benzer ancak dungeon-specific pool.

Foundation: [hunting-adventure/random-event.md](../hunting-adventure/random-event.md), [hunting-adventure/encounter.md](../hunting-adventure/encounter.md).

## 2. Dungeon Random Event Türleri

| Kind | Örnek | Etki |
|---|---|---|
| `treasure_boost` | "Ancient Cache" | Sonraki chest %50 daha büyük |
| `mob_wave` | "Skeleton Wave" | Ekstra mob spawn |
| `rare_summon` | "Rare Elite" | Rare elite mob spawn |
| `elite_ambush` | "Cursed Presence" | Random elite spawn ambush |
| `blessing` | "Whisper of the Ancients" | Party buff +10% dmg 5 min |
| `curse` | "Cursed Room" | Party debuff -5% dmg 5 min |
| `secret_door` (Phase 02) | Hidden room reveal | Bonus loot |
| `mimic` (Phase 02) | Chest is actually a mob | Combat trigger |

## 3. Content

```yaml
# content/dungeon_random_events/whisper_of_ancients.yaml
code: whisper_of_ancients
kind: blessing
dungeon_ids: [ashen_keep, silver_crypt]
trigger_chance_per_room_pct: 3           # %3 per random_event room
duration_ms: 300000                        # 5 min buff
effect:
  buff_code: ancients_blessing
  modifiers:
    - kind: mult_add
      target: damage_pct
      amount: 10
  duration_ms: 300000
```

## 4. Injection

Dungeon layout:
```yaml
rooms:
  - kind: enemy_pack
  - kind: random_event
    event_pool: [treasure_boost, blessing, curse, mob_wave]
  - kind: enemy_pack
  - kind: boss
```

Reveal trigger:
```
if room.kind == "random_event":
    pool = room.event_pool
    picked = weightedPick(pool, seed.sub("event-<room_id>"))
    triggerEvent(picked, party, seed)
```

## 5. Effect Apply

- **Buff/Curse**: apply status effect party-wide.
- **Treasure Boost**: next chest loot table × 1.5.
- **Mob Wave**: extra enemy pack next room.
- **Rare Summon**: rare elite mob (bkz. [hunting-adventure/rare-enemy.md](../hunting-adventure/rare-enemy.md)).

## 6. AFK Compute

Reveal (idle mode):
```
for each random_event room:
    event = weightedPick(pool, seed.sub("event-<room_id>"))
    apply_effect(event, snapshot, ...)
    reveal_summary.notable_encounters.append(event)
```

Analytical; deterministic.

## 7. Party vs Solo

- Event applies to all party members.
- Solo: only self.
- Party synergy: buffs/curses affect group.

## 8. Reveal Notification

Reveal panel:
```json
{
  "notable_encounters": [
    { "kind": "dungeon_random_event", "code": "whisper_of_ancients", "effect": "party_dmg_+10%_5min" }
  ]
}
```

## 9. Metrics
- `dungeon_random_event_triggered_total{code}`.
- `dungeon_random_event_buff_uptime_avg`.

## 10. Cross-Reference
- [Room](room.md), [Random Event](../hunting-adventure/random-event.md), [Encounter](encounter.md)
- [combat/status-effects.md](../combat/status-effects.md)

## 11. Bekleyen Kararlar
- Random event pool per dungeon (v1: shared pool; Phase 02 dungeon-özel).
- Mimic mekaniği Phase 02.
- Secret door Phase 02 (fog of war).
- Curse event opt-out (dispel Phase 02).
