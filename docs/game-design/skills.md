# Skills (Yetenekler)

## 1. Amaç
Karakterin muharebe ve etkileşim eylemlerini modellemek. Skill'ler class/spec/race'in **fiili çıkış aracıdır**.

## 2. Oyuncuya Fayda
- Yaptığı hareketin sonuca yansıdığını görmek.
- AFK rotation için kişiselleştirilebilir bir "otomatik strateji" oluşturmak.
- Yeni skill unlock etmek, ilerleme hissinin bir kanadı.

## 3. Gameplay Loop
- Session'da manuel muharebe / rotation editor.
- Meta'da skill upgrade + spec unlock.

## 4. Skill Türleri
- **Active** — kullanılır, resource harcar, cooldown vardır.
- **Passive** — sürekli etkili.
- **Utility** — muharebe dışı (travel, gather).
- **Racial** — race'ten gelir.
- **Ultimate** — uzun cooldown, yüksek etki.

## 5. Temel Kurallar
- Karakter aynı anda **max 12 aktif skill** slot'unda taşıyabilir (bar) + ekstra passive'ler.
- Cooldown ve resource server otoritesi.
- Skill upgrade seviyeleri: base → rank 1..5 (Phase 02'de rank sistemi genişleyebilir).
- Global Cooldown (GCD): 1.0 sn baz; haste ile azalır (min 0.5 sn).
- **Rotation Editor** var: oyuncu koşullu öncelik listesi kurar; AFK'de bu liste yürütülür.

## 6. Kullanıcı Akışı
```
Skill Panel → Sekmeler: Active | Passive | Utility | Racial | Ultimate
  ↓
Slot yerleşimi (drag & drop)
  ↓
Rotation Editor (koşul: HP<%X, cast X, else Y)
  ↓
Save → Server validate
```

## 7. İlişkiler
- [Classes](classes.md), [Specializations](specializations.md), [Races](races.md), [Combat](combat.md), [Stats](stats.md), [Items](items.md) (weapon proc)

## 8. Gereken Veriler
```
skill
  id, code, display_name (i18n)
  type (active/passive/utility/racial/ultimate)
  class_id (nullable — racial ise null)
  spec_id (nullable)
  race_id (nullable)
  resource_type, resource_cost
  cooldown_ms
  gcd_bypass (bool)
  cast_time_ms
  range (melee/short/mid/long/global)
  aoe (bool + shape)
  effects_json (damage, heal, buff, debuff, dot, hot, shield, teleport)
  scaling (json: coeffs per primary stat)
  unlock_condition (level, spec)
  rank_progression_ids [] (fk: skill_rank tablosu)
```

## 9. I/O
- Input: `use_skill(character_id, skill_id, target)`; `set_rotation(character_id, rotation_json)`.
- Output: server validated cast event, damage/heal event.

## 10. Progression
- Level up → yeni skill unlock (belirli level'da).
- Skill rank up: skill point + ilgili skill'i kullanma sayısı (use-based unlock).
- Talent point (Phase 02): passive/branch bonus.

## 11. Reward
Skill değil ödül üretir; ancak yeni skill unlock'u ilerleme ödülüdür.

## 12. Failure Durumları
- Insufficient resource: reject + client cooldown hint.
- Cooldown bitmemiş: reject + remaining ms.
- Out of range: reject + reason.
- Rotation editor sonsuz döngü riski: server rotation'ı tick'te maks N eval kontrol.

## 13. Edge Cases
- Skill nerf'i: ranklara yansır; ilgili karakterlere respec pencere açılabilir.
- Interrupt: cast time'lı skill, hasar alınca cast lock (Phase 02'de interrupt system detayı).

## 14. AFK Davranışı
Rotation editor'ün yazdığı priority list AFK'de yürütülür. Skill server tarafında tick'lenir; kaynak/cooldown deterministik simüle edilir.

## 15. Server Authoritative Gereksinimleri
- Cooldown + resource kontrol serverda.
- Rotation eval serverda; hile korumasız değil.
- RNG: crit/proc [Random System](../systems/random-system.md) üzerinden seed'li.

## 16. Genişletilebilirlik
- Yeni skill = data satırı + effect combiner.
- Effect combiner (damage/heal/buff/...) modüler; yeni effect tipi eklemek şema ekleme.
- "Skill trees" gelecek (Phase 02); mevcut şema uyumlu.

## 17. Örnek Senaryo
Warrior/Guardian: Shield Slam (rank 3, 8 sn CD) → Threat aura → Bulwark Stance (passive) → Rotation editor "HP<%50 → Second Wind (heal); else Shield Slam; else auto".

## 18. Balance Considerations
- Ultimate cooldown 90 sn; oyunun içeriğine göre 2 ult / dungeon avg.
- Skill scaling coeff'leri sabit stat bütçesi içinde; yüksek scaling düşük CD olamaz.
- AFK verimi ile manuel oyun verimi arasında en fazla ~%20 fark; manuel daha iyi ama zorunlu değil.
