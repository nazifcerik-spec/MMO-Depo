package wsgateway

import (
	"github.com/google/uuid"
)

// Publisher is the narrow interface downstream services import. Kept here so
// hunt/dungeon/etc do not depend on gorilla/websocket, only on this tiny
// contract — swap the Hub for a Redis-fanout bridge later without touching
// callers.
type Publisher interface {
	// Publish fans an event out to every open client for the account.
	// Non-blocking (dropped on slow readers). Returns the number of clients
	// the event was queued into.
	Publish(accountID uuid.UUID, ev Event) int
}

// Hub already satisfies Publisher via its Publish method — the compile-time
// assertion here catches an accidental signature drift.
var _ Publisher = (*Hub)(nil)

// NoopPublisher is a safe default when the ws gateway isn't wired (tests,
// worker binary). Every service that takes a Publisher should default to
// this rather than nil-check every call site.
type NoopPublisher struct{}

func (NoopPublisher) Publish(_ uuid.UUID, _ Event) int { return 0 }

// --- KindPayloadAdapter -----------------------------------------------------

// KindPayloadAdapter wraps a Hub so callers that don't want to build the Event
// envelope themselves can just pass (kind, payload). The wsgateway.Event
// timestamp gets stamped inside. This is the interface hunt/dungeon HTTP
// handlers depend on — they never see json.RawMessage or Event directly.
type KindPayloadAdapter struct {
	hub *Hub
}

// NewKindPayloadAdapter returns an adapter that satisfies the
// {hunt,dungeon}.EventPublisher interfaces (see those packages).
func NewKindPayloadAdapter(hub *Hub) *KindPayloadAdapter {
	return &KindPayloadAdapter{hub: hub}
}

// Publish marshals the payload and fans out to the account's clients.
func (a *KindPayloadAdapter) Publish(accountID uuid.UUID, kind string, payload any) int {
	ev := MarshalEvent(kind, payload, a.hub.log)
	return a.hub.Publish(accountID, ev)
}

// BroadcastKind marshals the payload and sends it to every connected client.
func (a *KindPayloadAdapter) BroadcastKind(kind string, payload any) int {
	return a.hub.Broadcast(MarshalEvent(kind, payload, a.hub.log))
}
