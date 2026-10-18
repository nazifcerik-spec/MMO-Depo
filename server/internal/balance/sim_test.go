package balance

import (
	"fmt"
	"sort"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/otter/server/internal/activity"
	"github.com/otter/server/internal/classskill"
	"github.com/otter/server/internal/combat"
	"github.com/otter/server/internal/contentpack"
	"github.com/otter/server/internal/domain/character"
	"github.com/otter/server/internal/dungeon"
	"github.com/otter/server/internal/hunt"
	"github.com/otter/server/internal/leveling"
	"github.com/otter/server/internal/skill"
)

var classes = []string{"warrior", "rogue", "ranger", "mage", "monk", "cleric", "paladin", "druid", "bard", "shaman"}

type world struct {
	items    []contentpack.Item
	mobs     []combat.Mob
	dungeons []dungeon.Template
	skills   *classskill.Catalog
}

type mobIndex map[string]combat.Mob

func (m mobIndex) Get(code string) (combat.Mob, bool) { v, ok := m[code]; return v, ok }

func (w *world) mobIndex() mobIndex {
	idx := mobIndex{}
	for _, m := range w.mobs {
		idx[m.Code] = m
	}
	return idx
}

func load(t *testing.T) *world {
	t.Helper()
	p, err := contentpack.Read("../../../content")
	require.NoError(t, err)
	require.NoError(t, leveling.Install(p.Progression))
	w := &world{items: p.Items}
	for _, m := range p.Mobs {
		w.mobs = append(w.mobs, combat.Mob{Code: m.Code, Level: m.Level, HPMax: m.HPMax, Attack: m.Attack,
			Defense: m.Defense, CritPct: int32(m.CritPct), HitPct: int32(m.HitPct), EvasionPct: int32(m.EvasionPct),
			XPReward: m.XPReward, GoldCopperMin: m.GoldCopperMin, GoldCopperMax: m.GoldCopperMax})
	}
	sort.Slice(w.mobs, func(i, j int) bool { return w.mobs[i].Level < w.mobs[j].Level })
	for _, d := range p.Dungeons {
		t := dungeon.Template{Code: d.Code, MinLevel: d.MinLevel, BossMobCode: d.BossMobCode, DurationSecs: d.DurationSeconds}
		for _, wv := range d.Waves {
			t.Waves = append(t.Waves, dungeon.Wave{MobCode: wv.MobCode, Count: wv.Count})
		}
		w.dungeons = append(w.dungeons, t)
	}
	var list []classskill.Skill
	for _, s := range p.Skills {
		cs := classskill.Skill{Code: s.Code, Class: s.Class, UnlockLevel: s.UnlockLevel, Kind: s.Kind,
			SortOrder: s.SortOrder, Effect: s.Effect, Power: s.Power, MPCost: s.MPCost, CooldownTurns: int(s.CooldownTurns),
			Spec: s.Spec, Config: classskill.Config{SelfHPBelow: s.AutoTrigger.SelfHPBelow,
				TargetHPBelow: s.AutoTrigger.TargetHPBelow, ResourceGain: s.ResourceGain, DurationTurns: s.DurationTurns}}
		for _, m := range s.Modifiers {
			cs.Modifiers = append(cs.Modifiers, skill.Modifier{Kind: skill.ModifierKind(m.Kind), Target: skill.Target(m.Target), Amount: m.Amount})
		}
		list = append(list, cs)
	}
	w.skills = classskill.NewCatalog(list)
	return w
}

// kit is the best gear the content offers a class at a level: per slot, the
// usable item with the highest level requirement.
func (w *world) kit(class string, level int) []contentpack.Item {
	best := map[string]contentpack.Item{}
	for _, it := range w.items {
		if it.EquipSlot == "" || int(it.MinLevel) > level || (it.ClassRestriction != "" && it.ClassRestriction != class) {
			continue
		}
		if cur, ok := best[it.EquipSlot]; !ok || it.MinLevel > cur.MinLevel ||
			(it.MinLevel == cur.MinLevel && it.BaseValueCopper > cur.BaseValueCopper) {
			best[it.EquipSlot] = it
		}
	}
	out := make([]contentpack.Item, 0, len(best))
	for _, it := range best {
		out = append(out, it)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].EquipSlot < out[j].EquipSlot })
	return out
}

// simAlloc spends a class's attribute points the way a player of its role
// would (DD-250, DD-252): 2 primary : 1 first secondary attribute.
func simAlloc(class string, level int) leveling.Stats {
	cfg := leveling.Current()
	k, ok := cfg.Class(class)
	if !ok || k.BaseStats == nil || k.PrimaryAttribute == "" {
		return leveling.Stats{}
	}
	n := cfg.AttributePointsEarned(level)
	var st leveling.Stats
	st.Add(k.PrimaryAttribute, n-n/3)
	if len(k.SecondaryAttributes) > 0 {
		st.Add(k.SecondaryAttributes[0], n/3)
	} else {
		st.Add(k.PrimaryAttribute, n/3)
	}
	return st
}

