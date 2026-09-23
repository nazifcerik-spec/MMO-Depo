package progression

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/otter/server/internal/economy"
	"github.com/otter/server/internal/notification"
)

type BonusEffect struct {
	Kind   string  `json:"kind"`
	Target string  `json:"target"`
	Amount float64 `json:"amount"`
}

type tierBonus struct {
	Tier    int16         `json:"tier"`
	Effects []BonusEffect `json:"bonus_effects"`
}

type track struct {
	Code, NameKey, Category string
	MaxTier                 int16
	Curve                   []int64
	Bonuses                 []tierBonus
	Sources                 map[string]int64 // event kind → xp per unit
}

type achievement struct {
	Code, NameKey, DescKey, Category, Trigger string
	Points                                    int16
	Hidden                                    bool
	Target                                    int32
	RewardTitle                               *string
	RewardGold                                int64
}

type faction struct {
	Code, NameKey string
	DailyCap      int
	Sources       map[string]int
}

// Catalog is the immutable progression content loaded at boot.
type Catalog struct {
	tracks       []track
	achievements []achievement
	factions     []faction
	titles       map[string]string // code → display key
}

func LoadCatalog(ctx context.Context, pool *pgxpool.Pool) (*Catalog, error) {
	c := &Catalog{titles: map[string]string{}}
	rows, err := pool.Query(ctx, `SELECT code, display_name_key, category, max_tier, tier_xp_curve, tier_bonuses, sources
	                                FROM content.mastery_track ORDER BY code`)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var t track
		var curve, bonuses, sources []byte
		if err := rows.Scan(&t.Code, &t.NameKey, &t.Category, &t.MaxTier, &curve, &bonuses, &sources); err != nil {
			rows.Close()
			return nil, err
		}
		var src []struct {
			Kind      string `json:"kind"`
			XPPerUnit int64  `json:"xp_per_unit"`
		}
		if err := json.Unmarshal(curve, &t.Curve); err != nil {
			rows.Close()
			return nil, fmt.Errorf("mastery %s curve: %w", t.Code, err)
		}
		if err := json.Unmarshal(bonuses, &t.Bonuses); err != nil {
			rows.Close()
			return nil, fmt.Errorf("mastery %s bonuses: %w", t.Code, err)
		}
		if err := json.Unmarshal(sources, &src); err != nil {
			rows.Close()
			return nil, fmt.Errorf("mastery %s sources: %w", t.Code, err)
		}
		t.Sources = map[string]int64{}
		for _, s := range src {
			t.Sources[s.Kind] = s.XPPerUnit
		}
		c.tracks = append(c.tracks, t)
	}
	rows.Close()

	rows, err = pool.Query(ctx, `SELECT code, display_name_key, description_key, category, point_value, is_hidden,
	                                    trigger_kind, progress_target, reward_title_code, reward_gold_copper
	                               FROM content.achievement_template ORDER BY category, progress_target`)
	if err != nil {
		return nil, err
	}
	c.achievements, err = pgx.CollectRows(rows, func(r pgx.CollectableRow) (achievement, error) {
		var a achievement
		err := r.Scan(&a.Code, &a.NameKey, &a.DescKey, &a.Category, &a.Points, &a.Hidden,
			&a.Trigger, &a.Target, &a.RewardTitle, &a.RewardGold)
		return a, err
	})
	if err != nil {
		return nil, err
	}

	rows, err = pool.Query(ctx, `SELECT code, display_name_key, daily_cap, sources FROM content.faction ORDER BY code`)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var f faction
		var sources []byte
		if err := rows.Scan(&f.Code, &f.NameKey, &f.DailyCap, &sources); err != nil {
			rows.Close()
			return nil, err
		}
		var src []struct {
			Kind  string `json:"kind"`
			Delta int    `json:"delta"`
		}
		if err := json.Unmarshal(sources, &src); err != nil {
			rows.Close()
			return nil, fmt.Errorf("faction %s sources: %w", f.Code, err)
		}
		f.Sources = map[string]int{}
		for _, s := range src {
			f.Sources[s.Kind] = s.Delta
		}
		c.factions = append(c.factions, f)
	}
	rows.Close()

	rows, err = pool.Query(ctx, `SELECT code, display_key FROM content.title_template`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var code, key string
		if err := rows.Scan(&code, &key); err != nil {
			return nil, err
		}
		c.titles[code] = key
	}
	return c, rows.Err()
}

