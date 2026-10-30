package pet

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/otter/server/internal/domain/character"
	"github.com/otter/server/internal/http/apierr"
	"github.com/otter/server/internal/httpctx"
)

type Handler struct{ svc *Service }

func NewHandler(svc *Service) *Handler { return &Handler{svc: svc} }

// Routes are mounted under /v1/character/{id} (ownership: router guard).
func (h *Handler) Routes(r chi.Router, idem func(http.Handler) http.Handler) {
	r.Get("/pets", h.list)
	r.With(idem).Post("/pets/hatch", h.hatch)
	r.With(idem).Post("/pets/{pet_id}/active", h.active)
	r.With(idem).Post("/pets/{pet_id}/hunt", h.hunt)
	r.With(idem).Post("/pets/{pet_id}/collect", h.collect)
}

// Species is the public species list (GET /v1/pets).
func (h *Handler) Species(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Cache-Control", "public, max-age=300")
	apierr.WriteJSON(w, http.StatusOK, map[string]any{"species": Current().All(), "max_level": MaxLevel,
		"power_per_level": PowerPerLevel, "mob_reach": MobReach})
}

func ids(w http.ResponseWriter, r *http.Request, param string) (char, account, id uuid.UUID, ok bool) {
	char, err := uuid.Parse(chi.URLParam(r, "id"))
	if err == nil && param != "" {
		id, err = uuid.Parse(chi.URLParam(r, param))
	}
	if err != nil {
		apierr.Write(w, r, http.StatusBadRequest, "E_INVALID_ID", "id")
		return
	}
	return char, httpctx.AccountID(r.Context()), id, true
}

func (h *Handler) list(w http.ResponseWriter, r *http.Request) {
	char, acct, _, ok := ids(w, r, "")
	if !ok {
		return
	}
	pets, eggs, err := h.svc.List(r.Context(), char, acct)
	if err != nil {
		writeErr(w, r, err)
		return
	}
	apierr.WriteJSON(w, http.StatusOK, map[string]any{"pets": pets, "eggs": eggs, "max_pets": MaxPets})
}

func (h *Handler) hatch(w http.ResponseWriter, r *http.Request) {
	char, acct, _, ok := ids(w, r, "")
	if !ok {
		return
	}
	var body struct {
		ItemID uuid.UUID `json:"item_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		apierr.Write(w, r, http.StatusBadRequest, "E_VALIDATION", "item_id")
		return
	}
	p, err := h.svc.Hatch(r.Context(), char, acct, body.ItemID)
	if err != nil {
		writeErr(w, r, err)
		return
	}
	apierr.WriteJSON(w, http.StatusOK, map[string]any{"pet": p})
}

func (h *Handler) active(w http.ResponseWriter, r *http.Request) {
	char, acct, id, ok := ids(w, r, "pet_id")
	if !ok {
		return
	}
	var body struct {
		Active bool `json:"active"`
	}
	_ = json.NewDecoder(r.Body).Decode(&body)
	if err := h.svc.SetActive(r.Context(), char, acct, id, body.Active); err != nil {
		writeErr(w, r, err)
		return
	}
	h.list(w, r)
}

func (h *Handler) hunt(w http.ResponseWriter, r *http.Request) {
	char, acct, id, ok := ids(w, r, "pet_id")
	if !ok {
		return
	}
	var body struct {
		Mob             string `json:"mob_code"`
		DurationSeconds int32  `json:"duration_seconds"`
	}
	_ = json.NewDecoder(r.Body).Decode(&body)
	p, err := h.svc.StartHunt(r.Context(), char, acct, id, body.Mob, body.DurationSeconds)
	if err != nil {
		writeErr(w, r, err)
		return
	}
	apierr.WriteJSON(w, http.StatusOK, map[string]any{"pet": p})
}

func (h *Handler) collect(w http.ResponseWriter, r *http.Request) {
	char, acct, id, ok := ids(w, r, "pet_id")
	if !ok {
		return
	}
	res, err := h.svc.CollectHunt(r.Context(), char, acct, id)
	if err != nil {
		writeErr(w, r, err)
		return
	}
	apierr.WriteJSON(w, http.StatusOK, res)
}

func writeErr(w http.ResponseWriter, r *http.Request, err error) {
	for _, m := range []struct {
		err    error
		status int
		code   string
	}{
		{ErrNotFound, 404, "E_NOT_FOUND"},
		{character.ErrNotFound, 404, "E_NOT_FOUND"},
		{character.ErrForbidden, 404, "E_NOT_FOUND"},
		{ErrNotEgg, 409, "E_NOT_EGG"},
		{ErrCap, 409, "E_PET_CAP"},
		{ErrHatching, 409, "E_PET_HATCHING"},
		{ErrHunting, 409, "E_PET_HUNTING"},
		{ErrIdle, 409, "E_PET_IDLE"},
		{ErrValidation, 400, "E_VALIDATION"},
	} {
		if errors.Is(err, m.err) {
			apierr.Write(w, r, m.status, m.code, err.Error())
			return
		}
	}
	apierr.Write(w, r, http.StatusInternalServerError, "E_INTERNAL", err.Error())
}