// player builds the combatant the way hunts do: base sheet (class base +
// allocated points) → gear + class passives → derived stats (class formula
// when configured); actives for the auto-cast; the class resource.
func (w *world) player(class string, level int, gear []contentpack.Item) combat.Combatant {
	cfg := leveling.Current()
	ch := &character.Character{ClassCode: class, Level: int16(level), Alloc: simAlloc(class, level)}
	raw := character.RawAttributes(ch)
	b := character.PrimaryStats{Strength: raw.Strength, Agility: raw.Dexterity, Intelligence: raw.Intelligence,
		Vitality: raw.Vitality, Charisma: raw.Wisdom, Luck: raw.Luck}
	var mods []skill.Modifier
	for _, it := range gear {
		for _, sb := range it.StatBonuses {
			mods = append(mods, skill.Modifier{Kind: skill.ModifierKind(sb.Kind), Target: skill.Target(sb.Target), Amount: sb.Amount})
		}
	}
	mods = append(mods, w.skills.Passives([]string{class}, level)...)
	mods = append(mods, skill.FromLeveling(cfg.SoloAdaptationModifiers(class, false))...) // hunting alone
	var formula *skill.Formula
	k, hasClass := cfg.Class(class)
	if hasClass && k.Derived != nil {
		formula = &skill.Formula{Derived: *k.Derived, Level: level}
	}
	s := skill.ComputeSheetWith(skill.Primary{STR: b.Strength, AGI: b.Agility, INT: b.Intelligence, VIT: b.Vitality,
		CHA: b.Charisma, Luck: b.Luck}, skill.Target(character.MainStat(class)), mods, formula)
	d := s.Derived
	p := combat.Combatant{Name: "player", HPMax: int32(s.Resources.HPMax), HPCurrent: int32(s.Resources.HPMax),
		Attack: int32(d.Attack), Defense: int32(d.Defense), CritPct: int32(d.CritPct),
		HitPct: int32(d.HitPct), EvasionPct: int32(d.EvasionPct), MPMax: int32(s.Resources.MPMax),
		Skills:            w.skills.Actives([]string{class}, level),
		Level:             int16(level),
		RegenPerEncounter: int32(float64(encSeconds/2) * character.HPRegenPer2s(s.Primary.VIT)),
		BlockPct:          int32(d.BlockPct), CritDmgPct: int32(d.CritDmgPct), ExecutePct: int32(d.ExecutePct),
		MissingHPDmgPct: int32(d.MissingHPDmgPct), LifestealPct: int32(d.LifestealPct), DmgReductionPct: int32(d.DmgReductionPct)}
	combat.ApplyResource(&p, class)
	return p
}

// mobFor is what the client suggests: the strongest mob at or below the level.
func (w *world) mobFor(level int) combat.Mob {
	m := w.mobs[0]
	for _, c := range w.mobs {
		if int(c.Level) <= level {
			m = c
		}
	}
	return m
}

type outcome struct {
	deathRate   float64 // share of hour-long hunts that end in a defeat
	winsPerHour float64
	xpPerHour   float64
	minPerLevel float64 // minutes of hunting to reach the next level
}

const (
	seeds      = 150
	encSeconds = 120
	encPerHour = 30 // DefaultEncounterSeconds = 120
	maxTurns   = 50
)

func simulate(p combat.Combatant, m combat.Mob, level int) outcome {
	var o outcome
	m.XPReward = leveling.Current().ScaleMobXP(hunt.AFKXPPerVictory(encSeconds), level, int(m.Level)) // as the reveal pays (DD-250, DD-254)
	for s := int64(1); s <= seeds; s++ {
		r := hunt.Roll(s*7919, p, m, encPerHour, maxTurns, func(seed int64) combat.RNG { return activity.NewRNG(seed) })
		if r.StoppedOnDeath {
			o.deathRate++
		}
		o.winsPerHour += float64(r.Victories)
		o.xpPerHour += float64(r.TotalXP)
	}
	o.deathRate /= seeds
	o.winsPerHour /= seeds
	o.xpPerHour /= seeds
	if o.xpPerHour > 0 {
		o.minPerLevel = float64(character.XPToNext(level)) / o.xpPerHour * 60
	}
	return o
}

var levels = []int{1, 3, 5, 8, 10, 12, 15, 18, 20, 24, 25, 30, 35, 40, 45, 50, 55, 60, 65, 70, 75, 80, 85, 90, 95, 100}

// toughGap is how far above the character a mob must be dangerous.
func toughGap(level int) int {
	if level < 25 {
		return 5
	}
	return 10
}

// toughFor is the first mob at least gap levels above (nil when none).
func (w *world) toughFor(level, gap int) *combat.Mob {
	for i := range w.mobs {
		if int(w.mobs[i].Level) >= level+gap {
			return &w.mobs[i]
		}
	}
	return nil
}

