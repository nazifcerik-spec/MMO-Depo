package contentpack

import (
	"errors"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"

	"github.com/otter/server/internal/domain/character"
	"github.com/otter/server/internal/leveling"
)

// loadProgressionFile reads progression.yaml (DD-250) and joins the class
// skill unlocks to it. A missing file keeps the engine's built-in defaults
// (nil config).
func loadProgressionFile(path string, skills []ClassSkill) (*leveling.Config, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}
		return nil, fmt.Errorf("read %s: %w", path, err)
	}
	var c leveling.Config
	if err := yaml.Unmarshal(raw, &c); err != nil {
		return nil, fmt.Errorf("parse %s: %w", path, err)
	}
	if c.LevelCap > character.LevelCap {
		return nil, fmt.Errorf("validate %s: level_cap %d above the database limit %d", path, c.LevelCap, character.LevelCap)
	}
	checkMods := func(owner string, mods []leveling.Modifier) error {
		for i, m := range mods {
			if !validModifierKinds(m.Kind) || !validStatTargets(m.Target) {
				return fmt.Errorf("%s: modifiers[%d] %s/%s unknown", owner, i, m.Kind, m.Target)
			}
		}
		return nil
	}
	for _, k := range c.Classes {
		if k.Parent == "" && !character.KnownClass(k.Code) {
			return nil, fmt.Errorf("validate %s: base class %s is not a playable class", path, k.Code)
		}
		if err := checkMods("class "+k.Code, k.StatModifiers); err != nil {
			return nil, err
		}
		if k.Mastery != nil {
			if err := checkMods("class "+k.Code+" mastery", k.Mastery.Modifiers); err != nil {
				return nil, err
			}
		}
	}
	for _, t := range c.Trees {
		for _, tl := range t.Talents {
			if err := checkMods("talent "+tl.Code, tl.Modifiers); err != nil {
				return nil, err
			}
		}
	}
	c.Unlocks = skillUnlocks(skills)
	if err := c.Prepare(); err != nil {
		return nil, fmt.Errorf("validate %s: %w", path, err)
	}
	// Each skill belongs to a node of its own class's tree (or, for classes
	// without one, to the class itself).
	for _, s := range skills {
		spec := s.specCode()
		k, ok := c.Class(spec)
		if _, hasTree := c.Class(s.Class); !hasTree {
			ok = spec == s.Class
		} else if !ok || k.Root() != s.Class {
			ok = false
		}
		if !ok {
			return nil, fmt.Errorf("validate skills.yaml: skill %s: spec %q is not a %s class", s.Code, spec, s.Class)
		}
		if k != nil && int(s.UnlockLevel) < k.RequiredLevel {
			return nil, fmt.Errorf("validate skills.yaml: skill %s unlocks before its class (%d < %d)",
				s.Code, s.UnlockLevel, k.RequiredLevel)
		}
	}
	return &c, nil
}

func (s ClassSkill) specCode() string {
	if s.Spec != "" {
		return s.Spec
	}
	return s.Class
}

func (s ClassSkill) category() string {
	switch {
	case s.Category != "":
		return s.Category
	case s.Kind == "passive":
		return "passive"
	}
	return "skill"
}

func skillUnlocks(skills []ClassSkill) []leveling.Unlock {
	out := make([]leveling.Unlock, 0, len(skills))
	for _, s := range skills {
		out = append(out, leveling.Unlock{Code: s.Code, Spec: s.specCode(), Level: int(s.UnlockLevel), Category: s.category()})
	}
	return out
}
