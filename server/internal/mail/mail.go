// Package mail delivers transactional e-mail (account recovery, DD-219) from
// the ops.email_outbox table. Uses only the standard library: net/smtp with
// STARTTLS (port 587) or implicit TLS (465).
package mail

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/tls"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"mime"
	"net"
	"net/http"
	"net/smtp"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Sender delivers one plain-text message.
type Sender interface {
	Send(ctx context.Context, to, subject, body string) error
}

// MaxAttempts before a message is given up on (its body is then cleared).
const MaxAttempts = 5

// giveUpAfter: recovery links expire within 24h, so older mail is pointless.
const giveUpAfter = 24 * time.Hour

// LogSender is the no-SMTP fallback. In dev it logs the full body so links can
// be clicked from the console; anywhere else it refuses (bodies carry
// one-time tokens that must never reach logs).
type LogSender struct {
	Dev bool
	Log *slog.Logger
}

func (s LogSender) Send(ctx context.Context, to, subject, body string) error {
	if !s.Dev {
		return fmt.Errorf("mail: SMTP_HOST not configured")
	}
	s.Log.InfoContext(ctx, "mail (dev, not sent)", "to", to, "subject", subject, "body", body)
	return nil
}

// BrevoSender sends through Brevo's HTTPS API (https://api.brevo.com), for
// hosts that block outbound SMTP ports (Render's free plan does).
type BrevoSender struct {
	APIKey, From, FromName string
	Endpoint               string       // default brevoEndpoint (tests override)
	Client                 *http.Client // default: 15 s timeout
}

const brevoEndpoint = "https://api.brevo.com/v3/smtp/email"

func (s BrevoSender) Send(ctx context.Context, to, subject, body string) error {
	if strings.ContainsAny(to, "\r\n") {
		return fmt.Errorf("mail: invalid recipient")
	}
	type addr struct {
		Email string `json:"email"`
		Name  string `json:"name,omitempty"`
	}
	payload, err := json.Marshal(struct {
		Sender      addr   `json:"sender"`
		To          []addr `json:"to"`
		Subject     string `json:"subject"`
		TextContent string `json:"textContent"`
	}{addr{s.From, s.FromName}, []addr{{Email: to}}, subject, body})
	if err != nil {
		return err
	}
	endpoint := s.Endpoint
	if endpoint == "" {
		endpoint = brevoEndpoint
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(payload))
	if err != nil {
		return err
	}
	req.Header.Set("api-key", s.APIKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	client := s.Client
	if client == nil {
		client = &http.Client{Timeout: 15 * time.Second}
	}
	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("mail: brevo: %w", err)
	}
	defer res.Body.Close()
	if res.StatusCode/100 != 2 {
		// The error body names the problem (bad key, unverified sender); it
		// never echoes the key or the message.
		msg, _ := io.ReadAll(io.LimitReader(res.Body, 512))
		return fmt.Errorf("mail: brevo: status %d: %s", res.StatusCode, strings.TrimSpace(string(msg)))
	}
	return nil
}

// SMTPSender sends through an authenticated SMTP relay.
type SMTPSender struct {
	Host, Username, Password, From string
	Port                           int
	Timeout                        time.Duration
}

