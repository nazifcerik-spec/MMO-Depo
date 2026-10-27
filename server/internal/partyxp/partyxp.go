// Package partyxp shares the XP of a finished job among the party members
// doing the same kind of work (DD-296). For every stretch of the job during
// which n members of the earner's party worked the same kind, the earner's
// XP of that stretch is pooled and split evenly over the n of them, plus 3%
// for each partner; the stretches nobody else worked stay the earner's own.
// Every member's own job is shared the same way when it ends, so two equal
// partners each end up with 103% of their solo XP.
package partyxp

import (
	"context"
	"sort"
	"time"

	"github.com/google/uuid"
)

// BonusPerPartner is the party bonus on a shared stretch, per partner.
const BonusPerPartner = 0.03

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
	Ref       uuid.UUID // the contract or job: keys the shares once
	Start     time.Time
	End       time.Time
	XP        int64
}

// Interval is a partner's same-kind job window.
type Interval struct{ Start, End time.Time }

// Share is what one partner received from a job.
type Share struct {
	CharacterID uuid.UUID `json:"character_id"`
	Name        string    `json:"name"`
	XP          int64     `json:"xp"`
}

// Outcome is the earner's side: the job's XP, what the earner keeps, and
// the partners' shares. Without partners Kept equals Earned.
type Outcome struct {
	Earned   int64   `json:"earned"`
	Kept     int64   `json:"kept"`
	Partners []Share `json:"partners,omitempty"`
}

// Solo is the outcome of a job nobody shared.
func Solo(xp int64) Outcome { return Outcome{Earned: xp, Kept: xp} }

// Split divides w.XP, earned evenly over [w.Start, w.End], between the
// earner and partners (each partner's windows of the same kind of work):
// on a stretch where n are working, each of them gets the stretch's XP / n
// × (1 + BonusPerPartner × (n−1)). Returns the earner's XP and each
// partner's, rounded.
func Split(w Work, partners [][]Interval) (keep int64, gifts []int64) {
	gifts = make([]int64, len(partners))
	total := w.End.Sub(w.Start).Seconds()
	if w.XP <= 0 || total <= 0 || len(partners) == 0 {
		return max(w.XP, 0), gifts
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
	keepF := 0.0
	giftF := make([]float64, len(partners))
	active := make([]bool, len(partners))
	for i := 1; i < len(cuts); i++ {
		a, b := cuts[i-1], cuts[i]
		if !b.After(a) {
			continue
		}
		mid := a.Add(b.Sub(a) / 2)
		n := 1
		for j, p := range partners {
			active[j] = covers(p, mid)
			if active[j] {
				n++
			}
		}
		part := float64(w.XP) * b.Sub(a).Seconds() / total
		if n == 1 {
			keepF += part
			continue
		}
		each := part / float64(n) * (1 + BonusPerPartner*float64(n-1))
		keepF += each
		for j := range partners {
			if active[j] {
				giftF[j] += each
			}
		}
	}
	for j, g := range giftF {
		gifts[j] = int64(g + 0.5)
	}
	return int64(keepF + 0.5), gifts
}

func covers(p []Interval, t time.Time) bool {
	for _, iv := range p {
		if !t.Before(iv.Start) && t.Before(iv.End) {
			return true
		}
	}
	return false
}

// Sharer shares a finished job's XP (Service satisfies).
type Sharer interface {
	Share(ctx context.Context, w Work) Outcome
}

// ShareIf shares through s, or keeps everything when there is no sharer.
func ShareIf(ctx context.Context, s Sharer, w Work) Outcome {
	if s == nil {
		return Solo(w.XP)
	}
	return s.Share(ctx, w)
}
