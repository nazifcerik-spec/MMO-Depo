//go:build integration

package integration

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/otter/server/internal/domain/character"
	"github.com/otter/server/internal/leveling"
)

// DD-252: the ten base classes, twenty promotions and forty
// specializations on the shared engine, against the real database.

var allBaseClasses = []string{"warrior", "rogue", "ranger", "mage", "monk", "cleric", "paladin", "druid", "bard", "shaman"}

// classFixture creates a character at (spec, level) for any class node —
// the reusable fixture factory of DD-252 (the class is the node's root).
func (w *testEnv) classFixture(spec string, level int16) player {
	w.t.Helper()
	path := leveling.Current().Path(spec)
	p := w.player(path[0], level)
	if spec != path[0] {
		_, err := w.p.Exec(w.ctx, `UPDATE character.character SET spec_code = $2 WHERE id = $1`, p.ID, spec)
		require.NoError(w.t, err)
	}
	return p
}

func firstChild(spec string) string {
	return leveling.Current().Choices(spec, 1000)[0].Code
}

// Every base class at the DD-252 checkpoints: 1, 99, 100, 299, 300, 599,
// 600, 699, 700, 849, 850, 999, 1000.
func TestAllClassLevelFixtures(t *testing.T) {
	w := newWorld(t)
	cfg := leveling.Current()
	for _, class := range allBaseClasses {
		promo := firstChild(class)
		spec := firstChild(promo)
		for _, f := range []struct {
			node      string
			level     int16
			promotion bool
			next      int
			title     string
		}{
			{class, 1, false, 10, ""},
			{class, 99, false, 100, ""},
			{class, 100, true, 100, "veteran"},
			{promo, 299, false, 300, "elite"},
			{promo, 300, true, 300, "champion"},
			{spec, 599, false, 600, "legendary"},
			{spec, 600, false, 700, "mythic"},
			{spec, 699, false, 700, "mythic"},
			{spec, 700, false, 800, "ancient"},
			{spec, 849, false, 850, "ascendant"},
			{spec, 850, false, 900, "ascendant"},
			{spec, 999, false, 1000, "immortal"},
			{spec, 1000, false, 0, "eternal"},
		} {
			p := w.classFixture(f.node, f.level)
			sum, err := w.chars.ProgressionSummary(w.ctx, p.ID, p.Account)
			name := fmt.Sprintf("%s %d", f.node, f.level)
			require.NoError(t, err, name)
			require.Equal(t, f.promotion, sum.Promotion.Available, name)
			require.Equal(t, f.title, sum.Title, name)
			require.Equal(t, f.level == 1000, sum.Mastery.Unlocked, name)
			if f.next == 0 {
				require.Nil(t, sum.NextMilestone, name)
			} else {
				require.NotNil(t, sum.NextMilestone, name)
				require.Equal(t, f.next, sum.NextMilestone.Level, name)
			}
			require.Equal(t, cfg.ResourceOf(class).Kind, sum.Resource.Kind, name)
			sheet, err := w.chars.Sheet(w.ctx, p.ID, p.Account)
			require.NoError(t, err, name)
			require.Equal(t, cfg.ResourceOf(class).Kind, sheet.ResourceKind, name)
			require.Positive(t, sheet.Derived.Attack, name)
		}
	}
}

// Overflow XP at both gates for every class, then the choice pays it out.
// Cross-class choices are refused.
func TestPromotionOverflowEveryClass(t *testing.T) {
	w := newWorld(t)
	for _, class := range allBaseClasses {
		p := w.classFixture(class, 99)
		res := w.grant(p, toLevel(99, 104)+100, class+"-afk")
		require.Equal(t, int16(100), res.Character.Level, class)
		require.True(t, res.PromotionAvailable, class)
		require.Equal(t, toLevel(100, 104)+100, res.Character.OverflowXP, class)
		promo := firstChild(class)
		g, err := w.chars.SelectPromotion(w.ctx, p.ID, p.Account, promo)
		require.NoError(t, err, class)
		require.Equal(t, int16(104), g.Character.Level, class)
		require.Zero(t, g.Character.OverflowXP, class)

		q := w.classFixture(promo, 299)
		res = w.grant(q, toLevel(299, 300)+77, class+"-gate2")
		require.Equal(t, int16(300), res.Character.Level, class)
		require.Equal(t, int64(77), res.Character.OverflowXP, class)
		g, err = w.chars.SelectPromotion(w.ctx, q.ID, q.Account, firstChild(promo))
		require.NoError(t, err, class)
		require.Equal(t, int64(77), g.Character.XP, class)
	}
	for _, bad := range []struct{ node, choice string }{
		{"mage", "guardian"}, {"cleric", "berserker"}, {"guardian", "pyromancer"},
		{"assassin", "beastmaster"}, {"rogue", "nightblade"},
	} {
		lvl := int16(100)
		if leveling.Current().Choices(bad.node, 1000)[0].RequiredLevel == 300 {
			lvl = 300
		}
		p := w.classFixture(bad.node, lvl)
		_, err := w.chars.SelectPromotion(w.ctx, p.ID, p.Account, bad.choice)
		require.ErrorIs(t, err, character.ErrInvalidPromotion, "%s → %s", bad.node, bad.choice)
	}
}

