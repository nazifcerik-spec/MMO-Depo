//go:build integration

package integration

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

// DD-360: the panel page is served without a session, the data behind it is not.
func TestHTTPAdminPanel(t *testing.T) {
	w := newWorld(t)
	p := newHTTPPlayer(t, "warrior")

	resp, err := http.Get(p.base + "/internal/admin/panel")
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.Contains(t, resp.Header.Get("Content-Security-Policy"), "default-src 'none'")

	anon := &hc{t: t, base: p.base}
	anon.fail("GET", "/internal/ops/overview", nil, 401, "")
	anon.fail("GET", "/internal/ops/market/listings", nil, 401, "")
	p.fail("GET", "/internal/ops/overview", nil, 401, "") // a player's token is not an operator's

	viewer := opsLogin(t, w, p.base, "viewer")
	o := viewer.ok("GET", "/internal/ops/overview", nil)
	require.GreaterOrEqual(t, o["accounts"].(float64), float64(1))
	require.GreaterOrEqual(t, o["characters"].(float64), float64(1))
	viewer.fail("GET", "/internal/ops/market/listings", nil, 403, "")

	agent := opsLogin(t, w, p.base, "support_l2")
	require.NotNil(t, agent.ok("GET", "/internal/ops/market/listings", nil))
}
