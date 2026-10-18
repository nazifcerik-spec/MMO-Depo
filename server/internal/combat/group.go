package combat

// Group is a party-vs-mob encounter (raids, DD-261). Each round every living
// player acts in slice order — heals, heal-over-time and barriers go to the
// most wounded living ally — then the mob swings at the toughest living
// player (highest defense: the tank holds it), or, every AoEEvery-th swing,
// hits every living player at AoEPct% power. Fallen players stay down.
type Group struct {
	Players    []Combatant
	Mob        Combatant
	MaxTurns   int
	AoEEvery   int   // 0 = no area attacks
	AoEPct     int32 // power of an area attack, percent of a basic swing
	StunImmune bool  // stun skills only deal their strike damage
}

// GroupResult is the outcome of Group.Resolve.
type GroupResult struct {
	Outcome      Outcome      `json:"outcome"`
	Turns        []TurnRecord `json:"turns"`
	PlayersHPEnd []int32      `json:"players_hp_end"`
	MobHPEnd     int32        `json:"mob_hp_end"`
	Duration     int          `json:"duration_turns"`
}

// SkillAoE marks an area attack in the log.
const SkillAoE = "aoe"

// Resolve runs the encounter deterministically (same inputs and seed, same
// log). Players are left in their end state.
func (g *Group) Resolve(rngNew func(seed int64) RNG, seed int64) (*GroupResult, error) {
	if g.Mob.HPCurrent <= 0 {
		return nil, invalid("mob HPCurrent must be > 0")
	}
	if alive(g.Players) == 0 {
		return nil, invalid("no living player")
	}
	if g.MaxTurns <= 0 {
		g.MaxTurns = DefaultMaxTurns
	}
	r := rngNew(seed)
	casters := make([]*caster, len(g.Players))
	for i := range g.Players {
		applyLevelGap(&g.Players[i], g.Mob.Level)
		casters[i] = newCaster(&g.Players[i])
		casters[i].noStun = g.StunImmune
	}
	log := make([]TurnRecord, 0, (len(g.Players)+1)*g.MaxTurns)
	end := func(o Outcome, turn int) *GroupResult {
		hp := make([]int32, len(g.Players))
		for i := range g.Players {
			hp[i] = max(0, g.Players[i].HPCurrent)
		}
		return &GroupResult{Outcome: o, Turns: log, PlayersHPEnd: hp, MobHPEnd: max(0, g.Mob.HPCurrent), Duration: turn}
	}
	swings := 0
	for turn := 1; turn <= g.MaxTurns; turn++ {
		for i := range g.Players {
			if g.Players[i].HPCurrent <= 0 {
				continue
			}
			log = append(log, casters[i].actWith(r, turn, &g.Players[i], &g.Mob, mostWounded(g.Players)))
			if g.Mob.HPCurrent <= 0 {
				return end(OutcomeVictory, turn), nil
			}
		}
		if stunnedBy(casters, g.Players, turn) {
			log = append(log, TurnRecord{Turn: turn, Actor: g.Mob.Name, Stunned: true})
			continue
		}
		swings++
		if g.AoEEvery > 0 && swings%g.AoEEvery == 0 {
			for i := range g.Players {
				if g.Players[i].HPCurrent > 0 {
					rec := swing(r, turn, &g.Mob, &g.Players[i], g.AoEPct)
					rec.Skill = SkillAoE
					log = append(log, rec)
				}
			}
		} else {
			log = append(log, resolveSwing(r, turn, &g.Mob, tank(g.Players)))
		}
		if alive(g.Players) == 0 {
			return end(OutcomeDefeat, turn), nil
		}
	}
	return end(OutcomeTimeout, g.MaxTurns), nil
}

func alive(ps []Combatant) int {
	n := 0
	for i := range ps {
		if ps[i].HPCurrent > 0 {
			n++
		}
	}
	return n
}

// mostWounded is the living player with the lowest HP share (first on ties).
func mostWounded(ps []Combatant) *Combatant {
	var best *Combatant
	for i := range ps {
		p := &ps[i]
		if p.HPCurrent <= 0 {
			continue
		}
		if best == nil || int64(p.HPCurrent)*int64(best.HPMax) < int64(best.HPCurrent)*int64(p.HPMax) {
			best = p
		}
	}
	return best
}

// tank is the living player with the highest defense (first on ties).
func tank(ps []Combatant) *Combatant {
	var best *Combatant
	for i := range ps {
		if p := &ps[i]; p.HPCurrent > 0 && (best == nil || p.Defense > best.Defense) {
			best = p
		}
	}
	return best
}

// stunnedBy reports whether a living player's stun holds the mob this turn.
func stunnedBy(cs []*caster, ps []Combatant, turn int) bool {
	for i, k := range cs {
		if ps[i].HPCurrent > 0 && k.stunned(turn) {
			return true
		}
	}
	return false
}
