# Quests (Görevler)

## 1. Amaç
Karakteri dünyaya, sistemlere ve narratif'e bağlayan structured hedefler. Progression'ın yönlendirici bileşeni.

## 2. Oyuncuya Fayda
- "Ne yapayım?" sorusuna cevap.
- Reward: XP, gold, item, faction rep.
- Story ile bağlanma.

## 3. Gameplay Loop
- Session'da daily quests + main story adımı.
- Meta'da long story arcs, faction quest lines.

## 4. Quest Türleri
- **Main Story** — sıralı; ana anlatı.
- **Side Quest** — bölgesel, opsiyonel.
- **Daily** — günlük havuz, günlük 3 seçim.
- **Weekly** — 5–10 haftalık büyük hedef.
- **Repeatable** — bounty benzeri; sınırlı ödül azalır.
- **Guild Quest** — guild-shared (Phase 02).
- **World Quest** — event-based (Phase 02).
- **Achievement Quest** — hidden trigger (Phase 02).

## 5. Temel Kurallar
- Quest bir veya birden çok objective adımından oluşur.
- Objective türleri: kill, collect, deliver, escort, discover, craft, use skill, dialog.
- Accept-required; ancak bazı quest'ler otomatik trigger (progress-based).
- Abandonment: mümkün ancak repeatable ise cooldown.
- Multi-character sharing (Phase 02): guild quest.

## 6. Kullanıcı Akışı
```
Quest Panel → tab (main/side/daily/weekly)
 ↓
Accept → tracker overlay
 ↓
Complete objectives → turn-in (NPC / auto)
 ↓
Reward window
```

## 7. İlişkiler
- [Reward System](../systems/reward-system.md), [World](world.md), [Combat](combat.md), [Professions](professions.md), [Guilds](guilds.md), [Notification System](../systems/notification-system.md)

## 8. Gereken Veriler
```
quest_template
  id, code, category, type
  giver_npc_id, region_id
  prerequisites (level, quest, item, rep)
  objectives [] {type, target, count}
  rewards (json: xp, gold, items, rep)
  time_limit (nullable)
  repeatable (bool), cooldown_ms

quest_progress
  character_id, quest_id, state (accepted/complete/turned_in/abandoned)
  objective_progress (json)
  accepted_at, completed_at
```

## 9. I/O
- Input: accept, abandon, turn in.
- Output: progress tick, reward payout.

## 10. Progression
- Quest completion → xp, unlock chain, faction rep.
- Achievement quest ile kozmetik unlock.

## 11. Reward
- XP + gold + item + occasional set piece.
- Long quest ile legendary progression triggers (Phase 02).

## 12. Failure Durumları
- Time-limited fail: expire; item quest ise item kaybı yok (soft fail).
- Abandonment: repeatable cooldown.

## 13. Edge Cases
- Chain interrupted (item lost): re-obtain path; NPC yardımı.
- Group turn-in: her karakter kendi reward'ını alır.

## 14. AFK Davranışı
Bazı quest objective'leri AFK aktivitelerle uyumlu (kill X mob → hunting AFK sırasında tick'lenir; gather X ore → mining aktivitesi).

## 15. Server Authoritative Gereksinimleri
- Progress tracking server.
- Reward payout transaction.

## 16. Genişletilebilirlik
- Branching quest (Phase 02).
- Player-authored notes (Phase 03).
- Faction rep grinding chains (Phase 02).

## 17. Örnek Senaryo
Daily: "Ash Wolf Pelts x10" → hunting AFK sırasında tick, login'de complete pop up, turn-in şehirde.

## 18. Balance Considerations
- Daily reward toplamı bir günün "ana ilerleme"sinin %20–35'i olmalı; günlükleri kaçırmak "yıkıcı" olmamalı.
- Ana story quest reward mid-tier item verir; bir kez oynanır kalıcı ilerlemeyi güçlendirir.
- Side quest yoğunluğu bölge başına 15–25 quest hedeflenir.
