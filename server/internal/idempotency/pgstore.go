package idempotency

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// PgStore reads and writes api.idempotency_record.
//
// Race note: two concurrent identical requests both MISS the initial SELECT,
// both execute the handler, and both attempt INSERT. UNIQUE(key, endpoint,
// account_id) makes the second INSERT a no-op via ON CONFLICT DO NOTHING; the
// first response is cached, the second one is not. Both clients receive their
// own response, so the effect is "handler ran twice." For paths where this is
// unacceptable (ledger writes, item transfer), the domain service takes its own
// SERIALIZABLE lock — idempotency stays a best-effort cache, DD-135 provides the
// authoritative uniqueness.
type PgStore struct {
	pool *pgxpool.Pool
}

func NewPgStore(pool *pgxpool.Pool) *PgStore {
	return &PgStore{pool: pool}
}

func (s *PgStore) Get(ctx context.Context, key, endpoint string, accountID uuid.UUID) (*Record, error) {
	const q = `
		SELECT key, endpoint, account_id, request_hash, response_status, response_body, created_at, expires_at
		  FROM api.idempotency_record
		 WHERE key = $1 AND endpoint = $2 AND account_id = $3 AND expires_at > now()`

	row := s.pool.QueryRow(ctx, q, key, endpoint, accountID)
	var r Record
	if err := row.Scan(
		&r.Key, &r.Endpoint, &r.AccountID, &r.RequestHash,
		&r.ResponseStatus, &r.ResponseBody, &r.CreatedAt, &r.ExpiresAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &r, nil
}

func (s *PgStore) Set(ctx context.Context, rec Record) error {
	const q = `
		INSERT INTO api.idempotency_record
			(key, endpoint, account_id, request_hash, response_status, response_body)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (key, endpoint, account_id) DO NOTHING`

	_, err := s.pool.Exec(ctx, q,
		rec.Key, rec.Endpoint, rec.AccountID, rec.RequestHash,
		rec.ResponseStatus, rec.ResponseBody,
	)
	return err
}
