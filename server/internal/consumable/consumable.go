// Package consumable applies a consumable item's use effect (DD-218).
//
// One transaction does everything: the item row and the character row are
// locked, the effect is applied against the character's current HP (stored +
// analytical regen, see character.EffectiveHP) and one unit is removed from the
// stack. The HTTP route sits behind the Idempotency middleware, so a retried
// request replays the first response instead of consuming a second unit.
package consumable

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/otter/server/internal/domain/character"
	"github.com/otter/server/internal/item"
)

var (
	ErrNotFound    = errors.New("consumable: item not found")
	ErrLocked      = errors.New("consumable: item is locked")
	ErrNotUsable   = errors.New("consumable: item has no use effect")
	ErrLevel       = errors.New("consumable: character level too low")
	ErrInCombat    = errors.New("consumable: character is in a hunt or dungeon")
	ErrNoEffect    = errors.New("consumable: effect would do nothing (hp already full)")
	errBadTemplate = errors.New("consumable: unknown effect kind")
)

// Catalog resolves item templates (item.Service).
type Catalog interface {
	Catalog() *item.Catalog
}

// SheetReader returns the live sheet (ownership check + max HP + vitality).
type SheetReader interface {
	Sheet(ctx context.Context, characterID, accountID uuid.UUID) (*character.Sheet, error)
}

type Service struct {
	pool   *pgxpool.Pool
	items  Catalog
	sheets SheetReader
	now    func() time.Time
	log    *slog.Logger
}

func NewService(pool *pgxpool.Pool, items Catalog, sheets SheetReader, logger *slog.Logger) *Service {
	if logger == nil {
		logger = slog.Default()
	}
	return &Service{pool: pool, items: items, sheets: sheets, now: time.Now, log: logger}
}

// Result is the HTTP response of a successful use.
type Result struct {
	TemplateCode   string `json:"template_code"`
	HPCurrent      int    `json:"hp_current"`
	HPMax          int    `json:"hp_max"`
	StackRemaining int32  `json:"stack_remaining"` // 0 = the stack is gone
}

// ApplyHeal is the pure effect step: new HP after healing `amount` from
// `current`, capped at max. ErrNoEffect when nothing would change.
func ApplyHeal(current, maxHP, amount int) (int, error) {
	if current >= maxHP {
		return current, ErrNoEffect
	}
	return min(maxHP, current+amount), nil
}

// Use consumes one unit of itemID for characterID (owned by accountID).
func (s *Service) Use(ctx context.Context, characterID, accountID, itemID uuid.UUID) (*Result, error) {
	sheet, err := s.sheets.Sheet(ctx, characterID, accountID)
	if err != nil {
		return nil, err
	}
	now := s.now().UTC()
	var res Result
	err = pgx.BeginFunc(ctx, s.pool, func(tx pgx.Tx) error {
		var (
			code   string
			stack  int32
			locked bool
			owner  *uuid.UUID
			slot   *int16
		)
		err := tx.QueryRow(ctx, `
			SELECT template_code, stack_quantity, is_locked, owner_character_id, slot
			  FROM inventory.item_instance WHERE id = $1 FOR UPDATE`, itemID).
			Scan(&code, &stack, &locked, &owner, &slot)
		if errors.Is(err, pgx.ErrNoRows) || (err == nil && (owner == nil || *owner != characterID || slot == nil)) {
			return ErrNotFound // equipped items (slot NULL) are not in the bag
		}
		if err != nil {
			return err
		}
		if locked {
			return ErrLocked
		}
		tmpl, err := s.items.Catalog().Get(code)
		if err != nil || tmpl.UseEffect == nil {
			return ErrNotUsable
		}
		if int(tmpl.MinLevel) > sheet.Level {
			return ErrLevel
		}

		var inFight bool
		if err := tx.QueryRow(ctx, `
			SELECT EXISTS (SELECT 1 FROM hunt.contract    WHERE character_id = $1 AND state = 'active')
			    OR EXISTS (SELECT 1 FROM dungeon.contract WHERE character_id = $1 AND state = 'active')`,
			characterID).Scan(&inFight); err != nil {
			return err
		}
		if inFight {
			return ErrInCombat // the reveal would overwrite the healed HP
		}

		var stored *int32
		var storedAt *time.Time
		if err := tx.QueryRow(ctx, `
			SELECT hp_current, hp_updated_at FROM character.character WHERE id = $1 FOR UPDATE`,
			characterID).Scan(&stored, &storedAt); err != nil {
			return err
		}
		maxHP := sheet.Resources.HPMax
		current := character.EffectiveHP(stored, storedAt, maxHP, sheet.Primary.Vitality, now)

		switch tmpl.UseEffect.Kind {
		case "heal_hp":
			hp, err := ApplyHeal(current, maxHP, tmpl.UseEffect.Amount)
			if err != nil {
				return err
			}
			if _, err := tx.Exec(ctx, `
				UPDATE character.character
				   SET hp_current = $2, hp_updated_at = $3, updated_at = now(), version = version + 1
				 WHERE id = $1`, characterID, hp, now); err != nil {
				return err
			}
			res.HPCurrent = hp
		default:
			return fmt.Errorf("%w %q", errBadTemplate, tmpl.UseEffect.Kind)
		}

		if stack > 1 {
			_, err = tx.Exec(ctx, `UPDATE inventory.item_instance SET stack_quantity = stack_quantity - 1 WHERE id = $1`, itemID)
		} else {
			_, err = tx.Exec(ctx, `DELETE FROM inventory.item_instance WHERE id = $1`, itemID)
		}
		if err != nil {
			return err
		}
		res.TemplateCode, res.HPMax, res.StackRemaining = code, maxHP, stack-1
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &res, nil
}
