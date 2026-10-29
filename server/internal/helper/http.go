package helper

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
	r.Get("/helpers", h.list)
	r.With(idem).Post("/helpers/{type}/hire", h.hire)
	r.With(idem).Post("/helpers/{type}/start", h.start)
	r.With(idem).Post("/helpers/{type}/collect", h.collect)
	r.With(idem).Post("/helpers/{type}/cancel", h.cancel)
}

func ids(w http.ResponseWriter, r *http.Request) (char, account uuid.UUID, ok bool) {
	char, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		apierr.Write(w, r, http.StatusBadRequest, "E_INVALID_ID", "id")
		return
	}
	return char, httpctx.AccountID(r.Context()), true
}

func (h *Handler) list(w http.ResponseWriter, r *http.Request) {
	char, acct, ok := ids(w, r)
	if !ok {
		return
	}
	v, err := h.svc.List(r.Context(), char, acct)
	if err != nil {
		writeErr(w, r, err)
		return
	}
	apierr.WriteJSON(w, http.StatusOK, map[string]any{"helpers": v, "half_pct": HalfPct})
}

func (h *Handler) hire(w http.ResponseWriter, r *http.Request) {
	if char, acct, ok := ids(w, r); ok {
		if err := h.svc.Hire(r.Context(), char, acct, chi.URLParam(r, "type")); err != nil {
			writeErr(w, r, err)
			return
		}
		h.list(w, r)
	}
}

func (h *Handler) start(w http.ResponseWriter, r *http.Request) {
	char, acct, ok := ids(w, r)
	if !ok {
		return
	}
	var body struct {
		Target          string `json:"target"`
		DurationSeconds int32  `json:"duration_seconds"`
		Recipe          string `json:"recipe_code"`
		Count           int    `json:"count"`
	}
	_ = json.NewDecoder(r.Body).Decode(&body)
	code := chi.URLParam(r, "type")
	t, _ := TypeOf(code)
	var err error
	if t.Kind == KindCraft {
		_, err = h.svc.StartCraft(r.Context(), char, acct, code, body.Recipe, body.Count)
	} else {
		_, err = h.svc.StartGather(r.Context(), char, acct, code, body.Target, body.DurationSeconds)
	}
	if err != nil {
		writeErr(w, r, err)
		return
	}
	h.list(w, r)
}

func (h *Handler) collect(w http.ResponseWriter, r *http.Request) {
	if char, acct, ok := ids(w, r); ok {
		res, err := h.svc.Collect(r.Context(), char, acct, chi.URLParam(r, "type"))
		if err != nil {
			writeErr(w, r, err)
			return
		}
		apierr.WriteJSON(w, http.StatusOK, map[string]any{"result": res})
	}
}

func (h *Handler) cancel(w http.ResponseWriter, r *http.Request) {
	if char, acct, ok := ids(w, r); ok {
		if err := h.svc.Cancel(r.Context(), char, acct, chi.URLParam(r, "type")); err != nil {
			writeErr(w, r, err)
			return
		}
		h.list(w, r)
	}
}

func writeErr(w http.ResponseWriter, r *http.Request, err error) {
	for _, m := range []struct {
		err    error
		status int
		code   string
	}{
		{ErrUnknownType, 404, "E_NOT_FOUND"},
		{character.ErrNotFound, 404, "E_NOT_FOUND"},
		{character.ErrForbidden, 404, "E_NOT_FOUND"},
		{ErrNotOwned, 409, "E_HELPER_NOT_HIRED"},
		{ErrOwned, 409, "E_HELPER_HIRED"},
		{ErrBusy, 409, "E_HELPER_BUSY"},
		{ErrIdle, 409, "E_HELPER_IDLE"},
		{ErrNotReady, 409, "E_HELPER_NOT_READY"},
		{ErrValidation, 400, "E_VALIDATION"},
		{ErrNoGold, 409, "E_INSUFFICIENT_GOLD"},
		{ErrSkill, 409, "E_SKILL_TOO_LOW"},
		{ErrMaterials, 409, "E_INSUFFICIENT_MATERIALS"},
	} {
		if errors.Is(err, m.err) {
			apierr.Write(w, r, m.status, m.code, err.Error())
			return
		}
	}
	apierr.Write(w, r, http.StatusInternalServerError, "E_INTERNAL", err.Error())
}
