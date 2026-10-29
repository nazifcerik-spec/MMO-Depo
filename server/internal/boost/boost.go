// Package boost is timed XP boosts (DD-320): +100% XP for two hours in one
// of three categories — gathering, adventure (hunts, dungeons, raids) and
// professions (crafting). A boost comes from watching a rewarded ad
// (verified by the ad network's server-side callback) or from the credit
// shop. A finished job's XP grows by the share of its time a boost of its
// category covered, so a boost switched on at the end of a long run adds
// nothing for the hours before it.
package boost

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	CatGather    = "gather"
	CatAdventure = "adventure"
	CatCraft     = "craft"

	// Duration and BonusPct are a boost's length and XP bonus.
	Duration = 2 * time.Hour
	BonusPct = 100
	// AdsPerDay caps ad boosts per character, category and UTC day.
	AdsPerDay = 3

	SourceAd     = "ad"
	SourceCredit = "credit"
)

// Categories are the boostable kinds of work.
var Categories = []string{CatGather, CatAdventure, CatCraft}

var (
	ErrCategory  = errors.New("boost: unknown category")
	ErrActive    = errors.New("boost: already running")
	ErrAdLimit   = errors.New("boost: no more ad boosts today")
	ErrNotFound  = errors.New("boost: character not found")
	ErrAdsOff    = errors.New("boost: rewarded ads are not set up")
	ErrBadReward = errors.New("boost: the ad reward did not verify")
)

func known(cat string) bool {
	for _, c := range Categories {
		if c == cat {
			return true
		}
	}
	return false
}

type Service struct {
	pool *pgxpool.Pool
	now  func() time.Time
}

func NewService(pool *pgxpool.Pool) *Service { return &Service{pool: pool, now: time.Now} }

// WithClock is for tests.
func (s *Service) WithClock(now func() time.Time) *Service { s.now = now; return s }

// Status is one category's boost for a character.
type Status struct {
	Category string     `json:"category"`
	EndsAt   *time.Time `json:"ends_at,omitempty"` // a running boost
	AdsLeft  int        `json:"ads_left"`          // ad boosts left today
}

// List is the character's three categories.
func (s *Service) List(ctx context.Context, char uuid.UUID) ([]Status, error) {
	now := s.now().UTC()
	day := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	out := make([]Status, 0, len(Categories))
	for _, c := range Categories {
		st := Status{Category: c, AdsLeft: AdsPerDay}
		var ends *time.Time
		var ads int
		if err := s.pool.QueryRow(ctx, `
			SELECT (SELECT max(ends_at) FROM boost.grant WHERE character_id = $1 AND category = $2 AND ends_at > $3),
			       (SELECT count(*) FROM boost.grant WHERE character_id = $1 AND category = $2 AND source = 'ad'
			         AND starts_at >= $4)`, char, c, now, day).Scan(&ends, &ads); err != nil {
			return nil, err
		}
		st.EndsAt, st.AdsLeft = ends, max(0, AdsPerDay-ads)
		out = append(out, st)
	}
	return out, nil
}

// Grant starts a boost of the category now. ref keys it once (the ad's
// transaction id, or the shop purchase), so a retried callback grants once.
func (s *Service) Grant(ctx context.Context, char uuid.UUID, cat, source, ref string) (*Status, error) {
	if !known(cat) {
		return nil, ErrCategory
	}
	now := s.now().UTC()
	day := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	err := pgx.BeginFunc(ctx, s.pool, func(tx pgx.Tx) error {
		if _, err := tx.Exec(ctx, `SELECT pg_advisory_xact_lock(hashtext('boost:' || $1::text))`, char); err != nil {
			return err
		}
		var dup bool
		if err := tx.QueryRow(ctx, `SELECT EXISTS (SELECT 1 FROM boost.grant WHERE ref = $1)`, ref).Scan(&dup); err != nil {
			return err
		}
		if dup {
			return nil // a retried callback or request
		}
		var active bool
		var ads int
		if err := tx.QueryRow(ctx, `
			SELECT EXISTS (SELECT 1 FROM boost.grant WHERE character_id = $1 AND category = $2 AND ends_at > $3),
			       (SELECT count(*) FROM boost.grant WHERE character_id = $1 AND category = $2 AND source = 'ad'
			         AND starts_at >= $4)`, char, cat, now, day).Scan(&active, &ads); err != nil {
			return err
		}
		if active {
			return ErrActive
		}
		if source == SourceAd && ads >= AdsPerDay {
			return ErrAdLimit
		}
		_, err := tx.Exec(ctx, `INSERT INTO boost.grant (character_id, category, source, ref, starts_at, ends_at)
			VALUES ($1, $2, $3, $4, $5, $6)`, char, cat, source, ref, now, now.Add(Duration))
		if isFK(err) {
			return ErrNotFound
		}
		return err
	})
	if err != nil {
		return nil, err
	}
	list, err := s.List(ctx, char)
	if err != nil {
		return nil, err
	}
	for _, st := range list {
		if st.Category == cat {
			return &st, nil
		}
	}
	return nil, fmt.Errorf("boost: %s missing", cat)
}

// Apply is xp grown by the boosts of cat covering [start, end): +BonusPct
// on the covered share. A job without length is boosted when a boost runs
// at its end.
func (s *Service) Apply(ctx context.Context, char uuid.UUID, cat string, start, end time.Time, xp int64) int64 {
	if xp <= 0 || s == nil {
		return xp
	}
	if !end.After(start) {
		start = end.Add(-time.Second)
	}
	rows, err := s.pool.Query(ctx, `SELECT starts_at, ends_at FROM boost.grant
		WHERE character_id = $1 AND category = $2 AND starts_at < $4 AND ends_at > $3`, char, cat, start, end)
	if err != nil {
		return xp // a boost never fails the reward it rides on
	}
	defer rows.Close()
	var covered time.Duration
	for rows.Next() {
		var a, b time.Time
		if rows.Scan(&a, &b) != nil {
			return xp
		}
		if b.After(end) {
			b = end
		}
		if a.Before(start) {
			a = start
		}
		covered += b.Sub(a)
	}
	share := min(1, float64(covered)/float64(end.Sub(start)))
	return xp + int64(float64(xp)*share*BonusPct/100+0.5)
}

// Applier is Apply as other packages take it (nil-safe through ApplyIf).
type Applier interface {
	Apply(ctx context.Context, char uuid.UUID, cat string, start, end time.Time, xp int64) int64
}

// ApplyIf boosts through b, or leaves xp as it is.
func ApplyIf(ctx context.Context, b Applier, char uuid.UUID, cat string, start, end time.Time, xp int64) int64 {
	if b == nil {
		return xp
	}
	return b.Apply(ctx, char, cat, start, end, xp)
}

func isFK(err error) bool {
	var pg interface{ SQLState() string }
	return errors.As(err, &pg) && pg.SQLState() == "23503"
}
