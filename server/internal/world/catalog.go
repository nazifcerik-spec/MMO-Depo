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
	links    []RegionLink            // every link once (the map's roads)
	linksBy  map[string][]RegionLink // code → all links touching this region
	starting string                  // code of the starting region (empty if none)
}

// NewCatalog builds the in-memory graph. Links get indexed under BOTH
// endpoints so `LinksFor(code)` is O(1) regardless of canonical ordering.
func NewCatalog(regions []Region, links []RegionLink) *Catalog {
	c := &Catalog{
		links:   links,
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

// Links lists every explicit link once (canonical order).
func (c *Catalog) Links() []RegionLink { return c.links }

func (c *Catalog) All() []Region {
	out := make([]Region, 0, len(c.regions))
	for _, r := range c.regions {
		out = append(out, r)
	}
	return out
}

// StartingCode returns the code marked is_starting=true, or empty when none.
func (c *Catalog) StartingCode() string { return c.starting }

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

// Routes is the shortest travel time from `from` to every region it can
// reach, with the path there (DD-280): Dijkstra over the travel edges
// (parent↔child and links). The origin itself is left out.
func (c *Catalog) Routes(from string) map[string]Route {
	if _, ok := c.regions[from]; !ok {
		return nil
	}
	best := map[string]int32{from: 0}
	prev := map[string]string{}
	done := map[string]bool{}
	for {
		cur, curT := "", int32(-1)
		for code, t := range best { // small graph: a linear scan beats a heap
			if !done[code] && (curT < 0 || t < curT || (t == curT && code < cur)) {
				cur, curT = code, t
			}
		}
		if curT < 0 {
			break
		}
		done[cur] = true
		for _, n := range c.Neighbors(cur) {
			t := curT + n.TravelSeconds
			if old, seen := best[n.Region.Code]; !seen || t < old {
				best[n.Region.Code], prev[n.Region.Code] = t, cur
			}
		}
	}
	out := make(map[string]Route, len(best)-1)
	for code, t := range best {
		if code == from {
			continue
		}
		path := []string{code}
		for p := prev[code]; p != ""; p = prev[p] {
			path = append([]string{p}, path...)
		}
		out[code] = Route{Seconds: t, Path: path}
	}
	return out
}

// Route is one destination's shortest journey: total seconds and the
// regions passed, origin first.
type Route struct {
	Seconds int32    `json:"seconds"`
	Path    []string `json:"path"`
}

// HuntRegions are the regions whose mobs a character standing in `code` may
// hunt (DD-273): the region itself and those linked to it by road. Parent
// and child edges do not count — a continent would otherwise reach every
// zone on it.
func (c *Catalog) HuntRegions(code string) []string {
	if _, ok := c.regions[code]; !ok {
		return nil
	}
	out := []string{code}
	for _, l := range c.linksBy[code] {
		other := l.RegionA
		if other == code {
			other = l.RegionB
		}
		out = append(out, other)
	}
	return out
}

// DungeonRegions lists the regions whose dungeons and raids a character in
// region code can see and enter (DD-302): the hunting reach plus the parent
// region, so a dungeon's entrance node (a child of its zone) counts too.
func (c *Catalog) DungeonRegions(code string) []string {
	out := c.HuntRegions(code)
	if r, ok := c.regions[code]; ok && r.ParentCode != nil {
		out = append(out, *r.ParentCode)
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
		       travel_seconds_from_parent, safe, is_starting, version,
		       COALESCE(map_x, 0), COALESCE(map_y, 0)
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
			&r.MapX, &r.MapY,
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
