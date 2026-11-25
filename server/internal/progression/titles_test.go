package progression

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"

	"github.com/otter/server/internal/leveling"
)

// Every title comes with its requirement and the character's progress
// (DD-372).
func TestTitleGoals(t *testing.T) {
	raw, err := os.ReadFile("../../../content/progression.yaml")
	require.NoError(t, err)
	cfg := &leveling.Config{}
	require.NoError(t, yaml.Unmarshal(raw, cfg))
	require.NoError(t, cfg.Prepare())

	slayer, battle := "the_slayer", "lord_of_battle"
	got := titleGoals([]AchievementView{
		{Code: "first_blood", Event: "mob_killed", Target: 1},
		{Code: "slayer_1000", Event: "mob_killed", Target: 1000, Progress: 12, TitleCode: &slayer, RewardGold: 25000, Points: 25},
		{Code: "mastery_warlord", Event: "class_mastery:warlord", Target: 1, TitleCode: &battle, Points: 100},
	}, cfg, 750, "warrior", "warlord")

	require.Len(t, got, 1+10+40)
	require.Equal(t, TitleGoal{Code: "the_slayer", Source: "achievement", Event: "mob_killed", Target: 1000, Progress: 12,
		Achievement: "slayer_1000", RewardGold: 25000, Points: 25}, got[0])
	require.Equal(t, TitleGoal{Code: "veteran", Source: "level", Event: "level", Target: 100, Progress: 100}, got[1])
	require.Equal(t, TitleGoal{Code: "eternal", Source: "level", Event: "level", Target: 1000, Progress: 750}, got[10])
	for i, g := range got[11:] {
		require.Equal(t, i < 4, g.OwnClass, "the character's own class comes first")
		require.Equal(t, int64(1000), g.Target)
	}
	require.Contains(t, got, TitleGoal{Code: "lord_of_battle", Source: "class_mastery", Event: "class_mastery:warlord",
		Target: 1000, Progress: 750, Achievement: "mastery_warlord", Points: 100, OwnClass: true})
}
