//go:build integration

package integration

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/otter/server/internal/leaderboard"
)

func entries(t *testing.T, m map[string]any) []map[string]any {
	t.Helper()
	return items(m, "entries")
}

// DD-357: the score counts what a character achieved, once per dungeon; tables
// rank by their own measure; a locked account leaves them.
func TestHTTPLeaderboard(t *testing.T) {
	w := newWorld(t)
	svc := leaderboard.NewService(w.p)
	a, b := newHTTPPlayer(t, "warrior"), newHTTPPlayer(t, "mage")
	_, err := w.p.Exec(w.ctx, `UPDATE character.character SET level = 900, achievement_points = 40 WHERE id = $1`, a.char)
	require.NoError(t, err)
	_, err = w.p.Exec(w.ctx, `UPDATE character.character SET level = 901, achievement_points = 0 WHERE id = $1`, b.char)
	require.NoError(t, err)
	for range 2 { // the same dungeon twice counts once
		_, err = w.p.Exec(w.ctx, `INSERT INTO dungeon.contract (character_id, dungeon_code, duration_seconds, expires_at, seed,
			state, player_snapshot, reveal_result) VALUES ($1, 'crypt_of_ashes', 60, now(), 1, 'revealed', '{}', '{"cleared": true}')`, a.char)
		require.NoError(t, err)
	}
	_, err = w.p.Exec(w.ctx, `INSERT INTO dungeon.contract (character_id, dungeon_code, duration_seconds, expires_at, seed,
		state, player_snapshot, reveal_result) VALUES ($1, 'forest_hollow', 60, now(), 1, 'revealed', '{}', '{"cleared": false}')`, a.char)
	require.NoError(t, err)
	_, err = w.p.Exec(w.ctx, `INSERT INTO dungeon.weekly_lockout (character_id, dungeon_code, week_key, contract_id) VALUES ($1, 'north_keep_siege', '2026-01', gen_random_uuid())`, a.char)
	require.NoError(t, err)
	_, err = svc.Refresh(w.ctx)
	require.NoError(t, err)

	mine := a.ok("GET", "/v1/leaderboard?board=score&character_id="+a.char, nil)["me"].(map[string]any)
	require.EqualValues(t, 900*10+40*5+1*20+1*100, mine["score"])
	require.EqualValues(t, 1, mine["dungeons_cleared"], "the same dungeon counts once, a failed run not at all")
	require.EqualValues(t, 1, mine["raids_cleared"])

	// The shared test database holds other characters, so compare ranks.
	rank := func(c *hc, board string) float64 {
		me := c.ok("GET", "/v1/leaderboard?board="+board+"&character_id="+c.char, nil)["me"].(map[string]any)
		return me["rank"].(float64)
	}
	require.Less(t, rank(b, "level"), rank(a, "level"), "901 ranks above 900")
	require.Less(t, rank(a, "score"), rank(b, "score"), "more achieved beats one more level")

	// My rank on a table is the row the page shows there.
	r := int(rank(a, "score"))
	page := entries(t, a.ok("GET", "/v1/leaderboard?board=score&limit=1&offset="+strconv.Itoa(r-1), nil))
	require.Len(t, page, 1)
	require.Equal(t, a.name, page[0]["name"])
	require.EqualValues(t, r, page[0]["rank"])
	classRank := int(a.ok("GET", "/v1/leaderboard?board=score&class=warrior&character_id="+a.char, nil)["me"].(map[string]any)["rank"].(float64))
	page = entries(t, a.ok("GET", "/v1/leaderboard?board=score&class=warrior&limit=1&offset="+strconv.Itoa(classRank-1), nil))
	require.Equal(t, a.name, page[0]["name"])
	other := a.ok("GET", "/v1/leaderboard?board=score&class=mage&character_id="+a.char, nil)
	require.Nil(t, other["me"], "not on another class's table")

	a.fail("GET", "/v1/leaderboard?board=nope", nil, 400, "E_VALIDATION")
	a.fail("GET", "/v1/leaderboard?class=nope", nil, 400, "E_VALIDATION")
	a.fail("GET", "/v1/leaderboard?character_id=nope", nil, 400, "E_INVALID_ID")

	// A locked account leaves the table.
	_, err = w.p.Exec(w.ctx, `UPDATE identity.account SET state = 'locked' WHERE id = (SELECT account_id FROM character.character WHERE id = $1)`, b.char)
	require.NoError(t, err)
	_, err = svc.Refresh(w.ctx)
	require.NoError(t, err)
	require.Nil(t, b.ok("GET", "/v1/leaderboard?character_id="+b.char, nil)["me"])
}
