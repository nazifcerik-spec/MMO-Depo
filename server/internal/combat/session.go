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
	// FleeBelowPct: the player breaks off (OutcomeFled) once a mob's swing
	// leaves it alive under this share of max HP. Zero = fight to the end.
	FleeBelowPct int32
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
				PlayerMPEnd: caster.mp,
			}, nil
		}
		// --- Mob swings (unless stunned) ---
		if caster.stunned(turn) {
			log = append(log, TurnRecord{Turn: turn, Actor: s.Mob.Name, Target: s.Player.Name, Stunned: true,
				TargetHPAft: s.Player.HPCurrent})
			continue
		}
		log = append(log, resolveSwing(r, turn, &s.Mob, &s.Player))
		if s.Player.HPCurrent <= 0 {
			return &Result{
				Outcome:     OutcomeDefeat,
				Turns:       log,
				PlayerHPEnd: 0,
				MobHPEnd:    s.Mob.HPCurrent,
				Duration:    turn,
				PlayerMPEnd: caster.mp,
			}, nil
		}
		if s.FleeBelowPct > 0 && int64(s.Player.HPCurrent)*100 < int64(s.Player.HPMax)*int64(s.FleeBelowPct) {
			return &Result{
				Outcome:     OutcomeFled,
				Turns:       log,
				PlayerHPEnd: s.Player.HPCurrent,
				MobHPEnd:    s.Mob.HPCurrent,
				Duration:    turn,
				PlayerMPEnd: caster.mp,
			}, nil
		}
	}
	return &Result{
		Outcome:     OutcomeTimeout,
		Turns:       log,
		PlayerHPEnd: s.Player.HPCurrent,
		MobHPEnd:    s.Mob.HPCurrent,
		Duration:    s.MaxTurns,
		PlayerMPEnd: caster.mp,
	}, nil
}

// Level-gap rule (DD-229): fighting a mob above your level costs accuracy
// and armor, so out-levelled content is dangerous whatever the gear. Past
// GapAttackGrace levels it also costs attack (DD-254), so healing classes
// cannot outlast far stronger mobs; dungeon bosses a few levels up are
// unaffected.
const (
	GapHitPerLevel     = 4  // hit chance points lost per level of gap
	GapDefensePerLevel = 6  // defense percent lost per level of gap
	GapAttackGrace     = 5  // levels of gap before attack is lost
	GapAttackPerLevel  = 10 // attack percent lost per level beyond the grace
)

func applyLevelGap(p *Combatant, mobLevel int16) {
	if p.Level <= 0 || mobLevel <= p.Level {
		return
	}
	gap := int32(mobLevel - p.Level)
	p.HitPct -= GapHitPerLevel * gap
	p.Defense = p.Defense * max(0, 100-GapDefensePerLevel*gap) / 100
	if over := gap - GapAttackGrace; over > 0 {
		p.Attack = max(1, p.Attack*max(10, 100-GapAttackPerLevel*over)/100)
	}
}

// caster tracks one combatant's resource, skill cooldowns and timed
// effects for an encounter.
type caster struct {
	mp     int32
	pooled bool  // a class resource (DD-252): generated in combat
	ready  []int // turn from which skill i may be cast again
	atk    buff
	def    buff
	weak   buff // on the target's defense
	dot    tick // on the target
	hot    tick
	shield int // turn the barrier fades (0 = none)
	stun   int // the target skips its swings while turn < stun

	// Group fights (DD-261): heals, heal-over-time and barriers go to an ally
	// (self in a duel); stun skills only strike a stun-immune target.
	hotOn    *Combatant
	shieldOn *Combatant
	noStun   bool
}

// buff is a timed stat bonus: delta is reverted when its turn comes.
type buff struct {
	until int
	delta int32
	on    bool
}

// tick is a per-turn effect with turns left.
type tick struct {
	left   int
	amount int32
}

