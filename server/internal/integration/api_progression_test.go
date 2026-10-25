//go:build integration

package integration

import (
	"math"
	"testing"

	"github.com/stretchr/testify/require"
)

// Progression through the real router (DD-250): summary, dev XP across the
// level-100 gate, promotion, attributes, talents, events, public config, and
// an AFK hunt paying through the canonical XP grant.
func TestHTTPProgression(t *testing.T) {
	w := newWorld(t)
	p := newHTTPPlayer(t, "warrior")
	_, err := w.p.Exec(w.ctx, `UPDATE character.character SET level = 99 WHERE id = $1`, p.char)
	require.NoError(t, err)

	sum := p.ok("GET", p.path("/leveling"), nil)
	require.Equal(t, 99.0, sum["level"])
	require.Equal(t, 1000.0, sum["level_cap"])
	require.Equal(t, "warrior", sum["spec_code"])
	require.Equal(t, "promotion", sum["next_milestone"].(map[string]any)["kind"])
	require.Equal(t, "rage", sum["resource"].(map[string]any)["kind"])

	res := p.ok("POST", "/v1/dev/character/"+p.char+"/xp", map[string]any{"amount": 200_000})
	require.Equal(t, 100.0, res["level"])
	require.Equal(t, true, res["promotion_available"])
	require.Positive(t, res["overflow_xp"])

	sum = p.ok("GET", p.path("/leveling"), nil)
	promo := sum["promotion"].(map[string]any)
	require.Equal(t, true, promo["available"])
	require.Len(t, promo["choices"], 2)
	require.Equal(t, "veteran", sum["title"])

	p.fail("POST", p.path("/promotion"), map[string]any{"choice": "warlord"}, 409, "E_PROMOTION_UNAVAILABLE")
	p.fail("POST", p.path("/promotion"), map[string]any{}, 400, "E_VALIDATION")
	pr := p.ok("POST", p.path("/promotion"), map[string]any{"choice": "executioner"})
	require.Equal(t, "executioner", pr["spec_code"])
	require.Greater(t, pr["level"], 100.0, "held XP paid out")

	p.fail("POST", p.path("/attributes"), map[string]any{"strength": 9999}, 409, "E_NOT_ENOUGH_POINTS")
	p.fail("POST", p.path("/attributes"), map[string]any{"charisma": 1}, 400, "E_UNKNOWN_ATTRIBUTE")
	al := p.ok("POST", p.path("/attributes"), map[string]any{"strength": 3, "vitality": 2})
	require.Equal(t, 3.0, al["allocated"].(map[string]any)["strength"])

	// What one more point adds to a warrior (progression.yaml formulas; raw
	// values are still in the full-value bracket).
	gains := p.ok("GET", p.path("/leveling"), nil)["attributes"].(map[string]any)["gains"].(map[string]any)
	require.Equal(t, map[string]any{"attack": 2.2, "block": 0.03}, gains["strength"])
	require.Equal(t, map[string]any{"hp": 18.0, "defense": 1.2, "hp_regen": 0.05}, gains["vitality"])
	require.Equal(t, map[string]any{"attack": 0.4, "crit": 0.05, "hit": 0.5, "dodge": 0.04}, gains["dexterity"])
	require.Empty(t, gains["intelligence"], "nothing for a warrior")

	tal := p.ok("GET", p.path("/talents"), nil)
	require.Len(t, tal["trees"], 3)
	avail := int(tal["available"].(float64))
	require.GreaterOrEqual(t, avail, 1, "a talent point every 20 levels from 100")
	for i := 0; i < avail; i++ {
		p.ok("POST", p.path("/talents/warrior_blood_cruelty"), map[string]any{})
	}
	p.fail("POST", p.path("/talents/warrior_blood_fury"), map[string]any{}, 409, "E_NOT_ENOUGH_POINTS")
	p.fail("POST", p.path("/talents/nope"), map[string]any{}, 404, "E_UNKNOWN_TALENT")
	p.fail("POST", p.path("/talents/warrior_blood_frenzy"), map[string]any{}, 409, "")

	ev := items(p.ok("GET", p.path("/leveling/events?limit=200"), nil), "events")
	kinds := map[string]bool{}
	for _, e := range ev {
		kinds[e["kind"].(string)] = true
	}
	for _, k := range []string{"XP_GAINED", "LEVEL_UP", "PROMOTION_AVAILABLE", "PROMOTION_SELECTED", "TITLE_UNLOCKED",
		"TALENT_POINT_GAINED", "ATTRIBUTES_ALLOCATED", "TALENT_LEARNED", "PASSIVE_UNLOCKED"} {
		require.True(t, kinds[k], "event %s recorded", k)
	}

	cfg := p.ok("GET", "/v1/progression/config", nil)
	require.Len(t, cfg["xp_table"], 999)

	p.ok("POST", "/v1/dev/character/"+p.char+"/respec", map[string]any{"spec": "guardian"})
	p.fail("POST", "/v1/dev/character/"+p.char+"/respec", map[string]any{"spec": "iron_fortress"}, 409, "E_SPEC_UNAVAILABLE")
	p.fail("POST", "/v1/dev/character/"+p.char+"/xp", map[string]any{"amount": 0}, 400, "E_VALIDATION")
}

