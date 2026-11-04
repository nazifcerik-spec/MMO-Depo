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

// GuildBankSlotCap is the size of a guild's bank (DD-358).
const GuildBankSlotCap = 100

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
type movable struct {
	tmpl, rarity    string
	prefix, suffix  *string
	enchant         *string
	stack, maxStack int32
	stackable       bool
}

// Moved is what a store move carried, for the guild's log.
type Moved struct {
	Template string
	Rarity   string
	Qty      int32
}

// store is a shared item store: the account bank or a guild's bank (DD-358).
// Its strings are constants of this file, never input.
type store struct {
	idCol, slotCol, lockSQL string
	id                      uuid.UUID
	cap                     int
	full                    error
}

func accountStore(id uuid.UUID) store {
	return store{"bank_account_id", "bank_slot", `SELECT 1 FROM identity.account WHERE id = $1 FOR UPDATE`, id, BankSlotCap, ErrBankFull}
}

func guildStore(id uuid.UUID) store {
	return store{"guild_bank_id", "guild_bank_slot", `SELECT 1 FROM social.guild WHERE id = $1 FOR UPDATE`, id, GuildBankSlotCap, ErrBankFull}
}

// Deposit moves qty (0 = all) of a bag item into the account bank.
func (s *Service) Deposit(ctx context.Context, characterID, accountID, itemID uuid.UUID, qty int32) error {
	_, err := s.depositTo(ctx, accountStore(accountID), characterID, itemID, qty)
	return err
}

// Withdraw moves qty (0 = all) of a banked item into the character's bag.
func (s *Service) Withdraw(ctx context.Context, characterID, accountID, itemID uuid.UUID, qty int32) error {
	_, err := s.withdrawFrom(ctx, accountStore(accountID), characterID, itemID, qty)
	return err
}

// GuildDeposit moves a bag item into a guild's bank.
func (s *Service) GuildDeposit(ctx context.Context, characterID, guildID, itemID uuid.UUID, qty int32) (Moved, error) {
	return s.depositTo(ctx, guildStore(guildID), characterID, itemID, qty)
}

// GuildWithdraw moves a guild's banked item into the character's bag.
func (s *Service) GuildWithdraw(ctx context.Context, characterID, guildID, itemID uuid.UUID, qty int32) (Moved, error) {
	return s.withdrawFrom(ctx, guildStore(guildID), characterID, itemID, qty)
}

// GuildBank lists a guild's banked items by slot.
func (s *Service) GuildBank(ctx context.Context, guildID uuid.UUID) (*BankView, error) {
	r, ok := s.repo.(*PgRepo)
	if !ok {
		return nil, errors.New("inventory: bank needs the Postgres repository")
	}
	sel := strings.Replace(bagSelect, "COALESCE(i.slot, -1)", "COALESCE(i.guild_bank_slot, -1)", 1)
	rows, err := r.pool.Query(ctx, sel+` WHERE i.guild_bank_id = $1 ORDER BY i.guild_bank_slot`, guildID)
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
	return &BankView{Items: items, SlotCap: GuildBankSlotCap, SlotsUsed: len(items)}, rows.Err()
}

