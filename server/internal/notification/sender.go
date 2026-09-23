package notification

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Push is what a sender needs to deliver one message to one device.
type Push struct {
	Token    string
	Priority string // high | normal | silent
	Payload  Payload
}

// ErrInvalidToken means the provider rejected the device token for good;
// the dispatcher counts it toward deactivation.
var ErrInvalidToken = errors.New("notification: invalid device token")

type Sender interface {
	Send(ctx context.Context, p Push) error
}

// LogSender is used when a provider has no credentials configured (dev).
type LogSender struct {
	Platform string
	Log      *slog.Logger
}

func (s LogSender) Send(ctx context.Context, p Push) error {
	s.Log.InfoContext(ctx, "notification: push (log sender)",
		"platform", s.Platform, "kind", p.Payload.Kind, "priority", p.Priority)
	return nil
}

var httpClient = &http.Client{Timeout: 10 * time.Second}

// --- FCM HTTP v1 ------------------------------------------------------------

type serviceAccount struct {
	ClientEmail string `json:"client_email"`
	PrivateKey  string `json:"private_key"`
	TokenURI    string `json:"token_uri"`
	ProjectID   string `json:"project_id"`
}

// FCMSender sends through the FCM HTTP v1 API using a service-account key
// exchanged for a cached OAuth2 access token.
type FCMSender struct {
	ProjectID string
	Endpoint  string // default https://fcm.googleapis.com
	sa        serviceAccount
	key       *rsa.PrivateKey

	mu      sync.Mutex
	token   string
	expires time.Time
}

// NewFCMSender loads a Google service-account JSON file.
func NewFCMSender(projectID, credentialsFile string) (*FCMSender, error) {
	raw, err := os.ReadFile(credentialsFile)
	if err != nil {
		return nil, err
	}
	var sa serviceAccount
	if err := json.Unmarshal(raw, &sa); err != nil {
		return nil, fmt.Errorf("fcm credentials: %w", err)
	}
	key, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(sa.PrivateKey))
	if err != nil {
		return nil, fmt.Errorf("fcm private key: %w", err)
	}
	if projectID == "" {
		projectID = sa.ProjectID
	}
	if sa.TokenURI == "" {
		sa.TokenURI = "https://oauth2.googleapis.com/token"
	}
	return &FCMSender{ProjectID: projectID, Endpoint: "https://fcm.googleapis.com", sa: sa, key: key}, nil
}

func (f *FCMSender) accessToken(ctx context.Context) (string, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.token != "" && time.Until(f.expires) > time.Minute {
		return f.token, nil
	}
	now := time.Now()
	assertion, err := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"iss":   f.sa.ClientEmail,
		"scope": "https://www.googleapis.com/auth/firebase.messaging",
		"aud":   f.sa.TokenURI,
		"iat":   now.Unix(),
		"exp":   now.Add(time.Hour).Unix(),
	}).SignedString(f.key)
	if err != nil {
		return "", err
	}
	form := url.Values{
		"grant_type": {"urn:ietf:params:oauth:grant-type:jwt-bearer"},
		"assertion":  {assertion},
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, f.sa.TokenURI, strings.NewReader(form.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	var out struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
	}
	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		return "", fmt.Errorf("fcm oauth: %d %s", resp.StatusCode, b)
	}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return "", err
	}
	f.token, f.expires = out.AccessToken, now.Add(time.Duration(out.ExpiresIn)*time.Second)
	return f.token, nil
}

