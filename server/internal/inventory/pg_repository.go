package inventory

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/otter/server/internal/item"
)

// PgRepo persists inventory operations on inventory.item_instance.
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

const bagSelect = `
	SELECT i.id, i.template_code, t.display_name_key, t.item_class, t.item_subclass,
	       COALESCE(i.slot, -1), i.equipment_slot,
	       i.stack_quantity, t.max_stack, i.rarity,
	       i.prefix_code, i.suffix_code, i.is_bound, i.is_locked, i.lock_owner,
	       t.base_weight_g, t.base_value_copper,
	       t.stat_bonuses, t.equip_slot, t.class_restriction, t.min_level, t.use_effect
	  FROM inventory.item_instance i
	  JOIN content.item_template t ON t.code = i.template_code`

func (r *PgRepo) scanBagItem(row pgx.Row) (BagItem, error) {
	var b BagItem
	err := row.Scan(
		&b.InstanceID, &b.TemplateCode, &b.DisplayNameKey, &b.ItemClass, &b.ItemSubclass,
		&b.Slot, &b.EquipmentSlot,
		&b.StackQuantity, &b.MaxStack, &b.Rarity,
		&b.Prefix, &b.Suffix, &b.IsBound, &b.IsLocked, &b.LockOwner,
		&b.WeightGramsPer, &b.BaseValueCopper,
		&b.StatBonusesRaw, &b.EquipSlotOnTmpl, &b.ClassRestriction, &b.MinLevel, &b.UseEffect,
	)
	return b, err
}

