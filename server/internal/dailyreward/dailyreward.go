// Package dailyreward is the daily login reward (DD-311): one claim per
// account and UTC day, paid to the character that claims it, on a
// 29-day cycle that restarts when a day is missed (DD-322). Gold grows with
// the week, the day and the character's level; each week's third and sixth
// days add healing potions of the character's tier.
package dailyreward

import (
	"context"
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/otter/server/internal/economy"
	"github.com/otter/server/internal/item"
)

// CycleDays is the length of the reward cycle.
const CycleDays = 29

var (
	ErrClaimed  = errors.New("dailyreward: already claimed today")
	ErrNotFound = errors.New("dailyreward: character not found")
)

// dayMult weighs a day's gold: week w (0..3) pays 1 + w/2 a day and three
// times that on its seventh day; day 29 closes the cycle at ten.
func dayMult(day int) float64 {
	if day == CycleDays {
		return 10
	}
	base := 1 + float64((day-1)/7)/2
	if day%7 == 0 {
		return 3 * base
	}
	return base
}

// potionQty is a day's healing potions (0 = none): three on each week's
// third day, five on its sixth.
func potionQty(day int) int32 {
	switch {
	case day == CycleDays:
		return 0
	case day%7 == 3:
		return 3
	case day%7 == 6:
		return 5
	}
	return 0
}

// potions are the healing potion tiers by the level they need.
var potions = []struct {
	level int
	code  string
}{
	{850, "potion_hp_supreme"}, {650, "potion_hp_huge"}, {450, "potion_hp_xlarge"},
	{250, "potion_hp_large"}, {100, "potion_hp_medium"}, {1, "potion_hp_small"},
}

// Reward is one day's reward for a character level.
type Reward struct {
	Day        int    `json:"day"`
	GoldCopper int64  `json:"gold_copper"`
	Potion     string `json:"potion,omitempty"`
	PotionQty  int32  `json:"potion_qty,omitempty"`
}

// RewardFor is day's reward at level: about five level-fitting kills' gold
// per weight (10 × L^1.5 copper each).
func RewardFor(day, level int) Reward {
	level = max(level, 1)
	r := Reward{Day: day, GoldCopper: int64(math.Round(dayMult(day) * 50 * math.Pow(float64(level), 1.5)))}
	if n := potionQty(day); n > 0 {
		for _, p := range potions {
			if level >= p.level {
				r.Potion, r.PotionQty = p.code, n
				break
			}
		}
	}
	return r
}

// Status is what the reward window shows.
type Status struct {
	Claimable bool      `json:"claimable"`
	Day       int       `json:"day"`    // the day the next claim pays (1..CycleDays)
	Streak    int       `json:"streak"` // days claimed in a row so far
	Rewards   []Reward  `json:"rewards"`
	ResetsAt  time.Time `json:"resets_at"` // next UTC midnight
}

// Minter mints a plain item (item.Service satisfies).
type Minter interface {
	MintExact(ctx context.Context, characterID uuid.UUID, templateCode string, qty int32,
		sourceKind string, sourceRefID uuid.UUID) (*item.Instance, error)
}

type Service struct {
	pool  *pgxpool.Pool
	txn   economy.TxnService
	items Minter
	now   func() time.Time
}

func NewService(pool *pgxpool.Pool, txn economy.TxnService, items Minter) *Service {
	return &Service{pool: pool, txn: txn, items: items, now: time.Now}
}

// WithClock is for tests.
func (s *Service) WithClock(now func() time.Time) *Service { s.now = now; return s }

// nextDay is the cycle day a claim made on today pays, given the last claim.
func nextDay(last time.Time, streak int, today time.Time) int {
	if last.Equal(today.AddDate(0, 0, -1)) {
		return streak%CycleDays + 1
	}
	return 1
}

func day(t time.Time) time.Time { return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC) }

// Status is the account's claim state, rewards priced for the character.
func (s *Service) Status(ctx context.Context, account, characterID uuid.UUID) (*Status, error) {
	level, err := s.level(ctx, account, characterID)
	if err != nil {
		return nil, err
	}
	today := day(s.now().UTC())
	var last time.Time
	var streak int
	err = s.pool.QueryRow(ctx, `SELECT last_day, streak FROM progression.daily_login WHERE account_id = $1`,
		account).Scan(&last, &streak)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return nil, err
	}
	st := &Status{Claimable: !day(last).Equal(today), Streak: streak, ResetsAt: today.AddDate(0, 0, 1)}
	if st.Claimable {
		st.Day = nextDay(day(last), streak, today)
		if st.Day == 1 {
			st.Streak = 0 // a missed day restarts the cycle
		}
	} else {
		st.Day = streak%CycleDays + 1 // tomorrow's
	}
	for d := 1; d <= CycleDays; d++ {
		st.Rewards = append(st.Rewards, RewardFor(d, level))
	}
	return st, nil
}

// Claim pays today's reward to the character. The account row is locked
// for the claim; the gold is keyed per account and day, so a retry after a
// failure pays it once.
func (s *Service) Claim(ctx context.Context, account, characterID uuid.UUID) (*Reward, error) {
	level, err := s.level(ctx, account, characterID)
	if err != nil {
		return nil, err
	}
	today := day(s.now().UTC())
	var out *Reward
	err = pgx.BeginFunc(ctx, s.pool, func(tx pgx.Tx) error {
		if _, err := tx.Exec(ctx, `
			INSERT INTO progression.daily_login (account_id, last_day, streak) VALUES ($1, DATE '1970-01-01', 0)
			ON CONFLICT DO NOTHING`, account); err != nil {
			return err
		}
		var last time.Time
		var streak int
		if err := tx.QueryRow(ctx, `
			SELECT last_day, streak FROM progression.daily_login WHERE account_id = $1 FOR UPDATE`,
			account).Scan(&last, &streak); err != nil {
			return err
		}
		if day(last).Equal(today) {
			return ErrClaimed
		}
		d := nextDay(day(last), streak, today)
		r := RewardFor(d, level)
		key := fmt.Sprintf("daily-%s-%s", account, today.Format("2006-01-02"))
		ref := uuid.NewSHA1(uuid.NameSpaceURL, []byte(key))
		if r.GoldCopper > 0 {
			if _, err := s.txn.Write(ctx, economy.Intent{
				Kind: "daily_reward", RefKind: "daily_login", RefID: ref, IdempotencyKey: key,
				Entries: []economy.Entry{
					{Currency: "gold", AccountKind: economy.AccountSystem, AccountRef: economy.SystemRewardSource, Delta: -r.GoldCopper},
					{Currency: "gold", AccountKind: economy.AccountCharacter, AccountRef: characterID, Delta: r.GoldCopper},
				},
			}); err != nil {
				return err
			}
		}
		if r.Potion != "" && s.items != nil {
			if _, err := s.items.MintExact(ctx, characterID, r.Potion, r.PotionQty, "daily_reward", ref); err != nil {
				return err
			}
		}
		if _, err := tx.Exec(ctx, `
			UPDATE progression.daily_login SET last_day = $2, streak = $3, claims = claims + 1, updated_at = now()
			 WHERE account_id = $1`, account, today, d); err != nil {
			return err
		}
		out = &r
		return nil
	})
	return out, err
}

func (s *Service) level(ctx context.Context, account, characterID uuid.UUID) (int, error) {
	var level int
	err := s.pool.QueryRow(ctx, `SELECT level FROM character.character WHERE id = $1 AND account_id = $2`,
		characterID, account).Scan(&level)
	if errors.Is(err, pgx.ErrNoRows) {
		return 0, ErrNotFound
	}
	return level, err
}
