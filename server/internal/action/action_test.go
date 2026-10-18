package action

import (
	"errors"
	"fmt"
	"testing"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/require"
)

func TestFromDB(t *testing.T) {
	busy := FromDB(fmt.Errorf("insert: %w", &pgconn.PgError{Code: "OT001", Hint: "hunt"}))
	require.ErrorIs(t, busy, ErrBusy)
	require.Equal(t, "hunt", Kind(busy))
	require.Equal(t, "action: character is busy with hunt", busy.Error())

	other := &pgconn.PgError{Code: "23505"}
	require.Same(t, other, FromDB(other))
	require.NoError(t, FromDB(nil))
	require.Empty(t, Kind(errors.New("x")))
}
