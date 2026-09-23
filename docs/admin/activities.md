# Activities (AFK Contract Inspect)

## 1. Amaç
Aktif AFK contract listesi + inspect + force-reveal (support case).

## 2. List View
- Filter: character, kind (mining/fishing/…), state, started_at range.
- Columns: contract_id · character · kind · started_at · expires_at · progress %.

## 3. Detail
```
Contract: ac_ab…
Character: Grim
Kind: mining
Params: { location: 'wildwood_iron_vein_L15', duration_seconds: 28800 }
Started: 2026-09-09 08:00 UTC
Expires: 2026-09-09 16:00 UTC (progress 45%)
State: active
Snapshot (frozen at start):
  gear_score, mastery_tier, gather_bonus, ...
Estimated yield (server analytical):
  180 iron ore, 12 rare gem, 2 legendary chance
```

## 4. Aksiyonlar

| Action | Permission | Step-up | Notes |
|---|---|---|---|
| `activity.force_reveal` | support_l2+ | no | Server bug workaround; reveal at current tick |
| `activity.cancel` | support_l2+ | no | Yield kayıp; reason zorunlu |
| `activity.snapshot_view` | support_l1+ | no | Read-only |

Force reveal ledger + audit; character notify.

## 5. Endpoint

- `GET /internal/admin/activity?…`
- `GET /internal/admin/activity/{id}`
- `POST /internal/admin/activity/{id}/force-reveal` `{ reason }`
- `POST /internal/admin/activity/{id}/cancel` `{ reason }`

## 6. Cross-Reference
- [systems/activity-engine.md](../systems/activity-engine.md), [systems/reward-engine.md](../systems/reward-engine.md)

## 7. Bekleyen Kararlar
- Force-reveal rate limit / dual approval Phase 02.
- Bulk force-reveal (server incident recovery) Phase 02.
