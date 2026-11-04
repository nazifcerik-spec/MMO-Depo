// Package classskill serves class skills (Sprint 42, DD-222; DD-250):
// passives feed the sheet modifier stack, actives are auto-cast in combat.
// What a character has follows from its class path (base class → promotions)
// and level alone — no per-character state. Of the actives that do the same
// job only the strongest stays (DD-352).
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
	"github.com/otter/server/internal/leveling"
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

	// DD-250: owning class node (base class or promotion) and data-driven
	// activation; the same definition serves manual and AFK play.
	Spec   string `json:"spec"`
	Config Config `json:"config"`

	// DD-352: what an active does and how strongly; set by NewCatalog.
	Job      string  `json:"job,omitempty"`
	Strength float64 `json:"strength,omitempty"`
}

// job is what an active skill does: its effect, kind, trigger and cooldown
// band. A basic attack, a regular skill and a heavy hit are different jobs, so
// a character keeps a rotation; a finisher is not an everyday strike.
// Awakening and ultimate skills are each their own job.
func (s Skill) job() string {
	if s.Config.Category == "awakening" || s.Config.Category == "ultimate" {
		return s.Code
	}
	trigger := "ready"
	switch {
	case s.Config.TargetHPBelow > 0:
		trigger = "target_hp"
	case s.Config.SelfHPBelow > 0 || s.Config.AllyHPBelow > 0:
		trigger = "hp"
	}
	band := "heavy"
	switch {
	case s.CooldownTurns <= 1:
		band = "basic"
	case s.CooldownTurns <= 4:
		band = "regular"
	case s.CooldownTurns >= 10:
		band = "long"
	}
	return s.Effect + "/" + s.Config.Type + "/" + trigger + "/" + band
}

// strength compares skills of one job: power over its whole duration.
func (s Skill) strength() float64 {
	return float64(s.Power) * float64(max(s.Config.DurationTurns, 1))
}

// strongest keeps, per job, the strongest active (the newest on a tie).
// Passives and order are untouched.
func strongest(list []Skill) []Skill {
	best := map[string]Skill{}
	for _, s := range list {
		if s.Kind != "active" {
			continue
		}
		b, seen := best[s.Job]
		if !seen || s.Strength > b.Strength || (s.Strength == b.Strength && s.UnlockLevel > b.UnlockLevel) {
			best[s.Job] = s
		}
	}
	out := make([]Skill, 0, len(list))
	for _, s := range list {
		if s.Kind != "active" || best[s.Job].Code == s.Code {
			out = append(out, s)
		}
	}
	return out
}

// Config is a skill's activation data (content.class_skill.config).
type Config struct {
	Category       string   `json:"category,omitempty"`        // skill | passive | awakening | ultimate
	ActivationMode string   `json:"activation_mode,omitempty"` // active | stance | passive
	SelfHPBelow    int32    `json:"self_hp_below,omitempty"`   // auto-trigger: own HP% <=
	TargetHPBelow  int32    `json:"target_hp_below,omitempty"` // auto-trigger: target HP% <=
	ResourceGain   int32    `json:"resource_gain,omitempty"`
	DurationTurns  int      `json:"duration_turns,omitempty"`
	Requires       []string `json:"requires,omitempty"` // equipment tags (e.g. shield)

	// Shared data model (DD-252).
	OnCooldown  bool            `json:"on_cooldown,omitempty"`   // auto-trigger: whenever ready
	AllyHPBelow int32           `json:"ally_hp_below,omitempty"` // auto-trigger: an ally's HP% <= (solo: self)
	Type        string          `json:"type,omitempty"`
	DamageType  string          `json:"damage_type,omitempty"`
	TargetType  string          `json:"target_type,omitempty"`
	Range       string          `json:"range,omitempty"`
	CastTime    int             `json:"cast_time,omitempty"`
	Tags        []string        `json:"tags,omitempty"`
	Passive     json.RawMessage `json:"passive,omitempty"`
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
		if s.Spec == "" {
			s.Spec = s.Class
		}
		if s.Kind == "active" {
			s.Job, s.Strength = s.job(), s.strength()
		}
		c.byClass[s.Class] = append(c.byClass[s.Class], s)
	}
	return c
}

