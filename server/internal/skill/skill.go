// Package skill implements the modifier stack (DD-107 pattern) — pure math for
// stat computation. Consumers pass base primary stats + a slice of Modifier
// (from equipped gear, buffs, class specs, mastery, ascension nodes), receive
// a computed Sheet.
//
// Apply order per target:
//
//  1. base + Σ(flat)                     — additive raw
//  2. × (1 + Σ(mult_add / 100))         — additive percentage
//  3. × Π(mult_prod values / 100 + 1)   — multiplicative percentage (rare)
//  4. clamp_min / clamp_max              — enforce bounds
//  5. set                                — hard override (last-write-wins)
//
// Primary stats compute first, then derived stats compute from FINAL primary
// (base derived formula) + their own modifier stack applies.
package skill

import (
	"encoding/json"
	"math"
)

// ModifierKind identifies how the modifier combines into the stack.
type ModifierKind string

const (
	KindFlat     ModifierKind = "flat"      // add raw value
	KindMultAdd  ModifierKind = "mult_add"  // additive percentage (0.10 = +10%)
	KindMultProd ModifierKind = "mult_prod" // multiplicative percentage
	KindClampMin ModifierKind = "clamp_min"
	KindClampMax ModifierKind = "clamp_max"
	KindSet      ModifierKind = "set" // hard override
)

// Target names the stat the modifier acts on.
type Target string

const (
	TargetSTR     Target = "str"
	TargetAGI     Target = "agi"
	TargetINT     Target = "int"
	TargetVIT     Target = "vit"
	TargetCHA     Target = "cha"
	TargetLuck    Target = "luck"
	TargetAttack  Target = "attack"
	TargetDefense Target = "defense"
	TargetCritPct Target = "crit_pct"
	TargetHitPct  Target = "hit_pct"
	TargetHPMax   Target = "hp_max"
	TargetMPMax   Target = "mp_max"
	TargetEnMax   Target = "energy_max"
	TargetStaMax  Target = "stamina_max"
)

// Modifier is one bonus/malus. Amount interpretation depends on Kind.
type Modifier struct {
	Kind   ModifierKind `json:"kind"`
	Target Target       `json:"target"`
	Amount float64      `json:"amount"`
}

