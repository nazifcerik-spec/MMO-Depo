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
	// Gearless: a margin over the level 6 dungeon so the run reaches the boss (DD-306).
	_, err := w.p.Exec(w.ctx, `UPDATE character.character SET level = 12 WHERE id = $1`, p.char)
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

	// Reward history (DD-242): one entry per collected contract (replays add
	// none), stored read, in its own list apart from the regular inbox.
	rewards := items(p.ok("GET", "/v1/notifications?rewards=1", nil), "notifications")
	require.Len(t, rewards, 2, "hunt + dungeon: %v", run)
	for _, n := range rewards {
		require.Equal(t, "reward_collected", n["kind"])
		require.NotNil(t, n["read_at"], "history must not raise the unread badge")
	}
	huntVars := rewards[1]["body"].(map[string]any)["body_vars"].(map[string]any)
	require.Equal(t, "hunt", huntVars["source"])
	require.Equal(t, "young_boar", huntVars["target"])
	require.Equal(t, hunt["total_xp"], huntVars["xp"])
	stacked := 0.0
	for _, it := range huntVars["items"].([]any) {
		stacked += it.(map[string]any)["count"].(float64)
	}
	drops, _ := res["items"].([]any)
	require.Equal(t, float64(len(drops)), stacked, "every drop is in a stack")
	require.Equal(t, "dungeon", rewards[0]["body"].(map[string]any)["body_vars"].(map[string]any)["source"])
	for _, n := range items(p.ok("GET", "/v1/notifications", nil), "notifications") {
		require.NotEqual(t, "reward_collected", n["kind"], "rewards stay out of the regular inbox")
	}

	// Deleting an entry (DD-245): own only; a foreign or repeated delete is a no-op.
	id := rewards[0]["id"].(string)
	other := newHTTPPlayer(t, "mage")
	require.Equal(t, 0.0, other.ok("DELETE", "/v1/notifications/"+id, nil)["deleted"])
	require.Equal(t, 1.0, p.ok("DELETE", "/v1/notifications/"+id, nil)["deleted"])
	require.Equal(t, 0.0, p.ok("DELETE", "/v1/notifications/"+id, nil)["deleted"])
	require.Len(t, items(p.ok("GET", "/v1/notifications?rewards=1", nil), "notifications"), 1)
	p.fail("DELETE", "/v1/notifications/nope", nil, 400, "E_INVALID_ID")
}
