package world

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

type Handler struct {
	svc     *Service
	catalog CatalogLookup
	log     *slog.Logger
}

func NewHandler(svc *Service, catalog CatalogLookup, logger *slog.Logger) *Handler {
	if logger == nil {
		logger = slog.Default()
	}
	return &Handler{svc: svc, catalog: catalog, log: logger}
}

// --- Public catalog reads ---------------------------------------------------

// Catalog returns the full region list + explicit links. Public read — clients
// use this to render a world map picker.
func (h *Handler) Catalog(w http.ResponseWriter, r *http.Request) {
	cat := h.catalog.Catalog()
	regions := cat.All()
	out := make([]map[string]any, 0, len(regions))
	for _, reg := range regions {
		row := regionRow(reg)
		out = append(out, row)
	}
	// Roads between regions besides parent/child (DD-279): the map draws them.
	links := make([]map[string]any, 0, len(cat.Links()))
	for _, l := range cat.Links() {
		links = append(links, map[string]any{"a": l.RegionA, "b": l.RegionB, "travel_seconds": l.TravelSeconds})
	}
	apierr.WriteJSON(w, http.StatusOK, map[string]any{
		"count":         len(out),
		"regions":       out,
		"links":         links,
		"starting_code": cat.StartingCode(),
	})
}

// --- Character-scoped -------------------------------------------------------

// Location returns the character's current region + neighbors + active travel.
func (h *Handler) Location(w http.ResponseWriter, r *http.Request) {
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
	view, err := h.svc.Location(r.Context(), charID, acct)
	if err != nil {
		writeErr(w, r, err)
		return
	}
	neighbors := make([]map[string]any, 0, len(view.Neighbors))
	for _, n := range view.Neighbors {
		row := regionRow(n.Region)
		row["travel_seconds"] = n.TravelSeconds
		row["via"] = n.Via
		neighbors = append(neighbors, row)
	}
	resp := map[string]any{
		"region":        regionRow(view.Region),
		"neighbors":     neighbors,
		"hunt_regions":  view.HuntRegions,
		"routes":        view.Routes,
		"starting_code": view.StartingCode,
	}
	if view.Travel != nil {
		row := travelRow(view.Travel, h.svc.cfg.Now())
		// The road the journey takes, for the map (DD-280).
		if r, ok := h.catalog.Catalog().Routes(view.Travel.FromRegionCode)[view.Travel.ToRegionCode]; ok {
			row["path"] = r.Path
		}
		resp["travel"] = row
	}
	apierr.WriteJSON(w, http.StatusOK, resp)
}

type startBody struct {
	ToRegionCode string `json:"to_region_code"`
}

// Start begins a travel from current region to `to_region_code`.
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
		CharacterID: charID, AccountID: acct, ToRegionCode: body.ToRegionCode,
	})
	if err != nil {
		writeErr(w, r, err)
		return
	}
	apierr.WriteJSON(w, http.StatusCreated, map[string]any{
		"travel":           travelRow(res.Travel, h.svc.cfg.Now()),
		"from_region_code": res.FromCode,
		"duration_seconds": res.Duration,
	})
}

// Complete arrives at the destination (rejected before arrives_at).
func (h *Handler) Complete(w http.ResponseWriter, r *http.Request) {
	acct := httpctx.AccountID(r.Context())
	if acct == uuid.Nil {
		apierr.Write(w, r, http.StatusUnauthorized, "E_UNAUTHENTICATED", "auth required")
		return
	}
	travelID, err := uuid.Parse(chi.URLParam(r, "travel_id"))
	if err != nil {
		apierr.Write(w, r, http.StatusBadRequest, "E_INVALID_ID", err.Error())
		return
	}
	t, err := h.svc.Complete(r.Context(), travelID, acct)
	if err != nil {
		writeErr(w, r, err)
		return
	}
	apierr.WriteJSON(w, http.StatusOK, map[string]any{
		"travel": travelRow(t, h.svc.cfg.Now()),
	})
}