// DecodeModifiers parses a JSONB stat_bonuses column.
func DecodeModifiers(raw json.RawMessage) ([]Modifier, error) {
	if len(raw) == 0 {
		return nil, nil
	}
	var out []Modifier
	if err := json.Unmarshal(raw, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// Primary is the six base attributes.
type Primary struct {
	STR  int `json:"str"`
	AGI  int `json:"agi"`
	INT  int `json:"int"`
	VIT  int `json:"vit"`
	CHA  int `json:"cha"`
	Luck int `json:"luck"`
}

// Derived is the four combat-derived stats.
type Derived struct {
	Attack  int `json:"attack"`
	Defense int `json:"defense"`
	CritPct int `json:"crit_pct"`
	HitPct  int `json:"hit_pct"`
}

// Resources are the runtime pools (max only — current is game-state).
type Resources struct {
	HPMax  int `json:"hp_max"`
	MPMax  int `json:"mp_max"`
	EnMax  int `json:"energy_max"`
	StaMax int `json:"stamina_max"`
}

// Sheet is the computed output — primary + derived + resource max caps.
type Sheet struct {
	Primary   Primary   `json:"primary"`
	Derived   Derived   `json:"derived"`
	Resources Resources `json:"resources"`
}

// ComputeSheet applies the modifier stack to base primary stats.
// Formulas for derived / resource maxes are placeholder v1 (Faz 00 spec —
// balance pass Phase 02). Content-driven formulas later.
func ComputeSheet(base Primary, mods []Modifier) Sheet {
	// 1. Primary stats.
	str := applyStack(float64(base.STR), filterTarget(mods, TargetSTR))
	agi := applyStack(float64(base.AGI), filterTarget(mods, TargetAGI))
	intel := applyStack(float64(base.INT), filterTarget(mods, TargetINT))
	vit := applyStack(float64(base.VIT), filterTarget(mods, TargetVIT))
	cha := applyStack(float64(base.CHA), filterTarget(mods, TargetCHA))
	luck := applyStack(float64(base.Luck), filterTarget(mods, TargetLuck))

	primary := Primary{
		STR:  intClamp(str, 1, 9999),
		AGI:  intClamp(agi, 1, 9999),
		INT:  intClamp(intel, 1, 9999),
		VIT:  intClamp(vit, 1, 9999),
		CHA:  intClamp(cha, 1, 9999),
		Luck: intClamp(luck, 1, 9999),
	}

	// 2. Derived — compute from FINAL primary, then apply their own mods.
	// Base formulas use integer division to give predictable half-values
	// (e.g. Luck 5 → +2 crit, not +2.5 that rounds to +3).
	attackBase := float64(primary.STR * 2)
	defenseBase := float64(primary.VIT * 2)
	critBase := float64(5 + primary.Luck/2)
	hitBase := float64(90 + primary.AGI/2)

	attack := applyStack(attackBase, filterTarget(mods, TargetAttack))
	defense := applyStack(defenseBase, filterTarget(mods, TargetDefense))
	crit := applyStack(critBase, filterTarget(mods, TargetCritPct))
	hit := applyStack(hitBase, filterTarget(mods, TargetHitPct))

	derived := Derived{
		Attack:  intClamp(attack, 0, 99999),
		Defense: intClamp(defense, 0, 99999),
		CritPct: intClamp(crit, 0, 100),
		HitPct:  intClamp(hit, 0, 100),
	}

	// 3. Resources — HP scales with VIT, MP with INT (placeholder).
	hpBase := 100.0 + float64(primary.VIT)*10.0
	mpBase := 50.0 + float64(primary.INT)*5.0
	enBase := 100.0
	staBase := 100.0

	hpMax := applyStack(hpBase, filterTarget(mods, TargetHPMax))
	mpMax := applyStack(mpBase, filterTarget(mods, TargetMPMax))
	enMax := applyStack(enBase, filterTarget(mods, TargetEnMax))
	staMax := applyStack(staBase, filterTarget(mods, TargetStaMax))

	resources := Resources{
		HPMax:  intClamp(hpMax, 1, 999999),
		MPMax:  intClamp(mpMax, 0, 999999),
		EnMax:  intClamp(enMax, 1, 100),
		StaMax: intClamp(staMax, 1, 100),
	}

	return Sheet{Primary: primary, Derived: derived, Resources: resources}
}

// applyStack computes the final value for one target by walking a filtered
// modifier list in canonical order.
func applyStack(base float64, mods []Modifier) float64 {
	// flat
	v := base
	for _, m := range mods {
		if m.Kind == KindFlat {
			v += m.Amount
		}
	}
	// mult_add: sum then apply once
	multAddSum := 0.0
	for _, m := range mods {
		if m.Kind == KindMultAdd {
			multAddSum += m.Amount
		}
	}
	v = v * (1.0 + multAddSum)
	// mult_prod: compound
	for _, m := range mods {
		if m.Kind == KindMultProd {
			v = v * (1.0 + m.Amount)
		}
	}
	// clamp_min / clamp_max
	for _, m := range mods {
		if m.Kind == KindClampMin && v < m.Amount {
			v = m.Amount
		}
	}
	for _, m := range mods {
		if m.Kind == KindClampMax && v > m.Amount {
			v = m.Amount
		}
	}
	// set — last-write-wins override
	for _, m := range mods {
		if m.Kind == KindSet {
			v = m.Amount
		}
	}
	return v
}

func filterTarget(mods []Modifier, t Target) []Modifier {
	out := make([]Modifier, 0, len(mods))
	for _, m := range mods {
		if m.Target == t {
			out = append(out, m)
		}
	}
	return out
}

func intClamp(v float64, min, max int) int {
	r := int(math.Round(v))
	if r < min {
		return min
	}
	if r > max {
		return max
	}
	return r
}
