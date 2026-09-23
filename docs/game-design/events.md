# Events (Etkinlikler)

## 1. Amaç
Sunucunun günlük ritmine ve seasonal takvimine event injection ile canlılık katmak. Community fokusu yaratmak.

## 2. Oyuncuya Fayda
- Sıra dışı aktivite, unique reward.
- Sunucu paylaşımı: aynı anda binlerce oyuncu bir olaya odaklanır.
- Kısa süreli hedef (FOMO değil, opt-in).

## 3. Gameplay Loop
- Session'da aktif event'e katıl.
- Meta'da seasonal event serileri; kozmetik unlock.

## 4. Event Türleri
- **Daily World Event** — 1–4 saatlik pencere; küçük ödül.
- **Weekly Event** — 24–48h pencere; kayda değer reward.
- **Seasonal Event** — 2–4 hafta; kozmetik + achievement.
- **Emergent World Event** — dünyada tetiklenir (world boss respawn cluster).
- **Guild Event** — guild-internal (Phase 02).

## 5. Temel Kurallar
- Event katılımı opsiyonel; kaçırmak "hayat mahvı" hissi vermez.
- Kozmetik reward event-özel; ancak next season aynı kategori olduğu için "arşiv"lenir, ileride re-run.
- Event tokens marketable değil (soulbound); event vendor'unda kozmetik satın alınır.

## 6. Kullanıcı Akışı
```
Event Board (city) veya Notification
 ↓
Detay okuyup katıl (opt-in)
 ↓
Event aktivite → tick, ödül
 ↓
Event vendor → token karşılığı ödül
```

## 7. İlişkiler
- [Notification System](../systems/notification-system.md), [World Bosses](world-bosses.md), [Reward System](../systems/reward-system.md), [Achievements](achievements.md), [Cities](cities.md)

## 8. Gereken Veriler
```
event_template
  id, code, name, category (daily/weekly/seasonal/emergent/guild)
  start_at, end_at (nullable for emergent)
  activity_hook_id
  reward_pool_id, token_currency_id
  narration_key (i18n)

event_participation
  character_id, event_id, contribution_score, rewards_claimed
```

## 9. I/O
- Input: join, contribute (activity), claim reward.
- Output: event progress, leaderboard (opsiyonel), reward payout.

## 10. Progression
- Season-long achievement (multiple event completion).
- Kozmetik unlock chain.

## 11. Reward
- Event token → cosmetic vendor.
- Rare event drops (mount, pet, transmog Phase 02).
- Small XP + gold.

## 12. Failure Durumları
- Missed window: soft-fail; makeup ödül yok, ama season sonundaki achievement kısmi kalır.
- Server load: soft-shard olası; opsiyonel queue.

## 13. Edge Cases
- Overlap events: iki eventte aynı anda contrib mümkün, ancak reward şeffaf.
- Cross-timezone: event ⏱ saatleri bir sonraki season için TZ friendly bir rotasyon.

## 14. AFK Davranışı
Bazı event'ler AFK aktivitelerle uyumlu (ör. haftalık "mining marathon" AFK mining tick'lerini sayar). Diğerleri kısa manuel etkileşim gerektirir.

## 15. Server Authoritative Gereksinimleri
- Event participation server ledger.
- Reward payout single-fire (idempotent).

## 16. Genişletilebilirlik
- Community meta-event (Phase 02): tüm sunucu ilerleme çubuğunu doldurur.
- Guild vs guild event (Phase 02+).
- Player-created event framework (Phase 03).

## 17. Örnek Senaryo
Weekly "Harvest Rush": profession XP +%20 buff aktif; %25 daha fazla materyal drop. Katılım pasif; harvest token collect vendor'da cosmetic hat.

## 18. Balance Considerations
- Event ödülleri **power değil**, konfor/kozmetik olmalı.
- Season'lar arasında breather zaman: her season sonrası 1 hafta boşluk.
- Emergent event'lerin tetikleyicileri random-ish ancak coğrafi rotasyon garantili.
