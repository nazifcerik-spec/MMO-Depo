# Guild Chat

## 1. Amaç
Guild ve guild officer kanalları. Guild member coordination — event heads-up, help request, casual talk.

## 2. Kanallar

| Kod | Erişim | Rate | Notlar |
|---|---|---|---|
| `guild` | Tüm guild members | 30/dk karakter | Ana kanal |
| `guild_officer` | Officer + Leader | 30/dk karakter | Officer koordinasyon |
| `guild_motd` | Server → herkes (join'de) | server-only | Message of the day |

## 3. Erişim

Guild membership zorunlu:
- `social.guild_member` tablosunda karakter var mı.
- Rank permission `chat_guild_send` (default: tüm rank).
- Officer kanalı için rank permission `chat_officer_send` (default: officer + leader).

Rank permission kaynağı: [guild/permissions.md](../guild/permissions.md).

## 4. Model

`chat.message` (channel_kind='guild' / 'guild_officer', channel_ref=guild_id) — 30 gün.
Redis rolling last 200 (yeni join scroll için).

## 5. Akış

```
POST /v1/chat/send
body: { channel_kind: "guild", channel_ref: "<guild_id>", body }
```

Server:
1. Middleware (auth, rate, mute, spam).
2. Membership check (`guild_member.character_id = self && guild_id = channel_ref`).
3. Officer channel için rank permission ek check.
4. Block check kanal içi — blocker'lar hedefsiz mesaj görmez ama mesaj gönderilir; hedef bazlı per-viewer filter.
5. Persist + Redis rolling + WS fanout `guild:{guild_id}` pubsub.

## 6. Guild MOTD

- Leader/officer set eder: `PUT /v1/guild/{id}/motd` body: `{ text }`.
- 512 char max.
- Karakter guild'e join'de veya login'de MOTD system message olarak `guild` kanalına push.
- Değişiklik anında tüm guild members WS "guild_motd_changed".

MOTD persist:
```sql
ALTER TABLE social.guild ADD COLUMN motd_text VARCHAR(512) NULL,
                        ADD COLUMN motd_set_at TIMESTAMPTZ NULL,
                        ADD COLUMN motd_set_by_character_id UUID NULL;
```

Değişiklik audit trail: `audit.guild_motd_change`.

## 7. Guild Bulletin (Phase 02)

MOTD tek satır; bulletin uzun form (blog gibi). Phase 02'de opt.

## 8. Guild Member Presence Feed

Guild UI'da "who is online" listesi:
- Presence key `presence:character:{id}`.
- Guild üyeleri için subscribe.
- Login/logout WS eventi guild geneli push (opt-in per user).

## 9. Rate Limit + Anti-Spam

Chat.md tabloya bkz.
Ek kural: guild chat'te 5 dk içinde 20+ mesaj → soft mute uyarı; devam ederse officer moderation prompt.

## 10. Endpoint

**POST `/v1/chat/send`** (kind='guild' | 'guild_officer').
**GET `/v1/chat/history?channel_kind=guild&channel_ref=<guild_id>&limit=100&cursor=`**.
**PUT `/v1/guild/{id}/motd`** — leader/officer.

## 11. Metrikler
- `chat_guild_message_total{guild_id}`
- `chat_guild_officer_message_total`
- `chat_guild_motd_change_total`

## 12. Cross-Reference
- [chat.md](chat.md), [guild/guild.md](../guild/guild.md), [guild/permissions.md](../guild/permissions.md)
- [database/social.md](../database/social.md)

## 13. Bekleyen Kararlar
- Guild bulletin Phase 02.
- Guild alliance chat Phase 03.
- Custom guild channel (leader oluşturur, rank permission ata) Phase 02.
