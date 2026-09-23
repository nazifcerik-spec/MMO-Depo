package market

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/otter/server/internal/domain/character"
	"github.com/otter/server/internal/economy"
	"github.com/otter/server/internal/http/apierr"
	"github.com/otter/server/internal/httpctx"
)

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

func listingID(w http.ResponseWriter, r *http.Request) (uuid.UUID, bool) {
	id, err := uuid.Parse(chi.URLParam(r, "listing_id"))
	if err != nil {
		apierr.Write(w, r, http.StatusBadRequest, "E_INVALID_ID", err.Error())
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

type sellBody struct {
	ItemInstanceID  uuid.UUID `json:"item_instance_id"`
	StackQuantity   int32     `json:"stack_quantity"`
	UnitPriceCopper int64     `json:"unit_price_copper"`
	DurationHours   int16     `json:"duration_hours"`
}

func (h *Handler) Sell(w http.ResponseWriter, r *http.Request) {
	actor, acct, ok := scope(w, r)
	if !ok {
		return
	}
	var b sellBody
	if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
		apierr.Write(w, r, http.StatusBadRequest, "E_INVALID_JSON", err.Error())
		return
	}
	if b.StackQuantity == 0 {
		b.StackQuantity = 1
	}
	l, err := h.svc.Sell(r.Context(), actor, acct, SellRequest{
		ItemID: b.ItemInstanceID, Quantity: b.StackQuantity, UnitPrice: b.UnitPriceCopper,
		DurationHours: b.DurationHours, IdempotencyKey: r.Header.Get("Idempotency-Key"),
	})
	reply(w, r, http.StatusCreated, l, err)
}

func (h *Handler) Buy(w http.ResponseWriter, r *http.Request) {
	actor, acct, ok := scope(w, r)
	if !ok {
		return
	}
	id, ok := listingID(w, r)
	if !ok {
		return
	}
	var b struct {
		ExpectedUnitPrice *int64 `json:"expected_unit_price_copper"`
	}
	if r.ContentLength != 0 {
		if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
			apierr.Write(w, r, http.StatusBadRequest, "E_INVALID_JSON", err.Error())
			return
		}
	}
	l, err := h.svc.Buy(r.Context(), actor, acct, id, b.ExpectedUnitPrice)
	reply(w, r, 200, l, err)
}

func (h *Handler) Cancel(w http.ResponseWriter, r *http.Request) {
	actor, acct, ok := scope(w, r)
	if !ok {
		return
	}
	id, ok := listingID(w, r)
	if !ok {
		return
	}
	l, err := h.svc.Cancel(r.Context(), actor, acct, id)
	reply(w, r, 200, l, err)
}

func (h *Handler) Mine(w http.ResponseWriter, r *http.Request) {
	actor, acct, ok := scope(w, r)
	if !ok {
		return
	}
	list, err := h.svc.MyListings(r.Context(), actor, acct)
	reply(w, r, 200, map[string]any{"count": len(list), "listings": list}, err)
}

func (h *Handler) Search(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	num := func(k string) int64 { n, _ := strconv.ParseInt(q.Get(k), 10, 64); return n }
	list, err := h.svc.Search(r.Context(), SearchQuery{
		TemplateCode: q.Get("template"), ItemClass: q.Get("class"), Rarity: q.Get("rarity"),
		MinPrice: num("min_price"), MaxPrice: num("max_price"), Sort: q.Get("sort"),
		Limit: int(num("limit")), Offset: int(num("offset")),
	})
	reply(w, r, 200, map[string]any{"count": len(list), "listings": list}, err)
}

func (h *Handler) History(w http.ResponseWriter, r *http.Request) {
	st, err := h.svc.History(r.Context(), chi.URLParam(r, "template_code"))
	reply(w, r, 200, st, err)
}

func writeErr(w http.ResponseWriter, r *http.Request, err error) {
	type m struct {
		err    error
		status int
		code   string
	}
	for _, c := range []m{
		{ErrNotFound, 404, "E_LISTING_NOT_FOUND"},
		{character.ErrNotFound, 404, "E_NOT_FOUND"},
		{character.ErrForbidden, 404, "E_NOT_FOUND"},
		{ErrUnavailable, 409, "E_LISTING_UNAVAILABLE"},
		{ErrItemNotFound, 404, "E_ITEM_NOT_FOUND"},
		{ErrNotTradable, 409, "E_ITEM_NOT_TRADABLE"},
		{ErrBadQuantity, 400, "E_BAD_QUANTITY"},
		{ErrBadDuration, 400, "E_BAD_DURATION"},
		{ErrPriceOutOfRange, 400, "E_PRICE_OUT_OF_RANGE"},
		{ErrPriceMismatch, 409, "E_PRICE_MISMATCH"},
		{ErrListingCap, 409, "E_MARKET_LISTING_CAP"},
		{ErrSelfBuy, 409, "E_SELF_BUY"},
		{ErrNotSeller, 403, "E_NOT_SELLER"},
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
