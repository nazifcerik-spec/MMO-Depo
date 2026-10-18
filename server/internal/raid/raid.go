// Package raid runs group dungeons (DD-261). A raid is a dungeon template
// with group_max above one, entered by a party:
//
//   - The party leader opens a lobby (forming); party members join, each
//     freezing a combat snapshot of their sheet (level, lockout, HP and bag
//     checked per member).
//   - The leader starts it once group_min members are in. From then on the
//     raid runs for the template's duration; nobody can join or leave.
//   - After it ends every member collects their own rewards. The fight is
//     resolved from the frozen snapshots and the run's seed, so every member
//     sees the same log; rewards are personal (each member gets the gold,
//     XP — scaled to their level — and their own drop rolls). A clear writes
//     the member's weekly lockout (shared with dungeons).
package raid

import (
	"errors"
	"strconv"
	"time"

	"github.com/google/uuid"
)

// Run states. A running raid whose expires_at has passed is finished.
const (
	StateForming   = "forming"
	StateRunning   = "running"
	StateCancelled = "cancelled"
)

// Run is one raid of a party.
type Run struct {
	ID        uuid.UUID  `json:"id"`
	PartyID   uuid.UUID  `json:"party_id"`
	RaidCode  string     `json:"raid_code"`
	LeaderID  uuid.UUID  `json:"leader_character_id"`
	State     string     `json:"state"`
	Duration  int32      `json:"duration_seconds"`
	StartedAt *time.Time `json:"started_at,omitempty"`
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	Members   []Member   `json:"members"`
	GroupMin  int16      `json:"group_min"`
	GroupMax  int16      `json:"group_max"`
	Seed      int64      `json:"-"`
}

// Member is a joined character. Its label in the fight log is "p<slot>".
type Member struct {
	ID          uuid.UUID `json:"-"`
	CharacterID uuid.UUID `json:"character_id"`
	Name        string    `json:"name"`
	ClassCode   string    `json:"class_code"`
	Level       int16     `json:"level"`
	Slot        int16     `json:"slot"`
	Label       string    `json:"label"` // the member's name in the fight log
	Collected   bool      `json:"collected"`
}

// Label is a member's name in the fight log.
func Label(slot int16) string { return "p" + strconv.Itoa(int(slot)) }

var (
	ErrUnknownRaid      = errors.New("raid: unknown raid_code")
	ErrNotFound         = errors.New("raid: not found")
	ErrNoParty          = errors.New("raid: character is not in a party")
	ErrNotLeader        = errors.New("raid: party leader only")
	ErrNotInParty       = errors.New("raid: character is not in the raid's party")
	ErrLobbyExists      = errors.New("raid: the party already has an open raid")
	ErrAlreadyInRaid    = errors.New("raid: character is already in a raid")
	ErrNotForming       = errors.New("raid: raid is no longer open")
	ErrRaidFull         = errors.New("raid: raid is full")
	ErrNotEnough        = errors.New("raid: not enough members")
	ErrNotRunning       = errors.New("raid: raid is not running")
	ErrNotReady         = errors.New("raid: raid has not finished")
	ErrNotMember        = errors.New("raid: character is not in this raid")
	ErrLevelRequirement = errors.New("raid: character level below min_level")
	ErrWeeklyLockout    = errors.New("raid: already cleared this raid this week")
	ErrTooWounded       = errors.New("raid: character has no hp left")
	ErrBagFull          = errors.New("raid: bag full")
)
