package dungeon

import (
	"github.com/otter/server/internal/combat"
)

// EncounterRecord is one wave-mob or boss combat outcome in the run log.
type EncounterRecord struct {
	Index       int            `json:"index"`
	Wave        int            `json:"wave"` // 1..N; 0 for boss
	IsBoss      bool           `json:"is_boss,omitempty"`
	MobCode     string         `json:"mob_code"`
	Outcome     combat.Outcome `json:"outcome"`
	Duration    int            `json:"duration_turns"`
	PlayerHPAft int32          `json:"player_hp_after"`
	MobHPAft    int32          `json:"mob_hp_after"`
	XPGained    int32          `json:"xp_gained"`
	GoldCopper  int64          `json:"gold_copper"`
}

// RunResult is the aggregate + per-encounter log for a dungeon run.
type RunResult struct {
	Cleared         bool              `json:"cleared"` // all waves + boss killed
	StoppedOnDeath  bool              `json:"stopped_on_death"`
	Victories       int               `json:"victories"`
	Defeats         int               `json:"defeats"`
	Timeouts        int               `json:"timeouts"`
	TotalXP         int64             `json:"total_xp"`
	TotalGoldCopper int64             `json:"total_gold_copper"`
	CommonDrops     int               `json:"common_drops"`
	RareDrops       int               `json:"rare_drops"`
	LegendaryDrops  int               `json:"legendary_drops"`
	Encounters      []EncounterRecord `json:"encounters"`
	PlayerHPFinal   int32             `json:"player_hp_final"`
	BossKilled      bool              `json:"boss_killed"`
}

// MobLookup is the subset of combat.Catalog the roller needs — kept as an
// interface so tests can pass a fake without touching the DB.
type MobLookup interface {
	Get(code string) (combat.Mob, bool)
}

// Roll walks the dungeon deterministically. Encounter N uses seed = baseSeed +
// N + 1 (matches hunt.Roll's stream discipline).
//
// Player HP carries over between every encounter (waves AND boss). A defeat
// halts the run immediately, sets StoppedOnDeath=true, and leaves Cleared and
// BossKilled false. A timeout is treated the same way — something is wrong
// with the balance choice, so drop the run rather than infinite-loop.
//
// Wave-mob drop odds are the same table hunt uses (60/8/1). The boss adds a
// guaranteed rare drop on kill + Template.BonusGoldCopper + Template.BonusXP.
// Item.Grant + XP.GrantXP happen upstream in the reveal path — this fn stays
// pure so tests assert exact counts.
func Roll(
	baseSeed int64,
	player combat.Combatant,
	tmpl Template,
	mobs MobLookup,
	maxTurnsPerEncounter int,
	rngFactory func(seed int64) combat.RNG,
) *RunResult {
	res := &RunResult{Encounters: make([]EncounterRecord, 0, tmpl.TotalEncounters())}
	if len(tmpl.Waves) == 0 {
		res.PlayerHPFinal = player.HPCurrent
		return res
	}

	seq := 0
	// --- Waves ---
	for wIdx, w := range tmpl.Waves {
		for i := 0; i < w.Count; i++ {
			seq++
			mob, ok := mobs.Get(w.MobCode)
			if !ok {
				// Content bug — treat as run failure so the reveal still returns
				// what it has instead of throwing.
				res.PlayerHPFinal = player.HPCurrent
				return res
			}
			rec, killed, ended := resolveOne(baseSeed, seq, player, mob, wIdx+1, false, maxTurnsPerEncounter, rngFactory)
			res.Encounters = append(res.Encounters, rec)
			switch rec.Outcome {
			case combat.OutcomeVictory:
				res.Victories++
				res.TotalXP += int64(rec.XPGained)
				res.TotalGoldCopper += rec.GoldCopper
				accumulateDrop(res, rngFactory(baseSeed+int64(seq)+8001))
				// A breather before the next pack (DD-229).
				player.HPCurrent = min(player.HPMax, rec.PlayerHPAft+player.RegenPerEncounter)
			case combat.OutcomeDefeat:
				res.Defeats++
				res.StoppedOnDeath = true
				res.PlayerHPFinal = 0
				_ = killed // unused; kept for symmetry with boss branch
				return res
			case combat.OutcomeTimeout:
				res.Timeouts++
				res.PlayerHPFinal = rec.PlayerHPAft
				player.HPCurrent = rec.PlayerHPAft
				return res
			}
			if ended {
				return res
			}
		}
	}

	// --- Boss ---
	bossMob, ok := mobs.Get(tmpl.BossMobCode)
	if !ok {
		res.PlayerHPFinal = player.HPCurrent
		return res
	}
	seq++
	rec, _, _ := resolveOne(baseSeed, seq, player, bossMob, 0, true, maxTurnsPerEncounter, rngFactory)
	res.Encounters = append(res.Encounters, rec)
	switch rec.Outcome {
	case combat.OutcomeVictory:
		res.Victories++
		res.TotalXP += int64(rec.XPGained) + int64(tmpl.BonusXP)
		res.TotalGoldCopper += rec.GoldCopper + tmpl.BonusGoldCopper
		res.BossKilled = true
		res.Cleared = true
		// Guaranteed rare from the boss + one wave-table roll for a bonus item.
		res.RareDrops++
		accumulateDrop(res, rngFactory(baseSeed+int64(seq)+8001))
		res.PlayerHPFinal = rec.PlayerHPAft
	case combat.OutcomeDefeat:
		res.Defeats++
		res.StoppedOnDeath = true
		res.PlayerHPFinal = 0
	case combat.OutcomeTimeout:
		res.Timeouts++
		res.PlayerHPFinal = rec.PlayerHPAft
	}
	return res
}

