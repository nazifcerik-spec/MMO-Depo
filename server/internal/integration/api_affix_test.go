//go:build integration

package integration

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

// Rarity and affixes make gear stronger (DD-275): the same sword adds more
// attack as an epic with "sharp" than as a common drop, and the bag shows the
// effective bonuses.
func TestHTTPRarityAndAffixesAddPower(t *testing.T) {
	w := newWorld(t)
	p := newHTTPPlayer(t, "warrior")
	me := player{ID: uuid.MustParse(p.char)}
	_, err := w.p.Exec(w.ctx, `UPDATE character.character SET level = 5 WHERE id = $1`, me.ID)
	require.NoError(t, err)
	attack := func() float64 {
		return p.ok("GET", p.path("/sheet"), nil)["derived"].(map[string]any)["attack"].(float64)
	}
	bare := attack()
	sword := w.give(me, "sword_iron", 1)
	p.ok("POST", p.path("/equip/"+sword.ID.String()), map[string]any{})
	common := attack()
	require.Greater(t, common, bare)

	_, err = w.p.Exec(w.ctx, `UPDATE inventory.item_instance SET rarity = 'epic', prefix_code = 'sharp' WHERE id = $1`, sword.ID)
	require.NoError(t, err)
	// sword_iron: +16 attack → ×1.2 = 19, sharp at level 5: round(2 + 0.5) = 3.
	require.Equal(t, common+6, attack())

	worn := p.ok("GET", p.path("/equipment"), nil)["slots"].(map[string]any)["mainhand"].(map[string]any)
	require.Equal(t, sword.ID.String(), worn["instance_id"])
	require.Len(t, worn["effective_bonuses"], 4, "attack, str, vit, sharp")
}
