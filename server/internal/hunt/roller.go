package hunt

import (
	"github.com/otter/server/internal/combat"
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
	// TurnLog is intentionally omitted here — the reveal response keeps
	// per-encounter detail behind a `verbose=true` flag (Phase 02). For now
	// only the summary is persisted.
}

// HuntResult is the aggregate + per-encounter log the reveal handler grants.
type HuntResult struct {
	PlannedEncounters int               `json:"planned_encounters"`
	Completed         int               `json:"completed_encounters"`
	Victories         int               `json:"victories"`
	Defeats           int               `json:"defeats"`
	Timeouts          int               `json:"timeouts"`
	StoppedOnDeath    bool              `json:"stopped_on_death"` // true if a defeat cut the hunt short
	TotalXP           int64             `json:"total_xp"`
	TotalGoldCopper   int64             `json:"total_gold_copper"`
	CommonDrops       int               `json:"common_drops"`
	RareDrops         int               `json:"rare_drops"`
	LegendaryDrops    int               `json:"legendary_drops"`
	Encounters        []EncounterRecord `json:"encounters"`
	PlayerHPFinal     int32             `json:"player_hp_final"`
}

// Roll runs `planned` encounters against the same mob template with a fresh
// combat.Session each time. Deterministic given (baseSeed, player, mob,
// planned): encounter N uses seed = baseSeed + N + 1 so the sequence is stable
// across retries.
//
// The player Combatant is NOT reset between encounters — HPCurrent carries
// over so a hunt is a stamina test, not a series of isolated fights. A single
// defeat cuts the hunt short with StoppedOnDeath = true and the remaining
// planned encounters are not executed (nor rewarded).
//
// Drop rolls are deliberately CHEAP here: one per victory, weighted from the
// per-victory rarity table. Item Grant + XP grant happen upstream in the
// reveal handler — this fn stays pure so tests can assert exact counts.
func Roll(
	baseSeed int64,
	player combat.Combatant,
	mob combat.Mob,
	planned int,
	maxTurnsPerEncounter int,
	rngFactory func(seed int64) combat.RNG,
) *HuntResult {
	res := &HuntResult{
		PlannedEncounters: planned,
		Encounters:        make([]EncounterRecord, 0, planned),
	}
	if planned <= 0 {
		res.PlayerHPFinal = player.HPCurrent
		return res
	}

	for i := 0; i < planned; i++ {
		enc := combat.Combatant{
			Name:       "player",
			HPMax:      player.HPMax,
			HPCurrent:  player.HPCurrent, // carry over
			Attack:     player.Attack,
			Defense:    player.Defense,
			CritPct:    player.CritPct,
			HitPct:     player.HitPct,
			EvasionPct: player.EvasionPct,
		}
		mobC := mob.AsCombatant()

		sess := &combat.Session{Player: enc, Mob: mobC, MaxTurns: maxTurnsPerEncounter}
		// Skip validate errors — Roll is only called by the service which has
		// already verified stats via Sheet + mob.AsCombatant.
		out, err := sess.Resolve(rngFactory, baseSeed+int64(i)+1)
		if err != nil {
			// Impossible in practice; break so partial log is still returned.
			break
		}

		rec := EncounterRecord{
			Index:       i + 1,
			Outcome:     out.Outcome,
			Duration:    out.Duration,
			PlayerHPAft: out.PlayerHPEnd,
			MobHPAft:    out.MobHPEnd,
		}

		switch out.Outcome {
		case combat.OutcomeVictory:
			res.Victories++
			rec.XPGained = mob.XPReward
			rec.GoldCopper = rollGoldInRange(rngFactory(baseSeed+int64(i)+9001), mob.GoldCopperMin, mob.GoldCopperMax)
			res.TotalXP += int64(rec.XPGained)
			res.TotalGoldCopper += rec.GoldCopper
			// Drop rarity table (per victory): 60 common / 8 rare / 1 legendary
			// out of 100 rolls. Nothing → 31% empty-hand chance.
			d := rngFactory(baseSeed+int64(i)+8001).Uint64() % 100
			switch {
			case d < 60:
				res.CommonDrops++
			case d < 68:
				res.RareDrops++
			case d < 69:
				res.LegendaryDrops++
			}
			// Carry over the surviving player HP.
			player.HPCurrent = out.PlayerHPEnd
		case combat.OutcomeDefeat:
			res.Defeats++
			res.StoppedOnDeath = true
			res.Encounters = append(res.Encounters, rec)
			res.Completed = i + 1
			res.PlayerHPFinal = 0
			return res
		case combat.OutcomeTimeout:
			res.Timeouts++
			// Timeout still ends the hunt: something is wrong with the balance
			// pick (mob was tanky, no lethal happened). Award nothing for this
			// encounter and stop.
			res.Encounters = append(res.Encounters, rec)
			res.Completed = i + 1
			res.PlayerHPFinal = out.PlayerHPEnd
			player.HPCurrent = out.PlayerHPEnd
			return res
		}

		res.Encounters = append(res.Encounters, rec)
		res.Completed = i + 1
	}
	res.PlayerHPFinal = player.HPCurrent
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
