// Package pet is the pet system (DD-328): eggs that cleared raids past
// level 400 drop hatch into pets, which level up on hunts of their own and
// add a bonus to their owner while active. Content (content/pets.yaml)
// installs the species catalog on every load, like the item sets; readers
// take Current().
package pet

import (
	"fmt"
	"math"
	"sort"
	"sync/atomic"

	"github.com/otter/server/internal/combat"
	"github.com/otter/server/internal/hunt"
)

// MaxLevel is a pet's level cap.
const MaxLevel = 100

// PowerPerLevel: a pet of level L fights like the mob curve at L × this.
const PowerPerLevel = 10

// MobReach is how far above its power a pet may hunt: the combat level-gap
// exemption, past which its attack falls off.
const MobReach = 5

// Modifier mirrors skill.Modifier (kind, target, amount) without the import.
type Modifier struct {
	Kind   string  `yaml:"kind"   json:"kind"`
	Target string  `yaml:"target" json:"target"`
	Amount float64 `yaml:"amount" json:"amount"`
}

// Species is one kind of pet and its egg.
type Species struct {
	Code           string   `yaml:"code"             json:"code"`
	DisplayNameKey string   `yaml:"display_name_key" json:"display_name_key"`
	Rank           int      `yaml:"rank"             json:"rank"`
	MinMobLevel    int      `yaml:"min_mob_level"    json:"min_mob_level"`
	EggPct         float64  `yaml:"egg_pct"          json:"egg_pct"`
	HatchSeconds   int      `yaml:"hatch_seconds"    json:"hatch_seconds"`
	PowerPct       int      `yaml:"power_pct"        json:"power_pct"`
	Bonus          Modifier `yaml:"bonus"            json:"bonus"` // at MaxLevel
}

// EggCode is the item template of the species' egg.
func (s Species) EggCode() string { return EggPrefix + s.Code }

// EggPrefix starts every egg's item code.
const EggPrefix = "pet_egg_"

// Catalog is the species list, best first.
type Catalog struct {
	species []Species
	byCode  map[string]int
	byEgg   map[string]int
}

var current atomic.Pointer[Catalog]

// Current is the installed catalog; empty before content loads.
func Current() *Catalog {
	if c := current.Load(); c != nil {
		return c
	}
	return &Catalog{}
}

// New validates species: unique codes and ranks, odds in (0, 100], a
// hatch time, a power and a bonus.
func New(species []Species) (*Catalog, error) {
	c := &Catalog{species: append([]Species(nil), species...), byCode: map[string]int{}, byEgg: map[string]int{}}
	sort.Slice(c.species, func(i, j int) bool { return c.species[i].Rank > c.species[j].Rank })
	ranks := map[int]bool{}
	for i, s := range c.species {
		if s.Code == "" || ranks[s.Rank] {
			return nil, fmt.Errorf("pet: empty code or duplicate rank %q/%d", s.Code, s.Rank)
		}
		if _, dup := c.byCode[s.Code]; dup {
			return nil, fmt.Errorf("pet: duplicate code %q", s.Code)
		}
		if s.EggPct <= 0 || s.EggPct > 100 || s.HatchSeconds <= 0 || s.PowerPct <= 0 || s.Bonus.Amount <= 0 {
			return nil, fmt.Errorf("pet %s: odds, hatch time, power and bonus must be positive", s.Code)
		}
		ranks[s.Rank] = true
		c.byCode[s.Code] = i
		c.byEgg[s.EggCode()] = i
	}
	return c, nil
}

// Install validates and swaps in the catalog content loaded.
func Install(species []Species) error {
	c, err := New(species)
	if err != nil {
		return err
	}
	current.Store(c)
	return nil
}

// All is every species, best first.
func (c *Catalog) All() []Species { return c.species }

// Get is the species of a code.
func (c *Catalog) Get(code string) (Species, bool) {
	i, ok := c.byCode[code]
	if !ok {
		return Species{}, false
	}
	return c.species[i], true
}

