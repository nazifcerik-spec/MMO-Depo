//go:build integration

package integration

import (
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
	p := newHTTPPlayer(t, "warrior")
	_, err := w.p.Exec(w.ctx, `UPDATE character.character SET level = 6 WHERE id = $1`, p.char)
	require.NoError(t, err)
	h := p.ok("POST", p.path("/hunt"), map[string]any{"mob_code": "young_boar", "duration_seconds": 3600})
	hid := h["contract"].(map[string]any)["id"].(string)
	_, err = w.p.Exec(w.ctx, `
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
	require.Equal(t, victories*2000, float64(amount), "AFK rate: 1000 XP per minute of 2-minute encounters")
}

// DD-254: a mob far below the character pays a tenth of the AFK rate.
func TestHTTPGrayMobPaysLess(t *testing.T) {
	w := newWorld(t)
	p := newHTTPPlayer(t, "warrior")
	_, err := w.p.Exec(w.ctx, `UPDATE character.character SET level = 40 WHERE id = $1`, p.char)
	require.NoError(t, err)
	mobs := p.ok("GET", "/v1/mobs", nil)
	require.Len(t, mobs["xp_bands"], 3)
	h := p.ok("POST", p.path("/hunt"), map[string]any{"mob_code": "young_boar", "duration_seconds": 3600})
	hid := h["contract"].(map[string]any)["id"].(string)
	_, err = w.p.Exec(w.ctx, `
		UPDATE hunt.contract SET started_at = now() - interval '61 minutes', expires_at = now() - interval '1 minute'
		 WHERE id = $1`, hid)
	require.NoError(t, err)
	res := p.ok("POST", "/v1/hunt/"+hid+"/reveal", map[string]any{})["hunt"].(map[string]any)
	require.Positive(t, res["victories"])
	require.Equal(t, res["victories"].(float64)*200, res["total_xp"], "level 40 vs level 1: 10% of 2000 per victory")
}
