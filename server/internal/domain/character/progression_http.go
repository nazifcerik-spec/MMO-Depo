package character

import (
	"context"
	"errors"
	"math"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/otter/server/internal/http/apierr"
	"github.com/otter/server/internal/httpctx"
	"github.com/otter/server/internal/leveling"
)

// Summary is everything the progression UI shows (DD-250 §43).
type Summary struct {
	Level          int                `json:"level"`
	LevelCap       int                `json:"level_cap"`
	XP             int64              `json:"xp"`
	XPToNext       int64              `json:"xp_to_next"`
	XPPct          float64            `json:"xp_pct"`
	Maxed          bool               `json:"maxed"`
	PostCapXP      int64              `json:"postcap_xp,omitempty"`
	Class          string             `json:"class_code"`
	Spec           string             `json:"spec_code"`
	Path           []string           `json:"path"`
	RoleTags       []string           `json:"role_tags"`
	Promotion      PromotionView      `json:"promotion"`
	Attributes     AttributeView      `json:"attributes"`
	Talents        PointsView         `json:"talents"`
	Title          string             `json:"title,omitempty"` // highest level title reached
	Mastery        MasteryView        `json:"mastery"`
	NextMilestone  *MilestoneView     `json:"next_milestone,omitempty"`
	MinutesToLevel float64            `json:"minutes_to_next_level"`
	Unlocked       []leveling.Unlock  `json:"unlocked"`
	Resource       *leveling.Resource `json:"resource,omitempty"`
}

type PromotionView struct {
	Available  bool         `json:"available"`
	OverflowXP int64        `json:"overflow_xp"`
	Choices    []ChoiceView `json:"choices,omitempty"`
}

type ChoiceView struct {
	Code     string   `json:"code"`
	RoleTags []string `json:"role_tags"`
}

type AttributeView struct {
	PointsView
	Raw       map[string]int     `json:"raw"`
	Allocated map[string]int     `json:"allocated"`
	Effective map[string]float64 `json:"effective"`
}

type PointsView struct {
	Earned    int `json:"earned"`
	Spent     int `json:"spent"`
	Available int `json:"available"`
}

type MasteryView struct {
	Level            int    `json:"level"`
	Unlocked         bool   `json:"unlocked"`
	Title            string `json:"title,omitempty"`
	Badge            string `json:"badge,omitempty"`
	Frame            string `json:"frame,omitempty"`
	Achievement      string `json:"achievement,omitempty"`
	PrestigeEligible bool   `json:"prestige_eligible,omitempty"`
}

type MilestoneView struct {
	leveling.Milestone
	Minutes float64 `json:"estimated_minutes"`
}

// RawAttributes are a character's attributes before gear: class base +
// allocated points (equipment requirements read these).
func RawAttributes(c *Character) leveling.Stats {
	p := basePrimary(c)
	return leveling.Stats{Strength: p.Strength, Vitality: p.Vitality, Dexterity: p.Agility,
		Intelligence: p.Intelligence, Wisdom: p.Charisma, Luck: p.Luck}
}

