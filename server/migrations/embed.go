// Package migrations embeds the SQL migrations so the API can apply them at
// start where the host has no separate migrate step (Render free, DD-239).
package migrations

import "embed"

//go:embed *.sql
var FS embed.FS
