package combat

// Session is one player-vs-mob encounter. Resolve mutates the internal HP
// pools; call Snapshot beforehand if the caller needs the original inputs.
//
// Turn order: player attacks first. On a killing blow the loop ends before the
// mob's counter — so a glass-cannon player CAN one-shot without eating a hit.
type Session struct {
	Player   Combatant
	Mob      Combatant
	MaxTurns int // hard cap; prevents pathological non-lethal loops from running forever
}

// DefaultMaxTurns keeps a session bounded when the caller supplies zero.
// 60 turns is enough for a level-adjacent hunt at the current damage scale;
// dungeons (Sprint 15) pass a higher cap.
const DefaultMaxTurns = 60

// Resolve runs the whole encounter deterministically. seed drives the RNG:
// same (session, seed) always produces the same log. rngNew is the factory —
// injecting it here (instead of hard-coding activity.NewRNG) keeps combat free
// of an upward dependency on activity. Callers pass activity.NewRNG.
//
// The returned Result includes the full turn log and outcome; the session's
// Combatants are left in their end state (HPCurrent floored at zero for the
// losing side).
func (s *Session) Resolve(rngNew func(seed int64) RNG, seed int64) (*Result, error) {
	if err := s.validate(); err != nil {
		return nil, err
	}
	if s.MaxTurns <= 0 {
		s.MaxTurns = DefaultMaxTurns
	}
	r := rngNew(seed)
	log := make([]TurnRecord, 0, 2*s.MaxTurns)
	caster := newCaster(&s.Player)

	for turn := 1; turn <= s.MaxTurns; turn++ {
		// --- Player acts (skill if one fits, else basic attack) ---
		log = append(log, caster.act(r, turn, &s.Player, &s.Mob))
		if s.Mob.HPCurrent <= 0 {
			return &Result{
				Outcome:     OutcomeVictory,
				Turns:       log,
				PlayerHPEnd: s.Player.HPCurrent,
				MobHPEnd:    0,
				Duration:    turn,
			}, nil
		}
		// --- Mob swings ---
		log = append(log, resolveSwing(r, turn, &s.Mob, &s.Player))
		if s.Player.HPCurrent <= 0 {
			return &Result{
				Outcome:     OutcomeDefeat,
				Turns:       log,
				PlayerHPEnd: 0,
				MobHPEnd:    s.Mob.HPCurrent,
				Duration:    turn,
			}, nil
		}
	}
	return &Result{
		Outcome:     OutcomeTimeout,
		Turns:       log,
		PlayerHPEnd: s.Player.HPCurrent,
		MobHPEnd:    s.Mob.HPCurrent,
		Duration:    s.MaxTurns,
	}, nil
}

// caster tracks one combatant's MP and skill cooldowns for an encounter.
type caster struct {
	mp    int32
	ready []int // turn from which skill i may be cast again
}

func newCaster(c *Combatant) *caster {
	return &caster{mp: c.MPMax, ready: make([]int, len(c.Skills))}
}

// act picks the first usable skill (list order = priority) or attacks.
func (k *caster) act(r RNG, turn int, self, target *Combatant) TurnRecord {
	for i, sk := range self.Skills {
		if turn < k.ready[i] || sk.MPCost > k.mp {
			continue
		}
		switch sk.Effect {
		case EffectHeal:
			if self.HPCurrent*2 > self.HPMax {
				continue // not worth it yet
			}
			heal := self.HPMax * sk.Power / 100
			if heal > self.HPMax-self.HPCurrent {
				heal = self.HPMax - self.HPCurrent
			}
			self.HPCurrent += heal
			k.spend(i, turn, sk)
			return TurnRecord{Turn: turn, Actor: self.Name, Target: self.Name, Hit: true, Skill: sk.Code,
				Heal: heal, TargetHPAft: self.HPCurrent}
		case EffectStrike:
			k.spend(i, turn, sk)
			rec := swing(r, turn, self, target, sk.Power)
			rec.Skill = sk.Code
			return rec
		}
	}
	return resolveSwing(r, turn, self, target)
}

func (k *caster) spend(i, turn int, sk Skill) {
	k.mp -= sk.MPCost
	k.ready[i] = turn + max(1, sk.CooldownTurns)
}

// resolveSwing is one attacker→defender exchange. Mutates defender.HPCurrent.
// Split out so the loop above stays readable — this is the only piece where
// TryHit + RollDamage are wired together.
func resolveSwing(r RNG, turn int, attacker, defender *Combatant) TurnRecord {
	return swing(r, turn, attacker, defender, 100)
}

// swing is an attack at powerPct% of normal damage (100 = basic attack).
func swing(r RNG, turn int, attacker, defender *Combatant, powerPct int32) TurnRecord {
	hit, crit := TryHit(r, attacker.HitPct, defender.EvasionPct, attacker.CritPct)
	rec := TurnRecord{
		Turn:   turn,
		Actor:  attacker.Name,
		Target: defender.Name,
		Hit:    hit,
		Crit:   crit,
	}
	if !hit {
		rec.TargetHPAft = defender.HPCurrent
		return rec
	}
	dmg := RollDamage(r, attacker.Attack, defender.Defense, crit)
	if powerPct != 100 {
		dmg = max(1, dmg*powerPct/100)
	}
	defender.HPCurrent -= dmg
	if defender.HPCurrent < 0 {
		defender.HPCurrent = 0
	}
	rec.Damage = dmg
	rec.TargetHPAft = defender.HPCurrent
	return rec
}

// validate refuses obviously broken sessions early. Balance-level "is this
// fair" is a Sprint 14+ concern; here we only reject inputs that would make
// the loop misbehave (both sides invulnerable, negative starting HP, ...).
func (s *Session) validate() error {
	if s.Player.HPCurrent <= 0 {
		return invalid("player HPCurrent must be > 0")
	}
	if s.Mob.HPCurrent <= 0 {
		return invalid("mob HPCurrent must be > 0")
	}
	if s.Player.Attack <= 0 && s.Mob.Attack <= 0 {
		return invalid("both combatants have zero attack — deadlock")
	}
	return nil
}
