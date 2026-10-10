package admin

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"

	"github.com/otter/server/internal/clienterr"
	"github.com/otter/server/internal/domain/auth"
	"github.com/otter/server/internal/economy"
	"github.com/otter/server/internal/http/apierr"
	"github.com/otter/server/internal/httpctx"
	"github.com/otter/server/internal/market"
	"github.com/otter/server/internal/support"
)

// MarketCanceller is the market capability admin needs.
type MarketCanceller interface {
	AdminCancel(ctx context.Context, listingID uuid.UUID, refundFee bool) (*market.Listing, error)
}

type Handler struct {
	svc        *Service
	market     MarketCanceller
	support    SupportDesk
	clientErrs ClientErrors
}

func NewHandler(svc *Service, m MarketCanceller) *Handler { return &Handler{svc: svc, market: m} }

// SupportDesk is the ticket queue the support routes work on (DD-234).
type SupportDesk interface {
	Queue(ctx context.Context, state string, limit int) ([]support.Ticket, error)
	Answer(ctx context.Context, admin, ticket uuid.UUID, reply string, closeIt bool) (*support.Ticket, error)
}

// WithSupport enables the support-ticket routes.
func (h *Handler) WithSupport(d SupportDesk) *Handler { h.support = d; return h }

// ClientErrors is the grouped app crash list (DD-235).
type ClientErrors interface {
	Recent(ctx context.Context, platform string, limit int) ([]clienterr.Group, error)
}

// WithClientErrors enables GET /internal/ops/client-errors.
func (h *Handler) WithClientErrors(c ClientErrors) *Handler { h.clientErrs = c; return h }

type ctxKey struct{}

func current(r *http.Request) *Admin { a, _ := r.Context().Value(ctxKey{}).(*Admin); return a }

// RequireAdmin authenticates the bearer token and stores the admin in ctx.
func (h *Handler) RequireAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
		a, err := h.svc.Authenticate(r.Context(), token)
		if err != nil {
			writeErr(w, r, err)
			return
		}
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxKey{}, a)))
	})
}

// allow checks a permission (and step-up when critical); writes 403 itself.
func (h *Handler) allow(w http.ResponseWriter, r *http.Request, ok bool, stepUp bool) (*Admin, bool) {
	a := current(r)
	if a == nil {
		writeErr(w, r, ErrUnauthorized)
		return nil, false
	}
	if !ok {
		h.svc.Audit(r.Context(), a, AuditEntry{Action: "denied:" + r.Method + " " + r.URL.Path, RequestID: chimw.GetReqID(r.Context())})
		writeErr(w, r, ErrForbidden)
		return nil, false
	}
	if stepUp && !a.SteppedUp(time.Now()) {
		writeErr(w, r, ErrStepUpRequired)
		return nil, false
	}
	return a, true
}

func perm(r *http.Request, p string) bool {
	a := current(r)
	return a != nil && a.Perms.Has(p)
}

type body struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	Code      string `json:"code"`
	Reason    string `json:"reason"`
	Hours     int    `json:"hours"`
	Days      int    `json:"days"`
	Role      string `json:"role"`
	Grant     bool   `json:"grant"`
	RefundFee bool   `json:"refund_fee"`
	StartsIn  int    `json:"starts_in_minutes"`
	Reply     string `json:"reply"`
	Close     bool   `json:"close"`
}

func decode(w http.ResponseWriter, r *http.Request) (body, bool) {
	var b body
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 1<<16)).Decode(&b); err != nil {
		apierr.Write(w, r, http.StatusBadRequest, "E_INVALID_JSON", err.Error())
		return b, false
	}
	return b, true
}

func pathID(w http.ResponseWriter, r *http.Request, key string) (uuid.UUID, bool) {
	id, err := uuid.Parse(chi.URLParam(r, key))
	if err != nil {
		apierr.Write(w, r, http.StatusBadRequest, "E_INVALID_ID", key+": "+err.Error())
	}
	return id, err == nil
}

func reply(w http.ResponseWriter, r *http.Request, v any, err error) {
	if err != nil {
		writeErr(w, r, err)
		return
	}
	apierr.WriteJSON(w, http.StatusOK, v)
}

func queryUUID(r *http.Request, k string) *uuid.UUID {
	if id, err := uuid.Parse(r.URL.Query().Get(k)); err == nil {
		return &id
	}
	return nil
}

// --- session ----------------------------------------------------------------

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	b, ok := decode(w, r)
	if !ok {
		return
	}
	sum := sha256.Sum256([]byte(httpctx.ClientIP(r)))
	token, err := h.svc.Login(r.Context(), b.Email, b.Password, b.Code, r.UserAgent(), sum[:])
	reply(w, r, map[string]any{"token": token, "expires_in": int(SessionTTL.Seconds())}, err)
}

