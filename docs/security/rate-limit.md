# Rate Limit

## 1. Amaç
Bot, DoS, abuse ve unintentional flood'a karşı endpoint bazlı istek sınırı. Sunucu kaynağı + user experience koruma.

## 2. Algoritma

- **Token bucket** (Redis + Lua script) — burst capacity + refill rate.
- Alternatif: **sliding window log** (audit endpointleri için exact count gerekliyse).
- Rate check middleware **her endpoint önünde** çalışır (allowlist yerine denylist yok — default allowlist).

## 3. Kısıt Boyutları

| Boyut | Örnek |
|---|---|
| **Per IP** | 300 req/dk (DDoS koruma; NAT'lı ağlar için gevşek) |
| **Per Account** | 600 req/dk (login sonrası, hard limit) |
| **Per Character** | Endpoint bazında (bkz. §4) |
| **Per Endpoint (Global)** | Absolute cap (server-wide) — mesela `/v1/market/search` 20k req/dk |

Bir istek tüm ilgili bucket'ları düşer; **herhangi biri tükenirse** 429.

## 4. Endpoint Bucket Kataloğu (v1)

| Endpoint / Grup | Karakter | Account | Notlar |
|---|---|---|---|
| `POST /v1/auth/login` | — | 10/dk / IP | 5 fail sonrası +15dk lock |
| `POST /v1/auth/register` | — | 3/dk / IP | Anti-account-flood |
| `POST /v1/auth/password/reset-request` | — | 5/gün / email | Anti-flood |
| `POST /v1/activity/start` | 30/dk | 60/dk | Contract creation |
| `POST /v1/activity/{id}/reveal` | 60/dk | — | Idempotent; cache hit rate limit'e sayılmaz |
| `POST /v1/vendor/{code}/buy` | 30/dk | 60/dk | — |
| `POST /v1/vendor/{code}/sell` | 30/dk | 60/dk | — |
| `POST /v1/market/listing` | 20/dk | 40/dk | Post cap 50 aktif |
| `POST /v1/market/listing/{id}/buy` | 30/dk | 60/dk | — |
| `GET /v1/market/search` | 20/dk | — | Bkz. [marketplace/search.md](../marketplace/search.md) |
| `POST /v1/trade/request` | 10/dk | 30/dk | — |
| `POST /v1/mail` | 20/dk | 60/gün | Anti-spam |
| `POST /v1/chat/send` | Bkz. [social/chat.md](../social/chat.md) tablo |
| `POST /v1/friends/request` | 10/dk | 30/gün | — |
| `POST /v1/block` | 30/dk | — | — |
| `POST /v1/character/{id}/prestige/node/*/choose` | 60/dk | — | — |
| `POST /v1/combat/{id}/skill` | 20/sn | — | Skill spam window |
| WebSocket connect | — | 5/dk / IP | Reconnect flood koruma |

## 5. Response

`429 Too Many Requests` + header:
```
Retry-After: <seconds>
X-RateLimit-Limit: <N>
X-RateLimit-Remaining: 0
X-RateLimit-Reset: <unix_ts>
```

Client UX: soft toast "Yavaşlayın, X sn sonra tekrar deneyin" + otomatik retry after backoff (jittered).

## 6. Anonymous vs Authenticated

- Anonim (public endpoint — landing, register form) yalnızca IP bazlı.
- Authenticated tüm 3 boyut (IP + account + character).

Log-out sonrası aynı IP kısa süre daha strict (5 dk).

## 7. Distributed Rate Limiter

- Redis Cluster tek tenant instance v1.
- Lua script atomic decrement + TTL.
- Cache miss (Redis down) → fail-open opt (mümkün olduğunca istek geçir; alert P0).

## 8. Circuit Breaker

Server-wide load metric üstünde:
- CPU > 85% 30sn → tüm non-critical endpoint'ler 429 (chat, market search geçici kesme).
- DB pool > 90% → aynı.
- Kritik (login, character sheet load, combat) her zaman geçer.

## 9. Metrikler
- `sec_rate_limit_hit_total{endpoint, boyut}`
- `sec_rate_limit_bucket_utilization{endpoint}`
- `sec_circuit_breaker_trip_total{reason}`

## 10. Cross-Reference
- [backend/middleware.md](../backend/middleware.md)
- [backend/redis.md](../backend/redis.md)

## 11. Bekleyen Kararlar
- Client-provided rate limit hint (adaptive Phase 02).
- Per-region rate override (yüksek pop region özel) Phase 02.
- Fair queue (bucket dolu ise sıraya al vs reject) Phase 03.
