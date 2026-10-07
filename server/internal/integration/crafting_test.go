//go:build integration

package integration

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/otter/server/internal/crafting"
	"github.com/otter/server/internal/inventory"
)

func newCrafting(t *testing.T, w *testEnv) *crafting.Service {
	t.Helper()
	cat, err := crafting.LoadCatalog(w.ctx, w.p)
	require.NoError(t, err)
	return crafting.NewService(crafting.NewPgRepo(w.p), crafting.NewCatalogService(cat, nil), w.inv, w.items, w.txn, w.chars,
		crafting.DefaultServiceConfig(), nil)
}

func (w *testEnv) count(p player, template string) int {
	w.t.Helper()
	var n int
	require.NoError(w.t, w.p.QueryRow(w.ctx, `
		SELECT COALESCE(sum(stack_quantity), 0) FROM inventory.item_instance
		 WHERE owner_character_id = $1 AND template_code = $2`, p.ID, template).Scan(&n))
	return n
}

// A craft that cannot run must cost nothing: no gold, no inputs.
func TestCraftingFailureCostsNothing(t *testing.T) {
	w := newWorld(t)
	svc := newCrafting(t, w)
	p := w.player("warrior", 1)
	w.gold(p, 1_000)
	_, err := svc.Learn(w.ctx, p.ID, p.Account, "alchemy")
	require.NoError(t, err)
	_, err = svc.Learn(w.ctx, p.ID, p.Account, "blacksmithing")
	require.NoError(t, err)

	_, err = svc.Craft(w.ctx, p.ID, p.Account, "recipe_potion_hp_small")
	require.ErrorIs(t, err, inventory.ErrInsufficient)
	require.Equal(t, int64(1_000), w.balance(p), "missing materials must not cost gold")

	// Two inputs, the second missing: the first must not be eaten either.
	_, err = w.p.Exec(w.ctx, `
		UPDATE crafting.character_profession SET level = 5 WHERE character_id = $1 AND profession_code = 'blacksmithing'`, p.ID)
	require.NoError(t, err)
	w.give(p, "iron_ore", 5)
	_, err = svc.Craft(w.ctx, p.ID, p.Account, "recipe_sword_iron")
	require.ErrorIs(t, err, inventory.ErrInsufficient)
	require.Equal(t, 5, w.count(p, "iron_ore"), "inputs are consumed all-or-nothing")
	require.Equal(t, int64(1_000), w.balance(p))

	w.give(p, "oak_log", 1)
	res, err := svc.Craft(w.ctx, p.ID, p.Account, "recipe_sword_iron")
	require.NoError(t, err)
	require.Equal(t, 0, w.count(p, "iron_ore"))
	require.Equal(t, 0, w.count(p, "oak_log"))
	require.Equal(t, 1, w.count(p, "sword_iron"))
	require.Less(t, w.balance(p), int64(1_000), "a real craft pays its gold cost")
	require.NotEmpty(t, res.OutputInstanceID)
}
