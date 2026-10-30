package pet

import (
	"math"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/otter/server/internal/activity"
	"github.com/otter/server/internal/combat"
)

func rng(seed int64) combat.RNG { return activity.NewRNG(seed) }

// species is the shipped ladder's shape: 20 ranks, 50% → 2%.
func species() []Species {
	out := make([]Species, 20)
	for i := range out {
		t := float64(i) / 19
		out[i] = Species{
			Code: string(rune('a' + i)), Rank: i + 1, MinMobLevel: 400 + 30*(i+1),
			EggPct: 50 * math.Pow(2.0/50, t), HatchSeconds: 3600, PowerPct: 200 + 3*i,
			Bonus: Modifier{Kind: "mult_add", Target: "attack", Amount: 0.03 + 0.09*t},
		}
	}
	return out
}

// curveMob is a normal-archetype mob of the level (gen_world.py formulas).
func curveMob(level int) combat.Mob {
	l := float64(level)
	return combat.Mob{Code: "m", Level: int16(level), HPMax: int32(150 + 25*l), Attack: int32(28 + 2.1*l),
		Defense: int32(1 + 0.8*l), CritPct: 5, HitPct: min(92, 82+int32(level)/5), EvasionPct: 5,
		XPReward: int32(level)}
}

func TestCatalogValidation(t *testing.T) {
	c, err := New(species())
	require.NoError(t, err)
	require.Equal(t, 20, c.All()[0].Rank, "best first")
	s, ok := c.ByEgg("pet_egg_a")
	require.True(t, ok)
	require.Equal(t, 1, s.Rank)

	bad := species()
	bad[1].Rank = 1
	_, err = New(bad)
	require.Error(t, err, "duplicate rank")
	bad = species()
	bad[0].EggPct = 0
	_, err = New(bad)
	require.Error(t, err)
}

// The best egg drops at its own odds, the weakest about at 50% where only it
// is in reach, never one where none is, and at most one per roll.
func TestRollEggOdds(t *testing.T) {
	c, err := New(species())
	require.NoError(t, err)
	const n = 200000
	count := func(level int) map[string]int {
		got := map[string]int{}
		r := rng(7)
		for i := 0; i < n; i++ {
			got[c.RollEgg(level, r)]++
		}
		return got
	}
	require.Equal(t, n, count(400)[""], "nothing below level 430")
	low := count(430)
	require.InDelta(t, 0.50, float64(low["pet_egg_a"])/n, 0.01)
	top := count(1000)
	require.InDelta(t, 0.02, float64(top["pet_egg_t"])/n, 0.002)
	mid := count(602) // the level 600 raid: ranks 1..6
	require.Zero(t, mid["pet_egg_g"])
	require.Positive(t, mid["pet_egg_f"])
}

func TestXPAndLevels(t *testing.T) {
	require.EqualValues(t, 1000, XPToNext(1))
	require.Zero(t, XPToNext(MaxLevel))
	l, xp := AddXP(1, 0, 1000+XPToNext(2)+5)
	require.Equal(t, 3, l)
	require.EqualValues(t, 5, xp)
	l, xp = AddXP(99, 0, 1<<50)
	require.Equal(t, MaxLevel, l)
	require.Zero(t, xp)

	s := species()[19]
	require.InDelta(t, 0.03, s.BonusAt(1).Amount, 0.0001)
	require.InDelta(t, 0.12, s.BonusAt(MaxLevel).Amount, 0.0001)
	require.True(t, CanHunt(1, curveMob(15)))
	require.False(t, CanHunt(1, curveMob(16)))
	require.False(t, CanHunt(100, combat.Mob{Level: 10, RaidBoss: true}))
}

// TestPetPace: an hour on even-level mobs earns about a level's XP at the
// start of the curve, for the weakest and the best species (PET_PACE=1 -v
// prints the table).
func TestPetPace(t *testing.T) {
	for _, s := range []Species{species()[0], species()[19]} {
		for _, level := range []int{1, 10, 50, 100} {
			mob := curveMob(Power(level))
			res := Roll(int64(level), s, level, mob, 3600, rng)
			hours := float64(XPToNext(min(level, MaxLevel-1))) / float64(res.TotalXP)
			if os.Getenv("PET_PACE") != "" {
				t.Logf("rank %2d level %3d: %3d wins/h, %v died, %.1f h per level", s.Rank, level, res.Victories,
					res.StoppedOnDeath, hours)
			}
			require.False(t, res.StoppedOnDeath, "rank %d level %d dies on even mobs", s.Rank, level)
			if level == 1 {
				require.InDelta(t, 1, hours, 0.6, "about an hour a level at the start")
			}
			// The top of its reach is still a fight it wins.
			top := Roll(int64(level), s, level, curveMob(Power(level)+MobReach), 3600, rng)
			require.Positive(t, top.Victories, "rank %d level %d at its reach", s.Rank, level)
		}
	}
}
