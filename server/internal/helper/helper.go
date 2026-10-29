// Package helper is hired NPC helpers (DD-319): a character buys helpers,
// one per gathering kind or crafting profession, and sends each on its own
// job beside the character's own work (helpers are not the character's one
// action). A helper works at HalfPct of the character — half the speed and
// half the XP, which goes to the character's profession — and each type
// adds its own bonus.
package helper

import (
	"errors"
	"math"
)

// HalfPct is a helper's speed and XP against the character's.
const HalfPct = 50

// Kinds of work.
const (
	KindGather = "gather"
	KindCraft  = "craft"
)

// Bonus kinds.
const (
	BonusYield    = "yield"     // gathering: more units
	BonusSpeed    = "speed"     // crafting: shorter jobs
	BonusXP       = "xp"        // more profession XP
	BonusGoldSave = "gold_save" // crafting: lower gold cost
	BonusDouble   = "double"    // crafting: a chance of a second output per unit
)

// Type is one kind of helper for sale.
type Type struct {
	Code        string `json:"code"`
	Kind        string `json:"kind"`       // gather | craft
	Profession  string `json:"profession"` // gathering kind or crafting profession
	Bonus       string `json:"bonus"`
	BonusPct    int    `json:"bonus_pct"`
	PriceCopper int64  `json:"price_copper"`
}

const gold = 10_000 // copper

// Types is the helper catalog: every gathering kind and crafting profession
// has one, each with a bonus of its own.
var Types = []Type{
	{"miner", KindGather, "mining", BonusYield, 15, 5 * gold},
	{"woodsman", KindGather, "woodcutting", BonusYield, 20, 5 * gold},
	{"angler", KindGather, "fishing", BonusXP, 25, 5 * gold},
	{"herbalist", KindGather, "herbalism", BonusXP, 20, 5 * gold},
	{"trapper", KindGather, "hunting", BonusYield, 10, 5 * gold},
	{"smith", KindCraft, "blacksmithing", BonusSpeed, 25, 10 * gold},
	{"tanner", KindCraft, "leatherworking", BonusXP, 25, 10 * gold},
	{"weaver", KindCraft, "tailoring", BonusGoldSave, 30, 10 * gold},
	{"carver", KindCraft, "woodworking", BonusSpeed, 20, 10 * gold},
	{"apothecary", KindCraft, "alchemy", BonusDouble, 10, 10 * gold},
	{"cook", KindCraft, "cooking", BonusDouble, 15, 10 * gold},
	{"smelter", KindCraft, "smelting", BonusSpeed, 30, 10 * gold},
	{"scribe", KindCraft, "enchanting", BonusXP, 30, 10 * gold},
	{"gemcutter", KindCraft, "jewelcrafting", BonusGoldSave, 25, 10 * gold},
}

// TypeOf finds a helper type by code.
func TypeOf(code string) (Type, bool) {
	for _, t := range Types {
		if t.Code == code {
			return t, true
		}
	}
	return Type{}, false
}

var (
	ErrUnknownType = errors.New("helper: unknown helper")
	ErrNotOwned    = errors.New("helper: not hired")
	ErrOwned       = errors.New("helper: already hired")
	ErrBusy        = errors.New("helper: already working")
	ErrIdle        = errors.New("helper: not working")
	ErrNotReady    = errors.New("helper: the job is not done")
	ErrValidation  = errors.New("helper: invalid job")
	ErrNoGold      = errors.New("helper: not enough gold")
	ErrSkill       = errors.New("helper: the character's level is too low for that")
	ErrMaterials   = errors.New("helper: not enough materials")
)

// pct applies a bonus percent to v when the type has that bonus.
func (t Type) pct(bonus string, v float64) float64 {
	if t.Bonus != bonus {
		return v
	}
	if bonus == BonusSpeed || bonus == BonusGoldSave {
		return v * float64(100-t.BonusPct) / 100 // less time / gold
	}
	return v * float64(100+t.BonusPct) / 100
}

// YieldMult is the multiplier on the character's gathering rate (activity
// snapshot YieldMult adds to 1): half the rate, then the yield bonus.
func (t Type) YieldMult() float64 { return t.pct(BonusYield, HalfPct/100.0) - 1 }

// XP is the profession XP a helper's work earns the character.
func (t Type) XP(full int) int { return int(math.Round(t.pct(BonusXP, float64(full)*HalfPct/100))) }

// CraftSeconds is a helper's time for count units of a recipe the
// character makes in unitSecs each: twice the time, less the speed bonus.
func (t Type) CraftSeconds(unitSecs int32, count int) int32 {
	return int32(math.Ceil(t.pct(BonusSpeed, float64(unitSecs)*float64(count)*100/HalfPct)))
}

// GoldCost is the gold a helper's craft costs.
func (t Type) GoldCost(perUnit int64, count int) int64 {
	return int64(math.Round(t.pct(BonusGoldSave, float64(perUnit*int64(count)))))
}

// Doubles is how many of count units come out twice, by a roll in 0..99
// per unit.
func (t Type) Doubles(count int, roll func(i int) int) int {
	if t.Bonus != BonusDouble {
		return 0
	}
	n := 0
	for i := 0; i < count; i++ {
		if roll(i) < t.BonusPct {
			n++
		}
	}
	return n
}
