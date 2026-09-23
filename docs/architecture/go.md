# Go (Application Server)

## 1. Amaç
Domain modüllerini barındıran, API isteklerini karşılayan, tick loop'larını yürüten stateless Go binary. Ölçekleme birimi.

## 2. Sürüm ve Ekosistem
- **Go**: en son stabil (Phase 01 başlangıcında pin; pattern: son iki minor sürümden birini destekle).
- **Module system**: `go mod`.
- **Web framework**: `chi` (minimal, middleware-friendly) veya `gin` — Phase 01 karar; öneri `chi`.
- **DB**: `pgx` (native Postgres, hızlı, prepared statement).
- **Migration**: `goose` veya `atlas`.
- **Redis client**: `redis/go-redis/v9`.
- **Config**: `spf13/viper` veya env-only (öneri: env-only + typed struct decoder).
- **Logging**: `log/slog` (stdlib, structured).
- **Metrics**: `prometheus/client_golang`.
- **Tracing**: `go.opentelemetry.io/otel`.
- **Testing**: `stretchr/testify`, `testcontainers-go` (Postgres/Redis).
- **Codegen**: `sqlc` (query-first typed SQL) veya raw pgx — öneri `sqlc`.
- **UUID**: `google/uuid`.
- **Cron / scheduling**: `robfig/cron/v3` (workers only).

## 3. Repo Yapısı (öneri)
```
/cmd
  /server         # API + WS entrypoint
  /worker         # background jobs entrypoint
  /migrate        # migration tool wrapper
/internal
  /character
  /race /class /specialization
  /skills /combat
  /activity /afk
  /inventory /equipment /items
  /economy /marketplace /ledger
  /professions /loot
  /world /travel /dungeons /world_bosses
  /party /guild /social
  /quests /achievements /mastery /ascension
  /notification
  /rng
  /platform
    /db        # pgx pool, sqlc queries
    /cache     # redis wrapper
    /bus       # event bus
    /auth      # jwt, session
    /time      # clock abstraction
    /telemetry # tracing/metrics/logging
    /httpx     # middleware, error responder
    /wsx       # websocket hub
    /config    # typed config loader
/api
  /openapi.yaml
  /proto        # (Phase 02 opsiyonel)
/deploy
  /docker
  /helm         # (Phase 02)
/docs
  ← already exists
/scripts
  /dev, /ci, /release
/test
  /integration
  /load
go.mod, go.sum, Makefile, .golangci.yml
```

## 4. Domain Module Sözleşmesi
Her modül şu iskeleti izler:
```
internal/<domain>/
  entity.go        // domain type'ları
  service.go       // use-case fonksiyonları
  repository.go    // interface tanımı
  repository_pg.go // Postgres impl (sqlc çağrılarını sarar)
  errors.go        // domain-specific error type'ları
  events.go        // yayınlanan event type'ları
  handler.go       // HTTP handler'ları (chi route grubu)
  handler_test.go
  service_test.go
```

Kurallar:
- **Domain paketi standart library ve sadece diğer domain interface'lerine bağımlı**.
- Postgres, Redis, HTTP framework detayları domain'e sızmaz.
- Cross-module: interface + dependency injection; domain-to-domain event bus ile.
- Test edilebilirlik zorunlu: her repository interface'inin in-memory impl'i mevcut.

## 5. Concurrency Kuralları
- Her request kendi context'ini alır; **context propagation zorunlu**.
- Goroutine leak yok: her `go` çağrısı `ctx.Done` respect eder.
- Long-lived worker: `context.WithCancel` + `errgroup`.
- Shared mutable state → yok (Redis'e taşı) veya `sync.RWMutex` explicit.
- Fan-out pattern: bounded worker pool; unbound goroutine spam yok.
- `time.After` yerine `time.NewTimer` (leak korumalı).

## 6. Error Handling
- Sentinel errors (`errors.Is`) domain seviyesinde tanımlı: `ErrNotFound`, `ErrConflict`, `ErrForbidden`, `ErrValidation`.
- Handler layer sentinel → HTTP status mapping (404, 409, 403, 400).
- Panic yalnızca gerçekten bug (assertion); recovery middleware var.
- Error wrap: `fmt.Errorf("%w: ...", ErrX, ...)` — cause zinciri korunur.

## 7. Observability
- **Logging**: `slog` JSON handler; request_id, character_id, activity_id her log satırında.
- **Metrics**: prometheus histogram (latency), counter (event), gauge (concurrent conns).
- **Tracing**: OTEL SDK; her HTTP handler ve DB call span üretir.
- Correlation id propagation: `X-Request-ID` header + context.

## 8. Time Abstraction
```go
type Clock interface {
    Now() time.Time
    Since(t time.Time) time.Duration
}
```
- Production: `system.Clock` (UTC).
- Tests: `fakeclock`; deterministik.
- **Domain kodu doğrudan `time.Now()` çağırmaz** — clock injection.

## 9. Deterministic RNG Adaptör
- `internal/rng` paketi seed'li stream sağlar.
- Tick simulator seed + tick_index alır; wall clock ile karışmaz.
- Test golden-replay ile korunur.

## 10. Code Quality
- **Formatter**: `gofmt` + `goimports` — commit hook.
- **Linter**: `golangci-lint` config'i (`errcheck`, `govet`, `staticcheck`, `revive`, `gosec`, `sqlclosecheck`).
- **Test coverage** hedefi: domain paketleri > %80; handler > %60.
- **Race detector** CI'da: `go test -race`.
- **Dependency audit**: `govulncheck` haftalık.

## 11. HTTP Sözleşmesi
- `application/json` payload.
- Error body: `{ "error": { "code": "...", "message": "...", "details": {...} } }`.
- Pagination: cursor-based (opaque token) — offset yok.
- Idempotency-Key header desteği (kritik POST).
- ETag opsiyonel (GET karakter sheet için Phase 02).

## 12. Anti-Patterns
- Global var mutable state.
- `init()` içinde network / DB çağrısı.
- Domain paketinden HTTP framework import.
- `panic()` iş akışı olarak.
- Test edilemeyen `time.Now()` bağımlılığı.
- Uzun function (>80 satır); büyük struct (>15 alan) → parçala.

## 13. Deployment
- Static binary: `CGO_ENABLED=0`.
- Multi-stage Dockerfile (build → distroless runtime).
- Graceful shutdown: `SIGTERM` → drain HTTP → wait max 30s.
- Readiness probe: `/healthz/ready` (DB + Redis ping).
- Liveness probe: `/healthz/live`.

## 14. İlişkiler
- [System Architecture](system-architecture.md), [PostgreSQL](postgresql.md), [Redis](redis.md), [Docker](docker.md), [Config](config.md), [Logging](logging.md), [CI/CD](cicd.md)

## 15. Bekleyen Kararlar
- HTTP router: `chi` mi `gin` mi (öneri chi).
- SQL codegen: `sqlc` vs raw pgx (öneri sqlc).
- Job queue impl: Postgres SKIP LOCKED vs Redis Stream vs NATS.
- gRPC ekleme (Phase 02 opsiyonel, service-to-service için).
