//go:build integration

package integration

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// Gathering professions (DD-244) through the real router: the resource list
// with levels, the level gate, a fast-forwarded haul of the chosen resource
// and the gathering XP it pays.
func TestHTTPGatheringProfessions(t *testing.T) {
	w := newWorld(t)
	p := newHTTPPlayer(t, "warrior")

	profs := items(p.ok("GET", p.path("/gathering"), nil), "professions")
	require.Len(t, profs, 5)
	mining := profs[0]
	require.Equal(t, "mining", mining["kind"])
	require.Equal(t, 1.0, mining["level"])
	nodes := mining["nodes"].([]any)
	require.Equal(t, "raw_stone", nodes[0].(map[string]any)["code"], "every profession starts at a level-1 resource")

	p.fail("POST", p.path("/activity"), map[string]any{"kind": "mining", "duration_seconds": 3600, "target": "adamantite_ore"}, 409, "E_GATHER_LEVEL")
	p.fail("POST", p.path("/activity"), map[string]any{"kind": "mining", "duration_seconds": 3600, "target": "oak_log"}, 400, "E_VALIDATION")
	act := p.ok("POST", p.path("/activity"), map[string]any{"kind": "mining", "duration_seconds": 3600, "target": "raw_stone"})
	id := act["contract"].(map[string]any)["id"].(string)
	_, err := w.p.Exec(w.ctx, `
		UPDATE activity.contract SET started_at = now() - interval '61 minutes', expires_at = now() - interval '1 minute'
		 WHERE id = $1`, id)
	require.NoError(t, err)
	res := p.ok("POST", "/v1/activity/"+id+"/reveal", map[string]any{})
	got := res["yield"].(map[string]any)["items"].([]any)
	require.NotEmpty(t, got, "%v", res)
	for _, it := range got {
		require.Equal(t, "raw_stone", it.(map[string]any)["template_code"])
	}
	gather := res["gather"].(map[string]any)
	require.Equal(t, float64(len(got)*8), gather["gained"], "tier 0 pays 8 XP a unit")
	require.True(t, gather["leveled_up"].(bool))

	again := p.ok("POST", "/v1/activity/"+id+"/reveal", map[string]any{})
	require.Equal(t, true, again["replayed"])
	after := items(p.ok("GET", p.path("/gathering"), nil), "professions")[0]
	require.Equal(t, gather["level"], after["level"], "a replay grants no XP")

	// Recipes describe what they make (workshop item details).
	var helm map[string]any
	for _, r := range items(p.ok("GET", "/v1/recipes?profession=leatherworking", nil), "recipes") { // leather armor (DD-274)
		if r["code"] == "recipe_helm_leather" {
			helm = r["output"].(map[string]any)
		}
	}
	require.NotNil(t, helm)
	require.Equal(t, "helm", helm["equip_slot"])
	require.NotEmpty(t, helm["stat_bonuses"])

	// Gathering professions are not learned by hand.
	p.fail("POST", p.path("/professions/learn"), map[string]any{"profession_code": "mining"}, 404, "E_UNKNOWN_PROFESSION")
}
