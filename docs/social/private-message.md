# Private Message (Whisper / DM)

## 1. Amaç
İki karakter arasında 1-to-1 anlık mesaj. Offline hedefe **queue** ile teslimat (online olunca push).

## 2. Karakteristik

- **Kanal kind**: `whisper`.
- **Persist**: 30 gün `chat.message` (personal scroll).
- **Offline queue**: hedef offline ise mesaj yine `chat.message`'a yazılır; hedef online olduğunda WS "unread whispers" event.
- **Kanal ref**: yok (from/to pair yeterli).
- **Length**: 1-512 char.
- **Rate**: 30 mesaj/dk unique target başına; account başına 60/dk toplam.

## 3. Akış

```
POST /v1/chat/whisper
body: { to_character_name, body }
```

Server:
1. Auth + rate limit + mute + spam filter.
2. Target character resolve (name → id + account_id).
3. Block check: `block.CheckBlock(from_account, to_account)` → block var **kendi tarafında** ise `E_TARGET_BLOCKED`; block **karşı tarafta** ise silent drop (mesaj `chat.message`'a yazılır ancak `moderation_flag='blocked_target'`, teslim edilmez).
4. Insert `chat.message` (channel_kind='whisper', from/to pair).
5. Target online ise WS push (`whisper:{to_character_id}`).
6. Metric.

## 4. Delivery Semantiği

- **At-least-once WS**: WS bağlantı stabil ise anlık; drop olursa Redis pending queue (5 dk retry).
- **At-least-once "unread indicator"**: WS ulaşmazsa target next login → `chat.message` unread count → in-game badge + optional push notification.

Unread hesaplama:
```
GET /v1/chat/whisper/unread
  → SELECT COUNT(*) FROM chat.message
     WHERE to_character_id = :self AND channel_kind='whisper'
       AND read_at IS NULL AND sent_at > :cutoff (30d);
```

`read_at` mark: client açtığında `POST /v1/chat/whisper/read` `{ from_character_id, up_to_sent_at }`.

```sql
-- chat.message.read_at kolonu whisper için opt:
ALTER TABLE chat.message ADD COLUMN read_at TIMESTAMPTZ NULL;
CREATE INDEX idx_chat_whisper_unread ON chat.message (to_character_id, channel_kind, read_at) WHERE channel_kind='whisper';
```

## 5. Presence-Aware UX (Client-Side)

- Target offline ise UI'da "**offline — teslim edildi, mesajı bir sonraki girişinde görecek**".
- Target dnd ise: "**do not disturb**", client uyarı (mesaj yine iletilir).
- Target no_one presence: state her zaman "offline" görünür (bkz. [friends.md](friends.md)).

## 6. Endpoint

**POST `/v1/chat/whisper`** — send.
**GET `/v1/chat/whisper/history?with_character_id=...&limit=100&cursor=`** — konuşma scroll.
**GET `/v1/chat/whisper/unread`** — badge sayaç.
**POST `/v1/chat/whisper/read`** — mark read.

## 7. Metrikler
- `whisper_sent_total`, `whisper_delivered_online_total`, `whisper_queued_offline_total`
- `whisper_ws_delivery_latency_ms`
- `whisper_unread_gauge_avg`

## 8. Cross-Reference
- [chat.md](chat.md), [friends.md](friends.md), [block.md](block.md)
- [systems/notification.md](../systems/notification.md) — push notify (mobile)

## 9. Bekleyen Kararlar
- Group whisper (3+ ad-hoc DM) Phase 02.
- Voice whisper Phase 03.
- Read receipt (opsiyonel) Phase 02 privacy setting.
