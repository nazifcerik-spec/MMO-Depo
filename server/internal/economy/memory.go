package economy

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
)

// MemoryTxnService is an in-memory TxnService for unit tests. It exercises the
// same invariants as PgTxnService (ValidateIntent + cap + non-negative +
// idempotency) but keeps state in maps guarded by a single mutex — no
// SERIALIZABLE semantics, no retry.
//
// Downstream domain packages (vendor, market, quest, activity…) use this in
// their unit tests to verify their ledger-write intents without spinning up a
// database.
type MemoryTxnService struct {
	mu         sync.Mutex
	currencies map[string]Currency
	balances   map[balanceKey]int64 // (kind,ref,currency) → balance
	txnByIdem  map[string]Result    // idempotency_key → cached Result
	entries    []LedgerRow          // append-only history
	nextID     int64
	now        func() time.Time
}

type balanceKey struct {
	Kind     AccountKind
	Ref      uuid.UUID
	Currency string
}

// NewMemoryTxnService returns an empty service seeded with the given currency
// catalog (nil → StaticCurrencies).
func NewMemoryTxnService(currencies map[string]Currency) *MemoryTxnService {
	if currencies == nil {
		currencies = StaticCurrencies()
	}
	return &MemoryTxnService{
		currencies: currencies,
		balances:   make(map[balanceKey]int64),
		txnByIdem:  make(map[string]Result),
		now:        time.Now,
	}
}

// SetClock overrides the clock (test-only).
func (m *MemoryTxnService) SetClock(now func() time.Time) { m.now = now }

func (m *MemoryTxnService) Balance(_ context.Context, kind AccountKind, ref uuid.UUID, currency string) (int64, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.balances[balanceKey{Kind: kind, Ref: ref, Currency: currency}], nil
}

func (m *MemoryTxnService) Write(_ context.Context, intent Intent) (Result, error) {
	if err := ValidateIntent(intent, m.currencies); err != nil {
		return Result{}, err
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	// Idempotency check
	if cached, ok := m.txnByIdem[intent.IdempotencyKey]; ok {
		copy := cached
		copy.Replayed = true
		return copy, nil
	}

	// Prevalidate cap + non-negative for non-system accounts.
	proposed := make(map[balanceKey]int64)
	for _, e := range intent.Entries {
		key := balanceKey{Kind: e.AccountKind, Ref: e.AccountRef, Currency: e.Currency}
		newBal := m.balances[key] + e.Delta
		if e.AccountKind != AccountSystem {
			cur := m.currencies[e.Currency]
			if newBal < 0 {
				return Result{}, fmt.Errorf("%w: currency=%s account=%s new=%d",
					ErrInsufficientBalance, e.Currency, e.AccountRef, newBal)
			}
			if newBal > cur.Cap {
				return Result{}, fmt.Errorf("%w: currency=%s cap=%d new=%d",
					ErrCurrencyCap, e.Currency, cur.Cap, newBal)
			}
		}
		proposed[key] = newBal
	}

	// Apply.
	txnID := uuid.New()
	rows := make([]LedgerRow, 0, len(intent.Entries))
	for i, e := range intent.Entries {
		key := balanceKey{Kind: e.AccountKind, Ref: e.AccountRef, Currency: e.Currency}
		newBal := m.balances[key] + e.Delta
		m.balances[key] = newBal

		m.nextID++
		var refIDPtr *uuid.UUID
		if intent.RefID != uuid.Nil {
			rid := intent.RefID
			refIDPtr = &rid
		}
		row := LedgerRow{
			ID:           m.nextID,
			TxnID:        txnID,
			EntrySeq:     int16(i),
			Kind:         intent.Kind,
			Currency:     e.Currency,
			AccountKind:  e.AccountKind,
			AccountRef:   e.AccountRef,
			Delta:        e.Delta,
			BalanceAfter: newBal,
			RefKind:      intent.RefKind,
			RefID:        refIDPtr,
			IdempKey:     intent.IdempotencyKey,
		}
		m.entries = append(m.entries, row)
		rows = append(rows, row)
	}

	result := Result{TxnID: txnID, Rows: rows, Replayed: false}
	m.txnByIdem[intent.IdempotencyKey] = result
	return result, nil
}

// Entries returns a copy of the ledger for assertion in tests.
func (m *MemoryTxnService) Entries() []LedgerRow {
	m.mu.Lock()
	defer m.mu.Unlock()
	out := make([]LedgerRow, len(m.entries))
	copy(out, m.entries)
	return out
}
