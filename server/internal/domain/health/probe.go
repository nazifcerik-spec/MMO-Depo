package health

import (
	"net/http"
	"strings"
	"time"
)

// Probe GETs url and returns a process exit code: 0 when it answers 2xx.
// The images are distroless (no shell, no wget), so their HEALTHCHECK runs
// the binary itself: `/app healthcheck` (DD-232).
func Probe(url string) int {
	c := http.Client{Timeout: 3 * time.Second}
	res, err := c.Get(url)
	if err != nil {
		return 1
	}
	_ = res.Body.Close()
	if res.StatusCode >= 200 && res.StatusCode < 300 {
		return 0
	}
	return 1
}

// LocalURL turns a listen address (":8080", "0.0.0.0:8080", empty = def)
// into a loopback URL for path.
func LocalURL(addr, def, path string) string {
	if addr == "" {
		addr = def
	}
	return "http://127.0.0.1:" + addr[strings.LastIndex(addr, ":")+1:] + path
}
