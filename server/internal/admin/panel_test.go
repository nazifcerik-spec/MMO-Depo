package admin

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
)

// DD-360: the panel's files are locked down and carry no way to write server
// data as HTML.
func TestPanelFilesAreLockedDown(t *testing.T) {
	r := chi.NewRouter()
	PanelRoutes(r)
	for path, kind := range map[string]string{
		"/internal/admin/panel":     "text/html",
		"/internal/admin/panel.js":  "text/javascript",
		"/internal/admin/panel.css": "text/css",
	} {
		w := httptest.NewRecorder()
		r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, path, nil))
		require.Equal(t, http.StatusOK, w.Code, path)
		require.Contains(t, w.Header().Get("Content-Type"), kind)
		csp := w.Header().Get("Content-Security-Policy")
		require.Contains(t, csp, "default-src 'none'")
		require.Contains(t, csp, "script-src 'self'")
		require.NotContains(t, csp, "unsafe")
		require.Contains(t, csp, "frame-ancestors 'none'")
		require.Equal(t, "no-store", w.Header().Get("Cache-Control"))
		require.Equal(t, "DENY", w.Header().Get("X-Frame-Options"))
	}
	js, err := panelFiles.ReadFile("panel/panel.js")
	require.NoError(t, err)
	for _, banned := range []string{"innerHTML", "outerHTML", "insertAdjacentHTML", "document.write", "eval(", "localStorage"} {
		require.False(t, strings.Contains(string(js), banned), banned)
	}
	html, err := panelFiles.ReadFile("panel/panel.html")
	require.NoError(t, err)
	require.False(t, strings.Contains(string(html), "<script>"), "no inline scripts: the CSP would refuse them")
	require.False(t, strings.Contains(string(html), "style="), "no inline styles")
}