// ByEgg is the species an egg item hatches into.
func (c *Catalog) ByEgg(template string) (Species, bool) {
	i, ok := c.byEgg[template]
	if !ok {
		return Species{}, false
	}
	return c.species[i], true
}

// RollEgg is the egg (item code) a member gets for a cleared raid whose
// boss is mobLevel, or "": every species the level reaches is rolled at
// its odds, best first, and the first hit drops — so the best egg comes at
// exactly its EggPct and at most one egg drops.
func (c *Catalog) RollEgg(mobLevel int, r combat.RNG) string {
	for _, s := range c.species {
		if s.MinMobLevel > mobLevel {
			continue
		}
		if float64(r.Uint64()>>11)/float64(1<<53)*100 < s.EggPct {
			return s.EggCode()
		}
	}
	return ""
}

// XPToNext is the XP a pet needs from level to level+1 (0 at the cap):
// about an hour of hunting even-level mobs at level 1, growing 2.5% a level.
func XPToNext(level int) int64 {
	if level >= MaxLevel {
		return 0
	}
	return int64(math.Round(XPPerHour * float64(level) * math.Pow(1.025, float64(level-1))))
}

// XPPerHour scales XPToNext: the XP per pet level a pet earns in an hour on
// mobs of its power (TestPetPace keeps it honest).
const XPPerHour = 1000

// Power is the mob level a pet of the level fights like.
func Power(level int) int { return level * PowerPerLevel }

// Combatant is a pet of the species and level as it fights: the content's
// mob curve (hp 150+25P, attack 28+2.1P, defense 1+0.8P, the gen_world.py
// formulas) at its power P, times its species' power.
func (s Species) Combatant(level int) combat.Combatant {
	p := float64(Power(level))
	f := float64(s.PowerPct) / 100
	hp := int32(math.Round((150 + 25*p) * f))
	return combat.Combatant{
		Name:       s.Code,
		HPMax:      hp,
		HPCurrent:  hp,
		Attack:     int32(math.Round((28 + 2.1*p) * f)),
		Defense:    int32(math.Round((1 + 0.8*p) * f)),
		CritPct:    5 + int32(s.Rank)/4,
		HitPct:     min(92, 82+int32(p)/5),
		EvasionPct: 5 + combat.LevelEvasion(int16(p)),
		Level:      int16(p),
	}
}

// BonusAt is the species bonus at a level: a quarter of it at level 1, all
// of it at MaxLevel.
func (s Species) BonusAt(level int) Modifier {
	level = max(1, min(level, MaxLevel))
	m := s.Bonus
	m.Amount = math.Round(m.Amount*(0.25+0.75*float64(level-1)/float64(MaxLevel-1))*10000) / 10000
	return m
}

// CanHunt says whether a pet of the level may hunt the mob.
func CanHunt(level int, mob combat.Mob) bool {
	return !mob.RaidBoss && int(mob.Level) <= Power(level)+MobReach
}

// AddXP adds xp to a pet at level/xp, rolling levels over; at the cap the
// XP stops counting.
func AddXP(level int, xp, gain int64) (int, int64) {
	xp += gain
	for level < MaxLevel {
		need := XPToNext(level)
		if xp < need {
			return level, xp
		}
		xp -= need
		level++
	}
	return MaxLevel, 0
}

// Roll is a pet's hunt (the character hunt's pacing, DD-281): mob fought
// back to back for budgetSecs, camping when wounded; each win is worth the
// mob's level in pet XP (Result.TotalXP). Deterministic given its inputs.
func Roll(seed int64, s Species, level int, mob combat.Mob, budgetSecs int32,
	rngNew func(int64) combat.RNG) *hunt.HuntResult {
	mob.XPReward = int32(mob.Level)
	return hunt.Roll(seed, s.Combatant(level), mob, budgetSecs, combat.DefaultMaxTurns, rngNew)
}
