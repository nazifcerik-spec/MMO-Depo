package worldevent

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/otter/server/internal/domain/character"
	"github.com/otter/server/internal/economy"
	"github.com/otter/server/internal/notification"
)

// Characters is what the service needs from the character domain.
// character.Service satisfies it.
type Characters interface {
	Owns(ctx context.Context, characterID, accountID uuid.UUID) error
	Sheet(ctx context.Context, characterID, accountID uuid.UUID) (*character.Sheet, error)
	GetLocation(ctx context.Context, characterID uuid.UUID) (string, error)
}

type Service struct {
	pool           *pgxpool.Pool
	chars          Characters
	txn            economy.TxnService
	startingRegion func() string
	now            func() time.Time
	log            *slog.Logger
}

// NewService: startingRegion resolves the fallback region for characters
// that never travelled (world catalog's starting region).
func NewService(pool *pgxpool.Pool, chars Characters, txn economy.TxnService, startingRegion func() string, logger *slog.Logger) *Service {
	if logger == nil {
		logger = slog.Default()
	}
	return &Service{pool: pool, chars: chars, txn: txn, startingRegion: startingRegion, now: time.Now, log: logger}
}

// --- reads ------------------------------------------------------------------

type SpawnView struct {
	ID             uuid.UUID  `json:"id"`
	TemplateCode   string     `json:"template_code"`
	Kind           string     `json:"kind"`
	DisplayNameKey string     `json:"display_name_key"`
	MinLevel       int16      `json:"min_level"`
	RegionCode     *string    `json:"region_code,omitempty"`
	GoalScore      int64      `json:"goal_score"`
	TotalScore     int64      `json:"total_score"`
	StartsAt       time.Time  `json:"starts_at"`
	EndsAt         time.Time  `json:"ends_at"`
	State          string     `json:"state"`
	Tiers          []Tier     `json:"tiers"`
	Enlisted       bool       `json:"enlisted"`
	MyPower        *int       `json:"my_power,omitempty"`
	MyScore        *int64     `json:"my_score,omitempty"`
	MyTier         *int16     `json:"my_reward_tier,omitempty"`
	RewardedAt     *time.Time `json:"rewarded_at,omitempty"`
}

// List shows scheduled/active spawns plus those closed in the last 24h,
// with the character's own attendance.
func (s *Service) List(ctx context.Context, char, account uuid.UUID) ([]SpawnView, error) {
	if err := s.chars.Owns(ctx, char, account); err != nil {
		return nil, err
	}
	rows, err := s.pool.Query(ctx, `
		SELECT sp.id, t.code, t.kind, t.display_name_key, t.min_level, t.region_code, t.hp_base,
		       sp.total_score, sp.starts_at, sp.ends_at, sp.state, t.tiers,
		       a.character_id IS NOT NULL, a.power, a.contribution_score, a.reward_tier, a.rewarded_at
		  FROM worldevent.spawn sp
		  JOIN content.world_event_template t ON t.code = sp.template_code
		  LEFT JOIN worldevent.attendance a ON a.spawn_id = sp.id AND a.character_id = $1
		 WHERE sp.state IN ('scheduled','active') OR sp.closed_at > now() - INTERVAL '24 hours'
		 ORDER BY sp.starts_at`, char)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, func(r pgx.CollectableRow) (SpawnView, error) {
		var v SpawnView
		var tiers []byte
		var score *int64
		err := r.Scan(&v.ID, &v.TemplateCode, &v.Kind, &v.DisplayNameKey, &v.MinLevel, &v.RegionCode, &v.GoalScore,
			&v.TotalScore, &v.StartsAt, &v.EndsAt, &v.State, &tiers,
			&v.Enlisted, &v.MyPower, &score, &v.MyTier, &v.RewardedAt)
		if err == nil {
			err = json.Unmarshal(tiers, &v.Tiers)
		}
		if v.Enlisted {
			v.MyScore = score
		}
		return v, err
	})
}