// ProgressionSummary builds the progression view of an owned character.
func (s *Service) ProgressionSummary(ctx context.Context, id, accountID uuid.UUID) (*Summary, error) {
	c, err := s.GetOwned(ctx, id, accountID)
	if err != nil {
		return nil, err
	}
	cfg := leveling.Current()
	spec, level := c.Spec(), int(c.Level)
	out := &Summary{
		Level: level, LevelCap: cfg.LevelCap, XP: c.XP, XPToNext: cfg.XPToNext(level),
		Maxed: level >= cfg.LevelCap, PostCapXP: c.PostCapXP,
		Class: c.ClassCode, Spec: spec, Path: cfg.Path(spec),
		Unlocked: cfg.UnlockedBy(cfg.Path(spec), level),
	}
	if out.Unlocked == nil {
		out.Unlocked = []leveling.Unlock{}
	}
	if out.XPToNext > 0 {
		out.XPPct = math.Round(float64(c.XP)/float64(out.XPToNext)*1000) / 10
		out.MinutesToLevel = cfg.EstimateMinutes(level, c.XP, level+1)
	}
	if k, ok := cfg.Class(spec); ok {
		out.RoleTags = k.RoleTags
	}
	if root, ok := cfg.Class(c.ClassCode); ok {
		out.Resource = root.Resource
	}
	choices := cfg.Choices(spec, level)
	out.Promotion = PromotionView{Available: len(choices) > 0, OverflowXP: c.OverflowXP}
	for _, k := range choices {
		out.Promotion.Choices = append(out.Promotion.Choices, ChoiceView{Code: k.Code, RoleTags: k.RoleTags})
	}

	raw := RawAttributes(c)
	attr := AttributeView{Raw: map[string]int{}, Allocated: map[string]int{}, Effective: map[string]float64{}}
	for _, a := range leveling.Attributes {
		attr.Raw[a] = raw.Get(a)
		attr.Allocated[a] = c.Alloc.Get(a)
		attr.Effective[a] = math.Round(cfg.EffectiveAttribute(raw.Get(a))*10) / 10
	}
	attr.Earned = cfg.AttributePointsEarned(level)
	attr.Spent = c.Alloc.Total()
	attr.Available = max(0, attr.Earned-attr.Spent)
	out.Attributes = attr

	spent := 0
	if ps, ok := s.repo.(ProgressionStore); ok {
		ranks, err := ps.Talents(ctx, id)
		if err != nil {
			return nil, err
		}
		for _, r := range ranks {
			spent += r
		}
	}
	earned := cfg.TalentPointsEarned(level)
	out.Talents = PointsView{Earned: earned, Spent: spent, Available: max(0, earned-spent)}

	for l, t := range cfg.Titles {
		if l <= level && (out.Title == "" || l > titleLevel(cfg, out.Title)) {
			out.Title = t
		}
	}
	out.Mastery.Level = cfg.Milestones.Mastery
	if k, ok := cfg.Class(spec); ok && k.Mastery != nil {
		out.Mastery.Title, out.Mastery.Badge, out.Mastery.Frame = k.Mastery.Title, k.Mastery.Badge, k.Mastery.Frame
		out.Mastery.Unlocked = level >= cfg.Milestones.Mastery
		out.Mastery.Achievement = k.Mastery.Achievement
		out.Mastery.PrestigeEligible = out.Mastery.Unlocked && k.Mastery.PrestigeEligible
	}
	if m, ok := cfg.NextMilestone(spec, level); ok {
		out.NextMilestone = &MilestoneView{Milestone: m, Minutes: cfg.EstimateMinutes(level, c.XP, m.Level)}
	}
	return out, nil
}

func titleLevel(cfg *leveling.Config, title string) int {
	for l, t := range cfg.Titles {
		if t == title {
			return l
		}
	}
	return 0
}

// --- HTTP --------------------------------------------------------------------

func (h *Handler) owned(w http.ResponseWriter, r *http.Request) (uuid.UUID, uuid.UUID, bool) {
	id, err := parseIDParam(r, "id")
	if err != nil {
		apierr.Write(w, r, http.StatusBadRequest, "E_INVALID_ID", err.Error())
		return uuid.Nil, uuid.Nil, false
	}
	return id, httpctx.AccountID(r.Context()), true
}

// Progression serves GET /v1/character/{id}/leveling.
func (h *Handler) Progression(w http.ResponseWriter, r *http.Request) {
	id, acct, ok := h.owned(w, r)
	if !ok {
		return
	}
	sum, err := h.svc.ProgressionSummary(r.Context(), id, acct)
	if err != nil {
		writeProgErr(w, r, err)
		return
	}
	apierr.WriteJSON(w, http.StatusOK, sum)
}

// Promote serves POST /v1/character/{id}/promotion {"choice": "guardian"}.
func (h *Handler) Promote(w http.ResponseWriter, r *http.Request) {
	id, acct, ok := h.owned(w, r)
	if !ok {
		return
	}
	var body struct {
		Choice string `json:"choice"`
	}
	if err := decodeJSON(r, &body); err != nil || body.Choice == "" {
		apierr.WriteField(w, r, http.StatusBadRequest, "E_VALIDATION", "choice required", "choice")
		return
	}
	g, err := h.svc.SelectPromotion(r.Context(), id, acct, body.Choice)
	if err != nil {
		writeProgErr(w, r, err)
		return
	}
	apierr.WriteJSON(w, http.StatusOK, map[string]any{
		"spec_code": g.Character.Spec(), "level": g.Character.Level, "xp": g.Character.XP,
		"levels_gained": len(g.LevelsGained), "promotion_available": g.PromotionAvailable, "events": g.Events,
	})
}

