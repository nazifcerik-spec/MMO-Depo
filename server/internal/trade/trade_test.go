package trade

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

// The user's example: 10 ore at a 100 gold average pay 80 gold.
func TestTax(t *testing.T) {
	const gold = 10_000 // copper
	require.EqualValues(t, 80*gold, Tax(100*gold, 10))
	require.EqualValues(t, 8, Tax(100, 1))
	require.Zero(t, Tax(0, 5))
}

func TestDomainErrorsMapToClientStatus(t *testing.T) {
	for _, e := range []error{ErrNotFound, ErrSelf, ErrBusy, ErrClosed, ErrItem, ErrTooMany, ErrBagFull, ErrNoGold,
		ErrNotReadyToGo, errChangedMidway} {
		rec := httptest.NewRecorder()
		writeErr(rec, httptest.NewRequest(http.MethodGet, "/", nil), fmt.Errorf("ctx: %w", e))
		require.Less(t, rec.Code, 500, e.Error())
	}
}