// ForClass lists every skill of a base class (all branches, locked ones included).
func (c *Catalog) ForClass(class string) []Skill { return c.byClass[class] }

// owned: skills of the base class whose owning node is on the path.
func (c *Catalog) owned(path []string) []Skill {
	if len(path) == 0 {
		return nil
	}
	var out []Skill
	for _, s := range c.byClass[path[0]] {
		for _, p := range path {
			if s.Spec == p {
				out = append(out, s)
				break
			}
		}
	}
	return out
}

// Passives returns the modifiers of the passives the path has unlocked at level.
func (c *Catalog) Passives(path []string, level int) []skill.Modifier {
	var out []skill.Modifier
	for _, s := range c.owned(path) {
		if s.Kind == "passive" && int(s.UnlockLevel) <= level {
			out = append(out, s.Modifiers...)
		}
	}
	return out
}

// Actives returns the path's unlocked active skills in cast-priority order.
func (c *Catalog) Actives(path []string, level int) []combat.Skill {
	var unlocked []Skill
	for _, s := range c.owned(path) {
		if s.Kind == "active" && int(s.UnlockLevel) <= level {
			unlocked = append(unlocked, s)
		}
	}
	var out []combat.Skill
	for _, s := range strongest(unlocked) {
		// Party healing is not simulated: alone, the character is the
		// ally its ally-HP rule watches.
		self := max(s.Config.SelfHPBelow, s.Config.AllyHPBelow)
		out = append(out, combat.Skill{Code: s.Code, Effect: s.Effect, Power: s.Power, MPCost: s.MPCost,
			CooldownTurns: s.CooldownTurns, Gain: s.Config.ResourceGain, SelfHPBelow: self,
			TargetHPBelow: s.Config.TargetHPBelow, DurationTurns: s.Config.DurationTurns})
	}
	return out
}

// pathOf is a sheet's class path from the leveling engine.
func pathOf(sheet *character.Sheet) []string {
	spec := sheet.Class
	if sheet.Spec != nil && *sheet.Spec != "" {
		spec = *sheet.Spec
	}
	path := leveling.Current().Path(spec)
	if len(path) == 0 || path[0] != sheet.Class {
		return []string{sheet.Class}
	}
	return path
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
		       COALESCE(effect, ''), power, mp_cost, cooldown_turns, modifiers,
		       COALESCE(spec_code, ''), config
		  FROM content.class_skill`)
	if err != nil {
		return err
	}
	list, err := pgx.CollectRows(rows, func(r pgx.CollectableRow) (Skill, error) {
		var k Skill
		var mods, cfg []byte
		if err := r.Scan(&k.Code, &k.NameKey, &k.Class, &k.UnlockLevel, &k.Kind, &k.SortOrder,
			&k.Effect, &k.Power, &k.MPCost, &k.CooldownTurns, &mods, &k.Spec, &cfg); err != nil {
			return k, err
		}
		if err := json.Unmarshal(cfg, &k.Config); err != nil {
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
// sheet's class path and level.
func (s *Service) PassiveModifiers(_ context.Context, _ uuid.UUID, sheet character.Sheet) ([]skill.Modifier, error) {
	return s.cat.Load().Passives(pathOf(&sheet), sheet.Level), nil
}

// ActiveSkills gives combat the unlocked actives for a sheet (hunt/dungeon).
func (s *Service) ActiveSkills(sheet *character.Sheet) []combat.Skill {
	return s.cat.Load().Actives(pathOf(sheet), sheet.Level)
}

// Unlocks lists every skill as a leveling unlock (contentpack installs them
// into the progression engine for milestones and unlock events).
func Unlocks(list []Skill) []leveling.Unlock {
	out := make([]leveling.Unlock, 0, len(list))
	for _, s := range list {
		spec := s.Spec
		if spec == "" {
			spec = s.Class
		}
		cat := s.Config.Category
		if cat == "" {
			cat = map[bool]string{true: "passive", false: "skill"}[s.Kind == "passive"]
		}
		out = append(out, leveling.Unlock{Code: s.Code, Spec: spec, Level: int(s.UnlockLevel), Category: cat})
	}
	return out
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
