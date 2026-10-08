package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/google/uuid"
)

// Account recovery (DD-219): password reset, e-mail verification, password
// change. Tokens are 32 random bytes, sent only by e-mail; the database keeps
// SHA-256(token). Consuming a token is a single conditional UPDATE, so a token
// works once even under concurrent clicks.

// TokenKind is the purpose of an account token.
type TokenKind string

const (
	TokenPasswordReset TokenKind = "password_reset"
	TokenEmailVerify   TokenKind = "email_verify"
)

var tokenTTL = map[TokenKind]time.Duration{
	TokenPasswordReset: time.Hour,      // docs/auth/password.md §6.1
	TokenEmailVerify:   24 * time.Hour, // docs/auth/register.md §8.2
}

// maxTokensPerHour caps reset/verify mails per account (password.md §8:
// 3/hour per e-mail). Extra requests still answer 202 but send nothing.
const maxTokensPerHour = 3

// AccountToken is a stored (hashed) single-use token.
type AccountToken struct {
	ID        uuid.UUID
	AccountID uuid.UUID
	Kind      TokenKind
	TokenHash []byte
	ExpiresAt time.Time
	CreatedAt time.Time
}

// OutboundEmail is one message for the e-mail outbox (worker email_dispatch).
type OutboundEmail struct {
	To      string
	Subject string
	Body    string
}

var (
	// ErrRecoveryTokenInvalid: unknown, already used or expired. One error for
	// all three so a token cannot be probed for state.
	ErrRecoveryTokenInvalid = errors.New("auth: recovery token invalid or expired")
	// ErrEmailAlreadyVerified: resend-verification on a verified account.
	ErrEmailAlreadyVerified = errors.New("auth: e-mail already verified")
)

func newRecoveryToken() (plain string, hash []byte, err error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", nil, fmt.Errorf("recovery token: %w", err)
	}
	plain = base64.RawURLEncoding.EncodeToString(b)
	return plain, hashRecoveryToken(plain), nil
}

func hashRecoveryToken(plain string) []byte {
	h := sha256.Sum256([]byte(plain))
	return h[:]
}

// issueToken stores a token and queues its e-mail in one transaction. Silently
// skips (nil error) when the hourly cap is reached — callers answer 202 anyway.
func (s *Service) issueToken(ctx context.Context, a *Account, kind TokenKind, now time.Time) error {
	n, err := s.repo.CountAccountTokens(ctx, a.ID, kind, now.Add(-time.Hour))
	if err != nil {
		return err
	}
	if n >= maxTokensPerHour {
		s.log.InfoContext(ctx, "auth: recovery mail skipped (hourly cap)", "kind", string(kind))
		return nil
	}
	plain, hash, err := newRecoveryToken()
	if err != nil {
		return err
	}
	tok := AccountToken{
		ID: uuid.New(), AccountID: a.ID, Kind: kind, TokenHash: hash,
		ExpiresAt: now.Add(tokenTTL[kind]), CreatedAt: now,
	}
	return s.repo.IssueAccountToken(ctx, tok, recoveryEmail(s.cfg.PublicBaseURL, a, kind, plain))
}

// RequestPasswordReset never reveals whether the e-mail exists: unknown,
// locked and capped accounts all return nil (HTTP 202).
func (s *Service) RequestPasswordReset(ctx context.Context, email string) error {
	a, err := s.repo.AccountByEmailNormalized(ctx, strings.ToLower(strings.TrimSpace(email)))
	if err != nil || a.State == AccountBanned {
		return nil
	}
	return s.issueToken(ctx, a, TokenPasswordReset, s.cfg.Now().UTC())
}

// ConfirmPasswordReset sets a new password and signs out every session.
func (s *Service) ConfirmPasswordReset(ctx context.Context, token, newPassword string) error {
	if err := validatePassword(newPassword, s.cfg.MinPasswordLen); err != nil {
		return err
	}
	hash, err := HashPassword(newPassword, s.cfg.Argon2Params)
	if err != nil {
		return fmt.Errorf("hash password: %w", err)
	}
	now := s.cfg.Now().UTC()
	accountID, err := s.repo.ConsumeAccountToken(ctx, TokenPasswordReset, hashRecoveryToken(token), now)
	if err != nil {
		return err
	}
	if err := s.repo.UpdatePassword(ctx, accountID, hash, now); err != nil {
		return err
	}
	return s.repo.RevokeAccountSessions(ctx, accountID, uuid.Nil, "password_reset", now)
}

