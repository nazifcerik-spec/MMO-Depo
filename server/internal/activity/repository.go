package activity

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// Repository is the persistence contract.
type Repository interface {
	Create(ctx context.Context, c Contract) error

	// GetByID returns a single contract; ErrNotFound when missing.
	GetByID(ctx context.Context, id uuid.UUID) (*Contract, error)

	// ActiveByCharacter returns the character's active contract if any, or
	// ErrNotFound.
	ActiveByCharacter(ctx context.Context, characterID uuid.UUID) (*Contract, error)

	// MarkRevealed atomically transitions active → revealed and stores the
	// cached result. Fails if the contract is not in state='active'.
	MarkRevealed(ctx context.Context, id uuid.UUID, revealResult json.RawMessage, revealTxnID uuid.UUID, at time.Time) error

	// MarkCancelled atomically transitions active → cancelled. Fails if not
	// in state='active'.
	MarkCancelled(ctx context.Context, id uuid.UUID, at time.Time) error

	// ListByCharacter returns recent contracts (any state), newest first.
	ListByCharacter(ctx context.Context, characterID uuid.UUID, limit int) ([]Contract, error)
}
