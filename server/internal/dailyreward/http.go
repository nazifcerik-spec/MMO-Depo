package dailyreward

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/otter/server/internal/http/apierr"
	"github.com/otter/server/internal/httpctx"
)

type Handler struct{ svc *Service }

func NewHandler(svc *Service) *Handler { return &Handler{svc: svc} }

// Status is GET /v1/character/{id}/daily-reward (ownership: router guard).
func (h *Handler) Status(w http.ResponseWriter, r *http.Request) {
	charID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		apierr.Write(w, r, http.StatusBadRequest, "E_INVALID_ID", "id")
		return
	}
	st, err := h.svc.Status(r.Context(), httpctx.AccountID(r.Context()), charID)
	if err != nil {
		writeErr(w, r, err)
		return
	}
	apierr.WriteJSON(w, http.StatusOK, st)
}

// Claim is POST /v1/character/{id}/daily-reward/claim: today's reward to
// this character, then the new status.
func (h *Handler) Claim(w http.ResponseWriter, r *http.Request) {
	charID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		apierr.Write(w, r, http.StatusBadRequest, "E_INVALID_ID", "id")
		return
	}
	acct := httpctx.AccountID(r.Context())
	rw, err := h.svc.Claim(r.Context(), acct, charID)
	if err != nil {
		writeErr(w, r, err)
		return
	}
	st, err := h.svc.Status(r.Context(), acct, charID)
	if err != nil {
		writeErr(w, r, err)
		return
	}
	apierr.WriteJSON(w, http.StatusOK, map[string]any{"reward": rw, "status": st})
}

func writeErr(w http.ResponseWriter, r *http.Request, err error) {
	switch {
	case errors.Is(err, ErrClaimed):
		apierr.Write(w, r, http.StatusConflict, "E_ALREADY_CLAIMED", "already claimed today")
	case errors.Is(err, ErrNotFound):
		apierr.Write(w, r, http.StatusNotFound, "E_NOT_FOUND", "character not found")
	default:
		apierr.Write(w, r, http.StatusInternalServerError, "E_INTERNAL", err.Error())
	}
}
