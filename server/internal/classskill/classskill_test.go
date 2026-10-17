package classskill

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/otter/server/internal/domain/character"
	"github.com/otter/server/internal/skill"
)

func catalog() []Skill {
	return []Skill{
		{Code: "bolt", Class: "mage", UnlockLevel: 1, Kind: "active", SortOrder: 20, Effect: "strike", Power: 200, MPCost: 10, CooldownTurns: 2},
		{Code: "blast", Class: "mage", UnlockLevel: 6, Kind: "active", SortOrder: 10, Effect: "strike", Power: 350, MPCost: 30, CooldownTurns: 5},
		{Code: "mind", Class: "mage", UnlockLevel: 3, Kind: "passive",
			Modifiers: []skill.Modifier{{Kind: skill.KindMultAdd, Target: skill.TargetMPMax, Amount: 0.2}}},
		{Code: "cleave", Class: "warrior", UnlockLevel: 1, Kind: "active", Effect: "strike", Power: 160, CooldownTurns: 3},
	}
}

func TestCatalog_UnlockAndPriority(t *testing.T) {
	c := NewCatalog(catalog())

	require.Empty(t, c.Passives([]string{"mage"}, 2))
	require.Len(t, c.Passives([]string{"mage"}, 3), 1)

	low := c.Actives([]string{"mage"}, 5)
	require.Len(t, low, 1)
	require.Equal(t, "bolt", low[0].Code)
	high := c.Actives([]string{"mage"}, 6)
	require.Equal(t, []string{"blast", "bolt"}, []string{high[0].Code, high[1].Code}, "sort_order is cast priority")
	require.Equal(t, 5, high[0].CooldownTurns)

	require.Len(t, c.ForClass("mage"), 3, "locked skills are listed too")
	require.Empty(t, c.Actives([]string{"bard"}, 60))
}

func TestService_SourcesAndList(t *testing.T) {
	s := NewService(nil)
	s.cat.Store(NewCatalog(catalog()))

	mods, err := s.PassiveModifiers(context.Background(), uuid.New(), character.Sheet{Class: "mage", Level: 10})
	require.NoError(t, err)
	require.Len(t, mods, 1)
	require.Len(t, s.ActiveSkills(&character.Sheet{Class: "warrior", Level: 1}), 1)
	require.Same(t, s.Catalog(), s.cat.Load())

	r := chi.NewRouter()
	r.Get("/classes/{class}/skills", s.List)
	for class, want := range map[string]string{"mage": `"code":"blast"`, "bard": `"skills":[]`} {
		w := httptest.NewRecorder()
		r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/classes/"+class+"/skills", nil))
		require.Equal(t, http.StatusOK, w.Code)
		require.Contains(t, w.Body.String(), want)
		require.NotContains(t, w.Body.String(), "sort_order")
	}
}

// DD-250: a promotion's skills belong to its class path only.
func TestCatalog_PathOwnsSkills(t *testing.T) {
	c := NewCatalog([]Skill{
		{Code: "w_strike", Class: "warrior", UnlockLevel: 1, Kind: "active", Effect: "strike", Power: 110},
		{Code: "g_wall", Class: "warrior", Spec: "guardian", UnlockLevel: 100, Kind: "active", Effect: "guard", Power: 30,
			Config: Config{SelfHPBelow: 40, DurationTurns: 3}},
		{Code: "e_rush", Class: "warrior", Spec: "executioner", UnlockLevel: 100, Kind: "active", Effect: "empower", Power: 20},
	})
	guardian := c.Actives([]string{"warrior", "guardian"}, 100)
	require.Len(t, guardian, 2)
	require.Equal(t, int32(40), guardian[1].SelfHPBelow)
	require.Equal(t, 3, guardian[1].DurationTurns)
	require.Len(t, c.Actives([]string{"warrior"}, 100), 1, "no promotion yet: base skills only")
	require.Len(t, c.Actives([]string{"warrior", "guardian"}, 99), 1, "promotion skills wait for their level")

	specs := map[string]string{}
	for _, u := range Unlocks(c.ForClass("warrior")) {
		specs[u.Code] = u.Spec
	}
	require.Equal(t, map[string]string{"w_strike": "warrior", "g_wall": "guardian", "e_rush": "executioner"}, specs)
}

func TestPathOfAndUnlockDefaults(t *testing.T) {
	spec := "guardian"
	// Without the progression tree installed the path falls back to the class.
	require.Equal(t, []string{"warrior"}, pathOf(&character.Sheet{Class: "warrior", Spec: &spec}))
	require.Equal(t, []string{"mage"}, pathOf(&character.Sheet{Class: "mage"}))
	require.Empty(t, NewCatalog(nil).Actives(nil, 10))

	u := Unlocks([]Skill{{Code: "p", Class: "mage", Kind: "passive", UnlockLevel: 3},
		{Code: "a", Class: "mage", Kind: "active", UnlockLevel: 1, Config: Config{Category: "ultimate"}}})
	require.Equal(t, "mage", u[0].Spec)
	require.Equal(t, "passive", u[0].Category)
	require.Equal(t, "ultimate", u[1].Category)
}
