// Package apierr is the single point that formats HTTP JSON responses so every
// handler stays consistent — status + body + Content-Type + correlation_id.
package apierr

import (
	"encoding/json"
	"log/slog"
	"net/http"

	chimw "github.com/go-chi/chi/v5/middleware"
)

// WriteJSON emits a success response.
func WriteJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}

// Error is the JSON error envelope every failed handler returns.
type Error struct {
	Code          string `json:"code"`
	Message       string `json:"message,omitempty"`
	Field         string `json:"field,omitempty"`
	CorrelationID string `json:"correlation_id,omitempty"`
}

// Write emits a structured error. Code is the machine-readable identifier
// (E_VALIDATION, E_IDEMPOTENT_MISMATCH…); message is a human hint.
func Write(w http.ResponseWriter, r *http.Request, status int, code, message string) {
	if status >= http.StatusInternalServerError {
		// 5xx detail (SQL, driver, stack context) is for operators, not players:
		// log it with the correlation id and send a generic message.
		slog.ErrorContext(r.Context(), "http 5xx",
			slog.String("correlation_id", chimw.GetReqID(r.Context())),
			slog.String("code", code), slog.String("detail", message))
		message = "internal error"
	}
	WriteJSON(w, status, map[string]Error{
		"error": {
			Code:          code,
			Message:       message,
			CorrelationID: chimw.GetReqID(r.Context()),
		},
	})
}

// WriteField is Write plus a validation field name.
func WriteField(w http.ResponseWriter, r *http.Request, status int, code, message, field string) {
	WriteJSON(w, status, map[string]Error{
		"error": {
			Code:          code,
			Message:       message,
			Field:         field,
			CorrelationID: chimw.GetReqID(r.Context()),
		},
	})
}
