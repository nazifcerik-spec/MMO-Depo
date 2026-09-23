package combat

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMobsHandler_SortedPublicCatalog(t *testing.T) {
	svc := NewService(NewCatalog([]Mob{
		{Code: "wolf", Level: 5, HPMax: 80},
		{Code: "rat", Level: 1, HPMax: 20},
		{Code: "bat", Level: 1, HPMax: 15},
	}), nil)
	rec := httptest.NewRecorder()
	NewHandler(svc, nil, nil).Mobs(rec, httptest.NewRequest("GET", "/v1/mobs", nil))
	require.Equal(t, 200, rec.Code)
	var body struct {
		Mobs []struct {
			Code  string `json:"code"`
			Level int    `json:"level"`
		} `json:"mobs"`
	}
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &body))
	require.Len(t, body.Mobs, 3)
	require.Equal(t, []string{"bat", "rat", "wolf"}, []string{body.Mobs[0].Code, body.Mobs[1].Code, body.Mobs[2].Code})
}
