package auth

import (
	"errors"
	"html/template"
	"net/http"
	"time"

	"github.com/google/uuid"

	"github.com/otter/server/internal/http/apierr"
	"github.com/otter/server/internal/httpctx"
)

// resetRequestFloor pads reset-request responses so "account exists" (DB
// writes) and "no such account" (nothing) take the same time (password.md §6.1).
const resetRequestFloor = 400 * time.Millisecond

// RequestPasswordReset — POST /v1/auth/password/reset-request {email} → 202.
func (h *Handler) RequestPasswordReset(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	var body struct {
		Email string `json:"email"`
	}
	if err := decodeJSON(r, &body); err != nil {
		apierr.Write(w, r, http.StatusBadRequest, "E_INVALID_JSON", err.Error())
		return
	}
	if err := h.svc.RequestPasswordReset(r.Context(), body.Email); err != nil {
		h.log.ErrorContext(r.Context(), "auth: reset request failed", "err", err.Error())
	}
	if d := resetRequestFloor - time.Since(start); d > 0 {
		time.Sleep(d)
	}
	apierr.WriteJSON(w, http.StatusAccepted, map[string]string{"status": "if the account exists, an e-mail is on its way"})
}

// ConfirmPasswordReset — POST /v1/auth/password/reset-confirm {token, new_password} → 204.
func (h *Handler) ConfirmPasswordReset(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Token       string `json:"token"`
		NewPassword string `json:"new_password"`
	}
	if err := decodeJSON(r, &body); err != nil {
		apierr.Write(w, r, http.StatusBadRequest, "E_INVALID_JSON", err.Error())
		return
	}
	if err := h.svc.ConfirmPasswordReset(r.Context(), body.Token, body.NewPassword); err != nil {
		writeAuthError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// VerifyEmail — POST /v1/auth/verify-email {token} → 200.
func (h *Handler) VerifyEmail(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Token string `json:"token"`
	}
	if err := decodeJSON(r, &body); err != nil {
		apierr.Write(w, r, http.StatusBadRequest, "E_INVALID_JSON", err.Error())
		return
	}
	if err := h.svc.VerifyEmail(r.Context(), body.Token); err != nil {
		writeAuthError(w, r, err)
		return
	}
	apierr.WriteJSON(w, http.StatusOK, map[string]bool{"email_verified": true})
}

// ResendVerification — POST /v1/auth/resend-verification (auth) → 202.
func (h *Handler) ResendVerification(w http.ResponseWriter, r *http.Request) {
	if err := h.svc.ResendVerification(r.Context(), httpctx.AccountID(r.Context())); err != nil {
		writeAuthError(w, r, err)
		return
	}
	apierr.WriteJSON(w, http.StatusAccepted, map[string]string{"status": "queued"})
}

