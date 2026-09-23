package dungeon

import (
	"encoding/json"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestContractJSON_HidesDeterminismInputs(t *testing.T) {
	b, err := json.Marshal(Contract{ID: uuid.New(), Seed: 424242, PlayerSnapshot: json.RawMessage(`{"attack":9}`)})
	require.NoError(t, err)
	require.NotContains(t, string(b), "424242")
	require.NotContains(t, string(b), "attack")
	require.Contains(t, string(b), `"dungeon_code":`)
}
