package world

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func ptr(s string) *string { return &s }

func fakeCatalog(t *testing.T) *Catalog {
	t.Helper()
	regions := []Region{
		{Code: "cont", Kind: KindContinent, TravelSecondsFromParent: 0},
		{Code: "city", Kind: KindCity, ParentCode: ptr("cont"), TravelSecondsFromParent: 60, IsStarting: true},
		{Code: "wild_a", Kind: KindWilderness, ParentCode: ptr("cont"), TravelSecondsFromParent: 120, MinLevel: 1},
		{Code: "wild_b", Kind: KindWilderness, ParentCode: ptr("cont"), TravelSecondsFromParent: 180, MinLevel: 8},
		{Code: "dung", Kind: KindDungeon, ParentCode: ptr("wild_a"), TravelSecondsFromParent: 200, MinLevel: 5},
	}
	links := []RegionLink{
		{RegionA: "wild_a", RegionB: "wild_b", TravelSeconds: 240},
	}
	return NewCatalog(regions, links)
}

func TestCatalog_StartingCode(t *testing.T) {
	c := fakeCatalog(t)
	require.Equal(t, "city", c.StartingCode())
}

// Any region a road leads to is one journey away, at the shortest route's
// time (DD-280): cont is the hub between city and the wilds.
func TestCatalog_Routes(t *testing.T) {
	c := fakeCatalog(t)
	r := c.Routes("city")
	require.Len(t, r, 4, "every other region")
	require.Equal(t, Route{Seconds: 60, Path: []string{"city", "cont"}}, r["cont"], "child → parent")
	require.Equal(t, Route{Seconds: 180, Path: []string{"city", "cont", "wild_a"}}, r["wild_a"])
	require.Equal(t, Route{Seconds: 240, Path: []string{"city", "cont", "wild_b"}}, r["wild_b"], "via the hub beats the link")
	require.Equal(t, Route{Seconds: 380, Path: []string{"city", "cont", "wild_a", "dung"}}, r["dung"])
	_, self := r["city"]
	require.False(t, self, "the origin is no destination")
	require.Equal(t, int32(240), c.Routes("wild_a")["wild_b"].Seconds, "the link: 240 < 120+180")
	require.Nil(t, c.Routes("ghost"))

	island := NewCatalog([]Region{{Code: "a"}, {Code: "b"}}, nil)
	require.Empty(t, island.Routes("a"), "no road, no route")
}

func TestCatalog_Neighbors_ContainsAllDirections(t *testing.T) {
	c := fakeCatalog(t)
	// wild_a has: parent cont, child dung, link wild_b → 3 neighbors.
	ns := c.Neighbors("wild_a")
	require.Len(t, ns, 3)
	codes := map[string]string{}
	for _, n := range ns {
		codes[n.Region.Code] = n.Via
	}
	require.Equal(t, "parent", codes["cont"])
	require.Equal(t, "child", codes["dung"])
	require.Equal(t, "link", codes["wild_b"])
}

func TestCatalog_Neighbors_ReturnsNilForUnknown(t *testing.T) {
	c := fakeCatalog(t)
	require.Nil(t, c.Neighbors("ghost"))
}

// Hunting reaches the region itself and its road links, never a parent or
// child (DD-273): a continent must not reach every zone on it.
func TestCatalog_HuntRegions(t *testing.T) {
	c := fakeCatalog(t)
	require.ElementsMatch(t, []string{"wild_a", "wild_b"}, c.HuntRegions("wild_a"))
	require.Equal(t, []string{"cont"}, c.HuntRegions("cont"))
	require.Equal(t, []string{"city"}, c.HuntRegions("city"))
	require.Nil(t, c.HuntRegions("nowhere"))
}

func TestCatalog_Links(t *testing.T) {
	c := fakeCatalog(t)
	require.Equal(t, []RegionLink{{RegionA: "wild_a", RegionB: "wild_b", TravelSeconds: 240}}, c.Links())
}

// Dungeons are in reach from their region, a linked one, or a child region
// such as the entrance node (DD-302).
func TestCatalog_DungeonRegions(t *testing.T) {
	c := fakeCatalog(t)
	require.ElementsMatch(t, []string{"dung", "wild_a"}, c.DungeonRegions("dung"), "the entrance sees its zone")
	require.Contains(t, c.DungeonRegions("wild_a"), "wild_b", "linked regions as in hunting")
	require.Nil(t, c.DungeonRegions("nope"))
}