func (f *FCMSender) Send(ctx context.Context, p Push) error {
	tok, err := f.accessToken(ctx)
	if err != nil {
		return err
	}
	data := map[string]string{"kind": p.Payload.Kind, "route": p.Payload.Route,
		"title_key": p.Payload.TitleKey, "body_key": p.Payload.BodyKey}
	if len(p.Payload.Vars) > 0 {
		vars, _ := json.Marshal(p.Payload.Vars)
		data["body_vars"] = string(vars)
	}
	android := map[string]any{"priority": "normal"}
	if p.Priority == "high" {
		android["priority"] = "high"
	}
	if p.Payload.CollapseKey != "" {
		android["collapse_key"] = p.Payload.CollapseKey
	}
	msg := map[string]any{"token": p.Token, "data": data, "android": android}
	if p.Priority != "silent" {
		// Client resolves the localisation keys; the notification block makes
		// the OS show it while the app is in the background.
		android["notification"] = map[string]any{
			"channel_id": p.Payload.Category, "title_loc_key": p.Payload.TitleKey, "body_loc_key": p.Payload.BodyKey,
		}
	}
	body, _ := json.Marshal(map[string]any{"message": msg})
	req, err := http.NewRequestWithContext(ctx, http.MethodPost,
		fmt.Sprintf("%s/v1/projects/%s/messages:send", f.Endpoint, f.ProjectID), bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+tok)
	req.Header.Set("Content-Type", "application/json")
	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		return nil
	}
	b, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
	if resp.StatusCode == http.StatusNotFound || bytes.Contains(b, []byte("UNREGISTERED")) ||
		(resp.StatusCode == http.StatusBadRequest && bytes.Contains(b, []byte("registration token"))) {
		return ErrInvalidToken
	}
	return fmt.Errorf("fcm send: %d %s", resp.StatusCode, b)
}

// --- APNs (token-based auth) --------------------------------------------------

// APNsSender sends through Apple's HTTP/2 provider API with an ES256 JWT.
type APNsSender struct {
	KeyID    string
	TeamID   string
	Topic    string // app bundle id
	Endpoint string // https://api.push.apple.com or https://api.sandbox.push.apple.com
	key      *ecdsa.PrivateKey

	mu     sync.Mutex
	jwt    string
	issued time.Time
}

// NewAPNsSender loads the .p8 key file from Apple.
func NewAPNsSender(keyFile, keyID, teamID, topic string, sandbox bool) (*APNsSender, error) {
	raw, err := os.ReadFile(keyFile)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(raw)
	if block == nil {
		return nil, errors.New("apns: key file is not PEM")
	}
	parsed, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("apns key: %w", err)
	}
	key, ok := parsed.(*ecdsa.PrivateKey)
	if !ok {
		return nil, errors.New("apns: key is not ECDSA")
	}
	endpoint := "https://api.push.apple.com"
	if sandbox {
		endpoint = "https://api.sandbox.push.apple.com"
	}
	return &APNsSender{KeyID: keyID, TeamID: teamID, Topic: topic, Endpoint: endpoint, key: key}, nil
}

// providerToken is refreshed every 50 min (Apple rejects tokens > 1h old).
func (a *APNsSender) providerToken() (string, error) {
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.jwt != "" && time.Since(a.issued) < 50*time.Minute {
		return a.jwt, nil
	}
	t := jwt.NewWithClaims(jwt.SigningMethodES256, jwt.MapClaims{"iss": a.TeamID, "iat": time.Now().Unix()})
	t.Header["kid"] = a.KeyID
	s, err := t.SignedString(a.key)
	if err != nil {
		return "", err
	}
	a.jwt, a.issued = s, time.Now()
	return s, nil
}

func (a *APNsSender) Send(ctx context.Context, p Push) error {
	tok, err := a.providerToken()
	if err != nil {
		return err
	}
	aps := map[string]any{"thread-id": p.Payload.Category}
	pushType, prio := "alert", "10"
	if p.Priority == "silent" {
		aps["content-available"] = 1
		pushType, prio = "background", "5"
	} else {
		aps["alert"] = map[string]any{"title-loc-key": p.Payload.TitleKey, "loc-key": p.Payload.BodyKey}
		aps["sound"] = "default"
		aps["category"] = p.Payload.Category
	}
	body, _ := json.Marshal(map[string]any{
		"aps": aps, "kind": p.Payload.Kind, "route": p.Payload.Route, "body_vars": p.Payload.Vars,
	})
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, a.Endpoint+"/3/device/"+url.PathEscape(p.Token), bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("authorization", "bearer "+tok)
	req.Header.Set("apns-topic", a.Topic)
	req.Header.Set("apns-push-type", pushType)
	req.Header.Set("apns-priority", prio)
	if p.Payload.CollapseKey != "" {
		req.Header.Set("apns-collapse-id", p.Payload.CollapseKey)
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		return nil
	}
	b, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
	if resp.StatusCode == http.StatusGone || bytes.Contains(b, []byte("BadDeviceToken")) ||
		bytes.Contains(b, []byte("Unregistered")) {
		return ErrInvalidToken
	}
	return fmt.Errorf("apns send: %d %s", resp.StatusCode, b)
}