// Cancel aborts an in-progress travel.
func (h *Handler) Cancel(w http.ResponseWriter, r *http.Request) {
	acct := httpctx.AccountID(r.Context())
	if acct == uuid.Nil {
		apierr.Write(w, r, http.StatusUnauthorized, "E_UNAUTHENTICATED", "auth required")
		return
	}
	travelID, err := uuid.Parse(chi.URLParam(r, "travel_id"))
	if err != nil {
		apierr.Write(w, r, http.StatusBadRequest, "E_INVALID_ID", err.Error())
		return
	}
	t, err := h.svc.Cancel(r.Context(), travelID, acct)
	if err != nil {
		writeErr(w, r, err)
		return
	}
	apierr.WriteJSON(w, http.StatusOK, map[string]any{
		"travel": travelRow(t, h.svc.cfg.Now()),
	})
}

// History lists past travels (newest first).
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
	now := h.svc.cfg.Now()
	for i := range list {
		out = append(out, travelRow(&list[i], now))
	}
	apierr.WriteJSON(w, http.StatusOK, map[string]any{"count": len(out), "travels": out})
}

// --- shared row builders ---------------------------------------------------

func regionRow(r Region) map[string]any {
	row := map[string]any{
		"code":                       r.Code,
		"display_name_key":           r.DisplayNameKey,
		"kind":                       string(r.Kind),
		"min_level":                  r.MinLevel,
		"travel_seconds_from_parent": r.TravelSecondsFromParent,
		"safe":                       r.Safe,
		"is_starting":                r.IsStarting,
		"map_x":                      r.MapX,
		"map_y":                      r.MapY,
	}
	if r.ParentCode != nil {
		row["parent_code"] = *r.ParentCode
	}
	return row
}

func travelRow(t *Travel, _ any) map[string]any {
	return map[string]any{
		"id":               t.ID.String(),
		"character_id":     t.CharacterID.String(),
		"from_region_code": t.FromRegionCode,
		"to_region_code":   t.ToRegionCode,
		"started_at":       t.StartedAt,
		"arrives_at":       t.ArrivesAt,
		"state":            string(t.State),
	}
}

func writeErr(w http.ResponseWriter, r *http.Request, err error) {
	if action.Write(w, r, err) {
		return
	}
	switch {
	case errors.Is(err, ErrNotFound):
		apierr.Write(w, r, http.StatusNotFound, "E_NOT_FOUND", "travel not found")
	case errors.Is(err, ErrForbidden), errors.Is(err, character.ErrForbidden), errors.Is(err, character.ErrNotFound):
		apierr.Write(w, r, http.StatusNotFound, "E_NOT_FOUND", "travel not found")
	case errors.Is(err, ErrUnknownRegion):
		apierr.Write(w, r, http.StatusNotFound, "E_UNKNOWN_REGION", "region_code not in catalog")
	case errors.Is(err, ErrTravelInProgress):
		apierr.Write(w, r, http.StatusConflict, "E_TRAVEL_IN_PROGRESS", "character already has an active travel")
	case errors.Is(err, ErrTravelNotActive):
		apierr.Write(w, r, http.StatusConflict, "E_TRAVEL_NOT_ACTIVE", "travel is not in progress")
	case errors.Is(err, ErrTravelNotReady):
		apierr.Write(w, r, http.StatusConflict, "E_TRAVEL_NOT_READY", "not yet arrived")
	case errors.Is(err, ErrSameRegion):
		apierr.Write(w, r, http.StatusBadRequest, "E_SAME_REGION", "already at destination")
	case errors.Is(err, ErrNoRoute):
		apierr.Write(w, r, http.StatusBadRequest, "E_NO_ROUTE", "no road leads to destination")
	case errors.Is(err, ErrLevelRequirement):
		apierr.Write(w, r, http.StatusForbidden, "E_LEVEL_REQUIREMENT", "character level below region min_level")
	case errors.Is(err, ErrNoStartingRegion):
		apierr.Write(w, r, http.StatusInternalServerError, "E_NO_STARTING_REGION", "world content missing starting region")
	case errors.Is(err, ErrValidation):
		apierr.Write(w, r, http.StatusBadRequest, "E_VALIDATION", "invalid input")
	default:
		apierr.Write(w, r, http.StatusInternalServerError, "E_INTERNAL", err.Error())
	}
}
