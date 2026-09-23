package activity

import (
	"encoding/json"
	"fmt"
	"math"
)

// YieldSummary is the machine-readable + human-readable outcome of a reveal.
// v1: gold_copper (via TxnService) + item drops (via ItemGranter, Sprint 8+).
type YieldSummary struct {
	Kind            Kind          `json:"kind"`
	ElapsedSeconds  int32         `json:"elapsed_seconds"`
	DurationSeconds int32         `json:"duration_seconds"`
	CompletionRatio float64       `json:"completion_ratio"` // 0.0-1.0
	GoldCopper      int64         `json:"gold_copper"`
	CommonCount     int           `json:"common_count"`
	RareCount       int           `json:"rare_count"`
	LegendaryCount  int           `json:"legendary_count"`
	Highlights      []string      `json:"highlights,omitempty"` // "1 legendary!", "3 rare gems"
	Items           []GrantedItem `json:"items,omitempty"`      // populated by Reveal after ItemGranter call
}

// yieldProfile is the per-kind rate table (v1 hardcoded placeholder).
// Content-driven table lands in Sprint 12 (content pack).
type yieldProfile struct {
	// Base resource rate per hour (Poisson λ for full duration).
	BasePerHour float64
	// Common item value in copper.
	CommonValueCopper int64
	// Rare drop rate per resource (Binomial p).
	RareChance float64
	// Rare item value in copper.
	RareValueCopper int64
	// Legendary drop rate per hour (Binomial n = hours ceil, p).
	LegendaryChancePerHour float64
	// Legendary value in copper.
	LegendaryValueCopper int64
}

// yieldProfiles — v1 placeholder balance. Numbers are illustrative; the real
// values come from the content pack + balance pass (Faz 00 doc: docs/economy).
var yieldProfiles = map[Kind]yieldProfile{
	KindMining: {
		BasePerHour: 40, CommonValueCopper: 5,
		RareChance: 0.01, RareValueCopper: 200,
		LegendaryChancePerHour: 0.02, LegendaryValueCopper: 5000,
	},
	KindFishing: {
		BasePerHour: 30, CommonValueCopper: 8,
		RareChance: 0.015, RareValueCopper: 250,
		LegendaryChancePerHour: 0.03, LegendaryValueCopper: 4000,
	},
	KindWoodcutting: {
		BasePerHour: 50, CommonValueCopper: 4,
		RareChance: 0.005, RareValueCopper: 150,
		LegendaryChancePerHour: 0.015, LegendaryValueCopper: 3500,
	},
	KindHerbalism: {
		BasePerHour: 35, CommonValueCopper: 7,
		RareChance: 0.02, RareValueCopper: 180,
		LegendaryChancePerHour: 0.025, LegendaryValueCopper: 4500,
	},
	KindHunting: {
		BasePerHour: 20, CommonValueCopper: 15,
		RareChance: 0.03, RareValueCopper: 400,
		LegendaryChancePerHour: 0.04, LegendaryValueCopper: 6000,
	},
}

// ComputeYield is the pure analytical yield calculation — no DB, no side
// effects. Deterministic given (contract, now). Callers use this to preview
// (dry-run) and to actually settle (reveal).
//
// Formula:
//
//	ratio = elapsed / duration              (0..1)
//	hours = elapsed / 3600
//	rate     = (base + mastery_flat) × (1 + mastery_mult)   (DD-220)
//	commons  = Poisson(rate × hours)
//	rares    = Binomial(commons, rare_p)    (rare given commons)
//	legendaries = Binomial(ceil(hours), legendary_chance_per_hour)
//	gold = commons×v_common + rares×v_rare + legendaries×v_legendary
func ComputeYield(c *Contract, elapsedSec int32) YieldSummary {
	prof, ok := yieldProfiles[c.Kind]
	if !ok {
		// Unknown kind — safe zero yield.
		return YieldSummary{Kind: c.Kind, ElapsedSeconds: elapsedSec, DurationSeconds: c.DurationSeconds}
	}

	if elapsedSec < 0 {
		elapsedSec = 0
	}
	if elapsedSec > c.DurationSeconds {
		elapsedSec = c.DurationSeconds
	}
	hours := float64(elapsedSec) / 3600.0
	ratio := 0.0
	if c.DurationSeconds > 0 {
		ratio = float64(elapsedSec) / float64(c.DurationSeconds)
	}

	rng := NewRNG(c.Seed)

	// Common resource count (Poisson) on the mastery-adjusted hourly rate.
	snap := parseSnapshot(c.Snapshot)
	lambda := (prof.BasePerHour + snap.YieldFlatPerHour) * (1 + snap.YieldMult) * hours
	commons := Poisson(rng, lambda)

	// Rare drops (Binomial on commons).
	rares := Binomial(rng, commons, prof.RareChance)

	// Legendary (Binomial on integer hours ceiling).
	hourTrials := int(math.Ceil(hours))
	if hourTrials < 0 {
		hourTrials = 0
	}
	legendaries := Binomial(rng, hourTrials, prof.LegendaryChancePerHour)

	// Gold value (v1: everything valued in gold; item instances Sprint 8).
	gold := int64(commons)*prof.CommonValueCopper +
		int64(rares)*prof.RareValueCopper +
		int64(legendaries)*prof.LegendaryValueCopper

	// Highlights for UI.
	var highlights []string
	if legendaries > 0 {
		highlights = append(highlights, fmt.Sprintf("✨ %d legendary drop%s", legendaries, plural(legendaries)))
	}
	if rares > 0 {
		highlights = append(highlights, fmt.Sprintf("★ %d rare drop%s", rares, plural(rares)))
	}

	return YieldSummary{
		Kind:            c.Kind,
		ElapsedSeconds:  elapsedSec,
		DurationSeconds: c.DurationSeconds,
		CompletionRatio: ratio,
		GoldCopper:      gold,
		CommonCount:     commons,
		RareCount:       rares,
		LegendaryCount:  legendaries,
		Highlights:      highlights,
	}
}

func plural(n int) string {
	if n == 1 {
		return ""
	}
	return "s"
}

// MarshalSummary serializes a YieldSummary as JSON for storage / response.
func MarshalSummary(y YieldSummary) json.RawMessage {
	b, _ := json.Marshal(y)
	return b
}
