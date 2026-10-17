package contentpack

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"

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
}

// AutoTrigger is the passive-only (AFK) cast rule.
type AutoTrigger struct {
	SelfHPBelow   int32 `yaml:"self_hp_below"`
	TargetHPBelow int32 `yaml:"target_hp_below"`
}

type skillsFile struct {
	Skills []ClassSkill `yaml:"skills"`
}

var (
	skillEffects     = map[string]bool{"strike": true, "heal": true, "guard": true, "empower": true}
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
		for _, p := range []int32{s.AutoTrigger.SelfHPBelow, s.AutoTrigger.TargetHPBelow} {
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
			if !skillEffects[s.Effect] {
				return fmt.Errorf("skill %s: effect %q not in strike|heal|guard|empower", s.Code, s.Effect)
			}
			if (s.Effect == "guard" || s.Effect == "empower") != (s.DurationTurns > 0) {
				return fmt.Errorf("skill %s: guard/empower need duration_turns (and only they)", s.Code)
			}
			if s.ResourceGain < 0 {
				return fmt.Errorf("skill %s: resource_gain must be >= 0", s.Code)
			}
			if s.Power < 1 || (s.Effect == "heal" && s.Power > 100) || s.Power > 1000 {
				return fmt.Errorf("skill %s: power out of range", s.Code)
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
			"self_hp_below": s.AutoTrigger.SelfHPBelow, "target_hp_below": s.AutoTrigger.TargetHPBelow,
			"resource_gain": s.ResourceGain, "duration_turns": s.DurationTurns, "requires": s.Requires,
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
