package combat

// Fight reports (DD-308): why an AFK run ended and how the two sides
// measured up over every fight in it, for the result screens.

// Why a hunt, dungeon or raid ended.
const (
	EndTimeUp  = "time_up"         // the planned time ran out
	EndEarly   = "collected_early" // collected before the planned end
	EndFled    = "fled"            // broke off under FleeBelowPct HP
	EndDefeat  = "defeat"          // the character (or the whole party) fell
	EndTimeout = "timeout"         // a fight dragged past its turn limit
	EndCleared = "cleared"         // every wave and the boss fell
)

// Stats are a combatant's numbers as it fought.
type Stats struct {
	Level      int16 `json:"level"`
	HPMax      int32 `json:"hp_max"`
	Attack     int32 `json:"attack"`
	Defense    int32 `json:"defense"`
	HitPct     int32 `json:"hit_pct"`
	EvasionPct int32 `json:"evasion_pct"`
	CritPct    int32 `json:"crit_pct"`
}

func StatsOf(c Combatant) Stats {
	return Stats{Level: c.Level, HPMax: c.HPMax, Attack: c.Attack, Defense: c.Defense,
		HitPct: c.HitPct, EvasionPct: c.EvasionPct, CritPct: c.CritPct}
}

// Side sums one side's swings.
type Side struct {
	Swings int   `json:"swings"`
	Hits   int   `json:"hits"`
	Crits  int   `json:"crits"`
	Damage int64 `json:"damage"`
	MaxHit int32 `json:"max_hit"`
	Healed int64 `json:"healed,omitempty"`
}

// add counts a turn of this side; atFoe: aimed at the other side (a heal
// on an ally is no swing).
func (s *Side) add(t TurnRecord, atFoe bool) {
	s.Healed += int64(t.Heal) + int64(t.TickHeal)
	if atFoe {
		s.Damage += int64(t.Tick) // damage over time lands without a swing
	}
	if t.Stunned || !atFoe || (t.Tick > 0 && t.Damage == 0 && t.Skill == "") {
		return
	}
	s.Swings++
	if !t.Hit {
		return
	}
	s.Hits++
	if t.Crit {
		s.Crits++
	}
	s.Damage += int64(t.Damage)
	s.MaxHit = max(s.MaxHit, t.Damage)
}

// Report is a run's fight report. You and Enemy are the last fight's
// sides: the hunted mob, or the dungeon or raid enemy where the run ended
// (the boss of a cleared run). In a raid, Yours sums the whole party.
type Report struct {
	End     string `json:"end"`
	EndMob  string `json:"end_mob,omitempty"` // who ended it (defeat, timeout)
	You     Stats  `json:"you"`
	Enemy   Stats  `json:"enemy"`
	EnemyID string `json:"enemy_code"`
	// Hit chances as TryHit computes them, you on the enemy and back.
	HitChance      int32 `json:"hit_chance"`
	EnemyHitChance int32 `json:"enemy_hit_chance"`
	Fights         int   `json:"fights"`
	Yours          Side  `json:"yours"`
	Theirs         Side  `json:"theirs"`
	// LowestHPPct is the lowest a player's HP fell inside a fight.
	LowestHPPct int `json:"lowest_hp_pct"`
}

// NewReport starts a report with nobody hurt yet.
func NewReport() *Report { return &Report{LowestHPPct: 100} }

// Fight adds one fight: you as fought, the party (you included; nil for a
// solo fight) to read HP shares by name, the enemy as fought and the log.
func (r *Report) Fight(you Combatant, party []Combatant, enemy Combatant, turns []TurnRecord) {
	r.You, r.Enemy, r.EnemyID = StatsOf(you), StatsOf(enemy), enemy.Name
	r.HitChance = HitChance(you.HitPct, enemy.EvasionPct)
	r.EnemyHitChance = HitChance(enemy.HitPct, you.EvasionPct)
	r.Fights++
	hpMax := map[string]int32{you.Name: you.HPMax}
	for _, p := range party {
		hpMax[p.Name] = p.HPMax
	}
	for _, t := range turns {
		if t.Actor == enemy.Name {
			r.Theirs.add(t, t.Target != enemy.Name)
			if m := hpMax[t.Target]; m > 0 && t.Hit {
				r.LowestHPPct = min(r.LowestHPPct, int(int64(t.TargetHPAft)*100/int64(m)))
			}
			continue
		}
		r.Yours.add(t, t.Target == enemy.Name)
	}
}

// HitChance is the chance a swing lands, as TryHit clamps it.
func HitChance(hit, evasion int32) int32 {
	return min(max(hit-evasion, MinHitPct), MaxHitPct)
}

// For is the report seen from party member c: its own stats and hit
// chances against the same enemy and fight sums.
func (r Report) For(c Combatant) *Report {
	r.You = StatsOf(c)
	r.HitChance = HitChance(c.HitPct, r.Enemy.EvasionPct)
	r.EnemyHitChance = HitChance(r.Enemy.HitPct, c.EvasionPct)
	return &r
}
