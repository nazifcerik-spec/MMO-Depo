package shop

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCatalogAndCodes(t *testing.T) {
	cat := Catalog()
	require.Len(t, cat, 3+(Backgrounds-FreeBackgrounds)+(Avatars-FreeAvatars))
	for _, c := range []struct {
		code, kind string
		free, ok   bool
	}{
		{"bg01", "background", true, true}, {"bg10", "background", true, true}, {"bg11", "background", false, true},
		{"bg30", "background", false, true}, {"bg31", "", false, false}, {"av40", "avatar", true, true},
		{"av41", "avatar", false, true}, {"av80", "avatar", false, true}, {"av00", "", false, false},
		{"av7", "", false, false}, {"boost_gather", "", false, false},
	} {
		kind, free, ok := cosmetic(c.code)
		require.Equal(t, c.ok, ok, c.code)
		require.Equal(t, c.kind, kind, c.code)
		require.Equal(t, c.free, free, c.code)
	}
	it, ok := find("boost_craft")
	require.True(t, ok)
	require.EqualValues(t, BoostPrice, it.Price)
	_, ok = find("bg01")
	require.False(t, ok, "free cosmetics are not for sale")
}

func TestDomainErrorsMapToClientStatus(t *testing.T) {
	for _, e := range []error{ErrUnknown, ErrOwned, ErrNoCredit, ErrLocked} {
		rec := httptest.NewRecorder()
		writeErr(rec, httptest.NewRequest(http.MethodGet, "/", nil), fmt.Errorf("ctx: %w", e))
		require.Less(t, rec.Code, 500, e.Error())
	}
}
