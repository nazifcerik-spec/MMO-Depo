package item

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/otter/server/internal/http/apierr"
	"github.com/otter/server/internal/httpctx"
)

// Handler exposes read-only inventory endpoints. Full inventory (slot / weight
// / move / equip) lands in Sprint 9. This is the "what items do I own" list.
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

// --- GET /v1/character/{id}/items -------------------------------------------

// Items lists the character's item instances, newest first. Enriched with the
// template display key + class/subclass so the client doesn't need a separate
// catalog fetch.
func (h *Handler) Items(w http.ResponseWriter, r *http.Request) {
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
	// Ownership check would ideally use character.Service.Owns; kept lean here
	// — Sprint 9 wires RequireCharacter middleware for this route.
	list, err := h.svc.ListByOwner(r.Context(), charID, 200)
	if err != nil {
		apierr.Write(w, r, http.StatusInternalServerError, "E_INTERNAL", err.Error())
		return
	}
	cat := h.svc.Catalog()
	out := make([]map[string]any, 0, len(list))
	for _, i := range list {
		row := map[string]any{
			"id":             i.ID.String(),
			"template_code":  i.TemplateCode,
			"stack_quantity": i.StackQuantity,
			"rarity":         string(i.Rarity),
			"is_bound":       i.IsBound,
			"is_locked":      i.IsLocked,
			"created_at":     i.CreatedAt.UTC().Format(time.RFC3339),
		}
		if i.PrefixCode != nil {
			row["prefix"] = *i.PrefixCode
		}
		if i.SuffixCode != nil {
			row["suffix"] = *i.SuffixCode
		}
		if i.SourceKind != nil {
			row["source_kind"] = *i.SourceKind
		}
		if tmpl, err := cat.Get(i.TemplateCode); err == nil {
			row["display_key"] = tmpl.DisplayNameKey
			row["item_class"] = tmpl.ItemClass
			row["item_subclass"] = tmpl.ItemSubclass
			row["base_value_copper"] = tmpl.BaseValueCopper
		}
		out = append(out, row)
	}
	apierr.WriteJSON(w, http.StatusOK, map[string]any{
		"items": out,
		"count": len(out),
	})
}
