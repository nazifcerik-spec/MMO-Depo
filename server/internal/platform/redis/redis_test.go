package redis

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOptions(t *testing.T) {
	o, err := Options("localhost:6379", "pw", 2)
	require.NoError(t, err)
	require.Equal(t, "localhost:6379", o.Addr)
	require.Equal(t, "pw", o.Password)
	require.Equal(t, 2, o.DB)
	require.Equal(t, 50, o.PoolSize)

	// Render Key Value hands out a connection string.
	o, err = Options("redis://red-abc123:6379", "", 0)
	require.NoError(t, err)
	require.Equal(t, "red-abc123:6379", o.Addr)
	o, err = Options("rediss://user:secret@cache.example.com:6380/3", "", 0)
	require.NoError(t, err)
	require.Equal(t, "cache.example.com:6380", o.Addr)
	require.Equal(t, "secret", o.Password)
	require.Equal(t, 3, o.DB)
	require.NotNil(t, o.TLSConfig)
	require.Equal(t, 50, o.PoolSize)

	_, err = Options("redis://bad host:x", "", 0)
	require.Error(t, err)
}
