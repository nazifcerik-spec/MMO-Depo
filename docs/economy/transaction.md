# Transaction Engine (Ledger)

## 1. Amaç
Ekonominin değişmez muhasebe defteri. Her değer transferi burada **çift-taraflı** (double-entry) ve **immutable** olarak yazılır. Ledger'a yazılmadan hiçbir wallet/balance projeksiyonu güncellenmez.

## 2. Prensipler

- **Double-entry**: her value flow bir `debit` (kaynak) + bir `credit` (hedef) çifti. Toplam == 0 (asset conservation).
- **Immutable append-only**: `UPDATE` / `DELETE` yasak; iptal `reversal` entry ile yapılır.
- **SERIALIZABLE isolation**: ledger yazımı ve balance projeksiyonu tek TX'te.
- **Idempotent**: her yazım `idempotency_key`'e bağlı; retry duplicate yazmaz.
- **Audit-first**: her ledger entry `audit.ledger_audit`'a mirror'lanır (aynı TX).
- **Outbox integration**: cross-store side effect (WS notify, worker job) outbox üzerinden (bkz. [backend/transactions.md](../backend/transactions.md)).

## 3. Model

```sql
CREATE TABLE economy.ledger_entry (
  id BIGSERIAL PRIMARY KEY,
  txn_id UUID NOT NULL,                     -- transaction batch id
  entry_seq SMALLINT NOT NULL,              -- txn içi sıra (0..N)
  kind VARCHAR(32) NOT NULL,                -- 'quest_reward'|'vendor_buy'|'trade'|'tax'|'fee'|'reversal'|...
  currency_code VARCHAR(48) NOT NULL,
  account_kind VARCHAR(16) NOT NULL,        -- 'character'|'account'|'guild'|'system'
  account_ref UUID NOT NULL,                -- character_id / account_id / guild_id / system-const-uuid
  delta BIGINT NOT NULL,                    -- + credit, - debit
  balance_after BIGINT NOT NULL,            -- post-entry projected balance (sanity check)
  ref_kind VARCHAR(32) NULL,                -- 'quest'|'trade_session'|'vendor'|'loot'|...
  ref_id UUID NULL,
  idempotency_key VARCHAR(96) NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  UNIQUE (idempotency_key, entry_seq)
);
CREATE INDEX idx_ledger_account ON economy.ledger_entry (account_kind, account_ref, created_at DESC);
CREATE INDEX idx_ledger_txn ON economy.ledger_entry (txn_id);
CREATE INDEX idx_ledger_ref ON economy.ledger_entry (ref_kind, ref_id);

-- Sistem hesabı sabit UUID'leri (constants):
--   'system.vendor_sink'
--   'system.tax_sink'
--   'system.fee_sink'
--   'system.reward_source'
--   'system.loot_source'
--   'system.mint_source'  (event / seasonal)
```

**Not:** `balance_after` alanı ledger içi consistency guard'dır (aynı `(account_kind, account_ref, currency_code)` üzerinden okurken); asıl hakikat entry stream'in kendisidir.

## 4. Servis API

```go
// Domain servisleri yalnızca bu servisi çağırır.
type TxnService interface {
    Write(ctx, intent TxnIntent) (TxnResult, error)
}

type TxnIntent struct {
    Kind            string       // 'vendor_buy' vs
    RefKind         string       // 'vendor'
    RefID           uuid.UUID
    IdempotencyKey  string       // gerekli
    Entries         []Entry      // en az 2 tane; sum(delta per currency) == 0
    ProjectionHints []Projection // wallet/balance update spesifik satırlar
}

type Entry struct {
    Currency    string
    AccountKind string        // 'character'|'account'|'guild'|'system'
    AccountRef  uuid.UUID
    Delta       int64         // + / -
}
```

Servis içinde:
1. Idempotency lookup (`idempotency_key` daha önce yazılmışsa cached result dön).
2. `BEGIN ISOLATION LEVEL SERIALIZABLE`.
3. Her account için mevcut balance oku (`FOR UPDATE`). Sistem hesabı balance kontrolüne dahil değildir (unbounded source/sink).
4. Sum(delta per currency) == 0 assert (asset conservation).
5. Her hedef balance ≥ 0 doğrula (çıkışta yetersiz → `E_INSUFFICIENT_BALANCE`).
6. `economy.ledger_entry` insert (batch).
7. Projeksiyon update:
   - `character` + `gold` → `economy.character_wallet`.
   - `character` + non-gold → `economy.character_currency_balance`.
   - `account` → `economy.account_currency_balance`.
   - `guild` → `guild.treasury` (bkz. [guild/bank.md](../guild/bank.md)).
   - Her projeksiyon UPDATE'i **optimistic lock** (`version`).
8. `audit.ledger_audit` insert.
9. Outbox insert (varsa WS / event bus için).
10. `COMMIT`.
11. Retry: `serialization_failure` → exp backoff 3 kez.

## 5. Reversal (İptal / İade)

Ledger append-only olduğundan iptal = ters yönde yeni entry group:

```
kind='reversal'
ref_kind='ledger_txn'
ref_id=<original_txn_id>
entries: her orijinal entry'nin negatifi
```

Reversal koşulları:
- Trade session cancel/finalize hatası (sistem içi rollback zaten commit'ten önce).
- Support-tool tarafından **operatör onaylı** rollback (audit trail zorunlu).

Kullanıcı UI'sinden reversal yoktur. Yalnızca admin/ops.

## 6. Reconciliation

Nightly worker (`economy_reconciler`):
- Her karakter için `SUM(ledger_entry.delta) per currency` hesapla.
- Wallet/balance projeksiyonu ile karşılaştır.
- Drift varsa:
  - < 1 copper: alarm sadece.
  - ≥ 1 copper: karakter freeze (transaction reddi) + ops paging (`P0`).
- Guild treasury için aynı check.

## 7. Retention

- `economy.ledger_entry`: **12 ay** hot; sonrası cold archive (S3 parquet). Kullanıcı history 12 ay görebilir.
- `audit.ledger_audit`: **24 ay** hot; sonrası cold archive (compliance).
- Reversal referansı için orijinal entry cold'a inmeden reversal olmalıdır.

## 8. Metrikler
- `economy_ledger_entry_total{kind, currency}`
- `economy_ledger_write_latency_ms`
- `economy_serialization_retry_total`
- `economy_ledger_reconcile_drift_total{severity}`
- `economy_ledger_reversal_total{reason}`

## 9. Endpoint (İç / Ops)

**GET `/internal/economy/ledger/{txn_id}`** — txn entries.
**GET `/v1/character/{id}/wallet/history`** — kişisel history (kendi karakteri).
**POST `/internal/ops/economy/reversal`** — op onayı ile reversal (admin JWT + 2FA).

## 10. Cross-Reference
- [backend/transactions.md](../backend/transactions.md), [backend/database.md](../backend/database.md)
- [database/economy.md](../database/economy.md), [database/audit.md](../database/audit.md)
- [gold.md](gold.md), [currency.md](currency.md), [tax.md](tax.md), [fee.md](fee.md)
- [systems/reward-engine.md](../systems/reward-engine.md) — reveal → txn

## 11. Bekleyen Kararlar
- Cold archive format (parquet vs jsonl).
- Multi-currency swap tek txn mı (v1: her currency ayrı invariant, tek TX içinde birden fazla currency block).
- Reversal SLA (operator flow) — Phase 02 support runbook.
