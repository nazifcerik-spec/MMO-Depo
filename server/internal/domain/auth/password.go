package auth

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"runtime"
	"strings"
	"time"

	"golang.org/x/crypto/argon2"
)

// Argon2Params encapsulates the tuning knobs. Defaults follow Faz 00 baseline;
// production benchmarks per DD-189 tune memory/time to keep verify < 100ms p95.
type Argon2Params struct {
	Memory      uint32 // KiB
	Iterations  uint32
	Parallelism uint8
	SaltLength  uint32
	KeyLength   uint32
}

// DefaultArgon2Params — OWASP's argon2id baseline: 19 MiB, 2 iterations,
// single lane, 16-byte salt, 32-byte key. Benchmarked in DD-238: the earlier
// 64 MiB × 3 cost ~140 ms per login and missed the login SLO under load.
// Existing hashes keep verifying with the parameters encoded in them.
var DefaultArgon2Params = Argon2Params{
	Memory:      19 * 1024,
	Iterations:  2,
	Parallelism: 1,
	SaltLength:  16,
	KeyLength:   32,
}

// ErrBusy: every key-derivation slot stayed taken for kdfWait (login storm).
// The short wait keeps a successful login inside the burst SLO (p95 < 800 ms:
// wait + ~140 ms hash) and sheds the rest fast.
// Callers answer 503 + Retry-After instead of queueing unbounded 64 MiB
// allocations (DD-238, docs/performance/load-login.md §6).
var ErrBusy = errors.New("auth: password hashing busy")

var (
	kdfSlots = make(chan struct{}, runtime.NumCPU())
	kdfWait  = 500 * time.Millisecond
)

// SetKDFConcurrency bounds concurrent argon2id runs, so peak memory is
// n × Memory. Call once at startup; n <= 0 keeps the default (CPU count).
func SetKDFConcurrency(n int) {
	if n > 0 {
		kdfSlots = make(chan struct{}, n)
	}
}

func deriveKey(password string, salt []byte, p Argon2Params) ([]byte, error) {
	slots := kdfSlots
	t := time.NewTimer(kdfWait)
	defer t.Stop()
	select {
	case slots <- struct{}{}:
	case <-t.C:
		return nil, ErrBusy
	}
	defer func() { <-slots }()
	return argon2.IDKey([]byte(password), salt, p.Iterations, p.Memory, p.Parallelism, p.KeyLength), nil
}

// dummyHashFor encodes a hash with p's cost whose password is unknowable
// (all-zero key): verifying against it always fails after the full work.
func dummyHashFor(p Argon2Params) string {
	return fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s", argon2.Version, p.Memory, p.Iterations, p.Parallelism,
		base64.RawStdEncoding.EncodeToString(make([]byte, 16)), base64.RawStdEncoding.EncodeToString(make([]byte, 32)))
}

// HashPassword returns the standard Argon2id encoded string, which embeds the
// algorithm, version, parameters, salt, and derived key. This is the format
// downstream libraries (including Passlib) can round-trip.
//
//	$argon2id$v=19$m=19456,t=2,p=1$<salt-b64>$<hash-b64>
func HashPassword(password string, p Argon2Params) (string, error) {
	if password == "" {
		return "", errors.New("password empty")
	}
	salt := make([]byte, p.SaltLength)
	if _, err := rand.Read(salt); err != nil {
		return "", fmt.Errorf("rand: %w", err)
	}

	hash, err := deriveKey(password, salt, p)
	if err != nil {
		return "", err
	}

	encoded := fmt.Sprintf(
		"$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version,
		p.Memory,
		p.Iterations,
		p.Parallelism,
		base64.RawStdEncoding.EncodeToString(salt),
		base64.RawStdEncoding.EncodeToString(hash),
	)
	return encoded, nil
}

// VerifyPassword returns nil when password matches the encoded hash, or a
// non-nil error (invalid encoding, wrong password, ErrBusy). It is constant-time
// comparing the derived key to defeat timing side-channels.
func VerifyPassword(password, encoded string) error {
	p, salt, want, err := decodeArgon2Encoded(encoded)
	if err != nil {
		return err
	}
	got, err := deriveKey(password, salt, p)
	if err != nil {
		return err
	}
	if subtle.ConstantTimeCompare(got, want) != 1 {
		return ErrInvalidCredentials
	}
	return nil
}

func decodeArgon2Encoded(encoded string) (Argon2Params, []byte, []byte, error) {
	// $argon2id$v=19$m=65536,t=3,p=1$<salt>$<hash>
	parts := strings.Split(encoded, "$")
	if len(parts) != 6 || parts[0] != "" || parts[1] != "argon2id" {
		return Argon2Params{}, nil, nil, errors.New("invalid argon2id encoding")
	}
	var version int
	if _, err := fmt.Sscanf(parts[2], "v=%d", &version); err != nil {
		return Argon2Params{}, nil, nil, fmt.Errorf("invalid version: %w", err)
	}
	if version != argon2.Version {
		return Argon2Params{}, nil, nil, fmt.Errorf("unsupported argon2 version %d", version)
	}
	var p Argon2Params
	if _, err := fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &p.Memory, &p.Iterations, &p.Parallelism); err != nil {
		return Argon2Params{}, nil, nil, fmt.Errorf("invalid params: %w", err)
	}
	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return Argon2Params{}, nil, nil, fmt.Errorf("invalid salt: %w", err)
	}
	hash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return Argon2Params{}, nil, nil, fmt.Errorf("invalid hash: %w", err)
	}
	p.SaltLength = uint32(len(salt))
	p.KeyLength = uint32(len(hash))
	return p, salt, hash, nil
}
