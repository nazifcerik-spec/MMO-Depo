// Package hunt is the AFK combat-encounter loop (Sprint 14).
//
// A hunt = repeated player-vs-mob encounters against a chosen mob template,
// resolved analytically at reveal time using the Sprint 13 combat pipeline.
// Everything is server-authoritative and deterministic: the same (seed, mob,
// player_sheet, elapsed_seconds) always produces the same encounter log +
// rewards.
//
// Design decisions:
//   - Own contract table (`hunt.contract`) rather than piggy-backing on
//     activity.contract. Hunt lifecycle differs — a hunt can end early on
//     player death mid-encounter, which the analytical activity kinds never
//     do — so a separate table keeps both stories readable.
//   - Fixed cadence: one encounter every EncounterSeconds. Full snapshot
//     (player sheet at start) frozen into the contract so mid-hunt gear
//     changes don't retroactively rebalance the log.
//   - Idempotent reveal via `source_ref_id = contract.id` on both the ledger
//     entry (existing pattern) and the item.Grant call.
//   - Player death halts the hunt at that encounter — rewards up to and
//     including the killing blow are still granted (partial), following the
//     "you fought bravely" pattern from Diablo-style dungeon runs.
package hunt

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
)

// State is the contract lifecycle marker.
type State string

const (
	StateActive    State = "active"
	StateRevealed  State = "revealed"
	StateCancelled State = "cancelled"
)

// Contract is one AFK hunt session.
type Contract struct {
	ID             uuid.UUID       `json:"id"`
	CharacterID    uuid.UUID       `json:"character_id"`
	MobCode        string          `json:"mob_code"`
	EncounterSecs  int32           `json:"encounter_seconds"` // fixed cadence; DefaultEncounterSeconds unless overridden
	DurationSecs   int32           `json:"duration_seconds"`
	StartedAt      time.Time       `json:"started_at"`
	ExpiresAt      time.Time       `json:"expires_at"`
	Seed           int64           `json:"-"`
	State          State           `json:"state"`
	PlayerSnapshot json.RawMessage `json:"-"` // frozen player Combatant so mid-hunt gear swaps don't warp the log
	RevealedAt     *time.Time      `json:"revealed_at,omitempty"`
	RevealResult   json.RawMessage `json:"-"`
	RevealTxnID    *uuid.UUID      `json:"reveal_txn_id,omitempty"`
	CancelledAt    *time.Time      `json:"cancelled_at,omitempty"`
	CreatedAt      time.Time       `json:"created_at"`
	UpdatedAt      time.Time       `json:"-"`
	Version        int32           `json:"-"`
}

// ElapsedSeconds is how many seconds of the contract have run at `now`,
// clamped to [0, DurationSecs].
func (c *Contract) ElapsedSeconds(now time.Time) int32 {
	dt := int32(now.Sub(c.StartedAt).Seconds())
	if dt < 0 {
		return 0
	}
	if dt > c.DurationSecs {
		return c.DurationSecs
	}
	return dt
}

// EncountersPlanned is how many encounters will have resolved by `elapsedSec`.
// Integer division — half-finished encounters do not grant partial rewards
// (they carry over into the next reveal on the same active contract if the
// player waits longer, or are forfeited on Stop).
func (c *Contract) EncountersPlanned(elapsedSec int32) int {
	if c.EncounterSecs <= 0 {
		return 0
	}
	return int(elapsedSec / c.EncounterSecs)
}

// --- Errors -----------------------------------------------------------------

var (
	ErrNotFound   = errors.New("hunt: contract not found")
	ErrForbidden  = errors.New("hunt: forbidden")
	ErrSlotBusy   = errors.New("hunt: character already has an active hunt")
	ErrNotActive  = errors.New("hunt: contract not active")
	ErrUnknownMob = errors.New("hunt: unknown mob_code")
	ErrTooWounded = errors.New("hunt: character has no hp left")
	ErrValidation = errors.New("hunt: validation")
)

// --- Tuning -----------------------------------------------------------------

// DefaultEncounterSeconds is the fallback pacing when the client sends 0.
// One encounter per 120s → 30 encounters/hour keeps the reveal log readable
// (max 240 encounters at 8h; the render trims further).
const DefaultEncounterSeconds = 120

// MinEncounterSeconds prevents pathologically fast cadences that would let a
// client bypass rate limits by requesting sub-second encounters.
const MinEncounterSeconds = 30

// MaxHuntSeconds caps a single hunt duration at 12 hours — matches the
// existing activity long-form cap and prevents runaway JSON blobs in the
// reveal log.
const MaxHuntSeconds = 12 * 60 * 60
