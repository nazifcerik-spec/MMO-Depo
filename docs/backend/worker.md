# Worker

## 1. Amaç
Arka plan işleri yürüten Go binary: reveal batch, marketplace expiry, world boss spawn, notification dispatch, economy metrics, cleanup/archive, scheduled tasks.

Kuyruk: `system.job_queue` (Postgres SKIP LOCKED — [DD-039](../decisions/design-decisions.md)). Alternatifler Phase 02.

## 2. Sözleşme
```go
package jobs

type Handler interface {
    Kind() string                                // "reveal_batch", "market_expiry", ...
    Handle(ctx context.Context, payload []byte) error
}

type Worker struct {
    pool       *db.Pool
    handlers   map[string]Handler
    log        *slog.Logger
    metrics    *WorkerMetrics
    cfg        Config
}

type Config struct {
    NodeID           string        // unique per instance
    Concurrency      int           // 8
    PollInterval     time.Duration // 500ms
    LockDuration     time.Duration // 5m (heartbeat renewed)
    MaxAttemptsDefault int         // 5
    ShutdownGrace    time.Duration // 30s
}

func NewWorker(cfg Config, pool *db.Pool, log *slog.Logger) *Worker
func (w *Worker) Register(h Handler)
func (w *Worker) Run(ctx context.Context) error
```

## 3. Job Lifecycle
```
enqueue → queued → picked (running + locked_by=node) → completed
                                                    ↘ failed (attempts < max → back to queued after delay)
                                                    ↘ dead (attempts >= max)
```

## 4. Pick Loop (SKIP LOCKED)
```go
func (w *Worker) pickBatch(ctx context.Context) ([]Job, error) {
    return w.repo.PickBatch(ctx, w.cfg.NodeID, w.cfg.Concurrency)
}
```
SQL:
```sql
UPDATE system.job_queue
SET state='running', locked_by=$1, locked_at=now(), attempts=attempts+1
WHERE id IN (
  SELECT id FROM system.job_queue
  WHERE state='queued' AND run_at <= now()
  ORDER BY run_at ASC
  LIMIT $2 FOR UPDATE SKIP LOCKED
)
RETURNING *;
```

## 5. Execution Pattern
```go
for _, job := range picked {
    go func(j Job) {
        ctx, cancel := context.WithTimeout(runCtx, jobTimeout(j.Kind))
        defer cancel()
        h, ok := w.handlers[j.Kind]
        if !ok {
            w.markDead(ctx, j, "no handler")
            return
        }
        err := h.Handle(ctx, j.Payload)
        switch {
        case err == nil:
            w.markCompleted(ctx, j.ID)
        case errors.Is(err, errs.ErrIdempotent):
            w.markCompleted(ctx, j.ID) // safe replay
        case j.Attempts >= j.MaxAttempts:
            w.markDead(ctx, j, err.Error())
        default:
            w.reschedule(ctx, j.ID, backoff(j.Attempts))
        }
    }(job)
}
```

## 6. Idempotency
- Enqueue optionally `idempotency_key` — unique index → duplicate no-op.
- Handler kendi domain idempotency katmanına yazar (ör. `economy.ledger_entry.idempotency_key`).

## 7. Backoff Policy
```
attempt 1 → immediate
attempt 2 → 30s
attempt 3 → 2m
attempt 4 → 10m
attempt 5 → 1h
> max → dead
```

## 8. Handler Kataloğu (v1)

### 8.1 `reveal_batch`
- **Amaç:** karakter offline'dan döndüğünde reveal hazırlığı; büyük computation için schedule.
- **Trigger:** login öncesi worker doldurabilir (opsiyonel) veya login sırasında async.
- **Payload:** `{"character_id": "..."}`.
- **Impl:** AFK Engine simulate → reveal_log write.

### 8.2 `market_expiry`
- **Amaç:** expire olmuş listing'leri iade + item mailbox'a.
- **Schedule:** her dakika bir kez.
- **Impl:** WHERE status='active' AND expires_at < now(); bulk update + item transfer + mail.

### 8.3 `world_boss_scheduler`
- **Amaç:** spawn window'ları içinde random spawn.
- **Schedule:** her 5 dk.
- **Impl:** template'lerin `spawn_windows_json` kontrol → seed'li spawn.

