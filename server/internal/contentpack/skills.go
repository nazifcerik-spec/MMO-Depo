package contentpack

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"slices"

	"gopkg.in/yaml.v3"

	"github.com/otter/server/internal/combat"
	"github.com/otter/server/internal/domain/character"
)

// ClassSkill mirrors content.class_skill (DD-222).
type ClassSkill struct {
	Code           string      `yaml:"code"`
	DisplayNameKey string      `yaml:"display_name_key"`
	Class          string      `yaml:"class"`
	UnlockLevel    int16       `yaml:"unlock_level"`
	Kind           string      `yaml:"kind"`
	SortOrder      int32       `yaml:"sort_order"`
	Effect         string      `yaml:"effect"`
	Power          int32       `yaml:"power"`
	MPCost         int32       `yaml:"mp_cost"`
	CooldownTurns  int32       `yaml:"cooldown_turns"`
	Modifiers      []StatBonus `yaml:"modifiers"`

	// DD-250: owning class node, category and data-driven activation.
	Spec           string      `yaml:"spec"`
	Category       string      `yaml:"category"`        // skill | passive | awakening | ultimate
	ActivationMode string      `yaml:"activation_mode"` // active | stance | passive
	AutoTrigger    AutoTrigger `yaml:"auto_trigger"`
	ResourceGain   int32       `yaml:"resource_gain"`
	DurationTurns  int         `yaml:"duration_turns"`
	Requires       []string    `yaml:"requires"`

	// Shared skill data model (DD-252): descriptive fields the resolver
	// does not need but admin, client and balance tools read.
	Type       string   `yaml:"type"`        // DAMAGE | HEAL | BARRIER | BUFF | DEBUFF | DOT | HOT | STUN | ...
	DamageType string   `yaml:"damage_type"` // physical | magic | holy | nature | fire | frost | ...
	TargetType string   `yaml:"target_type"` // enemy | self | ally | party | area
	Range      string   `yaml:"range"`       // melee | ranged
	CastTime   int      `yaml:"cast_time"`   // turns (0 = instant)
	Tags       []string `yaml:"tags"`
	Passive    *Passive `yaml:"passive"` // passive trigger model
}

// AutoTrigger is the passive-only (AFK) cast rule; one skill definition
// serves both modes. OnCooldown casts whenever ready and affordable.
type AutoTrigger struct {
	OnCooldown    bool  `yaml:"on_cooldown"`
	SelfHPBelow   int32 `yaml:"self_hp_below"`
	TargetHPBelow int32 `yaml:"target_hp_below"`
	AllyHPBelow   int32 `yaml:"ally_hp_below"` // solo: the character is its own ally
}

func (a AutoTrigger) defined() bool {
	return a.OnCooldown || a.SelfHPBelow > 0 || a.TargetHPBelow > 0 || a.AllyHPBelow > 0
}

// Passive is the shared passive model (DD-252). The applied effect is the
// skill's modifiers (steady-state value); the rest documents when it acts.
type Passive struct {
	Trigger          string `yaml:"trigger" json:"trigger"`
	Condition        string `yaml:"condition" json:"condition,omitempty"`
	InternalCooldown int    `yaml:"internal_cooldown" json:"internal_cooldown,omitempty"`
	StackLimit       int    `yaml:"stack_limit" json:"stack_limit,omitempty"`
	Duration         int    `yaml:"duration" json:"duration,omitempty"`
	Scaling          string `yaml:"scaling" json:"scaling,omitempty"`
}

type skillsFile struct {
	Skills []ClassSkill `yaml:"skills"`
}

// skillTypes maps each data-model effect type to the resolver effects that
// may implement it (DD-252).
var skillTypes = map[string][]string{
	"DAMAGE":             {"strike"},
	"HEAL":               {"heal"},
	"BARRIER":            {"barrier"},
	"BUFF":               {"empower", "guard"},
	"DEBUFF":             {"weaken"},
	"DOT":                {"dot"},
	"HOT":                {"hot"},
	"STUN":               {"stun"},
	"ROOT":               {"stun"},
	"SILENCE":            {"stun", "weaken"},
	"TAUNT":              {"guard"},
	"DISPEL":             {"heal", "barrier"},
	"RESOURCE_GAIN":      {"strike", "empower"},
	"RESOURCE_DRAIN":     {"weaken", "strike"},
	"COOLDOWN_REDUCTION": {"empower"},
	"STAT_MODIFIER":      {},
}

