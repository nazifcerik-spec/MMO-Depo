package worldevent

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestNextOccurrences(t *testing.T) {
	// 2026-09-22 is a Tuesday (weekday 2).
	now := time.Date(2026, 9, 22, 10, 0, 0, 0, time.UTC)
	got := NextOccurrences(now, 2, 19, 24*time.Hour)
	require.Equal(t, []time.Time{time.Date(2026, 9, 22, 19, 0, 0, 0, time.UTC)}, got)

	require.Empty(t, NextOccurrences(now, 2, 9, 24*time.Hour), "earlier today already passed, next week is out of horizon")
	require.Empty(t, NextOccurrences(now, 4, 19, 24*time.Hour), "Thursday is beyond 24h")

	got = NextOccurrences(now, 3, 5, 24*time.Hour) // Wednesday 05:00 is within 24h
	require.Equal(t, []time.Time{time.Date(2026, 9, 23, 5, 0, 0, 0, time.UTC)}, got)

	got = NextOccurrences(now, 2, 19, 8*24*time.Hour)
	require.Len(t, got, 2, "wider horizon picks up next week too")
}

func TestPowerAndContribution(t *testing.T) {
	require.Equal(t, 30*2+15+25, Power(30, 15, 250))
	require.Equal(t, 1, Power(0, 0, 0))
	require.Equal(t, int64(6000), Contribution(100, 60, 10000))
	require.Equal(t, int64(10000), Contribution(500, 60, 10000), "AFK capped at tier-3 floor")
	require.Equal(t, int64(0), Contribution(100, -5, 10000))
}

func TestAttendedMinutes(t *testing.T) {
	start := time.Date(2026, 1, 1, 19, 0, 0, 0, time.UTC)
	end := start.Add(time.Hour)
	require.Equal(t, int64(60), AttendedMinutes(start.Add(-time.Hour), start, end))
	require.Equal(t, int64(30), AttendedMinutes(start.Add(30*time.Minute), start, end), "late join")
	require.Equal(t, int64(0), AttendedMinutes(end.Add(time.Minute), start, end))
}

func TestTierFor(t *testing.T) {
	tiers := []Tier{{MinScore: 1}, {MinScore: 1000}, {MinScore: 10000}, {MinScore: 50000}}
	require.Equal(t, 0, TierFor(0, tiers, true))
	require.Equal(t, 1, TierFor(1, tiers, true))
	require.Equal(t, 2, TierFor(9999, tiers, true))
	require.Equal(t, 3, TierFor(10000, tiers, true))
	require.Equal(t, 4, TierFor(60000, tiers, true))
	require.Equal(t, FailedTierCap, TierFor(60000, tiers, false), "failed events cap rewards")
	require.Equal(t, 1, TierFor(5, tiers, false))
}
