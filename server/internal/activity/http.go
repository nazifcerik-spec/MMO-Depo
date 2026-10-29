package activity

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/otter/server/internal/action"
	"github.com/otter/server/internal/http/apierr"
	"github.com/otter/server/internal/httpctx"
)

// Handler bundles the activity endpoints.
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

// --- POST /v1/character/{id}/activity ---------------------------------------

type startReq struct {
	Kind            string          `json:"kind"`
	DurationSeconds int32           `json:"duration_seconds"`
	Target          string          `json:"target,omitempty"` // gathering resource (DD-244)
	Params          json.RawMessage `json:"params,omitempty"`
}

func (h *Handler) Start(w http.ResponseWriter, r *http.Request) {
	acct, char, ok := requireCharacterOwner(w, r)
	if !ok {
		return
	}
	var body startReq
	if err := decodeJSON(r, &body); err != nil {
		apierr.Write(w, r, http.StatusBadRequest, "E_INVALID_JSON", err.Error())
		return
	}
	res, err := h.svc.Start(r.Context(), StartRequest{
		CharacterID:     char,
		AccountID:       acct,
		Kind:            Kind(body.Kind),
		DurationSeconds: body.DurationSeconds,
		Target:          body.Target,
		Params:          body.Params,
	})
	if err != nil {
		writeActivityErr(w, r, err)
		return
	}
	apierr.WriteJSON(w, http.StatusCreated, map[string]any{
		"contract":        contractResponse(res.Contract),
		"full_yield_peek": res.FullYieldPeek,
	})
}

// --- GET /v1/character/{id}/activity/active ---------------------------------

func (h *Handler) Active(w http.ResponseWriter, r *http.Request) {
	acct, char, ok := requireCharacterOwner(w, r)
	if !ok {
		return
	}
	c, preview, err := h.svc.Active(r.Context(), char, acct)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			apierr.WriteJSON(w, http.StatusOK, map[string]any{"active": nil})
			return
		}
		writeActivityErr(w, r, err)
		return
	}
	apierr.WriteJSON(w, http.StatusOK, map[string]any{
		"contract":      contractResponse(c),
		"current_yield": preview,
	})
}

// --- GET /v1/character/{id}/activity/history --------------------------------

func (h *Handler) History(w http.ResponseWriter, r *http.Request) {
	acct, char, ok := requireCharacterOwner(w, r)
	if !ok {
		return
	}
	list, err := h.svc.History(r.Context(), char, acct, 20)
	if err != nil {
		writeActivityErr(w, r, err)
		return
	}
	out := make([]map[string]any, 0, len(list))
	for i := range list {
		out = append(out, contractResponse(&list[i]))
	}
	apierr.WriteJSON(w, http.StatusOK, map[string]any{"contracts": out})
}

// --- POST /v1/activity/{contract_id}/reveal ---------------------------------

func (h *Handler) Reveal(w http.ResponseWriter, r *http.Request) {
	acct := httpctx.AccountID(r.Context())
	id, err := parseUUIDParam(r, "contract_id")
	if err != nil {
		apierr.Write(w, r, http.StatusBadRequest, "E_INVALID_ID", err.Error())
		return
	}
	res, err := h.svc.Reveal(r.Context(), id, acct)
	if err != nil {
		writeActivityErr(w, r, err)
		return
	}
	apierr.WriteJSON(w, http.StatusOK, revealResponse(res))
}

func revealResponse(res *RevealResult) map[string]any {
	return map[string]any{
		"yield":         res.Yield,
		"ledger_txn_id": res.LedgerTxnID.String(),
		"replayed":      res.Replayed,
		"gather":        res.Gather,
		"party_share":   res.Party, // DD-305; null without a party bonus
	}
}

// --- POST /v1/activity/{contract_id}/restart --------------------------------

// Restart collects the activity so far and starts it again (DD-310): the
// reveal as /reveal answers it, plus the new contract (null and an error
// code when it could not start).
func (h *Handler) Restart(w http.ResponseWriter, r *http.Request) {
	acct := httpctx.AccountID(r.Context())
	id, err := parseUUIDParam(r, "contract_id")
	if err != nil {
		apierr.Write(w, r, http.StatusBadRequest, "E_INVALID_ID", err.Error())
		return
	}
	rev, st, err := h.svc.Restart(r.Context(), id, acct)
	if rev == nil {
		writeActivityErr(w, r, err)
		return
	}
	out := revealResponse(rev)
	out["contract"] = nil
	if st != nil {
		out["contract"] = contractResponse(st.Contract)
	}
	if err != nil {
		out["restart_error"] = apierr.CodeOf(r, func(w http.ResponseWriter, r *http.Request) { writeActivityErr(w, r, err) })
	}
	apierr.WriteJSON(w, http.StatusOK, out)
}

