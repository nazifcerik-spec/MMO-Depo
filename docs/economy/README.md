# Economy Engine — Katalog

> **Faz 25**: Gerçek MMO ekonomisi. Ledger-first, audit-first, server-authoritative.

## 1. Kapsam
Bu klasör, karakterler ve NPC/dünya arasında değer transferini yöneten bileşenleri kapsar. Item ekonomisi (loot, drop, crafting output) ayrı domainlerde tanımlıdır; bu klasör **değer akışını** ve **kayıt zincirini** tarif eder.

## 2. Altın Kural (DD-135)
Her değer transferi (currency, item, XP, reputation) tek bir kanonik yol üzerinden yazılır:

```
Domain Service
  → Transaction Service.WriteEntry(intent)
    → BEGIN (SERIALIZABLE)
    → INSERT economy.ledger_entry (double-entry)
    → UPDATE affected balance rows (optimistic lock)
    → INSERT audit.ledger_audit
    → COMMIT
```

Hiçbir modül `UPDATE character SET gold = gold + N` yazamaz. Yalnızca ledger yazar; balance projeksiyonu ledger'dan türetilir veya idempotent olarak güncellenir.

## 3. Dokümanlar
- [gold.md](gold.md) — Ana currency (copper base unit), cüzdan, cap, transfer.
- [currency.md](currency.md) — Currency taxonomy (gold, gems, honor, faction, event tokens).
- [vendor.md](vendor.md) — NPC vendor buy/sell mekaniği, buyback, fiyat formülü.
- [trade.md](trade.md) — Player-to-player trade window, güvenlik, iptal, audit.
- [transaction.md](transaction.md) — Ledger, double-entry, SERIALIZABLE isolation, outbox.
- [tax.md](tax.md) — Trade / market / auction / vendor spread tax modelleri.
- [fee.md](fee.md) — Service fee (mail postage, respec, teleport, guild create, rename).

## 4. Cross-Reference
- [database/economy.md](../database/economy.md) — SQL şema kaynağı.
- [systems/reward-engine.md](../systems/reward-engine.md) — Reveal → grant hattı.
- [inventory/item-transfer.md](../inventory/item-transfer.md) — Item hareketi.
- [guild/bank.md](../guild/bank.md) — Guild treasury/bank ledger.
- [systems/rng.md](../systems/rng.md) — Deterministic RNG (drop, price variance).

## 5. Metrikler (v1)
- `economy_ledger_entry_total{kind}`
- `economy_ledger_write_latency_ms{kind}`
- `economy_serialization_retry_total{kind}`
- `economy_sink_gold_total{sink}` — vendor sell, tax, fee, repair.
- `economy_source_gold_total{source}` — quest, loot, sale.

## 6. Bekleyen Kararlar
- Sink/source dengesi hedef oranı (Phase 02 balance pass).
- Currency conversion (gems→gold gibi) yasak/serbest — v1 yasak.
- Auction house — Faz 26 ayrı klasör.