func (s SMTPSender) Send(ctx context.Context, to, subject, body string) error {
	if strings.ContainsAny(to, "\r\n") {
		return fmt.Errorf("mail: invalid recipient") // header injection guard
	}
	addr := net.JoinHostPort(s.Host, fmt.Sprint(s.Port))
	timeout := s.Timeout
	if timeout <= 0 {
		timeout = 15 * time.Second
	}
	dialer := &net.Dialer{Timeout: timeout}
	tlsCfg := &tls.Config{ServerName: s.Host, MinVersion: tls.VersionTLS12}

	var conn net.Conn
	var err error
	if s.Port == 465 {
		conn, err = tls.DialWithDialer(dialer, "tcp", addr, tlsCfg)
	} else {
		conn, err = dialer.DialContext(ctx, "tcp", addr)
	}
	if err != nil {
		return fmt.Errorf("mail: dial: %w", err)
	}
	_ = conn.SetDeadline(time.Now().Add(timeout))
	c, err := smtp.NewClient(conn, s.Host)
	if err != nil {
		conn.Close()
		return fmt.Errorf("mail: handshake: %w", err)
	}
	defer c.Close()
	if s.Port != 465 {
		if ok, _ := c.Extension("STARTTLS"); !ok {
			return fmt.Errorf("mail: server does not offer STARTTLS") // never send credentials in clear
		}
		if err := c.StartTLS(tlsCfg); err != nil {
			return fmt.Errorf("mail: starttls: %w", err)
		}
	}
	if s.Username != "" {
		if err := c.Auth(smtp.PlainAuth("", s.Username, s.Password, s.Host)); err != nil {
			return fmt.Errorf("mail: auth: %w", err)
		}
	}
	if err := c.Mail(s.From); err != nil {
		return fmt.Errorf("mail: from: %w", err)
	}
	if err := c.Rcpt(to); err != nil {
		return fmt.Errorf("mail: rcpt: %w", err)
	}
	w, err := c.Data()
	if err != nil {
		return fmt.Errorf("mail: data: %w", err)
	}
	if _, err := w.Write(message(s.From, to, subject, body)); err != nil {
		return fmt.Errorf("mail: write: %w", err)
	}
	if err := w.Close(); err != nil {
		return fmt.Errorf("mail: close: %w", err)
	}
	return c.Quit()
}

// message renders an RFC 5322 plain-text message (UTF-8, Q-encoded subject).
func message(from, to, subject, body string) []byte {
	id := make([]byte, 12)
	_, _ = rand.Read(id)
	domain := from[strings.LastIndex(from, "@")+1:]
	var b strings.Builder
	fmt.Fprintf(&b, "From: %s\r\nTo: %s\r\nSubject: %s\r\nDate: %s\r\nMessage-ID: <%s@%s>\r\n",
		from, to, mime.QEncoding.Encode("utf-8", subject), time.Now().UTC().Format(time.RFC1123Z), hex.EncodeToString(id), domain)
	b.WriteString("MIME-Version: 1.0\r\nContent-Type: text/plain; charset=utf-8\r\nContent-Transfer-Encoding: 8bit\r\n\r\n")
	b.WriteString(strings.ReplaceAll(body, "\n", "\r\n"))
	return []byte(b.String())
}

// Dispatch sends up to `limit` pending outbox rows. The worker runs it under a
// per-job advisory lock, so no row-level locking is needed. Sent and abandoned
// rows have their body cleared (one-time links must not linger at rest).
func Dispatch(ctx context.Context, pool *pgxpool.Pool, sender Sender, limit int) (int64, error) {
	rows, err := pool.Query(ctx, `
		SELECT id, to_address, subject, body, attempts, created_at
		  FROM ops.email_outbox
		 WHERE sent_at IS NULL AND failed_at IS NULL
		 ORDER BY id LIMIT $1`, limit)
	if err != nil {
		return 0, err
	}
	type row struct {
		id        int64
		to, subj  string
		body      *string
		attempts  int
		createdAt time.Time
	}
	var batch []row
	for rows.Next() {
		var r row
		if err := rows.Scan(&r.id, &r.to, &r.subj, &r.body, &r.attempts, &r.createdAt); err != nil {
			rows.Close()
			return 0, err
		}
		batch = append(batch, r)
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return 0, err
	}

	var sent int64
	for _, r := range batch {
		if r.body == nil || time.Since(r.createdAt) > giveUpAfter {
			_, err = pool.Exec(ctx, `UPDATE ops.email_outbox SET failed_at = now(), body = NULL, last_error = 'expired' WHERE id = $1`, r.id)
		} else if serr := sender.Send(ctx, r.to, r.subj, *r.body); serr == nil {
			sent++
			_, err = pool.Exec(ctx, `UPDATE ops.email_outbox SET sent_at = now(), body = NULL, attempts = attempts + 1 WHERE id = $1`, r.id)
		} else {
			_, err = pool.Exec(ctx, `
				UPDATE ops.email_outbox
				   SET attempts = attempts + 1, last_error = left($2, 500),
				       failed_at = CASE WHEN attempts + 1 >= $3 THEN now() END,
				       body      = CASE WHEN attempts + 1 >= $3 THEN NULL ELSE body END
				 WHERE id = $1`, r.id, serr.Error(), MaxAttempts)
		}
		if err != nil {
			return sent, err
		}
	}
	return sent, nil
}
