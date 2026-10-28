package raid

import (
	"math"

	"github.com/otter/server/internal/activity"
	"github.com/otter/server/internal/combat"
	"github.com/otter/server/internal/dungeon"
	"github.com/otter/server/internal/leveling"
)

// Fight tuning (DD-261).
const (
	AddHPPct     = 300 // wave mobs are sized for a group
	BossAoEEvery = 3   // every third boss swing hits the whole group …
	BossAoEPct   = 60  // … at 60% power
	MaxTurns     = 150 // per encounter
	BossGear     = 2   // guaranteed gear pieces from the boss
	GearDropPct  = 4   // gear chance per wave kill (as in dungeons)
)

// GroupHPPct and GroupAttackPct size every enemy of a raid to the group
// that fights it (DD-298), indexed by member count (2..5): every size clears
// about 60% of runs at the raid's level (balance.TestRaidTargets), where a
// full party used to clear nearly all and a pair none.
var (
	GroupHPPct     = [6]int32{46, 46, 46, 77, 104, 127}
	GroupAttackPct = [6]int32{73, 73, 73, 89, 102, 113}
)

// groupScale is the scaling row for n members (clamped to 1..5).
func groupScale(n int) (hp, attack int32) {
	n = max(1, min(n, len(GroupHPPct)-1))
	return GroupHPPct[n], GroupAttackPct[n]
}

// MobLookup resolves mob templates (combat.Catalog satisfies).
type MobLookup interface {
	Get(code string) (combat.Mob, bool)
}

// Encounter is one fight of the raid log.
type Encounter struct {
	Index      int            `json:"index"`
	Wave       int            `json:"wave"` // 1..N; 0 for the boss
	IsBoss     bool           `json:"is_boss,omitempty"`
	MobCode    string         `json:"mob_code"`
	Outcome    combat.Outcome `json:"outcome"`
	Duration   int            `json:"duration_turns"`
	HPAfter    []int32        `json:"hp_after"` // per member, in slot order
	GoldCopper int64          `json:"gold_copper"`
}

// Result is the shared outcome of a raid: the same for every member.
type Result struct {
	Cleared    bool        `json:"cleared"`
	Wiped      bool        `json:"wiped"`
	Victories  int         `json:"victories"`
	BossKilled bool        `json:"boss_killed"`
	Encounters []Encounter `json:"encounters"`
	HPFinal    []int32     `json:"hp_final"`
}

// Roll fights the waves and then the boss with the whole group. HP carries
// over between encounters (the living rest RegenPerEncounter); the fallen
// stay down. The run stops at a wipe or a timeout. Encounter n draws from
// seed+n+1, so the same snapshots and seed always give the same raid.
func Roll(seed int64, players []combat.Combatant, tmpl dungeon.Template, mobs MobLookup, rngNew func(int64) combat.RNG) *Result {
	ps := make([]combat.Combatant, len(players))
	for i := range players {
		ps[i] = players[i].Clone()
		ps[i].Name = Label(int16(i + 1))
	}
	res := &Result{}
	hpPct, atkPct := groupScale(len(ps))
	seq := 0
	fight := func(mob combat.Mob, wave int, boss bool) bool {
		seq++
		enc := make([]combat.Combatant, len(ps))
		for i := range ps {
			enc[i] = ps[i].Clone() // buffs and the level gap stay inside the encounter
		}
		g := &combat.Group{Players: enc, Mob: mob.AsRaidCombatant(), MaxTurns: MaxTurns}
		if boss {
			g.AoEEvery, g.AoEPct, g.StunImmune = BossAoEEvery, BossAoEPct, true
		} else {
			g.Mob.HPMax = g.Mob.HPMax * AddHPPct / 100
		}
		g.Mob.HPMax = max(1, int32(int64(g.Mob.HPMax)*int64(hpPct)/100))
		g.Mob.HPCurrent = g.Mob.HPMax
		g.Mob.Attack = max(1, int32(int64(g.Mob.Attack)*int64(atkPct)/100))
		rec := Encounter{Index: seq, Wave: wave, IsBoss: boss, MobCode: mob.Code}
		out, err := g.Resolve(rngNew, seed+int64(seq)+1)
		if err != nil { // nobody standing: content or input bug, end the run
			rec.Outcome = combat.OutcomeTimeout
			rec.HPAfter = hpOf(ps)
			res.Encounters = append(res.Encounters, rec)
			return false
		}
		rec.Outcome, rec.Duration, rec.HPAfter = out.Outcome, out.Duration, out.PlayersHPEnd
		for i := range ps {
			ps[i].HPCurrent = out.PlayersHPEnd[i]
		}
		if out.Outcome == combat.OutcomeVictory {
			res.Victories++
			rec.GoldCopper = rollGold(rngNew(seed+int64(seq)+9001), mob.GoldCopperMin, mob.GoldCopperMax)
		}
		res.Encounters = append(res.Encounters, rec)
		if out.Outcome != combat.OutcomeVictory {
			res.Wiped = out.Outcome == combat.OutcomeDefeat
			return false
		}
		for i := range ps { // a breather before the next pack (DD-229)
			if ps[i].HPCurrent > 0 {
				ps[i].HPCurrent = min(ps[i].HPMax, ps[i].HPCurrent+ps[i].RegenPerEncounter)
			}
		}
		return true
	}
	done := func() *Result { res.HPFinal = hpOf(ps); return res }
	for w, wave := range tmpl.Waves {
		mob, ok := mobs.Get(wave.MobCode)
		if !ok {
			return done()
		}
		for i := 0; i < wave.Count; i++ {
			if !fight(mob, w+1, false) {
				return done()
			}
		}
	}
	boss, ok := mobs.Get(tmpl.BossMobCode)
	if !ok {
		return done()
	}
	if fight(boss, 0, true) {
		res.BossKilled, res.Cleared = true, true
	}
	return done()
}

