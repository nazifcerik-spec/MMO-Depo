# Listing (Order Book Entry)

## 1. Amaç
Bir oyuncunun bir itemi belirli fiyattan satışa çıkardığı, süresi olan, sunucu-otoriter bir kayıt.

## 2. Yaşam Döngüsü

```
draft (client-only, sunucuya gitmedi)
  → active (post edildi, escrow'da; alım için görünür)
    ├── sold        (alıcı bulundu; item alıcıya, gold satıcıya)
    ├── expired     (süre doldu; item satıcıya iade, listing fee YANMIŞ)
    └── cancelled   (satıcı iptal; item satıcıya iade; fee iade politikası bkz. cancel.md)
```

Terminal state'ler (`sold`, `expired`, `cancelled`) update edilmez; audit için 90 gün hot, sonrası cold.

## 3. Model

```sql
CREATE TABLE market.listing (
  id UUID PRIMARY KEY,
  seller_character_id UUID NOT NULL,
  seller_account_id UUID NOT NULL,             -- same-account guard için denorm
  item_instance_id UUID NULL,                  -- unique/non-stackable için
  item_template_id UUID NOT NULL,              -- search index anchor
  stack_quantity INT NOT NULL DEFAULT 1,       -- stackable için >1
  unit_price_copper BIGINT NOT NULL CHECK (unit_price_copper > 0),
  total_price_copper BIGINT GENERATED ALWAYS AS (unit_price_copper * stack_quantity) STORED,
  currency_code VARCHAR(48) NOT NULL DEFAULT 'gold',
  state VARCHAR(16) NOT NULL DEFAULT 'active',
  listed_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  expires_at TIMESTAMPTZ NOT NULL,
  closed_at TIMESTAMPTZ NULL,
  close_reason VARCHAR(24) NULL,               -- sold|expired|cancelled|admin
  listing_fee_paid_copper BIGINT NOT NULL DEFAULT 0,
  version INT NOT NULL DEFAULT 1
);

CREATE INDEX idx_market_listing_search ON market.listing (item_template_id, state, unit_price_copper)
  WHERE state = 'active';
CREATE INDEX idx_market_listing_seller ON market.listing (seller_character_id, state);
CREATE INDEX idx_market_listing_expires ON market.listing (expires_at) WHERE state = 'active';

CREATE TABLE market.listing_snapshot (
  listing_id UUID PRIMARY KEY,
  item_snapshot JSONB NOT NULL,   -- affix/rarity/socket/gem/enchant snapshot (search + display)
  captured_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
```

**Neden snapshot?** Item instance state (enchant, socket, gem) listing sırasında dondurulur; alıcı tam olarak gördüğü item'ı alır. Escrow item'ın state'i teoride değişmez ama snapshot search index + gösterim için primary source.

## 4. Escrow Modeli

Post sırasında item:
- `inventory.item_instance.owner_character_id = MARKET_ESCROW_UUID` (sabit sistem UUID'si) veya
- alternatif: `market.escrow_item (listing_id UUID PK, item_instance_id UUID NOT NULL UNIQUE)` join tablosu + item_instance `is_locked=true` + `lock_owner='market'`.

**Karar (v1):** join tablosu + item lock (bkz. [inventory/item-lock.md](../inventory/item-lock.md)); ownership değişmez, `is_locked=true` + `lock_owner='market'` + `locked_ref_id=listing_id`. Satıcı envanteri UI'da item "listed" badge ile grey görünür ama filter'da hariç.

```sql
CREATE TABLE market.escrow_item (
  listing_id UUID PRIMARY KEY,
  item_instance_id UUID NOT NULL UNIQUE,
  seller_character_id UUID NOT NULL,
  escrowed_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
```

Stackable itemler için: satıcı envanterinden `stack_quantity` düşülür ve escrow tablosunda `stack_reserved` sayısı tutulur (item_instance yok — stackable template referansı yeterli; ayrı tablo veya listing satırından türetilir).

## 5. Süre Seçenekleri (v1)

| Süre | Listing Fee Katsayısı |
|---|---|
| 12h | 1.0× |
| 24h | 1.5× |
| 48h | 2.5× |
| 168h (7 gün) | 5.0× |

Listing fee formülü [tax.md](tax.md) §Listing Fee'de.

## 6. Aktif Listing Cap

- Karakter başına max **50 aktif listing** (v1).
- Cap dolduğunda `E_MARKET_LISTING_CAP` → yeni post reddi.
- Cap sold/expired/cancel sonrası boşalır; anında sayılır.

## 7. Cross-Reference
- [sell.md](sell.md) — post akış
- [buy.md](buy.md) — instant buy akış
- [cancel.md](cancel.md) — cancel akış
- [tax.md](tax.md) — fee/cut
- [inventory/item-lock.md](../inventory/item-lock.md)
- [database/economy.md](../database/economy.md)

## 8. Bekleyen Kararlar
- Listing cap 50 v1 (Phase 02: guild tier / reputation ile artabilir).
- Süre yelpazesi (12/24/48/168) v1; Phase 02'de 336h (14d) opt.
- Snapshot vs live item state (v1: snapshot + escrow lock; değişebilirlik yok).
