// Package action is the one-action-at-a-time rule (DD-263): a character runs
// one of gathering, hunt, dungeon, raid, travel or crafting at a time. The
// database enforces it (character.claim_action trigger, migration 0039);
// this package turns the trigger's error into a domain error.
package action

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
)

// ErrBusy: the character is already doing something else.
var ErrBusy = errors.New("action: character is busy with another action")

// BusyError names what the character is busy with (gathering, hunt,
// dungeon, raid, travel or crafting).
type BusyError struct{ Kind string }

func (e *BusyError) Error() string        { return "action: character is busy with " + e.Kind }
func (e *BusyError) Is(target error) bool { return target == ErrBusy }

// sqlState is the claim_action trigger's error code.
const sqlState = "OT001"

// FromDB returns a *BusyError for the trigger's refusal and err otherwise.
func FromDB(err error) error {
	var pg *pgconn.PgError
	if errors.As(err, &pg) && pg.Code == sqlState {
		return &BusyError{Kind: pg.Hint}
	}
	return err
}

// Kind is the busy kind of err ("" when err is not a busy refusal).
func Kind(err error) string {
	var b *BusyError
	if errors.As(err, &b) {
		return b.Kind
	}
	return ""
}
