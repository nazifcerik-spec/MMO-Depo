# Load Test: WebSocket

## 1. Amaç
WS bağlantı sayısı, ping/pong RTT, fanout (chat/combat/party/guild), reconnect wave, memory footprint.

## 2. Senaryolar

### 2.1 Steady Connection Count
- Yük: T3'te **10.000 concurrent WS conn**.
- Send/receive rate: ~30-100 msg/sn per conn (varies by activity).
- Beklenti: memory/conn < 50 KB; RTT p95 < 50ms.

### 2.2 Connection Storm
- Server restart return → 5 dk içinde 10.000 CCU reconnect.
- ~33 conn/sn peak.
- Beklenti: connect handshake p95 < 150ms; no cascade fail.

### 2.3 Fanout Burst (Global Chat Event)
- 10.000 subscriber `global:region:{id}` kanalında.
- 20 msg/sn peak (event chatter).
- Fanout: 20 × 10.000 = 200.000 msg/sn peak (peak spike).
- Beklenti: fanout p95 < 400ms.

### 2.4 Party Combat Fanout
- 1.000 party (5 kişi) live combat aynı anda.
- Combat event/sn per party × 5 msg = 5000 party fanout msg/sn.
- Beklenti: fanout p95 < 200ms per party.

### 2.5 Idle Ping Load
- 10.000 conn × 1 ping/25sn (keep-alive) = 400 ping/sn.
- Beklenti: p95 < 50ms.

### 2.6 Long-Lived Session Soak
- 4 saat 10.000 conn sabit; drop/reconnect ratio; memory leak.
- Beklenti: drop < %1/saat; memory stable.

## 3. Kritik Path — Connect

1. HTTP Upgrade handshake.
2. Auth token verify (JWT + session lookup Redis).
3. Character bind (active_character_id).
4. Subscribe kanallar (region + party + guild + whisper self).
5. Presence set (Redis TTL 5m).
6. Response (upgrade success + welcome frame).

## 4. Kritik Path — Fanout

Redis pubsub:
```
message published → subscriber connection list lookup
  → iterate subscribers
    → serialize (JSON compact)
    → conn.WriteMessage
    → backpressure check (buffer full → drop or disconnect)
```

## 5. Connection Modeli

- **Go net/http + gorilla/websocket** (v1); alternatif nhooyr/websocket.
- Each conn = 1 goroutine (read) + 1 goroutine (write) + buffered channel.
- Memory footprint per conn: ~20-40 KB base + buffer.
- 10.000 conn × 40 KB = **400 MB** RAM per server node.

## 6. Sub Management

Connection state Redis:
- `ws:conn:{sid}` — active conn metadata.
- Redis SUBSCRIBE `chat:live:*` pattern (server-side worker fanout).
- Application map `sid → conn` — in-memory O(1) lookup.
- Server restart → connection lost; client exp backoff reconnect (bkz. [notification/mobile-lifecycle.md](../notification/mobile-lifecycle.md)).

## 7. Backpressure

Client slow (mobile bad network):
- Write buffer full → drop oldest (chat) or disconnect (combat critical).
- Metric: `ws_backpressure_drop_total{channel}`.
- Client reconnect + `resume_from` replay (5 dk pencere).

## 8. Metrics

- `ws_conn_count_gauge`.
- `ws_msg_sent_total{direction}`.
- `ws_msg_receive_total`.
- `ws_ping_rtt_ms`.
- `ws_fanout_latency_ms`.
- `ws_reconnect_total{reason}`.
- `ws_backpressure_drop_total`.
- `ws_bandwidth_bytes_per_sec`.

## 9. Reconnect Wave Handling

Server restart → 10.000 client reconnect:
- Client exp backoff (1s → 2s → 5s → 15s → 60s + jitter).
- Server: connect handshake < 150ms p95.
- Peak arrival: 30-60 sn içinde ~200-300 conn/sn.
- Auth cache warm-up (JWT verify cheap).

## 10. Load Balancing

- L4 (TCP) LB v1 (nginx/HAProxy).
- Sticky session yok (WS stateless — sub Redis pubsub via any node).
- Multi-node fanout: her node kendi conn'larına yazar; Redis pubsub cluster-wide.

## 11. Beklenen Sonuç

T3 steady:
- Ping RTT p95: ~30-40ms.
- Fanout p95: ~100ms.
- Memory: 400-500 MB per node (10k conn).
- CPU: 20-30% baseline.

T3 chat burst 200k fanout msg/sn:
- Fanout p95: ~250-400ms.
- Redis pubsub CPU spike ~60%.
- WS backpressure drop < 0.1%.

T4 (50.000 conn per node veya multi-node):
- Multi-node scale (2-3 node ile 50.000 conn dağıt).
- Node başına 15-20k conn (memory ~700-1000 MB).

## 12. Cross-Reference
- [backend/websocket.md](../backend/websocket.md), [backend/redis.md](../backend/redis.md)
- [notification/mobile-lifecycle.md](../notification/mobile-lifecycle.md)
- [social/chat.md](../social/chat.md), [social/global-chat.md](../social/global-chat.md)

## 13. Bekleyen Kararlar
- WS conn limit per node (v1: 15k; Phase 02 tune).
- Sticky session (LB stickiness) — v1 stateless yeter.
- NATS pubsub migration Phase 03 (Redis pubsub cluster scale limit).
- Message compression (permessage-deflate) Phase 02 bandwidth optimize.
