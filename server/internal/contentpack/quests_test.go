package contentpack

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestShippedQuestsValidate(t *testing.T) {
	dir := filepath.Join("..", "..", "..", "content")
	items, err := loadItemsFile(filepath.Join(dir, "items.yaml"))
	require.NoError(t, err)
	mobs, err := loadMobsFile(filepath.Join(dir, "mobs.yaml"))
	require.NoError(t, err)
	dungeons, err := loadDungeonsFile(filepath.Join(dir, "dungeons.yaml"))
	require.NoError(t, err)
	quests, err := loadQuestsFile(filepath.Join(dir, "quests.yaml"))
	require.NoError(t, err)
	require.NotEmpty(t, quests)
	require.NoError(t, validateQuests(quests, refs{
		mobs:     codeSet(mobs, func(m Mob) string { return m.Code }),
		dungeons: codeSet(dungeons, func(d Dungeon) string { return d.Code }),
		items:    codeSet(items, func(i Item) string { return i.Code }),
	}))
	daily := 0
	for _, q := range quests {
		if q.Category == "daily" && q.MinLevel == 1 {
			daily++
		}
	}
	require.GreaterOrEqual(t, daily, 3, "a level-1 character needs a full daily offer")
}

func TestValidateQuests_Rejects(t *testing.T) {
	r := refs{mobs: map[string]bool{"grey_wolf": true}, dungeons: map[string]bool{"crypt": true}, items: map[string]bool{}}
	ok := func() Quest {
		return Quest{Code: "q", DisplayNameKey: "k", Category: "main", MinLevel: 1,
			Objectives: []QuestObjective{{Event: "mob_killed:grey_wolf", Count: 3}}}
	}
	require.NoError(t, validateQuests([]Quest{ok()}, r))

	cases := map[string]func(q *Quest){
		"category":      func(q *Quest) { q.Category = "weekly" },
		"level":         func(q *Quest) { q.MinLevel = 0 },
		"no objectives": func(q *Quest) { q.Objectives = nil },
		"count":         func(q *Quest) { q.Objectives[0].Count = 0 },
		"unknown mob":   func(q *Quest) { q.Objectives[0].Event = "mob_killed:dragon" },
		"unknown kind":  func(q *Quest) { q.Objectives[0].Event = "player_killed" },
		"bad gather":    func(q *Quest) { q.Objectives[0].Event = "gather_farming" },
		"scoped gather": func(q *Quest) { q.Objectives[0].Event = "gather_mining:iron" },
		"scoped boss":   func(q *Quest) { q.Objectives[0].Event = "boss_killed:lich" },
		"unknown item":  func(q *Quest) { q.Objectives[0].Event = "item_crafted:sword" },
		"duplicate event": func(q *Quest) {
			q.Objectives = append(q.Objectives, q.Objectives[0])
		},
		"negative reward": func(q *Quest) { q.RewardGoldCopper = -1 },
		"unknown prereq":  func(q *Quest) { q.RequiresQuest = "nope" },
		"daily prereq": func(q *Quest) {
			q.Category = "daily"
			q.RequiresQuest = "q"
		},
	}
	for name, mut := range cases {
		q := ok()
		mut(&q)
		require.Error(t, validateQuests([]Quest{q}, r), name)
	}

	a, b := ok(), ok()
	b.Code = "b"
	a.RequiresQuest, b.RequiresQuest = "b", "q"
	require.ErrorContains(t, validateQuests([]Quest{a, b}, r), "cycle")

	a, b = ok(), ok()
	a.Code, b.Code = "a", "a"
	require.ErrorContains(t, validateQuests([]Quest{a, b}, r), "duplicate")

	d := ok()
	d.Code, d.Category = "d", "daily"
	m := ok()
	m.Code, m.RequiresQuest = "m", "d"
	require.ErrorContains(t, validateQuests([]Quest{d, m}, r), "daily")

	c := ok()
	c.Objectives = []QuestObjective{{Event: "dungeon_cleared:crypt", Count: 1}, {Event: "gather_fishing", Count: 5},
		{Event: "boss_killed", Count: 1}, {Event: "item_crafted", Count: 1}, {Event: "mob_killed", Count: 1}}
	require.NoError(t, validateQuests([]Quest{c}, r))
}
