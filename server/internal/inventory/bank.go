package inventory

import (
	"context"
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// The account bank (DD-314): BankSlotCap slots every character of an
// account shares. A banked item has no owner character; it belongs to the
// account (bank_account_id, bank_slot), so deleting a character keeps it.
const BankSlotCap = 200

var (
	ErrBankFull = errors.New("inventory: bank full")
	ErrBound    = errors.New("inventory: a bound item stays with its character")
	ErrEquipped = errors.New("inventory: take the item off first")
)

// BankView is GET /bank: the account's banked items.
type BankView struct {
	Items     []BagItem `json:"items"` // Slot is the bank slot
	SlotCap   int       `json:"slot_cap"`
	SlotsUsed int       `json:"slots_used"`
}

var bankSelect = strings.Replace(bagSelect, "COALESCE(i.slot, -1)", "COALESCE(i.bank_slot, -1)", 1)

// Bank lists the account's banked items by slot.
func (s *Service) Bank(ctx context.Context, accountID uuid.UUID) (*BankView, error) {
	r, ok := s.repo.(*PgRepo)
	if !ok {
		return nil, errors.New("inventory: bank needs the Postgres repository")
	}
	rows, err := r.pool.Query(ctx, bankSelect+` WHERE i.bank_account_id = $1 ORDER BY i.bank_slot`, accountID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []BagItem{}
	for rows.Next() {
		b, err := r.scanBagItem(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, b)
	}
	return &BankView{Items: items, SlotCap: BankSlotCap, SlotsUsed: len(items)}, rows.Err()
}

// bankLock serialises an account's bank moves (two characters at once).
func bankLock(ctx context.Context, tx pgx.Tx, accountID uuid.UUID) error {
	_, err := tx.Exec(ctx, `SELECT pg_advisory_xact_lock(hashtext('bank:' || $1::text))`, accountID)
	return err
}

type movable struct {
	tmpl, rarity    string
	prefix, suffix  *string
	enchant         *string
	stack, maxStack int32
	stackable       bool
}

// Deposit moves qty (0 = all) of a bag item into the account bank: onto a
// banked stack of the same item when it fits, else into the first free
// bank slot. Equipped, locked and bound items stay in the bag.
func (s *Service) Deposit(ctx context.Context, characterID, accountID, itemID uuid.UUID, qty int32) error {
	r, ok := s.repo.(*PgRepo)
	if !ok {
		return errors.New("inventory: bank needs the Postgres repository")
	}
	return pgx.BeginFunc(ctx, r.pool, func(tx pgx.Tx) error {
		if err := bankLock(ctx, tx, accountID); err != nil {
			return err
		}
		var m movable
		var equipped *string
		var locked, bound bool
		err := tx.QueryRow(ctx, `
			SELECT i.template_code, i.rarity, i.prefix_code, i.suffix_code, i.enchant_code, i.stack_quantity,
			       t.max_stack, t.stackable, i.equipment_slot, i.is_locked, i.is_bound
			  FROM inventory.item_instance i JOIN content.item_template t ON t.code = i.template_code
			 WHERE i.id = $1 AND i.owner_character_id = $2 FOR UPDATE OF i`, itemID, characterID).
			Scan(&m.tmpl, &m.rarity, &m.prefix, &m.suffix, &m.enchant, &m.stack, &m.maxStack, &m.stackable,
				&equipped, &locked, &bound)
		if err != nil {
			return notFoundOrErr(err)
		}
		switch {
		case equipped != nil:
			return ErrEquipped
		case locked:
			return ErrLocked
		case bound:
			return ErrBound
		}
		if qty <= 0 || qty > m.stack {
			qty = m.stack
		}
		// Onto a banked stack with room.
		if m.stackable {
			var into uuid.UUID
			err := tx.QueryRow(ctx, `
				SELECT id FROM inventory.item_instance
				 WHERE bank_account_id = $1 AND template_code = $2 AND rarity = $3
				   AND prefix_code IS NOT DISTINCT FROM $4 AND suffix_code IS NOT DISTINCT FROM $5
				   AND stack_quantity + $6 <= $7
				 ORDER BY bank_slot LIMIT 1 FOR UPDATE`,
				accountID, m.tmpl, m.rarity, m.prefix, m.suffix, qty, m.maxStack).Scan(&into)
			if err == nil {
				return moveQty(ctx, tx, itemID, into, qty, m.stack)
			}
			if !errors.Is(err, pgx.ErrNoRows) {
				return err
			}
		}
		slot, err := freeBankSlot(ctx, tx, accountID)
		if err != nil {
			return err
		}
		if qty == m.stack {
			_, err := tx.Exec(ctx, `
				UPDATE inventory.item_instance
				   SET owner_character_id = NULL, slot = NULL, bank_account_id = $2, bank_slot = $3,
				       updated_at = now(), version = version + 1
				 WHERE id = $1`, itemID, accountID, slot)
			return err
		}
		if _, err := tx.Exec(ctx, `UPDATE inventory.item_instance SET stack_quantity = stack_quantity - $2,
			updated_at = now(), version = version + 1 WHERE id = $1`, itemID, qty); err != nil {
			return err
		}
		_, err = tx.Exec(ctx, `
			INSERT INTO inventory.item_instance
			  (template_code, stack_quantity, rarity, prefix_code, suffix_code, is_bound, source_kind,
			   bank_account_id, bank_slot)
			SELECT template_code, $2, rarity, prefix_code, suffix_code, false, 'bank', $3, $4
			  FROM inventory.item_instance WHERE id = $1`, itemID, qty, accountID, slot)
		return err
	})
}

// Withdraw moves qty (0 = all) of a banked item into the character's bag;
// the bag must have room (a free slot or space on its stacks).
func (s *Service) Withdraw(ctx context.Context, characterID, accountID, itemID uuid.UUID, qty int32) error {
	r, ok := s.repo.(*PgRepo)
	if !ok {
		return errors.New("inventory: bank needs the Postgres repository")
	}
	var tmpl string
	var stack int32
	err := r.pool.QueryRow(ctx, `SELECT template_code, stack_quantity FROM inventory.item_instance
		WHERE id = $1 AND bank_account_id = $2`, itemID, accountID).Scan(&tmpl, &stack)
	if err != nil {
		return notFoundOrErr(err)
	}
	if qty <= 0 || qty > stack {
		qty = stack
	}
	if ok, err := s.HasRoom(ctx, characterID, tmpl, qty); err != nil {
		return err
	} else if !ok {
		return ErrBagFull
	}
	err = pgx.BeginFunc(ctx, r.pool, func(tx pgx.Tx) error {
		if err := bankLock(ctx, tx, accountID); err != nil {
			return err
		}
		var now int32
		if err := tx.QueryRow(ctx, `SELECT stack_quantity FROM inventory.item_instance
			WHERE id = $1 AND bank_account_id = $2 FOR UPDATE`, itemID, accountID).Scan(&now); err != nil {
			return notFoundOrErr(err)
		}
		if qty > now {
			qty = now
		}
		if qty == now {
			_, err := tx.Exec(ctx, `
				UPDATE inventory.item_instance
				   SET owner_character_id = $2, slot = NULL, bank_account_id = NULL, bank_slot = NULL,
				       updated_at = now(), version = version + 1
				 WHERE id = $1`, itemID, characterID)
			return err
		}
		if _, err := tx.Exec(ctx, `UPDATE inventory.item_instance SET stack_quantity = stack_quantity - $2,
			updated_at = now(), version = version + 1 WHERE id = $1`, itemID, qty); err != nil {
			return err
		}
		_, err := tx.Exec(ctx, `
			INSERT INTO inventory.item_instance
			  (template_code, owner_character_id, stack_quantity, rarity, prefix_code, suffix_code, is_bound, source_kind)
			SELECT template_code, $2, $3, rarity, prefix_code, suffix_code, false, 'bank'
			  FROM inventory.item_instance WHERE id = $1`, itemID, characterID, qty)
		return err
	})
	if err != nil {
		return err
	}
	// Into the bag: onto stacks, else a free slot (the overflow path).
	overflow, err := s.repo.Unslotted(ctx, characterID)
	if err != nil || len(overflow) == 0 {
		return err
	}
	return s.AssignSlots(ctx, characterID, overflow)
}

// moveQty moves qty of item from onto stack into; all of it deletes from.
func moveQty(ctx context.Context, tx pgx.Tx, from, into uuid.UUID, qty, fromStack int32) error {
	if _, err := tx.Exec(ctx, `UPDATE inventory.item_instance SET stack_quantity = stack_quantity + $2,
		updated_at = now(), version = version + 1 WHERE id = $1`, into, qty); err != nil {
		return err
	}
	if qty == fromStack {
		_, err := tx.Exec(ctx, `DELETE FROM inventory.item_instance WHERE id = $1`, from)
		return err
	}
	_, err := tx.Exec(ctx, `UPDATE inventory.item_instance SET stack_quantity = stack_quantity - $2,
		updated_at = now(), version = version + 1 WHERE id = $1`, from, qty)
	return err
}

func freeBankSlot(ctx context.Context, tx pgx.Tx, accountID uuid.UUID) (int16, error) {
	rows, err := tx.Query(ctx, `SELECT bank_slot FROM inventory.item_instance
		WHERE bank_account_id = $1 ORDER BY bank_slot`, accountID)
	if err != nil {
		return 0, err
	}
	used, err := pgx.CollectRows(rows, pgx.RowTo[int16])
	if err != nil {
		return 0, err
	}
	free := int16(0)
	for _, s := range used {
		if s != free {
			break
		}
		free++
	}
	if free >= BankSlotCap {
		return 0, ErrBankFull
	}
	return free, nil
}
