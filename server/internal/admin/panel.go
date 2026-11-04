package admin

import (
	"context"
	"embed"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/otter/server/internal/http/apierr"
)

// The admin panel (DD-360): three static files that carry no data and call the
// same JSON routes as everything else. Sign-in, permissions and audit stay on
// those routes.
//
//go:embed panel/panel.html panel/panel.js panel/panel.css
var panelFiles embed.FS

const panelCSP = "default-src 'none'; script-src 'self'; style-src 'self'; connect-src 'self'; " +
	"base-uri 'none'; form-action 'self'; frame-ancestors 'none'"

func servePanelFile(name, contentType string) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		raw, err := panelFiles.ReadFile("panel/" + name)
		if err != nil {
			http.NotFound(w, nil)
			return
		}
		h := w.Header()
		h.Set("Content-Type", contentType)
		h.Set("Content-Security-Policy", panelCSP)
		h.Set("X-Content-Type-Options", "nosniff")
		h.Set("X-Frame-Options", "DENY")
		h.Set("Referrer-Policy", "no-referrer")
		h.Set("Cache-Control", "no-store")
		_, _ = w.Write(raw)
	}
}

// PanelRoutes serves the panel's files. They are public: nothing in them is
// secret, and every call they make needs an admin session.
func PanelRoutes(r chi.Router) {
	r.Get("/internal/admin/panel", servePanelFile("panel.html", "text/html; charset=utf-8"))
	r.Get("/internal/admin/panel.js", servePanelFile("panel.js", "text/javascript; charset=utf-8"))
	r.Get("/internal/admin/panel.css", servePanelFile("panel.css", "text/css; charset=utf-8"))
}

// Overview is the panel's first page.
type Overview struct {
	Accounts        int64 `json:"accounts"`
	AccountsNew24h  int64 `json:"accounts_new_24h"`
	Logins24h       int64 `json:"logins_24h"`
	Characters      int64 `json:"characters"`
	Online          int   `json:"online"`
	OpenReports     int64 `json:"open_reports"`
	OpenTickets     int64 `json:"open_tickets"`
	ActiveListings  int64 `json:"active_listings"`
	ActiveAuctions  int64 `json:"active_auctions"`
	Guilds          int64 `json:"guilds"`
	ClientErrors24h int64 `json:"client_errors_24h"`
}

// Overview counts what an operator looks at first.
func (s *Service) Overview(ctx context.Context) (*Overview, error) {
	o := &Overview{}
	err := s.pool.QueryRow(ctx, `SELECT
		(SELECT count(*) FROM identity.account),
		(SELECT count(*) FROM identity.account WHERE created_at > now() - interval '24 hours'),
		(SELECT count(*) FROM identity.account WHERE last_login_at > now() - interval '24 hours'),
		(SELECT count(*) FROM character.character WHERE state = 'active'),
		(SELECT count(*) FROM ops.player_report WHERE state = 'open'),
		(SELECT count(*) FROM ops.support_ticket WHERE state = 'open'),
		(SELECT count(*) FROM market.listing WHERE state = 'active' AND kind = 'fixed'),
		(SELECT count(*) FROM market.listing WHERE state = 'active' AND kind = 'auction'),
		(SELECT count(*) FROM social.guild WHERE disbanded_at IS NULL),
		(SELECT COALESCE(sum(1), 0) FROM ops.client_error WHERE last_seen > now() - interval '24 hours')`).
		Scan(&o.Accounts, &o.AccountsNew24h, &o.Logins24h, &o.Characters, &o.OpenReports, &o.OpenTickets,
			&o.ActiveListings, &o.ActiveAuctions, &o.Guilds, &o.ClientErrors24h)
	return o, err
}

// OnlineCounter says how many live connections there are.
type OnlineCounter interface{ ClientCount() int }

// WithOnline shows the live connection count on the overview.
func (h *Handler) WithOnline(c OnlineCounter) *Handler { h.online = c; return h }

// Overview — GET /internal/ops/overview.
func (h *Handler) Overview(w http.ResponseWriter, r *http.Request) {
	if _, ok := h.allow(w, r, perm(r, PermPlayerRead), false); !ok {
		return
	}
	o, err := h.svc.Overview(r.Context())
	if err != nil {
		writeErr(w, r, err)
		return
	}
	if h.online != nil {
		o.Online = h.online.ClientCount()
	}
	apierr.WriteJSON(w, http.StatusOK, o)
}
