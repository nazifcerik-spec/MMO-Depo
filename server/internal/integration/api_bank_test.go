//go:build integration

package integration

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

// Account bank (DD-314): every character of an account shares 200 slots;
// stacks split and merge, bound and equipped items stay out.
func TestHTTPBank(t *testing.T) {
	w := newWorld(t)
	a := newHTTPPlayer(t, "warrior")
	me := player{ID: uuid.MustParse(a.char)}
	ore := w.give(me, "iron_ore", 10).ID.String()

	bank := a.ok("POST", a.path("/bank/deposit"), map[string]any{"item_id": ore, "quantity": 4})
	items := bank["items"].([]any)
	require.Len(t, items, 1)
	require.EqualValues(t, 4, items[0].(map[string]any)["stack_quantity"])
	require.EqualValues(t, 0, items[0].(map[string]any)["slot"])
	// The rest goes onto the same banked stack.
	bank = a.ok("POST", a.path("/bank/deposit"), map[string]any{"item_id": ore})
	items = bank["items"].([]any)
	require.Len(t, items, 1)
	require.EqualValues(t, 10, items[0].(map[string]any)["stack_quantity"])

	// A second character of the same account sees and takes it.
	var acct uuid.UUID
	require.NoError(t, w.p.QueryRow(w.ctx, `SELECT account_id FROM character.character WHERE id = $1`, a.char).Scan(&acct))
	b := a.ok("POST", "/v1/character", map[string]any{"name": "B" + randHex(5), "race_code": "human", "class_code": "mage"})
	bid := b["id"].(string)
	banked := items[0].(map[string]any)["instance_id"].(string)
	bv := a.ok("POST", "/v1/character/"+bid+"/bank/withdraw", map[string]any{"item_id": banked, "quantity": 3})
	require.EqualValues(t, 7, bv["items"].([]any)[0].(map[string]any)["stack_quantity"])
	var inBag int
	require.NoError(t, w.p.QueryRow(w.ctx, `SELECT COALESCE(SUM(stack_quantity), 0) FROM inventory.item_instance
		WHERE owner_character_id = $1 AND template_code = 'iron_ore' AND slot IS NOT NULL`, bid).Scan(&inBag))
	require.Equal(t, 3, inBag)

	// Deleting the depositing character keeps the bank.
	_, err := w.p.Exec(w.ctx, `DELETE FROM inventory.item_instance WHERE owner_character_id = $1`, a.char)
	require.NoError(t, err)
	require.Len(t, a.ok("GET", "/v1/character/"+bid+"/bank", nil)["items"], 1)

	// Bound items stay with their character; another account sees nothing.
	sword := w.give(me, "sword_iron", 1)
	_, err = w.p.Exec(w.ctx, `UPDATE inventory.item_instance SET is_bound = true WHERE id = $1`, sword.ID)
	require.NoError(t, err)
	a.fail("POST", a.path("/bank/deposit"), map[string]any{"item_id": sword.ID.String()}, 409, "E_ITEM_BOUND")
	c := newHTTPPlayer(t, "rogue")
	require.Empty(t, c.ok("GET", c.path("/bank"), nil)["items"])
	c.fail("POST", c.path("/bank/withdraw"), map[string]any{"item_id": banked}, 404, "E_NOT_FOUND")
	_ = acct
}