// VerifyEmail marks the account's e-mail verified.
func (s *Service) VerifyEmail(ctx context.Context, token string) error {
	now := s.cfg.Now().UTC()
	accountID, err := s.repo.ConsumeAccountToken(ctx, TokenEmailVerify, hashRecoveryToken(token), now)
	if err != nil {
		return err
	}
	return s.repo.MarkEmailVerified(ctx, accountID, now)
}

// ResendVerification queues a fresh verification mail (hourly cap applies).
func (s *Service) ResendVerification(ctx context.Context, accountID uuid.UUID) error {
	a, err := s.repo.AccountByID(ctx, accountID)
	if err != nil {
		return err
	}
	if a.EmailVerifiedAt != nil {
		return ErrEmailAlreadyVerified
	}
	return s.issueToken(ctx, a, TokenEmailVerify, s.cfg.Now().UTC())
}

// ChangePassword verifies the current password, stores the new one and signs
// out every other session (the caller's own session stays valid).
func (s *Service) ChangePassword(ctx context.Context, accountID, currentSession uuid.UUID, current, next string) error {
	a, err := s.repo.AccountByID(ctx, accountID)
	if err != nil {
		return err
	}
	if err := VerifyPassword(current, a.PasswordHash); err != nil {
		return ErrInvalidCredentials
	}
	if err := validatePassword(next, s.cfg.MinPasswordLen); err != nil {
		return err
	}
	hash, err := HashPassword(next, s.cfg.Argon2Params)
	if err != nil {
		return fmt.Errorf("hash password: %w", err)
	}
	now := s.cfg.Now().UTC()
	if err := s.repo.UpdatePassword(ctx, accountID, hash, now); err != nil {
		return err
	}
	return s.repo.RevokeAccountSessions(ctx, accountID, currentSession, "password_changed", now)
}

// recoveryEmail builds the bilingual (TR/EN) message with the one-time link.
func recoveryEmail(baseURL string, a *Account, kind TokenKind, token string) OutboundEmail {
	base := strings.TrimRight(baseURL, "/")
	q := url.Values{"token": {token}}.Encode()
	switch kind {
	case TokenPasswordReset:
		link := base + "/reset-password?" + q
		return OutboundEmail{
			To:      a.Email,
			Subject: "New World MMO — Şifre sıfırlama / Password reset",
			Body: fmt.Sprintf("Merhaba %s,\n\nHesabın için şifre sıfırlama istendi. Yeni şifre belirlemek için:\n%s\n\n"+
				"Bağlantı 1 saat geçerlidir. Bu isteği sen yapmadıysan bu e-postayı yok sayabilirsin.\n\n---\n\n"+
				"Hi %s,\n\nA password reset was requested for your account. Set a new password here:\n%s\n\n"+
				"This link expires in 1 hour. If you didn't request it, ignore this e-mail.\n",
				a.Username, link, a.Username, link),
		}
	default:
		link := base + "/verify-email?" + q
		return OutboundEmail{
			To:      a.Email,
			Subject: "New World MMO — E-posta doğrulama / Verify your e-mail",
			Body: fmt.Sprintf("Merhaba %s,\n\nE-posta adresini doğrulamak için:\n%s\n\nBağlantı 24 saat geçerlidir.\n\n---\n\n"+
				"Hi %s,\n\nConfirm your e-mail address here:\n%s\n\nThis link expires in 24 hours.\n",
				a.Username, link, a.Username, link),
		}
	}
}

// EmailVerified reports whether the account confirmed its e-mail (used by the
// optional character-creation gate, AUTH_REQUIRE_VERIFIED_EMAIL).
func (s *Service) EmailVerified(ctx context.Context, accountID uuid.UUID) (bool, error) {
	a, err := s.repo.AccountByID(ctx, accountID)
	if err != nil {
		return false, err
	}
	return a.EmailVerifiedAt != nil, nil
}
