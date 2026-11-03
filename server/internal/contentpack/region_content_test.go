package contentpack

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// Every placed monster, dungeon and resource lands on its region's card, and
// region-bound quests follow the monster or dungeon they name (DD-338).
func TestRegionContentShipped(t *testing.T) {
	p := shippedPack(t)
	rc := BuildRegionContent(p)
	require.Len(t, rc, len(p.Regions), "every region gets a card, even an empty one")

	count := func(f func(*RegionContent) int) (n int) {
		for _, c := range rc {
			n += f(c)
		}
		return n
	}
	mobs, nodes := 0, 0
	for _, m := range p.Mobs {
		if m.RegionCode != "" && m.Hunted() {
			mobs++
		}
	}
	for _, it := range p.Items {
		if it.Gather != nil && it.Gather.Region != "" {
			nodes++
		}
	}
	require.Equal(t, mobs, count(func(c *RegionContent) int { return len(c.Mobs) }))
	require.Equal(t, nodes, count(func(c *RegionContent) int { return len(c.Gather) }))
	for _, reg := range p.Regions {
		if reg.Kind == "dungeon" {
			require.NotEmpty(t, rc[reg.Code].Dungeons, "entrance %s shows its dungeons", reg.Code)
		}
	}
	require.NotZero(t, count(func(c *RegionContent) int { return len(c.Quests) }))

	for code, c := range rc {
		for i := 1; i < len(c.Mobs); i++ {
			require.LessOrEqual(t, c.Mobs[i-1].Level, c.Mobs[i].Level, code)
		}
	}
}

func TestRegionContentQuestFollowsItsTarget(t *testing.T) {
	p := &Pack{
		Regions:  []Region{{Code: "a"}, {Code: "b"}, {Code: "b_gate", Kind: "dungeon", ParentCode: "b"}},
		Mobs:     []Mob{{Code: "wolf", Level: 3, RegionCode: "a"}, {Code: "titan", RegionCode: "b", RaidBoss: true}},
		Dungeons: []Dungeon{{Code: "den", MinLevel: 5, RegionCode: "b"}, {Code: "keep", RegionCode: "b", GroupMax: 5}},
		Items:    []Item{{Code: "ore", Gather: &Gather{Kind: "mining", Level: 1, Region: "a"}}},
		Quests: []Quest{
			{Code: "hunt", Category: "side", Objectives: []QuestObjective{{Event: "gather_mining"}, {Event: "mob_killed:wolf"}}},
			{Code: "clear", Category: "main", Objectives: []QuestObjective{{Event: "dungeon_cleared:den"}}},
			{Code: "craft", Category: "daily", Objectives: []QuestObjective{{Event: "item_crafted:x"}}},
		},
	}
	rc := BuildRegionContent(p)
	require.Equal(t, []RegionMob{{Code: "wolf", Level: 3}}, rc["a"].Mobs)
	require.Empty(t, rc["b"].Mobs, "a raid boss shows as its raid")
	require.Equal(t, []RegionDungeon{{Code: "keep", Raid: true}, {Code: "den", MinLevel: 5}}, rc["b"].Dungeons)
	require.Equal(t, rc["b"].Dungeons, rc["b_gate"].Dungeons, "an entrance shows its zone's dungeons")
	require.Equal(t, []RegionNode{{Code: "ore", Kind: "mining", Level: 1}}, rc["a"].Gather)
	require.Equal(t, []RegionQuest{{Code: "hunt", Category: "side"}}, rc["a"].Quests)
	require.Equal(t, []RegionQuest{{Code: "clear", Category: "main"}}, rc["b"].Quests)
}
