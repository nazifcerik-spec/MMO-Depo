package auth_test

import (
	"testing"

	"github.com/otter/server/internal/domain/auth"
)

// go test -bench KDF -run ^$ ./internal/domain/auth/ — per-hash cost of the
// production parameters (DD-238 capacity math: logins/s ≈ slots / cost).
func BenchmarkKDFDefault(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if _, err := auth.HashPassword("benchmark-password", auth.DefaultArgon2Params); err != nil {
			b.Fatal(err)
		}
	}
}
