package social

import (
	"bytes"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestPostage_MatchesDocExample(t *testing.T) {
	// docs/social/mail.md §4: 5 items + 100g → 5 + 5 + 1000 = 1010c.
	require.Equal(t, int64(1010), Postage(5, 1_000_000))
	require.Equal(t, int64(5), Postage(0, 0))
	require.Equal(t, int64(6), Postage(0, 500), "0.5c rounds half up")
}

func TestCleanText(t *testing.T) {
	s, err := CleanText("  hi  ", 10)
	require.NoError(t, err)
	require.Equal(t, "hi", s)
	_, err = CleanText("   ", 10)
	require.ErrorIs(t, err, ErrEmpty)
	_, err = CleanText("çok uzun metin", 5)
	require.ErrorIs(t, err, ErrTooLong)
	_, err = CleanText("ğüşiöç", 6)
	require.NoError(t, err, "limit counts runes, not bytes")
}

func TestPair_CanonicalByteOrder(t *testing.T) {
	for i := 0; i < 50; i++ {
		x, y := uuid.New(), uuid.New()
		a, b := pair(x, y)
		a2, b2 := pair(y, x)
		require.Equal(t, a, a2)
		require.Equal(t, b, b2)
		require.Negative(t, bytes.Compare(a[:], b[:]))
	}
}
