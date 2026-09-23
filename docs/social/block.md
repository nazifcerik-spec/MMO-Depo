# Block

## 1. Amaç
Bir oyuncunun istenmeyen etkileşimleri (mesaj, whisper, invite, mail) tamamen sessizleştirmesi.

## 2. Karakteristik

- **Scope**: account-level block (block edilen account'un tüm karakterleri).
- **Cap**: 300 block karakter başına (v1).
- **One-sided**: block eden taraf hiçbir etkileşimi görmez; block edilen taraf silent (etkileşim mesajı "sent" görür ama iletilmez).
- **Etki alanı**: chat mesajları (tüm kanallar), private message, mail, friend request, party/guild/trade invite, marketplace listing gizleme (opt Phase 02).
- **Reversible**: unblock anytime; unblock sonrası ilişki resetlenmez (friend değil, block değil).

## 3. Model

```sql
CREATE TABLE social.block_list (
  blocker_account_id UUID NOT NULL,
  blocked_account_id UUID NOT NULL,
  reason VARCHAR(32) NULL,             -- 'spam'|'harassment'|'other'
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  PRIMARY KEY (blocker_account_id, blocked_account_id)
);
```

Block **account-account**; karakter-karakter değil. Karakter → account lookup her check'te var (denorm cache Redis ile).

## 4. Etki Matrisi

| Aksiyon | Block Etkisi |
|---|---|
| Global chat | Blocker filter (blocked mesaj gizli) |
| Guild chat | Aynı |
| Party chat | Aynı; blocker party join'de "blocked player detected" uyarı |
| Party invite | Silent drop (invite gitmez); target success görür |
| Private message | Silent drop |
| Mail | Silent drop (attachment varsa gönderene iade) |
| Friend request | Silent drop (bkz. friends.md) |
| Marketplace listing | v1: yok; Phase 02: opt "hide blocked seller listings" |
| Trade P2P | Zaten karşılıklı acceptli; block var ise `E_TARGET_BLOCKED` |
| Guild invite | Silent drop |

Block çift yönlü değil; **blocker → blocked** yönü tek yönlü olsa da etki iki taraflı (blocked kişi blocker'a mesaj gönderse blocker görmez).

## 5. Block ile Mail

Blocker → mail göndermek isterse `E_TARGET_BLOCKED` (kendisi block etmişse). Blocked → blocker mail göndermek isterse; mail gönderilir gibi görünür ama teslim edilmez (blocker'ın inbox'ında görünmez, gönderene iade **edilmez** — spam koruma; attachment sender'a 7 gün overflow queue).

## 6. Akış

### Block
`POST /v1/block` body: `{ target_character_name, reason? }`
- Target → account lookup.
- Zaten block → `E_ALREADY_BLOCKED`.
- Insert.
- Redis cache invalidate (`block:blocker:{account}` set).
- Aynı account içinde: eğer friend → auto-remove (friend + block çelişki).
- Aynı party/guild → block chat/invite effect ama üyelik silinmez (opsiyonel leave party UI hint).

### Unblock
`DELETE /v1/block/{target_account_id}` — clear.

## 7. Kontrol Noktası (Sunucu)

Her etkileşim öncesi middleware:
```
CheckBlock(from_account, to_account) → bool
```
- Redis cache lookup (`block:blocker:{to_account}` set contains from_account).
- Cache miss → SQL fallback + cache warm.

Cache TTL: 15 dk; block/unblock ile invalidate.

## 8. Endpoint

**GET `/v1/block`** — list.
**POST `/v1/block`** — add.
**DELETE `/v1/block/{target_account_id}`** — remove.

## 9. Rate Limit

- Block: 30/dk (abuse detection üstü ops flag).
- Read list: 30/dk.

## 10. Metrikler
- `block_added_total{reason}`
- `block_removed_total`
- `chat_message_dropped_block_total`
- `mail_dropped_block_total`

## 11. Cross-Reference
- [friends.md](friends.md), [chat.md](chat.md), [private-message.md](private-message.md), [mail.md](mail.md)
- [auth/account.md](../auth/account.md)

## 12. Bekleyen Kararlar
- Marketplace listing gizleme Phase 02.
- Aynı party'de block hint yapılırsa hangi aksiyon (leave prompt vs sadece mute) v1 seçim: sadece mute.
- Block sebep enum genişleme Phase 02 (moderation model input).
- Report akışı ile birleşme (bkz. [chat.md](chat.md) §Moderation).
