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
)

// huntStats is an hour of hunting measured over many seeds (DD-281): how
// long a kill takes, how much of the hour goes to camping, how often the
// hunt ends in death, and what the hour pays.
type huntStats struct {
	turnsPerKill float64 // turns of a won fight (the player acts once a turn)
	secsPerKill  float64 // hunting seconds per kill: fights plus camping
	restPct      float64 // share of the hunting time spent camping
	deathRate    float64 // share of hours that end in a defeat
	winsPerHour  float64
	xpPerHour    float64
	minPerLevel  float64
}

func measureHunt(p combat.Combatant, m combat.Mob, level, n int) huntStats {
	var s huntStats
	m.XPReward = hunt.HuntXP(m.XPReward, level, int(m.Level))
	var wins, turns, rest, elapsed float64
	for seed := int64(1); seed <= int64(n); seed++ {
		r := hunt.Roll(seed*7919, p, m, hourSecs, maxTurns, func(sd int64) combat.RNG { return activity.NewRNG(sd) })
		if r.StoppedOnDeath {
			s.deathRate++
		}
		wins += float64(r.Victories)
		turns += float64(r.FightTurns)
		rest += float64(r.RestSecs)
		elapsed += float64(r.ElapsedSecs)
		s.xpPerHour += float64(r.TotalXP)
	}
	s.deathRate /= float64(n)
	s.winsPerHour = wins / float64(n)
	s.xpPerHour /= float64(n)
	if wins > 0 {
		s.turnsPerKill, s.secsPerKill = turns/wins, elapsed/wins
	}
	if elapsed > 0 {
		s.restPct = rest / elapsed
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
	const n = 60
	var md, csv strings.Builder
	csv.WriteString("level,attack,defense,hp_max,crit_pct,hit_pct,mob,mob_level,mob_hp,mob_attack,mob_defense,mob_xp,turns_per_kill,secs_per_kill,rest_pct,death_pct,kills_per_hour,xp_per_hour,minutes_per_level\n")
	fmt.Fprintf(&md, "# Savaşçı av istatistiği — Sv 1–50\n\n")
	fmt.Fprintf(&md, "`WARRIOR_REPORT=1 go test ./internal/balance -run TestWarriorHuntReport`; üretim av kodu (`hunt.Roll`, DD-281), %d tohum, bir saatlik av: savaşlar peş peşe, tur %d sn, can %%%d'nin altına düşünce kamp (saniyede azami canın %%%d'i). Karakter: seviyenin ekipmanı (her yuvada ulaşılan en yüksek kademe, sıradan), dağıtılmış nitelik puanları, sınıf pasifleri ve aktifleri. Canavar: Sv−10 … Sv+5 arası her av canavarı; XP canavarın kendi değeri, seviye farkına göre düşer.\n\n",
		n, hunt.TurnSeconds, hunt.RestBelowPct, hunt.RestHealPctPerSec)
	md.WriteString("Sütunlar: **XP** = canavarın öldürme başı XP'si (seviye farkı uygulanmış); **Tur** = kazanılan savaşın ortalama turu; **Sn/öldürme** = savaş + kamp dahil öldürme başı av süresi; **Kamp** = sürenin kampta geçen payı; **Ölüm** = bir saatlik avın ölümle bitme oranı; **Öldürme/sa**, **XP/sa**; **Dk/sv** = bu canavarla bir sonraki seviyeye dakika.\n\n")

	type summary struct {
		lv, atk                 int
		best, deadly            string
		bestLv, deadlyLv        int
		bestXPH, bestSec, bestM float64
	}
	var sums []summary
	for lv := 1; lv <= 50; lv++ {
		p := w.player("warrior", lv, w.kit("warrior", lv))
		fmt.Fprintf(&md, "## Sv %d — saldırı %d · savunma %d · can %d · kritik %%%d · isabet %%%d\n\n", lv, p.Attack, p.Defense, p.HPMax, p.CritPct, p.HitPct)
		md.WriteString("| Canavar | Sv | Can | Sal | Sav | XP | Tur | Sn/öldürme | Kamp | Ölüm | Öldürme/sa | XP/sa | Dk/sv |\n|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|\n")
		sm := summary{lv: lv, atk: int(p.Attack)}
		for _, m := range w.mobs {
			if int(m.Level) < lv-10 || int(m.Level) > lv+5 {
				continue
			}
			s := measureHunt(p, m, lv, n)
			xp := hunt.HuntXP(m.XPReward, lv, int(m.Level))
			dk := "—"
			if s.minPerLevel > 0 {
				dk = fmt.Sprintf("%.0f", s.minPerLevel)
			}
			if s.deathRate <= 0.10 && s.xpPerHour > sm.bestXPH {
				sm.best, sm.bestLv, sm.bestXPH, sm.bestSec, sm.bestM = m.Code, int(m.Level), s.xpPerHour, s.secsPerKill, s.minPerLevel
			}
			if s.deathRate >= 0.5 && sm.deadly == "" {
				sm.deadly, sm.deadlyLv = m.Code, int(m.Level)
			}
			fmt.Fprintf(&md, "| %s | %d | %d | %d | %d | %d | %.1f | %.0f | %%%.0f | %%%.0f | %.0f | %.0f | %s |\n",
				m.Code, m.Level, m.HPMax, m.Attack, m.Defense, xp, s.turnsPerKill, s.secsPerKill, 100*s.restPct, 100*s.deathRate, s.winsPerHour, s.xpPerHour, dk)
			fmt.Fprintf(&csv, "%d,%d,%d,%d,%d,%d,%s,%d,%d,%d,%d,%d,%.2f,%.1f,%.1f,%.1f,%.1f,%.0f,%.1f\n", lv, p.Attack, p.Defense, p.HPMax, p.CritPct, p.HitPct,
				m.Code, m.Level, m.HPMax, m.Attack, m.Defense, xp, s.turnsPerKill, s.secsPerKill, 100*s.restPct, 100*s.deathRate, s.winsPerHour, s.xpPerHour, s.minPerLevel)
		}
		md.WriteString("\n")
		sums = append(sums, sm)
	}

	var head strings.Builder
	head.WriteString("## Özet\n\nEn verimli güvenli canavar = ölümle biten saat ≤ %10 olanlar içinde saatte en çok XP veren. Ölümcül = ≥ %50.\n\n")
	head.WriteString("| Sv | Saldırı | En verimli güvenli canavar | Sn/öldürme | XP/sa | Dk/sv | İlk ölümcül canavar |\n|---:|---:|---|---:|---:|---:|---|\n")
	for _, s := range sums {
		deadly := "—"
		if s.deadly != "" {
			deadly = fmt.Sprintf("%s (Sv %d)", s.deadly, s.deadlyLv)
		}
		fmt.Fprintf(&head, "| %d | %d | %s (Sv %d) | %.0f | %.0f | %.0f | %s |\n", s.lv, s.atk, s.best, s.bestLv, s.bestSec, s.bestXPH, s.bestM, deadly)
	}
	head.WriteString("\n## Saldırı gücü taraması\n\nAynı karakter, yalnız saldırı değeri çarpılarak; canavar seviyeye uygun en güçlü av canavarı.\n\n")
	head.WriteString("| Sv | Canavar | Çarpan | Saldırı | Tur | Sn/öldürme | Kamp | Ölüm | Öldürme/sa | XP/sa | Dk/sv |\n|---:|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|\n")
	for _, lv := range []int{10, 25, 50} {
		base := w.player("warrior", lv, w.kit("warrior", lv))
		m := w.mobFor(lv)
		for _, f := range []float64{0.6, 0.8, 1.0, 1.2, 1.4, 1.6} {
			p := base.Clone()
			p.Attack = int32(math.Round(float64(base.Attack) * f))
			s := measureHunt(p, m, lv, n)
			fmt.Fprintf(&head, "| %d | %s | ×%.1f | %d | %.1f | %.0f | %%%.0f | %%%.0f | %.0f | %.0f | %.0f |\n",
				lv, m.Code, f, p.Attack, s.turnsPerKill, s.secsPerKill, 100*s.restPct, 100*s.deathRate, s.winsPerHour, s.xpPerHour, s.minPerLevel)
			fmt.Fprintf(&csv, "sweep_%d_x%.1f,%d,%d,%d,%d,%d,%s,%d,%d,%d,%d,%d,%.2f,%.1f,%.1f,%.1f,%.1f,%.0f,%.1f\n", lv, f, p.Attack, p.Defense, p.HPMax, p.CritPct, p.HitPct,
				m.Code, m.Level, m.HPMax, m.Attack, m.Defense, hunt.HuntXP(m.XPReward, lv, int(m.Level)), s.turnsPerKill, s.secsPerKill, 100*s.restPct, 100*s.deathRate, s.winsPerHour, s.xpPerHour, s.minPerLevel)
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
