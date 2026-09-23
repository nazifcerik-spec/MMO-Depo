// Package classskill serves class skills v1 (Sprint 42, DD-222): passives
// feed the sheet modifier stack, actives are auto-cast in combat. What a
// character has follows from class + level alone — no per-character state.
package classskill

import (
	"context"
	"encoding/json"
	"net/http"
	"sort"
	"sync/atomic"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/otter/server/internal/combat"
	"github.com/otter/server/internal/domain/character"
	"github.com/otter/server/internal/http/apierr"
	"github.com/otter/server/internal/skill"
)

// Skill is one class skill definition.
type Skill struct {
	Code          string           `json:"code"`
	NameKey       string           `json:"display_name_key"`
	Class         string           `json:"class_code"`
	UnlockLevel   int16            `json:"unlock_level"`
	Kind          string           `json:"kind"` // active | passive
	SortOrder     int32            `json:"-"`
	Effect        string           `json:"effect,omitempty"`
	Power         int32            `json:"power,omitempty"`
	MPCost        int32            `json:"mp_cost,omitempty"`
	CooldownTurns int              `json:"cooldown_turns,omitempty"`
	Modifiers     []skill.Modifier `json:"modifiers,omitempty"`
}

// Catalog groups skills per class; actives in cast-priority order.
type Catalog struct{ byClass map[string][]Skill }

func NewCatalog(list []Skill) *Catalog {
	sort.SliceStable(list, func(i, j int) bool {
		if list[i].SortOrder != list[j].SortOrder {
			return list[i].SortOrder < list[j].SortOrder
		}
		if list[i].UnlockLevel != list[j].UnlockLevel {
			return list[i].UnlockLevel < list[j].UnlockLevel
		}
		return list[i].Code < list[j].Code
	})
	c := &Catalog{byClass: map[string][]Skill{}}
	for _, s := range list {
		c.byClass[s.Class] = append(c.byClass[s.Class], s)
	}
	return c
}

// ForClass lists every skill of a class (locked ones included).
func (c *Catalog) ForClass(class string) []Skill { return c.byClass[class] }

// Passives returns the modifiers of the passives unlocked at level.
func (c *Catalog) Passives(class string, level int) []skill.Modifier {
	var out []skill.Modifier
	for _, s := range c.byClass[class] {
		if s.Kind == "passive" && int(s.UnlockLevel) <= level {
			out = append(out, s.Modifiers...)
		}
	}
	return out
}

// Actives returns the unlocked active skills in cast-priority order.
func (c *Catalog) Actives(class string, level int) []combat.Skill {
	var out []combat.Skill
	for _, s := range c.byClass[class] {
		if s.Kind == "active" && int(s.UnlockLevel) <= level {
			out = append(out, combat.Skill{Code: s.Code, Effect: s.Effect, Power: s.Power, MPCost: s.MPCost,
				CooldownTurns: s.CooldownTurns})
		}
	}
	return out
}

// Service holds the hot-swappable catalog.
type Service struct {
	pool *pgxpool.Pool
	cat  atomic.Pointer[Catalog]
}

func NewService(pool *pgxpool.Pool) *Service {
	s := &Service{pool: pool}
	s.cat.Store(NewCatalog(nil))
	return s
}

func (s *Service) Catalog() *Catalog { return s.cat.Load() }

// Reload rebuilds the catalog from content.class_skill (contentpack.Reloader).
func (s *Service) Reload(ctx context.Context) error {
	rows, err := s.pool.Query(ctx, `
		SELECT code, display_name_key, class_code, unlock_level, kind, sort_order,
		       COALESCE(effect, ''), power, mp_cost, cooldown_turns, modifiers
		  FROM content.class_skill`)
	if err != nil {
		return err
	}
	list, err := pgx.CollectRows(rows, func(r pgx.CollectableRow) (Skill, error) {
		var k Skill
		var mods []byte
		if err := r.Scan(&k.Code, &k.NameKey, &k.Class, &k.UnlockLevel, &k.Kind, &k.SortOrder,
			&k.Effect, &k.Power, &k.MPCost, &k.CooldownTurns, &mods); err != nil {
			return k, err
		}
		return k, json.Unmarshal(mods, &k.Modifiers)
	})
	if err != nil {
		return err
	}
	s.cat.Store(NewCatalog(list))
	return nil
}

// PassiveModifiers is a skill.ModifierSource: unlocked passives for the
// sheet's class and level.
func (s *Service) PassiveModifiers(_ context.Context, _ uuid.UUID, sheet character.Sheet) ([]skill.Modifier, error) {
	return s.cat.Load().Passives(sheet.Class, sheet.Level), nil
}

// ActiveSkills gives combat the unlocked actives for a sheet (hunt/dungeon).
func (s *Service) ActiveSkills(sheet *character.Sheet) []combat.Skill {
	return s.cat.Load().Actives(sheet.Class, sheet.Level)
}

// List serves GET /v1/classes/{class}/skills (public; content only).
func (s *Service) List(w http.ResponseWriter, r *http.Request) {
	list := s.cat.Load().ForClass(chi.URLParam(r, "class"))
	if list == nil {
		list = []Skill{}
	}
	w.Header().Set("Cache-Control", "public, max-age=300")
	apierr.WriteJSON(w, http.StatusOK, map[string]any{"skills": list})
}
