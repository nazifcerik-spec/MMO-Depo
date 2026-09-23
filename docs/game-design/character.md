# Character

## 1. Amaç
Oyuncunun dünyayla etkileşimindeki avatarını tanımlamak. Karakter, ilerlemenin, sahipliğin ve kimlik hissinin merkezidir.

## 2. Oyuncuya Fayda
- Kalıcı bir kimlik: isim, ırk, sınıf, spec, kozmetik.
- Bir "hesap üstü" ilerleme ekseni: [Ascension](ascension.md) ile hesap seviyesi kalır.
- Sosyal görünürlük: guild, marketplace, party'de tanınırlık.

## 3. Gameplay Loop'taki Yer
- Session Loop'ta özne.
- Meta Loop'ta ilerlemenin taşıyıcısı.
- Sosyal loop'ta kimlik.

## 4. Temel Kurallar
- Hesap başına başlangıçta **2 karakter slotu**, ilerleyen sürümlerde kozmetik olarak genişletilebilir (bkz. GDD §13 monetization).
- Her karakterin **1 ırkı, 1 base class'ı** vardır; class değiştirilemez, spec değiştirilebilir (respec ile).
- Karakter silinirse **7 günlük soft-delete** window'u vardır; sonra kalıcı silinir.
- İsim rezervasyonu: karakter aktifken 12 ay pasif kalırsa isim yeniden havuza girebilir (Phase 02 politikası).

## 5. Kullanıcı Akışı
```
Create Character
  → İsim (unique, kurallar: 3-16 char, alfanumeric)
  → Race seçimi (bkz. Races)
  → Class seçimi (bkz. Classes) — spec ilk 5. levelde açılır
  → Görsel seçimler (Phase 02'de zenginleşecek)
  → Onay → Ana şehirde spawn
```

## 6. İlişkiler
- [Races](races.md), [Classes](classes.md), [Specializations](specializations.md)
- [Stats](stats.md), [Skills](skills.md)
- [Inventory](inventory.md), [Equipment](equipment.md)
- [Progression](progression.md), [Ascension](ascension.md), [Mastery](mastery.md)
- [Professions](professions.md) — karakter bazlı sayaçlar

## 7. Gereken Veriler (Data Model)

Yüksek seviyeli alanlar (Phase 01 şemasına referans):

```
character
  id (uuid)
  account_id
  name (unique, indexed)
  race_id (fk → races)
  class_id (fk → classes)
  specialization_ids [] (0..2 çakışık spec)
  level (int)
  xp (bigint)
  ascension_level (int)         -- default 0
  ascension_points (int)        -- default 0
  created_at, last_seen_at, last_reveal_at
  current_activity_id (nullable)
  current_zone_id (fk → territories)
  home_city_id (fk → cities)
  soft_deleted_at (nullable)
  cosmetic_json (blob)
  future:
    housing_id (nullable)
    faction_reputations (json)
    legendary_progress (json)
```

Alt tablolar: `character_stats`, `character_inventory`, `character_equipment`, `character_skills`, `character_professions`, `character_activity_log`.

## 8. I/O
- Input: create, rename (paid/limited), change cosmetics, respec, delete.
- Output: karakter profili (public), karakter sheet (private), activity contract owner.

## 9. Progression
- Level 1 → cap (Phase 01'de cap 50 önerisi; final cap Phase 02).
- Level başına stat + skill point (bkz. Stats, Skills).
- Cap sonrası Ascension kapısı açılır (bkz. Ascension).

## 10. Reward
Karakter kendi başına ödül üretmez; ödüllerin sahibidir. Level up görsel/UI kutlama tetikler, stat point verir.

## 11. Failure Durumları
- Karakter ölümü: karakter silinmez; XP debt, durability loss, kısa spawn timer.
- Create sırasında isim çakışması: server 409 döner, istemci öneri gösterir.
- Silme kilitlerken: guild leader ise "önce lider transferi" uyarısı.

## 12. Edge Cases
- Karakter, uzun offline sonrasında dönerse: sadece son aktivite hesaplanır; nesnesi değişmiş item template'ler için migration çalışır.
- Class silme kararı: temel prensip "asla". Balance amaçlı bir class'ın kaldırılması gerekirse mevcut karakterler "legacy" olarak korunur ve free respec ile başka class'a taşınır — **açıkça oyuncuya duyurulur.**

## 13. AFK Davranışı
Karakter offline iken **yalnızca son başlatılan aktivite** ile ilerler. Karakter, seçim yapamaz; sadece "önceden dediği" işi yapmaya devam eder. Karar odaklılık ilkesi.

## 14. Server Authoritative Gereksinimleri
- Karakter oluşturma race condition korumalı (isim unique index).
- Level up, XP, stat point tahsisi server hesabı.
- Cosmetic dahi server tarafında doğrulanır (kilitli görselin kilit durumu).

## 15. Genişletilebilirlik
- `cosmetic_json` esnek; ileri transmog için ayrılmış.
- `faction_reputations`, `housing_id`, `legendary_progress` alanları şimdiden yer tutar.
- Character-level flag'ler için `character_flags` k/v tablosu (unlock ve tutorial için).

## 16. Örnek Senaryolar
**A — Yeni oyuncu:** İnsan / Warrior yaratır. Level 5'te Guardian (defense) spec'ini açar. Level 10'da 2. spec katmanı olan Bulwark'u seçer.

**B — Deneyimli oyuncu:** Level cap'e ulaşır, Ascension'a girer. Level 1'e döner, karakteri kalır, ascension_level +1 olur ve kalıcı meta bonus kazanır.

**C — Kimlik hatası:** Oyuncu isim yazımını beğenmez; ücretli rename kullanır. İsim havuzu güncellenir.

## 17. Balance Considerations
- Erken level'da (1–10) hızlı stat point/skill point akışı: oyuncuyu class kimliğine bağlar.
- Cap sonrasında marjinal stat kazanımı azalır; artan derinlik Mastery + Ascension'dan gelir.
- Race + Class stat toplamı, tüm kombinasyonlar için ±%5 içinde tutulur (Phase 01 sim).
