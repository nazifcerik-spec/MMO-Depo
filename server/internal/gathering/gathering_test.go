package gathering

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/otter/server/internal/activity"
	"github.com/otter/server/internal/crafting"
	"github.com/otter/server/internal/item"
)

type memLevels struct {
	rows map[string]*crafting.CharacterProfession
}

func (m *memLevels) Get(_ context.Context, _ uuid.UUID, code string) (*crafting.CharacterProfession, error) {
	if cp, ok := m.rows[code]; ok {
		c := *cp
		return &c, nil
	}
	return nil, crafting.ErrNotFound
}

func (m *memLevels) List(context.Context, uuid.UUID) ([]crafting.CharacterProfession, error) {
	var out []crafting.CharacterProfession
	for _, cp := range m.rows {
		out = append(out, *cp)
	}
	return out, nil
}

func (m *memLevels) UpsertLearn(_ context.Context, _ uuid.UUID, code string, _ time.Time) (*crafting.CharacterProfession, error) {
	if _, ok := m.rows[code]; !ok {
		m.rows[code] = &crafting.CharacterProfession{ProfessionCode: code, Level: 1}
	}
	return m.rows[code], nil
}

func (m *memLevels) AddXP(_ context.Context, _ uuid.UUID, code string, delta int32, maxLevel int16) (*crafting.CharacterProfession, bool, error) {
	cp := m.rows[code]
	cp.XP += delta
	up := false
	for cp.Level < maxLevel && cp.XP >= crafting.ProfessionXPToNext(cp.Level) {
		cp.XP -= crafting.ProfessionXPToNext(cp.Level)
		cp.Level++
		up = true
	}
	return cp, up, nil
}

type staticCat struct{ c *item.Catalog }

func (s staticCat) Catalog() *item.Catalog { return s.c }

func newSvc() (*Service, *memLevels) {
	cat := item.NewCatalog([]item.Template{
		{Code: "iron_ore", GatherKind: "mining", GatherLevel: 16, GatherTier: 1},
		{Code: "raw_stone", GatherKind: "mining", GatherLevel: 1, GatherTier: 0},
	})
	lv := &memLevels{rows: map[string]*crafting.CharacterProfession{}}
	return NewService(staticCat{cat}, lv), lv
}

func TestLevelsStartAtOneAndGrowWithXP(t *testing.T) {
	svc, _ := newSvc()
	ctx, char := context.Background(), uuid.New()
	lvl, err := svc.Level(ctx, char, activity.KindMining)
	require.NoError(t, err)
	require.Equal(t, 1, lvl, "no row yet")

	p, err := svc.AddXP(ctx, char, activity.KindMining, 165) // 76 + 81 = 157 (DD-309)
	require.NoError(t, err)
	require.Equal(t, activity.GatherProgress{Level: 3, XP: 8, XPToNext: 86, Gained: 165, LeveledUp: true}, p)
	lvl, _ = svc.Level(ctx, char, activity.KindMining)
	require.Equal(t, 3, lvl)

	require.Equal(t, []activity.GatherNode{{Code: "raw_stone", Level: 1}, {Code: "iron_ore", Level: 16, Tier: 1}},
		svc.Nodes(activity.KindMining))
	require.Zero(t, xpToNext(MaxLevel), "nothing to gain at the cap")
}

func TestOverviewListsEveryProfession(t *testing.T) {
	svc, lv := newSvc()
	char := uuid.New()
	lv.rows["mining"] = &crafting.CharacterProfession{ProfessionCode: "mining", Level: 7, XP: 5}
	r := chi.NewRouter()
	r.Get("/c/{id}", NewHandler(svc).Overview)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/c/"+char.String(), nil))
	require.Equal(t, http.StatusOK, rec.Code)
	var body struct {
		Professions []struct {
			Kind  string                `json:"kind"`
			Level int                   `json:"level"`
			Nodes []activity.GatherNode `json:"nodes"`
		} `json:"professions"`
	}
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &body))
	require.Len(t, body.Professions, len(Kinds))
	require.Equal(t, "mining", body.Professions[0].Kind)
	require.Equal(t, 7, body.Professions[0].Level)
	require.Len(t, body.Professions[0].Nodes, 2)
	require.Equal(t, 1, body.Professions[1].Level, "untouched professions start at 1")

	rec = httptest.NewRecorder()
	r.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/c/nope", nil))
	require.Equal(t, http.StatusBadRequest, rec.Code)
}
