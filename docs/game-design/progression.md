# Progression (İlerleme Meta-Sistemi)

## 1. Amaç
Karakter ilerlemesini tanımlayan farklı eksenlerin **birlikte** nasıl işlediğini tarif etmek. Tek bir eksen değildir: level, profession, mastery, ascension paralel ve etkileşimli.

## 2. Oyuncuya Fayda
- Farklı oyun stillerinin ödüllendirilmesi.
- Sürekli anlamlı ilerleme (bir eksen dolduysa diğerinde ilerleme var).
- "Cap'i gördüm, ne yapayım?" sorusuna cevap.

## 3. Gameplay Loop
- Session'da hangi eksene katkı yaptığı sayacı görür.
- Meta'da eksen dengesizliği (hepsi profession, hiç combat) UI ipucu sunar.

## 4. Progression Eksenleri
- **Character Level** — birincil; combat gate.
- **Profession Level** — 3 slot; ekonomi ve craft.
- **Mastery** — track bazlı yatay derinlik.
- **Ascension Level** — cap sonrası dikey (bkz. [Ascension](ascension.md)).
- **Reputation / Faction** — Phase 02 (unlock, discount).
- **Guild Contribution** — sosyal ilerleme.
- **Collection** — mount/pet/transmog (Phase 02).

## 5. Temel Kurallar
- XP tipleri ayrı ledger'da tutulur; conversion yok (character xp ≠ profession xp).
- Level up tetikleyici: xp threshold; anında transaction.
- Eksenler arası "reward compensation": bir eksende ilerleyemeyen bir başka ekseni doldurur; balance.

## 6. Kullanıcı Akışı
```
Character sheet → progression tab
  Level bar / Profession bars / Mastery highlights / Ascension summary
```

## 7. İlişkiler
- [Character](character.md), [Stats](stats.md), [Classes](classes.md), [Skills](skills.md), [Professions](professions.md), [Mastery](mastery.md), [Ascension](ascension.md), [Reward System](../systems/reward-system.md), [Activity Engine](../systems/activity-engine.md)

## 8. Gereken Veriler
```
character.level, character.xp
character_profession.level, xp (per profession)
character_mastery (per track)
character_ascension
character.faction_reputations (Phase 02)
```

## 9. I/O
- Input: xp grant events.
- Output: level-up event, milestone unlock.

## 10. Progression Curve
- Character level: geometric curve; cap seviye 50 → 100 hedef (Phase 02).
- Profession: pace slower than character; grandmaster ~200 saat aktif kullanım.
- Mastery: geniş bant; her track uzun.
- Ascension: hafif exponential; ilk ascension ~60 saat aktif hedef.

## 11. Reward
Level ups: stat point + skill point + notification + rarely item.
Milestones: achievement + cosmetic.

## 12. Failure Durumları
- XP debt (ölüm) → temp negative delta; XP loss cap ile sınırlı.
- Level max iken xp: ascension pipeline'a redirect gösterilir (UI).

## 13. Edge Cases
- Uzun offline: XP tavan/floor ile hesap tam anlamlı olur.
- Rollback (exploit): affected characters manuel review.

## 14. AFK Davranışı
Tüm eksenler AFK aktivitelerden beslenir. AFK "sadece level" ilerletir değildir; profession + mastery + ascension eligibility hepsi ilerler.

## 15. Server Authoritative Gereksinimleri
- Tüm XP mutation transaction'lı.
- Curve tabloları server tarafında versiyonlu.

## 16. Genişletilebilirlik
- Yeni progression ekseni (ör. discovery) eklenebilir.
- Season pass tarzı geçici eksen (Phase 02) — kalıcı karakter değişmez.

## 17. Örnek Senaryo
Karakter cap 50, ascension 3, mining tier 62 (Expert), region mastery 20; günlük 4 saatlik AFK mining + 1 saatlik hunting; 2 hafta sonra ascension #4 için hazır.

## 18. Balance Considerations
- Eksenlerin XP yield oranı balanced: 1h AFK mining = tabii mining xp; ek olarak %5 character xp, %2 mastery xp.
- Level cap sonrası ödüllerin "boş" hissetmemesi için ascension gate açıkça yönlendirilir.
- Cross-eksen progression ile solo oyuncu kendini "her şeyi yapabilir" hisseder ancak spec/uzmanlık sinerjisi ödüllenir.
