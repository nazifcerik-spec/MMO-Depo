package character

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/otter/server/internal/http/apierr"
	"github.com/otter/server/internal/httpctx"
)

// Handler bundles the /v1/character endpoints. AccountID is read from ctx
// (set by RequireAuth); no character endpoint is reachable anonymously.
type Handler struct {
	svc *Service
	log *slog.Logger
}

func NewHandler(svc *Service, logger *slog.Logger) *Handler {
	if logger == nil {
		logger = slog.Default()
	}
	return &Handler{svc: svc, log: logger}
}

// --- Create -----------------------------------------------------------------

type createReq struct {
	Name      string `json:"name"`
	RaceCode  string `json:"race_code"`
	ClassCode string `json:"class_code"`
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	acct := httpctx.AccountID(r.Context())
	if acct == uuid.Nil {
		apierr.Write(w, r, http.StatusUnauthorized, "E_UNAUTHENTICATED", "auth required")
		return
	}
	var body createReq
	if err := decodeJSON(r, &body); err != nil {
		apierr.Write(w, r, http.StatusBadRequest, "E_INVALID_JSON", err.Error())
		return
	}
	c, err := h.svc.Create(r.Context(), CreateRequest{
		AccountID: acct,
		Name:      body.Name,
		RaceCode:  body.RaceCode,
		ClassCode: body.ClassCode,
	})
	if err != nil {
		writeCharErr(w, r, err)
		return
	}
	apierr.WriteJSON(w, http.StatusCreated, characterResponse(c))
}

// --- List -------------------------------------------------------------------

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	acct := httpctx.AccountID(r.Context())
	if acct == uuid.Nil {
		apierr.Write(w, r, http.StatusUnauthorized, "E_UNAUTHENTICATED", "auth required")
		return
	}
	list, err := h.svc.ListForAccount(r.Context(), acct)
	if err != nil {
		writeCharErr(w, r, err)
		return
	}
	out := make([]characterResp, 0, len(list))
	for _, c := range list {
		out = append(out, characterResponse(&c))
	}
	apierr.WriteJSON(w, http.StatusOK, map[string]any{
		"characters": out,
		"slot_cap":   h.svc.cfg.SlotCap,
	})
}

// --- Sheet (GET /v1/character/{id}/sheet) -----------------------------------

func (h *Handler) Sheet(w http.ResponseWriter, r *http.Request) {
	acct := httpctx.AccountID(r.Context())
	id, err := parseIDParam(r, "id")
	if err != nil {
		apierr.Write(w, r, http.StatusBadRequest, "E_INVALID_ID", err.Error())
		return
	}
	sheet, err := h.svc.Sheet(r.Context(), id, acct)
	if err != nil {
		writeCharErr(w, r, err)
		return
	}
	apierr.WriteJSON(w, http.StatusOK, sheet)
}

// --- Enter (POST /v1/character/{id}/enter) ----------------------------------

func (h *Handler) Enter(w http.ResponseWriter, r *http.Request) {
	acct := httpctx.AccountID(r.Context())
	id, err := parseIDParam(r, "id")
	if err != nil {
		apierr.Write(w, r, http.StatusBadRequest, "E_INVALID_ID", err.Error())
		return
	}
	sheet, err := h.svc.Enter(r.Context(), id, acct)
	if err != nil {
		writeCharErr(w, r, err)
		return
	}
	// After Enter, the client is expected to send `X-Character-Id: <id>` on
	// subsequent requests. The response echoes the id explicitly for client
	// convenience.
	apierr.WriteJSON(w, http.StatusOK, map[string]any{
		"character_id": id.String(),
		"sheet":        sheet,
	})
}

// --- Delete-request / Delete-cancel -----------------------------------------

func (h *Handler) RequestDelete(w http.ResponseWriter, r *http.Request) {
	acct := httpctx.AccountID(r.Context())
	id, err := parseIDParam(r, "id")
	if err != nil {
		apierr.Write(w, r, http.StatusBadRequest, "E_INVALID_ID", err.Error())
		return
	}
	c, err := h.svc.RequestDelete(r.Context(), id, acct)
	if err != nil {
		writeCharErr(w, r, err)
		return
	}
	apierr.WriteJSON(w, http.StatusOK, characterResponse(c))
}

func (h *Handler) CancelDelete(w http.ResponseWriter, r *http.Request) {
	acct := httpctx.AccountID(r.Context())
	id, err := parseIDParam(r, "id")
	if err != nil {
		apierr.Write(w, r, http.StatusBadRequest, "E_INVALID_ID", err.Error())
		return
	}
	c, err := h.svc.CancelDelete(r.Context(), id, acct)
	if err != nil {
		writeCharErr(w, r, err)
		return
	}
	apierr.WriteJSON(w, http.StatusOK, characterResponse(c))
}

