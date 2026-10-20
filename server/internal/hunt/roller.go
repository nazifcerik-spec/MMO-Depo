package hunt

import (
	"math"

	"github.com/otter/server/internal/combat"
)

// Hunt pacing (DD-281): fights run back to back, a turn takes TurnSeconds,
// so combat power sets how fast mobs fall. After a win a hunter camps until
// full at RestHealPctPerSec of max HP a second when its HP is under
// RestBelowPct of max or under three times what the fight cost (the next could
// kill); that time counts against the hunt too.
const (
	TurnSeconds       = 3
	RestBelowPct      = 50
	RestHealPctPerSec = 1
	// LootSeconds: gold and drop odds are per this much hunting time, so
	// faster kills do not multiply the economy (the pre-DD-281 cadence).
	LootSeconds = 120
	// MaxLoggedEncounters bounds the stored per-encounter log; the totals
	// cover every fight.
	MaxLoggedEncounters = 50
)

// EncounterRecord is one player-vs-mob outcome inside a hunt log.
type EncounterRecord struct {
	Index       int            `json:"index"`
	Outcome     combat.Outcome `json:"outcome"`
	Duration    int            `json:"duration_turns"`
	PlayerHPAft int32          `json:"player_hp_after"`
	MobHPAft    int32          `json:"mob_hp_after"`
	XPGained    int32          `json:"xp_gained"`
	GoldCopper  int64          `json:"gold_copper"`
}

// HuntResult is the aggregate + per-encounter log the reveal handler grants.
type HuntResult struct {
	Completed       int   `json:"completed_encounters"`
	Victories       int   `json:"victories"`
	Defeats         int   `json:"defeats"`
	Timeouts        int   `json:"timeouts"`
	StoppedOnDeath  bool  `json:"stopped_on_death"` // true if a defeat cut the hunt short
	TotalXP         int64 `json:"total_xp"`
	TotalGoldCopper int64 `json:"total_gold_copper"`
	CommonDrops     int   `json:"common_drops"`
	RareDrops       int   `json:"rare_drops"`
	LegendaryDrops  int   `json:"legendary_drops"`
	// Time spent: fighting plus camping, and the camping alone (DD-281).
	ElapsedSecs int32 `json:"elapsed_seconds"`
	RestSecs    int32 `json:"rest_seconds"`
	// FightTurns sums the turns of the won fights.
	FightTurns int `json:"fight_turns"`
	// LootRolls is the hunting time behind the wins in LootSeconds units:
	// the gear rolls the reveal makes.
	LootRolls     int               `json:"loot_rolls"`
	Encounters    []EncounterRecord `json:"encounters"` // the first MaxLoggedEncounters
	PlayerHPFinal int32             `json:"player_hp_final"`
}

// Roll fights `mob` back to back for budgetSecs of hunting time (DD-281):
// each fight takes its turns × TurnSeconds, a win pays the mob's XPReward
// (the caller applies the level-gap falloff), and a wounded hunter camps to
// full first (see RestBelowPct). A defeat or a timeout ends the hunt; a
// fight that would run past the budget is not fought. Deterministic given
// (baseSeed, player, mob, budgetSecs): fight N uses seed baseSeed + N + 1.
//
// HP carries over between fights (plus the player's RegenPerSec while it
// fights). Gold and drop odds scale with the time each win took, so an hour
// pays about the same loot however fast the mob dies.
func Roll(
	baseSeed int64,
	player combat.Combatant,
	mob combat.Mob,
	budgetSecs int32,
	maxTurnsPerEncounter int,
	rngFactory func(seed int64) combat.RNG,
) *HuntResult {
	res := &HuntResult{Encounters: []EncounterRecord{}}
	hp := player.HPCurrent
	var gold, loot float64
	for i := 0; res.ElapsedSecs < budgetSecs; i++ {
		before := hp
		enc := player.Clone() // skills + MP (full each fight)
		enc.Name, enc.HPCurrent = "player", hp
		sess := &combat.Session{Player: enc, Mob: mob.AsCombatant(), MaxTurns: maxTurnsPerEncounter}
		out, err := sess.Resolve(rngFactory, baseSeed+int64(i)+1)
		if err != nil {
			break // impossible in practice: the service validated the stats
		}
		fight := int32(out.Duration) * TurnSeconds
		if res.ElapsedSecs+fight > budgetSecs {
			break // the hunt ends mid-fight: that one does not count
		}
		res.ElapsedSecs += fight
		res.Completed++
		rec := EncounterRecord{Index: i + 1, Outcome: out.Outcome, Duration: out.Duration,
			PlayerHPAft: out.PlayerHPEnd, MobHPAft: out.MobHPEnd}

		stop := false
		switch out.Outcome {
		case combat.OutcomeVictory:
			res.Victories++
			res.FightTurns += out.Duration
			hp = min(player.HPMax, out.PlayerHPEnd+int32(player.RegenPerSec*float64(fight)))
			wall := fight
			if hp*100 < player.HPMax*RestBelowPct || hp < 3*(before-out.PlayerHPEnd) {
				rest := int32(math.Ceil(float64(player.HPMax-hp) * 100 / float64(player.HPMax*RestHealPctPerSec)))
				res.RestSecs += rest
				res.ElapsedSecs += rest
				wall += rest
				hp = player.HPMax
			}
			share := float64(wall) / LootSeconds
			loot += share
			rec.XPGained = mob.XPReward
			res.TotalXP += int64(mob.XPReward)
			g := float64(rollGoldInRange(rngFactory(baseSeed+int64(i)+9001), mob.GoldCopperMin, mob.GoldCopperMax)) * share
			gold += g
			rec.GoldCopper = int64(math.Round(g))
			// Drops (per LootSeconds of hunting): 60 common / 8 rare / 1
			// legendary out of 100, scaled by this win's share of that time.
			s := math.Min(1, share)
			switch d := float64(rngFactory(baseSeed+int64(i)+8001).Uint64() % 10000); {
			case d < 6000*s:
				res.CommonDrops++
			case d < 6800*s:
				res.RareDrops++
			case d < 6900*s:
				res.LegendaryDrops++
			}
		case combat.OutcomeDefeat:
			res.Defeats++
			res.StoppedOnDeath = true
			hp, stop = 0, true
		case combat.OutcomeTimeout:
			// Something is off with the pick (a mob nothing kills): stop
			// without a reward for it.
			res.Timeouts++
			hp, stop = out.PlayerHPEnd, true
		}
		if len(res.Encounters) < MaxLoggedEncounters {
			res.Encounters = append(res.Encounters, rec)
		}
		if stop {
			break
		}
	}
	res.TotalGoldCopper = int64(math.Round(gold))
	res.LootRolls = int(loot)
	res.PlayerHPFinal = hp
	return res
}

// rollGoldInRange returns a uniform int64 in [min, max] using one Uint64 roll.
// When min == max no roll is spent (returns min immediately).
func rollGoldInRange(r combat.RNG, min, max int64) int64 {
	if max <= min {
		return min
	}
	span := uint64(max - min + 1)
	return min + int64(r.Uint64()%span)
}
