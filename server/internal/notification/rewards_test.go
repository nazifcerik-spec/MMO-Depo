package notification

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestStackItems(t *testing.T) {
	type drop struct{ code, rarity, prefix string }
	key := func(d drop) RewardItem {
		ri := RewardItem{Code: d.code, Rarity: d.rarity}
		if d.prefix != "" {
			ri.Prefix = &d.prefix
		}
		return ri
	}
	got := StackItems([]drop{
		{"iron_ore", "common", ""}, {"copper_ore", "uncommon", ""}, {"iron_ore", "common", ""},
		{"iron_ore", "rare", ""}, {"iron_ore", "common", ""}, {"sword_iron", "rare", "ancient"},
		{"copper_ore", "uncommon", ""},
	}, key)
	simple := make([]string, len(got))
	for i, g := range got {
		simple[i] = g.Code + "/" + g.Rarity
		if g.Prefix != nil {
			simple[i] += "/" + *g.Prefix
		}
		simple[i] += "x" + string(rune('0'+g.Count))
	}
	// Rarest first; same rarity keeps first-seen order; affixes split stacks.
	require.Equal(t, []string{"iron_ore/rarex1", "sword_iron/rare/ancientx1", "copper_ore/uncommonx2", "iron_ore/commonx3"}, simple)
	require.Empty(t, StackItems([]drop{}, key))
}

func TestRewardNotice(t *testing.T) {
	ref := uuid.New()
	n := Reward{
		AccountID: uuid.New(), CharacterID: uuid.New(), Ref: ref, Source: "hunt", Target: "goblin_raider",
		GoldCopper: 120, XP: 45, NewLevel: 6, Extra: map[string]any{"victories": 3},
	}.Notice()
	require.Equal(t, KindReward, n.Kind)
	require.Equal(t, "reward:"+ref.String(), n.DedupKey, "one entry per contract, however often it is revealed")
	require.True(t, n.Read, "a history entry must not raise the unread badge")
	require.Equal(t, ModeInboxOnly, DefaultModes[KindReward], "never pushed")
	require.Equal(t, []RewardItem{}, n.Vars["items"], "items is always a list")
	require.Equal(t, 6, n.Vars["new_level"])
	require.Equal(t, 3, n.Vars["victories"])
	require.NotContains(t, Reward{}.Notice().Vars, "new_level")
}
