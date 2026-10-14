package activity

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

// Gathering professions (DD-244): every activity kind is a profession with
// its own level (1..100). A contract gathers one resource (node) the
// character's level unlocks; its tier sets how fast it comes and how much
// gathering XP each unit gives.

var (
	ErrUnknownTarget  = errors.New("activity: unknown gathering resource")
	ErrGatherLevelLow = errors.New("activity: gathering level too low for this resource")
)

// GatherNode is one resource of a gathering profession.
type GatherNode struct {
	Code  string `json:"code"`
	Level int    `json:"level"`
	Tier  int    `json:"tier"`
}

// GatherProgress is a profession's level after a reveal.
type GatherProgress struct {
	Level     int  `json:"level"`
	XP        int  `json:"xp"`
	XPToNext  int  `json:"xp_to_next"`
	Gained    int  `json:"gained"`
	LeveledUp bool `json:"leveled_up"`
}

// Gathering is the resource tables plus the character's gathering skill.
type Gathering interface {
	Nodes(kind Kind) []GatherNode // lowest unlock level first
	Level(ctx context.Context, characterID uuid.UUID, kind Kind) (int, error)
	AddXP(ctx context.Context, characterID uuid.UUID, kind Kind, xp int) (GatherProgress, error)
}

// WithGathering turns on gathering professions (DD-244). Without it a
// contract gathers nothing but gold (tests, transition).
func (s *Service) WithGathering(g Gathering) *Service {
	s.gather = g
	return s
}

// tierRate slows higher tiers: T0 at the kind's full hourly rate, T7 at 30%.
var tierRate = [8]float64{1.0, 0.9, 0.8, 0.7, 0.6, 0.5, 0.4, 0.3}

// tierXP is gathering XP per unit gathered.
var tierXP = [8]int{8, 12, 18, 26, 36, 50, 68, 90}

func tierIndex(t int) int { return min(7, max(0, t)) }

// pickNode resolves the requested resource, or the best one the level
// unlocks when none is named (older clients).
func pickNode(nodes []GatherNode, level int, target string) (GatherNode, error) {
	if target == "" {
		best := -1
		for i, n := range nodes {
			if n.Level <= level {
				best = i
			}
		}
		if best < 0 {
			return GatherNode{}, ErrUnknownTarget
		}
		return nodes[best], nil
	}
	for _, n := range nodes {
		if n.Code == target {
			if n.Level > level {
				return GatherNode{}, ErrGatherLevelLow
			}
			return n, nil
		}
	}
	return GatherNode{}, ErrUnknownTarget
}
