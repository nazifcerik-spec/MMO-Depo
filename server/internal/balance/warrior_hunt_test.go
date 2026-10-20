package balance

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/otter/server/internal/activity"
	"github.com/otter/server/internal/combat"
	"github.com/otter/server/internal/domain/character"
	"github.com/otter/server/internal/hunt"
	"github.com/otter/server/internal/leveling"
)

// huntStats is an hour-long AFK hunt measured over many seeds: how many
// combat turns a kill takes, how often the hunt ends in death, and what
// the hour pays.
type huntStats struct {
	turnsPerKill float64 // average turns of a won fight (the player acts once a turn)
	hpLossPct    float64 // average share of max HP a won fight costs
	deathRate    float64
	winsPerHour  float64
	xpPerHour    float64
	minPerLevel  float64
}

func measureHunt(p combat.Combatant, m combat.Mob, level, n int) huntStats {
	var s huntStats
	m.XPReward = leveling.Current().ScaleMobXP(hunt.AFKXPPerVictory(encSeconds), level, int(m.Level))
	var fights, turns, lost float64
	for seed := int64(1); seed <= int64(n); seed++ {
		r := hunt.Roll(seed*7919, p, m, encPerHour, maxTurns, func(sd int64) combat.RNG { return activity.NewRNG(sd) })
		if r.StoppedOnDeath {
			s.deathRate++
		}
		s.winsPerHour += float64(r.Victories)
		s.xpPerHour += float64(r.TotalXP)
		hp := float64(p.HPMax)
		for _, e := range r.Encounters {
			if e.Outcome == combat.OutcomeVictory {
				fights++
				turns += float64(e.Duration)
				lost += math.Max(0, hp-float64(e.PlayerHPAft)) / float64(p.HPMax)
			}
			hp = math.Min(float64(p.HPMax), float64(e.PlayerHPAft+p.RegenPerEncounter))
		}
	}
	s.deathRate /= float64(n)
	s.winsPerHour /= float64(n)
	s.xpPerHour /= float64(n)
	if fights > 0 {
		s.turnsPerKill, s.hpLossPct = turns/fights, lost/fights
	}
	if s.xpPerHour > 0 {
		s.minPerLevel = float64(character.XPToNext(level)) / s.xpPerHour * 60
	}
	return s
}

