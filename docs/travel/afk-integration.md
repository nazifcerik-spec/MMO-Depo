# AFK Integration

## 1. Amaç
Travel'ın offline devam etmesi: contract-based (activity engine), reveal ile arrival algılama.

## 2. Contract Persistence

Travel `activity.contract` üstünde (bkz. [activity-engine/travel.md](../activity-engine/travel.md)):
- `activity_type='travel'`.
- `channel='primary'`.
- `params_json.expected_duration_ms`, `expected_arrival_at`.
- `started_at`, `last_reveal_at`.

Karakter offline olsa da contract yaşar.

## 3. Offline Continuation

Login → reveal:
```go
elapsed_ms = now - contract.last_reveal_at
elapsed_total = now - contract.started_at
if elapsed_total >= contract.params.expected_duration_ms {
    // Arrived
    character.CurrentLocationID = contract.params.destination_location_id
    contract.Status = "completed"
    yield = Yield{XPCharacter: discoveryXP(origin, destination)}
    emit event 'travel.arrived'
}
```

## 4. Soft Cap

Travel süresi genelde <24h; soft cap (bkz. [activity-engine/offline-progress.md §3](../activity-engine/offline-progress.md)) uygulanır ama nadir travel süresi soft cap'e ulaşır (ekstrem cross-continent).

Uygulama:
- `elapsed_effective = softCap(elapsed_since_last_reveal)`.
- Travel bitmişse `elapsed_effective >= expected_duration - last_reveal_progress`.

## 5. Reveal Kararları

Reveal 3 senaryo:
1. **Arrived**: elapsed >= expected → complete + arrival XP.
2. **Still traveling**: elapsed < expected → no yield; contract active.
3. **Long overdue**: elapsed >> expected (30+ gün) → auto-cancel + safe fallback (character teleport home city).

## 6. Notifications

- Travel start: WS push "Travel started; arrival at ..." (client'a).
- Progress update (opsiyonel): every 25% progress → toast.
- Arrival: push notification (opt-in mobile) "You have arrived at Emerald Coast".
- Long overdue: e-mail (Phase 02).

## 7. Auto-Next Activity (Phase 02)

Travel + queued next activity:
```json
{
  "auto_next_activity": {
    "activity_type": "fishing",
    "params": { "spot_id": "loc-emerald-coast" }
  }
}
```

Arrival → auto-start next contract (aynı reveal içinde chained transaction).

v1'de kapalı; kullanıcı arrival sonra manual start.

## 8. Snapshot Determinism

Travel snapshot içinde:
```json
{
  "race_speed_bonus_pct": 5,
  "mount_speed_multiplier": 2.0,
  "equipment_speed_bonus_pct": 12,
  "buff_codes_active": ["speed_potion_minor"]
}
```

Snapshot al; travel süresince değişimler etkilemez (deterministic hesap).

Karakter mid-travel yeni item takarsa veya buff bitse:
- Contract'ın expected_duration_ms değişmez.
- Yeni travel yeni snapshot ile.

## 9. Party Travel

v1'de her karakter kendi travel'ı; party synchronization yok.

Phase 02:
- Party travel: leader'ın hızı bağlayıcı (yavaş üye tüm party'yi yavaşlatır).
- Karavan simülasyonu.

## 10. Cancel Mid-Flight

User cancel:
- v1: character.current_location_id = origin (rollback).
- Ledger: waygate/ship cost non-refundable; walking/mount cost yok.
- Contract status='interrupted'.

Phase 02: midway waypoint (halfway location).

## 11. Interrupt Events (Phase 02)

- Ambush → combat trigger → travel pause.
- Weather event mid-travel → speed re-eval (Phase 02 dynamic).
- Emergency teleport consumable → travel cancel + instant to bind point.

## 12. Metrics
- `travel_completed_offline_total` — arrived while offline.
- `travel_completed_online_total`.
- `travel_long_overdue_total{days}`.

## 13. Audit
- `travel.completed_via_reveal` (offline arrival).
- `travel.auto_cancelled_overdue`.

## 14. Test Stratejisi
- Offline arrival simulate (elapsed > expected).
- Long overdue safe fallback.
- Snapshot immutability across reveal.
- Multi-reveal (partial then complete).

## 15. Cross-Reference
- [Travel](travel.md), [Base Time](base-time.md), [Distance](distance.md)
- [activity-engine/travel.md](../activity-engine/travel.md), [activity-engine/offline-progress.md](../activity-engine/offline-progress.md), [activity-engine/lifecycle.md](../activity-engine/lifecycle.md)
- [systems/notification-system.md](../systems/notification-system.md)

## 16. Bekleyen Kararlar
- Long overdue auto-cancel window (v1: 30 gün + safe fallback).
- Auto-next activity Phase 02.
- Party travel Phase 02.
- Interrupt event framework Phase 02.
