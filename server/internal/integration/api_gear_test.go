//go:build integration

package integration

import (
	"sync"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

// Equipped gear feeds the sheet; equipping into a filled slot sends the old
// piece back to the bag.
func TestHTTPGearAffectsSheet(t *testing.T) {
	w := newWorld(t)
	p := newHTTPPlayer(t, "warrior")
	me := player{ID: uuid.MustParse(p.char)}
	_, err := w.p.Exec(w.ctx, `UPDATE character.character SET level = 5 WHERE id = $1`, me.ID)
	require.NoError(t, err)
	stat := func(group, key string) float64 {
		return p.ok("GET", p.path("/sheet"), nil)[group].(map[string]any)[key].(float64)
	}
	attack, str := stat("derived", "attack"), stat("primary", "str")

	first := w.give(me, "sword_iron", 1).ID.String()
	second := w.give(me, "sword_iron", 1).ID.String()
	p.ok("POST", p.path("/equip/"+first), map[string]any{})
	require.Greater(t, stat("derived", "attack"), attack, "sword adds attack")
	require.Equal(t, str+3, stat("primary", "str"), "sword_iron: +3 str")

	p.ok("POST", p.path("/equip/"+second), map[string]any{})
	eq := p.ok("GET", p.path("/equipment"), nil)["slots"].(map[string]any)
	require.Equal(t, second, eq["mainhand"].(map[string]any)["instance_id"])
	require.NotEqual(t, -1.0, p.slotOf(first), "the replaced sword is back in the bag")
	require.Equal(t, str+3, stat("primary", "str"), "bonuses don't stack across swaps")

	p.ok("POST", p.path("/unequip/mainhand"), map[string]any{})
	require.Equal(t, attack, stat("derived", "attack"))
}

// A craft the character can't pay for consumes nothing.
func TestHTTPCraftNeedsGold(t *testing.T) {
	w := newWorld(t)
	p := newHTTPPlayer(t, "warrior")
	me := player{ID: uuid.MustParse(p.char)}
	p.ok("POST", p.path("/professions/learn"), map[string]any{"profession_code": "alchemy"})
	w.give(me, "linseed", 2)
	p.fail("POST", p.path("/craft"), map[string]any{"recipe_code": "recipe_potion_hp_small"}, 409, "E_INSUFFICIENT_GOLD")
	require.Equal(t, 2.0, p.stackOf("linseed"))
	p.mint(5)
	p.ok("POST", p.path("/craft"), map[string]any{"recipe_code": "recipe_potion_hp_small"})
	require.Equal(t, 0.0, p.stackOf("linseed"))
	profs := items(p.ok("GET", p.path("/professions"), nil), "professions")
	require.Equal(t, "alchemy", profs[0]["profession_code"])
}

// Concurrent reveals of one finished dungeon run pay out once and agree.
func TestHTTPDungeonConcurrentReveal(t *testing.T) {
	w := newWorld(t)
	p := newHTTPPlayer(t, "warrior")
	_, err := w.p.Exec(w.ctx, `UPDATE character.character SET level = 10 WHERE id = $1`, p.char)
	require.NoError(t, err)
	id := p.ok("POST", p.path("/dungeon"), map[string]any{"dungeon_code": "crypt_of_ashes"})["contract"].(map[string]any)["id"].(string)
	_, err = w.p.Exec(w.ctx, `
		UPDATE dungeon.contract SET started_at = now() - interval '2 hours', expires_at = now() - interval '1 hour'
		 WHERE id = $1`, id)
	require.NoError(t, err)

	const n = 6
	gold := make([]float64, n)
	codes := make([]int, n)
	var wg sync.WaitGroup
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			var out map[string]any
			codes[i], out = p.call("POST", "/v1/dungeon/"+id+"/reveal", map[string]any{})
			if run, ok := out["run"].(map[string]any); ok {
				gold[i] = run["total_gold_copper"].(float64)
			}
		}()
	}
	wg.Wait()
	for i := range codes {
		require.Equal(t, 200, codes[i], "reveal %d", i)
		require.Equal(t, gold[0], gold[i], "every reveal reports the same run")
	}
	var credits int
	require.NoError(t, w.p.QueryRow(w.ctx, `
		SELECT count(DISTINCT txn_id) FROM economy.ledger_entry WHERE idempotency_key = $1`, "dungeon-reveal-"+id).Scan(&credits))
	require.LessOrEqual(t, credits, 1, "paid at most once")
	var paid int64
	require.NoError(t, w.p.QueryRow(w.ctx, `SELECT coalesce(sum(delta), 0) FROM economy.ledger_entry
		 WHERE idempotency_key = $1 AND account_ref = $2`, "dungeon-reveal-"+id, p.char).Scan(&paid))
	require.Equal(t, int64(gold[0]), paid, "the reported gold is what was credited")
}