func (r *PgRepo) Bag(ctx context.Context, characterID uuid.UUID) ([]BagItem, error) {
	rows, err := r.pool.Query(ctx,
		bagSelect+`
		 WHERE i.owner_character_id = $1 AND i.slot IS NOT NULL
		 ORDER BY i.slot`, characterID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []BagItem{}
	for rows.Next() {
		b, err := r.scanBagItem(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, b)
	}
	return out, rows.Err()
}

func (r *PgRepo) Unslotted(ctx context.Context, characterID uuid.UUID) ([]item.Instance, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, template_code, stack_quantity FROM inventory.item_instance
		 WHERE owner_character_id = $1 AND slot IS NULL AND equipment_slot IS NULL AND is_locked = false
		 ORDER BY created_at`, characterID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []item.Instance{}
	for rows.Next() {
		var in item.Instance
		if err := rows.Scan(&in.ID, &in.TemplateCode, &in.StackQuantity); err != nil {
			return nil, err
		}
		out = append(out, in)
	}
	return out, rows.Err()
}

func (r *PgRepo) OccupiedSlots(ctx context.Context, characterID uuid.UUID) ([]Slot, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT slot FROM inventory.item_instance
		  WHERE owner_character_id = $1 AND slot IS NOT NULL
		  ORDER BY slot`, characterID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []Slot{}
	for rows.Next() {
		var s Slot
		if err := rows.Scan(&s); err != nil {
			return nil, err
		}
		out = append(out, s)
	}
	return out, rows.Err()
}

func (r *PgRepo) StackableMergeCandidates(ctx context.Context, characterID uuid.UUID, templateCode string) ([]BagItem, error) {
	rows, err := r.pool.Query(ctx,
		bagSelect+`
		 WHERE i.owner_character_id = $1
		   AND i.template_code = $2
		   AND i.slot IS NOT NULL
		   AND i.is_locked = false
		   AND t.stackable = true
		   AND i.stack_quantity < t.max_stack
		 ORDER BY i.slot`,
		characterID, templateCode)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []BagItem{}
	for rows.Next() {
		b, err := r.scanBagItem(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, b)
	}
	return out, rows.Err()
}

func (r *PgRepo) AssignSlot(ctx context.Context, itemID, ownerID uuid.UUID, slot Slot) error {
	const q = `UPDATE inventory.item_instance
	              SET slot = $3, updated_at = now(), version = version + 1
	            WHERE id = $1 AND owner_character_id = $2 AND slot IS NULL`
	tag, err := r.pool.Exec(ctx, q, itemID, ownerID, slot)
	if err != nil {
		if pgUniqueViolation(err) {
			return ErrSlotOccupied
		}
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *PgRepo) AddStack(ctx context.Context, itemID uuid.UUID, delta int32) error {
	const q = `UPDATE inventory.item_instance
	              SET stack_quantity = stack_quantity + $2,
	                  updated_at = now(), version = version + 1
	            WHERE id = $1 AND is_locked = false`
	tag, err := r.pool.Exec(ctx, q, itemID, delta)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrLocked
	}
	return nil
}

func (r *PgRepo) SwapSlots(ctx context.Context, ownerID, itemA, itemB uuid.UUID) error {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable})
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	// Read both items' slots inside the TX.
	var slotA, slotB Slot
	var lockedA, lockedB bool
	err = tx.QueryRow(ctx,
		`SELECT slot, is_locked FROM inventory.item_instance WHERE id=$1 AND owner_character_id=$2`,
		itemA, ownerID).Scan(&slotA, &lockedA)
	if err != nil {
		return notFoundOrErr(err)
	}
	err = tx.QueryRow(ctx,
		`SELECT slot, is_locked FROM inventory.item_instance WHERE id=$1 AND owner_character_id=$2`,
		itemB, ownerID).Scan(&slotB, &lockedB)
	if err != nil {
		return notFoundOrErr(err)
	}
	if lockedA || lockedB {
		return ErrLocked
	}

	// Two-phase swap: move A to a sentinel slot (negative) to bypass unique
	// constraint, then B → A's slot, then A → B's original slot.
	const sentinel Slot = 200 // outside bag range; still valid CHECK
	if _, err := tx.Exec(ctx,
		`UPDATE inventory.item_instance SET slot=$2 WHERE id=$1`, itemA, sentinel); err != nil {
		return err
	}
	if _, err := tx.Exec(ctx,
		`UPDATE inventory.item_instance SET slot=$2 WHERE id=$1`, itemB, slotA); err != nil {
		return err
	}
	if _, err := tx.Exec(ctx,
		`UPDATE inventory.item_instance SET slot=$2 WHERE id=$1`, itemA, slotB); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (r *PgRepo) MoveToEmpty(ctx context.Context, itemID, ownerID uuid.UUID, target Slot) error {
	const q = `UPDATE inventory.item_instance
	              SET slot = $3, updated_at = now(), version = version + 1
	            WHERE id = $1 AND owner_character_id = $2 AND is_locked = false`
	tag, err := r.pool.Exec(ctx, q, itemID, ownerID, target)
	if err != nil {
		if pgUniqueViolation(err) {
			return ErrSlotOccupied
		}
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *PgRepo) SplitStack(ctx context.Context, sourceID, ownerID uuid.UUID, amount int32, target Slot) (uuid.UUID, error) {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable})
	if err != nil {
		return uuid.Nil, err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	// Read source with lock.
	var stack int32
	var stackable bool
	var locked bool
	var tmplCode, rarity string
	var prefix, suffix *string
	var isBound bool
	err = tx.QueryRow(ctx,
		`SELECT i.stack_quantity, t.stackable, i.is_locked, i.template_code,
		        i.rarity, i.prefix_code, i.suffix_code, i.is_bound
		   FROM inventory.item_instance i
		   JOIN content.item_template t ON t.code = i.template_code
		  WHERE i.id = $1 AND i.owner_character_id = $2 FOR UPDATE`,
		sourceID, ownerID).Scan(&stack, &stackable, &locked, &tmplCode, &rarity, &prefix, &suffix, &isBound)
	if err != nil {
		return uuid.Nil, notFoundOrErr(err)
	}
	if locked {
		return uuid.Nil, ErrLocked
	}
	if !stackable {
		return uuid.Nil, ErrNotStackable
	}
	if amount <= 0 || amount >= stack {
		return uuid.Nil, ErrInvalidSplit
	}

	// Decrement source.
	if _, err := tx.Exec(ctx,
		`UPDATE inventory.item_instance
		    SET stack_quantity = stack_quantity - $2,
		        updated_at = now(), version = version + 1
		  WHERE id = $1`, sourceID, amount); err != nil {
		return uuid.Nil, err
	}

	// Insert new instance at target slot.
	newID := uuid.New()
	if _, err := tx.Exec(ctx, `
		INSERT INTO inventory.item_instance
		  (id, template_code, owner_character_id, stack_quantity, rarity,
		   prefix_code, suffix_code, is_bound, source_kind, slot)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, 'split', $9)`,
		newID, tmplCode, ownerID, amount, rarity, prefix, suffix, isBound, target); err != nil {
		if pgUniqueViolation(err) {
			return uuid.Nil, ErrSlotOccupied
		}
		return uuid.Nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return uuid.Nil, err
	}
	return newID, nil
}

// --- Equipment methods (EquippedRepo) ---------------------------------------

