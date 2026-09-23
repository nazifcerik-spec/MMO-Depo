# Inventory (Envanter)

## 1. Amaç
Karakter, guild, mailbox ve stash düzeylerinde eşya sahipliğini modellemek. Envanter, oyun ekonomisinin ve loot doyumunun mekanik kabuğudur.

## 2. Oyuncuya Fayda
- Kazandıklarını görme, düzenleme, kullanma alanı.
- Kapasitenin, aktiviteyi disipline eden bir sınırlayıcı olması ("bag full → decision needed").
- Filtreleme, sıralama, wishlist gibi konforlar.

## 3. Gameplay Loop
- Session'da hep açılır: envanter temizliği, salvage, junk sell.
- Meta'da bag upgrade, stash expansion.

## 4. Envanter Katmanları
1. **Bag** — karakter üzerinde taşınan; slot bazlı; kapasite item weight veya slot sayısı ile sınırlı.
2. **Bank / Stash** — şehir bankası; guild bank ayrı.
3. **Mailbox** — geçici; 7–14 gün item retention.
4. **Loadout Storage** — set of equipment presets (Phase 02).
5. **Currency Wallet** — item değil, currency ayrı ledger.

## 5. Temel Kurallar
- Baz slot sayısı: 40. Bag item'ları slot'u genişletebilir (max 120).
- Stack rules: aynı template + aynı önemli flag (rarity, bound status) → stack; instance-only item'lar stack olmaz.
- Ağırlık sistemi opsiyonel; başlangıçta **sadece slot sınırı** — sadelik.
- Slot dolduğunda: activity `paused` state'ine geçer, notification tetikler.

## 6. Kullanıcı Akışı
```
Inventory Panel
 ├─ Filters: type, rarity, usable, tradable
 ├─ Sort: rarity | level | name | value
 ├─ Bulk actions: salvage, sell, mail
 └─ Context menu: use, equip, split, list, destroy (double confirm)
```

## 7. İlişkiler
- [Items](items.md), [Equipment](equipment.md), [Economy](economy.md), [Marketplace](marketplace.md), [Activity Engine](../systems/activity-engine.md), [AFK Engine](../systems/afk-engine.md)

## 8. Gereken Veriler
```
character_inventory
  character_id, slot_index, item_instance_id, stack_count
  auto_pause_on_full (bool, default true)

bank_inventory
  character_id, tab_index, slot_index, item_instance_id, stack_count

guild_bank
  guild_id, tab_index, slot_index, item_instance_id, stack_count

mailbox
  character_id, mail_id, items_json, expires_at
```

## 9. I/O
- Input: move, split, stack, salvage, sell, mail.
- Output: inventory diff event, notification.

## 10. Progression
- Bag upgrade: quest reward veya crafted bag item.
- Bank tab unlock: gold sink (5k, 20k, 100k, ... artan maliyet).
- Guild bank tab unlock: guild progression.

## 11. Reward
Envanterin kendisi ödül üretmez; ödüllerin muhafızıdır.

## 12. Failure Durumları
- Slot dolu → activity pause + drop overflow mailbox'a.
- Mailbox dolu → yeni mail reddedilir (client visible warning); expire olmuş item'lar arşive.
- Concurrent modification (aynı slot iki işlem): transaction rollback; UI refresh.

## 13. Edge Cases
- Item silme UX'i: rare+ için 3 sn hold-to-confirm.
- Migration: item template değişince stack merge logic re-run.
- Guild bank hak yönetimi: rank bazlı read/write; audit log zorunlu.

## 14. AFK Davranışı
Envanter durumu, AFK Engine'in "activity yield stopping condition"udur. Full slot → yield stop (yield discard değil; sadece durur).

## 15. Server Authoritative Gereksinimleri
- Tüm envanter işlemleri transaction'lı.
- Concurrent atomic; optimistic locking (version alanı).
- Ekonomi ledger'a event yansır.

## 16. Genişletilebilirlik
- Auto-salvage kural setleri (Phase 02): rarity < X, level < Y → auto.
- Wishlist ve auto-mail (Phase 02).
- Cross-character stash (BoA sonrası) planlı.

## 17. Örnek Senaryo
AFK mining, envanter 40/40 dolduğunda activity pause. Player login, bulk salvage (uncommon+altı) → 12 slot açıldı, activity resume opsiyonu popup.

## 18. Balance Considerations
- Envanter darlığı bir gameplay tension'dır; ancak zorlayıcı olmamalı. Baz 40 slot ile 8 saat AFK aktivite tam yönetilebilir olmalı.
- Bank tab maliyetleri gold sink olarak ekonomiye hizmet eder.
- Salvage → material dönüşüm oranı, ekonomiye materyal enflasyonu getirmeyecek şekilde ayarlanır.
