package character

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

// PgRepo persists character aggregates in Postgres.
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

func (r *PgRepo) Create(ctx context.Context, c Character) error {
	const q = `
		INSERT INTO character.character
		  (id, account_id, name, name_normalized, race_code, class_code, spec_code,
		   level, xp, state, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)`
	_, err := r.pool.Exec(ctx, q,
		c.ID, c.AccountID, c.Name, c.NameNormalized, c.RaceCode, c.ClassCode, c.SpecCode,
		c.Level, c.XP, string(c.State), c.CreatedAt, c.UpdatedAt,
	)
	if err != nil {
		if pgUniqueViolation(err) {
			return ErrNameTaken
		}
		return err
	}
	return nil
}

func (r *PgRepo) scan(row pgx.Row) (*Character, error) {
	var c Character
	var state string
	var locRegion *string
	err := row.Scan(
		&c.ID, &c.AccountID, &c.Name, &c.NameNormalized, &c.RaceCode, &c.ClassCode,
		&c.SpecCode, &c.Level, &c.XP, &state, &c.LocationID, &locRegion,
		&c.CreatedAt, &c.UpdatedAt, &c.LastPlayedAt, &c.DeleteGraceAt, &c.Version,
		&c.HPStored, &c.HPUpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	c.State = State(state)
	if locRegion != nil {
		c.LocationRegionCode = *locRegion
	}
	return &c, nil
}

const baseSelect = `SELECT id, account_id, name, name_normalized, race_code, class_code,
	                       spec_code, level, xp, state, location_id, location_region_code,
	                       created_at, updated_at, last_played_at, delete_grace_at, version,
	                       hp_current, hp_updated_at
	                  FROM character.character`

func (r *PgRepo) Get(ctx context.Context, id uuid.UUID) (*Character, error) {
	return r.scan(r.pool.QueryRow(ctx, baseSelect+" WHERE id = $1", id))
}

func (r *PgRepo) ListByAccount(ctx context.Context, accountID uuid.UUID) ([]Character, error) {
	rows, err := r.pool.Query(ctx,
		baseSelect+" WHERE account_id = $1 AND state IN ('active','delete_grace') ORDER BY created_at ASC",
		accountID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := []Character{}
	for rows.Next() {
		c, err := r.scan(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, *c)
	}
	return out, rows.Err()
}

func (r *PgRepo) CountActiveByAccount(ctx context.Context, accountID uuid.UUID) (int, error) {
	var n int
	err := r.pool.QueryRow(ctx,
		`SELECT count(*) FROM character.character
		  WHERE account_id = $1 AND state IN ('active','delete_grace')`,
		accountID).Scan(&n)
	return n, err
}

func (r *PgRepo) Rename(ctx context.Context, id uuid.UUID, name, normalized string) error {
	const q = `UPDATE character.character
	              SET name = $2, name_normalized = $3,
	                  updated_at = now(), version = version + 1
	            WHERE id = $1`
	tag, err := r.pool.Exec(ctx, q, id, name, normalized)
	if err != nil {
		if pgUniqueViolation(err) {
			return ErrNameTaken
		}
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *PgRepo) TouchLastPlayed(ctx context.Context, id uuid.UUID, ts time.Time) error {
	const q = `UPDATE character.character
	              SET last_played_at = $2
	            WHERE id = $1 AND state = 'active'`
	_, err := r.pool.Exec(ctx, q, id, ts)
	return err
}

func (r *PgRepo) SetDeleteGrace(ctx context.Context, id uuid.UUID, expiresAt time.Time) error {
	const q = `UPDATE character.character
	              SET state = 'delete_grace', delete_grace_at = $2,
	                  updated_at = now(), version = version + 1
	            WHERE id = $1 AND state = 'active'`
	tag, err := r.pool.Exec(ctx, q, id, expiresAt)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		// Either not found or not in active state.
		exists, err := r.Get(ctx, id)
		if err != nil {
			return err
		}
		if exists.State != StateActive {
			return ErrNotPlayable
		}
		return ErrNotFound
	}
	return nil
}

// AddXP is transactional: locks the row, computes rollover into the new
// bracket in application code (Postgres has no native "compute next threshold"),
// then writes level+xp+version in one UPDATE. Repeated rollovers happen inside
// the loop so a huge XP grant can leap several levels at once.
func (r *PgRepo) AddXP(ctx context.Context, id uuid.UUID, delta int64) (*Character, bool, error) {
	if delta < 0 {
		return nil, false, ErrInvalidXP
	}
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted})
	if err != nil {
		return nil, false, err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	var (
		level int16
		xp    int64
	)
	err = tx.QueryRow(ctx, `SELECT level, xp FROM character.character WHERE id=$1 FOR UPDATE`, id).
		Scan(&level, &xp)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, false, ErrNotFound
		}
		return nil, false, err
	}

	leveled := false
	xp += delta
	for {
		needed := XPToNext(int(level))
		if xp < needed || int(level) >= LevelCap {
			break
		}
		xp -= needed
		level++
		leveled = true
	}
	// At cap the excess XP is intentionally kept — the client can still show
	// "full bar"; Phase 02 prestige/paragon consumes the carry.

	_, err = tx.Exec(ctx, `
		UPDATE character.character
		   SET level = $2, xp = $3, updated_at = now(), version = version + 1
		 WHERE id = $1`, id, level, xp)
	if err != nil {
		return nil, false, err
	}
	if err := tx.Commit(ctx); err != nil {
		return nil, false, err
	}
	// Re-fetch to return a full aggregate (state, timestamps, version).
	c, gerr := r.Get(ctx, id)
	if gerr != nil {
		return nil, false, gerr
	}
	return c, leveled, nil
}

// SetLocation writes location_region_code. Passing empty string clears the
// column (NULL) — reserved for admin/reset flows; travel.Complete never
// clears, only assigns the destination.
func (r *PgRepo) SetLocation(ctx context.Context, id uuid.UUID, regionCode string) error {
	var val any
	if regionCode == "" {
		val = nil
	} else {
		val = regionCode
	}
	const q = `UPDATE character.character
	              SET location_region_code = $2,
	                  updated_at = now(), version = version + 1
	            WHERE id = $1`
	tag, err := r.pool.Exec(ctx, q, id, val)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

// GetLocation returns the current region_code or empty string when NULL.
func (r *PgRepo) GetLocation(ctx context.Context, id uuid.UUID) (string, error) {
	var loc *string
	err := r.pool.QueryRow(ctx,
		`SELECT location_region_code FROM character.character WHERE id = $1`, id).
		Scan(&loc)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", ErrNotFound
		}
		return "", err
	}
	if loc == nil {
		return "", nil
	}
	return *loc, nil
}

func (r *PgRepo) CancelDeleteGrace(ctx context.Context, id uuid.UUID) error {
	const q = `UPDATE character.character
	              SET state = 'active', delete_grace_at = NULL,
	                  updated_at = now(), version = version + 1
	            WHERE id = $1 AND state = 'delete_grace'`
	tag, err := r.pool.Exec(ctx, q, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotInDeleteGrace
	}
	return nil
}

func (r *PgRepo) SetHP(ctx context.Context, id uuid.UUID, hp int32, at time.Time) error {
	tag, err := r.pool.Exec(ctx, `
		UPDATE character.character
		   SET hp_current = $2, hp_updated_at = $3, updated_at = now(), version = version + 1
		 WHERE id = $1`, id, max(hp, 0), at)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}