// Allocate serves POST /v1/character/{id}/attributes {"strength": 2, ...}.
func (h *Handler) Allocate(w http.ResponseWriter, r *http.Request) {
	id, acct, ok := h.owned(w, r)
	if !ok {
		return
	}
	var body map[string]int
	if err := decodeJSON(r, &body); err != nil {
		apierr.Write(w, r, http.StatusBadRequest, "E_INVALID_JSON", "invalid body")
		return
	}
	c, err := h.svc.AllocateAttributes(r.Context(), id, acct, body)
	if err != nil {
		writeProgErr(w, r, err)
		return
	}
	alloc := map[string]int{}
	for _, a := range leveling.Attributes {
		alloc[a] = c.Alloc.Get(a)
	}
	apierr.WriteJSON(w, http.StatusOK, map[string]any{"allocated": alloc})
}

// TalentsView serves GET /v1/character/{id}/talents.
func (h *Handler) TalentsView(w http.ResponseWriter, r *http.Request) {
	id, acct, ok := h.owned(w, r)
	if !ok {
		return
	}
	ranks, err := h.svc.Talents(r.Context(), id, acct)
	if err != nil {
		writeProgErr(w, r, err)
		return
	}
	c, err := h.svc.GetOwned(r.Context(), id, acct)
	if err != nil {
		writeProgErr(w, r, err)
		return
	}
	cfg := leveling.Current()
	type talentRow struct {
		leveling.Talent
		Rank int `json:"rank"`
	}
	type treeRow struct {
		Code    string      `json:"code"`
		Points  int         `json:"points"`
		Talents []talentRow `json:"talents"`
	}
	spent := 0
	trees := []treeRow{}
	for _, t := range cfg.TreesFor(c.Spec()) {
		row := treeRow{Code: t.Code, Talents: []talentRow{}}
		for _, tl := range t.Talents {
			n := ranks[tl.Code]
			spent += n
			if !tl.Capstone {
				row.Points += n
			}
			row.Talents = append(row.Talents, talentRow{Talent: tl, Rank: n})
		}
		trees = append(trees, row)
	}
	earned := cfg.TalentPointsEarned(int(c.Level))
	apierr.WriteJSON(w, http.StatusOK, map[string]any{
		"earned": earned, "spent": spent, "available": max(0, earned-spent),
		"tree_points": cfg.TalentTreePoints, "capstone_level": cfg.Milestones.Capstone,
		"level": c.Level, "trees": trees,
	})
}

// LearnTalentHTTP serves POST /v1/character/{id}/talents/{code}.
func (h *Handler) LearnTalentHTTP(w http.ResponseWriter, r *http.Request) {
	id, acct, ok := h.owned(w, r)
	if !ok {
		return
	}
	rank, err := h.svc.LearnTalent(r.Context(), id, acct, chi.URLParam(r, "code"))
	if err != nil {
		writeProgErr(w, r, err)
		return
	}
	apierr.WriteJSON(w, http.StatusOK, map[string]any{"code": chi.URLParam(r, "code"), "rank": rank})
}

// EventsHTTP serves GET /v1/character/{id}/leveling/events?limit=.
func (h *Handler) EventsHTTP(w http.ResponseWriter, r *http.Request) {
	id, acct, ok := h.owned(w, r)
	if !ok {
		return
	}
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	list, err := h.svc.Events(r.Context(), id, acct, limit)
	if err != nil {
		writeProgErr(w, r, err)
		return
	}
	apierr.WriteJSON(w, http.StatusOK, map[string]any{"events": list})
}

// ProgressionConfig serves the progression definitions (public, cacheable;
// admin tooling reads the same): curve, XP table, rewards, classes, talents,
// unlocks.
func ProgressionConfig(w http.ResponseWriter, _ *http.Request) {
	cfg := leveling.Current()
	w.Header().Set("Cache-Control", "public, max-age=300")
	apierr.WriteJSON(w, http.StatusOK, map[string]any{"config": cfg, "xp_table": cfg.Table()})
}

