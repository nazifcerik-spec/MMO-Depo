package inventory

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/otter/server/internal/http/apierr"
	"github.com/otter/server/internal/httpctx"
)

// EquipmentHandler bundles the equipment HTTP routes.
type EquipmentHandler struct {
	svc *EquipmentService
	log *slog.Logger
}

func NewEquipmentHandler(svc *EquipmentService, logger *slog.Logger) *EquipmentHandler {
	if logger == nil {
		logger = slog.Default()
	}
	return &EquipmentHandler{svc: svc, log: logger}
}

// --- GET /v1/character/{id}/equipment ---------------------------------------

func (h *EquipmentHandler) View(w http.ResponseWriter, r *http.Request) {
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
	view, err := h.svc.View(r.Context(), charID)
	if err != nil {
		writeEquipErr(w, r, err)
		return
	}
	// Simplify wire — flatten slot map with template metadata.
	slots := map[string]map[string]any{}
	for slot, b := range view.Slots {
		row := map[string]any{
			"instance_id":   b.InstanceID.String(),
			"template_code": b.TemplateCode,
			"display_name":  b.DisplayNameKey,
			"rarity":        b.Rarity,
			"item_class":    b.ItemClass,
			"item_subclass": b.ItemSubclass,
			"is_bound":      b.IsBound,
		}
		if b.Prefix != nil {
			row["prefix"] = *b.Prefix
		}
		if b.Suffix != nil {
			row["suffix"] = *b.Suffix
		}
		if len(b.EffectiveBonuses) > 0 {
			row["effective_bonuses"] = b.EffectiveBonuses // rarity + affixes (DD-275)
		}
		slots[slot] = row
	}
	apierr.WriteJSON(w, http.StatusOK, map[string]any{
		"character_id": charID.String(),
		"slots":        slots,
	})
}

// --- POST /v1/character/{id}/items/{item_id}/enchant (DD-293) ---------------

type enchantBody struct {
	ScrollID uuid.UUID `json:"scroll_id"`
}

// Enchant puts a scroll's enchant on the item.
func (h *EquipmentHandler) Enchant(w http.ResponseWriter, r *http.Request) {
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
	itemID, err := uuid.Parse(chi.URLParam(r, "item_id"))
	if err != nil {
		apierr.Write(w, r, http.StatusBadRequest, "E_INVALID_ID", err.Error())
		return
	}
	var body enchantBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.ScrollID == uuid.Nil {
		apierr.Write(w, r, http.StatusBadRequest, "E_INVALID_JSON", "scroll_id required")
		return
	}
	res, err := h.svc.Enchant(r.Context(), charID, itemID, body.ScrollID)
	if err != nil {
		writeEquipErr(w, r, err)
		return
	}
	apierr.WriteJSON(w, http.StatusOK, res)
}

// --- POST /v1/character/{id}/equip/{item_id} --------------------------------

func (h *EquipmentHandler) Equip(w http.ResponseWriter, r *http.Request) {
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
	itemID, err := uuid.Parse(chi.URLParam(r, "item_id"))
	if err != nil {
		apierr.Write(w, r, http.StatusBadRequest, "E_INVALID_ID", "item_id: "+err.Error())
		return
	}
	slot, err := h.svc.Equip(r.Context(), charID, itemID)
	if err != nil {
		writeEquipErr(w, r, err)
		return
	}
	apierr.WriteJSON(w, http.StatusOK, map[string]any{
		"equipped_slot": slot,
		"item_id":       itemID.String(),
	})
}

// --- POST /v1/character/{id}/unequip/{slot} ---------------------------------

func (h *EquipmentHandler) Unequip(w http.ResponseWriter, r *http.Request) {
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
	slot := chi.URLParam(r, "slot")
	itemID, err := h.svc.Unequip(r.Context(), charID, slot)
	if err != nil {
		writeEquipErr(w, r, err)
		return
	}
	apierr.WriteJSON(w, http.StatusOK, map[string]any{
		"unequipped_slot": slot,
		"item_id":         itemID.String(),
	})
}

func writeEquipErr(w http.ResponseWriter, r *http.Request, err error) {
	switch {
	case errors.Is(err, ErrNotFound):
		apierr.Write(w, r, http.StatusNotFound, "E_NOT_FOUND", "item not found")
	case errors.Is(err, ErrNotEquippable):
		apierr.Write(w, r, http.StatusBadRequest, "E_NOT_EQUIPPABLE", "item cannot be equipped")
	case errors.Is(err, ErrClassRestriction):
		apierr.Write(w, r, http.StatusForbidden, "E_CLASS_RESTRICTION", "your class cannot equip this item")
	case errors.Is(err, ErrLevelRequirement):
		apierr.Write(w, r, http.StatusForbidden, "E_LEVEL_REQUIREMENT", "character level too low")
	case errors.Is(err, ErrSpecRequirement):
		apierr.Write(w, r, http.StatusForbidden, "E_SPEC_REQUIREMENT", "your specialization cannot equip this item")
	case errors.Is(err, ErrEquipmentPermission):
		apierr.Write(w, r, http.StatusForbidden, "E_EQUIPMENT_PERMISSION", "your class cannot use this equipment type")
	case errors.Is(err, ErrAttributeRequirement):
		apierr.Write(w, r, http.StatusForbidden, "E_ATTRIBUTE_REQUIREMENT", "attributes too low for this item")
	case errors.Is(err, ErrUnknownEquipSlot):
		apierr.Write(w, r, http.StatusBadRequest, "E_UNKNOWN_SLOT", "unknown equipment slot")
	case errors.Is(err, ErrNothingEquipped):
		apierr.Write(w, r, http.StatusNotFound, "E_NOTHING_EQUIPPED", "slot is empty")
	case errors.Is(err, ErrBagFull):
		apierr.Write(w, r, http.StatusConflict, "E_BAG_FULL", "bag full — cannot unequip")
	case errors.Is(err, ErrLocked):
		apierr.Write(w, r, http.StatusConflict, "E_ITEM_LOCKED", "item is locked")
	case errors.Is(err, ErrNotEnchantScroll):
		apierr.Write(w, r, http.StatusBadRequest, "E_NOT_ENCHANT", "that is not an enchanting scroll")
	case errors.Is(err, ErrEnchantSlot):
		apierr.Write(w, r, http.StatusConflict, "E_ENCHANT_SLOT", "the enchant does not fit that gear")
	default:
		apierr.Write(w, r, http.StatusInternalServerError, "E_INTERNAL", err.Error())
	}
}