func (h *Handler) Me(w http.ResponseWriter, r *http.Request) {
	a := current(r)
	perms := make([]string, 0, len(a.Perms))
	for p := range a.Perms {
		perms = append(perms, p)
	}
	apierr.WriteJSON(w, http.StatusOK, map[string]any{"admin": a, "permissions": perms})
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	reply(w, r, map[string]any{"logged_out": true}, h.svc.Logout(r.Context(), current(r)))
}

func (h *Handler) StepUp(w http.ResponseWriter, r *http.Request) {
	b, ok := decode(w, r)
	if !ok {
		return
	}
	until, err := h.svc.StepUp(r.Context(), current(r), b.Code)
	reply(w, r, map[string]any{"stepup_until": until}, err)
}

func (h *Handler) Users(w http.ResponseWriter, r *http.Request) {
	if _, ok := h.allow(w, r, perm(r, PermRoleGrant), false); !ok {
		return
	}
	list, err := h.svc.Users(r.Context())
	reply(w, r, map[string]any{"users": list}, err)
}

func (h *Handler) SetRole(w http.ResponseWriter, r *http.Request) {
	a, ok := h.allow(w, r, perm(r, PermRoleGrant), true)
	if !ok {
		return
	}
	id, ok := pathID(w, r, "admin_id")
	if !ok {
		return
	}
	b, ok := decode(w, r)
	if !ok {
		return
	}
	reply(w, r, map[string]any{"role": b.Role, "granted": b.Grant}, h.svc.SetRole(r.Context(), a, id, b.Role, b.Grant, b.Reason))
}

func (h *Handler) Audit(w http.ResponseWriter, r *http.Request) {
	if _, ok := h.allow(w, r, perm(r, PermAuditRead), false); !ok {
		return
	}
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	rows, err := h.svc.AuditLog(r.Context(), queryUUID(r, "actor"), queryUUID(r, "target"), limit)
	reply(w, r, map[string]any{"count": len(rows), "actions": rows}, err)
}

// --- players / characters ---------------------------------------------------

func (h *Handler) SearchPlayers(w http.ResponseWriter, r *http.Request) {
	a, ok := h.allow(w, r, perm(r, PermPlayerRead), false)
	if !ok {
		return
	}
	list, err := h.svc.SearchPlayers(r.Context(), a, r.URL.Query().Get("q"))
	reply(w, r, map[string]any{"count": len(list), "players": list}, err)
}

func (h *Handler) Player(w http.ResponseWriter, r *http.Request) {
	a, ok := h.allow(w, r, perm(r, PermPlayerRead), false)
	if !ok {
		return
	}
	id, ok := pathID(w, r, "account_id")
	if !ok {
		return
	}
	d, err := h.svc.Player(r.Context(), a, id)
	reply(w, r, d, err)
}

func (h *Handler) Mute(w http.ResponseWriter, r *http.Request) {
	b, ok := decode(w, r)
	if !ok {
		return
	}
	a, ok := h.allow(w, r, current(r) != nil && current(r).Perms.CanMute(b.Hours), b.Hours <= 0)
	if !ok {
		return
	}
	id, ok := pathID(w, r, "account_id")
	if !ok {
		return
	}
	s, err := h.svc.Sanction(r.Context(), a, id, "mute", time.Duration(b.Hours)*time.Hour, b.Reason)
	reply(w, r, s, err)
}

func (h *Handler) Ban(w http.ResponseWriter, r *http.Request) {
	b, ok := decode(w, r)
	if !ok {
		return
	}
	a, ok := h.allow(w, r, current(r) != nil && current(r).Perms.CanBan(b.Days), b.Days <= 0)
	if !ok {
		return
	}
	id, ok := pathID(w, r, "account_id")
	if !ok {
		return
	}
	s, err := h.svc.Sanction(r.Context(), a, id, "ban", time.Duration(b.Days)*24*time.Hour, b.Reason)
	reply(w, r, s, err)
}

func (h *Handler) Lift(kind, needed string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		a, ok := h.allow(w, r, perm(r, needed), false)
		if !ok {
			return
		}
		id, ok := pathID(w, r, "account_id")
		if !ok {
			return
		}
		b, ok := decode(w, r)
		if !ok {
			return
		}
		n, err := h.svc.LiftSanctions(r.Context(), a, id, kind, b.Reason)
		reply(w, r, map[string]any{"revoked": n}, err)
	}
}

func (h *Handler) Character(w http.ResponseWriter, r *http.Request) {
	a, ok := h.allow(w, r, perm(r, PermCharacterRead), false)
	if !ok {
		return
	}
	id, ok := pathID(w, r, "character_id")
	if !ok {
		return
	}
	d, err := h.svc.Character(r.Context(), a, id)
	reply(w, r, d, err)
}