func (r *PgRepo) Equipped(ctx context.Context, characterID uuid.UUID) ([]BagItem, error) {
	rows, err := r.pool.Query(ctx,
		bagSelect+`
		 WHERE i.owner_character_id = $1 AND i.equipment_slot IS NOT NULL
		 ORDER BY i.equipment_slot`, characterID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []BagItem{}
	for rows.Next() {
		b, err := r.scanBagItem(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, b)
	}
	return out, rows.Err()
}

func (r *PgRepo) EquipItem(ctx context.Context, itemID, ownerID uuid.UUID, equipSlot string) error {
	const q = `UPDATE inventory.item_instance
	              SET slot = NULL, equipment_slot = $3,
	                  is_bound = CASE WHEN (SELECT bind_kind FROM content.item_template WHERE code = template_code) = 'boe'
	                                  THEN true ELSE is_bound END,
	                  updated_at = now(), version = version + 1
	            WHERE id = $1 AND owner_character_id = $2 AND is_locked = false`
	tag, err := r.pool.Exec(ctx, q, itemID, ownerID, equipSlot)
	if err != nil {
		if pgUniqueViolation(err) {
			return ErrSlotOccupied
		}
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *PgRepo) UnequipItem(ctx context.Context, itemID, ownerID uuid.UUID, targetBagSlot Slot) error {
	const q = `UPDATE inventory.item_instance
	              SET equipment_slot = NULL, slot = $3,
	                  updated_at = now(), version = version + 1
	            WHERE id = $1 AND owner_character_id = $2 AND is_locked = false`
	tag, err := r.pool.Exec(ctx, q, itemID, ownerID, targetBagSlot)
	if err != nil {
		if pgUniqueViolation(err) {
			return ErrSlotOccupied
		}
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *PgRepo) FindEquippedInSlot(ctx context.Context, ownerID uuid.UUID, equipSlot string) (uuid.UUID, error) {
	var id uuid.UUID
	err := r.pool.QueryRow(ctx,
		`SELECT id FROM inventory.item_instance
		  WHERE owner_character_id = $1 AND equipment_slot = $2`,
		ownerID, equipSlot).Scan(&id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return uuid.Nil, ErrNotFound
		}
		return uuid.Nil, err
	}
	return id, nil
}

func (r *PgRepo) DeleteInstance(ctx context.Context, itemID, ownerID uuid.UUID) error {
	const q = `DELETE FROM inventory.item_instance
	            WHERE id = $1 AND owner_character_id = $2 AND is_locked = false`
	tag, err := r.pool.Exec(ctx, q, itemID, ownerID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

// ConsumeByTemplate reduces `quantity` units of an item template from the
// character's unlocked, un-equipped bag stacks. All work happens in one
// transaction; ErrInsufficient rolls back so no partial consume is visible.
//
// Selection order: smallest stacks first (asc stack_quantity). This drains
// fragments before touching pristine full stacks, keeping bag layout tidy.
func (r *PgRepo) ConsumeByTemplate(ctx context.Context, ownerID uuid.UUID, templateCode string, quantity int32) error {
	if quantity <= 0 {
		return nil
	}
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted})
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	// Lock candidate rows FOR UPDATE so two concurrent crafts don't over-eat.
	rows, err := tx.Query(ctx, `
		SELECT id, stack_quantity
		  FROM inventory.item_instance
		 WHERE owner_character_id = $1
		   AND template_code = $2
		   AND is_locked = false
		   AND equipment_slot IS NULL
		   AND slot IS NOT NULL
		 ORDER BY stack_quantity ASC
		 FOR UPDATE`, ownerID, templateCode)
	if err != nil {
		return err
	}
	type row struct {
		id  uuid.UUID
		qty int32
	}
	var picks []row
	var total int32
	for rows.Next() {
		var pr row
		if err := rows.Scan(&pr.id, &pr.qty); err != nil {
			rows.Close()
			return err
		}
		picks = append(picks, pr)
		total += pr.qty
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return err
	}
	if total < quantity {
		return ErrInsufficient
	}
	remaining := quantity
	for _, p := range picks {
		if remaining <= 0 {
			break
		}
		take := p.qty
		if take > remaining {
			take = remaining
		}
		if take == p.qty {
			// Stack fully drained → delete.
			if _, err := tx.Exec(ctx,
				`DELETE FROM inventory.item_instance WHERE id = $1`, p.id); err != nil {
				return err
			}
		} else {
			if _, err := tx.Exec(ctx, `
				UPDATE inventory.item_instance
				   SET stack_quantity = stack_quantity - $2, updated_at = now(), version = version + 1
				 WHERE id = $1`, p.id, take); err != nil {
				return err
			}
		}
		remaining -= take
	}
	return tx.Commit(ctx)
}

func notFoundOrErr(err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrNotFound
	}
	return err
}
