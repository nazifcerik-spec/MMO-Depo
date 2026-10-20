package balance

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/otter/server/internal/activity"
	"github.com/otter/server/internal/combat"
	"github.com/otter/server/internal/domain/character"
	"github.com/otter/server/internal/hunt"
	"github.com/otter/server/internal/leveling"
	"github.com/otter/server/internal/skill"
)

// Specialization comparison (DD-252): every final specialization at high
// levels against synthetic mobs calibrated to the field's average, so the
// numbers compare specializations with each other (the shipped mob content
// ends at level 25).

var specLevels = []int{700, 1000}

const specSeeds = 60

// specRole groups a specialization for the comparisons by its first role tag.
func specRole(k *leveling.Class) string {
	switch k.RoleTags[0] {
	case "tank":
		return "tank"
	case "healer", "buff_support", "barrier_support", "group_support", "debuff_support":
		return "support"
	}
	return "dps"
}

var roleTargets = map[string]map[string]bool{
	"tank":    {"defense": true, "hp_max": true, "block_pct": true, "dmg_reduction_pct": true},
	"support": {"hp_max": true, "lifesteal_pct": true, "mp_max": true, "dmg_reduction_pct": true},
	"dps":     {"attack": true, "crit_pct": true, "crit_dmg_pct": true, "execute_pct": true, "hit_pct": true, "missing_hp_dmg_pct": true},
}

// talentMods spends the level's talent points the way the role would: the
// best-fitting tree to 25 (+ capstone from its level), the rest in the next.
func talentMods(cfg *leveling.Config, spec, role string, level int) []leveling.Modifier {
	trees := cfg.TreesFor(spec)
	score := func(t leveling.Tree) int {
		n := 0
		for _, tl := range t.Talents {
			for _, m := range tl.Modifiers {
				if roleTargets[role][m.Target] {
					n++
				}
			}
		}
		return n
	}
	sort.SliceStable(trees, func(i, j int) bool { return score(trees[i]) > score(trees[j]) })
	points := cfg.TalentPointsEarned(level)
	var out []leveling.Modifier
	for ti, t := range trees {
		inTree := 0
		for _, tl := range t.Talents {
			if tl.Capstone {
				if ti == 0 && inTree >= cfg.TalentTreePoints && level >= cfg.Milestones.Capstone && points > 0 {
					out = append(out, tl.Modifiers...)
					points--
				}
				continue
			}
			for r := 0; r < tl.MaxRank && points > 0; r++ {
				out = append(out, tl.Modifiers...)
				points--
				inTree++
			}
		}
	}
	return out
}

// specPlayer is the sim player of a class path at a level: base class
// formula + allocation, the best shipped gear, path passives, talents,
// mastery, Solo Adaptation; path actives and the class resource.
func (w *world) specPlayer(spec string, level int) combat.Combatant {
	cfg := leveling.Current()
	path := cfg.Path(spec)
	class := path[0]
	k, _ := cfg.Class(spec)
	ch := &character.Character{ClassCode: class, Level: int16(level), Alloc: simAlloc(class, level)}
	raw := character.RawAttributes(ch)
	var mods []skill.Modifier
	for _, it := range w.kit(class, 25) {
		for _, sb := range it.StatBonuses {
			mods = append(mods, skill.Modifier{Kind: skill.ModifierKind(sb.Kind), Target: skill.Target(sb.Target), Amount: sb.Amount})
		}
	}
	mods = append(mods, w.skills.Passives(path, level)...)
	var lm []leveling.Modifier
	for _, code := range path {
		n, _ := cfg.Class(code)
		lm = append(lm, n.StatModifiers...)
	}
	lm = append(lm, talentMods(cfg, spec, specRole(k), level)...)
	if k.Mastery != nil && level >= cfg.Milestones.Mastery {
		lm = append(lm, k.Mastery.Modifiers...)
	}
	lm = append(lm, cfg.SoloAdaptationModifiers(spec, false)...)
	mods = append(mods, skill.FromLeveling(lm)...)
	base, _ := cfg.Class(class)
	s := skill.ComputeSheetWith(skill.Primary{STR: raw.Strength, AGI: raw.Dexterity, INT: raw.Intelligence,
		VIT: raw.Vitality, CHA: raw.Wisdom, Luck: raw.Luck}, skill.Target(character.MainStat(class)), mods,
		&skill.Formula{Derived: *base.Derived, Level: level})
	d := s.Derived
	p := combat.Combatant{Name: "player", HPMax: int32(s.Resources.HPMax), HPCurrent: int32(s.Resources.HPMax),
		Attack: int32(d.Attack), Defense: int32(d.Defense), CritPct: int32(d.CritPct), HitPct: int32(d.HitPct),
		EvasionPct: int32(d.EvasionPct), MPMax: int32(s.Resources.MPMax), Skills: w.skills.Actives(path, level),
		Level: int16(level), BlockPct: int32(d.BlockPct), CritDmgPct: int32(d.CritDmgPct), ExecutePct: int32(d.ExecutePct),
		MissingHPDmgPct: int32(d.MissingHPDmgPct), LifestealPct: int32(d.LifestealPct), DmgReductionPct: int32(d.DmgReductionPct)}
	combat.ApplyResource(&p, class)
	return p
}