func (h *Handler) Freeze(freeze bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		a, ok := h.allow(w, r, perm(r, PermCharacterFreeze), false)
		if !ok {
			return
		}
		id, ok := pathID(w, r, "character_id")
		if !ok {
			return
		}
		b, ok := decode(w, r)
		if !ok {
			return
		}
		reply(w, r, map[string]any{"frozen": freeze}, h.svc.SetFrozen(r.Context(), a, id, freeze, b.Reason))
	}
}

// --- economy / market / guild / events ---------------------------------------

func (h *Handler) Ledger(w http.ResponseWriter, r *http.Request) {
	if _, ok := h.allow(w, r, perm(r, PermLedgerRead), false); !ok {
		return
	}
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	list, err := h.svc.Ledger(r.Context(), queryUUID(r, "account"), queryUUID(r, "txn"), limit)
	reply(w, r, map[string]any{"count": len(list), "entries": list}, err)
}

func (h *Handler) Reverse(w http.ResponseWriter, r *http.Request) {
	ok := perm(r, PermLedgerReversal100k) || perm(r, PermLedgerReversalAny)
	a, ok := h.allow(w, r, ok, true)
	if !ok {
		return
	}
	id, ok := pathID(w, r, "txn_id")
	if !ok {
		return
	}
	b, ok := decode(w, r)
	if !ok {
		return
	}
	txn, err := h.svc.Reverse(r.Context(), a, id, b.Reason)
	reply(w, r, map[string]any{"reversal_txn_id": txn}, err)
}

func (h *Handler) CancelListing(w http.ResponseWriter, r *http.Request) {
	a, ok := h.allow(w, r, perm(r, PermMarketCancel), false)
	if !ok {
		return
	}
	id, ok := pathID(w, r, "listing_id")
	if !ok {
		return
	}
	b, ok := decode(w, r)
	if !ok {
		return
	}
	if strings.TrimSpace(b.Reason) == "" {
		writeErr(w, r, ErrReasonRequired)
		return
	}
	l, err := h.market.AdminCancel(r.Context(), id, b.RefundFee)
	if err == nil {
		h.svc.Audit(r.Context(), a, AuditEntry{Action: "market.listing_cancel", TargetKind: "listing", TargetID: &id,
			After: map[string]any{"refund_fee": b.RefundFee}, Reason: b.Reason})
	}
	reply(w, r, l, err)
}

func (h *Handler) DisbandGuild(w http.ResponseWriter, r *http.Request) {
	a, ok := h.allow(w, r, perm(r, PermGuildDisband), true)
	if !ok {
		return
	}
	id, ok := pathID(w, r, "guild_id")
	if !ok {
		return
	}
	b, ok := decode(w, r)
	if !ok {
		return
	}
	reply(w, r, map[string]any{"disbanded": id}, h.svc.DisbandGuild(r.Context(), a, id, b.Reason))
}

func (h *Handler) TriggerEvent(w http.ResponseWriter, r *http.Request) {
	a, ok := h.allow(w, r, perm(r, PermEventTrigger), false)
	if !ok {
		return
	}
	b, ok := decode(w, r)
	if !ok {
		return
	}
	id, err := h.svc.TriggerWorldEvent(r.Context(), a, chi.URLParam(r, "template_code"),
		time.Duration(b.StartsIn)*time.Minute, b.Reason)
	reply(w, r, map[string]any{"spawn_id": id}, err)
}

func writeErr(w http.ResponseWriter, r *http.Request, err error) {
	type m struct {
		err    error
		status int
		code   string
	}
	for _, c := range []m{
		{ErrBadCredentials, 401, "E_BAD_CREDENTIALS"},
		{auth.ErrBusy, 503, "E_BUSY"},
		{ErrUnauthorized, 401, "E_UNAUTHENTICATED"},
		{ErrForbidden, 403, "E_FORBIDDEN"},
		{ErrStepUpRequired, 403, "E_STEPUP_REQUIRED"},
		{ErrBadCode, 401, "E_BAD_CODE"},
		{ErrReasonRequired, 400, "E_REASON_REQUIRED"},
		{ErrNotFound, 404, "E_NOT_FOUND"},
		{ErrInvalid, 400, "E_INVALID"},
		{market.ErrNotFound, 404, "E_NOT_FOUND"},
		{market.ErrUnavailable, 409, "E_LISTING_UNAVAILABLE"},
		{economy.ErrInsufficientBalance, 409, "E_INSUFFICIENT_BALANCE"},
	} {
		if errors.Is(err, c.err) {
			apierr.Write(w, r, c.status, c.code, err.Error())
			return
		}
	}
	apierr.Write(w, r, http.StatusInternalServerError, "E_INTERNAL", err.Error())
}

