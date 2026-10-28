// Package partyxp is the party bonus on XP (DD-305): for every stretch of
// a finished job during which n ≥ 2 members of the earner's party (the
// earner included) worked the same kind of work, the earner's XP of that
// stretch grows by 3% per member — five together earn 115%. Nobody gives
// anything away; each member's own job gets the bonus when it ends.
package partyxp

import (
	"context"
	"sort"
	"time"

	"github.com/google/uuid"
)

// BonusPerMember is the party bonus on a stretch, per member working it.
const BonusPerMember = 0.03

// Factor is the XP multiplier of n members working together.
func Factor(n int) float64 {
	if n < 2 {
		return 1
	}
	return 1 + BonusPerMember*float64(n)
}

// Work kinds: a gathering kind and a crafting profession are kinds of their
// own ("gather:fishing", "craft:smelting"); hunts and dungeons are one each.
const (
	KindHunt     = "hunt"
	KindDungeon  = "dungeon"
	GatherPrefix = "gather:"
	CraftPrefix  = "craft:"
)

// Work is one finished job: whose, what kind, when it ran, its XP.
type Work struct {
	Kind      string
	Character uuid.UUID
	Start     time.Time
	End       time.Time
	XP        int64
}

// Interval is a partner's same-kind job window.
type Interval struct{ Start, End time.Time }

// Partner is a party member who worked alongside.
type Partner struct {
	CharacterID uuid.UUID `json:"character_id"`
	Name        string    `json:"name"`
}

// Outcome is the job's XP before (Earned) and after (Kept) the party bonus.
// Without partners Kept equals Earned.
type Outcome struct {
	Earned   int64     `json:"earned"`
	Kept     int64     `json:"kept"`
	Partners []Partner `json:"partners,omitempty"`
}

// Solo is the outcome of a job nobody worked alongside.
func Solo(xp int64) Outcome { return Outcome{Earned: xp, Kept: xp} }

// Boost returns w.XP, earned evenly over [w.Start, w.End], with the party
// bonus of every stretch the partners (each one's windows of the same kind
// of work) worked alongside, rounded.
func Boost(w Work, partners [][]Interval) int64 {
	total := w.End.Sub(w.Start).Seconds()
	if w.XP <= 0 || total <= 0 || len(partners) == 0 {
		return max(w.XP, 0)
	}
	cuts := []time.Time{w.Start, w.End}
	for _, p := range partners {
		for _, iv := range p {
			for _, t := range []time.Time{iv.Start, iv.End} {
				if t.After(w.Start) && t.Before(w.End) {
					cuts = append(cuts, t)
				}
			}
		}
	}
	sort.Slice(cuts, func(i, j int) bool { return cuts[i].Before(cuts[j]) })
	xp := 0.0
	for i := 1; i < len(cuts); i++ {
		a, b := cuts[i-1], cuts[i]
		if !b.After(a) {
			continue
		}
		mid := a.Add(b.Sub(a) / 2)
		n := 1
		for _, p := range partners {
			if covers(p, mid) {
				n++
			}
		}
		xp += float64(w.XP) * b.Sub(a).Seconds() / total * Factor(n)
	}
	return int64(xp + 0.5)
}

func covers(p []Interval, t time.Time) bool {
	for _, iv := range p {
		if !t.Before(iv.Start) && t.Before(iv.End) {
			return true
		}
	}
	return false
}

// Booster applies the party bonus to a finished job (Service satisfies).
type Booster interface {
	Boost(ctx context.Context, w Work) Outcome
}

// BoostIf boosts through b, or leaves the XP as it is without one.
func BoostIf(ctx context.Context, b Booster, w Work) Outcome {
	if b == nil {
		return Solo(w.XP)
	}
	return b.Boost(ctx, w)
}