func hpOf(ps []combat.Combatant) []int32 {
	out := make([]int32, len(ps))
	for i := range ps {
		out[i] = max(0, ps[i].HPCurrent)
	}
	return out
}

func rollGold(r combat.RNG, lo, hi int64) int64 {
	if hi <= lo {
		return lo
	}
	return lo + int64(r.Uint64()%uint64(hi-lo+1))
}

// Reward is one member's personal share of a raid.
type Reward struct {
	XP         int64             `json:"xp"`
	GoldCopper int64             `json:"gold_copper"`
	Common     int               `json:"common_drops"`
	Rare       int               `json:"rare_drops"`
	Legendary  int               `json:"legendary_drops"`
	Gear       activity.GearDrop `json:"-"`
}

// MemberReward is what the member in slot (1-based) at level earns: every
// kill pays the AFK combat XP rate scaled to the member's level (DD-250,
// DD-254) and the mob's gold; the boss adds the template's bonus, a rare
// item and BossGear guaranteed gear pieces. Drop rolls are the member's own.
func MemberReward(res *Result, slot int16, level int, class string, seed int64, tmpl dungeon.Template,
	mobs MobLookup, rngNew func(int64) combat.RNG) Reward {
	var rw Reward
	cfg := leveling.Current()
	perKill := afkXPPerKill(tmpl)
	r := rngNew(seed + int64(slot)*7919 + 8001)
	waveKills := 0
	for _, e := range res.Encounters {
		if e.Outcome != combat.OutcomeVictory {
			continue
		}
		mob, _ := mobs.Get(e.MobCode)
		rw.XP += int64(cfg.ScaleDungeonXP(perKill, level, int(mob.Level))) // each member's own level (DD-303)
		rw.GoldCopper += e.GoldCopper
		switch d := r.Uint64() % 100; {
		case d < 60:
			rw.Common++
		case d < 68:
			rw.Rare++
		case d < 69:
			rw.Legendary++
		}
		if e.IsBoss {
			rw.XP += int64(cfg.ScaleDungeonXP(tmpl.BonusXP, level, int(mob.Level)))
			rw.GoldCopper += tmpl.BonusGoldCopper
			rw.Rare++
		} else {
			waveKills++
		}
	}
	if res.Victories > 0 {
		boss, _ := mobs.Get(tmpl.BossMobCode)
		rw.Gear = activity.GearDrop{Class: class, Level: activity.DropLevel(int(boss.Level), level),
			Rolls: waveKills, Pct: GearDropPct}
		if res.BossKilled {
			rw.Gear.Guaranteed, rw.Gear.Boss = BossGear, true
			rw.Gear.SetSource = tmpl.Code // and a piece of the raid's set (DD-277)
		}
	}
	return rw
}

// afkXPPerKill spreads the raid's minutes × afk_combat_xp_per_minute over
// its encounters (as dungeons do).
func afkXPPerKill(tmpl dungeon.Template) int32 {
	n := tmpl.TotalEncounters()
	if n <= 0 {
		return 0
	}
	return int32(math.Round(leveling.Current().AFKCombatXPPerMinute * float64(tmpl.DurationSecs) / 60 / float64(n)))
}