// specMetrics are the simulated numbers of one specialization at a level.
type specMetrics struct {
	KillTurns   float64 `json:"kill_turns"`    // mean turns to kill a field mob (lower = faster)
	WinRate     float64 `json:"win_rate"`      // field mob, one encounter
	Survival    float64 `json:"survival"`      // turns survived against an unkillable hard-hitting mob
	BossWinRate float64 `json:"boss_win_rate"` // boss (6x HP, 1.1x attack), one encounter
	BossTurns   float64 `json:"boss_turns"`    // mean turns of won boss fights
	AFKSafety   float64 `json:"afk_safety"`    // share of 30-encounter chains finished alive
}

type field struct{ mob, hard, boss combat.Mob }

// calibrate builds the synthetic mobs of a level from the mean player.
func calibrate(players []combat.Combatant, level int) field {
	var atk, def, hp float64
	for _, p := range players {
		atk += float64(p.Attack)
		def += float64(p.Defense)
		hp += float64(p.HPMax)
	}
	n := float64(len(players))
	atk, def, hp = atk/n, def/n, hp/n
	m := combat.Mob{Code: "field", Level: int16(level), Defense: int32(0.4 * atk), HitPct: 90, CritPct: 5, EvasionPct: 5}
	m.HPMax = int32(8 * 0.8 * atk)
	m.Attack = int32(0.07*hp + 0.5*def)
	hard, boss := m, m
	hard.Code, hard.Attack = "hard", int32(float64(m.Attack)*1.4)
	boss.Code, boss.HPMax, boss.Attack = "boss", m.HPMax*6, int32(float64(m.Attack)*1.1)
	return field{m, hard, boss}
}

func rng(seed int64) combat.RNG { return activity.NewRNG(seed) }

func measure(p combat.Combatant, f field) specMetrics {
	var m specMetrics
	wins, bossWins := 0, 0
	turns, bossTurns := 0, 0
	for s := int64(1); s <= specSeeds; s++ {
		sess := combat.Session{Player: p.Clone(), Mob: f.mob.AsCombatant(), MaxTurns: 60}
		if r, err := sess.Resolve(rng, s*31); err == nil && r.Outcome == combat.OutcomeVictory {
			wins++
			turns += r.Duration
		}
		boss := combat.Session{Player: p.Clone(), Mob: f.boss.AsCombatant(), MaxTurns: 200}
		if r, err := boss.Resolve(rng, s*37); err == nil && r.Outcome == combat.OutcomeVictory {
			bossWins++
			bossTurns += r.Duration
		}
		wall := f.hard.AsCombatant()
		wall.HPMax, wall.HPCurrent = math.MaxInt32/2, math.MaxInt32/2
		if r, err := (&combat.Session{Player: p.Clone(), Mob: wall, MaxTurns: 300}).Resolve(rng, s*41); err == nil {
			m.Survival += float64(r.Duration)
		}
		chain := p.Clone()
		chain.RegenPerEncounter = p.HPMax / 20
		if r := hunt.Roll(s*43, chain, f.mob, hourSecs, 60, rng); !r.StoppedOnDeath {
			m.AFKSafety++
		}
	}
	m.WinRate = float64(wins) / specSeeds
	m.KillTurns = float64(turns) / math.Max(1, float64(wins))
	m.BossWinRate = float64(bossWins) / specSeeds
	m.BossTurns = float64(bossTurns) / math.Max(1, float64(bossWins))
	m.Survival /= specSeeds
	m.AFKSafety /= specSeeds
	return m
}

