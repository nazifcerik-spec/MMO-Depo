//go:build integration

package integration

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/otter/server/internal/inventory"
)

// A full bag (100/100, DD-246) refuses everything that would add an item —
// buying, buying back, crafting, market buys, mail item claims and starting
// a gathering / hunt / dungeon run — before any gold moves. Stackables that
// fit on an existing stack still go through.
func TestHTTPFullBagRefusesNewItems(t *testing.T) {
	w := newWorld(t)
	p := newHTTPPlayer(t, "warrior")
	me := player{ID: uuid.MustParse(p.char)}
	_, err := w.p.Exec(w.ctx, `UPDATE character.character SET level = 6 WHERE id = $1`, me.ID)
	require.NoError(t, err)
	p.mint(1_000_000)

	// Before the bag fills: inputs for a craft, a potion stack, a sold item to
	// buy back, and an item listed / mailed by another player.
	p.ok("POST", p.path("/professions/learn"), map[string]any{"profession_code": "leatherworking"})
	w.give(me, "forest_rat_pelt", 3)
	w.give(me, "common_grass", 1)
	w.give(me, "potion_hp_small", 1)
	sold := w.give(me, "sword_training", 1)
	bb := p.ok("POST", "/v1/vendor/general_ironhold/sell", map[string]any{"item_instance_id": sold.ID, "character_id": p.char})

	o := newHTTPPlayer(t, "mage")
	other := player{ID: uuid.MustParse(o.char)}
	o.mint(10_000)
	listed := w.give(other, "sword_training", 1)
	l := o.ok("POST", o.path("/market/sell"), map[string]any{"item_instance_id": listed.ID, "stack_quantity": 1,
		"unit_price_copper": 100, "duration_hours": 12})
	gift := w.give(other, "helm_iron", 1)
	o.ok("POST", o.path("/mail"), map[string]any{"recipient_name": p.name, "subject": "gift", "body": "x",
		"item_ids": []string{gift.ID.String()}})
	mailID := items(p.ok("GET", p.path("/mail"), nil), "mail")[0]["id"].(string)

	// Fill the remaining slots.
	view := p.ok("GET", p.path("/inventory"), nil)
	var filler string
	for i := int(view["slots_used"].(float64)); i < inventory.BagSlotCap; i++ {
		filler = w.give(me, "sword_training", 1).ID.String()
	}
	require.Equal(t, float64(inventory.BagSlotCap), p.ok("GET", p.path("/inventory"), nil)["slots_used"])
	gold := w.balance(me)

	full := func(method, path string, body map[string]any) { p.fail(method, path, body, 409, "E_BAG_FULL") }
	full("POST", "/v1/vendor/general_ironhold/buy", map[string]any{"template_code": "sword_training", "character_id": p.char})
	full("POST", "/v1/vendor/general_ironhold/buyback/"+bb["buyback_id"].(string), map[string]any{"character_id": p.char})
	full("POST", p.path("/craft"), map[string]any{"recipe_code": "recipe_helm_leather"})
	full("POST", p.path("/market/buy/"+l["id"].(string)), map[string]any{"expected_unit_price_copper": 100})
	full("POST", p.path("/mail/"+mailID+"/claim"), map[string]any{})
	full("POST", p.path("/activity"), map[string]any{"kind": "mining", "duration_seconds": 3600})
	full("POST", p.path("/hunt"), map[string]any{"mob_code": "young_boar", "duration_seconds": 600})
	full("POST", p.path("/dungeon"), map[string]any{"dungeon_code": "crypt_of_ashes"})
	require.Equal(t, gold, w.balance(me), "nothing was charged")
	require.Equal(t, 3.0, p.stackOf("forest_rat_pelt"), "nothing was consumed")

	// A stackable that fits on an existing stack is still fine.
	p.ok("POST", "/v1/vendor/general_ironhold/buy", map[string]any{"template_code": "potion_hp_small", "character_id": p.char})
	require.Equal(t, 2.0, p.stackOf("potion_hp_small"))

	// Freeing one slot lets the next purchase in.
	p.ok("DELETE", p.path("/inventory/"+filler), nil)
	p.ok("POST", "/v1/vendor/general_ironhold/buy", map[string]any{"template_code": "sword_training", "character_id": p.char})
}
