//go:build integration

package integration

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// Party XP bonus (DD-305): two party members fishing at the same time each
// keep their own XP plus 3% per member working (+6%); nobody gives XP away.
func TestHTTPPartyXPBonus(t *testing.T) {
	w := newWorld(t)
	a, b := newHTTPPlayer(t, "warrior"), newHTTPPlayer(t, "mage")
	inv := a.ok("POST", a.path("/party/invite"), map[string]any{"target_name": b.name})["id"].(string)
	b.ok("POST", b.path("/party/invites/"+inv+"/accept"), map[string]any{})

	start := func(p *hc) string {
		run := p.ok("POST", p.path("/activity"), map[string]any{"kind": "fishing", "duration_seconds": 3600})
		id := run["contract"].(map[string]any)["id"].(string)
		// Both ran the same full hour.
		_, err := w.p.Exec(w.ctx, `UPDATE activity.contract SET started_at = now() - interval '61 minutes',
			expires_at = now() - interval '1 minute' WHERE id = $1`, id)
		require.NoError(t, err)
		return id
	}
	fishingXP := func(p *hc) int64 {
		var xp, level int64
		err := w.p.QueryRow(w.ctx, `SELECT xp, level FROM crafting.character_profession
			WHERE character_id = $1 AND profession_code = 'fishing'`, p.char).Scan(&xp, &level)
		if err != nil {
			return 0
		}
		return xp + level*1_000_000 // order-preserving: a level-up still counts as more
	}
	bonus := func(res map[string]any) map[string]any {
		share := res["party_share"].(map[string]any)
		earned, kept := share["earned"].(float64), share["kept"].(float64)
		require.Positive(t, earned)
		require.InDelta(t, earned*1.06, kept, 1, "two members: +6%")
		return share
	}
	ra, rb := start(a), start(b)
	// The party shows what each member is doing (DD-301).
	for _, m := range a.ok("GET", a.path("/party"), nil)["party"].(map[string]any)["members"].([]any) {
		act := m.(map[string]any)["activity"].(map[string]any)
		require.Equal(t, "gathering", act["kind"])
		require.Equal(t, "fishing", act["detail"])
	}

	share := bonus(a.ok("POST", "/v1/activity/"+ra+"/reveal", map[string]any{}))
	partners := share["partners"].([]any)
	require.Len(t, partners, 1)
	require.Equal(t, b.char, partners[0].(map[string]any)["character_id"])
	require.Zero(t, fishingXP(b), "nothing is given to the partner")

	before := fishingXP(a)
	bonus(b.ok("POST", "/v1/activity/"+rb+"/reveal", map[string]any{}))
	require.Equal(t, before, fishingXP(a), "B's haul gives A nothing either")

	// Hunts get the bonus on character XP the same way.
	hunt := func(p *hc) string {
		h := p.ok("POST", p.path("/hunt"), map[string]any{"mob_code": "young_boar", "duration_seconds": 600})
		id := h["contract"].(map[string]any)["id"].(string)
		_, err := w.p.Exec(w.ctx, `UPDATE hunt.contract SET started_at = now() - interval '11 minutes',
			expires_at = now() - interval '1 minute' WHERE id = $1`, id)
		require.NoError(t, err)
		return id
	}
	ha := hunt(a)
	hunt(b)
	hres := a.ok("POST", "/v1/hunt/"+ha+"/reveal", map[string]any{})
	require.NotNil(t, hres["party_share"], "%v", hres)
	bonus(hres)

	// Someone outside the party gets no bonus.
	c := newHTTPPlayer(t, "rogue")
	rc := start(c)
	require.Nil(t, c.ok("POST", "/v1/activity/"+rc+"/reveal", map[string]any{})["party_share"])
}