// TestBalanceMatrix prints the class × level matrix (run with -v): death rate
// of an hour-long hunt at the suggested mob with the level's gear, without
// gear, and against a mob 5+ levels up; hunting minutes per level.
func TestBalanceMatrix(t *testing.T) {
	if !testing.Verbose() {
		t.Skip("prints a matrix: run with -v")
	}
	w := load(t)
	var b strings.Builder
	fmt.Fprintf(&b, "%-4s %-16s", "lvl", "mob")
	for _, c := range classes {
		fmt.Fprintf(&b, " | %-22s", c+" kit/bare/+5 min")
	}
	b.WriteString("\n")
	for _, lv := range levels {
		m := w.mobFor(lv)
		fmt.Fprintf(&b, "%-4d %-16s", lv, m.Code)
		for _, c := range classes {
			kit := simulate(w.player(c, lv, w.kit(c, lv)), m, lv)
			bare := simulate(w.player(c, lv, nil), m, lv)
			up := "  -"
			if tm := w.toughFor(lv, toughGap(lv)); tm != nil {
				up = fmt.Sprintf("%3.0f", simulate(w.player(c, lv, w.kit(c, lv)), *tm, lv).deathRate*100)
			}
			fmt.Fprintf(&b, " | %3.0f %3.0f %s %6.0f    ", kit.deathRate*100, bare.deathRate*100, up, kit.minPerLevel)
		}
		b.WriteString("\n")
	}
	t.Log("\n" + b.String())
}

// clearRate is the share of runs that kill the boss.
func (w *world) clearRate(class string, level int, d dungeon.Template) float64 {
	p, idx, n := w.player(class, level, w.kit(class, level)), w.mobIndex(), 0
	p.RegenPerEncounter = int32(float64(d.DurationSecs/int32(d.TotalEncounters())/6) * float64(p.RegenPerEncounter) / float64(encSeconds/2))
	for s := int64(1); s <= seeds; s++ {
		if dungeon.Roll(s*104729, p, d, idx, combat.DefaultMaxTurns, func(seed int64) combat.RNG { return activity.NewRNG(seed) }).Cleared {
			n++
		}
	}
	return float64(n) / seeds
}

// TestDungeonMatrix prints boss clear rates at the dungeon's level and above.
func TestDungeonMatrix(t *testing.T) {
	if !testing.Verbose() {
		t.Skip("prints a matrix: run with -v")
	}
	w := load(t)
	var b strings.Builder
	fmt.Fprintf(&b, "%-16s %-4s", "dungeon", "min")
	for _, c := range classes {
		fmt.Fprintf(&b, " | %-17s", c+" +0/+2/+4")
	}
	b.WriteString("\n")
	for _, d := range w.dungeons {
		fmt.Fprintf(&b, "%-16s %-4d", d.Code, d.MinLevel)
		for _, c := range classes {
			fmt.Fprintf(&b, " | %4.0f %4.0f %4.0f  ", 100*w.clearRate(c, int(d.MinLevel), d),
				100*w.clearRate(c, int(d.MinLevel)+2, d), 100*w.clearRate(c, int(d.MinLevel)+4, d))
		}
		b.WriteString("\n")
	}
	t.Log("\n" + b.String())
}

// TestBalanceTargets pins the DD-229 targets so content or formula changes
// that break the curve fail CI (tune with the matrices above).
func TestBalanceTargets(t *testing.T) {
	w := load(t)
	prevMin := 0.0
	for _, lv := range levels {
		m := w.mobFor(lv)
		for _, c := range classes {
			kit := simulate(w.player(c, lv, w.kit(c, lv)), m, lv)
			require.LessOrEqual(t, kit.deathRate, 0.10, "%s L%d: an hour at %s with level gear must be safe", c, lv, m.Code)
			// At-level AFK hunting earns close to the curve's baseline rate
			// (DD-250: the curve is timed at baseline_xp_per_minute).
			baseline := float64(character.XPToNext(lv)) / leveling.Current().Curve.BaselineXPPerMinute
			require.Greater(t, kit.minPerLevel, 0.0)
			require.LessOrEqual(t, kit.minPerLevel, baseline*1.25, "%s L%d: minutes per level vs baseline", c, lv)
			// Out-levelled mobs must be dangerous: 5 levels up in the starting
			// content; past it, sustain skills carry supports through +5, so
			// the bar is 10 levels up (DD-254).
			if tm := w.toughFor(lv, toughGap(lv)); tm != nil {
				up := simulate(w.player(c, lv, w.kit(c, lv)), *tm, lv)
				require.GreaterOrEqual(t, up.deathRate, 0.25, "%s L%d vs %s: out-levelled mobs must be dangerous", c, lv, tm.Code)
			}
			prevMin = max(prevMin, kit.minPerLevel)
		}
	}
	require.Greater(t, prevMin, 0.0)

	for _, d := range w.dungeons {
		var lo, hi float64 = 1, 0
		for _, c := range classes {
			at := w.clearRate(c, int(d.MinLevel), d)
			lo, hi = min(lo, at), max(hi, at)
			require.GreaterOrEqual(t, at, 0.30, "%s %s at min level: clearable", c, d.Code)
			require.GreaterOrEqual(t, w.clearRate(c, int(d.MinLevel)+4, d), 0.90, "%s %s four levels up: reliable", c, d.Code)
		}
		require.LessOrEqual(t, hi-lo, 0.40, "%s: class spread at min level", d.Code)
	}
}
