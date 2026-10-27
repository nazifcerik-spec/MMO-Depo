//go:build integration

package integration

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

// Enchanting through the real router (DD-293): a scroll puts its enchant on
// fitting gear, the bonus reaches the sheet, a new enchant replaces the old
// one, a sale and buyback keep it, and a gathering enchant on worn jewelry
// raises a new run's yield.
func TestHTTPEnchanting(t *testing.T) {
	w := newWorld(t)
	p := newHTTPPlayer(t, "warrior")
	me := player{ID: uuid.MustParse(p.char)}
	_, err := w.p.Exec(w.ctx, `UPDATE character.character SET level = 30 WHERE id = $1`, p.char)
	require.NoError(t, err)
	sword := w.give(me, "sword_training", 1).ID.String()
	helm := w.give(me, "helm_novice", 1).ID.String()
	necklace := w.give(me, "necklace_runed", 1).ID.String()
	potion := w.give(me, "potion_hp_small", 1).ID.String()
	w.give(me, "enchant_weapon_attack_1", 2)
	w.give(me, "enchant_weapon_crit_1", 1)
	w.give(me, "enchant_weapon_attack_2", 1) // level 100
	w.give(me, "enchant_jewelry_mining_1", 1)
	withScroll := func(code string) map[string]any { return map[string]any{"scroll_id": p.bagItem(code)} }
	path := func(item string) string { return p.path("/items/" + item + "/enchant") }
	bagRow := func(id string) map[string]any {
		for _, it := range items(p.ok("GET", p.path("/inventory"), nil), "items") {
			if it["instance_id"] == id {
				return it
			}
		}
		return nil
	}

	// Refusals change nothing.
	p.fail("POST", path(helm), withScroll("enchant_weapon_attack_1"), 409, "E_ENCHANT_SLOT")
	p.fail("POST", path(sword), map[string]any{"scroll_id": potion}, 400, "E_NOT_ENCHANT")
	p.fail("POST", path(sword), withScroll("enchant_weapon_attack_2"), 403, "E_LEVEL_REQUIREMENT")
	p.fail("POST", path(sword), map[string]any{}, 400, "E_INVALID_JSON")
	p.fail("POST", path(uuid.NewString()), withScroll("enchant_weapon_attack_1"), 404, "E_NOT_FOUND")
	require.Equal(t, 2.0, p.stackOf("enchant_weapon_attack_1"))

	// The scroll shows what it does; the enchant joins the sword's bonuses.
	scroll := bagRow(p.bagItem("enchant_weapon_attack_1"))
	require.Equal(t, []any{"weapon"}, scroll["enchant"].(map[string]any)["slots"])
	res := p.ok("POST", path(sword), withScroll("enchant_weapon_attack_1"))
	require.Equal(t, "enchant_weapon_attack_1", res["enchant_code"])
	require.Nil(t, res["replaced"])
	require.Equal(t, 1.0, p.stackOf("enchant_weapon_attack_1"), "one scroll used")
	row := bagRow(sword)
	require.Equal(t, "enchant_weapon_attack_1", row["enchant_code"])
	require.Contains(t, row["effective_bonuses"], map[string]any{"kind": "mult_add", "target": "attack", "amount": 0.01})
	attack := []any{map[string]any{"kind": "mult_add", "target": "attack", "amount": 0.01}}
	require.Equal(t, attack, row["enchant_effect"].(map[string]any)["modifiers"], "what the enchant adds")

	// Worn, a new enchant replaces the old one and reaches the sheet.
	p.ok("POST", p.path("/equip/"+sword), map[string]any{})
	crit := func() float64 {
		return p.ok("GET", p.path("/sheet"), nil)["derived"].(map[string]any)["crit_pct"].(float64)
	}
	before := crit()
	res = p.ok("POST", path(sword), withScroll("enchant_weapon_crit_1"))
	require.Equal(t, "enchant_weapon_attack_1", res["replaced"])
	require.Equal(t, before+1, crit(), "+1 crit from the enchant")
	worn := p.ok("GET", p.path("/equipment"), nil)["slots"].(map[string]any)["mainhand"].(map[string]any)
	require.Equal(t, "enchant_weapon_crit_1", worn["enchant_code"])
	require.Equal(t, []any{"weapon"}, worn["enchant_effect"].(map[string]any)["slots"])

	// A scroll recipe says what its scroll does.
	var recipe map[string]any
	for _, r := range items(p.ok("GET", "/v1/recipes?profession=enchanting", nil), "recipes") {
		if r["output_template_code"] == "enchant_weapon_attack_1" {
			recipe = r
		}
	}
	require.NotNil(t, recipe)
	require.Equal(t, attack, recipe["output"].(map[string]any)["enchant"].(map[string]any)["modifiers"])

	// A sale and its buyback keep the enchant.
	p.ok("POST", p.path("/unequip/mainhand"), map[string]any{})
	bb := p.ok("POST", "/v1/vendor/general_ironhold/sell", map[string]any{"item_instance_id": sword, "character_id": p.char})
	back := p.ok("POST", "/v1/vendor/general_ironhold/buyback/"+bb["buyback_id"].(string), map[string]any{"character_id": p.char})
	require.Equal(t, "enchant_weapon_crit_1", bagRow(back["instance_id"].(string))["enchant_code"])

	// A mining enchant on worn jewelry raises a new mining run's yield.
	p.ok("POST", path(necklace), withScroll("enchant_jewelry_mining_1"))
	p.ok("POST", p.path("/equip/"+necklace), map[string]any{})
	run := p.ok("POST", p.path("/activity"), map[string]any{"kind": "mining", "duration_seconds": 3600})
	var mult float64
	require.NoError(t, w.p.QueryRow(w.ctx, `SELECT COALESCE((snapshot->>'yield_mult')::float8, 0) FROM activity.contract WHERE id = $1`,
		run["contract"].(map[string]any)["id"]).Scan(&mult))
	require.InDelta(t, 0.02, mult, 1e-9)
}
