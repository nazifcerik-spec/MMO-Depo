package contentpack

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestShippedSkillsValidate(t *testing.T) {
	skills, err := loadSkillsFile(filepath.Join("..", "..", "..", "content", "skills.yaml"))
	require.NoError(t, err)
	perClassStarter := map[string]bool{}
	for _, s := range skills {
		if s.Kind == "active" && s.UnlockLevel == 1 {
			perClassStarter[s.Class] = true
		}
	}
	require.Len(t, perClassStarter, 5, "every class starts with an active skill")
}

func TestValidateSkills_Rejects(t *testing.T) {
	active := func() ClassSkill {
		return ClassSkill{Code: "a", DisplayNameKey: "k", Class: "mage", UnlockLevel: 1, Kind: "active",
			Effect: "strike", Power: 150, MPCost: 5, CooldownTurns: 2}
	}
	passive := func() ClassSkill {
		return ClassSkill{Code: "p", DisplayNameKey: "k", Class: "rogue", UnlockLevel: 3, Kind: "passive",
			Modifiers: []StatBonus{{Kind: "flat", Target: "crit_pct", Amount: 5}}}
	}
	require.NoError(t, validateSkills([]ClassSkill{active(), passive()}))

	cases := map[string]ClassSkill{}
	add := func(name string, base func() ClassSkill, mut func(*ClassSkill)) {
		s := base()
		mut(&s)
		cases[name] = s
	}
	add("class", active, func(s *ClassSkill) { s.Class = "bard" })
	add("level", active, func(s *ClassSkill) { s.UnlockLevel = 0 })
	add("kind", active, func(s *ClassSkill) { s.Kind = "toggle" })
	add("effect", active, func(s *ClassSkill) { s.Effect = "stun" })
	add("power", active, func(s *ClassSkill) { s.Power = 0 })
	add("heal power", active, func(s *ClassSkill) { s.Effect, s.Power = "heal", 150 })
	add("cooldown", active, func(s *ClassSkill) { s.CooldownTurns = 0 })
	add("active mods", active, func(s *ClassSkill) { s.Modifiers = passive().Modifiers })
	add("passive effect", passive, func(s *ClassSkill) { s.Effect = "strike" })
	add("passive empty", passive, func(s *ClassSkill) { s.Modifiers = nil })
	add("en_max typo", passive, func(s *ClassSkill) { s.Modifiers[0].Target = "en_max" })
	add("bad kind", passive, func(s *ClassSkill) { s.Modifiers[0].Kind = "add" })
	add("no name", passive, func(s *ClassSkill) { s.DisplayNameKey = "" })
	for name, s := range cases {
		require.Error(t, validateSkills([]ClassSkill{s}), name)
	}
	require.ErrorContains(t, validateSkills([]ClassSkill{active(), active()}), "duplicate")
}
