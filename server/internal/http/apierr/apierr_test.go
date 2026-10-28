package apierr

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWrite_HidesServerErrorDetail(t *testing.T) {
	rec := httptest.NewRecorder()
	Write(rec, httptest.NewRequest("GET", "/", nil), http.StatusInternalServerError, "E_INTERNAL",
		errors.New(`pq: relation "ledger.entry" does not exist`).Error())
	require.Equal(t, 500, rec.Code)
	require.NotContains(t, rec.Body.String(), "ledger.entry")
	require.Contains(t, rec.Body.String(), `"message":"internal error"`)

	rec = httptest.NewRecorder()
	Write(rec, httptest.NewRequest("GET", "/", nil), http.StatusConflict, "E_BAG_FULL", "bag is full")
	require.Contains(t, rec.Body.String(), "bag is full", "4xx messages are player-facing and kept")
}

func TestCodeOf(t *testing.T) {
	r := httptest.NewRequest("GET", "/", nil)
	require.Equal(t, "E_BAG_FULL", CodeOf(r, func(w http.ResponseWriter, r *http.Request) {
		Write(w, r, http.StatusConflict, "E_BAG_FULL", "bag is full")
	}))
	require.Equal(t, "E_VALIDATION", CodeOf(r, func(w http.ResponseWriter, r *http.Request) {
		WriteField(w, r, http.StatusBadRequest, "E_VALIDATION", "bad", "target")
	}))
}
