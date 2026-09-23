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
}

type skillsFile struct {
	Skills []ClassSkill `yaml:"skills"`
}

var skillEffects = map[string]bool{"strike": true, "heal": true}

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
		if s.UnlockLevel < 1 || s.UnlockLevel > 60 {
			return fmt.Errorf("skill %s: unlock_level must be 1..60", s.Code)
		}
		switch s.Kind {
		case "active":
			if !skillEffects[s.Effect] {
				return fmt.Errorf("skill %s: effect %q not in strike|heal", s.Code, s.Effect)
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
			if s.Effect != "" || s.Power != 0 || s.MPCost != 0 {
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
  effect, power, mp_cost, cooldown_turns, modifiers
) VALUES ($1,$2,$3,$4,$5,$6,NULLIF($7,''),$8,$9,$10,$11)
ON CONFLICT (code) DO UPDATE SET
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
	for _, s := range skills {
		mods := s.Modifiers
		if mods == nil {
			mods = []StatBonus{}
		}
		raw, err := json.Marshal(mods)
		if err != nil {
			return fmt.Errorf("skill %s: marshal modifiers: %w", s.Code, err)
		}
		cd := max(s.CooldownTurns, 1)
		if _, err := tx.Exec(ctx, q, s.Code, s.DisplayNameKey, s.Class, s.UnlockLevel, s.Kind, s.SortOrder,
			s.Effect, s.Power, s.MPCost, cd, raw); err != nil {
			return fmt.Errorf("upsert skill %s: %w", s.Code, err)
		}
	}
	return tx.Commit(ctx)
}
