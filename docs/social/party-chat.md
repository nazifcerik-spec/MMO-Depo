# Party Chat

## 1. Amaç
Party (2-5 üye) canlı chat. En düşük latency hedefi (< 200ms p95).

## 2. Karakteristik

- **Kanal kind**: `party`.
- **Channel ref**: party_id.
- **Persist**: 30 gün `chat.message` (personal scroll); Redis rolling 100.
- **Rate**: 60/dk karakter (yoğun coordination).
- **Access**: sadece o party'nin aktif üyesi.
- **Auto-leave**: party leave / disband → sub unsubscribe; scroll geçmişi kişisel kalır.

## 3. Akış

```
POST /v1/chat/send
body: { channel_kind: "party", channel_ref: "<party_id>", body }
```

Server:
1. Middleware.
2. Membership check (`party.member` tablosunda + state='active').
3. Block filter (per-viewer).
4. Persist + Redis rolling + WS fanout `party:{party_id}`.

## 4. Party Notification Messages

Party life-cycle system messages `party` kanalına system message olarak yazılır:
- "X joined the party."
- "X left the party."
- "X promoted to leader."
- "X was kicked."
- "Party disbanded."
- Combat events opt-in: "X used <skill> on <target>" (v1 kapalı; Phase 02 verbose mode).

System mesajları `chat.message.from_character_id=SYSTEM_UUID` + `moderation_flag='system'`.

## 5. Voice (Phase 03)

V1'de yok. Push-to-talk / VC integration Phase 03 tartışma.

## 6. Endpoint

**POST `/v1/chat/send`** (kind='party').
**GET `/v1/chat/history?channel_kind=party&channel_ref=<party_id>&limit=100`**.

## 7. Metrikler
- `chat_party_message_total`
- `chat_party_ws_latency_ms`

## 8. Cross-Reference
- [chat.md](chat.md), [party/party.md](../party/party.md), [party/member.md](../party/member.md)

## 9. Bekleyen Kararlar
- Voice chat Phase 03.
- Party history persist after party disband (v1: personal scroll kalır 30 gün; party ref'i null'lanır — search yavaş).
