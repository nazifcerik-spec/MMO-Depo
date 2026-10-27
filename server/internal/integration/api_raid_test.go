//go:build integration

package integration

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// Raids (DD-261): the party leader opens a lobby, members join, the leader
// starts it; once it has run, every member collects their own rewards once,
// and a clear locks the raid for the week.
func TestHTTPRaidFlow(t *testing.T) {
	w := newWorld(t)
	a, b, c := newHTTPPlayer(t, "warrior"), newHTTPPlayer(t, "cleric"), newHTTPPlayer(t, "mage")
	outsider := newHTTPPlayer(t, "rogue")
	for _, p := range []*hc{a, b, c, outsider} { // far above the level-100 raid: a sure clear
		_, err := w.p.Exec(w.ctx, `UPDATE character.character SET level = 400 WHERE id = $1`, p.char)
		require.NoError(t, err)
	}
	const code = "north_keep_siege"

	cat := items(a.ok("GET", "/v1/raids", nil), "raids")
	require.Len(t, cat, 6)
	require.Equal(t, code, cat[0]["code"])
	require.Equal(t, 2.0, cat[0]["group_min"])
	for _, d := range items(a.ok("GET", "/v1/dungeon/catalog", nil), "dungeons") {
		require.NotEqual(t, code, d["code"], "raids are not solo dungeons")
	}
	a.fail("POST", a.path("/dungeon"), map[string]any{"dungeon_code": code}, 404, "E_UNKNOWN_DUNGEON")
	for _, m := range items(a.ok("GET", "/v1/mobs", nil), "mobs") {
		require.NotEqual(t, "varn_undying", m["code"], "raid bosses are not hunted")
	}
	a.fail("POST", a.path("/hunt"), map[string]any{"mob_code": "varn_undying", "duration_seconds": 3600}, 400, "E_UNKNOWN_MOB")

	// A party first; only its leader opens a lobby.
	require.Equal(t, false, a.ok("GET", a.path("/raid"), nil)["in_party"])
	a.fail("POST", a.path("/raid"), map[string]any{"raid_code": code}, 409, "E_NOT_IN_PARTY")
	a.ok("POST", a.path("/party"), map[string]any{})
	for _, p := range []*hc{b, c} {
		a.ok("POST", a.path("/party/invite"), map[string]any{"target_name": p.name})
		inv := items(p.ok("GET", p.path("/party/invites"), nil), "invites")[0]["id"].(string)
		p.ok("POST", p.path("/party/invites/"+inv+"/accept"), map[string]any{})
	}
	b.fail("POST", b.path("/raid"), map[string]any{"raid_code": code}, 403, "E_NOT_LEADER")
	a.fail("POST", a.path("/raid"), map[string]any{"raid_code": "crypt_of_ashes"}, 404, "E_UNKNOWN_RAID")

	// Calling a lobby off frees the party for a new one.
	// Opened only where the raid is (DD-302).
	a.fail("POST", a.path("/raid"), map[string]any{"raid_code": code}, 409, "E_DUNGEON_REGION")
	w.placeAt(t, "northern_keep_ruins", a.char)
	first := a.ok("POST", a.path("/raid"), map[string]any{"raid_code": code})["raid"].(map[string]any)["id"].(string)
	a.fail("POST", a.path("/raid"), map[string]any{"raid_code": code}, 409, "")
	a.ok("POST", a.path("/raid/"+first+"/leave"), map[string]any{})
	b.fail("POST", b.path("/raid/"+first+"/join"), map[string]any{}, 409, "E_RAID_NOT_OPEN")

	raid := a.ok("POST", a.path("/raid"), map[string]any{"raid_code": code})["raid"].(map[string]any)
	id := raid["id"].(string)
	require.Len(t, raid["members"], 1)
	a.fail("POST", a.path("/raid/"+id+"/start"), map[string]any{}, 409, "E_RAID_NOT_ENOUGH")
	lobby := b.ok("GET", b.path("/raid"), nil)
	require.Nil(t, lobby["raid"])
	require.Equal(t, id, lobby["lobby"].(map[string]any)["id"])
	b.ok("POST", b.path("/raid/"+id+"/join"), map[string]any{})
	b.fail("POST", b.path("/raid/"+id+"/join"), map[string]any{}, 409, "E_ALREADY_IN_RAID")
	outsider.fail("POST", outsider.path("/raid/"+id+"/join"), map[string]any{}, 409, "E_NOT_IN_PARTY")
	c.ok("POST", c.path("/raid/"+id+"/join"), map[string]any{})
	c.ok("POST", c.path("/raid/"+id+"/leave"), map[string]any{})
	c.fail("POST", c.path("/raid/"+id+"/leave"), map[string]any{}, 404, "E_NOT_FOUND")
	c.ok("POST", c.path("/raid/"+id+"/join"), map[string]any{})
	b.fail("POST", b.path("/raid/"+id+"/start"), map[string]any{}, 403, "E_NOT_LEADER")
	b.fail("POST", b.path("/raid/"+id+"/reveal"), map[string]any{}, 409, "E_RAID_NOT_RUNNING")

	started := a.ok("POST", a.path("/raid/"+id+"/start"), map[string]any{})["raid"].(map[string]any)
	require.Equal(t, "running", started["state"])
	members := started["members"].([]any)
	require.Len(t, members, 3)
	require.Equal(t, "p3", members[2].(map[string]any)["label"])
	c.fail("POST", c.path("/raid/"+id+"/leave"), map[string]any{}, 409, "E_RAID_NOT_OPEN")
	a.fail("POST", a.path("/raid/"+id+"/reveal"), map[string]any{}, 409, "E_RAID_NOT_READY")
	outsider.fail("POST", outsider.path("/raid/"+id+"/reveal"), map[string]any{}, 404, "E_NOT_FOUND")

	_, err := w.p.Exec(w.ctx, `UPDATE raid.run SET started_at = now() - interval '2 hours', expires_at = now() - interval '1 hour'
		WHERE id = $1`, id)
	require.NoError(t, err)

	var log []any
	for _, p := range []*hc{a, b, c} {
		res := p.ok("POST", p.path("/raid/"+id+"/reveal"), map[string]any{})
		result := res["result"].(map[string]any)
		require.Equal(t, true, result["cleared"], "%v", result)
		if log == nil {
			log = result["encounters"].([]any)
		}
		require.Equal(t, log, result["encounters"], "every member sees the same raid")
		reward := res["reward"].(map[string]any)
		require.Positive(t, reward["gold_copper"].(float64))
		gear := 0
		for _, it := range items(res, "items") {
			var slot *string
			require.NoError(t, w.p.QueryRow(w.ctx, `SELECT equip_slot FROM content.item_template WHERE code = $1`,
				it["template_code"]).Scan(&slot))
			if slot != nil {
				gear++
			}
		}
		require.GreaterOrEqual(t, gear, 2, "the boss drops two gear pieces")

		again := p.ok("POST", p.path("/raid/"+id+"/reveal"), map[string]any{})
		require.Equal(t, true, again["replayed"])
		require.Equal(t, reward["gold_copper"], again["reward"].(map[string]any)["gold_copper"])
		var paid int64
		require.NoError(t, w.p.QueryRow(w.ctx, `
			SELECT coalesce(sum(e.delta), 0) FROM economy.ledger_entry e
			  JOIN raid.member m ON e.idempotency_key = 'raid-reveal-' || m.id::text
			 WHERE m.run_id = $1 AND m.character_id = $2 AND e.account_ref = $2`, id, p.char).Scan(&paid))
		require.Equal(t, int64(reward["gold_copper"].(float64)), paid, "paid once")
	}

	cur := a.ok("GET", a.path("/raid"), nil)
	require.Nil(t, cur["raid"], "collected: out of the raid")
	require.Equal(t, []any{code}, cur["lockouts"])
	a.fail("POST", a.path("/raid"), map[string]any{"raid_code": code}, 409, "E_RAID_WEEKLY_LOCKED")
}