// ChangePassword — POST /v1/auth/password/change (auth) {current_password, new_password} → 204.
// Every other session is signed out; this one stays.
func (h *Handler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Current string `json:"current_password"`
		New     string `json:"new_password"`
	}
	if err := decodeJSON(r, &body); err != nil {
		apierr.Write(w, r, http.StatusBadRequest, "E_INVALID_JSON", err.Error())
		return
	}
	acct, sid := httpctx.AccountID(r.Context()), httpctx.SessionID(r.Context())
	if acct == uuid.Nil {
		apierr.Write(w, r, http.StatusUnauthorized, "E_UNAUTHENTICATED", "auth required")
		return
	}
	if err := h.svc.ChangePassword(r.Context(), acct, sid, body.Current, body.New); err != nil {
		writeAuthError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// --- E-mail link pages -----------------------------------------------------------
//
// Until the app registers deep links, the links in recovery e-mails open
// these self-contained pages (no external assets; strict CSP).

var pageTmpl = template.Must(template.New("page").Parse(`<!doctype html>
<html lang="tr"><head><meta charset="utf-8"><meta name="viewport" content="width=device-width,initial-scale=1">
<title>Otter</title><style>
body{margin:0;min-height:100vh;display:flex;align-items:center;justify-content:center;background:#0F1013;color:#E8EAEE;font:15px/1.5 system-ui,sans-serif}
main{width:min(360px,90vw);padding:24px;background:#1A1C22;border:1px solid #2E3240;border-radius:16px}
h1{font-size:20px;margin:0 0 12px}p{color:#A8ACB8}input,button{width:100%;box-sizing:border-box;padding:12px;margin:6px 0;border-radius:10px;font:inherit}
input{background:#0F1013;color:#E8EAEE;border:1px solid #2E3240}button{background:#E4A756;color:#0F1013;border:0;font-weight:700;cursor:pointer}
.ok{color:#5FB973}.err{color:#E85B62}</style></head><body><main>
{{if .Form}}<h1>Yeni şifre / New password</h1>
<form id="f"><input type="password" id="p" minlength="8" required placeholder="Yeni şifre / New password" autocomplete="new-password">
<input type="password" id="c" minlength="8" required placeholder="Tekrar / Confirm" autocomplete="new-password">
<button type="submit">Kaydet / Save</button></form><p id="m"></p>
<script nonce="{{.Nonce}}">
document.getElementById('f').onsubmit=async function(e){e.preventDefault();var m=document.getElementById('m'),p=document.getElementById('p').value;
if(p!==document.getElementById('c').value){m.className='err';m.textContent='Şifreler eşleşmiyor / Passwords do not match';return}
var r=await fetch('/v1/auth/password/reset-confirm',{method:'POST',headers:{'Content-Type':'application/json','Idempotency-Key':crypto.randomUUID()},
body:JSON.stringify({token:new URLSearchParams(location.search).get('token')||'',new_password:p})});
if(r.status===204){this.remove();m.className='ok';m.textContent='Şifren değişti, uygulamadan giriş yapabilirsin. / Password changed — log in from the app.'}
else{var j=await r.json().catch(function(){return{}});m.className='err';m.textContent=(j.error&&j.error.message)||'Hata / Error'}};
</script>{{else}}<h1>{{.Title}}</h1><p class="{{.Class}}">{{.Message}}</p>{{end}}
</main></body></html>`))

type page struct {
	Form                  bool
	Nonce                 string
	Title, Message, Class string
}

func renderPage(w http.ResponseWriter, status int, p page) {
	csp := "default-src 'none'; style-src 'unsafe-inline'; connect-src 'self'; form-action 'none'; frame-ancestors 'none'"
	if p.Nonce != "" {
		csp += "; script-src 'nonce-" + p.Nonce + "'"
	}
	w.Header().Set("Content-Security-Policy", csp)
	w.Header().Set("Referrer-Policy", "no-referrer") // the token is in the URL
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)
	_ = pageTmpl.Execute(w, p)
}

// VerifyEmailPage — GET /verify-email?token=… (link in the verification mail).
func (h *Handler) VerifyEmailPage(w http.ResponseWriter, r *http.Request) {
	err := h.svc.VerifyEmail(r.Context(), r.URL.Query().Get("token"))
	switch {
	case err == nil:
		renderPage(w, http.StatusOK, page{Title: "E-posta doğrulandı / E-mail verified",
			Message: "Uygulamaya dönebilirsin. / You can return to the app.", Class: "ok"})
	case errors.Is(err, ErrRecoveryTokenInvalid):
		renderPage(w, http.StatusBadRequest, page{Title: "Bağlantı geçersiz / Invalid link",
			Message: "Bağlantının süresi dolmuş ya da kullanılmış. Uygulamadan yeni bağlantı iste. / The link expired or was used — request a new one in the app.", Class: "err"})
	default:
		h.log.ErrorContext(r.Context(), "auth: verify page failed", "err", err.Error())
		renderPage(w, http.StatusInternalServerError, page{Title: "Hata / Error", Message: "Tekrar dene. / Please try again.", Class: "err"})
	}
}

// ResetPasswordPage — GET /reset-password?token=… serves the new-password form;
// the token is only consumed when the form posts to reset-confirm.
func (h *Handler) ResetPasswordPage(w http.ResponseWriter, _ *http.Request) {
	nonce, _, err := newRecoveryToken()
	if err != nil {
		renderPage(w, http.StatusInternalServerError, page{Title: "Hata / Error", Class: "err"})
		return
	}
	renderPage(w, http.StatusOK, page{Form: true, Nonce: nonce})
}
