# Races (Irklar)

## 1. Amaç
Oyuncuya karakterinin **köken kimliği** hissini vermek; class'tan bağımsız ama tamamlayıcı bir seçim katmanı sunmak. Irklar sadece kozmetik değildir; ancak class seçimini zorlayacak kadar da baskın değildir.

## 2. Oyuncuya Fayda
- Karaktere kalıcı bir "renk" verir (kültür, doğa, özellik).
- Küçük ama hissedilebilir mekanik farklar (stat, passive, profession affinity).
- Sosyal kimlik: guild composition, RP, kozmetik uyum.

## 3. Gameplay Loop'taki Yer
- Character creation aşamasında **tek seferlik** seçim.
- Session Loop'u doğrudan değiştirmez, ancak aktivite tercihini yönlendirir (profession bonus vs.).
- Meta Loop'ta racial passive'ler yeni içerikte hâlâ önemli olacak şekilde tasarlanır.

## 4. Temel Kurallar
- Race seçimi **karakter yaratımında** yapılır; sonradan değiştirilemez (Phase 02'de ücretli race change ihtimali açık).
- Race seçimi **hiçbir class'ı kilitlemez**.
- Her ırkın hedef stat sapması ±%3–5 civarında tutulur.
- Her ırkın: base stats modifier, 1 passive racial, 1 aktif racial ability, 1 profession affinity, 1 unique flavor mekaniği vardır.
- Racial ability'ler class ability'lerinden **etki bütçesi olarak daha küçüktür**; sinerji öneri katmanıdır, ana silah değil.

## 5. Kullanıcı Akışı
```
Character Creation
  ↓
Race listesi ekranı
  ├─ Görsel + kısa lore
  ├─ Mekanik özet: stat delta, passive, active, affinity
  ↓
Seç → Onay
```

## 6. İlişkiler
- [Character](character.md), [Classes](classes.md), [Stats](stats.md), [Skills](skills.md), [Professions](professions.md), [Combat](combat.md)

## 7. Gereken Veriler
```
race
  id
  code (string, unique)
  display_name (i18n key)
  lore_key (i18n)
  stat_modifiers (json: {str:+1, dex:-1, ...})
  passive_ability_id (fk → abilities)
  active_ability_id (fk → abilities)
  profession_affinity (fk → professions, bonus_pct)
  unique_mechanic_flag (enum)
```

## 8. I/O
- Input: creation sırasında `race_id` seçimi.
- Output: karakter sheet'ine stat modifier yansıması, ability listesine racial abilities eklenmesi.

## 9. Progression
- Race seviyesi yoktur.
- Racial passive'ler karakter leveline **hafifçe** ölçeklenebilir (ör. +1 stat per 10 level, opsiyonel).

## 10. Reward
Race, aktif ödül dağıtmaz; ödüllerin efektif değerini ince modüle eder (profession'da +%X kaynak drop, vs.).

## 11. Failure Durumları
- Race data'sı yüklenemezse: karakter yaratma disabled; UI net hata gösterir.
- Racial ability serverdan reddedilirse: silent fallback + telemetri.

## 12. Edge Cases
- Bir ırkın nerf'i: mevcut karakterlere free "Racial Rebalance Token" verilir (compensating token modeli).
- Yeni ırk eklenmesi: mevcut karakterlere unlock quest / one-time race change window sunulabilir.

## 13. AFK Davranışı
Racial passive'ler AFK aktivite hesabına dahildir: örneğin "mining +%3" varsa AFK mining ödülü hesabına girer.

## 14. Server Authoritative Gereksinimleri
- Stat modifier'lar server tarafında uygulanır.
- Ability trigger'ları server tarafında doğrulanır.

## 15. Genişletilebilirlik
- Yeni ırk = yeni satır + yeni ability id'leri. Kod değişikliği gerektirmez.
- Race sub-culture eklemek istersek: `race_variant` alt tablo (ör. Aynı ırktan farklı klanlar).

## 16. Örnek Race Şablonu (Placeholder, isim/ton Phase 01'de finalize)

> Bu ırklar örnektir; final isim ve lore Phase 01 içerik pass'inde belirlenir. Sayılar yaklaşık başlangıç önerisidir.

| Race (kod) | Kimlik | Stat modifiers | Passive | Active | Profession Affinity |
|---|---|---|---|---|---|
| `human` | Denge, adaptasyon | +1 herşey | +5% XP gain | Rally: kısa süreli +5% party stat | Trading |
| `dwarf` | Dağ, zanaat, direnç | +2 STA, +1 STR, -1 DEX | +5% mining yield | Stone Skin (defense burst) | Mining |
| `elf` | Doğa, hız, isabet | +2 DEX, +1 INT, -1 STR | +5% woodcutting yield | Wind Step (movement) | Woodcutting |
| `orc` | Vahşi güç | +2 STR, +1 STA, -1 INT | +5% melee crit dmg | Blood Roar (rage) | Hunting |
| `merfolk` | Kıyı, su, esneklik | +2 DEX, +1 INT, -1 STR | +5% fishing yield | Tide Grace (evasion) | Fishing |
| `arcanborn` | Sihir yatkınlığı | +2 INT, +1 WIS, -1 STA | +5% mana regen | Arcane Surge (spell burst) | Enchanting |

Her ırk için ayrıca "unique flavor" bir mekaniği vardır (ör. `dwarf`: gece görüşü ile dungeon vision +1 tile; `merfolk`: su altı travel maliyeti -%50).

## 17. Balance Considerations
- Her class için "en iyi ırk" olmamalı; her ırk 2–3 class'la öne çıkabilecek meta sinerjilere sahip olmalı.
- Passive stat modifier toplamı bütçe: net toplam +2 stat (bir negatif yumuşatır).
- Racial active'in cooldown'u genelde 3–5 dakika arası; combat game-changer olmasın.
- Profession affinity bonusu %3–%8 aralığı; ekonomiyi bozacak seviyede değil.
