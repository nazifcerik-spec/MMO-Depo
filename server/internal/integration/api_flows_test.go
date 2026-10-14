//go:build integration

package integration

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func (c *hc) mint(copper int) {
	c.ok("POST", "/v1/dev/mint", map[string]any{"character_id": c.char, "amount": copper, "currency": "gold"})
}

func (c *hc) bagItem(code string) string {
	for _, it := range items(c.ok("GET", c.path("/inventory"), nil), "items") {
		if it["template_code"] == code {
			return it["instance_id"].(string)
		}
	}
	return ""
}

// Success paths the gameplay loop test only reaches as errors.
func TestHTTPSuccessFlows(t *testing.T) {
	w := newWorld(t)
	p := newHTTPPlayer(t, "warrior")
	id := uuid.MustParse(p.char)
	me := player{ID: id}
	p.mint(100_000)

	// Crafting: learn, gather the inputs, craft, profession XP.
	p.ok("POST", p.path("/professions/learn"), map[string]any{"profession_code": "alchemy"})
	p.ok("POST", p.path("/professions/learn"), map[string]any{"profession_code": "alchemy"}) // idempotent
	w.give(me, "common_grass", 4)
	crafted := p.ok("POST", p.path("/craft"), map[string]any{"recipe_code": "recipe_potion_hp_small"})
	require.NotEmpty(t, crafted)
	profs := items(p.ok("GET", p.path("/professions"), nil), "professions")
	require.Positive(t, profs[0]["xp"].(float64)+profs[0]["level"].(float64)-1, "%v", profs)
	potion := p.bagItem("potion_hp_small")
	require.NotEmpty(t, potion)

	// Consumables: heal when wounded, not during a fight, not for gear.
	_, err := w.p.Exec(w.ctx, `UPDATE character.character SET hp_current = 5, hp_updated_at = now() WHERE id = $1`, id)
	require.NoError(t, err)
	used := p.ok("POST", p.path("/items/"+potion+"/use"), map[string]any{})
	require.Greater(t, used["hp_current"].(float64), 5.0, "%v", used)
	w.give(me, "sword_iron", 1)
	p.fail("POST", p.path("/items/"+p.bagItem("sword_iron")+"/use"), map[string]any{}, 400, "E_NOT_USABLE")
	_, err = w.p.Exec(w.ctx, `UPDATE character.character SET hp_current = 50, hp_updated_at = now() WHERE id = $1`, id)
	require.NoError(t, err)
	h := p.ok("POST", p.path("/hunt"), map[string]any{"mob_code": "young_boar", "duration_seconds": 600})
	if second := p.bagItem("potion_hp_small"); second != "" {
		p.fail("POST", p.path("/items/"+second+"/use"), map[string]any{}, 409, "E_IN_COMBAT")
	}
	p.ok("POST", "/v1/hunt/"+h["contract"].(map[string]any)["id"].(string)+"/stop", map[string]any{})

	// Equip gates: level and class.
	p.fail("POST", p.path("/equip/"+p.bagItem("sword_iron")), map[string]any{}, 403, "E_LEVEL_REQUIREMENT")
	w.give(me, "staff_apprentice", 1)
	_, err = w.p.Exec(w.ctx, `UPDATE character.character SET level = 5 WHERE id = $1`, id)
	require.NoError(t, err)
	p.fail("POST", p.path("/equip/"+p.bagItem("staff_apprentice")), map[string]any{}, 403, "E_CLASS_RESTRICTION")
	p.ok("POST", p.path("/equip/"+p.bagItem("sword_iron")), map[string]any{})
	p.fail("POST", p.path("/equip/"+uuid.NewString()), map[string]any{}, 404, "")

	// Travel arrives.
	loc := p.ok("GET", p.path("/location"), nil)
	from := loc["region"].(map[string]any)["code"].(string)
	var dest string
	for _, n := range items(loc, "neighbors") {
		dest = n["code"].(string)
		break
	}
	tr := p.ok("POST", p.path("/travel"), map[string]any{"to_region_code": dest})
	trID := tr["travel"].(map[string]any)["id"].(string)
	p.fail("POST", p.path("/travel"), map[string]any{"to_region_code": dest}, 409, "")
	_, err = w.p.Exec(w.ctx, `UPDATE world.travel SET arrives_at = now() - interval '1 second' WHERE id = $1`, trID)
	require.NoError(t, err)
	p.ok("POST", "/v1/travel/"+trID+"/complete", map[string]any{})
	now := p.ok("GET", p.path("/location"), nil)
	require.Equal(t, dest, now["region"].(map[string]any)["code"])
	require.NotEqual(t, from, dest)
	p.fail("POST", p.path("/travel"), map[string]any{"to_region_code": "nowhere"}, 404, "")

	// Market purchase by a second player.
	b := newHTTPPlayer(t, "mage")
	b.mint(50_000)
	l := p.ok("POST", p.path("/market/sell"), map[string]any{"item_instance_id": p.bagItem("staff_apprentice"),
		"stack_quantity": 1, "unit_price_copper": 3000, "duration_hours": 24})
	b.ok("POST", b.path("/market/buy/"+l["id"].(string)), map[string]any{"expected_unit_price_copper": 3000})
	require.NotEmpty(t, b.bagItem("staff_apprentice"))
	b.fail("POST", b.path("/market/buy/"+l["id"].(string)), map[string]any{}, 409, "")
	b.fail("POST", b.path("/market/cancel/"+l["id"].(string)), map[string]any{}, 403, "")

	// Character name clash and delete-grace rules.
	p.fail("POST", p.path("/rename"), map[string]any{"name": b.name}, 409, "")
	p.ok("POST", p.path("/delete-request"), map[string]any{})
	p.fail("POST", p.path("/delete-request"), map[string]any{}, 409, "")
	p.ok("POST", p.path("/delete-cancel"), map[string]any{})
	p.fail("GET", "/v1/character/"+b.char+"/wallet", nil, 404, "")
	p.fail("POST", "/v1/dev/mint", map[string]any{"character_id": p.char, "amount": 1, "currency": "doubloons"}, 400, "")
}