type specRow struct {
	Code      string                 `json:"code"`
	Class     string                 `json:"class"`
	Promotion string                 `json:"promotion"`
	Role      string                 `json:"role"`
	RoleTags  []string               `json:"role_tags"`
	Support   float64                `json:"support_power"`
	Profile   leveling.AFKProfile    `json:"afk_profile"`
	StrongIn  []string               `json:"strong_in"`
	WeakIn    []string               `json:"weak_in"`
	Sim       map[string]specMetrics `json:"sim"`
	// Scores (0..1 within the field at the top level): simulated kill
	// speed, survival, boss and AFK safety; design AOE and group value.
	Scores map[string]float64 `json:"scores"`
}

func finalSpecs(cfg *leveling.Config) []*leveling.Class {
	var out []*leveling.Class
	for i := range cfg.Classes {
		k, _ := cfg.Class(cfg.Classes[i].Code)
		if k.Tier() == 2 {
			out = append(out, k)
		}
	}
	return out
}

func buildSpecReport(t *testing.T) ([]specRow, map[string]any) {
	w := load(t)
	cfg := leveling.Current()
	specs := finalSpecs(cfg)
	require.Len(t, specs, 40)
	rows := make([]specRow, len(specs))
	for i, k := range specs {
		path := cfg.Path(k.Code)
		rows[i] = specRow{Code: k.Code, Class: path[0], Promotion: path[1], Role: specRole(k), RoleTags: k.RoleTags,
			Support: cfg.SupportPower(k.Code), Profile: *k.AFKProfile, StrongIn: k.StrongIn, WeakIn: k.WeakIn,
			Sim: map[string]specMetrics{}}
	}
	for _, lv := range specLevels {
		players := make([]combat.Combatant, len(specs))
		for i, k := range specs {
			players[i] = w.specPlayer(k.Code, lv)
		}
		f := calibrate(players, lv)
		for i := range specs {
			rows[i].Sim[fmt.Sprint(lv)] = measure(players[i], f)
		}
	}
	top := fmt.Sprint(specLevels[len(specLevels)-1])
	metric := map[string]func(r specRow) float64{
		"kill_speed":  func(r specRow) float64 { return 1 / r.Sim[top].KillTurns },
		"survival":    func(r specRow) float64 { return r.Sim[top].Survival },
		"boss":        func(r specRow) float64 { return r.Sim[top].BossWinRate / math.Max(1, r.Sim[top].BossTurns) },
		"afk_safety":  func(r specRow) float64 { return r.Sim[top].AFKSafety },
		"aoe":         func(r specRow) float64 { return float64(r.Profile.AOEValue) },
		"group_value": func(r specRow) float64 { return float64(r.Profile.GroupValue) + r.Support },
	}
	for name, f := range metric {
		lo, hi := math.Inf(1), math.Inf(-1)
		for _, r := range rows {
			lo, hi = math.Min(lo, f(r)), math.Max(hi, f(r))
		}
		for i := range rows {
			if rows[i].Scores == nil {
				rows[i].Scores = map[string]float64{}
			}
			rows[i].Scores[name] = math.Round((f(rows[i])-lo)/math.Max(hi-lo, 1e-9)*1000) / 1000
		}
	}
	mean := func(role, name string) float64 {
		s, n := 0.0, 0
		for _, r := range rows {
			if r.Role == role {
				s += r.Scores[name]
				n++
			}
		}
		return s / float64(n)
	}
	checks := map[string]any{}
	for _, role := range []string{"dps", "tank", "support"} {
		for name := range metric {
			checks[role+"_"+name] = math.Round(mean(role, name)*1000) / 1000
		}
	}
	return rows, checks
}

