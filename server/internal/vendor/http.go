package vendor

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/otter/server/internal/domain/character"
	"github.com/otter/server/internal/economy"
	"github.com/otter/server/internal/http/apierr"
	"github.com/otter/server/internal/httpctx"
	"github.com/otter/server/internal/inventory"
)

type Handler struct {
	svc   *Service
	owner CharacterOwner
	log   *slog.Logger
}

// CharacterOwner verifies a character belongs to an account (character.Service).
// Vendor routes take character_id in the body, so the router-level path
// ownership guard does not cover them.
type CharacterOwner interface {
	Owns(ctx context.Context, characterID, accountID uuid.UUID) error
}

func NewHandler(svc *Service, owner CharacterOwner, logger *slog.Logger) *Handler {
	if logger == nil {
		logger = slog.Default()
	}
	return &Handler{svc: svc, owner: owner, log: logger}
}

// ownedCharacter parses a body character_id and rejects ids the caller does
// not own with 404 (existence is not leaked).
func (h *Handler) ownedCharacter(w http.ResponseWriter, r *http.Request, acct uuid.UUID, raw string) (uuid.UUID, bool) {
	id, err := uuid.Parse(raw)
	if err != nil {
		apierr.Write(w, r, http.StatusBadRequest, "E_INVALID_ID", "character_id: "+err.Error())
		return uuid.Nil, false
	}
	if err := h.owner.Owns(r.Context(), id, acct); err != nil {
		if errors.Is(err, character.ErrNotFound) || errors.Is(err, character.ErrForbidden) {
			apierr.Write(w, r, http.StatusNotFound, "E_NOT_FOUND", "character not found")
		} else {
			h.log.WarnContext(r.Context(), "vendor: ownership check failed", slog.String("err", err.Error()))
			apierr.Write(w, r, http.StatusInternalServerError, "E_INTERNAL", "internal error")
		}
		return uuid.Nil, false
	}
	return id, true
}

// --- GET /v1/vendor/{code} --------------------------------------------------

func (h *Handler) Info(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")
	info, err := h.svc.Info(r.Context(), code)
	if err != nil {
		writeErr(w, r, err)
		return
	}
	apierr.WriteJSON(w, http.StatusOK, info)
}

// --- POST /v1/vendor/{code}/buy ---------------------------------------------

type buyReq struct {
	TemplateCode        string `json:"template_code"`
	ExpectedPriceCopper *int64 `json:"expected_price_copper,omitempty"`
	CharacterID         string `json:"character_id"`
}

func (h *Handler) Buy(w http.ResponseWriter, r *http.Request) {
	acct := httpctx.AccountID(r.Context())
	if acct == uuid.Nil {
		apierr.Write(w, r, http.StatusUnauthorized, "E_UNAUTHENTICATED", "auth required")
		return
	}
	var body buyReq
	if err := decodeJSON(r, &body); err != nil {
		apierr.Write(w, r, http.StatusBadRequest, "E_INVALID_JSON", err.Error())
		return
	}
	charID, ok := h.ownedCharacter(w, r, acct, body.CharacterID)
	if !ok {
		return
	}
	code := chi.URLParam(r, "code")
	idem := r.Header.Get("Idempotency-Key")
	res, err := h.svc.Buy(r.Context(), BuyRequest{
		VendorCode:          code,
		CharacterID:         charID,
		TemplateCode:        body.TemplateCode,
		ExpectedPriceCopper: body.ExpectedPriceCopper,
		IdempotencyKey:      idem,
	})
	if err != nil {
		writeErr(w, r, err)
		return
	}
	apierr.WriteJSON(w, http.StatusOK, res)
}

// --- POST /v1/vendor/{code}/sell --------------------------------------------

type sellReq struct {
	ItemInstanceID string `json:"item_instance_id"`
	CharacterID    string `json:"character_id"`
}

func (h *Handler) Sell(w http.ResponseWriter, r *http.Request) {
	acct := httpctx.AccountID(r.Context())
	if acct == uuid.Nil {
		apierr.Write(w, r, http.StatusUnauthorized, "E_UNAUTHENTICATED", "auth required")
		return
	}
	var body sellReq
	if err := decodeJSON(r, &body); err != nil {
		apierr.Write(w, r, http.StatusBadRequest, "E_INVALID_JSON", err.Error())
		return
	}
	charID, ok := h.ownedCharacter(w, r, acct, body.CharacterID)
	if !ok {
		return
	}
	instID, err := uuid.Parse(body.ItemInstanceID)
	if err != nil {
		apierr.Write(w, r, http.StatusBadRequest, "E_INVALID_ID", "item_instance_id: "+err.Error())
		return
	}
	res, err := h.svc.Sell(r.Context(), SellRequest{
		VendorCode:     chi.URLParam(r, "code"),
		CharacterID:    charID,
		ItemInstanceID: instID,
		IdempotencyKey: r.Header.Get("Idempotency-Key"),
	})
	if err != nil {
		writeErr(w, r, err)
		return
	}
	apierr.WriteJSON(w, http.StatusOK, res)
}

