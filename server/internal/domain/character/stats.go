package character

import "time"

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

// basePrimaryForClass — v1 hardcoded per class. Content-driven (content.class)
// replaces this in Faz 01 sprint 8. Numbers are placeholder for balance.
func basePrimaryForClass(classCode string) PrimaryStats {
	switch classCode {
	case "warrior":
		return PrimaryStats{Strength: 15, Agility: 10, Intelligence: 8, Vitality: 15, Charisma: 10, Luck: 5}
	case "mage":
		return PrimaryStats{Strength: 8, Agility: 10, Intelligence: 15, Vitality: 10, Charisma: 12, Luck: 5}
	case "hunter":
		return PrimaryStats{Strength: 12, Agility: 15, Intelligence: 8, Vitality: 10, Charisma: 10, Luck: 5}
	case "priest":
		return PrimaryStats{Strength: 8, Agility: 8, Intelligence: 12, Vitality: 12, Charisma: 15, Luck: 5}
	case "rogue":
		return PrimaryStats{Strength: 10, Agility: 15, Intelligence: 10, Vitality: 10, Charisma: 8, Luck: 12}
	default:
		return PrimaryStats{Strength: 10, Agility: 10, Intelligence: 10, Vitality: 10, Charisma: 10, Luck: 10}
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

// deriveStats computes attack/defense/crit/hit from PrimaryStats.
func deriveStats(p PrimaryStats) DerivedStats {
	return DerivedStats{
		Attack:  p.Strength * 2,
		Defense: p.Vitality * 2,
		CritPct: 5 + p.Luck/2,
		HitPct:  90 + p.Agility/2,
	}
}

// xpToNext returns the XP required to hit level+1. Placeholder linear curve
// (`level × 100`) — real curve is content.xp_curve in Faz 01 sprint 6+.
func xpToNext(level int) int64 {
	if level < 1 {
		level = 1
	}
	return int64(level) * 100
}

// buildSheet renders a Sheet from the stored Character. Stateless — no repo
// lookups; caller passes the resolved character in.
func buildSheet(c *Character) Sheet {
	p := basePrimaryForClass(c.ClassCode)
	return Sheet{
		ID:        c.ID.String(),
		Name:      c.Name,
		Race:      c.RaceCode,
		Class:     c.ClassCode,
		Spec:      c.SpecCode,
		Level:     int(c.Level),
		XP:        c.XP,
		XPToNext:  xpToNext(int(c.Level)),
		Primary:   p,
		Derived:   deriveStats(p),
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
