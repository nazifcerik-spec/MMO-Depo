package vendor

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// Repository is the persistence contract for vendor templates and buybacks.
// Vendor templates load once at boot (like item.Catalog).
type Repository interface {
	// Catalog loads all vendor templates (called at boot).
	Catalog(ctx context.Context) ([]Template, error)

	// InsertBuyback persists a fresh buyback entry.
	InsertBuyback(ctx context.Context, b BuybackEntry) error

	// GetBuyback returns a buyback for restore. Ownership + expiry check is
	// the service's job.
	GetBuyback(ctx context.Context, id uuid.UUID) (*BuybackEntry, error)

	// ClaimBuyback atomically marks a buyback claimed (state transition
	// active → claimed). Fails on already-claimed / expired.
	ClaimBuyback(ctx context.Context, id uuid.UUID, at time.Time) error

	// ListBuybacks returns character's active (unclaimed + not expired)
	// buybacks, newest first.
	ListBuybacks(ctx context.Context, characterID uuid.UUID, limit int) ([]BuybackEntry, error)
}