// --- Rename -----------------------------------------------------------------

type renameReq struct {
	Name string `json:"name"`
}

func (h *Handler) Rename(w http.ResponseWriter, r *http.Request) {
	acct := httpctx.AccountID(r.Context())
	id, err := parseIDParam(r, "id")
	if err != nil {
		apierr.Write(w, r, http.StatusBadRequest, "E_INVALID_ID", err.Error())
		return
	}
	var body renameReq
	if err := decodeJSON(r, &body); err != nil {
		apierr.Write(w, r, http.StatusBadRequest, "E_INVALID_JSON", err.Error())
		return
	}
	c, err := h.svc.Rename(r.Context(), id, acct, body.Name)
	if err != nil {
		writeCharErr(w, r, err)
		return
	}
	apierr.WriteJSON(w, http.StatusOK, characterResponse(c))
}

// --- response shape ---------------------------------------------------------

type characterResp struct {
	ID            string     `json:"id"`
	Name          string     `json:"name"`
	RaceCode      string     `json:"race_code"`
	ClassCode     string     `json:"class_code"`
	SpecCode      *string    `json:"spec_code,omitempty"`
	Level         int        `json:"level"`
	XP            int64      `json:"xp"`
	State         string     `json:"state"`
	CreatedAt     string     `json:"created_at"`
	LastPlayedAt  *time.Time `json:"last_played_at,omitempty"`
	DeleteGraceAt *time.Time `json:"delete_grace_at,omitempty"`
}

func characterResponse(c *Character) characterResp {
	return characterResp{
		ID:            c.ID.String(),
		Name:          c.Name,
		RaceCode:      c.RaceCode,
		ClassCode:     c.ClassCode,
		SpecCode:      c.SpecCode,
		Level:         int(c.Level),
		XP:            c.XP,
		State:         string(c.State),
		CreatedAt:     c.CreatedAt.UTC().Format(time.RFC3339),
		LastPlayedAt:  c.LastPlayedAt,
		DeleteGraceAt: c.DeleteGraceAt,
	}
}

// --- helpers ----------------------------------------------------------------

func parseIDParam(r *http.Request, name string) (uuid.UUID, error) {
	raw := chi.URLParam(r, name)
	if raw == "" {
		return uuid.Nil, errors.New("missing id")
	}
	return uuid.Parse(raw)
}

func decodeJSON(r *http.Request, dst any) error {
	body, err := io.ReadAll(io.LimitReader(r.Body, 16<<10)) // 16 KB character payloads
	if err != nil {
		return err
	}
	if len(body) == 0 {
		return errors.New("empty body")
	}
	return json.Unmarshal(body, dst)
}

func writeCharErr(w http.ResponseWriter, r *http.Request, err error) {
	var ve *ValidationError
	switch {
	case errors.As(err, &ve):
		apierr.WriteField(w, r, http.StatusBadRequest, "E_VALIDATION", ve.Message, ve.Field)
	case errors.Is(err, ErrNotFound):
		// Do not distinguish from Forbidden — avoid leaking existence.
		apierr.Write(w, r, http.StatusNotFound, "E_NOT_FOUND", "character not found")
	case errors.Is(err, ErrForbidden):
		apierr.Write(w, r, http.StatusNotFound, "E_NOT_FOUND", "character not found")
	case errors.Is(err, ErrNameTaken):
		apierr.Write(w, r, http.StatusConflict, "E_NAME_TAKEN", "name already taken")
	case errors.Is(err, ErrSlotCap):
		apierr.Write(w, r, http.StatusConflict, "E_SLOT_CAP", "character slot cap reached")
	case errors.Is(err, ErrEmailNotVerified):
		apierr.Write(w, r, http.StatusForbidden, "E_EMAIL_NOT_VERIFIED", "verify your e-mail first")
	case errors.Is(err, ErrNotPlayable):
		apierr.Write(w, r, http.StatusConflict, "E_NOT_PLAYABLE", "character is not in an active state")
	case errors.Is(err, ErrNotInDeleteGrace):
		apierr.Write(w, r, http.StatusConflict, "E_NOT_IN_DELETE_GRACE", "character is not in delete_grace state")
	default:
		apierr.Write(w, r, http.StatusInternalServerError, "E_INTERNAL", "internal error")
	}
}
