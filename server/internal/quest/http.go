package quest

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/otter/server/internal/domain/character"
	"github.com/otter/server/internal/http/apierr"
	"github.com/otter/server/internal/httpctx"
)

// Handler serves /v1/character/{id}/quests[/{code}/{accept|abandon|turn-in}].
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

var errStatus = []struct {
	err    error
	status int
	code   string
}{
	{character.ErrNotFound, http.StatusNotFound, "E_NOT_FOUND"},
	{character.ErrForbidden, http.StatusNotFound, "E_NOT_FOUND"},
	{ErrNotFound, http.StatusNotFound, "E_NOT_FOUND"},
	{ErrLevel, http.StatusForbidden, "E_LEVEL_REQUIREMENT"},
	{ErrLocked, http.StatusConflict, "E_QUEST_LOCKED"},
	{ErrTaken, http.StatusConflict, "E_QUEST_TAKEN"},
	{ErrNotOffered, http.StatusConflict, "E_QUEST_NOT_OFFERED"},
	{ErrOpenCap, http.StatusConflict, "E_QUEST_CAP"},
	{ErrNotActive, http.StatusConflict, "E_QUEST_NOT_ACTIVE"},
	{ErrIncomplete, http.StatusConflict, "E_QUEST_INCOMPLETE"},
	{ErrNotPlayable, http.StatusConflict, "E_NOT_PLAYABLE"},
}

func reply(w http.ResponseWriter, r *http.Request, v any, err error) {
	if err == nil {
		apierr.WriteJSON(w, http.StatusOK, v)
		return
	}
	for _, e := range errStatus {
		if errors.Is(err, e.err) {
			apierr.Write(w, r, e.status, e.code, err.Error())
			return
		}
	}
	apierr.Write(w, r, http.StatusInternalServerError, "E_INTERNAL", err.Error())
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	char, acct, ok := scope(w, r)
	if !ok {
		return
	}
	b, err := h.svc.List(r.Context(), char, acct)
	reply(w, r, b, err)
}

func (h *Handler) Accept(w http.ResponseWriter, r *http.Request) {
	char, acct, ok := scope(w, r)
	if !ok {
		return
	}
	v, err := h.svc.Accept(r.Context(), char, acct, chi.URLParam(r, "code"))
	reply(w, r, v, err)
}

func (h *Handler) Abandon(w http.ResponseWriter, r *http.Request) {
	char, acct, ok := scope(w, r)
	if !ok {
		return
	}
	err := h.svc.Abandon(r.Context(), char, acct, chi.URLParam(r, "code"))
	reply(w, r, map[string]any{"abandoned": true}, err)
}

func (h *Handler) TurnIn(w http.ResponseWriter, r *http.Request) {
	char, acct, ok := scope(w, r)
	if !ok {
		return
	}
	res, err := h.svc.TurnIn(r.Context(), char, acct, chi.URLParam(r, "code"))
	reply(w, r, res, err)
}
