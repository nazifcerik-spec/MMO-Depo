package story

import (
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidateAndInstall(t *testing.T) {
	require.NoError(t, Install(nil))
	require.Empty(t, All())
	ok := []Chapter{{Code: "c1", MinLevel: 1}, {Code: "c2", MinLevel: 10, Requires: "c1"}}
	require.NoError(t, Install(ok))
	require.Equal(t, ok, All())
	for name, chs := range map[string][]Chapter{
		"no code":   {{MinLevel: 1}},
		"duplicate": {{Code: "a", MinLevel: 1}, {Code: "a", MinLevel: 2}},
		"level":     {{Code: "a", MinLevel: 0}},
		"requires":  {{Code: "a", MinLevel: 1, Requires: "b"}, {Code: "b", MinLevel: 1}},
	} {
		require.Error(t, Install(chs), name)
	}
	rec := httptest.NewRecorder()
	CatalogHandler(rec, httptest.NewRequest("GET", "/v1/story", nil))
	require.Contains(t, rec.Body.String(), `"code":"c2"`)
}
