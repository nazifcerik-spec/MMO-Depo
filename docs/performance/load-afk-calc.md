# Load Test: AFK Calculation (Reveal Throughput)

## 1. Amaç
Reveal endpoint + worker batch throughput. AFK-first MMO'da **her player'ın 2-4 saatte bir reveal** — server yükü consistent.

## 2. Senaryolar

### 2.1 Steady Reveal
- Yük: T3'te 10.000 CCU × 2 reveal/saat = **~5.5 reveal/sn** system-wide.
- Analitik hesap O(1) (DD-084).
- Beklenti: reveal p95 < 250ms.

### 2.2 Reveal Wave (Mass Login)
- Server maintenance return: 5 dk içinde 10.000 CCU login + reveal 8h AFK contract.
- ~30 reveal/sn peak.
- Beklenti: batch worker throughput handle; p95 < 500ms tolere.

### 2.3 Contract Expire Batch (Worker)
- Worker `activity_expire` her 30 sn: aktif contract'lardan expired olanları scan + tag.
- T3'te 8h contract 3000 aktif → 30 sn'de expire olan sayısı ~1.
- T4 peak: ~5-10 contract/sn expire.

### 2.4 Notification Push Fanout
- Reveal complete → push notification queue insert.
- Rate: reveal rate ile eşit ~5.5/sn (T3).
- FCM/APNs delivery async; queue lag < 5 sn hedef.

## 3. Kritik Path — Reveal

1. Auth + idempotency lookup.
2. Contract fetch (DB).
3. State check (active/expired).
4. **Analitik yield hesap** (Poisson/Binomial sampling — CPU pure, seed-based).
5. Snapshot + current character state + regional event etkileri.
6. Ledger yaz (SERIALIZABLE TX): gold + XP + mastery XP + reputation gain.
7. Inventory item insert (batch).
8. Achievement trigger event.
9. Cache idempotency response.
10. COMMIT.
11. Outbox: WS notify + push notification.

## 4. Bottleneck Analiz

- **Analitik hesap CPU**: 8h contract, 8*3600 tick simulation değil; O(1) sampler → < 5ms.
- **Ledger TX**: SERIALIZABLE + 5-10 entry → 20-50ms.
- **Inventory insert batch**: 10-30 item → single INSERT ... VALUES ... = 10-30ms.
- **Idempotency cache**: Redis GET/SET < 5ms.

Total p95 target: 200-250ms tolerable.

## 5. Analitik Sampler Örnek

```
1h mining, base_rate 40 ore/h, character bonus +25%
expected = 40 × 1.25 = 50 ore/h
For 8h contract: expected 400
Sampler: Poisson(400, seed=<contract_id_hash>)
CPU: nanoseconds
```

Rare drop 1%: Binomial(N=8*3600, p=0.01, seed).
Test: 10.000 reveal sampler → distribution chi-square pass.

## 6. Metrics

- `activity_reveal_p50/p95/p99_ms`.
- `activity_reveal_analytical_calc_ms`.
- `activity_reveal_ledger_write_ms`.
- `activity_reveal_serialization_retry_total`.
- `activity_expire_batch_size`.
- `notif_push_queue_lag_seconds`.

## 7. Worker Batch Model

`activity_reveal_dispatcher`:
- Idle-triggered from notification push (server-timer or user reveal call).
- User-initiated → sync path (ana test).
- Auto-reveal (contract expire silently done, sonuç bekle) — v1'de yok; kullanıcı reveal etmezse contract "waiting" state (max 30 gün mail-safe overflow — Phase 02 opt).

## 8. Notification Push Pipeline

Reveal complete → `event.notification.dispatch` outbox → `push_dispatcher` worker:
- Rate limit check.
- Provider (FCM/APNs) API call (~50-200ms external).
- Async; queue lag < 5 sn.

Test: reveal burst → notification queue lag stabil.

## 9. Beklenen Sonuç

T3 steady:
- Reveal p50: ~80ms.
- Reveal p95: ~200ms.
- Reveal p99: ~600ms.
- Analitik CPU: negligible (< 5ms).
- Push queue lag < 2 sn.

T4 peak (30 reveal/sn):
- p95: ~350ms.
- Worker parallelism = CPU core count; Go goroutine efficient.
- DB write throughput OK (SERIALIZABLE serializable).

## 10. Cross-Reference
- [systems/activity-engine.md](../systems/activity-engine.md), [systems/reward-engine.md](../systems/reward-engine.md)
- [notification/activity-completed.md](../notification/activity-completed.md), [notification/push-notification.md](../notification/push-notification.md)
- [testing/afk.md](../testing/afk.md)

## 11. Bekleyen Kararlar
- Auto-reveal (kullanıcı offline unclaimed contract) Phase 02.
- Batch reveal (birden fazla concurrent contract tek response) Phase 02.
- Analitik sampler benchmark (target < 10ms) Phase 01 doğrula.
