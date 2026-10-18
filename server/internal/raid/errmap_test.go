package raid

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/otter/server/internal/domain/character"
)

// Every domain error must map to a client status: an unmapped one reaches
// players as a 500 "server error". Errors are wrapped like services do.
func TestDomainErrorsMapToClientStatus(t *testing.T) {
	for _, e := range []error{character.ErrForbidden, character.ErrNotFound,
		ErrUnknownRaid, ErrNotFound, ErrNoParty, ErrNotLeader, ErrNotInParty, ErrLobbyExists, ErrAlreadyInRaid,
		ErrNotForming, ErrRaidFull, ErrNotEnough, ErrNotRunning, ErrNotReady, ErrNotMember,
		ErrLevelRequirement, ErrWeeklyLockout, ErrTooWounded, ErrBagFull} {
		rec := httptest.NewRecorder()
		writeErr(rec, httptest.NewRequest(http.MethodGet, "/", nil), fmt.Errorf("context: %w", e))
		require.Less(t, rec.Code, http.StatusInternalServerError, e.Error())
	}
	rec := httptest.NewRecorder()
	writeErr(rec, httptest.NewRequest(http.MethodGet, "/", nil), errors.New("boom"))
	require.Equal(t, http.StatusInternalServerError, rec.Code)
}
