//go:build integration

package integration

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

// Contract endpoints: empty state, ownership, bad ids, double stop.
func TestHTTPContractEdges(t *testing.T) {
	w := newWorld(t)
	p := newHTTPPlayer(t, "warrior")
	other := newHTTPPlayer(t, "mage")
	_, err := w.p.Exec(w.ctx, `UPDATE character.character SET level = 6 WHERE id = $1`, p.char)
	require.NoError(t, err)

	for _, kind := range []string{"hunt", "dungeon"} {
		require.Nil(t, p.ok("GET", p.path("/"+kind+"/active"), nil)["contract"], kind)
		p.fail("POST", "/v1/"+kind+"/nope/reveal", map[string]any{}, 400, "E_INVALID_ID")
		p.fail("POST", "/v1/"+kind+"/nope/stop", map[string]any{}, 400, "E_INVALID_ID")
		p.fail("POST", "/v1/"+kind+"/"+uuid.NewString()+"/reveal", map[string]any{}, 404, "")
		p.fail("POST", "/v1/"+kind+"/"+uuid.NewString()+"/stop", map[string]any{}, 404, "")
		other.fail("GET", "/v1/character/"+p.char+"/"+kind+"/active", nil, 404, "")
	}

	start := map[string]map[string]any{
		"hunt":    {"mob_code": "young_boar", "duration_seconds": 600},
		"dungeon": {"dungeon_code": "crypt_of_ashes"},
	}
	for _, kind := range []string{"hunt", "dungeon"} {
		id := p.ok("POST", p.path("/"+kind), start[kind])["contract"].(map[string]any)["id"].(string)
		active := p.ok("GET", p.path("/"+kind+"/active"), nil)["contract"].(map[string]any)
		require.Equal(t, id, active["id"], kind)
		require.NotContains(t, active, "seed", "the RNG seed never leaves the server")
		p.fail("POST", p.path("/"+kind), start[kind], 409, "")
		// Someone else's contract is invisible.
		other.fail("POST", "/v1/"+kind+"/"+id+"/stop", map[string]any{}, 404, "")
		other.fail("POST", "/v1/"+kind+"/"+id+"/reveal", map[string]any{}, 404, "")
		stopped := p.ok("POST", "/v1/"+kind+"/"+id+"/stop", map[string]any{})
		require.Equal(t, id, stopped["contract_id"])
		require.Equal(t, "cancelled", p.ok("POST", "/v1/"+kind+"/"+id+"/stop", map[string]any{})["state"], "stop is idempotent")
		p.fail("POST", "/v1/"+kind+"/"+id+"/reveal", map[string]any{}, 409, "")
	}

	// Travel: someone else's trip is invisible too.
	dest := items(p.ok("GET", p.path("/location"), nil), "neighbors")[0]["code"]
	trip := p.ok("POST", p.path("/travel"), map[string]any{"to_region_code": dest})["travel"].(map[string]any)["id"].(string)
	other.fail("POST", "/v1/travel/"+trip+"/cancel", map[string]any{}, 404, "")
	other.fail("POST", "/v1/travel/"+trip+"/complete", map[string]any{}, 404, "")
	p.ok("POST", "/v1/travel/"+trip+"/cancel", map[string]any{})

	p.fail("POST", p.path("/hunt"), map[string]any{"mob_code": "young_boar", "duration_seconds": 1}, 400, "")
	p.fail("POST", p.path("/hunt"), "{", 400, "")
	p.fail("POST", p.path("/dungeon"), map[string]any{"dungeon_code": "nowhere"}, 404, "")
	p.fail("POST", p.path("/dungeon"), "{", 400, "")

	// Crafting learn errors.
	p.fail("POST", p.path("/professions/learn"), map[string]any{"profession_code": "juggling"}, 404, "")
	p.fail("POST", p.path("/professions/learn"), "{", 400, "")
	p.fail("POST", p.path("/craft"), map[string]any{"recipe_code": "recipe_potion_hp_small"}, 409, "")
}
