package item

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Catalog is the in-memory template lookup — loaded once at boot from
// content.item_template rows. Sprint 12 YAML loader replaces the DB read with
// a versioned content pack hot-reload.
type Catalog struct {
	byCode map[string]Template
	nodes  map[string][]Template // gather kind → nodes by unlock level (DD-244)
}

func NewCatalog(templates []Template) *Catalog {
	m := make(map[string]Template, len(templates))
	nodes := map[string][]Template{}
	for _, t := range templates {
		m[t.Code] = t
		if t.GatherKind != "" {
			nodes[t.GatherKind] = append(nodes[t.GatherKind], t)
		}
	}
	for _, list := range nodes {
		sort.Slice(list, func(i, j int) bool {
			if list[i].GatherLevel != list[j].GatherLevel {
				return list[i].GatherLevel < list[j].GatherLevel
			}
			return list[i].Code < list[j].Code
		})
	}
	return &Catalog{byCode: m, nodes: nodes}
}

// Nodes lists a gathering profession's resources, lowest unlock level first.
func (c *Catalog) Nodes(kind string) []Template { return c.nodes[kind] }

// Get returns a template by code; ErrUnknownTemplate on miss.
func (c *Catalog) Get(code string) (Template, error) {
	t, ok := c.byCode[code]
	if !ok {
		return Template{}, fmt.Errorf("%w: %s", ErrUnknownTemplate, code)
	}
	return t, nil
}

// All returns a copy of the templates slice.
func (c *Catalog) All() []Template {
	out := make([]Template, 0, len(c.byCode))
	for _, t := range c.byCode {
		out = append(out, t)
	}
	return out
}

// LoadCatalog reads content.item_template rows and returns a populated Catalog.
func LoadCatalog(ctx context.Context, pool *pgxpool.Pool) (*Catalog, error) {
	const q = `
		SELECT code, display_name_key, item_class, item_subclass, min_level,
		       base_value_copper, base_weight_g, stackable, max_stack, bind_kind,
		       no_trade, no_sell, affixable,
		       equip_slot, class_restriction, stat_bonuses, use_effect,
		       COALESCE(gather_kind, ''), COALESCE(gather_level, 0), COALESCE(gather_tier, 0)
		  FROM content.item_template`
	rows, err := pool.Query(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("load item templates: %w", err)
	}
	defer rows.Close()

	list := []Template{}
	for rows.Next() {
		var t Template
		var bind string
		var bonuses []byte
		var useEffect []byte
		if err := rows.Scan(
			&t.Code, &t.DisplayNameKey, &t.ItemClass, &t.ItemSubclass, &t.MinLevel,
			&t.BaseValueCopper, &t.BaseWeightG, &t.Stackable, &t.MaxStack, &bind,
			&t.NoTrade, &t.NoSell, &t.Affixable,
			&t.EquipSlot, &t.ClassRestriction, &bonuses, &useEffect,
			&t.GatherKind, &t.GatherLevel, &t.GatherTier,
		); err != nil {
			return nil, fmt.Errorf("scan item template: %w", err)
		}
		t.BindKind = BindKind(bind)
		t.StatBonuses = bonuses
		if len(useEffect) > 0 {
			t.UseEffect = new(UseEffect)
			if err := json.Unmarshal(useEffect, t.UseEffect); err != nil {
				return nil, fmt.Errorf("item %s: bad use_effect: %w", t.Code, err)
			}
		}
		list = append(list, t)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if len(list) == 0 {
		return nil, fmt.Errorf("empty item catalog — did 0007_item migration run?")
	}
	return NewCatalog(list), nil
}
