package partyxp

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

var t0 = time.Date(2026, 10, 27, 12, 0, 0, 0, time.UTC)

func hours(h float64) time.Time { return t0.Add(time.Duration(h * float64(time.Hour))) }

func work(from, to float64, xp int64) Work {
	return Work{Kind: GatherPrefix + "fishing", Character: uuid.New(), Start: hours(from), End: hours(to), XP: xp}
}

func TestFactor(t *testing.T) {
	require.InDelta(t, 1, Factor(1), 1e-9, "alone: no bonus")
	require.InDelta(t, 1.06, Factor(2), 1e-9)
	require.InDelta(t, 1.15, Factor(5), 1e-9, "the user's example: five members, +15%")
}

func TestBoostAlone(t *testing.T) {
	require.EqualValues(t, 600, Boost(work(0, 1, 600), nil))
	// A partner working at another time adds nothing.
	require.EqualValues(t, 600, Boost(work(0, 1, 600), [][]Interval{{{hours(2), hours(3)}}}))
}

func TestBoostFullOverlap(t *testing.T) {
	// Two together: +6% each, nobody gives anything away.
	require.EqualValues(t, 1060, Boost(work(0, 1, 1000), [][]Interval{{{hours(0), hours(1)}}}))
	// Five together: +15%.
	all := []Interval{{hours(-1), hours(2)}}
	require.EqualValues(t, 1150, Boost(work(0, 1, 1000), [][]Interval{all, all, all, all}))
}

// A six-hour job beside a one-hour one gets the bonus only on the hour they
// overlap.
func TestBoostPartialOverlap(t *testing.T) {
	require.EqualValues(t, 606, Boost(work(0, 6, 600), [][]Interval{{{hours(2), hours(3)}}})) // 500 + 100 × 1.06
}

func TestBoostStaggered(t *testing.T) {
	// B overlaps the first half, C the second; never three at once.
	b := []Interval{{hours(-1), hours(0.5)}}
	c := []Interval{{hours(0.5), hours(2)}}
	require.EqualValues(t, 1060, Boost(work(0, 1, 1000), [][]Interval{b, c}))
	// Two, then three, then two, then alone: 106 + 109 + 106 + 100.
	b = []Interval{{hours(0), hours(2)}}
	c = []Interval{{hours(1), hours(3)}}
	require.EqualValues(t, 421, Boost(work(0, 4, 400), [][]Interval{b, c}))
}

func TestBoostNothing(t *testing.T) {
	p := [][]Interval{{{hours(0), hours(1)}}}
	require.EqualValues(t, 0, Boost(work(0, 1, 0), p))
	require.EqualValues(t, 50, Boost(work(1, 1, 50), p), "a job without length keeps its XP")
}

func TestBoostIfWithoutBooster(t *testing.T) {
	require.Equal(t, Solo(70), BoostIf(context.Background(), nil, work(0, 1, 70)))
}

func TestWindowsQuery(t *testing.T) {
	for kind, arg := range map[string]string{KindHunt: "", KindDungeon: "", GatherPrefix + "mining": "mining",
		CraftPrefix + "smelting": "smelting"} {
		sql, got, err := windowsQuery(kind)
		require.NoError(t, err, kind)
		require.NotEmpty(t, sql)
		require.Equal(t, arg, got, kind)
	}
	_, _, err := windowsQuery("raid")
	require.Error(t, err)
}