// An AFK hunt reveal pays through the canonical grant: one xp_grant row
// keyed by the contract, none added by a replay.
func TestHTTPHuntXPIsCanonical(t *testing.T) {
	w := newWorld(t)
	p := newHTTPPlayer(t, "warrior") // level 1, as the boar
	h := p.ok("POST", p.path("/hunt"), map[string]any{"mob_code": "young_boar", "duration_seconds": 3600})
	hid := h["contract"].(map[string]any)["id"].(string)
	_, err := w.p.Exec(w.ctx, `
		UPDATE hunt.contract SET started_at = now() - interval '61 minutes', expires_at = now() - interval '1 minute'
		 WHERE id = $1`, hid)
	require.NoError(t, err)
	res := p.ok("POST", "/v1/hunt/"+hid+"/reveal", map[string]any{})
	p.ok("POST", "/v1/hunt/"+hid+"/reveal", map[string]any{})
	var n int
	var amount int64
	require.NoError(t, w.p.QueryRow(w.ctx, `
		SELECT count(*), COALESCE(sum(amount), 0) FROM character.xp_grant
		 WHERE character_id = $1 AND source = 'hunt' AND source_id = $2`, p.char, hid).Scan(&n, &amount))
	require.Equal(t, 1, n)
	require.Equal(t, res["hunt"].(map[string]any)["total_xp"], float64(amount))
	victories := res["hunt"].(map[string]any)["victories"].(float64)
	require.Greater(t, victories, 30.0, "fights run back to back (DD-281)")
	require.Equal(t, victories*mobXP(p, "young_boar"), float64(amount), "level 1 vs 1: the mob's full XP a kill")
}

// mobXP is a hunted mob's own XP per kill from the public catalog.
func mobXP(p *hc, code string) float64 {
	for _, m := range items(p.ok("GET", "/v1/mobs", nil), "mobs") {
		if m["code"] == code {
			return m["xp_reward"].(float64)
		}
	}
	p.t.Fatalf("no mob %s", code)
	return 0
}

// DD-281: a mob far from the character's level pays a tenth of its XP.
func TestHTTPGrayMobPaysLess(t *testing.T) {
	w := newWorld(t)
	p := newHTTPPlayer(t, "warrior")
	_, err := w.p.Exec(w.ctx, `UPDATE character.character SET level = 40 WHERE id = $1`, p.char)
	require.NoError(t, err)
	mobs := p.ok("GET", "/v1/mobs", nil)
	require.Equal(t, map[string]any{"free_pct": 0.05, "power": 4.0, "floor": 0.1}, mobs["hunt_level_xp"])
	h := p.ok("POST", p.path("/hunt"), map[string]any{"mob_code": "young_boar", "duration_seconds": 3600})
	hid := h["contract"].(map[string]any)["id"].(string)
	_, err = w.p.Exec(w.ctx, `
		UPDATE hunt.contract SET started_at = now() - interval '61 minutes', expires_at = now() - interval '1 minute'
		 WHERE id = $1`, hid)
	require.NoError(t, err)
	res := p.ok("POST", "/v1/hunt/"+hid+"/reveal", map[string]any{})["hunt"].(map[string]any)
	require.Positive(t, res["victories"])
	require.Equal(t, res["victories"].(float64)*math.Round(mobXP(p, "young_boar")*0.1), res["total_xp"], "level 40 vs level 1: 10% of the mob's XP")
}