// resolveOne runs one encounter and returns its record + control flags. `wave`
// is 0 for the boss to match EncounterRecord.Wave semantics.
// Boss scaling over the mob template (percent).
const (
	BossHPPct     = 200
	BossAttackPct = 120
)

func resolveOne(
	baseSeed int64, seq int,
	player combat.Combatant, mob combat.Mob,
	wave int, isBoss bool, maxTurns int,
	rngFactory func(seed int64) combat.RNG,
) (rec EncounterRecord, killed bool, ended bool) {
	enc := player.Clone() // HP carries over; skills + MP (full each fight) too
	enc.Name = "player"
	mobC := mob.AsCombatant()
	if isBoss { // a dungeon boss outclasses the same mob in the wild (DD-229)
		mobC.HPMax = mobC.HPMax * BossHPPct / 100
		mobC.HPCurrent = mobC.HPMax
		mobC.Attack = mobC.Attack * BossAttackPct / 100
	}
	sess := &combat.Session{Player: enc, Mob: mobC, MaxTurns: maxTurns}
	out, err := sess.Resolve(rngFactory, baseSeed+int64(seq)+1)
	if err != nil {
		return EncounterRecord{
			Index: seq, Wave: wave, IsBoss: isBoss, MobCode: mob.Code,
			Outcome:     combat.OutcomeTimeout, // treat validate errors as timeouts
			PlayerHPAft: player.HPCurrent, MobHPAft: mob.HPMax,
		}, false, true
	}
	rec = EncounterRecord{
		Index: seq, Wave: wave, IsBoss: isBoss, MobCode: mob.Code,
		Outcome: out.Outcome, Duration: out.Duration,
		PlayerHPAft: out.PlayerHPEnd, MobHPAft: out.MobHPEnd,
	}
	if out.Outcome == combat.OutcomeVictory {
		rec.XPGained = mob.XPReward
		rec.GoldCopper = rollGoldInRange(rngFactory(baseSeed+int64(seq)+9001), mob.GoldCopperMin, mob.GoldCopperMax)
		killed = true
	}
	return rec, killed, false
}

// accumulateDrop mirrors hunt's per-victory rarity table (60/8/1) so drops
// across the two AFK loops feel equivalent per kill.
func accumulateDrop(res *RunResult, r combat.RNG) {
	d := r.Uint64() % 100
	switch {
	case d < 60:
		res.CommonDrops++
	case d < 68:
		res.RareDrops++
	case d < 69:
		res.LegendaryDrops++
	}
}

func rollGoldInRange(r combat.RNG, min, max int64) int64 {
	if max <= min {
		return min
	}
	span := uint64(max - min + 1)
	return min + int64(r.Uint64()%span)
}