// depositTo moves qty (0 = all) of a bag item into a store: onto a stored stack
// of the same item when it fits, else into the first free slot. Equipped,
// locked and bound items stay in the bag.
func (s *Service) depositTo(ctx context.Context, to store, characterID, itemID uuid.UUID, qty int32) (Moved, error) {
	var moved Moved
	r, ok := s.repo.(*PgRepo)
	if !ok {
		return moved, errors.New("inventory: bank needs the Postgres repository")
	}
	err := pgx.BeginFunc(ctx, r.pool, func(tx pgx.Tx) error {
		if _, err := tx.Exec(ctx, to.lockSQL, to.id); err != nil {
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
		moved = Moved{m.tmpl, m.rarity, qty}
		// Onto a stored stack with room.
		if m.stackable {
			var into uuid.UUID
			err := tx.QueryRow(ctx, `
				SELECT id FROM inventory.item_instance
				 WHERE `+to.idCol+` = $1 AND template_code = $2 AND rarity = $3
				   AND prefix_code IS NOT DISTINCT FROM $4 AND suffix_code IS NOT DISTINCT FROM $5
				   AND stack_quantity + $6 <= $7
				 ORDER BY `+to.slotCol+` LIMIT 1 FOR UPDATE`,
				to.id, m.tmpl, m.rarity, m.prefix, m.suffix, qty, m.maxStack).Scan(&into)
			if err == nil {
				return moveQty(ctx, tx, itemID, into, qty, m.stack)
			}
			if !errors.Is(err, pgx.ErrNoRows) {
				return err
			}
		}
		slot, err := freeStoreSlot(ctx, tx, to)
		if err != nil {
			return err
		}
		if qty == m.stack {
			_, err := tx.Exec(ctx, `
				UPDATE inventory.item_instance
				   SET owner_character_id = NULL, slot = NULL, `+to.idCol+` = $2, `+to.slotCol+` = $3,
				       updated_at = now(), version = version + 1
				 WHERE id = $1`, itemID, to.id, slot)
			return err
		}
		if _, err := tx.Exec(ctx, `UPDATE inventory.item_instance SET stack_quantity = stack_quantity - $2,
			updated_at = now(), version = version + 1 WHERE id = $1`, itemID, qty); err != nil {
			return err
		}
		_, err = tx.Exec(ctx, `
			INSERT INTO inventory.item_instance
			  (template_code, stack_quantity, rarity, prefix_code, suffix_code, is_bound, source_kind,
			   `+to.idCol+`, `+to.slotCol+`)
			SELECT template_code, $2, rarity, prefix_code, suffix_code, false, 'bank', $3, $4
			  FROM inventory.item_instance WHERE id = $1`, itemID, qty, to.id, slot)
		return err
	})
	return moved, err
}

// withdrawFrom moves qty (0 = all) of a stored item into the character's bag;
// the bag must have room (a free slot or space on its stacks).
func (s *Service) withdrawFrom(ctx context.Context, from store, characterID, itemID uuid.UUID, qty int32) (Moved, error) {
	var moved Moved
	r, ok := s.repo.(*PgRepo)
	if !ok {
		return moved, errors.New("inventory: bank needs the Postgres repository")
	}
	var tmpl string
	var stack int32
	err := r.pool.QueryRow(ctx, `SELECT template_code, stack_quantity FROM inventory.item_instance
		WHERE id = $1 AND `+from.idCol+` = $2`, itemID, from.id).Scan(&tmpl, &stack)
	if err != nil {
		return moved, notFoundOrErr(err)
	}
	if qty <= 0 || qty > stack {
		qty = stack
	}
	if ok, err := s.HasRoom(ctx, characterID, tmpl, qty); err != nil {
		return moved, err
	} else if !ok {
		return moved, ErrBagFull
	}
	err = pgx.BeginFunc(ctx, r.pool, func(tx pgx.Tx) error {
		if _, err := tx.Exec(ctx, from.lockSQL, from.id); err != nil {
			return err
		}
		var now int32
		if err := tx.QueryRow(ctx, `SELECT stack_quantity, rarity FROM inventory.item_instance
			WHERE id = $1 AND `+from.idCol+` = $2 FOR UPDATE`, itemID, from.id).Scan(&now, &moved.Rarity); err != nil {
			return notFoundOrErr(err)
		}
		if qty > now {
			qty = now
		}
		moved.Template, moved.Qty = tmpl, qty
		if qty == now {
			_, err := tx.Exec(ctx, `
				UPDATE inventory.item_instance
				   SET owner_character_id = $2, slot = NULL, `+from.idCol+` = NULL, `+from.slotCol+` = NULL,
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
		return moved, err
	}
	return moved, s.PlaceOverflow(ctx, characterID)
}

// PlaceOverflow puts the character's unslotted items into the bag: onto
// stacks, else free slots (items that arrived from the bank or a trade).
func (s *Service) PlaceOverflow(ctx context.Context, characterID uuid.UUID) error {
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

func freeStoreSlot(ctx context.Context, tx pgx.Tx, in store) (int16, error) {
	rows, err := tx.Query(ctx, `SELECT `+in.slotCol+` FROM inventory.item_instance
		WHERE `+in.idCol+` = $1 ORDER BY `+in.slotCol, in.id)
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
	if int(free) >= in.cap {
		return 0, in.full
	}
	return free, nil
}
