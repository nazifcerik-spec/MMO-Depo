package character

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/otter/server/internal/leveling"
)

// Repository is the persistence contract for the character domain.
type Repository interface {
	Create(ctx context.Context, c Character) error
	Get(ctx context.Context, id uuid.UUID) (*Character, error)
	ListByAccount(ctx context.Context, accountID uuid.UUID) ([]Character, error)
	CountActiveByAccount(ctx context.Context, accountID uuid.UUID) (int, error)

	// Rename atomically updates name + name_normalized. Returns ErrNameTaken
	// on unique-violation.
	Rename(ctx context.Context, id uuid.UUID, name, normalized string) error

	// TouchLastPlayed sets last_played_at = ts. No-op on non-active characters.
	TouchLastPlayed(ctx context.Context, id uuid.UUID, ts time.Time) error

	// SetDeleteGrace transitions state='delete_grace' with the given expire
	// time. Guard: only active characters may enter this state.
	SetDeleteGrace(ctx context.Context, id uuid.UUID, expiresAt time.Time) error

	// CancelDeleteGrace transitions state='delete_grace' → 'active'.
	CancelDeleteGrace(ctx context.Context, id uuid.UUID) error

	// GrantXP applies an XP grant in one transaction (DD-250): source
	// idempotency, level-ups through the leveling engine, gate overflow, level
	// rewards / unlocks / titles / events. See Service.GrantXP.
	GrantXP(ctx context.Context, id uuid.UUID, amount int64, src XPSource, cfg *leveling.Config) (*Grant, error)

	// SetLocation writes the character's `location_region_code`. Sprint 16
	// world/travel completion path is the only writer today.
	SetLocation(ctx context.Context, id uuid.UUID, regionCode string) error

	// GetLocation returns the character's current region_code, or empty when
	// unset (fresh character before first travel completion — caller should
	// treat as "at the starting region" configured in the world catalog).
	GetLocation(ctx context.Context, id uuid.UUID) (string, error)

	// SetHP stores current HP as of `at` (DD-218). Regen from then on is
	// computed by readers (EffectiveHP); nothing ticks.
	SetHP(ctx context.Context, id uuid.UUID, hp int32, at time.Time) error
}
