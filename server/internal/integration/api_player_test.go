//go:build integration

package integration

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

// DD-370: what the character is doing, from whichever action table holds it.
func TestHTTPCurrentAction(t *testing.T) {
	w := newWorld(t)
	a, other := newHTTPPlayer(t, "warrior"), newHTTPPlayer(t, "mage")
	require.Nil(t, a.ok("GET", a.path("/action"), nil)["action"], "idle")

	a.ok("POST", a.path("/activity"), map[string]any{"kind": "mining", "duration_seconds": 3600})
	act := a.ok("GET", a.path("/action"), nil)["action"].(map[string]any)
	require.Equal(t, "gathering", act["kind"])
	require.Equal(t, "mining", act["detail"])
	require.NotEmpty(t, act["started_at"])
	require.NotEmpty(t, act["ends_at"])

	// A raid lobby has no times yet.
	run := uuid.New()
	_, err := w.p.Exec(w.ctx, `INSERT INTO raid.run (id, party_id, raid_code, leader_character_id, state, duration_seconds)
		VALUES ($1, gen_random_uuid(), 'north_keep_siege', $2, 'forming', 3600)`, run, other.char)
	require.NoError(t, err)
	_, err = w.p.Exec(w.ctx, `INSERT INTO raid.member (id, run_id, character_id, slot, snapshot) VALUES (gen_random_uuid(), $1, $2, 0, '{}')`, run, other.char)
	require.NoError(t, err)
	lobby := other.ok("GET", other.path("/action"), nil)["action"].(map[string]any)
	require.Equal(t, "raid", lobby["kind"])
	require.Nil(t, lobby["ends_at"])

	// Someone else's character is not mine to ask about.
	a.fail("GET", "/v1/character/"+other.char+"/action", nil, 404, "")
}

// DD-371: the public page shows the character, never the account.
func TestHTTPPlayerProfile(t *testing.T) {
	w := newWorld(t)
	a, b := newHTTPPlayer(t, "warrior"), newHTTPPlayer(t, "mage")
	id := uuid.MustParse(a.char)
	in := w.give(player{ID: id}, "sword_iron", 1)
	_, err := w.p.Exec(w.ctx, `UPDATE character.character SET level = 60 WHERE id = $1`, id)
	require.NoError(t, err)
	a.ok("POST", a.path("/equip/"+in.ID.String()), map[string]any{})
	_, err = w.p.Exec(w.ctx, `INSERT INTO crafting.character_profession (character_id, profession_code, level, xp) VALUES ($1, 'mining', 12, 0)
		ON CONFLICT (character_id, profession_code) DO UPDATE SET level = 12`, id)
	require.NoError(t, err)
	_, err = w.p.Exec(w.ctx, `INSERT INTO cosmetic.profile (account_id, avatar, background)
		SELECT account_id, 'av07', 'bg03' FROM character.character WHERE id = $1`, id)
	require.NoError(t, err)
	_, err = w.p.Exec(w.ctx, `INSERT INTO pet.pet (id, character_id, species_code, level, hatch_ready_at, active)
		VALUES (gen_random_uuid(), $1, 'flamewing', 5, now(), true)`, id)
	require.NoError(t, err)

	p := b.ok("GET", "/v1/players/"+a.name, nil)
	require.Equal(t, a.name, p["name"])
	require.Equal(t, "warrior", p["class_code"])
	require.Equal(t, "av07", p["avatar"])
	require.Equal(t, "bg03", p["background"])
	gear := p["equipment"].(map[string]any)["mainhand"].(map[string]any)
	require.Equal(t, "sword_iron", gear["template_code"])
	require.NotContains(t, gear, "instance_id")
	require.Equal(t, "mining", items(p, "professions")[0]["code"])
	require.EqualValues(t, 12, items(p, "professions")[0]["level"])
	require.EqualValues(t, 5, p["pet"].(map[string]any)["level"])
	for _, k := range []string{"account_id", "email", "username"} {
		require.NotContains(t, p, k)
	}

	// Names match as everywhere else: case does not matter.
	require.Equal(t, a.name, b.ok("GET", "/v1/players/"+"h"+a.name[1:], nil)["name"])
	b.fail("GET", "/v1/players/Nobody"+randHex(3), nil, 404, "E_NOT_FOUND")
	(&hc{t: t, base: a.base}).fail("GET", "/v1/players/"+a.name, nil, 401, "")

	// A locked account is not shown.
	_, err = w.p.Exec(w.ctx, `UPDATE identity.account SET state = 'locked' WHERE id = (SELECT account_id FROM character.character WHERE id = $1)`, id)
	require.NoError(t, err)
	b.fail("GET", "/v1/players/"+a.name, nil, 404, "E_NOT_FOUND")
}
