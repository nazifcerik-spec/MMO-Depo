package character

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"

	"github.com/google/uuid"

	"github.com/otter/server/internal/leveling"
)

// Progression errors (DD-250).
var (
	ErrInvalidPromotion  = errors.New("character: promotion not available")
	ErrNotEnoughPoints   = errors.New("character: not enough points")
	ErrUnknownAttribute  = errors.New("character: unknown attribute")
	ErrUnknownTalent     = errors.New("character: unknown talent")
	ErrTalentMaxed       = errors.New("character: talent at max rank")
	ErrCapstoneLocked    = errors.New("character: capstone needs level and a full tree")
	ErrBusy              = errors.New("character: busy (combat, AFK or dungeon)")
	ErrInvalidSpec       = errors.New("character: specialization not allowed")
	ErrProgressionUnsupp = errors.New("character: progression store unavailable")
)

// Progression event kinds (character.progression_event.kind).
const (
	EventXPGained          = "XP_GAINED"
	EventLevelUp           = "LEVEL_UP"
	EventPromotionAvail    = "PROMOTION_AVAILABLE"
	EventPromotionSelected = "PROMOTION_SELECTED"
	EventSkillUnlocked     = "SKILL_UNLOCKED"
	EventPassiveUnlocked   = "PASSIVE_UNLOCKED"
	EventAttributePoints   = "ATTRIBUTE_POINTS_GAINED"
	EventTalentPoint       = "TALENT_POINT_GAINED"
	EventTitleUnlocked     = "TITLE_UNLOCKED"
	EventAwakening         = "AWAKENING_UNLOCKED"
	EventUltimate          = "ULTIMATE_UNLOCKED"
	EventCapstone          = "CAPSTONE_UNLOCKED"
	EventMastery           = "CLASS_MASTERY_UNLOCKED"
	EventTalentLearned     = "TALENT_LEARNED"
	EventAttributesSpent   = "ATTRIBUTES_ALLOCATED"
	EventRespec            = "SPECIALIZATION_CHANGED"
)

// XPSource names where XP came from. A non-empty ID makes the grant
// idempotent: the same (Kind, ID) pays once (character.xp_grant).
type XPSource struct {
	Kind     string         // hunt | dungeon | quest | combat | achievement | event | boss | pvp | dev
	ID       string         // contract / quest / event id
	Metadata map[string]any // audit only — never sensitive data
}

// Grant is the progression summary of one XP grant (or promotion).
type Grant struct {
	Character          *Character
	FromLevel          int
	LevelsGained       []int
	PromotionAvailable bool
	Events             []string // event kinds recorded by this call
	Replayed           bool     // the source had already paid
}

// XPGrantResult is the summary handed to reveal / turn-in callers.
type XPGrantResult struct {
	Character          *Character
	LeveledUp          bool
	FromLevel          int
	LevelsGained       []int
	PromotionAvailable bool
}

// MarshalJSON keeps the reveal payload to what clients need — never the full
// character row (account id, timestamps, internal state).
func (r XPGrantResult) MarshalJSON() ([]byte, error) {
	out := struct {
		Level              int16 `json:"level"`
		XP                 int64 `json:"xp"`
		LeveledUp          bool  `json:"leveled_up"`
		LevelsGained       int   `json:"levels_gained"`
		PromotionAvailable bool  `json:"promotion_available,omitempty"`
		OverflowXP         int64 `json:"overflow_xp,omitempty"`
	}{LeveledUp: r.LeveledUp, LevelsGained: len(r.LevelsGained), PromotionAvailable: r.PromotionAvailable}
	if r.Character != nil {
		out.Level, out.XP, out.OverflowXP = r.Character.Level, r.Character.XP, r.Character.OverflowXP
	}
	return json.Marshal(out)
}

// ProgressionStore is the transactional progression persistence (PgRepo).
type ProgressionStore interface {
	Promote(ctx context.Context, id uuid.UUID, choice string, cfg *leveling.Config) (*Grant, error)
	AllocateAttributes(ctx context.Context, id uuid.UUID, add leveling.Stats, cfg *leveling.Config) (*Character, error)
	Talents(ctx context.Context, id uuid.UUID) (map[string]int, error)
	LearnTalent(ctx context.Context, id uuid.UUID, code string, cfg *leveling.Config) (int, error)
	ChangeSpec(ctx context.Context, id uuid.UUID, spec string, cfg *leveling.Config) (*Character, error)
	Events(ctx context.Context, id uuid.UUID, limit int) ([]ProgressionEvent, error)
}

// ProgressionEvent is one audit row.
type ProgressionEvent struct {
	Kind     string         `json:"kind"`
	Ref      string         `json:"ref"`
	Level    int            `json:"level"`
	Metadata map[string]any `json:"metadata,omitempty"`
	At       string         `json:"at"`
}

func (s *Service) store() (ProgressionStore, error) {
	ps, ok := s.repo.(ProgressionStore)
	if !ok {
		return nil, ErrProgressionUnsupp
	}
	return ps, nil
}

// GrantXP is the one canonical way to award character XP (DD-250 §4):
// validates, applies the XP modifier, levels up as far as it reaches,
// stops at promotion gates (rest kept as overflow) and the cap, records
// level rewards / unlocks / titles / events — all in one transaction. A
// source with an ID pays once.
func (s *Service) GrantXP(ctx context.Context, characterID uuid.UUID, amount int64, src XPSource) (*XPGrantResult, error) {
	if amount < 0 {
		return nil, ErrInvalidXP
	}
	cfg := leveling.Current()
	amount = int64(float64(amount) * cfg.XPMultiplier)
	g, err := s.repo.GrantXP(ctx, characterID, amount, src, cfg)
	if err != nil {
		return nil, err
	}
	if len(g.LevelsGained) > 0 {
		s.log.InfoContext(ctx, "character leveled up",
			slog.String("character_id", characterID.String()),
			slog.Int("new_level", int(g.Character.Level)),
			slog.String("source", src.Kind),
		)
	}
	return &XPGrantResult{Character: g.Character, LeveledUp: len(g.LevelsGained) > 0, FromLevel: g.FromLevel,
		LevelsGained: g.LevelsGained, PromotionAvailable: g.PromotionAvailable}, nil
}

