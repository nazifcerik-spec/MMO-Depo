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

// Rings and charms are worn in pairs (DD-276): the first ring fills ring1,
// the second ring2, a third swaps ring1. Jewelry has no class.
func TestHTTPJewelryPairs(t *testing.T) {
	w := newWorld(t)
	p := newHTTPPlayer(t, "mage")
	me := player{ID: uuid.MustParse(p.char)}
	_, err := w.p.Exec(w.ctx, `UPDATE character.character SET level = 25 WHERE id = $1`, me.ID)
	require.NoError(t, err)
	slotOf := func(id uuid.UUID) string {
		for slot, it := range p.ok("GET", p.path("/equipment"), nil)["slots"].(map[string]any) {
			if it.(map[string]any)["instance_id"] == id.String() {
				return slot
			}
		}
		return ""
	}
	a, b, c := w.give(me, "ring_runed", 1), w.give(me, "ring_runed", 1), w.give(me, "ring_runed", 1)
	for _, r := range []uuid.UUID{a.ID, b.ID, c.ID} {
		p.ok("POST", p.path("/equip/"+r.String()), map[string]any{})
	}
	require.Equal(t, "", slotOf(a.ID), "the third ring swapped the first back to the bag")
	require.Equal(t, "ring2", slotOf(b.ID))
	require.Equal(t, "ring1", slotOf(c.ID))
	charm := w.give(me, "charm_runed", 1)
	p.ok("POST", p.path("/equip/"+charm.ID.String()), map[string]any{})
	require.Equal(t, "trinket1", slotOf(charm.ID))
	neck := w.give(me, "necklace_runed", 1)
	p.ok("POST", p.path("/equip/"+neck.ID.String()), map[string]any{})
	require.Equal(t, "necklace", slotOf(neck.ID))

	// Guardian shields need the branch.
	shield := w.give(me, "shield_royal", 1)
	_, err = w.p.Exec(w.ctx, `UPDATE character.character SET level = 110 WHERE id = $1`, me.ID)
	require.NoError(t, err)
	p.fail("POST", p.path("/equip/"+shield.ID.String()), map[string]any{}, 403, "")
}

// Worn set pieces add the set's bonuses (DD-277): two pieces of the Undying
// Plate add 4% HP on top of what each adds alone; the catalog is public.
func TestHTTPSetBonus(t *testing.T) {
	w := newWorld(t)
	p := newHTTPPlayer(t, "warrior")
	me := player{ID: uuid.MustParse(p.char)}
	_, err := w.p.Exec(w.ctx, `UPDATE character.character SET level = 100 WHERE id = $1`, me.ID)
	require.NoError(t, err)
	hp := func() float64 {
		return p.ok("GET", p.path("/sheet"), nil)["resources"].(map[string]any)["hp_max"].(float64)
	}
	bare := hp()
	helm, chest := w.give(me, "undying_helm", 1), w.give(me, "undying_chestplate", 1)
	p.ok("POST", p.path("/equip/"+helm.ID.String()), map[string]any{})
	withHelm := hp()
	p.ok("POST", p.path("/unequip/helm"), map[string]any{})
	p.ok("POST", p.path("/equip/"+chest.ID.String()), map[string]any{})
	withChest := hp()
	p.ok("POST", p.path("/equip/"+helm.ID.String()), map[string]any{})
	both := hp()
	alone := withHelm + withChest - bare // what the two add without the set
	require.InDelta(t, alone*1.04, both, 3, "2 pieces: +4%% HP (alone %v, both %v)", alone, both)

	sets := items(p.ok("GET", "/v1/item-sets", nil), "sets")
	require.Len(t, sets, 18)
}