// TestWarriorHuntReport writes reports/warrior_hunt_report.{md,csv}: a
// warrior with the level's gear at levels 1-50 against every mob from ten
// levels below to five above, and an attack sweep at levels 10, 25 and 50.
// Heavy; runs only with WARRIOR_REPORT=1.
func TestWarriorHuntReport(t *testing.T) {
	if os.Getenv("WARRIOR_REPORT") == "" {
		t.Skip("set WARRIOR_REPORT=1 to write reports/warrior_hunt_report.*")
	}
	w := load(t)
	const n = 100
	var md, csv strings.Builder
	csv.WriteString("level,attack,defense,hp_max,crit_pct,hit_pct,mob,mob_level,mob_hp,mob_attack,mob_defense,turns_per_kill,hp_loss_pct,death_pct,kills_per_hour,xp_per_hour,minutes_per_level\n")
	fmt.Fprintf(&md, "# Savaşçı av istatistiği — Sv 1–50\n\n")
	fmt.Fprintf(&md, "`WARRIOR_REPORT=1 go test ./internal/balance -run TestWarriorHuntReport`; üretim savaş kodu (`hunt.Roll`), %d tohum, saatlik AFK av (120 sn'de bir karşılaşma, saatte en çok 30), karakter seviyesinin ekipmanı (her yuvada ulaşılan en yüksek kademe), seviyeye göre dağıtılmış nitelik puanları, sınıf pasifleri ve aktifleri. Canavar: Sv−10 … Sv+5 arası her av canavarı.\n\n", n)
	md.WriteString("Sütunlar: **Tur/öldürme** = kazanılan savaşın ortalama turu (oyuncu her tur bir kez vurur); **Can kaybı** = kazanılan savaşta azami canın ortalama yüzdesi; **Ölüm** = bir saatlik avın ölümle bitme oranı; **Öldürme/sa**, **XP/sa**; **Dk/sv** = bu canavarla bir sonraki seviyeye dakika.\n\n")

	// Per level: the strongest mob hunted safely (death ≤ 10%) and the first
	// that kills the warrior in most hours (death ≥ 50%).
	type summary struct {
		lv, atk          int
		safe, deadly     string
		safeLv, deadlyLv int
		turns, hpLoss    float64
	}
	var sums []summary
	for lv := 1; lv <= 50; lv++ {
		p := w.player("warrior", lv, w.kit("warrior", lv))
		fmt.Fprintf(&md, "## Sv %d — saldırı %d · savunma %d · can %d · kritik %%%d · isabet %%%d\n\n", lv, p.Attack, p.Defense, p.HPMax, p.CritPct, p.HitPct)
		md.WriteString("| Canavar | Sv | Can | Sal | Sav | Tur/öldürme | Can kaybı | Ölüm | Öldürme/sa | XP/sa | Dk/sv |\n|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|\n")
		sm := summary{lv: lv, atk: int(p.Attack)}
		for _, m := range w.mobs {
			if int(m.Level) < lv-10 || int(m.Level) > lv+5 {
				continue
			}
			s := measureHunt(p, m, lv, n)
			dk := "—"
			if s.minPerLevel > 0 {
				dk = fmt.Sprintf("%.0f", s.minPerLevel)
			}
			if s.deathRate <= 0.10 && int(m.Level) >= sm.safeLv {
				sm.safe, sm.safeLv, sm.turns, sm.hpLoss = m.Code, int(m.Level), s.turnsPerKill, s.hpLossPct
			}
			if s.deathRate >= 0.5 && sm.deadly == "" {
				sm.deadly, sm.deadlyLv = m.Code, int(m.Level)
			}
			fmt.Fprintf(&md, "| %s | %d | %d | %d | %d | %.1f | %%%.0f | %%%.0f | %.1f | %.0f | %s |\n",
				m.Code, m.Level, m.HPMax, m.Attack, m.Defense, s.turnsPerKill, 100*s.hpLossPct, 100*s.deathRate, s.winsPerHour, s.xpPerHour, dk)
			fmt.Fprintf(&csv, "%d,%d,%d,%d,%d,%d,%s,%d,%d,%d,%d,%.2f,%.1f,%.1f,%.2f,%.0f,%.1f\n", lv, p.Attack, p.Defense, p.HPMax, p.CritPct, p.HitPct,
				m.Code, m.Level, m.HPMax, m.Attack, m.Defense, s.turnsPerKill, 100*s.hpLossPct, 100*s.deathRate, s.winsPerHour, s.xpPerHour, s.minPerLevel)
		}
		md.WriteString("\n")
		sums = append(sums, sm)
	}

	var head strings.Builder
	head.WriteString("## Özet\n\nGüvenli = bir saatlik avın ölümle bitme oranı ≤ %10; ölümcül = ≥ %50. Güvenli her canavar aynı saatlik XP'yi verir; saldırı gücü savaşın süresini ve can kaybını belirler.\n\n")
	head.WriteString("| Sv | Saldırı | En güçlü güvenli canavar | Tur/öldürme | Can kaybı | İlk ölümcül canavar |\n|---:|---:|---|---:|---:|---|\n")
	for _, s := range sums {
		deadly := "—"
		if s.deadly != "" {
			deadly = fmt.Sprintf("%s (Sv %d)", s.deadly, s.deadlyLv)
		}
		fmt.Fprintf(&head, "| %d | %d | %s (Sv %d) | %.1f | %%%.0f | %s |\n", s.lv, s.atk, s.safe, s.safeLv, s.turns, 100*s.hpLoss, deadly)
	}
	head.WriteString("\n## Saldırı gücü taraması\n\nAynı karakter, yalnız saldırı değeri çarpılarak; canavar seviyeye uygun en güçlü av canavarı.\n\n")
	head.WriteString("| Sv | Canavar | Çarpan | Saldırı | Tur/öldürme | Can kaybı | Ölüm | XP/sa | Dk/sv |\n|---:|---|---:|---:|---:|---:|---:|---:|---:|\n")
	for _, lv := range []int{10, 25, 50} {
		base := w.player("warrior", lv, w.kit("warrior", lv))
		m := w.mobFor(lv)
		for _, f := range []float64{0.6, 0.8, 1.0, 1.2, 1.4, 1.6} {
			p := base.Clone()
			p.Attack = int32(math.Round(float64(base.Attack) * f))
			s := measureHunt(p, m, lv, n)
			fmt.Fprintf(&head, "| %d | %s | ×%.1f | %d | %.1f | %%%.0f | %%%.0f | %.0f | %.0f |\n",
				lv, m.Code, f, p.Attack, s.turnsPerKill, 100*s.hpLossPct, 100*s.deathRate, s.xpPerHour, s.minPerLevel)
			fmt.Fprintf(&csv, "sweep_%d_x%.1f,%d,%d,%d,%d,%d,%s,%d,%d,%d,%d,%.2f,%.1f,%.1f,%.2f,%.0f,%.1f\n", lv, f, p.Attack, p.Defense, p.HPMax, p.CritPct, p.HitPct,
				m.Code, m.Level, m.HPMax, m.Attack, m.Defense, s.turnsPerKill, 100*s.hpLossPct, 100*s.deathRate, s.winsPerHour, s.xpPerHour, s.minPerLevel)
		}
	}
	head.WriteString("\n")

	out := md.String()
	i := strings.Index(out, "## Sv 1 ")
	out = out[:i] + head.String() + out[i:]
	dir := filepath.Join("..", "..", "..", "reports")
	require.NoError(t, os.MkdirAll(dir, 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(dir, "warrior_hunt_report.md"), []byte(out), 0o644))
	require.NoError(t, os.WriteFile(filepath.Join(dir, "warrior_hunt_report.csv"), []byte(csv.String()), 0o644))
}
