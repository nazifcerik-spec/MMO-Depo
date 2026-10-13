package notification

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Notice is one event to deliver. DedupKey makes delivery idempotent
// (e.g. "activity_completed:<contract_id>").
type Notice struct {
	AccountID   uuid.UUID
	CharacterID *uuid.UUID
	Kind        string
	DedupKey    string
	TitleKey    string
	BodyKey     string
	Vars        map[string]any
	Route       string
	Category    string // android channel / iOS category: activity|event|social|system
	CollapseKey string
	High        bool
	Read        bool // stored already read: a history entry, not an alert
}

// Notifier is the narrow interface producers depend on.
type Notifier interface {
	Notify(ctx context.Context, n Notice) error
}

// Send is a nil-safe helper for producers: failures are logged, never
// returned, because a missed notification must not fail the game action.
func Send(ctx context.Context, n Notifier, log *slog.Logger, notice Notice) {
	if n == nil {
		return
	}
	if err := n.Notify(ctx, notice); err != nil && log != nil {
		log.WarnContext(ctx, "notification: notify failed", "kind", notice.Kind, "err", err.Error())
	}
}

// Payload is the provider-agnostic body stored on push/inbox rows.
type Payload struct {
	Kind        string         `json:"kind"`
	TitleKey    string         `json:"title_key"`
	BodyKey     string         `json:"body_key"`
	Vars        map[string]any `json:"body_vars,omitempty"`
	Category    string         `json:"category"`
	Route       string         `json:"route,omitempty"`
	CollapseKey string         `json:"collapse_key,omitempty"`
}

var (
	ErrUnknownKind  = errors.New("notification: unknown kind")
	ErrBadMode      = errors.New("notification: bad mode")
	ErrBadPlatform  = errors.New("notification: platform must be android or ios")
	ErrBadToken     = errors.New("notification: push_token required")
	ErrBadClock     = errors.New("notification: quiet hours must be HH:MM")
	ErrBadTimezone  = errors.New("notification: unknown timezone")
	ErrBadMuteHours = errors.New("notification: mute hours must be 1..168")
	ErrNotFound     = errors.New("notification: not found")
)

type Service struct {
	pool *pgxpool.Pool
	now  func() time.Time
	log  *slog.Logger
}

func NewService(pool *pgxpool.Pool, logger *slog.Logger) *Service {
	if logger == nil {
		logger = slog.Default()
	}
	return &Service{pool: pool, now: time.Now, log: logger}
}

// Notify applies preferences and rate limits, then writes the inbox row
// and/or queues the push. Safe to call repeatedly with the same DedupKey.
func (s *Service) Notify(ctx context.Context, n Notice) error {
	if _, ok := DefaultModes[n.Kind]; !ok {
		return ErrUnknownKind
	}
	if n.DedupKey == "" {
		n.DedupKey = n.Kind + ":" + uuid.NewString()
	}
	if n.Category == "" {
		n.Category = "system"
	}
	now := s.now().UTC()

	mode := DefaultModes[n.Kind]
	var custom *string
	var qStart, qEnd, qTZ *string
	var muteUntil *time.Time
	err := s.pool.QueryRow(ctx, `
		SELECT (SELECT mode FROM notification.preference_kind WHERE account_id = $1 AND kind = $2),
		       p.quiet_hours_start, p.quiet_hours_end, p.quiet_hours_tz, p.global_mute_until
		  FROM (SELECT 1) one
		  LEFT JOIN notification.preference p ON p.account_id = $1`, n.AccountID, n.Kind).
		Scan(&custom, &qStart, &qEnd, &qTZ, &muteUntil)
	if err != nil {
		return err
	}
	if custom != nil {
		mode = Mode(*custom)
	}
	muted := muteUntil != nil && muteUntil.After(now)
	quiet := InQuietHours(now, deref(qStart), deref(qEnd), deref(qTZ))
	d := Decide(n.Kind, mode, muted, quiet, n.High)

	body, err := json.Marshal(Payload{
		Kind: n.Kind, TitleKey: n.TitleKey, BodyKey: n.BodyKey, Vars: n.Vars,
		Category: n.Category, Route: n.Route, CollapseKey: n.CollapseKey,
	})
	if err != nil {
		return err
	}
	if d.Inbox {
		if _, err := s.pool.Exec(ctx, `
			INSERT INTO notification.inbox_message (account_id, character_id, kind, dedup_key, body, route, read_at)
			VALUES ($1, $2, $3, $4, $5, NULLIF($6, ''), CASE WHEN $7::boolean THEN now() END)
			ON CONFLICT (account_id, dedup_key) DO NOTHING`,
			n.AccountID, n.CharacterID, n.Kind, n.DedupKey, body, n.Route, n.Read); err != nil {
			return err
		}
	}
	if !d.Push {
		return nil
	}
	if n.Kind != KindSystemBroadcast {
		var sent int
		var last *time.Time
		if err := s.pool.QueryRow(ctx, `
			SELECT count(*), max(created_at) FROM notification.push_message
			 WHERE account_id = $1 AND kind = $2 AND created_at > $3`,
			n.AccountID, n.Kind, now.Add(-time.Hour)).Scan(&sent, &last); err != nil {
			return err
		}
		if !RateAllowed(n.Kind, sent, last, now) {
			s.log.DebugContext(ctx, "notification: push rate limited", "kind", n.Kind)
			return nil
		}
	}
	_, err = s.pool.Exec(ctx, `
		INSERT INTO notification.push_message (account_id, character_id, kind, dedup_key, priority, payload)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (account_id, dedup_key) DO NOTHING`,
		n.AccountID, n.CharacterID, n.Kind, n.DedupKey, d.Priority, body)
	return err
}

