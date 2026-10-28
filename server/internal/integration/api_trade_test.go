//go:build integration

package integration

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

// Player trading (DD-315): both sides put items up, both get ready, the
// items change hands and each side pays 8% of what it gives at the average
// market price.
func TestHTTPTrade(t *testing.T) {
	w := newWorld(t)
	a, b := newHTTPPlayer(t, "warrior"), newHTTPPlayer(t, "mage")
	pa, pb := player{ID: uuid.MustParse(a.char)}, player{ID: uuid.MustParse(b.char)}
	a.mint(2_000_000)
	b.mint(2_000_000)
	ore := w.give(pa, "iron_ore", 10).ID.String()
	potion := w.give(pb, "potion_hp_small", 3).ID.String()
	// Some iron ore sold at 100 gold each.
	_, err := w.p.Exec(w.ctx, `INSERT INTO market.price_history (template_code, unit_price_copper, quantity)
		VALUES ('iron_ore', 1000000, 2), ('iron_ore', 1000000, 1)`)
	require.NoError(t, err)

	tr := a.ok("POST", a.path("/trade"), map[string]any{"target_name": b.name})["trade"].(map[string]any)
	id := tr["id"].(string)
	b.fail("POST", b.path("/trade"), map[string]any{"target_name": a.name}, 409, "E_TRADE_BUSY")
	// Both ready with nothing up: refused, the marks clear.
	a.ok("POST", a.path("/trade/"+id+"/ready"), map[string]any{})
	b.fail("POST", b.path("/trade/"+id+"/ready"), map[string]any{}, 409, "E_TRADE_EMPTY")

	tr = a.ok("POST", a.path("/trade/"+id+"/offer"), map[string]any{"item_id": ore, "quantity": 10})["trade"].(map[string]any)
	me := tr["me"].(map[string]any)
	// Other tests sell iron ore too: the average comes from the same history.
	var avg int64
	require.NoError(t, w.p.QueryRow(w.ctx, `SELECT round(sum(unit_price_copper::numeric * quantity) / sum(quantity))
		FROM market.price_history WHERE template_code = 'iron_ore' AND sold_at > now() - interval '30 days'`).Scan(&avg))
	require.EqualValues(t, avg, me["items"].([]any)[0].(map[string]any)["unit_value_copper"])
	taxA := avg * 10 * 8 / 100
	require.InDelta(t, taxA, me["tax_copper"], 1, "average × 10 × 8%")
	b.ok("POST", b.path("/trade/"+id+"/offer"), map[string]any{"item_id": potion, "quantity": 2})

	goldA, goldB := a.gold(), b.gold()
	a.ok("POST", a.path("/trade/"+id+"/ready"), map[string]any{})
	// A change clears both ready marks.
	b.ok("POST", b.path("/trade/"+id+"/offer"), map[string]any{"item_id": potion, "quantity": 1})
	require.Equal(t, false, b.ok("GET", b.path("/trade"), nil)["trade"].(map[string]any)["them"].(map[string]any)["ready"])
	a.ok("POST", a.path("/trade/"+id+"/ready"), map[string]any{})
	done := b.ok("POST", b.path("/trade/"+id+"/ready"), map[string]any{})["trade"].(map[string]any)
	require.Equal(t, "done", done["state"])
	taxB := int64(done["me"].(map[string]any)["tax_copper"].(float64))
	require.Positive(t, taxB)

	require.InDelta(t, goldA-taxA, a.gold(), 1)
	require.Equal(t, goldB-taxB, b.gold())
	count := func(p *hc, tmpl string) (n int) {
		require.NoError(t, w.p.QueryRow(w.ctx, `SELECT COALESCE(SUM(stack_quantity), 0) FROM inventory.item_instance
			WHERE owner_character_id = $1 AND template_code = $2 AND slot IS NOT NULL`, p.char, tmpl).Scan(&n))
		return n
	}
	require.Equal(t, 10, count(b, "iron_ore"))
	require.Zero(t, count(a, "iron_ore"))
	require.Equal(t, 1, count(a, "potion_hp_small"))
	require.Equal(t, 2, count(b, "potion_hp_small"))
	require.Nil(t, a.ok("GET", a.path("/trade"), nil)["trade"])

	// Bound items cannot go up; a cancelled trade moves nothing.
	sword := w.give(pa, "sword_iron", 1)
	_, err = w.p.Exec(w.ctx, `UPDATE inventory.item_instance SET is_bound = true WHERE id = $1`, sword.ID)
	require.NoError(t, err)
	id = a.ok("POST", a.path("/trade"), map[string]any{"target_name": b.name})["trade"].(map[string]any)["id"].(string)
	a.fail("POST", a.path("/trade/"+id+"/offer"), map[string]any{"item_id": sword.ID.String()}, 409, "E_TRADE_ITEM")
	b.ok("POST", b.path("/trade/"+id+"/cancel"), map[string]any{})
	a.fail("POST", a.path("/trade/"+id+"/ready"), map[string]any{}, 409, "E_TRADE_CLOSED")
	_ = pb
}
