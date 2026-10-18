package hunt

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/otter/server/internal/action"
	"github.com/otter/server/internal/domain/character"
	"github.com/otter/server/internal/http/apierr"
	"github.com/otter/server/internal/httpctx"
)

// EventPublisher is the narrow ws-fanout surface the handler uses to notify
// the account's live sessions when a hunt reveals. Kept as an interface so
// hunt has no direct dependency on the wsgateway package.
type EventPublisher interface {
	Publish(accountID uuid.UUID, kind string, payload any) int
}

// Handler bundles the hunt HTTP routes.
type Handler struct {
	svc *Service
	log *slog.Logger
	pub EventPublisher // nil-safe: no-op when unset
}

func NewHandler(svc *Service, logger *slog.Logger) *Handler {
	if logger == nil {
		logger = slog.Default()
	}
	return &Handler{svc: svc, log: logger}
}

// WithPublisher wires the ws event pipe. Returns the handler for chaining.
func (h *Handler) WithPublisher(p EventPublisher) *Handler {
	h.pub = p
	return h
}

// publish is the nil-safe entry point; every hook goes through here.
func (h *Handler) publish(accountID uuid.UUID, kind string, payload any) {
	if h.pub == nil {
		return
	}
	h.pub.Publish(accountID, kind, payload)
}

// startBody is the POST /v1/character/{id}/hunt body.
type startBody struct {
	MobCode          string `json:"mob_code"`
	DurationSeconds  int32  `json:"duration_seconds"`
	EncounterSeconds int32  `json:"encounter_seconds,omitempty"`
}

// Start creates a new hunt contract.
func (h *Handler) Start(w http.ResponseWriter, r *http.Request) {
	acct := httpctx.AccountID(r.Context())
	if acct == uuid.Nil {
		apierr.Write(w, r, http.StatusUnauthorized, "E_UNAUTHENTICATED", "auth required")
		return
	}
	charID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		apierr.Write(w, r, http.StatusBadRequest, "E_INVALID_ID", err.Error())
		return
	}
	var body startBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		apierr.Write(w, r, http.StatusBadRequest, "E_INVALID_JSON", err.Error())
		return
	}
	res, err := h.svc.Start(r.Context(), StartRequest{
		CharacterID:     charID,
		AccountID:       acct,
		MobCode:         body.MobCode,
		DurationSeconds: body.DurationSeconds,
		EncounterSecs:   body.EncounterSeconds,
	})
	if err != nil {
		writeErr(w, r, err)
		return
	}
	apierr.WriteJSON(w, http.StatusCreated, res)
}

// Reveal grants rewards for elapsed encounters (idempotent).
func (h *Handler) Reveal(w http.ResponseWriter, r *http.Request) {
	acct := httpctx.AccountID(r.Context())
	if acct == uuid.Nil {
		apierr.Write(w, r, http.StatusUnauthorized, "E_UNAUTHENTICATED", "auth required")
		return
	}
	contractID, err := uuid.Parse(chi.URLParam(r, "contract_id"))
	if err != nil {
		apierr.Write(w, r, http.StatusBadRequest, "E_INVALID_ID", err.Error())
		return
	}
	res, err := h.svc.Reveal(r.Context(), contractID, acct)
	if err != nil {
		writeErr(w, r, err)
		return
	}
	// Fan-out to any live sessions of this account. Non-blocking; a slow
	// client gets dropped in the hub.
	h.publish(acct, "hunt.reveal", map[string]any{
		"contract_id": contractID.String(),
		"gold_copper": res.Hunt.TotalGoldCopper,
		"xp":          res.Hunt.TotalXP,
		"victories":   res.Hunt.Victories,
		"defeats":     res.Hunt.Defeats,
		"drops":       res.Hunt.CommonDrops + res.Hunt.RareDrops + res.Hunt.LegendaryDrops,
		"leveled_up":  res.XPGain != nil && res.XPGain.LeveledUp,
	})
	apierr.WriteJSON(w, http.StatusOK, res)
}

// Stop cancels the active hunt (no reward).
func (h *Handler) Stop(w http.ResponseWriter, r *http.Request) {
	acct := httpctx.AccountID(r.Context())
	if acct == uuid.Nil {
		apierr.Write(w, r, http.StatusUnauthorized, "E_UNAUTHENTICATED", "auth required")
		return
	}
	contractID, err := uuid.Parse(chi.URLParam(r, "contract_id"))
	if err != nil {
		apierr.Write(w, r, http.StatusBadRequest, "E_INVALID_ID", err.Error())
		return
	}
	c, err := h.svc.Stop(r.Context(), contractID, acct)
	if err != nil {
		writeErr(w, r, err)
		return
	}
	apierr.WriteJSON(w, http.StatusOK, map[string]any{
		"contract_id": c.ID.String(),
		"state":       string(c.State),
	})
}

