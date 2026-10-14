// Package gathering backs the gathering professions (DD-244): each activity
// kind is a profession levelled 1..100 by what it gathers. Levels live in
// crafting.character_profession (created on first XP); the resource tables
// are the item catalog's gather nodes.
package gathering

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/otter/server/internal/activity"
	"github.com/otter/server/internal/crafting"
	"github.com/otter/server/internal/http/apierr"
	"github.com/otter/server/internal/item"
)

// MaxLevel caps every gathering profession (matches gather_level's range).
const MaxLevel = 100

// Kinds is the display order of the gathering professions.
var Kinds = []activity.Kind{
	activity.KindMining, activity.KindHerbalism, activity.KindWoodcutting,
	activity.KindFishing, activity.KindHunting,
}

// Levels is the slice of crafting.PgRepo used here.
type Levels interface {
	Get(ctx context.Context, characterID uuid.UUID, professionCode string) (*crafting.CharacterProfession, error)
	List(ctx context.Context, characterID uuid.UUID) ([]crafting.CharacterProfession, error)
	UpsertLearn(ctx context.Context, characterID uuid.UUID, professionCode string, now time.Time) (*crafting.CharacterProfession, error)
	AddXP(ctx context.Context, characterID uuid.UUID, professionCode string, delta int32, maxLevel int16) (*crafting.CharacterProfession, bool, error)
}

// Catalog yields the live item catalog (hot-reloaded by the content pack).
type Catalog interface {
	Catalog() *item.Catalog
}

// Service implements activity.Gathering.
type Service struct {
	items  Catalog
	levels Levels
	now    func() time.Time
}

func NewService(items Catalog, levels Levels) *Service {
	return &Service{items: items, levels: levels, now: time.Now}
}

var _ activity.Gathering = (*Service)(nil)

func (s *Service) Nodes(kind activity.Kind) []activity.GatherNode {
	tmpls := s.items.Catalog().Nodes(string(kind))
	out := make([]activity.GatherNode, len(tmpls))
	for i, t := range tmpls {
		out[i] = activity.GatherNode{Code: t.Code, Level: int(t.GatherLevel), Tier: int(t.GatherTier)}
	}
	return out
}

// Level is the character's level in kind; 1 before the first haul.
func (s *Service) Level(ctx context.Context, characterID uuid.UUID, kind activity.Kind) (int, error) {
	cp, err := s.levels.Get(ctx, characterID, string(kind))
	if errors.Is(err, crafting.ErrNotFound) {
		return 1, nil
	}
	if err != nil {
		return 0, err
	}
	return int(cp.Level), nil
}

func (s *Service) AddXP(ctx context.Context, characterID uuid.UUID, kind activity.Kind, xp int) (activity.GatherProgress, error) {
	if _, err := s.levels.UpsertLearn(ctx, characterID, string(kind), s.now().UTC()); err != nil {
		return activity.GatherProgress{}, err
	}
	cp, leveled, err := s.levels.AddXP(ctx, characterID, string(kind), int32(xp), MaxLevel)
	if err != nil {
		return activity.GatherProgress{}, err
	}
	return activity.GatherProgress{
		Level: int(cp.Level), XP: int(cp.XP), XPToNext: xpToNext(cp.Level),
		Gained: xp, LeveledUp: leveled,
	}, nil
}

func xpToNext(level int16) int {
	if level >= MaxLevel {
		return 0
	}
	return int(crafting.ProfessionXPToNext(level))
}

// Handler serves GET /v1/character/{id}/gathering (ownership checked by the
// route group): every gathering profession's level plus its resources.
type Handler struct{ svc *Service }

func NewHandler(svc *Service) *Handler { return &Handler{svc: svc} }

func (h *Handler) Overview(w http.ResponseWriter, r *http.Request) {
	charID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		apierr.Write(w, r, http.StatusBadRequest, "E_INVALID_ID", "id: invalid uuid")
		return
	}
	rows, err := h.svc.levels.List(r.Context(), charID)
	if err != nil {
		apierr.Write(w, r, http.StatusInternalServerError, "E_INTERNAL", err.Error())
		return
	}
	have := make(map[string]crafting.CharacterProfession, len(rows))
	for _, cp := range rows {
		have[cp.ProfessionCode] = cp
	}
	out := make([]map[string]any, 0, len(Kinds))
	for _, k := range Kinds {
		cp, ok := have[string(k)]
		if !ok {
			cp = crafting.CharacterProfession{Level: 1}
		}
		out = append(out, map[string]any{
			"kind": k, "level": cp.Level, "xp": cp.XP, "xp_to_next": xpToNext(cp.Level),
			"nodes": h.svc.Nodes(k),
		})
	}
	apierr.WriteJSON(w, http.StatusOK, map[string]any{"professions": out})
}
