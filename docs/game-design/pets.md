# Pets (Evcil Hayvanlar)

## 1. Amaç
Karakterin yanında çalışan pasif/aktif companion sistemi. Combat, gathering ve travel'da yardımcı; koleksiyon ve kimlik boyutu.

## 2. Oyuncuya Fayda
- Sürekli companion (görsel + mekanik).
- Combat kısmı için minor dps/support katkı.
- Gathering aktivitelerinde bonus yield.
- Koleksiyon boyutu (Phase 02 pet compendium).

## 3. Gameplay Loop
- Pet aktive et → summon.
- Pet aktivitesi seç (combat auto / gather assist / idle).
- Pet leveling passive; kullanıldıkça büyür.

## 4. Pet Türleri
- **Combat Pets** — hunter/warlock-like; damage/tank/heal minor.
- **Gathering Pets** — mining/foraging/fishing bonus.
- **Utility Pets** — inventory extend (bag pet), speed pet, luck pet.
- **Cosmetic Pets** — sadece görsel (Phase 02 shop).

## 5. Temel Kurallar
- Aktif pet: 1 combat pet + 1 utility pet (Phase 00 baseline).
- Pet level cap = karakter cap (soft link).
- Pet ölmez; combat'ta downed → kısa cooldown sonra revive.
- Pet food (consumable) buff verir; opsiyonel.
- Pet başka oyuncuya trade edilemez (bound to character).

## 6. Kullanıcı Akışı
```
Pet Panel → owned pets grid
 ├─ Activate (set as active)
 ├─ Rename
 ├─ Feed (buff)
 └─ Retire (stable, not delete)
```

## 7. İlişkiler
- [Combat](combat.md), [Activity Engine](../systems/activity-engine.md), [Professions](professions.md), [Items](items.md), [AFK Engine](../systems/afk-engine.md)

## 8. Gereken Veriler
```
pet_template
  id, code, species, type (combat/gather/utility/cosmetic)
  base_stats (json)
  active_ability_id (nullable)
  passive_ability_id (nullable)
  rarity

pet_instance
  id, owner_char_id, template_id
  level, xp, name (character-set)
  active (bool)
  bond_level (Phase 02: kullanım süresi ile artar)
```

## 9. I/O
- Input: summon, dismiss, feed, retire.
- Output: pet stat contrib, activity bonus.

## 10. Progression
- Pet leveling: karakter aktivitesine bağlı XP paylaşımı.
- Bond level (Phase 02): pet-owner arasında bonus unlock.

## 11. Reward
- Combat contribution.
- Gathering yield bonus.
- Special rare pet drop (world boss, dungeon, achievement reward).

## 12. Failure Durumları
- Pet downed: kısa süre revive cooldown; combat continues without.
- Envanter dolu (gather pet): pet contribution paused.

## 13. Edge Cases
- Pet type mismatch (gather pet aktif ama karakter combat aktivitesinde): pet passive stat sağlar, gather bonusu inaktif.
- Legendary pet loss risk yok (bound to char).

## 14. AFK Davranışı
Pet, karakter AFK aktivitesinin passive assist katmanı. Deterministik seed dahilinde.

## 15. Server Authoritative Gereksinimleri
- Pet state server.
- Pet activation transaction.

## 16. Genişletilebilirlik
- Pet compendium (Phase 02).
- Pet arena / battle mini-game (Phase 02+).
- Cross-account BoA pets (Phase 03).

## 17. Örnek Senaryo
Ranger'ın Wolf Cub'ı active; hunting spot'ta AFK. Pet damage %8 boost + occasional flank crit; 8 saatte %5 ekstra loot.

## 18. Balance Considerations
- Pet damage contribution %5–15 aralığı; pet olmadan da oyun oynanabilir olmalı.
- Gathering pet bonus %5–10; profession'ın kendi yield'ının üstünde marjinal.
- Legendary pet stat farkı %5 ile sınırlı; rarity daha çok görsel ve unique effect.
