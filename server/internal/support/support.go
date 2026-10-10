// Package support is the in-app support desk (DD-234,
// docs/release/closed-beta.md §5): players file tickets, operators answer
// them from the admin panel and the player gets a notification.
package support

import (
	"context"
	"errors"
	"log/slog"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/otter/server/internal/notification"
)

// Limits.
const (
	MaxOpenPerAccount = 10
	MaxBody           = 4000
	MaxSubject        = 120
)

var Categories = map[string]bool{"bug": true, "account": true, "abuse": true, "feedback": true, "other": true}

// clientInfoKeys are what the app may attach (no free-form data).
var clientInfoKeys = map[string]bool{"platform": true, "locale": true, "flavor": true, "app_version": true}

var (
	ErrInvalid          = errors.New("support: invalid ticket")
	ErrTooManyOpen      = errors.New("support: too many open tickets")
	ErrNotFound         = errors.New("support: ticket not found")
	ErrClosed           = errors.New("support: ticket closed")
	ErrReplyRequired    = errors.New("support: reply required")
	ErrForeignCharacter = errors.New("support: character not on this account")
)

// Ticket is one support request with its answer.
type Ticket struct {
	ID          uuid.UUID         `json:"id"`
	AccountID   uuid.UUID         `json:"-"`
	CharacterID *uuid.UUID        `json:"character_id,omitempty"`
	Category    string            `json:"category"`
	Subject     string            `json:"subject"`
	Body        string            `json:"body"`
	ClientInfo  map[string]string `json:"client_info,omitempty"`
	State       string            `json:"state"`
	Reply       *string           `json:"reply,omitempty"`
	RepliedAt   *time.Time        `json:"replied_at,omitempty"`
	ClosedAt    *time.Time        `json:"closed_at,omitempty"`
	CreatedAt   time.Time         `json:"created_at"`
}

// NewTicket is what a player submits.
type NewTicket struct {
	Category    string            `json:"category"`
	Subject     string            `json:"subject"`
	Body        string            `json:"body"`
	CharacterID *uuid.UUID        `json:"character_id"`
	ClientInfo  map[string]string `json:"client_info"`
}

type Service struct {
	pool     *pgxpool.Pool
	notifier notification.Notifier
	log      *slog.Logger
}

func NewService(pool *pgxpool.Pool, notifier notification.Notifier, logger *slog.Logger) *Service {
	if logger == nil {
		logger = slog.Default()
	}
	return &Service{pool: pool, notifier: notifier, log: logger}
}

func clean(t NewTicket) (NewTicket, error) {
	t.Subject = strings.TrimSpace(t.Subject)
	t.Body = strings.TrimSpace(t.Body)
	if !Categories[t.Category] || utf8.RuneCountInString(t.Subject) < 3 || utf8.RuneCountInString(t.Subject) > MaxSubject ||
		utf8.RuneCountInString(t.Body) < 10 || utf8.RuneCountInString(t.Body) > MaxBody {
		return t, ErrInvalid
	}
	info := map[string]string{}
	for k, v := range t.ClientInfo {
		if clientInfoKeys[k] {
			if r := []rune(strings.TrimSpace(v)); len(r) > 64 {
				v = string(r[:64])
			}
			info[k] = strings.TrimSpace(v)
		}
	}
	t.ClientInfo = info
	return t, nil
}

