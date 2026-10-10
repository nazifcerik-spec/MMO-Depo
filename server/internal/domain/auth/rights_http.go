package auth

import (
	"net/http"

	"github.com/google/uuid"

	"github.com/otter/server/internal/http/apierr"
	"github.com/otter/server/internal/httpctx"
)

// DeleteAccount — DELETE /v1/me {current_password} → 202 {purge_after}.
// Signs every device out; signing in before purge_after cancels (DD-231).
func (h *Handler) DeleteAccount(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Current string `json:"current_password"`
	}
	if err := decodeJSON(r, &body); err != nil {
		apierr.Write(w, r, http.StatusBadRequest, "E_INVALID_JSON", err.Error())
		return
	}
	acct := httpctx.AccountID(r.Context())
	if acct == uuid.Nil {
		apierr.Write(w, r, http.StatusUnauthorized, "E_UNAUTHENTICATED", "auth required")
		return
	}
	at, err := h.svc.RequestDeletion(r.Context(), acct, body.Current)
	if err != nil {
		writeAuthError(w, r, err)
		return
	}
	apierr.WriteJSON(w, http.StatusAccepted, map[string]any{"purge_after": at})
}

// ExportLink — POST /v1/me/export-link → 201 {url, expires_at}: a one-time
// link the app opens in the browser to download the account's data.
func (h *Handler) ExportLink(w http.ResponseWriter, r *http.Request) {
	acct := httpctx.AccountID(r.Context())
	if acct == uuid.Nil {
		apierr.Write(w, r, http.StatusUnauthorized, "E_UNAUTHENTICATED", "auth required")
		return
	}
	url, exp, err := h.svc.IssueExportLink(r.Context(), acct)
	if err != nil {
		writeAuthError(w, r, err)
		return
	}
	apierr.WriteJSON(w, http.StatusCreated, map[string]any{"url": url, "expires_at": exp})
}
