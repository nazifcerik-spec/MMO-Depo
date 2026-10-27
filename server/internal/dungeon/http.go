package dungeon

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

// EventPublisher is the narrow ws-fanout surface used to notify the account's
// live sessions when a dungeon reveals. Nil-safe — set via WithPublisher.
type EventPublisher interface {
	Publish(accountID uuid.UUID, kind string, payload any) int
}

type Handler struct {
	svc *Service
	log *slog.Logger
	pub EventPublisher
}

func NewHandler(svc *Service, logger *slog.Logger) *Handler {
	if logger == nil {
		logger = slog.Default()
	}
	return &Handler{svc: svc, log: logger}
}

// WithPublisher wires the ws event pipe.
func (h *Handler) WithPublisher(p EventPublisher) *Handler {
	h.pub = p
	return h
}

func (h *Handler) publish(accountID uuid.UUID, kind string, payload any) {
	if h.pub == nil {
		return
	}
	h.pub.Publish(accountID, kind, payload)
}

type startBody struct {
	DungeonCode string `json:"dungeon_code"`
}

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
		CharacterID: charID, AccountID: acct, DungeonCode: body.DungeonCode,
	})
	if err != nil {
		writeErr(w, r, err)
		return
	}
	apierr.WriteJSON(w, http.StatusCreated, res)
}

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
	h.publish(acct, "dungeon.reveal", map[string]any{
		"contract_id":      contractID.String(),
		"cleared":          res.Run.Cleared,
		"boss_killed":      res.Run.BossKilled,
		"gold_copper":      res.Run.TotalGoldCopper,
		"xp":               res.Run.TotalXP,
		"lockout":          res.Lockout,
		"week_key":         res.WeekKey,
		"stopped_on_death": res.Run.StoppedOnDeath,
	})
	apierr.WriteJSON(w, http.StatusOK, res)
}

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
	apierr.WriteJSON(w, http.StatusOK, map[string]any{"contract": summarize(c)})
}

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
	apierr.WriteJSON(w, http.StatusOK, map[string]any{"count": len(out), "runs": out})
}

// Catalog is a public read: list of available dungeons + preview counts.
func (h *Handler) Catalog(w http.ResponseWriter, r *http.Request) {
	out := make([]map[string]any, 0, 8)
	for _, t := range h.svc.ListCatalog() {
		out = append(out, map[string]any{
			"code":              t.Code,
			"display_name_key":  t.DisplayNameKey,
			"min_level":         t.MinLevel,
			"duration_seconds":  t.DurationSecs,
			"encounters":        t.TotalEncounters(),
			"boss_mob_code":     t.BossMobCode,
			"bonus_gold_copper": t.BonusGoldCopper,
			"bonus_xp":          t.BonusXP,
			"xp":                t.BaseXP(), // at its own level; members' levels scale it (DD-303)
			"region_code":       t.RegionCode,
		})
	}
	apierr.WriteJSON(w, http.StatusOK, map[string]any{"count": len(out), "dungeons": out})
}

func summarize(c *Contract) map[string]any {
	m := map[string]any{
		"id":               c.ID.String(),
		"character_id":     c.CharacterID.String(),
		"dungeon_code":     c.DungeonCode,
		"duration_seconds": c.DurationSecs,
		"started_at":       c.StartedAt,
		"expires_at":       c.ExpiresAt,
		"state":            string(c.State),
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
		apierr.Write(w, r, http.StatusNotFound, "E_NOT_FOUND", "dungeon run not found")
	case errors.Is(err, ErrForbidden), errors.Is(err, character.ErrForbidden), errors.Is(err, character.ErrNotFound):
		apierr.Write(w, r, http.StatusNotFound, "E_NOT_FOUND", "dungeon run not found")
	case errors.Is(err, ErrSlotBusy):
		apierr.Write(w, r, http.StatusConflict, "E_DUNGEON_ACTIVE", "character already has an active dungeon run")
	case errors.Is(err, ErrNotActive):
		apierr.Write(w, r, http.StatusConflict, "E_DUNGEON_NOT_ACTIVE", "dungeon run is not active")
	case errors.Is(err, ErrNotReady):
		apierr.Write(w, r, http.StatusConflict, "E_DUNGEON_NOT_READY", "not enough time elapsed to reveal")
	case errors.Is(err, ErrWeeklyLockout):
		apierr.Write(w, r, http.StatusConflict, "E_DUNGEON_WEEKLY_LOCKED", "already cleared this dungeon this week")
	case errors.Is(err, ErrUnknownDungeon):
		apierr.Write(w, r, http.StatusNotFound, "E_UNKNOWN_DUNGEON", "dungeon_code not in catalog")
	case errors.Is(err, ErrWrongRegion):
		apierr.Write(w, r, http.StatusConflict, "E_DUNGEON_REGION", "enter this dungeon from its region")
	case errors.Is(err, ErrLevelRequirement):
		apierr.Write(w, r, http.StatusForbidden, "E_LEVEL_REQUIREMENT", "character level below min_level")
	case errors.Is(err, ErrBagFull):
		apierr.Write(w, r, http.StatusConflict, "E_BAG_FULL", "bag is full")
	case errors.Is(err, ErrTooWounded):
		apierr.Write(w, r, http.StatusConflict, "E_TOO_WOUNDED", "rest or heal before fighting")
	case errors.Is(err, ErrUnknownMob):
		apierr.Write(w, r, http.StatusInternalServerError, "E_CONTENT_BROKEN", err.Error())
	case errors.Is(err, ErrValidation):
		apierr.Write(w, r, http.StatusBadRequest, "E_VALIDATION", "invalid input")
	default:
		apierr.Write(w, r, http.StatusInternalServerError, "E_INTERNAL", err.Error())
	}
}
