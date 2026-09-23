package world

import (
	"context"
	"fmt"
	"log/slog"
	"sync/atomic"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Catalog is the in-memory region graph. Immutable — Reload builds a fresh
// Catalog and CatalogService atomically swaps it.
type Catalog struct {
	regions  map[string]Region       // code → region
	linksBy  map[string][]RegionLink // code → all links touching this region
	starting string                  // code of the starting region (empty if none)
}

// NewCatalog builds the in-memory graph. Links get indexed under BOTH
// endpoints so `LinksFor(code)` is O(1) regardless of canonical ordering.
func NewCatalog(regions []Region, links []RegionLink) *Catalog {
	c := &Catalog{
		regions: make(map[string]Region, len(regions)),
		linksBy: make(map[string][]RegionLink, len(regions)),
	}
	for _, r := range regions {
		c.regions[r.Code] = r
		if r.IsStarting {
			// Migration + validation enforce a single starting region; a race
			// during hot reload would leave the LAST one seen as canonical.
			c.starting = r.Code
		}
	}
	for _, l := range links {
		c.linksBy[l.RegionA] = append(c.linksBy[l.RegionA], l)
		c.linksBy[l.RegionB] = append(c.linksBy[l.RegionB], l)
	}
	return c
}

func (c *Catalog) Get(code string) (Region, bool) { r, ok := c.regions[code]; return r, ok }

func (c *Catalog) All() []Region {
	out := make([]Region, 0, len(c.regions))
	for _, r := range c.regions {
		out = append(out, r)
	}
	return out
}

// StartingCode returns the code marked is_starting=true, or empty when none.
func (c *Catalog) StartingCode() string { return c.starting }

// LinksFor returns every non-parent link touching `code`.
func (c *Catalog) LinksFor(code string) []RegionLink { return c.linksBy[code] }

// TravelDuration returns the direct-edge duration between two region codes,
// or (0, false) when no direct edge exists.
//
// A direct edge is: parent↔child in the tree, OR a matching entry in
// world.region_link. Multi-hop routing is intentionally NOT supported in
// MVP — the caller sees ErrNoRoute and must break the trip into segments.
func (c *Catalog) TravelDuration(fromCode, toCode string) (int32, bool) {
	if fromCode == toCode {
		return 0, false
	}
	from, okF := c.regions[fromCode]
	to, okT := c.regions[toCode]
	if !okF || !okT {
		return 0, false
	}
	// Parent↔child: symmetric — weight is the CHILD's travel_seconds_from_parent
	// (parent has no such weight because it may have many children).
	if to.ParentCode != nil && *to.ParentCode == fromCode {
		return to.TravelSecondsFromParent, true
	}
	if from.ParentCode != nil && *from.ParentCode == toCode {
		return from.TravelSecondsFromParent, true
	}
	// Explicit link (either direction).
	for _, l := range c.linksBy[fromCode] {
		if (l.RegionA == fromCode && l.RegionB == toCode) ||
			(l.RegionB == fromCode && l.RegionA == toCode) {
			return l.TravelSeconds, true
		}
	}
	return 0, false
}

// Neighbors returns every region directly reachable from `code` — children,
// parent, and explicit links. Used by GET /location to render "where can I
// go from here".
type NeighborInfo struct {
	Region        Region
	TravelSeconds int32
	Via           string // "child" | "parent" | "link"
}

func (c *Catalog) Neighbors(code string) []NeighborInfo {
	self, ok := c.regions[code]
	if !ok {
		return nil
	}
	out := make([]NeighborInfo, 0, 4)
	// Parent
	if self.ParentCode != nil {
		if p, ok := c.regions[*self.ParentCode]; ok {
			out = append(out, NeighborInfo{
				Region: p, TravelSeconds: self.TravelSecondsFromParent, Via: "parent",
			})
		}
	}
	// Children
	for _, r := range c.regions {
		if r.ParentCode != nil && *r.ParentCode == code {
			out = append(out, NeighborInfo{
				Region: r, TravelSeconds: r.TravelSecondsFromParent, Via: "child",
			})
		}
	}
	// Links
	for _, l := range c.linksBy[code] {
		other := l.RegionA
		if other == code {
			other = l.RegionB
		}
		if r, ok := c.regions[other]; ok {
			out = append(out, NeighborInfo{
				Region: r, TravelSeconds: l.TravelSeconds, Via: "link",
			})
		}
	}
	return out
}

// --- CatalogService ---------------------------------------------------------

// CatalogService wraps the catalog behind an atomic pointer so contentpack
// hot-reload can swap the graph without blocking readers.
type CatalogService struct {
	catPtr atomic.Pointer[Catalog]
	log    *slog.Logger
}

func NewCatalogService(cat *Catalog, logger *slog.Logger) *CatalogService {
	if logger == nil {
		logger = slog.Default()
	}
	s := &CatalogService{log: logger}
	s.catPtr.Store(cat)
	return s
}

func (s *CatalogService) Catalog() *Catalog { return s.catPtr.Load() }

// ReloadCatalog re-reads world.region + world.region_link and swaps the
// catalog. Called by contentpack after YAML upsert.
func (s *CatalogService) ReloadCatalog(ctx context.Context, pool *pgxpool.Pool) error {
	fresh, err := LoadCatalog(ctx, pool)
	if err != nil {
		return err
	}
	s.catPtr.Store(fresh)
	s.log.InfoContext(ctx, "world catalog reloaded",
		slog.Int("regions", len(fresh.All())),
		slog.String("starting", fresh.StartingCode()),
	)
	return nil
}

// LoadCatalog reads world.region + world.region_link and returns a fresh
// Catalog. Two queries kept simple — small tables (~ hundreds max).
func LoadCatalog(ctx context.Context, pool *pgxpool.Pool) (*Catalog, error) {
	regRows, err := pool.Query(ctx, `
		SELECT code, display_name_key, kind, parent_code, min_level,
		       travel_seconds_from_parent, safe, is_starting, version
		  FROM world.region
		 ORDER BY code`)
	if err != nil {
		return nil, fmt.Errorf("region query: %w", err)
	}
	defer regRows.Close()
	var regions []Region
	for regRows.Next() {
		var r Region
		var kind string
		if err := regRows.Scan(
			&r.Code, &r.DisplayNameKey, &kind, &r.ParentCode, &r.MinLevel,
			&r.TravelSecondsFromParent, &r.Safe, &r.IsStarting, &r.Version,
		); err != nil {
			return nil, fmt.Errorf("region scan: %w", err)
		}
		r.Kind = Kind(kind)
		regions = append(regions, r)
	}
	if err := regRows.Err(); err != nil {
		return nil, err
	}

	linkRows, err := pool.Query(ctx, `
		SELECT region_a_code, region_b_code, travel_seconds FROM world.region_link`)
	if err != nil {
		return nil, fmt.Errorf("region_link query: %w", err)
	}
	defer linkRows.Close()
	var links []RegionLink
	for linkRows.Next() {
		var l RegionLink
		if err := linkRows.Scan(&l.RegionA, &l.RegionB, &l.TravelSeconds); err != nil {
			return nil, fmt.Errorf("region_link scan: %w", err)
		}
		links = append(links, l)
	}
	if err := linkRows.Err(); err != nil {
		return nil, err
	}
	return NewCatalog(regions, links), nil
}
