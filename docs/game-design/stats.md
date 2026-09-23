# Stats (İstatistikler)

## 1. Amaç
Muharebe, meslek ve aktivite hesaplamalarını sürdüren sayısal iskelet. Tüm mekanik ödünleşimlerin (trade-off) para birimi.

## 2. Oyuncuya Fayda
- Ekipman, race, class, spec, buff, temporary state'in etkilerini tek yerde görmek.
- "Neyi optimize ediyorum?" sorusuna cevap verecek net metrikler.

## 3. Gameplay Loop
- Session Loop'ta character sheet ile etkileşim (item takma/çıkarma).
- Meta Loop'ta ilerlemenin metriği.

## 4. Stat Kategorileri

### 4.1 Primary Stats (birincil)
| Kod | Ad | Etki |
|---|---|---|
| `STR` | Strength | Melee damage, taşıma kapasitesi |
| `DEX` | Dexterity | Ranged damage, evasion, crit chance |
| `INT` | Intelligence | Spell damage, mana pool |
| `WIS` | Wisdom | Mana regen, healing power, meslek verimliliği |
| `STA` | Stamina | HP pool, physical resist |
| `CHA` | Charisma | Vendor fiyatları, quest reward %, party buff scaling |

### 4.2 Secondary / Derived Stats
Primary'lerden hesaplanır ama doğrudan da artırılabilir:
- **Attack Power** (fizik hasar taban)
- **Spell Power** (büyü hasar taban)
- **Crit Chance / Crit Damage**
- **Haste** (attack/cast speed)
- **Accuracy / Evasion**
- **Armor / Magic Resist**
- **Block Chance / Parry Chance**
- **HP / Mana / Resource Regen**
- **Life Steal / Spell Vamp**
- **Cooldown Reduction (CDR)** (hard cap %30)

### 4.3 Tertiary / Utility Stats
- **Movement Speed** (soft cap Phase 01'de belirlenir)
- **Carry Capacity**
- **Gathering Yield %**
- **Crafting Success %**
- **XP Gain %**
- **Gold Find %**
- **Rarity Bonus %**
- **AFK Efficiency %** (hedef: kozmetik ödül, güç değil)

## 5. Temel Kurallar
- Level up → sabit sayıda **stat point** (öneri: 5/level).
- Class + Race baz + item + buff → nihai stat.
- **Diminishing returns**: crit chance %75'ten sonra logaritmik, CDR %30 hard cap.
- Negative stat mümkün değil (0 taban).
- **Stat ekranı UI**: primary + derived + tertiary üç sekmede.

## 6. Kullanıcı Akışı
```
Character sheet açılır
  ↓
Primary sekme (STR/DEX/... — level up allocation)
  ↓
Item takıldığında delta popup gösterilir
  ↓
Buff/Debuff'lar renkli olarak eklenir
```

## 7. İlişkiler
- [Character](character.md), [Races](races.md), [Classes](classes.md), [Specializations](specializations.md)
- [Combat](combat.md), [Items](items.md), [Equipment](equipment.md), [Skills](skills.md)
- [Professions](professions.md), [Activity Engine](../systems/activity-engine.md)

## 8. Gereken Veriler
```
character_stats
  character_id
  primary { str, dex, int, wis, sta, cha }
  allocated_points_total
  cached_derived (server-side hesaplanır, cache olarak tutulur)
```
Not: derived stat'lar cache'lenir; item change/buff apply anında invalidate edilir.

## 9. I/O
- Input: `allocate_stat_point(character_id, stat, count)`.
- Output: yeni primary + recomputed derived.

## 10. Progression
- Level up başına 5 stat point.
- Ascension bonusu: her ascension level +1 kalıcı bütçe (Phase 02 balance).
- Mastery: belirli mastery track'leri belirli stat cap'lerini artırır (Phase 02).

## 11. Reward
Stat'lar dolaylı ödül; her item'ın stat rolü ekipman ödülünün cazibesini artırır.

## 12. Failure Durumları
- Allocation client hesabı server ile uyuşmazsa: server otoritedir; client patch alır ve UI reset.
- Invalid stat gönderirse: 400 hatası; hiçbir işlem yapılmaz.

## 13. Edge Cases
- Negative stat modifier (ör. debuff): stat 0 altına düşmez.
- Sonsuz item stacking: item slot sayısı sabit olduğu için doğal sınırlı.
- Stat cap değişimi (patch): cache invalidate + recompute.

## 14. AFK Davranışı
Stat'lar AFK hesabında (activity yield, combat outcome) doğrudan kullanılır. Örneğin `gathering_yield_pct` AFK mining ödülüne dahil.

## 15. Server Authoritative Gereksinimleri
- Tüm stat hesabı server'da.
- Point allocation transaction; race condition safe.
- Buff/Debuff timer server ticker'ı.

## 16. Genişletilebilirlik
- Yeni derived stat eklemek → hesaplama fonksiyonuna kayıt, UI listesine ekleme.
- Yeni primary stat eklemek **büyük iş**; class/race/item şemalarını etkiler. Bu ihtimal en aza indirilir.
- "Resistance types" ileride ateş/soğuk/nature/holy/shadow olarak zenginleşebilir; şu an yalnız `armor` + `magic_resist`.

## 17. Örnek Senaryo
Warrior level 30, STR 60, STA 40, item bonusları +20 STR / +15 STA:
- Attack Power = f(STR + item AP) → 250
- HP = f(STA) → 1200
- Armor = f(equipment) → 300
- Crit chance = f(DEX + item) → %8
Rage kaynağı = STR-dependent gain rate.

## 18. Balance Considerations
- Diminishing returns eğrisi düz olmamalı; ancak "duvara çarpma" hissi de vermemeli.
- Stat point atamada "yanlış yatırım" hafif olmalı (respec ucuz).
- Utility stat'lar (gold find, gather yield) itici olmamalı; sadece kozmetik ivme.
