package auth

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// AccessClaims are the JWT payload for the access token. Kept compact — every
// mutating request parses it.
type AccessClaims struct {
	AccountID uuid.UUID `json:"sub"`
	SessionID uuid.UUID `json:"sid"`
	jwt.RegisteredClaims
}

// GetSubject satisfies jwt.Claims when RegisteredClaims.Subject is empty.
func (c AccessClaims) GetSubject() (string, error) {
	return c.AccountID.String(), nil
}

// Signer signs and verifies EdDSA (Ed25519) JWTs. One instance per process.
type Signer struct {
	privateKey ed25519.PrivateKey
	publicKey  ed25519.PublicKey
	issuer     string
	accessTTL  time.Duration
}

// NewSigner constructs a Signer from a 32-byte Ed25519 seed. The public key
// is derived deterministically. Issuer identifies this server in the JWT.
func NewSigner(seed []byte, issuer string, accessTTL time.Duration) (*Signer, error) {
	if len(seed) != ed25519.SeedSize {
		return nil, fmt.Errorf("ed25519 seed must be %d bytes, got %d", ed25519.SeedSize, len(seed))
	}
	priv := ed25519.NewKeyFromSeed(seed)
	pub, _ := priv.Public().(ed25519.PublicKey)
	if accessTTL <= 0 {
		accessTTL = 15 * time.Minute
	}
	return &Signer{
		privateKey: priv,
		publicKey:  pub,
		issuer:     issuer,
		accessTTL:  accessTTL,
	}, nil
}

// GenerateSeed returns a fresh Ed25519 seed (32 bytes). Dev-only bootstrap;
// production wires the seed from a secret manager (Vault / AWS SM).
func GenerateSeed() ([]byte, error) {
	seed := make([]byte, ed25519.SeedSize)
	_, err := rand.Read(seed)
	return seed, err
}

// EncodeSeed base64s a seed for env var storage.
func EncodeSeed(seed []byte) string {
	return base64.StdEncoding.EncodeToString(seed)
}

// DecodeSeed reverses EncodeSeed.
func DecodeSeed(encoded string) ([]byte, error) {
	// Hosts that generate secrets (Render generateValue) may use the URL-safe
	// alphabet or drop padding; the seed bytes are what matter.
	for _, enc := range []*base64.Encoding{base64.StdEncoding, base64.URLEncoding, base64.RawStdEncoding, base64.RawURLEncoding} {
		if b, err := enc.DecodeString(encoded); err == nil {
			return b, nil
		}
	}
	return base64.StdEncoding.DecodeString(encoded)
}

// SignAccessToken issues a fresh access JWT for (account, session).
func (s *Signer) SignAccessToken(accountID, sessionID uuid.UUID, now time.Time) (string, time.Time, error) {
	exp := now.Add(s.accessTTL)
	claims := AccessClaims{
		AccountID: accountID,
		SessionID: sessionID,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    s.issuer,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(exp),
			NotBefore: jwt.NewNumericDate(now),
			ID:        uuid.NewString(),
		},
	}
	tok := jwt.NewWithClaims(jwt.SigningMethodEdDSA, claims)
	signed, err := tok.SignedString(s.privateKey)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("sign access: %w", err)
	}
	return signed, exp, nil
}

// VerifyAccessToken parses and validates the token. Returns claims or one of
// ErrTokenExpired / ErrTokenInvalid.
func (s *Signer) VerifyAccessToken(raw string) (*AccessClaims, error) {
	parser := jwt.NewParser(
		jwt.WithValidMethods([]string{jwt.SigningMethodEdDSA.Alg()}),
		jwt.WithIssuer(s.issuer),
	)
	claims := &AccessClaims{}
	tok, err := parser.ParseWithClaims(raw, claims, func(t *jwt.Token) (interface{}, error) {
		return s.publicKey, nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrTokenExpired
		}
		return nil, fmt.Errorf("%w: %v", ErrTokenInvalid, err)
	}
	if !tok.Valid {
		return nil, ErrTokenInvalid
	}
	if claims.AccountID == uuid.Nil || claims.SessionID == uuid.Nil {
		return nil, fmt.Errorf("%w: missing sub/sid", ErrTokenInvalid)
	}
	return claims, nil
}
