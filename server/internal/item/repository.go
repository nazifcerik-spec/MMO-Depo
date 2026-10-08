package item

import (
	"context"

	"github.com/google/uuid"
)

// Repository is the persistence contract for item_instance rows.
// Template catalog is separate (Catalog struct, loaded once at boot).
type Repository interface {
	// CreateInstancesBatch inserts N instances in one round trip — used by
	// activity reveal which mints many drops at once.
	CreateInstancesBatch(ctx context.Context, instances []Instance) error

	Get(ctx context.Context, id uuid.UUID) (*Instance, error)
	ListByOwner(ctx context.Context, characterID uuid.UUID, limit int) ([]Instance, error)
	ListBySource(ctx context.Context, sourceRefID uuid.UUID) ([]Instance, error)
}
