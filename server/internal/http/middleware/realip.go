package middleware

import (
	"fmt"
	"net"
	"net/http"
	"net/netip"
	"strings"

	"github.com/otter/server/internal/httpctx"
)

// ParseTrustedProxies parses TRUSTED_PROXIES entries (CIDRs or single IPs).
func ParseTrustedProxies(entries []string) ([]netip.Prefix, error) {
	out := make([]netip.Prefix, 0, len(entries))
	for _, e := range entries {
		e = strings.TrimSpace(e)
		if e == "" {
			continue
		}
		if p, err := netip.ParsePrefix(e); err == nil {
			out = append(out, p.Masked())
			continue
		}
		a, err := netip.ParseAddr(e)
		if err != nil {
			return nil, fmt.Errorf("trusted proxy %q: not an IP or CIDR", e)
		}
		out = append(out, netip.PrefixFrom(a.Unmap(), a.Unmap().BitLen()))
	}
	return out, nil
}

// RealIP resolves the caller IP once per request (docs/backend/middleware.md
// §3.2). Forwarded headers are client-writable, so they are only honoured when
// the peer is a trusted proxy; X-Forwarded-For is then read right to left,
// skipping trusted hops, and the first untrusted address is the client. With
// no trusted proxies the peer address is always used — otherwise anyone could
// pick a fresh "IP" per request and walk past the per-IP rate limits.
func RealIP(trusted []netip.Prefix) func(http.Handler) http.Handler {
	isTrusted := func(a netip.Addr) bool {
		a = a.Unmap()
		for _, p := range trusted {
			if p.Contains(a) {
				return true
			}
		}
		return false
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := peer(r)
			if a, err := netip.ParseAddr(ip); err == nil && isTrusted(a) {
				ip = forwarded(r, ip, isTrusted)
			}
			next.ServeHTTP(w, r.WithContext(httpctx.WithClientIP(r.Context(), ip)))
		})
	}
}

func peer(r *http.Request) string {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}

func forwarded(r *http.Request, fallback string, isTrusted func(netip.Addr) bool) string {
	hops := strings.Split(strings.Join(r.Header.Values("X-Forwarded-For"), ","), ",")
	for i := len(hops) - 1; i >= 0; i-- {
		a, err := netip.ParseAddr(strings.TrimSpace(hops[i]))
		if err != nil {
			break // garbage from the client side: stop at the last good hop
		}
		if !isTrusted(a) {
			return a.Unmap().String()
		}
		fallback = a.Unmap().String()
	}
	if strings.TrimSpace(r.Header.Get("X-Forwarded-For")) == "" {
		if a, err := netip.ParseAddr(strings.TrimSpace(r.Header.Get("X-Real-IP"))); err == nil {
			return a.Unmap().String()
		}
	}
	return fallback
}
