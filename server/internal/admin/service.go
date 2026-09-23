package admin

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"log/slog"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/otter/server/internal/domain/auth"
	"github.com/otter/server/internal/economy"
)

const (
	SessionTTL = 8 * time.Hour
	StepUpTTL  = 5 * time.Minute
)

var (
	ErrBadCredentials = errors.New("admin: invalid credentials")
	ErrUnauthorized   = errors.New("admin: not authenticated")
	ErrForbidden      = errors.New("admin: permission denied")
	ErrStepUpRequired = errors.New("admin: step-up verification required")
	ErrBadCode        = errors.New("admin: invalid or reused code")
	ErrReasonRequired = errors.New("admin: reason required")
	ErrNotFound       = errors.New("admin: not found")
	ErrInvalid        = errors.New("admin: invalid request")
)

// Admin is the authenticated operator for one request.
type Admin struct {
	ID          uuid.UUID     `json:"id"`
	Email       string        `json:"email"`
	DisplayName string        `json:"display_name"`
	Roles       []string      `json:"roles"`
	Perms       PermissionSet `json:"-"`
	SessionID   uuid.UUID     `json:"-"`
	StepUpUntil *time.Time    `json:"stepup_until,omitempty"`
}

func (a *Admin) SteppedUp(now time.Time) bool {
	return a.StepUpUntil != nil && a.StepUpUntil.After(now)
}

type Service struct {
	pool *pgxpool.Pool
	txn  economy.TxnService
	now  func() time.Time
	log  *slog.Logger
}

func NewService(pool *pgxpool.Pool, txn economy.TxnService, logger *slog.Logger) *Service {
	if logger == nil {
		logger = slog.Default()
	}
	return &Service{pool: pool, txn: txn, now: time.Now, log: logger}
}

func hashToken(t string) []byte { h := sha256.Sum256([]byte(t)); return h[:] }

var (
	dummyOnce sync.Once
	dummyHash string
)

// --- account bootstrap (cmd/adminctl) ---------------------------------------

// CreateAdmin inserts an admin with a fresh TOTP secret and the given roles.
func (s *Service) CreateAdmin(ctx context.Context, email, name, password string, roles []string) (uuid.UUID, string, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	if email == "" || len(password) < 12 {
		return uuid.Nil, "", ErrInvalid
	}
	for _, r := range roles {
		if !ValidRole(r) {
			return uuid.Nil, "", ErrInvalid
		}
	}
	hash, err := auth.HashPassword(password, auth.DefaultArgon2Params)
	if err != nil {
		return uuid.Nil, "", err
	}
	secret, err := GenerateTOTPSecret()
	if err != nil {
		return uuid.Nil, "", err
	}
	var id uuid.UUID
	err = pgx.BeginFunc(ctx, s.pool, func(tx pgx.Tx) error {
		if err := tx.QueryRow(ctx, `
			INSERT INTO ops.admin_user (email, display_name, password_hash, totp_secret)
			VALUES ($1, $2, $3, $4) RETURNING id`, email, name, hash, secret).Scan(&id); err != nil {
			return err
		}
		for _, r := range roles {
			if _, err := tx.Exec(ctx, `
				INSERT INTO ops.admin_user_role (admin_user_id, role_code) VALUES ($1, $2)`, id, r); err != nil {
				return err
			}
		}
		return nil
	})
	return id, secret, err
}

// --- sessions ---------------------------------------------------------------

// consumeTOTP verifies a code for the admin and records its step so the same
// code cannot be used twice.
func (s *Service) consumeTOTP(ctx context.Context, adminID uuid.UUID, code string) error {
	return pgx.BeginFunc(ctx, s.pool, func(tx pgx.Tx) error {
		var secret string
		var last int64
		if err := tx.QueryRow(ctx, `
			SELECT totp_secret, last_totp_step FROM ops.admin_user WHERE id = $1 FOR UPDATE`, adminID).
			Scan(&secret, &last); err != nil {
			return err
		}
		step, ok := VerifyTOTP(secret, code, s.now(), last)
		if !ok {
			return ErrBadCode
		}
		_, err := tx.Exec(ctx, `UPDATE ops.admin_user SET last_totp_step = $2 WHERE id = $1`, adminID, step)
		return err
	})
}

