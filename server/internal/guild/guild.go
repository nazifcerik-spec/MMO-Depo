// Package guild implements guilds (Sprint 20, docs/guild/*): creation with a
// ledger-debited fee, invites, a fixed 5-rank permission matrix, leadership
// transfer, MOTD and 48h scheduled disband (finalised by the worker).
// Treasury, bank, research and war are later sprints.
package guild

import (
	"errors"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
)

const (
	CreationFeeCopper = 100_000 // docs/guild/guild.md §3 — 10 gold
	MemberCap         = 50
	InviteTTL         = 48 * time.Hour
	DisbandGrace      = 48 * time.Hour
	MaxMOTD           = 256
	MaxDescription    = 512
)

// Rank indexes (docs/guild/roles.md). Lower is more senior.
const (
	RankLeader   int16 = 0
	RankOfficer  int16 = 1
	RankVeteran  int16 = 2
	RankMember   int16 = 3
	RankInitiate int16 = 4
)

var RankNames = map[int16]string{
	RankLeader: "leader", RankOfficer: "officer", RankVeteran: "veteran",
	RankMember: "member", RankInitiate: "initiate",
}

// --- permission matrix (docs/guild/permissions.md §3) -----------------------

func CanInvite(rank int16) bool  { return rank <= RankVeteran }
func CanSetMOTD(rank int16) bool { return rank <= RankOfficer }

// CanKick: leader/officer may remove strictly lower ranks.
func CanKick(actor, target int16) bool { return actor <= RankOfficer && target > actor }

// CanSetRank: leader may move anyone below them to ranks 1..4; officers may
// move members below officer to ranks veteran..initiate. Rank 0 only moves
// via TransferLeadership.
func CanSetRank(actor, target, newRank int16) bool {
	if newRank < RankOfficer || newRank > RankInitiate || target <= actor {
		return false
	}
	switch actor {
	case RankLeader:
		return true
	case RankOfficer:
		return newRank >= RankVeteran
	}
	return false
}

// --- validation -------------------------------------------------------------

var (
	nameRe = regexp.MustCompile(`^[\p{L}\p{N}]+( [\p{L}\p{N}]+)*$`)
	tagRe  = regexp.MustCompile(`^[A-Za-z0-9]{2,5}$`)
)

// NormalizeName validates a guild name and returns (display, normalized).
func NormalizeName(raw string) (string, string, error) {
	s := strings.TrimSpace(raw)
	if n := len([]rune(s)); n < 3 || n > 24 || !nameRe.MatchString(s) {
		return "", "", ErrInvalidName
	}
	return s, strings.ToLower(s), nil
}

// NormalizeTag validates a tag and returns (display, normalized).
func NormalizeTag(raw string) (string, string, error) {
	s := strings.TrimSpace(raw)
	if !tagRe.MatchString(s) {
		return "", "", ErrInvalidTag
	}
	return strings.ToUpper(s), strings.ToLower(s), nil
}

// --- types ------------------------------------------------------------------

type Member struct {
	CharacterID uuid.UUID `json:"character_id"`
	Name        string    `json:"name"`
	Level       int16     `json:"level"`
	ClassCode   string    `json:"class_code"`
	Rank        int16     `json:"rank_index"`
	RankName    string    `json:"rank_name"`
	JoinedAt    time.Time `json:"joined_at"`
}

type Guild struct {
	ID                 uuid.UUID  `json:"id"`
	Name               string     `json:"name"`
	Tag                string     `json:"tag"`
	LeaderCharacterID  uuid.UUID  `json:"leader_character_id"`
	MOTD               string     `json:"motd"`
	Description        string     `json:"description"`
	Level              int16      `json:"level"`
	XP                 int64      `json:"xp"`
	MemberCount        int        `json:"member_count"`
	MemberCap          int16      `json:"member_cap"`
	CreatedAt          time.Time  `json:"created_at"`
	DisbandScheduledAt *time.Time `json:"disband_scheduled_at,omitempty"`
	Members            []Member   `json:"members,omitempty"`
	MyRank             *int16     `json:"my_rank,omitempty"`
}

type Invite struct {
	ID            uuid.UUID `json:"id"`
	GuildID       uuid.UUID `json:"guild_id"`
	GuildName     string    `json:"guild_name"`
	GuildTag      string    `json:"guild_tag"`
	InvitedByName string    `json:"invited_by_name"`
	ExpiresAt     time.Time `json:"expires_at"`
}

var (
	ErrNotFound           = errors.New("guild: not found")
	ErrNotInGuild         = errors.New("guild: character is not in a guild")
	ErrAlreadyInGuild     = errors.New("guild: character is already in a guild")
	ErrForbidden          = errors.New("guild: rank lacks permission")
	ErrInvalidName        = errors.New("guild: name must be 3-24 letters/digits, single spaces")
	ErrInvalidTag         = errors.New("guild: tag must be 2-5 letters/digits")
	ErrNameTaken          = errors.New("guild: name already taken")
	ErrTagTaken           = errors.New("guild: tag already taken")
	ErrGuildFull          = errors.New("guild: member cap reached")
	ErrTargetNotFound     = errors.New("guild: target character not found")
	ErrSelfTarget         = errors.New("guild: cannot target yourself")
	ErrTargetInGuild      = errors.New("guild: target is already in a guild")
	ErrAlreadyInvited     = errors.New("guild: target already has a pending invite")
	ErrInviteInvalid      = errors.New("guild: invite is not pending or has expired")
	ErrNotMember          = errors.New("guild: target is not a member")
	ErrLeaderCannotLeave  = errors.New("guild: leader must transfer leadership or disband first")
	ErrDisbandScheduled   = errors.New("guild: disband already scheduled")
	ErrDisbandNotPending  = errors.New("guild: no disband scheduled")
	ErrTooLong            = errors.New("guild: text too long")
	ErrIdempotencyMissing = errors.New("guild: Idempotency-Key header required")
)
