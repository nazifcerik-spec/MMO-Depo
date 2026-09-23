# Gold (Primary Currency)

## 1. Amaç
Oyunun ana değer birimi. Vendor, trade, tax, fee ve player-to-player transactionların referans birimidir.

## 2. Birim ve Skala

- **Base unit: copper (bakır)** — integer, `BIGINT`.
- Görüntüde `1 gold = 100 silver = 10.000 copper`.
- Server ve DB **her zaman copper cinsinden** çalışır; UI dönüşümü client tarafındadır.
- Neden integer copper? Floating point drift + rounding hatası ledger'da tolere edilemez.

**Cap (v1):**
- Karakter cüzdanı: `2^53 - 1` copper (JS güvenli sayı) → uygulama softcap **10^12 copper** (10 milyar gold).
- Aşımı deneyen credit `E_GOLD_CAP_EXCEEDED` ile red; artan miktar mail'e overflow (Faz 26).

## 3. Cüzdan Modeli

```sql
-- economy.character_wallet: hızlı okuma için projeksiyon.
CREATE TABLE economy.character_wallet (
  character_id UUID PRIMARY KEY REFERENCES character.character(id),
  gold_copper BIGINT NOT NULL DEFAULT 0 CHECK (gold_copper >= 0),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  version INT NOT NULL DEFAULT 1
);
```

Cüzdan projeksiyondur; **hakikat kaynağı ledger'dır** (`economy.ledger_entry`). Reconciliation job (nightly) projeksiyonu ledger toplamı ile doğrular.

## 4. Kaynaklar (Source — sistemin gold ürettiği yerler)

| Kaynak | Açıklama | Ölçek |
|---|---|---|
| Quest reward | Quest completion | Küçük–orta |
| Loot drop | Enemy/chest gold pile | Küçük |
| Vendor sell | Player → NPC item satışı | Orta |
| Daily bonus | Faz 02 login bonus | Küçük |
| Event reward | Seasonal event | Orta |

## 5. Sinkler (Sink — sistemin gold yaktığı yerler)

| Sink | Açıklama |
|---|---|
| Vendor buy | NPC'den item alışı |
| Trade tax | P2P trade komisyonu |
| Auction fee | Auction house (Faz 26) |
| Repair | Item durability tamiri |
| Service fee | Teleport, mail, respec, rename, guild create |
| Craft fee | Bazı reçetelerde vendor material cost |

**Denge hedefi:** haftalık `sink / source ≈ 0.85 – 0.95` (Phase 02 balance pass).

## 6. Transfer Operasyonları

Karakter cüzdanı üzerindeki her mutasyon tek bir servise gider:

```
POST /internal/economy/wallet/apply
  body: {
    character_id,
    delta_copper,       -- +/-
    kind,               -- 'quest_reward'|'vendor_sell'|'vendor_buy'|'trade_out'|...
    ref_kind, ref_id,   -- kaynağa referans (quest_id, trade_id, ...)
    idempotency_key
  }
```

Servis içeride ledger yazar + wallet projeksiyonunu günceller (bkz. [transaction.md](transaction.md)).

## 7. Görüntüleme (Client)
- Ana ekranda cüzdan pill: `12.345 g 67 s 89 c`.
- Trade / vendor / auction ekranlarında hep 3 sütun (gold/silver/copper).
- Copper altı için virgül yok; drift olmasın.

## 8. Endpoint (Public)

**GET `/v1/character/{id}/wallet`** — cüzdan snapshot.
**GET `/v1/character/{id}/wallet/history?kind=&from=&to=&cursor=`** — ledger sayfalı okuma (kişisel).

## 9. Cross-Reference
- [transaction.md](transaction.md), [tax.md](tax.md), [fee.md](fee.md), [vendor.md](vendor.md)
- [database/economy.md](../database/economy.md)
- [inventory/item-transfer.md](../inventory/item-transfer.md)

## 10. Bekleyen Kararlar
- Softcap yükseltme skalası (level ile ölçekleme mi, sabit mi).
- Cüzdan overflow → mail item mi, otomatik banka mı (Faz 26).
- Currency exchange yasağı v1 kesin — Phase 02'de premium currency review.
