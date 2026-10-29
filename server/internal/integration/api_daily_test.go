//go:build integration

package integration

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// Daily login reward (DD-311): once per account and day, on a 29-day cycle (DD-322).
func TestHTTPDailyReward(t *testing.T) {
	w := newWorld(t)
	p := newHTTPPlayer(t, "mage")
	st := p.ok("GET", p.path("/daily-reward"), nil)
	require.Equal(t, true, st["claimable"])
	require.EqualValues(t, 1, st["day"])
	require.Len(t, st["rewards"], 29)

	before := p.gold()
	res := p.ok("POST", p.path("/daily-reward/claim"), map[string]any{})
	rw := res["reward"].(map[string]any)
	require.EqualValues(t, 1, rw["day"])
	require.Equal(t, before+int64(rw["gold_copper"].(float64)), p.gold())
	require.Equal(t, false, res["status"].(map[string]any)["claimable"])
	p.fail("POST", p.path("/daily-reward/claim"), map[string]any{}, 409, "E_ALREADY_CLAIMED")

	// Claimed yesterday at day 2: today pays day 3 with potions.
	_, err := w.p.Exec(w.ctx, `UPDATE progression.daily_login SET last_day = current_date - 1, streak = 2
		WHERE account_id = (SELECT account_id FROM character.character WHERE id = $1)`, p.char)
	require.NoError(t, err)
	rw = p.ok("POST", p.path("/daily-reward/claim"), map[string]any{})["reward"].(map[string]any)
	require.EqualValues(t, 3, rw["day"])
	require.Equal(t, "potion_hp_small", rw["potion"])
	var n int
	require.NoError(t, w.p.QueryRow(w.ctx, `SELECT COALESCE(SUM(stack_quantity), 0) FROM inventory.item_instance
		WHERE owner_character_id = $1 AND template_code = 'potion_hp_small'`, p.char).Scan(&n))
	require.Equal(t, 3, n)
}