func newCaster(c *Combatant) *caster {
	k := &caster{mp: c.MPMax, ready: make([]int, len(c.Skills)), pooled: c.Resource != ""}
	if k.pooled {
		k.mp = min(max(c.ResourceStart, 0), c.MPMax) // Rage, Chi, ... build up during the fight
	}
	if c.MPStart != nil {
		k.mp = min(max(*c.MPStart, 0), c.MPMax)
	}
	return k
}

func (k *caster) stunned(turn int) bool { return turn < k.stun }

func (k *caster) expire(turn int, self, target *Combatant) {
	if k.atk.delta != 0 && turn >= k.atk.until {
		self.Attack -= k.atk.delta
		k.atk = buff{}
	}
	if k.def.delta != 0 && turn >= k.def.until {
		self.Defense -= k.def.delta
		k.def = buff{}
	}
	if k.weak.on && turn >= k.weak.until {
		target.Defense += k.weak.delta
		k.weak = buff{}
	}
	if k.shield != 0 && turn >= k.shield {
		k.shieldOn.Barrier, k.shield = 0, 0
	}
}

// ticks applies damage- and heal-over-time at the start of a turn.
func (k *caster) ticks(target *Combatant) (dmg, heal int32) {
	if k.dot.left > 0 {
		dmg = min(k.dot.amount, target.HPCurrent)
		target.HPCurrent -= dmg
		k.dot.left--
	}
	if k.hot.left > 0 {
		if on := k.hotOn; on.HPCurrent > 0 { // the fallen are not healed back up
			heal = min(k.hot.amount, on.HPMax-on.HPCurrent)
			on.HPCurrent += heal
		}
		k.hot.left--
	}
	return dmg, heal
}

