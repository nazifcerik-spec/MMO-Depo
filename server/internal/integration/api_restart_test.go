//go:build integration

package integration

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// Restart (DD-310): collect what a gathering run or a hunt has so far and
// start the same one again for the same planned time.
func TestHTTPRestartGatheringAndHunt(t *testing.T) {
	w := newWorld(t)
	p := newHTTPPlayer(t, "warrior")

	run := p.ok("POST", p.path("/activity"), map[string]any{"kind": "fishing", "duration_seconds": 3600})
	id := run["contract"].(map[string]any)["id"].(string)
	_, err := w.p.Exec(w.ctx, `UPDATE activity.contract SET started_at = now() - interval '30 minutes',
		expires_at = now() + interval '30 minutes' WHERE id = $1`, id)
	require.NoError(t, err)
	res := p.ok("POST", "/v1/activity/"+id+"/restart", map[string]any{})
	require.NotNil(t, res["yield"], "the half hour so far is collected")
	next := res["contract"].(map[string]any)
	require.NotEqual(t, id, next["id"])
	require.Equal(t, "fishing", next["kind"])
	require.EqualValues(t, 3600, next["duration_seconds"], "the same planned time again")
	require.Nil(t, res["restart_error"])
	// The old one is collected: no second restart of it.
	p.fail("POST", "/v1/activity/"+id+"/restart", map[string]any{}, 409, "E_NOT_ACTIVE")
	p.ok("POST", "/v1/activity/"+next["id"].(string)+"/stop", map[string]any{})

	h := p.ok("POST", p.path("/hunt"), map[string]any{"mob_code": "young_boar", "duration_seconds": 3600})
	hid := h["contract"].(map[string]any)["id"].(string)
	_, err = w.p.Exec(w.ctx, `UPDATE hunt.contract SET started_at = now() - interval '20 minutes',
		expires_at = now() + interval '40 minutes' WHERE id = $1`, hid)
	require.NoError(t, err)
	hr := p.ok("POST", "/v1/hunt/"+hid+"/restart", map[string]any{})
	rep := hr["reveal"].(map[string]any)["hunt"].(map[string]any)["report"].(map[string]any)
	require.Equal(t, "collected_early", rep["end"])
	nh := hr["contract"].(map[string]any)
	require.NotEqual(t, hid, nh["id"])
	require.Equal(t, "young_boar", nh["mob_code"])
	require.EqualValues(t, 3600, nh["duration_seconds"])
	p.ok("POST", "/v1/hunt/"+nh["id"].(string)+"/stop", map[string]any{})
}
