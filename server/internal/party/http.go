package party

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/otter/server/internal/domain/character"
	"github.com/otter/server/internal/http/apierr"
	"github.com/otter/server/internal/httpctx"
)

// Handler serves /v1/character/{id}/party/*.
type Handler struct{ svc *Service }

func NewHandler(svc *Service) *Handler { return &Handler{svc: svc} }

// scope extracts (actor character, account) or writes the error response.
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

type targetBody struct {
	TargetName        string    `json:"target_name"`
	TargetCharacterID uuid.UUID `json:"target_character_id"`
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	actor, acct, ok := scope(w, r)
	if !ok {
		return
	}
	p, err := h.svc.Get(r.Context(), actor, acct)
	if errors.Is(err, ErrNotInParty) {
		apierr.WriteJSON(w, http.StatusOK, map[string]any{"party": nil})
		return
	}
	if err != nil {
		writeErr(w, r, err)
		return
	}
	apierr.WriteJSON(w, http.StatusOK, map[string]any{"party": p})
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	actor, acct, ok := scope(w, r)
	if !ok {
		return
	}
	p, err := h.svc.Create(r.Context(), actor, acct)
	if err != nil {
		writeErr(w, r, err)
		return
	}
	apierr.WriteJSON(w, http.StatusCreated, map[string]any{"party": p})
}

func (h *Handler) Invite(w http.ResponseWriter, r *http.Request) {
	actor, acct, ok := scope(w, r)
	if !ok {
		return
	}
	var b targetBody
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
	inviteID, err := uuid.Parse(chi.URLParam(r, "invite_id"))
	if err != nil {
		apierr.Write(w, r, http.StatusBadRequest, "E_INVALID_ID", err.Error())
		return
	}
	p, err := h.svc.Accept(r.Context(), actor, acct, inviteID)
	if err != nil {
		writeErr(w, r, err)
		return
	}
	apierr.WriteJSON(w, http.StatusOK, map[string]any{"party": p})
}

func (h *Handler) Decline(w http.ResponseWriter, r *http.Request) {
	actor, acct, ok := scope(w, r)
	if !ok {
		return
	}
	inviteID, err := uuid.Parse(chi.URLParam(r, "invite_id"))
	if err != nil {
		apierr.Write(w, r, http.StatusBadRequest, "E_INVALID_ID", err.Error())
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
	var b targetBody
	if !decode(w, r, &b) {
		return
	}
	if err := h.svc.Kick(r.Context(), actor, acct, b.TargetCharacterID); err != nil {
		writeErr(w, r, err)
		return
	}
	apierr.WriteJSON(w, http.StatusOK, map[string]any{"kicked": b.TargetCharacterID})
}

func (h *Handler) Promote(w http.ResponseWriter, r *http.Request) {
	actor, acct, ok := scope(w, r)
	if !ok {
		return
	}
	var b targetBody
	if !decode(w, r, &b) {
		return
	}
	if err := h.svc.Promote(r.Context(), actor, acct, b.TargetCharacterID); err != nil {
		writeErr(w, r, err)
		return
	}
	apierr.WriteJSON(w, http.StatusOK, map[string]any{"leader_character_id": b.TargetCharacterID})
}

func (h *Handler) Disband(w http.ResponseWriter, r *http.Request) {
	actor, acct, ok := scope(w, r)
	if !ok {
		return
	}
	if err := h.svc.Disband(r.Context(), actor, acct); err != nil {
		writeErr(w, r, err)
		return
	}
	apierr.WriteJSON(w, http.StatusOK, map[string]any{"disbanded": true})
}

func writeErr(w http.ResponseWriter, r *http.Request, err error) {
	switch {
	case errors.Is(err, ErrNotFound), errors.Is(err, character.ErrNotFound), errors.Is(err, character.ErrForbidden):
		apierr.Write(w, r, http.StatusNotFound, "E_NOT_FOUND", "not found")
	case errors.Is(err, ErrNotInParty):
		apierr.Write(w, r, http.StatusConflict, "E_NOT_IN_PARTY", err.Error())
	case errors.Is(err, ErrAlreadyInParty):
		apierr.Write(w, r, http.StatusConflict, "E_ALREADY_IN_PARTY", err.Error())
	case errors.Is(err, ErrNotLeader):
		apierr.Write(w, r, http.StatusForbidden, "E_NOT_LEADER", err.Error())
	case errors.Is(err, ErrPartyFull):
		apierr.Write(w, r, http.StatusConflict, "E_PARTY_FULL", err.Error())
	case errors.Is(err, ErrTargetNotFound):
		apierr.Write(w, r, http.StatusNotFound, "E_TARGET_NOT_FOUND", err.Error())
	case errors.Is(err, ErrSelfTarget):
		apierr.Write(w, r, http.StatusBadRequest, "E_SELF_TARGET", err.Error())
	case errors.Is(err, ErrTargetInParty):
		apierr.Write(w, r, http.StatusConflict, "E_TARGET_IN_PARTY", err.Error())
	case errors.Is(err, ErrAlreadyInvited):
		apierr.Write(w, r, http.StatusConflict, "E_ALREADY_INVITED", err.Error())
	case errors.Is(err, ErrInviteInvalid):
		apierr.Write(w, r, http.StatusConflict, "E_INVITE_INVALID", err.Error())
	case errors.Is(err, ErrNotMember):
		apierr.Write(w, r, http.StatusNotFound, "E_NOT_MEMBER", err.Error())
	default:
		apierr.Write(w, r, http.StatusInternalServerError, "E_INTERNAL", err.Error())
	}
}
