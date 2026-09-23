# Adventure (Macera)

## 1. Amaç
"Hikaye + keşif" tabanlı, AFK-uyumlu aktivite kategorisini modellemek. Hunting kadar combat-focus, mining kadar gather-focus değildir; hikaye/olay dizisi ile ilerler.

## 2. Oyuncuya Fayda
- Dünya lore'unu keşfetme.
- Mixed reward (XP + item + faction rep + kozmetik).
- Solo veya party ile deneyimlenebilir mini olay serileri.

## 3. Gameplay Loop
- Adventure seçimi (bölge, tema, tier).
- AFK zaman ilerledikçe olay olay ilerler (stage'ler).
- Login'de "hikaye özeti" reveal (bkz. [Notification System](../systems/notification-system.md)).

## 4. Temel Kurallar
- Adventure = birden fazla `stage`'den oluşan zaman-tabanlı akış.
- Her stage: kombine bir mini "check": combat check | skill check | resource check | luck check.
- Başarısız stage: alternative branch veya penalty (küçük yield loss, XP loss yok).
- Adventure süresi: 15 dk / 1 sa / 4 sa / 8 sa varyantları.
- Aynı adventure aynı gün max N kez tekrar (rest window).

## 5. Kullanıcı Akışı
```
Adventure Panel
 ├─ Bölge filtresi + tier
 ├─ Adventure seçenekleri (kart görünümü)
 ├─ Required stats (info görüntüsü)
 ├─ Party opt-in
 └─ Start
```

## 6. İlişkiler
- [World](world.md), [Combat](combat.md), [Skills](skills.md), [Loot System](../systems/loot-system.md), [Quests](quests.md), [Party](party.md), [Notification System](../systems/notification-system.md), [Activity Engine](../systems/activity-engine.md)

## 7. Gereken Veriler
```
adventure_template
  id, code, region_id, tier
  duration_bucket (short/mid/long)
  stage_ids [] (ordered)
  reward_pool_id
  faction_id (nullable)

adventure_stage
  id, adventure_id, order_index
  check_type (combat/skill/resource/luck)
  requirement_json
  success_effect_json, failure_effect_json
  narration_key (i18n)

adventure_run
  id, character_id (or party_id), adventure_template_id
  start_at, end_at, current_stage
  seed, result_json
```

## 8. I/O
- Input: start adventure, cancel.
- Output: reveal panel with narration + rewards.

## 9. Progression
- Bir bölgenin adventure listesi tamamlanınca "region mastery" +1 (Phase 02).
- Adventure difficulty tier'ları unlock: bölge unlock progression'a bağlı.

## 10. Reward
- XP + item + material + faction rep + occasional cosmetic drop.
- Story unlock (dünya bilgisi, achievement).

## 11. Failure Durumları
- Tüm stage fail: adventure fail; %50 base reward + narration.
- Party dropout: solo devam; başarı oranı düşer.
- Server crash mid-run: son deterministik snapshot'tan resume (Phase 01 tasarım).

## 12. Edge Cases
- Character level mismatch: tier gating.
- Interrupt: cancel edilebilir; %50 iade değil, o ana kadar birikmiş partial yield alınır.
- Party split: her karakter kendi rewardını alır (share değil).

## 13. AFK Davranışı
Doğuştan AFK. Login'de her stage için kısa narration + summary reveal.

## 14. Server Authoritative Gereksinimleri
- Stage rolls seed'li deterministik.
- Narration text i18n key set; client tarafında localize edilir.

## 15. Genişletilebilirlik
- Branching narrative (Phase 02): oyuncu seçimleri stage yönünü değiştirir.
- Coop-only adventures (Phase 02).
- Seasonal adventure event line (Phase 02).

## 16. Örnek Senaryo
"Ancient Grove Whispers" adventure (4 saat, tier 20): 5 stage — enter grove (skill), find idol (luck), guardian fight (combat), decipher runes (skill), retrieve relic (resource). 4/5 başarı; relic + rare herb + 8k XP.

## 17. Balance Considerations
- Adventure yield hunting/mining altında/üstünde değil; genelde XP-heavy, material-light.
- Rare cosmetic drops adventurin özgün lezzeti; %1–2 aralığı.
- Party bonusu +%10 yield, ama personal share doğrudan verilir.
