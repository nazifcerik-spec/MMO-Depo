# Loot System (Ganimet)

## 1. Amaç
Item drop dağıtımını modellemek. Loot table, rarity roll, quantity, dedup, party share, personal roll.

## 2. Oyuncuya Fayda
- Adil ve şeffaf drop rate hissi.
- Rare drop dopaminini korumak.
- Party içinde adil paylaşım.

## 3. Kavramlar
- **Loot Table**: kaynak (mob, boss, node, chest) → items + weights.
- **Rarity Roll**: rarity band'a göre item pick.
- **Personal Loot**: her karakter kendi roll'unu alır (party share yok, klasik MMO'nun dışında).
- **Master Loot**: raid için opsiyonel (Phase 02).
- **Dedup**: aynı item aynı source'tan back-to-back drop olmasın (soft de-dup, opsiyonel).

## 4. Kurallar
1. Loot roll seed'li (bkz. [Random System](random-system.md)).
2. Loot table = list of `{item_template, weight, min_qty, max_qty, rarity_band}`.
3. Party'de default personal loot (herkes kendi roll'unu alır).
4. Master loot: leader dağıtım kararı; audit-log zorunlu.
5. BoP item drop olursa: doğrudan alıcının inventory'sine (mail overflow rule uygulanır).
6. Rare drop notification party-wide (Phase 02 opsiyonel).
7. Guaranteed drops (quest item) rarity roll'dan bağımsız.

## 5. Kullanıcı Akışı
```
Kill / node harvest / chest open
 ↓
Roll seed → items selected
 ↓
Distribute per rule → Reward System → Ledger
```

## 6. İlişkiler
- [Items](../game-design/items.md), [Combat](../game-design/combat.md), [Dungeons](../game-design/dungeons.md), [World Bosses](../game-design/world-bosses.md), [Professions](../game-design/professions.md), [Random System](random-system.md), [Reward System](reward-system.md)

## 7. Data Model
```
loot_table
  id, code, source_type (mob/boss/node/chest/quest)
  entries [] { item_template_id, weight, min_qty, max_qty, guaranteed (bool) }
  rarity_distribution (json)
  meta (json: e.g. anti_dedup_window, boss_bonus)

loot_roll
  id, source_event_id, seed, roll_result_json, distributed_to []
```

## 8. I/O
- Input: `roll_loot(source, char/party, seed)`.
- Output: item instances created, distributed.

## 9. Failure Durumları
- Loot table missing: safe fallback drop (junk).
- Duplicate roll (idempotency): no-op.
- Inventory full: mail overflow.

## 10. Edge Cases
- BoP + full inv: mail delivery + BoP preserved.
- Party wipe on kill hit: pending roll credited if kill_credit awarded.
- Drop chance ≥ 100% (guaranteed): min_qty enforced.

## 11. Server Authoritative Gereksinimleri
- Loot roll server; client sadece animasyon.
- Audit log tüm rolls.

## 12. Genişletilebilirlik
- Rarity boost buff (event).
- Item smart drop (Phase 02): karakter class'ına uygun item bias.
- Bad luck protection (Phase 02): rare item'ın drop şansı sürede artar.

## 13. Örnek
Ashen Wolf loot table:
```
common: leather_scrap ×2 (weight 60, guaranteed 1)
uncommon: wolf_fang (weight 25, min 0, max 1)
rare: alpha_pelt (weight 4)
epic: shadow_essence (weight 0.5)
legendary: ashen_fang (weight 0.02)
```

## 14. Balance Considerations
- Rarity dağılımı GDD ile uyumlu.
- Guaranteed drop garanti; junk floor engellenir.
- Bad luck protection ekleneceğinde audit ile validate.
