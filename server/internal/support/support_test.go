package support

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestClean(t *testing.T) {
	ok := NewTicket{Category: "bug", Subject: "  A bug  ", Body: "  Something broke here.  ",
		ClientInfo: map[string]string{"platform": "android", "evil": "x", "locale": strings.Repeat("y", 100)}}
	got, err := clean(ok)
	require.NoError(t, err)
	require.Equal(t, "A bug", got.Subject)
	require.Equal(t, "Something broke here.", got.Body)
	require.Equal(t, map[string]string{"platform": "android", "locale": strings.Repeat("y", 64)}, got.ClientInfo)

	for name, mut := range map[string]func(*NewTicket){
		"category":      func(t *NewTicket) { t.Category = "spam" },
		"short subject": func(t *NewTicket) { t.Subject = "ab" },
		"long subject":  func(t *NewTicket) { t.Subject = strings.Repeat("s", MaxSubject+1) },
		"short body":    func(t *NewTicket) { t.Body = "too short" },
		"long body":     func(t *NewTicket) { t.Body = strings.Repeat("b", MaxBody+1) },
	} {
		bad := ok
		mut(&bad)
		_, err := clean(bad)
		require.ErrorIs(t, err, ErrInvalid, name)
	}
}

// Every support error reaches the client as a 4xx.
func TestWriteErr(t *testing.T) {
	for _, e := range []error{ErrInvalid, ErrReplyRequired, ErrTooManyOpen, ErrForeignCharacter, ErrNotFound, ErrClosed} {
		rec := httptest.NewRecorder()
		WriteErr(rec, httptest.NewRequest(http.MethodGet, "/", nil), fmt.Errorf("ctx: %w", e))
		require.Less(t, rec.Code, 500, e.Error())
	}
}