// TestSpecializationBalance pins the DD-252 principles: pure DPS kill
// fastest, tanks survive longest, supports bring the most group value, and
// no specialization tops every category. BALANCE_REPORT=1 also writes
// reports/class_balance_report.{json,md}.
func TestSpecializationBalance(t *testing.T) {
	if raceDetector {
		t.Skip("pure simulation: runs in the coverage job (no -race)")
	}
	rows, checks := buildSpecReport(t)
	f := func(k string) float64 { return checks[k].(float64) }
	require.Greater(t, f("dps_kill_speed"), f("tank_kill_speed"), "DPS kill faster than tanks")
	require.Greater(t, f("dps_kill_speed"), f("support_kill_speed"), "DPS kill faster than supports")
	require.Greater(t, f("tank_survival"), f("dps_survival"), "tanks survive longer than DPS")
	require.Greater(t, f("support_group_value"), f("dps_group_value"), "supports bring more group value")
	categories := []string{"kill_speed", "survival", "boss", "aoe", "group_value", "afk_safety"}
	for _, r := range rows {
		top := 0
		for _, c := range categories {
			if r.Scores[c] >= 0.75 {
				top++
			}
		}
		require.Less(t, top, len(categories), "%s must not dominate every category", r.Code)
		for _, lv := range specLevels {
			require.GreaterOrEqual(t, r.Sim[fmt.Sprint(lv)].WinRate, 0.9, "%s L%d: viable against field mobs", r.Code, lv)
		}
	}
	if os.Getenv("BALANCE_REPORT") != "" {
		writeReport(t, rows, checks)
	}
}

func writeReport(t *testing.T, rows []specRow, checks map[string]any) {
	dir := filepath.Join("..", "..", "..", "reports")
	require.NoError(t, os.MkdirAll(dir, 0o755))
	doc := map[string]any{
		"description": "Development-only class balance report (DD-252): every final specialization simulated with the " +
			"production combat engine against synthetic mobs calibrated to the field average at each level. Scores " +
			"are 0..1 within the field at level 1000; aoe and group_value come from the design profile (the engine " +
			"resolves single-target encounters).",
		"generated_by": "go test ./internal/balance -run TestSpecializationBalance (BALANCE_REPORT=1)",
		"levels":       specLevels, "seeds": specSeeds, "role_means": checks, "specializations": rows,
	}
	raw, err := json.MarshalIndent(doc, "", "  ")
	require.NoError(t, err)
	require.NoError(t, os.WriteFile(filepath.Join(dir, "class_balance_report.json"), append(raw, '\n'), 0o644))

	var b strings.Builder
	b.WriteString("# Class Balance Report\n\nDevelopment-only (DD-252). Generated by `BALANCE_REPORT=1 go test ./internal/balance -run TestSpecializationBalance`.\n\n")
	b.WriteString("Each final specialization is simulated with the production combat engine (skills, passives, talents, mastery, Solo Adaptation, class resource) against synthetic mobs calibrated to the field average. Scores are 0–1 within the field at level 1000; **aoe** and **group** come from the design profile because encounters are single-target.\n\n")
	b.WriteString("## Role means (level 1000)\n\n| role | kill speed | survival | boss | AFK safety | aoe | group |\n|---|---|---|---|---|---|---|\n")
	for _, role := range []string{"dps", "tank", "support"} {
		fmt.Fprintf(&b, "| %s | %.2f | %.2f | %.2f | %.2f | %.2f | %.2f |\n", role, checks[role+"_kill_speed"], checks[role+"_survival"],
			checks[role+"_boss"], checks[role+"_afk_safety"], checks[role+"_aoe"], checks[role+"_group_value"])
	}
	b.WriteString("\n## Specializations\n\n| class | promotion | specialization | role | L700 kill turns | L1000 kill turns | L1000 survival turns | boss win | AFK safety | kill | surv | boss | aoe | group | afk | strong in | weak in |\n|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|\n")
	sorted := slices.Clone(rows)
	sort.SliceStable(sorted, func(i, j int) bool { return sorted[i].Class < sorted[j].Class })
	for _, r := range sorted {
		a, z := r.Sim["700"], r.Sim["1000"]
		fmt.Fprintf(&b, "| %s | %s | %s | %s | %.1f | %.1f | %.1f | %.0f%% | %.0f%% | %.2f | %.2f | %.2f | %.2f | %.2f | %.2f | %s | %s |\n",
			r.Class, r.Promotion, r.Code, r.Role, a.KillTurns, z.KillTurns, z.Survival, z.BossWinRate*100, z.AFKSafety*100,
			r.Scores["kill_speed"], r.Scores["survival"], r.Scores["boss"], r.Scores["aoe"], r.Scores["group_value"], r.Scores["afk_safety"],
			strings.Join(r.StrongIn, ", "), strings.Join(r.WeakIn, ", "))
	}
	require.NoError(t, os.WriteFile(filepath.Join(dir, "class_balance_report.md"), []byte(b.String()), 0o644))
}
