package item

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/otter/server/internal/activity"
)

func TestRollAffix_FollowsRarityRules(t *testing.T) {
	sword := Template{Code: "sword", ItemClass: "weapon", Affixable: true}
	inPool := func(pool []string, v *string) bool {
		for _, p := range pool {
			if v != nil && *v == p {
				return true
			}
		}
		return false
	}

	// Not affixable or common: never.
	for i := int64(0); i < 50; i++ {
		p, s := rollAffix(activity.NewRNG(i), Template{ItemClass: "weapon"}, RarityLegendary)
		require.Nil(t, p)
		require.Nil(t, s)
		p, s = rollAffix(activity.NewRNG(i), sword, RarityCommon)
		require.Nil(t, p)
		require.Nil(t, s)
	}

	// Legendary: always both, from the class pools; optional tiers sometimes.
	var optPrefix, optNone int
	for i := int64(0); i < 200; i++ {
		p, s := rollAffix(activity.NewRNG(i), sword, RarityLegendary)
		require.True(t, inPool(prefixPool["weapon"], p), "prefix %v", p)
		require.True(t, inPool(suffixPool["weapon"], s), "suffix %v", s)

		p, _ = rollAffix(activity.NewRNG(i), sword, RarityEpic)
		require.NotNil(t, p, "epic always has a prefix")

		p, s = rollAffix(activity.NewRNG(i), sword, RarityRare)
		if p != nil {
			optPrefix++
		}
		if p == nil && s == nil {
			optNone++
		}
	}
	require.Positive(t, optPrefix, "rare rolls a prefix sometimes")
	require.Positive(t, optNone, "rare rolls nothing sometimes")

	// Class without pools (e.g. affixable consumable) stays bare.
	p, s := rollAffix(activity.NewRNG(1), Template{ItemClass: "consumable", Affixable: true}, RarityLegendary)
	require.Nil(t, p)
	require.Nil(t, s)
}

func TestRollAffix_Deterministic(t *testing.T) {
	sword := Template{ItemClass: "armor", Affixable: true}
	p1, s1 := rollAffix(activity.NewRNG(42), sword, RarityMythic)
	p2, s2 := rollAffix(activity.NewRNG(42), sword, RarityMythic)
	require.Equal(t, *p1, *p2)
	require.Equal(t, *s1, *s2)
}

// Every affix that can roll adds a stat, growing with the item's level
// (DD-275).
func TestAffixBonuses(t *testing.T) {
	for _, pools := range []map[string][]string{prefixPool, suffixPool} {
		for class, codes := range pools {
			for _, code := range codes {
				c := code
				got := AffixBonuses(&c, nil, 500)
				require.Len(t, got, 1, "%s %s adds a stat", class, code)
				require.Greater(t, got[0].Amount, AffixBonuses(&c, nil, 1)[0].Amount, "%s grows with level", code)
			}
		}
	}
	sharp, bear, nope := "sharp", "of_the_bear", "nope"
	require.Equal(t, []StatBonus{{"attack", 12}, {"str", 9}}, AffixBonuses(&sharp, &bear, 100))
	require.Equal(t, []StatBonus{{"attack", 2}}, AffixBonuses(&sharp, &nope, 0), "unknown codes add nothing")
	require.Empty(t, AffixBonuses(nil, nil, 100))
	require.Equal(t, 1.0, StatMultiplier(RarityCommon))
	require.Greater(t, StatMultiplier(RarityMythic), StatMultiplier(RarityLegendary))
}
