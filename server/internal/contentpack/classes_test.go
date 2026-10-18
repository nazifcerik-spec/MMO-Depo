package contentpack

import (
	"fmt"
	"slices"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/otter/server/internal/domain/character"
	"github.com/otter/server/internal/leveling"
)

// DD-252: the shipped class content — ten base classes, twenty promotions,
// forty specializations on one engine.

var baseClasses = []string{"warrior", "rogue", "ranger", "mage", "monk", "cleric", "paladin", "druid", "bard", "shaman"}

func shippedPack(t *testing.T) *Pack {
	t.Helper()
	p, err := Read("../../../content")
	require.NoError(t, err)
	require.NoError(t, leveling.Install(p.Progression))
	return p
}

func tierNodes(c *leveling.Config, tier int) []*leveling.Class {
	var out []*leveling.Class
	for _, k := range c.Classes {
		n, _ := c.Class(k.Code)
		if n.Tier() == tier {
			out = append(out, n)
		}
	}
	return out
}

func TestClassTreeShipped(t *testing.T) {
	p := shippedPack(t)
	c := p.Progression
	bases, promos, specs := tierNodes(c, 0), tierNodes(c, 1), tierNodes(c, 2)
	require.Len(t, bases, 10, "all 10 base classes load")
	require.Len(t, promos, 20, "all 20 promotions load")
	require.Len(t, specs, 40, "all 40 final specializations load")
	require.Empty(t, tierNodes(c, 3), "no class branches after level 300")
	for _, b := range bases {
		require.True(t, character.KnownClass(b.Code), b.Code)
		require.True(t, slices.Contains(baseClasses, b.Code), b.Code)
		kids := c.Choices(b.Code, 100)
		require.Len(t, kids, 2, "%s: two promotions at 100", b.Code)
		require.Empty(t, c.Choices(b.Code, 99))
		for _, k := range kids {
			require.Equal(t, 100, k.RequiredLevel)
			require.Equal(t, b.Code, k.Parent)
			gs := c.Choices(k.Code, 300)
			require.Len(t, gs, 2, "%s: two specializations at 300", k.Code)
			for _, g := range gs {
				require.Equal(t, []string{b.Code, k.Code, g.Code}, c.Path(g.Code), "parent chain (no cycles)")
				require.Empty(t, c.Choices(g.Code, 1000))
			}
		}
	}
}

func TestInvalidClassCombinationsRejected(t *testing.T) {
	c := shippedPack(t).Progression
	codes := func(list []*leveling.Class) []string {
		var out []string
		for _, k := range list {
			out = append(out, k.Code)
		}
		return out
	}
	require.NotContains(t, codes(c.Choices("mage", 100)), "guardian", "Mage → Guardian")
	require.NotContains(t, codes(c.Choices("cleric", 100)), "berserker", "Cleric → Berserker")
	for _, promo := range c.Choices("warrior", 100) {
		require.NotContains(t, codes(c.Choices(promo.Code, 300)), "pyromancer", "Warrior → Pyromancer")
	}
	for _, promo := range c.Choices("rogue", 100) {
		require.NotContains(t, codes(c.Choices(promo.Code, 300)), "beastmaster", "Rogue → Beastmaster")
	}
	require.NotEqual(t, "warrior", c.Path("pyromancer")[0])
}