func deref(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// --- devices ----------------------------------------------------------------

type DeviceRequest struct {
	DeviceID   uuid.UUID `json:"device_id"`
	Platform   string    `json:"platform"`
	PushToken  string    `json:"push_token"`
	Locale     string    `json:"locale"`
	AppVersion string    `json:"app_version"`
}

// RegisterDevice upserts the device token; re-registration reactivates it
// and resets the failure counter.
func (s *Service) RegisterDevice(ctx context.Context, account uuid.UUID, r DeviceRequest) error {
	if r.Platform != "android" && r.Platform != "ios" {
		return ErrBadPlatform
	}
	if r.PushToken == "" || len(r.PushToken) > 4096 || r.DeviceID == uuid.Nil {
		return ErrBadToken
	}
	_, err := s.pool.Exec(ctx, `
		INSERT INTO notification.device_push_token (account_id, device_id, platform, push_token, locale, app_version)
		VALUES ($1, $2, $3, $4, NULLIF($5, ''), NULLIF($6, ''))
		ON CONFLICT (device_id, platform) DO UPDATE SET
		  account_id = EXCLUDED.account_id, push_token = EXCLUDED.push_token,
		  locale = EXCLUDED.locale, app_version = EXCLUDED.app_version,
		  is_active = true, failure_count = 0, last_used_at = now()`,
		account, r.DeviceID, r.Platform, r.PushToken, r.Locale, r.AppVersion)
	return err
}

func (s *Service) UnregisterDevice(ctx context.Context, account, deviceID uuid.UUID) error {
	tag, err := s.pool.Exec(ctx, `
		UPDATE notification.device_push_token SET is_active = false
		 WHERE account_id = $1 AND device_id = $2`, account, deviceID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

// --- inbox ------------------------------------------------------------------

type InboxItem struct {
	ID          uuid.UUID       `json:"id"`
	CharacterID *uuid.UUID      `json:"character_id,omitempty"`
	Kind        string          `json:"kind"`
	Body        json.RawMessage `json:"body"`
	Route       *string         `json:"route,omitempty"`
	CreatedAt   time.Time       `json:"created_at"`
	ReadAt      *time.Time      `json:"read_at,omitempty"`
}

// Inbox lists the newest 100 entries. Reward history (KindReward) is its own
// list so frequent collections never push other notifications out (DD-242).
func (s *Service) Inbox(ctx context.Context, account uuid.UUID, unreadOnly, rewards bool) ([]InboxItem, int, error) {
	rows, err := s.pool.Query(ctx, `
		SELECT id, character_id, kind, body, route, created_at, read_at
		  FROM notification.inbox_message
		 WHERE account_id = $1 AND expires_at > now() AND ($2 = false OR read_at IS NULL)
		   AND (kind = $4) = $3
		 ORDER BY created_at DESC LIMIT 100`, account, unreadOnly, rewards, KindReward)
	if err != nil {
		return nil, 0, err
	}
	items, err := pgx.CollectRows(rows, func(r pgx.CollectableRow) (InboxItem, error) {
		var it InboxItem
		err := r.Scan(&it.ID, &it.CharacterID, &it.Kind, &it.Body, &it.Route, &it.CreatedAt, &it.ReadAt)
		return it, err
	})
	if err != nil {
		return nil, 0, err
	}
	var unread int
	err = s.pool.QueryRow(ctx, `
		SELECT count(*) FROM notification.inbox_message
		 WHERE account_id = $1 AND read_at IS NULL AND expires_at > now()`, account).Scan(&unread)
	return items, unread, err
}

// MarkRead marks one inbox item (or all when id is uuid.Nil) as read.
func (s *Service) MarkRead(ctx context.Context, account, id uuid.UUID) (int64, error) {
	tag, err := s.pool.Exec(ctx, `
		UPDATE notification.inbox_message SET read_at = now()
		 WHERE account_id = $1 AND read_at IS NULL AND ($2 = '00000000-0000-0000-0000-000000000000'::uuid OR id = $2)`,
		account, id)
	return tag.RowsAffected(), err
}

// --- preferences ------------------------------------------------------------

type Preferences struct {
	QuietHoursStart *string         `json:"quiet_hours_start"`
	QuietHoursEnd   *string         `json:"quiet_hours_end"`
	QuietHoursTZ    *string         `json:"quiet_hours_tz"`
	GlobalMuteUntil *time.Time      `json:"global_mute_until"`
	Kinds           map[string]Mode `json:"kinds"`
}

func (s *Service) Preferences(ctx context.Context, account uuid.UUID) (*Preferences, error) {
	p := &Preferences{Kinds: map[string]Mode{}}
	for k, m := range DefaultModes {
		p.Kinds[k] = m
	}
	err := s.pool.QueryRow(ctx, `
		SELECT quiet_hours_start, quiet_hours_end, quiet_hours_tz, global_mute_until
		  FROM notification.preference WHERE account_id = $1`, account).
		Scan(&p.QuietHoursStart, &p.QuietHoursEnd, &p.QuietHoursTZ, &p.GlobalMuteUntil)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return nil, err
	}
	rows, err := s.pool.Query(ctx, `SELECT kind, mode FROM notification.preference_kind WHERE account_id = $1`, account)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var k, m string
		if err := rows.Scan(&k, &m); err != nil {
			return nil, err
		}
		if _, known := DefaultModes[k]; known {
			p.Kinds[k] = Mode(m)
		}
	}
	return p, rows.Err()
}

type PreferenceUpdate struct {
	QuietHoursStart *string         `json:"quiet_hours_start"`
	QuietHoursEnd   *string         `json:"quiet_hours_end"`
	QuietHoursTZ    *string         `json:"quiet_hours_tz"`
	Kinds           map[string]Mode `json:"kinds"`
}

// UpdatePreferences replaces quiet hours (all three or none) and upserts the
// given kind modes. system_broadcast cannot be changed.
func (s *Service) UpdatePreferences(ctx context.Context, account uuid.UUID, u PreferenceUpdate) (*Preferences, error) {
	start, end, tz := deref(u.QuietHoursStart), deref(u.QuietHoursEnd), deref(u.QuietHoursTZ)
	if start != "" || end != "" {
		if _, err := ParseClock(start); err != nil {
			return nil, ErrBadClock
		}
		if _, err := ParseClock(end); err != nil {
			return nil, ErrBadClock
		}
	}
	if tz != "" {
		if _, err := time.LoadLocation(tz); err != nil {
			return nil, ErrBadTimezone
		}
	}
	for k, m := range u.Kinds {
		if _, ok := DefaultModes[k]; !ok || k == KindSystemBroadcast {
			return nil, ErrUnknownKind
		}
		if !m.Valid() {
			return nil, ErrBadMode
		}
	}
	err := pgx.BeginFunc(ctx, s.pool, func(tx pgx.Tx) error {
		if _, err := tx.Exec(ctx, `
			INSERT INTO notification.preference (account_id, quiet_hours_start, quiet_hours_end, quiet_hours_tz)
			VALUES ($1, NULLIF($2, ''), NULLIF($3, ''), NULLIF($4, ''))
			ON CONFLICT (account_id) DO UPDATE SET
			  quiet_hours_start = EXCLUDED.quiet_hours_start, quiet_hours_end = EXCLUDED.quiet_hours_end,
			  quiet_hours_tz = EXCLUDED.quiet_hours_tz, updated_at = now(),
			  version = notification.preference.version + 1`,
			account, start, end, tz); err != nil {
			return err
		}
		for k, m := range u.Kinds {
			if _, err := tx.Exec(ctx, `
				INSERT INTO notification.preference_kind (account_id, kind, mode) VALUES ($1, $2, $3)
				ON CONFLICT (account_id, kind) DO UPDATE SET mode = EXCLUDED.mode, updated_at = now()`,
				account, k, string(m)); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return s.Preferences(ctx, account)
}

// SetGlobalMute silences pushes for `hours` (0 clears). Inbox still fills.
func (s *Service) SetGlobalMute(ctx context.Context, account uuid.UUID, hours int) (*time.Time, error) {
	if hours < 0 || hours > 168 {
		return nil, ErrBadMuteHours
	}
	var until *time.Time
	if hours > 0 {
		t := s.now().UTC().Add(time.Duration(hours) * time.Hour)
		until = &t
	}
	_, err := s.pool.Exec(ctx, `
		INSERT INTO notification.preference (account_id, global_mute_until) VALUES ($1, $2)
		ON CONFLICT (account_id) DO UPDATE SET global_mute_until = EXCLUDED.global_mute_until,
		  updated_at = now(), version = notification.preference.version + 1`, account, until)
	return until, err
}
