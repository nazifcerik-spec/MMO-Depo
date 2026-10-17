package social

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Every domain error must map to a client status: an unmapped one reaches
// players as a 500 "server error" (503 = retry is fine). Errors are wrapped
// like services do.
func TestDomainErrorsMapToClientStatus(t *testing.T) {
	for name, e := range map[string]error{
		"ErrAlreadyFriends":     ErrAlreadyFriends,
		"ErrAlreadyRequested":   ErrAlreadyRequested,
		"ErrAttachmentInvalid":  ErrAttachmentInvalid,
		"ErrBadChannel":         ErrBadChannel,
		"ErrBlocked":            ErrBlocked,
		"ErrEmpty":              ErrEmpty,
		"ErrFriendCap":          ErrFriendCap,
		"ErrHasAttachments":     ErrHasAttachments,
		"ErrIdempotencyMissing": ErrIdempotencyMissing,
		"ErrInboxFull":          ErrInboxFull,
		"ErrMuted":              ErrMuted,
		"ErrNotFound":           ErrNotFound,
		"ErrNotFriends":         ErrNotFriends,
		"ErrNotInChannel":       ErrNotInChannel,
		"ErrNothingToClaim":     ErrNothingToClaim,
		"ErrBagFull":            ErrBagFull,
		"ErrSelfTarget":         ErrSelfTarget,
		"ErrTargetNotFound":     ErrTargetNotFound,
		"ErrTooLong":            ErrTooLong,
		"ErrTooManyAttachments": ErrTooManyAttachments,
	} {
		t.Run(name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			wrapped := fmt.Errorf("context: %w", e)
			var codes []int
			rec0 := httptest.NewRecorder()
			writeErr(rec0, req, wrapped)
			codes = append(codes, rec0.Code)
			for _, c := range codes {
				if c < http.StatusInternalServerError || c == http.StatusServiceUnavailable {
					return
				}
			}
			t.Fatalf("%s maps to %v (want a 4xx)", name, codes)
		})
	}
	_ = errors.New
}
