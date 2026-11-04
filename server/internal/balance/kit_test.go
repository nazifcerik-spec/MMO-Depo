package balance

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/otter/server/internal/classskill"
	"github.com/otter/server/internal/leveling"
)

// DD-352 on the shipped content: along every path a character holds one active
// per job, always has something to cast, and the ultimate is never replaced.
func TestShippedKitsHoldOneOfEachJob(t *testing.T) {
	w := load(t)
	cfg := leveling.Current()
	for _, k := range cfg.Classes {
		if n, _ := cfg.Class(k.Code); n.Tier() != 2 {
			continue
		}
		path := cfg.Path(k.Code)
		jobs, ultimate := map[string]string{}, map[string]bool{}
		for _, s := range w.skills.ForClass(path[0]) {
			jobs[s.Code] = s.Job
			ultimate[s.Code] = s.Config.Category == "ultimate"
		}
		for _, level := range []int{1, 25, 75, 299, 300, 500, 699, 700, 850, 950, 1000} {
			at := level
			p := path
			switch {
			case level < 300:
				p = path[:1]
			case level < 700:
				p = path[:2]
			}
			seen, hasUltimate := map[string]bool{}, false
			kit := w.skills.Actives(p, at, "shield")
			require.NotEmpty(t, kit, "%v at %d", p, at)
			for _, s := range kit {
				require.False(t, seen[jobs[s.Code]], "%v at %d: two skills do %s", p, at, jobs[s.Code])
				seen[jobs[s.Code]] = true
				hasUltimate = hasUltimate || ultimate[s.Code]
			}
			require.Equal(t, at >= 950, hasUltimate, "%v at %d", p, at)
		}
	}
}

// DD-362: every base class gains a skill at levels 100, 150, 200 and 250. They
// take nothing from the base kit, and a gap skill never outclasses a later skill
// of the same job, or the promotion would lose what it grants.
func TestGapSkillsLeadToThePromotion(t *testing.T) {
	w := load(t)
	cfg := leveling.Current()
	for _, k := range cfg.Classes {
		if k.Parent != "" {
			continue
		}
		var gap, base, later []classskill.Skill
		for _, s := range w.skills.ForClass(k.Code) {
			switch {
			case s.Spec != k.Code:
				later = append(later, s)
			case s.UnlockLevel > 75:
				gap = append(gap, s)
			default:
				base = append(base, s)
			}
		}
		levels := map[int]int{}
		for _, s := range gap {
			levels[int(s.UnlockLevel)]++
		}
		require.Equal(t, map[int]int{100: 1, 150: 1, 200: 1, 250: 1}, levels, k.Code)

		kit := map[string]bool{}
		for _, a := range w.skills.Actives(cfg.Path(k.Code), 299) {
			kit[a.Code] = true
		}
		for _, s := range append(append([]classskill.Skill{}, base...), gap...) {
			if s.Kind == "active" {
				require.True(t, kit[s.Code], "%s dropped from the kit before the promotion", s.Code)
			}
		}
		for _, g := range gap {
			if g.Kind != "active" {
				continue
			}
			for _, l := range later {
				if l.Kind == "active" && l.Job == g.Job {
					require.Less(t, g.Strength, l.Strength, "%s would outclass %s", g.Code, l.Code)
				}
			}
		}
	}
}

// DD-353: every active earns its place: on some path through its class node it
// is in the kit when it unlocks, so no branch's skill is dead on arrival.
func TestEveryActiveIsUsedSomewhere(t *testing.T) {
	w := load(t)
	cfg := leveling.Current()
	used := map[string]bool{}
	all := map[string]bool{}
	for _, k := range cfg.Classes {
		n, _ := cfg.Class(k.Code)
		if n.Tier() != 2 {
			continue
		}
		path := cfg.Path(k.Code)
		for _, s := range w.skills.ForClass(path[0]) {
			if s.Kind != "active" {
				continue
			}
			all[s.Code] = true
			for _, a := range w.skills.Actives(path, int(s.UnlockLevel), "shield") {
				used[a.Code] = true
			}
		}
	}
	for code := range all {
		require.True(t, used[code], "%s is outclassed on every path", code)
	}
}
