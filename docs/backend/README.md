# Backend Core

Go application server'ının temel yapı taşları. Domain modülleri bu katmanların üstünde çalışır. Her doküman:
- **Amaç ve sözleşme** (interface).
- **Örnek Go iskeleti** (pseudo — Phase 01'de gerçek pakete dönüşür).
- **Fail modları** ve **test stratejisi**.

Domain kodunun bu katmanlara doğrudan bağımlı olmaması — hepsi `internal/platform/*` altında adapter olarak bulunur; domain paketi sadece interface bilir (bkz. [Go Server](../architecture/go.md)).

## Doküman Haritası

| Katman | Amaç | Doküman |
|---|---|---|
| **HTTP Server** | Bootstrap, TLS, graceful shutdown, timeouts | [http-server.md](http-server.md) |
| **Router** | Route grouplama, version, path convention, healthz | [router.md](router.md) |
| **Middleware** | Auth, request-id, logging, metrics, rate-limit, recover, panic, CORS | [middleware.md](middleware.md) |
| **Error System** | Sentinel errors, wrap, mapping (HTTP/WS), user-safe messages | [error-system.md](error-system.md) |
| **Logger** | slog, structured, correlation, sampling | [logger.md](logger.md) |
| **Database** | pgx pool, config, health, retry, sqlc pattern | [database.md](database.md) |
| **Redis** | client pool, health, key convention, Lua scripts | [redis.md](redis.md) |
| **Transactions** | Tx wrapper, isolation, retry, cross-store 2PC yok | [transactions.md](transactions.md) |
| **WebSocket** | Hub, connection, multiplex, heartbeat, backpressure | [websocket.md](websocket.md) |
| **Worker** | Job kuyruk consumer, scheduler, cron | [worker.md](worker.md) |
| **Event Bus** | Domain event yayınlama + subscribe (in-process, ileride pub/sub) | [event-bus.md](event-bus.md) |

## Ortak İlkeler

1. **Context-first**: her fonksiyon `ctx context.Context` alır; iptal edilebilir; deadline'lı.
2. **Interface-first**: domain paketi interface bilir; adapter'lar `platform/*`.
3. **No hidden globals**: singleton yok; DI (constructor injection).
4. **Fail-fast**: config invalid → boot fail; runtime error explicit.
5. **Observability builtin**: her katman `slog` + OTEL trace + Prometheus metric üretir.
6. **Testable**: her katmanın in-memory / fake implementation'ı var; hızlı unit test.
7. **Server authoritative**: hesaplama, RNG, time — daima server tarafında (bkz. GDD §15).

## Paket Konumları (planlı)

```
internal/platform/
  httpx/       — server + router + middleware + error responder
  wsx/         — websocket hub
  db/          — pgx pool + tx helper + sqlc query wrapper
  cache/       — redis client + convention
  bus/         — event bus (in-process → v2 Redis pub/sub)
  jobs/        — worker queue (Postgres SKIP LOCKED → v2 alt)
  auth/        — jwt, session, blocklist
  telemetry/   — otel setup, metrics, tracing
  logging/     — slog handler + helpers
  time/        — Clock interface
  config/      — typed config loader
  rng/         — deterministic RNG streams
  errors/      — sentinel + wrap + mapping
```

## Bekleyen Kararlar
- HTTP router impl: `chi` vs `gin` — öneri **chi**.
- JSON codec: stdlib vs `sonic`/`go-json` — v1 stdlib yeterli.
- Event bus v2: Redis pub/sub vs NATS.
- WebSocket lib: `nhooyr.io/websocket` vs `gorilla/websocket` — öneri **nhooyr** (context-aware, modern).
