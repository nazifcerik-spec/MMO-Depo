# Router

## 1. Amaç
Route tanımı, versioning, path convention, middleware zinciri. `chi` üzerine ince bir organizasyon katmanı.

## 2. Sözleşme
```go
package httpx

// RouteRegistrar: her domain modülü bunu implement eder.
type RouteRegistrar interface {
    Register(r chi.Router)  // kendi route grubunu ekler
}

// Build: tüm modüllerin route'larını tek chi.Mux üstünde toplar.
func BuildRouter(mws Middlewares, modules ...RouteRegistrar) *chi.Mux
```

Her domain paketi `handler.go`'da:
```go
func (h *CharacterHandler) Register(r chi.Router) {
    r.Route("/characters", func(r chi.Router) {
        r.Post("/", h.create)
        r.Get("/{id}", h.get)
        r.Post("/{id}/allocate-stats", h.allocateStats)
    })
}
```

## 3. Path Convention
- Prefix version: `/v1/...`.
- REST-benzeri isim kaynak; verb'ler nadir (`/reveal`, `/refresh` gibi RPC-ish uç noktalar açık olarak isimlendirilir).
- Kaynak plural: `/characters`, `/items`, `/guilds`.
- Nested resource max 2 seviye; ötesi query param.
- Kebab-case yerine snake yok, path'te lowercase.
- Query param: `snake_case`.

Örnekler:
```
GET  /v1/healthz/live
GET  /v1/healthz/ready
POST /v1/auth/login
POST /v1/auth/refresh
POST /v1/auth/logout
GET  /v1/me
GET  /v1/characters
POST /v1/characters
GET  /v1/characters/{id}
POST /v1/characters/{id}/allocate-stats
POST /v1/characters/{id}/respec
POST /v1/characters/{id}/ascend
POST /v1/characters/{id}/reveal
POST /v1/characters/{id}/activities
DELETE /v1/characters/{id}/activities/{cid}
GET  /v1/characters/{id}/inventory
POST /v1/characters/{id}/equipment
GET  /v1/market/listings?item_template_id=...&sort=price_asc&cursor=...
POST /v1/market/listings
DELETE /v1/market/listings/{id}
POST /v1/market/listings/{id}/buy
GET  /v1/guilds/{id}
POST /v1/guilds
GET  /v1/party/{id}
POST /v1/party
GET  /v1/notifications
```

WebSocket:
```
GET  /v1/ws   → upgrade
```

Metrics/health (public within cluster):
```
GET  /metrics  (port 9090)
GET  /healthz/live
GET  /healthz/ready
```

## 4. Middleware Zinciri (order)
```
1. RequestID
2. RealIP (behind proxy)
3. Recover (panic → 500)
4. OtelTrace
5. Logger
6. Metrics
7. RateLimit (opsiyonel per-route)
8. Compress
9. AuthOptional (parses JWT into ctx if present, doesn't require)
10. RouteAuth (kritik path'lerde middleware zorunlu)
11. Handler
```

Detay: [middleware.md](middleware.md).

## 5. Versioning
- Prefix `/v1/`; breaking change → `/v2/`.
- Route registrar'lar version bilir:
```go
r.Route("/v1", func(r chi.Router) {
    for _, m := range modules { m.Register(r) }
})
```
- Bir handler iki versiyonu destekleyebilir; ancak v2 ayrı kod path'i tercih (basitlik).

## 6. Content Negotiation
- Default `application/json`.
- Accept-Encoding: gzip (Phase 02 brotli).
- Request/response Content-Type zorunlu (POST/PUT/PATCH).

## 7. Idempotency
- `Idempotency-Key` header kritik POST'larda (create character, create listing, buy).
- Header yoksa: server autogen ve reject etme yok (opsiyonel).
- Application layer key → correlation_id + ledger idempotency_key.

## 8. Pagination
- Cursor-based:
  - Request: `?cursor=<opaque>&limit=50`.
  - Response: `{ items: [...], next_cursor: "<opaque or null>" }`.
- Offset yasak (deep offset perf sorunu + inconsistency).

## 9. Error Response
Bkz. [Error System](error-system.md).
```
{
  "error": {
    "code": "character.not_found",
    "message": "Character not found.",
    "request_id": "01J8T...",
    "details": { }
  }
}
```

## 10. CORS
- Origin allowlist config'ten.
- Preflight cache 10 dk.
- Credentials: `include` sadece web client için allowed origin.

## 11. Test Stratejisi
- Handler'lar `httptest.NewRecorder()` ile unit test.
- Route registration `chi.Walk` ile "declared routes" listesi CI check.
- OpenAPI (Phase 02): auto-generated route katalog ile drift check.

## 12. İlişkiler
- [HTTP Server](http-server.md), [Middleware](middleware.md), [Error System](error-system.md), [WebSocket](websocket.md), [Logger](logger.md)

## 13. Bekleyen Kararlar
- `chi` vs `gin` final — öneri chi (net docs, minimal).
- OpenAPI spec-first vs code-first — v1 code-first + CI generator; Phase 02 spec-first değerlendir.
- gRPC service-to-service (Phase 02) — REST yanı sıra.
