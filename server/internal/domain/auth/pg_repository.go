package auth

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

// PgRepo persists auth aggregates in Postgres. Every method takes ctx and
// respects statement timeouts set on the pool.
type PgRepo struct {
	pool *pgxpool.Pool
}

func NewPgRepo(pool *pgxpool.Pool) *PgRepo {
	return &PgRepo{pool: pool}
}

// pgUniqueViolation returns true when err is a Postgres 23505 unique constraint
// violation — matched against email/username uniqueness on CreateAccount.
func pgUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "23505"
}

// --- Account -----------------------------------------------------------------

func (r *PgRepo) CreateAccount(ctx context.Context, a Account) error {
	const q = `
		INSERT INTO auth.account
		  (id, email, email_normalized, username, username_normalized, password_hash,
		   state, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)`
	_, err := r.pool.Exec(ctx, q,
		a.ID, a.Email, a.EmailNormalized, a.Username, a.UsernameNormalized,
		a.PasswordHash, string(a.State), a.CreatedAt, a.UpdatedAt,
	)
	if err != nil {
		if pgUniqueViolation(err) {
			return ErrAccountExists
		}
		return err
	}
	return nil
}

func (r *PgRepo) accountBase() string {
	return `SELECT id, email, email_normalized, username, username_normalized, password_hash,
	               state, email_verified_at, last_login_at, created_at, updated_at, version
	          FROM auth.account`
}

func (r *PgRepo) scanAccount(row pgx.Row) (*Account, error) {
	var a Account
	var state string
	err := row.Scan(&a.ID, &a.Email, &a.EmailNormalized, &a.Username, &a.UsernameNormalized,
		&a.PasswordHash, &state, &a.EmailVerifiedAt, &a.LastLoginAt,
		&a.CreatedAt, &a.UpdatedAt, &a.Version)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}
	a.State = AccountState(state)
	return &a, nil
}

func (r *PgRepo) AccountByEmailNormalized(ctx context.Context, emailN string) (*Account, error) {
	return r.scanAccount(r.pool.QueryRow(ctx, r.accountBase()+" WHERE email_normalized = $1", emailN))
}

func (r *PgRepo) AccountByUsernameNormalized(ctx context.Context, userN string) (*Account, error) {
	return r.scanAccount(r.pool.QueryRow(ctx, r.accountBase()+" WHERE username_normalized = $1", userN))
}

func (r *PgRepo) AccountByID(ctx context.Context, id uuid.UUID) (*Account, error) {
	return r.scanAccount(r.pool.QueryRow(ctx, r.accountBase()+" WHERE id = $1", id))
}

func (r *PgRepo) UpdateAccountLastLogin(ctx context.Context, id uuid.UUID, ts time.Time) error {
	const q = `UPDATE auth.account SET last_login_at = $2, updated_at = $2 WHERE id = $1`
	_, err := r.pool.Exec(ctx, q, id, ts)
	return err
}

// --- Session -----------------------------------------------------------------

func (r *PgRepo) CreateSession(ctx context.Context, s Session) error {
	const q = `
		INSERT INTO auth.session
		  (id, account_id, device_id, ip_hash, user_agent, created_at, last_active_at, expires_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`
	_, err := r.pool.Exec(ctx, q,
		s.ID, s.AccountID, s.DeviceID, s.IPHash, s.UserAgent,
		s.CreatedAt, s.LastActiveAt, s.ExpiresAt,
	)
	return err
}

func (r *PgRepo) Session(ctx context.Context, id uuid.UUID) (*Session, error) {
	const q = `SELECT id, account_id, device_id, ip_hash, user_agent,
	                  created_at, last_active_at, expires_at, revoked_at, revoke_reason
	             FROM auth.session WHERE id = $1`
	var s Session
	var reason *string
	err := r.pool.QueryRow(ctx, q, id).Scan(
		&s.ID, &s.AccountID, &s.DeviceID, &s.IPHash, &s.UserAgent,
		&s.CreatedAt, &s.LastActiveAt, &s.ExpiresAt, &s.RevokedAt, &reason,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrSessionRevoked
		}
		return nil, err
	}
	if reason != nil {
		s.RevokeReason = *reason
	}
	return &s, nil
}

func (r *PgRepo) RevokeSession(ctx context.Context, id uuid.UUID, reason string, at time.Time) error {
	const q = `UPDATE auth.session
	              SET revoked_at = $2, revoke_reason = $3
	            WHERE id = $1 AND revoked_at IS NULL`
	_, err := r.pool.Exec(ctx, q, id, at, reason)
	return err
}

func (r *PgRepo) TouchSession(ctx context.Context, id uuid.UUID, at time.Time) error {
	const q = `UPDATE auth.session SET last_active_at = $2 WHERE id = $1 AND revoked_at IS NULL`
	_, err := r.pool.Exec(ctx, q, id, at)
	return err
}

// --- Refresh token -----------------------------------------------------------

func (r *PgRepo) CreateRefreshToken(ctx context.Context, t RefreshToken) error {
	const q = `
		INSERT INTO auth.refresh_token
		  (id, session_id, token_hash, parent_token_id, issued_at, expires_at)
		VALUES ($1,$2,$3,$4,$5,$6)`
	_, err := r.pool.Exec(ctx, q,
		t.ID, t.SessionID, t.TokenHash, t.ParentTokenID, t.IssuedAt, t.ExpiresAt,
	)
	return err
}

