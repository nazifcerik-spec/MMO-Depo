package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"log/slog"
	"net/mail"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
)

// ServiceConfig tunes the auth service. Clock injection keeps tests
// deterministic. Argon2Params are per-process; changing them re-hashes on next
// successful login (opt Phase 02).
type ServiceConfig struct {
	Argon2Params Argon2Params
	AccessTTL    time.Duration
	RefreshTTL   time.Duration
	SessionTTL   time.Duration
	Now          func() time.Time // clock injection for tests

	// MinPasswordLen: hard minimum per NIST — no complexity mandate.
	MinPasswordLen int

	// PublicBaseURL prefixes links in recovery e-mails (DD-219).
	PublicBaseURL string

	// TermsVersion is recorded when a player accepts the terms at sign-up
	// (DD-231); bump it when the published terms change.
	TermsVersion string
}

// DefaultServiceConfig returns Faz 00 baseline values.
func DefaultServiceConfig() ServiceConfig {
	return ServiceConfig{
		Argon2Params:   DefaultArgon2Params,
		AccessTTL:      15 * time.Minute,
		RefreshTTL:     30 * 24 * time.Hour,
		SessionTTL:     30 * 24 * time.Hour,
		Now:            time.Now,
		MinPasswordLen: 8,
		TermsVersion:   DefaultTermsVersion,
	}
}

// Service is the auth business logic. Depends only on the Repository interface
// (so tests use MemoryRepo) and the Signer (Ed25519 JWT).
type Service struct {
	repo   Repository
	signer *Signer
	cfg    ServiceConfig
	log    *slog.Logger
	// dummyHash burns the same argon2id cost as a real account on
	// unknown-account logins, so timing does not reveal which names exist.
	// Built from cfg's parameters: a fixed-cost dummy would drift from them.
	dummyHash string
}

func NewService(repo Repository, signer *Signer, cfg ServiceConfig, logger *slog.Logger) *Service {
	if cfg.Now == nil {
		cfg.Now = time.Now
	}
	if cfg.TermsVersion == "" {
		cfg.TermsVersion = DefaultTermsVersion
	}
	if logger == nil {
		logger = slog.Default()
	}
	return &Service{repo: repo, signer: signer, cfg: cfg, log: logger, dummyHash: dummyHashFor(cfg.Argon2Params)}
}

// --- Requests / results ------------------------------------------------------

// RegisterRequest is the sanitized register payload.
type RegisterRequest struct {
	Email     string
	Username  string
	Password  string
	IPHash    []byte
	UserAgent string
	// AcceptTerms: the player confirmed being 13+ and accepted the terms and
	// privacy policy (DD-231). Required.
	AcceptTerms bool
}

// LoginRequest is the sanitized login payload.
type LoginRequest struct {
	EmailOrUsername string
	Password        string
	IPHash          []byte
	UserAgent       string
}

// --- Register ---------------------------------------------------------------

var (
	usernamePattern = regexp.MustCompile(`^[a-zA-Z0-9_.]{3,32}$`)
)

func (s *Service) Register(ctx context.Context, req RegisterRequest) (*LoginResult, error) {
	email, err := validateEmail(req.Email)
	if err != nil {
		return nil, err
	}
	username, err := validateUsername(req.Username)
	if err != nil {
		return nil, err
	}
	if err := validatePassword(req.Password, s.cfg.MinPasswordLen); err != nil {
		return nil, err
	}
	if !req.AcceptTerms {
		return nil, &ValidationError{Field: "accept_terms", Message: "confirm you are 13+ and accept the terms"}
	}

	hash, err := HashPassword(req.Password, s.cfg.Argon2Params)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	now := s.cfg.Now().UTC()
	account := Account{
		ID:                 uuid.New(),
		Email:              email,
		EmailNormalized:    strings.ToLower(email),
		Username:           username,
		UsernameNormalized: strings.ToLower(username),
		PasswordHash:       hash,
		State:              AccountActive,
		CreatedAt:          now,
		UpdatedAt:          now,
		Version:            1,
		TermsVersion:       s.cfg.TermsVersion,
		TermsAcceptedAt:    &now,
	}
	if err := s.repo.CreateAccount(ctx, account); err != nil {
		return nil, err
	}
	// Verification mail is best effort: the account works regardless and the
	// player can resend from the app (DD-219).
	if err := s.issueToken(ctx, &account, TokenEmailVerify, now); err != nil {
		s.log.WarnContext(ctx, "auth: verification mail not queued", "err", err.Error())
	}

	return s.issueSession(ctx, account, req.IPHash, req.UserAgent, now)
}

// --- Login ------------------------------------------------------------------

