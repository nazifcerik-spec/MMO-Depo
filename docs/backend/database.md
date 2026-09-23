# Database (pgx pool)

## 1. Amaç
PostgreSQL bağlantı havuzu; sqlc-generated query'ler için wrapper; health check; retry; migration entry point.

Şema tanımları: [/docs/database/](../database/README.md). Bu doküman Go tarafını anlatır.

## 2. Sözleşme
```go
package db

type Pool struct {
    pgx *pgxpool.Pool
    log *slog.Logger
}

type Config struct {
    URL              string
    MaxConns         int32         // 50
    MinConns         int32         // 5
    MaxConnLifetime  time.Duration // 30m
    MaxConnIdleTime  time.Duration // 5m
    HealthCheckPeriod time.Duration // 30s
    StatementTimeout time.Duration // 3s
    AppName          string        // "mmo-server-<sha>"
}

func New(ctx context.Context, cfg Config, log *slog.Logger) (*Pool, error)
func (p *Pool) Close()
func (p *Pool) Ping(ctx context.Context) error

// Query helpers
func (p *Pool) Q(ctx context.Context) DBTX          // pgx-compatible interface for sqlc
func (p *Pool) WithTx(ctx context.Context, opts TxOptions, fn func(TX) error) error
```

`DBTX` sqlc'nin generated `Querier`'ı için gerekli interface.

## 3. Boot & Config
```go
cfg, err := pgxpool.ParseConfig(url)
if err != nil { return nil, err }
cfg.MaxConns = c.MaxConns
cfg.MinConns = c.MinConns
cfg.MaxConnLifetime = c.MaxConnLifetime
cfg.MaxConnIdleTime = c.MaxConnIdleTime
cfg.HealthCheckPeriod = c.HealthCheckPeriod
cfg.ConnConfig.RuntimeParams["application_name"] = c.AppName
cfg.ConnConfig.RuntimeParams["statement_timeout"] = strconv.Itoa(int(c.StatementTimeout / time.Millisecond))
cfg.ConnConfig.Tracer = otelpgx.NewTracer()  // OTEL integration

pool, err := pgxpool.NewWithConfig(ctx, cfg)
if err != nil { return nil, err }
if err := pool.Ping(ctx); err != nil { return nil, err }
```

## 4. Repository Pattern (Domain)
```go
package character

type Repository interface {
    ByID(ctx context.Context, id uuid.UUID) (Character, error)
    ByName(ctx context.Context, name string) (Character, error)
    Create(ctx context.Context, c Character) error
    UpdateOptimistic(ctx context.Context, c Character, expectedVersion int32) error
}

// Postgres impl
type pgRepo struct{ q Querier }

func NewPGRepo(q Querier) Repository { return &pgRepo{q: q} }

func (r *pgRepo) ByID(ctx context.Context, id uuid.UUID) (Character, error) {
    row, err := r.q.CharacterByID(ctx, id)
    if err != nil {
        if errors.Is(err, pgx.ErrNoRows) { return Character{}, errs.ErrNotFound }
        return Character{}, fmt.Errorf("character.byid: %w", err)
    }
    return mapRow(row), nil
}
```

## 5. sqlc Pattern
- SQL query dosyaları `internal/character/query.sql`.
- Codegen: `sqlc generate` (yaml config).
- Generated tipler `internal/character/db.go` (private).
- Domain paketi kendi tip'ini mapper ile döndürür.

Örnek `query.sql`:
```sql
-- name: CharacterByID :one
SELECT id, account_id, name, race_id, class_id, level, xp, version, created_at
FROM character.character
WHERE id = $1 AND soft_deleted_at IS NULL;

-- name: CharacterCreate :exec
INSERT INTO character.character (id, account_id, name, race_id, class_id, seed)
VALUES ($1, $2, $3, $4, $5, $6);

-- name: CharacterUpdate :execrows
UPDATE character.character
SET level=$2, xp=$3, updated_at=now(), version=version+1
WHERE id=$1 AND version=$4;
```