type Owner interface {
	Owns(ctx context.Context, characterID, accountID uuid.UUID) error
}

type EventPublisher interface {
	Publish(accountID uuid.UUID, kind string, payload any) int
}

type Service struct {
	pool   *pgxpool.Pool
	cat    *Catalog
	chars  Owner
	txn    economy.TxnService
	pub    EventPublisher
	notify notification.Notifier
	now    func() time.Time
	log    *slog.Logger
}

func NewService(pool *pgxpool.Pool, cat *Catalog, chars Owner, txn economy.TxnService, logger *slog.Logger) *Service {
	if logger == nil {
		logger = slog.Default()
	}
	return &Service{pool: pool, cat: cat, chars: chars, txn: txn, now: time.Now, log: logger}
}

func (s *Service) WithPublisher(p EventPublisher) *Service       { s.pub = p; return s }
func (s *Service) WithNotifier(n notification.Notifier) *Service { s.notify = n; return s }

// Unlock describes a newly unlocked achievement.
type Unlock struct {
	Code       string  `json:"code"`
	Points     int16   `json:"points"`
	TitleCode  *string `json:"title_code,omitempty"`
	RewardGold int64   `json:"reward_gold_copper"`
}

// Record feeds events into mastery, achievements and reputation. It never
// fails the caller's game action: errors are logged.
func (s *Service) Record(ctx context.Context, characterID uuid.UUID, events ...Event) {
	var unlocks []Unlock
	var tierUps []map[string]any
	err := pgx.BeginFunc(ctx, s.pool, func(tx pgx.Tx) error {
		var err error
		if tierUps, err = s.recordMastery(ctx, tx, characterID, events); err != nil {
			return err
		}
		if unlocks, err = s.recordAchievements(ctx, tx, characterID, events); err != nil {
			return err
		}
		return s.recordReputation(ctx, tx, characterID, events)
	})
	if err != nil {
		s.log.WarnContext(ctx, "progression: record failed", "character_id", characterID.String(), "err", err.Error())
		return
	}
	if len(unlocks) == 0 && len(tierUps) == 0 {
		return
	}
	var account uuid.UUID
	if err := s.pool.QueryRow(ctx, `SELECT account_id FROM character.character WHERE id = $1`, characterID).Scan(&account); err != nil {
		return
	}
	for _, u := range unlocks {
		if u.RewardGold > 0 {
			if _, err := s.txn.Write(ctx, economy.Intent{
				Kind: "achievement_reward", RefKind: "achievement",
				IdempotencyKey: "achievement-reward-" + characterID.String() + "-" + u.Code,
				Entries: []economy.Entry{
					{Currency: "gold", AccountKind: economy.AccountSystem, AccountRef: economy.SystemRewardSource, Delta: -u.RewardGold},
					{Currency: "gold", AccountKind: economy.AccountCharacter, AccountRef: characterID, Delta: u.RewardGold},
				},
			}); err != nil {
				s.log.ErrorContext(ctx, "progression: achievement gold reward failed", "code", u.Code, "err", err.Error())
			}
		}
		if s.pub != nil {
			s.pub.Publish(account, "achievement.unlocked", u)
		}
		ch := characterID
		notification.Send(ctx, s.notify, s.log, notification.Notice{
			AccountID: account, CharacterID: &ch, Kind: "achievement_unlocked",
			DedupKey: "achievement_unlocked:" + characterID.String() + ":" + u.Code,
			TitleKey: "notif.achievement_unlocked.title", BodyKey: "notif.achievement_unlocked.body",
			Vars: map[string]any{"code": u.Code, "points": u.Points}, Route: "character/achievements",
			Category: "system",
		})
	}
	if s.pub != nil {
		for _, t := range tierUps {
			s.pub.Publish(account, "mastery.tier_up", t)
		}
	}
}