func (s *Service) Login(ctx context.Context, req LoginRequest) (*LoginResult, error) {
	if req.EmailOrUsername == "" {
		return nil, &ValidationError{Field: "email_or_username", Message: "required"}
	}
	if req.Password == "" {
		return nil, &ValidationError{Field: "password", Message: "required"}
	}

	// Try email first, then username (both normalized lowercase).
	needle := strings.ToLower(strings.TrimSpace(req.EmailOrUsername))
	var account *Account
	var err error
	if strings.Contains(needle, "@") {
		account, err = s.repo.AccountByEmailNormalized(ctx, needle)
	} else {
		account, err = s.repo.AccountByUsernameNormalized(ctx, needle)
	}
	if err != nil {
		// Constant-time-ish: still spend an Argon2 verify against a dummy hash
		// so the response time doesn't leak "unknown account" vs "wrong pw".
		if errors.Is(VerifyPassword(req.Password, s.dummyHash), ErrBusy) {
			return nil, ErrBusy // same answer as a known account under load
		}
		return nil, ErrInvalidCredentials
	}

	switch account.State {
	case AccountLocked:
		return nil, ErrAccountLocked
	case AccountBanned:
		return nil, ErrAccountBanned
	case AccountDeleted:
		return nil, ErrInvalidCredentials
	}

	if err := VerifyPassword(req.Password, account.PasswordHash); err != nil {
		if errors.Is(err, ErrBusy) {
			return nil, err
		}
		return nil, ErrInvalidCredentials
	}

	now := s.cfg.Now().UTC()
	_ = s.repo.UpdateAccountLastLogin(ctx, account.ID, now)

	// Signing in during the deletion grace cancels the deletion (DD-231).
	cancelled := false
	if account.DeletionRequestedAt != nil {
		if err := s.repo.SetDeletionRequested(ctx, account.ID, nil); err != nil {
			return nil, err
		}
		cancelled = true
	}
	res, err := s.issueSession(ctx, *account, req.IPHash, req.UserAgent, now)
	if res != nil {
		res.DeletionCancelled = cancelled
	}
	return res, err
}

// --- Refresh ----------------------------------------------------------------

func (s *Service) Refresh(ctx context.Context, refreshToken string) (*LoginResult, error) {
	if refreshToken == "" {
		return nil, &ValidationError{Field: "refresh_token", Message: "required"}
	}
	hash := sha256.Sum256([]byte(refreshToken))
	tok, err := s.repo.RefreshTokenByHash(ctx, hash[:])
	if err != nil {
		return nil, err
	}

	now := s.cfg.Now().UTC()

	// Revoked? treat like not-found — attacker gets nothing back.
	if tok.RevokedAt != nil {
		return nil, ErrSessionRevoked
	}

	// Expired?
	if tok.ExpiresAt.Before(now) {
		return nil, ErrTokenExpired
	}

	// REUSE DETECTION — DD-189. If used_at is set, this is the compromise
	// path: attacker replayed a rotated token. Revoke the whole session tree.
	if tok.UsedAt != nil {
		_ = s.repo.RevokeSession(ctx, tok.SessionID, "refresh_reuse", now)
		_ = s.repo.RevokeRefreshTokensForSession(ctx, tok.SessionID, "refresh_reuse", now)
		s.log.WarnContext(ctx, "refresh token reuse detected — session revoked",
			slog.String("session_id", tok.SessionID.String()))
		return nil, ErrRefreshTokenReuse
	}

	// Validate session
	session, err := s.repo.Session(ctx, tok.SessionID)
	if err != nil {
		return nil, err
	}
	if session.RevokedAt != nil {
		return nil, ErrSessionRevoked
	}
	if session.ExpiresAt.Before(now) {
		return nil, ErrSessionExpired
	}

	account, err := s.repo.AccountByID(ctx, session.AccountID)
	if err != nil {
		return nil, err
	}
	switch account.State {
	case AccountLocked:
		return nil, ErrAccountLocked
	case AccountBanned:
		return nil, ErrAccountBanned
	}

	// Mark this token used, issue the next pair.
	if err := s.repo.MarkRefreshTokenUsed(ctx, tok.ID, now); err != nil {
		return nil, err
	}
	_ = s.repo.TouchSession(ctx, session.ID, now)

	return s.issueTokensForSession(ctx, *session, *account, &tok.ID, now)
}

// --- Logout -----------------------------------------------------------------

// Logout revokes the current session and all its refresh tokens. Idempotent —
// calling twice is safe.
func (s *Service) Logout(ctx context.Context, sessionID uuid.UUID) error {
	now := s.cfg.Now().UTC()
	if err := s.repo.RevokeSession(ctx, sessionID, "user_logout", now); err != nil {
		return err
	}
	return s.repo.RevokeRefreshTokensForSession(ctx, sessionID, "user_logout", now)
}