// --- attend -----------------------------------------------------------------

// Attend enlists the character with a power snapshot (DD-088 style freeze).
// Open from 24h before start until the event ends.
func (s *Service) Attend(ctx context.Context, char, account, spawnID uuid.UUID) (*SpawnView, error) {
	if err := s.chars.Owns(ctx, char, account); err != nil {
		return nil, err
	}
	var state, kind string
	var start, end time.Time
	var minLevel int16
	var region *string
	err := s.pool.QueryRow(ctx, `
		SELECT sp.state, sp.starts_at, sp.ends_at, t.kind, t.min_level, t.region_code
		  FROM worldevent.spawn sp JOIN content.world_event_template t ON t.code = sp.template_code
		 WHERE sp.id = $1`, spawnID).Scan(&state, &start, &end, &kind, &minLevel, &region)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	now := s.now()
	open := (state == "scheduled" && start.Sub(now) <= RegistrationWindow) || (state == "active" && now.Before(end))
	if !open {
		return nil, ErrNotOpen
	}
	sheet, err := s.chars.Sheet(ctx, char, account)
	if err != nil {
		return nil, err
	}
	if int16(sheet.Level) < minLevel {
		return nil, ErrLevelRequirement
	}
	if region != nil {
		loc, err := s.chars.GetLocation(ctx, char)
		if err != nil {
			return nil, err
		}
		if loc == "" && s.startingRegion != nil {
			loc = s.startingRegion()
		}
		if loc != *region {
			return nil, ErrWrongRegion
		}
	}
	power := Power(sheet.Derived.Attack, sheet.Derived.Defense, sheet.Resources.HPMax)
	_, err = s.pool.Exec(ctx, `
		INSERT INTO worldevent.attendance (spawn_id, character_id, account_id, registered_at, power)
		VALUES ($1, $2, $3, $4, $5)`, spawnID, char, account, now.UTC(), power)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, ErrAlreadyAttending
		}
		return nil, err
	}
	list, err := s.List(ctx, char, account)
	if err != nil {
		return nil, err
	}
	for i := range list {
		if list[i].ID == spawnID {
			return &list[i], nil
		}
	}
	return nil, ErrNotFound
}

// --- seasonal ---------------------------------------------------------------

type SeasonalView struct {
	Code           string    `json:"code"`
	DisplayNameKey string    `json:"display_name_key"`
	StartsAt       time.Time `json:"starts_at"`
	EndsAt         time.Time `json:"ends_at"`
	DailyTokens    int       `json:"daily_login_tokens"`
	ClaimedToday   bool      `json:"claimed_today"`
	TokenBalance   int64     `json:"event_token_balance"`
}

func (s *Service) Seasonal(ctx context.Context, account uuid.UUID) ([]SeasonalView, error) {
	today := s.now().UTC().Format("2006-01-02")
	rows, err := s.pool.Query(ctx, `
		SELECT e.code, e.display_name_key, e.starts_at, e.ends_at, e.daily_login_tokens,
		       EXISTS (SELECT 1 FROM worldevent.seasonal_claim c
		                WHERE c.event_code = e.code AND c.account_id = $1 AND c.claim_date = $2::date),
		       COALESCE((SELECT amount FROM economy.account_currency_balance
		                  WHERE account_id = $1 AND currency_code = 'event_token'), 0)
		  FROM content.seasonal_event e
		 WHERE now() BETWEEN e.starts_at AND e.ends_at
		 ORDER BY e.ends_at`, account, today)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, func(r pgx.CollectableRow) (SeasonalView, error) {
		var v SeasonalView
		err := r.Scan(&v.Code, &v.DisplayNameKey, &v.StartsAt, &v.EndsAt, &v.DailyTokens, &v.ClaimedToday, &v.TokenBalance)
		return v, err
	})
}

