# Economy (Ledger + Wallet + Sink/Source)

## 1. Amaç
Ekonominin sağlığını izlemek, drift teşhisi, reversal aksiyon.

## 2. Dashboard

Üst özet:
- **Total gold in circulation** (character + guild wallet toplam).
- **24h source total** (quest, loot, vendor sell).
- **24h sink total** (vendor buy, tax, fee, market cut).
- **Sink/source ratio** (hedef 0.85-0.95).
- **Active listings gold value**.
- **Bank total gold**.

Chart: gold circulation 30d line + sink/source stacked bar.

## 3. Wallet View

Character wallet:
- Search: character_id.
- Columns: character · wallet_copper · last_ledger_ts · reconciled_at.
- Drift flag: `SUM(ledger.delta) - wallet_copper != 0` → red.

Guild treasury benzer.

## 4. Ledger View

Query builder:
```
kind: [any / market_buy / market_listing_fee / vendor_sell / …]
currency: gold / event_<code> / ...
account_kind: character / account / guild / system
account_ref: <uuid> (opsiyonel)
from / to date
```

Result grid: entry_id · txn_id · kind · account_kind · account_ref · delta · balance_after · created_at · ref_kind · ref_id.

Detail: txn compound (aynı txn_id tüm entries; sum=0 assertion visible).

## 5. Reversal

Kritik aksiyon:
```
Reversal Modal
──────────────
Original txn: txn_abc12…
Kind: market_buy
Entries:
  buyer  -10,000c
  seller +9,500c
  tax_sink +500c
──────────────
Reason: (required 32+ char)

[Confirm Step-Up]
```

Confirm → `POST /internal/ops/economy/reversal` (bkz. [economy/transaction.md](../economy/transaction.md) §Reversal). SLA + support runbook Phase 02.

Rate limit: reversal max 5/gün per operator (aşım = admin escalation).

## 6. Reconciliation

- Nightly worker `economy_reconciler` sonucu dashboard'da.
- Drift alert:
  - < 1c: warning.
  - ≥ 1c: character freeze + P0 page + admin dashboard highlight.
- Manual reconcile trigger: `POST /internal/ops/economy/reconcile` (admin).

## 7. Currency Balance

Account currency (gem/honor/event_<code>/faction_<code>) her currency için toplam supply + kullanıcı başına ortalama.

## 8. Endpoint

- `GET /internal/admin/economy/dashboard`
- `GET /internal/admin/economy/wallet?character_id=`
- `GET /internal/admin/economy/ledger?…`
- `GET /internal/admin/economy/ledger/{txn_id}` (compound)
- `POST /internal/ops/economy/reversal` (step-up)
- `POST /internal/ops/economy/reconcile` (admin)

## 9. Metrikler
- `admin_economy_ledger_query_total`
- `admin_economy_reversal_total{operator}`
- `admin_economy_drift_alert_total{severity}`

## 10. Cross-Reference
- [economy/transaction.md](../economy/transaction.md), [economy/gold.md](../economy/gold.md)
- [database/audit.md](../database/audit.md)

## 11. Bekleyen Kararlar
- Reversal SLA + support runbook Phase 02.
- Dashboard aggregate live vs 5dk cache (v1: 5dk cache).
- Multi-currency dashboard tab Phase 02.
