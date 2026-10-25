package hunt

import (
	"math"

	"github.com/otter/server/internal/combat"
)

// Hunt pacing (DD-281): fights run back to back, a turn takes TurnSeconds,
// so combat power sets how fast mobs fall. After a win a hunter is wounded
// when its HP is under RestBelowPct of max or under three times what the
// fight cost (the next could kill). It drinks its potions first (DD-288);
// with none left it camps until full at RestHealPctPerSec of max a second,
// and that time counts against the hunt too. A carried pool (mana, energy)
// is spent across fights: it refills in combat as usual, BreatherMPPct of it
// between two fights, from mana potions under MPPotionBelowPct and during a
// camp; when dry the hunter fights on with basic attacks. A fight that leaves it alive under FleeBelowPct HP is
// broken off and ends the hunt.
const (
	TurnSeconds       = 3
	RestBelowPct      = 50
	RestHealPctPerSec = 1
	FleeBelowPct      = 5
	MPPotionBelowPct  = 30
	BreatherMPPct     = 10
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

// Supply kinds a hunt takes along (DD-288).
const (
	SupplyHP   = "hp"   // a healing potion
	SupplyMP   = "mp"   // a mana potion: refills the carried resource
	SupplyFood = "food" // a dish: heals and starts a timed attack/defense buff
)

// Supply is a stack of consumables taken on a hunt, its effect frozen at the
// start so a content change never rewrites a replay.
type Supply struct {
	Code     string `json:"code"`
	Kind     string `json:"kind"`
	Qty      int32  `json:"qty"`
	HealPct  int32  `json:"heal_pct,omitempty"`  // share of max HP
	HealFlat int32  `json:"heal_flat,omitempty"` // flat HP
	MPPct    int32  `json:"mp_pct,omitempty"`    // share of the pool
	BuffAtk  int32  `json:"buff_attack_pct,omitempty"`
	BuffDef  int32  `json:"buff_defense_pct,omitempty"`
	BuffSecs int32  `json:"buff_seconds,omitempty"`
}

// heal is the HP one unit restores for a hunter of maxHP.
func (s Supply) heal(maxHP int32) int32 {
	return s.HealFlat + int32((int64(maxHP)*int64(s.HealPct)+99)/100)
}

// Buff is a food buff already running when the hunt starts.
type Buff struct {
	AttackPct  int32 `json:"attack_pct"`
	DefensePct int32 `json:"defense_pct"`
	Secs       int32 `json:"seconds"` // left of it at the start
}

// Loadout is what a hunt brings along (DD-288).
type Loadout struct {
	Supplies []Supply `json:"supplies,omitempty"`
	Buff     *Buff    `json:"buff,omitempty"`
}

// SupplyUse is how many units of a supply the hunt used.
type SupplyUse struct {
	Code string `json:"code"`
	Qty  int32  `json:"qty"`
}

// HuntResult is the aggregate + per-encounter log the reveal handler grants.
type HuntResult struct {
	Completed       int   `json:"completed_encounters"`
	Victories       int   `json:"victories"`
	Defeats         int   `json:"defeats"`
	Timeouts        int   `json:"timeouts"`
	StoppedOnDeath  bool  `json:"stopped_on_death"` // true if a defeat cut the hunt short
	Fled            bool  `json:"fled"`             // broke off under FleeBelowPct HP (DD-288)
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
	// SuppliesUsed: units used per supply, in loadout order (DD-288).
	SuppliesUsed []SupplyUse `json:"supplies_used,omitempty"`
}

// Roll fights `mob` back to back for budgetSecs of hunting time (DD-281):
// each fight takes its turns × TurnSeconds, a win pays the mob's XPReward
// (the caller applies the level-gap falloff), and a wounded hunter camps to
// full first (see RestBelowPct). A defeat or a timeout ends the hunt; a
// fight that would run past the budget is not fought. Deterministic given
// (baseSeed, player, mob, budgetSecs): fight N uses seed baseSeed + N + 1.
//
// HP carries over between fights (plus the player's RegenPerSec while it
// fights), and so does a carried pool (DD-288). Gold and drop odds scale
// with the time each win took, so an hour pays about the same loot however
// fast the mob dies.
func Roll(
	baseSeed int64,
	player combat.Combatant,
	mob combat.Mob,
	budgetSecs int32,
	maxTurnsPerEncounter int,
	rngFactory func(seed int64) combat.RNG,
) *HuntResult {
	return RollWith(baseSeed, player, mob, budgetSecs, maxTurnsPerEncounter, rngFactory, Loadout{})
}

// RollWith is Roll with supplies (DD-288). Before the first fight the hunter
// tops up with the potions that would not overheal; after each win it eats
// when no food buff runs, drinks potions while wounded (strongest first) and
// a mana potion while its pool is under MPPotionBelowPct; only then camps.
// A running food buff raises attack and defense for the fights it covers.
func RollWith(
	baseSeed int64,
	player combat.Combatant,
	mob combat.Mob,
	budgetSecs int32,
	maxTurnsPerEncounter int,
	rngFactory func(seed int64) combat.RNG,
	load Loadout,
) *HuntResult {
	res := &HuntResult{Encounters: []EncounterRecord{}}
	k := newKit(player, load)
	k.prepare(res, 0, true)
	var gold, loot float64
	for i := 0; res.ElapsedSecs < budgetSecs; i++ {
		before := k.hp
		enc := player.Clone() // skills; MP full each fight unless the pool is carried
		enc.Name, enc.HPCurrent = "player", k.hp
		if k.carry {
			mp := k.mp
			enc.MPStart = &mp
		}
		if res.ElapsedSecs < k.buffUntil {
			enc.Attack = enc.Attack * (100 + k.buffAtk) / 100
			enc.Defense = enc.Defense * (100 + k.buffDef) / 100
		}
		sess := &combat.Session{Player: enc, Mob: mob.AsCombatant(), MaxTurns: maxTurnsPerEncounter,
			FleeBelowPct: FleeBelowPct}
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
		if k.carry {
			k.mp = min(player.MPMax, out.PlayerMPEnd+player.MPMax*BreatherMPPct/100)
		}
		rec := EncounterRecord{Index: i + 1, Outcome: out.Outcome, Duration: out.Duration,
			PlayerHPAft: out.PlayerHPEnd, MobHPAft: out.MobHPEnd}

		stop := false
		switch out.Outcome {
		case combat.OutcomeVictory:
			res.Victories++
			res.FightTurns += out.Duration
			k.hp = min(player.HPMax, out.PlayerHPEnd+int32(player.RegenPerSec*float64(fight)))
			k.lastCost = before - out.PlayerHPEnd
			wall := fight + k.prepare(res, res.ElapsedSecs, false)
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
			k.hp, stop = 0, true
		case combat.OutcomeFled:
			// Nearly dead: the hunter breaks off and goes home (DD-288).
			res.Fled = true
			k.hp, stop = out.PlayerHPEnd, true
		case combat.OutcomeTimeout:
			// Something is off with the pick (a mob nothing kills): stop
			// without a reward for it.
			res.Timeouts++
			k.hp, stop = out.PlayerHPEnd, true
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
	res.PlayerHPFinal = k.hp
	for i, s := range load.Supplies {
		if n := s.Qty - k.left[i]; n > 0 {
			res.SuppliesUsed = append(res.SuppliesUsed, SupplyUse{Code: s.Code, Qty: n})
		}
	}
	return res
}

// kit is a hunter's running state between fights: HP, the carried pool,
// the food buff and what is left of each supply.
type kit struct {
	player           combat.Combatant
	supplies         []Supply
	left             []int32
	hp, mp, lastCost int32
	carry            bool
	buffAtk, buffDef int32
	buffUntil        int32 // hunt second the food buff ends
}

func newKit(player combat.Combatant, load Loadout) *kit {
	k := &kit{player: player, supplies: load.Supplies, hp: player.HPCurrent,
		carry: player.ResourceCarry && player.MPMax > 0, mp: player.MPMax}
	k.left = make([]int32, len(load.Supplies))
	for i, s := range load.Supplies {
		k.left[i] = max(0, s.Qty)
	}
	if b := load.Buff; b != nil && b.Secs > 0 {
		k.buffAtk, k.buffDef, k.buffUntil = b.AttackPct, b.DefensePct, b.Secs
	}
	return k
}

// pick returns the index of the best supply of a kind still in stock
// (better: more HP healed / pool refilled / buff), or -1. fits, when set,
// limits the choice to units that would not overheal.
func (k *kit) pick(kind string, fits bool) int {
	best, bestV := -1, int32(-1)
	for i, s := range k.supplies {
		if s.Kind != kind || k.left[i] <= 0 {
			continue
		}
		v := s.heal(k.player.HPMax)
		switch kind {
		case SupplyMP:
			v = s.MPPct
		case SupplyFood:
			v = s.BuffAtk + s.BuffDef
		}
		if fits && v > k.player.HPMax-k.hp {
			continue
		}
		if v > bestV {
			best, bestV = i, v
		}
	}
	return best
}

func (k *kit) use(i int) {
	s := k.supplies[i]
	k.left[i]--
	k.hp = min(k.player.HPMax, k.hp+s.heal(k.player.HPMax))
	if s.MPPct > 0 {
		k.mp = min(k.player.MPMax, k.mp+(k.player.MPMax*s.MPPct+99)/100)
	}
}

func (k *kit) wounded() bool {
	return k.hp*100 < k.player.HPMax*RestBelowPct || (k.hp < 3*k.lastCost && k.hp < k.player.HPMax)
}

// prepare readies the hunter at hunt second now and returns the seconds it
// camped (added to the hunt). first: before the first fight, potions only
// top HP up where they would not overheal (or the hunter is wounded).
func (k *kit) prepare(res *HuntResult, now int32, first bool) int32 {
	if now >= k.buffUntil {
		if i := k.pick(SupplyFood, false); i >= 0 {
			s := k.supplies[i]
			k.use(i)
			k.buffAtk, k.buffDef, k.buffUntil = s.BuffAtk, s.BuffDef, now+s.BuffSecs
		}
	}
	for {
		i := -1
		if k.wounded() {
			i = k.pick(SupplyHP, false)
		} else if first {
			i = k.pick(SupplyHP, true)
		}
		if i < 0 {
			break
		}
		k.use(i)
	}
	for k.carry && k.mp*100 < k.player.MPMax*MPPotionBelowPct {
		i := k.pick(SupplyMP, false)
		if i < 0 {
			break
		}
		k.use(i)
	}
	if !k.wounded() {
		return 0
	}
	// Camp: HP and the carried pool refill at the same pace.
	missing := max(float64(k.player.HPMax-k.hp)/float64(k.player.HPMax), 0)
	if k.carry {
		missing = max(missing, float64(k.player.MPMax-k.mp)/float64(k.player.MPMax))
	}
	rest := int32(math.Ceil(missing * 100 / RestHealPctPerSec))
	res.RestSecs += rest
	res.ElapsedSecs += rest
	k.hp, k.mp = k.player.HPMax, k.player.MPMax
	return rest
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
