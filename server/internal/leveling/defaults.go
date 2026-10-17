package leveling

// defaultConfig is used until content/progression.yaml is installed: the
// curve and point rules of DD-250 without class trees (no promotion gates).
var defaultConfig = func() *Config {
	c := &Config{
		LevelCap:             1000,
		Curve:                Curve{Kind: "formula", BaselineXPPerMinute: 1000, A: 3.75006670, B: 0.00121846427, C: 1.69337967},
		XPMultiplier:         1,
		AFKCombatXPPerMinute: 1000,
		AttributePoints:      AttributeRule{Every: 2, Amount: 1, Milestones: map[int]int{100: 5, 300: 5, 600: 5, 850: 5, 1000: 5}},
		TalentPoints:         TalentRule{From: 100, Every: 20, Until: 980, Amount: 1},
		TalentTreePoints:     25,
		Milestones:           Milestones{Awakening: 600, Ultimate: 700, Capstone: 850, Mastery: 1000},
		Scaling:              []Bracket{{100, 1}, {250, 0.8}, {400, 0.55}, {0, 0.35}},
		Caps:                 map[string]float64{"crit_pct": 60, "evasion_pct": 40, "block_pct": 60, "dmg_reduction_pct": 75},
	}
	if err := c.Prepare(); err != nil {
		panic(err)
	}
	return c
}()
