package worldevent

import (
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

func reply(w http.ResponseWriter, r *http.Request, v any, err error) {
	type m struct {
		err    error
		status int
		code   string
	}
	if err == nil {
		apierr.WriteJSON(w, http.StatusOK, v)
		return
	}
	for _, c := range []m{
		{ErrNotFound, 404, "E_NOT_FOUND"},
		{character.ErrNotFound, 404, "E_NOT_FOUND"},
		{character.ErrForbidden, 404, "E_NOT_FOUND"},
		{ErrNotOpen, 409, "E_EVENT_NOT_OPEN"},
		{ErrLevelRequirement, 403, "E_LEVEL_REQUIREMENT"},
		{ErrWrongRegion, 409, "E_WRONG_REGION"},
		{ErrAlreadyAttending, 409, "E_ALREADY_ENLISTED"},
		{ErrEventNotActive, 409, "E_EVENT_NOT_ACTIVE"},
		{ErrAlreadyClaimed, 409, "E_ALREADY_CLAIMED"},
	} {
		if errors.Is(err, c.err) {
			apierr.Write(w, r, c.status, c.code, err.Error())
			return
		}
	}
	apierr.Write(w, r, http.StatusInternalServerError, "E_INTERNAL", err.Error())
}

func ids(w http.ResponseWriter, r *http.Request, keys ...string) ([]uuid.UUID, uuid.UUID, bool) {
	acct := httpctx.AccountID(r.Context())
	if acct == uuid.Nil {
		apierr.Write(w, r, http.StatusUnauthorized, "E_UNAUTHENTICATED", "auth required")
		return nil, uuid.Nil, false
	}
	out := make([]uuid.UUID, 0, len(keys))
	for _, k := range keys {
		id, err := uuid.Parse(chi.URLParam(r, k))
		if err != nil {
			apierr.Write(w, r, http.StatusBadRequest, "E_INVALID_ID", k+": "+err.Error())
			return nil, uuid.Nil, false
		}
		out = append(out, id)
	}
	return out, acct, true
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	p, acct, ok := ids(w, r, "id")
	if !ok {
		return
	}
	list, err := h.svc.List(r.Context(), p[0], acct)
	reply(w, r, map[string]any{"count": len(list), "events": list}, err)
}

func (h *Handler) Attend(w http.ResponseWriter, r *http.Request) {
	p, acct, ok := ids(w, r, "id", "spawn_id")
	if !ok {
		return
	}
	v, err := h.svc.Attend(r.Context(), p[0], acct, p[1])
	reply(w, r, v, err)
}

func (h *Handler) Seasonal(w http.ResponseWriter, r *http.Request) {
	_, acct, ok := ids(w, r)
	if !ok {
		return
	}
	list, err := h.svc.Seasonal(r.Context(), acct)
	reply(w, r, map[string]any{"count": len(list), "events": list}, err)
}

func (h *Handler) ClaimDaily(w http.ResponseWriter, r *http.Request) {
	_, acct, ok := ids(w, r)
	if !ok {
		return
	}
	n, err := h.svc.ClaimDaily(r.Context(), acct, chi.URLParam(r, "code"))
	reply(w, r, map[string]any{"event_tokens": n}, err)
}
