package balance

import (
	"fmt"
	"os"
	"testing"
)

// TestMobCurve prints the field the content's mob curve is calibrated to
// (MOBCURVE=1): the class-average player at each level with the best gear
// that level can use, and the base mob that keeps the level-25 ratios.
func TestMobCurve(t *testing.T) {
	if os.Getenv("MOBCURVE") == "" {
		t.Skip()
	}
	w := load(t)
	mean := func(level int) (atk, def, hp float64) {
		for _, c := range classes {
			p := w.player(c, level, w.kit(c, level))
			atk += float64(p.Attack)
			def += float64(p.Defense)
			hp += float64(p.HPMax)
		}
		n := float64(len(classes))
		return atk / n, def / n, hp / n
	}
	base := func(l int) (float64, float64, float64) {
		return 150 + 25*float64(l), 28 + 2.1*float64(l), 1 + 0.8*float64(l)
	}
	a25, d25, h25 := mean(25)
	bh, ba, bd := base(25)
	kh, kd := bh/a25, bd/a25
	ka := (ba - 0.5*d25) / h25
	for _, l := range []int{5, 10, 15, 20, 25, 30, 35, 40, 45, 50, 55, 60, 65, 70, 75, 80, 85, 90, 95, 100} {
		a, d, h := mean(l)
		ch, ca, cd := base(l)
		fmt.Printf("L%-3d player atk %5.0f def %5.0f hp %6.0f | fitted mob hp %6.0f atk %5.0f def %4.0f | v1 curve hp %5.0f atk %4.0f def %3.0f\n",
			l, a, d, h, kh*a, ka*h+0.5*d, kd*a, ch, ca, cd)
	}
}
