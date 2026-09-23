package economy

import (
	"context"

	"github.com/google/uuid"
)

// TxnService is the write + read contract for ledger operations. Every value
// transfer in Otter passes through Write; balance snapshots come from Balance.
type TxnService interface {
	// Write applies an Intent atomically. On idempotency HIT the cached result
	// is returned with Replayed=true.
	Write(ctx context.Context, intent Intent) (Result, error)

	// Balance returns the current projection for the (kind, ref, currency)
	// triple. Returns 0 for accounts that have never been written to.
	Balance(ctx context.Context, kind AccountKind, ref uuid.UUID, currency string) (int64, error)
}
