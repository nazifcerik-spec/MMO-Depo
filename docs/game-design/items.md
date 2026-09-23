# Items (Eşyalar)

## 1. Amaç
Dünyadaki tüm sahiplenilebilir nesneleri modellemek. Ekipman, kaynak, consumable, quest item, kozmetik, currency-item... hepsi tek soyutlamadan türetilir.

## 2. Oyuncuya Fayda
- Loot doyumu: rarity, affix, unique effect.
- Kişisel ilerleme: item'lar en görünür progression.
- Ticaret: değeri paylaşılan bir dil.

## 3. Gameplay Loop
- Session'da inventory, equip, sat.
- Meta'da tier atlama, set bonus, enchant, upgrade path.

## 4. Item Model Katmanları
**Template ↔ Instance** ayrımı kritik:
- **Item Template**: sabit tanım (id, name, base stats, art, rarity band). Data-driven.
- **Item Instance**: her drop bir instance; unique id, rolled stats, affixes, durability, socket, owner.

## 5. Item Türleri
- Equipment (weapon, armor, accessory, offhand)
- Consumable (potion, food, scroll)
- Material (ore, herb, wood, leather, cloth, gem, dust)
- Quest Item (soulbound, no trade)
- Currency Item (event token, marks)
- Cosmetic (dye, transmog appearance — Phase 02 için hazır)
- Utility (bag, key, map)
- Recipe (öğretici parşömen)

## 6. Rarity (Nadirlik) Sistemi
| Kod | Ad | Renk (öneri) | Affix bütçesi |
|---|---|---|---|
| `poor` | Bozuk | Gri | 0 |
| `common` | Sıradan | Beyaz | 0–1 |
| `uncommon` | Nadir | Yeşil | 1–2 |
| `rare` | Ender | Mavi | 2–3 |
| `epic` | Destansı | Mor | 3–4 + 1 socket |
| `legendary` | Efsanevi | Turuncu | 4 + unique effect |
| `mythic` | Efsanevi (server-first) | Kırmızı | Custom + progression gated |

## 7. Affix Sistemi
Item instance'a random atanır:
- **Prefix** (max 1) — genelde primary/derived stat.
- **Suffix** (max 1) — genelde utility (e.g. "of the Bear", +STA).
- **Enchant** (max 1) — enchanter tarafından uygulanır, geçici veya kalıcı.
- **Socket** — 0..N gem slot; gem'ler jewelcrafter tarafından üretilir.
- **Set bonus** — item bir "set" grubuna aitse, N item takıldığında bonuslar açılır.
- **Unique effect** — sadece legendary+ item'lar; skill-benzeri proc.

## 8. Temel Kurallar
- Instance stat roll'u seed + template range; server tarafında.
- Binding türleri:
  - `bop` — Bind on Pickup
  - `boe` — Bind on Equip
  - `boa` — Bind on Account (Phase 02)
  - `no_bind` — özgür ticaret
- Durability: kullanımda azalır; 0'da item passive'lerini kaybeder (item silinmez).
- Requirements: `min_level`, `class_id`, `spec_id`, `race_id` (opsiyonel), `stat_min`.
- Trading: BoE değil ve BoP değilse marketplace'e listelenebilir.

## 9. Kullanıcı Akışı
```
Drop → Notification
  → Inventory
    → Equip / Salvage / Sell / List
  → Equipment slot check (req'ler karşılanıyor mu)
  → Server transaction
```

## 10. İlişkiler
- [Inventory](inventory.md), [Equipment](equipment.md), [Loot System](../systems/loot-system.md), [Economy](economy.md), [Marketplace](marketplace.md), [Professions](professions.md), [Combat](combat.md)

## 11. Gereken Veriler
```
item_template
  id, code, name (i18n), description (i18n)
  category (enum), subcategory
  rarity (enum)
  base_stats (json)
  affix_pool_ids [] (prefix/suffix havuzları)
  socket_slots (int)
  set_id (nullable)
  binding_type
  requirements (json)
  base_durability (int, nullable)
  weight (int)
  art_ref
  unique_effect_id (nullable)

item_instance
  id (uuid)
  template_id
  owner_id (character/guild/mailbox/marketplace)
  rolled_stats (json)
  affixes (json)
  sockets (json: [{gem_instance_id, slot}])
  enchant_id (nullable)
  durability_current
  bound_to (nullable char id)
  created_at, seed (RNG replay için)
```

## 12. I/O
- Input: pick up (auto), equip, unequip, salvage, list, buy, mail.
- Output: inventory diff, equipment diff, stat recompute.

## 13. Progression
- Tier ekipman: level bantlarına bağlı (Tier I..V öneri).
- Legendary progression Phase 02 sistemi (uzun quest).
- Enchant path: base → +1..+10; upgrade risk vardır (Phase 02).

## 14. Reward
Item = dünyanın kendi ödül dili. Loot event'i player için primary dopamine.

## 15. Failure Durumları
- Envanter dolu → mailbox'a düşer (7 gün) veya "drop refused" işaretli node'da beklet (Phase 02).
- Requirement karşılanmıyor → equip reddedilir; UI görsel red frame.
- Durability 0 → passive/proc devre dışı; onarım için repair NPC.

## 16. Edge Cases
- Item template değişimi (patch): mevcut instance'lar `template_snapshot` alanına eski template hash'ini tutar; migration seçenekli.
- Item silme: sadece kullanıcı onaylı; UI çift onay + rare item için özel confirm.
- Duplication exploit önleme: item id'ler unique index; transfer atomik transaction.

## 17. AFK Davranışı
Item drop AFK sırasında activity yield'ına bağlanır; deterministik seed ile.

## 18. Server Authoritative Gereksinimleri
- Tüm item state (owner, stats, durability) server otoritesi.
- İstemci "envanterimde şu var" iddiasına asla güvenmez.
- Ekonomi audit log her item transferini yakalar.

## 19. Genişletilebilirlik
- Yeni rarity: enum extension + tablolar.
- Yeni affix pool: data.
- Transmog: `appearance_override` alanı planlı.
- Legendary: `unique_effect_id` + `legendary_progress_state` alanları.

## 20. Örnek Senaryo
Dungeon drop: `steel_longsword_of_the_bear` (rare, prefix "Steel" +5 STR, suffix "of the Bear" +10 STA, socket x1, durability 100, min level 20). Item instance seed'li, owner karakter, BoE.

## 21. Balance Considerations
- Rarity dağılımı: common %60, uncommon %25, rare %10, epic %4, legendary %0.9, mythic %0.1 (dungeon/world bosstan varyasyon).
- Affix budget kontrolü: prefix+suffix stat toplamı belirlenen band aralığında.
- Enchant/upgrade risk sistemi kullanıcı için "yıkıcı" hissetmemeli; fail state için materyal geri kazanımı olmalı.
