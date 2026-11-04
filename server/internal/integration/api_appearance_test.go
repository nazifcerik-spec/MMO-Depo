//go:build integration

package integration

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// DD-349: a race starts only as its classes, and its passives reach the sheet.
func TestHTTPRaceClassesAndPassives(t *testing.T) {
	w := newWorld(t)
	p := newHTTPPlayer(t, "warrior")
	create := func(race, class string) map[string]any {
		return map[string]any{"name": "R" + randHex(5), "race_code": race, "class_code": class}
	}
	p.fail("POST", "/v1/character", create("human", "shaman"), 400, "E_VALIDATION")
	p.fail("POST", "/v1/character", create("gnome", "warrior"), 400, "E_VALIDATION")
	orc := p.ok("POST", "/v1/character", create("orc", "warrior"))

	races := p.ok("GET", "/v1/classes", nil)["races"].(map[string]any)
	require.Len(t, races, 10)
	require.Contains(t, races["orc"].(map[string]any)["classes"], "shaman")

	// The same points give an orc (+2%) more strength than a human (+1%).
	str := func(id string) float64 {
		_, err := w.p.Exec(w.ctx, `UPDATE character.character SET attr_str = 1000 WHERE id = $1`, id)
		require.NoError(t, err)
		sheet := p.ok("GET", "/v1/character/"+id+"/sheet", nil)
		return sheet["primary"].(map[string]any)["str"].(float64)
	}
	require.Greater(t, str(orc["id"].(string)), str(p.char))
}

// DD-346: the look is saved at creation and stays editable through the API.
func TestHTTPAppearance(t *testing.T) {
	p := newHTTPPlayer(t, "warrior")
	sheet := p.ok("GET", p.path("/sheet"), nil)
	require.Equal(t, "male", sheet["body"])
	require.EqualValues(t, 1, sheet["face_id"])

	got := p.ok("POST", p.path("/appearance"), map[string]any{"face_id": 3, "hair_style_id": 4, "hair_color_id": 2, "eye_color_id": 5})
	require.EqualValues(t, 3, got["face_id"])
	sheet = p.ok("GET", p.path("/sheet"), nil)
	require.EqualValues(t, 4, sheet["hair_style_id"])
	require.EqualValues(t, 2, sheet["hair_color_id"])
	require.EqualValues(t, 5, sheet["eye_color_id"])

	got = p.ok("POST", p.path("/appearance"), map[string]any{"body": "female"})
	require.Equal(t, "female", got["body"])
	require.EqualValues(t, 1, got["face_id"], "faces belong to a body")
	require.EqualValues(t, 2, got["hair_color_id"], "colours carry over")

	p.fail("POST", p.path("/appearance"), map[string]any{"face_id": 99}, 400, "E_VALIDATION")
}
