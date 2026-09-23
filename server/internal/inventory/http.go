package inventory

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/otter/server/internal/http/apierr"
	"github.com/otter/server/internal/httpctx"
)

// Handler bundles the inventory endpoints.
type Handler struct {
	svc *Service
	log *slog.Logger
}

func NewHandler(svc *Service, logger *slog.Logger) *Handler {
	if logger == nil {
		logger = slog.Default()
	}
	return &Handler{svc: svc, log: logger}
}

// --- GET /v1/character/{id}/inventory ---------------------------------------

func (h *Handler) View(w http.ResponseWriter, r *http.Request) {
	acct, char, ok := requireCharacterOwner(w, r)
	if !ok {
		return
	}
	_ = acct
	view, err := h.svc.View(r.Context(), char)
	if err != nil {
		writeErr(w, r, err)
		return
	}
	apierr.WriteJSON(w, http.StatusOK, view)
}

// --- POST /v1/character/{id}/inventory/{item_id}/move -----------------------

type moveReq struct {
	TargetSlot int16 `json:"target_slot"`
}

func (h *Handler) Move(w http.ResponseWriter, r *http.Request) {
	acct, char, ok := requireCharacterOwner(w, r)
	if !ok {
		return
	}
	_ = acct
	itemID, err := uuid.Parse(chi.URLParam(r, "item_id"))
	if err != nil {
		apierr.Write(w, r, http.StatusBadRequest, "E_INVALID_ID", "item_id: "+err.Error())
		return
	}
	var body moveReq
	if err := decodeJSON(r, &body); err != nil {
		apierr.Write(w, r, http.StatusBadRequest, "E_INVALID_JSON", err.Error())
		return
	}
	if err := h.svc.Move(r.Context(), char, itemID, body.TargetSlot); err != nil {
		writeErr(w, r, err)
		return
	}
	// Return updated view for client convenience.
	view, _ := h.svc.View(r.Context(), char)
	apierr.WriteJSON(w, http.StatusOK, view)
}

// --- POST /v1/character/{id}/inventory/{item_id}/split ----------------------

type splitReq struct {
	Amount     int32 `json:"amount"`
	TargetSlot int16 `json:"target_slot"`
}

func (h *Handler) Split(w http.ResponseWriter, r *http.Request) {
	acct, char, ok := requireCharacterOwner(w, r)
	if !ok {
		return
	}
	_ = acct
	itemID, err := uuid.Parse(chi.URLParam(r, "item_id"))
	if err != nil {
		apierr.Write(w, r, http.StatusBadRequest, "E_INVALID_ID", "item_id: "+err.Error())
		return
	}
	var body splitReq
	if err := decodeJSON(r, &body); err != nil {
		apierr.Write(w, r, http.StatusBadRequest, "E_INVALID_JSON", err.Error())
		return
	}
	newID, err := h.svc.Split(r.Context(), char, itemID, body.Amount, body.TargetSlot)
	if err != nil {
		writeErr(w, r, err)
		return
	}
	apierr.WriteJSON(w, http.StatusOK, map[string]any{
		"new_instance_id": newID.String(),
	})
}

// --- DELETE /v1/character/{id}/inventory/{item_id} --------------------------

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	acct, char, ok := requireCharacterOwner(w, r)
	if !ok {
		return
	}
	_ = acct
	itemID, err := uuid.Parse(chi.URLParam(r, "item_id"))
	if err != nil {
		apierr.Write(w, r, http.StatusBadRequest, "E_INVALID_ID", "item_id: "+err.Error())
		return
	}
	if err := h.svc.Delete(r.Context(), char, itemID); err != nil {
		writeErr(w, r, err)
		return
	}
	apierr.WriteJSON(w, http.StatusOK, map[string]bool{"ok": true})
}

// --- helpers ----------------------------------------------------------------

func requireCharacterOwner(w http.ResponseWriter, r *http.Request) (uuid.UUID, uuid.UUID, bool) {
	acct := httpctx.AccountID(r.Context())
	if acct == uuid.Nil {
		apierr.Write(w, r, http.StatusUnauthorized, "E_UNAUTHENTICATED", "auth required")
		return uuid.Nil, uuid.Nil, false
	}
	char, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		apierr.Write(w, r, http.StatusBadRequest, "E_INVALID_ID", err.Error())
		return uuid.Nil, uuid.Nil, false
	}
	return acct, char, true
}

func decodeJSON(r *http.Request, dst any) error {
	body, err := io.ReadAll(io.LimitReader(r.Body, 8<<10))
	if err != nil {
		return err
	}
	if len(body) == 0 {
		return errors.New("empty body")
	}
	return json.Unmarshal(body, dst)
}

func writeErr(w http.ResponseWriter, r *http.Request, err error) {
	switch {
	case errors.Is(err, ErrNotFound):
		apierr.Write(w, r, http.StatusNotFound, "E_NOT_FOUND", "item not found")
	case errors.Is(err, ErrLocked):
		apierr.Write(w, r, http.StatusConflict, "E_ITEM_LOCKED", "item is locked (escrow)")
	case errors.Is(err, ErrBagFull):
		apierr.Write(w, r, http.StatusConflict, "E_BAG_FULL", "bag is full")
	case errors.Is(err, ErrSlotOccupied):
		apierr.Write(w, r, http.StatusConflict, "E_SLOT_OCCUPIED", "target slot occupied")
	case errors.Is(err, ErrSlotOutOfRange):
		apierr.Write(w, r, http.StatusBadRequest, "E_SLOT_OUT_OF_RANGE", "slot must be 0..59")
	case errors.Is(err, ErrNotStackable):
		apierr.Write(w, r, http.StatusConflict, "E_NOT_STACKABLE", "item cannot be split")
	case errors.Is(err, ErrInvalidSplit):
		apierr.Write(w, r, http.StatusBadRequest, "E_INVALID_SPLIT", "invalid split amount")
	default:
		apierr.Write(w, r, http.StatusInternalServerError, "E_INTERNAL", err.Error())
	}
}
