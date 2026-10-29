package economy

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/otter/server/internal/http/apierr"
	"github.com/otter/server/internal/httpctx"
)

// Handler exposes the read (wallet) endpoint plus a dev-only mint helper.
// Mounted under /v1/character and /v1/dev by server.go.
type Handler struct {
	svc TxnService
	log *slog.Logger
	dev bool // enable /v1/dev/* endpoints when true
}

func NewHandler(svc TxnService, dev bool, logger *slog.Logger) *Handler {
	if logger == nil {
		logger = slog.Default()
	}
	return &Handler{svc: svc, dev: dev, log: logger}
}

// --- GET /v1/character/{id}/wallet ------------------------------------------

func (h *Handler) Wallet(w http.ResponseWriter, r *http.Request) {
	acct := httpctx.AccountID(r.Context())
	if acct == uuid.Nil {
		apierr.Write(w, r, http.StatusUnauthorized, "E_UNAUTHENTICATED", "auth required")
		return
	}
	idParam := chi.URLParam(r, "id")
	charID, err := uuid.Parse(idParam)
	if err != nil {
		apierr.Write(w, r, http.StatusBadRequest, "E_INVALID_ID", err.Error())
		return
	}
	// Ownership is enforced by the character middleware chain in server.go, but
	// this endpoint is not currently gated by it — hitting the DB for both is
	// cheap. TODO: attach RequireCharacter for /wallet routes too.
	wallet, err := ReadWallet(r.Context(), h.svc, charID, acct)
	if err != nil {
		apierr.Write(w, r, http.StatusInternalServerError, "E_INTERNAL", err.Error())
		return
	}
	apierr.WriteJSON(w, http.StatusOK, wallet)
}

// --- POST /v1/dev/mint (dev only) -------------------------------------------

type mintReq struct {
	CharacterID string `json:"character_id"`
	Amount      int64  `json:"amount"`
	Currency    string `json:"currency"` // default "gold"
}

// Mint is a dev-only helper for smoke testing. Returns 404 when the server is
// not running with ENV=dev — never exposed in staging/prod.
func (h *Handler) Mint(w http.ResponseWriter, r *http.Request) {
	if !h.dev {
		http.NotFound(w, r)
		return
	}
	var body mintReq
	if err := decodeJSON(r, &body); err != nil {
		apierr.Write(w, r, http.StatusBadRequest, "E_INVALID_JSON", err.Error())
		return
	}
	charID, err := uuid.Parse(body.CharacterID)
	if err != nil {
		apierr.Write(w, r, http.StatusBadRequest, "E_INVALID_ID", "character_id: "+err.Error())
		return
	}
	if body.Amount <= 0 {
		apierr.Write(w, r, http.StatusBadRequest, "E_VALIDATION", "amount must be positive")
		return
	}
	currency := body.Currency
	if currency == "" {
		currency = "gold"
	}

	// Mint = system.mint_source → character. AccountKind for the character
	// side must match the currency scope. v1 only accepts gold minting to
	// character.
	kind := AccountCharacter
	ref := charID
	if currency == "gem" || currency == "credit" { // account currencies (DD-321)
		kind = AccountAccount
		ref = httpctx.AccountID(r.Context())
	}

	// Idempotency key from the request header (idempotency middleware already
	// enforced its presence). Falls back to a synthetic key from the URL if
	// missing — dev-only.
	idemKey := r.Header.Get("Idempotency-Key")
	if idemKey == "" {
		idemKey = "dev-mint-" + charID.String() + "-" + body.Currency
	}

	intent := Intent{
		Kind:           "dev_mint",
		RefKind:        "dev",
		IdempotencyKey: idemKey,
		Entries: []Entry{
			{Currency: currency, AccountKind: AccountSystem, AccountRef: SystemMintSource, Delta: -body.Amount},
			{Currency: currency, AccountKind: kind, AccountRef: ref, Delta: body.Amount},
		},
	}
	res, err := h.svc.Write(r.Context(), intent)
	if err != nil {
		writeEconomyErr(w, r, err)
		return
	}
	apierr.WriteJSON(w, http.StatusOK, map[string]any{
		"txn_id":   res.TxnID.String(),
		"replayed": res.Replayed,
		"entries":  len(res.Rows),
	})
}

func decodeJSON(r *http.Request, dst any) error {
	body, err := io.ReadAll(io.LimitReader(r.Body, 16<<10))
	if err != nil {
		return err
	}
	if len(body) == 0 {
		return errors.New("empty body")
	}
	return json.Unmarshal(body, dst)
}

func writeEconomyErr(w http.ResponseWriter, r *http.Request, err error) {
	var ve *ValidationError
	switch {
	case errors.As(err, &ve):
		apierr.WriteField(w, r, http.StatusBadRequest, "E_VALIDATION", ve.Message, ve.Field)
	case errors.Is(err, ErrInsufficientBalance):
		apierr.Write(w, r, http.StatusConflict, "E_INSUFFICIENT_BALANCE", err.Error())
	case errors.Is(err, ErrCurrencyCap):
		apierr.Write(w, r, http.StatusConflict, "E_CURRENCY_CAP", err.Error())
	case errors.Is(err, ErrScopeMismatch):
		apierr.Write(w, r, http.StatusBadRequest, "E_CURRENCY_SCOPE", err.Error())
	case errors.Is(err, ErrUnknownCurrency):
		apierr.Write(w, r, http.StatusBadRequest, "E_UNKNOWN_CURRENCY", err.Error())
	case errors.Is(err, ErrSerialization):
		apierr.Write(w, r, http.StatusServiceUnavailable, "E_SERIALIZATION_RETRY", "please retry")
	default:
		apierr.Write(w, r, http.StatusInternalServerError, "E_INTERNAL", err.Error())
	}
}
