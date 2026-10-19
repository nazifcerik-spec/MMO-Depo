//go:build integration

package integration

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// Hunting is regional (DD-273): from the starting city the road to
// Greenwood is in reach, Emberpeak is not until the character stands in it
// or in a region linked to it. Every huntable mob has a region.
func TestHTTPHuntRegionGate(t *testing.T) {
	w := newWorld(t)
	p := newHTTPPlayer(t, "warrior")
	loc := p.ok("GET", p.path("/location"), nil)
	require.Equal(t, "ironhold_city", loc["region"].(map[string]any)["code"])
	require.ElementsMatch(t, []any{"ironhold_city", "greenwood"}, loc["hunt_regions"])

	mobs := items(p.ok("GET", "/v1/mobs", nil), "mobs")
	for _, m := range mobs {
		require.NotEmpty(t, m["region_code"], "%s has a region", m["code"])
	}

	hunt := func(mob string) (int, map[string]any) {
		return p.call("POST", p.path("/hunt"), map[string]any{"mob_code": mob, "duration_seconds": 600})
	}
	p.fail("POST", p.path("/hunt"), map[string]any{"mob_code": "rock_ogre", "duration_seconds": 600}, 409, "E_MOB_REGION")
	code, out := hunt("young_boar")
	require.Equal(t, 201, code, "%v", out)
	p.ok("POST", "/v1/hunt/"+out["contract"].(map[string]any)["id"].(string)+"/stop", map[string]any{})

	// Whispering Reach is linked to Emberpeak by road.
	_, err := w.p.Exec(w.ctx, `UPDATE character.character SET location_region_code = 'whispering_reach' WHERE id = $1`, p.char)
	require.NoError(t, err)
	code, out = hunt("rock_ogre")
	require.Equal(t, 201, code, "%v", out)
	p.ok("POST", "/v1/hunt/"+out["contract"].(map[string]any)["id"].(string)+"/stop", map[string]any{})
	// Greenwood is linked too; the Misty Marshes are two roads away.
	code, out = hunt("young_boar")
	require.Equal(t, 201, code, "%v", out)
	p.ok("POST", "/v1/hunt/"+out["contract"].(map[string]any)["id"].(string)+"/stop", map[string]any{})
	p.fail("POST", p.path("/hunt"), map[string]any{"mob_code": "mist_leech", "duration_seconds": 600}, 409, "E_MOB_REGION")
}
