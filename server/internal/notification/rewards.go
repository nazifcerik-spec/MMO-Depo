package notification

import (
	"sort"

	"github.com/google/uuid"
)

// KindReward is the history entry for what a gathering, hunt or dungeon
// collection paid out (DD-242). Inbox only and written already read: the
// player has just seen the loot sheet, so it is a record, not an alert.
const KindReward = "reward_collected"

// RewardItem is one stack of identical drops ("Iron Ore × 5").
type RewardItem struct {
	Code   string  `json:"code"`
	Rarity string  `json:"rarity"`
	Prefix *string `json:"prefix,omitempty"`
	Suffix *string `json:"suffix,omitempty"`
	Count  int     `json:"count"`
}

// Reward describes one collected contract.
type Reward struct {
	AccountID   uuid.UUID
	CharacterID uuid.UUID
	Ref         uuid.UUID // contract id: one entry per contract, however often it is revealed
	Source      string    // activity | hunt | dungeon
	Target      string    // activity kind, mob code or dungeon code
	GoldCopper  int64
	XP          int64
	Items       []RewardItem
	NewLevel    int            // 0 = no level-up
	Extra       map[string]any // source-specific flags (e.g. dungeon cleared)
}

// Notice turns the reward into its inbox entry.
func (r Reward) Notice() Notice {
	char := r.CharacterID
	items := r.Items
	if items == nil {
		items = []RewardItem{}
	}
	vars := map[string]any{
		"source": r.Source, "target": r.Target,
		"gold_copper": r.GoldCopper, "xp": r.XP, "items": items,
	}
	if r.NewLevel > 0 {
		vars["new_level"] = r.NewLevel
	}
	for k, v := range r.Extra {
		vars[k] = v
	}
	return Notice{
		AccountID: r.AccountID, CharacterID: &char, Kind: KindReward,
		DedupKey: "reward:" + r.Ref.String(),
		TitleKey: "notif.reward.title", BodyKey: "notif.reward.body",
		Vars: vars, Category: "activity", Read: true,
	}
}

var rarityRank = map[string]int{"common": 0, "uncommon": 1, "rare": 2, "epic": 3, "legendary": 4, "mythic": 5}

// StackItems groups identical drops (template, rarity and affixes) into
// counted stacks, rarest first, then in order of first appearance.
func StackItems[T any](items []T, key func(T) RewardItem) []RewardItem {
	out := []RewardItem{}
	index := map[string]int{}
	for _, it := range items {
		ri := key(it)
		k := ri.Code + "|" + ri.Rarity + "|" + deref(ri.Prefix) + "|" + deref(ri.Suffix)
		if i, ok := index[k]; ok {
			out[i].Count++
			continue
		}
		ri.Count = 1
		index[k] = len(out)
		out = append(out, ri)
	}
	sort.SliceStable(out, func(i, j int) bool { return rarityRank[out[i].Rarity] > rarityRank[out[j].Rarity] })
	return out
}