// --- POST /v1/vendor/{code}/buyback/{buyback_id} ----------------------------

type buybackReq struct {
	CharacterID string `json:"character_id"`
}

func (h *Handler) Buyback(w http.ResponseWriter, r *http.Request) {
	acct := httpctx.AccountID(r.Context())
	if acct == uuid.Nil {
		apierr.Write(w, r, http.StatusUnauthorized, "E_UNAUTHENTICATED", "auth required")
		return
	}
	var body buybackReq
	if err := decodeJSON(r, &body); err != nil {
		apierr.Write(w, r, http.StatusBadRequest, "E_INVALID_JSON", err.Error())
		return
	}
	charID, ok := h.ownedCharacter(w, r, acct, body.CharacterID)
	if !ok {
		return
	}
	buybackID, err := uuid.Parse(chi.URLParam(r, "buyback_id"))
	if err != nil {
		apierr.Write(w, r, http.StatusBadRequest, "E_INVALID_ID", "buyback_id: "+err.Error())
		return
	}
	res, err := h.svc.Buyback(r.Context(), charID, buybackID, r.Header.Get("Idempotency-Key"))
	if err != nil {
		writeErr(w, r, err)
		return
	}
	apierr.WriteJSON(w, http.StatusOK, res)
}

// --- GET /v1/character/{id}/buybacks ----------------------------------------

func (h *Handler) ListBuybacks(w http.ResponseWriter, r *http.Request) {
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
	entries, err := h.svc.ListBuybacks(r.Context(), charID)
	if err != nil {
		writeErr(w, r, err)
		return
	}
	out := make([]map[string]any, 0, len(entries))
	for _, b := range entries {
		var snap ItemSnapshot
		_ = json.Unmarshal(b.ItemSnapshot, &snap)
		out = append(out, map[string]any{
			"id":                b.ID.String(),
			"vendor_code":       b.VendorCode,
			"template_code":     snap.TemplateCode,
			"rarity":            snap.Rarity,
			"prefix":            snap.Prefix,
			"suffix":            snap.Suffix,
			"stack_quantity":    snap.StackQuantity,
			"sold_price_copper": b.SoldPriceCopper,
			"sold_at":           b.SoldAt.UTC().Format(time.RFC3339),
			"expires_at":        b.ExpiresAt.UTC().Format(time.RFC3339),
		})
	}
	apierr.WriteJSON(w, http.StatusOK, map[string]any{
		"buybacks": out,
		"count":    len(out),
	})
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
	var ve *ValidationError
	switch {
	case errors.As(err, &ve):
		apierr.WriteField(w, r, http.StatusBadRequest, "E_VALIDATION", ve.Message, ve.Field)
	case errors.Is(err, ErrVendorNotFound):
		apierr.Write(w, r, http.StatusNotFound, "E_VENDOR_NOT_FOUND", "vendor not found")
	case errors.Is(err, ErrItemNotSold):
		apierr.Write(w, r, http.StatusBadRequest, "E_ITEM_NOT_SOLD", "vendor does not sell this item")
	case errors.Is(err, ErrItemNoSell):
		apierr.Write(w, r, http.StatusConflict, "E_ITEM_NO_SELL", "item cannot be sold")
	case errors.Is(err, ErrItemBound):
		apierr.Write(w, r, http.StatusConflict, "E_ITEM_BOUND", "bound items cannot be sold")
	case errors.Is(err, ErrItemLocked):
		apierr.Write(w, r, http.StatusConflict, "E_ITEM_LOCKED", "item is in escrow")
	case errors.Is(err, ErrItemNotOwned):
		apierr.Write(w, r, http.StatusNotFound, "E_NOT_FOUND", "item not found")
	case errors.Is(err, ErrBuybackNotFound):
		apierr.Write(w, r, http.StatusNotFound, "E_BUYBACK_NOT_FOUND", "buyback not found")
	case errors.Is(err, ErrBuybackExpired):
		apierr.Write(w, r, http.StatusConflict, "E_BUYBACK_EXPIRED", "buyback window expired")
	case errors.Is(err, ErrBuybackClaimed):
		apierr.Write(w, r, http.StatusConflict, "E_BUYBACK_CLAIMED", "buyback already claimed")
	case errors.Is(err, economy.ErrInsufficientBalance):
		apierr.Write(w, r, http.StatusConflict, "E_INSUFFICIENT_GOLD", "not enough gold")
	case errors.Is(err, inventory.ErrBagFull):
		apierr.Write(w, r, http.StatusConflict, "E_BAG_FULL", "bag is full")
	default:
		// Internal detail goes to the log, never to the client.
		slog.ErrorContext(r.Context(), "vendor: unhandled error", slog.String("err", err.Error()))
		apierr.Write(w, r, http.StatusInternalServerError, "E_INTERNAL", "internal error")
	}
}
