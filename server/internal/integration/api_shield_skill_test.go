//go:build integration

package integration

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

// DD-364: a skill that needs a shield is only cast with one on. Without it the
// Guardian keeps the skill it would otherwise replace.
func TestHTTPShieldSkillNeedsAShield(t *testing.T) {
	w := newWorld(t)
	p := newHTTPPlayer(t, "warrior")
	me := player{ID: uuid.MustParse(p.char)}
	_, err := w.p.Exec(w.ctx, `UPDATE character.character SET level = 300, spec_code = 'guardian' WHERE id = $1`, me.ID)
	require.NoError(t, err)

	// The skills frozen into a hunt's start snapshot.
	hunt := func() string {
		h := p.ok("POST", p.path("/hunt"), map[string]any{"mob_code": "young_boar", "duration_seconds": 600})
		id := h["contract"].(map[string]any)["id"].(string)
		var snapshot string
		require.NoError(t, w.p.QueryRow(w.ctx, `SELECT player_snapshot::text FROM hunt.contract WHERE id = $1`, id).Scan(&snapshot))
		p.ok("POST", "/v1/hunt/"+id+"/stop", map[string]any{})
		return snapshot
	}
	tags := func() any { return p.ok("GET", p.path("/sheet"), nil)["equipment_tags"] }

	require.Nil(t, tags())
	bare := hunt()
	require.Contains(t, bare, "warrior_heavy_slash")
	require.NotContains(t, bare, "guardian_shield_slam")

	shield := w.give(me, "shield_sovereign", 1)
	p.ok("POST", p.path("/equip/"+shield.ID.String()), map[string]any{})
	require.Equal(t, []any{"shield"}, tags())
	armed := hunt()
	require.Contains(t, armed, "guardian_shield_slam")
	require.NotContains(t, armed, "warrior_heavy_slash", "the stronger skill of the job replaces it")

	p.ok("POST", p.path("/unequip/offhand"), map[string]any{})
	require.Nil(t, tags())
}