// SelectPromotion picks a class promotion at a reached gate, then processes
// the XP held while waiting (DD-250 §6).
func (s *Service) SelectPromotion(ctx context.Context, id, accountID uuid.UUID, choice string) (*Grant, error) {
	if _, err := s.GetOwned(ctx, id, accountID); err != nil {
		return nil, err
	}
	ps, err := s.store()
	if err != nil {
		return nil, err
	}
	return ps.Promote(ctx, id, choice, leveling.Current())
}

// AllocateAttributes spends attribute points (server-authoritative).
func (s *Service) AllocateAttributes(ctx context.Context, id, accountID uuid.UUID, add map[string]int) (*Character, error) {
	if _, err := s.GetOwned(ctx, id, accountID); err != nil {
		return nil, err
	}
	var st leveling.Stats
	total := 0
	for k, v := range add {
		if v < 0 {
			return nil, &ValidationError{Field: k, Message: "must be >= 0"}
		}
		if !st.Add(k, v) {
			return nil, ErrUnknownAttribute
		}
		total += v
	}
	if total == 0 {
		return nil, &ValidationError{Field: "points", Message: "allocate at least one point"}
	}
	ps, err := s.store()
	if err != nil {
		return nil, err
	}
	return ps.AllocateAttributes(ctx, id, st, leveling.Current())
}

// Talents returns the character's talent ranks.
func (s *Service) Talents(ctx context.Context, id, accountID uuid.UUID) (map[string]int, error) {
	if _, err := s.GetOwned(ctx, id, accountID); err != nil {
		return nil, err
	}
	ps, err := s.store()
	if err != nil {
		return nil, err
	}
	return ps.Talents(ctx, id)
}

// LearnTalent spends one talent point on code.
func (s *Service) LearnTalent(ctx context.Context, id, accountID uuid.UUID, code string) (int, error) {
	if _, err := s.GetOwned(ctx, id, accountID); err != nil {
		return 0, err
	}
	ps, err := s.store()
	if err != nil {
		return 0, err
	}
	return ps.LearnTalent(ctx, id, code, leveling.Current())
}

// BusyChecker reports whether the character is in combat / AFK / a dungeon.
type BusyChecker interface {
	Busy(ctx context.Context, characterID uuid.UUID) (bool, error)
}

// GearValidator unequips gear the character may no longer wear.
type GearValidator interface {
	UnequipInvalid(ctx context.Context, characterID uuid.UUID) (int, error)
}

// ChangeSpecialization is the respec boundary (DD-250 §35; no monetization
// yet): same tier, not while busy; keeps level, XP and earned titles;
// unequips gear the new specialization may not use.
func (s *Service) ChangeSpecialization(ctx context.Context, id uuid.UUID, spec string, busy BusyChecker, gear GearValidator) (*Character, int, error) {
	if busy != nil {
		b, err := busy.Busy(ctx, id)
		if err != nil {
			return nil, 0, err
		}
		if b {
			return nil, 0, ErrBusy
		}
	}
	ps, err := s.store()
	if err != nil {
		return nil, 0, err
	}
	c, err := ps.ChangeSpec(ctx, id, spec, leveling.Current())
	if err != nil {
		return nil, 0, err
	}
	removed := 0
	if gear != nil {
		if removed, err = gear.UnequipInvalid(ctx, id); err != nil {
			return c, 0, err
		}
	}
	return c, removed, nil
}

// ProgressionModifiers are the sheet modifiers of a character's progression:
// promotion stat modifiers along the class path, talent ranks and the class
// mastery at the level cap. (Skill passives come from classskill.)
func (s *Service) ProgressionModifiers(ctx context.Context, id uuid.UUID, sheet Sheet) ([]leveling.Modifier, error) {
	cfg := leveling.Current()
	spec := sheet.Class
	if sheet.Spec != nil && *sheet.Spec != "" {
		spec = *sheet.Spec
	}
	var out []leveling.Modifier
	for _, code := range cfg.Path(spec) {
		if k, ok := cfg.Class(code); ok {
			out = append(out, k.StatModifiers...)
		}
	}
	if k, ok := cfg.Class(spec); ok && k.Mastery != nil && sheet.Level >= cfg.Milestones.Mastery {
		out = append(out, k.Mastery.Modifiers...)
	}
	ps, ok := s.repo.(ProgressionStore)
	if !ok {
		return out, nil
	}
	ranks, err := ps.Talents(ctx, id)
	if err != nil {
		return nil, err
	}
	for code, rank := range ranks {
		t, _, ok := cfg.Talent(code)
		if !ok {
			continue // talent retired from content: no effect
		}
		for i := 0; i < rank; i++ {
			out = append(out, t.Modifiers...)
		}
	}
	return out, nil
}

// Events lists the latest progression events of an owned character.
func (s *Service) Events(ctx context.Context, id, accountID uuid.UUID, limit int) ([]ProgressionEvent, error) {
	if _, err := s.GetOwned(ctx, id, accountID); err != nil {
		return nil, err
	}
	ps, err := s.store()
	if err != nil {
		return nil, err
	}
	return ps.Events(ctx, id, limit)
}
