//go:build integration

package integration

import (
	"testing"

	"github.com/stretchr/testify/require"
)

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
