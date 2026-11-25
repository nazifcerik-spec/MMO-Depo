package contentpack

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// The codex (DD-372) over the shipped content: every tag opens a page that
// exists and has a name.
func TestCodexShipped(t *testing.T) {
	p := shippedPack(t)
	seeds := CodexSeeds{
		Achievements: []SeedAchievement{
			{Code: "slayer_1000", Trigger: "mob_killed", Title: "the_slayer", Target: 1000, Gold: 25000, Points: 25},
			{Code: "miner_500", Trigger: "gather_mining", Title: "prospector", Target: 500},
			{Code: "mastery_warlord", Trigger: "class_mastery:warlord", Title: "lord_of_battle", Target: 1},
		},
		Factions: []SeedSource{{Code: "ironhold_guard", Events: []string{"dungeon_cleared:crypt_of_ashes", "mob_killed"}}},
		Tracks:   []SeedTrack{{SeedSource: SeedSource{Code: "mining", Events: []string{"gather_mining"}}, MaxTier: 10}},
		Events:   []SeedEvent{{Code: "greenwood_incursion", Kind: "invasion", Region: "greenwood", MinLevel: 5}},
	}
	c := BuildCodex(p, p.Progression, seeds)
	en := p.Strings.ByLang["en"]
	for kind, pages := range c {
		for code, page := range pages {
			switch kind {
			case "list", "race": // named by the interface strings
			default:
				require.Contains(t, en, kind+"."+code)
			}
			for _, g := range page.Groups {
				require.NotEmpty(t, g.Refs, "%s/%s %s", kind, code, g.Key)
				seen := map[string]bool{}
				for _, r := range g.Refs {
					require.NotNil(t, c.Page(r.Kind, r.Code), "%s/%s %s → %s/%s", kind, code, g.Key, r.Kind, r.Code)
					require.False(t, seen[r.Kind+r.Code], "%s/%s lists %s twice", kind, code, r.Code)
					seen[r.Kind+r.Code] = true
				}
			}
			for _, g := range page.Goals {
				if g.Event != "level" {
					require.NotNil(t, c.eventPage(g.Event), "%s/%s goal %s opens a page", kind, code, g.Event)
				}
			}
		}
	}

	refs := func(kind, code, key string) []string {
		var out []string
		for _, g := range c.Page(kind, code).Groups {
			if g.Key == key {
				for _, r := range g.Refs {
					out = append(out, r.Code)
				}
			}
		}
		return out
	}
	// Both ends of a connection name each other.
	require.Contains(t, refs("item", "potion_hp_small", "sold_by"), "general_ironhold")
	require.Contains(t, refs("vendor", "general_ironhold", "sells"), "potion_hp_small")
	require.Contains(t, refs("item", "copper_ore", "gathered"), "mining")
	require.Contains(t, refs("profession", "mining", "nodes"), "copper_ore")
	for _, d := range p.Dungeons {
		require.Equal(t, []string{d.BossMobCode}, refs("dungeon", d.Code, "boss"))
		require.Contains(t, refs("mob", d.BossMobCode, "boss_of"), d.Code)
	}
	for _, q := range p.Quests {
		require.Len(t, c.Page("quest", q.Code).Goals, len(q.Objectives))
		for _, o := range q.Objectives {
			target := c.eventPage(o.Event)
			require.Contains(t, refs(target.Kind, target.Code, "quests"), q.Code)
		}
	}
	// Titles say what earns them.
	require.Equal(t, []CodexGoal{{Event: "mob_killed", Count: 1000}}, c.Page("title", "the_slayer").Goals)
	require.Equal(t, []string{"slayer_1000"}, refs("title", "the_slayer", "achievements"))
	require.Contains(t, refs("profession", "mining", "titles"), "prospector")
	require.Equal(t, []CodexGoal{{Event: "level", Count: 100}}, c.Page("title", "veteran").Goals)
	require.Equal(t, []CodexGoal{{Event: "class_mastery:warlord", Count: 1000}}, c.Page("title", "lord_of_battle").Goals)
	require.Contains(t, refs("class", "warlord", "titles"), "lord_of_battle")
	require.Len(t, refs("list", "mobs", "mobs"), len(huntedSet(p.Mobs)))
	require.Contains(t, refs("dungeon", "crypt_of_ashes", "reputation"), "ironhold_guard")
	// A raid does not spoil the pet index (DD-342); the egg names its raids.
	for _, s := range p.Pets {
		raids := refs("item", s.EggCode(), "dropped_in")
		require.NotEmpty(t, raids)
		require.NotContains(t, refs("dungeon", raids[0], "drops"), s.EggCode())
	}
}
