package auth_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/otter/server/internal/domain/auth"
)

// Reduced params for fast tests; production uses DefaultArgon2Params.
var fastParams = auth.Argon2Params{
	Memory:      8 * 1024, // 8 MiB
	Iterations:  1,
	Parallelism: 1,
	SaltLength:  16,
	KeyLength:   32,
}

func TestHashPasswordEncoding(t *testing.T) {
	h, err := auth.HashPassword("hunter22", fastParams)
	require.NoError(t, err)
	require.True(t, strings.HasPrefix(h, "$argon2id$v=19$"))
	// $argon2id$v=19$m=8192,t=1,p=1$<salt>$<hash>
	parts := strings.Split(h, "$")
	require.Len(t, parts, 6)
	require.Equal(t, "argon2id", parts[1])
}

func TestHashPassword_UniqueSaltEachCall(t *testing.T) {
	h1, err := auth.HashPassword("hunter22", fastParams)
	require.NoError(t, err)
	h2, err := auth.HashPassword("hunter22", fastParams)
	require.NoError(t, err)
	require.NotEqual(t, h1, h2, "salt must be random each call")
}

func TestVerifyPassword_Success(t *testing.T) {
	h, err := auth.HashPassword("MyP@ssw0rd", fastParams)
	require.NoError(t, err)
	require.NoError(t, auth.VerifyPassword("MyP@ssw0rd", h))
}

func TestVerifyPassword_Wrong(t *testing.T) {
	h, err := auth.HashPassword("MyP@ssw0rd", fastParams)
	require.NoError(t, err)
	err = auth.VerifyPassword("wrong", h)
	require.ErrorIs(t, err, auth.ErrInvalidCredentials)
}

func TestVerifyPassword_InvalidEncoding(t *testing.T) {
	err := auth.VerifyPassword("x", "not-argon2-string")
	require.Error(t, err)
	require.NotErrorIs(t, err, auth.ErrInvalidCredentials, "encoding error must not be swallowed as bad password")
}

func TestHashPassword_EmptyRejected(t *testing.T) {
	_, err := auth.HashPassword("", fastParams)
	require.Error(t, err)
}
