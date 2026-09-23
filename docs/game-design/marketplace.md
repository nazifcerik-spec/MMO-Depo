# Marketplace (Pazar / Auction House)

## 1. Amaç
Oyuncular arası eşya ve materyal ticaretini güvenli, ölçeklenebilir bir platformda buluşturmak.

## 2. Oyuncuya Fayda
- Kazandıklarını likit hale çevirme.
- Kısa yoldan aradığı ekipmanı bulma.
- Uzun soluklu ekonomi oyunu (arbitrage, crafting business).

## 3. Gameplay Loop
- Session'da listing kur, satılanları topla.
- Meta'da fiyat trend takibi, tedarik zinciri planlaması.

## 4. Temel Kurallar
- Listing süresi: 12h / 24h / 48h; longer = higher listing fee.
- Fee model: listing fee (upfront, non-refundable) + sale fee (satışta kesinti).
- Fee önerisi: listing %2, sale %5 (Phase 01 balance test).
- Bought items → mailbox'a (safe delivery).
- Cancel listing: item mailbox'a döner; listing fee yanmaz (opsiyonel: %50 iade).
- Global market (kıta bazlı, tek dünya). Region market ileride.

## 5. Kullanıcı Akışı
```
Market Panel → Search / Sell tabs
 Search:
   filter (category, rarity, level, class, price, stat)
   Sort (price, unit price, expiry)
   Buy → mailbox
 Sell:
   pick item from inventory → set price/qty/duration
   pay listing fee → live
```

## 6. İlişkiler
- [Economy](economy.md), [Items](items.md), [Inventory](inventory.md), [Notification System](../systems/notification-system.md), [Cities](cities.md)

## 7. Gereken Veriler
```
market_listing
  id, seller_char_id, item_instance_id, qty, unit_price
  listing_fee_paid, sale_fee_pct
  created_at, expires_at, status (active/sold/expired/cancelled)

market_transaction
  id, buyer_char_id, seller_char_id, listing_id
  unit_price, qty, total_gross, sale_fee, net_to_seller
  executed_at

market_watchlist (Phase 02)
  character_id, item_template_id, max_price
```

## 8. I/O
- Input: create/cancel listing, buy.
- Output: mailbox delivery, listing status.

## 9. Progression
Personal reputation (opsiyonel Phase 02): iyi seller "trusted" olur, fee'de küçük indirim.

## 10. Reward
Trade fiyat farkı; kişisel ekonomik ilerleme.

## 11. Failure Durumları
- Insufficient funds/inventory: reject.
- Race condition (aynı listing iki alıcı): ilk transaction kazanır; ikinci "already sold" mesajı.
- Listing expire: item mailbox'a.

## 12. Edge Cases
- Underpricing exploit: fiyat < X ise auto-flag (opsiyonel).
- Market manipulation (fake listings): audit + rate limit + investigation.
- Price gouging (rare item single seller): "ask price ↑" sinyali için community info panel.

## 13. AFK Davranışı
Listing 24 saat default; oyuncu offline iken satılırsa mailbox'a düşer. Login'de reveal.

## 14. Server Authoritative Gereksinimleri
- Transaction atomic (buyer gold−, seller mailbox+, item ownership transfer).
- Concurrent buyers: optimistic locking.
- Fee calculation server.

## 15. Genişletilebilirlik
- Watchlist alerts (Phase 02).
- Bid/auction format (Phase 02 rare drop).
- Cross-region arbitrage system (Phase 03).

## 16. Örnek Senaryo
Karakter, rare item'ı 24h listing 5000g fiyattan koyar (100g listing fee); 8 saat sonra satış → 5000g − 250g sale fee = 4750g mailbox'a.

## 17. Balance Considerations
- Fee toplamı %7 civarı: yumuşak sink, spam engelleyici.
- Listing süresi çeşitliliği zaman-para trade-off.
- UI iyileştirmesi ekonomiye sağlıklı katılım için kritiktir; power-user features (bulk list, saved searches) Phase 02.
