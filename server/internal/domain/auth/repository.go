package auth

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// Repository is the persistence contract for the auth domain. One aggregated
// interface (rather than three fine-grained) keeps the service wiring simple;
// consumers can still adapt / spy specific methods.
type Repository interface {
	// --- Account
	CreateAccount(ctx context.Context, a Account) error
	AccountByEmailNormalized(ctx context.Context, emailN string) (*Account, error)
	AccountByUsernameNormalized(ctx context.Context, userN string) (*Account, error)
	AccountByID(ctx context.Context, id uuid.UUID) (*Account, error)
	UpdateAccountLastLogin(ctx context.Context, id uuid.UUID, ts time.Time) error

	// --- Session
	CreateSession(ctx context.Context, s Session) error
	Session(ctx context.Context, id uuid.UUID) (*Session, error)
	RevokeSession(ctx context.Context, id uuid.UUID, reason string, at time.Time) error
	TouchSession(ctx context.Context, id uuid.UUID, at time.Time) error

	// --- Refresh token
	CreateRefreshToken(ctx context.Context, t RefreshToken) error
	RefreshTokenByHash(ctx context.Context, hash []byte) (*RefreshToken, error)
	MarkRefreshTokenUsed(ctx context.Context, id uuid.UUID, at time.Time) error
	RevokeRefreshTokensForSession(ctx context.Context, sessionID uuid.UUID, reason string, at time.Time) error

	// --- Account recovery (DD-219)
	CountAccountTokens(ctx context.Context, accountID uuid.UUID, kind TokenKind, since time.Time) (int, error)
	// IssueAccountToken stores the hashed token and queues its e-mail atomically.
	IssueAccountToken(ctx context.Context, t AccountToken, mail OutboundEmail) error
	// ConsumeAccountToken marks a live token used (and burns the account's other
	// open tokens of that kind); ErrRecoveryTokenInvalid when none matches.
	ConsumeAccountToken(ctx context.Context, kind TokenKind, hash []byte, now time.Time) (uuid.UUID, error)
	UpdatePassword(ctx context.Context, accountID uuid.UUID, hash string, at time.Time) error
	MarkEmailVerified(ctx context.Context, accountID uuid.UUID, at time.Time) error
	// RevokeAccountSessions revokes every live session (+ refresh tokens) of the
	// account except `keep` (uuid.Nil = revoke all).
	RevokeAccountSessions(ctx context.Context, accountID, keep uuid.UUID, reason string, at time.Time) error
}
