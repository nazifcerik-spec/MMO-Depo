package consumable

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

// Use is POST /v1/character/{id}/items/{item_id}/use. Ownership of {id} is
// enforced by the router guard; the service re-checks item ownership.
func (h *Handler) Use(w http.ResponseWriter, r *http.Request) {
	charID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		apierr.Write(w, r, http.StatusBadRequest, "E_INVALID_ID", "id")
		return
	}
	itemID, err := uuid.Parse(chi.URLParam(r, "item_id"))
	if err != nil {
		apierr.Write(w, r, http.StatusBadRequest, "E_INVALID_ID", "item_id")
		return
	}
	res, err := h.svc.Use(r.Context(), charID, httpctx.AccountID(r.Context()), itemID)
	if err != nil {
		writeErr(w, r, err)
		return
	}
	apierr.WriteJSON(w, http.StatusOK, res)
}

// Active is GET /v1/character/{id}/buff: the running food buff or null
// (DD-269). Ownership of {id} is enforced by the router guard.
func (h *Handler) Active(w http.ResponseWriter, r *http.Request) {
	charID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		apierr.Write(w, r, http.StatusBadRequest, "E_INVALID_ID", "id")
		return
	}
	b, err := h.svc.Active(r.Context(), charID)
	if err != nil {
		writeErr(w, r, err)
		return
	}
	apierr.WriteJSON(w, http.StatusOK, map[string]any{"buff": b})
}

func writeErr(w http.ResponseWriter, r *http.Request, err error) {
	switch {
	case errors.Is(err, ErrNotFound), errors.Is(err, character.ErrNotFound), errors.Is(err, character.ErrForbidden):
		apierr.Write(w, r, http.StatusNotFound, "E_NOT_FOUND", "item not found")
	case errors.Is(err, ErrLocked):
		apierr.Write(w, r, http.StatusConflict, "E_ITEM_LOCKED", "item is locked")
	case errors.Is(err, ErrNotUsable):
		apierr.Write(w, r, http.StatusBadRequest, "E_NOT_USABLE", "item cannot be used")
	case errors.Is(err, ErrLevel):
		apierr.Write(w, r, http.StatusForbidden, "E_LEVEL_REQUIREMENT", "character level too low")
	case errors.Is(err, ErrInCombat):
		apierr.Write(w, r, http.StatusConflict, "E_IN_COMBAT", "finish the hunt or dungeon first")
	case errors.Is(err, ErrHuntOnly):
		apierr.Write(w, r, http.StatusConflict, "E_HUNT_ONLY", "take it on a hunt as a supply")
	case errors.Is(err, ErrNoEffect):
		apierr.Write(w, r, http.StatusConflict, "E_HP_FULL", "hp is already full")
	default:
		apierr.Write(w, r, http.StatusInternalServerError, "E_INTERNAL", err.Error())
	}
}
