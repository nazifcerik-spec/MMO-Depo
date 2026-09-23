package wsgateway

import (
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"

	"github.com/otter/server/internal/domain/auth"
	"github.com/otter/server/internal/http/apierr"
)

// Upgrader is shared across requests. CheckOrigin defaults to permissive here
// because auth is enforced via the token query param — the cross-origin risk
// is limited to "someone leaks their own token", which is user-side. Prod
// deployment should tighten this to an allowlist matching the client origin.
var upgrader = websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

// Handler bundles the /v1/ws upgrade endpoint + presence read.
type Handler struct {
	hub  *Hub
	auth *auth.Service
	log  *slog.Logger
}

func NewHandler(hub *Hub, authSvc *auth.Service, logger *slog.Logger) *Handler {
	if logger == nil {
		logger = slog.Default()
	}
	return &Handler{hub: hub, auth: authSvc, log: logger}
}

// Upgrade accepts the WebSocket handshake at GET /v1/ws?token=<jwt>. Auth
// happens via the token query param because browser WebSocket clients cannot
// send custom headers on the handshake — a well-established workaround.
//
// After the handshake succeeds the client is registered in the Hub and its
// read/write pumps run until the connection dies. A welcome event is queued
// so the client sees "you are online" immediately.
func (h *Handler) Upgrade(w http.ResponseWriter, r *http.Request) {
	tok := r.URL.Query().Get("token")
	if tok == "" {
		apierr.Write(w, r, http.StatusUnauthorized, "E_UNAUTHENTICATED", "missing token")
		return
	}
	accountID, sessionID, err := h.auth.Verify(r.Context(), tok)
	if err != nil {
		apierr.Write(w, r, http.StatusUnauthorized, "E_UNAUTHENTICATED", "invalid token")
		return
	}

	// Handshake — after this, ANY error must be logged (no ResponseWriter).
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		h.log.WarnContext(r.Context(), "wsgateway: upgrade failed",
			"account", accountID.String(), "err", err.Error())
		return
	}

	c := NewClient(conn, accountID, h.log)
	unregister := h.hub.Register(accountID, c)
	defer unregister()

	welcome := MarshalEvent("system.welcome", map[string]any{
		"account_id": accountID.String(),
		"session_id": sessionID.String(),
		"client_id":  c.ID().String(),
	}, h.log)

	c.Run(r.Context(), welcome)
}

// Presence returns a JSON snapshot of online account ids + client count.
// Dev/admin diagnostic — a real "who is online" list for the friend system
// arrives with Sprint 21.
func (h *Handler) Presence(w http.ResponseWriter, r *http.Request) {
	online := h.hub.OnlineAccounts()
	ids := make([]string, len(online))
	for i, id := range online {
		ids[i] = id.String()
	}
	apierr.WriteJSON(w, http.StatusOK, map[string]any{
		"count":    len(ids),
		"clients":  h.hub.ClientCount(),
		"accounts": ids,
	})
}

// --- Dev broadcast helper ---------------------------------------------------

// DevBroadcast is the dev-only "poke the hub" endpoint — sends a system event
// to a specific account (or everyone if account_id absent) so the dev console
// can prove the pipe end-to-end without a real domain event.
func (h *Handler) DevBroadcast(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	msg := q.Get("message")
	if msg == "" {
		msg = "hello from server"
	}
	ev := MarshalEvent("system.dev_broadcast", map[string]any{"message": msg}, h.log)

	var delivered int
	if accountStr := q.Get("account_id"); accountStr != "" {
		id, err := uuid.Parse(accountStr)
		if err != nil {
			apierr.Write(w, r, http.StatusBadRequest, "E_INVALID_ID", "account_id: "+err.Error())
			return
		}
		delivered = h.hub.Publish(id, ev)
	} else {
		delivered = h.hub.Broadcast(ev)
	}
	apierr.WriteJSON(w, http.StatusOK, map[string]any{
		"delivered": delivered,
		"kind":      ev.Kind,
	})
}
