// Package activity implements the AFK Engine per DD-084 — analytical reveal
// (Poisson + Binomial samplers, no per-tick simulation).
//
// Layering:
//
//	sampler.go       — deterministic RNG (xoshiro256++) + Poisson/Binomial.
//	yield.go         — kind-specific yield formula (v1 placeholder: gold only).
//	snapshot.go      — freeze character state at contract start (DD-088).
//	repository.go    — Repository interface + MemoryRepo (tests).
//	pg_repository.go — Postgres implementation.
//	service.go       — Start / Reveal / Stop / Active with TxnService integration.
//	http.go          — HTTP handlers.
//
// Reveal flow (idempotent):
//  1. Load contract; return cached reveal_result if state='revealed'.
//  2. Compute elapsed_seconds = min(now - started_at, duration).
//  3. sampler(seed) → yield summary (gold + placeholder loot).
//  4. TxnService.Write mint intent (system.reward_source → character.wallet).
//  5. Persist reveal_result + reveal_txn_id + state='revealed'.
//  6. Return summary.
package activity

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
)

// Kind is a machine-readable activity type. Content-driven allowlist would land
// with the content pack loader (Sprint 12); v1 uses this small enum.
type Kind string

const (
	KindMining      Kind = "mining"
	KindFishing     Kind = "fishing"
	KindWoodcutting Kind = "woodcutting"
	KindHerbalism   Kind = "herbalism"
	KindHunting     Kind = "hunting"
)

// KnownKinds is the v1 allowlist. Phase 02 replaces with content.activity_kind.
var KnownKinds = map[Kind]bool{
	KindMining:      true,
	KindFishing:     true,
	KindWoodcutting: true,
	KindHerbalism:   true,
	KindHunting:     true,
}

// State is the contract lifecycle marker.
type State string

const (
	StateActive    State = "active"
	StateRevealed  State = "revealed"
	StateCancelled State = "cancelled"
)

// Contract is one active or historical AFK activity.
type Contract struct {
	ID              uuid.UUID
	CharacterID     uuid.UUID
	Kind            Kind
	DurationSeconds int32
	StartedAt       time.Time
	ExpiresAt       time.Time
	Seed            int64
	State           State
	Snapshot        json.RawMessage
	Params          json.RawMessage

	RevealedAt   *time.Time
	RevealResult json.RawMessage
	RevealTxnID  *uuid.UUID

	CancelledAt *time.Time

	CreatedAt time.Time
	UpdatedAt time.Time
	Version   int32
}

// ElapsedSeconds computes how many seconds of the contract have run at `now`,
// clamped to [0, DurationSeconds].
func (c *Contract) ElapsedSeconds(now time.Time) int32 {
	dt := int32(now.Sub(c.StartedAt).Seconds())
	if dt < 0 {
		return 0
	}
	if dt > c.DurationSeconds {
		return c.DurationSeconds
	}
	return dt
}

// --- Errors ------------------------------------------------------------------

var (
	// ErrNotFound is returned when the contract id does not resolve.
	ErrNotFound = errors.New("activity: contract not found")

	// ErrForbidden is returned when the caller does not own the target character.
	ErrForbidden = errors.New("activity: forbidden")

	// ErrSlotBusy is returned by Start when the character already has an
	// active contract (v1 concurrent cap = 1; DD-OPEN-3).
	ErrSlotBusy = errors.New("activity: character already has an active contract")

	// ErrNotActive is returned when reveal/stop targets a non-active contract
	// (already revealed or cancelled).
	ErrNotActive = errors.New("activity: contract not active")

	// ErrValidation wraps input-shape rejection.
	ErrValidation = errors.New("activity: validation")
)

// ValidationError attaches the offending field for HTTP mapping.
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return "activity: validation: " + e.Field + ": " + e.Message
}
func (e *ValidationError) Unwrap() error { return ErrValidation }

// Duration limits (v1). Config Phase 02.
const (
	MinDurationSeconds = int32(60)        // 1 minute
	MaxDurationSeconds = int32(12 * 3600) // 12 hours
)
