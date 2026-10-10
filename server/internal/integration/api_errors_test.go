//go:build integration

package integration

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

// Error paths of the remaining handlers: malformed input, foreign ids,
// state gates. Each must be a precise 4xx, never a 500.
func TestHTTPHandlerErrors(t *testing.T) {
	w := newWorld(t)
	p := newHTTPPlayer(t, "warrior")
	o := newHTTPPlayer(t, "mage")
	me := player{ID: uuid.MustParse(p.char)}
	p.mint(200_000)

	// Vendor.
	v := "/v1/vendor/general_ironhold"
	p.fail("GET", "/v1/vendor/nowhere", nil, 404, "")
	p.fail("POST", v+"/buy", "{", 400, "E_INVALID_JSON")
	p.fail("POST", v+"/buy", map[string]any{"character_id": "nope", "template_code": "potion_hp_small"}, 400, "E_INVALID_ID")
	p.fail("POST", v+"/buy", map[string]any{"character_id": o.char, "template_code": "potion_hp_small"}, 404, "E_NOT_FOUND")
	p.fail("POST", v+"/buy", map[string]any{"character_id": p.char, "template_code": "dragon_egg"}, 400, "E_ITEM_NOT_SOLD")
	p.fail("POST", v+"/sell", map[string]any{"character_id": p.char, "item_instance_id": "nope"}, 400, "")
	p.fail("POST", v+"/sell", map[string]any{"character_id": p.char, "item_instance_id": uuid.NewString()}, 404, "")
	p.fail("POST", v+"/buyback/"+uuid.NewString(), "{", 400, "E_INVALID_JSON")
	p.fail("POST", v+"/buyback/nope", map[string]any{"character_id": p.char}, 400, "E_INVALID_ID")
	p.fail("POST", v+"/buyback/"+uuid.NewString(), map[string]any{"character_id": p.char}, 404, "")
	sold := w.give(me, "linseed", 3)
	p.ok("POST", v+"/sell", map[string]any{"character_id": p.char, "item_instance_id": sold.ID.String()})
	bb := items(p.ok("GET", p.path("/buybacks"), nil), "buybacks")
	require.Len(t, bb, 1)
	o.fail("POST", v+"/buyback/"+bb[0]["id"].(string), map[string]any{"character_id": o.char}, 404, "")
	p.ok("POST", v+"/buyback/"+bb[0]["id"].(string), map[string]any{"character_id": p.char})

	// Consumables.
	p.fail("POST", p.path("/items/nope/use"), map[string]any{}, 400, "E_INVALID_ID")
	p.fail("POST", p.path("/items/"+uuid.NewString()+"/use"), map[string]any{}, 404, "E_NOT_FOUND")
	potion := w.give(player{ID: uuid.MustParse(o.char)}, "potion_hp_small", 1)
	p.fail("POST", p.path("/items/"+potion.ID.String()+"/use"), map[string]any{}, 404, "E_NOT_FOUND")
	o.fail("POST", o.path("/items/"+potion.ID.String()+"/use"), map[string]any{}, 409, "E_HP_FULL")

	// Guild: member management without a guild / without rights.
	p.fail("POST", p.path("/guild/kick"), map[string]any{"target_character_id": o.char}, 409, "E_NOT_IN_GUILD")
	p.fail("POST", p.path("/guild/transfer"), map[string]any{"target_character_id": o.char}, 409, "E_NOT_IN_GUILD")
	p.fail("POST", p.path("/guild/motd"), map[string]any{"motd": "hi"}, 409, "E_NOT_IN_GUILD")
	p.fail("POST", p.path("/guild/invite"), map[string]any{"target_name": o.name}, 409, "E_NOT_IN_GUILD")
	p.fail("POST", p.path("/guild/kick"), "{", 400, "E_INVALID_JSON")
	p.fail("POST", p.path("/guild/transfer"), "{", 400, "E_INVALID_JSON")
	p.fail("POST", p.path("/guild/motd"), "{", 400, "E_INVALID_JSON")
	p.fail("POST", p.path("/guild/invite"), "{", 400, "E_INVALID_JSON")
	p.fail("POST", p.path("/guild"), "{", 400, "E_INVALID_JSON")
	require.Equal(t, 0.0, p.ok("GET", p.path("/guild/invites"), nil)["count"])
	p.ok("POST", p.path("/guild"), map[string]any{"name": "Err " + randHex(4), "tag": "E" + randHex(2)})
	p.fail("POST", p.path("/guild/kick"), map[string]any{"target_character_id": o.char}, 404, "E_NOT_MEMBER")
	p.fail("POST", p.path("/guild/transfer"), map[string]any{"target_character_id": o.char}, 404, "E_NOT_MEMBER")
	p.fail("POST", p.path("/guild/invite"), map[string]any{"target_name": p.name}, 400, "E_SELF_TARGET")
	p.fail("POST", p.path("/guild/invite"), map[string]any{"target_name": "Nobody" + randHex(3)}, 404, "E_TARGET_NOT_FOUND")
	p.ok("POST", p.path("/guild/motd"), map[string]any{"motd": "Welcome"})

	// World.
	dest := items(p.ok("GET", p.path("/location"), nil), "neighbors")[0]["code"]
	trip := p.ok("POST", p.path("/travel"), map[string]any{"to_region_code": dest})["travel"].(map[string]any)["id"].(string)
	p.fail("POST", "/v1/travel/"+trip+"/complete", map[string]any{}, 409, "")
	p.fail("POST", "/v1/travel/nope/complete", map[string]any{}, 400, "")
	p.fail("POST", "/v1/travel/nope/cancel", map[string]any{}, 400, "")
	p.ok("POST", "/v1/travel/"+trip+"/cancel", map[string]any{})
	p.ok("POST", "/v1/travel/"+trip+"/cancel", map[string]any{}) // idempotent
	p.fail("POST", p.path("/travel"), "{", 400, "")
	require.NotEmpty(t, items(p.ok("GET", p.path("/travel/history?limit=5"), nil), "travels"))

	// Crafting / character / market input errors.
	p.fail("POST", p.path("/craft"), "{", 400, "")
	p.fail("POST", p.path("/craft"), map[string]any{"recipe_code": "recipe_ghost"}, 404, "")
	p.ok("GET", "/v1/recipes?profession=alchemy", nil)
	p.fail("POST", p.path("/rename"), "{", 400, "")
	p.fail("POST", p.path("/delete-cancel"), map[string]any{}, 409, "")
	p.fail("POST", "/v1/character", "{", 400, "")
	p.fail("POST", p.path("/market/sell"), "{", 400, "")
	p.fail("POST", p.path("/market/buy/nope"), map[string]any{}, 400, "")
	p.ok("GET", "/v1/market/history/linseed", nil)
	p.ok("GET", p.path("/market/listings"), nil)
	p.ok("GET", p.path("/dungeon/history?limit=3"), nil)
	p.ok("GET", p.path("/hunt/history?limit=3"), nil)
}
