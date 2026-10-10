package clienterr

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/otter/server/internal/http/apierr"
)

// Handler serves POST /v1/client/errors (public: crashes happen before login).
type Handler struct{ svc *Service }

func NewHandler(svc *Service) *Handler { return &Handler{svc: svc} }

func (h *Handler) Report(w http.ResponseWriter, r *http.Request) {
	var req Report
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 16<<10)).Decode(&req); err != nil {
		apierr.Write(w, r, http.StatusBadRequest, "E_INVALID_JSON", err.Error())
		return
	}
	if err := h.svc.Record(r.Context(), req); err != nil {
		WriteErr(w, r, err)
		return
	}
	w.WriteHeader(http.StatusAccepted)
}

// WriteErr maps clienterr errors (shared with the admin handler).
func WriteErr(w http.ResponseWriter, r *http.Request, err error) {
	if errors.Is(err, ErrInvalid) {
		apierr.Write(w, r, http.StatusBadRequest, "E_VALIDATION", err.Error())
		return
	}
	apierr.Write(w, r, http.StatusInternalServerError, "E_INTERNAL", err.Error())
}
