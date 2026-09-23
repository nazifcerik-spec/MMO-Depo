package vendor

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/otter/server/internal/domain/character"
	"github.com/otter/server/internal/httpctx"
)

type denyOwner struct{}

func (denyOwner) Owns(context.Context, uuid.UUID, uuid.UUID) error { return character.ErrForbidden }

// Vendor routes take character_id in the body; spending another player's
// gold or selling their items must be impossible. The service is nil: any
// call past the ownership check would panic the test.
func TestVendorRejectsForeignCharacter(t *testing.T) {
	h := NewHandler(nil, denyOwner{}, nil)
	r := chi.NewRouter()
	r.Post("/v1/vendor/{code}/buy", h.Buy)
	r.Post("/v1/vendor/{code}/sell", h.Sell)
	r.Post("/v1/vendor/{code}/buyback/{buyback_id}", h.Buyback)

	victim := uuid.NewString()
	for path, body := range map[string]string{
		"/v1/vendor/general_ironhold/buy":                         `{"template_code":"potion_hp_small","character_id":"` + victim + `"}`,
		"/v1/vendor/general_ironhold/sell":                        `{"item_instance_id":"` + uuid.NewString() + `","character_id":"` + victim + `"}`,
		"/v1/vendor/general_ironhold/buyback/" + uuid.NewString(): `{"character_id":"` + victim + `"}`,
	} {
		req := httptest.NewRequest(http.MethodPost, path, strings.NewReader(body))
		req = req.WithContext(httpctx.WithAccountID(req.Context(), uuid.New()))
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)
		require.Equal(t, http.StatusNotFound, rec.Code, path)
		require.Contains(t, rec.Body.String(), "E_NOT_FOUND", path)
	}
}
