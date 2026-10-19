//go:build integration

package integration

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

// Food (DD-269): eating heals and feeds — attack and defense rise for an
// hour, even at full HP; a potion at full HP is still refused.
func TestHTTPFoodBuff(t *testing.T) {
	w := newWorld(t)
	p := newHTTPPlayer(t, "warrior")
	me := player{ID: uuid.MustParse(p.char)}
	_, err := w.p.Exec(w.ctx, `UPDATE character.character SET level = 100, hp_current = NULL WHERE id = $1`, me.ID)
	require.NoError(t, err)
	stat := func(key string) float64 {
		return p.ok("GET", p.path("/sheet"), nil)["derived"].(map[string]any)[key].(float64)
	}
	attack, defense := stat("attack"), stat("defense")
	require.Nil(t, p.ok("GET", p.path("/buff"), nil)["buff"])

	w.give(me, "potion_hp_small", 1)
	p.fail("POST", p.path("/items/"+p.bagItem("potion_hp_small")+"/use"), map[string]any{}, 409, "E_HP_FULL")

	w.give(me, "food_leviathan_banquet", 1)
	res := p.ok("POST", p.path("/items/"+p.bagItem("food_leviathan_banquet")+"/use"), map[string]any{})
	buff := res["buff"].(map[string]any)
	require.Equal(t, 9.0, buff["attack_pct"])
	require.Equal(t, 9.0, buff["defense_pct"])
	require.Equal(t, "food_leviathan_banquet", p.ok("GET", p.path("/buff"), nil)["buff"].(map[string]any)["source_template"])
	require.Greater(t, stat("attack"), attack, "fed: more attack")
	require.Greater(t, stat("defense"), defense, "fed: more defense")

	// A new dish replaces the buff.
	w.give(me, "food_grilled_guppy", 1)
	res = p.ok("POST", p.path("/items/"+p.bagItem("food_grilled_guppy")+"/use"), map[string]any{})
	require.Equal(t, 2.0, res["buff"].(map[string]any)["attack_pct"])
}
