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

	"github.com/otter/server/internal/leveling"
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

	// Combat extras (DD-250): flat percentage points, capped by the leveling
	// configuration. Zero unless a class formula, passive, talent or item
	// gives them.
	TargetEvasionPct      Target = "evasion_pct"
	TargetBlockPct        Target = "block_pct"
	TargetCritDmgPct      Target = "crit_dmg_pct"
	TargetExecutePct      Target = "execute_pct"
	TargetMissingHPDmgPct Target = "missing_hp_dmg_pct"
	TargetLifestealPct    Target = "lifesteal_pct"
	TargetDmgReductionPct Target = "dmg_reduction_pct"
)

// KnownKinds / KnownTargets are what ComputeSheet understands; content
// validation uses them so a typo can never become a silently ignored bonus.
var (
	KnownKinds = map[ModifierKind]bool{
		KindFlat: true, KindMultAdd: true, KindMultProd: true, KindClampMin: true, KindClampMax: true, KindSet: true,
	}
	KnownTargets = map[Target]bool{
		TargetSTR: true, TargetAGI: true, TargetINT: true, TargetVIT: true, TargetCHA: true, TargetLuck: true,
		TargetAttack: true, TargetDefense: true, TargetCritPct: true, TargetHitPct: true,
		TargetHPMax: true, TargetMPMax: true, TargetEnMax: true, TargetStaMax: true,
		TargetEvasionPct: true, TargetBlockPct: true, TargetCritDmgPct: true, TargetExecutePct: true,
		TargetMissingHPDmgPct: true, TargetLifestealPct: true, TargetDmgReductionPct: true,
	}
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

// Get returns a primary stat by its leveling attribute name.
func (p Primary) Get(name string) int {
	switch name {
	case "strength":
		return p.STR
	case "vitality":
		return p.VIT
	case "dexterity":
		return p.AGI
	case "intelligence":
		return p.INT
	case "wisdom":
		return p.CHA
	case "luck":
		return p.Luck
	}
	return 0
}

// Derived is the combat-derived stats.
type Derived struct {
	Attack  int `json:"attack"`
	Defense int `json:"defense"`
	CritPct int `json:"crit_pct"`
	HitPct  int `json:"hit_pct"`

	EvasionPct      int `json:"evasion_pct"`
	BlockPct        int `json:"block_pct"`
	CritDmgPct      int `json:"crit_dmg_pct"`
	ExecutePct      int `json:"execute_pct"`
	MissingHPDmgPct int `json:"missing_hp_dmg_pct"`
	LifestealPct    int `json:"lifesteal_pct"`
	DmgReductionPct int `json:"dmg_reduction_pct"`
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

// ComputeSheet is ComputeSheetFor with attack driven by STR.
func ComputeSheet(base Primary, mods []Modifier) Sheet { return ComputeSheetFor(base, TargetSTR, mods) }

// ComputeSheetFor applies the modifier stack to base primary stats; attack
// comes from the class's main stat (DD-229: TargetSTR, TargetAGI or TargetINT).
// Formulas for derived / resource maxes are placeholder v1.
func ComputeSheetFor(base Primary, main Target, mods []Modifier) Sheet {
	return ComputeSheetWith(base, main, mods, nil)
}

// Formula is a class's configuration-driven derived-stat formula at a level
// (DD-250); nil keeps the v1 formulas.
type Formula struct {
	leveling.Derived
	Level int
}

// FromLeveling converts progression modifiers.
func FromLeveling(list []leveling.Modifier) []Modifier {
	out := make([]Modifier, len(list))
	for i, m := range list {
		out[i] = Modifier{Kind: ModifierKind(m.Kind), Target: Target(m.Target), Amount: m.Amount}
	}
	return out
}

// ComputeSheetWith is ComputeSheetFor with an optional class formula. Raw
// attributes scale through the configured diminishing returns (DD-250 §11)
// — identical below the first bracket — and percent stats respect caps.
func ComputeSheetWith(base Primary, main Target, mods []Modifier, f *Formula) Sheet {
	cfg := leveling.Current()
	eff := cfg.EffectiveAttribute
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
	whole := func(v int) int { return int(eff(v)) }
	attackBase := float64(whole(primary.STR) * 2)
	switch main {
	case TargetAGI:
		attackBase = float64(whole(primary.AGI) * 2)
	case TargetINT:
		attackBase = float64(whole(primary.INT) * 2)
	}
	defenseBase := float64(whole(primary.VIT))
	critBase := float64(5 + whole(primary.Luck)/2)
	hitBase := float64(90 + whole(primary.AGI)/2)
	evasionBase := float64(min(whole(primary.AGI)/2, 40))
	blockBase := 0.0
	hpBase := 100.0 + float64(whole(primary.VIT))*10.0
	mpBase := 50.0 + float64(whole(primary.INT))*5.0
	if f != nil {
		attr := func(name string) float64 { return eff(primary.Get(name)) }
		attackBase = f.Attack.Eval(f.Level, attr)
		defenseBase = f.Armor.Eval(f.Level, attr)
		critBase = f.Crit.Eval(f.Level, attr)
		hitBase = f.Hit.Eval(f.Level, attr)
		evasionBase = f.Dodge.Eval(f.Level, attr)
		blockBase = f.Block.Eval(f.Level, attr)
		hpBase = f.HP.Eval(f.Level, attr)
		if f.MP != nil {
			mpBase = f.MP.Eval(f.Level, attr)
		}
	}
	pct := func(t Target, base float64) int {
		return intClamp(cfg.Cap(string(t), applyStack(base, filterTarget(mods, t))), 0, 100)
	}

	attack := applyStack(attackBase, filterTarget(mods, TargetAttack))
	defense := applyStack(defenseBase, filterTarget(mods, TargetDefense))
	crit := applyStack(critBase, filterTarget(mods, TargetCritPct))
	hit := applyStack(hitBase, filterTarget(mods, TargetHitPct))

	derived := Derived{
		Attack:  intClamp(attack, 0, 99999),
		Defense: intClamp(defense, 0, 99999),
		CritPct: intClamp(cfg.Cap(string(TargetCritPct), crit), 0, 100),
		HitPct:  intClamp(cfg.Cap(string(TargetHitPct), hit), 0, 100),

		EvasionPct:      pct(TargetEvasionPct, evasionBase),
		BlockPct:        pct(TargetBlockPct, blockBase),
		CritDmgPct:      intClamp(applyStack(0, filterTarget(mods, TargetCritDmgPct)), 0, 1000),
		ExecutePct:      intClamp(applyStack(0, filterTarget(mods, TargetExecutePct)), 0, 1000),
		MissingHPDmgPct: intClamp(applyStack(0, filterTarget(mods, TargetMissingHPDmgPct)), 0, 1000),
		LifestealPct:    pct(TargetLifestealPct, 0),
		DmgReductionPct: pct(TargetDmgReductionPct, 0),
	}

	// 3. Resources — HP from the class formula, MP from INT unless the
	// class formula has an MP term.
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
