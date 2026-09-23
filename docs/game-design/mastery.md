# Mastery (Ustalık)

## 1. Amaç
Karakter level cap'ine ulaştıktan sonra yatay derinleşme sunmak. Belirli aktivite/kategoride "usta" olmak.

## 2. Oyuncuya Fayda
- Level cap sonrası anlamlı ilerleme.
- Kişisel identifier ("Ben Grandmaster Miner").
- Konfor + küçük gameplay bonusları.

## 3. Gameplay Loop
- Meta'da ilerler; günlük aktiviteler pasif olarak mastery XP verir.

## 4. Mastery Track'ları
- **Weapon Mastery** — silah türüne göre (sword, bow, staff, ...): küçük dmg bonus + unique proc.
- **Combat Discipline** — role bazlı (tanking, healing, dps): rol-spesifik bonuslar.
- **Profession Mastery** — profession bazlı: extra yield, unique recipe.
- **Region Mastery** — bir bölgede geçirilen zaman: familiar bonus (drop, gather).
- **Bestiary Mastery** — mob type kill mastery: dmg bonus vs that family.
- **Collection Mastery** — mount/pet/transmog completion (Phase 02).

## 5. Temel Kurallar
- Mastery cap: her track için 100 level; slow progression.
- Mastery point: certain rank'lerde spend edilecek passive pick (Phase 02).
- Cross-character: karakter başına ayrı (Phase 02'de account-wide bazı track'lar).

## 6. Kullanıcı Akışı
```
Mastery Panel → tabs (weapon/combat/profession/region/bestiary/collection)
 ↓
Track detay: current level, next milestone, passive
```

## 7. İlişkiler
- [Progression](progression.md), [Combat](combat.md), [Professions](professions.md), [Items](items.md), [World](world.md)

## 8. Gereken Veriler
```
mastery_track
  id, code, category, max_level, xp_curve, milestones_json

character_mastery
  character_id, track_id, level, xp
  milestones_claimed []
```

## 9. I/O
- Input: internal xp gain events.
- Output: milestone unlock notification.

## 10. Progression
Sabit XP curve; ancak diminishing returns yok — mastery uzun soluklu ama sürekli.

## 11. Reward
Milestone: passive bonus, title, cosmetic.

## 12. Failure Durumları
Mastery loss yoktur; delta yalnız pozitif.

## 13. Edge Cases
- Track rebalance: XP migrate (yeni curve).
- Retro credit yeni track eklendiğinde: son 30 günlük aktivite üzerinden estimate.

## 14. AFK Davranışı
AFK aktivitelerden mastery XP kazanılır; passive.

## 15. Server Authoritative Gereksinimleri
- XP grant server.
- Milestone unlock idempotent.

## 16. Genişletilebilirlik
- Yeni track eklemek şema stabil.
- Account-wide mastery (Phase 02).

## 17. Örnek Senaryo
Miner mastery 82 → milestone: %2 chance double node yield. Uzun bir mining life sonucu.

## 18. Balance Considerations
- Milestone bonusları %1–5 aralığı; yatay derinlik hissi vermeli, meta-changer olmamalı.
- Bir track cap 200 saat aktif engagement civarı.
- Bilinen exploit: hızlandırma botu; use-per-hour cap ile önlenir.
