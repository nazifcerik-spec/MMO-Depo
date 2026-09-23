package hunt

import (
	"encoding/json"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

// The seed + frozen snapshot fully determine the outcome; a client holding
// them could simulate the hunt and reroll bad ones. They must never be sent.
func TestContractJSON_HidesDeterminismInputs(t *testing.T) {
	b, err := json.Marshal(Contract{ID: uuid.New(), Seed: 424242, PlayerSnapshot: json.RawMessage(`{"attack":9}`)})
	require.NoError(t, err)
	s := string(b)
	require.NotContains(t, s, "424242")
	require.NotContains(t, s, "attack")
	require.Contains(t, s, `"id":`)
	require.Contains(t, s, `"mob_code":`)
}
