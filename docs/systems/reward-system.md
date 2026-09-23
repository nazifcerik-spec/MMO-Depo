# Reward System (Ödül Sistemi)

## 1. Amaç
Farklı ödül tiplerinin (XP, gold, item, currency, mastery xp, faction rep) tek bir soyutlama altında dağıtılması. Idempotent ve auditable.

## 2. Oyuncuya Fayda
- Ödüller şeffaf, tutarlı, sürprizsiz.
- Ödül dağıtımında hata olsa da retry safe.

## 3. Kavramlar
- **Reward Payload**: bir ödül olay birimi.
- **Reward Source**: activity, quest, event, achievement, dungeon, boss, market.
- **Reward Ledger**: append-only kayıt.
- **Idempotency Key**: source event id + character id.

## 4. Reward Türleri
- XP (character)
- XP (profession)
- XP (mastery)
- Gold
- Currency token
- Item instance (bkz. [Loot System](loot-system.md))
- Faction rep (Phase 02)
- Ascension point
- Cosmetic unlock

## 5. Kurallar
1. Her ödül idempotency key ile yazılır.
2. Ödül transaction'lıdır; kısmi başarı yok.
3. Envanter dolu ise item mailbox'a (auto).
4. Gold cap altında (üst sınıra çarpınca overflow: warning + user opt).
5. Faction reputation cap: default 21000 (Exalted-style, Phase 02).
6. Aynı source event iki kez ödül üretemez.

## 6. Kullanıcı Akışı
```
Event triggers → RewardPayload built → dispatch
 ↓
Ledger write → apply to character
 ↓
Notification (optional)
```

## 7. İlişkiler
- Tüm aktivite/quest/event sistemleri; [Economy](../game-design/economy.md), [Items](../game-design/items.md), [Progression](../game-design/progression.md)

## 8. Data Model
```
reward_ledger
  id, character_id, source_type, source_id
  payload_json (xp, gold, items, tokens, mastery_xp, rep)
  idempotency_key (unique)
  applied_at
```

## 9. I/O
- Input: `grant(character_id, reward_payload, source, idempotency_key)`.
- Output: ledger entry + character state deltas.

## 10. Failure Durumları
- Idempotency key exists: no-op (safe).
- Character not found: reject.
- Inventory overflow: mail routing.

## 11. Edge Cases
- Multi-currency payload: her currency ayrı ledger row veya JSON payload; audit-friendly.
- Faction rep cap: yüksek rep'te absorb; overflow notification.

## 12. Server Authoritative Gereksinimleri
- Ledger append-only.
- Idempotency enforced.
- Rate limit yok (aktivite tick sıklığı zaten dolaylı sınırlar).

## 13. Genişletilebilirlik
- Yeni reward tipi (housing, mount, kozmetik): payload extension.
- Reward multiplier (event 2×): calculator layer.

## 14. Örnek Senaryo
Boss kill event → reward: 400 gold, 1200 xp, 1 rare item, 30 dungeon_token, 50 mastery_xp[weapon_sword]. Idempotency key = `boss:ashen_wyrm:spawn_id_1234:char_777`.

## 15. Balance Considerations
- Reward ölçeği aktivite süresine göre; 1 saat AFK ≈ 1 dungeon completion range.
- Rare drop chances güvenilir loot table üstünden gelir.
- Reward inflasyon monitör (bkz. [Economy](../game-design/economy.md)).