// PassiveTriggers are the passive model's triggers.
var PassiveTriggers = map[string]bool{
	"ON_ATTACK": true, "ON_HIT": true, "ON_CRIT": true, "ON_BLOCK": true, "ON_DODGE": true, "ON_HEAL": true,
	"ON_DAMAGE_TAKEN": true, "ON_KILL": true, "ON_DEBUFF_APPLIED": true, "ON_BUFF_APPLIED": true,
	"ON_RESOURCE_SPENT": true, "ON_RESOURCE_GAINED": true, "ON_SKILL_CAST": true, "ON_LOW_HP": true,
	"ON_COMBAT_START": true, "ON_COMBAT_END": true,
}

var (
	skillCategories  = map[string]bool{"": true, "skill": true, "passive": true, "awakening": true, "ultimate": true}
	activationModes  = map[string]bool{"": true, "active": true, "stance": true, "passive": true}
	equipRequirement = map[string]bool{"shield": true, "two_handed": true}
)

func loadSkillsFile(path string) ([]ClassSkill, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}
		return nil, fmt.Errorf("read %s: %w", path, err)
	}
	var f skillsFile
	if err := yaml.Unmarshal(raw, &f); err != nil {
		return nil, fmt.Errorf("parse %s: %w", path, err)
	}
	if err := validateSkills(f.Skills); err != nil {
		return nil, fmt.Errorf("validate %s: %w", path, err)
	}
	return f.Skills, nil
}

func validateSkills(skills []ClassSkill) error {
	seen := make(map[string]bool, len(skills))
	for i, s := range skills {
		if s.Code == "" {
			return fmt.Errorf("skill[%d]: empty code", i)
		}
		if seen[s.Code] {
			return fmt.Errorf("skill[%d] %s: duplicate code", i, s.Code)
		}
		seen[s.Code] = true
		if s.DisplayNameKey == "" {
			return fmt.Errorf("skill %s: display_name_key required", s.Code)
		}
		if !character.KnownClass(s.Class) {
			return fmt.Errorf("skill %s: unknown class %q", s.Code, s.Class)
		}
		if s.UnlockLevel < 1 || int(s.UnlockLevel) > character.LevelCap {
			return fmt.Errorf("skill %s: unlock_level must be 1..%d", s.Code, character.LevelCap)
		}
		if !skillCategories[s.Category] || !activationModes[s.ActivationMode] {
			return fmt.Errorf("skill %s: unknown category %q or activation_mode %q", s.Code, s.Category, s.ActivationMode)
		}
		if s.Type != "" {
			allowed, ok := skillTypes[s.Type]
			if !ok || (s.Kind == "active" && !slices.Contains(allowed, s.Effect)) || (s.Kind == "passive" && s.Type != "STAT_MODIFIER") {
				return fmt.Errorf("skill %s: type %q does not fit %s/%s", s.Code, s.Type, s.Kind, s.Effect)
			}
		}
		for _, p := range []int32{s.AutoTrigger.SelfHPBelow, s.AutoTrigger.TargetHPBelow, s.AutoTrigger.AllyHPBelow} {
			if p < 0 || p > 100 {
				return fmt.Errorf("skill %s: auto_trigger thresholds are 0..100", s.Code)
			}
		}
		for _, r := range s.Requires {
			if !equipRequirement[r] {
				return fmt.Errorf("skill %s: unknown requirement %q", s.Code, r)
			}
		}
		switch s.Kind {
		case "active":
			if !combat.Effects[s.Effect] {
				return fmt.Errorf("skill %s: unknown effect %q", s.Code, s.Effect)
			}
			if combat.Timed(s.Effect) != (s.DurationTurns > 0) {
				return fmt.Errorf("skill %s: timed effects need duration_turns (and only they)", s.Code)
			}
			if !s.AutoTrigger.defined() {
				return fmt.Errorf("skill %s: active skills need an auto_trigger rule (passive-only play)", s.Code)
			}
			if s.Passive != nil {
				return fmt.Errorf("skill %s: actives take no passive model", s.Code)
			}
			if s.ResourceGain < 0 {
				return fmt.Errorf("skill %s: resource_gain must be >= 0", s.Code)
			}
			bounded := s.Effect == "heal" || s.Effect == "hot" || s.Effect == "barrier" || s.Effect == "weaken"
			if s.Power < 1 || (bounded && s.Power > 100) || s.Power > 1000 {
				return fmt.Errorf("skill %s: power out of range", s.Code)
			}
			if s.Effect == "stun" && int(s.CooldownTurns) < s.DurationTurns+2 {
				return fmt.Errorf("skill %s: stun cooldown must exceed its duration by 2 (no stunlock)", s.Code)
			}
			if s.MPCost < 0 || s.CooldownTurns < 1 {
				return fmt.Errorf("skill %s: mp_cost >= 0 and cooldown_turns >= 1 required", s.Code)
			}
			if len(s.Modifiers) > 0 {
				return fmt.Errorf("skill %s: actives take no modifiers", s.Code)
			}
		case "passive":
			if s.Effect != "" || s.Power != 0 || s.MPCost != 0 || s.DurationTurns != 0 || s.ResourceGain != 0 {
				return fmt.Errorf("skill %s: passives take only modifiers", s.Code)
			}
			if len(s.Modifiers) == 0 {
				return fmt.Errorf("skill %s: passive needs modifiers", s.Code)
			}
			if s.AutoTrigger.defined() {
				return fmt.Errorf("skill %s: passives take no auto_trigger", s.Code)
			}
			if s.Passive != nil && (!PassiveTriggers[s.Passive.Trigger] || s.Passive.StackLimit < 0 ||
				s.Passive.InternalCooldown < 0 || s.Passive.Duration < 0) {
				return fmt.Errorf("skill %s: invalid passive model", s.Code)
			}
			for j, m := range s.Modifiers {
				if !validModifierKinds(m.Kind) || !validStatTargets(m.Target) {
					return fmt.Errorf("skill %s: modifiers[%d] %s/%s unknown", s.Code, j, m.Kind, m.Target)
				}
			}
		default:
			return fmt.Errorf("skill %s: kind %q not in active|passive", s.Code, s.Kind)
		}
	}
	return nil
}

