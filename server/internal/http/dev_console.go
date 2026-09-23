package http

import (
	_ "embed"
	"net/http"
)

// Embedded dev console HTML — served only when ENV=dev.
// A single-page vanilla-JS tool for exercising register → character → wallet →
// mint flows against the local API. Not exposed in staging/prod.

//go:embed dev_console.html
var devConsoleHTML []byte

func devConsoleHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	// Small dev convenience — allow the page to reload without hitting caches.
	w.Header().Set("Cache-Control", "no-store")
	_, _ = w.Write(devConsoleHTML)
}
