package character

import (
	"math"
	"time"
)

// PrimaryStats are the six per-character attributes. Full modifier stack (flat
// → mult_add → mult_prod → clamp) lands with the skill-engine sprint (Faz 01
// sprint 9). For now these come purely from the class base.
type PrimaryStats struct {
	Strength     int `json:"str"`
	Agility      int `json:"agi"`
	Intelligence int `json:"int"`
	Vitality     int `json:"vit"`
	Charisma     int `json:"cha"`
	Luck         int `json:"luck"`
}

// Resources are the runtime pools. Current == Max at spawn / after full rest.
type Resources struct {
	HPCurrent      int `json:"hp_current"`
	HPMax          int `json:"hp_max"`
	MPCurrent      int `json:"mp_current"`
	MPMax          int `json:"mp_max"`
	EnergyCurrent  int `json:"energy_current"`
	EnergyMax      int `json:"energy_max"`
	StaminaCurrent int `json:"stamina_current"`
	StaminaMax     int `json:"stamina_max"`
}

// DerivedStats are computed from PrimaryStats. Values are rounded down.
type DerivedStats struct {
	Attack  int `json:"attack"`
	Defense int `json:"defense"`
	CritPct int `json:"crit_pct"` // percentage points (5 = 5%)
	HitPct  int `json:"hit_pct"`
}

// Sheet is the compound "read one screen" response for the character sheet UI.
// XPToNext is included so the client can render the level bar without a second
// call.
type Sheet struct {
	ID        string       `json:"id"`
	Name      string       `json:"name"`
	Race      string       `json:"race_code"`
	Class     string       `json:"class_code"`
	Spec      *string      `json:"spec_code,omitempty"`
	Level     int          `json:"level"`
	XP        int64        `json:"xp"`
	XPToNext  int64        `json:"xp_to_next"`
	Primary   PrimaryStats `json:"primary"`
	Derived   DerivedStats `json:"derived"`
	Resources Resources    `json:"resources"`
	State     string       `json:"state"`
}

// classBase is each class's level-1 primary block; growthPerLevel is what
// every level after the first adds, in half points (3 = +1.5/level; DD-229). Placeholder numbers tuned with
// the balance simulation (internal/balance).
var classBase = map[string]PrimaryStats{
	"warrior": {Strength: 15, Agility: 10, Intelligence: 8, Vitality: 15, Charisma: 10, Luck: 5},
	"mage":    {Strength: 8, Agility: 10, Intelligence: 15, Vitality: 10, Charisma: 12, Luck: 5},
	"hunter":  {Strength: 12, Agility: 15, Intelligence: 8, Vitality: 10, Charisma: 10, Luck: 5},
	"priest":  {Strength: 8, Agility: 8, Intelligence: 12, Vitality: 12, Charisma: 15, Luck: 5},
	"rogue":   {Strength: 10, Agility: 15, Intelligence: 10, Vitality: 10, Charisma: 8, Luck: 12},
}

var growthPerLevel = map[string]PrimaryStats{
	"warrior": {Strength: 4, Agility: 1, Vitality: 2},
	"mage":    {Intelligence: 5, Agility: 1, Vitality: 1, Charisma: 1},
	"hunter":  {Agility: 4, Strength: 1, Vitality: 2, Luck: 1},
	"priest":  {Intelligence: 4, Vitality: 2, Charisma: 2},
	"rogue":   {Agility: 4, Vitality: 1, Luck: 2},
}

// BasePrimary is a class's primary stats at a level, before gear and skills.
func BasePrimary(class string, level int) PrimaryStats {
	b, ok := classBase[class]
	if !ok {
		b = PrimaryStats{Strength: 10, Agility: 10, Intelligence: 10, Vitality: 10, Charisma: 10, Luck: 10}
	}
	g, n := growthPerLevel[class], max(level, 1)-1 // half points: g*n/2
	return PrimaryStats{
		Strength:     b.Strength + g.Strength*n/2,
		Agility:      b.Agility + g.Agility*n/2,
		Intelligence: b.Intelligence + g.Intelligence*n/2,
		Vitality:     b.Vitality + g.Vitality*n/2,
		Charisma:     b.Charisma + g.Charisma*n/2,
		Luck:         b.Luck + g.Luck*n/2,
	}
}

// MainStat names the primary stat behind a class's attack (DD-229):
// "str", "agi" or "int".
func MainStat(class string) string {
	switch class {
	case "mage", "priest":
		return "int"
	case "hunter", "rogue":
		return "agi"
	default:
		return "str"
	}
}

func mainStatValue(class string, p PrimaryStats) int {
	switch MainStat(class) {
	case "agi":
		return p.Agility
	case "int":
		return p.Intelligence
	default:
		return p.Strength
	}
}

// deriveResources computes HP/MP/Energy/Stamina at full from PrimaryStats.
// Placeholder formulas — replaced by skill-engine modifier stack.
func deriveResources(p PrimaryStats) Resources {
	hp := 100 + p.Vitality*10
	mp := 50 + p.Intelligence*5
	return Resources{
		HPCurrent: hp, HPMax: hp,
		MPCurrent: mp, MPMax: mp,
		EnergyCurrent: 100, EnergyMax: 100,
		StaminaCurrent: 100, StaminaMax: 100,
	}
}

// deriveStats computes attack/defense/crit/hit from PrimaryStats (mirrors
// skill.ComputeSheetFor without modifiers).
func deriveStats(class string, p PrimaryStats) DerivedStats {
	return DerivedStats{
		Attack:  mainStatValue(class, p) * 2,
		Defense: p.Vitality,
		CritPct: 5 + p.Luck/2,
		HitPct:  90 + p.Agility/2,
	}
}

// XPToNext returns the XP required to reach level+1.
func XPToNext(level int) int64 {
	if level < 1 {
		level = 1
	}
	return int64(math.Round(50 * math.Pow(float64(level), 2.2)))
}

// buildSheet renders a Sheet from the stored Character. Stateless — no repo
// lookups; caller passes the resolved character in.
func buildSheet(c *Character) Sheet {
	p := BasePrimary(c.ClassCode, int(c.Level))
	return Sheet{
		ID:        c.ID.String(),
		Name:      c.Name,
		Race:      c.RaceCode,
		Class:     c.ClassCode,
		Spec:      c.SpecCode,
		Level:     int(c.Level),
		XP:        c.XP,
		XPToNext:  XPToNext(int(c.Level)),
		Primary:   p,
		Derived:   deriveStats(c.ClassCode, p),
		Resources: deriveResources(p),
		State:     string(c.State),
	}
}

// HPRegenPer2s is the out-of-combat regen per 2-second step
// (docs/character/health.md §3.2, VIT standing in for STA).
func HPRegenPer2s(vitality int) float64 { return 1 + 0.1*float64(vitality) }

// EffectiveHP resolves a character's current HP at `now` without any ticking:
// the stored value plus analytical regen since it was written, capped at max.
// A nil stored value means the character has never been damaged (full HP).
func EffectiveHP(stored *int32, updatedAt *time.Time, maxHP, vitality int, now time.Time) int {
	if stored == nil || updatedAt == nil {
		return maxHP
	}
	hp := int(*stored)
	if elapsed := now.Sub(*updatedAt); elapsed > 0 {
		steps := int64(elapsed / (2 * time.Second))
		hp += int(float64(steps) * HPRegenPer2s(vitality))
	}
	return max(0, min(hp, maxHP))
}
