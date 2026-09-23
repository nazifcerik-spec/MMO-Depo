package item

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Catalog is the in-memory template lookup — loaded once at boot from
// content.item_template rows. Sprint 12 YAML loader replaces the DB read with
// a versioned content pack hot-reload.
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

// Get returns a template by code; ErrUnknownTemplate on miss.
func (c *Catalog) Get(code string) (Template, error) {
	t, ok := c.byCode[code]
	if !ok {
		return Template{}, fmt.Errorf("%w: %s", ErrUnknownTemplate, code)
	}
	return t, nil
}

// Has returns true when the code is in the catalog.
func (c *Catalog) Has(code string) bool {
	_, ok := c.byCode[code]
	return ok
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
		       equip_slot, class_restriction, stat_bonuses
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
		if err := rows.Scan(
			&t.Code, &t.DisplayNameKey, &t.ItemClass, &t.ItemSubclass, &t.MinLevel,
			&t.BaseValueCopper, &t.BaseWeightG, &t.Stackable, &t.MaxStack, &bind,
			&t.NoTrade, &t.NoSell, &t.Affixable,
			&t.EquipSlot, &t.ClassRestriction, &bonuses,
		); err != nil {
			return nil, fmt.Errorf("scan item template: %w", err)
		}
		t.BindKind = BindKind(bind)
		t.StatBonuses = bonuses
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