// One fixture per final specialization: reaching 1000 grants the mastery
// title, achievement and event exactly once; awakening and ultimate are
// unlocked on the way.
func TestAllSpecializationFixtures(t *testing.T) {
	w := newWorld(t)
	cfg := leveling.Current()
	n := 0
	for _, k := range cfg.Classes {
		node, _ := cfg.Class(k.Code)
		if node.Tier() != 2 {
			continue
		}
		n++
		p := w.classFixture(k.Code, 599)
		w.grant(p, toLevel(599, 1000)+5, k.Code+"-to-cap")
		c := w.char(p)
		require.Equal(t, int16(1000), c.Level, k.Code)
		require.Equal(t, int64(5), c.PostCapXP, k.Code)
		require.Len(t, w.events(p, character.EventMastery), 1, k.Code)
		require.Len(t, w.events(p, character.EventAwakening), 1, "%s awakening", k.Code)
		require.Len(t, w.events(p, character.EventUltimate), 1, "%s ultimate", k.Code)
		require.True(t, w.titles(p)[node.Mastery.Title], k.Code)
		var unlocked bool
		require.NoError(t, w.p.QueryRow(w.ctx, `
			SELECT unlocked_at IS NOT NULL FROM progression.achievement_progress
			 WHERE character_id = $1 AND achievement_code = $2`, p.ID, node.Mastery.Achievement).Scan(&unlocked))
		require.True(t, unlocked, k.Code)
		w.grant(p, 1_000_000, k.Code+"-more")
		require.Len(t, w.events(p, character.EventMastery), 1, "%s: mastery unlocks once", k.Code)
		sum, err := w.chars.ProgressionSummary(w.ctx, p.ID, p.Account)
		require.NoError(t, err)
		require.True(t, sum.Mastery.PrestigeEligible, k.Code)
		require.Equal(t, node.Mastery.Achievement, sum.Mastery.Achievement)
	}
	require.Equal(t, 40, n)
}

// Solo Adaptation over the production router: a support character hits
// harder alone than in a party (DD-252).
func TestHTTPSoloAdaptation(t *testing.T) {
	cleric := newHTTPPlayer(t, "cleric")
	friend := newHTTPPlayer(t, "warrior")
	solo := cleric.ok("GET", cleric.path("/sheet"), nil)["derived"].(map[string]any)["attack"].(float64)
	cleric.ok("POST", cleric.path("/party"), map[string]any{})
	inv := cleric.ok("POST", cleric.path("/party/invite"), map[string]any{"target_name": friend.name})["id"].(string)
	friend.ok("POST", friend.path("/party/invites/"+inv+"/accept"), map[string]any{})
	grouped := cleric.ok("GET", cleric.path("/sheet"), nil)["derived"].(map[string]any)["attack"].(float64)
	require.Greater(t, solo, grouped, "support power converts to attack only alone")
	w := friend.ok("GET", friend.path("/sheet"), nil)["derived"].(map[string]any)["attack"].(float64)
	require.Positive(t, w)
}

// AFK hunts resolve for every class through the real router.
func TestHTTPAFKHuntEveryClass(t *testing.T) {
	w := newWorld(t)
	for _, class := range allBaseClasses {
		p := newHTTPPlayer(t, class)
		_, err := w.p.Exec(w.ctx, `UPDATE character.character SET level = 6 WHERE id = $1`, p.char)
		require.NoError(t, err)
		h := p.ok("POST", p.path("/hunt"), map[string]any{"mob_code": "young_boar", "duration_seconds": 3600})
		hid := h["contract"].(map[string]any)["id"].(string)
		_, err = w.p.Exec(w.ctx, `
			UPDATE hunt.contract SET started_at = now() - interval '61 minutes', expires_at = now() - interval '1 minute'
			 WHERE id = $1`, hid)
		require.NoError(t, err)
		res := p.ok("POST", "/v1/hunt/"+hid+"/reveal", map[string]any{})["hunt"].(map[string]any)
		require.Positive(t, res["victories"], class)
		require.Equal(t, false, res["stopped_on_death"], class)
	}
}
