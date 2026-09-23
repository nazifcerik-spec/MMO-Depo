package progression

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

// Handler serves /v1/character/{id}/{progression,achievements,titles}.
type Handler struct{ svc *Service }

func NewHandler(svc *Service) *Handler { return &Handler{svc: svc} }

func scope(w http.ResponseWriter, r *http.Request) (uuid.UUID, uuid.UUID, bool) {
	acct := httpctx.AccountID(r.Context())
	if acct == uuid.Nil {
		apierr.Write(w, r, http.StatusUnauthorized, "E_UNAUTHENTICATED", "auth required")
		return uuid.Nil, uuid.Nil, false
	}
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		apierr.Write(w, r, http.StatusBadRequest, "E_INVALID_ID", err.Error())
		return uuid.Nil, uuid.Nil, false
	}
	return id, acct, true
}

func reply(w http.ResponseWriter, r *http.Request, v any, err error) {
	switch {
	case err == nil:
		apierr.WriteJSON(w, http.StatusOK, v)
	case errors.Is(err, character.ErrNotFound), errors.Is(err, character.ErrForbidden):
		apierr.Write(w, r, http.StatusNotFound, "E_NOT_FOUND", "not found")
	case errors.Is(err, ErrTitleNotOwned):
		apierr.Write(w, r, http.StatusConflict, "E_TITLE_NOT_OWNED", err.Error())
	default:
		apierr.Write(w, r, http.StatusInternalServerError, "E_INTERNAL", err.Error())
	}
}

func (h *Handler) Summary(w http.ResponseWriter, r *http.Request) {
	char, acct, ok := scope(w, r)
	if !ok {
		return
	}
	s, err := h.svc.Summary(r.Context(), char, acct)
	reply(w, r, s, err)
}

func (h *Handler) Achievements(w http.ResponseWriter, r *http.Request) {
	char, acct, ok := scope(w, r)
	if !ok {
		return
	}
	list, err := h.svc.Achievements(r.Context(), char, acct)
	reply(w, r, map[string]any{"count": len(list), "achievements": list}, err)
}

func (h *Handler) Titles(w http.ResponseWriter, r *http.Request) {
	char, acct, ok := scope(w, r)
	if !ok {
		return
	}
	list, err := h.svc.Titles(r.Context(), char, acct)
	reply(w, r, map[string]any{"count": len(list), "titles": list}, err)
}

func (h *Handler) SetActiveTitle(w http.ResponseWriter, r *http.Request) {
	char, acct, ok := scope(w, r)
	if !ok {
		return
	}
	var b struct {
		TitleCode string `json:"title_code"`
	}
	if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
		apierr.Write(w, r, http.StatusBadRequest, "E_INVALID_JSON", err.Error())
		return
	}
	err := h.svc.SetActiveTitle(r.Context(), char, acct, b.TitleCode)
	reply(w, r, map[string]any{"active_title_code": b.TitleCode}, err)
}
