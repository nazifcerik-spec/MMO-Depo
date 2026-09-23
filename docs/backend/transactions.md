# Transactions

## 1. Amaç
Postgres transaction'larını disipline etmek: isolation seçimi, retry, cross-store atomicity (2PC yok — outbox pattern), event yayın timing'i.

## 2. Sözleşme
```go
package db

type TxOptions struct {
    Iso         pgx.TxIsoLevel  // ReadCommitted (default), Serializable
    AccessMode  pgx.TxAccessMode // ReadWrite (default), ReadOnly
    DeferrableMode pgx.TxDeferrableMode
}

type TX interface {
    // pgx.Tx-compatible for sqlc Querier
}

func (p *Pool) WithTx(ctx context.Context, opts TxOptions, fn func(TX) error) error
```

Otomatik retry serializable failures için:
```go
func (p *Pool) WithSerializableTx(ctx context.Context, fn func(TX) error) error {
    return retrySerializable(ctx, func() error {
        return p.WithTx(ctx, TxOptions{Iso: pgx.Serializable}, fn)
    })
}
```

## 3. Isolation Level Politikası

| Use Case | Isolation | Neden |
|---|---|---|
| Basit read (character sheet) | Read Committed | Default; hızlı |
| Character mutation (level up, allocate stats) | Read Committed + optimistic lock | Version alanı çakışmayı yakalar |
| Marketplace buy / sell | **Serializable** | Currency + item + listing state atomik |
| Currency ledger insert | **Serializable** | Balance drift ihlali önleme |
| Guild bank withdraw | **Serializable** | Race condition riski |
| Activity contract change | Read Committed + partial unique | 1 primary contract garantisi DB constraint'le |
| Reveal write | Read Committed + idempotency | Idempotency key duplicate önler |
| Content read | Read Committed | Cache-friendly |

## 4. WithTx Örneği
```go
err := pool.WithSerializableTx(ctx, func(tx db.TX) error {
    q := character.NewQuerier(tx)
    econ := economy.NewQuerier(tx)

    // 1. Load listing FOR UPDATE
    l, err := econ.MarketListingLock(ctx, listingID)
    if err != nil { return err }
    if l.Status != "active" { return economy.ErrListingUnavailable }

    // 2. Buyer wallet debit (ledger append)
    if err := econ.LedgerInsert(ctx, ...); err != nil { return err }

    // 3. Update listing → sold
    n, err := econ.MarketListingMarkSold(ctx, listingID, l.Version)
    if err != nil { return err }
    if n == 0 { return errs.ErrConflict }

    // 4. Insert transaction
    if err := econ.MarketTransactionInsert(ctx, ...); err != nil { return err }

    // 5. Item ownership transfer
    if err := econ.ItemOwnerUpdate(ctx, itemInstanceID, "mailbox", mailID); err != nil { return err }

    // 6. Mail create
    if err := econ.MailInsert(ctx, ...); err != nil { return err }

    // 7. Seller ledger credit
    if err := econ.LedgerInsert(ctx, ...); err != nil { return err }

    // 8. Outbox: event 'market.transaction' yazılır
    if err := outbox.Insert(ctx, tx, event); err != nil { return err }

    return nil
})
```

## 5. Retry Politikası
- **40001** (serialization_failure): 3 retry, jitter'lı exponential (10ms/50ms/250ms).
- **40P01** (deadlock_detected): 2 retry.
- Diğer error → propagate.
- Context cancel → retry durur.

## 6. Timeout Kuralları
- Handler ctx deadline < pool statement_timeout.
- Tx içinde external I/O yasak (HTTP call, Redis busy loop, sleep).
- Tx max süre: 500ms hedef; > 1s warning; > 5s statement_timeout kesme.

## 7. Cross-Store Atomicity (Outbox Pattern)

Postgres + Redis + external notification — 2PC yok. Yerine **outbox**:

1. Tx içinde `audit.event_log` veya dedicated `outbox` tablosuna event insert.
2. Tx commit sonrası worker outbox okur → Redis pub, notification service, external webhook.
3. Idempotency: outbox event id.
4. Retry safe: worker delivery attempt loop.

```go
// Domain service
tx içinde:
  economy tx başarılı → outbox.Insert(event="market.transaction.completed", payload={...})
commit
// Worker
polls system.job_queue veya outbox → publishes to Redis pub/sub
```

## 8. Optimistic Locking Pattern
```go
// Repository interface
UpdateOptimistic(ctx, entity, expectedVersion) error  // ErrConflict if version mismatch

// Service level retry
for i := 0; i < 3; i++ {
    e, _ := repo.ByID(ctx, id)
    e.Mutate()
    err := repo.UpdateOptimistic(ctx, e, e.Version)
    if err == nil { return nil }
    if !errors.Is(err, errs.ErrConflict) { return err }
    // reload & retry
}
return errs.ErrConflict
```

## 9. Advisory Locks
- Kısa süreli (migration, scheduled job) için Postgres `pg_advisory_lock`.
- Redis lock (bkz. [Redis](redis.md)) short-hold multi-node için.

## 10. Anti-Patterns
- Tx içinde `time.Sleep` / `http.Get`.
- Tx handler'ın ömrünü aşan (goroutine leak).
- Serializable her yerde (performans katili).
- Optimistic lock retry olmadan üzerine yazmak.
- Outbox atlamak (tx commit + external side-effect direct).

## 11. Test Stratejisi
- Retry policy: table-driven inject 40001 → expected retry count.
- Concurrency test: goroutine race + reconcile assert.
- Golden path integration: full buy flow, ledger reconciled.

## 12. İlişkiler
- [Database](database.md), [Redis](redis.md), [Event Bus](event-bus.md), [Error System](error-system.md), [architecture/postgresql.md](../architecture/postgresql.md), [database/economy-market.md](../database/economy-market.md)

## 13. Bekleyen Kararlar
- Outbox tablosu ayrı mı `audit.event_log` yeter mi — öneri: v1'de `audit.event_log` yeter; v2 dedicated `outbox`.
- Deferrable transaction kullanımı (nadir; long analytical read).
- CDC (change data capture) benimseme (Phase 03).
