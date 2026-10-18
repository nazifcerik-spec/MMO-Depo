package dungeon

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"sync/atomic"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Catalog is the in-memory dungeon template lookup. Hot-swapped by
// contentpack after YAML upsert (same pattern as combat.Catalog).
type Catalog struct {
	byCode map[string]Template
}

func NewCatalog(templates []Template) *Catalog {
	m := make(map[string]Template, len(templates))
	for _, t := range templates {
		m[t.Code] = t
	}
	return &Catalog{byCode: m}
}

func (c *Catalog) Get(code string) (Template, bool) {
	t, ok := c.byCode[code]
	return t, ok
}

func (c *Catalog) All() []Template {
	out := make([]Template, 0, len(c.byCode))
	for _, t := range c.byCode {
		out = append(out, t)
	}
	return out
}

// --- Service ---------------------------------------------------------------

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

// ReloadCatalog re-reads content.dungeon_template rows and swaps the
// in-memory catalog. Called by contentpack loader after YAML upsert.
func (s *CatalogService) ReloadCatalog(ctx context.Context, pool *pgxpool.Pool) error {
	fresh, err := LoadCatalog(ctx, pool)
	if err != nil {
		return err
	}
	s.catPtr.Store(fresh)
	s.log.InfoContext(ctx, "dungeon catalog reloaded", slog.Int("dungeons", len(fresh.All())))
	return nil
}

// LoadCatalog reads content.dungeon_template rows and returns a populated
// Catalog. Waves stored as JSONB [{mob_code, count}, …].
func LoadCatalog(ctx context.Context, pool *pgxpool.Pool) (*Catalog, error) {
	rows, err := pool.Query(ctx, `
		SELECT code, display_name_key, min_level, duration_seconds,
		       waves, boss_mob_code, bonus_gold_copper, bonus_xp, group_min, group_max
		  FROM content.dungeon_template
		 ORDER BY min_level, code`)
	if err != nil {
		return nil, fmt.Errorf("dungeon catalog query: %w", err)
	}
	defer rows.Close()

	var out []Template
	for rows.Next() {
		var t Template
		var wavesJSON []byte
		if err := rows.Scan(
			&t.Code, &t.DisplayNameKey, &t.MinLevel, &t.DurationSecs,
			&wavesJSON, &t.BossMobCode, &t.BonusGoldCopper, &t.BonusXP, &t.GroupMin, &t.GroupMax,
		); err != nil {
			return nil, fmt.Errorf("dungeon scan: %w", err)
		}
		if err := json.Unmarshal(wavesJSON, &t.Waves); err != nil {
			return nil, fmt.Errorf("dungeon %s: unmarshal waves: %w", t.Code, err)
		}
		out = append(out, t)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return NewCatalog(out), nil
}