// --- Verify (used by auth middleware) ---------------------------------------

// Verify parses the access token, checks the session is still live, and returns
// the resolved account + session ids. Cheap enough for every request; session
// lookup can be cached by the caller if needed.
func (s *Service) Verify(ctx context.Context, accessToken string) (accountID, sessionID uuid.UUID, err error) {
	claims, err := s.signer.VerifyAccessToken(accessToken)
	if err != nil {
		return uuid.Nil, uuid.Nil, err
	}
	session, err := s.repo.Session(ctx, claims.SessionID)
	if err != nil {
		return uuid.Nil, uuid.Nil, err
	}
	if session.RevokedAt != nil {
		return uuid.Nil, uuid.Nil, ErrSessionRevoked
	}
	if session.ExpiresAt.Before(s.cfg.Now().UTC()) {
		return uuid.Nil, uuid.Nil, ErrSessionExpired
	}
	return claims.AccountID, claims.SessionID, nil
}

// --- helpers ----------------------------------------------------------------

// issueSession creates a fresh session + refresh token + access token.
func (s *Service) issueSession(ctx context.Context, a Account, ipHash []byte, ua string, now time.Time) (*LoginResult, error) {
	session := Session{
		ID:           uuid.New(),
		AccountID:    a.ID,
		IPHash:       ipHash,
		UserAgent:    ua,
		CreatedAt:    now,
		LastActiveAt: now,
		ExpiresAt:    now.Add(s.cfg.SessionTTL),
	}
	if err := s.repo.CreateSession(ctx, session); err != nil {
		return nil, err
	}
	return s.issueTokensForSession(ctx, session, a, nil, now)
}

// issueTokensForSession mints a new refresh token linked to session, plus an
// access JWT. When parentID is non-nil the new token is chained to the caller
// (rotation from Refresh); nil means it is the root token of the session.
func (s *Service) issueTokensForSession(ctx context.Context, session Session, account Account, parentID *uuid.UUID, now time.Time) (*LoginResult, error) {
	refreshPlain, err := generateRefreshToken()
	if err != nil {
		return nil, err
	}
	refreshHash := sha256.Sum256([]byte(refreshPlain))
	refreshExpires := now.Add(s.cfg.RefreshTTL)

	rt := RefreshToken{
		ID:            uuid.New(),
		SessionID:     session.ID,
		TokenHash:     refreshHash[:],
		ParentTokenID: parentID,
		IssuedAt:      now,
		ExpiresAt:     refreshExpires,
	}
	if err := s.repo.CreateRefreshToken(ctx, rt); err != nil {
		return nil, err
	}

	accessToken, accessExp, err := s.signer.SignAccessToken(account.ID, session.ID, now)
	if err != nil {
		return nil, err
	}

	return &LoginResult{
		AccountID:       account.ID,
		SessionID:       session.ID,
		AccessToken:     accessToken,
		AccessExpiresAt: accessExp,
		RefreshToken:    refreshPlain,
		RefreshExpires:  refreshExpires,
		EmailVerified:   account.EmailVerifiedAt != nil,
	}, nil
}

// generateRefreshToken returns 32 random bytes base64url-encoded (43 chars,
// unpadded). URL-safe so clients can carry it in headers or query without
// escaping.
func generateRefreshToken() (string, error) {
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		return "", fmt.Errorf("rand: %w", err)
	}
	return base64.RawURLEncoding.EncodeToString(buf), nil
}

// --- validation -------------------------------------------------------------

func validateEmail(raw string) (string, error) {
	s := strings.TrimSpace(raw)
	if s == "" {
		return "", &ValidationError{Field: "email", Message: "required"}
	}
	if len(s) > 320 {
		return "", &ValidationError{Field: "email", Message: "too long"}
	}
	if _, err := mail.ParseAddress(s); err != nil {
		return "", &ValidationError{Field: "email", Message: "invalid format"}
	}
	return s, nil
}

func validateUsername(raw string) (string, error) {
	s := strings.TrimSpace(raw)
	if s == "" {
		return "", &ValidationError{Field: "username", Message: "required"}
	}
	if !usernamePattern.MatchString(s) {
		return "", &ValidationError{Field: "username", Message: "3-32 chars, letters/digits/_/. only"}
	}
	return s, nil
}

func validatePassword(pw string, minLen int) error {
	if len(pw) < minLen {
		return &ValidationError{Field: "password", Message: fmt.Sprintf("must be at least %d chars", minLen)}
	}
	if len(pw) > 4096 {
		return &ValidationError{Field: "password", Message: "too long"}
	}
	return nil
}
