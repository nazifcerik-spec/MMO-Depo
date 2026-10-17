package leveling_test

import (
	"math"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/otter/server/internal/contentpack"
	"github.com/otter/server/internal/leveling"
)

// shipped is the real content configuration (content/progression.yaml).
func shipped(t *testing.T) *leveling.Config {
	t.Helper()
	p, err := contentpack.Read("../../../content")
	require.NoError(t, err)
	require.NotNil(t, p.Progression)
	return p.Progression
}

func TestXPFormula(t *testing.T) {
	c := shipped(t)
	for level, want := range map[int]int64{1: 3752, 50: 4669, 100: 6719, 200: 13352, 300: 22828, 500: 49061,
		800: 104177, 900: 126344, 999: 150041} {
		require.Equal(t, want, c.XPToNext(level), "level %d", level)
	}
	require.Zero(t, c.XPToNext(1000), "no next level at the cap")
	require.Equal(t, c.XPToNext(1), c.XPToNext(0), "below 1 reads as level 1")
	require.Len(t, c.Table(), 999)
}

// TestProgressionSimulator walks a Warrior from 1 to 1000 at the baseline
// 1000 XP/min (development check, not gameplay code) and pins the DD-250
// cumulative timing targets.
func TestProgressionSimulator(t *testing.T) {
	c := shipped(t)
	targets := map[int]float64{50: 3, 100: 8, 200: 24.5, 300: 54, 400: 102, 500: 171, 600: 266, 700: 390,
		800: 546, 900: 738, 1000: 968}
	rate := c.Curve.BaselineXPPerMinute
	st, minutes := leveling.State{Spec: "warrior", Level: 1}, 0.0
	for st.Level < c.LevelCap {
		need := c.XPToNext(st.Level) - st.XP
		minutes += float64(need) / rate
		var gained []int
		st, gained = c.Advance(st, need)
		require.Equal(t, []int{st.Level}, gained)
		if c.PromotionPending(st.Spec, st.Level) {
			choice := c.Choices(st.Spec, st.Level)[0].Code
			st.Spec, st.XP, st.Overflow = choice, st.Overflow, 0
		}
		if want, ok := targets[st.Level]; ok {
			hours := minutes / 60
			require.InDelta(t, want, hours, math.Max(0.5, want*0.02), "level %d after %.1f h", st.Level, hours)
		}
	}
	final, ok := c.Class(st.Spec)
	require.True(t, ok)
	require.Equal(t, 2, final.Tier(), "the walk took both promotions")
	days := minutes / 60 / 9
	require.GreaterOrEqual(t, days, 107.0)
	require.LessOrEqual(t, days, 110.0)
}

func TestAdvance_SingleMultipleOverflowCap(t *testing.T) {
	c := shipped(t)
	// Single level, remainder kept.
	st, gained := c.Advance(leveling.State{Spec: "warrior", Level: 1}, c.XPToNext(1)+10)
	require.Equal(t, []int{2}, gained)
	require.Equal(t, int64(10), st.XP)

	// One big grant crosses several levels; nothing is lost.
	start := leveling.State{Spec: "warrior", Level: 24, XP: c.XPToNext(24) * 9 / 10}
	amount := int64(40_000)
	st, gained = c.Advance(start, amount)
	require.Equal(t, []int{25, 26, 27, 28, 29, 30, 31, 32, 33, 34}, gained[:10])
	total := start.XP + amount
	for l := 24; l < st.Level; l++ {
		total -= c.XPToNext(l)
	}
	require.Equal(t, total, st.XP, "overflow XP carried into the new level")

	// Promotion gate: stops at 100, the rest is held — even across calls.
	st, gained = c.Advance(leveling.State{Spec: "warrior", Level: 99}, 50_000)
	require.Equal(t, []int{100}, gained)
	require.Equal(t, 100, st.Level)
	require.Zero(t, st.XP)
	require.Equal(t, int64(50_000-c.XPToNext(99)), st.Overflow)
	st2, gained := c.Advance(st, 1234)
	require.Empty(t, gained)
	require.Equal(t, st.Overflow+1234, st2.Overflow)
	require.True(t, c.PromotionPending("warrior", 100))
	require.False(t, c.PromotionPending("guardian", 100))

	// Level 300 is the final specialization gate.
	st, _ = c.Advance(leveling.State{Spec: "guardian", Level: 299}, c.XPToNext(299)+5)
	require.Equal(t, 300, st.Level)
	require.Equal(t, int64(5), st.Overflow)
	require.ElementsMatch(t, []string{"iron_fortress", "warlord"}, codes(c.Choices("guardian", 300)))
	require.ElementsMatch(t, []string{"berserker", "weapon_master"}, codes(c.Choices("executioner", 300)))
	require.Empty(t, c.Choices("iron_fortress", 1000), "no class branch after 300")

	// Cap: 999 → 1000, never 1001; extra XP is kept as post-cap.
	st, gained = c.Advance(leveling.State{Spec: "iron_fortress", Level: 999}, c.XPToNext(999)*3)
	require.Equal(t, []int{1000}, gained)
	require.Equal(t, 1000, st.Level)
	require.Zero(t, st.XP)
	require.Equal(t, c.XPToNext(999)*2, st.PostCap)
	st, gained = c.Advance(st, 10)
	require.Empty(t, gained)
	require.Equal(t, 1000, st.Level)
}

func codes(list []*leveling.Class) []string {
	out := make([]string, len(list))
	for i, k := range list {
		out[i] = k.Code
	}
	return out
}

