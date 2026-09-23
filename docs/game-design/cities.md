# Cities (Şehirler)

## 1. Amaç
Sosyal ve ekonomik yoğunlaşma noktaları. Her şehir bir "hub" olduğu kadar bir "kimlik" taşır.

## 2. Oyuncuya Fayda
- Vendor, bank, marketplace, quest hub erişimi.
- Guild hall (Phase 02) merkezi.
- Sosyal buluşma (chat, LFG board).

## 3. Gameplay Loop
- Session'da sık dönülür (bank, market, quest).
- Meta'da farklı şehirlerin farklı uzmanlıkları belirginleşir.

## 4. Şehir Türleri
- **Capital City** (kıta başkent) — tam servis.
- **Regional City** — çoğu servis.
- **Village** — bazı servis + yerel kaynaklar.
- **Outpost** — spesifik servis (crafting station, quest giver).

## 5. Temel Kurallar
- Şehirler farklı **uzmanlıklara** sahip: bir şehir smith odaklı olabilir, diğeri alchemy.
- Her şehirde bir **city faction** olabilir (Phase 02 rep grinding).
- Marketplace: kıta bazlı global, ancak lokal "market fee" farkı olabilir.
- Şehirler safe zone; PvP (Phase 02) burada kapalı.

## 6. Kullanıcı Akışı
```
City girişi → Notice board (event, quest, guild)
 ├─ Bank
 ├─ Marketplace
 ├─ Vendors (weapon, armor, potion, food, materials)
 ├─ Trainers (profession, class, spec)
 ├─ Quest Givers
 ├─ Guild Hall entrance
 └─ Inn (rest for +XP buff — Phase 02)
```

## 7. İlişkiler
- [World](world.md), [Map](map.md), [Economy](economy.md), [Marketplace](marketplace.md), [Quests](quests.md), [Guilds](guilds.md), [Professions](professions.md)

## 8. Gereken Veriler
```
city
  id, territory_id, name, size (capital/regional/village/outpost)
  services (json: bank, market, trainer_ids, vendor_ids, faction_id, quest_hub_id)
  ambient_theme
  spawn_point_id (respawn location)
```

## 9. I/O
- Input: enter city, interact with npc, use service.
- Output: service UIs (bank, market, vendor).

## 10. Progression
- City fame / reputation (Phase 02).
- Home city set: karakterin "ana üssü" tanımlanır; respawn burada.

## 11. Reward
Şehir kendisi ödül değil; ödüllerin ticaretinin yapıldığı yer.

## 12. Failure Durumları
- Servis npc göç etti (patch): oyuncu compensate + yönlendirme.
- Instance yoğunluğu (peak time): şehir layer'lara ayrılır (soft sharding, opsiyonel).

## 13. Edge Cases
- Home city değiştirme: cooldown ile mümkün.
- Under attack event (Phase 02): şehir geçici safe-off; guild defense.

## 14. AFK Davranışı
Şehirde AFK sadece "rest" niteliği taşır; combat aktivitesi yok. Rest inn (Phase 02) küçük XP buff birikimi.

## 15. Server Authoritative Gereksinimleri
- NPC/service state server.
- City-specific pricing modulation server tarafında.

## 16. Genişletilebilirlik
- Guild hall integration (Phase 02).
- Player housing plot (Phase 02+).
- City events (parade, market day) Phase 02.

## 17. Örnek Senaryo
Capital "Ashfield" ana şehir; blacksmith odaklı; marketplace fee 4%; guild halls burada; oyuncunun home city'si.

## 18. Balance Considerations
- Şehir uzmanlıkları oyuncuyu **birden çok şehri ziyaret etmeye teşvik etmeli**; ancak tek bir şehir "kral" olmamalı.
- Service kalitesi tier'lı: capital > regional > village > outpost.
