# Middleware

## 1. Amaç
Cross-cutting concerns: request id, real ip, panic recovery, tracing, logging, metrics, rate limit, compression, auth. Handler kodunu temiz tutmak, tek yerden yönetmek.

## 2. Standart Zincir (Global)
```
RequestID → RealIP → Recover → OtelTrace → Logger → Metrics → Compress
                                                                  ↓
                                             (route-scoped) → RateLimit → Auth → Handler
```

## 3. Middleware Kataloğu

### 3.1 RequestID
- `X-Request-ID` header var → context'e koy.
- Yoksa `ulid.Make()` üret, response header'a echo.
- Downstream (DB, log, trace) hepsi context'ten okur.

### 3.2 RealIP
- Proxy arkasında `X-Forwarded-For` / `X-Real-IP` parse.
- Config allowlist (proxy IP'leri); dışarıdan spoofing engel.

### 3.3 Recover
```go
func Recover(log *slog.Logger) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            defer func() {
                if rec := recover(); rec != nil {
                    stack := debug.Stack()
                    log.ErrorContext(r.Context(), "panic recovered",
                        "err", fmt.Sprint(rec),
                        "stack", string(stack),
                        "path", r.URL.Path,
                    )
                    metrics.PanicTotal.WithLabelValues(r.URL.Path).Inc()
                    httpx.WriteError(w, r, errors.ErrInternal)
                }
            }()
            next.ServeHTTP(w, r)
        })
    }
}
```

### 3.4 OtelTrace
- `otelhttp.NewMiddleware("mmo-api")`.
- Span attributes: route, method, status, request_id.
- Trace context propagated to DB/Redis calls.

### 3.5 Logger
- Structured request start/finish log:
```
{level=info, msg="http request", method=POST, path=/v1/characters/xxx/reveal, status=200, dur_ms=42, request_id=..., account_id=..., character_id=... (if scoped)}
```
- Body log YOK (PII / size). Sadece meta.
- Debug seviyesinde request headers whitelisted subset log.

### 3.6 Metrics
```
http_requests_total{method, route_template, status}
http_request_duration_seconds{method, route_template} (histogram)
http_request_body_bytes{method, route_template} (histogram, sampled)
http_response_body_bytes{method, route_template} (histogram, sampled)
http_in_flight_requests
```
Route template `chi.RouteContext(r.Context()).RoutePattern()`.

### 3.7 RateLimit
```go
type RateLimiter interface {
    Allow(ctx context.Context, subject, bucket string, cost int) (allowed bool, remaining int, resetAt time.Time, err error)
}
```
- Token bucket (Redis Lua atomic).
- Subject key: `account_id` (auth'ed) veya `ip_hash` (anonymous).
- Bucket: per endpoint kod (örn. `market.list`, `auth.login`).
- Response header:
  - `X-RateLimit-Remaining`
  - `X-RateLimit-Reset` (unix seconds)
- 429 döndüğünde `Retry-After` header.

Fail policy:
- Redis unreachable → fail-closed on `auth.*`, `market.*`, `char.mutate`; fail-open on read-only browse.

### 3.8 Compress
- Gzip encoder; Accept-Encoding'e göre.
- Response > 1KB için etkin (küçük payload boşuna).

### 3.9 AuthOptional
- `Authorization: Bearer <jwt>` parse.
- Signature verify (jwt_key registry, kid).
- Claim'leri `ctx`'e koy: `account_id`, `session_id`, `roles`, `character_id?`.
- Blocklist (Redis) kontrol.
- Hatalı token → no-op (endpoint auth zorunlu değilse); expired → 401 ile açık ret.

### 3.10 AuthRequired
- ctx içinde account_id yoksa → 401.
- Role/permission check (opsiyonel per-route).

### 3.11 CharacterScope
- Route param `{id}` (character id) veya body → **karakter sahiplenmesi** kontrol.
- Karakter account_id'ye ait değilse → 403.
- Karakter soft-deleted → 404.

### 3.12 CORS
- `chi/cors` ile allowlist origin.
- Preflight 10 dk cache.

### 3.13 BodyLimit
- Content-Length > 1MB → 413.
- Handler-özel override (bir dungeon replay dump route'u için higher).

### 3.14 IdempotencyKey
- POST kritik endpoint'lerde `Idempotency-Key` header parse; ctx'e koy.
- Handler kendi idempotency ledger'ına yazar (economy.ledger_entry, quest_reward_claim).

## 4. Zincir Uygulama Örneği
```go
r := chi.NewRouter()
r.Use(mws.RequestID)
r.Use(mws.RealIP)
r.Use(mws.Recover)
r.Use(mws.OtelTrace)
r.Use(mws.Logger)
r.Use(mws.Metrics)
r.Use(mws.Compress)
r.Use(mws.CORS)
r.Use(mws.BodyLimit(1 << 20))

r.Route("/v1", func(r chi.Router) {
    // Public
    r.Group(func(r chi.Router) {
        r.Use(mws.RateLimit("public", 60))
        r.Post("/auth/login", authH.login)
        r.Post("/auth/refresh", authH.refresh)
    })

    // Authed
    r.Group(func(r chi.Router) {
        r.Use(mws.AuthRequired)
        r.Use(mws.RateLimit("authed", 300))
        r.Get("/me", meH.get)

        // Character-scoped
        r.Route("/characters/{id}", func(r chi.Router) {
            r.Use(mws.CharacterScope)
            r.Post("/reveal", activityH.reveal)
            r.Post("/activities", activityH.start)
            r.Get("/inventory", inventoryH.list)
        })
    })
})
```

## 5. Fail Modları
- Middleware error handler'a "bypass" bırakmaz; kendi error response'unu yazar.
- Panic → Recover; log stack; 500.
- Auth token invalid → 401 + `WWW-Authenticate: Bearer`.
- RateLimit exceeded → 429 + Retry-After.
- Body too large → 413.

## 6. Context Kanonik Anahtarları
```go
type ctxKey int
const (
    ctxRequestID ctxKey = iota
    ctxAccountID
    ctxSessionID
    ctxCharacterID
    ctxIdempotencyKey
    ctxUserRoles
    ctxRealIPHash
)
```
Domain kodu bunları helper (`ctxutil.CharacterID(ctx)`) ile okur.

## 7. Test Stratejisi
- Middleware unit test: fake handler + `httptest.NewRecorder`.
- Chain integration test: kompozit route + assertion (header, status).
- Rate limit test: table-driven (allow, block, reset).

## 8. İlişkiler
- [Router](router.md), [Error System](error-system.md), [Logger](logger.md), [Auth (Go Server)](../architecture/go.md), [Redis](redis.md) (rate limit backend)

## 9. Bekleyen Kararlar
- CORS allowlist final (env config).
- Body limit default per route family.
- Idempotency-Key server-side deduplication catalog (Phase 02 opsiyonel).
- CAPTCHA / bot mitigation middleware (Phase 02+).
