# Travel (Seyahat)

## 1. Amaç
Konumlar arası geçişi zamana ve maliyete bağlamak. Zamanı, oyunun bir sistem parçası olarak kullanmak; AFK travel'ı bir aktivite tipi yapmak.

## 2. Oyuncuya Fayda
- Yolculuk sırasında karakter uyumaz; aktivite bir "yol" olarak devam eder.
- Farklı travel metodları farklı ödünleşimler getirir (fast paid vs slow free).

## 3. Gameplay Loop
- Session'da: "hunt spot'a git" → travel start → AFK'de tamamlanır.
- Meta'da: mount upgrade, teleport unlock, cross-continent gemi.

## 4. Travel Türleri
- **Walking** — free, en yavaş.
- **Mount** — mount tipi ve tier'a bağlı hız.
- **Road** — road üstünde bonus hız.
- **Wagon / Caravan Line** — belirli city pairs arasında; ucuz, orta hız, party-friendly.
- **Ship** — cross-continent; program saatleri (in-game).
- **Teleport / Waygate** — gold + item cost; instant veya kısa timer; kullanım cooldown.

## 5. Temel Kurallar
- Travel bir aktivitedir (Activity Engine kontratı).
- Interrupt: karakteri travel'dan çıkarır; hedefe varmamış sayılır (partial time compensation opsiyonel).
- Combat interrupt: yol üzerinde ambush (Phase 02) → travel pause + combat.
- Route hesabı server; tehlike bölgelerinden geçiş longer path opsiyonu.
- Environment: yağmur/kar hız modifieri (Phase 02).

## 6. Kullanıcı Akışı
```
Destination select
 ↓
Route preview (süre + maliyet + tehlike)
 ↓
Start → AFK
 ↓
Arrive → notification / auto-next activity (opsiyonel)
```

## 7. İlişkiler
- [World](world.md), [Map](map.md), [Cities](cities.md), [Mounts](mounts.md), [Activity Engine](../systems/activity-engine.md), [AFK Engine](../systems/afk-engine.md), [Economy](economy.md)

## 8. Gereken Veriler
```
travel_activity
  character_id, start_location, end_location, method, path_json
  start_at, expected_end_at, seed
  cost_paid (gold + items)

teleport_binding
  character_id, waypoint_id
```

## 9. I/O
- Input: request travel(method, destination), interrupt.
- Output: activity ticker; arrival event.

## 10. Progression
- Mount tier progression.
- Teleport unlock: quest, level, item.
- Faction rep unlock ile bazı waygate açılır (Phase 02).

## 11. Reward
Travel doğrudan ödül vermez; ancak "auto-next activity" ile arrival'da yeni aktivite başlayabilir (workflow).

## 12. Failure Durumları
- Insufficient funds: travel reject.
- Blocked route (patch/event): reroute veya iptal.

## 13. Edge Cases
- Alternate route hesabı: fastest vs safest.
- Cross-continent (ship): scheduled → oyuncu vardığında ship kalkmışsa waitlist.

## 14. AFK Davranışı
Travel doğuştan AFK. Server ticker vardığında notification.

## 15. Server Authoritative Gereksinimleri
- Route + duration hesabı server.
- Cost transaction server.

## 16. Genişletilebilirlik
- Ambush events (Phase 02).
- Guild caravan (Phase 02).
- Weather / seasons (Phase 02+).

## 17. Örnek Senaryo
Ashfield → Emerald Coast, mount tier 2, 24 dk; oyuncu logout, arrival'da auto-fishing activity başlar.

## 18. Balance Considerations
- Teleport her zaman gold sink; ekonomi için önemli.
- Walking ile mount arasında fark hissedilir olmalı (2×+).
- Cross-continent süre 30 dk üzerinde olmalı; anlık teleport ile MMO ölçek hissi kaybolur.
