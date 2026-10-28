package raid

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/otter/server/internal/activity"
	"github.com/otter/server/internal/combat"
	"github.com/otter/server/internal/dungeon"
	"github.com/otter/server/internal/leveling"
)

type mobs map[string]combat.Mob

func (m mobs) Get(code string) (combat.Mob, bool) { v, ok := m[code]; return v, ok }

func rng(seed int64) combat.RNG { return activity.NewRNG(seed) }

var testMobs = mobs{
	"add":  {Code: "add", Level: 10, HPMax: 100, Attack: 20, Defense: 2, HitPct: 85, GoldCopperMin: 10, GoldCopperMax: 20},
	"boss": {Code: "boss", Level: 12, HPMax: 3000, Attack: 60, Defense: 5, HitPct: 90, RaidBoss: true},
}

var testRaid = dungeon.Template{Code: "r", MinLevel: 10, DurationSecs: 3600, GroupMin: 3, GroupMax: 5,
	Waves: []dungeon.Wave{{MobCode: "add", Count: 2}}, BossMobCode: "boss", BonusGoldCopper: 500, BonusXP: 1000}

func hero(hp, atk int32) combat.Combatant {
	return combat.Combatant{Name: "player", HPMax: hp, HPCurrent: hp, Attack: atk, Defense: 10, HitPct: 90, Level: 10,
		RegenPerEncounter: 20}
}

func TestRollClearsAndIsDeterministic(t *testing.T) {
	party := []combat.Combatant{hero(800, 60), hero(600, 70), hero(500, 80)}
	a := Roll(42, party, testRaid, testMobs, rng)
	b := Roll(42, party, testRaid, testMobs, rng)
	require.Equal(t, a, b, "same snapshots and seed, same raid")
	require.True(t, a.Cleared)
	require.True(t, a.BossKilled)
	require.False(t, a.Wiped)
	require.Equal(t, 3, a.Victories)
	require.Len(t, a.Encounters, 3)
	require.True(t, a.Encounters[2].IsBoss)
	require.Len(t, a.HPFinal, 3)
	require.Equal(t, int32(800), party[0].HPCurrent, "the inputs are not mutated")
	for _, e := range a.Encounters[:2] {
		require.GreaterOrEqual(t, e.GoldCopper, int64(10))
	}
	require.Equal(t, "p2", Label(2))
	require.Equal(t, combat.EndCleared, a.Report.End)
	require.Equal(t, 3, a.Report.Fights)
	// Each member sees its own stats against the same enemy (DD-308).
	me := a.Report.For(party[2])
	require.Equal(t, party[2].HPMax, me.You.HPMax)
	require.Equal(t, a.Report.Yours, me.Yours)
}

func TestRollWipeStopsTheRun(t *testing.T) {
	weak := []combat.Combatant{hero(50, 5), hero(50, 5), hero(50, 5)}
	res := Roll(7, weak, testRaid, testMobs, rng)
	require.True(t, res.Wiped)
	require.False(t, res.Cleared)
	require.Equal(t, combat.EndDefeat, res.Report.End)
	require.NotEmpty(t, res.Report.EndMob)
	require.Equal(t, []int32{0, 0, 0}, res.HPFinal)
	require.Len(t, res.Encounters, 1)

	// An unknown mob ends the run with what it has.
	broken := testRaid
	broken.BossMobCode = "missing"
	res = Roll(42, []combat.Combatant{hero(800, 60), hero(600, 70), hero(500, 80)}, broken, testMobs, rng)
	require.False(t, res.Cleared)
	require.Equal(t, 2, res.Victories)
	broken.Waves = []dungeon.Wave{{MobCode: "missing", Count: 1}}
	require.Empty(t, Roll(1, []combat.Combatant{hero(800, 60)}, broken, testMobs, rng).Encounters)
	// Nobody standing: the encounter cannot run.
	res = Roll(1, []combat.Combatant{{HPMax: 10}}, testRaid, testMobs, rng)
	require.Equal(t, combat.OutcomeTimeout, res.Encounters[0].Outcome)
}

func TestMemberReward(t *testing.T) {
	res := Roll(42, []combat.Combatant{hero(800, 60), hero(600, 70), hero(500, 80)}, testRaid, testMobs, rng)
	require.True(t, res.Cleared)
	rw := MemberReward(res, 1, 10, "warrior", 42, testRaid, testMobs, rng)
	perKill := int64(afkXPPerKill(testRaid))
	require.Equal(t, int64(leveling.Current().AFKCombatXPPerMinute*60), perKill*3)
	require.Equal(t, 3*perKill+1000, rw.XP, "three kills at the AFK rate plus the boss bonus")
	gold := res.Encounters[0].GoldCopper + res.Encounters[1].GoldCopper + res.Encounters[2].GoldCopper + 500
	require.Equal(t, gold, rw.GoldCopper)
	require.GreaterOrEqual(t, rw.Rare, 1, "the boss always drops a rare")
	require.Equal(t, activity.GearDrop{Class: "warrior", Level: 10, Rolls: 2, Pct: GearDropPct, Guaranteed: BossGear, Boss: true, SetSource: "r"}, rw.Gear) // + a set piece (DD-277)

	other := MemberReward(res, 2, 10, "mage", 42, testRaid, testMobs, rng)
	require.Equal(t, rw.GoldCopper, other.GoldCopper, "gold is not split")
	require.Equal(t, rw.XP, other.XP)

	// No kills, nothing.
	wiped := Roll(7, []combat.Combatant{hero(50, 5)}, testRaid, testMobs, rng)
	require.Equal(t, Reward{}, MemberReward(wiped, 1, 10, "warrior", 7, testRaid, testMobs, rng))
	require.Zero(t, afkXPPerKill(dungeon.Template{}))
}
