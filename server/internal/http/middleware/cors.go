package middleware

import (
	"net/http"
	"net/url"
	"strings"
)

// CORS lets browser clients (Flutter web) call the API. Only listed origins
// are echoed back — never "*", because requests carry bearer tokens. In dev,
// any http://localhost:<port> / 127.0.0.1 origin is accepted as well since
// `flutter run -d chrome` picks a random port. Native apps don't send Origin
// and are unaffected.
func CORS(allowed []string, allowLocalhost bool) func(http.Handler) http.Handler {
	set := make(map[string]bool, len(allowed))
	for _, o := range allowed {
		if o = strings.TrimRight(strings.TrimSpace(o), "/"); o != "" {
			set[o] = true
		}
	}
	ok := func(origin string) bool {
		if set[origin] {
			return true
		}
		if !allowLocalhost {
			return false
		}
		u, err := url.Parse(origin)
		if err != nil || u.Scheme != "http" {
			return false
		}
		h := u.Hostname()
		return h == "localhost" || h == "127.0.0.1"
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")
			if origin == "" {
				next.ServeHTTP(w, r)
				return
			}
			w.Header().Add("Vary", "Origin")
			if !ok(origin) {
				if r.Method == http.MethodOptions {
					w.WriteHeader(http.StatusForbidden)
					return
				}
				next.ServeHTTP(w, r) // browser blocks the response itself
				return
			}
			h := w.Header()
			h.Set("Access-Control-Allow-Origin", origin)
			h.Set("Access-Control-Expose-Headers", "X-Request-Id, traceparent, Retry-After, X-RateLimit-Remaining")
			if r.Method == http.MethodOptions && r.Header.Get("Access-Control-Request-Method") != "" {
				h.Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE")
				h.Set("Access-Control-Allow-Headers",
					"Authorization, Content-Type, Idempotency-Key, X-Request-Id, X-Character-Id, traceparent")
				h.Set("Access-Control-Max-Age", "600")
				w.WriteHeader(http.StatusNoContent)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
