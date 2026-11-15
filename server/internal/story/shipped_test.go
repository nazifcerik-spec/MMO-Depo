package story

import (
	"math/rand/v2"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

// Every path through the shipped story (DD-367) ends: random answers, seeded,
// thousands of times. The ending scene shows exactly one ending, the epilogue
// exactly one last line, the final choice always offers the three open
// endings, and every ending variant is reached by some path.
func TestShippedStoryEveryPathEnds(t *testing.T) {
	raw, err := os.ReadFile("../../../content/story.yaml")
	require.NoError(t, err)
	var b Book
	require.NoError(t, yaml.Unmarshal(raw, &b))
	require.NotEmpty(t, b.Scenes)

	isMain := func(k string) bool {
		return (strings.HasPrefix(k, "chain_") || strings.HasPrefix(k, "crown_") ||
			strings.HasPrefix(k, "breath_") || k == "lullaby") && !strings.HasSuffix(k, "_uriel")
	}
	seen := map[string]bool{}
	rng := rand.New(rand.NewPCG(367, 1))
	races := []string{"human", "dwarf", "undead", "drakari", "dark_elf"}
	for run := 0; run < 6000; run++ {
		done := map[string]string{}
		race := races[run%len(races)]
		for steps := 0; ; steps++ {
			require.Less(t, steps, len(b.Scenes)+1, "the walk must end")
			st, cur := b.walk(done, func(*Scene, State) {})
			if cur == nil {
				break
			}
			before, after := extras(cur, st, race)
			all := append(before, after...)
			switch cur.Code {
			case "e_ending":
				var main []string
				for _, k := range all {
					if isMain(k) {
						main = append(main, k)
						seen[k] = true
					}
				}
				require.Len(t, main, 1, "one ending, state %v", st)
			case "e_epilogue":
				var last []string
				for _, k := range all {
					if strings.HasPrefix(k, "last_") {
						last = append(last, k)
					}
				}
				require.Len(t, last, 1, "one last line, state %v", st)
			case "e_choice":
				require.GreaterOrEqual(t, len(cur.Visible(st)), 3)
			}
			opts := cur.Visible(st)
			pick := ""
			if len(opts) > 0 {
				pick = opts[rng.IntN(len(opts))].Code
			}
			done[cur.Code] = pick
		}
	}
	for _, k := range []string{"chain_strong", "chain_cracked", "crown_keeper", "crown_ash", "crown_saved",
		"crown_half", "breath_full", "breath_half", "lullaby"} {
		require.True(t, seen[k], "ending %s never reached", k)
	}
}
