// Package world implements the region hierarchy + travel loop (Sprint 16).
//
// Concepts:
//   - Region: a node in a shallow parent/child tree (continent > zone/city >
//     wilderness/dungeon). Every character sits at exactly one region.
//   - RegionLink: a non-parent shortcut adjacency (e.g. a road between two
//     sibling wildernesses).
//   - Travel: an in-transit state a character enters when leaving one region
//     for another. Duration is the shortest route over the parent↔child edge
//     weights and RegionLink travel_seconds (DD-280).
//
// Design decisions:
//   - Any region a road leads to is one journey (Catalog.Routes, Dijkstra);
//     only the destination's level gate applies (DD-280).
//   - Travel is server-authoritative and read-only until the arrival time.
//     Complete refuses when `now < arrives_at` (E_TRAVEL_IN_PROGRESS);
//     Cancel returns the character to `from_region_code` with no reward or
//     refund — travel is intentionally cheap in v1 (no gold cost).
//   - Character location is written ONLY on completion. Cancel is a no-op on
//     location, so nothing changes.
//   - Slot cap 1: a character with an active travel row cannot start a hunt,
//     dungeon, activity, or another travel. Enforced by
//     `ux_travel_in_progress` partial UNIQUE.
package world

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// Kind matches the CHECK constraint on world.region.kind.
type Kind string

const (
	KindContinent  Kind = "continent"
	KindZone       Kind = "zone"
	KindCity       Kind = "city"
	KindWilderness Kind = "wilderness"
	KindDungeon    Kind = "dungeon"
)

// Region is one node of the world tree.
type Region struct {
	Code                    string
	DisplayNameKey          string
	Kind                    Kind
	ParentCode              *string
	MinLevel                int16
	TravelSecondsFromParent int32
	Safe                    bool
	IsStarting              bool
	Version                 int32
	MapX, MapY              int16 // position on the travel map (DD-279)
}

// RegionLink is one non-parent adjacency. Codes are stored canonically
// (RegionA < RegionB) — helpers pick min/max on lookup.
type RegionLink struct {
	RegionA       string
	RegionB       string
	TravelSeconds int32
}

// TravelState mirrors the CHECK constraint on world.travel.state.
type TravelState string

const (
	TravelInProgress TravelState = "in_progress"
	TravelCompleted  TravelState = "completed"
	TravelCancelled  TravelState = "cancelled"
)

// Travel is one in-transit or historical journey.
type Travel struct {
	ID             uuid.UUID
	CharacterID    uuid.UUID
	FromRegionCode string
	ToRegionCode   string
	StartedAt      time.Time
	ArrivesAt      time.Time
	State          TravelState
	CompletedAt    *time.Time
	CancelledAt    *time.Time
	CreatedAt      time.Time
	UpdatedAt      time.Time
	Version        int32
}

// --- Errors -----------------------------------------------------------------

var (
	ErrNotFound         = errors.New("world: not found")
	ErrForbidden        = errors.New("world: forbidden")
	ErrUnknownRegion    = errors.New("world: unknown region_code")
	ErrTravelInProgress = errors.New("world: character already has an active travel")
	ErrTravelNotActive  = errors.New("world: travel is not in progress")
	ErrTravelNotReady   = errors.New("world: travel not yet arrived")
	ErrSameRegion       = errors.New("world: destination equals current region")
	ErrNoRoute          = errors.New("world: no road between regions")
	ErrLevelRequirement = errors.New("world: character level below region min_level")
	ErrValidation       = errors.New("world: validation")
	ErrNoStartingRegion = errors.New("world: no starting region defined")
)