// gain adds resource after an action: per landed hit plus the skill's own
// gain, capped at the pool size. The v1 mana pool does not regenerate.
func (k *caster) gain(self *Combatant, rec TurnRecord, sk *Skill) {
	if !k.pooled {
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

// setupFloorPct: attack buffs and armor debuffs are not cast against a
// target at or below this HP share (DD-252).
const setupFloorPct = 40

func pctAtOrBelow(cur, max, pct int32) bool { return int64(cur)*100 <= int64(max)*int64(pct) }

// act runs the start-of-turn effects, then picks the first usable skill
// (list order = priority) or attacks.
func (k *caster) act(r RNG, turn int, self, target *Combatant) TurnRecord {
	return k.actWith(r, turn, self, target, self)
}

// actWith is act with ally as the receiver of heals and barriers.
func (k *caster) actWith(r RNG, turn int, self, target, ally *Combatant) TurnRecord {
	k.expire(turn, self, target)
	if k.pooled && self.ResourcePerTurn > 0 {
		k.mp = min(k.mp+self.ResourcePerTurn, self.MPMax)
	}
	dmg, heal := k.ticks(target)
	if target.HPCurrent <= 0 {
		return TurnRecord{Turn: turn, Actor: self.Name, Target: target.Name, Hit: true, Tick: dmg, TickHeal: heal}
	}
	rec := k.cast(r, turn, self, target, ally)
	rec.Tick, rec.TickHeal = dmg, heal
	return rec
}

// supportive effects land on the ally rather than the caster.
func supportive(effect string) bool {
	return effect == EffectHeal || effect == EffectHot || effect == EffectBarrier
}

// cast uses the first usable skill or attacks.
func (k *caster) cast(r RNG, turn int, self, target, ally *Combatant) TurnRecord {
	for i := range self.Skills {
		sk := &self.Skills[i]
		if turn < k.ready[i] || sk.MPCost > k.mp {
			continue
		}
		benef := self
		if supportive(sk.Effect) {
			benef = ally
		}
		if sk.SelfHPBelow > 0 && !pctAtOrBelow(benef.HPCurrent, benef.HPMax, sk.SelfHPBelow) {
			continue
		}
		if sk.TargetHPBelow > 0 && !pctAtOrBelow(target.HPCurrent, target.HPMax, sk.TargetHPBelow) {
			continue
		}
		if (sk.Effect == EffectEmpower || sk.Effect == EffectWeaken) && pctAtOrBelow(target.HPCurrent, target.HPMax, setupFloorPct) {
			continue // the target is nearly down: finish it instead of setting up
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
			if sk.SelfHPBelow == 0 && benef.HPCurrent*2 > benef.HPMax {
				continue // not worth it yet
			}
			heal := benef.HPMax * sk.Power / 100
			if heal > benef.HPMax-benef.HPCurrent {
				heal = benef.HPMax - benef.HPCurrent
			}
			benef.HPCurrent += heal
			k.spend(i, turn, *sk)
			return TurnRecord{Turn: turn, Actor: self.Name, Target: benef.Name, Hit: true, Skill: sk.Code,
				Heal: heal, TargetHPAft: benef.HPCurrent}
		case EffectStrike, EffectStun:
			if sk.Effect == EffectStun && !k.noStun && k.stunned(turn+1) {
				continue // still stunned: no chaining
			}
			k.spend(i, turn, *sk)
			rec := swing(r, turn, self, target, sk.Power)
			rec.Skill = sk.Code
			if sk.Effect == EffectStun && rec.Hit && !k.noStun {
				k.stun = turn + max(1, sk.DurationTurns)
			}
			k.gain(self, rec, sk)
			return rec
		case EffectBarrier:
			if k.shield != 0 {
				continue
			}
			benef.Barrier = max(1, benef.HPMax*sk.Power/100)
			k.shield, k.shieldOn = turn+max(1, sk.DurationTurns), benef
			k.spend(i, turn, *sk)
			k.gain(self, TurnRecord{}, sk)
			return TurnRecord{Turn: turn, Actor: self.Name, Target: benef.Name, Hit: true, Skill: sk.Code,
				TargetHPAft: benef.HPCurrent}
		case EffectHot:
			if k.hot.left > 0 || (sk.SelfHPBelow == 0 && benef.HPCurrent*10 > benef.HPMax*7) {
				continue // running already, or not hurt enough (default: <= 70% HP)
			}
			amount := max(1, benef.HPMax*sk.Power/100)
			heal := min(amount, benef.HPMax-benef.HPCurrent)
			benef.HPCurrent += heal
			k.hot, k.hotOn = tick{left: max(1, sk.DurationTurns) - 1, amount: amount}, benef
			k.spend(i, turn, *sk)
			k.gain(self, TurnRecord{}, sk)
			return TurnRecord{Turn: turn, Actor: self.Name, Target: benef.Name, Hit: true, Skill: sk.Code,
				Heal: heal, TargetHPAft: benef.HPCurrent}
		case EffectDot:
			if k.dot.left > 0 {
				continue
			}
			raw := float64(self.Attack) - float64(target.Defense)*defenseMitigation
			amount := max(1, int32(max(raw, MinDamageFloor)*float64(sk.Power)/100))
			first := min(amount, target.HPCurrent)
			target.HPCurrent -= first
			k.dot = tick{left: max(1, sk.DurationTurns) - 1, amount: amount}
			k.spend(i, turn, *sk)
			rec := TurnRecord{Turn: turn, Actor: self.Name, Target: target.Name, Hit: true, Skill: sk.Code,
				Damage: first, TargetHPAft: target.HPCurrent}
			k.gain(self, rec, sk)
			return rec
		case EffectWeaken:
			if k.weak.on {
				continue
			}
			delta := target.Defense * sk.Power / 100
			target.Defense -= delta
			k.weak = buff{until: turn + max(1, sk.DurationTurns), delta: delta, on: true}
			k.spend(i, turn, *sk)
			k.gain(self, TurnRecord{}, sk)
			return TurnRecord{Turn: turn, Actor: self.Name, Target: target.Name, Hit: true, Skill: sk.Code,
				TargetHPAft: target.HPCurrent}
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
	if defender.Barrier > 0 {
		rec.Absorbed = min(defender.Barrier, dmg)
		defender.Barrier -= rec.Absorbed
		dmg -= rec.Absorbed
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
