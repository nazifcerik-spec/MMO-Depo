package privacy

import (
	"context"
	"encoding/json"
	"errors"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/otter/server/internal/domain/auth"
)

// TokenRedeemer turns a one-time export link token into an account id.
type TokenRedeemer interface {
	ConsumeExportToken(ctx context.Context, token string) (uuid.UUID, error)
}

// Legal identifies the operator on the policy pages; empty values render a
// draft banner (the texts are drafts until legally reviewed).
type Legal struct {
	Operator     string
	SupportEmail string
	TermsVersion string
}

// Handler serves the export download and the legal pages.
type Handler struct {
	pool       *pgxpool.Pool
	tokens     TokenRedeemer
	contentDir string
	legal      Legal
	log        *slog.Logger
	now        func() time.Time
}

func NewHandler(pool *pgxpool.Pool, tokens TokenRedeemer, contentDir string, legal Legal, logger *slog.Logger) *Handler {
	if logger == nil {
		logger = slog.Default()
	}
	return &Handler{pool: pool, tokens: tokens, contentDir: contentDir, legal: legal, log: logger, now: time.Now}
}

// Export — GET /v1/me/export?token=… (opened in the browser from the app):
// the account's data as a JSON download. The token is single-use.
func (h *Handler) Export(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Referrer-Policy", "no-referrer") // the token is in the URL
	acct, err := h.tokens.ConsumeExportToken(r.Context(), r.URL.Query().Get("token"))
	if err != nil {
		status, msg := http.StatusBadRequest,
			"Bağlantı geçersiz ya da süresi doldu; uygulamadan yeni bağlantı isteyin.\n"+
				"The link is invalid or expired; request a new one in the app.\n"
		if !errors.Is(err, auth.ErrRecoveryTokenInvalid) {
			h.log.ErrorContext(r.Context(), "privacy: export token", "err", err.Error())
			status, msg = http.StatusInternalServerError, "Hata / Error\n"
		}
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(status)
		_, _ = w.Write([]byte(msg))
		return
	}
	doc, err := Export(r.Context(), h.pool, acct, h.now())
	if err != nil {
		h.log.ErrorContext(r.Context(), "privacy: export", "err", err.Error())
		http.Error(w, "Hata / Error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Content-Disposition", `attachment; filename="new-world-mmo-data.json"`)
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	_ = enc.Encode(doc)
}

// Terms — GET /terms; Privacy — GET /privacy (?lang=tr|en, else
// Accept-Language). Texts live in the content pack: legal/<doc>.<lang>.html.
func (h *Handler) Terms(w http.ResponseWriter, r *http.Request)   { h.page(w, r, "terms") }
func (h *Handler) Privacy(w http.ResponseWriter, r *http.Request) { h.page(w, r, "privacy") }

var legalTitles = map[string]map[string]string{
	"terms":   {"tr": "Kullanım Şartları", "en": "Terms of Service"},
	"privacy": {"tr": "Gizlilik Politikası", "en": "Privacy Policy"},
}

func pageLang(r *http.Request) string {
	if l := r.URL.Query().Get("lang"); l == "tr" || l == "en" {
		return l
	}
	if strings.HasPrefix(strings.ToLower(r.Header.Get("Accept-Language")), "tr") {
		return "tr"
	}
	return "en"
}

func (h *Handler) page(w http.ResponseWriter, r *http.Request, doc string) {
	lang := pageLang(r)
	// doc and lang come from fixed sets: no user input reaches the path.
	raw, err := os.ReadFile(filepath.Join(h.contentDir, "legal", doc+"."+lang+".html"))
	if err != nil {
		h.log.ErrorContext(r.Context(), "privacy: legal page", "doc", doc, "err", err.Error())
		http.Error(w, "Not available", http.StatusServiceUnavailable)
		return
	}
	operator, email := h.legal.Operator, h.legal.SupportEmail
	draft := operator == "" || email == ""
	if operator == "" {
		operator = "[OPERATOR]"
	}
	if email == "" {
		email = "[SUPPORT EMAIL]"
	}
	body := strings.NewReplacer("{{operator}}", template.HTMLEscapeString(operator),
		"{{support_email}}", template.HTMLEscapeString(email),
		"{{version}}", template.HTMLEscapeString(h.legal.TermsVersion)).Replace(string(raw))
	w.Header().Set("Content-Security-Policy", "default-src 'none'; style-src 'unsafe-inline'; frame-ancestors 'none'")
	w.Header().Set("Cache-Control", "public, max-age=3600")
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_ = legalTmpl.Execute(w, map[string]any{
		"Lang": lang, "Title": legalTitles[doc][lang], "Body": template.HTML(body), "Draft": draft, "Doc": doc,
	})
}

var legalTmpl = template.Must(template.New("legal").Parse(`<!doctype html>
<html lang="{{.Lang}}"><head><meta charset="utf-8"><meta name="viewport" content="width=device-width,initial-scale=1">
<title>{{.Title}} — New World MMO</title><style>
body{margin:0;background:#0F1013;color:#E8EAEE;font:15px/1.6 system-ui,sans-serif}
main{max-width:760px;margin:0 auto;padding:24px 16px 64px}h1{font-size:24px}h2{font-size:18px;margin-top:28px}
a{color:#E4A756}.draft{background:#3A2A12;border:1px solid #E4A756;border-radius:8px;padding:10px 12px;margin:12px 0}
nav{display:flex;gap:16px;font-size:14px}</style></head><body><main>
<nav><a href="/terms?lang={{.Lang}}">{{if eq .Lang "tr"}}Kullanım Şartları{{else}}Terms{{end}}</a>
<a href="/privacy?lang={{.Lang}}">{{if eq .Lang "tr"}}Gizlilik{{else}}Privacy{{end}}</a>
<a href="/{{.Doc}}?lang={{if eq .Lang "tr"}}en{{else}}tr{{end}}">{{if eq .Lang "tr"}}English{{else}}Türkçe{{end}}</a></nav>
{{if .Draft}}<div class="draft">{{if eq .Lang "tr"}}Taslak: işletmeci bilgileri girilmedi, metin hukuki inceleme bekliyor.{{else}}Draft: operator details not set; awaiting legal review.{{end}}</div>{{end}}
<h1>{{.Title}}</h1>
{{.Body}}
</main></body></html>`))
