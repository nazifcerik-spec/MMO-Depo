package dungeon

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Repository is the persistence contract for the dungeon domain.
type Repository interface {
	Create(ctx context.Context, c *Contract) error
	GetByID(ctx context.Context, id uuid.UUID) (*Contract, error)
	Active(ctx context.Context, characterID uuid.UUID) (*Contract, error)
	History(ctx context.Context, characterID uuid.UUID, limit int) ([]Contract, error)

	// MarkRevealed writes the final result. If insertLockout is true, the
	// weekly_lockout row is inserted atomically in the same tx — a duplicate
	// lockout row (someone raced two reveals through) surfaces as
	// ErrWeeklyLockout so the caller can return the cached result.
	MarkRevealed(
		ctx context.Context, id uuid.UUID,
		result json.RawMessage, txnID uuid.UUID, now time.Time,
		insertLockout bool, characterID uuid.UUID, dungeonCode, weekKey string,
	) error

	MarkCancelled(ctx context.Context, id uuid.UUID, now time.Time) error

	// HasLockout returns true when (character, dungeon, week) already exists.
	HasLockout(ctx context.Context, characterID uuid.UUID, dungeonCode, weekKey string) (bool, error)
}

// --- Postgres ---------------------------------------------------------------

type PgRepo struct{ pool *pgxpool.Pool }

func NewPgRepo(pool *pgxpool.Pool) *PgRepo { return &PgRepo{pool: pool} }

func pgUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "23505"
}

const selectBase = `
	SELECT id, character_id, dungeon_code, duration_seconds,
	       started_at, expires_at, seed, state, player_snapshot,
	       revealed_at, reveal_result, reveal_txn_id,
	       cancelled_at, created_at, updated_at, version
	  FROM dungeon.contract`

func (r *PgRepo) scan(row pgx.Row) (*Contract, error) {
	var c Contract
	var state string
	err := row.Scan(
		&c.ID, &c.CharacterID, &c.DungeonCode, &c.DurationSecs,
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
		INSERT INTO dungeon.contract (
		  id, character_id, dungeon_code, duration_seconds,
		  started_at, expires_at, seed, state, player_snapshot,
		  created_at, updated_at
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)`
	_, err := r.pool.Exec(ctx, q,
		c.ID, c.CharacterID, c.DungeonCode, c.DurationSecs,
		c.StartedAt, c.ExpiresAt, c.Seed, string(c.State), c.PlayerSnapshot,
		c.CreatedAt, c.UpdatedAt,
	)
	if err != nil {
		if pgUniqueViolation(err) {
			return ErrSlotBusy
		}
		return err
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

func (r *PgRepo) MarkRevealed(
	ctx context.Context, id uuid.UUID,
	result json.RawMessage, txnID uuid.UUID, now time.Time,
	insertLockout bool, characterID uuid.UUID, dungeonCode, weekKey string,
) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	const upd = `
		UPDATE dungeon.contract
		   SET state = 'revealed', revealed_at = $2, reveal_result = $3,
		       reveal_txn_id = NULLIF($4, '00000000-0000-0000-0000-000000000000'::uuid),
		       updated_at = $2, version = version + 1
		 WHERE id = $1 AND state = 'active'`
	tag, err := tx.Exec(ctx, upd, id, now, result, txnID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotActive
	}
	if insertLockout {
		_, err = tx.Exec(ctx, `
			INSERT INTO dungeon.weekly_lockout
			  (character_id, dungeon_code, week_key, cleared_at, contract_id)
			VALUES ($1,$2,$3,$4,$5)`,
			characterID, dungeonCode, weekKey, now, id)
		if err != nil {
			if pgUniqueViolation(err) {
				return ErrWeeklyLockout
			}
			return err
		}
	}
	return tx.Commit(ctx)
}

func (r *PgRepo) MarkCancelled(ctx context.Context, id uuid.UUID, now time.Time) error {
	const q = `
		UPDATE dungeon.contract
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

func (r *PgRepo) HasLockout(ctx context.Context, characterID uuid.UUID, dungeonCode, weekKey string) (bool, error) {
	var one int
	err := r.pool.QueryRow(ctx,
		`SELECT 1 FROM dungeon.weekly_lockout
		  WHERE character_id = $1 AND dungeon_code = $2 AND week_key = $3`,
		characterID, dungeonCode, weekKey,
	).Scan(&one)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