// Login checks password + TOTP and returns an opaque session token.
func (s *Service) Login(ctx context.Context, email, password, code, userAgent string, ipHash []byte) (string, error) {
	var id uuid.UUID
	var hash string
	var active bool
	err := s.pool.QueryRow(ctx, `
		SELECT id, password_hash, is_active FROM ops.admin_user WHERE email = $1`,
		strings.ToLower(strings.TrimSpace(email))).Scan(&id, &hash, &active)
	if errors.Is(err, pgx.ErrNoRows) {
		// Burn comparable time so account existence does not leak via timing.
		dummyOnce.Do(func() { dummyHash, _ = auth.HashPassword("otter-dummy", auth.DefaultArgon2Params) })
		_ = auth.VerifyPassword(password, dummyHash)
		return "", ErrBadCredentials
	}
	if err != nil {
		return "", err
	}
	if !active || auth.VerifyPassword(password, hash) != nil {
		return "", ErrBadCredentials
	}
	if err := s.consumeTOTP(ctx, id, code); err != nil {
		if errors.Is(err, ErrBadCode) {
			return "", ErrBadCredentials
		}
		return "", err
	}
	raw := make([]byte, 32)
	if _, err := rand.Read(raw); err != nil {
		return "", err
	}
	token := base64.RawURLEncoding.EncodeToString(raw)
	if len(userAgent) > 256 {
		userAgent = userAgent[:256]
	}
	if _, err := s.pool.Exec(ctx, `
		INSERT INTO ops.admin_session (admin_user_id, token_hash, expires_at, ip_hash, user_agent)
		VALUES ($1, $2, $3, $4, $5)`, id, hashToken(token), s.now().Add(SessionTTL), ipHash, userAgent); err != nil {
		return "", err
	}
	s.Audit(ctx, &Admin{ID: id}, AuditEntry{Action: "admin.login"})
	return token, nil
}

// Authenticate resolves a bearer token to the admin and their permissions.
func (s *Service) Authenticate(ctx context.Context, token string) (*Admin, error) {
	if token == "" {
		return nil, ErrUnauthorized
	}
	a := &Admin{}
	err := s.pool.QueryRow(ctx, `
		SELECT u.id, u.email, u.display_name, se.id, se.stepup_until,
		       COALESCE(array_agg(r.role_code) FILTER (WHERE r.role_code IS NOT NULL), '{}')
		  FROM ops.admin_session se
		  JOIN ops.admin_user u ON u.id = se.admin_user_id AND u.is_active
		  LEFT JOIN ops.admin_user_role r ON r.admin_user_id = u.id
		 WHERE se.token_hash = $1 AND se.revoked_at IS NULL AND se.expires_at > now()
		 GROUP BY u.id, se.id`, hashToken(token)).
		Scan(&a.ID, &a.Email, &a.DisplayName, &a.SessionID, &a.StepUpUntil, &a.Roles)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrUnauthorized
	}
	if err != nil {
		return nil, err
	}
	a.Perms = Permissions(a.Roles)
	return a, nil
}

func (s *Service) Logout(ctx context.Context, a *Admin) error {
	_, err := s.pool.Exec(ctx, `UPDATE ops.admin_session SET revoked_at = now() WHERE id = $1`, a.SessionID)
	s.Audit(ctx, a, AuditEntry{Action: "admin.logout"})
	return err
}

// StepUp re-verifies a fresh TOTP code and opens a 5-minute window for
// critical actions on this session.
func (s *Service) StepUp(ctx context.Context, a *Admin, code string) (time.Time, error) {
	if err := s.consumeTOTP(ctx, a.ID, code); err != nil {
		return time.Time{}, err
	}
	until := s.now().Add(StepUpTTL)
	if _, err := s.pool.Exec(ctx, `UPDATE ops.admin_session SET stepup_until = $2 WHERE id = $1`, a.SessionID, until); err != nil {
		return time.Time{}, err
	}
	a.StepUpUntil = &until
	s.Audit(ctx, a, AuditEntry{Action: "admin.stepup"})
	return until, nil
}

// --- audit ------------------------------------------------------------------

type AuditEntry struct {
	Action     string
	TargetKind string
	TargetID   *uuid.UUID
	Before     any
	After      any
	Reason     string
	StepUp     bool
	RequestID  string
}

