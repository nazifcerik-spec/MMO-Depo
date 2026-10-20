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
		&c.Alloc.Strength, &c.Alloc.Vitality, &c.Alloc.Dexterity, &c.Alloc.Intelligence, &c.Alloc.Wisdom, &c.Alloc.Luck,
		&c.OverflowXP, &c.PostCapXP,
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
	                       hp_current, hp_updated_at,
	                       attr_str, attr_vit, attr_dex, attr_int, attr_wis, attr_luck,
	                       promotion_overflow_xp, postcap_xp
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

func (r *PgRepo) Purge(ctx context.Context, id uuid.UUID) error {
	return pgx.BeginFunc(ctx, r.pool, func(tx pgx.Tx) error {
		var locked uuid.UUID
		err := tx.QueryRow(ctx, `SELECT id FROM character.character WHERE id = $1 FOR UPDATE`, id).Scan(&locked)
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrNotFound
		}
		if err != nil {
			return err
		}
		_, err = PurgeTx(ctx, tx, []uuid.UUID{id})
		return err
	})
}

// PurgeTx hard-deletes characters inside tx (rows must be locked by the
// caller). Guilds they lead pass to the most senior, earliest remaining
// member; guilds with no one else left are removed. Items go first; every
// other character-owned table cascades. Ledger rows stay (no FK) for audit.
func PurgeTx(ctx context.Context, tx pgx.Tx, ids []uuid.UUID) (int64, error) {
	if _, err := tx.Exec(ctx, `
		UPDATE social.guild g SET leader_character_id = s.character_id, version = g.version + 1
		  FROM (SELECT DISTINCT ON (m.guild_id) m.guild_id, m.character_id
		          FROM social.guild_member m
		          JOIN social.guild g2 ON g2.id = m.guild_id
		         WHERE g2.leader_character_id = ANY($1) AND NOT (m.character_id = ANY($1))
		         ORDER BY m.guild_id, m.rank_index, m.joined_at) s
		 WHERE g.id = s.guild_id`, ids); err != nil {
		return 0, err
	}
	if _, err := tx.Exec(ctx, `
		UPDATE social.guild_member m SET rank_index = 0
		  FROM social.guild g
		 WHERE g.id = m.guild_id AND g.leader_character_id = m.character_id AND m.rank_index <> 0`); err != nil {
		return 0, err
	}
	if _, err := tx.Exec(ctx, `DELETE FROM social.guild WHERE leader_character_id = ANY($1)`, ids); err != nil {
		return 0, err
	}
	if _, err := tx.Exec(ctx, `DELETE FROM inventory.item_instance WHERE owner_character_id = ANY($1)`, ids); err != nil {
		return 0, err
	}
	tag, err := tx.Exec(ctx, `DELETE FROM character.character WHERE id = ANY($1)`, ids)
	return tag.RowsAffected(), err
}

// AddXP is transactional: locks the row, computes rollover into the new
// bracket in application code (Postgres has no native "compute next threshold"),
// then writes level+xp+version in one UPDATE. Repeated rollovers happen inside
// the loop so a huge XP grant can leap several levels at once.
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
