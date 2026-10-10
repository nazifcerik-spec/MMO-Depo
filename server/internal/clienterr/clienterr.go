// Package clienterr collects uncaught client errors (DD-235). Reports are
// grouped by a fingerprint into one bounded row each; no account or IP is
// stored and obvious secrets (e-mails, tokens) are masked before storage.
package clienterr

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Limits.
const (
	MaxMessage   = 1000
	MaxStack     = 8000
	MaxStackRows = 40
	// MaxDistinct caps the table: new fingerprints beyond it are dropped
	// (known ones keep counting) so a flood cannot grow storage.
	MaxDistinct = 10_000
)

var (
	kinds     = map[string]bool{"flutter": true, "platform": true, "zone": true}
	platforms = map[string]bool{"android": true, "ios": true, "web": true, "windows": true, "macos": true, "linux": true}

	ErrInvalid = errors.New("clienterr: invalid report")
)

// Report is what the app sends.
type Report struct {
	Kind       string `json:"kind"`
	Message    string `json:"message"`
	Stack      string `json:"stack"`
	Platform   string `json:"platform"`
	AppVersion string `json:"app_version"`
	Flavor     string `json:"flavor"`
}

// Group is one fingerprint as operators see it.
type Group struct {
	Fingerprint string    `json:"fingerprint"`
	Kind        string    `json:"kind"`
	Message     string    `json:"message"`
	Stack       string    `json:"stack"`
	Platform    string    `json:"platform"`
	AppVersion  string    `json:"app_version"`
	Flavor      string    `json:"flavor"`
	Count       int64     `json:"count"`
	FirstSeen   time.Time `json:"first_seen"`
	LastSeen    time.Time `json:"last_seen"`
}

type Service struct {
	pool    *pgxpool.Pool
	counter func(platform, kind string) // metric hook; nil = none
}

func NewService(pool *pgxpool.Pool, counter func(platform, kind string)) *Service {
	return &Service{pool: pool, counter: counter}
}

var (
	reEmail   = regexp.MustCompile(`[\w.+-]+@[\w-]+\.[\w.-]+`)
	reBearer  = regexp.MustCompile(`(?i)bearer\s+\S+`)
	reJWT     = regexp.MustCompile(`eyJ[\w-]+\.[\w-]+\.[\w-]+`)
	reQuery   = regexp.MustCompile(`(?i)([?&](?:token|code|password|key)=)[^&\s]+`)
	reVersion = regexp.MustCompile(`^[\w.+-]*$`)

	// Fingerprint normalisation: frame numbers, line:col, addresses, numbers.
	reFrameNo = regexp.MustCompile(`^#\d+\s+`)
	reLineCol = regexp.MustCompile(`:\d+(:\d+)?`)
	reHex     = regexp.MustCompile(`0x[0-9a-fA-F]+`)
	reDigits  = regexp.MustCompile(`\d+`)
)

func redact(s string) string {
	s = reJWT.ReplaceAllString(s, "***")
	s = reBearer.ReplaceAllString(s, "Bearer ***")
	s = reQuery.ReplaceAllString(s, "${1}***")
	return reEmail.ReplaceAllString(s, "***@***")
}

func truncate(s string, n int) string {
	if utf8.RuneCountInString(s) <= n {
		return s
	}
	return string([]rune(s)[:n])
}

func clean(r Report) (Report, error) {
	r.Kind = strings.TrimSpace(r.Kind)
	if !kinds[r.Kind] {
		return r, ErrInvalid
	}
	r.Message = truncate(redact(strings.TrimSpace(r.Message)), MaxMessage)
	if r.Message == "" {
		return r, ErrInvalid
	}
	lines := strings.Split(strings.TrimSpace(r.Stack), "\n")
	if len(lines) > MaxStackRows {
		lines = lines[:MaxStackRows]
	}
	r.Stack = truncate(redact(strings.Join(lines, "\n")), MaxStack)
	if !platforms[r.Platform] {
		r.Platform = "other"
	}
	if len(r.AppVersion) > 32 || !reVersion.MatchString(r.AppVersion) {
		r.AppVersion = ""
	}
	if len(r.Flavor) > 16 || !reVersion.MatchString(r.Flavor) {
		r.Flavor = ""
	}
	return r, nil
}

func normalise(s string) string {
	return reDigits.ReplaceAllString(reHex.ReplaceAllString(s, "#"), "#")
}

// fingerprint groups the same error across devices and builds: kind, the
// message's first line and the top three frames, numbers stripped.
func fingerprint(r Report) string {
	first, _, _ := strings.Cut(r.Message, "\n")
	parts := []string{r.Kind, normalise(first)}
	for _, l := range strings.Split(r.Stack, "\n") {
		if l = strings.TrimSpace(l); l == "" {
			continue
		}
		parts = append(parts, normalise(reLineCol.ReplaceAllString(reFrameNo.ReplaceAllString(l, ""), "")))
		if len(parts) == 5 {
			break
		}
	}
	sum := sha256.Sum256([]byte(strings.Join(parts, "\n")))
	return hex.EncodeToString(sum[:])
}

// Record stores one report (count+1 on a known fingerprint).
func (s *Service) Record(ctx context.Context, r Report) error {
	r, err := clean(r)
	if err != nil {
		return err
	}
	if s.counter != nil {
		s.counter(r.Platform, r.Kind)
	}
	fp := fingerprint(r)
	tag, err := s.pool.Exec(ctx, `
UPDATE ops.client_error SET count = count + 1, last_seen = now(),
       platform = $2, app_version = $3, flavor = $4
 WHERE fingerprint = $1`, fp, r.Platform, r.AppVersion, r.Flavor)
	if err != nil || tag.RowsAffected() > 0 {
		return err
	}
	_, err = s.pool.Exec(ctx, `
INSERT INTO ops.client_error (fingerprint, kind, message, stack, platform, app_version, flavor)
SELECT $1, $2, $3, $4, $5, $6, $7
 WHERE (SELECT count(*) FROM ops.client_error) < $8
ON CONFLICT (fingerprint) DO UPDATE SET count = ops.client_error.count + 1, last_seen = now()`,
		fp, r.Kind, r.Message, r.Stack, r.Platform, r.AppVersion, r.Flavor, MaxDistinct)
	return err
}

// Recent lists the most recently seen groups (operators).
func (s *Service) Recent(ctx context.Context, platform string, limit int) ([]Group, error) {
	if limit <= 0 || limit > 200 {
		limit = 50
	}
	rows, err := s.pool.Query(ctx, `
SELECT fingerprint, kind, message, stack, platform, app_version, flavor, count, first_seen, last_seen
  FROM ops.client_error
 WHERE $1::text = '' OR platform = $1::text
 ORDER BY last_seen DESC LIMIT $2`, platform, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []Group{}
	for rows.Next() {
		var g Group
		if err := rows.Scan(&g.Fingerprint, &g.Kind, &g.Message, &g.Stack, &g.Platform, &g.AppVersion,
			&g.Flavor, &g.Count, &g.FirstSeen, &g.LastSeen); err != nil {
			return nil, err
		}
		out = append(out, g)
	}
	return out, rows.Err()
}
