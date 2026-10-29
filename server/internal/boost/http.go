package boost

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/otter/server/internal/http/apierr"
)

type Handler struct {
	svc  *Service
	keys KeySource
	// devAds lets a development build "watch" an ad without the network.
	devAds bool
}

func NewHandler(svc *Service, keys KeySource, devAds bool) *Handler {
	return &Handler{svc: svc, keys: keys, devAds: devAds}
}

// Routes are mounted under /v1/character/{id} (ownership: router guard).
func (h *Handler) Routes(r chi.Router, idem func(http.Handler) http.Handler) {
	r.Get("/boosts", h.list)
	r.With(idem).Post("/boosts/{category}/dev-ad", h.devAd)
}

func (h *Handler) list(w http.ResponseWriter, r *http.Request) {
	char, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		apierr.Write(w, r, http.StatusBadRequest, "E_INVALID_ID", "id")
		return
	}
	list, err := h.svc.List(r.Context(), char)
	if err != nil {
		writeErr(w, r, err)
		return
	}
	apierr.WriteJSON(w, http.StatusOK, map[string]any{"boosts": list, "ads_enabled": h.devAds,
		"duration_seconds": int(Duration.Seconds()), "bonus_pct": BonusPct})
}

// devAd stands in for a watched ad in development builds only.
func (h *Handler) devAd(w http.ResponseWriter, r *http.Request) {
	if !h.devAds {
		writeErr(w, r, ErrAdsOff)
		return
	}
	char, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		apierr.Write(w, r, http.StatusBadRequest, "E_INVALID_ID", "id")
		return
	}
	st, err := h.svc.Grant(r.Context(), char, chi.URLParam(r, "category"), SourceAd, "dev:"+uuid.NewString())
	if err != nil {
		writeErr(w, r, err)
		return
	}
	apierr.WriteJSON(w, http.StatusOK, st)
}

// SSV is GET /v1/ads/ssv: AdMob's rewarded-ad callback (public; trusted
// only through its signature).
func (h *Handler) SSV(w http.ResponseWriter, r *http.Request) {
	rw, err := VerifySSV(r.Context(), r.URL.RawQuery, h.keys)
	if err != nil {
		apierr.Write(w, r, http.StatusBadRequest, "E_BAD_REWARD", err.Error())
		return
	}
	char, err := uuid.Parse(rw.CharacterID)
	if err != nil {
		apierr.Write(w, r, http.StatusBadRequest, "E_BAD_REWARD", "user_id")
		return
	}
	if _, err := h.svc.Grant(r.Context(), char, rw.Category, SourceAd, "ad:"+rw.TransactionID); err != nil &&
		!errors.Is(err, ErrActive) && !errors.Is(err, ErrAdLimit) {
		writeErr(w, r, err)
		return
	}
	apierr.WriteJSON(w, http.StatusOK, map[string]bool{"ok": true})
}

func writeErr(w http.ResponseWriter, r *http.Request, err error) {
	switch {
	case errors.Is(err, ErrCategory), errors.Is(err, ErrNotFound):
		apierr.Write(w, r, http.StatusNotFound, "E_NOT_FOUND", err.Error())
	case errors.Is(err, ErrActive):
		apierr.Write(w, r, http.StatusConflict, "E_BOOST_ACTIVE", err.Error())
	case errors.Is(err, ErrAdLimit):
		apierr.Write(w, r, http.StatusConflict, "E_AD_LIMIT", err.Error())
	case errors.Is(err, ErrAdsOff):
		apierr.Write(w, r, http.StatusConflict, "E_ADS_UNAVAILABLE", err.Error())
	case errors.Is(err, ErrBadReward):
		apierr.Write(w, r, http.StatusBadRequest, "E_BAD_REWARD", err.Error())
	default:
		apierr.Write(w, r, http.StatusInternalServerError, "E_INTERNAL", err.Error())
	}
}
