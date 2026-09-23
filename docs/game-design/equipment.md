# Equipment (Ekipman)

## 1. Amaç
Karakterin üzerinde taşıdığı item'ları modellemek; stat, görsel, set etkileri, weapon skill etkileşimlerini sunmak.

## 2. Oyuncuya Fayda
- Karakter kimliğinin görsel ve mekanik ifadesi.
- Loot planlaması (upgrade path).
- Set bonusu için hedefli farming.

## 3. Gameplay Loop
- Session'da yeni drop → equip → stat delta değerlendirme.
- Meta'da tier ilerlemesi, set completion, enchant path.

## 4. Slot Yapısı
| Slot | Örnek |
|---|---|
| Head | Miğfer, kukulet |
| Neck | Kolye |
| Shoulder | Omuzluk |
| Back | Pelerin |
| Chest | Zırh gövdesi |
| Wrist | Bilek koruma |
| Hands | Eldiven |
| Waist | Kemer |
| Legs | Bacaklık |
| Feet | Ayakkabı/pota |
| Ring | Yüzük ×2 |
| Trinket | Muska ×2 |
| Main Hand | Ana silah |
| Off Hand | Kalkan / sekonder silah / focus (mage) |
| Ranged | Yay/tatar (Ranger için main hand olabilir) |

## 5. Temel Kurallar
- Equip için req'ler karşılanmalı (level, class, spec, race, stat).
- Dual-wield / 2-handed weapon: class + spec bağlı.
- **Weapon skill** ayrı bir passive proficiency değil; class weapon affinity ile modellenir.
- Loadout preset (Phase 02): kayıtlı gear + rotation seti.
- Transmog (Phase 02): görsel override; stat item'dan gelir.

## 6. Kullanıcı Akışı
```
Character Sheet → Equipment panel (görsel doll)
  ↓
Slot'a click → uyumlu item filtreli inventory picker
  ↓
Equip → stat delta → server confirm
```

## 7. İlişkiler
- [Items](items.md), [Inventory](inventory.md), [Stats](stats.md), [Classes](classes.md), [Specializations](specializations.md), [Combat](combat.md)

## 8. Gereken Veriler
```
character_equipment
  character_id
  slot (enum) → item_instance_id
  loadout_id (nullable Phase 02)
```

## 9. I/O
- Input: equip/unequip.
- Output: stat recompute event.

## 10. Progression
- Tier'lar ekipman ile gelir; oyuncunun ilerleme metriği ekipman tier'ı ve enchant/socket derinliği.
- Set bonusları (2/4/6 parça) mid-endgame'de belirginleşir.

## 11. Reward
Equip fiili "ödülü giymek"tir; loot doyumunun aktivasyon anı.

## 12. Failure Durumları
- Requirement fail → net UI mesajı.
- Two-hand equip iken off-hand doluysa: off-hand otomatik unequip önerisi (confirm).
- Durability 0: passive/proc çalışmaz; base stat kalır (Phase 01 karar; alternatif tamamen deactivate).

## 13. Edge Cases
- Set item'ın 4/6 bonusu için 5 parça takılı: bir sonraki eşyayı takana kadar 4-parça bonusu aktif.
- Race change (Phase 02): race-bound item'lar geçici olarak grey out edilir.

## 14. AFK Davranışı
Ekipman stat AFK combat hesabına doğrudan girer. AFK'de equip değiştirilemez (login gerekir).

## 15. Server Authoritative Gereksinimleri
- Equip transaction; stat recompute server.
- Loadout switch transaction; instant değil, 1 sn "swap animation" (exploit önleme).

## 16. Genişletilebilirlik
- Transmog `appearance_id` alanı hazır.
- Runes / catalysts (Phase 02) ek slot olarak eklenebilir.
- "Ornament slot" (kozmetik aksesuar) eklenebilir.

## 17. Örnek Senaryo
Warrior/Bulwark: 2H sword slot'a takılmak istenirken shield equip'te; UI "Shield unequip olacak, devam?" onayı → equip → stat delta gösterilir.

## 18. Balance Considerations
- Set bonus tasarımı: 2p konfor, 4p sinergi, 6p signature.
- Legendary sayısı sınırlı; oyuncu 2 legendary'yi aynı anda kullanabilir (soft cap).
- Weapon speed vs damage: fast weapon = düşük burst, yüksek proc; slow weapon = yüksek burst, düşük proc.