func (s *Service) recordMastery(ctx context.Context, tx pgx.Tx, char uuid.UUID, events []Event) ([]map[string]any, error) {
	var ups []map[string]any
	for _, t := range s.cat.tracks {
		var gain int64
		for _, e := range events {
			gain += t.Sources[e.Kind] * e.Amount
		}
		if gain <= 0 {
			continue
		}
		var xp int64
		var tier int16
		if err := tx.QueryRow(ctx, `
			INSERT INTO progression.character_mastery (character_id, track_code, xp, last_gain_at)
			VALUES ($1, $2, $3, now())
			ON CONFLICT (character_id, track_code) DO UPDATE
			  SET xp = progression.character_mastery.xp + EXCLUDED.xp, last_gain_at = now()
			RETURNING xp, current_tier`, char, t.Code, gain).Scan(&xp, &tier); err != nil {
			return nil, err
		}
		if next := MasteryTier(xp, t.Curve, t.MaxTier); next > tier {
			if _, err := tx.Exec(ctx, `
				UPDATE progression.character_mastery SET current_tier = $3
				 WHERE character_id = $1 AND track_code = $2 AND current_tier < $3`, char, t.Code, next); err != nil {
				return nil, err
			}
			ups = append(ups, map[string]any{"track": t.Code, "tier": next})
		}
	}
	return ups, nil
}

func (s *Service) recordAchievements(ctx context.Context, tx pgx.Tx, char uuid.UUID, events []Event) ([]Unlock, error) {
	var unlocks []Unlock
	for _, a := range s.cat.achievements {
		var gain int64
		for _, e := range events {
			if e.Kind == a.Trigger {
				gain += e.Amount
			}
		}
		if gain <= 0 {
			continue
		}
		if gain > int64(a.Target) {
			gain = int64(a.Target)
		}
		tag, err := tx.Exec(ctx, `
			INSERT INTO progression.achievement_progress (character_id, achievement_code, progress_current)
			VALUES ($1, $2, LEAST($3::int, $4::int))
			ON CONFLICT (character_id, achievement_code) DO UPDATE
			  SET progress_current = LEAST(progression.achievement_progress.progress_current + $3::int, $4::int)
			WHERE progression.achievement_progress.unlocked_at IS NULL`, char, a.Code, gain, a.Target)
		if err != nil {
			return nil, err
		}
		if tag.RowsAffected() == 0 {
			continue // already unlocked
		}
		tag, err = tx.Exec(ctx, `
			UPDATE progression.achievement_progress SET unlocked_at = now()
			 WHERE character_id = $1 AND achievement_code = $2 AND unlocked_at IS NULL AND progress_current >= $3`,
			char, a.Code, a.Target)
		if err != nil {
			return nil, err
		}
		if tag.RowsAffected() == 0 {
			continue
		}
		if _, err := tx.Exec(ctx, `
			UPDATE character.character SET achievement_points = achievement_points + $2 WHERE id = $1`,
			char, a.Points); err != nil {
			return nil, err
		}
		if a.RewardTitle != nil {
			if _, err := tx.Exec(ctx, `
				INSERT INTO progression.character_title (character_id, title_code, source_ref_kind)
				VALUES ($1, $2, 'achievement') ON CONFLICT DO NOTHING`, char, *a.RewardTitle); err != nil {
				return nil, err
			}
		}
		unlocks = append(unlocks, Unlock{Code: a.Code, Points: a.Points, TitleCode: a.RewardTitle, RewardGold: a.RewardGold})
	}
	return unlocks, nil
}