// Create files a ticket for the account.
func (s *Service) Create(ctx context.Context, account uuid.UUID, req NewTicket) (*Ticket, error) {
	req, err := clean(req)
	if err != nil {
		return nil, err
	}
	var out *Ticket
	err = pgx.BeginFunc(ctx, s.pool, func(tx pgx.Tx) error {
		// Serialise per account so the open-ticket cap holds under concurrency.
		if _, err := tx.Exec(ctx, `SELECT 1 FROM auth.account WHERE id = $1 FOR UPDATE`, account); err != nil {
			return err
		}
		var open int
		if err := tx.QueryRow(ctx, `SELECT count(*) FROM ops.support_ticket WHERE account_id = $1 AND state <> 'closed'`,
			account).Scan(&open); err != nil {
			return err
		}
		if open >= MaxOpenPerAccount {
			return ErrTooManyOpen
		}
		if req.CharacterID != nil {
			var owner uuid.UUID
			err := tx.QueryRow(ctx, `SELECT account_id FROM character.character WHERE id = $1`, *req.CharacterID).Scan(&owner)
			if errors.Is(err, pgx.ErrNoRows) || (err == nil && owner != account) {
				return ErrForeignCharacter
			}
			if err != nil {
				return err
			}
		}
		rows, err := tx.Query(ctx, `
			INSERT INTO ops.support_ticket (account_id, character_id, category, subject, body, client_info)
			VALUES ($1, $2, $3, $4, $5, $6) RETURNING `+cols, account, req.CharacterID, req.Category, req.Subject, req.Body, req.ClientInfo)
		if err != nil {
			return err
		}
		t, err := pgx.CollectOneRow(rows, scan)
		out = &t
		return err
	})
	return out, err
}

// Mine lists the account's latest tickets, newest first.
func (s *Service) Mine(ctx context.Context, account uuid.UUID) ([]Ticket, error) {
	rows, err := s.pool.Query(ctx, `SELECT `+cols+` FROM ops.support_ticket WHERE account_id = $1
	                                 ORDER BY created_at DESC LIMIT 50`, account)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, scan)
}

// Queue lists tickets in a state for operators, oldest first ("" = open).
func (s *Service) Queue(ctx context.Context, state string, limit int) ([]Ticket, error) {
	if state == "" {
		state = "open"
	}
	if limit <= 0 || limit > 200 {
		limit = 50
	}
	rows, err := s.pool.Query(ctx, `SELECT `+cols+` FROM ops.support_ticket WHERE state = $1
	                                 ORDER BY created_at LIMIT $2`, state, limit)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, scan)
}

// Answer records an operator's reply (and optionally closes the ticket) and
// notifies the player. Closed tickets cannot be answered again.
func (s *Service) Answer(ctx context.Context, admin, ticket uuid.UUID, reply string, closeIt bool) (*Ticket, error) {
	reply = strings.TrimSpace(reply)
	if reply == "" || utf8.RuneCountInString(reply) > MaxBody {
		return nil, ErrReplyRequired
	}
	state := "answered"
	if closeIt {
		state = "closed"
	}
	rows, err := s.pool.Query(ctx, `
		UPDATE ops.support_ticket
		   SET reply = $3, replied_by = $2, replied_at = now(), state = $4::text, updated_at = now(),
		       closed_at = CASE WHEN $4::text = 'closed' THEN now() END
		 WHERE id = $1 AND state <> 'closed'
		RETURNING `+cols, ticket, admin, reply, state)
	if err != nil {
		return nil, err
	}
	t, err := pgx.CollectOneRow(rows, scan)
	if errors.Is(err, pgx.ErrNoRows) {
		var exists bool
		if qerr := s.pool.QueryRow(ctx, `SELECT EXISTS (SELECT 1 FROM ops.support_ticket WHERE id = $1)`, ticket).
			Scan(&exists); qerr != nil {
			return nil, qerr
		}
		if exists {
			return nil, ErrClosed
		}
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	notification.Send(ctx, s.notifier, s.log, notification.Notice{
		AccountID: t.AccountID, Kind: "support_reply", DedupKey: "support_reply:" + t.ID.String() + ":" + time.Now().UTC().Format(time.RFC3339Nano),
		TitleKey: "notif.support_reply.title", BodyKey: "notif.support_reply.body",
		Vars: map[string]any{"subject": t.Subject}, Route: "/support", Category: "system",
	})
	return &t, nil
}

const cols = `id, account_id, character_id, category, subject, body, client_info, state, reply, replied_at, closed_at, created_at`

func scan(row pgx.CollectableRow) (Ticket, error) {
	var t Ticket
	err := row.Scan(&t.ID, &t.AccountID, &t.CharacterID, &t.Category, &t.Subject, &t.Body, &t.ClientInfo,
		&t.State, &t.Reply, &t.RepliedAt, &t.ClosedAt, &t.CreatedAt)
	return t, err
}
