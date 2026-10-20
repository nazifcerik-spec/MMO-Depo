// Package character implements the character lifecycle: create (with slot cap),
// list, enter (bind to session), delete (immediate, name-confirmed, DD-283) / cancel a pending
// account-purge deletion, rename.
//
// Stats + sheet computation is a placeholder here (base values per class); the
// full modifier stack lands with the skill-engine sprint and item-engine.
//
// Layering matches auth:
//
//	http.go        — HTTP handlers.
//	service.go     — business logic + slot cap + grace policy.
//	repository.go  — Repository interface + MemoryRepo.
//	pg_repository.go — Postgres implementation.
//	validation.go  — race/class allowlist + name pattern.
//	stats.go       — base stats per class + derived stats.
package character

import (
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/otter/server/internal/leveling"
)

// State is the character lifecycle marker.
type State string

const (
	StateActive      State = "active"
	StateDeleteGrace State = "delete_grace"
	StateFrozen      State = "frozen"
)

// Character is the persisted character record.
type Character struct {
	ID                 uuid.UUID
	AccountID          uuid.UUID
	Name               string
	NameNormalized     string
	RaceCode           string
	ClassCode          string
	SpecCode           *string
	Level              int16
	XP                 int64
	State              State
	LocationID         *uuid.UUID
	LocationRegionCode string // Sprint 16 — world region code; empty = "at starting region"
	CreatedAt          time.Time
	UpdatedAt          time.Time
	LastPlayedAt       *time.Time
	DeleteGraceAt      *time.Time
	Version            int32

	// HP persistence (DD-218). HPStored nil = full. See EffectiveHP.
	HPStored    *int32
	HPUpdatedAt *time.Time

	// Progression (DD-250): allocated raw attribute points, XP held at a
	// promotion gate, XP earned past the level cap.
	Alloc      leveling.Stats
	OverflowXP int64
	PostCapXP  int64
}

// Spec is the character's current class node (the base class until promoted).
func (c *Character) Spec() string {
	if c.SpecCode != nil && *c.SpecCode != "" {
		return *c.SpecCode
	}
	return c.ClassCode
}

// IsPlayable returns true when the character can be used in-game. Frozen or
// delete_grace characters are read-only from the domain's perspective.
func (c *Character) IsPlayable() bool {
	return c.State == StateActive
}

// --- Errors ------------------------------------------------------------------

var (
	// ErrNotFound is returned when the character id does not match any row.
	ErrNotFound = errors.New("character: not found")

	// ErrForbidden is returned when the caller does not own the target character.
	ErrForbidden = errors.New("character: forbidden")

	// ErrNameTaken is returned when the requested name is already in use.
	ErrNameTaken = errors.New("character: name taken")

	// ErrSlotCap is returned when the account is already at the character-slot
	// cap (DD-OPEN-8; v1 default = 2).
	ErrSlotCap = errors.New("character: slot cap reached")
	// ErrEmailNotVerified: character creation needs a verified e-mail (DD-219).
	ErrEmailNotVerified = errors.New("character: e-mail not verified")

	// ErrNotPlayable is returned for enter/mutation attempts against a frozen
	// or delete_grace character.
	ErrNotPlayable = errors.New("character: not playable")

	// ErrNameMismatch: a delete request must repeat the character's name.
	ErrNameMismatch = errors.New("character: confirmation name does not match")

	// ErrNotInDeleteGrace is returned when cancel-delete is called on a
	// character that is not currently in delete_grace state.
	ErrNotInDeleteGrace = errors.New("character: not in delete grace")

	// ErrValidation wraps input-shape rejections; ValidationError carries the
	// offending field name.
	ErrValidation = errors.New("character: validation")

	// ErrInvalidXP is returned by AddXP when a negative delta is passed.
	// Callers must never subtract XP through this path — XP loss on death
	// (Phase 02) will get its own explicit method.
	ErrInvalidXP = errors.New("character: invalid xp delta")
)

// LevelCap matches the CHECK constraint on the character.level column
// (migration 0036, DD-250). The leveling engine's configured cap must not
// exceed it.
const LevelCap = 1000

// ValidationError attaches a field to ErrValidation for HTTP error mapping.
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return "character: validation: " + e.Field + " " + e.Message
}
func (e *ValidationError) Unwrap() error { return ErrValidation }
