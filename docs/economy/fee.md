# Fee (Servis Ücretleri)

## 1. Amaç
Sistemin sink'lerinden biri. Servis (respec, teleport, rename, guild create, repair, mail) karşılığında alınır. `system.fee_sink` hesabına gider.

## 2. Fee Kalemleri (v1 + Faz 02)

| Fee | Currency | Formül | Aktif |
|---|---|---|---|
| **Guild create** | gold | flat 10.000 g | v1 (bkz. [guild/guild.md](../guild/guild.md)) |
| **Character rename** | gold | flat 50 g (v1: bir kez ücretsiz, sonra 50 g) | v1 |
| **Respec** (spec point reset) | gold | seviyeye göre `level² × 100 c` (cap 500 g) | v1 (bkz. [class](../class/README.md)) |
| **Fast travel / teleport** | gold | mesafe bazlı `base + distance × factor` | Phase 02 (bkz. [travel](../travel/README.md)) |
| **Repair** | gold | `∑(item.repair_cost)` | Phase 02 |
| **Mail postage** | gold | `base 5 c + attached_item_count × 1 c + gold_leg × 0.001` | Faz 26 |
| **Guild bank tab unlock** | gold | tab1=25k, tab2=100k, tab3=500k, tab4=2M | v1 (bkz. [guild/bank.md](../guild/bank.md)) |
| **Auction listing** | gold | reserve × 5% | Faz 26 |

## 3. Fee Hesaplama Kuralı

- Formül server-side; client sadece preview görür.
- Preview endpoint (`GET /v1/fee/preview?kind=...`) client'a gösterilecek fiyatı verir.
- Client `POST` isteğinde `expected_fee_copper` parametresi gönderir; server'daki fiyatla mismatch → `E_FEE_MISMATCH` (fiyat tablosu update edildi ise).

## 4. Ledger Şeması

```
kind='fee_<name>' (örn. fee_respec, fee_rename)
entries:
  -amount → character.wallet.gold
  +amount → system.fee_sink
ref_kind='service'
ref_id=<service_ref_uuid>
```

## 5. Reputation İndirimi (Phase 02)

Faction fee indirimi (bkz. [character/reputation.md](../character/reputation.md)):
```
final_fee = round(base_fee × (1 - rep_tier × 0.05))    -- cap max %25
```
Sadece belirli fee türlerine uygulanır (repair, mail postage; guild create ve rename hariç).

## 6. Konfig

```yaml
# config/economy/fee.yml
guild_create_copper: 100000000       # 10.000 g
rename:
  first_free: true
  subsequent_copper: 500000          # 50 g
respec:
  formula: "level*level*100"          # c
  cap_copper: 5000000                # 500 g
repair:
  base_rate_per_durability_missing: 2
mail:
  base_copper: 5
  per_item_copper: 1
  gold_leg_rate: 0.001
```

Hot reload — bkz. [architecture/config.md](../architecture/config.md).

## 7. Metrikler
- `economy_fee_collected_copper_total{fee_kind}`
- `economy_fee_events_total{fee_kind}`
- `economy_fee_preview_mismatch_total{fee_kind}`

## 8. Endpoint

**GET `/v1/fee/preview?kind=respec&character_id=...`** — server-authoritative fiyat preview.
**Aktif fee'nin bulunduğu servis endpoint'i içinde** ödeme birlikte yapılır (respec endpoint'i, guild create endpoint'i vs).

## 9. Cross-Reference
- [transaction.md](transaction.md), [tax.md](tax.md), [gold.md](gold.md)
- [guild/guild.md](../guild/guild.md), [guild/bank.md](../guild/bank.md)
- [class/respec.md](../class/respec.md)
- [travel/README.md](../travel/README.md)
- [character/reputation.md](../character/reputation.md)

## 10. Bekleyen Kararlar
- Rename first-free politikası v1 (Phase 02: cosmetic charge review).
- Fast travel fee eğrisi Phase 02.
- Mail postage detay Faz 26.
- Faction fee indirim cap Phase 02.
