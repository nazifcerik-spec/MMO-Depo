//go:build integration

package integration

import (
	"fmt"
	"strings"
	"sync"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/otter/server/internal/inventory"
)

func (c *hc) slotOf(id string) float64 {
	for _, it := range items(c.ok("GET", c.path("/inventory"), nil), "items") {
		if it["instance_id"] == id {
			return it["slot"].(float64)
		}
	}
	return -1
}

func (c *hc) stackOf(code string) float64 {
	n := 0.0
	for _, it := range items(c.ok("GET", c.path("/inventory"), nil), "items") {
		if it["template_code"] == code {
			n += it["stack_quantity"].(float64)
		}
	}
	return n
}

func (c *hc) gold() int64 {
	return int64(c.ok("GET", c.path("/wallet"), nil)["gold_copper"].(float64))
}

func TestHTTPInventoryOps(t *testing.T) {
	w := newWorld(t)
	p := newHTTPPlayer(t, "warrior")
	me := player{ID: uuid.MustParse(p.char)}
	p.mint(100_000)
	ore := w.give(me, "iron_ore", 10).ID.String()
	sword := w.give(me, "sword_iron", 1).ID.String()
	inv := func(id, op string) string { return p.path("/inventory/" + id + op) }

	// Split part of a stack into an empty slot.
	part := p.ok("POST", inv(ore, "/split"), map[string]any{"amount": 3, "target_slot": 40})["new_instance_id"].(string)
	require.Equal(t, 40.0, p.slotOf(part))
	p.fail("POST", inv(ore, "/split"), map[string]any{"amount": 0, "target_slot": 41}, 400, "E_INVALID_SPLIT")
	p.fail("POST", inv(ore, "/split"), map[string]any{"amount": 7, "target_slot": 41}, 400, "E_INVALID_SPLIT")
	p.fail("POST", inv(sword, "/split"), map[string]any{"amount": 1, "target_slot": 41}, 409, "E_NOT_STACKABLE")
	p.fail("POST", inv(ore, "/split"), map[string]any{"amount": 1, "target_slot": 40}, 409, "E_SLOT_OCCUPIED")
	p.fail("POST", inv(ore, "/split"), map[string]any{"amount": 1, "target_slot": inventory.BagSlotCap}, 400, "E_SLOT_OUT_OF_RANGE")
	p.fail("POST", inv(uuid.NewString(), "/split"), map[string]any{"amount": 1, "target_slot": 41}, 404, "E_NOT_FOUND")
	p.fail("POST", inv("nope", "/split"), map[string]any{"amount": 1, "target_slot": 41}, 400, "E_INVALID_ID")
	p.fail("POST", inv(ore, "/split"), "{", 400, "E_INVALID_JSON")

	// Move: to an empty slot, swap with an occupied one, onto itself.
	p.ok("POST", inv(part, "/move"), map[string]any{"target_slot": 50})
	require.Equal(t, 50.0, p.slotOf(part))
	swordSlot := p.slotOf(sword)
	p.ok("POST", inv(part, "/move"), map[string]any{"target_slot": swordSlot})
	require.Equal(t, swordSlot, p.slotOf(part))
	require.Equal(t, 50.0, p.slotOf(sword))
	p.ok("POST", inv(part, "/move"), map[string]any{"target_slot": swordSlot})
	p.fail("POST", inv(part, "/move"), map[string]any{"target_slot": -1}, 400, "E_SLOT_OUT_OF_RANGE")
	p.fail("POST", inv(uuid.NewString(), "/move"), map[string]any{"target_slot": 1}, 404, "E_NOT_FOUND")
	p.fail("POST", inv("nope", "/move"), map[string]any{"target_slot": 1}, 400, "E_INVALID_ID")
	p.fail("POST", inv(part, "/move"), "", 400, "E_INVALID_JSON")

	// A listed stack leaves the bag (escrow) and comes back when cancelled.
	l := p.ok("POST", p.path("/market/sell"), map[string]any{"item_instance_id": ore,
		"stack_quantity": 7, "unit_price_copper": 50, "duration_hours": 24})
	require.Equal(t, -1.0, p.slotOf(ore))
	p.fail("POST", inv(ore, "/move"), map[string]any{"target_slot": 55}, 404, "E_NOT_FOUND")
	p.fail("DELETE", inv(ore, ""), nil, 404, "E_NOT_FOUND")
	// Cancelling merges the stack back into `part` (deleting the listed
	// instance); the listing record must survive that.
	p.ok("POST", p.path("/market/cancel/"+l["id"].(string)), map[string]any{})
	require.Equal(t, 10.0, p.stackOf("iron_ore"))

	// Same for a buyer who already holds the template.
	b := newHTTPPlayer(t, "mage")
	b.mint(10_000)
	w.give(player{ID: uuid.MustParse(b.char)}, "iron_ore", 2)
	l = p.ok("POST", p.path("/market/sell"), map[string]any{"item_instance_id": part,
		"stack_quantity": 4, "unit_price_copper": 50, "duration_hours": 24})
	sold := b.ok("POST", b.path("/market/buy/"+l["id"].(string)), map[string]any{"expected_unit_price_copper": 50})
	require.Equal(t, "sold", sold["state"], "%v", sold)
	require.Equal(t, 6.0, b.stackOf("iron_ore"))

	// Delete is final.
	p.ok("DELETE", inv(part, ""), nil)
	p.fail("DELETE", inv(part, ""), nil, 404, "E_NOT_FOUND")
	p.fail("DELETE", inv("nope", ""), nil, 400, "E_INVALID_ID")

	// Someone else's bag is invisible.
	b.fail("POST", "/v1/character/"+p.char+"/inventory/"+sword+"/move", map[string]any{"target_slot": 1}, 404, "")
	b.fail("GET", "/v1/character/"+p.char+"/inventory", nil, 404, "")

	// Equip / unequip.
	_, err := w.p.Exec(w.ctx, `UPDATE character.character SET level = 5 WHERE id = $1`, me.ID)
	require.NoError(t, err)
	p.fail("POST", p.path("/equip/"+w.give(me, "linseed", 1).ID.String()), map[string]any{}, 400, "E_NOT_EQUIPPABLE")
	p.ok("POST", p.path("/equip/"+sword), map[string]any{})
	eq := p.ok("GET", p.path("/equipment"), nil)["slots"].(map[string]any)
	require.Equal(t, sword, eq["mainhand"].(map[string]any)["instance_id"])
	p.ok("POST", p.path("/unequip/mainhand"), map[string]any{})
	p.fail("POST", p.path("/unequip/mainhand"), map[string]any{}, 404, "E_NOTHING_EQUIPPED")
	p.fail("POST", p.path("/unequip/tail"), map[string]any{}, 400, "E_UNKNOWN_SLOT")
	p.ok("POST", p.path("/equip/"+sword), map[string]any{})

	// Full bag: unequipping has nowhere to go; overflow waits until there is room.
	for n := len(items(p.ok("GET", p.path("/inventory"), nil), "items")); n < inventory.BagSlotCap; n++ {
		w.give(me, "helm_iron", 1)
	}
	p.fail("POST", p.path("/unequip/mainhand"), map[string]any{}, 409, "E_BAG_FULL")
	extra := w.give(me, "helm_iron", 1)
	require.False(t, w.slotted(extra.ID))
	p.ok("DELETE", inv(p.bagItem("helm_iron"), ""), nil)
	require.NotEqual(t, -1.0, p.slotOf(extra.ID.String()), "overflow is slotted once there is room")
}

