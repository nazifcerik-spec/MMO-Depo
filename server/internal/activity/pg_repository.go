package activity

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/otter/server/internal/action"
)

// PgRepo persists activity.contract rows.
type PgRepo struct {
	pool *pgxpool.Pool
}

func NewPgRepo(pool *pgxpool.Pool) *PgRepo {
	return &PgRepo{pool: pool}
}

func pgUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "23505"
}

func (r *PgRepo) Create(ctx context.Context, c Contract) error {
	const q = `
		INSERT INTO activity.contract
		  (id, character_id, kind, duration_seconds, started_at, expires_at,
		   seed, state, snapshot, params, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)`
	_, err := r.pool.Exec(ctx, q,
		c.ID, c.CharacterID, string(c.Kind), c.DurationSeconds,
		c.StartedAt, c.ExpiresAt, c.Seed, string(c.State),
		[]byte(c.Snapshot), []byte(c.Params), c.CreatedAt, c.UpdatedAt,
	)
	if err != nil {
		// UNIQUE partial index (character_id WHERE state='active') → slot busy.
		if pgUniqueViolation(err) {
			return ErrSlotBusy
		}
		return action.FromDB(err) // busy with another action (DD-263)
	}
	return nil
}

const baseSelect = `
	SELECT id, character_id, kind, duration_seconds, started_at, expires_at,
	       seed, state, snapshot, params,
	       revealed_at, reveal_result, reveal_txn_id,
	       cancelled_at, created_at, updated_at, version
	  FROM activity.contract`

func (r *PgRepo) scan(row pgx.Row) (*Contract, error) {
	var c Contract
	var (
		kind, state string
		snapshot    []byte
		params      []byte
		reveal      []byte
	)
	err := row.Scan(
		&c.ID, &c.CharacterID, &kind, &c.DurationSeconds, &c.StartedAt, &c.ExpiresAt,
		&c.Seed, &state, &snapshot, &params,
		&c.RevealedAt, &reveal, &c.RevealTxnID,
		&c.CancelledAt, &c.CreatedAt, &c.UpdatedAt, &c.Version,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	c.Kind = Kind(kind)
	c.State = State(state)
	c.Snapshot = json.RawMessage(snapshot)
	c.Params = json.RawMessage(params)
	if reveal != nil {
		c.RevealResult = json.RawMessage(reveal)
	}
	return &c, nil
}

func (r *PgRepo) GetByID(ctx context.Context, id uuid.UUID) (*Contract, error) {
	return r.scan(r.pool.QueryRow(ctx, baseSelect+" WHERE id = $1", id))
}

func (r *PgRepo) ActiveByCharacter(ctx context.Context, characterID uuid.UUID) (*Contract, error) {
	return r.scan(r.pool.QueryRow(ctx,
		baseSelect+" WHERE character_id = $1 AND state = 'active'", characterID))
}

func (r *PgRepo) MarkRevealed(ctx context.Context, id uuid.UUID, result json.RawMessage, txnID uuid.UUID, at time.Time) error {
	const q = `UPDATE activity.contract
	              SET state = 'revealed',
	                  reveal_result = $2,
	                  reveal_txn_id = $3,
	                  revealed_at = $4,
	                  updated_at = $4,
	                  version = version + 1
	            WHERE id = $1 AND state = 'active'`
	tag, err := r.pool.Exec(ctx, q, id, []byte(result), txnID, at)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotActive
	}
	return nil
}

func (r *PgRepo) MarkCancelled(ctx context.Context, id uuid.UUID, at time.Time) error {
	const q = `UPDATE activity.contract
	              SET state = 'cancelled',
	                  cancelled_at = $2,
	                  updated_at = $2,
	                  version = version + 1
	            WHERE id = $1 AND state = 'active'`
	tag, err := r.pool.Exec(ctx, q, id, at)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotActive
	}
	return nil
}

func (r *PgRepo) ListByCharacter(ctx context.Context, characterID uuid.UUID, limit int) ([]Contract, error) {
	if limit <= 0 {
		limit = 20
	}
	rows, err := r.pool.Query(ctx,
		baseSelect+` WHERE character_id = $1 ORDER BY started_at DESC LIMIT $2`,
		characterID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []Contract{}
	for rows.Next() {
		c, err := r.scan(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, *c)
	}
	return out, rows.Err()
}