// ClaimDaily grants the once-per-UTC-day login tokens to the account. The
// ledger key is per (event, account, day), so concurrent claims pay once.
func (s *Service) ClaimDaily(ctx context.Context, account uuid.UUID, code string) (int, error) {
	now := s.now().UTC()
	today := now.Format("2006-01-02")
	var tokens int
	var start, end time.Time
	err := s.pool.QueryRow(ctx, `
		SELECT daily_login_tokens, starts_at, ends_at FROM content.seasonal_event WHERE code = $1`, code).
		Scan(&tokens, &start, &end)
	if errors.Is(err, pgx.ErrNoRows) {
		return 0, ErrNotFound
	}
	if err != nil {
		return 0, err
	}
	if now.Before(start) || !now.Before(end) {
		return 0, ErrEventNotActive
	}
	var claimed bool
	if err := s.pool.QueryRow(ctx, `
		SELECT EXISTS (SELECT 1 FROM worldevent.seasonal_claim
		                WHERE event_code = $1 AND account_id = $2 AND claim_date = $3::date)`,
		code, account, today).Scan(&claimed); err != nil {
		return 0, err
	}
	if claimed {
		return 0, ErrAlreadyClaimed
	}
	if _, err := s.txn.Write(ctx, economy.Intent{
		Kind: "seasonal_daily", RefKind: "seasonal_event",
		IdempotencyKey: fmt.Sprintf("seasonal-%s-%s-%s", code, account, today),
		Entries: []economy.Entry{
			{Currency: "event_token", AccountKind: economy.AccountSystem, AccountRef: economy.SystemRewardSource, Delta: -int64(tokens)},
			{Currency: "event_token", AccountKind: economy.AccountAccount, AccountRef: account, Delta: int64(tokens)},
		},
	}); err != nil {
		return 0, err
	}
	if _, err := s.pool.Exec(ctx, `
		INSERT INTO worldevent.seasonal_claim (event_code, account_id, claim_date, tokens)
		VALUES ($1, $2, $3::date, $4) ON CONFLICT DO NOTHING`, code, account, today, tokens); err != nil {
		return 0, err
	}
	return tokens, nil
}

// --- worker -----------------------------------------------------------------

// Tick is the worker entry point: schedule upcoming spawns from the weekly
// rotation, activate due ones, close finished ones and pay pending rewards.
func Tick(ctx context.Context, pool *pgxpool.Pool, txn economy.TxnService, notify notification.Notifier, log *slog.Logger) (int64, error) {
	var n int64
	steps := []func(context.Context, *pgxpool.Pool, economy.TxnService, notification.Notifier, *slog.Logger) (int64, error){
		schedule, activate, closeDue, payRewards,
	}
	for _, step := range steps {
		k, err := step(ctx, pool, txn, notify, log)
		n += k
		if err != nil {
			return n, err
		}
	}
	return n, nil
}

func schedule(ctx context.Context, pool *pgxpool.Pool, _ economy.TxnService, _ notification.Notifier, _ *slog.Logger) (int64, error) {
	rows, err := pool.Query(ctx, `
		SELECT r.template_code, r.day_of_week, r.hour_utc, t.duration_seconds
		  FROM content.world_event_rotation r JOIN content.world_event_template t ON t.code = r.template_code`)
	if err != nil {
		return 0, err
	}
	type slot struct {
		code            string
		dow, hour, secs int
	}
	slots, err := pgx.CollectRows(rows, func(r pgx.CollectableRow) (slot, error) {
		var s slot
		err := r.Scan(&s.code, &s.dow, &s.hour, &s.secs)
		return s, err
	})
	if err != nil {
		return 0, err
	}
	var n int64
	now := time.Now()
	for _, sl := range slots {
		for _, start := range NextOccurrences(now, sl.dow, sl.hour, RegistrationWindow) {
			tag, err := pool.Exec(ctx, `
				INSERT INTO worldevent.spawn (template_code, starts_at, ends_at) VALUES ($1, $2, $3)
				ON CONFLICT (template_code, starts_at) DO NOTHING`,
				sl.code, start, start.Add(time.Duration(sl.secs)*time.Second))
			if err != nil {
				return n, err
			}
			n += tag.RowsAffected()
		}
	}
	return n, nil
}

