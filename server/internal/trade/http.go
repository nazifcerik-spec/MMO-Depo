package trade

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/otter/server/internal/http/apierr"
	"github.com/otter/server/internal/httpctx"
)

type Handler struct{ svc *Service }

func NewHandler(svc *Service) *Handler { return &Handler{svc: svc} }

// Routes are mounted under /v1/character/{id}/trade (ownership: router guard).
func (h *Handler) Routes(r chi.Router, idem func(http.Handler) http.Handler) {
	r.Get("/trade", h.current)
	r.With(idem).Post("/trade", h.open)
	r.With(idem).Post("/trade/{trade_id}/offer", h.offer)
	r.With(idem).Post("/trade/{trade_id}/remove", h.remove)
	r.With(idem).Post("/trade/{trade_id}/ready", h.ready)
	r.With(idem).Post("/trade/{trade_id}/cancel", h.cancel)
}

func ids(w http.ResponseWriter, r *http.Request, withTrade bool) (actor, account, trade uuid.UUID, ok bool) {
	var err error
	if actor, err = uuid.Parse(chi.URLParam(r, "id")); err != nil {
		apierr.Write(w, r, http.StatusBadRequest, "E_INVALID_ID", "id")
		return
	}
	if withTrade {
		if trade, err = uuid.Parse(chi.URLParam(r, "trade_id")); err != nil {
			apierr.Write(w, r, http.StatusBadRequest, "E_INVALID_ID", "trade_id")
			return
		}
	}
	return actor, httpctx.AccountID(r.Context()), trade, true
}

func reply(w http.ResponseWriter, r *http.Request, v *View, err error) {
	if err != nil {
		writeErr(w, r, err)
		return
	}
	apierr.WriteJSON(w, http.StatusOK, map[string]any{"trade": v})
}

func (h *Handler) current(w http.ResponseWriter, r *http.Request) {
	if actor, acct, _, ok := ids(w, r, false); ok {
		v, err := h.svc.Current(r.Context(), actor, acct)
		reply(w, r, v, err)
	}
}

func (h *Handler) open(w http.ResponseWriter, r *http.Request) {
	actor, acct, _, ok := ids(w, r, false)
	if !ok {
		return
	}
	var body struct {
		TargetName string `json:"target_name"`
	}
	_ = json.NewDecoder(r.Body).Decode(&body)
	v, err := h.svc.Open(r.Context(), actor, acct, body.TargetName)
	reply(w, r, v, err)
}

type itemBody struct {
	ItemID   uuid.UUID `json:"item_id"`
	Quantity int32     `json:"quantity"`
}

func (h *Handler) offer(w http.ResponseWriter, r *http.Request) {
	actor, acct, id, ok := ids(w, r, true)
	if !ok {
		return
	}
	var body itemBody
	_ = json.NewDecoder(r.Body).Decode(&body)
	v, err := h.svc.Offer(r.Context(), actor, acct, id, body.ItemID, body.Quantity)
	reply(w, r, v, err)
}

func (h *Handler) remove(w http.ResponseWriter, r *http.Request) {
	actor, acct, id, ok := ids(w, r, true)
	if !ok {
		return
	}
	var body itemBody
	_ = json.NewDecoder(r.Body).Decode(&body)
	v, err := h.svc.Remove(r.Context(), actor, acct, id, body.ItemID)
	reply(w, r, v, err)
}

func (h *Handler) ready(w http.ResponseWriter, r *http.Request) {
	if actor, acct, id, ok := ids(w, r, true); ok {
		v, err := h.svc.Ready(r.Context(), actor, acct, id)
		reply(w, r, v, err)
	}
}

func (h *Handler) cancel(w http.ResponseWriter, r *http.Request) {
	if actor, acct, id, ok := ids(w, r, true); ok {
		reply(w, r, nil, h.svc.Cancel(r.Context(), actor, acct, id))
	}
}

func writeErr(w http.ResponseWriter, r *http.Request, err error) {
	for _, m := range []struct {
		err    error
		status int
		code   string
	}{
		{ErrNotFound, 404, "E_NOT_FOUND"},
		{ErrSelf, 400, "E_SELF_TARGET"},
		{ErrBusy, 409, "E_TRADE_BUSY"},
		{ErrClosed, 409, "E_TRADE_CLOSED"},
		{ErrItem, 409, "E_TRADE_ITEM"},
		{ErrTooMany, 409, "E_TRADE_TOO_MANY"},
		{ErrBagFull, 409, "E_BAG_FULL"},
		{ErrNoGold, 409, "E_INSUFFICIENT_GOLD"},
		{ErrNotReadyToGo, 409, "E_TRADE_EMPTY"},
		{errChangedMidway, 409, "E_TRADE_ITEM"},
	} {
		if errors.Is(err, m.err) {
			apierr.Write(w, r, m.status, m.code, err.Error())
			return
		}
	}
	apierr.Write(w, r, http.StatusInternalServerError, "E_INTERNAL", err.Error())
}
