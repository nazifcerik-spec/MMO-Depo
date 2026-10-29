package boost

import (
	"context"
	"crypto/ecdsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"
)

// AdMob rewarded-ad server-side verification (SSV): after a player watches
// a rewarded ad the network calls the game's callback URL with the reward
// and an ECDSA signature over the query; only a callback whose signature
// verifies against Google's published keys grants a boost. The ad's
// user_id is the character id, custom_data the boost category.

// GoogleKeysURL lists the keys AdMob signs callbacks with.
const GoogleKeysURL = "https://www.gstatic.com/admob/reward/verifier-keys.json"

// KeySource looks up a verifier key by id.
type KeySource interface {
	Key(ctx context.Context, id string) (*ecdsa.PublicKey, error)
}

// GoogleKeys fetches and caches AdMob's verifier keys (refreshed daily or
// when an unknown key id shows up).
type GoogleKeys struct {
	URL    string
	Client *http.Client

	mu      sync.Mutex
	keys    map[string]*ecdsa.PublicKey
	fetched time.Time
}

func (g *GoogleKeys) Key(ctx context.Context, id string) (*ecdsa.PublicKey, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	if k, ok := g.keys[id]; ok && time.Since(g.fetched) < 24*time.Hour {
		return k, nil
	}
	if err := g.fetch(ctx); err != nil {
		return nil, err
	}
	if k, ok := g.keys[id]; ok {
		return k, nil
	}
	return nil, fmt.Errorf("%w: unknown key %s", ErrBadReward, id)
}

func (g *GoogleKeys) fetch(ctx context.Context) error {
	u := g.URL
	if u == "" {
		u = GoogleKeysURL
	}
	c := g.Client
	if c == nil {
		c = &http.Client{Timeout: 10 * time.Second}
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return err
	}
	res, err := c.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	var body struct {
		Keys []struct {
			KeyID json.Number `json:"keyId"`
			PEM   string      `json:"pem"`
		} `json:"keys"`
	}
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		return err
	}
	keys := map[string]*ecdsa.PublicKey{}
	for _, k := range body.Keys {
		block, _ := pem.Decode([]byte(k.PEM))
		if block == nil {
			continue
		}
		pub, err := x509.ParsePKIXPublicKey(block.Bytes)
		if err != nil {
			continue
		}
		if ec, ok := pub.(*ecdsa.PublicKey); ok {
			keys[k.KeyID.String()] = ec
		}
	}
	g.keys, g.fetched = keys, time.Now()
	return nil
}

// Reward is a verified ad reward.
type Reward struct {
	TransactionID string
	CharacterID   string
	Category      string
	At            time.Time
}

// VerifySSV checks an AdMob callback's raw query (as received: the signed
// content is everything before "&signature=") and returns its reward.
func VerifySSV(ctx context.Context, rawQuery string, keys KeySource) (*Reward, error) {
	i := strings.Index(rawQuery, "&signature=")
	if i < 0 {
		return nil, fmt.Errorf("%w: no signature", ErrBadReward)
	}
	content := rawQuery[:i]
	q, err := url.ParseQuery(rawQuery)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrBadReward, err)
	}
	sig, err := base64.RawURLEncoding.DecodeString(strings.TrimRight(q.Get("signature"), "="))
	if err != nil {
		return nil, fmt.Errorf("%w: signature encoding", ErrBadReward)
	}
	key, err := keys.Key(ctx, q.Get("key_id"))
	if err != nil {
		return nil, err
	}
	sum := sha256.Sum256([]byte(content))
	if !ecdsa.VerifyASN1(key, sum[:], sig) {
		return nil, fmt.Errorf("%w: bad signature", ErrBadReward)
	}
	r := &Reward{TransactionID: q.Get("transaction_id"), CharacterID: q.Get("user_id"), Category: q.Get("custom_data")}
	if ms, err := strconv.ParseInt(q.Get("timestamp"), 10, 64); err == nil {
		r.At = time.UnixMilli(ms)
	}
	if r.TransactionID == "" || r.CharacterID == "" {
		return nil, errors.New("boost: ad reward without transaction or user")
	}
	return r, nil
}
