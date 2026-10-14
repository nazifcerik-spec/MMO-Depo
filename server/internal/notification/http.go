package notification

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/otter/server/internal/http/apierr"
	"github.com/otter/server/internal/httpctx"
)

// Handler serves account-scoped notification routes.
type Handler struct{ svc *Service }

func NewHandler(svc *Service) *Handler { return &Handler{svc: svc} }

func account(w http.ResponseWriter, r *http.Request) (uuid.UUID, bool) {
	a := httpctx.AccountID(r.Context())
	if a == uuid.Nil {
		apierr.Write(w, r, http.StatusUnauthorized, "E_UNAUTHENTICATED", "auth required")
	}
	return a, a != uuid.Nil
}

func decode(w http.ResponseWriter, r *http.Request, v any) bool {
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		apierr.Write(w, r, http.StatusBadRequest, "E_INVALID_JSON", err.Error())
		return false
	}
	return true
}

func reply(w http.ResponseWriter, r *http.Request, v any, err error) {
	if err != nil {
		writeErr(w, r, err)
		return
	}
	apierr.WriteJSON(w, http.StatusOK, v)
}

func (h *Handler) RegisterDevice(w http.ResponseWriter, r *http.Request) {
	a, ok := account(w, r)
	if !ok {
		return
	}
	var b DeviceRequest
	if !decode(w, r, &b) {
		return
	}
	reply(w, r, map[string]any{"registered": true}, h.svc.RegisterDevice(r.Context(), a, b))
}

func (h *Handler) UnregisterDevice(w http.ResponseWriter, r *http.Request) {
	a, ok := account(w, r)
	if !ok {
		return
	}
	id, err := uuid.Parse(chi.URLParam(r, "device_id"))
	if err != nil {
		apierr.Write(w, r, http.StatusBadRequest, "E_INVALID_ID", err.Error())
		return
	}
	reply(w, r, map[string]any{"unregistered": true}, h.svc.UnregisterDevice(r.Context(), a, id))
}

func (h *Handler) Inbox(w http.ResponseWriter, r *http.Request) {
	a, ok := account(w, r)
	if !ok {
		return
	}
	items, unread, err := h.svc.Inbox(r.Context(), a, r.URL.Query().Get("unread") == "1", r.URL.Query().Get("rewards") == "1")
	reply(w, r, map[string]any{"count": len(items), "unread": unread, "notifications": items}, err)
}

func (h *Handler) MarkRead(w http.ResponseWriter, r *http.Request) {
	a, ok := account(w, r)
	if !ok {
		return
	}
	id := uuid.Nil
	if p := chi.URLParam(r, "notification_id"); p != "" {
		var err error
		if id, err = uuid.Parse(p); err != nil {
			apierr.Write(w, r, http.StatusBadRequest, "E_INVALID_ID", err.Error())
			return
		}
	}
	n, err := h.svc.MarkRead(r.Context(), a, id)
	reply(w, r, map[string]any{"marked": n}, err)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	a, ok := account(w, r)
	if !ok {
		return
	}
	id, err := uuid.Parse(chi.URLParam(r, "notification_id"))
	if err != nil {
		apierr.Write(w, r, http.StatusBadRequest, "E_INVALID_ID", err.Error())
		return
	}
	n, err := h.svc.Delete(r.Context(), a, id)
	reply(w, r, map[string]any{"deleted": n}, err)
}

func (h *Handler) GetPreferences(w http.ResponseWriter, r *http.Request) {
	a, ok := account(w, r)
	if !ok {
		return
	}
	p, err := h.svc.Preferences(r.Context(), a)
	reply(w, r, p, err)
}

func (h *Handler) PutPreferences(w http.ResponseWriter, r *http.Request) {
	a, ok := account(w, r)
	if !ok {
		return
	}
	var b PreferenceUpdate
	if !decode(w, r, &b) {
		return
	}
	p, err := h.svc.UpdatePreferences(r.Context(), a, b)
	reply(w, r, p, err)
}

func (h *Handler) Mute(w http.ResponseWriter, r *http.Request) {
	a, ok := account(w, r)
	if !ok {
		return
	}
	var b struct {
		Hours int `json:"hours"`
	}
	if !decode(w, r, &b) {
		return
	}
	if b.Hours == 0 {
		b.Hours = 24
	}
	until, err := h.svc.SetGlobalMute(r.Context(), a, b.Hours)
	reply(w, r, map[string]any{"global_mute_until": until}, err)
}

func (h *Handler) Unmute(w http.ResponseWriter, r *http.Request) {
	a, ok := account(w, r)
	if !ok {
		return
	}
	_, err := h.svc.SetGlobalMute(r.Context(), a, 0)
	reply(w, r, map[string]any{"global_mute_until": nil}, err)
}

func writeErr(w http.ResponseWriter, r *http.Request, err error) {
	for _, e := range []error{ErrUnknownKind, ErrBadMode, ErrBadPlatform, ErrBadToken, ErrBadClock, ErrBadTimezone, ErrBadMuteHours} {
		if errors.Is(err, e) {
			apierr.Write(w, r, http.StatusBadRequest, "E_VALIDATION", err.Error())
			return
		}
	}
	if errors.Is(err, ErrNotFound) {
		apierr.Write(w, r, http.StatusNotFound, "E_NOT_FOUND", err.Error())
		return
	}
	apierr.Write(w, r, http.StatusInternalServerError, "E_INTERNAL", err.Error())
}