func (s *Service) recordReputation(ctx context.Context, tx pgx.Tx, char uuid.UUID, events []Event) error {
	now := s.now().UTC()
	for _, f := range s.cat.factions {
		delta := 0
		for _, e := range events {
			delta += f.Sources[e.Kind] * int(e.Amount)
		}
		if delta == 0 {
			continue
		}
		var rep, used int
		var resetOn time.Time
		if _, err := tx.Exec(ctx, `
			INSERT INTO progression.character_reputation (character_id, faction_code)
			VALUES ($1, $2) ON CONFLICT DO NOTHING`, char, f.Code); err != nil {
			return err
		}
		if err := tx.QueryRow(ctx, `
			SELECT rep_xp, daily_gain_used, daily_reset_on FROM progression.character_reputation
			 WHERE character_id = $1 AND faction_code = $2 FOR UPDATE`, char, f.Code).Scan(&rep, &used, &resetOn); err != nil {
			return err
		}
		if !SameUTCDay(resetOn, now) {
			used = 0
		}
		rep, used = ApplyRepGain(rep, used, delta, f.DailyCap)
		if _, err := tx.Exec(ctx, `
			UPDATE progression.character_reputation
			   SET rep_xp = $3, current_tier = $4, daily_gain_used = $5, daily_reset_on = $6, last_gain_at = now()
			 WHERE character_id = $1 AND faction_code = $2`,
			char, f.Code, rep, RepTier(rep), used, now); err != nil {
			return err
		}
	}
	return nil
}

// --- reads ------------------------------------------------------------------

type MasteryView struct {
	Track       string        `json:"track"`
	NameKey     string        `json:"display_name_key"`
	Category    string        `json:"category"`
	XP          int64         `json:"xp"`
	Tier        int16         `json:"tier"`
	MaxTier     int16         `json:"max_tier"`
	NextTierXP  int64         `json:"next_tier_xp"`
	ActiveBonus []BonusEffect `json:"active_bonuses"`
}

type ReputationView struct {
	Faction   string `json:"faction"`
	NameKey   string `json:"display_name_key"`
	Rep       int    `json:"rep_xp"`
	Tier      int16  `json:"tier"`
	TierName  string `json:"tier_name"`
	DailyLeft int    `json:"daily_gain_left"`
}

type Summary struct {
	AchievementPoints int              `json:"achievement_points"`
	ActiveTitle       *string          `json:"active_title_code"`
	Mastery           []MasteryView    `json:"mastery"`
	Reputation        []ReputationView `json:"reputation"`
}

func (s *Service) Summary(ctx context.Context, char, account uuid.UUID) (*Summary, error) {
	if err := s.chars.Owns(ctx, char, account); err != nil {
		return nil, err
	}
	sum := &Summary{}
	if err := s.pool.QueryRow(ctx, `
		SELECT achievement_points, active_title_code FROM character.character WHERE id = $1`, char).
		Scan(&sum.AchievementPoints, &sum.ActiveTitle); err != nil {
		return nil, err
	}
	xp := map[string]int64{}
	rows, err := s.pool.Query(ctx, `SELECT track_code, xp FROM progression.character_mastery WHERE character_id = $1`, char)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var code string
		var v int64
		if err := rows.Scan(&code, &v); err != nil {
			rows.Close()
			return nil, err
		}
		xp[code] = v
	}
	rows.Close()
	for _, t := range s.cat.tracks {
		tier := MasteryTier(xp[t.Code], t.Curve, t.MaxTier)
		v := MasteryView{Track: t.Code, NameKey: t.NameKey, Category: t.Category, XP: xp[t.Code],
			Tier: tier, MaxTier: t.MaxTier, NextTierXP: NextTierXP(tier, t.Curve, t.MaxTier), ActiveBonus: []BonusEffect{}}
		for _, b := range t.Bonuses {
			if b.Tier <= tier {
				v.ActiveBonus = append(v.ActiveBonus, b.Effects...)
			}
		}
		sum.Mastery = append(sum.Mastery, v)
	}

	type repRow struct {
		rep, used int
		resetOn   time.Time
	}
	reps := map[string]repRow{}
	rows, err = s.pool.Query(ctx, `
		SELECT faction_code, rep_xp, daily_gain_used, daily_reset_on FROM progression.character_reputation
		 WHERE character_id = $1`, char)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var code string
		var r repRow
		if err := rows.Scan(&code, &r.rep, &r.used, &r.resetOn); err != nil {
			rows.Close()
			return nil, err
		}
		reps[code] = r
	}
	rows.Close()
	now := s.now()
	for _, f := range s.cat.factions {
		r := reps[f.Code]
		used := r.used
		if !SameUTCDay(r.resetOn, now) {
			used = 0
		}
		tier := RepTier(r.rep)
		sum.Reputation = append(sum.Reputation, ReputationView{
			Faction: f.Code, NameKey: f.NameKey, Rep: r.rep, Tier: tier, TierName: RepTierNames[tier],
			DailyLeft: f.DailyCap - used,
		})
	}
	return sum, nil
}

