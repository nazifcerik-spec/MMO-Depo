package quest

import (
	"math"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/otter/server/internal/crafting"
	"github.com/otter/server/internal/leveling"
)

func TestRewardForScalesWithLevels(t *testing.T) {
	daily := &Template{Category: CategoryDaily, RewardXP: 100, RewardGold: 300,
		Objectives: []Objective{{Event: "gather_mining", Count: 40}}}
	at := func(level int16) Reward { return RewardFor(daily, level, map[string]int16{"mining": 30}) }
	r := at(100)
	require.EqualValues(t, math.Round(float64(leveling.Current().XPToNext(100))*0.10), r.XP, "10% of the level")
	require.EqualValues(t, math.Round(10*math.Pow(100, 1.5)*5), r.GoldCopper, "five kills' gold")
	require.Equal(t, []ProfessionReward{{Kind: "mining",
		XP: int64(math.Round(float64(crafting.ProfessionXPToNext(30)) * 0.2))}}, r.Professions)
	require.Greater(t, at(500).XP, r.XP, "grows with the level")
	require.Greater(t, at(500).GoldCopper, r.GoldCopper)
	require.EqualValues(t, math.Round(float64(crafting.ProfessionXPToNext(1))*0.2),
		RewardFor(daily, 5, nil).Professions[0].XP, "an untrained profession counts as level 1")

	side := &Template{Category: CategorySide}
	require.EqualValues(t, math.Round(float64(leveling.Current().XPToNext(10))*0.25), RewardFor(side, 10, nil).XP)

	main := &Template{Category: CategoryMain, RewardXP: 777, RewardGold: 55}
	require.Equal(t, Reward{XP: 777, GoldCopper: 55}, RewardFor(main, 400, nil), "main quests keep their rewards")
}
