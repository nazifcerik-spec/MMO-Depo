// Package economy is the ledger-first transaction service (DD-135).
//
// Every value transfer in Otter flows through TxnService.Write. The invariants:
//
//   - Double-entry: SUM(delta) per currency == 0 across every Intent.
//   - Idempotency: (idempotency_key, entry_seq) UNIQUE in the ledger.
//   - SERIALIZABLE: PgTxnService wraps writes in SERIALIZABLE isolation with
//     up to 8 attempts with jittered backoff on serialization_failure.
//   - Cap + non-negative: character/account/guild scope balances are clamped
//     [0, currency.cap]; system-scope is unbounded (source/sink accounting).
//
// Domain modules (vendor, market, quest reward, activity reveal…) call Write
// exclusively — they never mutate wallet projection rows directly.
package economy

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
)

// AccountKind identifies which projection table a balance lives in and which
// invariant tier applies (system = unbounded, others = clamped).
type AccountKind string

const (
	AccountCharacter AccountKind = "character"
	AccountAccount   AccountKind = "account"
	AccountGuild     AccountKind = "guild"
	AccountSystem    AccountKind = "system"
)

// Entry is a single line in a double-entry set. Delta is signed: positive =
// credit, negative = debit. For a valid Intent, sum(Delta) per currency must
// equal 0.
type Entry struct {
	Currency    string
	AccountKind AccountKind
	AccountRef  uuid.UUID
	Delta       int64
}

// Intent is what a domain module hands to TxnService.Write. Kind is a
// machine-readable label (`vendor_buy`, `mint`, `quest_reward`); RefKind + RefID
// point back to the source aggregate for audit tracing.
type Intent struct {
	Kind           string
	RefKind        string
	RefID          uuid.UUID
	IdempotencyKey string
	Entries        []Entry
}

// LedgerRow is the persisted representation of a ledger entry. Returned in
// Result so callers can build follow-up projections + emit events.
type LedgerRow struct {
	ID           int64
	TxnID        uuid.UUID
	EntrySeq     int16
	Kind         string
	Currency     string
	AccountKind  AccountKind
	AccountRef   uuid.UUID
	Delta        int64
	BalanceAfter int64
	RefKind      string
	RefID        *uuid.UUID
	IdempKey     string
}

// Result is what Write returns. Replayed indicates the response came from the
// idempotency cache (no new writes happened).
type Result struct {
	TxnID    uuid.UUID
	Rows     []LedgerRow
	Replayed bool
}

// --- Errors ------------------------------------------------------------------

var (
	// ErrValidation covers invalid Intent shape (missing key, sum != 0, unknown
	// currency, empty entries).
	ErrValidation = errors.New("economy: validation")

	// ErrInsufficientBalance is returned when a debit would drive a non-system
	// balance below zero.
	ErrInsufficientBalance = errors.New("economy: insufficient balance")

	// ErrCurrencyCap is returned when a credit would exceed the currency cap
	// for a non-system account.
	ErrCurrencyCap = errors.New("economy: currency cap exceeded")

	// ErrScopeMismatch is returned when the AccountKind of an entry doesn't
	// match the currency's declared scope (e.g. account-scope 'gem' credited
	// to a character).
	ErrScopeMismatch = errors.New("economy: currency scope mismatch")

	// ErrUnknownCurrency is returned when Entry.Currency isn't in the catalog.
	ErrUnknownCurrency = errors.New("economy: unknown currency")

	// ErrSerialization is returned when SERIALIZABLE retry budget is
	// exhausted. Callers may retry with backoff.
	ErrSerialization = errors.New("economy: serialization retry exhausted")
)

// ValidationError attaches the offending field to ErrValidation.
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return "economy: validation: " + e.Field + ": " + e.Message
}
func (e *ValidationError) Unwrap() error { return ErrValidation }

// ValidateIntent enforces the shape rules that both PgTxnService and
// MemoryTxnService need to check before any persistence work. Returns nil when
// valid; a ValidationError otherwise.
func ValidateIntent(intent Intent, currencies map[string]Currency) error {
	if intent.IdempotencyKey == "" {
		return &ValidationError{Field: "idempotency_key", Message: "required"}
	}
	if len(intent.Kind) == 0 {
		return &ValidationError{Field: "kind", Message: "required"}
	}
	if len(intent.Entries) < 2 {
		return &ValidationError{Field: "entries", Message: "at least two entries required (double-entry)"}
	}

	// Currency existence + scope compatibility + delta != 0.
	sums := make(map[string]int64)
	for i, e := range intent.Entries {
		if e.Delta == 0 {
			return &ValidationError{Field: fmt.Sprintf("entries[%d].delta", i), Message: "must be non-zero"}
		}
		if e.AccountRef == uuid.Nil {
			return &ValidationError{Field: fmt.Sprintf("entries[%d].account_ref", i), Message: "required"}
		}
		cur, ok := currencies[e.Currency]
		if !ok {
			return fmt.Errorf("%w: %s", ErrUnknownCurrency, e.Currency)
		}
		// Scope check — system may credit/debit any currency; non-system must
		// match declared currency scope.
		if e.AccountKind != AccountSystem && string(e.AccountKind) != cur.Scope {
			return fmt.Errorf("%w: currency %s is scope %s, entry account_kind %s",
				ErrScopeMismatch, cur.Code, cur.Scope, e.AccountKind)
		}
		sums[e.Currency] += e.Delta
	}
	// Sum invariant per currency.
	for cur, s := range sums {
		if s != 0 {
			return &ValidationError{
				Field:   "entries",
				Message: fmt.Sprintf("sum of delta for currency %s must be zero, got %d", cur, s),
			}
		}
	}
	return nil
}
