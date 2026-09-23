# Lifecycle — Start / Pause / Resume / Cancel / Complete

## 1. Amaç
Contract state machine + transition endpoint'leri + event akışı.

## 2. State Machine

```
     ┌─────────┐
     │ (none)  │
     └────┬────┘
          │ Start
          ▼
     ┌─────────┐   Pause    ┌─────────┐
     │ active  │ ─────────► │ paused  │
     │         │ ◄───────── │         │
     └────┬────┘   Resume   └────┬────┘
          │                       │
          │  Complete             │  Cancel
          │  (natural end)        │  (manual)
          ▼                       ▼
     ┌─────────┐            ┌───────────┐
     │completed│            │cancelled  │
     └─────────┘            └───────────┘
                                  │
                             (both →)
                                  ▼
                              ┌────────┐
                              │archived│  (30 gün sonra worker)
                              └────────┘
```

Enum: `activity.contract_status` (bkz. [database/activity.md §1](../database/activity.md)):
- `active`, `paused`, `completed`, `interrupted`, `archived`.

`cancelled` = `interrupted`. Terminoloji birleştirme:
- **completed** — natural end (crafting iterations done, travel arrived).
- **interrupted** — cancelled or forced stop.
- **archived** — 30+ days old; storage optimization.

## 3. Endpoint

### 3.1 Start

**POST `/v1/characters/{id}/activities/start`**
Auth + CharacterScope.

Request:
```json
{
  "activity_type": "mining",
  "channel": "primary",             // primary|pet|mount_travel; default primary
  "params": { "spot_id": "loc-iron-valley", "autoloot": true }
}
```

Server:
1. Load character FOR UPDATE.
2. Validate:
   - Character alive + not dead.
   - Activity type registered.
   - Params schema valid.
   - Channel available (mevcut contract yok o kanalda; varsa auto-close previous).
3. `Activity.ValidateStart(char, params, loc)` — activity-özel gerekler.
4. Transaction:
   - Auto-close previous contract on channel (if any) → Complete or Cancel.
   - `Activity.Snapshot(char, params)` → snapshot row insert.
   - `activity.contract` insert (status='active', last_reveal_at=now(), seed=generate).
   - Emit event `activity.contract_started`.
   - Audit.
5. Response 201: created contract.

Same-channel conflict:
```
Current channel active contract → auto-complete (reveal partial) → new contract start.
```

### 3.2 Pause

**POST `/v1/characters/{id}/activities/{contract_id}/pause`**
Auth + CharacterScope.

Server:
1. Load contract FOR UPDATE.
2. Validate: status='active'.
3. Transaction:
   - Do reveal (partial yield settlement up to now).
   - Update `status='paused'`, `paused_at=now()`, `pause_reason='user'`.
   - Emit event `activity.contract_paused`.
4. Response 200.

Note: pause reveal içerir; yield dispatch edilir. Sonraki resume yeni segmentte hesap.

### 3.3 Resume

**POST `/v1/characters/{id}/activities/{contract_id}/resume`**
Auth + CharacterScope.

Server:
1. Load contract.
2. Validate: status='paused'.
3. Validate: pause_reason'a göre resume koşulları:
   - `inventory_full` → envanter boş slot var mı?
   - `stamina_zero` → stamina > 0?
   - `bait_exhausted` → bait item envanterde?
   - `death` → char alive?
