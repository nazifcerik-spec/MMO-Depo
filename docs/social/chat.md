# Chat Engine

## 1. Amaç
Tüm kanalların (global, region, party, guild, whisper) ortak mesaj altyapısı. Server-relay + WS delivery + rate limit + moderation.

## 2. Kanal Taksonomisi

| Kanal | Kod | Kapsam | Persist | Latency Hedefi |
|---|---|---|---|---|
| Global (world) | `global` | Bölge veya dünya (v1: region) | Redis last 200; no long-term | < 500ms p95 |
| Trade | `trade` | Region (opsiyonel opt-in) | Aynı | < 500ms |
| LFG (looking for group) | `lfg` | Region | Aynı | < 500ms |
| Region local | `region` | Region | Aynı | < 500ms |
| Party | `party` | 2-5 member | Redis 100 msg + personal scroll | < 200ms |
| Guild | `guild` | Guild members | Redis 200 msg + personal scroll | < 200ms |
| Guild officer | `guild_officer` | Guild officers only | Aynı | < 200ms |
| Whisper (DM) | `whisper` | 1-to-1 | Personal scroll only (30d) | < 200ms |
| System | `system` | Server → char | Personal scroll | server push |

Persistence:
- **Redis rolling window** kanal başına (last 100-200) yeni join'e history vermek için.
- **Personal scroll** (kişisel client cache + server-side son 30 gün whisper/guild/party) sadece o karakterin **kendi** okuduğu mesajlar için.
- **Global public kanallar long-term persist edilmez** (compliance / hacim); yalnızca live WS + Redis rolling.

## 3. Domain Model

```sql
-- Kalıcı: whisper + system + guild + party (personal scroll için)
CREATE TABLE chat.message (
  id BIGSERIAL PRIMARY KEY,
  channel_kind VARCHAR(16) NOT NULL,
  channel_ref UUID NULL,               -- party_id, guild_id, whisper: NULL (from/to pair)
  from_account_id UUID NOT NULL,
  from_character_id UUID NOT NULL,
  from_display_name VARCHAR(64) NOT NULL,
  to_account_id UUID NULL,             -- whisper only
  to_character_id UUID NULL,           -- whisper only
  body TEXT NOT NULL,                  -- max 512 char
  sent_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  trust_score SMALLINT NOT NULL DEFAULT 0,
  moderation_flag VARCHAR(24) NULL     -- 'clean'|'auto_muted'|'reported'
);
CREATE INDEX idx_chat_msg_channel_ts ON chat.message (channel_kind, channel_ref, sent_at DESC);
CREATE INDEX idx_chat_msg_whisper ON chat.message (channel_kind, from_character_id, to_character_id, sent_at DESC)
  WHERE channel_kind = 'whisper';

-- Redis rolling window (public global/trade/lfg/region): LIST veya STREAM
-- key: chat:live:{channel_kind}:{scope_id}, MAXLEN 200
```

**Global/trade/lfg/region kanalları `chat.message`'a yazılmaz** (log volume); yalnızca Redis rolling + hot-path WS. Compliance için opt-in ops sampling (Phase 02).

## 4. Delivery Akışı

```
Client → POST /v1/chat/send { channel_kind, channel_ref?, body, to? }
  → API middleware:
    1. Auth + character bind
    2. Mute check (auth.account.chat_mute_until)
    3. Block check (per channel)
    4. Rate limit (channel + character + account)
    5. Spam filter (substring, link filter)
  → chat.Service.Send()
    6. Persist (channel kind'a göre)
    7. Redis rolling push (public kanallar için)
    8. WS fanout (pubsub subscribe by channel key)
    9. Metrics + trust score update
  → 202 Accepted (fire-and-forget WS delivery; sender kendi eventini de görür)
```

## 5. WS Subscription Modeli

Client login'de default sub'lar:
- `region:{region_id}` (region + global + trade + lfg — aynı region key altında filter kind ile)
- `party:{party_id}` (varsa)
- `guild:{guild_id}` (varsa)
- `whisper:{character_id}` (self kanal — DM receive)

Sub yönetimi Redis pubsub `SUBSCRIBE chat:{channel_key}`.
Region değiştiğinde eski sub drop + yeni sub add (bkz. [travel](../travel/README.md)).

## 6. Rate Limit

| Kanal | Karakter/dk | Account/dk | Global/dk |
|---|---|---|---|
| global | 5 | 8 | — |
| trade | 3 | 5 | — |
| lfg | 5 | 8 | — |
| region | 10 | 20 | — |
| party | 60 | — | — |
| guild | 30 | — | — |
| whisper | 30 (unique target) | 60 total | — |
| system | server-only | — | — |

Aşım → `E_RATE_LIMIT` + 30 sn cooldown badge (client display).

## 7. Spam / Filter

- **Text length**: 1-512 char (whisper 512, global 256).
- **Substring blacklist** (küfür, RMT link keywords) config'de. Match → auto-mute 10 dk + `moderation_flag='auto_muted'`.
- **Link filter**: yalnızca whitelist domain (`game domain, wiki, twitter/x` v1). Diğer link → filtered `[link_blocked]`.
- **Duplicate spam**: aynı mesaj 3 kez ardışık → cooldown 60 sn.
- **Trust score**: yeni account -20, 7 gün ago +10, 30 gün ago +30, reported -50; global kanal min trust 0.

Config `chat.filter.yml` hot reload.

## 8. Moderation

- **Report**: `POST /v1/chat/report` `{ message_id, reason }` → `audit.chat_report` + trust score -50 kasıtsız; aggregation ops queue.
- **Mute**: admin action `auth.account.chat_mute_until` set → chat send reddedilir.
- **Ban**: chat + tüm oyun aksiyonları (bkz. [auth/account.md](../auth/account.md)).
- **Auto-mute**: substring blacklist match 3 kez / 24h → 24h mute.

## 9. Endpoint

**POST `/v1/chat/send`**
**GET `/v1/chat/history?channel_kind=guild&channel_ref=...&limit=100&cursor=`** — persistent kanallar.
**GET `/v1/chat/whisper/history?with_character_id=...&limit=100`**
**POST `/v1/chat/report`**

## 10. Metrikler
- `chat_message_sent_total{channel_kind}`
- `chat_message_delivered_total{channel_kind}`
- `chat_message_dropped_total{reason}`
- `chat_rate_limit_hit_total{channel_kind}`
- `chat_auto_mute_total`
- `chat_report_filed_total{reason}`
- `chat_ws_broadcast_latency_ms`

## 11. Cross-Reference
- [private-message.md](private-message.md), [guild-chat.md](guild-chat.md), [party-chat.md](party-chat.md), [global-chat.md](global-chat.md)
- [block.md](block.md)
- [backend/websocket.md](../backend/websocket.md), [backend/redis.md](../backend/redis.md)
- [database/audit.md](../database/audit.md)
- [auth/account.md](../auth/account.md)

## 12. Bekleyen Kararlar
- Global kanal cross-region (v1: region; Phase 03: cross-region opt-in).
- Trust score parametreleri Phase 02 tune.
- Emoji whitelist (v1 küçük set; Phase 02 genişleme).
- Rich text (color, bold) — v1 yok; Phase 02'de guild custom.
- Long-term chat log retention (compliance) — Phase 02 ops karar.