// Active returns the currently-running hunt for the character, or 404.
func (h *Handler) Active(w http.ResponseWriter, r *http.Request) {
	acct := httpctx.AccountID(r.Context())
	if acct == uuid.Nil {
		apierr.Write(w, r, http.StatusUnauthorized, "E_UNAUTHENTICATED", "auth required")
		return
	}
	charID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		apierr.Write(w, r, http.StatusBadRequest, "E_INVALID_ID", err.Error())
		return
	}
	c, err := h.svc.Active(r.Context(), charID, acct)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			apierr.WriteJSON(w, http.StatusOK, map[string]any{"contract": nil})
			return
		}
		writeErr(w, r, err)
		return
	}
	apierr.WriteJSON(w, http.StatusOK, map[string]any{
		"contract": summarize(c),
	})
}

// History lists the character's past hunts (newest first).
func (h *Handler) History(w http.ResponseWriter, r *http.Request) {
	acct := httpctx.AccountID(r.Context())
	if acct == uuid.Nil {
		apierr.Write(w, r, http.StatusUnauthorized, "E_UNAUTHENTICATED", "auth required")
		return
	}
	charID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		apierr.Write(w, r, http.StatusBadRequest, "E_INVALID_ID", err.Error())
		return
	}
	limit := 20
	if q := r.URL.Query().Get("limit"); q != "" {
		if n, err := strconv.Atoi(q); err == nil {
			limit = n
		}
	}
	list, err := h.svc.History(r.Context(), charID, acct, limit)
	if err != nil {
		writeErr(w, r, err)
		return
	}
	out := make([]map[string]any, 0, len(list))
	for i := range list {
		out = append(out, summarize(&list[i]))
	}
	apierr.WriteJSON(w, http.StatusOK, map[string]any{"count": len(out), "hunts": out})
}

// summarize is the wire shape for a Contract — hides player_snapshot bytes.
func summarize(c *Contract) map[string]any {
	m := map[string]any{
		"id":                c.ID.String(),
		"character_id":      c.CharacterID.String(),
		"mob_code":          c.MobCode,
		"encounter_seconds": c.EncounterSecs,
		"duration_seconds":  c.DurationSecs,
		"started_at":        c.StartedAt,
		"expires_at":        c.ExpiresAt,
		"state":             string(c.State),
	}
	if c.RevealedAt != nil {
		m["revealed_at"] = *c.RevealedAt
	}
	if c.RevealTxnID != nil {
		m["reveal_txn_id"] = c.RevealTxnID.String()
	}
	if c.CancelledAt != nil {
		m["cancelled_at"] = *c.CancelledAt
	}
	return m
}

func writeErr(w http.ResponseWriter, r *http.Request, err error) {
	if action.Write(w, r, err) {
		return
	}
	switch {
	case errors.Is(err, ErrNotFound):
		apierr.Write(w, r, http.StatusNotFound, "E_NOT_FOUND", "hunt not found")
	case errors.Is(err, ErrForbidden), errors.Is(err, character.ErrForbidden), errors.Is(err, character.ErrNotFound):
		apierr.Write(w, r, http.StatusNotFound, "E_NOT_FOUND", "hunt not found") // leak-safe
	case errors.Is(err, ErrSlotBusy):
		apierr.Write(w, r, http.StatusConflict, "E_HUNT_ACTIVE", "character already has an active hunt")
	case errors.Is(err, ErrNotActive):
		apierr.Write(w, r, http.StatusConflict, "E_HUNT_NOT_ACTIVE", "hunt is not active")
	case errors.Is(err, ErrBagFull):
		apierr.Write(w, r, http.StatusConflict, "E_BAG_FULL", "bag is full")
	case errors.Is(err, ErrTooWounded):
		apierr.Write(w, r, http.StatusConflict, "E_TOO_WOUNDED", "rest or heal before fighting")
	case errors.Is(err, ErrUnknownMob):
		apierr.Write(w, r, http.StatusBadRequest, "E_UNKNOWN_MOB", "mob_code not in catalog")
	case errors.Is(err, ErrValidation):
		apierr.Write(w, r, http.StatusBadRequest, "E_VALIDATION", "invalid input")
	default:
		apierr.Write(w, r, http.StatusInternalServerError, "E_INTERNAL", err.Error())
	}
}
