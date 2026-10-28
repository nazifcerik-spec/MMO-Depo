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
	"github.com/otter/server/internal/inventory"
	"github.com/otter/server/internal/leveling"
	"github.com/otter/server/internal/skill"
)

var classes = []string{"warrior", "rogue", "ranger", "mage", "monk", "cleric", "paladin", "druid", "bard", "shaman"}

type world struct {
	items    []contentpack.Item
	mobs     []combat.Mob
	dungeons []dungeon.Template
	raids    []dungeon.Template // group dungeons (DD-261)
	raidMobs map[string]combat.Mob
	skills   *classskill.Catalog
}

type mobIndex map[string]combat.Mob

func (m mobIndex) Get(code string) (combat.Mob, bool) { v, ok := m[code]; return v, ok }

func (w *world) mobIndex() mobIndex {
	idx := mobIndex{}
	for _, m := range w.mobs {
		idx[m.Code] = m
	}
	for code, m := range w.raidMobs {
		idx[code] = m
	}
	return idx
}

func load(t *testing.T) *world {
	t.Helper()
	p, err := contentpack.Read("../../../content")
	require.NoError(t, err)
	require.NoError(t, leveling.Install(p.Progression))
	// Level gear is what drops and is crafted: raid set pieces are a reward
	// on top, not the baseline the targets are tuned for (DD-277).
	inSet := map[string]bool{}
	for _, s := range p.ItemSets {
		for _, piece := range s.Pieces {
			inSet[piece] = true
		}
	}
	w := &world{}
	for _, it := range p.Items {
		if !inSet[it.Code] {
			w.items = append(w.items, it)
		}
	}
	for _, m := range p.Mobs {
		w.mobs = append(w.mobs, combat.Mob{Code: m.Code, Level: m.Level, HPMax: m.HPMax, Attack: m.Attack,
			Defense: m.Defense, CritPct: int32(m.CritPct), HitPct: int32(m.HitPct), EvasionPct: int32(m.EvasionPct),
			XPReward: m.XPReward, GoldCopperMin: m.GoldCopperMin, GoldCopperMax: m.GoldCopperMax, RaidBoss: m.RaidBoss})
	}
	w.raidMobs = map[string]combat.Mob{}
	for i := 0; i < len(w.mobs); i++ { // raid bosses are not hunted (DD-261)
		if w.mobs[i].RaidBoss {
			w.raidMobs[w.mobs[i].Code] = w.mobs[i]
			w.mobs = append(w.mobs[:i], w.mobs[i+1:]...)
			i--
		}
	}
	sort.Slice(w.mobs, func(i, j int) bool { return w.mobs[i].Level < w.mobs[j].Level })
	for _, d := range p.Dungeons {
		t := dungeon.Template{Code: d.Code, MinLevel: d.MinLevel, BossMobCode: d.BossMobCode, DurationSecs: d.DurationSeconds,
			GroupMin: d.GroupMin, GroupMax: d.GroupMax, BonusXP: d.BonusXP, BonusGoldCopper: d.BonusGoldCopper, PowerPct: d.PowerPct}
		for _, wv := range d.Waves {
			t.Waves = append(t.Waves, dungeon.Wave{MobCode: wv.MobCode, Count: wv.Count})
		}
		if t.IsRaid() {
			w.raids = append(w.raids, t)
			continue
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
func (w *world) kit(build string, level int) []contentpack.Item {
	class, _, _ := strings.Cut(build, "/")
	best := map[string]contentpack.Item{}
	for _, it := range w.items {
		if it.EquipSlot == "" || int(it.MinLevel) > level || (it.ClassRestriction != "" && it.ClassRestriction != class) ||
			!leveling.Current().CanUse(class, it.ItemSubclass) {
			continue
		}
		if cur, ok := best[it.EquipSlot]; !ok || it.MinLevel > cur.MinLevel ||
			(it.MinLevel == cur.MinLevel && it.BaseValueCopper > cur.BaseValueCopper) {
			best[it.EquipSlot] = it
		}
	}
	out := make([]contentpack.Item, 0, len(best)+2)
	for slot, it := range best {
		out = append(out, it)
		if _, paired := inventory.PairedSlots[slot]; paired {
			out = append(out, it) // two rings, two trinkets (DD-276)
		}
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

// builds are what the targets test at a level: every class path open at it
// — the base class below 100, its two promotions from 100
// ("class/promotion"), the four specializations from 300 (DD-255, DD-256).
func builds(level int) []string {
	cfg := leveling.Current()
	var expand func(path string, node string) []string
	expand = func(path, node string) []string {
		kids := cfg.Choices(node, level)
		if len(kids) == 0 {
			return []string{path}
		}
		var out []string
		for _, k := range kids {
			out = append(out, expand(path+"/"+k.Code, k.Code)...)
		}
		return out
	}
	var out []string
	for _, c := range classes {
		out = append(out, expand(c, c)...)
	}
	return out
}

// player builds the combatant the way hunts do: base sheet (class base +
// allocated points) → gear + path passives → derived stats (class formula
// when configured); path actives for the auto-cast; the class resource.
// build is a class or "class/promotion".
func (w *world) player(build string, level int, gear []contentpack.Item) combat.Combatant {
	return w.playerIn(build, level, gear, false)
}

// playerIn is player with the party state Solo Adaptation reads (DD-252).
func (w *world) playerIn(build string, level int, gear []contentpack.Item, inParty bool) combat.Combatant {
	cfg := leveling.Current()
	path := strings.Split(build, "/")
	class := path[0]
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
	mods = append(mods, w.skills.Passives(path, level)...)
	mods = append(mods, skill.FromLeveling(cfg.SoloAdaptationModifiers(class, inParty))...)
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
		Skills:            w.skills.Actives(path, level),
		Level:             int16(level),
		RegenPerEncounter: int32(float64(encSeconds/2) * character.HPRegenPer2s(s.Primary.VIT)),
		RegenPerSec:       character.HPRegenPer2s(s.Primary.VIT) / 2,
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
	seeds = 150
	// encSeconds is the dungeon/raid rest base: player() fills
	// RegenPerEncounter for it and the dungeon sims rescale it.
	encSeconds = 120
	maxTurns   = 50
	hourSecs   = 3600
	// huntSeeds: an hour of back-to-back hunting holds a few hundred fights
	// (DD-281), so far fewer hours than dungeon runs give the same sample.
	huntSeeds = 20
)

// simulate is an hour of hunting the way the reveal pays it (DD-281):
// fights back to back, each kill the mob's own XP with the level-gap
// falloff.
func simulate(p combat.Combatant, m combat.Mob, level int) outcome {
	var o outcome
	m.XPReward = hunt.HuntXP(m.XPReward, level, int(m.Level))
	for s := int64(1); s <= huntSeeds; s++ {
		r := hunt.Roll(s*7919, p, m, hourSecs, maxTurns, func(seed int64) combat.RNG { return activity.NewRNG(seed) })
		if r.StoppedOnDeath {
			o.deathRate++
		}
		o.winsPerHour += float64(r.Victories)
		o.xpPerHour += float64(r.TotalXP)
	}
	o.deathRate /= huntSeeds
	o.winsPerHour /= huntSeeds
	o.xpPerHour /= huntSeeds
	if o.xpPerHour > 0 {
		o.minPerLevel = float64(character.XPToNext(level)) / o.xpPerHour * 60
	}
	return o
}

var levels = []int{1, 3, 5, 8, 10, 12, 15, 18, 20, 24, 25, 30, 35, 40, 45, 50, 55, 60, 65, 70, 75, 80, 85, 90, 95, 100,
	110, 120, 130, 140, 150, 160, 170, 180, 190, 200, 210, 220, 230, 240, 250, 260, 270, 280, 290,
	310, 330, 350, 370, 390, 410, 430, 450, 470, 490, 510, 530, 550, 570, 590,
	610, 630, 650, 670, 690, 710, 730, 750, 770, 790, 810, 830,
	850, 870, 890, 910, 930, 950, 970, 990}

// toughGap is how far above the character a mob must be dangerous.
func toughGap(int) int { return 5 }

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
	if raceDetector {
		t.Skip("pure simulation: runs in the coverage job (no -race)")
	}
	w := load(t)
	prevMin := 0.0
	for _, lv := range levels {
		m := w.mobFor(lv)
		for _, c := range builds(lv) {
			kit := simulate(w.player(c, lv, w.kit(c, lv)), m, lv)
			require.LessOrEqual(t, kit.deathRate, 0.10, "%s L%d: an hour at %s with level gear must be safe", c, lv, m.Code)
			// Common level gear levels within three times the curve's time;
			// the mid-quality pace itself is TestHuntPace's (DD-281).
			baseline := float64(character.XPToNext(lv)) / leveling.Current().Curve.BaselineXPPerMinute
			require.Greater(t, kit.minPerLevel, 0.0)
			require.LessOrEqual(t, kit.minPerLevel, baseline*3, "%s L%d: minutes per level vs baseline", c, lv)
			// Out-levelled mobs must be dangerous in the starting content.
			// Past it, sustain skills let tanks and healers outlast stronger
			// mobs, so out-levelling is made unprofitable instead: XP pays per
			// minute and loot is capped at the character's level (DD-255).
			if tm := w.toughFor(lv, toughGap(lv)); lv < 25 && tm != nil {
				// Camping (DD-281) can make them survivable; then they must pay
				// less an hour than the level's own mob.
				up := simulate(w.player(c, lv, w.kit(c, lv)), *tm, lv)
				require.True(t, up.deathRate >= 0.25 || up.xpPerHour < kit.xpPerHour,
					"%s L%d vs %s: out-levelled mobs must be dangerous or slower (death %.2f, XP/h %.0f vs %.0f)",
					c, lv, tm.Code, up.deathRate, up.xpPerHour, kit.xpPerHour)
			}
			prevMin = max(prevMin, kit.minPerLevel)
		}
	}
	require.Greater(t, prevMin, 0.0)

	for _, d := range w.dungeons {
		var lo, hi float64 = 1, 0
		for _, c := range builds(int(d.MinLevel)) {
			at := w.clearRate(c, int(d.MinLevel), d)
			lo, hi = min(lo, at), max(hi, at)
			require.GreaterOrEqual(t, at, 0.30, "%s %s at min level: clearable", c, d.Code)
			require.GreaterOrEqual(t, w.clearRate(c, int(d.MinLevel)+4, d), 0.90, "%s %s four levels up: reliable", c, d.Code)
		}
		require.LessOrEqual(t, hi-lo, 0.40, "%s: class spread at min level", d.Code)
	}
}
