package action

import (
	"errors"
	"net/http"

	"github.com/otter/server/internal/http/apierr"
)

// Write answers a busy refusal: 409 E_CHARACTER_BUSY with the busy kind as
// the error field, so the client can say what to finish first. It returns
// false (writing nothing) for any other error.
func Write(w http.ResponseWriter, r *http.Request, err error) bool {
	if !errors.Is(err, ErrBusy) {
		return false
	}
	apierr.WriteField(w, r, http.StatusConflict, "E_CHARACTER_BUSY", "finish the current action first", Kind(err))
	return true
}
