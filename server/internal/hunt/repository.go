package hunt

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

// Repository is the persistence contract for the hunt domain.
type Repository interface {
	Create(ctx context.Context, c *Contract) error
	GetByID(ctx context.Context, id uuid.UUID) (*Contract, error)
	Active(ctx context.Context, characterID uuid.UUID) (*Contract, error)
	History(ctx context.Context, characterID uuid.UUID, limit int) ([]Contract, error)
	MarkRevealed(ctx context.Context, id uuid.UUID, result json.RawMessage, txnID uuid.UUID, now time.Time) error
	MarkCancelled(ctx context.Context, id uuid.UUID, now time.Time) error
}

// --- Postgres ---------------------------------------------------------------

type PgRepo struct{ pool *pgxpool.Pool }

func NewPgRepo(pool *pgxpool.Pool) *PgRepo { return &PgRepo{pool: pool} }

func pgUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "23505"
}

const selectBase = `
	SELECT id, character_id, mob_code, encounter_seconds, duration_seconds,
	       started_at, expires_at, seed, state, player_snapshot,
	       revealed_at, reveal_result, reveal_txn_id,
	       cancelled_at, created_at, updated_at, version
	  FROM hunt.contract`

func (r *PgRepo) scan(row pgx.Row) (*Contract, error) {
	var c Contract
	var state string
	err := row.Scan(
		&c.ID, &c.CharacterID, &c.MobCode, &c.EncounterSecs, &c.DurationSecs,
		&c.StartedAt, &c.ExpiresAt, &c.Seed, &state, &c.PlayerSnapshot,
		&c.RevealedAt, &c.RevealResult, &c.RevealTxnID,
		&c.CancelledAt, &c.CreatedAt, &c.UpdatedAt, &c.Version,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	c.State = State(state)
	return &c, nil
}

func (r *PgRepo) Create(ctx context.Context, c *Contract) error {
	const q = `
		INSERT INTO hunt.contract (
		  id, character_id, mob_code, encounter_seconds, duration_seconds,
		  started_at, expires_at, seed, state, player_snapshot,
		  created_at, updated_at
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)`
	_, err := r.pool.Exec(ctx, q,
		c.ID, c.CharacterID, c.MobCode, c.EncounterSecs, c.DurationSecs,
		c.StartedAt, c.ExpiresAt, c.Seed, string(c.State), c.PlayerSnapshot,
		c.CreatedAt, c.UpdatedAt,
	)
	if err != nil {
		if pgUniqueViolation(err) {
			return ErrSlotBusy
		}
		return action.FromDB(err) // busy with another action (DD-263)
	}
	return nil
}

func (r *PgRepo) GetByID(ctx context.Context, id uuid.UUID) (*Contract, error) {
	return r.scan(r.pool.QueryRow(ctx, selectBase+` WHERE id = $1`, id))
}

func (r *PgRepo) Active(ctx context.Context, characterID uuid.UUID) (*Contract, error) {
	return r.scan(r.pool.QueryRow(ctx,
		selectBase+` WHERE character_id = $1 AND state = 'active' LIMIT 1`, characterID))
}

func (r *PgRepo) History(ctx context.Context, characterID uuid.UUID, limit int) ([]Contract, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	rows, err := r.pool.Query(ctx,
		selectBase+` WHERE character_id = $1 ORDER BY created_at DESC LIMIT $2`,
		characterID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Contract
	for rows.Next() {
		c, err := r.scan(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, *c)
	}
	return out, rows.Err()
}

func (r *PgRepo) MarkRevealed(ctx context.Context, id uuid.UUID, result json.RawMessage, txnID uuid.UUID, now time.Time) error {
	const q = `
		UPDATE hunt.contract
		   SET state = 'revealed', revealed_at = $2, reveal_result = $3,
		       reveal_txn_id = NULLIF($4, '00000000-0000-0000-0000-000000000000'::uuid),
		       updated_at = $2, version = version + 1
		 WHERE id = $1 AND state = 'active'`
	tag, err := r.pool.Exec(ctx, q, id, now, result, txnID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotActive
	}
	return nil
}

func (r *PgRepo) MarkCancelled(ctx context.Context, id uuid.UUID, now time.Time) error {
	const q = `
		UPDATE hunt.contract
		   SET state = 'cancelled', cancelled_at = $2, updated_at = $2, version = version + 1
		 WHERE id = $1 AND state = 'active'`
	tag, err := r.pool.Exec(ctx, q, id, now)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotActive
	}
	return nil
}
