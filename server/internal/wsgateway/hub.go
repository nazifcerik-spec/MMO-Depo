// Package wsgateway is the WebSocket presence + fanout hub (Sprint 18).
//
// Design:
//   - One shared Hub per API process. Each authenticated WebSocket becomes a
//     Client registered under (accountID → set of clients) so a user with
//     multiple tabs / devices receives every event exactly once per tab.
//   - Publish is fire-and-forget with a small per-client send channel; a slow
//     reader gets dropped rather than blocking the sender. This preserves
//     backpressure at the leaf, not at the source.
//   - Presence = `hub.OnlineAccounts()`, a snapshot of which account ids have
//     at least one open connection. Cheap enough (map iteration under lock)
//     for a debug endpoint; a broadcast presence stream is a Phase 02 concern.
//   - Multi-instance fanout (Redis pub/sub bridge) is intentionally NOT here.
//     The Hub is authoritative for one process; a bridge that mirrors Publish
//     across pods lands with Sprint 25 (world events) when it first matters.
package wsgateway

import (
	"encoding/json"
	"log/slog"
	"sync"
	"time"

	"github.com/google/uuid"
)

// Event is the wire shape server → client. Kind is a short tag ("hunt.reveal",
// "wallet.updated", "system.ping"); Payload is any JSON-serialisable struct.
type Event struct {
	Kind      string          `json:"kind"`
	Payload   json.RawMessage `json:"payload,omitempty"`
	Timestamp time.Time       `json:"timestamp"`
}

// Marshal converts a typed payload into an Event ready for Publish. Returns
// an event with an empty payload + logs the marshal error rather than failing
// the whole Publish path — a malformed one-off payload should not silence the
// broadcast pipe.
func MarshalEvent(kind string, payload any, logger *slog.Logger) Event {
	if payload == nil {
		return Event{Kind: kind, Timestamp: time.Now().UTC()}
	}
	raw, err := json.Marshal(payload)
	if err != nil {
		if logger != nil {
			logger.Warn("wsgateway: marshal payload failed", "kind", kind, "err", err.Error())
		}
		return Event{Kind: kind, Timestamp: time.Now().UTC()}
	}
	return Event{Kind: kind, Payload: raw, Timestamp: time.Now().UTC()}
}

// Hub is the in-process fan-out. Safe for concurrent Publish + Register.
type Hub struct {
	mu        sync.RWMutex
	byAccount map[uuid.UUID]map[*Client]struct{}
	log       *slog.Logger
}

// NewHub constructs an empty Hub.
func NewHub(logger *slog.Logger) *Hub {
	if logger == nil {
		logger = slog.Default()
	}
	return &Hub{
		byAccount: make(map[uuid.UUID]map[*Client]struct{}),
		log:       logger,
	}
}

// Register adds a client for the given accountID. Returns an unregister fn
// callers store and defer.
func (h *Hub) Register(accountID uuid.UUID, c *Client) func() {
	h.mu.Lock()
	set, ok := h.byAccount[accountID]
	if !ok {
		set = make(map[*Client]struct{})
		h.byAccount[accountID] = set
	}
	set[c] = struct{}{}
	h.mu.Unlock()

	return func() {
		h.mu.Lock()
		if set, ok := h.byAccount[accountID]; ok {
			delete(set, c)
			if len(set) == 0 {
				delete(h.byAccount, accountID)
			}
		}
		h.mu.Unlock()
	}
}

// Publish delivers the event to every open client for the given account. A
// client whose send buffer is full is dropped (its writer goroutine will exit
// on the next read pump loop). Publish never blocks — safe to call from
// request handlers.
func (h *Hub) Publish(accountID uuid.UUID, ev Event) int {
	h.mu.RLock()
	set := h.byAccount[accountID]
	// Snapshot to avoid holding the lock through the send loop.
	targets := make([]*Client, 0, len(set))
	for c := range set {
		targets = append(targets, c)
	}
	h.mu.RUnlock()

	delivered := 0
	for _, c := range targets {
		select {
		case c.send <- ev:
			delivered++
		default:
			// Buffer full — drop this client. It will exit on next read cycle.
			h.log.Warn("wsgateway: dropping slow client",
				"account", accountID.String(),
				"client_id", c.id.String(),
			)
			c.closeSlow()
		}
	}
	return delivered
}

// Broadcast delivers the event to EVERY connected client. Used for global
// system events (server-wide announcement, dev poke). Prefer Publish when the
// target is a specific account.
func (h *Hub) Broadcast(ev Event) int {
	h.mu.RLock()
	targets := make([]*Client, 0, 64)
	for _, set := range h.byAccount {
		for c := range set {
			targets = append(targets, c)
		}
	}
	h.mu.RUnlock()

	delivered := 0
	for _, c := range targets {
		select {
		case c.send <- ev:
			delivered++
		default:
			h.log.Warn("wsgateway: dropping slow client on broadcast",
				"client_id", c.id.String())
			c.closeSlow()
		}
	}
	return delivered
}

// OnlineAccounts snapshots the account ids with at least one open connection.
// Used by the presence endpoint; the returned slice is a fresh copy.
func (h *Hub) OnlineAccounts() []uuid.UUID {
	h.mu.RLock()
	defer h.mu.RUnlock()
	out := make([]uuid.UUID, 0, len(h.byAccount))
	for id := range h.byAccount {
		out = append(out, id)
	}
	return out
}

// ClientCount returns the total number of active connections across all
// accounts. Cheap — a sum over one map layer.
func (h *Hub) ClientCount() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	n := 0
	for _, set := range h.byAccount {
		n += len(set)
	}
	return n
}
