# Vendor (NPC Merchant)

## 1. Amaç
NPC üzerinden gold ↔ item değişimi. Alt yapı; player driven market (auction, trade) yanında **stabil dip fiyat / stabil sink** görevi görür.

## 2. Vendor Türleri

| Type | Açıklama |
|---|---|
| `general` | Genel emtia (yiyecek, potion, mat) |
| `weapon` / `armor` | Ekipman satar |
| `craft_supply` | Meslek malzemesi (bkz. [professions](../professions/README.md)) |
| `special` | Faction / event / reputation gated |
| `buyback_only` | Sadece geri alım (özel event) |

## 3. Model

```sql
CREATE TABLE content.vendor_template (
  id UUID PRIMARY KEY,
  code VARCHAR(48) UNIQUE NOT NULL,
  display_name_key VARCHAR(64) NOT NULL,
  vendor_type VARCHAR(24) NOT NULL,
  location_id UUID NOT NULL,               -- world/location
  faction_code VARCHAR(48) NULL,
  min_reputation SMALLINT NULL,
  buy_multiplier NUMERIC(5,3) DEFAULT 1.0,  -- vendor→player alışta çarpan
  sell_multiplier NUMERIC(5,3) DEFAULT 0.25,-- player→vendor satışta çarpan (sink)
  version INT DEFAULT 1
);

CREATE TABLE content.vendor_inventory (
  vendor_template_id UUID NOT NULL,
  item_template_id UUID NOT NULL,
  price_override_copper BIGINT NULL,       -- boş → item.base_price * buy_multiplier
  stock_kind VARCHAR(16) NOT NULL DEFAULT 'infinite', -- 'infinite'|'daily'|'weekly'
  stock_capacity INT NULL,
  currency_code VARCHAR(48) NOT NULL DEFAULT 'gold',
  PRIMARY KEY (vendor_template_id, item_template_id)
);

CREATE TABLE economy.vendor_buyback (
  id UUID PRIMARY KEY,
  character_id UUID NOT NULL,
  vendor_template_id UUID NOT NULL,
  item_instance_snapshot JSONB NOT NULL,   -- rebuy için dondurulmuş instance
  sold_price_copper BIGINT NOT NULL,
  sold_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  expires_at TIMESTAMPTZ NOT NULL          -- default now()+2h
);
```

## 4. Fiyat Formülü

**Buy (vendor→player):**
```
final_price = round(item.base_price_copper × vendor.buy_multiplier × reputation_discount × affix_multiplier)
reputation_discount = clamp(1 - repTier × 0.02, 0.85, 1.0)
```

**Sell (player→vendor):**
```
final_price = round(item.base_price_copper × vendor.sell_multiplier × rarity_factor)
rarity_factor: common=1.0, uncommon=1.2, rare=1.5, epic=1.8, legendary=2.0
```

Bind-on-Pickup (BoP) item satılabilir; auction'a gitmesin diye vendor sink ana yol. Bound/quest/soulbound item satılamaz (`is_no_sell` flag).

## 5. Akış — Player → Vendor Alış

```
POST /v1/vendor/{vendor_template_code}/buy
  body: {
    character_id,
    item_template_code,
    quantity,
    idempotency_key,
    expected_price_copper  -- client tarafından okuduğu fiyat (price race guard)
  }
```

Server:
1. Vendor aktif mi, karakter menzilde mi (v1 herhangi bir yer, Phase 02 location-gated).
2. `expected_price_copper == server_price` doğrula (aksi hâlde `E_PRICE_MISMATCH`).
3. Karakter cüzdanı yeterli mi (`gold ≥ price × qty`).
4. Envanter slot / weight yeterli mi.
5. Stock (daily/weekly) yeterli mi.
6. Tek TX: wallet-debit + inventory-insert + stock-decrement + ledger + audit.

## 6. Akış — Player → Vendor Satış

```
POST /v1/vendor/{vendor_template_code}/sell
  body: {
    character_id,
    item_instance_ids: [uuid, ...],
    idempotency_key
  }
```

Server:
1. Her item için lock (bkz. [inventory/item-lock.md](../inventory/item-lock.md)).
2. `is_no_sell=false` doğrula.
3. Satış fiyatı hesapla.
4. Tek TX:
   - inventory-remove
   - wallet-credit
   - `economy.vendor_buyback` insert (2h TTL)
   - ledger + audit
5. Response: total_credited + buyback_id list.

## 7. Buyback

2 saat içinde son satılan itemler geri alınabilir:

```
POST /v1/vendor/{vendor_template_code}/buyback
  body: { buyback_id, idempotency_key }
```

- Fiyat = satışta alınan `sold_price_copper` (indirim/artış yok).
- Buyback expired → item silinir (persist).

**Not:** buyback snapshot instance'ı, orijinal instance UUID'sini yeniden kullanamaz (yeni UUID); item history ledger'da linkli.

## 8. Repair (İlişki)
Bazı vendor'lar `repair` özelliğine sahiptir (`vendor_template.can_repair=true`). Bkz. `item-engine/durability.md` (Phase 02).

## 9. Endpoint

**GET `/v1/vendor/{code}`** — vendor inventory + fiyat listesi.
**GET `/v1/vendor/{code}/stock`** — daily/weekly stock durumu.
**POST `/v1/vendor/{code}/buy`**
**POST `/v1/vendor/{code}/sell`**
**GET `/v1/character/{id}/buyback`** — aktif buyback list.
**POST `/v1/vendor/{code}/buyback`**

## 10. Metrikler
- `vendor_buy_total{vendor_code}` / `vendor_sell_total`
- `vendor_gold_sink_copper` / `vendor_gold_source_copper`
- `vendor_price_mismatch_total`
- `vendor_buyback_used_total`

## 11. Cross-Reference
- [gold.md](gold.md), [transaction.md](transaction.md), [tax.md](tax.md)
- [item-engine/item-template.md](../item-engine/item-template.md) — `base_price_copper`
- [character/reputation.md](../character/reputation.md) — reputation discount

## 12. Bekleyen Kararlar
- Location-gated vendor (menzil kısıtı) — Phase 02.
- Dynamic pricing (arz/talep) — Phase 03; v1 statik.
- Vendor stock global mi karakter-özel mi (Phase 02: daily-limited stock karakter-özel).
- Repair mekaniği detay — item-engine/durability.md.
