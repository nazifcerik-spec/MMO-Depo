# Social — Katalog

> **Faz 28**: Oyuncular arası iletişim ve ilişki katmanı — friends, block, chat kanalları, mail.

## 1. Kapsam
Bu klasör, oyuncu-oyuncu iletişimini sunucu-otoriter, moderasyona hazır ve rate-limited biçimde ele alır. AFK MMO olduğu için **asenkron iletişim (mail) online sync (chat) kadar önemli**.

## 2. Ana Bileşenler

- **Friends / Block** — presence + interaction gating (block chat + trade + party invite).
- **Chat kanalları** — global, region, party, guild, private message; ortak `chat.message` domain'i.
- **Mail** — offline mesaj + gold/item attachment + escrow + moderation.
- **Rate limit + spam filter** — server-side, message trust score.
- **Moderation** — report akışı, mute, ban integration (bkz. [auth/account.md](../auth/account.md)).
- **Persistence** — chat mesajları kısa TTL (kişisel scroll için 30 gün), mail 30 gün + 7 gün overflow.

## 3. Dokümanlar
- [friends.md](friends.md) — Friends list + presence + invite/accept.
- [block.md](block.md) — Block list + interaction gating.
- [chat.md](chat.md) — Chat engine + channel taxonomy + delivery (WS) + rate limit.
- [private-message.md](private-message.md) — Whisper/DM 1-to-1 + offline queue.
- [guild-chat.md](guild-chat.md) — Guild + guild officer channels.
- [party-chat.md](party-chat.md) — Party (2-5) canlı chat.
- [global-chat.md](global-chat.md) — World/trade/lfg channels + throttle + region split.
- [mail.md](mail.md) — Mail + attachment + fee + expiration.

## 4. Cross-Reference
- [auth/account.md](../auth/account.md) — Mute / ban state
- [character/character.md](../character/character.md) — Display name
- [guild/guild.md](../guild/guild.md), [party/party.md](../party/party.md)
- [economy/fee.md](../economy/fee.md) — Mail postage
- [inventory/item-transfer.md](../inventory/item-transfer.md) — Mail attach lock
- [systems/notification.md](../systems/notification.md) — Mail/private message push

## 5. Prensipler

- **Server-authoritative**: her mesaj server relay; client-to-client direkt yok.
- **Persist for scroll** yalnızca kişisel scroll (30 gün TTL kişisel inbox); global public kanallar **persist edilmez** (yalnızca live WS + last 200 in Redis).
- **Rate limit** kanal başına + karakter/account başına.
- **Anti-spam**: substring blacklist + link filter (whitelist domain) + trust score.
- **Block gate**: block edilen taraftan gelen chat/mail/invite tamamen düşer (silent drop).
- **Moderation trail**: raporlama ve mute audit (bkz. [database/audit.md](../database/audit.md)).

## 6. Metrikler
- `chat_message_total{channel_kind}`
- `chat_message_dropped_total{reason}` — rate limit / block / mute
- `mail_sent_total`, `mail_delivered_total`, `mail_expired_total`
- `friends_count_avg`, `block_count_avg`
- `chat_report_filed_total{reason}`

## 7. Bekleyen Kararlar
- Cross-region global chat (v1: region-scoped; Phase 03 global).
- Voice chat (Phase 03+ opt).
- Rich text / emoji policy (v1: text only + limited emoji set).
- Chat log retention politikası (compliance).
