# Mounts (Binekler)

## 1. Amaç
Travel süresini azaltmak, dünyayı geniş hissettirmek, koleksiyon katmanı sunmak.

## 2. Oyuncuya Fayda
- Hızlı seyahat.
- Görsel prestij (rare mount).
- Kimi mount'lar aktivite bonusu (kargo, gather).

## 3. Gameplay Loop
- Mount summon → travel aktivitesi başlar → offline sürebilir.

## 4. Mount Türleri
- **Ground** — tüm oyunculara açık; base tier.
- **Cargo** — envanter genişlemesi (travel esnasında).
- **Aquatic** — su travel süresi kısaltıcı.
- **Winged** (Phase 02) — belirli bölgelere hızlı.
- **Cosmetic-only** — hız aynı, kozmetik farklı.

## 5. Temel Kurallar
- Aktif mount: 1.
- Mount collection unlocked ama bir seferde 1 tanesi aktif.
- Mount speed tier: 1.5×, 2×, 2.5× walk speed.
- Mount ölmez, combat sırasında dismount olur.
- Bazı mount'lar unlock condition ile gelir (quest, achievement, drop).

## 6. Kullanıcı Akışı
```
Mount Panel → owned mounts grid
 ├─ Set active
 └─ Summon (in-world)
```

## 7. İlişkiler
- [Travel](travel.md), [World](world.md), [Activity Engine](../systems/activity-engine.md), [Items](items.md), [Combat](combat.md)

## 8. Gereken Veriler
```
mount_template
  id, code, type (ground/cargo/aquatic/winged/cosmetic)
  speed_multiplier, cargo_bonus_slots, rarity, unlock_condition_json

mount_instance
  id, owner_char_id, template_id, active (bool)
```

## 9. I/O
- Input: activate mount, summon, dismount.
- Output: travel speed applied.

## 10. Progression
- Mount koleksiyonu genişledikçe achievement.
- Speed tier progression: level 20 unlock 2×, level 40 unlock 2.5×.

## 11. Reward
Mount kendisi ödül; travel süresi azaltıcı.

## 12. Failure Durumları
- Summon zorunlu şart (sadece safe zone dışı allowed).
- Combat başlarsa dismount + kısa summon cooldown.

## 13. Edge Cases
- Water crossing on ground mount: aquatic yok ise auto-swim, hız yavaş.
- Instance boundary crossing: dismount + walking teslim.

## 14. AFK Davranışı
Mount aktivitesi AFK olarak travel eder; seyahat aktivitesi tick'lenir.

## 15. Server Authoritative Gereksinimleri
- Speed multiplier server hesabı.
- Mount unlock condition kontrol.

## 16. Genişletilebilirlik
- Mount equipment (Phase 02): mount barding.
- Mount racing (Phase 02+).
- Guild-shared mount (Phase 03).

## 17. Örnek Senaryo
Karakter "Steel Warhorse" (2× speed, cargo +5 slot) aktif; şehirden hunt spot'a 8 dk yerine 4 dk'da varır.

## 18. Balance Considerations
- Mount olmadan da oyun akışı ciddi zarar görmemeli; travel süresi absorb edilebilir.
- Rare mount drop rate düşük; grinding değil, event/achievement/coleksiyon.
- Cargo mount bonusu belirli scenario'da; envanter genişletici alternatif olarak bank tab satın al mümkün.
