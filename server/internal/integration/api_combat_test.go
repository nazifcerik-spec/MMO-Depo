//go:build integration

package integration

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// Fast-forwarded AFK contracts through the real router: a hunt with real
// encounters (rewards, progression and quest ticks) and a full dungeon run.
func TestHTTPCombatRewind(t *testing.T) {
	w := newWorld(t)
	p := newHTTPPlayer(t, "warrior")
	_, err := w.p.Exec(w.ctx, `UPDATE character.character SET level = 3 WHERE id = $1`, p.char)
	require.NoError(t, err)

	p.ok("POST", p.path("/quests/main_01_first_hunt/accept"), map[string]any{})
	h := p.ok("POST", p.path("/hunt"), map[string]any{"mob_code": "young_boar", "duration_seconds": 600})
	hid := h["contract"].(map[string]any)["id"].(string)
	_, err = w.p.Exec(w.ctx, `
		UPDATE hunt.contract SET started_at = now() - interval '11 minutes', expires_at = now() - interval '1 minute'
		 WHERE id = $1`, hid)
	require.NoError(t, err)
	res := p.ok("POST", "/v1/hunt/"+hid+"/reveal", map[string]any{})
	hunt := res["hunt"].(map[string]any)
	require.Positive(t, hunt["victories"].(float64), "%v", res)
	again := p.ok("POST", "/v1/hunt/"+hid+"/reveal", map[string]any{})
	require.Equal(t, hunt["victories"], again["hunt"].(map[string]any)["victories"], "reveal replays")
	q := items(p.ok("GET", p.path("/quests"), nil), "open")
	require.Positive(t, q[0]["objectives"].([]any)[0].(map[string]any)["progress"].(float64), "kills ticked the quest")
	ach := items(p.ok("GET", p.path("/achievements"), nil), "achievements")
	unlocked := false
	for _, a := range ach {
		if a["code"] == "first_blood" && a["unlocked_at"] != nil {
			unlocked = true
		}
	}
	require.True(t, unlocked)

	// Full dungeon run; a cleared run locks the dungeon for the week.
	d := p.ok("POST", p.path("/dungeon"), map[string]any{"dungeon_code": "crypt_of_ashes"})
	did := d["contract"].(map[string]any)["id"].(string)
	p.fail("POST", "/v1/dungeon/"+did+"/reveal", map[string]any{}, 409, "")
	p.ok("POST", "/v1/dungeon/"+did+"/stop", map[string]any{})
	d = p.ok("POST", p.path("/dungeon"), map[string]any{"dungeon_code": "crypt_of_ashes"})
	did = d["contract"].(map[string]any)["id"].(string)
	_, err = w.p.Exec(w.ctx, `
		UPDATE dungeon.contract SET started_at = now() - interval '2 hours', expires_at = now() - interval '1 hour'
		 WHERE id = $1`, did)
	require.NoError(t, err)
	run := p.ok("POST", "/v1/dungeon/"+did+"/reveal", map[string]any{})
	r := run["run"].(map[string]any)
	p.ok("POST", "/v1/dungeon/"+did+"/reveal", map[string]any{}) // replay
	if r["cleared"] == true {
		p.fail("POST", p.path("/dungeon"), map[string]any{"dungeon_code": "crypt_of_ashes"}, 409, "")
	}
	hist := items(p.ok("GET", p.path("/dungeon/history"), nil), "runs")
	require.NotEmpty(t, hist)
}
