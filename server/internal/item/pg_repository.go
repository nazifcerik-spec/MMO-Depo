package item

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// PgRepo persists inventory.item_instance rows.
type PgRepo struct {
	pool *pgxpool.Pool
}

func NewPgRepo(pool *pgxpool.Pool) *PgRepo {
	return &PgRepo{pool: pool}
}

const insertQ = `
	INSERT INTO inventory.item_instance
	  (id, template_code, owner_character_id, stack_quantity, rarity,
	   prefix_code, suffix_code, is_bound, source_kind, source_ref_id,
	   created_at, updated_at)
	VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)`

// CreateInstancesBatch uses pgx.Batch to insert many rows in a single round
// trip — critical for activity reveal which may create 10-50 items at once.
func (r *PgRepo) CreateInstancesBatch(ctx context.Context, instances []Instance) error {
	if len(instances) == 0 {
		return nil
	}
	batch := &pgx.Batch{}
	for _, i := range instances {
		batch.Queue(insertQ,
			i.ID, i.TemplateCode, i.OwnerCharacterID, i.StackQuantity, string(i.Rarity),
			i.PrefixCode, i.SuffixCode, i.IsBound, i.SourceKind, i.SourceRefID,
			i.CreatedAt, i.UpdatedAt,
		)
	}
	br := r.pool.SendBatch(ctx, batch)
	defer br.Close()
	for range instances {
		if _, err := br.Exec(); err != nil {
			return err
		}
	}
	return nil
}

const selectBase = `
	SELECT id, template_code, owner_character_id, stack_quantity, rarity,
	       prefix_code, suffix_code, is_bound, is_locked, lock_owner,
	       locked_ref_id, source_kind, source_ref_id, equipment_slot,
	       created_at, updated_at, version
	  FROM inventory.item_instance`

func (r *PgRepo) scan(row pgx.Row) (*Instance, error) {
	var i Instance
	var rarity string
	if err := row.Scan(
		&i.ID, &i.TemplateCode, &i.OwnerCharacterID, &i.StackQuantity, &rarity,
		&i.PrefixCode, &i.SuffixCode, &i.IsBound, &i.IsLocked, &i.LockOwner,
		&i.LockedRefID, &i.SourceKind, &i.SourceRefID, &i.EquipmentSlot,
		&i.CreatedAt, &i.UpdatedAt, &i.Version,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	i.Rarity = Rarity(rarity)
	return &i, nil
}

func (r *PgRepo) Get(ctx context.Context, id uuid.UUID) (*Instance, error) {
	return r.scan(r.pool.QueryRow(ctx, selectBase+" WHERE id = $1", id))
}

func (r *PgRepo) ListByOwner(ctx context.Context, characterID uuid.UUID, limit int) ([]Instance, error) {
	if limit <= 0 {
		limit = 100
	}
	rows, err := r.pool.Query(ctx,
		selectBase+" WHERE owner_character_id = $1 ORDER BY created_at DESC LIMIT $2",
		characterID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []Instance{}
	for rows.Next() {
		i, err := r.scan(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, *i)
	}
	return out, rows.Err()
}

func (r *PgRepo) ListBySource(ctx context.Context, sourceRefID uuid.UUID) ([]Instance, error) {
	rows, err := r.pool.Query(ctx,
		selectBase+" WHERE source_ref_id = $1 ORDER BY created_at ASC",
		sourceRefID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []Instance{}
	for rows.Next() {
		i, err := r.scan(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, *i)
	}
	return out, rows.Err()
}
