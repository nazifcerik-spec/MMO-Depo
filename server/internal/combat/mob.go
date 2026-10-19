package combat

import (
	"context"
	"fmt"
	"log/slog"
	"sync/atomic"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Mob is one enemy template loaded from content pack. Level + XP reward drive
// combat scaling; the stats block feeds Combatant. Sprint 14 (hunting) rolls
// a per-encounter Mob out of this template + world scaling.
type Mob struct {
	Code           string
	DisplayNameKey string
	Level          int16
	HPMax          int32
	Attack         int32
	Defense        int32
	CritPct        int32
	HitPct         int32
	EvasionPct     int32
	XPReward       int32
	GoldCopperMin  int64
	GoldCopperMax  int64
	RaidBoss       bool   // fought only in raids (DD-261)
	RegionCode     string // where it is hunted (DD-273); empty for raid bosses
}

// AsCombatant builds a fresh Combatant with full HP for a live encounter.
func (m Mob) AsCombatant() Combatant {
	return Combatant{
		Name:       m.Code,
		HPMax:      m.HPMax,
		HPCurrent:  m.HPMax,
		Attack:     m.Attack,
		Defense:    m.Defense,
		CritPct:    m.CritPct,
		HitPct:     m.HitPct,
		EvasionPct: m.EvasionPct,
		Level:      m.Level,
	}
}

// --- Catalog ----------------------------------------------------------------

// Catalog is the in-memory mob lookup. Hot-swapped by contentpack via
// atomic.Pointer so reads stay lock-free.
type Catalog struct {
	byCode map[string]Mob
}

func NewCatalog(mobs []Mob) *Catalog {
	m := make(map[string]Mob, len(mobs))
	for _, mob := range mobs {
		m[mob.Code] = mob
	}
	return &Catalog{byCode: m}
}

func (c *Catalog) Get(code string) (Mob, bool) {
	m, ok := c.byCode[code]
	return m, ok
}

func (c *Catalog) All() []Mob {
	out := make([]Mob, 0, len(c.byCode))
	for _, m := range c.byCode {
		out = append(out, m)
	}
	return out
}

// --- Service ---------------------------------------------------------------

// Service wraps the mob catalog + provides ReloadCatalog for contentpack.
type Service struct {
	catPtr atomic.Pointer[Catalog]
	log    *slog.Logger
}

func NewService(cat *Catalog, logger *slog.Logger) *Service {
	if logger == nil {
		logger = slog.Default()
	}
	s := &Service{log: logger}
	s.catPtr.Store(cat)
	return s
}

func (s *Service) Catalog() *Catalog { return s.catPtr.Load() }

// ReloadCatalog re-reads content.mob_template rows and atomically swaps the
// in-memory catalog. Called by contentpack after YAML upsert.
func (s *Service) ReloadCatalog(ctx context.Context, pool *pgxpool.Pool) error {
	fresh, err := LoadCatalog(ctx, pool)
	if err != nil {
		return err
	}
	s.catPtr.Store(fresh)
	s.log.InfoContext(ctx, "mob catalog reloaded", slog.Int("mobs", len(fresh.All())))
	return nil
}

// LoadCatalog reads content.mob_template rows and returns a populated Catalog.
func LoadCatalog(ctx context.Context, pool *pgxpool.Pool) (*Catalog, error) {
	rows, err := pool.Query(ctx, `
		SELECT code, display_name_key, level,
		       hp_max, attack, defense, crit_pct, hit_pct, evasion_pct,
		       xp_reward, gold_copper_min, gold_copper_max, raid_boss, COALESCE(region_code, '')
		FROM content.mob_template
		ORDER BY level, code`)
	if err != nil {
		return nil, fmt.Errorf("mob catalog query: %w", err)
	}
	defer rows.Close()

	var mobs []Mob
	for rows.Next() {
		var m Mob
		if err := rows.Scan(
			&m.Code, &m.DisplayNameKey, &m.Level,
			&m.HPMax, &m.Attack, &m.Defense, &m.CritPct, &m.HitPct, &m.EvasionPct,
			&m.XPReward, &m.GoldCopperMin, &m.GoldCopperMax, &m.RaidBoss, &m.RegionCode,
		); err != nil {
			return nil, fmt.Errorf("mob scan: %w", err)
		}
		mobs = append(mobs, m)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return NewCatalog(mobs), nil
}