// ClassNodeView is one node of the class tree for the character creation
// screen (DD-253): what a player needs to choose a class and see its paths.
type ClassNodeView struct {
	Code                string               `json:"code"`
	Parent              string               `json:"parent,omitempty"`
	RequiredLevel       int                  `json:"required_level,omitempty"`
	Tier                int                  `json:"tier"`
	RoleTags            []string             `json:"role_tags"`
	Resource            *leveling.Resource   `json:"resource,omitempty"`
	BaseStats           *leveling.Stats      `json:"base_stats,omitempty"`
	PrimaryAttribute    string               `json:"primary_attribute,omitempty"`
	SecondaryAttributes []string             `json:"secondary_attributes,omitempty"`
	Equipment           []string             `json:"equipment"`
	SupportPower        float64              `json:"support_power,omitempty"`
	AFKProfile          *leveling.AFKProfile `json:"afk_profile,omitempty"`
	StrongIn            []string             `json:"strong_in,omitempty"`
	WeakIn              []string             `json:"weak_in,omitempty"`
	MasteryTitle        string               `json:"mastery_title,omitempty"`
}

// ClassCatalog serves GET /v1/classes: the class tree with base stats,
// resources, effective equipment permissions and specialization metadata.
func ClassCatalog(w http.ResponseWriter, _ *http.Request) {
	cfg := leveling.Current()
	out := make([]ClassNodeView, 0, len(cfg.Classes))
	for _, c := range cfg.Classes {
		k, _ := cfg.Class(c.Code)
		v := ClassNodeView{Code: k.Code, Parent: k.Parent, RequiredLevel: k.RequiredLevel, Tier: k.Tier(),
			RoleTags: k.RoleTags, Resource: k.Resource, BaseStats: k.BaseStats, PrimaryAttribute: k.PrimaryAttribute,
			SecondaryAttributes: k.SecondaryAttributes, Equipment: cfg.Permissions(k.Code), SupportPower: k.SupportPower,
			AFKProfile: k.AFKProfile, StrongIn: k.StrongIn, WeakIn: k.WeakIn}
		if k.Mastery != nil {
			v.MasteryTitle = k.Mastery.Title
		}
		if v.Equipment == nil {
			v.Equipment = []string{}
		}
		out = append(out, v)
	}
	w.Header().Set("Cache-Control", "public, max-age=300")
	apierr.WriteJSON(w, http.StatusOK, map[string]any{"classes": out,
		"attribute_points_total": cfg.AttributePointsEarned(cfg.LevelCap), "milestones": cfg.Milestones})
}

func writeProgErr(w http.ResponseWriter, r *http.Request, err error) {
	switch {
	case errors.Is(err, ErrInvalidPromotion):
		apierr.Write(w, r, http.StatusConflict, "E_PROMOTION_UNAVAILABLE", "that promotion is not available")
	case errors.Is(err, ErrNotEnoughPoints):
		apierr.Write(w, r, http.StatusConflict, "E_NOT_ENOUGH_POINTS", "not enough points")
	case errors.Is(err, ErrUnknownAttribute):
		apierr.Write(w, r, http.StatusBadRequest, "E_UNKNOWN_ATTRIBUTE", "unknown attribute")
	case errors.Is(err, ErrUnknownTalent):
		apierr.Write(w, r, http.StatusNotFound, "E_UNKNOWN_TALENT", "unknown talent")
	case errors.Is(err, ErrTalentMaxed):
		apierr.Write(w, r, http.StatusConflict, "E_TALENT_MAXED", "talent at max rank")
	case errors.Is(err, ErrCapstoneLocked):
		apierr.Write(w, r, http.StatusConflict, "E_CAPSTONE_LOCKED", "capstone needs the capstone level and a full tree")
	case errors.Is(err, ErrBusy):
		apierr.Write(w, r, http.StatusConflict, "E_BUSY_COMBAT", "not while in combat, AFK or a dungeon")
	case errors.Is(err, ErrInvalidSpec):
		apierr.Write(w, r, http.StatusConflict, "E_SPEC_UNAVAILABLE", "that specialization is not available")
	case errors.Is(err, ErrProgressionUnsupp):
		apierr.Write(w, r, http.StatusServiceUnavailable, "E_UNAVAILABLE", "progression unavailable")
	default:
		writeCharErr(w, r, err)
	}
}
