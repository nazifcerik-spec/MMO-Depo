package guild

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/otter/server/internal/domain/character"
	"github.com/otter/server/internal/economy"
	"github.com/otter/server/internal/http/apierr"
	"github.com/otter/server/internal/httpctx"
)

// Handler serves /v1/character/{id}/guild/* and GET /v1/guilds/{guild_id}.
type Handler struct{ svc *Service }

func NewHandler(svc *Service) *Handler { return &Handler{svc: svc} }

func scope(w http.ResponseWriter, r *http.Request) (uuid.UUID, uuid.UUID, bool) {
	acct := httpctx.AccountID(r.Context())
	if acct == uuid.Nil {
		apierr.Write(w, r, http.StatusUnauthorized, "E_UNAUTHENTICATED", "auth required")
		return uuid.Nil, uuid.Nil, false
	}
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		apierr.Write(w, r, http.StatusBadRequest, "E_INVALID_ID", err.Error())
		return uuid.Nil, uuid.Nil, false
	}
	return id, acct, true
}

func decode(w http.ResponseWriter, r *http.Request, v any) bool {
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		apierr.Write(w, r, http.StatusBadRequest, "E_INVALID_JSON", err.Error())
		return false
	}
	return true
}

func pathUUID(w http.ResponseWriter, r *http.Request, key string) (uuid.UUID, bool) {
	id, err := uuid.Parse(chi.URLParam(r, key))
	if err != nil {
		apierr.Write(w, r, http.StatusBadRequest, "E_INVALID_ID", key+": "+err.Error())
		return uuid.Nil, false
	}
	return id, true
}

type body struct {
	Name              string    `json:"name"`
	Tag               string    `json:"tag"`
	Description       string    `json:"description"`
	TargetName        string    `json:"target_name"`
	TargetCharacterID uuid.UUID `json:"target_character_id"`
	RankIndex         *int16    `json:"rank_index"`
	MOTD              string    `json:"motd"`
}

func (h *Handler) Mine(w http.ResponseWriter, r *http.Request) {
	actor, acct, ok := scope(w, r)
	if !ok {
		return
	}
	g, err := h.svc.Mine(r.Context(), actor, acct)
	if errors.Is(err, ErrNotInGuild) {
		apierr.WriteJSON(w, http.StatusOK, map[string]any{"guild": nil})
		return
	}
	if err != nil {
		writeErr(w, r, err)
		return
	}
	apierr.WriteJSON(w, http.StatusOK, map[string]any{"guild": g})
}

func (h *Handler) Public(w http.ResponseWriter, r *http.Request) {
	id, ok := pathUUID(w, r, "guild_id")
	if !ok {
		return
	}
	g, err := h.svc.Public(r.Context(), id)
	if err != nil {
		writeErr(w, r, err)
		return
	}
	apierr.WriteJSON(w, http.StatusOK, map[string]any{"guild": g})
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	actor, acct, ok := scope(w, r)
	if !ok {
		return
	}
	var b body
	if !decode(w, r, &b) {
		return
	}
	g, err := h.svc.Create(r.Context(), actor, acct, CreateRequest{
		Name: b.Name, Tag: b.Tag, Description: b.Description,
		IdempotencyKey: r.Header.Get("Idempotency-Key"),
	})
	if err != nil {
		writeErr(w, r, err)
		return
	}
	apierr.WriteJSON(w, http.StatusCreated, map[string]any{"guild": g, "fee_copper": CreationFeeCopper})
}

func (h *Handler) Invite(w http.ResponseWriter, r *http.Request) {
	actor, acct, ok := scope(w, r)
	if !ok {
		return
	}
	var b body
	if !decode(w, r, &b) {
		return
	}
	inv, err := h.svc.Invite(r.Context(), actor, acct, b.TargetName)
	if err != nil {
		writeErr(w, r, err)
		return
	}
	apierr.WriteJSON(w, http.StatusCreated, inv)
}

func (h *Handler) Invites(w http.ResponseWriter, r *http.Request) {
	actor, acct, ok := scope(w, r)
	if !ok {
		return
	}
	list, err := h.svc.Invites(r.Context(), actor, acct)
	if err != nil {
		writeErr(w, r, err)
		return
	}
	apierr.WriteJSON(w, http.StatusOK, map[string]any{"count": len(list), "invites": list})
}

func (h *Handler) Accept(w http.ResponseWriter, r *http.Request) {
	actor, acct, ok := scope(w, r)
	if !ok {
		return
	}
	inviteID, ok := pathUUID(w, r, "invite_id")
	if !ok {
		return
	}
	g, err := h.svc.Accept(r.Context(), actor, acct, inviteID)
	if err != nil {
		writeErr(w, r, err)
		return
	}
	apierr.WriteJSON(w, http.StatusOK, map[string]any{"guild": g})
}

func (h *Handler) Decline(w http.ResponseWriter, r *http.Request) {
	actor, acct, ok := scope(w, r)
	if !ok {
		return
	}
	inviteID, ok := pathUUID(w, r, "invite_id")
	if !ok {
		return
	}
	if err := h.svc.Decline(r.Context(), actor, acct, inviteID); err != nil {
		writeErr(w, r, err)
		return
	}
	apierr.WriteJSON(w, http.StatusOK, map[string]any{"declined": true})
}

func (h *Handler) Leave(w http.ResponseWriter, r *http.Request) {
	actor, acct, ok := scope(w, r)
	if !ok {
		return
	}
	if err := h.svc.Leave(r.Context(), actor, acct); err != nil {
		writeErr(w, r, err)
		return
	}
	apierr.WriteJSON(w, http.StatusOK, map[string]any{"left": true})
}

