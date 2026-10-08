package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	appmw "github.com/otter/server/internal/http/middleware"
	"github.com/otter/server/internal/httpctx"
)

func resolve(t *testing.T, trusted []string, remote string, headers map[string]string) string {
	t.Helper()
	prefixes, err := appmw.ParseTrustedProxies(trusted)
	require.NoError(t, err)
	var got string
	h := appmw.RealIP(prefixes)(http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
		got = httpctx.ClientIP(r)
	}))
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.RemoteAddr = remote
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	h.ServeHTTP(httptest.NewRecorder(), req)
	return got
}

func TestRealIP(t *testing.T) {
	lb := []string{"10.0.0.0/8"}
	for name, c := range map[string]struct {
		trusted []string
		remote  string
		headers map[string]string
		want    string
	}{
		"direct: headers ignored": {nil, "198.51.100.7:5000",
			map[string]string{"X-Forwarded-For": "1.2.3.4", "X-Real-IP": "1.2.3.4"}, "198.51.100.7"},
		"untrusted peer: headers ignored": {lb, "198.51.100.7:5000",
			map[string]string{"X-Forwarded-For": "1.2.3.4"}, "198.51.100.7"},
		"trusted LB: client appended by LB": {lb, "10.0.0.2:5000",
			map[string]string{"X-Forwarded-For": "203.0.113.9"}, "203.0.113.9"},
		"spoofed prefix is skipped": {lb, "10.0.0.2:5000",
			map[string]string{"X-Forwarded-For": "1.2.3.4, 203.0.113.9"}, "203.0.113.9"},
		"trusted hops skipped": {lb, "10.0.0.2:5000",
			map[string]string{"X-Forwarded-For": "203.0.113.9, 10.1.1.1"}, "203.0.113.9"},
		"garbage stops the walk": {lb, "10.0.0.2:5000",
			map[string]string{"X-Forwarded-For": "203.0.113.9, junk, 10.1.1.1"}, "10.1.1.1"},
		"real-ip from trusted peer": {lb, "10.0.0.2:5000",
			map[string]string{"X-Real-IP": "203.0.113.5"}, "203.0.113.5"},
		"single trusted IP entry": {[]string{"10.0.0.2"}, "10.0.0.2:5000",
			map[string]string{"X-Forwarded-For": "203.0.113.9"}, "203.0.113.9"},
		"ipv6 peer": {nil, "[2001:db8::1]:5000", nil, "2001:db8::1"},
	} {
		t.Run(name, func(t *testing.T) {
			require.Equal(t, c.want, resolve(t, c.trusted, c.remote, c.headers))
		})
	}
}

func TestParseTrustedProxies_RejectsJunk(t *testing.T) {
	_, err := appmw.ParseTrustedProxies([]string{"10.0.0.0/8", "lb.internal"})
	require.Error(t, err)
	p, err := appmw.ParseTrustedProxies([]string{" ", "192.168.1.0/24"})
	require.NoError(t, err)
	require.Len(t, p, 1)
}

// Without the middleware (or before it ran) the peer address is used.
func TestClientIP_DefaultsToPeer(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.RemoteAddr = "10.0.0.5:54321"
	req.Header.Set("X-Forwarded-For", "1.2.3.4")
	require.Equal(t, "10.0.0.5", httpctx.ClientIP(req))
}
