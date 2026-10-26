//go:build integration

package integration

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

// The bag marks what the character can wear (DD-290): a warrior wears
// plate, not cloth or leather nor a mage's staff, and plate above its level
// only later.
func TestHTTPBagWearMarks(t *testing.T) {
	w := newWorld(t)
	p := newHTTPPlayer(t, "warrior")
	me := player{ID: uuid.MustParse(p.char)}
	for _, code := range []string{"helm_novice", "hood_novice", "helm_leather", "staff_willow", "helm_adamant", "potion_hp_small"} {
		w.give(me, code, 1)
	}
	wear := map[string]any{}
	for _, it := range items(p.ok("GET", p.path("/inventory"), nil), "items") {
		wear[it["template_code"].(string)] = it["wear"]
	}
	require.Equal(t, map[string]any{
		"helm_novice":     "ok",
		"hood_novice":     "class", // cloth
		"helm_leather":    "class",
		"staff_willow":    "class", // a mage's weapon
		"helm_adamant":    "level", // plate, level 50
		"potion_hp_small": nil,     // not equipment
	}, wear)
	// The equip refusals agree with the marks.
	p.fail("POST", p.path("/equip/"+p.bagItem("hood_novice")), map[string]any{}, 403, "")
	p.ok("POST", p.path("/equip/"+p.bagItem("helm_novice")), map[string]any{})
}
