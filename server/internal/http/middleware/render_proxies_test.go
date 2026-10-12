package middleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/otter/server/internal/httpctx"
)

// The TRUSTED_PROXIES value in render.yaml resolves the player behind
// Cloudflare and Render's loopback edge, and a client-written
// X-Forwarded-For entry cannot pose as the player.
func TestRenderTrustedProxies(t *testing.T) {
	trusted, err := ParseTrustedProxies(strings.Split("::1/128,127.0.0.0/8,10.0.0.0/8,173.245.48.0/20,103.21.244.0/22,103.22.200.0/22,103.31.4.0/22,141.101.64.0/18,108.162.192.0/18,190.93.240.0/20,188.114.96.0/20,197.234.240.0/22,198.41.128.0/17,162.158.0.0/15,104.16.0.0/13,104.24.0.0/14,172.64.0.0/13,131.0.72.0/22,2400:cb00::/32,2606:4700::/32,2803:f800::/32,2405:b500::/32,2405:8100::/32,2a06:98c0::/29,2c0f:f248::/32", ","))
	require.NoError(t, err)
	var got string
	h := RealIP(trusted)(http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) { got = httpctx.ClientIP(r) }))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.RemoteAddr = "[::1]:43740"
	req.Header.Set("X-Forwarded-For", "6.6.6.6, 203.0.113.7, 162.158.95.228")
	h.ServeHTTP(httptest.NewRecorder(), req)
	require.Equal(t, "203.0.113.7", got)
}
