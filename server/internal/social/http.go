package social

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

// Handler serves /v1/character/{id}/{friends,blocks,mail,chat}.
type Handler struct{ svc *Service }

func NewHandler(svc *Service) *Handler { return &Handler{svc: svc} }

type body struct {
	TargetName        string      `json:"target_name"`
	TargetCharacterID uuid.UUID   `json:"target_character_id"`
	RecipientName     string      `json:"recipient_name"`
	Subject           string      `json:"subject"`
	Body              string      `json:"body"`
	ItemIDs           []uuid.UUID `json:"item_ids"`
	GoldCopper        int64       `json:"gold_copper"`
	Channel           string      `json:"channel"`
}

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

// withBody wraps handlers that need (actor, account, decoded body).
func withBody(fn func(w http.ResponseWriter, r *http.Request, actor, acct uuid.UUID, b body)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		actor, acct, ok := scope(w, r)
		if !ok {
			return
		}
		var b body
		if r.Method == http.MethodPost && r.ContentLength != 0 {
			if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
				apierr.Write(w, r, http.StatusBadRequest, "E_INVALID_JSON", err.Error())
				return
			}
		}
		fn(w, r, actor, acct, b)
	}
}

func pathID(w http.ResponseWriter, r *http.Request, key string) (uuid.UUID, bool) {
	id, err := uuid.Parse(chi.URLParam(r, key))
	if err != nil {
		apierr.Write(w, r, http.StatusBadRequest, "E_INVALID_ID", key+": "+err.Error())
	}
	return id, err == nil
}

func reply(w http.ResponseWriter, r *http.Request, status int, v any, err error) {
	if err != nil {
		writeErr(w, r, err)
		return
	}
	apierr.WriteJSON(w, status, v)
}

// --- friends ----------------------------------------------------------------

func (h *Handler) Friends() http.HandlerFunc {
	return withBody(func(w http.ResponseWriter, r *http.Request, actor, acct uuid.UUID, _ body) {
		list, err := h.svc.Friends(r.Context(), actor, acct)
		reply(w, r, 200, map[string]any{"count": len(list), "friends": list}, err)
	})
}

func (h *Handler) FriendRequests() http.HandlerFunc {
	return withBody(func(w http.ResponseWriter, r *http.Request, actor, acct uuid.UUID, _ body) {
		list, err := h.svc.FriendRequests(r.Context(), actor, acct)
		reply(w, r, 200, map[string]any{"count": len(list), "requests": list}, err)
	})
}

func (h *Handler) RequestFriend() http.HandlerFunc {
	return withBody(func(w http.ResponseWriter, r *http.Request, actor, acct uuid.UUID, b body) {
		accepted, err := h.svc.RequestFriend(r.Context(), actor, acct, b.TargetName)
		reply(w, r, 200, map[string]any{"friends": accepted, "requested": !accepted}, err)
	})
}

func (h *Handler) AnswerFriend(accept bool) http.HandlerFunc {
	return withBody(func(w http.ResponseWriter, r *http.Request, actor, acct uuid.UUID, _ body) {
		id, ok := pathID(w, r, "request_id")
		if !ok {
			return
		}
		var err error
		if accept {
			err = h.svc.AcceptFriend(r.Context(), actor, acct, id)
		} else {
			err = h.svc.DeclineFriend(r.Context(), actor, acct, id)
		}
		reply(w, r, 200, map[string]any{"accepted": accept}, err)
	})
}

func (h *Handler) RemoveFriend() http.HandlerFunc {
	return withBody(func(w http.ResponseWriter, r *http.Request, actor, acct uuid.UUID, b body) {
		err := h.svc.RemoveFriend(r.Context(), actor, acct, b.TargetCharacterID)
		reply(w, r, 200, map[string]any{"removed": b.TargetCharacterID}, err)
	})
}

// --- blocks -----------------------------------------------------------------

func (h *Handler) Blocks() http.HandlerFunc {
	return withBody(func(w http.ResponseWriter, r *http.Request, actor, acct uuid.UUID, _ body) {
		list, err := h.svc.BlockList(r.Context(), actor, acct)
		reply(w, r, 200, map[string]any{"count": len(list), "blocked": list}, err)
	})
}

func (h *Handler) Block() http.HandlerFunc {
	return withBody(func(w http.ResponseWriter, r *http.Request, actor, acct uuid.UUID, b body) {
		id, err := h.svc.Block(r.Context(), actor, acct, b.TargetName)
		reply(w, r, 200, map[string]any{"blocked": id}, err)
	})
}

func (h *Handler) Unblock() http.HandlerFunc {
	return withBody(func(w http.ResponseWriter, r *http.Request, actor, acct uuid.UUID, b body) {
		err := h.svc.Unblock(r.Context(), actor, acct, b.TargetCharacterID)
		reply(w, r, 200, map[string]any{"unblocked": b.TargetCharacterID}, err)
	})
}

