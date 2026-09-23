# Transaction Validation

## 1. Amaç
Ledger transaction'larının atomik, tutarlı, idempotent, konsolide olması. Ekonomik durumun her koşulda **kaynak-hedef sıfır toplam** kuralına uyması.

## 2. Prensipler (DD-135 uzantısı)

- Tüm ledger yazımı `TxnService.Write` — domain modülleri direkt SQL yazamaz.
- **SERIALIZABLE isolation** — race condition tolere edilmez.
- **Double-entry sum=0** her currency için invariant (`sum(delta per currency) == 0`).
- **Idempotency**: `idempotency_key` UNIQUE constraint (bkz. [replay-protection.md](replay-protection.md)).
- **Optimistic locking** (`version INT`) wallet/balance projeksiyonu için.
- **Retry policy**: `serialization_failure` → exp backoff max 3 (aksi 500 + client retry).

## 3. Assertion'lar (Server-Side)

Her TX içinde:

1. **Idempotency lookup** — daha önce yazılmış mı → cached result.
2. **Actor authorization** — endpoint permission + character ownership.
3. **Business rule** — vendor buy fiyat mismatch, listing state active, envanter kapasite vs.
4. **Balance sufficiency** — çıkış ≥ 0 (sistem hesabı hariç).
5. **Currency validity** — currency_code var ve tradable/scope uyumu.
6. **Sum invariant** — asset conservation.
7. **Cap check** — wallet cap, currency cap.
8. **Version check** — projeksiyon version match; mismatch → serialization retry.

Ret → domain-specific error kodu (`E_INSUFFICIENT_BALANCE`, `E_ITEM_NOT_OWNED`, `E_PRICE_MISMATCH`).

## 4. Cross-Store Atomicity

Ledger + Redis (WS notify, cache invalidate) + external side effect için **outbox pattern**:
1. TX içinde `outbox_event` insert.
2. COMMIT.
3. Worker outbox → pubsub / cache invalidate.
4. Outbox marked processed.

Böylece: DB commit → tüm side effect kesin gerçekleşir; DB rollback → hiçbiri.

Bkz. [backend/transactions.md](../backend/transactions.md).

## 5. Reversal Kanalı

Kullanıcı reversal yolu yok; sadece admin `/internal/ops/economy/reversal` (bkz. [economy/transaction.md](../economy/transaction.md) + [admin/economy.md](../admin/economy.md)).

## 6. Reconciliation

- Nightly `economy_reconciler` worker.
- Drift ≥ 1c → character freeze + P0 page + admin dashboard.
- Rebuild path: ledger stream'den projeksiyon regenerate (event sourcing partial).

## 7. Fraud Signal (Anti-Cheat)

- **Impossibly rapid earn**: 1 saatte 10x sink limit gold earn → flag.
- **Zero-loss trade cycle**: A→B→C→A gold döngüsü → wash trade flag.
- **Ledger replay attempt**: aynı idempotency_key farklı content ile geldi → security alert.
- **Sum violation attempt**: sum≠0 write (bug veya intrusion) → P0 alert + rollback.

Bu sinyaller `audit.security_event` + Anomaly detection (Phase 03 ML).

## 8. Endpoint Katalog (Ekonomiye Etki Eden)

Her biri TxnService kullanmalı:
- vendor/*/buy, vendor/*/sell
- market/listing (post), listing/{id}/buy, listing/{id}/cancel (fee refund), listing expire (worker)
- trade/{session_id}/confirm (finalize)
- guild/*/bank/deposit, withdraw
- guild/*/boss/summon (treasury cost)
- fee/* (respec, rename, guild create, teleport, mail postage, repair)
- reward-engine grant (quest, loot, event token)

Non-ledger domain (character mutation gibi) TxnService kullanmayabilir ama audit yazar.

## 9. Metrikler
- `sec_ledger_serialization_retry_total`
- `sec_ledger_sum_violation_total` (P0 alert)
- `sec_ledger_idempotency_replay_total`
- `sec_reconcile_drift_total{severity}`

## 10. Cross-Reference
- [economy/transaction.md](../economy/transaction.md), [backend/transactions.md](../backend/transactions.md)
- [replay-protection.md](replay-protection.md), [audit.md](audit.md)

## 11. Bekleyen Kararlar
- Wash trade detection algoritma Phase 02.
- ML anomaly detection Phase 03.
- Reversal SLA runbook Phase 02.
