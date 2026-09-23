# Core Gameplay

## 1. Amaç
Oyunun oyuncunun eline aldığında bir "günlük ritim" ve uzun vadede bir "meta ilerleme" hissi vermesini sağlamak. Session, day, week döngülerini net şekilde tanımlamak.

## 2. Oyuncuya Fayda
- Kısa oturumla anlamlı ilerleme.
- Uzak durduğunda karakterin çalışmaya devam etmesi.
- Yaptığı seçimin (aktivite, spec, sosyal katılım) sonuca yansıması.

## 3. Gameplay Loop
Üç iç içe döngü:

### 3.1 Session Loop (~5–15 dk)
1. **Login** → server AFK Engine ile offline sonucu hesaplar.
2. **Reveal** → oyuncuya "sen yokken şu oldu" özeti gösterilir (bkz. [Notification System](../systems/notification-system.md), [AFK Engine](../systems/afk-engine.md)).
3. **Housekeeping** → envanter sıkıştır, junk sat, market listing kontrol.
4. **Meaningful choice** → aktif aktiviteyi değiştir, quest kabul et, dungeon kur, guild etkileşimi.
5. **Optional active play** → kısa manuel dungeon run, boss ping, party hunt.
6. **Set & Forget** → yeni aktivite başlat, logout.

### 3.2 Daily Loop
- Günlük görevler (rotasyonlu havuzdan 3 tane).
- Dungeon reset'i (2–3 dungeon per gün önerisi).
- Market rotasyonu (24 saatlik listing süresi ile örtüşür).
- Event katılımı (varsa günlük mini event).

### 3.3 Meta Loop (haftalar–aylar)
- Class specialization tırmanışı.
- Mastery levelleri (koleksiyon, keşif).
- Ekipman upgrade (tier atlama).
- Guild progression (ortak hedefler).
- Ascension hazırlığı (bkz. [Ascension](ascension.md)).

## 4. Temel Kurallar
- Oyuncu her an **tek bir ana aktivite** çalıştırır (paralel eklentiler: pet, mount travel).
- Bir aktivite başlatıldığında server "activity contract" oluşturur; interrupt olmadıkça devam eder.
- AFK süreleri için soft cap uygulanır (bkz. [AFK Engine](../systems/afk-engine.md)).
- Sunucu, oyuncu offline iken yeni aktivite başlatamaz; en son seçilen aktivite tick eder.

## 5. Kullanıcı Akışı
```
Login
  ↓
Offline Report (dismiss)
  ↓
Character Screen
  ↓
Activity Panel  ←→  Inventory  ←→  Market  ←→  Guild
  ↓
Confirm activity → Server persists → Logout safe
```

## 6. İlişkiler
- [AFK Engine](../systems/afk-engine.md): Ana yürütücü.
- [Activity Engine](../systems/activity-engine.md): Aktivite tipi soyutlaması.
- [Progression](progression.md): Kısa/orta/uzun döngü verimliliği.
- [Character](character.md): Loop'un öznesi.

## 7. Gereken Veriler
- `character.current_activity` (aktif aktivite referansı ya da null).
- `character.last_seen_at`, `character.last_reveal_at`.
- `activity_state` snapshot (bkz. Activity Engine).

## 8. I/O
- Input: aktivite değiştirme isteği, quest kabul, market işlemi.
- Output: reveal panel, envanter deltası, XP artışı, notification.

## 9. Progression
Loop'un kendisi ilerlemez; ama tüm ilerlemenin taşıyıcısıdır. Session frekansı arttıkça oyuncu daha çok karar noktasına dokunur; ancak sistem "her gün bağlan yoksa kaybet" cezası uygulamaz.

## 10. Reward
Loop, ödül üretmez; ödülleri **transfer eden** katmandır. Ödül üretimi Activity ve Combat sistemlerine aittir.

## 11. Failure Durumları
- Login sırasında AFK Engine hesaplaması hata verirse: state donmaz; retry policy uygulanır ve son geçerli snapshot gösterilir.
- Aktivite kontratı bozulursa (patch, veri değişikliği): AFK Engine "safe fallback" ödül tablosuna düşer (bkz. AFK Engine §Fallback).

## 12. Edge Cases
- Uzun offline (30+ gün): soft cap uygulanır, veri arşivlenir, güncelleme sırasında migration çalışır.
- Aynı anda iki cihazdan login: session takedown, son cihaz kazanır.
- Server saat dilimi değişikliği: her şey UTC epoch üzerinden hesaplanır, TZ sadece sunum.

## 13. AFK Davranışı
Core gameplay AFK-first tasarlanmıştır. "Aktif oynamak" bir bonus katmandır; ana ilerleme AFK yürüyor kabul edilir.

## 14. Server Authoritative Gereksinimleri
- Aktivite değiştirme: **her zaman** server doğrular; istemci sadece tetikleyicidir.
- Reveal hesabı: server tarafında tek noktada, deterministik.
- İstemcinin gösterdiği XP/gold sadece görseldir; kaynak-of-truth server'dır.

## 15. Genişletilebilirlik
Yeni aktivite tipi eklemek core gameplay'i bozmaz; Activity Engine sözleşmesine uyar. Yeni bir session-scope özellik (örn. daily wheel) eklemek için loop'a nokta genişletmesi (extension point) tanımlanmalı.

## 16. Örnek Senaryolar
**Senaryo A — Sabah 5 dakika:** Oyuncu login olur, gece boyunca yürüyen mining aktivitesi + hunting-pet aktivitesinin sonucunu görür. Envanteri düzenler, envanter dolduğu için mining'i yerine woodcutting koyar, çıkar.

**Senaryo B — Akşam 45 dakika:** Oyuncu login olur, günlük görevleri kabul eder, party kurup dungeon'a girer, döner, market'e dungeon drop'unu koyar, gece için crafting aktivitesi başlatır.

**Senaryo C — 3 haftalık ara:** Oyuncu döner, AFK soft cap yüzünden 30 saatlik ödül alır (24 gün offline'in kayda değer kısmı), inbox'ta expire olmuş listing'lerini geri alır, kısa oryantasyon panelinden neyin değiştiğini okur.

## 17. Balance Considerations
- Bir günlük aktif oyuncu ile bir günlük "sadece login+değiştir+çık" oyuncusu arasındaki ilerleme farkı **yaklaşık 1.5x**'i geçmemeli; aktif oyunun ödülü konfor ve olay çeşitliliğidir.
- Session Loop'un tek başına 15 dakikayı geçmesi gerekmemeli; geçiyorsa UI sürtünmesi vardır demektir.
