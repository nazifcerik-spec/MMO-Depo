package crafting

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Repository is the persistence contract for character professions + audit log.
type Repository interface {
	// UpsertLearn ensures a (character, profession) row exists at level 1, xp 0.
	// Idempotent — re-learning is a no-op that returns the current row.
	UpsertLearn(ctx context.Context, characterID uuid.UUID, professionCode string, now time.Time) (*CharacterProfession, error)

	// Get returns the character's row for a specific profession, or ErrNotFound.
	Get(ctx context.Context, characterID uuid.UUID, professionCode string) (*CharacterProfession, error)

	// List returns every profession the character has learned.
	List(ctx context.Context, characterID uuid.UUID) ([]CharacterProfession, error)

	// AddXP atomically adds delta to xp, rolls over into levels using the
	// engine curve (ProfessionXPToNext) up to maxLevel, and returns the row
	// AFTER the update + leveledUp flag. Caller passes maxLevel from catalog.
	AddXP(ctx context.Context, characterID uuid.UUID, professionCode string, delta int32, maxLevel int16) (*CharacterProfession, bool, error)

	// InsertCraftLog appends an audit row. Non-blocking on failure would be
	// acceptable but v1 returns the error so the client sees any DB issue.
	InsertCraftLog(ctx context.Context, characterID uuid.UUID, recipeCode string, outputID *uuid.UUID, goldTxnID *uuid.UUID, now time.Time) error
}

type PgRepo struct{ pool *pgxpool.Pool }

func NewPgRepo(pool *pgxpool.Pool) *PgRepo { return &PgRepo{pool: pool} }

func (r *PgRepo) UpsertLearn(ctx context.Context, characterID uuid.UUID, professionCode string, now time.Time) (*CharacterProfession, error) {
	const q = `
		INSERT INTO crafting.character_profession
		  (character_id, profession_code, level, xp, learned_at, updated_at)
		VALUES ($1, $2, 1, 0, $3, $3)
		ON CONFLICT (character_id, profession_code) DO NOTHING`
	if _, err := r.pool.Exec(ctx, q, characterID, professionCode, now); err != nil {
		return nil, err
	}
	return r.Get(ctx, characterID, professionCode)
}

func (r *PgRepo) Get(ctx context.Context, characterID uuid.UUID, professionCode string) (*CharacterProfession, error) {
	var cp CharacterProfession
	err := r.pool.QueryRow(ctx, `
		SELECT character_id, profession_code, level, xp, learned_at, updated_at, version
		  FROM crafting.character_profession
		 WHERE character_id = $1 AND profession_code = $2`, characterID, professionCode).
		Scan(&cp.CharacterID, &cp.ProfessionCode, &cp.Level, &cp.XP, &cp.LearnedAt, &cp.UpdatedAt, &cp.Version)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &cp, nil
}

func (r *PgRepo) List(ctx context.Context, characterID uuid.UUID) ([]CharacterProfession, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT character_id, profession_code, level, xp, learned_at, updated_at, version
		  FROM crafting.character_profession
		 WHERE character_id = $1
		 ORDER BY profession_code`, characterID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []CharacterProfession
	for rows.Next() {
		var cp CharacterProfession
		if err := rows.Scan(&cp.CharacterID, &cp.ProfessionCode, &cp.Level, &cp.XP, &cp.LearnedAt, &cp.UpdatedAt, &cp.Version); err != nil {
			return nil, err
		}
		out = append(out, cp)
	}
	return out, rows.Err()
}

func (r *PgRepo) AddXP(ctx context.Context, characterID uuid.UUID, professionCode string, delta int32, maxLevel int16) (*CharacterProfession, bool, error) {
	if delta < 0 {
		return nil, false, ErrValidation
	}
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted})
	if err != nil {
		return nil, false, err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	var (
		level int16
		xp    int32
	)
	err = tx.QueryRow(ctx,
		`SELECT level, xp FROM crafting.character_profession
		  WHERE character_id = $1 AND profession_code = $2 FOR UPDATE`,
		characterID, professionCode,
	).Scan(&level, &xp)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, false, ErrProfessionNotLearned
		}
		return nil, false, err
	}
	leveled := false
	xp += delta
	for {
		needed := ProfessionXPToNext(level)
		if xp < needed || level >= maxLevel {
			break
		}
		xp -= needed
		level++
		leveled = true
	}
	_, err = tx.Exec(ctx, `
		UPDATE crafting.character_profession
		   SET level = $3, xp = $4, updated_at = now(), version = version + 1
		 WHERE character_id = $1 AND profession_code = $2`,
		characterID, professionCode, level, xp)
	if err != nil {
		return nil, false, err
	}
	if err := tx.Commit(ctx); err != nil {
		return nil, false, err
	}
	return r.getAfter(ctx, characterID, professionCode, leveled)
}

func (r *PgRepo) getAfter(ctx context.Context, characterID uuid.UUID, professionCode string, leveled bool) (*CharacterProfession, bool, error) {
	cp, err := r.Get(ctx, characterID, professionCode)
	if err != nil {
		return nil, false, err
	}
	return cp, leveled, nil
}

func (r *PgRepo) InsertCraftLog(ctx context.Context, characterID uuid.UUID, recipeCode string, outputID *uuid.UUID, goldTxnID *uuid.UUID, now time.Time) error {
	const q = `
		INSERT INTO crafting.craft_log
		  (character_id, recipe_code, output_instance_id, gold_txn_id, crafted_at)
		VALUES ($1, $2, $3, $4, $5)`
	_, err := r.pool.Exec(ctx, q, characterID, recipeCode, outputID, goldTxnID, now)
	return err
}
