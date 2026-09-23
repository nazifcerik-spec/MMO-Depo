# Party (Parti)

## 1. Amaç
Küçük ölçekli (2–5 kişi) geçici sosyal birim. Dungeon, hunting, adventure gibi aktivitelerin coop varyantı.

## 2. Oyuncuya Fayda
- Sosyal etkileşim, coop tatmini.
- Zorluk paylaşımı; rol dağıtımı (tank/dps/heal).
- Ödül çoğalması (share değil, çoğaltma).

## 3. Gameplay Loop
- Session'da party kur → aktivite yürüt → dağıl.
- Guild party ile daha kalıcı.

## 4. Temel Kurallar
- Party size: 2–5.
- Party leader: aktivite başlatır, kick yetkisi.
- Loot rule seçenekleri: FFA / Round Robin / Need Before Greed / Master Loot (guild raid Phase 02).
- Vote-kick (2/3 çoğunluk).
- Cross-region party mümkün ancak aktivite hepsinin bulunduğu bölgede olmalı.

## 5. Kullanıcı Akışı
```
Party Panel → create / invite / LFG
 ↓
Members ready check → start activity
 ↓
Aktivite tamam → party persist (dissolve isteğe bağlı)
```

## 6. İlişkiler
- [Combat](combat.md), [Dungeons](dungeons.md), [Hunting](hunting.md), [Adventure](adventure.md), [Social](social.md), [Guilds](guilds.md)

## 7. Gereken Veriler
```
party
  id, leader_char_id, created_at
  members [] { character_id, role, joined_at }
  loot_rule (enum)
  activity_id (nullable)
```

## 8. I/O
- Input: create, invite, kick, promote, disband, loot rule.
- Output: party UI updates, chat channel, activity events.

## 9. Progression
Party kendisi ilerlemez; ancak party içinde yapılan aktivite karakter progression'a katkı verir.

## 10. Reward
- Kill loot: rule'a göre.
- Party bonus XP (+5–10%).
- Rare drop shared roll (dungeon boss).

## 11. Failure Durumları
- Member disconnect: 60 sn grace + auto-reconnect; sonra idle drop, party solo devam.
- Party wipe: dungeon reset kuralı.
- Toxic player: report + vote-kick.

## 12. Edge Cases
- Cross-timezone party: aktivite başlatma tolerans.
- Level mismatch: dungeon otomatik level scaling (Phase 02 opsiyonel).

## 13. AFK Davranışı
Party idle-run dungeon aktivitesi AFK. Solo AFK hunt party ile paralel çalıştırılamaz (tek aktivite kuralı).

## 14. Server Authoritative Gereksinimleri
- Party state server; leader migration otomatik.
- Loot roll server + audit.

## 15. Genişletilebilirlik
- Raid (10+ oyuncu) alt sistem.
- Cross-guild party (default).
- Persistent friend party (Phase 02): saved compositions.

## 16. Örnek Senaryo
3 kişi guild party, dungeon farm, need before greed, 12 dk, 3 rare drop rolled.

## 17. Balance Considerations
- Party bonus %10 XP, %5 loot chance; solo oyunla karşılaştırılabilir seviyede.
- Loot rule default "Need Before Greed" — solo oyuncular için adaletli.
- Cross-continent party mümkün ancak aktivite lokasyona bağlı olduğundan doğal engel.
