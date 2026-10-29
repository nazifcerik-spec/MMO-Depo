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
	// ErrGatherRegion: the resource grows on another map (DD-317).
	ErrGatherRegion = errors.New("activity: the resource grows on another map")
)

// GatherNode is one resource of a gathering profession.
type GatherNode struct {
	Code  string `json:"code"`
	Level int    `json:"level"`
	Tier  int    `json:"tier"`
	// Region is the map it grows on (DD-317); empty = anywhere.
	Region string `json:"region,omitempty"`
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

// BagRoom tells whether the bag can take more items (DD-246);
// inventory.Service satisfies. An empty templateCode asks for an empty slot.
type BagRoom interface {
	HasRoom(ctx context.Context, characterID uuid.UUID, templateCode string, qty int32) (bool, error)
}

// WithBagRoom refuses to start a run while the bag has no room (DD-246).
func (s *Service) WithBagRoom(b BagRoom) *Service {
	s.bag = b
	return s
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

// TierXP is the gathering XP one unit of a tier pays (helpers use it, DD-319).
func TierXP(tier int) int { return tierXP[tierIndex(tier)] }

// pickNode resolves the requested resource, or the best one the level
// unlocks when none is named (older clients).
// pickNode picks the target, or with none the best unlocked resource; in
// is the regions the character gathers in (nil = everywhere, DD-317).
func pickNode(nodes []GatherNode, level int, target string, in func(region string) bool) (GatherNode, error) {
	here := func(n GatherNode) bool { return in == nil || n.Region == "" || in(n.Region) }
	if target == "" {
		best, unlocked := -1, false
		for i, n := range nodes {
			if n.Level <= level {
				unlocked = true
				if here(n) {
					best = i
				}
			}
		}
		if best < 0 {
			if unlocked {
				return GatherNode{}, ErrGatherRegion
			}
			return GatherNode{}, ErrUnknownTarget
		}
		return nodes[best], nil
	}
	for _, n := range nodes {
		if n.Code == target {
			if n.Level > level {
				return GatherNode{}, ErrGatherLevelLow
			}
			if !here(n) {
				return GatherNode{}, ErrGatherRegion
			}
			return n, nil
		}
	}
	return GatherNode{}, ErrUnknownTarget
}

// RegionGate is the regions a character reaches from where it stands
// (world.Service satisfies, as for hunting, DD-273).
type RegionGate interface {
	HuntRegions(ctx context.Context, characterID uuid.UUID) ([]string, error)
}

// WithRegionGate gathers each resource only in reach of its map (DD-317).
func (s *Service) WithRegionGate(g RegionGate) *Service {
	s.regions = g
	return s
}
