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
	return Work{Kind: GatherPrefix + "fishing", Character: uuid.New(), Ref: uuid.New(), Start: hours(from), End: hours(to), XP: xp}
}

func TestSplitAlone(t *testing.T) {
	keep, gifts := Split(work(0, 1, 600), nil)
	require.EqualValues(t, 600, keep)
	require.Empty(t, gifts)
	// A partner working another time shares nothing.
	keep, gifts = Split(work(0, 1, 600), [][]Interval{{{hours(2), hours(3)}}})
	require.EqualValues(t, 600, keep)
	require.Equal(t, []int64{0}, gifts)
}

func TestSplitFullOverlap(t *testing.T) {
	// Two: half each, plus 3%.
	keep, gifts := Split(work(0, 1, 1000), [][]Interval{{{hours(0), hours(1)}}})
	require.EqualValues(t, 515, keep)
	require.Equal(t, []int64{515}, gifts)
	// Four: a quarter each, plus 9%.
	all := []Interval{{hours(-1), hours(2)}}
	keep, gifts = Split(work(0, 1, 1000), [][]Interval{all, all, all})
	require.EqualValues(t, 273, keep) // 250 × 1.09 = 272.5
	require.Equal(t, []int64{273, 273, 273}, gifts)
}

// The user's example: a six-hour job beside a one-hour one shares only the
// hour they overlap; the other five stay the earner's own.
func TestSplitPartialOverlap(t *testing.T) {
	keep, gifts := Split(work(0, 6, 600), [][]Interval{{{hours(2), hours(3)}}})
	// 500 alone + the shared hour's 100 / 2 × 1.03 = 51.5.
	require.EqualValues(t, 552, keep)
	require.Equal(t, []int64{52}, gifts)
}

func TestSplitStaggered(t *testing.T) {
	// B overlaps the first half, C the second; at no time are three working.
	b := []Interval{{hours(-1), hours(0.5)}}
	c := []Interval{{hours(0.5), hours(2)}}
	keep, gifts := Split(work(0, 1, 1000), [][]Interval{b, c})
	require.EqualValues(t, 515, keep)
	require.Equal(t, []int64{258, 258}, gifts) // 257.5 each
	// Two windows of one partner both count.
	keep, gifts = Split(work(0, 4, 400), [][]Interval{{{hours(0), hours(1)}, {hours(3), hours(4)}}})
	require.EqualValues(t, 303, keep) // 200 + 2 × 51.5
	require.Equal(t, []int64{103}, gifts)
}

func TestSplitNothingToShare(t *testing.T) {
	p := [][]Interval{{{hours(0), hours(1)}}}
	keep, gifts := Split(work(0, 1, 0), p)
	require.EqualValues(t, 0, keep)
	require.Equal(t, []int64{0}, gifts)
	keep, _ = Split(work(1, 1, 50), p)
	require.EqualValues(t, 50, keep, "a job without length keeps its XP")
}

func TestShareIfWithoutSharer(t *testing.T) {
	out := ShareIf(context.Background(), nil, work(0, 1, 70))
	require.Equal(t, Solo(70), out)
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
