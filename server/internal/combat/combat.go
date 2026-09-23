// Package combat implements the damage pipeline and session resolver.
//
// Design (Sprint 13):
//   - Server-authoritative. All rolls happen server-side; client only receives
//     the outcome + turn log.
//   - Deterministic. Given the same seed + inputs, Resolve returns the exact
//     same TurnRecord slice. RNG is xoshiro256++ (via activity/rng — same
//     algorithm the analytical sampler uses so behaviors line up).
//   - Pure pipeline. TryHit and RollDamage take an RNG and no other state;
//     the Session wraps them into an alternating-turn loop.
//   - Session-based. Idle MMO does not need per-turn player input — a whole
//     encounter resolves in one call. Sprint 14 hunting + Sprint 15 dungeon use
//     this same primitive.
package combat

import "fmt"

// RNG is the minimal deterministic PRNG surface combat needs. activity.RNG
// satisfies it. Kept as an interface here so combat has no upward dependency
// on the activity package.
type RNG interface {
	Uint64() uint64
}

// --- Types ------------------------------------------------------------------

// Combatant is the minimal input a single-encounter resolver needs. Built from
// character.Sheet.Derived / Resources for the player side; from Mob for the
// enemy side.
type Combatant struct {
	Name       string // display / debug
	HPMax      int32  // full pool
	HPCurrent  int32  // live pool; mutated during Resolve
	Attack     int32
	Defense    int32
	CritPct    int32 // 0..100 (integer percent)
	HitPct     int32 // 0..100 base to-hit before defender evasion
	EvasionPct int32 // 0..100 subtracted from attacker HitPct

	// Class skills (DD-222). MP starts full every encounter; mobs have none.
	MPMax  int32
	Skills []Skill
}

// Skill effects.
const (
	EffectStrike = "strike" // hit roll as usual, damage × Power/100
	EffectHeal   = "heal"   // restore Power% of HPMax; only used at <= 50% HP
)

// Skill is one active class skill, auto-cast on the combatant's turn: the
// first skill in list order that is off cooldown, affordable and useful is
// cast; otherwise a basic attack.
type Skill struct {
	Code          string
	Effect        string
	Power         int32 // strike: damage percent (170 = 1.7×); heal: % of HPMax
	MPCost        int32
	CooldownTurns int // turns until it can be cast again (>= 1)
}

// Clone returns a copy whose Skills slice is not shared.
func (c Combatant) Clone() Combatant {
	c.Skills = append([]Skill(nil), c.Skills...)
	return c
}

// TurnRecord is one line of the combat log — enough for the client to render
// a play-by-play + enough for debugging.
type TurnRecord struct {
	Turn        int    `json:"turn"`
	Actor       string `json:"actor"` // attacker name
	Target      string `json:"target"`
	Hit         bool   `json:"hit"`
	Crit        bool   `json:"crit,omitempty"`
	Skill       string `json:"skill,omitempty"` // class skill used this swing
	Heal        int32  `json:"heal,omitempty"`  // HP restored (heal skills)
	Damage      int32  `json:"damage"`
	TargetHPAft int32  `json:"target_hp_after"` // post-application, floored at 0
}

// Outcome is the final result of a Session.Resolve.
type Outcome string

const (
	OutcomeVictory Outcome = "victory" // player killed the mob
	OutcomeDefeat  Outcome = "defeat"  // mob killed the player
	OutcomeTimeout Outcome = "timeout" // hit MaxTurns with both alive
)

// Result bundles what Resolve produces.
type Result struct {
	Outcome     Outcome      `json:"outcome"`
	Turns       []TurnRecord `json:"turns"`
	PlayerHPEnd int32        `json:"player_hp_end"`
	MobHPEnd    int32        `json:"mob_hp_end"`
	Duration    int          `json:"duration_turns"`
}

// --- Pipeline (pure) --------------------------------------------------------

// TryHit returns (didHit, wasCrit). Miss when a uniform 0..99 roll >= effective
// hit chance. Crit is a second uniform roll over the attacker's crit_pct.
//
// Effective hit = clamp(attackerHit - defenderEvasion, minHitFloor..maxHitCap).
// Floor prevents "unhittable" defenders (>95 evasion) turning combat into
// stalemate; cap prevents "always hits" from bypassing an intended miss chance.
func TryHit(r RNG, attackerHitPct, defenderEvasionPct, attackerCritPct int32) (didHit bool, wasCrit bool) { //nolint:revive
	eff := attackerHitPct - defenderEvasionPct
	if eff < MinHitPct {
		eff = MinHitPct
	}
	if eff > MaxHitPct {
		eff = MaxHitPct
	}
	roll := int32(r.Uint64() % 100)
	if roll >= eff {
		return false, false
	}
	if attackerCritPct <= 0 {
		return true, false
	}
	critRoll := int32(r.Uint64() % 100)
	return true, critRoll < attackerCritPct
}

// RollDamage returns the applied damage after defense mitigation. Never below
// MinDamageFloor even if defense >= attack (so combat cannot deadlock at 0-dmg).
//
//	raw   = max(minDamageFloor, attack - defense * defenseMitigation)
//	varied = raw * (0.90 .. 1.10)   uniform variance
//	final  = varied * (crit ? 2 : 1)
func RollDamage(r RNG, attack, defense int32, crit bool) int32 {
	raw := float64(attack) - float64(defense)*defenseMitigation
	if raw < float64(MinDamageFloor) {
		raw = float64(MinDamageFloor)
	}
	// Uniform 0.90 .. 1.10 → cheap: sample a 0..20 int, add 90, /100.
	v := 90 + int(r.Uint64()%21)
	varied := raw * float64(v) / 100.0
	if crit {
		varied *= critMultiplier
	}
	if varied < float64(MinDamageFloor) {
		varied = float64(MinDamageFloor)
	}
	return int32(varied)
}

// --- Tuning constants -------------------------------------------------------
//
// Kept in code (not YAML) for Sprint 13 v1 — moving to a `combat.yaml` content
// pack is Phase 02 when the sim tools that drive balancing exist.

const (
	// MinHitPct floors effective to-hit — even a min-agi player still lands
	// occasional swings against a nimble mob.
	MinHitPct = 15
	// MaxHitPct caps effective to-hit — even a max-hit player misses sometimes.
	MaxHitPct = 95
	// MinDamageFloor guarantees positive damage per hit; without it a
	// defender with defense >> attacker.attack could deadlock the session.
	MinDamageFloor = 1
	// Defense converts to a fractional mitigation of raw attack. 0.5 means
	// each defense point cancels half an attack point. Combat balance knob.
	defenseMitigation = 0.5
	// Crit multiplies final damage. 2.0 is the classic Diablo-style default.
	critMultiplier = 2.0
)

// --- Errors -----------------------------------------------------------------

// ErrInvalidCombatant is returned when Resolve receives a Combatant with
// non-sensical stats (negative HP, zero attack on both sides, etc.). The
// intent is to fail early rather than compute a garbage log.
type CombatError struct {
	Reason string
}

func (e *CombatError) Error() string { return "combat: " + e.Reason }

func invalid(msg string, args ...any) error {
	return &CombatError{Reason: fmt.Sprintf(msg, args...)}
}
