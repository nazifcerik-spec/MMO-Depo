//go:build integration

package integration

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

// Pets (DD-328): an egg from the bag hatches after its species' time, the
// pet hunts on its own to level up, and the active pet adds its bonus.
func TestHTTPPets(t *testing.T) {
	w := newWorld(t)
	p := newHTTPPlayer(t, "warrior")
	me := player{ID: uuid.MustParse(p.char)}

	species := p.ok("GET", "/v1/pets", nil)["species"].([]any)
	require.Len(t, species, 20)
	require.Equal(t, "eternal_kirin", species[0].(map[string]any)["code"], "best first")

	w.give(me, "pet_egg_eternal_kirin", 2)
	got := p.ok("GET", p.path("/pets"), nil)
	require.Empty(t, got["pets"])
	eggs := got["eggs"].([]any)
	require.Len(t, eggs, 1)
	egg := eggs[0].(map[string]any)
	require.Equal(t, "eternal_kirin", egg["species"])
	require.EqualValues(t, 2, egg["quantity"])

	// Hatching takes one egg; the pet is in it for a day.
	pt := p.ok("POST", p.path("/pets/hatch"), map[string]any{"item_id": egg["item_id"]})["pet"].(map[string]any)
	require.Equal(t, false, pt["hatched"])
	require.Equal(t, 1, w.count(me, "pet_egg_eternal_kirin"))
	id := pt["id"].(string)
	p.fail("POST", p.path("/pets/"+id+"/active"), map[string]any{"active": true}, 409, "E_PET_HATCHING")
	p.fail("POST", p.path("/pets/"+id+"/hunt"), map[string]any{"mob_code": "marsh_lurker", "duration_seconds": 3600},
		409, "E_PET_HATCHING")
	p.fail("POST", p.path("/pets/hatch"), map[string]any{"item_id": uuid.NewString()}, 409, "E_NOT_EGG")

	_, err := w.p.Exec(w.ctx, `UPDATE pet.pet SET hatch_ready_at = now() - interval '1 second', level = 100
		WHERE id = $1`, id)
	require.NoError(t, err)

	// Active: its bonus (12% defense at level 100) joins the sheet.
	defense := func() float64 {
		return p.ok("GET", p.path("/sheet"), nil)["derived"].(map[string]any)["defense"].(float64)
	}
	before := defense()
	p.ok("POST", p.path("/pets/"+id+"/active"), map[string]any{"active": true})
	require.Greater(t, defense(), before)
	p.ok("POST", p.path("/pets/"+id+"/active"), map[string]any{"active": false})
	require.Equal(t, before, defense())
	_, err = w.p.Exec(w.ctx, `UPDATE pet.pet SET level = 1 WHERE id = $1`, id)
	require.NoError(t, err)

	// A hunt: raid bosses and mobs past its reach are refused.
	p.fail("POST", p.path("/pets/"+id+"/hunt"), map[string]any{"mob_code": "aethon_unbound", "duration_seconds": 3600},
		400, "E_VALIDATION")
	p.fail("POST", p.path("/pets/"+id+"/collect"), map[string]any{}, 409, "E_PET_IDLE")
	p.ok("POST", p.path("/pets/"+id+"/hunt"), map[string]any{"mob_code": "marsh_lurker", "duration_seconds": 3600})
	p.fail("POST", p.path("/pets/"+id+"/hunt"), map[string]any{"mob_code": "marsh_lurker", "duration_seconds": 3600},
		409, "E_PET_HUNTING")
	_, err = w.p.Exec(w.ctx, `UPDATE pet.pet SET hunt_started_at = now() - interval '61 minutes',
		hunt_ready_at = now() - interval '1 minute' WHERE id = $1`, id)
	require.NoError(t, err)
	res := p.ok("POST", p.path("/pets/"+id+"/collect"), map[string]any{})
	require.Positive(t, res["victories"])
	require.Positive(t, res["xp"])
	require.Equal(t, true, res["leveled_up"], "an hour on level 10 mobs lifts a level 1 pet")
	require.Greater(t, res["pet"].(map[string]any)["level"].(float64), 1.0)
	require.Nil(t, res["pet"].(map[string]any)["hunt"])
}
