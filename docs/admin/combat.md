# Combat (Admin View)

## 1. Amaç
Aktif combat session gözlem + geçmiş forensic + suspicious pattern tespiti (kill farm bot, dmg exploit).

## 2. Active Sessions View
- Filter: character, enemy_type, dungeon_run_id, started_at.
- Columns: session_id · character · enemy(s) · duration · state · skill_cast_count.

## 3. Session Detail
```
Combat session: cs_ab…
Character: Grim (Warrior 42)
Started: 2026-09-09 12:34 UTC
Duration: 3m 28s
State: active

Participants:
  - Grim (self)
  - Enemy: Goblin Raider Lvl 15 HP 720/1000

Skill casts (timeline):
  00:00 auto_attack (dmg 42)
  00:02 whirlwind (dmg 88 crit)
  ...

Damage summary:
  Dealt: 1420
  Taken: 340
  Crit rate: 18%
```

Combat log 24h hot retention (`analytics.combat_event`); 12 ay cold arşiv.

## 4. Suspicious Pattern

- **Kill farm bot**: aynı enemy type 100+ kill / saat → flag.
- **Damage exploit**: outlier dmg (10× expected).
- **Skill spam**: cooldown bypass tespit (client-side try).

Alert dashboard'da; drill-down replay (Phase 02).

## 5. Combat Replay (Phase 02)

Event stream'den replay UI — timeline scrubber + participants HP graph. V1'de sadece log.

## 6. Aksiyonlar

Combat sürerken admin intervention **yok** (server-authoritative; end önce müdahale etme). Post-mortem sadece read + audit.

Farm detection sonucu → account/character review (bkz. players.md, characters.md).

## 7. Endpoint

- `GET /internal/admin/combat/session?…`
- `GET /internal/admin/combat/session/{id}` compound
- `GET /internal/admin/combat/suspicious?…` (Phase 02)

## 8. Cross-Reference
- [combat](../combat/README.md)
- [database/audit.md](../database/audit.md)

## 9. Bekleyen Kararlar
- Combat replay UI Phase 02.
- ML-based anomaly detection Phase 03.
- Real-time intervention (kick from combat) — v1'de yok (risky).
