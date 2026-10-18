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
	require.Len(t, perClassStarter, 10, "every class starts with an active skill")
}

func TestValidateSkills_Rejects(t *testing.T) {
	active := func() ClassSkill {
		return ClassSkill{Code: "a", DisplayNameKey: "k", Class: "mage", UnlockLevel: 1, Kind: "active",
			Effect: "strike", Power: 150, MPCost: 5, CooldownTurns: 2, AutoTrigger: AutoTrigger{OnCooldown: true}}
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
	add("class", active, func(s *ClassSkill) { s.Class = "necromancer" })
	add("level", active, func(s *ClassSkill) { s.UnlockLevel = 0 })
	add("kind", active, func(s *ClassSkill) { s.Kind = "toggle" })
	add("effect", active, func(s *ClassSkill) { s.Effect = "explode" })
	add("untimed stun", active, func(s *ClassSkill) { s.Effect = "stun" })
	add("stunlock", active, func(s *ClassSkill) { s.Effect, s.DurationTurns, s.CooldownTurns = "stun", 2, 3 })
	add("no auto rule", active, func(s *ClassSkill) { s.AutoTrigger = AutoTrigger{} })
	add("type mismatch", active, func(s *ClassSkill) { s.Type = "HEAL" })
	add("unknown type", active, func(s *ClassSkill) { s.Type = "MIND_CONTROL" })
	add("barrier power", active, func(s *ClassSkill) { s.Effect, s.Power, s.DurationTurns = "barrier", 150, 3 })
	add("active passive model", active, func(s *ClassSkill) { s.Passive = &Passive{Trigger: "ON_HIT"} })
	add("passive trigger", passive, func(s *ClassSkill) { s.Passive = &Passive{Trigger: "ON_SNEEZE"} })
	add("passive auto rule", passive, func(s *ClassSkill) { s.AutoTrigger.OnCooldown = true })
	add("passive type", passive, func(s *ClassSkill) { s.Type = "DAMAGE" })
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
