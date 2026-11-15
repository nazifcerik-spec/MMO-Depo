//go:build integration

package integration

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

// DD-367 through the production router: the story walks scene by scene,
// scene quests are offered only while their scene is current and gate it,
// choices set flags, skip or reshape later scenes and grant reputation, and
// answers are final.
func TestHTTPStoryChoices(t *testing.T) {
	c := newHTTPPlayer(t, "warrior")
	db, ctx := pool(t)

	journal := func() (cur map[string]any, state map[string]any, chapters []map[string]any) {
		out := c.ok("GET", c.path("/story"), nil)
		scenes := items(out, "scenes")
		require.NotEmpty(t, scenes)
		return scenes[len(scenes)-1], out["state"].(map[string]any), items(out, "chapters")
	}
	advance := func(scene, option string) {
		c.ok("POST", c.path("/story/"+scene+"/advance"), map[string]any{"option": option})
	}
	offered := func(code string) bool {
		for _, q := range items(c.ok("GET", c.path("/quests"), nil), "available") {
			if q["code"] == code {
				return true
			}
		}
		return false
	}
	finishQuest := func(code string) { // objectives are proven by the quest suite
		_, err := db.Exec(ctx, `UPDATE progression.character_quest SET state = 'completed'
			WHERE character_id = $1 AND quest_code = $2`, c.char, code)
		require.NoError(t, err)
		c.ok("POST", c.path("/quests/"+code+"/turn-in"), map[string]any{})
	}
	pass := func(scenes map[string]string) { // jump ahead: scenes passed with these answers
		for sc, opt := range scenes {
			_, err := db.Exec(ctx, `INSERT INTO progression.story_scene (character_id, scene_code, option_code)
				VALUES ($1, $2, $3)`, c.char, sc, opt)
			require.NoError(t, err)
		}
	}
	setLevel := func(level int) {
		_, err := db.Exec(ctx, `UPDATE character.character SET level = $2 WHERE id = $1`, c.char, level)
		require.NoError(t, err)
	}

	cur, _, chapters := journal()
	require.Equal(t, "p_star", cur["code"])
	require.Equal(t, "current", cur["state"])
	require.Empty(t, cur["extras"], "the dwarf/undead/big lines follow the race")
	require.Equal(t, "open", chapters[0]["state"])
	require.Equal(t, "locked", chapters[1]["state"])
	require.False(t, offered("story_p_bread"), "its scene is not current yet")
	c.fail("POST", c.path("/quests/story_p_bread/accept"), map[string]any{}, 409, "E_QUEST_LOCKED")
	c.fail("POST", c.path("/story/p_ironhold/advance"), map[string]any{}, 409, "E_STORY_NOT_CURRENT")

	advance("p_star", "")
	advance("p_star", "") // replay: no-op
	c.fail("POST", c.path("/story/p_star/advance"), map[string]any{"option": "x"}, 409, "E_STORY_CHOSEN")

	cur, _, _ = journal()
	require.Equal(t, "p_ironhold", cur["code"])
	require.Equal(t, false, cur["ready"])
	c.fail("POST", c.path("/story/p_ironhold/advance"), map[string]any{}, 409, "E_STORY_QUESTS")
	require.True(t, offered("story_p_bread"))
	c.ok("POST", c.path("/quests/story_p_bread/accept"), map[string]any{})
	finishQuest("story_p_bread")
	cur, _, _ = journal()
	require.Equal(t, true, cur["ready"])
	require.Equal(t, "turned_in", items(cur, "quests")[0]["state"])
	advance("p_ironhold", "")

	cur, _, _ = journal()
	require.Equal(t, "p_voice", cur["code"])
	require.Len(t, items(cur, "options"), 2)
	c.fail("POST", c.path("/story/p_voice/advance"), map[string]any{}, 400, "E_STORY_OPTION")
	c.fail("POST", c.path("/story/p_voice/advance"), map[string]any{"option": "bogus"}, 400, "E_STORY_OPTION")
	advance("p_voice", "tell")

	c.ok("POST", c.path("/quests/story_p_kettle/accept"), map[string]any{})
	finishQuest("story_p_kettle")
	advance("p_inn", "")
	cur, _, chapters = journal()
	require.Equal(t, "c1_wolves", cur["code"])
	require.Equal(t, "locked", cur["state"], "level 1 < 2")
	require.Equal(t, "done", chapters[0]["state"])
	require.False(t, offered("story_c1_wolves"), "a level-locked scene offers nothing")
	c.fail("POST", c.path("/story/c1_wolves/advance"), map[string]any{}, 403, "E_LEVEL_REQUIREMENT")

	// The faction choice: flags from earlier answers shape the scene, the
	// answer grants reputation.
	setLevel(15)
	pass(map[string]string{"c1_wolves": "", "c1_cage": "companion", "c1_crypt": "", "c1_caravan": "", "c1_ledger": ""})
	cur, state, _ := journal()
	require.Equal(t, "c1_doors", cur["code"])
	require.Equal(t, []any{"hedda"}, cur["extras"], "told Hedda in the prologue")
	flags := state["flags"].(map[string]any)
	require.Equal(t, "anlatti", flags["ses_sir"])
	require.Equal(t, "yoldas", flags["mirel"])
	require.Equal(t, float64(1), state["counters"].(map[string]any)["merhamet"])
	advance("c1_doors", "guard")
	var rep int
	require.NoError(t, db.QueryRow(ctx, `SELECT rep_xp FROM progression.character_reputation
		WHERE character_id = $1 AND faction_code = 'ironhold_guard'`, c.char).Scan(&rep))
	require.Equal(t, 500, rep)

	// A choice only the Covenant gets: for the Guard the scene just continues.
	setLevel(23)
	pass(map[string]string{"c1_peak": "", "c1_spire": "purify"})
	cur, _, _ = journal()
	require.Equal(t, "c1_descent", cur["code"])
	require.Equal(t, []any{"guard"}, cur["extras"])
	require.Empty(t, items(cur, "options"))
	advance("c1_descent", "")

	// An option that needs an optional quest.
	setLevel(63)
	pass(map[string]string{"c1_marsh": "", "c1_temple": "take", "c1_hedda": ""})
	cur, _, _ = journal()
	require.Equal(t, "c1_iron_king", cur["code"])
	require.Equal(t, "story_c1_anvil", items(cur, "options")[0]["requires_quest"])
	c.fail("POST", c.path("/story/c1_iron_king/advance"), map[string]any{"option": "save"}, 409, "E_STORY_QUESTS")
	_, err := db.Exec(ctx, `INSERT INTO progression.character_quest (id, character_id, quest_code, state, progress, accepted_at)
		VALUES ($1, $2, 'story_c1_mine', 'turned_in', '[0]', now())`, uuid.New(), c.char)
	require.NoError(t, err)
	c.fail("POST", c.path("/story/c1_iron_king/advance"), map[string]any{"option": "save"}, 409, "E_STORY_OPTION_LOCKED")
	advance("c1_iron_king", "shatter")
	_, state, _ = journal()
	require.Equal(t, "yikildi", state["flags"].(map[string]any)["hedda_klan"])

	// Text: complete in en + tr, other languages fall back to en.
	text := c.ok("GET", "/v1/story/text?lang=tr", nil)["text"].(map[string]any)
	require.Equal(t, "Düşen Yıldız", text["p_star"])
	text = c.ok("GET", "/v1/story/text?lang=ja", nil)["text"].(map[string]any)
	require.Equal(t, "The Falling Star", text["p_star"])

	other := newHTTPPlayer(t, "mage")
	other.fail("GET", c.path("/story"), nil, 404, "")
	other.fail("POST", c.path("/story/c1_keep/advance"), map[string]any{}, 404, "")
}