## 6. Transaction Helper
Bkz. [transactions.md](transactions.md). Örneği:
```go
err := pool.WithTx(ctx, db.TxOptions{Iso: pgx.Serializable}, func(tx db.TX) error {
    q := character.NewQuerier(tx)
    // ... multiple queries
    return nil
})
```

## 7. Health Check
- `Ping` her 30 sn health probe endpoint tarafından.
- Pool `Stat()` metrics'e yansır:
  - `db_pool_acquired`, `db_pool_idle`, `db_pool_total`, `db_pool_wait_count`.

## 8. Retry Policy
- pgx-level retry yok (transient error application'a).
- **Serialization failure** (`40001`): 3 retry, exponential jitter (10ms, 50ms, 250ms).
- **Deadlock** (`40P01`): 2 retry.
- Diğer: propagate.

```go
func retrySerializable(ctx context.Context, fn func() error) error {
    delays := []time.Duration{10*time.Millisecond, 50*time.Millisecond, 250*time.Millisecond}
    var err error
    for i := 0; ; i++ {
        err = fn()
        if err == nil { return nil }
        var pgErr *pgconn.PgError
        if !errors.As(err, &pgErr) || (pgErr.Code != "40001" && pgErr.Code != "40P01") {
            return err
        }
        if i >= len(delays) { return err }
        select {
        case <-time.After(delays[i] + jitter()):
        case <-ctx.Done(): return ctx.Err()
        }
    }
}
```

## 9. Timeout Katmanları
- pool cfg `statement_timeout` = 3s.
- Handler kendi `ctx` deadline'ı (2s tercih; DB call bunu geçmez).
- Uzun query (analytics, worker) explicit `WithTimeout(ctx, 30s)`.

## 10. Enum / JSON Mapping
- Postgres ENUM → Go string type + validator.
- JSONB → domain struct (`json.Unmarshal`); sqlc buna destek verir.
- UUID → `github.com/google/uuid`.
- `TIMESTAMPTZ` → `time.Time` (UTC).

## 11. Multi-Schema Kullanımı
- Search path: `SET search_path TO auth, character, content, world, activity, economy, social, companion, progression, audit, system, public;`.
- Explicit prefix her query'de tercih edilir; okunabilirlik.
- Migration'lar per-schema klasörlü.

## 12. Migration Runner
- `cmd/migrate` binary; goose veya atlas.
- Deploy sırasında migration önce (bkz. [CI/CD](../architecture/cicd.md) §Migration Deploy).
- Idempotent, forward-only.
- CLI:
  ```
  migrate up
  migrate status
  migrate create <name>
  ```

## 13. Fail Modları
- Pool exhausted → new request `context.DeadlineExceeded`; metric spike.
- Connection lost mid-tx → pgx returns; retry policy uygulanır.
- Migration lock (advisory) çakışma → boot fail with clear message.

## 14. Test Stratejisi
- **Unit**: repository interface mock; domain fonksiyon test.
- **Integration**: `testcontainers-go` Postgres 16; her test suite kendi migrations up + seed + teardown.
- **Golden query**: sqlc-generated query'lerin plan'ları (EXPLAIN) beklenen index kullanır.

## 15. Metrics
- pool state.
- Query duration histogram (per prepared statement name).
- Serialization failure counter.
- Slow query (> 500ms) counter + log.

## 16. Anti-Patterns
- Domain paketinden doğrudan `pgxpool` çağırmak (interface bypass).
- `SELECT *`.
- N+1 (loop içinde query).
- Long-running transaction (> 500ms warning).
- ORM lazy load / magic — kabul yok, sqlc explicit.

## 17. İlişkiler
- [Transactions](transactions.md), [Router](router.md), [Error System](error-system.md), [Logger](logger.md), [architecture/postgresql.md](../architecture/postgresql.md), [Worker](worker.md)

## 18. Bekleyen Kararlar
- sqlc vs raw pgx final — öneri sqlc.
- Migration tool: goose vs atlas.
- Read replica routing (Phase 02); read-only pool ayrı.
- Prepared statement cache stratejisi.
