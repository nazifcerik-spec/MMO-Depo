// Package crafting implements professions + recipes (Sprint 17).
//
// A profession is a per-character skill (blacksmithing, alchemy, cooking) that
// gates access to recipes. Learning is idempotent (no-op on re-learn); levelling
// happens automatically via XP grants at each successful craft.
//
// Crafting is server-authoritative and TIMED (DD-264): a job of count units,
// each taking the recipe's craft_seconds.
//
//	Start:   skill gate → materials + gold for the whole job checked → job row
//	         claimed (one action at a time, DD-263) → gold debited → inputs
//	         consumed (a failure undoes the job and refunds the gold).
//	Collect: any time. The units finished so far are minted and pay their
//	         profession XP; the unfinished rest is refunded (gold + inputs).
//	         Idempotent: the job row flips to collected once and stores the
//	         result.
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

	"github.com/otter/server/internal/partyxp"
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
	CraftSeconds       int32 // per unit (DD-264)
	Version            int32
}

// Recipe difficulty (DD-294): how far the character's skill is past a
// recipe's minimum decides the color shown and the share of its profession
// XP a craft pays. Grey recipes teach nothing more.
const (
	DifficultyOrange = "orange"
	DifficultyYellow = "yellow"
	DifficultyGreen  = "green"
	DifficultyGrey   = "grey"
)

// DifficultyBand: a skill gap below Below is Color and pays XPFactor.
type DifficultyBand struct {
	Below    int16   `json:"below"`
	Color    string  `json:"color"`
	XPFactor float64 `json:"xp_factor"`
}

// DifficultyBands, closest first; a gap past the last is grey (no XP). The
// recipe ladders keep a non-grey recipe within reach at every skill
// (TestRecipeLaddersNeverGoGrey).
var DifficultyBands = []DifficultyBand{
	{Below: 5, Color: DifficultyOrange, XPFactor: 1},
	{Below: 10, Color: DifficultyYellow, XPFactor: 0.6},
	{Below: 20, Color: DifficultyGreen, XPFactor: 0.25},
}

// Difficulty is a recipe's color and XP share at a skill.
func Difficulty(minSkill, skill int16) (string, float64) {
	gap := skill - minSkill
	for _, b := range DifficultyBands {
		if gap < b.Below {
			return b.Color, b.XPFactor
		}
	}
	return DifficultyGrey, 0
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
	Count              int       `json:"count"`    // units the job was started with (DD-245)
	Made               int       `json:"made"`     // units delivered (DD-264)
	Refunded           int       `json:"refunded"` // unfinished units whose gold + inputs came back
	// Party is the party bonus on the XP (DD-305); nil = none.
	Party    *partyxp.Outcome `json:"party_share,omitempty"`
	Replayed bool             `json:"replayed,omitempty"`
}

// Job states.
const (
	JobActive    = "active"
	JobCollected = "collected"
)

// Job is a timed crafting run (DD-264). Recipe is frozen at start.
type Job struct {
	ID          uuid.UUID
	CharacterID uuid.UUID
	Recipe      Recipe
	Count       int
	UnitSeconds int32
	StartedAt   time.Time
	ReadyAt     time.Time
	State       string
	Result      []byte
}

// Done is how many units are finished at now.
func (j *Job) Done(now time.Time) int {
	if j.UnitSeconds <= 0 {
		return j.Count
	}
	n := int(now.Sub(j.StartedAt) / (time.Duration(j.UnitSeconds) * time.Second))
	return max(0, min(j.Count, n))
}

// MaxJobSeconds caps one job, like the AFK activities (6 h).
const MaxJobSeconds = 6 * 60 * 60

// DefaultCraftSeconds applies to a recipe without a time.
const DefaultCraftSeconds = 60

// MaxCraftCount caps one batch request (DD-245).
const MaxCraftCount = 100

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
	ErrJobActive            = errors.New("crafting: a crafting job is already running")
	ErrJobTooLong           = errors.New("crafting: job longer than MaxJobSeconds")
)

// --- Tuning -----------------------------------------------------------------

// ProfessionXPToNext is the shared profession curve (DD-244), 20% steeper
// since DD-309: level N → N+1 needs round(72 × 1.06^N) XP — 76 at 1,
// ~1 330 at 50, ~23 000 at 99.
// Kept in code (not YAML) since balance tuning belongs with the engine.
func ProfessionXPToNext(level int16) int32 {
	if level < 1 {
		level = 1
	}
	return int32(math.Round(72 * math.Pow(1.06, float64(level))))
}
