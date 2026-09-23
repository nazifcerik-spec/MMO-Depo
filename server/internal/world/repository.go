package world

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

// TravelRepository is the persistence contract for the travel domain.
type TravelRepository interface {
	Create(ctx context.Context, t *Travel) error
	GetByID(ctx context.Context, id uuid.UUID) (*Travel, error)
	Active(ctx context.Context, characterID uuid.UUID) (*Travel, error)
	History(ctx context.Context, characterID uuid.UUID, limit int) ([]Travel, error)
	MarkCompleted(ctx context.Context, id uuid.UUID, now time.Time) error
	MarkCancelled(ctx context.Context, id uuid.UUID, now time.Time) error
}

type PgRepo struct{ pool *pgxpool.Pool }

func NewPgRepo(pool *pgxpool.Pool) *PgRepo { return &PgRepo{pool: pool} }

func pgUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "23505"
}

const selectBase = `
	SELECT id, character_id, from_region_code, to_region_code,
	       started_at, arrives_at, state,
	       completed_at, cancelled_at, created_at, updated_at, version
	  FROM world.travel`

func (r *PgRepo) scan(row pgx.Row) (*Travel, error) {
	var t Travel
	var state string
	err := row.Scan(
		&t.ID, &t.CharacterID, &t.FromRegionCode, &t.ToRegionCode,
		&t.StartedAt, &t.ArrivesAt, &state,
		&t.CompletedAt, &t.CancelledAt, &t.CreatedAt, &t.UpdatedAt, &t.Version,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	t.State = TravelState(state)
	return &t, nil
}

func (r *PgRepo) Create(ctx context.Context, t *Travel) error {
	const q = `
		INSERT INTO world.travel (
		  id, character_id, from_region_code, to_region_code,
		  started_at, arrives_at, state, created_at, updated_at
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)`
	_, err := r.pool.Exec(ctx, q,
		t.ID, t.CharacterID, t.FromRegionCode, t.ToRegionCode,
		t.StartedAt, t.ArrivesAt, string(t.State), t.CreatedAt, t.UpdatedAt,
	)
	if err != nil {
		if pgUniqueViolation(err) {
			return ErrTravelInProgress
		}
		return err
	}
	return nil
}

func (r *PgRepo) GetByID(ctx context.Context, id uuid.UUID) (*Travel, error) {
	return r.scan(r.pool.QueryRow(ctx, selectBase+` WHERE id = $1`, id))
}

func (r *PgRepo) Active(ctx context.Context, characterID uuid.UUID) (*Travel, error) {
	return r.scan(r.pool.QueryRow(ctx,
		selectBase+` WHERE character_id = $1 AND state = 'in_progress' LIMIT 1`, characterID))
}

func (r *PgRepo) History(ctx context.Context, characterID uuid.UUID, limit int) ([]Travel, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	rows, err := r.pool.Query(ctx,
		selectBase+` WHERE character_id = $1 ORDER BY started_at DESC LIMIT $2`,
		characterID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Travel
	for rows.Next() {
		t, err := r.scan(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, *t)
	}
	return out, rows.Err()
}

func (r *PgRepo) MarkCompleted(ctx context.Context, id uuid.UUID, now time.Time) error {
	const q = `
		UPDATE world.travel
		   SET state = 'completed', completed_at = $2,
		       updated_at = $2, version = version + 1
		 WHERE id = $1 AND state = 'in_progress'`
	tag, err := r.pool.Exec(ctx, q, id, now)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrTravelNotActive
	}
	return nil
}

func (r *PgRepo) MarkCancelled(ctx context.Context, id uuid.UUID, now time.Time) error {
	const q = `
		UPDATE world.travel
		   SET state = 'cancelled', cancelled_at = $2,
		       updated_at = $2, version = version + 1
		 WHERE id = $1 AND state = 'in_progress'`
	tag, err := r.pool.Exec(ctx, q, id, now)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrTravelNotActive
	}
	return nil
}
