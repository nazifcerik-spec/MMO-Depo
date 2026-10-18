package crafting

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/otter/server/internal/action"
	"github.com/otter/server/internal/domain/character"
	"github.com/otter/server/internal/http/apierr"
	"github.com/otter/server/internal/httpctx"
	"github.com/otter/server/internal/inventory"
	"github.com/otter/server/internal/item"
)

type Handler struct {
	svc   *Service
	log   *slog.Logger
	items ItemCatalog // nil = recipes without output details
}

// ItemCatalog describes recipe outputs (item.Service satisfies).
type ItemCatalog interface {
	Catalog() *item.Catalog
}

// WithItems adds each recipe's output properties to GET /recipes, so the
// workshop can show what an item does before it is made.
func (h *Handler) WithItems(c ItemCatalog) *Handler {
	h.items = c
	return h
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
	var cat *item.Catalog
	if h.items != nil {
		cat = h.items.Catalog()
	}
	for _, rec := range all {
		if filterProf != "" && rec.ProfessionCode != filterProf {
			continue
		}
		row := recipeRow(rec)
		if cat != nil {
			if t, err := cat.Get(rec.OutputTemplateCode); err == nil {
				row["output"] = outputRow(t)
			}
		}
		out = append(out, row)
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

// Craft starts a timed job (DD-264): 201 with the job.
func (h *Handler) Craft(w http.ResponseWriter, r *http.Request) {
	acct, charID, ok := ids(w, r)
	if !ok {
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
	job, err := h.svc.Start(r.Context(), charID, acct, body.RecipeCode, body.Count)
	if err != nil {
		writeErr(w, r, err)
		return
	}
	apierr.WriteJSON(w, http.StatusCreated, map[string]any{"job": jobView(job, h.svc.cfg.Now())})
}

// ActiveJob is the running job, or {"job": null}.
func (h *Handler) ActiveJob(w http.ResponseWriter, r *http.Request) {
	acct, charID, ok := ids(w, r)
	if !ok {
		return
	}
	job, err := h.svc.Active(r.Context(), charID, acct)
	if err != nil {
		writeErr(w, r, err)
		return
	}
	var view map[string]any
	if job != nil {
		view = jobView(job, h.svc.cfg.Now())
	}
	apierr.WriteJSON(w, http.StatusOK, map[string]any{"job": view})
}

// Collect ends a job: finished units delivered, the rest refunded.
func (h *Handler) Collect(w http.ResponseWriter, r *http.Request) {
	acct, charID, ok := ids(w, r)
	if !ok {
		return
	}
	jobID, err := uuid.Parse(chi.URLParam(r, "job_id"))
	if err != nil {
		apierr.Write(w, r, http.StatusBadRequest, "E_INVALID_ID", "invalid job id")
		return
	}
	res, err := h.svc.Collect(r.Context(), charID, acct, jobID)
	if err != nil {
		writeErr(w, r, err)
		return
	}
	apierr.WriteJSON(w, http.StatusOK, res)
}

func jobView(j *Job, now time.Time) map[string]any {
	return map[string]any{
		"id":                   j.ID.String(),
		"recipe_code":          j.Recipe.Code,
		"profession_code":      j.Recipe.ProfessionCode,
		"output_template_code": j.Recipe.OutputTemplateCode,
		"output_stack":         j.Recipe.OutputStack,
		"count":                j.Count,
		"unit_seconds":         j.UnitSeconds,
		"started_at":           j.StartedAt,
		"ready_at":             j.ReadyAt,
		"done":                 j.Done(now),
		"state":                j.State,
	}
}

func ids(w http.ResponseWriter, r *http.Request) (acct, char uuid.UUID, ok bool) {
	acct = httpctx.AccountID(r.Context())
	if acct == uuid.Nil {
		apierr.Write(w, r, http.StatusUnauthorized, "E_UNAUTHENTICATED", "auth required")
		return acct, char, false
	}
	char, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		apierr.Write(w, r, http.StatusBadRequest, "E_INVALID_ID", err.Error())
		return acct, char, false
	}
	return acct, char, true
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
		"craft_seconds":        r.CraftSeconds,
	}
}

// outputRow is what a recipe makes: class, requirements, bonuses, effect.
func outputRow(t item.Template) map[string]any {
	row := map[string]any{
		"item_class":        t.ItemClass,
		"item_subclass":     t.ItemSubclass,
		"min_level":         t.MinLevel,
		"base_value_copper": t.BaseValueCopper,
		"weight_grams":      t.BaseWeightG,
		"max_stack":         t.MaxStack,
		"equip_slot":        t.EquipSlot,
		"class_restriction": t.ClassRestriction,
	}
	if len(t.StatBonuses) > 0 {
		row["stat_bonuses"] = json.RawMessage(t.StatBonuses)
	}
	if t.UseEffect != nil {
		row["use_effect"] = t.UseEffect
	}
	return row
}

func writeErr(w http.ResponseWriter, r *http.Request, err error) {
	if action.Write(w, r, err) {
		return
	}
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
	case errors.Is(err, ErrJobActive):
		apierr.Write(w, r, http.StatusConflict, "E_CRAFT_ACTIVE", "a crafting job is already running")
	case errors.Is(err, ErrJobTooLong):
		apierr.Write(w, r, http.StatusBadRequest, "E_CRAFT_TOO_LONG", "job longer than 6 hours")
	case errors.Is(err, ErrValidation):
		apierr.Write(w, r, http.StatusBadRequest, "E_VALIDATION", "invalid input")
	case errors.Is(err, inventory.ErrBagFull):
		apierr.Write(w, r, http.StatusConflict, "E_BAG_FULL", "bag is full")
	default:
		if errors.Is(err, inventory.ErrInsufficient) {
			apierr.Write(w, r, http.StatusConflict, "E_INSUFFICIENT_MATERIALS", err.Error())
			return
		}
		apierr.Write(w, r, http.StatusInternalServerError, "E_INTERNAL", err.Error())
	}
}
