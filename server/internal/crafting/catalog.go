package crafting

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"sync/atomic"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Catalog is the in-memory profession + recipe lookup, hot-swapped by
// contentpack after YAML upsert.
type Catalog struct {
	professions map[string]Profession
	recipes     map[string]Recipe
	// recipesByProfession is the index for the "list all my professions'
	// available recipes" client view — small: dozens of entries per profession.
	recipesByProfession map[string][]Recipe
}

func NewCatalog(professions []Profession, recipes []Recipe) *Catalog {
	c := &Catalog{
		professions:         make(map[string]Profession, len(professions)),
		recipes:             make(map[string]Recipe, len(recipes)),
		recipesByProfession: make(map[string][]Recipe),
	}
	for _, p := range professions {
		c.professions[p.Code] = p
	}
	for _, r := range recipes {
		c.recipes[r.Code] = r
		c.recipesByProfession[r.ProfessionCode] = append(c.recipesByProfession[r.ProfessionCode], r)
	}
	return c
}

func (c *Catalog) Profession(code string) (Profession, bool) {
	p, ok := c.professions[code]
	return p, ok
}

func (c *Catalog) Recipe(code string) (Recipe, bool) {
	r, ok := c.recipes[code]
	return r, ok
}

func (c *Catalog) AllProfessions() []Profession {
	out := make([]Profession, 0, len(c.professions))
	for _, p := range c.professions {
		out = append(out, p)
	}
	return out
}

func (c *Catalog) AllRecipes() []Recipe {
	out := make([]Recipe, 0, len(c.recipes))
	for _, r := range c.recipes {
		out = append(out, r)
	}
	return out
}

// --- CatalogService ---------------------------------------------------------

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

// ReloadCatalog re-reads content.profession + content.recipe and swaps.
func (s *CatalogService) ReloadCatalog(ctx context.Context, pool *pgxpool.Pool) error {
	fresh, err := LoadCatalog(ctx, pool)
	if err != nil {
		return err
	}
	s.catPtr.Store(fresh)
	s.log.InfoContext(ctx, "crafting catalog reloaded",
		slog.Int("professions", len(fresh.AllProfessions())),
		slog.Int("recipes", len(fresh.AllRecipes())),
	)
	return nil
}

// LoadCatalog reads both tables into an in-memory Catalog.
func LoadCatalog(ctx context.Context, pool *pgxpool.Pool) (*Catalog, error) {
	profRows, err := pool.Query(ctx, `
		SELECT code, display_name_key, max_level, version
		  FROM content.profession
		 ORDER BY code`)
	if err != nil {
		return nil, fmt.Errorf("profession query: %w", err)
	}
	defer profRows.Close()
	var professions []Profession
	for profRows.Next() {
		var p Profession
		if err := profRows.Scan(&p.Code, &p.DisplayNameKey, &p.MaxLevel, &p.Version); err != nil {
			return nil, fmt.Errorf("profession scan: %w", err)
		}
		professions = append(professions, p)
	}
	if err := profRows.Err(); err != nil {
		return nil, err
	}

	recRows, err := pool.Query(ctx, `
		SELECT code, display_name_key, profession_code, min_skill,
		       inputs, gold_cost_copper,
		       output_template_code, output_stack, profession_xp_reward, version
		  FROM content.recipe
		 ORDER BY profession_code, min_skill, code`)
	if err != nil {
		return nil, fmt.Errorf("recipe query: %w", err)
	}
	defer recRows.Close()
	var recipes []Recipe
	for recRows.Next() {
		var r Recipe
		var inputsJSON []byte
		if err := recRows.Scan(
			&r.Code, &r.DisplayNameKey, &r.ProfessionCode, &r.MinSkill,
			&inputsJSON, &r.GoldCostCopper,
			&r.OutputTemplateCode, &r.OutputStack, &r.ProfessionXPReward, &r.Version,
		); err != nil {
			return nil, fmt.Errorf("recipe scan: %w", err)
		}
		if err := json.Unmarshal(inputsJSON, &r.Inputs); err != nil {
			return nil, fmt.Errorf("recipe %s: inputs unmarshal: %w", r.Code, err)
		}
		recipes = append(recipes, r)
	}
	if err := recRows.Err(); err != nil {
		return nil, err
	}
	return NewCatalog(professions, recipes), nil
}
