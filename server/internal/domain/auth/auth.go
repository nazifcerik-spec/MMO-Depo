// Package auth implements user registration, login, session, refresh-rotation
// (DD-189), and Bearer verification.
//
// Layering:
//
//	http.go        HTTP handlers (POST /v1/auth/*).
//	service.go     Business logic (Register/Login/Refresh/Logout/Verify).
//	repository.go  Persistence contract + MemoryRepo (tests).
//	pg_repo.go     Postgres implementation.
//	password.go    Argon2id encode/verify.
//	jwt.go         Ed25519 signer.
package auth

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// AccountState is the machine-readable lifecycle marker.
type AccountState string

const (
	AccountActive AccountState = "active"
	AccountLocked AccountState = "locked" // temporary — auto expire or admin unlock
	AccountBanned AccountState = "banned" // permanent per admin action
)

// Account is the persisted user record. PasswordHash is the full Argon2id
// encoded string (algorithm + params + salt + hash) — never plaintext.
type Account struct {
	ID                 uuid.UUID
	Email              string
	EmailNormalized    string
	Username           string
	UsernameNormalized string
	PasswordHash       string
	State              AccountState
	EmailVerifiedAt    *time.Time
	LastLoginAt        *time.Time
	CreatedAt          time.Time
	UpdatedAt          time.Time
	Version            int32
}

// Session represents one logged-in device instance. Refresh tokens rotate
// against the session; if any rotation detects reuse, the whole session and
// its token chain are revoked (DD-189).
type Session struct {
	ID           uuid.UUID
	AccountID    uuid.UUID
	DeviceID     *uuid.UUID
	IPHash       []byte
	UserAgent    string
	CreatedAt    time.Time
	LastActiveAt time.Time
	ExpiresAt    time.Time
	RevokedAt    *time.Time
	RevokeReason string
}

// RefreshToken is a single link in the rotation chain. token_hash is
// SHA-256(plaintext) — the plaintext is only returned once to the client at
// issuance. ParentTokenID lets reuse detection walk the chain forward if
// needed.
type RefreshToken struct {
	ID            uuid.UUID
	SessionID     uuid.UUID
	TokenHash     []byte
	ParentTokenID *uuid.UUID
	IssuedAt      time.Time
	ExpiresAt     time.Time
	UsedAt        *time.Time
	RevokedAt     *time.Time
	RevokeReason  string
}

// LoginResult is the response every Register/Login/Refresh call returns.
// AccessToken and RefreshToken are plaintext — the client stores them and the
// server retains only the hash.
type LoginResult struct {
	AccountID       uuid.UUID
	SessionID       uuid.UUID
	AccessToken     string
	AccessExpiresAt time.Time
	RefreshToken    string
	RefreshExpires  time.Time
	EmailVerified   bool // DD-219: client shows a verify banner when false
}

// --- Errors ------------------------------------------------------------------

var (
	// ErrAccountExists is returned by Register when email or username is taken.
	ErrAccountExists = errors.New("auth: account already exists")

	// ErrInvalidCredentials is returned for login attempts that fail auth. The
	// same error covers "unknown email" and "wrong password" so the caller
	// cannot enumerate accounts.
	ErrInvalidCredentials = errors.New("auth: invalid credentials")

	// ErrAccountLocked / ErrAccountBanned reject login for a non-active state.
	ErrAccountLocked = errors.New("auth: account locked")
	ErrAccountBanned = errors.New("auth: account banned")

	// ErrSessionRevoked means the session was revoked — logout, admin action,
	// or refresh reuse cascade.
	ErrSessionRevoked = errors.New("auth: session revoked")

	// ErrSessionExpired means expires_at is in the past.
	ErrSessionExpired = errors.New("auth: session expired")

	// ErrRefreshTokenNotFound is returned when the refresh token does not match
	// any known hash (unknown or already garbage-collected).
	ErrRefreshTokenNotFound = errors.New("auth: refresh token not found")

	// ErrRefreshTokenReuse is the compromise signal — a token already used_at
	// has been presented again. Session and its chain are revoked (DD-189).
	ErrRefreshTokenReuse = errors.New("auth: refresh token reuse detected")

	// ErrTokenExpired covers access + refresh expiration.
	ErrTokenExpired = errors.New("auth: token expired")

	// ErrTokenInvalid covers signature mismatch, malformed JWT, unknown key.
	ErrTokenInvalid = errors.New("auth: token invalid")

	// ErrValidation is a request-shape rejection (email format, weak password,
	// empty username). HTTP layer maps to 400 with the field name.
	ErrValidation = errors.New("auth: validation")
)

// ValidationError enriches ErrValidation with the offending field.
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string { return "auth: validation: " + e.Field + " " + e.Message }
func (e *ValidationError) Unwrap() error { return ErrValidation }
