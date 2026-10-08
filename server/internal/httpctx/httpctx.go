// Package httpctx defines the typed context keys shared between middleware and
// handlers. Middleware writes; handlers read. Never expose the raw key.
package httpctx

import (
	"context"
	"net"
	"net/http"

	"github.com/google/uuid"
)

type ctxKey int

const (
	accountIDKey ctxKey = iota + 1
	sessionIDKey
	characterIDKey
	clientIPKey
)

// WithAccountID stores the authenticated account id in ctx. Auth middleware
// calls this after a valid session is resolved.
func WithAccountID(ctx context.Context, id uuid.UUID) context.Context {
	return context.WithValue(ctx, accountIDKey, id)
}

// AccountID returns the authenticated account id, or uuid.Nil for anonymous.
// Anonymous requests use uuid.Nil as the idempotency / rate-limit scope key.
func AccountID(ctx context.Context) uuid.UUID {
	if v, ok := ctx.Value(accountIDKey).(uuid.UUID); ok {
		return v
	}
	return uuid.Nil
}

// WithSessionID stores the JWT session id in ctx. Used by logout.
func WithSessionID(ctx context.Context, id uuid.UUID) context.Context {
	return context.WithValue(ctx, sessionIDKey, id)
}

// SessionID returns the authenticated session id, or uuid.Nil for anonymous.
func SessionID(ctx context.Context) uuid.UUID {
	if v, ok := ctx.Value(sessionIDKey).(uuid.UUID); ok {
		return v
	}
	return uuid.Nil
}

// WithCharacterID stores the character bound to this request (set by the
// character selection endpoint or a per-request X-Character-Id header).
// Character-scope rate limits and audit rely on this.
func WithCharacterID(ctx context.Context, id uuid.UUID) context.Context {
	return context.WithValue(ctx, characterIDKey, id)
}

// CharacterID returns the bound character id or uuid.Nil.
func CharacterID(ctx context.Context) uuid.UUID {
	if v, ok := ctx.Value(characterIDKey).(uuid.UUID); ok {
		return v
	}
	return uuid.Nil
}

// WithClientIP stores the caller IP resolved by the RealIP middleware.
func WithClientIP(ctx context.Context, ip string) context.Context {
	return context.WithValue(ctx, clientIPKey, ip)
}

// ClientIP returns the caller IP: the one RealIP resolved (forwarded headers
// honoured only from trusted proxies), else the connection's peer address.
func ClientIP(r *http.Request) string {
	if v, ok := r.Context().Value(clientIPKey).(string); ok && v != "" {
		return v
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}