// Routes mounts /internal/admin and /internal/ops on r.
func (h *Handler) Routes(r chi.Router, limit func(http.Handler) http.Handler, loginLimit func(http.Handler) http.Handler) {
	r.With(loginLimit).Post("/internal/admin/login", h.Login)
	r.Group(func(r chi.Router) {
		r.Use(limit, h.RequireAdmin)
		r.Get("/internal/admin/session/me", h.Me)
		r.Post("/internal/admin/logout", h.Logout)
		r.Post("/internal/admin/stepup", h.StepUp)
		r.Get("/internal/admin/users", h.Users)
		r.Post("/internal/admin/users/{admin_id}/role", h.SetRole)

		r.Get("/internal/ops/audit", h.Audit)
		r.Get("/internal/ops/players", h.SearchPlayers)
		r.Get("/internal/ops/players/{account_id}", h.Player)
		r.Post("/internal/ops/players/{account_id}/mute", h.Mute)
		r.Post("/internal/ops/players/{account_id}/ban", h.Ban)
		r.Post("/internal/ops/players/{account_id}/unmute", h.Lift("mute", PermPlayerMute24h))
		r.Post("/internal/ops/players/{account_id}/unban", h.Lift("ban", PermPlayerBan7d))
		r.Get("/internal/ops/characters/{character_id}", h.Character)
		r.Post("/internal/ops/characters/{character_id}/freeze", h.Freeze(true))
		r.Post("/internal/ops/characters/{character_id}/unfreeze", h.Freeze(false))
		r.Get("/internal/ops/ledger", h.Ledger)
		r.Post("/internal/ops/ledger/{txn_id}/reverse", h.Reverse)
		r.Post("/internal/ops/market/{listing_id}/cancel", h.CancelListing)
		r.Post("/internal/ops/guilds/{guild_id}/disband", h.DisbandGuild)
		r.Post("/internal/ops/world-events/{template_code}/trigger", h.TriggerEvent)
		if h.support != nil {
			r.Get("/internal/ops/support/tickets", h.SupportQueue)
			r.Post("/internal/ops/support/tickets/{ticket_id}/reply", h.SupportReply)
		}
		if h.clientErrs != nil {
			r.Get("/internal/ops/client-errors", h.ClientErrorList)
		}
	})
}

// ClientErrorList — GET /internal/ops/client-errors?platform=android&limit=N.
func (h *Handler) ClientErrorList(w http.ResponseWriter, r *http.Request) {
	a, ok := h.allow(w, r, perm(r, PermClientErrors), false)
	if !ok {
		return
	}
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	list, err := h.clientErrs.Recent(r.Context(), r.URL.Query().Get("platform"), limit)
	if err != nil {
		clienterr.WriteErr(w, r, err)
		return
	}
	h.svc.Audit(r.Context(), a, AuditEntry{Action: "client_errors.read", RequestID: chimw.GetReqID(r.Context())})
	apierr.WriteJSON(w, http.StatusOK, map[string]any{"errors": list})
}

// SupportQueue — GET /internal/ops/support/tickets?state=open|answered|closed&limit=N.
func (h *Handler) SupportQueue(w http.ResponseWriter, r *http.Request) {
	a, ok := h.allow(w, r, perm(r, PermSupportTickets), false)
	if !ok {
		return
	}
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	list, err := h.support.Queue(r.Context(), r.URL.Query().Get("state"), limit)
	if err != nil {
		support.WriteErr(w, r, err)
		return
	}
	h.svc.Audit(r.Context(), a, AuditEntry{Action: "support.queue", RequestID: chimw.GetReqID(r.Context())})
	apierr.WriteJSON(w, http.StatusOK, map[string]any{"tickets": list})
}

// SupportReply — POST /internal/ops/support/tickets/{ticket_id}/reply {reply, close}.
func (h *Handler) SupportReply(w http.ResponseWriter, r *http.Request) {
	b, ok := decode(w, r)
	if !ok {
		return
	}
	a, ok := h.allow(w, r, perm(r, PermSupportTickets), false)
	if !ok {
		return
	}
	id, ok := pathID(w, r, "ticket_id")
	if !ok {
		return
	}
	t, err := h.support.Answer(r.Context(), a.ID, id, b.Reply, b.Close)
	if err != nil {
		support.WriteErr(w, r, err)
		return
	}
	h.svc.Audit(r.Context(), a, AuditEntry{Action: "support.reply", TargetKind: "support_ticket", TargetID: &id,
		After: map[string]any{"state": t.State}, RequestID: chimw.GetReqID(r.Context())})
	apierr.WriteJSON(w, http.StatusOK, t)
}