// --- POST /v1/activity/{contract_id}/stop -----------------------------------

func (h *Handler) Stop(w http.ResponseWriter, r *http.Request) {
	acct := httpctx.AccountID(r.Context())
	id, err := parseUUIDParam(r, "contract_id")
	if err != nil {
		apierr.Write(w, r, http.StatusBadRequest, "E_INVALID_ID", err.Error())
		return
	}
	c, err := h.svc.Stop(r.Context(), id, acct)
	if err != nil {
		writeActivityErr(w, r, err)
		return
	}
	apierr.WriteJSON(w, http.StatusOK, contractResponse(c))
}

// --- helpers ---------------------------------------------------------------

// requireCharacterOwner reads account id from ctx + character id from URL and
// verifies both are present. Actual ownership check happens in the service.
func requireCharacterOwner(w http.ResponseWriter, r *http.Request) (uuid.UUID, uuid.UUID, bool) {
	acct := httpctx.AccountID(r.Context())
	if acct == uuid.Nil {
		apierr.Write(w, r, http.StatusUnauthorized, "E_UNAUTHENTICATED", "auth required")
		return uuid.Nil, uuid.Nil, false
	}
	char, err := parseUUIDParam(r, "id")
	if err != nil {
		apierr.Write(w, r, http.StatusBadRequest, "E_INVALID_ID", err.Error())
		return uuid.Nil, uuid.Nil, false
	}
	return acct, char, true
}

func parseUUIDParam(r *http.Request, name string) (uuid.UUID, error) {
	raw := chi.URLParam(r, name)
	if raw == "" {
		return uuid.Nil, errors.New("missing " + name)
	}
	return uuid.Parse(raw)
}

func decodeJSON(r *http.Request, dst any) error {
	body, err := io.ReadAll(io.LimitReader(r.Body, 16<<10))
	if err != nil {
		return err
	}
	if len(body) == 0 {
		return errors.New("empty body")
	}
	return json.Unmarshal(body, dst)
}

// contractResponse is the wire shape for a Contract; unmarshalled JSONB blobs
// are inlined for client convenience.
func contractResponse(c *Contract) map[string]any {
	out := map[string]any{
		"id":               c.ID.String(),
		"character_id":     c.CharacterID.String(),
		"kind":             string(c.Kind),
		"duration_seconds": c.DurationSeconds,
		"started_at":       c.StartedAt.UTC().Format(time.RFC3339),
		"expires_at":       c.ExpiresAt.UTC().Format(time.RFC3339),
		"state":            string(c.State),
	}
	if c.RevealedAt != nil {
		out["revealed_at"] = c.RevealedAt.UTC().Format(time.RFC3339)
	}
	if c.CancelledAt != nil {
		out["cancelled_at"] = c.CancelledAt.UTC().Format(time.RFC3339)
	}
	if len(c.RevealResult) > 0 {
		var y YieldSummary
		if err := json.Unmarshal(c.RevealResult, &y); err == nil {
			out["reveal_yield"] = y
		}
	}
	return out
}

func writeActivityErr(w http.ResponseWriter, r *http.Request, err error) {
	if action.Write(w, r, err) {
		return
	}
	var ve *ValidationError
	switch {
	case errors.As(err, &ve):
		apierr.WriteField(w, r, http.StatusBadRequest, "E_VALIDATION", ve.Message, ve.Field)
	case errors.Is(err, ErrNotFound):
		apierr.Write(w, r, http.StatusNotFound, "E_NOT_FOUND", "contract not found")
	case errors.Is(err, ErrForbidden):
		// Existence-leak-safe: same shape as not-found.
		apierr.Write(w, r, http.StatusNotFound, "E_NOT_FOUND", "contract not found")
	case errors.Is(err, ErrSlotBusy):
		apierr.Write(w, r, http.StatusConflict, "E_SLOT_BUSY", "character already has an active contract")
	case errors.Is(err, ErrNotActive):
		apierr.Write(w, r, http.StatusConflict, "E_NOT_ACTIVE", "contract is not active")
	case errors.Is(err, ErrUnknownTarget):
		apierr.WriteField(w, r, http.StatusBadRequest, "E_VALIDATION", "unknown gathering resource", "target")
	case errors.Is(err, ErrBagFull):
		apierr.Write(w, r, http.StatusConflict, "E_BAG_FULL", "bag is full")
	case errors.Is(err, ErrGatherRegion):
		apierr.Write(w, r, http.StatusConflict, "E_GATHER_REGION", "the resource grows on another map")
	case errors.Is(err, ErrGatherLevelLow):
		apierr.Write(w, r, http.StatusConflict, "E_GATHER_LEVEL", "gathering level too low for this resource")
	default:
		apierr.Write(w, r, http.StatusInternalServerError, "E_INTERNAL", err.Error())
	}
}
