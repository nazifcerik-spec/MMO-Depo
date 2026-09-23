// Package party implements ephemeral groups of up to 5 characters
// (Sprint 20, docs/party/*). v1: leader-only invites, 60s invite TTL,
// loot rule fixed to 'personal', no activity binding yet.
package party

import (
	"errors"
	"sort"
	"time"

	"github.com/google/uuid"
)

const (
	SizeCap   = 5
	InviteTTL = 60 * time.Second
)

type Member struct {
	CharacterID uuid.UUID `json:"character_id"`
	Name        string    `json:"name"`
	Level       int16     `json:"level"`
	ClassCode   string    `json:"class_code"`
	RoleHint    *string   `json:"role_hint,omitempty"`
	JoinedAt    time.Time `json:"joined_at"`
}

type Party struct {
	ID                uuid.UUID `json:"id"`
	LeaderCharacterID uuid.UUID `json:"leader_character_id"`
	SizeCap           int16     `json:"size_cap"`
	LootRule          string    `json:"loot_rule"`
	CreatedAt         time.Time `json:"created_at"`
	Members           []Member  `json:"members"`
}

type Invite struct {
	ID              uuid.UUID `json:"id"`
	PartyID         uuid.UUID `json:"party_id"`
	InvitedBy       uuid.UUID `json:"invited_by_character_id"`
	InvitedByName   string    `json:"invited_by_name"`
	TargetCharacter uuid.UUID `json:"target_character_id"`
	ExpiresAt       time.Time `json:"expires_at"`
}

var (
	ErrNotFound       = errors.New("party: not found")
	ErrNotInParty     = errors.New("party: character is not in a party")
	ErrAlreadyInParty = errors.New("party: character is already in a party")
	ErrNotLeader      = errors.New("party: leader only")
	ErrPartyFull      = errors.New("party: party is full")
	ErrTargetNotFound = errors.New("party: target character not found")
	ErrSelfTarget     = errors.New("party: cannot target yourself")
	ErrTargetInParty  = errors.New("party: target is already in a party")
	ErrAlreadyInvited = errors.New("party: target already has a pending invite")
	ErrInviteInvalid  = errors.New("party: invite is not pending or has expired")
	ErrNotMember      = errors.New("party: target is not a member")
)

// NextLeader picks who inherits leadership when `leaving` goes: the member
// who joined earliest among the rest. ok=false means nobody is left.
func NextLeader(members []Member, leaving uuid.UUID) (uuid.UUID, bool) {
	rest := make([]Member, 0, len(members))
	for _, m := range members {
		if m.CharacterID != leaving {
			rest = append(rest, m)
		}
	}
	if len(rest) == 0 {
		return uuid.Nil, false
	}
	sort.SliceStable(rest, func(i, j int) bool { return rest[i].JoinedAt.Before(rest[j].JoinedAt) })
	return rest[0].CharacterID, true
}