func activate(ctx context.Context, pool *pgxpool.Pool, _ economy.TxnService, notify notification.Notifier, log *slog.Logger) (int64, error) {
	rows, err := pool.Query(ctx, `
		UPDATE worldevent.spawn SET state = 'active'
		 WHERE state = 'scheduled' AND starts_at <= now()
		RETURNING id, template_code`)
	if err != nil {
		return 0, err
	}
	type act struct {
		id   uuid.UUID
		code string
	}
	acts, err := pgx.CollectRows(rows, func(r pgx.CollectableRow) (act, error) {
		var a act
		err := r.Scan(&a.id, &a.code)
		return a, err
	})
	if err != nil {
		return 0, err
	}
	for _, a := range acts {
		arows, err := pool.Query(ctx, `
			SELECT character_id, account_id FROM worldevent.attendance WHERE spawn_id = $1`, a.id)
		if err != nil {
			return 0, err
		}
		type att struct{ char, account uuid.UUID }
		atts, err := pgx.CollectRows(arows, func(r pgx.CollectableRow) (att, error) {
			var x att
			err := r.Scan(&x.char, &x.account)
			return x, err
		})
		if err != nil {
			return 0, err
		}
		for _, x := range atts {
			ch := x.char
			notification.Send(ctx, notify, log, notification.Notice{
				AccountID: x.account, CharacterID: &ch, Kind: "world_event_active",
				DedupKey: "world_event_active:" + a.id.String(),
				TitleKey: "notif.world_event_active.title", BodyKey: "notif.world_event_active.body",
				Vars: map[string]any{"event": a.code}, Route: "world-events/" + a.id.String(),
				Category: "event", High: true,
			})
		}
	}
	return int64(len(acts)), nil
}

// closeDue scores every attendee, decides success and assigns reward tiers
// in one transaction per spawn. Payment happens in payRewards so a crash
// between the two never loses or doubles a reward.
func closeDue(ctx context.Context, pool *pgxpool.Pool, _ economy.TxnService, _ notification.Notifier, _ *slog.Logger) (int64, error) {
	rows, err := pool.Query(ctx, `SELECT id FROM worldevent.spawn WHERE state = 'active' AND ends_at <= now()`)
	if err != nil {
		return 0, err
	}
	ids, err := pgx.CollectRows(rows, pgx.RowTo[uuid.UUID])
	if err != nil {
		return 0, err
	}
	var n int64
	for _, id := range ids {
		err := pgx.BeginFunc(ctx, pool, func(tx pgx.Tx) error {
			var start, end time.Time
			var goal int64
			var tiersRaw []byte
			var capTier int16
			err := tx.QueryRow(ctx, `
				SELECT sp.starts_at, sp.ends_at, t.hp_base, t.tiers, t.afk_cap_tier
				  FROM worldevent.spawn sp JOIN content.world_event_template t ON t.code = sp.template_code
				 WHERE sp.id = $1 AND sp.state = 'active' FOR UPDATE OF sp SKIP LOCKED`, id).
				Scan(&start, &end, &goal, &tiersRaw, &capTier)
			if errors.Is(err, pgx.ErrNoRows) {
				return nil
			}
			if err != nil {
				return err
			}
			var tiers []Tier
			if err := json.Unmarshal(tiersRaw, &tiers); err != nil {
				return err
			}
			var capScore int64
			if capTier >= 1 && int(capTier) <= len(tiers) {
				capScore = tiers[capTier-1].MinScore
			}
			arows, err := tx.Query(ctx, `
				SELECT character_id, registered_at, power FROM worldevent.attendance WHERE spawn_id = $1`, id)
			if err != nil {
				return err
			}
			type att struct {
				char  uuid.UUID
				reg   time.Time
				power int
				score int64
			}
			atts, err := pgx.CollectRows(arows, func(r pgx.CollectableRow) (att, error) {
				var a att
				err := r.Scan(&a.char, &a.reg, &a.power)
				return a, err
			})
			if err != nil {
				return err
			}
			var total int64
			for i := range atts {
				atts[i].score = Contribution(atts[i].power, AttendedMinutes(atts[i].reg, start, end), capScore)
				total += atts[i].score
			}
			succeeded := total >= goal
			for _, a := range atts {
				if _, err := tx.Exec(ctx, `
					UPDATE worldevent.attendance SET contribution_score = $3, reward_tier = $4
					 WHERE spawn_id = $1 AND character_id = $2`,
					id, a.char, a.score, TierFor(a.score, tiers, succeeded)); err != nil {
					return err
				}
			}
			state := "failed"
			if succeeded {
				state = "succeeded"
			}
			_, err = tx.Exec(ctx, `
				UPDATE worldevent.spawn SET state = $2, total_score = $3, closed_at = now() WHERE id = $1`,
				id, state, total)
			if err == nil {
				n++
			}
			return err
		})
		if err != nil {
			return n, err
		}
	}
	return n, nil
}

