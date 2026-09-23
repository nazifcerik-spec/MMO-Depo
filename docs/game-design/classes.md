# Classes (Sınıflar)

## 1. Amaç
Karakterin **oynanış kimliğini** tanımlamak: silah tercihleri, muharebe stili, ana kaynak (mana/rage/focus), rol yatkınlığı (tank/dps/support). Class, iki katmanlı Specialization ağacının kökü olarak çalışır.

## 2. Oyuncuya Fayda
- "Ben bir X'im" hissi.
- Class'a özgü skill'ler ve ekipman tercihleri.
- Specialization ile derinleşen kişiselleştirme.

## 3. Gameplay Loop'taki Yer
- Character creation'da kritik seçim.
- Session Loop'ta muharebe kararlarını şekillendirir.
- Meta Loop'ta spec'ler ve mastery ile derinleşir.

## 4. Temel Kurallar
- Class **kalıcıdır**. Değiştirmek istersen yeni karakter yaratman gerekir. (Phase 02'de ücretli class change tartışılabilir; şu an ilke: hayır.)
- **Spec** değiştirmek mümkündür (respec token / gold maliyeti ile).
- Class = base class; **her karakterin 1 tane** base class'ı vardır.
- Spec ağacı **iki katmanlıdır**: Level 5'te 2 seçenekten 1, Level 10'da 3 seçenekten 1.
- Class kodu değil, class **data'sı** yeni eklenmesine izin verir (data-driven).

## 5. Kullanıcı Akışı
```
Karakter Yaratımı
  → Class listesi (görsel + kısa açıklama + rol tag'leri)
  → Seç
  → Level 5: Spec I ekranı açılır
  → Level 10: Spec II ekranı açılır
  → Talent (opsiyonel Phase 02): pasif noktalar
```

## 6. İlişkiler
- [Character](character.md), [Races](races.md), [Specializations](specializations.md), [Skills](skills.md), [Stats](stats.md), [Combat](combat.md), [Equipment](equipment.md)

## 7. Base Class Havuzu (Öneri)

Phase 00'da 6 base class ile başlanır. Rol dağılımı: 2 tank-yatkın, 2 dps-yatkın, 1 support-yatkın, 1 hibrit.

| Kod | Ad (placeholder) | Rol yatkınlığı | Kaynak | Silah tercihleri | Not |
|---|---|---|---|---|---|
| `warrior` | Savaşçı | Tank / Melee DPS | Rage | Kılıç, balta, kalkan | Klasik ön saf |
| `ranger` | Okçu | Ranged DPS / Skirmish | Focus | Yay, tatar yayı, hafif zırh | Hunting yatkınlığı |
| `mage` | Büyücü | Ranged Caster DPS | Mana | Asa, staff, kumaş | Enchanting yatkınlığı |
| `cleric` | Rahip | Support / Off-heal | Mana | Kalkan, topuz | Party viability |
| `rogue` | Suikastçı | Melee DPS / Burst | Energy | Hançer, hafif zırh | Yüksek burst, düşük dayanıklılık |
| `druid` | Doğacı | Hybrid (dps/support) | Nature | Kısa saldırı, doğa büyüsü | Adventure/exploration bonus |

Her class'ın 2 katmanlı specialization ağacı vardır (bkz. [Specializations](specializations.md)).

## 8. Gereken Veriler
```
class
  id
  code (unique)
  display_name (i18n)
  resource_type (enum: rage/mana/energy/focus/nature)
  weapon_affinities (json)
  armor_affinity (enum: cloth/leather/mail/plate)
  role_tags (array: tank/dps/support/hybrid)
  starter_skill_ids [] (fk)
  spec_tier1_ids [] (2 spec)
  spec_tier2_ids [] (3 spec)
  base_stat_bias (json: {STR:+2, ...})
```

## 9. I/O
- Input: creation'da `class_id`.
- Output: initial skill set, resource type, stat bias, ekipman filtre önerisi.

## 10. Progression
- 1 → cap arasında sabit XP eğrisi.
- Level 5 → Spec I unlock.
- Level 10 → Spec II unlock.
- Level 20/30/40/... → passive skill point (talent, Phase 02).
- Cap sonrası: Mastery + Ascension.

## 11. Reward
Class kendi başına ödül üretmez; skill unlock ve stat bias sağlar.

## 12. Failure Durumları
- Spec ağacı verisi eksikse: unlock ertelenir, kullanıcıya bakım mesajı.
- Respec sırasında bakiye yoksa: işlem reddedilir.

## 13. Edge Cases
- Class rework: mevcut karakterlere free respec verilir ve rebalance duyurulur.
- Spec silinmesi: legacy spec olarak korunur veya en yakın yeni spec'e taşınır (oyuncu onayı ile).

## 14. AFK Davranışı
Class, AFK muharebe rotation'ının seçim havuzunu belirler. Her class için önerilen "default rotation" tanımlıdır; oyuncu isterse özelleştirir.

## 15. Server Authoritative Gereksinimleri
- Skill kullanımı server otoritesi (cooldown, kaynak).
- Spec seçimi ve respec işlemi server transaction.
- Class + Race stat toplamı server hesabı.

## 16. Genişletilebilirlik
- Yeni class = data satırı + yeni skill/spec kayıtları. Kod değişikliği yok.
- Farklı sub-role (ör. warrior→spellblade hybrid) spec katmanında modellenir; yeni base class açmak yerine spec üretmek tercih edilir.

## 17. Örnek Senaryolar
**A — Warrior/Guardian/Bulwark:** Tank yönünde iki katman; party dungeon'unda main tank.

**B — Mage/Firecaller/Cinderborn:** Alan hasarı ve DoT odaklı; hunting AFK'de yüksek verim.

**C — Rogue/Shadow/Duelist:** Burst ve tekil hedef; PvP (Phase 02) için hazırlıklı.

## 18. Balance Considerations
- Her class, tek başına AFK dungeon-lite bir muhtevayı geçebilecek kendine yeterlilikte olmalı.
- Class stat bias toplamı: net +3 stat bütçesi; class-race birleşiminde ±5 stat sapması kabul edilir.
- Aşırı kompleks resource sistemi tercih edilmez; her class'ın kaynağı **tek bir metafora** oturur (rage = alınan/verilen hasar, focus = zaman, mana = doldurulan havuz).