func (r *PgRepo) RefreshTokenByHash(ctx context.Context, hash []byte) (*RefreshToken, error) {
	const q = `SELECT id, session_id, token_hash, parent_token_id,
	                  issued_at, expires_at, used_at, revoked_at, revoke_reason
	             FROM auth.refresh_token
	            WHERE token_hash = $1`
	var t RefreshToken
	var reason *string
	err := r.pool.QueryRow(ctx, q, hash).Scan(
		&t.ID, &t.SessionID, &t.TokenHash, &t.ParentTokenID,
		&t.IssuedAt, &t.ExpiresAt, &t.UsedAt, &t.RevokedAt, &reason,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrRefreshTokenNotFound
		}
		return nil, err
	}
	if reason != nil {
		t.RevokeReason = *reason
	}
	return &t, nil
}

func (r *PgRepo) MarkRefreshTokenUsed(ctx context.Context, id uuid.UUID, at time.Time) error {
	const q = `UPDATE auth.refresh_token SET used_at = $2 WHERE id = $1 AND used_at IS NULL`
	_, err := r.pool.Exec(ctx, q, id, at)
	return err
}

func (r *PgRepo) RevokeRefreshTokensForSession(ctx context.Context, sessionID uuid.UUID, reason string, at time.Time) error {
	const q = `UPDATE auth.refresh_token
	              SET revoked_at = $3, revoke_reason = $2
	            WHERE session_id = $1 AND revoked_at IS NULL`
	_, err := r.pool.Exec(ctx, q, sessionID, reason, at)
	return err
}

// --- Account recovery (DD-219) -----------------------------------------------

func (r *PgRepo) CountAccountTokens(ctx context.Context, accountID uuid.UUID, kind TokenKind, since time.Time) (int, error) {
	var n int
	err := r.pool.QueryRow(ctx, `
		SELECT count(*) FROM auth.account_token
		 WHERE account_id = $1 AND kind = $2 AND created_at >= $3`, accountID, string(kind), since).Scan(&n)
	return n, err
}

func (r *PgRepo) IssueAccountToken(ctx context.Context, t AccountToken, mail OutboundEmail) error {
	return pgx.BeginFunc(ctx, r.pool, func(tx pgx.Tx) error {
		if _, err := tx.Exec(ctx, `
			INSERT INTO auth.account_token (id, account_id, kind, token_hash, expires_at, created_at)
			VALUES ($1, $2, $3, $4, $5, $6)`,
			t.ID, t.AccountID, string(t.Kind), t.TokenHash, t.ExpiresAt, t.CreatedAt); err != nil {
			return err
		}
		_, err := tx.Exec(ctx, `INSERT INTO ops.email_outbox (to_address, subject, body) VALUES ($1, $2, $3)`,
			mail.To, mail.Subject, mail.Body)
		return err
	})
}

func (r *PgRepo) ConsumeAccountToken(ctx context.Context, kind TokenKind, hash []byte, now time.Time) (uuid.UUID, error) {
	var accountID uuid.UUID
	err := pgx.BeginFunc(ctx, r.pool, func(tx pgx.Tx) error {
		err := tx.QueryRow(ctx, `
			UPDATE auth.account_token SET used_at = $3
			 WHERE token_hash = $1 AND kind = $2 AND used_at IS NULL AND expires_at > $3
			RETURNING account_id`, hash, string(kind), now).Scan(&accountID)
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrRecoveryTokenInvalid
		}
		if err != nil {
			return err
		}
		// Older links of the same kind stop working too.
		_, err = tx.Exec(ctx, `
			UPDATE auth.account_token SET used_at = $3
			 WHERE account_id = $1 AND kind = $2 AND used_at IS NULL`, accountID, string(kind), now)
		return err
	})
	return accountID, err
}

func (r *PgRepo) UpdatePassword(ctx context.Context, accountID uuid.UUID, hash string, at time.Time) error {
	_, err := r.pool.Exec(ctx, `
		UPDATE auth.account
		   SET password_hash = $2, password_updated_at = $3, updated_at = $3, version = version + 1
		 WHERE id = $1`, accountID, hash, at)
	return err
}

func (r *PgRepo) MarkEmailVerified(ctx context.Context, accountID uuid.UUID, at time.Time) error {
	_, err := r.pool.Exec(ctx, `
		UPDATE auth.account
		   SET email_verified_at = COALESCE(email_verified_at, $2), updated_at = $2, version = version + 1
		 WHERE id = $1`, accountID, at)
	return err
}

func (r *PgRepo) RevokeAccountSessions(ctx context.Context, accountID, keep uuid.UUID, reason string, at time.Time) error {
	return pgx.BeginFunc(ctx, r.pool, func(tx pgx.Tx) error {
		if _, err := tx.Exec(ctx, `
			UPDATE auth.refresh_token rt SET revoked_at = $3, revoke_reason = $4
			  FROM auth.session s
			 WHERE rt.session_id = s.id AND s.account_id = $1 AND s.id <> $2
			   AND rt.revoked_at IS NULL`, accountID, keep, at, reason); err != nil {
			return err
		}
		_, err := tx.Exec(ctx, `
			UPDATE auth.session SET revoked_at = $3, revoke_reason = $4
			 WHERE account_id = $1 AND id <> $2 AND revoked_at IS NULL`, accountID, keep, at, reason)
		return err
	})
}
