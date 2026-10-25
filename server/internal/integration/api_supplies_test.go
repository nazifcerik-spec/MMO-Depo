//go:build integration

package integration

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

// Hunt supplies through the real router (DD-288): potions and dishes are
// taken from the bag and held for the hunt; collecting spends what the
// fights used and gives the rest back, stopping gives everything back.
func TestHTTPHuntSupplies(t *testing.T) {
	w := newWorld(t)
	p := newHTTPPlayer(t, "warrior")
	me := player{ID: uuid.MustParse(p.char)}
	_, err := w.p.Exec(w.ctx, `UPDATE character.character SET level = 6 WHERE id = $1`, p.char)
	require.NoError(t, err)
	w.give(me, "potion_hp_small", 15)
	w.give(me, "potion_hp_small", 7)
	w.give(me, "food_grilled_guppy", 4)
	w.give(me, "potion_mp_small", 3)

	// units counts a template's units: free in the bag (slotted or not), or
	// held by a hunt.
	units := func(code string, held bool) int {
		var n int
		require.NoError(t, w.p.QueryRow(w.ctx, `
			SELECT COALESCE(SUM(stack_quantity), 0) FROM inventory.item_instance
			 WHERE owner_character_id = $1 AND template_code = $2
			   AND (CASE WHEN $3 THEN lock_owner = 'hunt' ELSE NOT is_locked END)`,
			p.char, code, held).Scan(&n))
		return n
	}
	hunt := func(supplies ...map[string]any) map[string]any {
		return map[string]any{"mob_code": "young_boar", "duration_seconds": 600, "supplies": supplies}
	}
	s := func(code string, qty int) map[string]any {
		return map[string]any{"template_code": code, "quantity": qty}
	}

	// A mana potion only works on a hunt.
	p.fail("POST", p.path("/items/"+p.bagItem("potion_mp_small")+"/use"), map[string]any{}, 409, "E_HUNT_ONLY")

	// Bad requests change nothing.
	p.fail("POST", p.path("/hunt"), hunt(s("dragon_scale", 1)), 400, "E_BAD_SUPPLY")
	p.fail("POST", p.path("/hunt"), hunt(s("common_grass", 1)), 400, "E_BAD_SUPPLY") // not a consumable effect
	p.fail("POST", p.path("/hunt"), hunt(s("potion_hp_small", 0)), 400, "E_BAD_SUPPLY")
	p.fail("POST", p.path("/hunt"), hunt(s("potion_hp_small", 1), s("potion_hp_small", 1)), 400, "E_BAD_SUPPLY")
	p.fail("POST", p.path("/hunt"), hunt(s("potion_hp_huge", 1)), 400, "E_BAD_SUPPLY") // level 650
	p.fail("POST", p.path("/hunt"), hunt(s("potion_hp_small", 23)), 409, "E_NOT_ENOUGH_SUPPLIES")
	require.Equal(t, 22, units("potion_hp_small", false))

	// Stop: everything goes back.
	h := p.ok("POST", p.path("/hunt"), hunt(s("potion_hp_small", 10), s("potion_mp_small", 1)))
	hid := h["contract"].(map[string]any)["id"].(string)
	require.Equal(t, 10, units("potion_hp_small", true), "a partial stack is split off")
	require.Equal(t, 12, units("potion_hp_small", false))
	p.ok("POST", "/v1/hunt/"+hid+"/stop", map[string]any{})
	require.Zero(t, units("potion_hp_small", true))
	require.Equal(t, 22, units("potion_hp_small", false))
	require.Equal(t, 3, units("potion_mp_small", false))

	// Collect: the fights' share is spent, the rest comes back.
	_, err = w.p.Exec(w.ctx, `UPDATE character.character SET hp_current = 40, hp_updated_at = now() WHERE id = $1`, p.char)
	require.NoError(t, err)
	h = p.ok("POST", p.path("/hunt"), hunt(s("potion_hp_small", 20), s("food_grilled_guppy", 2), s("potion_mp_small", 2)))
	c := h["contract"].(map[string]any)
	hid = c["id"].(string)
	require.Len(t, c["loadout"].(map[string]any)["supplies"], 3, "the loadout is frozen on the contract")
	require.Equal(t, 20, units("potion_hp_small", true))
	require.Equal(t, 2, units("potion_hp_small", false))
	_, err = w.p.Exec(w.ctx, `
		UPDATE hunt.contract SET started_at = now() - interval '11 minutes', expires_at = now() - interval '1 minute'
		 WHERE id = $1`, hid)
	require.NoError(t, err)
	res := p.ok("POST", "/v1/hunt/"+hid+"/reveal", map[string]any{})["hunt"].(map[string]any)
	usedOf := map[string]int{}
	for _, u := range res["supplies_used"].([]any) {
		m := u.(map[string]any)
		usedOf[m["code"].(string)] = int(m["qty"].(float64))
	}
	require.Positive(t, usedOf["potion_hp_small"], "a wounded hunter drinks before the first fight: %v", res)
	require.Equal(t, 1, usedOf["food_grilled_guppy"], "one dish: its buff outlasts a 10 minute hunt")
	require.Zero(t, units("potion_hp_small", true), "nothing stays held")
	require.Equal(t, 22-usedOf["potion_hp_small"], units("potion_hp_small", false))
	require.Equal(t, 4-usedOf["food_grilled_guppy"], units("food_grilled_guppy", false))
	require.Equal(t, 3, units("potion_mp_small", false), "a warrior's rage is not carried: no mana potion used")
	p.ok("POST", "/v1/hunt/"+hid+"/reveal", map[string]any{}) // the replay spends nothing more
	require.Equal(t, 22-usedOf["potion_hp_small"], units("potion_hp_small", false))

	// The bag view slots what came back.
	p.ok("GET", p.path("/inventory"), nil)
	var loose int
	require.NoError(t, w.p.QueryRow(w.ctx, `
		SELECT count(*) FROM inventory.item_instance
		 WHERE owner_character_id = $1 AND slot IS NULL AND equipment_slot IS NULL AND NOT is_locked`, p.char).Scan(&loose))
	require.Zero(t, loose)
}
