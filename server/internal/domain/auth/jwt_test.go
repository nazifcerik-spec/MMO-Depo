package auth_test

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/otter/server/internal/domain/auth"
)

func newTestSigner(t *testing.T, ttl time.Duration) *auth.Signer {
	t.Helper()
	seed, err := auth.GenerateSeed()
	require.NoError(t, err)
	s, err := auth.NewSigner(seed, "otter-test", ttl)
	require.NoError(t, err)
	return s
}

func TestSignAndVerify_Roundtrip(t *testing.T) {
	signer := newTestSigner(t, 15*time.Minute)
	acct := uuid.New()
	sess := uuid.New()
	now := time.Now()

	tok, exp, err := signer.SignAccessToken(acct, sess, now)
	require.NoError(t, err)
	require.NotEmpty(t, tok)
	require.True(t, exp.After(now))

	claims, err := signer.VerifyAccessToken(tok)
	require.NoError(t, err)
	require.Equal(t, acct, claims.AccountID)
	require.Equal(t, sess, claims.SessionID)
}

func TestVerify_ExpiredToken(t *testing.T) {
	signer := newTestSigner(t, 1*time.Millisecond)
	tok, _, err := signer.SignAccessToken(uuid.New(), uuid.New(), time.Now().Add(-1*time.Hour))
	require.NoError(t, err)

	// Give it a moment to be well past expiration.
	time.Sleep(5 * time.Millisecond)
	_, err = signer.VerifyAccessToken(tok)
	require.ErrorIs(t, err, auth.ErrTokenExpired)
}

func TestVerify_WrongIssuer(t *testing.T) {
	seed, _ := auth.GenerateSeed()
	a, _ := auth.NewSigner(seed, "otter-a", time.Minute)
	b, _ := auth.NewSigner(seed, "otter-b", time.Minute)

	tok, _, err := a.SignAccessToken(uuid.New(), uuid.New(), time.Now())
	require.NoError(t, err)

	_, err = b.VerifyAccessToken(tok)
	require.ErrorIs(t, err, auth.ErrTokenInvalid)
}

func TestVerify_WrongKey(t *testing.T) {
	seed1, _ := auth.GenerateSeed()
	seed2, _ := auth.GenerateSeed()
	a, _ := auth.NewSigner(seed1, "otter", time.Minute)
	b, _ := auth.NewSigner(seed2, "otter", time.Minute)

	tok, _, err := a.SignAccessToken(uuid.New(), uuid.New(), time.Now())
	require.NoError(t, err)

	_, err = b.VerifyAccessToken(tok)
	require.ErrorIs(t, err, auth.ErrTokenInvalid)
}

func TestVerify_TamperedToken(t *testing.T) {
	signer := newTestSigner(t, time.Minute)
	tok, _, err := signer.SignAccessToken(uuid.New(), uuid.New(), time.Now())
	require.NoError(t, err)

	// Flip a byte in the payload segment
	tampered := tok[:len(tok)-4] + "AAAA"
	_, err = signer.VerifyAccessToken(tampered)
	require.ErrorIs(t, err, auth.ErrTokenInvalid)
}

func TestVerify_RejectsAlgNone(t *testing.T) {
	signer := newTestSigner(t, time.Minute)
	// Craft a "none" token; verifier must reject on algorithm mismatch.
	claims := auth.AccessClaims{
		AccountID: uuid.New(),
		SessionID: uuid.New(),
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "otter-test",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
	}
	tok := jwt.NewWithClaims(jwt.SigningMethodNone, claims)
	raw, err := tok.SignedString(jwt.UnsafeAllowNoneSignatureType)
	require.NoError(t, err)

	_, err = signer.VerifyAccessToken(raw)
	require.ErrorIs(t, err, auth.ErrTokenInvalid)
}

func TestEncodeDecodeSeed_Roundtrip(t *testing.T) {
	s, _ := auth.GenerateSeed()
	encoded := auth.EncodeSeed(s)
	decoded, err := auth.DecodeSeed(encoded)
	require.NoError(t, err)
	require.Equal(t, s, decoded)
}

func TestNewSigner_RejectsShortSeed(t *testing.T) {
	_, err := auth.NewSigner(make([]byte, 16), "otter", time.Minute)
	require.Error(t, err)
}
