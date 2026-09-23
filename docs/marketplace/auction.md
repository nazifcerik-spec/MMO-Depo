# Auction (Bid-Based)

> ⚠️ **v1**: Kapsam dışı (öneri). **Phase 02 primary.**
> V1'de sadece fixed-price listing (bkz. [listing.md](listing.md)). Auction burada tasarım hazır.

## 1. Amaç
Rare / prestige item için piyasanın fiyat keşfetmesine izin veren, zaman-sınırlı, artışlı teklif akışı.

## 2. Karakteristik

- **Reserve price** — minimum kabul fiyatı; ulaşılmazsa satılmaz.
- **Buyout price (opsiyonel)** — instant close fiyatı.
- **Bid increment** — min artış (yüzde veya flat).
- **Anti-snipe** — son 5 dk içinde bid → 5 dk süre uzatımı (max toplam +30 dk).
- **Bid escrow** — bidder'ın gold'u bid anında escrow'a alınır (aşırı bid engeli).
- **Outbid** — daha yüksek bid gelince önceki bidder'a gold otomatik iade.

## 3. Model

```sql
CREATE TABLE market.auction (
  id UUID PRIMARY KEY,
  seller_character_id UUID NOT NULL,
  seller_account_id UUID NOT NULL,
  item_instance_id UUID NULL,
  item_template_id UUID NOT NULL,
  stack_quantity INT NOT NULL DEFAULT 1,
  reserve_price_copper BIGINT NOT NULL,
  buyout_price_copper BIGINT NULL,
  min_bid_increment_copper BIGINT NOT NULL,     -- veya %
  current_bid_copper BIGINT NULL,
  current_bidder_character_id UUID NULL,
  bid_count INT NOT NULL DEFAULT 0,
  started_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  ends_at TIMESTAMPTZ NOT NULL,
  original_ends_at TIMESTAMPTZ NOT NULL,        -- anti-snipe extend base
  extension_count SMALLINT NOT NULL DEFAULT 0,
  max_extensions SMALLINT NOT NULL DEFAULT 6,   -- 6 × 5min = 30 dk cap
  state VARCHAR(16) NOT NULL DEFAULT 'active',
  -- active | sold | unsold_reserve_not_met | cancelled | expired
  closed_at TIMESTAMPTZ NULL,
  version INT NOT NULL DEFAULT 1
);

CREATE TABLE market.auction_bid (
  id BIGSERIAL PRIMARY KEY,
  auction_id UUID NOT NULL,
  bidder_character_id UUID NOT NULL,
  bidder_account_id UUID NOT NULL,
  bid_copper BIGINT NOT NULL,
  bid_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  outbid_at TIMESTAMPTZ NULL,                   -- refund süresi
  refund_txn_id UUID NULL
);
```

## 4. Akış — Post Auction

**POST `/v1/market/auction`**
```json
{
  "item_instance_id": "...",
  "reserve_price_copper": 50000,
  "buyout_price_copper": 200000,
  "duration_hours": 24,
  "min_bid_increment_copper": 500
}
```

Server: fixed-listing sell akışına benzer (escrow lock, listing fee tahsili, cap kontrolü). Auction fee formülü aynı listing fee ama +50% (bid overhead) — Phase 02'de tune.

## 5. Akış — Bid

**POST `/v1/market/auction/{id}/bid`**
```json
{
  "bid_copper": 55000,
  "auto_refund_previous": true,
  "idempotency_key": "..."
}
```

Server (SERIALIZABLE):
1. Auction lock; state='active'; ends_at > now().
2. Same-account guard.
3. `bid >= max(reserve, current_bid + min_increment)`.
4. Bidder wallet ≥ bid.
5. Bidder gold escrow'a (ledger: `-bid → bidder, +bid → system.auction_escrow`).
6. Önceki bidder varsa iade (ledger: `-previous_bid → system.auction_escrow, +previous_bid → previous_bidder`).
7. `auction.current_bid=bid, current_bidder=bidder, bid_count++, updated_at=now()`.
8. Anti-snipe: `now() > ends_at - 5min` ve `extension_count < max_extensions` → `ends_at = now() + 5min; extension_count++`.
9. `market.auction_bid` insert; önceki bidder'a `outbid_at=now(), refund_txn_id`.
10. Outbox: WS "outbid" bildirim + WS "auction_updated" watchers.

## 6. Akış — Buyout

Buyout price varsa herhangi bir zamanda instant close:

**POST `/v1/market/auction/{id}/buyout`**
- current_bid iade (varsa).
- buyout_price alınır (buyer wallet -→ seller wallet + sale_cut → tax_sink).
- Item transfer; state='sold'; closed_at=now().

## 7. Auto Close (Worker)

Cron `auction_close` her 1 dk:
```sql
SELECT id FROM market.auction
 WHERE state='active' AND ends_at <= now()
 FOR UPDATE SKIP LOCKED
 LIMIT 500;
```

Her batch:
- Reserve karşılandı mı (`current_bid >= reserve_price`)?
  - Evet: state='sold'; ledger settle (escrow'daki gold → seller net + sale cut); item → winner.
  - Hayır: state='unsold_reserve_not_met'; escrow'daki gold varsa (min zorunlu bid ≥ reserve olduğu için current_bid > 0 ise satıldı → burada sadece bid_count=0 keys); item satıcıya iade.

## 8. Cancel Kısıtı

Satıcı auction'ı cancel edebilir **sadece** bid gelmediyse (`bid_count=0`). Bid varsa cancel yok (bidder tricking koruma). Fee iade politikası [cancel.md](cancel.md) ile aynı.

## 9. Endpoint Özet

**POST `/v1/market/auction`** — create.
**POST `/v1/market/auction/{id}/bid`**
**POST `/v1/market/auction/{id}/buyout`**
**DELETE `/v1/market/auction/{id}`** — cancel (bid_count=0 only).
**GET `/v1/market/auction/{id}`** — detail.
**GET `/v1/market/auction/search`** — search (filter benzer).
**GET `/v1/character/{id}/auction/active`** — kişisel aktif auction (seller + as bidder).

## 10. Metrikler
- `auction_created_total`
- `auction_bid_total` / `auction_bid_reject_total{reason}`
- `auction_snipe_extension_total`
- `auction_sold_total`, `auction_unsold_total`
- `auction_gold_escrowed_copper_gauge`

## 11. Cross-Reference
- [listing.md](listing.md), [sell.md](sell.md), [buy.md](buy.md), [tax.md](tax.md)
- [economy/transaction.md](../economy/transaction.md)
- [systems/notification.md](../systems/notification.md) — outbid + won push

## 12. Bekleyen Kararlar
- Auction v1'e alınacak mı (öneri: **Phase 02**; v1 sadece fixed).
- Anti-snipe süre + cap ayarı (5dk × 6 = 30dk) Phase 02.
- Bid escrow ayrı sistem hesabı mı `system.tax_sink` mi (öneri: ayrı `system.auction_escrow`).
- Proxy bidding (auto bid up to X) Phase 03.
- Auction listing fee premium (%50 üstü) Phase 02.
