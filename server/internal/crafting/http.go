package crafting

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/otter/server/internal/domain/character"
	"github.com/otter/server/internal/http/apierr"
	"github.com/otter/server/internal/httpctx"
	"github.com/otter/server/internal/inventory"
)

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

// --- Public catalog ---------------------------------------------------------

func (h *Handler) Professions(w http.ResponseWriter, r *http.Request) {
	out := make([]map[string]any, 0)
	for _, p := range h.svc.ListProfessions() {
		out = append(out, map[string]any{
			"code":             p.Code,
			"display_name_key": p.DisplayNameKey,
			"max_level":        p.MaxLevel,
			"kind":             p.Kind,
		})
	}
	apierr.WriteJSON(w, http.StatusOK, map[string]any{"count": len(out), "professions": out})
}

func (h *Handler) Recipes(w http.ResponseWriter, r *http.Request) {
	filterProf := r.URL.Query().Get("profession")
	all := h.svc.ListRecipes()
	out := make([]map[string]any, 0, len(all))
	for _, rec := range all {
		if filterProf != "" && rec.ProfessionCode != filterProf {
			continue
		}
		out = append(out, recipeRow(rec))
	}
	apierr.WriteJSON(w, http.StatusOK, map[string]any{"count": len(out), "recipes": out})
}

// --- Character-scoped -------------------------------------------------------

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
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
	list, err := h.svc.List(r.Context(), charID, acct)
	if err != nil {
		writeErr(w, r, err)
		return
	}
	out := make([]map[string]any, 0, len(list))
	for _, cp := range list {
		out = append(out, map[string]any{
			"profession_code": cp.ProfessionCode,
			"level":           cp.Level,
			"xp":              cp.XP,
			"xp_to_next":      ProfessionXPToNext(cp.Level),
			"learned_at":      cp.LearnedAt,
		})
	}
	apierr.WriteJSON(w, http.StatusOK, map[string]any{"count": len(out), "professions": out})
}

type learnBody struct {
	ProfessionCode string `json:"profession_code"`
}

func (h *Handler) Learn(w http.ResponseWriter, r *http.Request) {
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
	var body learnBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		apierr.Write(w, r, http.StatusBadRequest, "E_INVALID_JSON", err.Error())
		return
	}
	cp, err := h.svc.Learn(r.Context(), charID, acct, body.ProfessionCode)
	if err != nil {
		writeErr(w, r, err)
		return
	}
	apierr.WriteJSON(w, http.StatusOK, map[string]any{
		"profession_code": cp.ProfessionCode,
		"level":           cp.Level,
		"xp":              cp.XP,
		"learned_at":      cp.LearnedAt,
	})
}

type craftBody struct {
	RecipeCode string `json:"recipe_code"`
	Count      int    `json:"count,omitempty"` // default 1 (DD-245)
}

func (h *Handler) Craft(w http.ResponseWriter, r *http.Request) {
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
	var body craftBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		apierr.Write(w, r, http.StatusBadRequest, "E_INVALID_JSON", err.Error())
		return
	}
	if body.Count == 0 {
		body.Count = 1
	}
	res, err := h.svc.CraftMany(r.Context(), charID, acct, body.RecipeCode, body.Count)
	if err != nil {
		writeErr(w, r, err)
		return
	}
	apierr.WriteJSON(w, http.StatusOK, res)
}

// --- shared -----------------------------------------------------------------

func recipeRow(r Recipe) map[string]any {
	inputs := make([]map[string]any, 0, len(r.Inputs))
	for _, in := range r.Inputs {
		inputs = append(inputs, map[string]any{
			"template_code": in.TemplateCode,
			"quantity":      in.Quantity,
		})
	}
	return map[string]any{
		"code":                 r.Code,
		"display_name_key":     r.DisplayNameKey,
		"profession_code":      r.ProfessionCode,
		"min_skill":            r.MinSkill,
		"inputs":               inputs,
		"gold_cost_copper":     r.GoldCostCopper,
		"output_template_code": r.OutputTemplateCode,
		"output_stack":         r.OutputStack,
		"profession_xp_reward": r.ProfessionXPReward,
	}
}

func writeErr(w http.ResponseWriter, r *http.Request, err error) {
	switch {
	case errors.Is(err, ErrNotFound):
		apierr.Write(w, r, http.StatusNotFound, "E_NOT_FOUND", "not found")
	case errors.Is(err, ErrForbidden), errors.Is(err, character.ErrForbidden), errors.Is(err, character.ErrNotFound):
		apierr.Write(w, r, http.StatusNotFound, "E_NOT_FOUND", "not found")
	case errors.Is(err, ErrUnknownProfession):
		apierr.Write(w, r, http.StatusNotFound, "E_UNKNOWN_PROFESSION", "profession_code not in catalog")
	case errors.Is(err, ErrUnknownRecipe):
		apierr.Write(w, r, http.StatusNotFound, "E_UNKNOWN_RECIPE", "recipe_code not in catalog")
	case errors.Is(err, ErrProfessionNotLearned):
		apierr.Write(w, r, http.StatusConflict, "E_PROFESSION_NOT_LEARNED", "learn the profession first")
	case errors.Is(err, ErrSkillTooLow):
		apierr.Write(w, r, http.StatusConflict, "E_SKILL_TOO_LOW", "skill level below recipe min_skill")
	case errors.Is(err, ErrInsufficientGold):
		apierr.Write(w, r, http.StatusConflict, "E_INSUFFICIENT_GOLD", "not enough gold")
	case errors.Is(err, ErrValidation):
		apierr.Write(w, r, http.StatusBadRequest, "E_VALIDATION", "invalid input")
	default:
		if errors.Is(err, inventory.ErrInsufficient) {
			apierr.Write(w, r, http.StatusConflict, "E_INSUFFICIENT_MATERIALS", err.Error())
			return
		}
		apierr.Write(w, r, http.StatusInternalServerError, "E_INTERNAL", err.Error())
	}
}
