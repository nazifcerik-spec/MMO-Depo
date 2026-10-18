// Package dungeon implements the AFK dungeon run (Sprint 15).
//
// A dungeon is an ordered list of mob waves + a boss encounter, resolved
// analytically at reveal time using the Sprint 13 combat pipeline. Solo
// dungeons can be run any number of times (DD-262): like hunts, their XP,
// boss bonus and loot follow the character's level (mobs far below pay
// less, loot is capped at the character's level). Raids (group templates,
// DD-261) keep the weekly lockout in `dungeon.weekly_lockout`.
//
// Design decisions:
//   - Template = ordered waves + boss code. Each wave has (mob_code, count);
//     mobs inside a wave are fought sequentially (HP carry-over) same as
//     hunts, so an attritional dungeon can wear the player down before boss.
//   - Frozen `player_snapshot` at Start (same reason as hunt: mid-run gear
//     swaps must not warp the log).
//   - Reveal waits for `duration_seconds` to elapse. This models the AFK
//     "you're running the dungeon over the next 15 minutes" pattern; early
//     reveal returns E_NOT_READY.
//   - Player death mid-run halts the reveal at that encounter. Rewards up
//     to (but not including) the fatal wave are still granted, matching the
//     "you fought bravely" pattern.
//   - Boss kill grants a guaranteed rare drop tier bonus on top of the
//     accumulated wave loot.
package dungeon

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

// Contract is one AFK dungeon run.
type Contract struct {
	ID             uuid.UUID       `json:"id"`
	CharacterID    uuid.UUID       `json:"character_id"`
	DungeonCode    string          `json:"dungeon_code"`
	DurationSecs   int32           `json:"duration_seconds"`
	StartedAt      time.Time       `json:"started_at"`
	ExpiresAt      time.Time       `json:"expires_at"`
	Seed           int64           `json:"-"`
	State          State           `json:"state"`
	PlayerSnapshot json.RawMessage `json:"-"`
	RevealedAt     *time.Time      `json:"revealed_at,omitempty"`
	RevealResult   json.RawMessage `json:"-"`
	RevealTxnID    *uuid.UUID      `json:"reveal_txn_id,omitempty"`
	CancelledAt    *time.Time      `json:"cancelled_at,omitempty"`
	CreatedAt      time.Time       `json:"created_at"`
	UpdatedAt      time.Time       `json:"-"`
	Version        int32           `json:"-"`
}

// ElapsedSeconds clamped to [0, DurationSecs].
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

// IsReady returns true when enough time has elapsed to reveal (== full
// duration). Dungeons are all-or-nothing on completion.
func (c *Contract) IsReady(now time.Time) bool {
	return c.ElapsedSeconds(now) >= c.DurationSecs
}

// --- Template ---------------------------------------------------------------

// Wave is one segment of a dungeon. All `Count` mobs of `MobCode` are fought
// sequentially before the next wave starts. Count>=1.
type Wave struct {
	MobCode string `json:"mob_code"`
	Count   int    `json:"count"`
}

// Template is the persistent dungeon shape (loaded from content pack).
type Template struct {
	Code            string
	DisplayNameKey  string
	MinLevel        int16
	DurationSecs    int32
	Waves           []Wave
	BossMobCode     string
	BonusGoldCopper int64 // paid on top of accumulated wave gold when boss dies
	BonusXP         int32 // paid on top of accumulated wave xp when boss dies
	GroupMin        int16 // party size bounds (DD-261); 1/1 = a solo dungeon
	GroupMax        int16
}

// IsRaid reports whether the template is a raid (a group run, DD-261).
func (t *Template) IsRaid() bool { return t.GroupMax > 1 }

// TotalEncounters — wave sum + 1 boss. Used for planned-encounter preview.
func (t *Template) TotalEncounters() int {
	n := 0
	for _, w := range t.Waves {
		n += w.Count
	}
	return n + 1 // + boss
}

// --- Weekly lockout ---------------------------------------------------------

// ISOWeek returns the "YYYY-WW" week key used as the lockout key. Keeps the
// week arithmetic in one place — swapping to a different rollover (Tuesday
// reset, calendar-quarter reset, ...) touches only this fn.
func ISOWeek(t time.Time) string {
	y, w := t.UTC().ISOWeek()
	// Pad two digits so lexical sort matches chronological sort.
	if w < 10 {
		return itoa4(y) + "-0" + itoa1(w)
	}
	return itoa4(y) + "-" + itoa2(w)
}

// itoa helpers — avoid pulling strconv for four-digit / one-two-digit ints.
func itoa4(n int) string {
	b := [4]byte{'0', '0', '0', '0'}
	i := 3
	for n > 0 && i >= 0 {
		b[i] = byte('0' + n%10)
		n /= 10
		i--
	}
	return string(b[:])
}
func itoa2(n int) string {
	b := [2]byte{byte('0' + n/10), byte('0' + n%10)}
	return string(b[:])
}
func itoa1(n int) string { return string([]byte{byte('0' + n)}) }

// --- Errors -----------------------------------------------------------------

var (
	ErrNotFound         = errors.New("dungeon: contract not found")
	ErrForbidden        = errors.New("dungeon: forbidden")
	ErrSlotBusy         = errors.New("dungeon: character already has an active dungeon run")
	ErrNotActive        = errors.New("dungeon: contract not active")
	ErrNotReady         = errors.New("dungeon: not enough time elapsed for reveal")
	ErrWeeklyLockout    = errors.New("dungeon: character already cleared this dungeon this week")
	ErrUnknownDungeon   = errors.New("dungeon: unknown dungeon_code")
	ErrUnknownMob       = errors.New("dungeon: wave references unknown mob")
	ErrTooWounded       = errors.New("dungeon: character has no hp left")
	ErrBagFull          = errors.New("dungeon: bag full")
	ErrLevelRequirement = errors.New("dungeon: character level below min_level")
	ErrValidation       = errors.New("dungeon: validation")
)

// --- Tuning -----------------------------------------------------------------

// MaxDungeonSeconds caps a single run at 1 hour — dungeons are meant to be
// shorter than long hunts.
const MaxDungeonSeconds = 60 * 60