func TestHTTPGuildLifecycle(t *testing.T) {
	a := newHTTPPlayer(t, "warrior")
	b := newHTTPPlayer(t, "mage")
	a.mint(200_000)
	require.Nil(t, a.ok("GET", a.path("/guild"), nil)["guild"])

	g := a.ok("POST", a.path("/guild"), map[string]any{"name": "Guild " + randHex(4), "tag": strings.ToUpper(randHex(2))})
	guildID := g["guild"].(map[string]any)["id"].(string)

	inv := a.ok("POST", a.path("/guild/invite"), map[string]any{"target_name": b.name})
	b.ok("POST", b.path("/guild/invites/"+inv["id"].(string)+"/decline"), map[string]any{})
	b.fail("POST", b.path("/guild/invites/"+inv["id"].(string)+"/decline"), map[string]any{}, 404, "E_NOT_FOUND")
	inv = a.ok("POST", a.path("/guild/invite"), map[string]any{"target_name": b.name})
	b.ok("POST", b.path("/guild/invites/"+inv["id"].(string)+"/accept"), map[string]any{})

	a.fail("POST", a.path("/guild/rank"), map[string]any{"target_character_id": b.char}, 400, "E_VALIDATION")
	a.fail("POST", a.path("/guild/rank"), "{", 400, "E_INVALID_JSON")
	a.ok("POST", a.path("/guild/rank"), map[string]any{"target_character_id": b.char, "rank_index": 2})
	b.fail("POST", b.path("/guild/rank"), map[string]any{"target_character_id": a.char, "rank_index": 4}, 403, "E_GUILD_FORBIDDEN")

	a.fail("POST", a.path("/guild/leave"), map[string]any{}, 409, "E_LEADER_CANNOT_LEAVE")
	b.ok("POST", b.path("/guild/leave"), map[string]any{})
	b.fail("POST", b.path("/guild/leave"), map[string]any{}, 409, "E_NOT_IN_GUILD")

	a.fail("POST", a.path("/guild/disband/cancel"), map[string]any{}, 409, "E_DISBAND_NOT_PENDING")
	a.ok("POST", a.path("/guild/disband"), map[string]any{})
	a.fail("POST", a.path("/guild/disband"), map[string]any{}, 409, "E_DISBAND_SCHEDULED")
	a.ok("POST", a.path("/guild/disband/cancel"), map[string]any{})

	a.ok("GET", "/v1/guilds/"+guildID, nil)
	a.fail("GET", "/v1/guilds/nope", nil, 400, "E_INVALID_ID")
	a.fail("POST", a.path("/guild/invites/nope/accept"), map[string]any{}, 400, "E_INVALID_ID")
	b.fail("GET", "/v1/character/"+a.char+"/guild", nil, 404, "")
}

