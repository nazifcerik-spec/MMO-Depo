# Professions (Meslekler)

## 1. Amaç
Karakteri savaş dışı ekonomik döngüye bağlamak. Hem kaynak üretimi hem de kişisel ilerleme ekseni. AFK aktivitelerin çoğunun kabuğu.

## 2. Oyuncuya Fayda
- Ekonomiye katılım.
- Kendi ekipmanını üretebilme.
- Class'tan bağımsız uzun soluklu ilerleme.
- Offline zamanı verimli kılan aktivite akışı.

## 3. Gameplay Loop
- Session'da profession değiştirme ve node seçimi.
- Daily'de günlük profession quest.
- Meta'da profession tier atlama, unique recipe unlock.

## 4. Profession Türleri (Phase 00 seti)

### Toplayıcı (Gathering)
- **Mining** — cevher, taş.
- **Woodcutting** — odun, reçine.
- **Fishing** — balık, deniz mahsulü.
- **Foraging** — bitki, mantar, otsu materyal.
- **Skinning** — deri, kürk (av + skin).

### Üretici (Crafting)
- **Blacksmithing** — silah, ağır zırh.
- **Leatherworking** — orta/hafif zırh, aksesuar.
- **Tailoring** — kumaş zırh, tekstil.
- **Alchemy** — potion, elixir, tonic.
- **Cooking** — buff yiyeceği, uzun süreli passive.
- **Enchanting** — item enchant, disenchant.
- **Jewelcrafting** — gem cut, aksesuar.

### Servis (utility)
- **Cartography** (map notes) — Phase 02.
- **Fletching** (bow/arrow) — Phase 02.

## 5. Temel Kurallar
- **3 aktif profession slot** (Phase 00). Meslek unutup yenisini öğrenme cooldown'u vardır (Phase 02 opsiyonel).
- Meslek seviyeleri: 1–100 tier'lar (Novice/Apprentice/Journeyman/Expert/Master/Grandmaster).
- Node respawn süresi bölgeye + rarity'ye bağlı.
- Craft başarı şansı stat'a (WIS + tool + recipe rank) bağlı.
- **Failure yumuşak**: crit fail'de bile bir miktar materyal kurtarılır.

## 6. Kullanıcı Akışı
```
Profession panel
 ↓
Aktif profession seç
 ↓
Node listesi (bölge bazlı) → seç
 ↓
Aktivite başlat → Activity Engine'e kayıt
 ↓
AFK'de tick'lenir → login'de reveal
```

## 7. İlişkiler
- [Items](items.md), [Economy](economy.md), [Activity Engine](../systems/activity-engine.md), [AFK Engine](../systems/afk-engine.md), [World](world.md), [Cities](cities.md)

## 8. Gereken Veriler
```
profession
  id, code, category
  max_tier, xp_curve

character_profession
  character_id, profession_id, level, xp, active (bool), unlocked_at

recipe
  id, profession_id, tier_required
  inputs [] {item_id, qty}
  outputs [] {item_id, qty, chance}
  time_ms, xp_reward
  station_required (nullable)

resource_node
  id, location_id, resource_type
  tier, respawn_ms
  yield_curve (per skill level)
```

## 9. I/O
- Input: `start_activity(profession_activity)`, `craft(recipe_id, qty)`.
- Output: inventory delta, profession xp, activity tick log.

## 10. Progression
- XP kazanımı = base + rank * multiplier.
- Level cap'te unique recipe (grandmaster) unlock.
- Bazı recipe'ler quest ile açılır; bazıları rare drop; bazıları world event.

## 11. Reward
- Materyal + intermediate + craft output.
- XP + occasional rare procs (double yield, bonus rarity).
- Achievements ve title'lar (Phase 02).

## 12. Failure Durumları
- Node exhausted: aktivite otomatik komşu node'a atlar veya durur (oyuncu tercih).
- Envanter dolu: activity pause + notification (bkz. AFK Engine).
- Craft fail: input'lar (kısmi geri kazanım) düşer, xp yarısı verilir.

## 13. Edge Cases
- Aynı bölgede birden çok gatherer: node share vs personal instance (Phase 01 karar — öneri: personal instance; MMO yoğunluğu daha temiz).
- Rare node (tier üstü): server-wide anons; ilk gelen 1–3 kişi drop.
- Recipe silinmesi: tarih arşivlenir; oyuncu envanterinde item legacy olarak kalır.

## 14. AFK Davranışı
Tüm profession'lar AFK öncelikli. Node yield deterministik tick + seed; envanter kapasitesine ve stamina'ya (varsa) göre durur.

## 15. Server Authoritative Gereksinimleri
- Tüm yield hesabı server.
- Recipe input/output atomic transaction.
- Node respawn timer server ticker.

## 16. Genişletilebilirlik
- Yeni profession = data satırı + node/recipe seed.
- Servis meslekleri (map, fletch) sonra eklenebilir.
- Profession'lar arası crosscraft (ör. blacksmith + enchant kombine crafts) recipe içinde tanımlı.

## 17. Örnek Senaryo
Karakter mining tier 40; iron node'da AFK; 8 saatlik offline; login'de: iron ore x120, stone x50, %2 rare gem x2 drop. Envanter %80'de durmuş.

## 18. Balance Considerations
- Kaynak fiyatlarının stabilitesi ekonominin sağlığı. Faucet monitoring (bkz. [Economy](economy.md)).
- Gathering profession'ları crafting'in altını doldurabilecek yeterlilikte olmalı; ancak "her şey sat, para bas" ekonomisi engellenmeli (sink'ler önemli).
- Recipe rarity dağılımı: %70 common, %25 uncommon (dungeon/quest), %5 rare (event/world boss drop).