// --- mail -------------------------------------------------------------------

func (h *Handler) Inbox() http.HandlerFunc {
	return withBody(func(w http.ResponseWriter, r *http.Request, actor, acct uuid.UUID, _ body) {
		list, err := h.svc.Inbox(r.Context(), actor, acct)
		reply(w, r, 200, map[string]any{"count": len(list), "mail": list}, err)
	})
}

func (h *Handler) ReadMail() http.HandlerFunc {
	return withBody(func(w http.ResponseWriter, r *http.Request, actor, acct uuid.UUID, _ body) {
		id, ok := pathID(w, r, "mail_id")
		if !ok {
			return
		}
		m, err := h.svc.ReadMail(r.Context(), actor, acct, id)
		reply(w, r, 200, m, err)
	})
}

func (h *Handler) SendMail() http.HandlerFunc {
	return withBody(func(w http.ResponseWriter, r *http.Request, actor, acct uuid.UUID, b body) {
		res, err := h.svc.SendMail(r.Context(), actor, acct, SendRequest{
			RecipientName: b.RecipientName, Subject: b.Subject, Body: b.Body,
			ItemIDs: b.ItemIDs, GoldCopper: b.GoldCopper,
			IdempotencyKey: r.Header.Get("Idempotency-Key"),
		})
		reply(w, r, http.StatusCreated, res, err)
	})
}

func (h *Handler) ClaimMail() http.HandlerFunc {
	return withBody(func(w http.ResponseWriter, r *http.Request, actor, acct uuid.UUID, _ body) {
		id, ok := pathID(w, r, "mail_id")
		if !ok {
			return
		}
		m, err := h.svc.ClaimMail(r.Context(), actor, acct, id)
		reply(w, r, 200, m, err)
	})
}

func (h *Handler) DeleteMail() http.HandlerFunc {
	return withBody(func(w http.ResponseWriter, r *http.Request, actor, acct uuid.UUID, _ body) {
		id, ok := pathID(w, r, "mail_id")
		if !ok {
			return
		}
		err := h.svc.DeleteMail(r.Context(), actor, acct, id)
		reply(w, r, 200, map[string]any{"deleted": id}, err)
	})
}

// --- chat -------------------------------------------------------------------

func (h *Handler) SendChat() http.HandlerFunc {
	return withBody(func(w http.ResponseWriter, r *http.Request, actor, acct uuid.UUID, b body) {
		msg, err := h.svc.SendChat(r.Context(), actor, acct, b.Channel, b.TargetName, b.Body)
		reply(w, r, http.StatusCreated, msg, err)
	})
}

func (h *Handler) ChatHistory() http.HandlerFunc {
	return withBody(func(w http.ResponseWriter, r *http.Request, actor, acct uuid.UUID, _ body) {
		msgs, err := h.svc.ChatHistory(r.Context(), actor, acct, chi.URLParam(r, "channel"), r.URL.Query().Get("target"))
		reply(w, r, 200, map[string]any{"count": len(msgs), "messages": msgs}, err)
	})
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
		{ErrTargetNotFound, 404, "E_TARGET_NOT_FOUND"},
		{ErrSelfTarget, 400, "E_SELF_TARGET"},
		{ErrBlocked, 403, "E_BLOCKED"},
		{ErrAlreadyFriends, 409, "E_ALREADY_FRIENDS"},
		{ErrAlreadyRequested, 409, "E_ALREADY_REQUESTED"},
		{ErrFriendCap, 409, "E_FRIEND_CAP"},
		{ErrNotFriends, 404, "E_NOT_FRIENDS"},
		{ErrInboxFull, 409, "E_INBOX_FULL"},
		{ErrTooLong, 400, "E_TOO_LONG"},
		{ErrEmpty, 400, "E_EMPTY"},
		{ErrTooManyAttachments, 400, "E_TOO_MANY_ATTACHMENTS"},
		{ErrAttachmentInvalid, 400, "E_ATTACHMENT_INVALID"},
		{ErrHasAttachments, 409, "E_HAS_ATTACHMENTS"},
		{ErrNothingToClaim, 409, "E_NOTHING_TO_CLAIM"},
		{ErrBagFull, 409, "E_BAG_FULL"},
		{ErrBadChannel, 400, "E_BAD_CHANNEL"},
		{ErrNotInChannel, 409, "E_NOT_IN_CHANNEL"},
		{ErrIdempotencyMissing, 400, "E_IDEMPOTENCY_KEY_REQUIRED"},
		{ErrMuted, 403, "E_MUTED"},
		{economy.ErrInsufficientBalance, 409, "E_INSUFFICIENT_GOLD"},
	} {
		if errors.Is(err, c.err) {
			apierr.Write(w, r, c.status, c.code, err.Error())
			return
		}
	}
	apierr.Write(w, r, http.StatusInternalServerError, "E_INTERNAL", err.Error())
}
