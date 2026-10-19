//go:build integration

package integration

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/otter/server/internal/crafting"
	"github.com/otter/server/internal/inventory"
)

func newCrafting(t *testing.T, w *testEnv, now *time.Time) *crafting.Service {
	t.Helper()
	cat, err := crafting.LoadCatalog(w.ctx, w.p)
	require.NoError(t, err)
	return crafting.NewService(crafting.NewPgRepo(w.p), crafting.NewCatalogService(cat, nil), w.inv, w.items, w.txn, w.chars,
		crafting.ServiceConfig{Now: func() time.Time { return *now }}, nil)
}

func (w *testEnv) count(p player, template string) int {
	w.t.Helper()
	var n int
	require.NoError(w.t, w.p.QueryRow(w.ctx, `
		SELECT COALESCE(sum(stack_quantity), 0) FROM inventory.item_instance
		 WHERE owner_character_id = $1 AND template_code = $2`, p.ID, template).Scan(&n))
	return n
}

// A craft that cannot run must cost nothing: no gold, no inputs. A timed
// job (DD-264) takes everything at start and delivers at collection.
func TestCraftingFailureCostsNothing(t *testing.T) {
	w := newWorld(t)
	now := time.Now().UTC()
	svc := newCrafting(t, w, &now)
	p := w.player("warrior", 1)
	w.gold(p, 1_000)
	_, err := svc.Learn(w.ctx, p.ID, p.Account, "alchemy")
	require.NoError(t, err)
	_, err = svc.Learn(w.ctx, p.ID, p.Account, "blacksmithing")
	require.NoError(t, err)

	_, err = svc.Start(w.ctx, p.ID, p.Account, "recipe_potion_hp_small", 1)
	require.ErrorIs(t, err, inventory.ErrInsufficient)
	require.Equal(t, int64(1_000), w.balance(p), "missing materials must not cost gold")

	// Two inputs, the second missing: the first must not be eaten either.
	_, err = w.p.Exec(w.ctx, `
		UPDATE crafting.character_profession SET level = 5 WHERE character_id = $1 AND profession_code = 'blacksmithing'`, p.ID)
	require.NoError(t, err)
	w.give(p, "aluminium_bar", 4) // metal gear is made of bars (DD-274)
	_, err = svc.Start(w.ctx, p.ID, p.Account, "recipe_sword_iron", 1)
	require.ErrorIs(t, err, inventory.ErrInsufficient)
	require.Equal(t, 4, w.count(p, "aluminium_bar"), "inputs are consumed all-or-nothing")
	require.Equal(t, int64(1_000), w.balance(p))

	w.give(p, "young_sapling", 2)
	job, err := svc.Start(w.ctx, p.ID, p.Account, "recipe_sword_iron", 1)
	require.NoError(t, err)
	require.Equal(t, 0, w.count(p, "aluminium_bar"), "inputs are taken at start")
	require.Equal(t, 0, w.count(p, "young_sapling"))
	require.Less(t, w.balance(p), int64(1_000), "a real craft pays its gold cost")
	require.Equal(t, 0, w.count(p, "sword_iron"), "not made yet")

	now = job.ReadyAt
	res, err := svc.Collect(w.ctx, p.ID, p.Account, job.ID)
	require.NoError(t, err)
	require.Equal(t, 1, res.Made)
	require.Equal(t, 1, w.count(p, "sword_iron"))
	require.NotEmpty(t, res.OutputInstanceID)
}

// Collecting early delivers the finished units and refunds the rest.
func TestCraftingEarlyCollectRefunds(t *testing.T) {
	w := newWorld(t)
	now := time.Now().UTC()
	svc := newCrafting(t, w, &now)
	p := w.player("warrior", 1)
	w.gold(p, 1_000)
	_, err := svc.Learn(w.ctx, p.ID, p.Account, "alchemy")
	require.NoError(t, err)
	w.give(p, "common_grass", 12)
	job, err := svc.Start(w.ctx, p.ID, p.Account, "recipe_potion_hp_small", 4) // 3 grass, 5c each
	require.NoError(t, err)
	require.Equal(t, 0, w.count(p, "common_grass"))
	require.Equal(t, int64(1_000-20), w.balance(p))

	now = now.Add(time.Duration(job.UnitSeconds) * time.Second) // one unit done
	res, err := svc.Collect(w.ctx, p.ID, p.Account, job.ID)
	require.NoError(t, err)
	require.Equal(t, 1, res.Made)
	require.Equal(t, 3, res.Refunded)
	require.Equal(t, 1, w.count(p, "potion_hp_small"))
	require.Equal(t, 9, w.count(p, "common_grass"), "three units of grass back")
	require.Equal(t, int64(1_000-5), w.balance(p), "three units of gold back")

	again, err := svc.Collect(w.ctx, p.ID, p.Account, job.ID)
	require.NoError(t, err)
	require.True(t, again.Replayed)
	require.Equal(t, 1, w.count(p, "potion_hp_small"), "collected once")
}

// Cooking (DD-268): a fish, a herb and a log make two dishes that heal.
func TestCookingMakesFood(t *testing.T) {
	w := newWorld(t)
	now := time.Now().UTC()
	svc := newCrafting(t, w, &now)
	p := w.player("warrior", 1)
	w.gold(p, 1_000)
	_, err := svc.Learn(w.ctx, p.ID, p.Account, "cooking")
	require.NoError(t, err)
	w.give(p, "tiny_guppy", 2)
	w.give(p, "common_grass", 1)
	w.give(p, "young_sapling", 1)
	job, err := svc.Start(w.ctx, p.ID, p.Account, "recipe_cook_grilled_guppy", 1)
	require.NoError(t, err)
	now = job.ReadyAt
	res, err := svc.Collect(w.ctx, p.ID, p.Account, job.ID)
	require.NoError(t, err)
	require.Equal(t, 1, res.Made)
	require.Equal(t, 2, w.count(p, "food_grilled_guppy"))
	var effect []byte
	require.NoError(t, w.p.QueryRow(w.ctx, `SELECT use_effect FROM content.item_template WHERE code = 'food_grilled_guppy'`).Scan(&effect))
	require.Contains(t, string(effect), "heal_pct") // a share of max HP (DD-274)
}
