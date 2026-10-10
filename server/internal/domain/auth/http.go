package auth

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"

	"github.com/otter/server/internal/http/apierr"
	"github.com/otter/server/internal/httpctx"
)

// Handler bundles the /v1/auth/* endpoints. IPHasher is optional — when set,
// the client IP is SHA-256'd for audit persistence (no plaintext IP in DB).
type Handler struct {
	svc    *Service
	log    *slog.Logger
	ipHash func(string) []byte
}

func NewHandler(svc *Service, logger *slog.Logger) *Handler {
	if logger == nil {
		logger = slog.Default()
	}
	return &Handler{
		svc: svc,
		log: logger,
		ipHash: func(ip string) []byte {
			s := sha256.Sum256([]byte(ip))
			return s[:]
		},
	}
}

// Mount attaches the handler to the given chi.Router under /auth (usually
// mounted from /v1). Public endpoints (register/login/refresh) are gated only
// by rate limit + idempotency middleware; logout requires auth.
//
// Callers that need per-endpoint middleware chains (as server.go does) can
// use the Register/Login/Refresh/Logout methods directly instead.
func (h *Handler) Mount(r chi.Router) {
	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", h.Register)
		r.Post("/login", h.Login)
		r.Post("/refresh", h.Refresh)
		r.Post("/logout", h.Logout) // caller must add RequireAuth middleware
	})
}

// --- register ---------------------------------------------------------------

type registerReq struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
	// AcceptTerms: 13+ and terms/privacy accepted (DD-231).
	AcceptTerms bool `json:"accept_terms"`
}

// Register is the POST /v1/auth/register handler.
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var body registerReq
	if err := decodeJSON(r, &body); err != nil {
		apierr.Write(w, r, http.StatusBadRequest, "E_INVALID_JSON", err.Error())
		return
	}
	res, err := h.svc.Register(r.Context(), RegisterRequest{
		Email:       body.Email,
		Username:    body.Username,
		Password:    body.Password,
		IPHash:      h.ipHash(httpctx.ClientIP(r)),
		UserAgent:   r.UserAgent(),
		AcceptTerms: body.AcceptTerms,
	})
	if err != nil {
		writeAuthError(w, r, err)
		return
	}
	apierr.WriteJSON(w, http.StatusCreated, sessionResponse(res))
}

// --- login ------------------------------------------------------------------

type loginReq struct {
	EmailOrUsername string `json:"email_or_username"`
	Password        string `json:"password"`
}

// Login is the POST /v1/auth/login handler.
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var body loginReq
	if err := decodeJSON(r, &body); err != nil {
		apierr.Write(w, r, http.StatusBadRequest, "E_INVALID_JSON", err.Error())
		return
	}
	res, err := h.svc.Login(r.Context(), LoginRequest{
		EmailOrUsername: body.EmailOrUsername,
		Password:        body.Password,
		IPHash:          h.ipHash(httpctx.ClientIP(r)),
		UserAgent:       r.UserAgent(),
	})
	if err != nil {
		writeAuthError(w, r, err)
		return
	}
	apierr.WriteJSON(w, http.StatusOK, sessionResponse(res))
}

// --- refresh ----------------------------------------------------------------

type refreshReq struct {
	RefreshToken string `json:"refresh_token"`
}

// Refresh is the POST /v1/auth/refresh handler.
func (h *Handler) Refresh(w http.ResponseWriter, r *http.Request) {
	var body refreshReq
	if err := decodeJSON(r, &body); err != nil {
		apierr.Write(w, r, http.StatusBadRequest, "E_INVALID_JSON", err.Error())
		return
	}
	res, err := h.svc.Refresh(r.Context(), body.RefreshToken)
	if err != nil {
		writeAuthError(w, r, err)
		return
	}
	apierr.WriteJSON(w, http.StatusOK, sessionResponse(res))
}

// --- logout -----------------------------------------------------------------

// Logout is the POST /v1/auth/logout handler; requires an authenticated session.
func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	sid := httpctx.SessionID(r.Context())
	if sid == uuid.Nil {
		apierr.Write(w, r, http.StatusUnauthorized, "E_UNAUTHENTICATED", "no active session")
		return
	}
	if err := h.svc.Logout(r.Context(), sid); err != nil {
		h.log.WarnContext(r.Context(), "logout failed",
			slog.String("correlation_id", chimw.GetReqID(r.Context())),
			slog.String("err", err.Error()))
		apierr.Write(w, r, http.StatusInternalServerError, "E_INTERNAL", "logout failed")
		return
	}
	apierr.WriteJSON(w, http.StatusOK, map[string]bool{"ok": true})
}