// Every milestone of the shared progression has its skill (DD-252).
func TestRequiredSkillsPerMilestone(t *testing.T) {
	p := shippedPack(t)
	c := p.Progression
	type key struct {
		spec  string
		level int
	}
	cats := map[key][]string{}
	kinds := map[key][]string{}
	seen := map[string]bool{}
	for _, s := range p.Skills {
		require.False(t, seen[s.Code], "no skill unlocks twice: %s", s.Code)
		seen[s.Code] = true
		k := key{s.specCode(), int(s.UnlockLevel)}
		cats[k] = append(cats[k], s.category())
		kinds[k] = append(kinds[k], s.Kind)
		if s.Kind == "active" {
			require.True(t, s.AutoTrigger.defined(), "%s: passive-only rule", s.Code)
		}
	}
	has := func(spec string, level int, kind, cat string) bool {
		k := key{spec, level}
		for i := range kinds[k] {
			if (kind == "" || kinds[k][i] == kind) && (cat == "" || cats[k][i] == cat) {
				return true
			}
		}
		return false
	}
	for _, b := range tierNodes(c, 0) {
		for _, lv := range []int{1, 10, 25, 75} {
			require.True(t, has(b.Code, lv, "active", ""), "%s Lv%d skill", b.Code, lv)
		}
		require.True(t, has(b.Code, 50, "passive", ""), "%s Lv50 passive", b.Code)
	}
	for _, pr := range tierNodes(c, 1) {
		require.True(t, has(pr.Code, 100, "passive", ""), "%s promotion passive", pr.Code)
		for _, lv := range []int{100, 150, 200, 250} {
			require.True(t, has(pr.Code, lv, "", ""), "%s Lv%d skill", pr.Code, lv)
		}
	}
	for _, sp := range tierNodes(c, 2) {
		for _, lv := range []int{300, 450, 550} {
			require.True(t, has(sp.Code, lv, "passive", ""), "%s Lv%d passive", sp.Code, lv)
		}
		for _, lv := range []int{350, 400, 500} {
			require.True(t, has(sp.Code, lv, "", ""), "%s Lv%d skill", sp.Code, lv)
		}
		require.True(t, has(sp.Code, 600, "", "awakening"), "%s awakening", sp.Code)
		require.True(t, has(sp.Code, 700, "active", "ultimate"), "%s ultimate", sp.Code)
		n := 0
		for k, list := range cats {
			if k.spec == sp.Code {
				for _, cat := range list {
					if cat == "ultimate" {
						n++
					}
				}
			}
		}
		require.Equal(t, 1, n, "%s: exactly one ultimate", sp.Code)
		// No unlock reaches a path twice.
		path := c.Path(sp.Code)
		codes := map[string]bool{}
		for _, u := range c.UnlockedBy(path, 1000) {
			require.False(t, codes[u.Code], u.Code)
			codes[u.Code] = true
		}
	}
	for _, s := range p.Skills {
		if s.category() == "ultimate" {
			require.GreaterOrEqual(t, s.CooldownTurns, int32(15), "%s: long cooldown", s.Code)
		}
	}
}

func TestTalentTreesAndMasteries(t *testing.T) {
	c := shippedPack(t).Progression
	for _, b := range tierNodes(c, 0) {
		trees := c.TreesFor(b.Code)
		require.Len(t, trees, 3, "%s: three talent trees", b.Code)
		for _, tr := range trees {
			points, caps := 0, 0
			for _, tl := range tr.Talents {
				if tl.Capstone {
					caps++
				} else {
					points += tl.MaxRank
				}
			}
			require.Equal(t, 25, points, tr.Code)
			require.Equal(t, 1, caps, "%s: capstone", tr.Code)
		}
	}
	require.Equal(t, 45, c.TalentPointsEarned(1000), "a character cannot fill all three trees")
	masteries := map[string]bool{}
	for _, sp := range tierNodes(c, 2) {
		m := sp.Mastery
		require.NotNil(t, m, "%s: Lv1000 mastery", sp.Code)
		require.NotEmpty(t, m.Title)
		require.NotEmpty(t, m.Badge)
		require.NotEmpty(t, m.Achievement)
		require.True(t, m.PrestigeEligible)
		require.False(t, masteries[m.Title], "no mastery title twice: %s", m.Title)
		masteries[m.Title] = true
		require.NotEmpty(t, m.Modifiers, sp.Code)
		for _, mod := range m.Modifiers {
			v := mod.Amount
			if mod.Kind == "mult_add" {
				v *= 100
			}
			require.True(t, v > 0 && v <= 3, "%s: small mastery bonus (%s %.1f)", sp.Code, mod.Target, v)
		}
		require.NotNil(t, sp.AFKProfile, "%s: AFK profile", sp.Code)
		require.NotEmpty(t, sp.StrongIn, sp.Code)
		require.NotEmpty(t, sp.WeakIn, sp.Code)
	}
	n := 0
	for _, r := range c.RewardsAt(1000) {
		if r.Type == leveling.RewardMastery {
			n++
		}
	}
	require.Equal(t, 1, n, "mastery unlocks once")
}

