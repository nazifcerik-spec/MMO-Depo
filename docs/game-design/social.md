# Social (Sosyal Sistemler)

## 1. Amaç
Party ve Guild dışındaki sosyal etkileşim kanallarını modellemek: chat, arkadaş listesi, mail, ignore, report.

## 2. Oyuncuya Fayda
- İnsanlarla iletişim.
- Toxic karşı korunma (ignore, report).
- Kalıcı bağlantılar (friends, mail).

## 3. Gameplay Loop
- Session'da chat'e bakma; mail açma; friend online check.
- Meta'da guild + party sosyal döngü.

## 4. Kanallar
- **Say / Local** — aynı location.
- **Zone** — region-wide.
- **City** — city-scope.
- **Guild** — guild members.
- **Party** — party members.
- **Whisper / DM** — 1:1.
- **Trade Channel** — city-based, ticaret.
- **LFG Channel** — party arayışı.
- **Global** (kapalı by default; opt-in).

## 5. Temel Kurallar
- Chat sansür/filtre server; profane list moderate.
- Rate limit spam.
- Friend list max 100 (Phase 02 expand).
- Ignore list max 100 (kişiselleştirilebilir Phase 02).
- Mail: attach items/gold; expires 30 gün; cost gold (küçük).

## 6. Kullanıcı Akışı
```
Chat window → channel tab
 ↓
Friends panel → add/remove/online status
 ↓
Mailbox → inbox/outbox/attach
```

## 7. İlişkiler
- [Party](party.md), [Guilds](guilds.md), [Economy](economy.md) (mail transfer), [Notification System](../systems/notification-system.md)

## 8. Gereken Veriler
```
chat_message
  id, channel, sender_char_id, target_id (nullable), text, at

friendship
  char_a, char_b, since, note

ignore_list
  owner_char, target_char, added_at

mail
  id, sender_char_id, receiver_char_id
  subject, body, attachments_json (items+gold)
  cost_gold_paid, expires_at, read (bool)

report
  id, reporter_char, reported_char, reason, evidence_json, status
```

## 9. I/O
- Input: send message, add friend, ignore, mail.
- Output: chat feed, notifications.

## 10. Progression
Sosyal ilerleme metriği: friend count, guild role. Achievement bazlı.

## 11. Reward
Yumuşak ödüller (achievement); doğrudan power değil.

## 12. Failure Durumları
- Rate limit hit: temp mute (5–15 min).
- Report abuse (mass false report): reporter warning.

## 13. Edge Cases
- Cross-language chat (Phase 02 translation feature).
- Guild bulk mail: sınırlı.

## 14. AFK Davranışı
Mesaj/mailler biriktirilir; login'de reveal + notification.

## 15. Server Authoritative Gereksinimleri
- Chat filtre server.
- Mail transaction atomic (item transfer safe).
- Report audit log.

## 16. Genişletilebilirlik
- Voice chat opt-in (Phase 03).
- Chat channel custom (Phase 02).
- Community forum integration (Phase 03).

## 17. Örnek Senaryo
Guild chat aktif; friend online; mail'de listing satışı 4750g attach; oyuncu mail'i açar, gold ledger'a düşer.

## 18. Balance Considerations
- Chat spam pena: 3 uyarı sonra 24h mute.
- Mail cost small (10g) — spam engelleyici, casual friendly.
- Report bir "denge kalkanı"; oyun sağlıklı topluluk için önemli.
