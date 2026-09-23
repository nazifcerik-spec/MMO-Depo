package wsgateway

import (
	"context"
	"encoding/json"
	"log/slog"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// clientSendBuf is the per-connection outgoing buffer. Chosen small so a slow
// reader is dropped quickly rather than piling MB of undelivered events.
const clientSendBuf = 64

// Timings — kept tight to catch dead connections fast without spamming pings.
const (
	writeWait     = 10 * time.Second
	pongWait      = 60 * time.Second
	pingInterval  = 25 * time.Second
	readSizeLimit = 1 << 16 // 64 KiB; client messages are short subscribe / ping frames
)

// Client wraps one open WebSocket. Owns a read pump (server → discard) and a
// write pump (Hub-fanout channel → wire).
type Client struct {
	id        uuid.UUID
	accountID uuid.UUID
	conn      *websocket.Conn
	send      chan Event
	log       *slog.Logger

	closeOnce sync.Once
}

// NewClient wires a fresh client around an established websocket connection.
// The caller is responsible for calling client.Run to start the read + write
// loops; on return, the connection is closed and the caller should invoke
// the unregister fn returned by Hub.Register.
func NewClient(conn *websocket.Conn, accountID uuid.UUID, logger *slog.Logger) *Client {
	if logger == nil {
		logger = slog.Default()
	}
	return &Client{
		id:        uuid.New(),
		accountID: accountID,
		conn:      conn,
		send:      make(chan Event, clientSendBuf),
		log:       logger,
	}
}

// ID exposes the connection-scoped id for logging + presence diagnostics.
func (c *Client) ID() uuid.UUID { return c.id }

// AccountID exposes the authenticated account this connection belongs to.
func (c *Client) AccountID() uuid.UUID { return c.accountID }

// closeSlow is called by the Hub when this client's send buffer overflows.
// Closes the connection so both pumps unwind; safe to call multiple times.
func (c *Client) closeSlow() {
	c.closeOnce.Do(func() {
		_ = c.conn.WriteControl(
			websocket.CloseMessage,
			websocket.FormatCloseMessage(websocket.ClosePolicyViolation, "slow client"),
			time.Now().Add(writeWait),
		)
		_ = c.conn.Close()
	})
}

// Run blocks until the read pump exits. Starts the write pump in a goroutine.
// When Run returns, both pumps have stopped and the connection is closed.
func (c *Client) Run(ctx context.Context, welcome Event) {
	writeDone := make(chan struct{})
	go c.writePump(writeDone, welcome)

	c.readPump(ctx)

	// Signal write pump to exit by closing send channel — safe from here
	// because readPump has exited so nothing else references send.
	close(c.send)
	<-writeDone
}

// readPump receives frames from the client. v1 discards all inbound messages
// (client → server messaging is Phase 02: subscribe filters, presence hints).
// Keeps the connection alive by extending the pong deadline on every pong.
func (c *Client) readPump(ctx context.Context) {
	defer c.closeSlow()

	c.conn.SetReadLimit(readSizeLimit)
	_ = c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		_ = c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}
		if _, _, err := c.conn.ReadMessage(); err != nil {
			// Normal close, timeout, or ungraceful drop — all end the pump.
			if websocket.IsUnexpectedCloseError(err,
				websocket.CloseNormalClosure,
				websocket.CloseGoingAway,
				websocket.CloseAbnormalClosure) {
				c.log.Debug("wsgateway: read error", "client", c.id.String(), "err", err.Error())
			}
			return
		}
	}
}

// writePump drains send + heartbeats. Exits when send is closed.
func (c *Client) writePump(done chan<- struct{}, welcome Event) {
	defer close(done)
	ticker := time.NewTicker(pingInterval)
	defer ticker.Stop()

	// Welcome frame first so the client immediately knows the connection is
	// live and authenticated. Best-effort — if the socket is already broken,
	// the send loop below picks up the same error and exits.
	if welcome.Kind != "" {
		if err := c.writeEvent(welcome); err != nil {
			return
		}
	}

	for {
		select {
		case ev, ok := <-c.send:
			if !ok {
				// Hub closed our channel — send a normal close frame and stop.
				_ = c.conn.WriteControl(
					websocket.CloseMessage,
					websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""),
					time.Now().Add(writeWait),
				)
				return
			}
			if err := c.writeEvent(ev); err != nil {
				return
			}
		case <-ticker.C:
			_ = c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (c *Client) writeEvent(ev Event) error {
	_ = c.conn.SetWriteDeadline(time.Now().Add(writeWait))
	// WriteJSON handles the JSON encoding + text frame. Payload is already
	// json.RawMessage, so this is one marshal at the outer envelope only.
	if err := c.conn.WriteJSON(ev); err != nil {
		return err
	}
	return nil
}

// Silence unused import in some builds — json.RawMessage lives on Event.
var _ = json.RawMessage(nil)
