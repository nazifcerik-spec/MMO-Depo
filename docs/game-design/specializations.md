# Specializations (Uzmanlık Dalları)

## 1. Amaç
Base class'ı iki katmanlı seçimle derinleştirmek. Aynı base class'tan farklı **oynanış yönlerini** kanıtlamak. Sistem hard-coded değil, **data-driven** olmalı ki gelecek içerik pack'leri yeni spec ekleyebilsin.

## 2. Oyuncuya Fayda
- Aynı class'ın "ben nasıl oynuyorum" hissini kişiselleştirmek.
- Party içinde net rol kimliği (main tank vs off-tank; burst dps vs sustained dps).
- Respec'le denemeye açık, taahhüde bağlı olmayan derinlik.

## 3. Gameplay Loop
- Character creation'ın hemen sonrası: Level 5 & 10 spec unlock kararları.
- Meta Loop'ta talent/mastery ile beslenir.
- Yeni spec eklendiğinde bir "spec preview event" ile oyuncuya sunulur (Phase 02).

## 4. Temel Kurallar
- **Spec I** (tier 1): Level 5'te açılır. 2 seçenekten 1.
- **Spec II** (tier 2): Level 10'da açılır. **Spec I'e bağlı olarak** 3 seçenekten 1. Yani her Spec I'in 3 alt Spec II'si var.
- Yani base class başına toplam: **2 Spec I × 3 Spec II = 6 uç dal**.
- Respec: gold + zaman cooldown; ilk respec bedava (öğrenme fırsatı).
- Yeni class eklenirken en az 2 Spec I + 6 Spec II ile birlikte gelmeli.

## 5. Kullanıcı Akışı
```
Level 5 → "Yolunu seç" ekranı → Spec I onay
Level 10 → Spec I'e bağlı 3 seçenek → Spec II onay
Herhangi bir zaman → Respec panel (maliyet gösterir)
```

## 6. İlişkiler
- [Classes](classes.md), [Skills](skills.md), [Combat](combat.md), [Stats](stats.md)

## 7. Gereken Veriler
```
specialization
  id
  code (unique)
  class_id (fk)
  tier (1 or 2)
  parent_spec_id (nullable for tier 1)
  display_name (i18n)
  passive_ids [] (fk)
  active_skill_unlock_ids [] (fk)
  stat_bias (json)
  playstyle_tags (array)
```

## 8. I/O
- Input: `select_specialization(character_id, spec_id)`.
- Output: skill unlocks, passive activations, stat bias adjustments.

## 9. Progression
- Tier 1 → Tier 2 açılması level odaklı.
- Spec içinde ilerleme talent puanı üzerinden (Phase 02) veya use-based mastery (Phase 02) olabilir; Phase 00'da spec kendisi sabit set.

## 10. Reward
Spec seçimi anında: yeni skill(ler), yeni passive(ler), stat bias.

## 11. Failure Durumları
- Spec seçimi validation'ı: parent spec I ≠ tier 2 base ise reddet.
- Respec sırasında insufficient funds: transaction rollback.

## 12. Edge Cases
- Yeni Spec II eklendiğinde eski oyunculara "seçim penceresi" verilir (yeni açılmışsa mevcut Spec II değişebilir, ücretsiz).
- Spec rebalance: rework öncesi ve sonrası mekaniği log'lanır; ciddi rework varsa free respec pencere açılır.

## 13. AFK Davranışı
Spec'in aktif skill'leri AFK rotation'ına dahil olur; passive'ler AFK muharebe hesabında uygulanır.

## 14. Server Authoritative Gereksinimleri
- Spec seçim ve değişim işlemleri transaction'lı.
- Skill unlock durumu server'da tutulur; istemci güvenilir değildir.

## 15. Genişletilebilirlik
- Tier eklenmesi (tier 3, tier 4) sonradan mümkün — Phase 02+ konusu; şema `tier` int alanını taşır.
- Cross-class spec (multi-class) şu aşamada **planlanmıyor**. Tasarımsal olarak class kimliğini bulanıklaştırır.

## 16. Örnek Senaryolar
**Warrior:**
- Spec I: `guardian` (defense) | `berserker` (offense)
- Spec II (guardian): `bulwark` (block master) | `sentinel` (crowd control tank) | `warden` (party protector)
- Spec II (berserker): `reaver` (dual wield) | `titan` (2-handed) | `bloodletter` (execute/heal)

**Mage:**
- Spec I: `elementalist` (dps) | `mystic` (utility/support)
- Spec II (elementalist): `firecaller` | `stormbringer` | `frostweaver`
- Spec II (mystic): `chronomancer` | `illusionist` | `arcanist`

## 17. Balance Considerations
- Her Spec II'nin en az bir "signature skill"i olmalı.
- Aynı Spec I'in altındaki 3 Spec II arasında rol farkı olmalı (aynı build değişkeni olmamalı).
- Passive stat bias toplamı Spec I için +3, Spec II için +2 bütçesi.
