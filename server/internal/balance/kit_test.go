package balance

import (
	"testing"

	"github.com/stretchr/testify/require"

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
			kit := w.skills.Actives(p, at)
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
