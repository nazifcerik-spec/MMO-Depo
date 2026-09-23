// Package progression implements mastery, achievements, titles and faction
// reputation (Sprint 24, docs/progression/*). Game systems report Events;
// content decides which counters each event feeds.
package progression

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

// Recorder is what game systems depend on to report events. Service
// satisfies it; nil recorders are skipped by RecordIf.
type Recorder interface {
	Record(ctx context.Context, characterID uuid.UUID, events ...Event)
}

// RecordIf is a nil-safe helper that drops zero-amount events.
func RecordIf(ctx context.Context, r Recorder, characterID uuid.UUID, events ...Event) {
	if r == nil {
		return
	}
	kept := events[:0]
	for _, e := range events {
		if e.Amount > 0 {
			kept = append(kept, e)
		}
	}
	if len(kept) > 0 {
		r.Record(ctx, characterID, kept...)
	}
}

// Event is one thing that happened to a character, e.g. {mob_killed, 12}.
type Event struct {
	Kind   string
	Amount int64
}

// Event kinds emitted by game systems.
const (
	EventMobKilled      = "mob_killed"
	EventBossKilled     = "boss_killed"
	EventDungeonCleared = "dungeon_cleared" // also "dungeon_cleared:<code>"
	EventItemCrafted    = "item_crafted"
	EventGatherPrefix   = "gather_" // gather_mining, gather_fishing, ...
)

// MasteryTier returns the tier reached with `xp` given per-tier increments
// (docs/progression/mastery.md §4), capped at maxTier.
func MasteryTier(xp int64, curve []int64, maxTier int16) int16 {
	var cum int64
	var tier int16
	for _, step := range curve {
		if tier >= maxTier {
			break
		}
		cum += step
		if xp < cum {
			break
		}
		tier++
	}
	return tier
}

// NextTierXP is the cumulative xp needed for the next tier, or -1 at max.
func NextTierXP(tier int16, curve []int64, maxTier int16) int64 {
	if tier >= maxTier || int(tier) >= len(curve) {
		return -1
	}
	var cum int64
	for i := 0; i <= int(tier); i++ {
		cum += curve[i]
	}
	return cum
}

// Reputation tiers (docs/progression/reputation.md §3), cumulative thresholds.
var RepTierNames = []string{"hated", "hostile", "unfriendly", "neutral", "friendly", "honored", "revered", "exalted"}
var repThresholds = []int{-6000, -3000, -1000, 0, 3000, 9000, 21000, 42000}

const (
	RepMin = -6000
	RepMax = 42000
)

// RepTier maps signed rep xp to a tier index 0..7.
func RepTier(rep int) int16 {
	tier := int16(0)
	for i, t := range repThresholds {
		if rep >= t {
			tier = int16(i)
		}
	}
	return tier
}

// ApplyRepGain applies delta under the daily cap for positive gains
// (losses are never capped) and clamps to [RepMin, RepMax]. It returns the
// new rep and new daily usage.
func ApplyRepGain(rep, used, delta, dailyCap int) (newRep, newUsed int) {
	if delta > 0 {
		room := dailyCap - used
		if room <= 0 {
			return rep, used
		}
		if delta > room {
			delta = room
		}
		used += delta
	}
	rep += delta
	if rep < RepMin {
		rep = RepMin
	}
	if rep > RepMax {
		rep = RepMax
	}
	return rep, used
}

// SameUTCDay reports whether a and b fall on the same UTC calendar day.
func SameUTCDay(a, b time.Time) bool {
	ay, am, ad := a.UTC().Date()
	by, bm, bd := b.UTC().Date()
	return ay == by && am == bm && ad == bd
}

var (
	ErrTitleNotOwned = errors.New("progression: title not unlocked")
)
