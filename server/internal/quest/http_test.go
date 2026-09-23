package quest

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/otter/server/internal/domain/character"
	"github.com/otter/server/internal/httpctx"
)

type fakeChars struct {
	c   *character.Character
	err error
}

func (f fakeChars) GetOwned(context.Context, uuid.UUID, uuid.UUID) (*character.Character, error) {
	return f.c, f.err
}

func serve(t *testing.T, h *Handler, method, path string, account uuid.UUID) *httptest.ResponseRecorder {
	t.Helper()
	r := chi.NewRouter()
	r.Route("/character/{id}/quests", func(r chi.Router) {
		r.Get("/", h.List)
		r.Post("/{code}/accept", h.Accept)
		r.Post("/{code}/abandon", h.Abandon)
		r.Post("/{code}/turn-in", h.TurnIn)
	})
	req := httptest.NewRequest(method, path, nil)
	if account != uuid.Nil {
		req = req.WithContext(httpctx.WithAccountID(req.Context(), account))
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

// Everything below is decided before any DB access, so no pool is needed.
func TestHandler_GuardsAndErrorMapping(t *testing.T) {
	acct := uuid.New()
	char := "/character/" + uuid.NewString() + "/quests"
	foreign := NewHandler(NewService(nil, fakeChars{err: character.ErrForbidden}, nil, nil, nil))

	w := serve(t, foreign, http.MethodGet, char+"/", uuid.Nil)
	require.Equal(t, http.StatusUnauthorized, w.Code)
	w = serve(t, foreign, http.MethodGet, "/character/not-a-uuid/quests/", acct)
	require.Equal(t, http.StatusBadRequest, w.Code)
	for _, path := range []string{char + "/", char + "/q/accept", char + "/q/abandon", char + "/q/turn-in"} {
		method := http.MethodPost
		if path == char+"/" {
			method = http.MethodGet
		}
		w = serve(t, foreign, method, path, acct)
		require.Equal(t, http.StatusNotFound, w.Code, path)
		require.Contains(t, w.Body.String(), "E_NOT_FOUND")
	}

	low := &character.Character{Level: 1, State: character.StateActive}
	svc := NewService(nil, fakeChars{c: low}, nil, nil, nil)
	svc.cat.Store(NewCatalog([]Template{{Code: "big", Category: CategoryMain, MinLevel: 30,
		Objectives: []Objective{{Event: "mob_killed", Count: 1}}}}))
	h := NewHandler(svc)
	w = serve(t, h, http.MethodPost, char+"/nope/accept", acct)
	require.Equal(t, http.StatusNotFound, w.Code)
	w = serve(t, h, http.MethodPost, char+"/big/accept", acct)
	require.Equal(t, http.StatusForbidden, w.Code)
	require.Contains(t, w.Body.String(), "E_LEVEL_REQUIREMENT")
	w = serve(t, h, http.MethodPost, char+"/nope/turn-in", acct)
	require.Equal(t, http.StatusNotFound, w.Code)

	frozen := NewService(nil, fakeChars{c: &character.Character{Level: 60, State: "frozen"}}, nil, nil, nil)
	w = serve(t, NewHandler(frozen), http.MethodPost, char+"/big/accept", acct)
	require.Equal(t, http.StatusConflict, w.Code)
	require.Contains(t, w.Body.String(), "E_NOT_PLAYABLE")
}

func TestAccept_DailyMustBeOffered(t *testing.T) {
	svc := NewService(nil, fakeChars{c: &character.Character{Level: 1, State: character.StateActive}}, nil, nil, nil)
	svc.cat.Store(NewCatalog([]Template{{Code: "d_high", Category: CategoryDaily, MinLevel: 1}}))
	_, err := svc.Accept(context.Background(), uuid.New(), uuid.New(), "d_other")
	require.ErrorIs(t, err, ErrNotFound)
	svc.cat.Store(NewCatalog([]Template{
		{Code: "d1", Category: CategoryDaily, MinLevel: 1}, {Code: "d2", Category: CategoryDaily, MinLevel: 1},
		{Code: "d3", Category: CategoryDaily, MinLevel: 1}, {Code: "d4", Category: CategoryDaily, MinLevel: 1},
	}))
	char := uuid.New()
	offer := svc.cat.Load().DailyOffer(char, 1, svc.now())
	for _, code := range []string{"d1", "d2", "d3", "d4"} {
		if !contains(offer, code) {
			_, err := svc.Accept(context.Background(), char, uuid.New(), code)
			require.ErrorIs(t, err, ErrNotOffered)
			return
		}
	}
	t.Fatal("offer should leave one daily out")
}

func contains(list []string, s string) bool {
	for _, v := range list {
		if v == s {
			return true
		}
	}
	return false
}
