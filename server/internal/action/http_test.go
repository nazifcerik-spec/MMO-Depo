package action

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWrite(t *testing.T) {
	rec := httptest.NewRecorder()
	require.True(t, Write(rec, httptest.NewRequest(http.MethodPost, "/", nil), &BusyError{Kind: "dungeon"}))
	require.Equal(t, http.StatusConflict, rec.Code)
	var body struct {
		Error struct{ Code, Field string }
	}
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &body))
	require.Equal(t, "E_CHARACTER_BUSY", body.Error.Code)
	require.Equal(t, "dungeon", body.Error.Field)

	rec = httptest.NewRecorder()
	require.False(t, Write(rec, httptest.NewRequest(http.MethodPost, "/", nil), errors.New("other")))
	require.Equal(t, 200, rec.Code, "nothing written")
}
