//go:build integration

package integration

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/otter/server/internal/worker"
)

// Client error reports (DD-235): anonymous intake grouped by fingerprint,
// secrets masked, read by operators with ops.client_errors.
func TestHTTPClientErrors(t *testing.T) {
	w := newWorld(t)
	p := newHTTPPlayer(t, "warrior")
	anon := &hc{t: t, base: p.base}

	tag := randHex(6)
	report := func(line int, version string) map[string]any {
		return map[string]any{"kind": "flutter", "platform": "android", "app_version": version, "flavor": "prod",
			"message": "Null check operator used on a null value " + tag + " for mail bob@example.com",
			"stack": "#0      BagGrid.build (package:otter/bag_" + tag + ".dart:" + strings.Repeat("1", line) + ":5)\n" +
				"#1      StatelessElement.build (package:flutter/src/widgets/framework.dart:5:1)"}
	}
	anon.ok("POST", "/v1/client/errors", report(2, "1.0.0+1"))
	anon.ok("POST", "/v1/client/errors", report(3, "1.0.1+2")) // same bug, new build → same group
	anon.fail("POST", "/v1/client/errors", map[string]any{"kind": "flutter", "message": " "}, 400, "E_VALIDATION")
	anon.fail("POST", "/v1/client/errors", "{", 400, "E_INVALID_JSON")

	opsLogin(t, w, p.base, "support_l1").fail("GET", "/internal/ops/client-errors", nil, 403, "")
	ops := opsLogin(t, w, p.base, "ops_engineer")
	var mine []map[string]any
	for _, g := range items(ops.ok("GET", "/internal/ops/client-errors?platform=android&limit=200", nil), "errors") {
		if strings.Contains(g["message"].(string), tag) {
			mine = append(mine, g)
		}
	}
	require.Len(t, mine, 1, "both reports share one fingerprint")
	require.Equal(t, 2.0, mine[0]["count"])
	require.Equal(t, "1.0.1+2", mine[0]["app_version"], "latest build is kept")
	require.NotContains(t, mine[0]["message"], "bob@example.com", "e-mails are masked")
	for _, g := range items(ops.ok("GET", "/internal/ops/client-errors?platform=ios&limit=200", nil), "errors") {
		require.NotContains(t, g["message"], tag, "platform filter")
	}

	// GC: groups unseen for 30 days are dropped.
	_, err := w.p.Exec(w.ctx, `UPDATE ops.client_error SET last_seen = now() - INTERVAL '31 days' WHERE message LIKE '%' || $1 || '%'`, tag)
	require.NoError(t, err)
	_, err = worker.ClientErrorGC(w.p).Run(w.ctx)
	require.NoError(t, err)
	var left int
	require.NoError(t, w.p.QueryRow(w.ctx, `SELECT count(*) FROM ops.client_error WHERE message LIKE '%' || $1 || '%'`, tag).Scan(&left))
	require.Zero(t, left)
}