func payRewards(ctx context.Context, pool *pgxpool.Pool, txn economy.TxnService, _ notification.Notifier, log *slog.Logger) (int64, error) {
	rows, err := pool.Query(ctx, `
		SELECT a.spawn_id, a.character_id, a.account_id, a.reward_tier, t.tiers
		  FROM worldevent.attendance a
		  JOIN worldevent.spawn sp ON sp.id = a.spawn_id
		  JOIN content.world_event_template t ON t.code = sp.template_code
		 WHERE a.reward_tier IS NOT NULL AND a.rewarded_at IS NULL
		 LIMIT 500`)
	if err != nil {
		return 0, err
	}
	type due struct {
		spawn, char, account uuid.UUID
		tier                 int16
		tiers                []byte
	}
	dues, err := pgx.CollectRows(rows, func(r pgx.CollectableRow) (due, error) {
		var d due
		err := r.Scan(&d.spawn, &d.char, &d.account, &d.tier, &d.tiers)
		return d, err
	})
	if err != nil {
		return 0, err
	}
	var n int64
	for _, d := range dues {
		var tiers []Tier
		if err := json.Unmarshal(d.tiers, &tiers); err != nil {
			return n, err
		}
		if d.tier >= 1 && int(d.tier) <= len(tiers) {
			t := tiers[d.tier-1]
			var entries []economy.Entry
			if t.GoldCopper > 0 {
				entries = append(entries,
					economy.Entry{Currency: "gold", AccountKind: economy.AccountSystem, AccountRef: economy.SystemRewardSource, Delta: -t.GoldCopper},
					economy.Entry{Currency: "gold", AccountKind: economy.AccountCharacter, AccountRef: d.char, Delta: t.GoldCopper})
			}
			if t.EventTokens > 0 {
				entries = append(entries,
					economy.Entry{Currency: "event_token", AccountKind: economy.AccountSystem, AccountRef: economy.SystemRewardSource, Delta: -t.EventTokens},
					economy.Entry{Currency: "event_token", AccountKind: economy.AccountAccount, AccountRef: d.account, Delta: t.EventTokens})
			}
			if len(entries) > 0 {
				if _, err := txn.Write(ctx, economy.Intent{
					Kind: "world_event_reward", RefKind: "world_event_spawn", RefID: d.spawn,
					IdempotencyKey: "worldevent-" + d.spawn.String() + "-" + d.char.String(),
					Entries:        entries,
				}); err != nil {
					log.ErrorContext(ctx, "worldevent: reward failed", "spawn", d.spawn.String(), "err", err.Error())
					continue
				}
			}
		}
		if _, err := pool.Exec(ctx, `
			UPDATE worldevent.attendance SET rewarded_at = now() WHERE spawn_id = $1 AND character_id = $2`,
			d.spawn, d.char); err != nil {
			return n, err
		}
		n++
	}
	return n, nil
}
