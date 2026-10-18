//go:build integration

package integration

import (
	"net/http"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/require"

	"github.com/otter/server/internal/progression"
)

// Content names over HTTP (signed in, like the other catalogs), and every
// code seeded in the database (not in the YAML the loader validates) has a
// name in every language (DD-230).
func TestHTTPContentStrings(t *testing.T) {
	w := newWorld(t)
	p := newHTTPPlayer(t, "warrior")
	tr := p.ok("GET", "/v1/content/strings?lang=tr", nil)
	require.Equal(t, "tr", tr["lang"])
	trs := tr["strings"].(map[string]any)
	require.Equal(t, "Demir Kılıç", trs["item.sword_iron"])
	en := p.ok("GET", "/v1/content/strings?lang=xx", nil)
	require.Equal(t, "en", en["lang"], "unknown language: English")
	ens := en["strings"].(map[string]any)
	require.Equal(t, "Iron Sword", ens["item.sword_iron"])

	req, _ := http.NewRequest(http.MethodGet, p.base+"/v1/content/strings?lang=tr", nil)
	req.Header.Set("Authorization", "Bearer "+p.token)
	req.Header.Set("If-None-Match", `"`+tr["version"].(string)+`-tr"`)
	res, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	res.Body.Close()
	require.Equal(t, http.StatusNotModified, res.StatusCode)

	seeded := map[string]string{
		"achievement": `SELECT code FROM content.achievement_template`,
		"category":    `SELECT DISTINCT category FROM content.achievement_template`,
		"title":       `SELECT code FROM content.title_template`,
		"faction":     `SELECT code FROM content.faction`,
		"track":       `SELECT code FROM content.mastery_track`,
		"worldevent":  `SELECT code FROM content.world_event_template WHERE code NOT LIKE 'it\_%'`,
		"seasonal":    `SELECT code FROM content.seasonal_event WHERE code NOT LIKE 'it\_%'`,
	}
	for kind, q := range seeded {
		rows, err := w.p.Query(w.ctx, q)
		require.NoError(t, err)
		codes, err := pgx.CollectRows(rows, pgx.RowTo[string]) // closes rows before any assertion
		require.NoError(t, err)
		for _, code := range codes {
			for _, s := range []map[string]any{trs, ens} {
				require.NotEmpty(t, s[kind+"."+code], "%s.%s", kind, code)
			}
		}
	}
	for _, tier := range progression.RepTierNames {
		require.NotEmpty(t, trs["reptier."+tier], tier)
	}
}