// Audit writes one audit.admin_action row. Failures are logged at error
// level: audit must never silently disappear.
func (s *Service) Audit(ctx context.Context, a *Admin, e AuditEntry) {
	marshal := func(v any) []byte {
		if v == nil {
			return nil
		}
		b, err := json.Marshal(v)
		if err != nil {
			return nil
		}
		return b
	}
	var stepAt *time.Time
	if e.StepUp && a.StepUpUntil != nil {
		t := a.StepUpUntil.Add(-StepUpTTL)
		stepAt = &t
	}
	if _, err := s.pool.Exec(ctx, `
		INSERT INTO audit.admin_action (admin_user_id, roles_at_time, action, target_kind, target_id,
		                                before_state, after_state, reason, requires_stepup, stepup_verified_at, request_id)
		VALUES ($1, $2, $3, NULLIF($4, ''), $5, $6, $7, NULLIF($8, ''), $9, $10, NULLIF($11, ''))`,
		a.ID, strings.Join(a.Roles, ","), e.Action, e.TargetKind, e.TargetID,
		marshal(e.Before), marshal(e.After), e.Reason, e.StepUp, stepAt, e.RequestID); err != nil {
		s.log.ErrorContext(ctx, "admin: audit write failed", "action", e.Action, "err", err.Error())
	}
}

type AuditRow struct {
	ID          int64           `json:"id"`
	AdminUserID uuid.UUID       `json:"admin_user_id"`
	Roles       string          `json:"roles_at_time"`
	Action      string          `json:"action"`
	TargetKind  *string         `json:"target_kind,omitempty"`
	TargetID    *uuid.UUID      `json:"target_id,omitempty"`
	Before      json.RawMessage `json:"before_state,omitempty"`
	After       json.RawMessage `json:"after_state,omitempty"`
	Reason      *string         `json:"reason,omitempty"`
	StepUp      bool            `json:"requires_stepup"`
	CreatedAt   time.Time       `json:"created_at"`
}

func (s *Service) AuditLog(ctx context.Context, actor, target *uuid.UUID, limit int) ([]AuditRow, error) {
	if limit <= 0 || limit > 200 {
		limit = 100
	}
	rows, err := s.pool.Query(ctx, `
		SELECT id, admin_user_id, roles_at_time, action, target_kind, target_id,
		       before_state, after_state, reason, requires_stepup, created_at
		  FROM audit.admin_action
		 WHERE ($1::uuid IS NULL OR admin_user_id = $1) AND ($2::uuid IS NULL OR target_id = $2)
		 ORDER BY created_at DESC LIMIT $3`, actor, target, limit)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, func(r pgx.CollectableRow) (AuditRow, error) {
		var a AuditRow
		err := r.Scan(&a.ID, &a.AdminUserID, &a.Roles, &a.Action, &a.TargetKind, &a.TargetID,
			&a.Before, &a.After, &a.Reason, &a.StepUp, &a.CreatedAt)
		return a, err
	})
}

// --- admin users ------------------------------------------------------------

type UserRow struct {
	ID          uuid.UUID `json:"id"`
	Email       string    `json:"email"`
	DisplayName string    `json:"display_name"`
	Active      bool      `json:"is_active"`
	Roles       []string  `json:"roles"`
}

func (s *Service) Users(ctx context.Context) ([]UserRow, error) {
	rows, err := s.pool.Query(ctx, `
		SELECT u.id, u.email, u.display_name, u.is_active,
		       COALESCE(array_agg(r.role_code ORDER BY r.role_code) FILTER (WHERE r.role_code IS NOT NULL), '{}')
		  FROM ops.admin_user u LEFT JOIN ops.admin_user_role r ON r.admin_user_id = u.id
		 GROUP BY u.id ORDER BY u.email`)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, func(r pgx.CollectableRow) (UserRow, error) {
		var u UserRow
		err := r.Scan(&u.ID, &u.Email, &u.DisplayName, &u.Active, &u.Roles)
		return u, err
	})
}

// SetRole grants or revokes a role. Admins cannot change their own roles.
func (s *Service) SetRole(ctx context.Context, a *Admin, target uuid.UUID, role string, grant bool, reason string) error {
	if !ValidRole(role) || target == a.ID {
		return ErrInvalid
	}
	var err error
	if grant {
		_, err = s.pool.Exec(ctx, `
			INSERT INTO ops.admin_user_role (admin_user_id, role_code, granted_by_admin_id)
			VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`, target, role, a.ID)
	} else {
		_, err = s.pool.Exec(ctx, `DELETE FROM ops.admin_user_role WHERE admin_user_id = $1 AND role_code = $2`, target, role)
	}
	if err != nil {
		return err
	}
	action := "admin.role_grant"
	if !grant {
		action = "admin.role_revoke"
	}
	s.Audit(ctx, a, AuditEntry{Action: action, TargetKind: "admin_user", TargetID: &target,
		After: map[string]any{"role": role}, Reason: reason, StepUp: true})
	return nil
}
