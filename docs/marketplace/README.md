# Marketplace — Katalog

> **Faz 26**: Player Market + Auction. Fixed-price listing (v1) + Bid-based auction (v1 opt / Phase 02).

## 1. Kapsam
Bu klasör oyuncu-oyuncu değer transferini **asenkron** (aynı anda online olmadan) mümkün kılan altyapıyı tarif eder. Trade (P2P eş zamanlı) → [economy/trade.md](../economy/trade.md); Marketplace → **buradaki** asenkron order book.

## 2. Ana Bileşenler

- **Fixed-price listing** (v1 primary) — satıcı fiyat belirler, alıcı instant buy.
- **Auction** (v1 opt / Phase 02) — bid-based; timed close.
- **Escrow** — listing süresince item marketplace escrow'unda; satıcı erişemez.
- **Server-authoritative order matching** — client sadece intent gönderir.
- **Ledger-first** — tüm gold hareketi `economy.ledger_entry` üzerinden (DD-135).
- **Tax & fee** — vendor değil, marketplace kendi listing fee + sale cut (bkz. [tax.md](tax.md)).

## 3. Dokümanlar
- [listing.md](listing.md) — Listing lifecycle (post → active → sold/expired/cancelled).
- [sell.md](sell.md) — Post-to-market akış + item escrow + validasyon.
- [buy.md](buy.md) — Instant buy akış + price race guard + atomic swap.
- [cancel.md](cancel.md) — Iptal akışı + fee iade politikası.
- [search.md](search.md) — Query API + sıralama + sayfalama + cache.
- [filter.md](filter.md) — Filter dimensions (rarity, level, stat, socket, class, price range).
- [price.md](price.md) — Fiyat girişi + suggested price + min/max clamp.
- [tax.md](tax.md) — Listing fee + sale cut + AH sink.
- [auction.md](auction.md) — Bid-based auction (Phase 02 primary) + anti-snipe.
- [market-history.md](market-history.md) — Historical price data + median/percentile.

## 4. Kısıtlar (v1)

- **Global marketplace tek** (dünya çapında; shard/region ayrımı yok v1).
- **Sadece `gold` currency** (DD-137).
- **BoP item listelenemez**; bind rule kontrolü zorunlu.
- **Karakter başına aktif listing cap**: 50 (opt tune).
- **Listing süresi**: 12h / 24h / 48h / 168h (7 gün). Süre uzatıldıkça listing fee artar.
- **Same-account listing kendi karakterine satış**: v1 **engelli** (self-fund koruma; trade'den farklı).
- **Search rate limit**: karakter başına 20 req/dk.
- **Order matching**: instant buy (auction ayrı akış).

## 5. Cross-Reference
- [economy/README.md](../economy/README.md) — Ledger, wallet, tax altyapısı
- [economy/trade.md](../economy/trade.md) — Eş zamanlı P2P alternatif
- [inventory/item-lock.md](../inventory/item-lock.md) — Escrow lock
- [systems/notification.md](../systems/notification.md) — Sold/expired notify
- [database/economy.md](../database/economy.md) — SQL şema

## 6. Metrikler (v1)
- `market_listing_active_gauge`
- `market_listing_posted_total{item_class}`
- `market_listing_sold_total`, `market_listing_expired_total`, `market_listing_cancelled_total`
- `market_search_qps`, `market_search_latency_ms`
- `market_gold_traded_copper_total`
- `market_price_mismatch_total`

## 7. Bekleyen Kararlar
- Auction v1 mi Phase 02 mi (öneri: Phase 02; v1'de sadece fixed-price).
- Same-account listing anti-farm review (Phase 02).
- Global search cache TTL (v1: 30 sn öneri).
- Cross-region marketplace Phase 03.