func TestPointSchedules(t *testing.T) {
	c := shipped(t)
	require.Equal(t, 525, c.AttributePointsEarned(1000))
	require.Equal(t, 0, c.AttributePointsEarned(1))
	require.Equal(t, 1, c.AttributePointsEarned(2))
	require.Equal(t, 49, c.AttributePointsEarned(99))
	require.Equal(t, 55, c.AttributePointsEarned(100), "50 + the level-100 milestone")
	require.Equal(t, 0, c.TalentPointsEarned(99))
	require.Equal(t, 1, c.TalentPointsEarned(100))
	require.Equal(t, 1, c.TalentPointsEarned(119))
	require.Equal(t, 2, c.TalentPointsEarned(120))
	require.Equal(t, 45, c.TalentPointsEarned(980))
	require.Equal(t, 45, c.TalentPointsEarned(1000), "no talent points after 980")

	// The registry hands out exactly those totals, level by level.
	attr, talent := 0, 0
	for l := 1; l <= c.LevelCap; l++ {
		for _, r := range c.RewardsAt(l) {
			switch r.Type {
			case leveling.RewardAttributePoints:
				attr += r.Amount
			case leveling.RewardTalentPoint:
				talent += r.Amount
			}
		}
	}
	require.Equal(t, 525, attr)
	require.Equal(t, 45, talent)

	types := func(level int) []string {
		var out []string
		for _, r := range c.RewardsAt(level) {
			out = append(out, r.Type+":"+r.Ref)
		}
		return out
	}
	require.ElementsMatch(t, []string{"ATTRIBUTE_POINTS:", "PROMOTION_UNLOCK:warrior", "TALENT_POINT:", "TITLE:veteran"}, types(100))
	require.Contains(t, types(600), "AWAKENING:")
	require.Contains(t, types(700), "ULTIMATE:")
	require.Contains(t, types(850), "CAPSTONE_ELIGIBLE:")
	require.Contains(t, types(1000), "CLASS_MASTERY:")
	require.Contains(t, types(1000), "TITLE:eternal")
}

func TestEffectiveAttribute(t *testing.T) {
	c := shipped(t)
	for raw, want := range map[int]float64{0: 0, 50: 50, 100: 100, 250: 220, 400: 302.5, 500: 337.5} {
		require.InDelta(t, want, c.EffectiveAttribute(raw), 1e-9, "raw %d", raw)
	}
	require.Equal(t, 60.0, c.Cap("crit_pct", 95), "percent stats are capped")
	require.Equal(t, 12.0, c.Cap("crit_pct", 12))
}

func TestNextMilestone(t *testing.T) {
	c := shipped(t)
	for _, tc := range []struct {
		spec  string
		level int
		want  int
		kind  string
	}{
		{"warrior", 73, 75, "skill"},
		{"warrior", 99, 100, "promotion"},
		{"warrior", 100, 100, "promotion"}, // pending choice is the objective
		{"guardian", 287, 300, "promotion"},
		{"iron_fortress", 584, 600, "awakening"},
		{"iron_fortress", 671, 700, "ultimate"},
		{"iron_fortress", 823, 850, "capstone"},
		{"iron_fortress", 963, 1000, "mastery"},
	} {
		m, ok := c.NextMilestone(tc.spec, tc.level)
		require.True(t, ok, "%s %d", tc.spec, tc.level)
		require.Equal(t, tc.want, m.Level, "%s %d", tc.spec, tc.level)
		require.Equal(t, tc.kind, m.Kind, "%s %d", tc.spec, tc.level)
	}
	m, _ := c.NextMilestone("warrior", 73)
	require.Equal(t, "warrior_execute", m.Ref)
	_, ok := c.NextMilestone("iron_fortress", 1000)
	require.False(t, ok, "nothing after mastery")
	require.Equal(t, []string{"warrior", "guardian", "iron_fortress"}, c.Path("iron_fortress"))
}

// A table curve replaces the formula without touching gameplay code.
func TestTableCurveAndValidation(t *testing.T) {
	cfg := &leveling.Config{LevelCap: 4, Curve: leveling.Curve{Kind: "table", Table: []int64{10, 20, 30}},
		AttributePoints: leveling.AttributeRule{Every: 2, Amount: 1}, TalentPoints: leveling.TalentRule{From: 2, Every: 1, Until: 3, Amount: 1}}
	require.NoError(t, cfg.Prepare())
	st, gained := cfg.Advance(leveling.State{Spec: "x", Level: 1}, 35)
	require.Equal(t, []int{2, 3}, gained)
	require.Equal(t, int64(5), st.XP)

	bad := *cfg
	bad.Curve.Table = []int64{10}
	require.Error(t, bad.Prepare(), "table must cover every level")
	bad = *cfg
	bad.Classes = []leveling.Class{{Code: "a"}, {Code: "b", Parent: "zz", RequiredLevel: 2}}
	require.Error(t, bad.Prepare(), "unknown parent")
	bad = *cfg
	bad.Classes = []leveling.Class{{Code: "a"}}
	bad.TalentTreePoints = 25
	bad.Trees = []leveling.Tree{{Code: "t", Class: "a", Talents: []leveling.Talent{{Code: "x", MaxRank: 5}}}}
	require.Error(t, bad.Prepare(), "a tree needs its spendable points and a capstone")
	require.NoError(t, leveling.Install(cfg))
	require.Equal(t, int64(10), leveling.Current().XPToNext(1))
	require.NoError(t, leveling.Install(shipped(t)))
}
