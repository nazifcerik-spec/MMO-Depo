package raid

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/otter/server/internal/action"
	"github.com/otter/server/internal/domain/character"
	"github.com/otter/server/internal/http/apierr"
	"github.com/otter/server/internal/httpctx"
)

// Handler serves the raid endpoints. Character routes sit under
// /character/{id} (ownership enforced by the router).
type Handler struct{ svc *Service }

func NewHandler(svc *Service) *Handler { return &Handler{svc: svc} }

// Catalog is the public raid list.
func (h *Handler) Catalog(w http.ResponseWriter, _ *http.Request) {
	list := h.svc.Catalog()
	out := make([]map[string]any, 0, len(list))
	for _, t := range list {
		out = append(out, map[string]any{
			"code": t.Code, "display_name_key": t.DisplayNameKey, "min_level": t.MinLevel,
			"duration_seconds": t.DurationSecs, "encounters": t.TotalEncounters(), "boss_mob_code": t.BossMobCode,
			"bonus_gold_copper": t.BonusGoldCopper, "bonus_xp": t.BonusXP,
			"group_min": t.GroupMin, "group_max": t.GroupMax, "boss_gear": BossGear,
		})
	}
	apierr.WriteJSON(w, http.StatusOK, map[string]any{"count": len(out), "raids": out})
}

func (h *Handler) Current(w http.ResponseWriter, r *http.Request) {
	acct, char, ok := ids(w, r)
	if !ok {
		return
	}
	cur, err := h.svc.Current(r.Context(), char, acct)
	if err != nil {
		writeErr(w, r, err)
		return
	}
	apierr.WriteJSON(w, http.StatusOK, cur)
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	acct, char, ok := ids(w, r)
	if !ok {
		return
	}
	var body struct {
		RaidCode string `json:"raid_code"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		apierr.Write(w, r, http.StatusBadRequest, "E_INVALID_JSON", err.Error())
		return
	}
	run, err := h.svc.Create(r.Context(), char, acct, body.RaidCode)
	if err != nil {
		writeErr(w, r, err)
		return
	}
	apierr.WriteJSON(w, http.StatusCreated, map[string]any{"raid": run})
}

// Join, Leave, Start and Reveal act on {raid_id}.
func (h *Handler) Join(w http.ResponseWriter, r *http.Request) {
	h.runAction(w, r, func(acct, char, run uuid.UUID) (any, error) {
		res, err := h.svc.Join(r.Context(), char, acct, run)
		return map[string]any{"raid": res}, err
	})
}

func (h *Handler) Leave(w http.ResponseWriter, r *http.Request) {
	h.runAction(w, r, func(acct, char, run uuid.UUID) (any, error) {
		return map[string]any{"left": true}, h.svc.Leave(r.Context(), char, acct, run)
	})
}

func (h *Handler) Start(w http.ResponseWriter, r *http.Request) {
	h.runAction(w, r, func(acct, char, run uuid.UUID) (any, error) {
		res, err := h.svc.Start(r.Context(), char, acct, run)
		return map[string]any{"raid": res}, err
	})
}

func (h *Handler) Reveal(w http.ResponseWriter, r *http.Request) {
	h.runAction(w, r, func(acct, char, run uuid.UUID) (any, error) {
		return h.svc.Reveal(r.Context(), char, acct, run)
	})
}

func (h *Handler) runAction(w http.ResponseWriter, r *http.Request, fn func(acct, char, run uuid.UUID) (any, error)) {
	acct, char, ok := ids(w, r)
	if !ok {
		return
	}
	run, err := uuid.Parse(chi.URLParam(r, "raid_id"))
	if err != nil {
		apierr.Write(w, r, http.StatusBadRequest, "E_INVALID_ID", "invalid raid id")
		return
	}
	out, err := fn(acct, char, run)
	if err != nil {
		writeErr(w, r, err)
		return
	}
	apierr.WriteJSON(w, http.StatusOK, out)
}

func ids(w http.ResponseWriter, r *http.Request) (acct, char uuid.UUID, ok bool) {
	acct = httpctx.AccountID(r.Context())
	if acct == uuid.Nil {
		apierr.Write(w, r, http.StatusUnauthorized, "E_UNAUTHENTICATED", "auth required")
		return acct, char, false
	}
	char, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		apierr.Write(w, r, http.StatusBadRequest, "E_INVALID_ID", "invalid character id")
		return acct, char, false
	}
	return acct, char, true
}

func writeErr(w http.ResponseWriter, r *http.Request, err error) {
	if action.Write(w, r, err) {
		return
	}
	type mapping struct {
		status int
		code   string
	}
	for e, m := range map[error]mapping{
		ErrUnknownRaid:         {http.StatusNotFound, "E_UNKNOWN_RAID"},
		ErrNotFound:            {http.StatusNotFound, "E_NOT_FOUND"},
		ErrNoParty:             {http.StatusConflict, "E_NOT_IN_PARTY"},
		ErrNotInParty:          {http.StatusConflict, "E_NOT_IN_PARTY"},
		ErrNotLeader:           {http.StatusForbidden, "E_NOT_LEADER"},
		ErrLobbyExists:         {http.StatusConflict, "E_RAID_EXISTS"},
		ErrAlreadyInRaid:       {http.StatusConflict, "E_ALREADY_IN_RAID"},
		ErrNotForming:          {http.StatusConflict, "E_RAID_NOT_OPEN"},
		ErrRaidFull:            {http.StatusConflict, "E_RAID_FULL"},
		ErrNotEnough:           {http.StatusConflict, "E_RAID_NOT_ENOUGH"},
		ErrNotRunning:          {http.StatusConflict, "E_RAID_NOT_RUNNING"},
		ErrNotReady:            {http.StatusConflict, "E_RAID_NOT_READY"},
		ErrNotMember:           {http.StatusNotFound, "E_NOT_FOUND"},
		ErrLevelRequirement:    {http.StatusForbidden, "E_LEVEL_REQUIREMENT"},
		ErrWeeklyLockout:       {http.StatusConflict, "E_RAID_WEEKLY_LOCKED"},
		ErrTooWounded:          {http.StatusConflict, "E_TOO_WOUNDED"},
		ErrBagFull:             {http.StatusConflict, "E_BAG_FULL"},
		character.ErrForbidden: {http.StatusNotFound, "E_NOT_FOUND"},
		character.ErrNotFound:  {http.StatusNotFound, "E_NOT_FOUND"},
	} {
		if errors.Is(err, e) {
			apierr.Write(w, r, m.status, m.code, err.Error())
			return
		}
	}
	apierr.Write(w, r, http.StatusInternalServerError, "E_INTERNAL", err.Error())
}
