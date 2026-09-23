// Package worldevent implements scheduled world bosses and region invasions
// with AFK attendance, plus seasonal daily-login tokens (Sprint 25,
// docs/world-events/*).
package worldevent

import (
	"errors"
	"sort"
	"time"
)

// Tier is one reward bracket (world-boss.md §6).
type Tier struct {
	MinScore    int64 `json:"min_score"`
	GoldCopper  int64 `json:"gold_copper"`
	EventTokens int64 `json:"event_tokens"`
}

const (
	RegistrationWindow = 24 * time.Hour
	// FailedTierCap: when the community misses the goal, attendees still get
	// rewards but at most tier 2.
	FailedTierCap = 2
)

// NextOccurrences lists the weekly slot (day_of_week 0=Sunday, hour UTC)
// starting times in (now, now+horizon].
func NextOccurrences(now time.Time, dow, hour int, horizon time.Duration) []time.Time {
	now = now.UTC()
	day := time.Date(now.Year(), now.Month(), now.Day(), hour, 0, 0, 0, time.UTC)
	var out []time.Time
	for i := 0; i <= int(horizon/(24*time.Hour))+7; i++ {
		t := day.AddDate(0, 0, i)
		if int(t.Weekday()) != dow || !t.After(now) {
			continue
		}
		if t.After(now.Add(horizon)) {
			break
		}
		out = append(out, t)
	}
	return out
}

// Power is the AFK damage proxy snapshotted at enlist time.
func Power(attack, defense, hpMax int) int {
	p := attack*2 + defense + hpMax/10
	if p < 1 {
		return 1
	}
	return p
}

// Contribution = power × minutes attended, capped so AFK alone cannot pass
// the cap score (world-boss.md §5).
func Contribution(power int, minutes int64, capScore int64) int64 {
	if minutes < 0 {
		minutes = 0
	}
	s := int64(power) * minutes
	if capScore > 0 && s > capScore {
		return capScore
	}
	return s
}

// TierFor returns the 1-based tier reached by score (0 = none). Tiers must be
// ascending by MinScore; failed events are capped at FailedTierCap.
func TierFor(score int64, tiers []Tier, succeeded bool) int {
	sorted := append([]Tier(nil), tiers...)
	sort.Slice(sorted, func(i, j int) bool { return sorted[i].MinScore < sorted[j].MinScore })
	tier := 0
	for i, t := range sorted {
		if score >= t.MinScore {
			tier = i + 1
		}
	}
	if !succeeded && tier > FailedTierCap {
		tier = FailedTierCap
	}
	return tier
}

// AttendedMinutes counts whole minutes from max(registered, start) to end.
func AttendedMinutes(registered, start, end time.Time) int64 {
	from := start
	if registered.After(from) {
		from = registered
	}
	if !end.After(from) {
		return 0
	}
	return int64(end.Sub(from) / time.Minute)
}

var (
	ErrNotFound         = errors.New("worldevent: not found")
	ErrNotOpen          = errors.New("worldevent: registration is not open")
	ErrLevelRequirement = errors.New("worldevent: character level too low")
	ErrWrongRegion      = errors.New("worldevent: travel to the event region first")
	ErrAlreadyAttending = errors.New("worldevent: already enlisted")
	ErrEventNotActive   = errors.New("worldevent: seasonal event is not running")
	ErrAlreadyClaimed   = errors.New("worldevent: today's reward already claimed")
)
