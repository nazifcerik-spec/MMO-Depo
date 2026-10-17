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
	applyLevelGap(&s.Player, s.Mob.Level)
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

// Level-gap rule (DD-229): fighting a mob above your level costs accuracy
// and armor, so out-levelled content is dangerous whatever the gear.
const (
	GapHitPerLevel     = 4 // hit chance points lost per level of gap
	GapDefensePerLevel = 6 // defense percent lost per level of gap
)

func applyLevelGap(p *Combatant, mobLevel int16) {
	if p.Level <= 0 || mobLevel <= p.Level {
		return
	}
	gap := int32(mobLevel - p.Level)
	p.HitPct -= GapHitPerLevel * gap
	p.Defense = p.Defense * max(0, 100-GapDefensePerLevel*gap) / 100
}

// caster tracks one combatant's resource (MP or Rage), skill cooldowns and
// timed buffs for an encounter.
type caster struct {
	mp    int32
	rage  bool
	ready []int // turn from which skill i may be cast again
	atk   buff
	def   buff
}

// buff is a timed stat bonus: delta is reverted when its turn comes.
type buff struct {
	until int
	delta int32
}

func newCaster(c *Combatant) *caster {
	k := &caster{mp: c.MPMax, ready: make([]int, len(c.Skills)), rage: c.Resource == "rage"}
	if k.rage {
		k.mp = 0 // Rage builds up during the fight
	}
	return k
}

func (k *caster) expire(turn int, self *Combatant) {
	if k.atk.delta != 0 && turn >= k.atk.until {
		self.Attack -= k.atk.delta
		k.atk = buff{}
	}
	if k.def.delta != 0 && turn >= k.def.until {
		self.Defense -= k.def.delta
		k.def = buff{}
	}
}

// gain adds resource after an action: Rage per landed hit plus the skill's
// own gain, capped at the pool size.
func (k *caster) gain(self *Combatant, rec TurnRecord, sk *Skill) {
	if !k.rage {
		return
	}
	if rec.Hit && rec.Damage > 0 {
		k.mp += self.ResourcePerHit
	}
	if sk != nil {
		k.mp += sk.Gain
	}
	k.mp = min(k.mp, self.MPMax)
}

func pctAtOrBelow(cur, max, pct int32) bool { return int64(cur)*100 <= int64(max)*int64(pct) }

// act picks the first usable skill (list order = priority) or attacks.
func (k *caster) act(r RNG, turn int, self, target *Combatant) TurnRecord {
	k.expire(turn, self)
	for i := range self.Skills {
		sk := &self.Skills[i]
		if turn < k.ready[i] || sk.MPCost > k.mp {
			continue
		}
		if sk.SelfHPBelow > 0 && !pctAtOrBelow(self.HPCurrent, self.HPMax, sk.SelfHPBelow) {
			continue
		}
		if sk.TargetHPBelow > 0 && !pctAtOrBelow(target.HPCurrent, target.HPMax, sk.TargetHPBelow) {
			continue
		}
		switch sk.Effect {
		case EffectGuard, EffectEmpower:
			b := &k.def
			stat := &self.Defense
			if sk.Effect == EffectEmpower {
				b, stat = &k.atk, &self.Attack
			}
			if b.delta != 0 {
				continue // already up
			}
			delta := max(1, *stat*sk.Power/100)
			*stat += delta
			*b = buff{until: turn + max(1, sk.DurationTurns), delta: delta}
			k.spend(i, turn, *sk)
			k.gain(self, TurnRecord{}, sk)
			return TurnRecord{Turn: turn, Actor: self.Name, Target: self.Name, Hit: true, Skill: sk.Code,
				TargetHPAft: self.HPCurrent}
		case EffectHeal:
			if sk.SelfHPBelow == 0 && self.HPCurrent*2 > self.HPMax {
				continue // not worth it yet
			}
			heal := self.HPMax * sk.Power / 100
			if heal > self.HPMax-self.HPCurrent {
				heal = self.HPMax - self.HPCurrent
			}
			self.HPCurrent += heal
			k.spend(i, turn, *sk)
			return TurnRecord{Turn: turn, Actor: self.Name, Target: self.Name, Hit: true, Skill: sk.Code,
				Heal: heal, TargetHPAft: self.HPCurrent}
		case EffectStrike:
			k.spend(i, turn, *sk)
			rec := swing(r, turn, self, target, sk.Power)
			rec.Skill = sk.Code
			k.gain(self, rec, sk)
			return rec
		}
	}
	rec := resolveSwing(r, turn, self, target)
	k.gain(self, rec, nil)
	return rec
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
	dmg = progressionDamage(r, dmg, crit, attacker, defender)
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

// progressionDamage applies the DD-250 combat extras in a fixed order. With
// all of them zero it returns dmg untouched and draws no RNG.
func progressionDamage(r RNG, dmg int32, crit bool, a, d *Combatant) int32 {
	if a.CritDmgPct == 0 && a.ExecutePct == 0 && a.MissingHPDmgPct == 0 && a.LifestealPct == 0 &&
		d.DmgReductionPct == 0 && d.BlockPct == 0 {
		return dmg
	}
	bonus := int64(100)
	if crit && a.CritDmgPct > 0 {
		bonus += int64(a.CritDmgPct)
	}
	if a.ExecutePct > 0 && pctAtOrBelow(d.HPCurrent, d.HPMax, ExecuteThresholdPct) {
		bonus += int64(a.ExecutePct)
	}
	if a.MissingHPDmgPct > 0 && a.HPMax > 0 {
		bonus += int64(a.MissingHPDmgPct) * int64(a.HPMax-a.HPCurrent) / int64(a.HPMax)
	}
	out := int64(dmg) * bonus / 100
	if d.DmgReductionPct > 0 {
		out = out * int64(100-min(d.DmgReductionPct, 90)) / 100
	}
	if d.BlockPct > 0 && int32(r.Uint64()%100) < d.BlockPct {
		out /= 2
	}
	out = max(out, 1)
	if a.LifestealPct > 0 {
		heal := int32(out * int64(a.LifestealPct) / 100)
		a.HPCurrent = min(a.HPMax, a.HPCurrent+heal)
	}
	return int32(out)
}
