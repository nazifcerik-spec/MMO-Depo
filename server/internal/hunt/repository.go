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
	// MarkRevealed also settles the held supplies: used units are spent,
	// the rest goes back to the bag (DD-288).
	MarkRevealed(ctx context.Context, id uuid.UUID, result json.RawMessage, txnID uuid.UUID, now time.Time, used []SupplyUse) error
	// MarkCancelled returns every held supply to the bag.
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
	       started_at, expires_at, seed, state, player_snapshot, loadout,
	       revealed_at, reveal_result, reveal_txn_id,
	       cancelled_at, created_at, updated_at, version
	  FROM hunt.contract`

func (r *PgRepo) scan(row pgx.Row) (*Contract, error) {
	var c Contract
	var state string
	var loadout []byte
	err := row.Scan(
		&c.ID, &c.CharacterID, &c.MobCode, &c.EncounterSecs, &c.DurationSecs,
		&c.StartedAt, &c.ExpiresAt, &c.Seed, &state, &c.PlayerSnapshot, &loadout,
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
	if len(loadout) > 0 {
		c.Loadout = &Loadout{}
		if err := json.Unmarshal(loadout, c.Loadout); err != nil {
			return nil, err
		}
	}
	return &c, nil
}

func (r *PgRepo) Create(ctx context.Context, c *Contract) error {
	const q = `
		INSERT INTO hunt.contract (
		  id, character_id, mob_code, encounter_seconds, duration_seconds,
		  started_at, expires_at, seed, state, player_snapshot, loadout,
		  created_at, updated_at
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13)`
	var loadout []byte
	if c.Loadout != nil {
		var err error
		if loadout, err = json.Marshal(c.Loadout); err != nil {
			return err
		}
	}
	err := pgx.BeginFunc(ctx, r.pool, func(tx pgx.Tx) error {
		if _, err := tx.Exec(ctx, q,
			c.ID, c.CharacterID, c.MobCode, c.EncounterSecs, c.DurationSecs,
			c.StartedAt, c.ExpiresAt, c.Seed, string(c.State), c.PlayerSnapshot, loadout,
			c.CreatedAt, c.UpdatedAt,
		); err != nil {
			return err
		}
		if c.Loadout == nil {
			return nil
		}
		for _, s := range c.Loadout.Supplies {
			if err := holdSupply(ctx, tx, c.CharacterID, c.ID, s.Code, s.Qty); err != nil {
				return err
			}
		}
		return nil
	})
	switch {
	case err == nil, errors.Is(err, ErrNoSupplies):
		return err
	case pgUniqueViolation(err):
		return ErrSlotBusy
	default:
		return action.FromDB(err) // busy with another action (DD-263)
	}
}

// holdSupply takes qty units of a template out of the character's bag
// (smallest stacks first) into instances locked to the hunt (DD-288): whole
// stacks are locked, a partial one is split.
//
// SupplyStock is how many units of a consumable the bag holds free.
func (r *PgRepo) SupplyStock(ctx context.Context, characterID uuid.UUID, code string) (int32, error) {
	var n int32
	err := r.pool.QueryRow(ctx, `
		SELECT COALESCE(SUM(stack_quantity), 0)::int FROM inventory.item_instance
		 WHERE owner_character_id = $1 AND template_code = $2 AND is_locked = false
		   AND equipment_slot IS NULL`, characterID, code).Scan(&n)
	return n, err
}

func holdSupply(ctx context.Context, tx pgx.Tx, characterID, huntID uuid.UUID, code string, qty int32) error {
	rows, err := tx.Query(ctx, `
		SELECT id, stack_quantity FROM inventory.item_instance
		 WHERE owner_character_id = $1 AND template_code = $2 AND is_locked = false
		   AND equipment_slot IS NULL
		 ORDER BY stack_quantity, id FOR UPDATE`, characterID, code)
	if err != nil {
		return err
	}
	type stack struct {
		id  uuid.UUID
		qty int32
	}
	stacks, err := pgx.CollectRows(rows, func(r pgx.CollectableRow) (stack, error) {
		var s stack
		return s, r.Scan(&s.id, &s.qty)
	})
	if err != nil {
		return err
	}
	need := qty
	for _, s := range stacks {
		if need == 0 {
			break
		}
		if s.qty <= need {
			if _, err := tx.Exec(ctx, `
				UPDATE inventory.item_instance
				   SET is_locked = true, lock_owner = 'hunt', locked_ref_id = $2, slot = NULL,
				       updated_at = now(), version = version + 1
				 WHERE id = $1`, s.id, huntID); err != nil {
				return err
			}
			need -= s.qty
			continue
		}
		if _, err := tx.Exec(ctx, `
			UPDATE inventory.item_instance SET stack_quantity = stack_quantity - $2, updated_at = now(),
			       version = version + 1 WHERE id = $1`, s.id, need); err != nil {
			return err
		}
		if _, err := tx.Exec(ctx, `
			INSERT INTO inventory.item_instance
			  (template_code, owner_character_id, stack_quantity, rarity, prefix_code, suffix_code,
			   is_bound, is_locked, lock_owner, locked_ref_id, source_kind, source_ref_id)
			SELECT template_code, owner_character_id, $2, rarity, prefix_code, suffix_code,
			       is_bound, true, 'hunt', $3, source_kind, source_ref_id
			  FROM inventory.item_instance WHERE id = $1`, s.id, need, huntID); err != nil {
			return err
		}
		need = 0
	}
	if need > 0 {
		return ErrNoSupplies
	}
	return nil
}

// settleSupplies spends the used units of a hunt's held supplies and
// unlocks the rest back to the bag (placed on the next bag view).
func settleSupplies(ctx context.Context, tx pgx.Tx, huntID uuid.UUID, used []SupplyUse) error {
	for _, u := range used {
		rows, err := tx.Query(ctx, `
			SELECT id, stack_quantity FROM inventory.item_instance
			 WHERE lock_owner = 'hunt' AND locked_ref_id = $1 AND template_code = $2 AND is_locked
			 ORDER BY stack_quantity, id FOR UPDATE`, huntID, u.Code)
		if err != nil {
			return err
		}
		type stack struct {
			id  uuid.UUID
			qty int32
		}
		stacks, err := pgx.CollectRows(rows, func(r pgx.CollectableRow) (stack, error) {
			var s stack
			return s, r.Scan(&s.id, &s.qty)
		})
		if err != nil {
			return err
		}
		need := u.Qty
		for _, s := range stacks {
			if need == 0 {
				break
			}
			if s.qty <= need {
				_, err = tx.Exec(ctx, `DELETE FROM inventory.item_instance WHERE id = $1`, s.id)
				need -= s.qty
			} else {
				_, err = tx.Exec(ctx, `UPDATE inventory.item_instance SET stack_quantity = stack_quantity - $2,
					updated_at = now(), version = version + 1 WHERE id = $1`, s.id, need)
				need = 0
			}
			if err != nil {
				return err
			}
		}
	}
	_, err := tx.Exec(ctx, `
		UPDATE inventory.item_instance
		   SET is_locked = false, lock_owner = NULL, locked_ref_id = NULL, slot = NULL,
		       updated_at = now(), version = version + 1
		 WHERE lock_owner = 'hunt' AND locked_ref_id = $1`, huntID)
	return err
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

func (r *PgRepo) MarkRevealed(ctx context.Context, id uuid.UUID, result json.RawMessage, txnID uuid.UUID, now time.Time, used []SupplyUse) error {
	const q = `
		UPDATE hunt.contract
		   SET state = 'revealed', revealed_at = $2, reveal_result = $3,
		       reveal_txn_id = NULLIF($4, '00000000-0000-0000-0000-000000000000'::uuid),
		       updated_at = $2, version = version + 1
		 WHERE id = $1 AND state = 'active'`
	return pgx.BeginFunc(ctx, r.pool, func(tx pgx.Tx) error {
		tag, err := tx.Exec(ctx, q, id, now, result, txnID)
		if err != nil {
			return err
		}
		if tag.RowsAffected() == 0 {
			return ErrNotActive
		}
		return settleSupplies(ctx, tx, id, used)
	})
}

func (r *PgRepo) MarkCancelled(ctx context.Context, id uuid.UUID, now time.Time) error {
	const q = `
		UPDATE hunt.contract
		   SET state = 'cancelled', cancelled_at = $2, updated_at = $2, version = version + 1
		 WHERE id = $1 AND state = 'active'`
	return pgx.BeginFunc(ctx, r.pool, func(tx pgx.Tx) error {
		tag, err := tx.Exec(ctx, q, id, now)
		if err != nil {
			return err
		}
		if tag.RowsAffected() == 0 {
			return ErrNotActive
		}
		return settleSupplies(ctx, tx, id, nil)
	})
}
