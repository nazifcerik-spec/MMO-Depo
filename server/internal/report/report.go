// Package report lets players report each other (DD-355). Operators work the
// queue from the admin panel; the reporter learns only that it was reviewed.
package report

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/otter/server/internal/domain/character"
	"github.com/otter/server/internal/http/apierr"
	"github.com/otter/server/internal/httpctx"
	"github.com/otter/server/internal/notification"
)

// MaxNote bounds what a player or an operator writes.
const MaxNote = 500

// Categories are the reasons a player can give.
var Categories = map[string]bool{"harassment": true, "cheating": true, "name": true, "spam": true, "scam": true, "other": true}

var (
	ErrInvalid  = errors.New("report: invalid report")
	ErrSelf     = errors.New("report: own character")
	ErrNotFound = errors.New("report: not found")
	ErrResolved = errors.New("report: already resolved")
)

// New is what a player submits.
type New struct {
	CharacterID   *uuid.UUID `json:"character_id,omitempty"`
	TargetName    string     `json:"target_name"`
	Category      string     `json:"category"`
	Note          string     `json:"note"`
	ChatMessageID *int64     `json:"chat_message_id,omitempty"`
}

// Report is one row of the operators' queue.
type Report struct {
	ID            uuid.UUID  `json:"id"`
	ReporterAcct  uuid.UUID  `json:"reporter_account"`
	ReporterName  *string    `json:"reporter_name,omitempty"`
	TargetAccount *uuid.UUID `json:"target_account,omitempty"`
	TargetID      *uuid.UUID `json:"target_character,omitempty"`
	TargetName    string     `json:"target_name"`
	Category      string     `json:"category"`
	Note          string     `json:"note"`
	ChatExcerpt   *string    `json:"chat_excerpt,omitempty"`
	State         string     `json:"state"`
	Resolution    *string    `json:"resolution_note,omitempty"`
	ResolvedAt    *time.Time `json:"resolved_at,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
	// OpenAgainst is how many open reports the target's account has.
	OpenAgainst int `json:"open_against"`
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

// Create files a report. The same reporter, target and reason while one is
// still open changes nothing and succeeds.
func (s *Service) Create(ctx context.Context, account uuid.UUID, req New) error {
	note := strings.TrimSpace(req.Note)
	if !Categories[req.Category] || utf8.RuneCountInString(note) > MaxNote || strings.TrimSpace(req.TargetName) == "" {
		return ErrInvalid
	}
	var target, targetAccount uuid.UUID
	var targetName string
	err := s.pool.QueryRow(ctx, `SELECT id, account_id, name FROM character.character WHERE name_normalized = $1`,
		character.NameKey(req.TargetName)).Scan(&target, &targetAccount, &targetName)
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrNotFound
	}
	if err != nil {
		return err
	}
	if targetAccount == account {
		return ErrSelf
	}
	if req.CharacterID != nil {
		var owner uuid.UUID
		err := s.pool.QueryRow(ctx, `SELECT account_id FROM character.character WHERE id = $1`, *req.CharacterID).Scan(&owner)
		if errors.Is(err, pgx.ErrNoRows) || (err == nil && owner != account) {
			return ErrNotFound
		}
		if err != nil {
			return err
		}
	}
	// The excerpt comes from our own log, never from the client, and only if
	// the reported character wrote it.
	var excerpt *string
	if req.ChatMessageID != nil {
		var body string
		err := s.pool.QueryRow(ctx, `SELECT body FROM social.chat_message WHERE id = $1 AND sender_character_id = $2`,
			*req.ChatMessageID, target).Scan(&body)
		if err == nil {
			excerpt = &body
		} else if !errors.Is(err, pgx.ErrNoRows) {
			return err
		}
	}
	_, err = s.pool.Exec(ctx, `
		INSERT INTO ops.player_report (reporter_account, reporter_character, target_account, target_character, target_name,
		                               category, note, chat_message_id, chat_excerpt)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		ON CONFLICT (reporter_account, target_character, category) WHERE state = 'open' DO NOTHING`,
		account, req.CharacterID, targetAccount, target, targetName, req.Category, note, req.ChatMessageID, excerpt)
	return err
}

const cols = `r.id, r.reporter_account, rc.name, r.target_account, r.target_character, r.target_name, r.category, r.note,
	r.chat_excerpt, r.state, r.resolution_note, r.resolved_at, r.created_at,
	(SELECT count(*) FROM ops.player_report o WHERE o.target_account = r.target_account AND o.state = 'open')`

func scan(row pgx.CollectableRow) (Report, error) {
	var r Report
	err := row.Scan(&r.ID, &r.ReporterAcct, &r.ReporterName, &r.TargetAccount, &r.TargetID, &r.TargetName, &r.Category, &r.Note,
		&r.ChatExcerpt, &r.State, &r.Resolution, &r.ResolvedAt, &r.CreatedAt, &r.OpenAgainst)
	return r, err
}

// Queue lists reports by state, oldest first.
func (s *Service) Queue(ctx context.Context, state string, limit int) ([]Report, error) {
	if state == "" {
		state = "open"
	}
	if limit <= 0 || limit > 200 {
		limit = 50
	}
	rows, err := s.pool.Query(ctx, `SELECT `+cols+`
		FROM ops.player_report r LEFT JOIN character.character rc ON rc.id = r.reporter_character
		WHERE r.state = $1 ORDER BY r.created_at LIMIT $2`, state, limit)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, scan)
}

// Resolve closes a report as actioned or dismissed and tells the reporter it
// was reviewed, not what happened.
func (s *Service) Resolve(ctx context.Context, admin, id uuid.UUID, state, note string) error {
	note = strings.TrimSpace(note)
	if (state != "actioned" && state != "dismissed") || utf8.RuneCountInString(note) > MaxNote {
		return ErrInvalid
	}
	var reporter uuid.UUID
	var target string
	err := s.pool.QueryRow(ctx, `
		UPDATE ops.player_report SET state = $2, resolved_by = $3, resolved_at = now(), resolution_note = $4
		 WHERE id = $1 AND state = 'open' RETURNING reporter_account, target_name`, id, state, admin, note).Scan(&reporter, &target)
	if errors.Is(err, pgx.ErrNoRows) {
		var exists bool
		if qerr := s.pool.QueryRow(ctx, `SELECT EXISTS (SELECT 1 FROM ops.player_report WHERE id = $1)`, id).Scan(&exists); qerr != nil {
			return qerr
		}
		if exists {
			return ErrResolved
		}
		return ErrNotFound
	}
	if err != nil {
		return err
	}
	notification.Send(ctx, s.notifier, s.log, notification.Notice{
		AccountID: reporter, Kind: "report_reviewed", DedupKey: "report_reviewed:" + id.String(),
		TitleKey: "notif.report_reviewed.title", BodyKey: "notif.report_reviewed.body",
		Vars: map[string]any{"target": target}, Category: "system",
	})
	return nil
}

// Handler serves the player side: POST /v1/reports.
type Handler struct{ svc *Service }

func NewHandler(svc *Service) *Handler { return &Handler{svc: svc} }

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	acct := httpctx.AccountID(r.Context())
	if acct == uuid.Nil {
		apierr.Write(w, r, http.StatusUnauthorized, "E_UNAUTHENTICATED", "auth required")
		return
	}
	var req New
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 8<<10)).Decode(&req); err != nil {
		apierr.Write(w, r, http.StatusBadRequest, "E_INVALID_JSON", err.Error())
		return
	}
	if err := h.svc.Create(r.Context(), acct, req); err != nil {
		WriteErr(w, r, err)
		return
	}
	apierr.WriteJSON(w, http.StatusAccepted, map[string]any{"status": "received"})
}

// WriteErr maps report errors (shared with the admin handlers).
func WriteErr(w http.ResponseWriter, r *http.Request, err error) {
	switch {
	case errors.Is(err, ErrInvalid):
		apierr.Write(w, r, http.StatusBadRequest, "E_VALIDATION", "a known category, a target and a note up to 500 characters")
	case errors.Is(err, ErrSelf):
		apierr.Write(w, r, http.StatusConflict, "E_REPORT_SELF", "cannot report your own character")
	case errors.Is(err, ErrNotFound):
		apierr.Write(w, r, http.StatusNotFound, "E_NOT_FOUND", "not found")
	case errors.Is(err, ErrResolved):
		apierr.Write(w, r, http.StatusConflict, "E_REPORT_RESOLVED", "already resolved")
	default:
		apierr.Write(w, r, http.StatusInternalServerError, "E_INTERNAL", err.Error())
	}
}