func TestResourcesAndEquipmentPermissions(t *testing.T) {
	p := shippedPack(t)
	c := p.Progression
	kinds := map[string]bool{}
	for _, b := range tierNodes(c, 0) {
		r := c.ResourceOf(b.Code)
		require.NotNil(t, r, b.Code)
		require.True(t, leveling.ResourceKinds[r.Kind], "%s resource %s", b.Code, r.Kind)
		kinds[r.Kind] = true
		require.NotNil(t, b.BaseStats, b.Code)
		require.NotNil(t, b.Derived, b.Code)
		require.NotEmpty(t, b.PrimaryAttribute, b.Code)
		armor, weapon := false, false
		for _, perm := range b.EquipmentPermissions {
			if c.IsFamily(perm) {
				weapon = true
			} else {
				armor = true
			}
		}
		require.True(t, armor && weapon, "%s: armor type and weapon families", b.Code)
	}
	require.Len(t, kinds, 9, "rage energy focus mana chi conviction nature_energy rhythm spirit")
	// Every class's own weapons are usable by its whole tree.
	for _, it := range p.Items {
		if it.ClassRestriction == "" || it.ItemClass != "weapon" {
			continue
		}
		require.True(t, c.IsFamily(it.ItemSubclass), "%s: governed family", it.Code)
		for _, k := range c.Classes {
			n, _ := c.Class(k.Code)
			if n.Root() == it.ClassRestriction {
				require.True(t, c.CanUse(k.Code, it.ItemSubclass), "%s cannot use %s", k.Code, it.Code)
			}
		}
	}
	require.False(t, c.CanUse("mage", "sword"))
	require.True(t, c.CanUse("mage", "helm"), "armor pieces are not governed")
}

func TestSoloAdaptation(t *testing.T) {
	c := shippedPack(t).Progression
	require.Equal(t, 0.35, c.SoloAdaptation.Ratio)
	solo := c.SoloAdaptationModifiers("virtuoso", false)
	require.Len(t, solo, 1)
	require.Equal(t, "attack", solo[0].Target)
	require.InDelta(t, 0.35*c.SupportPower("virtuoso"), solo[0].Amount, 1e-9)
	require.Empty(t, c.SoloAdaptationModifiers("virtuoso", true), "gone in a party")
	require.Empty(t, c.SoloAdaptationModifiers("berserker", false), "pure DPS has none")
	require.Greater(t, c.SupportPower("cleric"), 0.0, "base support classes adapt before promotion")
}

// Localization keys for every class node, skill, talent, title and
// resource in all four languages (DD-252).
func TestClassLocalizationComplete(t *testing.T) {
	p := shippedPack(t)
	keys := progressionKeys(p)
	require.NotEmpty(t, keys)
	for _, lang := range []string{"en", "tr", "zh-cn", "es"} {
		m := p.Strings.ByLang[lang]
		require.NotNil(t, m, lang)
		for _, k := range keys {
			require.NotEmpty(t, m[k], "%s: %s", lang, k)
		}
	}
	// Spot checks: translated, not English fallbacks.
	for lang, want := range map[string]string{"en": "Nightblade", "tr": "Gece Bıçağı", "zh-cn": "夜刃", "es": "Hoja Nocturna"} {
		require.Equal(t, want, p.Strings.ByLang[lang]["class.nightblade"], lang)
	}
	for lang, want := range map[string]string{"en": "Time Stop", "tr": "Zamanı Durdur", "zh-cn": "时间停止", "es": "Detener el Tiempo"} {
		require.Equal(t, want, p.Strings.ByLang[lang]["skill.chronomancer_time_stop"], lang)
	}
	require.Contains(t, keys, "skill.nightblade_shadow_strike.desc")
	require.Contains(t, keys, "title.master_of_shadows")
	require.Contains(t, keys, "resource.nature_energy")
	require.Contains(t, keys, fmt.Sprintf("achievement.mastery_%s", "hexer"))
}