func (l *Loader) upsertSkills(ctx context.Context, skills []ClassSkill) error {
	tx, err := l.DB.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	const q = `
INSERT INTO content.class_skill (
  code, display_name_key, class_code, unlock_level, kind, sort_order,
  effect, power, mp_cost, cooldown_turns, modifiers, spec_code, config
) VALUES ($1,$2,$3,$4,$5,$6,NULLIF($7,''),$8,$9,$10,$11,NULLIF($12,''),$13)
ON CONFLICT (code) DO UPDATE SET
  spec_code        = EXCLUDED.spec_code,
  config           = EXCLUDED.config,
  display_name_key = EXCLUDED.display_name_key,
  class_code       = EXCLUDED.class_code,
  unlock_level     = EXCLUDED.unlock_level,
  kind             = EXCLUDED.kind,
  sort_order       = EXCLUDED.sort_order,
  effect           = EXCLUDED.effect,
  power            = EXCLUDED.power,
  mp_cost          = EXCLUDED.mp_cost,
  cooldown_turns   = EXCLUDED.cooldown_turns,
  modifiers        = EXCLUDED.modifiers,
  version          = content.class_skill.version + 1`
	codes := make([]string, 0, len(skills))
	for _, s := range skills {
		mods := s.Modifiers
		if mods == nil {
			mods = []StatBonus{}
		}
		raw, err := json.Marshal(mods)
		if err != nil {
			return fmt.Errorf("skill %s: marshal modifiers: %w", s.Code, err)
		}
		cfg, err := json.Marshal(map[string]any{
			"category": s.category(), "activation_mode": s.ActivationMode,
			"on_cooldown": s.AutoTrigger.OnCooldown, "ally_hp_below": s.AutoTrigger.AllyHPBelow,
			"self_hp_below": s.AutoTrigger.SelfHPBelow, "target_hp_below": s.AutoTrigger.TargetHPBelow,
			"resource_gain": s.ResourceGain, "duration_turns": s.DurationTurns, "requires": s.Requires,
			"type": s.Type, "damage_type": s.DamageType, "target_type": s.TargetType, "range": s.Range,
			"cast_time": s.CastTime, "tags": s.Tags, "passive": s.Passive,
		})
		if err != nil {
			return err
		}
		cd := max(s.CooldownTurns, 1)
		if _, err := tx.Exec(ctx, q, s.Code, s.DisplayNameKey, s.Class, s.UnlockLevel, s.Kind, s.SortOrder,
			s.Effect, s.Power, s.MPCost, cd, raw, s.Spec, cfg); err != nil {
			return fmt.Errorf("upsert skill %s: %w", s.Code, err)
		}
		codes = append(codes, s.Code)
	}
	// Skills removed from content stop existing (no per-character state refers to them).
	if _, err := tx.Exec(ctx, `DELETE FROM content.class_skill WHERE code <> ALL($1)`, codes); err != nil {
		return fmt.Errorf("prune skills: %w", err)
	}
	return tx.Commit(ctx)
}
