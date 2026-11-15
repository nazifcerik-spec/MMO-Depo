package story

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

// Handler serves /v1/character/{id}/story[/{scene}/advance].
type Handler struct{ svc *Service }

func NewHandler(svc *Service) *Handler { return &Handler{svc: svc} }

var errStatus = []struct {
	err    error
	status int
	code   string
}{
	{character.ErrNotFound, http.StatusNotFound, "E_NOT_FOUND"},
	{character.ErrForbidden, http.StatusNotFound, "E_NOT_FOUND"},
	{ErrNotCurrent, http.StatusConflict, "E_STORY_NOT_CURRENT"},
	{ErrLevel, http.StatusForbidden, "E_LEVEL_REQUIREMENT"},
	{ErrQuests, http.StatusConflict, "E_STORY_QUESTS"},
	{ErrOption, http.StatusBadRequest, "E_STORY_OPTION"},
	{ErrLocked, http.StatusConflict, "E_STORY_OPTION_LOCKED"},
	{ErrChosen, http.StatusConflict, "E_STORY_CHOSEN"},
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

func (h *Handler) Journal(w http.ResponseWriter, r *http.Request) {
	char, acct, ok := scope(w, r)
	if !ok {
		return
	}
	j, err := h.svc.Journal(r.Context(), char, acct)
	reply(w, r, j, err)
}

func (h *Handler) Advance(w http.ResponseWriter, r *http.Request) {
	char, acct, ok := scope(w, r)
	if !ok {
		return
	}
	var body struct {
		Option string `json:"option"`
	}
	if r.ContentLength != 0 {
		if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 1024)).Decode(&body); err != nil {
			apierr.Write(w, r, http.StatusBadRequest, "E_VALIDATION", "bad body")
			return
		}
	}
	j, err := h.svc.Advance(r.Context(), char, acct, chi.URLParam(r, "scene"), body.Option)
	reply(w, r, j, err)
}
