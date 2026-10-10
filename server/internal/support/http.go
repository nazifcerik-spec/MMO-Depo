package support

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/google/uuid"

	"github.com/otter/server/internal/http/apierr"
	"github.com/otter/server/internal/httpctx"
)

// Handler serves the player side: POST/GET /v1/support/tickets.
type Handler struct{ svc *Service }

func NewHandler(svc *Service) *Handler { return &Handler{svc: svc} }

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	acct := httpctx.AccountID(r.Context())
	if acct == uuid.Nil {
		apierr.Write(w, r, http.StatusUnauthorized, "E_UNAUTHENTICATED", "auth required")
		return
	}
	var req NewTicket
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 16<<10)).Decode(&req); err != nil {
		apierr.Write(w, r, http.StatusBadRequest, "E_INVALID_JSON", err.Error())
		return
	}
	t, err := h.svc.Create(r.Context(), acct, req)
	if err != nil {
		WriteErr(w, r, err)
		return
	}
	apierr.WriteJSON(w, http.StatusCreated, t)
}

func (h *Handler) Mine(w http.ResponseWriter, r *http.Request) {
	acct := httpctx.AccountID(r.Context())
	if acct == uuid.Nil {
		apierr.Write(w, r, http.StatusUnauthorized, "E_UNAUTHENTICATED", "auth required")
		return
	}
	list, err := h.svc.Mine(r.Context(), acct)
	if err != nil {
		WriteErr(w, r, err)
		return
	}
	apierr.WriteJSON(w, http.StatusOK, map[string]any{"tickets": list})
}

// WriteErr maps support errors (shared with the admin handlers).
func WriteErr(w http.ResponseWriter, r *http.Request, err error) {
	switch {
	case errors.Is(err, ErrInvalid):
		apierr.Write(w, r, http.StatusBadRequest, "E_VALIDATION",
			"category, a subject (3-120 characters) and a message (10-4000) are required")
	case errors.Is(err, ErrReplyRequired):
		apierr.Write(w, r, http.StatusBadRequest, "E_VALIDATION", "reply required (up to 4000 characters)")
	case errors.Is(err, ErrTooManyOpen):
		apierr.Write(w, r, http.StatusConflict, "E_TOO_MANY_TICKETS", "too many open tickets")
	case errors.Is(err, ErrForeignCharacter), errors.Is(err, ErrNotFound):
		apierr.Write(w, r, http.StatusNotFound, "E_NOT_FOUND", "not found")
	case errors.Is(err, ErrClosed):
		apierr.Write(w, r, http.StatusConflict, "E_TICKET_CLOSED", "ticket is closed")
	default:
		apierr.Write(w, r, http.StatusInternalServerError, "E_INTERNAL", err.Error())
	}
}