4. Transaction:
   - `Activity.Snapshot` refresh (opsiyonel; v1'de mevcut snapshot kalır).
   - `status='active'`, `last_reveal_at=now()`, `paused_at=NULL`, `pause_reason=NULL`.
   - Emit event `activity.contract_resumed`.
5. Response 200.

### 3.4 Cancel

**DELETE `/v1/characters/{id}/activities/{contract_id}`**
Auth + CharacterScope.

Request query params:
```
?mode=partial      # partial yield settle (default)
?mode=discard      # no yield (throw away)
?mode=full         # complete as if finished (only if allowed)
```

Server:
1. Load contract.
2. Validate: status IN ('active','paused').
3. Transaction:
   - Do reveal (partial yield) unless mode=discard.
   - Update `status='interrupted'`, `updated_at=now()`.
   - `Activity.OnCancel(contract, mode)` → activity-özel cleanup.
   - Emit event `activity.contract_cancelled`.
   - Audit.
4. Response 200: Settlement.

### 3.5 Complete

Natural completion (crafting iterations done, travel arrived):
- Automatic during reveal.
- Not user-triggered directly; complete happens when threshold "complete condition" hit.

Reveal flow triggers complete internally:
- `Activity.DetectThreshold` → `{Kind: "travel_complete"}` or `{Kind: "crafting_iterations_done"}`.
- Complete transaction: yield settle, status='completed', emit event.

Alternative: explicit `POST /v1/characters/{id}/activities/{contract_id}/complete` (v1 opsiyonel; rare).

## 4. Auto-Complete Coşullar

Bazı aktiviteler auto-complete olur:
- **Travel**: destination'a varış.
- **Crafting**: iterations tamamlanma.
- **Adventure**: stage'ler tamamlanma.
- **Dungeon idle**: bounded clears cap.

Mining/Fishing/Hunting **süresiz aktivite** (contract user tarafından cancel edilir veya threshold pauselar).

## 5. Auto-Close on New Start

Same-channel çakışma:
```
Channel: primary
Current: mining (active)
Request: start hunting

Server:
1. Detect conflict.
2. Auto-complete/cancel current (reveal + settle).
3. Start new contract.
```

Client onay flow (opsiyonel v1):
```
Response: 409 { code: "activity.channel_conflict", details: { current_contract_id: "..." } }
Client asks user: "Aktif mining aktivitesini durdur ve hunting başlat?"
Client → POST with force=true.
Server auto-close + start.
```

Default: **auto-close** (Phase 01 UX friction az); Phase 02 confirmation opsiyonu.

## 6. Event Payloads

```json
// activity.contract_started
{
  "_v": 1,
  "contract_id": "...",
  "character_id": "...",
  "activity_type": "mining",
  "channel": "primary",
  "params": { ... },
  "started_at": "..."
}

// activity.contract_paused
{
  "contract_id": "...",
  "pause_reason": "inventory_full",
  "paused_at": "..."
}

// activity.contract_resumed
{
  "contract_id": "...",
  "resumed_at": "..."
}

// activity.contract_cancelled
{
  "contract_id": "...",
  "mode": "partial",
  "cancelled_at": "..."
}

// activity.contract_completed
{
  "contract_id": "...",
  "completed_at": "...",
  "settlement_summary": { ... }
}
```

Subscribers:
- `notification.activity_state`.
- `achievement.activity_progression`.
- Party notification (party members) — Phase 02.

## 7. Concurrency

- Contract updates optimistic lock (`version`).
- Same-channel unique constraint (partial unique index).
- Concurrent Start on same channel → server auto-close previous (first-in-wins).

## 8. Audit

- `activity.contract_started`
- `activity.contract_paused` (severity=info; pause_reason attached)
- `activity.contract_resumed`
- `activity.contract_cancelled`
- `activity.contract_completed`
- `activity.contract_archived` (worker)

## 9. Reveal Trigger

Reveal her contract için otomatik:
- Login (all contracts).
- Pause action (partial reveal).
- Cancel action (partial reveal).
- Complete action (final reveal).
- Explicit reveal call (`POST /v1/characters/{id}/reveal`).

Idempotent (`idempotency_key = hash(contract_id, correlation)`).

## 10. Test Stratejisi

- Start happy path (new contract).
- Start with existing channel → auto-close.
- Pause + resume round-trip.
- Cancel modes (partial/discard/full).
- Auto-complete (travel arrival).
- Concurrent start (race).
- Snapshot immutability across pause/resume.

## 11. Metrics
- `activity_lifecycle_transition_total{from,to,reason}`.
- `activity_pause_reason_total{reason}`.
- `activity_channel_conflict_total`.

## 12. Cross-Reference
- [Activity](activity.md), [Offline Progress](offline-progress.md), [Rewards](rewards.md)
- [database/activity.md](../database/activity.md)
- [systems/activity-engine.md](../systems/activity-engine.md)
- [backend/event-bus.md](../backend/event-bus.md)

## 13. Bekleyen Kararlar
- Auto-close vs confirmation UX (v1 auto; Phase 02 flag).
- Auto-complete explicit endpoint gerekli mi.
- Contract archival worker frequency.
- Cross-contract migration (patch mid-contract).
