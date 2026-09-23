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

func TestCatalog_TravelDuration_ParentChildBothWays(t *testing.T) {
	c := fakeCatalog(t)
	// parent → child
	d, ok := c.TravelDuration("cont", "city")
	require.True(t, ok)
	require.Equal(t, int32(60), d, "child's travel_seconds_from_parent used for parent→child")
	// child → parent (symmetric)
	d, ok = c.TravelDuration("city", "cont")
	require.True(t, ok)
	require.Equal(t, int32(60), d)
}

func TestCatalog_TravelDuration_ExplicitLink(t *testing.T) {
	c := fakeCatalog(t)
	d, ok := c.TravelDuration("wild_a", "wild_b")
	require.True(t, ok)
	require.Equal(t, int32(240), d)
	// reverse
	d, ok = c.TravelDuration("wild_b", "wild_a")
	require.True(t, ok)
	require.Equal(t, int32(240), d)
}

func TestCatalog_TravelDuration_NoRouteBetweenUnrelated(t *testing.T) {
	c := fakeCatalog(t)
	_, ok := c.TravelDuration("city", "wild_b")
	require.False(t, ok, "siblings without an explicit link cannot travel directly")
}

func TestCatalog_TravelDuration_SelfIsNotARoute(t *testing.T) {
	c := fakeCatalog(t)
	_, ok := c.TravelDuration("city", "city")
	require.False(t, ok)
}

func TestCatalog_TravelDuration_UnknownRegions(t *testing.T) {
	c := fakeCatalog(t)
	_, ok := c.TravelDuration("ghost", "city")
	require.False(t, ok)
	_, ok = c.TravelDuration("city", "ghost")
	require.False(t, ok)
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
