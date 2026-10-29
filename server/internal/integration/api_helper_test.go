//go:build integration

package integration

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

// Hired helpers (DD-319): bought per character, working beside the
// character at half speed and XP, each with its own bonus.
func TestHTTPHelpers(t *testing.T) {
	w := newWorld(t)
	p := newHTTPPlayer(t, "warrior")
	me := player{ID: uuid.MustParse(p.char)}
	p.mint(500_000)
	list := p.ok("GET", p.path("/helpers"), nil)["helpers"].([]any)
	require.Len(t, list, 14)
	p.fail("POST", p.path("/helpers/miner/start"), map[string]any{"duration_seconds": 3600}, 409, "E_HELPER_NOT_HIRED")

	gold := p.gold()
	p.ok("POST", p.path("/helpers/miner/hire"), map[string]any{})
	require.Equal(t, gold-50_000, p.gold(), "5 gold")
	p.fail("POST", p.path("/helpers/miner/hire"), map[string]any{}, 409, "E_HELPER_HIRED")

	// A gathering run beside the character's own work.
	p.ok("POST", p.path("/helpers/miner/start"), map[string]any{"duration_seconds": 3600})
	p.ok("POST", p.path("/activity"), map[string]any{"kind": "fishing", "duration_seconds": 3600})
	p.fail("POST", p.path("/helpers/miner/start"), map[string]any{"duration_seconds": 3600}, 409, "E_HELPER_BUSY")
	_, err := w.p.Exec(w.ctx, `UPDATE helper.job SET started_at = now() - interval '61 minutes', ready_at = now() - interval '1 minute'
		WHERE character_id = $1`, p.char)
	require.NoError(t, err)
	res := p.ok("POST", p.path("/helpers/miner/collect"), map[string]any{})["result"].(map[string]any)
	require.Equal(t, "mining", res["profession"])
	require.NotEmpty(t, res["items"], "an hour of half-speed mining brings ore")
	require.Positive(t, res["xp"])
	p.fail("POST", p.path("/helpers/miner/collect"), map[string]any{}, 409, "E_HELPER_IDLE")

	// A smelting batch: materials and gold now, the bars when done.
	p.ok("POST", p.path("/helpers/smelter/hire"), map[string]any{})
	_, err = w.p.Exec(w.ctx, `INSERT INTO crafting.character_profession (character_id, profession_code, level, xp)
		VALUES ($1, 'smelting', 1, 0) ON CONFLICT DO NOTHING`, p.char)
	require.NoError(t, err)
	var recipe, ore, fuel string
	var oreQty, fuelQty int32
	require.NoError(t, w.p.QueryRow(w.ctx, `
		SELECT code, inputs->0->>'template_code', (inputs->0->>'quantity')::int, inputs->1->>'template_code',
		       (inputs->1->>'quantity')::int
		  FROM content.recipe WHERE profession_code = 'smelting' AND min_skill = 1 ORDER BY code LIMIT 1`).
		Scan(&recipe, &ore, &oreQty, &fuel, &fuelQty))
	p.fail("POST", p.path("/helpers/smelter/start"), map[string]any{"recipe_code": recipe, "count": 2}, 409, "E_INSUFFICIENT_MATERIALS")
	w.give(me, ore, oreQty*2)
	w.give(me, fuel, fuelQty*2)
	p.ok("POST", p.path("/helpers/smelter/start"), map[string]any{"recipe_code": recipe, "count": 2})
	require.Zero(t, w.count(me, ore), "materials are taken at start")
	p.fail("POST", p.path("/helpers/smelter/collect"), map[string]any{}, 409, "E_HELPER_NOT_READY")
	_, err = w.p.Exec(w.ctx, `UPDATE helper.job SET ready_at = now() - interval '1 second' WHERE character_id = $1
		AND type_code = 'smelter'`, p.char)
	require.NoError(t, err)
	res = p.ok("POST", p.path("/helpers/smelter/collect"), map[string]any{})["result"].(map[string]any)
	require.EqualValues(t, 2, res["made"])

	// Cancelling a craft gives everything back.
	w.give(me, ore, oreQty)
	w.give(me, fuel, fuelQty)
	p.ok("POST", p.path("/helpers/smelter/start"), map[string]any{"recipe_code": recipe, "count": 1})
	p.ok("POST", p.path("/helpers/smelter/cancel"), map[string]any{})
	require.Equal(t, int(oreQty), w.count(me, ore))
}
