//go:build integration

package integration

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

// One action at a time (DD-263): whatever the character is doing blocks
// every other kind of action until it is collected or stopped.
func TestHTTPOneActionAtATime(t *testing.T) {
	w := newWorld(t)
	p := newHTTPPlayer(t, "warrior")
	_, err := w.p.Exec(w.ctx, `UPDATE character.character SET level = 10 WHERE id = $1`, p.char)
	require.NoError(t, err)
	busy := func(path string, body map[string]any, kind string) {
		t.Helper()
		code, out := p.call("POST", path, body)
		require.Equal(t, 409, code, "%s %v", path, out)
		e := out["error"].(map[string]any)
		require.Equal(t, "E_CHARACTER_BUSY", e["code"], "%s", path)
		require.Equal(t, kind, e["field"], "%s", path)
	}
	hunt := map[string]any{"mob_code": "young_boar", "duration_seconds": 3600}
	gather := map[string]any{"kind": "mining", "duration_seconds": 3600}
	dungeon := map[string]any{"dungeon_code": "crypt_of_ashes"}
	loc := p.ok("GET", p.path("/location"), nil)
	travel := map[string]any{"to_region_code": items(loc, "neighbors")[0]["code"]}

	h := p.ok("POST", p.path("/hunt"), hunt)["contract"].(map[string]any)["id"].(string)
	busy(p.path("/activity"), gather, "hunt")
	busy(p.path("/dungeon"), dungeon, "hunt")
	busy(p.path("/travel"), travel, "hunt")
	p.ok("POST", "/v1/hunt/"+h+"/stop", map[string]any{})

	d := p.ok("POST", p.path("/dungeon"), dungeon)["contract"].(map[string]any)["id"].(string)
	busy(p.path("/hunt"), hunt, "dungeon")
	busy(p.path("/activity"), gather, "dungeon")
	p.ok("POST", "/v1/dungeon/"+d+"/stop", map[string]any{})

	tr := p.ok("POST", p.path("/travel"), travel)["travel"].(map[string]any)["id"].(string)
	busy(p.path("/hunt"), hunt, "travel")
	p.ok("POST", "/v1/travel/"+tr+"/cancel", map[string]any{})

	// Different actions started at once: exactly one wins.
	starts := []struct {
		path string
		body map[string]any
	}{{p.path("/hunt"), hunt}, {p.path("/activity"), gather}, {p.path("/dungeon"), dungeon}, {p.path("/travel"), travel}}
	codes := make([]int, len(starts))
	var wg sync.WaitGroup
	for i, s := range starts {
		wg.Add(1)
		go func() {
			defer wg.Done()
			codes[i], _ = p.call("POST", s.path, s.body)
		}()
	}
	wg.Wait()
	won := 0
	for _, c := range codes {
		if c == 200 || c == 201 {
			won++
		} else {
			require.Equal(t, 409, c, "%v", codes)
		}
	}
	require.Equal(t, 1, won, "%v", codes)
	var n int
	require.NoError(t, w.p.QueryRow(w.ctx, `SELECT count(*) FROM character.active_action WHERE character_id = $1`, p.char).Scan(&n))
	require.Equal(t, 1, n)
}

// Solo dungeons have no weekly limit (DD-262), and a level gate keeps
// under-levelled characters out of dungeons and raids.
func TestHTTPDungeonsUnlimitedAndLevelGated(t *testing.T) {
	w := newWorld(t)
	p := newHTTPPlayer(t, "warrior")
	p.fail("POST", p.path("/dungeon"), map[string]any{"dungeon_code": "crypt_of_ashes"}, 403, "E_LEVEL_REQUIREMENT")
	_, err := w.p.Exec(w.ctx, `UPDATE character.character SET level = 25 WHERE id = $1`, p.char)
	require.NoError(t, err)
	for run := 0; run < 2; run++ {
		id := p.ok("POST", p.path("/dungeon"), map[string]any{"dungeon_code": "crypt_of_ashes"})["contract"].(map[string]any)["id"].(string)
		_, err := w.p.Exec(w.ctx, `UPDATE dungeon.contract SET started_at = now() - interval '2 hours',
			expires_at = now() - interval '1 hour' WHERE id = $1`, id)
		require.NoError(t, err)
		res := p.ok("POST", "/v1/dungeon/"+id+"/reveal", map[string]any{})
		require.Equal(t, true, res["run"].(map[string]any)["cleared"], "run %d", run)
		require.Equal(t, false, res["lockout_written"])
	}

	// Dungeons are entered where they are (DD-302): from the region, a linked
	// one, or the entrance node; the catalog says where and what a run pays.
	for _, d := range items(p.ok("GET", "/v1/dungeon/catalog", nil), "dungeons") {
		if d["code"] == "forgotten_mine" {
			require.Equal(t, "iron_mountains", d["region_code"])
			require.Positive(t, d["xp"].(float64))
		}
	}
	_, err = w.p.Exec(w.ctx, `UPDATE character.character SET level = 70 WHERE id = $1`, p.char)
	require.NoError(t, err)
	p.fail("POST", p.path("/dungeon"), map[string]any{"dungeon_code": "forgotten_mine"}, 409, "E_DUNGEON_REGION")
	w.placeAt(t, "forgotten_mine", p.char)
	require.Contains(t, p.ok("GET", p.path("/location"), nil)["dungeon_regions"], "iron_mountains")
	mine := p.ok("POST", p.path("/dungeon"), map[string]any{"dungeon_code": "forgotten_mine"})["contract"].(map[string]any)["id"].(string)
	p.ok("POST", "/v1/dungeon/"+mine+"/stop", map[string]any{})

	// Raids keep their level gate for every member.
	lead, low := newHTTPPlayer(t, "warrior"), newHTTPPlayer(t, "cleric")
	_, err = w.p.Exec(w.ctx, `UPDATE character.character SET level = 100 WHERE id = $1`, lead.char)
	require.NoError(t, err)
	lead.ok("POST", lead.path("/party"), map[string]any{})
	lead.ok("POST", lead.path("/party/invite"), map[string]any{"target_name": low.name})
	inv := items(low.ok("GET", low.path("/party/invites"), nil), "invites")[0]["id"].(string)
	low.ok("POST", low.path("/party/invites/"+inv+"/accept"), map[string]any{})
	w.placeAt(t, "northern_keep_ruins", lead.char) // the raid's region (DD-302)
	id := lead.ok("POST", lead.path("/raid"), map[string]any{"raid_code": "north_keep_siege"})["raid"].(map[string]any)["id"].(string)
	low.fail("POST", low.path("/raid/"+id+"/join"), map[string]any{}, 403, "E_LEVEL_REQUIREMENT")
}
