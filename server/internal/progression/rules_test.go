package progression

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

var curve = []int64{100, 500, 2000, 8000, 30000, 100000, 300000, 800000, 2000000, 5000000}

func TestMasteryTier_CumulativeCurve(t *testing.T) {
	require.Equal(t, int16(0), MasteryTier(0, curve, 10))
	require.Equal(t, int16(0), MasteryTier(99, curve, 10))
	require.Equal(t, int16(1), MasteryTier(100, curve, 10))
	require.Equal(t, int16(1), MasteryTier(599, curve, 10))
	require.Equal(t, int16(2), MasteryTier(600, curve, 10))
	require.Equal(t, int16(10), MasteryTier(9_000_000, curve, 10))
	require.Equal(t, int16(3), MasteryTier(9_000_000, curve, 3), "capped at max tier")
}

func TestNextTierXP(t *testing.T) {
	require.Equal(t, int64(100), NextTierXP(0, curve, 10))
	require.Equal(t, int64(600), NextTierXP(1, curve, 10))
	require.Equal(t, int64(-1), NextTierXP(10, curve, 10))
}

func TestRepTier(t *testing.T) {
	require.Equal(t, "neutral", RepTierNames[RepTier(0)])
	require.Equal(t, "neutral", RepTierNames[RepTier(2999)])
	require.Equal(t, "friendly", RepTierNames[RepTier(3000)])
	require.Equal(t, "exalted", RepTierNames[RepTier(RepMax)])
	require.Equal(t, "unfriendly", RepTierNames[RepTier(-1)])
	require.Equal(t, "hated", RepTierNames[RepTier(RepMin)])
}

func TestApplyRepGain_DailyCapAndClamp(t *testing.T) {
	rep, used := ApplyRepGain(0, 1400, 300, 1500)
	require.Equal(t, 100, rep, "only the remaining 100 of today's cap applies")
	require.Equal(t, 1500, used)

	rep, used = ApplyRepGain(100, 1500, 50, 1500)
	require.Equal(t, 100, rep, "cap exhausted")
	require.Equal(t, 1500, used)

	rep, used = ApplyRepGain(100, 1500, -500, 1500)
	require.Equal(t, -400, rep, "losses ignore the cap")
	require.Equal(t, 1500, used)

	rep, _ = ApplyRepGain(RepMax-10, 0, 100, 1500)
	require.Equal(t, RepMax, rep)
	rep, _ = ApplyRepGain(RepMin+10, 0, -100, 1500)
	require.Equal(t, RepMin, rep)
}

func TestSameUTCDay(t *testing.T) {
	a := time.Date(2026, 3, 1, 23, 59, 0, 0, time.UTC)
	require.True(t, SameUTCDay(a, a.Add(30*time.Second)))
	require.False(t, SameUTCDay(a, a.Add(2*time.Minute)))
}

type spy struct{ got []Event }

func (s *spy) Record(_ context.Context, _ uuid.UUID, e ...Event) { s.got = append(s.got, e...) }

func TestRecordIf_DropsZeroAndNil(t *testing.T) {
	RecordIf(context.Background(), nil, uuid.New(), Event{Kind: "x", Amount: 1}) // no panic
	s := &spy{}
	RecordIf(context.Background(), s, uuid.New(), Event{Kind: "a", Amount: 0}, Event{Kind: "b", Amount: 2})
	require.Equal(t, []Event{{Kind: "b", Amount: 2}}, s.got)
	s.got = nil
	RecordIf(context.Background(), s, uuid.New(), Event{Kind: "a", Amount: 0})
	require.Empty(t, s.got, "nothing recorded when all amounts are zero")
}

func TestActiveEffectsAndSumTarget(t *testing.T) {
	curve := []int64{100, 500, 2000}
	cat := &Catalog{tracks: []track{
		{Code: "mining", MaxTier: 3, Curve: curve, Bonuses: []tierBonus{
			{Tier: 1, Effects: []BonusEffect{{Kind: "flat", Target: "gather_yield_mining", Amount: 1}}},
			{Tier: 3, Effects: []BonusEffect{{Kind: "mult_add", Target: "gather_yield_mining", Amount: 0.1}}},
		}},
		{Code: "combat", MaxTier: 3, Curve: curve, Bonuses: []tierBonus{
			{Tier: 2, Effects: []BonusEffect{{Kind: "flat", Target: "attack", Amount: 2}}},
		}},
	}}

	require.Empty(t, cat.ActiveEffects(nil), "no xp → nothing unlocked")

	got := cat.ActiveEffects(map[string]int64{"mining": 2600, "combat": 599})
	require.Len(t, got, 2, "mining tier 3 unlocks both entries; combat tier 1 none")
	flat, mult := SumTarget(got, "gather_yield_mining")
	require.Equal(t, 1.0, flat)
	require.InDelta(t, 0.1, mult, 1e-9)

	got = cat.ActiveEffects(map[string]int64{"mining": 150, "combat": 600})
	flat, mult = SumTarget(got, "gather_yield_mining")
	require.Equal(t, 1.0, flat)
	require.Zero(t, mult)
	flat, _ = SumTarget(got, "attack")
	require.Equal(t, 2.0, flat)
	flat, _ = SumTarget(got, "gather_yield_fishing")
	require.Zero(t, flat)
}
