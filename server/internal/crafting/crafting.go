// Package crafting implements professions + recipes (Sprint 17).
//
// A profession is a per-character skill (blacksmithing, alchemy, cooking) that
// gates access to recipes. Learning is idempotent (no-op on re-learn); levelling
// happens automatically via XP grants at each successful craft.
//
// Crafting is server-authoritative and INSTANT (no AFK contract in v1):
//  1. Verify profession learned + skill >= recipe.min_skill.
//  2. Verify bag has all input stacks (ConsumeByTemplate refuses on shortage).
//  3. Debit gold via ledger (if recipe.gold_cost > 0) — vendor-sink pattern.
//  4. Consume inputs (atomic; rolls back on shortage).
//  5. Mint output via item.Service.MintExact (auto-slot).
//  6. Grant profession XP (level rollover with cap).
//  7. Append audit row to crafting.craft_log.
//
// Design decisions:
//   - Content lives in `content.profession` + `content.recipe`; catalog is
//     read once at boot and hot-reloadable via contentpack (Sprint 12 pattern).
//   - No skill loss on failure — v1 crafts never "fail". Failure states (broke
//     the anvil, ruined the potion) are Phase 02.
//   - Inputs must be stackable; non-stackable inputs (a specific magic sword
//     as a mat, say) would need reservation-by-instance-id, deferred to Phase 02.
package crafting

import (
	"errors"
	"math"
	"time"

	"github.com/google/uuid"
)

// Profession is a top-level meslek category. Loaded from content pack.
type Profession struct {
	Code           string
	DisplayNameKey string
	MaxLevel       int16
	Kind           string // "crafting" or "gathering" (DD-244)
	Version        int32
}

// Profession kinds. Gathering professions level from activities and are
// never learned by hand.
const (
	KindCrafting  = "crafting"
	KindGathering = "gathering"
)

// RecipeInput is one required material for a recipe.
type RecipeInput struct {
	TemplateCode string `json:"template_code"`
	Quantity     int32  `json:"quantity"`
}

// Recipe is a craftable — inputs + gold cost + output + profession req + xp.
type Recipe struct {
	Code               string
	DisplayNameKey     string
	ProfessionCode     string
	MinSkill           int16
	Inputs             []RecipeInput
	GoldCostCopper     int64
	OutputTemplateCode string
	OutputStack        int32
	ProfessionXPReward int32
	Version            int32
}

// CharacterProfession is one row in crafting.character_profession — the
// character's skill in a specific profession.
type CharacterProfession struct {
	CharacterID    uuid.UUID
	ProfessionCode string
	Level          int16
	XP             int32
	LearnedAt      time.Time
	UpdatedAt      time.Time
	Version        int32
}

// CraftResult is the reveal-shaped response for POST /craft — enough for the
// client to render "you made X, gained N xp, level → M".
type CraftResult struct {
	OutputInstanceID   uuid.UUID `json:"output_instance_id"`
	OutputTemplateCode string    `json:"output_template_code"`
	OutputStack        int32     `json:"output_stack"`
	GoldTxnID          uuid.UUID `json:"gold_txn_id,omitempty"`
	ProfessionCode     string    `json:"profession_code"`
	SkillBefore        int16     `json:"skill_before"`
	SkillAfter         int16     `json:"skill_after"`
	XPGained           int32     `json:"xp_gained"`
	LeveledUp          bool      `json:"leveled_up"`
}

// --- Errors -----------------------------------------------------------------

var (
	ErrNotFound             = errors.New("crafting: not found")
	ErrForbidden            = errors.New("crafting: forbidden")
	ErrUnknownProfession    = errors.New("crafting: unknown profession_code")
	ErrUnknownRecipe        = errors.New("crafting: unknown recipe_code")
	ErrProfessionNotLearned = errors.New("crafting: profession not learned")
	ErrSkillTooLow          = errors.New("crafting: skill level below min_skill")
	ErrInsufficientGold     = errors.New("crafting: insufficient gold")
	ErrValidation           = errors.New("crafting: validation")
)

// --- Tuning -----------------------------------------------------------------

// ProfessionXPToNext is the shared profession curve (DD-244): level N → N+1
// needs round(60 × 1.06^N) XP — 64 at 1, ~1 100 at 50, ~19 200 at 99.
// Kept in code (not YAML) since balance tuning belongs with the engine.
func ProfessionXPToNext(level int16) int32 {
	if level < 1 {
		level = 1
	}
	return int32(math.Round(60 * math.Pow(1.06, float64(level))))
}