### 8.4 `event_scheduler`
- **Amaç:** daily/weekly/seasonal event start/stop.
- **Schedule:** her dakika.
- **Impl:** event_template + tarih kontrol → event_participation rows aktive.

### 8.5 `notification_dispatch`
- **Amaç:** notification queue → push/email/in-game.
- **Trigger:** notification insert edildiğinde outbox → job enqueue.
- **Impl:** provider adapter (FCM, SES/SMTP, in-game via Redis pub).

### 8.6 `economy_metrics`
- **Schedule:** günlük (00:15 UTC).
- **Impl:** commodity_price_index + currency_supply_snapshot doldur.

### 8.7 `cleanup_archive`
- **Schedule:** haftalık.
- **Impl:**
  - `soft_deleted_at` 7 gün geçmiş characters → hard delete.
  - expired mail → sender'a iade veya vault.
  - contract hard cap 30 gün → archived.
  - chat_message 30 gün → cold storage.

### 8.8 `weekly_reset`
- **Schedule:** cron Salı 06:00 UTC (bkz. DD-016).
- **Impl:** dungeon_lockout, weekly quest reset, event rotation.

### 8.9 `ledger_reconcile`
- **Schedule:** günlük.
- **Impl:** SUM(delta) per wallet vs `character.wallet.balance`; fark → alarm + admin.

### 8.10 `audit_ship`
- **Schedule:** hourly.
- **Impl:** `audit.event_log` cold-storage export (Phase 02 dedicated).

## 9. Scheduled Task Runner
Ayrı bileşen: `system.scheduled_task` tablosunu okur, cron expr'e göre `system.job_queue`'ya enqueue eder.

```go
type Scheduler struct { pool *db.Pool; log *slog.Logger }
func (s *Scheduler) Run(ctx context.Context) error {
    // tick every 30s
    // for each due task → INSERT INTO job_queue ... ON CONFLICT DO NOTHING
    // update next_run_at
}
```

## 10. Multi-Node Coordination
- Advisory lock (pg_advisory_lock) scheduler'da (tek node scheduler garantisi).
- Job pick SKIP LOCKED zaten multi-node safe.
- Node ID unique (env-based veya UUID).

## 11. Graceful Shutdown
- SIGTERM → ctx cancel.
- Aktif job'lar `ShutdownGrace` içinde tamamlanmaya bırakılır.
- Grace sonrası incomplete: `state='queued', locked_by=NULL` (next pick).

## 12. Fail Modları
- Handler panic → recover → mark failed + stack log.
- DB unreachable → poll loop retry with backoff; alert threshold.
- Handler not registered → dead directly.

## 13. Metrics
- `worker_jobs_processed_total{kind,result}` (result: completed/failed/dead).
- `worker_job_duration_seconds{kind}` histogram.
- `worker_queue_depth{state}` gauge.
- `worker_active_jobs` gauge.

## 14. Observability Örneği
```
{level=info, msg="job completed", kind="reveal_batch", job_id="...", dur_ms=182}
{level=warn, msg="job failed, retrying", kind="market_expiry", job_id="...", attempt=2, next_run="..."}
{level=error, msg="job dead", kind="notification_dispatch", job_id="...", err="..."}
```

## 15. Test Stratejisi
- Unit: handler bazlı; fake queue.
- Integration: testcontainers Postgres + real queue; enqueue → poll → assert state transitions.
- Chaos: kill worker mid-job → next pick sees stale lock (heartbeat expired) → re-picks.

## 16. Anti-Patterns
- Handler'da HTTP call'un timeout'suz olması.
- Job payload büyük (> 32KB) — payload'a küçük id koy, DB'den yükle.
- Long-running job (> 5m) tek job olarak; parçala + resume ile.
- Kritik business logic sadece worker'da (server-driven flow bozulabilir).

## 17. İlişkiler
- [Database](database.md), [Transactions](transactions.md), [Event Bus](event-bus.md), [Error System](error-system.md), [Logger](logger.md), [database/audit-system.md §job_queue](../database/audit-system.md), [systems/afk-engine.md](../systems/afk-engine.md), [systems/notification-system.md](../systems/notification-system.md)

## 18. Bekleyen Kararlar
- Queue impl v2: Redis Streams veya NATS.
- Dead job vault ekranı / admin panel.
- Priority queue (kritik reveal vs low-priority cleanup).
- Job trace propagation (parent request trace_id → job trace continuation).
