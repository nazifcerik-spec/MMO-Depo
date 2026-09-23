# Travel (Entity + Lifecycle + Endpoint)

## 1. Amaç
Travel'ın somut model + API + state machine tanımı. Activity contract'ının travel-özel yüzü.

Ortak activity lifecycle: [activity-engine/lifecycle.md](../activity-engine/lifecycle.md). Bu doküman travel-özel farkları anlatır.

## 2. Travel Params Şeması

```json
{
  "destination_location_id": "loc-emerald-coast",
  "method": "mount",
  "route_id": null,                    // opsiyonel; fixed route (wagon/ship)
  "waygate_code": null                  // opsiyonel; waygate transport
}
```

Server compute:
```json
{
  "destination_location_id": "...",
  "method": "mount",
  "route_json": [
    { "from": "loc-ashfield", "to": "loc-emerald-coast", "distance_units": 500, "terrain": "grass_road" }
  ],
  "cost_paid_copper": 0,
  "expected_duration_ms": 720000,
  "expected_arrival_at": "2026-09-08T20:12:00Z",
  "route_seed": "..."
}
```

Route computed at start; contract'a persist.

## 3. Endpoints

### 3.1 Estimate

**POST `/v1/characters/{id}/travel/estimate`**
Auth + CharacterScope.

Request:
```json
{
  "destination_location_id": "loc-emerald-coast",
  "method": "mount"     // opsiyonel; default = best available
}
```

Response 200:
```json
{
  "route": [
    { "from": "loc-ashfield", "to": "loc-emerald-coast", "distance_units": 500, "terrain": "grass_road" }
  ],
  "expected_duration_ms": 720000,
  "cost_copper": 0,
  "modifiers_applied": {
    "base_speed": 1.0,
    "mount_multiplier": 2.0,
    "race_multiplier": 1.05,
    "terrain_multiplier": 0.9,
    "equipment_multiplier": 1.08,
    "road_bonus": 1.15,
    "weather_multiplier": 1.0
  }
}
```

### 3.2 Start Travel

**POST `/v1/characters/{id}/travel`**
Auth + CharacterScope.

Request:
```json
{
  "destination_location_id": "...",
  "method": "mount",
  "waygate_code": null,
  "accept_cost_copper": 0
}
```

Server:
1. Load character FOR UPDATE.
2. Validate:
   - Character alive.
   - Destination exists + unlocked.
   - Method valid + prerequisites (mount active for method=mount, waygate unlocked for waygate).
   - Sufficient funds (waygate/ship cost).
   - `accept_cost_copper >= computed_cost` (opsiyonel guard).
3. Snapshot compute (race, mount, equipment, current location).
4. Compute expected_duration + route.
5. Auto-close current primary channel contract (auto-close DD-defined).
6. Create `activity.contract`:
   - activity_type='travel'.
   - params_json = full route + expected values.
   - channel='primary'.
   - status='active'.
7. Debit cost via ledger (waygate/ship).
8. Emit event `travel.started`.
9. Response 201: contract + expected_arrival_at.

### 3.3 Current Travel

**GET `/v1/characters/{id}/travel/current`**
```json
{
  "contract_id": "...",
  "origin_location_id": "...",
  "destination_location_id": "...",
  "method": "mount",
  "started_at": "...",
  "expected_arrival_at": "...",
  "elapsed_ms": 300000,
  "remaining_ms": 420000,
  "progress_pct": 41.6
}
```

### 3.4 Cancel

**POST `/v1/characters/{id}/travel/{contract_id}/cancel`**
Server:
1. Load contract.
2. Validate: status IN ('active', 'paused').
3. Handle:
   - **Walk/mount**: character.current_location_id remains origin (v1: no midway).
   - **Waygate**: cost non-refundable; contract cancelled.
   - **Ship/wagon**: partial refund (Phase 02).
4. Contract status='interrupted'.
5. Emit event `travel.cancelled`.
6. Response 200.

v1 basitleştirme: cancel → origin geri; hiç ilerleme yok.
Phase 02: midway waypoint (halfway location).

## 4. Complete (Arrival)

Reveal içinde threshold detection:
```go
if elapsedMs >= contract.ExpectedDurationMs {
    // Complete
    character.CurrentLocationID = params.DestinationLocationID
    contract.Status = "completed"
    yield = Yield{XPCharacter: discoveryXP(origin, destination)}
    emit event 'travel.arrived'
    // Auto-next activity (Phase 02): oyuncu queue tanımlamışsa auto-start.
}
```

## 5. Route Persistence

`params.route_json` — travel başlangıcında server compute.

Route değişimi (patch mid-contract):
- Content update olsa da mevcut contract kararlı (snapshot pattern).
- Yeni travel yeni route ile.

## 6. Fail Modları

- `travel.destination_unavailable` (deprecated/locked).
- `travel.method_unavailable` (mount yok, waygate unlocked değil).
- `travel.insufficient_funds`.
- `travel.route_not_found` (no path calculable).
- `travel.already_at_destination`.

## 7. Concurrency

- Character primary channel unique constraint (bkz. [database/activity.md §1](../database/activity.md)).
- Concurrent travel start → auto-close previous.

## 8. Auto-Next Activity (Phase 02)

Client rows:
```json
{ "auto_next_activity": { "activity_type": "mining", "params": {...} } }
```

Arrival → next contract auto-start.

v1'de yok; user explicit next action.

## 9. Multi-Character Travel (Party)

- v1: her karakter kendi travel.
- Party formation birleşik mount (Phase 03): karavan simülasyonu.

## 10. Metrics
- `travel_started_total{method}`.
- `travel_completed_total{method}`.
- `travel_cancelled_total{method,reason}`.
- `travel_duration_avg{method}` histogram.
- `travel_route_units_avg` gauge.

## 11. Audit
- `travel.started` (severity=info).
- `travel.arrived`.
- `travel.cancelled`.

## 12. Test Stratejisi
- Estimate → start → arrival roundtrip.
- Method mismatch (mount without mount) → 400.
- Waygate cost + revert on rollback.
- Cancel mid-travel → origin.
- Reveal with elapsed > expected → complete.

## 13. Cross-Reference
- [Distance](distance.md), [Base Time](base-time.md), [Mount](mount.md), [Race](race.md), [Terrain](terrain.md), [Equipment](equipment.md), [Weather](weather.md), [AFK Integration](afk-integration.md)
- [activity-engine/travel.md](../activity-engine/travel.md), [activity-engine/lifecycle.md](../activity-engine/lifecycle.md)
- [world/location.md](../world/location.md)

## 14. Bekleyen Kararlar
- Midway waypoint (partial cancel) Phase 02.
- Auto-next activity Phase 02.
- Party travel coordination Phase 03.
- Emergency teleport (rare consumable) Phase 02.
