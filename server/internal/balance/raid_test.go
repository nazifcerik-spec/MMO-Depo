package balance

import (
	"fmt"
	"strings"
	"testing"

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

// TestRaidMatrix prints clear rates per group (run with -v): five at the
// raid's level and five levels up, three at its level.
func TestRaidMatrix(t *testing.T) {
	if !testing.Verbose() {
		t.Skip("prints a matrix: run with -v")
	}
	w := load(t)
	var b strings.Builder
	for _, r := range w.raids {
		lv := int(r.MinLevel)
		fmt.Fprintf(&b, "%-18s %4d |", r.Code, lv)
		for _, g := range raidGroups(lv, 5) {
			fmt.Fprintf(&b, " %3.0f", 100*w.raidClearRate(g, lv, r))
		}
		b.WriteString(" | +5:")
		for _, g := range raidGroups(min(lv+5, 1000), 5) {
			fmt.Fprintf(&b, " %3.0f", 100*w.raidClearRate(g, min(lv+5, 1000), r))
		}
		b.WriteString(" | 3p:")
		for _, g := range raidGroups(lv, 3) {
			fmt.Fprintf(&b, " %3.0f", 100*w.raidClearRate(g, lv, r))
		}
		b.WriteString("\n")
	}
	t.Log("\n" + b.String())
}

// TestRaidTargets pins the DD-261 targets: a full mixed party at the raid's
// level usually clears it, a minimal party of three rarely does.
func TestRaidTargets(t *testing.T) {
	if raceDetector {
		t.Skip("pure simulation: runs in the coverage job (no -race)")
	}
	w := load(t)
	require.Len(t, w.raids, 6)
	for _, r := range w.raids {
		lv := int(r.MinLevel)
		var full, three float64
		for _, g := range raidGroups(lv, 5) {
			at := w.raidClearRate(g, lv, r)
			require.GreaterOrEqual(t, at, 0.40, "%s %v at min level: clearable", r.Code, g)
			full += at
		}
		for _, g := range raidGroups(lv, 3) {
			three += w.raidClearRate(g, lv, r)
		}
		full, three = full/10, three/10
		require.GreaterOrEqual(t, full, 0.70, "%s: five at min level clear most runs", r.Code)
		require.LessOrEqual(t, three, 0.30, "%s: three at min level rarely clear", r.Code)
	}
}