func (h *Handler) Kick(w http.ResponseWriter, r *http.Request) {
	actor, acct, ok := scope(w, r)
	if !ok {
		return
	}
	var b body
	if !decode(w, r, &b) {
		return
	}
	if err := h.svc.Kick(r.Context(), actor, acct, b.TargetCharacterID); err != nil {
		writeErr(w, r, err)
		return
	}
	apierr.WriteJSON(w, http.StatusOK, map[string]any{"kicked": b.TargetCharacterID})
}

func (h *Handler) SetRank(w http.ResponseWriter, r *http.Request) {
	actor, acct, ok := scope(w, r)
	if !ok {
		return
	}
	var b body
	if !decode(w, r, &b) {
		return
	}
	if b.RankIndex == nil {
		apierr.Write(w, r, http.StatusBadRequest, "E_VALIDATION", "rank_index required")
		return
	}
	if err := h.svc.SetRank(r.Context(), actor, acct, b.TargetCharacterID, *b.RankIndex); err != nil {
		writeErr(w, r, err)
		return
	}
	apierr.WriteJSON(w, http.StatusOK, map[string]any{"character_id": b.TargetCharacterID, "rank_index": *b.RankIndex})
}

func (h *Handler) Transfer(w http.ResponseWriter, r *http.Request) {
	actor, acct, ok := scope(w, r)
	if !ok {
		return
	}
	var b body
	if !decode(w, r, &b) {
		return
	}
	if err := h.svc.TransferLeadership(r.Context(), actor, acct, b.TargetCharacterID); err != nil {
		writeErr(w, r, err)
		return
	}
	apierr.WriteJSON(w, http.StatusOK, map[string]any{"leader_character_id": b.TargetCharacterID})
}

func (h *Handler) MOTD(w http.ResponseWriter, r *http.Request) {
	actor, acct, ok := scope(w, r)
	if !ok {
		return
	}
	var b body
	if !decode(w, r, &b) {
		return
	}
	if err := h.svc.SetMOTD(r.Context(), actor, acct, b.MOTD); err != nil {
		writeErr(w, r, err)
		return
	}
	apierr.WriteJSON(w, http.StatusOK, map[string]any{"motd": b.MOTD})
}

func (h *Handler) Disband(w http.ResponseWriter, r *http.Request) {
	actor, acct, ok := scope(w, r)
	if !ok {
		return
	}
	at, err := h.svc.ScheduleDisband(r.Context(), actor, acct)
	if err != nil {
		writeErr(w, r, err)
		return
	}
	apierr.WriteJSON(w, http.StatusOK, map[string]any{"disband_scheduled_at": at})
}

func (h *Handler) CancelDisband(w http.ResponseWriter, r *http.Request) {
	actor, acct, ok := scope(w, r)
	if !ok {
		return
	}
	if err := h.svc.CancelDisband(r.Context(), actor, acct); err != nil {
		writeErr(w, r, err)
		return
	}
	apierr.WriteJSON(w, http.StatusOK, map[string]any{"disband_scheduled_at": nil})
}

func writeErr(w http.ResponseWriter, r *http.Request, err error) {
	type m struct {
		err    error
		status int
		code   string
	}
	for _, c := range []m{
		{ErrNotFound, 404, "E_NOT_FOUND"},
		{character.ErrNotFound, 404, "E_NOT_FOUND"},
		{character.ErrForbidden, 404, "E_NOT_FOUND"},
		{ErrNotInGuild, 409, "E_NOT_IN_GUILD"},
		{ErrAlreadyInGuild, 409, "E_ALREADY_IN_GUILD"},
		{ErrForbidden, 403, "E_GUILD_FORBIDDEN"},
		{ErrInvalidName, 400, "E_INVALID_NAME"},
		{ErrInvalidTag, 400, "E_INVALID_TAG"},
		{ErrNameTaken, 409, "E_NAME_TAKEN"},
		{ErrTagTaken, 409, "E_TAG_TAKEN"},
		{ErrGuildFull, 409, "E_GUILD_FULL"},
		{ErrTargetNotFound, 404, "E_TARGET_NOT_FOUND"},
		{ErrSelfTarget, 400, "E_SELF_TARGET"},
		{ErrTargetInGuild, 409, "E_TARGET_IN_GUILD"},
		{ErrAlreadyInvited, 409, "E_ALREADY_INVITED"},
		{ErrInviteInvalid, 409, "E_INVITE_INVALID"},
		{ErrNotMember, 404, "E_NOT_MEMBER"},
		{ErrLeaderCannotLeave, 409, "E_LEADER_CANNOT_LEAVE"},
		{ErrDisbandScheduled, 409, "E_DISBAND_SCHEDULED"},
		{ErrDisbandNotPending, 409, "E_DISBAND_NOT_PENDING"},
		{ErrTooLong, 400, "E_TOO_LONG"},
		{ErrIdempotencyMissing, 400, "E_IDEMPOTENCY_KEY_REQUIRED"},
		{economy.ErrInsufficientBalance, 409, "E_INSUFFICIENT_GOLD"},
	} {
		if errors.Is(err, c.err) {
			apierr.Write(w, r, c.status, c.code, err.Error())
			return
		}
	}
	apierr.Write(w, r, http.StatusInternalServerError, "E_INTERNAL", err.Error())
}
