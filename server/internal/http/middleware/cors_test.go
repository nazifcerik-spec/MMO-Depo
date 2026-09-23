package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	appmw "github.com/otter/server/internal/http/middleware"
)

func TestCORS(t *testing.T) {
	var hits int
	h := appmw.CORS([]string{"https://play.otter.gg/"}, true)(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		hits++
		w.WriteHeader(200)
	}))
	do := func(method, origin string, preflight bool) *httptest.ResponseRecorder {
		req := httptest.NewRequest(method, "/v1/character", nil)
		if origin != "" {
			req.Header.Set("Origin", origin)
		}
		if preflight {
			req.Header.Set("Access-Control-Request-Method", "POST")
		}
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, req)
		return rec
	}

	rec := do("OPTIONS", "https://play.otter.gg", true)
	require.Equal(t, 204, rec.Code)
	require.Equal(t, "https://play.otter.gg", rec.Header().Get("Access-Control-Allow-Origin"))
	require.Contains(t, rec.Header().Get("Access-Control-Allow-Headers"), "Idempotency-Key")
	require.Zero(t, hits, "preflight never reaches handlers")

	rec = do("GET", "http://localhost:53172", false)
	require.Equal(t, "http://localhost:53172", rec.Header().Get("Access-Control-Allow-Origin"))

	rec = do("OPTIONS", "https://evil.example", true)
	require.Equal(t, 403, rec.Code)
	require.Empty(t, rec.Header().Get("Access-Control-Allow-Origin"))

	rec = do("GET", "https://localhost.evil.example", false)
	require.Empty(t, rec.Header().Get("Access-Control-Allow-Origin"))

	rec = do("GET", "", false) // native client
	require.Equal(t, 200, rec.Code)
	require.Empty(t, rec.Header().Get("Access-Control-Allow-Origin"))

	strict := appmw.CORS(nil, false)(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {}))
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Origin", "http://localhost:3000")
	rec = httptest.NewRecorder()
	strict.ServeHTTP(rec, req)
	require.Empty(t, rec.Header().Get("Access-Control-Allow-Origin"), "localhost only in dev")
}
