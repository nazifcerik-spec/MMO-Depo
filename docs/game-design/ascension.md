# Ascension (Yükseliş)

## 1. Amaç
Cap sonrası dikey ilerleme sunmak. Karakteri sıfırlamadan, ancak "reset karşılığında kalıcı meta bonus" ile.

## 2. Oyuncuya Fayda
- Cap sonrası hedef.
- Kalıcı meta ilerleme (her ascension "seviye üstü seviye").
- Karakter kimliği korunur.

## 3. Gameplay Loop
- Meta ölçekli; her ascension aralığı haftalar/aylar.

## 4. Temel Kurallar
- Ascension kapı koşulu: level cap + belirli quest chain + item.
- Ascension işlemi: level 1'e döner; envanter, ekipman, guild, kozmetik korunur.
- Ascension level +1: kalıcı **Ascension Point** kazanılır.
- Ascension point → skill tree'de kalıcı bonuslar (small stat, small yield, small proc).
- Reset arası cooldown: minimum 3 gün (spam engeli).
- Diminishing returns: ilk 10 ascension güçlü; sonrası küçük ama sürekli.

## 5. Kullanıcı Akışı
```
Level cap
 ↓
Ascension quest chain
 ↓
Ascension confirm (irreversible warning)
 ↓
Level 1 respawn, ascension_level +1, ascension_points +N
 ↓
Ascension Tree UI'da point spend
```

## 6. İlişkiler
- [Progression](progression.md), [Character](character.md), [Stats](stats.md), [Skills](skills.md), [Mastery](mastery.md)

## 7. Gereken Veriler
```
ascension_state
  character_id, ascension_level, ascension_points_total, spent_points
  spent_json (tree node → points)

ascension_tree_node
  id, tier, category (combat/gather/social/utility)
  cost, effect_json
  prerequisite_ids []
```

## 8. I/O
- Input: ascension confirm; ascension point spend.
- Output: character reset event; permanent bonuses applied.

## 9. Progression
- Ascension level 0..N (soft cap 50; Phase 02 rework).
- XP curve her ascension'da hafif artar (%5 per level).

## 10. Reward
Ascension point + achievement + kozmetik title.

## 11. Failure Durumları
- Quest chain uncomplete: ascension reddedilir.
- Ascension mid-confirm interrupt: transaction rollback; ascension olmaz.

## 12. Edge Cases
- Ekipman scaling: cap seviye item'lar level 1 karakterde equip edilemez → temp bag'e; tekrar cap'e ulaşınca equip.
- Guild kick ve leadership: ascension işlemi guild etkisi yok.
- Marketplace listing: ascension esnasında listing'ler kalır.

## 13. AFK Davranışı
Ascension işlemi manuel; AFK ile ilgisi yok. Sonrasında normal AFK loop devam.

## 14. Server Authoritative Gereksinimleri
- Ascension transaction atomic + auditable.
- Ascension point spend transaction.
- Reset validation: level cap + quest + item check.

## 15. Genişletilebilirlik
- Yeni tree kategorisi eklenebilir.
- Legendary ascension (Phase 03) — sonsuz.

## 16. Örnek Senaryo
Warrior level cap 50, ascension #1: level 1, +3 ascension point → gather node yield %1 kalıcı boost node harcanır. XP eğrisi %5 arttı.

## 17. Balance Considerations
- İlk 10 ascension: her biri ~%1–3 stat/yield artışı; 10. sonrası daha küçük ama sürekli.
- Ascension'lar aktif oyuncu için 60–100 saatlik aktivite sonrası hedeflenir.
- Ascension power creep sıkı kontrol: cap sonrası oyunun challenge'ı ascension level ile hafifçe scale eder.