// --- helpers ----------------------------------------------------------------

// sessionResponse is the wire shape for Register/Login/Refresh. Timestamps in
// RFC3339 (client parses to Date).
type sessionResp struct {
	AccountID       string `json:"account_id"`
	SessionID       string `json:"session_id"`
	AccessToken     string `json:"access_token"`
	AccessExpiresAt string `json:"access_expires_at"`
	RefreshToken    string `json:"refresh_token"`
	RefreshExpires  string `json:"refresh_expires_at"`
	EmailVerified   bool   `json:"email_verified"`
	// DeletionCancelled: this sign-in cancelled a pending account deletion.
	DeletionCancelled bool `json:"deletion_cancelled,omitempty"`
}

func sessionResponse(r *LoginResult) sessionResp {
	return sessionResp{
		AccountID:         r.AccountID.String(),
		SessionID:         r.SessionID.String(),
		AccessToken:       r.AccessToken,
		AccessExpiresAt:   r.AccessExpiresAt.UTC().Format(time.RFC3339),
		RefreshToken:      r.RefreshToken,
		RefreshExpires:    r.RefreshExpires.UTC().Format(time.RFC3339),
		EmailVerified:     r.EmailVerified,
		DeletionCancelled: r.DeletionCancelled,
	}
}

func decodeJSON(r *http.Request, dst any) error {
	body, err := io.ReadAll(io.LimitReader(r.Body, 32<<10)) // 32 KB auth payloads
	if err != nil {
		return err
	}
	if len(body) == 0 {
		return errors.New("empty body")
	}
	return json.Unmarshal(body, dst)
}

// writeAuthError maps auth domain errors to HTTP status codes.
func writeAuthError(w http.ResponseWriter, r *http.Request, err error) {
	var ve *ValidationError
	switch {
	case errors.As(err, &ve):
		apierr.WriteField(w, r, http.StatusBadRequest, "E_VALIDATION", ve.Message, ve.Field)
	case errors.Is(err, ErrAccountExists):
		apierr.Write(w, r, http.StatusConflict, "E_ACCOUNT_EXISTS", "email or username already taken")
	case errors.Is(err, ErrInvalidCredentials):
		apierr.Write(w, r, http.StatusUnauthorized, "E_INVALID_CREDENTIALS", "invalid credentials")
	case errors.Is(err, ErrAccountLocked):
		apierr.Write(w, r, http.StatusForbidden, "E_ACCOUNT_LOCKED", "account locked")
	case errors.Is(err, ErrAccountBanned):
		apierr.Write(w, r, http.StatusForbidden, "E_ACCOUNT_BANNED", "account banned")
	case errors.Is(err, ErrRefreshTokenNotFound):
		apierr.Write(w, r, http.StatusUnauthorized, "E_REFRESH_TOKEN_UNKNOWN", "unknown refresh token")
	case errors.Is(err, ErrRefreshTokenReuse):
		apierr.Write(w, r, http.StatusUnauthorized, "E_REFRESH_TOKEN_REUSE", "refresh token replayed — session revoked")
	case errors.Is(err, ErrTokenExpired):
		apierr.Write(w, r, http.StatusUnauthorized, "E_TOKEN_EXPIRED", "token expired")
	case errors.Is(err, ErrTokenInvalid):
		apierr.Write(w, r, http.StatusUnauthorized, "E_TOKEN_INVALID", "token invalid")
	case errors.Is(err, ErrSessionRevoked):
		apierr.Write(w, r, http.StatusUnauthorized, "E_SESSION_REVOKED", "session revoked")
	case errors.Is(err, ErrSessionExpired):
		apierr.Write(w, r, http.StatusUnauthorized, "E_SESSION_EXPIRED", "session expired")
	case errors.Is(err, ErrRecoveryTokenInvalid):
		apierr.Write(w, r, http.StatusBadRequest, "E_RECOVERY_TOKEN_INVALID", "link is invalid or expired")
	case errors.Is(err, ErrExportLimit):
		apierr.Write(w, r, http.StatusTooManyRequests, "E_RATE_LIMITED", "too many export links, try again later")
	case errors.Is(err, ErrBusy):
		w.Header().Set("Retry-After", "1")
		apierr.Write(w, r, http.StatusServiceUnavailable, "E_BUSY", "server busy, retry shortly")
	case errors.Is(err, ErrEmailAlreadyVerified):
		apierr.Write(w, r, http.StatusConflict, "E_ALREADY_VERIFIED", "e-mail already verified")
	default:
		apierr.Write(w, r, http.StatusInternalServerError, "E_INTERNAL", "internal error")
	}
}