// Racing founders of the same name: exactly one guild, and every loser keeps
// their gold (either never charged or refunded).
func TestHTTPGuildCreateRace(t *testing.T) {
	const n = 6
	ps := make([]*hc, n)
	for i := range ps {
		ps[i] = newHTTPPlayer(t, "warrior")
		ps[i].mint(150_000)
	}
	name := "Race " + randHex(4)
	codes := make([]int, n)
	var wg sync.WaitGroup
	for i, p := range ps {
		wg.Add(1)
		go func() {
			defer wg.Done()
			codes[i], _ = p.call("POST", p.path("/guild"), map[string]any{"name": name, "tag": strings.ToUpper(randHex(2))})
		}()
	}
	wg.Wait()
	won := 0
	for i, p := range ps {
		if codes[i] == 201 {
			won++
			require.Equal(t, int64(150_000-100_000), p.gold())
			continue
		}
		require.Equal(t, 409, codes[i])
		require.Equal(t, int64(150_000), p.gold(), "loser %d must keep its gold", i)
	}
	require.Equal(t, 1, won, "%v", codes)
}

// Every class gets its own base attributes; bad names and classes are refused.
func TestHTTPCharacterClasses(t *testing.T) {
	seen := map[string]bool{}
	for _, class := range []string{"warrior", "rogue", "ranger", "mage", "monk", "cleric", "paladin", "druid", "bard", "shaman"} {
		p := newHTTPPlayer(t, class)
		seen[fmt.Sprint(p.ok("GET", p.path("/sheet"), nil)["primary"])] = true
	}
	require.Len(t, seen, 10, "each class has distinct base attributes")
	p := newHTTPPlayer(t, "warrior")
	p.fail("POST", "/v1/character", map[string]any{"name": "x", "race_code": "human", "class_code": "warrior"}, 400, "")
	p.fail("POST", "/v1/character", map[string]any{"name": "Valid" + randHex(3), "race_code": "human", "class_code": "hunter"}, 400, "")
}
