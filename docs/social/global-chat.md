# Global Chat (World / Trade / LFG / Region)

## 1. Amaç
Bir region ölçeğinde geniş kitleye erişebilen kanallar. Trade / LFG / world talk.

V1'de "gerçekten global" değil — **region-scoped** (Phase 03'te cross-region opt).

## 2. Kanallar

| Kod | Amaç | Scope | Rate/dk | Notlar |
|---|---|---|---|---|
| `global` | Genel sohbet | Region | 5 | Trust threshold 0 |
| `trade` | Alım/satım duyuru | Region | 3 | Otomatik spam filter yüksek |
| `lfg` | Grup arama | Region | 5 | Party UI'dan tek-tık integration |
| `region` | Bölge lokal | Region | 10 | Coğrafi (yeni oyuncu welcome) |

## 3. Persistence

- **Persist edilmez** (log volume). `chat.message`'a yazılmaz.
- **Redis rolling** kanal + scope key altında last 200 mesaj (MAXLEN); TTL 24h.
- **Compliance opt-in sampling** (Phase 02): %1 örnek `audit.chat_sample`'a alınır — abuse forensic.

Yeni join client history:
```
GET /v1/chat/live?channel_kind=global&scope_key=region:{region_id}&limit=100
```
→ Redis LRANGE.

## 4. WS Fanout

Redis pubsub channel key:
```
chat:live:global:region:{region_id}
chat:live:trade:region:{region_id}
chat:live:lfg:region:{region_id}
chat:live:region:region:{region_id}
```

Client region değişiminde önceki region key'lerinden unsub + yeni region key'lerine sub.

Fanout worker Redis pubsub → WS connection map → client push.

## 5. Anti-Spam Katmanı

Global kanallar bot / farm hedefi. Ek kurallar:

- **Trust score min**:
  - `global`: trust ≥ 0
  - `trade`: trust ≥ 10 (7 gün eski account)
  - `lfg`: trust ≥ 0
  - `region`: trust ≥ 0
- **Slow mode**: server-wide config `chat.slow_mode_global_sec` (default 12 sn per karakter, aşım → 30 sn cooldown).
- **Substring blacklist** agresif.
- **Link filter** whitelist only.
- **Duplicate detection**: aynı içerik 5 dk penceresinde reddedilir.
- **Report threshold**: 3+ report / 1 saat → auto-mute 24h.

## 6. Endpoint

**POST `/v1/chat/send`** (kind='global'|'trade'|'lfg'|'region', `channel_ref` region_id).
**GET `/v1/chat/live?channel_kind=trade&scope_key=region:{id}&limit=100`** — canlı rolling scroll.

## 7. Global Announcement (System)

Server → tüm oyuncular (world boss spawn, event start, patch notice):
- Channel kind `system` + fanout tüm region'lara.
- Persist: `chat.message` (system) — 30 gün her karakter için (küçük hacim).

`POST /internal/ops/chat/announce` (admin).

## 8. Metrikler
- `chat_global_message_total{channel_kind, region_id}`
- `chat_global_ws_fanout_latency_ms`
- `chat_global_slow_mode_hit_total`
- `chat_global_auto_mute_total`

## 9. Cross-Reference
- [chat.md](chat.md), [block.md](block.md)
- [backend/redis.md](../backend/redis.md), [backend/websocket.md](../backend/websocket.md)
- [world/region.md](../world/region.md) — region key kaynağı

## 10. Bekleyen Kararlar
- Cross-region trade channel (v1 region-only; Phase 03 opt).
- Slow mode saniyesi (12 v1) Phase 02 tune.
- Public chat 1% sample retention Phase 02 (compliance karar).
- Trade channel için post format template (item link) Phase 02.
