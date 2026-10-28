package balance

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/otter/server/internal/activity"
	"github.com/otter/server/internal/combat"
	"github.com/otter/server/internal/dungeon"
	"github.com/otter/server/internal/raid"
)

const raidSeeds = 40

// raidGroups are the parties the raid targets test at a level: ten mixed
// groups of size n, each n consecutive base classes (wrapping) on their first
// open path, so every class is in n of them.
func raidGroups(level, n int) [][]string {
	first := map[string]string{}
	for _, b := range builds(level) {
		c, _, _ := strings.Cut(b, "/")
		if _, ok := first[c]; !ok {
			first[c] = b
		}
	}
	var out [][]string
	for i := range classes {
		g := make([]string, n)
		for j := range g {
			g[j] = first[classes[(i+j)%len(classes)]]
		}
		out = append(out, g)
	}
	return out
}

// raidClearRate is the share of runs in which the group kills the boss.
func (w *world) raidClearRate(group []string, level int, t dungeon.Template) float64 {
	players := make([]combat.Combatant, len(group))
	for i, b := range group {
		p := w.playerIn(b, level, w.kit(b, level), true)
		p.RegenPerEncounter = int32(float64(t.DurationSecs/int32(t.TotalEncounters())/6) * float64(p.RegenPerEncounter) / float64(encSeconds/2))
		players[i] = p
	}
	idx, n := w.mobIndex(), 0
	for s := int64(1); s <= raidSeeds; s++ {
		if raid.Roll(s*104729, players, t, idx, func(seed int64) combat.RNG { return activity.NewRNG(seed) }).Cleared {
			n++
		}
	}
	return float64(n) / raidSeeds
}

// TestRaidMatrix prints the clear rate per group size (run with -v): the
// average over ten mixed groups and the worst group, at the raid's level,
// then a full party five levels up.
func TestRaidMatrix(t *testing.T) {
	if !testing.Verbose() {
		t.Skip("prints a matrix: run with -v")
	}
	w := load(t)
	var b strings.Builder
	for _, r := range w.raids {
		lv := int(r.MinLevel)
		fmt.Fprintf(&b, "%-18s %4d |", r.Code, lv)
		for n := 2; n <= 5; n++ {
			avg, worst := w.raidSize(lv, n, r)
			fmt.Fprintf(&b, " %dp %3.0f/%3.0f", n, 100*avg, 100*worst)
		}
		up, _ := w.raidSize(min(lv+5, 1000), 5, r)
		fmt.Fprintf(&b, " | +5: %3.0f\n", 100*up)
	}
	t.Log("\n" + b.String())
}

// raidSize is the average clear rate of ten mixed groups of n and the worst.
func (w *world) raidSize(level, n int, t dungeon.Template) (avg, worst float64) {
	worst = 1
	for _, g := range raidGroups(level, n) {
		c := w.raidClearRate(g, level, t)
		avg += c
		worst = min(worst, c)
	}
	return avg / 10, worst
}

// TestRaidTargets pins the DD-298 targets, made harder by DD-307: raids
// take a party of two to five, sized to it, so every size at the raid's
// level clears it about as often (30..70% on average) and a full party five
// levels up clears at least 55% (level 1000 has no higher level).
func TestRaidTargets(t *testing.T) {
	if raceDetector {
		t.Skip("pure simulation: runs in the coverage job (no -race)")
	}
	w := load(t)
	require.Len(t, w.raids, 6)
	for _, r := range w.raids {
		lv := int(r.MinLevel)
		require.EqualValues(t, 2, r.GroupMin, r.Code)
		require.EqualValues(t, 5, r.GroupMax, r.Code)
		for n := 2; n <= 5; n++ {
			avg, worst := w.raidSize(lv, n, r)
			t.Logf("%s Sv %d, %d members: %.2f (worst %.2f)", r.Code, lv, n, avg, worst)
			assert.GreaterOrEqual(t, avg, 0.30, "%s with %d: clearable", r.Code, n)
			assert.LessOrEqual(t, avg, 0.70, "%s with %d: still a challenge", r.Code, n)
		}
		if lv < 1000 {
			up, _ := w.raidSize(lv+5, 5, r)
			assert.GreaterOrEqual(t, up, 0.55, "%s: five levels up a full party clears most runs", r.Code)
		}
	}
}
