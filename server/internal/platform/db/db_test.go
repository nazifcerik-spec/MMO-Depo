package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

// A malformed DATABASE_URL — e.g. only the password pasted in — must not
// be echoed into the error (and so into the logs).
func TestConnectNeverEchoesTheURL(t *testing.T) {
	for _, bad := range []string{"s3cretPassw0rd", "postgres://u:s3cretPassw0rd@host:notaport/db"} {
		_, err := Connect(context.Background(), bad, 1)
		require.Error(t, err)
		require.NotContains(t, err.Error(), "s3cretPassw0rd")
	}
}
