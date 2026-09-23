# Currency Taxonomy

## 1. Amaç
Oyunda birden fazla currency türü olacak. Bunları homojen bir ledger üzerinden yönetmek için taksonomi ve kısıtlar burada tanımlıdır.

## 2. Currency Türleri (v1)

| Code | Adı | Base Unit | Trade | Cap | Kaynak |
|---|---|---|---|---|---|
| `gold` | Gold | copper (10⁴) | ✅ | 10¹² copper | Quest, loot, vendor sell |
| `gem` | Gem (premium) | 1 | ❌ | 10⁷ | Store (Phase 02) |
| `honor` | Honor (PvP) | 1 | ❌ | 10⁵ | PvP kill/win (Phase 02) |
| `faction_<code>` | Faction token | 1 | ❌ | 10⁵ per faction | Faction quest |
| `event_<code>` | Event token | 1 | ❌ | 10⁵ per event | Event participation |
| `guild_contribution` | Guild contribution | 1 | ❌ | 10⁶ per member | Guild activity (guild-scoped) |

**v1'de sadece `gold` trade edilebilir.** Diğer currency'ler **Bind-on-Account** (BoA) — hesaba bağlı, karakter arası (Phase 02+ opt) veya kesin karakter-özel.

## 3. Model

```sql
CREATE TABLE economy.currency (
  code VARCHAR(48) PRIMARY KEY,
  display_name_key VARCHAR(64) NOT NULL,
  scope VARCHAR(16) NOT NULL,        -- 'character' | 'account' | 'guild'
  tradable BOOLEAN NOT NULL,
  cap BIGINT NOT NULL,
  base_unit_ratio INT NOT NULL DEFAULT 1,  -- copper için 10000, diğerleri 1
  version INT NOT NULL DEFAULT 1
);

CREATE TABLE economy.character_currency_balance (
  character_id UUID NOT NULL,
  currency_code VARCHAR(48) NOT NULL,
  amount BIGINT NOT NULL DEFAULT 0 CHECK (amount >= 0),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  version INT NOT NULL DEFAULT 1,
  PRIMARY KEY (character_id, currency_code)
);

CREATE TABLE economy.account_currency_balance (
  account_id UUID NOT NULL,
  currency_code VARCHAR(48) NOT NULL,
  amount BIGINT NOT NULL DEFAULT 0 CHECK (amount >= 0),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  version INT NOT NULL DEFAULT 1,
  PRIMARY KEY (account_id, currency_code)
);
```

**Not:** `gold` özel bir currency olarak `economy.character_wallet.gold_copper` alanında da tutulur (hot path). Ledger her iki tabloyu tek transaction'da senkron tutar (idempotent projeksiyon).

## 4. Scope Kuralları

- `scope='character'` → `character_currency_balance` veya `character_wallet`.
- `scope='account'` → `account_currency_balance` (tüm karakterler paylaşır).
- `scope='guild'` → `guild.treasury` (bkz. [guild/bank.md](../guild/bank.md)).

Karakter transferi yalnızca `scope='character'` ve `tradable=true` currency için mümkündür → **v1'de sadece gold**.

## 5. Currency Conversion (v1: YOK)
V1'de her currency izole. Gem → gold gibi otomatik dönüşüm **yok**. Gerekçe: sink/source dengesi bozulur.

Phase 02'de sınırlı vendor exchange (event token → altın veya kozmetik) düşünülebilir; her exchange yolu explicit karar (DD-###) gerektirir.

## 6. Cross-Reference
- [gold.md](gold.md), [transaction.md](transaction.md)
- [guild/bank.md](../guild/bank.md)
- [character/reputation.md](../character/reputation.md) — faction seviye ↔ token ilişkisi

## 7. Bekleyen Kararlar
- Honor Phase 02'de aktif.
- Faction token listesi (Phase 02 faction pass).
- Event token TTL — event sonrası ne kadar geçerli (varsayılan +14 gün).
- Store gem SKU → currency map (Phase 02 monetization).
