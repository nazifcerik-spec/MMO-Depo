# Dungeons (Active Runs + Template)

## 1. Amaç
Aktif dungeon run gözlem + template rollout + weekly lockout gözlem.

## 2. Active Runs
- Filter: character, dungeon_template, difficulty, state (active/succeeded/failed).
- Columns: run_id · character/party · template · difficulty · floor/total · started_at · state.

## 3. Run Detail
```
Run: dr_ab…
Party: Grim + Aria + Kael
Template: iron_mine (heroic)
Started: 2026-09-09 14:12 UTC
Floor: 3/5
State: active (afk mode)
Loot so far: (personal per member)
  Grim: 2 uncommon, 1 rare
  Aria: 3 uncommon
  Kael: 1 rare
```

## 4. Template Browser

Content-driven; template list + preview (encounter, loot table, difficulty tiers). Read-only.

Yeni template rollout content pipeline (Phase 02 tool).

## 5. Weekly Lockout Overview
- Character bazında lockout counter (weekly reset Salı 06:00 UTC — DD-016).
- Support use case: lockout reset (rare bug case) admin step-up.

## 6. Aksiyonlar

| Action | Permission | Step-up |
|---|---|---|
| `dungeon.run_cancel` | support_l2+ | no |
| `dungeon.force_complete` | admin | yes | Server bug; reward as-is |
| `dungeon.lockout_reset` | admin | yes |

## 7. Endpoint

- `GET /internal/admin/dungeon/template`
- `GET /internal/admin/dungeon/run?…`
- `GET /internal/admin/dungeon/run/{id}`
- `POST /internal/admin/dungeon/run/{id}/cancel`
- `POST /internal/admin/dungeon/run/{id}/force-complete`
- `POST /internal/admin/dungeon/lockout/{character_id}/{template_code}/reset`

## 8. Cross-Reference
- [dungeon](../dungeon/README.md)

## 9. Bekleyen Kararlar
- Force-complete reward calculation (mid-progress) Phase 02.
- Bulk lockout reset (server incident) Phase 02.