type AchievementView struct {
	Code       string     `json:"code"`
	NameKey    string     `json:"display_name_key"`
	DescKey    string     `json:"description_key"`
	Category   string     `json:"category"`
	Points     int16      `json:"points"`
	Progress   int32      `json:"progress"`
	Target     int32      `json:"target"`
	UnlockedAt *time.Time `json:"unlocked_at,omitempty"`
	TitleCode  *string    `json:"reward_title_code,omitempty"`
	RewardGold int64      `json:"reward_gold_copper"`
}

func (s *Service) Achievements(ctx context.Context, char, account uuid.UUID) ([]AchievementView, error) {
	if err := s.chars.Owns(ctx, char, account); err != nil {
		return nil, err
	}
	type prog struct {
		cur int32
		at  *time.Time
	}
	state := map[string]prog{}
	rows, err := s.pool.Query(ctx, `
		SELECT achievement_code, progress_current, unlocked_at FROM progression.achievement_progress
		 WHERE character_id = $1`, char)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var code string
		var p prog
		if err := rows.Scan(&code, &p.cur, &p.at); err != nil {
			rows.Close()
			return nil, err
		}
		state[code] = p
	}
	rows.Close()
	out := []AchievementView{}
	for _, a := range s.cat.achievements {
		p := state[a.Code]
		if a.Hidden && p.at == nil {
			continue
		}
		out = append(out, AchievementView{Code: a.Code, NameKey: a.NameKey, DescKey: a.DescKey, Category: a.Category,
			Points: a.Points, Progress: p.cur, Target: a.Target, UnlockedAt: p.at, TitleCode: a.RewardTitle,
			RewardGold: a.RewardGold})
	}
	return out, nil
}

type TitleView struct {
	Code       string    `json:"code"`
	DisplayKey string    `json:"display_key"`
	UnlockedAt time.Time `json:"unlocked_at"`
	Active     bool      `json:"active"`
}

func (s *Service) Titles(ctx context.Context, char, account uuid.UUID) ([]TitleView, error) {
	if err := s.chars.Owns(ctx, char, account); err != nil {
		return nil, err
	}
	rows, err := s.pool.Query(ctx, `
		SELECT t.title_code, t.unlocked_at, COALESCE(c.active_title_code = t.title_code, false)
		  FROM progression.character_title t JOIN character.character c ON c.id = t.character_id
		 WHERE t.character_id = $1 ORDER BY t.unlocked_at`, char)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, func(r pgx.CollectableRow) (TitleView, error) {
		var t TitleView
		err := r.Scan(&t.Code, &t.UnlockedAt, &t.Active)
		t.DisplayKey = s.cat.titles[t.Code]
		return t, err
	})
}

// SetActiveTitle equips an unlocked title; empty code clears it.
func (s *Service) SetActiveTitle(ctx context.Context, char, account uuid.UUID, code string) error {
	if err := s.chars.Owns(ctx, char, account); err != nil {
		return err
	}
	code = strings.TrimSpace(code)
	if code == "" {
		_, err := s.pool.Exec(ctx, `UPDATE character.character SET active_title_code = NULL WHERE id = $1`, char)
		return err
	}
	tag, err := s.pool.Exec(ctx, `
		UPDATE character.character SET active_title_code = $2
		 WHERE id = $1 AND EXISTS (SELECT 1 FROM progression.character_title WHERE character_id = $1 AND title_code = $2)`,
		char, code)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrTitleNotOwned
	}
	return nil
}
