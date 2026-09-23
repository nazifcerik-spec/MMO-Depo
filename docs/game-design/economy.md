# Economy (Ekonomi)

## 1. Amaç
Oyunun uzun vadeli sağlığını belirleyen kaynak akış sistemi. Para, materyal, item, currency çeşitliliğinin dolaşım kurallarını yazmak. Faucet & sink dengesi.

## 2. Oyuncuya Fayda
- Adil, öngörülebilir bir ticaret ortamı.
- Değer üretebilen aktiviteler.
- Enflasyon nedeniyle emeğin buhar olmadığı sürdürülebilir bir dünya.

## 3. Gameplay Loop
- Session'da vendor, market, bank.
- Meta'da market trend analizi, spec supplier, price arbitrage.

## 4. Para Birimleri
- **Gold** — primary currency; her yerde geçer.
- **Silver / Copper** — gold alt birimleri (görsel/sunum katmanı).
- **Currency Items** — dungeon token, event token, faction mark, honor point (Phase 02 PvP), guild seal.
- Currency wallet: server-side ledger; item değil.

## 5. Faucet (Para Girişi) Kaynakları
- Mob kill (küçük).
- Quest reward.
- Vendor sale (junk sell).
- Marketplace sale (peer-to-peer transfer; net yeni para değil, sadece transfer).
- Achievement / event reward.
- Dungeon completion (dungeon token → vendor).

## 6. Sink (Para Çıkışı) Kaynakları
- NPC vendor purchase (potion, food, repair, mount).
- Repair (durability restore).
- Teleport / waygate fees.
- Marketplace listing fee + sale fee.
- Respec fee.
- Guild bank tab unlock.
- Cosmetic transmog fee (Phase 02).
- Legendary crafting materials (Phase 02).

## 7. Temel Kurallar
- **Faucet > Sink olmamalı.** Aylık monitoring; %10 üzeri hesaba sink patch.
- Currency ledger append-only; audit trail her transferi tutar.
- Yeni item ekonomiye girmeden önce **sink hesabı** onaylanmalı (item template review checklist).
- Duplication exploit → auto rollback + geçici trade freeze.
- Gold cap: çok yüksek (bigint safe); oyuncu başına.

## 8. Kullanıcı Akışı
```
Vendor → sell junk → gold gain
Repair NPC → cost gold
Marketplace listing → fee upfront + sale fee
```

## 9. İlişkiler
- [Items](items.md), [Marketplace](marketplace.md), [Professions](professions.md), [Loot System](../systems/loot-system.md), [Quests](quests.md), [Cities](cities.md)

## 10. Gereken Veriler
```
currency_ledger
  id, character_id, currency_type, delta, balance_after
  reason (enum: quest, vendor_sale, market_sale, repair, ...)
  correlation_id (activity/order/transaction)
  timestamp

commodity_price_index (server metric)
  day, commodity_code, avg_price, volume, active_listings
```

## 11. I/O
- Input: purchase, sell, transfer, list, cancel listing.
- Output: ledger event, price index update.

## 12. Progression
Ekonomik ilerleme: net worth; oyuncu her aşamada daha büyük alışverişlere erişir.

## 13. Failure Durumları
- Insufficient funds: transaction reject.
- Duplication attempt: rollback + investigation flag.
- Marketplace overloading: rate limit.

## 14. Edge Cases
- Botting: pattern detection + soft ban.
- Real-money-trading (RMT): audit trail + report system.
- Market crash (event): temp price floor sistem safeguard (opsiyonel, tartışmalı).

## 15. AFK Davranışı
Ekonomi doğrudan AFK etkilenmez; ancak AFK aktiviteleri (mining, hunting) faucet üretir. Marketplace listing süresi 24 saat default (AFK-friendly).

## 16. Server Authoritative Gereksinimleri
- Tüm currency mutation transaction'lı, append-only ledger.
- Idempotency key her transaction'da.
- Auditors: her gün otomatik total supply raporu.

## 17. Genişletilebilirlik
- Region-specific currency (Phase 02 optional).
- Player-run banks (Phase 03 fikir).
- Cross-server auction (Phase 03).

## 18. Örnek Senaryo
Karakter 1 hafta içinde 12k gold kazanır (quest+drop+market), 8k gold harcar (repair, teleport, potion). Net +4k. Guild bank'a %10 katkı.

## 19. Balance Considerations
- **Aylık gold arzı büyümesi ≤ %2** hedefi (steady-state).
- **Temel commodity fiyatları** (potion, iron ore, bread) sabit bir bant içinde tutulur; drift alarm.
- Sink noktaları oyuncu için "gereksiz ceza" hissi vermemeli; **konfor karşılığı** olmalı.
- Yeni content patch'i sink/faucet dengeli getirir.
