# Hunting (Avcılık / Mob Farming)

## 1. Amaç
Karakterin belirli bir bölgede/mob tipinde AFK combat yapmasını modellemek. Loot, XP ve skinning kaynağı olarak kullanılır.

## 2. Oyuncuya Fayda
- AFK combat verimi.
- Materyal + XP + rare drop çeşitliliği.
- Class rotation'ın "gerçek koşullarda" işlemesini test etme.

## 3. Gameplay Loop
- Aktivite: hunt spot seç → mob type seç → aktif olarak veya AFK'de dövüş.
- Session'da manuel hunt daha yüksek verim ancak zaman maliyeti.
- Meta'da rare mob spawn timer'ları, boss avı, world quest.

## 4. Temel Kurallar
- Hunt spot: `location_id`'a bağlı, karakter oraya seyahat etmiş olmalı.
- Mob tipi seçimi: birden fazla mob tipini "hunt list"e ekleyebilir.
- Danger tier: kolay/orta/zor/elite → önerilen level bandı vardır; underlevel play sıra ile yield azalır.
- Kill rate = f(character DPS, mob HP, respawn timer, party count).
- Loot: mob loot_table_id ile roll edilir.
- **Rest / Recovery**: HP/mana regen aktivite loop'una dahil; potion tüketimi opsiyonel.

## 5. Kullanıcı Akışı
```
Travel → Hunt Spot
 ↓
Hunt Panel → mob type list → limit set
 ↓
Start → Activity Engine
 ↓ (offline)
Reveal
```

## 6. İlişkiler
- [Combat](combat.md), [World](world.md), [Map](map.md), [Skills](skills.md), [Items](items.md), [Loot System](../systems/loot-system.md), [Activity Engine](../systems/activity-engine.md), [Pets](pets.md)

## 7. Gereken Veriler
```
hunt_spot
  id, location_id, difficulty_tier
  mob_pool_ids [] (fk npc_template)
  respawn_timer_ms

hunt_activity_state
  character_id, spot_id, target_mob_type_ids []
  autoloot (bool), use_potions (bool)
  stop_conditions (json: hp_below, inventory_full, mob_kill_target)
```

## 8. I/O
- Input: start hunt, adjust filters, stop.
- Output: kill log, loot yield, XP yield, activity result.

## 9. Progression
- Hunter mastery track (Phase 02): unique mobs kill sayısı, bestiary unlock.
- Rare mob spawn → server anonsu (opsiyonel), kill credit.

## 10. Reward
- XP per kill.
- Loot drops (common..epic).
- Skinning materyali (skinning profession aktif ise).
- Rare drop shot at legendary (çok düşük şans).

## 11. Failure Durumları
- Death: HP 0 → respawn + XP debt + activity paused.
- Envanter dolu → activity paused.
- Party disband ederse: solo devam eder (varsayılan).

## 12. Edge Cases
- Overkill mob (aşırı level altı): yield tamamen kesilir (%0), oyuncu warning görür.
- Undertier (overleveled) hunt: yield %70'e iner (challenge minimum eşiği).
- Mob spawn contention: personal instance kuralı (Phase 01 karar).

## 13. AFK Davranışı
Hunting flagship AFK aktivitedir. Rotation + item + stat + spec + racial all deterministik simüle edilir.

## 14. Server Authoritative Gereksinimleri
- Kill event server; loot roll server.
- Anti-farm exploit: aynı mob kısa sürede birden fazla kez farklı hesaplarca solo farm ediliyorsa flag.

## 15. Genişletilebilirlik
- "Bounty hunt" (haftalık özel target) sonra eklenebilir.
- Coop hunt (2–5 kişi) party sistemine köprü.
- Faction-hostile mob'lar (Phase 02).

## 16. Örnek Senaryo
Ranger, Whispering Forest Wolf spot'unda AFK; 6 saat, 620 wolf kill, 82 leather, 3 rare wolf fang, 14000 XP.

## 17. Balance Considerations
- Mob HP scaling level ile lineer, damage biraz sub-lineer.
- Solo AFK hunt yield vs manuel hunt yield ~%80.
- Party hunt yield: kişi başına düşen loot düşer ama tempo artar → net kişi başı yield ~%85–95 solo'ya benzer; social bonus per kill +%5.
