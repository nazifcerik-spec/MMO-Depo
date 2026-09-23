# Notification System (Bildirimler)

## 1. Amaç
Oyuncuya event'leri farklı kanallardan (in-game, push, e-posta) ulaştırmak. Öncelik ve opt-in kuralları ile spam engellemek.

## 2. Oyuncuya Fayda
- Önemli olayları kaçırmama.
- Konfor: "aktivitem bitti" push.
- Kontrol: hangi kanalı hangi olay için istiyor.

## 3. Kavramlar
- **Channel**: in-game toast, in-game inbox, push, e-mail, external webhook (Discord bot Phase 02).
- **Category**: system, activity, market, guild, event, boss, social.
- **Priority**: low, normal, high, critical.
- **Preferences**: karakter/hesap bazlı opt-in.

## 4. Kurallar
- Push notifications opt-in.
- Rate limit per user per category (spam).
- Critical (guild disband request, security alert) her zaman iletilir.
- E-mail: aylık max 8 (varsayılan).
- User preferences kalıcı (per account).

## 5. Kullanıcı Akışı
```
System triggers event
 ↓
Router applies rules per user preferences
 ↓
Deliver via channel(s)
 ↓
Persist in inbox (if in-game)
```

## 6. Data Model
```
notification
  id, character_id, category, priority
  title (i18n key), body (i18n key), params_json
  channels_delivered [], created_at, read_at

notification_pref
  character_id, category, channel, enabled (bool)
```

## 7. I/O
- Input: internal events.
- Output: delivered notification, inbox update.

## 8. Failure Durumları
- Push delivery fail: fallback to inbox.
- E-mail bounce: retry with exponential back-off; disable after N.

## 9. Edge Cases
- Offline uzun süre: notifications inbox'ta biriker; summary compact edilebilir.
- Localization missing: fallback language.

## 10. Server Authoritative Gereksinimleri
- Preferences validation server.
- Rate limit enforcement.

## 11. Genişletilebilirlik
- Discord bot integration (Phase 02).
- SMS opt-in (Phase 03).
- Community webhook (Phase 03).

## 12. Örnek
Activity ends (AFK mining stopped, inventory full):
- Channel: in-game toast + push (if opt-in) + inbox entry.
- Priority: normal.
- Payload: item counts, time active.

## 13. Balance Considerations
- Default settings tolerant: kritik olaylar hep açık; ötesi opt-in.
- Push spam en büyük risk; per-category daily limit ile korunur.
