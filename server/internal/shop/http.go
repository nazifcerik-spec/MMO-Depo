package shop

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/otter/server/internal/boost"
	"github.com/otter/server/internal/http/apierr"
	"github.com/otter/server/internal/httpctx"
)

type Handler struct{ svc *Service }

func NewHandler(svc *Service) *Handler { return &Handler{svc: svc} }

// Routes are mounted under /v1/character/{id} (ownership: router guard);
// credits and cosmetics are the account's.
func (h *Handler) Routes(r chi.Router, idem func(http.Handler) http.Handler) {
	r.Get("/shop", h.shop)
	r.With(idem).Post("/shop/{code}/buy", h.buy)
	r.Get("/cosmetics", h.cosmetics)
	r.With(idem).Post("/cosmetics/select", h.selectCosmetic)
}

func (h *Handler) shop(w http.ResponseWriter, r *http.Request) {
	v, err := h.svc.Shop(r.Context(), httpctx.AccountID(r.Context()))
	if err != nil {
		writeErr(w, r, err)
		return
	}
	apierr.WriteJSON(w, http.StatusOK, v)
}

func (h *Handler) buy(w http.ResponseWriter, r *http.Request) {
	char, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		apierr.Write(w, r, http.StatusBadRequest, "E_INVALID_ID", "id")
		return
	}
	key := r.Header.Get("Idempotency-Key")
	if key == "" {
		key = uuid.NewString()
	}
	if err := h.svc.Buy(r.Context(), httpctx.AccountID(r.Context()), char, chi.URLParam(r, "code"), key); err != nil {
		writeErr(w, r, err)
		return
	}
	h.shop(w, r)
}

func (h *Handler) cosmetics(w http.ResponseWriter, r *http.Request) {
	c, err := h.svc.Cosmetics(r.Context(), httpctx.AccountID(r.Context()))
	if err != nil {
		writeErr(w, r, err)
		return
	}
	apierr.WriteJSON(w, http.StatusOK, c)
}

func (h *Handler) selectCosmetic(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Code string `json:"code"`
	}
	_ = json.NewDecoder(r.Body).Decode(&body)
	c, err := h.svc.Select(r.Context(), httpctx.AccountID(r.Context()), body.Code)
	if err != nil {
		writeErr(w, r, err)
		return
	}
	apierr.WriteJSON(w, http.StatusOK, c)
}

func writeErr(w http.ResponseWriter, r *http.Request, err error) {
	switch {
	case errors.Is(err, ErrUnknown), errors.Is(err, boost.ErrNotFound), errors.Is(err, boost.ErrCategory):
		apierr.Write(w, r, http.StatusNotFound, "E_NOT_FOUND", err.Error())
	case errors.Is(err, ErrOwned):
		apierr.Write(w, r, http.StatusConflict, "E_OWNED", err.Error())
	case errors.Is(err, ErrNoCredit):
		apierr.Write(w, r, http.StatusConflict, "E_INSUFFICIENT_CREDITS", err.Error())
	case errors.Is(err, ErrLocked):
		apierr.Write(w, r, http.StatusConflict, "E_COSMETIC_LOCKED", err.Error())
	case errors.Is(err, boost.ErrActive):
		apierr.Write(w, r, http.StatusConflict, "E_BOOST_ACTIVE", err.Error())
	default:
		apierr.Write(w, r, http.StatusInternalServerError, "E_INTERNAL", err.Error())
	}
}
