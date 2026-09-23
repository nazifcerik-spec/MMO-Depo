package quest

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/otter/server/internal/progression"
)

func dailies(levels ...int16) *Catalog {
	var list []Template
	for i, lv := range levels {
		list = append(list, Template{Code: string(rune('a' + i)), Category: CategoryDaily, MinLevel: lv})
	}
	list = append(list, Template{Code: "main", Category: CategoryMain, MinLevel: 1})
	return NewCatalog(list)
}

func TestDailyOffer_StablePerDayAndLevelGated(t *testing.T) {
	cat := dailies(1, 1, 1, 1, 1, 1, 10)
	char := uuid.MustParse("6f1c1e8e-0000-4000-8000-000000000001")
	day := time.Date(2026, 9, 23, 13, 0, 0, 0, time.UTC)

	offer := cat.DailyOffer(char, 5, day)
	require.Len(t, offer, DailySlots)
	require.Equal(t, offer, cat.DailyOffer(char, 5, day.Add(10*time.Hour)), "same UTC day → same offer")
	require.NotContains(t, offer, "g", "level-10 daily is not offered at level 5")
	require.NotContains(t, offer, "main")

	// Across days and characters the offer varies (not a fixed prefix).
	seen := map[string]bool{}
	for d := 0; d < 30; d++ {
		for _, c := range cat.DailyOffer(uuid.New(), 5, day.AddDate(0, 0, d)) {
			seen[c] = true
		}
	}
	require.Len(t, seen, 6)

	require.Len(t, dailies(1, 1).DailyOffer(char, 1, day), 2, "small pool → fewer slots")
	require.Empty(t, dailies(20).DailyOffer(char, 1, day))
}

func TestDayKey(t *testing.T) {
	loc := time.FixedZone("UTC+3", 3*3600)
	require.Equal(t, time.Date(2026, 9, 22, 0, 0, 0, 0, time.UTC), DayKey(time.Date(2026, 9, 23, 1, 0, 0, 0, loc)))
}

func TestAdvance(t *testing.T) {
	tmpl := &Template{Objectives: []Objective{{Event: "mob_killed:grey_wolf", Count: 5}, {Event: "gather_mining", Count: 10}}}

	next, changed, done := Advance(tmpl, []int64{0, 0}, []progression.Event{
		{Kind: "mob_killed", Amount: 9}, {Kind: "mob_killed:grey_wolf", Amount: 3},
	})
	require.Equal(t, []int64{3, 0}, next, "generic mob_killed does not count toward a scoped objective")
	require.True(t, changed)
	require.False(t, done)

	next, changed, done = Advance(tmpl, next, []progression.Event{
		{Kind: "mob_killed:grey_wolf", Amount: 50}, {Kind: "gather_mining", Amount: 10},
	})
	require.Equal(t, []int64{5, 10}, next, "capped at count")
	require.True(t, changed)
	require.True(t, done)

	again, changed, done := Advance(tmpl, next, []progression.Event{{Kind: "gather_mining", Amount: 4}})
	require.Equal(t, next, again)
	require.False(t, changed, "finished objectives do not change")
	require.True(t, done)

	_, changed, _ = Advance(tmpl, []int64{0, 0}, []progression.Event{{Kind: "gather_fishing", Amount: 4}})
	require.False(t, changed)
}

func TestCatalogOrder(t *testing.T) {
	cat := NewCatalog([]Template{{Code: "b", SortOrder: 2}, {Code: "z", SortOrder: 1}, {Code: "a", SortOrder: 2}})
	require.Equal(t, []string{"z", "a", "b"}, []string{cat.list[0].Code, cat.list[1].Code, cat.list[2].Code})
	got, ok := cat.Get("a")
	require.True(t, ok)
	require.Equal(t, int32(2), got.SortOrder)
	_, ok = cat.Get("nope")
	require.False(t, ok)
}
