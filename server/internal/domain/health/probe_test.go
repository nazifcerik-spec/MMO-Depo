package health

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestProbe(t *testing.T) {
	ok := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {}))
	defer ok.Close()
	bad := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
	}))
	defer bad.Close()
	require.Equal(t, 0, Probe(ok.URL))
	require.Equal(t, 1, Probe(bad.URL))
	require.Equal(t, 1, Probe("http://127.0.0.1:1/live"), "nothing listening")
}

func TestLocalURL(t *testing.T) {
	require.Equal(t, "http://127.0.0.1:8080/live", LocalURL(":8080", ":1", "/live"))
	require.Equal(t, "http://127.0.0.1:9091/metrics", LocalURL("0.0.0.0:9091", ":1", "/metrics"))
	require.Equal(t, "http://127.0.0.1:8080/live", LocalURL("", ":8080", "/live"))
}
