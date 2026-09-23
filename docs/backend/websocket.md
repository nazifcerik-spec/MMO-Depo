# WebSocket

## 1. Amaç
Tek connection üstünden multiplex mesajlaşma: activity tick, notification, chat, world event, reveal ready, presence. Server-authoritative push kanalı.

## 2. Sözleşme

### 2.1 Client-side (Flutter)
Tek `/v1/ws` upgrade; JWT `Authorization: Bearer` header ile.

### 2.2 Server hub
```go
package wsx

type Hub struct {
    conns   *shardedMap[uuid.UUID, *Conn]  // characterID → Conn
    log     *slog.Logger
    cache   *cache.Client                   // pub/sub
    cfg     Config
}

type Conn struct {
    id           string   // ulid
    characterID  uuid.UUID
    accountID    uuid.UUID
    send         chan Frame
    closed       chan struct{}
    ws           *websocket.Conn // nhooyr.io/websocket
    lastPingAt   time.Time
}

func NewHub(cfg Config, cache *cache.Client, log *slog.Logger) *Hub
func (h *Hub) HandleUpgrade(w http.ResponseWriter, r *http.Request)
func (h *Hub) Push(ctx context.Context, characterID uuid.UUID, msg Frame) error
func (h *Hub) Broadcast(ctx context.Context, channel string, msg Frame) error
func (h *Hub) Shutdown(ctx context.Context) error
```

## 3. Konfigürasyon
```go
type Config struct {
    ReadBufferSize   int           // 4KB
    WriteBufferSize  int           // 4KB
    HandshakeTimeout time.Duration // 5s
    ReadTimeout      time.Duration // 60s (heartbeat window)
    WriteTimeout     time.Duration // 5s
    PingInterval     time.Duration // 25s
    SendQueueSize    int           // 64 frames
    MaxMessageBytes  int           // 32KB
    ShutdownGrace    time.Duration // 5s
}
```

## 4. Mesaj Formatı
JSON envelope:
```json
{
  "type": "activity_tick",
  "id": "01J8T...",
  "at": "2026-09-08T12:00:00.000Z",
  "payload": { ... }
}
```

Server → Client tipleri:
| Type | Amaç |
|---|---|
| `hello` | Bağlantı kuruldu, initial state hint |
| `activity_tick` | Aktif contract progress ping (opsiyonel; sparse) |
| `reveal_ready` | AFK reveal hesaplandı, refresh önerisi |
| `notification` | Genel notification (yeni chat, market sold, boss spawn) |
| `chat_message` | Chat channel'a mesaj push |
| `world_event` | Server-wide olay (boss spawn, event start) |
| `error` | Hata (bkz. Error System §8) |
| `pong` | Ping cevabı |

Client → Server tipleri:
| Type | Amaç |
|---|---|
| `subscribe` | Chat channel opt-in (`{"channels":["guild:...", "party:..."]}`) |
| `unsubscribe` | |
| `chat_send` | `{"channel":"guild:xxx","body":"..."}` |
| `ping` | Heartbeat |
| `presence` | Client visibility state (foreground/background) |

## 5. Bağlantı Akışı
```
1. Upgrade request → validate JWT (Auth middleware yeniden çalışır)
2. Load session → verify not revoked
3. Character selection: client "?character_id=..." veya sonra "select" mesajı
4. Hub.register(conn) → conns map + Redis WS conn registry (ws:conn:<node>:<char>)
5. Send "hello" with server_time + initial state hint
6. Start read/write goroutines
7. Subscribe hub-level channels for this character (notif + guild + party + zone)
8. On disconnect: unregister, cleanup Redis key
```

## 6. Heartbeat
- Server her 25s `ping` (WebSocket control frame).
- Client 60s içinde `pong` yoksa → close (idle).
- Client uygulama seviyesi `ping/pong` da destekli (ancak stdlib control frame yeterli).

## 7. Backpressure
- Send channel bounded (`SendQueueSize`).
- Dolu ise: yeni frame drop (with metric increment) veya connection kick, tercih:
  - **Broadcast (chat/world_event)**: drop + count.
  - **Personal critical (reveal_ready)**: block up to 500ms, sonra kick.

## 8. Fan-Out (Chat + Notification)
- Chat: `pub:chat:<kind>:<ref>` (Redis pub/sub).
- Notification: `pub:notif:<character_id>`.
- Hub bir subscriber goroutine → mesajlar → ilgili conn'lara push.
- **Multi-node fanout**: Redis pub/sub sayesinde tüm app node'larda ilgili conn'lara ulaşır.

## 9. Auth Rotation
- Access token expire → server 401 error frame + graceful close.
- Client refresh flow (HTTP) → yeni access token ile reconnect.
- Blocklist güncellenirse: `session_id` blocklist Redis key sub → aktif conn kill.

## 10. Multi-Cihaz Politikası
- Aynı karakterle iki cihaz? İki bağlantı kabul (senaryo: telefon + desktop).
- Aynı **session** iki cihaz? Session tekildir; ilk conn kick.

## 11. Rate Limit
- Client → Server mesaj: per-connection 20 msg/sec.
- Chat send: per-channel 3 msg/2sec + global 10/min anti-spam.
- Bounce mesajları hub tarafında; log.

## 12. Fail Modları
- Redis pub/sub down → notification/chat degraded; personal push (activity ready) DB polling fallback (worker'dan).
- Send timeout → conn kick.
- Malformed JSON → error frame + counter.
- Message too large → close 1009.

## 13. Metrics
- `ws_connections_active` gauge.
- `ws_messages_sent_total{type}` counter.
- `ws_messages_recv_total{type}` counter.
- `ws_send_dropped_total{type}` counter.
- `ws_connection_duration_seconds` histogram.

## 14. Shutdown
- SIGTERM → `Hub.Shutdown(ctx)` → aktif conn'lara `close` frame (1001 going away) → drain send buffers → grace 5s → force close.
- Redis subscription iptal.

## 15. Anti-Patterns
- HTTP handler içinde upgrade sonrası WriteTimeout sızıntısı — `nhooyr` ctx-aware, dikkat.
- Frame içinde büyük binary blob (32KB üstü) — dosyalar HTTP.
- Business logic hub'da; hub sadece transport (business logic domain service'te).
- Broadcast'i tek goroutine'de yapıp yavaş peer'ın herkesi bloklaması (per-conn buffer).

## 16. Test Stratejisi
- Unit: Hub.Push mock conn.
- Integration: real WS server + client, message round-trip.
- Load: 20k concurrent conn, msg fanout timing.

## 17. Güvenlik
- Origin check (browser client için).
- Message size limit.
- Auth token in header (query string yasak — leak log/proxy).

## 18. İlişkiler
- [HTTP Server](http-server.md), [Middleware](middleware.md), [Redis](redis.md), [Event Bus](event-bus.md), [Error System](error-system.md), [Logger](logger.md), [systems/notification-system.md](../systems/notification-system.md)

## 19. Bekleyen Kararlar
- WS lib final: `nhooyr.io/websocket` (öneri) vs `gorilla/websocket`.
- Binary framing (msgpack) mi JSON mı — v1 JSON; v2 msgpack optional.
- Compression (permessage-deflate) — Phase 02.
- Sticky session gerekli mi (Redis pub/sub sayesinde hayır).
