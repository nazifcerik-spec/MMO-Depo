// Package social implements friends, the block list, mail and chat
// (Sprint 21, docs/social/*).
package social

import (
	"bytes"
	"context"
	"errors"
	"log/slog"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/otter/server/internal/economy"
	"github.com/otter/server/internal/item"
	"github.com/otter/server/internal/notification"
)

const (
	FriendCap        = 200
	InboxCap         = 100
	MailTTL          = 30 * 24 * time.Hour
	MaxAttachments   = 12
	MaxSubject       = 64
	MaxMailBody      = 2000
	MaxChatBody      = 512
	PostageBase      = 5 // copper — docs/social/mail.md §4
	PostagePerItem   = 1
	PostageGoldRate  = 0.001
	ChatHistoryLimit = 50
)

var (
	ErrNotFound           = errors.New("social: not found")
	ErrTargetNotFound     = errors.New("social: target character not found")
	ErrSelfTarget         = errors.New("social: cannot target yourself")
	ErrBlocked            = errors.New("social: blocked")
	ErrAlreadyFriends     = errors.New("social: already friends")
	ErrAlreadyRequested   = errors.New("social: request already pending")
	ErrFriendCap          = errors.New("social: friend list is full")
	ErrNotFriends         = errors.New("social: not friends")
	ErrInboxFull          = errors.New("social: recipient inbox is full")
	ErrTooLong            = errors.New("social: text too long")
	ErrEmpty              = errors.New("social: text is empty")
	ErrTooManyAttachments = errors.New("social: at most 12 item attachments")
	ErrAttachmentInvalid  = errors.New("social: attachment must be your own unlocked, unbound, unequipped item")
	ErrHasAttachments     = errors.New("social: claim attachments before deleting")
	ErrNothingToClaim     = errors.New("social: nothing to claim")
	ErrBadChannel         = errors.New("social: unknown chat channel")
	ErrNotInChannel       = errors.New("social: not a member of that channel")
	ErrIdempotencyMissing = errors.New("social: Idempotency-Key header required")
	ErrMuted              = errors.New("social: account is muted")
)

// Owner verifies an account owns a character. character.Service satisfies.
type Owner interface {
	Owns(ctx context.Context, characterID, accountID uuid.UUID) error
}

// SlotAssigner places claimed mail items into the bag. inventory.Service satisfies.
type SlotAssigner interface {
	AssignSlots(ctx context.Context, ownerID uuid.UUID, instances []item.Instance) error
}

// Publisher delivers WS events. wsgateway.KindPayloadAdapter satisfies.
type Publisher interface {
	Publish(accountID uuid.UUID, kind string, payload any) int
	BroadcastKind(kind string, payload any) int
}

type Service struct {
	pool   *pgxpool.Pool
	chars  Owner
	txn    economy.TxnService
	slots  SlotAssigner
	pub    Publisher
	notify notification.Notifier
	now    func() time.Time
	log    *slog.Logger
}

func NewService(pool *pgxpool.Pool, chars Owner, txn economy.TxnService, slots SlotAssigner, logger *slog.Logger) *Service {
	if logger == nil {
		logger = slog.Default()
	}
	return &Service{pool: pool, chars: chars, txn: txn, slots: slots, now: time.Now, log: logger}
}

func (s *Service) WithPublisher(p Publisher) *Service { s.pub = p; return s }

func (s *Service) WithNotifier(n notification.Notifier) *Service { s.notify = n; return s }

// Postage returns the mail fee in copper (docs/social/mail.md §4).
func Postage(items int, gold int64) int64 {
	return PostageBase + int64(items)*PostagePerItem + int64(float64(gold)*PostageGoldRate+0.5)
}

// CleanText trims and enforces 1..max runes.
func CleanText(s string, max int) (string, error) {
	s = strings.TrimSpace(s)
	n := len([]rune(s))
	if n == 0 {
		return "", ErrEmpty
	}
	if n > max {
		return "", ErrTooLong
	}
	return s, nil
}

// --- shared helpers ---------------------------------------------------------

type querier interface {
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
}

type charRef struct {
	ID      uuid.UUID `json:"character_id"`
	Name    string    `json:"name"`
	Account uuid.UUID `json:"-"`
}

func characterByName(ctx context.Context, q querier, name string) (charRef, error) {
	var c charRef
	err := q.QueryRow(ctx, `
		SELECT id, name, account_id FROM character.character
		 WHERE name_normalized = $1 AND state = 'active'`, strings.ToLower(strings.TrimSpace(name))).
		Scan(&c.ID, &c.Name, &c.Account)
	if errors.Is(err, pgx.ErrNoRows) {
		return c, ErrTargetNotFound
	}
	return c, err
}

func accountOf(ctx context.Context, q querier, character uuid.UUID) (uuid.UUID, error) {
	var a uuid.UUID
	err := q.QueryRow(ctx, `SELECT account_id FROM character.character WHERE id = $1`, character).Scan(&a)
	return a, err
}

// blockedEither is true when either side has blocked the other.
func blockedEither(ctx context.Context, q querier, a, b uuid.UUID) (bool, error) {
	var blocked bool
	err := q.QueryRow(ctx, `
		SELECT EXISTS (SELECT 1 FROM social.ignore_list
		                WHERE (owner_character_id = $1 AND target_character_id = $2)
		                   OR (owner_character_id = $2 AND target_character_id = $1))`, a, b).Scan(&blocked)
	return blocked, err
}

// pair orders two ids the way Postgres compares uuids (byte-wise), matching
// the friendship CHECK (character_a < character_b).
func pair(a, b uuid.UUID) (uuid.UUID, uuid.UUID) {
	if bytes.Compare(a[:], b[:]) < 0 {
		return a, b
	}
	return b, a
}

func (s *Service) publish(account uuid.UUID, kind string, payload any) {
	if s.pub != nil && account != uuid.Nil {
		s.pub.Publish(account, kind, payload)
	}
}

func isUnique(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "23505"
}