// Raid finding (DD-299): the queue forms a random party with its lobby; the
// board lets anyone join an advertising party.
func TestHTTPRaidFinding(t *testing.T) {
	w := newWorld(t)
	const code = "north_keep_siege" // level 100
	ps := make([]*hc, 9)
	for i := range ps {
		ps[i] = newHTTPPlayer(t, []string{"warrior", "cleric", "mage", "rogue", "ranger"}[i%5])
		_, err := w.p.Exec(w.ctx, `UPDATE character.character SET level = 150 WHERE id = $1`, ps[i].char)
		require.NoError(t, err)
	}
	low := newHTTPPlayer(t, "warrior")
	q := func(p *hc) string { return p.path("/raid/queue") }
	// Queueing and advertising happen in the raid's region (DD-302).
	ps[0].fail("POST", q(ps[0]), map[string]any{"raid_code": code}, 409, "E_DUNGEON_REGION")
	for _, p := range ps {
		w.placeAt(t, "northern_keep_ruins", p.char)
	}
	cur := func(p *hc) map[string]any { return p.ok("GET", p.path("/raid"), nil) }

	// Queue: level gate, unknown raid, not from a party.
	low.fail("POST", q(low), map[string]any{"raid_code": code}, 403, "E_LEVEL_REQUIREMENT")
	ps[0].fail("POST", q(ps[0]), map[string]any{"raid_code": "nope"}, 404, "E_UNKNOWN_RAID")
	ps[0].fail("POST", q(ps[0]), "{", 400, "E_INVALID_JSON")
	ps[0].fail("POST", ps[0].path("/raid/queue/leave"), map[string]any{}, 404, "E_NOT_QUEUED")

	// Two wait; nothing forms until the first has waited five minutes.
	st := ps[0].ok("POST", q(ps[0]), map[string]any{"raid_code": code})["queue"].(map[string]any)
	require.Equal(t, 1.0, st["waiting"])
	st = ps[1].ok("POST", q(ps[1]), map[string]any{"raid_code": code})["queue"].(map[string]any)
	require.Equal(t, 2.0, st["waiting"])
	require.Equal(t, false, cur(ps[1])["in_party"])
	_, err := w.p.Exec(w.ctx, `UPDATE raid.queue SET queued_at = now() - interval '6 minutes' WHERE character_id = $1`, ps[0].char)
	require.NoError(t, err)
	got := cur(ps[1])
	require.Equal(t, true, got["in_party"], "a look at the queue formed the party")
	require.Nil(t, got["queue"])
	raid := got["raid"].(map[string]any)
	require.Len(t, raid["members"], 2, "both are in the lobby")
	require.Equal(t, ps[0].char, raid["leader_character_id"], "the longest waiting leads")
	ps[0].fail("POST", q(ps[0]), map[string]any{"raid_code": code}, 409, "E_ALREADY_IN_PARTY")
	ps[0].ok("POST", ps[0].path("/raid/"+raid["id"].(string)+"/start"), map[string]any{})

	// Five waiting form a party at once.
	for _, p := range ps[2:6] {
		p.ok("POST", q(p), map[string]any{"raid_code": code})
	}
	require.Equal(t, 4.0, cur(ps[5])["queue"].(map[string]any)["waiting"])
	ps[6].ok("POST", q(ps[6]), map[string]any{"raid_code": code})
	full := cur(ps[6])["raid"].(map[string]any)
	require.Len(t, full["members"], 5)

	// The board: a lone leader advertises, others join the party (and its lobby).
	lead, j1, j2 := ps[7], ps[8], newHTTPPlayer(t, "cleric")
	_, err = w.p.Exec(w.ctx, `UPDATE character.character SET level = 150 WHERE id = $1`, j2.char)
	require.NoError(t, err)
	l := lead.ok("POST", lead.path("/raid/listing"), map[string]any{"raid_code": code})["listing"].(map[string]any)
	lid := l["id"].(string)
	require.Len(t, l["members"], 1)
	require.Equal(t, lid, cur(lead)["listing"].(map[string]any)["id"])
	board := items(j1.ok("GET", j1.path("/raid/listings?raid="+code), nil), "listings")
	require.Equal(t, lid, board[0]["id"], "newest first")
	low.fail("POST", low.path("/raid/listings/"+lid+"/join"), map[string]any{}, 403, "E_LEVEL_REQUIREMENT")
	j1.fail("POST", j1.path("/raid/listings/nope/join"), map[string]any{}, 400, "E_INVALID_ID")
	j1.ok("POST", j1.path("/raid/listings/"+lid+"/join"), map[string]any{})
	j1.fail("POST", j1.path("/raid/listing"), map[string]any{"raid_code": code}, 403, "E_NOT_LEADER")
	j1.fail("POST", j1.path("/raid/listing/close"), map[string]any{}, 403, "E_NOT_LEADER")
	ps[0].fail("POST", ps[0].path("/raid/listings/"+lid+"/join"), map[string]any{}, 409, "E_ALREADY_IN_PARTY")
	lobby := lead.ok("POST", lead.path("/raid"), map[string]any{"raid_code": code})["raid"].(map[string]any)
	j2.ok("POST", j2.path("/raid/listings/"+lid+"/join"), map[string]any{})
	require.Len(t, cur(j2)["raid"].(map[string]any)["members"], 2, "a joiner enters the open lobby")
	require.Equal(t, lobby["id"], cur(j2)["raid"].(map[string]any)["id"])

	// Closing takes it off the board; starting the raid does too.
	lead.ok("POST", lead.path("/raid/listing/close"), map[string]any{})
	lead.fail("POST", lead.path("/raid/listing/close"), map[string]any{}, 404, "E_NOT_FOUND")
	low.fail("POST", low.path("/raid/listings/"+lid+"/join"), map[string]any{}, 409, "E_LISTING_GONE")
	for _, x := range items(j1.ok("GET", j1.path("/raid/listings"), nil), "listings") {
		require.NotEqual(t, lid, x["id"])
	}
	lead.ok("POST", lead.path("/raid/listing"), map[string]any{"raid_code": code})
	lead.ok("POST", lead.path("/raid/"+lobby["id"].(string)+"/start"), map[string]any{})
	require.Nil(t, cur(lead)["listing"], "under way: off the board")
}
