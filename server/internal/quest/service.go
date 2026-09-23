package quest

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/otter/server/internal/domain/character"
	"github.com/otter/server/internal/economy"
	"github.com/otter/server/internal/progression"
)

var (
	ErrNotFound    = errors.New("quest: not found")
	ErrLevel       = errors.New("quest: level too low")
	ErrLocked      = errors.New("quest: prerequisite not done")
	ErrTaken       = errors.New("quest: already taken")
	ErrNotOffered  = errors.New("quest: daily not offered today")
	ErrOpenCap     = errors.New("quest: too many open quests")
	ErrNotActive   = errors.New("quest: not open")
	ErrIncomplete  = errors.New("quest: objectives not complete")
	ErrNotPlayable = errors.New("quest: character not playable")
)

// CharacterFetcher is the slice of character.Service quests need.
type CharacterFetcher interface {
	GetOwned(ctx context.Context, characterID, accountID uuid.UUID) (*character.Character, error)
}

// XPGranter grants character XP (character.Service).
type XPGranter interface {
	GrantXP(ctx context.Context, characterID uuid.UUID, delta int64) (*character.XPGrantResult, error)
}

// EventPublisher pushes WS events (wsgateway adapter); nil-safe.
type EventPublisher interface {
	Publish(accountID uuid.UUID, kind string, payload any) int
}

type Service struct {
	pool *pgxpool.Pool
	cat  atomic.Pointer[Catalog]
	char CharacterFetcher
	txn  economy.TxnService
	xp   XPGranter
	pub  EventPublisher
	now  func() time.Time
	log  *slog.Logger
}

func NewService(pool *pgxpool.Pool, chars CharacterFetcher, txn economy.TxnService, xp XPGranter, logger *slog.Logger) *Service {
	if logger == nil {
		logger = slog.Default()
	}
	s := &Service{pool: pool, char: chars, txn: txn, xp: xp, now: time.Now, log: logger}
	s.cat.Store(NewCatalog(nil))
	return s
}

func (s *Service) WithPublisher(p EventPublisher) *Service { s.pub = p; return s }

// WithClock injects the clock (tests: day rollover).
func (s *Service) WithClock(now func() time.Time) *Service { s.now = now; return s }

// Reload rebuilds the catalog from content.quest_template (contentpack.Reloader).
func (s *Service) Reload(ctx context.Context) error {
	rows, err := s.pool.Query(ctx, `
		SELECT code, display_name_key, category, min_level, COALESCE(requires_quest_code, ''),
		       objectives, reward_xp, reward_gold_copper, sort_order
		  FROM content.quest_template`)
	if err != nil {
		return err
	}
	list, err := pgx.CollectRows(rows, func(r pgx.CollectableRow) (Template, error) {
		var t Template
		var obj []byte
		if err := r.Scan(&t.Code, &t.NameKey, &t.Category, &t.MinLevel, &t.Requires,
			&obj, &t.RewardXP, &t.RewardGold, &t.SortOrder); err != nil {
			return t, err
		}
		return t, json.Unmarshal(obj, &t.Objectives)
	})
	if err != nil {
		return err
	}
	s.cat.Store(NewCatalog(list))
	return nil
}

// --- views -----------------------------------------------------------------

type ObjectiveView struct {
	Event    string `json:"event"`
	Count    int64  `json:"count"`
	Progress int64  `json:"progress"`
}

type View struct {
	Code       string          `json:"code"`
	NameKey    string          `json:"display_name_key"`
	Category   string          `json:"category"`
	MinLevel   int16           `json:"min_level"`
	State      string          `json:"state"` // available | active | completed
	Objectives []ObjectiveView `json:"objectives"`
	RewardXP   int64           `json:"reward_xp"`
	RewardGold int64           `json:"reward_gold_copper"`
	ExpiresAt  *time.Time      `json:"expires_at,omitempty"` // dailies: next UTC midnight
}

type Board struct {
	Open          []View    `json:"open"`
	Available     []View    `json:"available"`
	OpenCap       int       `json:"open_cap"`
	DailyResetsAt time.Time `json:"daily_resets_at"`
}

func view(t *Template, state string, progress []int64, today time.Time) View {
	v := View{Code: t.Code, NameKey: t.NameKey, Category: t.Category, MinLevel: t.MinLevel, State: state,
		RewardXP: t.RewardXP, RewardGold: t.RewardGold, Objectives: make([]ObjectiveView, len(t.Objectives))}
	for i, o := range t.Objectives {
		v.Objectives[i] = ObjectiveView{Event: o.Event, Count: o.Count}
		if i < len(progress) {
			v.Objectives[i].Progress = progress[i]
		}
	}
	if t.Category == CategoryDaily {
		end := today.AddDate(0, 0, 1)
		v.ExpiresAt = &end
	}
	return v
}

// openRow is one open (active/completed) quest valid today.
type openRow struct {
	id       uuid.UUID
	code     string
	state    string
	progress []int64
}

// q is satisfied by both the pool and a tx.
type q interface {
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
}

func openRows(ctx context.Context, db q, char uuid.UUID, today time.Time, lock bool) ([]openRow, error) {
	sql := `SELECT id, quest_code, state, progress FROM progression.character_quest
	         WHERE character_id = $1 AND state IN ('active','completed')
	           AND (day_key IS NULL OR day_key = $2)`
	if lock {
		sql += ` FOR UPDATE`
	}
	rows, err := db.Query(ctx, sql, char, today)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, func(r pgx.CollectableRow) (openRow, error) {
		var o openRow
		var raw []byte
		if err := r.Scan(&o.id, &o.code, &o.state, &raw); err != nil {
			return o, err
		}
		return o, json.Unmarshal(raw, &o.progress)
	})
}

// List returns the character's open quests plus what can be accepted now.
func (s *Service) List(ctx context.Context, char, account uuid.UUID) (*Board, error) {
	c, err := s.char.GetOwned(ctx, char, account)
	if err != nil {
		return nil, err
	}
	cat := s.cat.Load()
	today := DayKey(s.now())
	open, err := openRows(ctx, s.pool, char, today, false)
	if err != nil {
		return nil, err
	}
	// Quests closed for good (main/side turned in) or taken today (dailies).
	rows, err := s.pool.Query(ctx, `
		SELECT quest_code FROM progression.character_quest
		 WHERE character_id = $1 AND ((state = 'turned_in' AND day_key IS NULL) OR day_key = $2)`, char, today)
	if err != nil {
		return nil, err
	}
	closedList, err := pgx.CollectRows(rows, pgx.RowTo[string])
	if err != nil {
		return nil, err
	}
	done := map[string]bool{}
	for _, code := range closedList {
		done[code] = true
	}
	b := &Board{Open: []View{}, Available: []View{}, OpenCap: OpenCap, DailyResetsAt: today.AddDate(0, 0, 1)}
	taken := map[string]bool{}
	for _, o := range open {
		taken[o.code] = true
		if t, ok := cat.Get(o.code); ok {
			b.Open = append(b.Open, view(t, o.state, o.progress, today))
		}
	}
	offered := map[string]bool{}
	for _, code := range cat.DailyOffer(char, c.Level, today) {
		offered[code] = true
	}
	for i := range cat.list {
		t := &cat.list[i]
		if taken[t.Code] || done[t.Code] || t.MinLevel > c.Level {
			continue
		}
		if t.Category == CategoryDaily && !offered[t.Code] {
			continue
		}
		if t.Requires != "" && !done[t.Requires] { // prerequisites are never dailies
			continue
		}
		b.Available = append(b.Available, view(t, "available", nil, today))
	}
	return b, nil
}

// Accept starts a quest. Progress starts at zero (no backfill).
func (s *Service) Accept(ctx context.Context, char, account uuid.UUID, code string) (*View, error) {
	c, err := s.char.GetOwned(ctx, char, account)
	if err != nil {
		return nil, err
	}
	if c.State != character.StateActive {
		return nil, ErrNotPlayable
	}
	cat := s.cat.Load()
	t, ok := cat.Get(code)
	if !ok {
		return nil, ErrNotFound
	}
	if t.MinLevel > c.Level {
		return nil, ErrLevel
	}
	now := s.now().UTC()
	today := DayKey(now)
	var dayKey *time.Time
	if t.Category == CategoryDaily {
		offered := false
		for _, o := range cat.DailyOffer(char, c.Level, today) {
			offered = offered || o == code
		}
		if !offered {
			return nil, ErrNotOffered
		}
		dayKey = &today
	}
	progress := make([]int64, len(t.Objectives))
	raw, _ := json.Marshal(progress)
	err = pgx.BeginFunc(ctx, s.pool, func(tx pgx.Tx) error {
		// Past-day dailies left open expire now (missed daily = missed reward).
		if _, err := tx.Exec(ctx, `
			UPDATE progression.character_quest SET state = 'expired', closed_at = $3
			 WHERE character_id = $1 AND day_key < $2 AND state IN ('active','completed')`, char, today, now); err != nil {
			return err
		}
		var open int
		if err := tx.QueryRow(ctx, `
			SELECT count(*) FROM progression.character_quest
			 WHERE character_id = $1 AND state IN ('active','completed')`, char).Scan(&open); err != nil {
			return err
		}
		if open >= OpenCap {
			return ErrOpenCap
		}
		if t.Category != CategoryDaily {
			var done bool
			if err := tx.QueryRow(ctx, `
				SELECT EXISTS (SELECT 1 FROM progression.character_quest
				                WHERE character_id = $1 AND quest_code = $2 AND state = 'turned_in')`,
				char, code).Scan(&done); err != nil {
				return err
			}
			if done {
				return ErrTaken
			}
		}
		if t.Requires != "" {
			var ok bool
			if err := tx.QueryRow(ctx, `
				SELECT EXISTS (SELECT 1 FROM progression.character_quest
				                WHERE character_id = $1 AND quest_code = $2 AND state = 'turned_in')`,
				char, t.Requires).Scan(&ok); err != nil {
				return err
			}
			if !ok {
				return ErrLocked
			}
		}
		_, err := tx.Exec(ctx, `
			INSERT INTO progression.character_quest (id, character_id, quest_code, state, progress, day_key, accepted_at)
			VALUES ($1, $2, $3, 'active', $4, $5, $6)`, uuid.New(), char, code, raw, dayKey, now)
		if isUnique(err) {
			return ErrTaken
		}
		return err
	})
	if err != nil {
		return nil, err
	}
	v := view(t, "active", progress, today)
	return &v, nil
}

// Abandon drops an open quest (a daily stays used for today).
func (s *Service) Abandon(ctx context.Context, char, account uuid.UUID, code string) error {
	if _, err := s.char.GetOwned(ctx, char, account); err != nil {
		return err
	}
	now := s.now().UTC()
	tag, err := s.pool.Exec(ctx, `
		UPDATE progression.character_quest SET state = 'abandoned', closed_at = $4
		 WHERE character_id = $1 AND quest_code = $2 AND state IN ('active','completed')
		   AND (day_key IS NULL OR day_key = $3)`, char, code, DayKey(now), now)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotActive
	}
	return nil
}

type TurnInResult struct {
	Code       string                   `json:"code"`
	RewardXP   int64                    `json:"reward_xp"`
	RewardGold int64                    `json:"reward_gold_copper"`
	XP         *character.XPGrantResult `json:"xp,omitempty"`
}

// TurnIn pays a completed quest. Order is chosen so a retry can never pay
// twice: idempotent ledger write → conditional state flip → XP only for the
// request that flipped it.
func (s *Service) TurnIn(ctx context.Context, char, account uuid.UUID, code string) (*TurnInResult, error) {
	if _, err := s.char.GetOwned(ctx, char, account); err != nil {
		return nil, err
	}
	t, ok := s.cat.Load().Get(code)
	if !ok {
		return nil, ErrNotFound
	}
	now := s.now().UTC()
	var id uuid.UUID
	var state string
	err := s.pool.QueryRow(ctx, `
		SELECT id, state FROM progression.character_quest
		 WHERE character_id = $1 AND quest_code = $2 AND state IN ('active','completed')
		   AND (day_key IS NULL OR day_key = $3)`, char, code, DayKey(now)).Scan(&id, &state)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotActive
	}
	if err != nil {
		return nil, err
	}
	if state != "completed" {
		return nil, ErrIncomplete
	}
	var txnID *uuid.UUID
	if t.RewardGold > 0 {
		res, err := s.txn.Write(ctx, economy.Intent{
			Kind: "quest_reward", RefKind: "quest", RefID: id,
			IdempotencyKey: "quest-reward-" + id.String(),
			Entries: []economy.Entry{
				{Currency: "gold", AccountKind: economy.AccountSystem, AccountRef: economy.SystemRewardSource, Delta: -t.RewardGold},
				{Currency: "gold", AccountKind: economy.AccountCharacter, AccountRef: char, Delta: t.RewardGold},
			},
		})
		if err != nil {
			return nil, fmt.Errorf("quest: reward ledger: %w", err)
		}
		txnID = &res.TxnID
	}
	tag, err := s.pool.Exec(ctx, `
		UPDATE progression.character_quest SET state = 'turned_in', closed_at = $2, reward_txn_id = $3
		 WHERE id = $1 AND state = 'completed'`, id, now, txnID)
	if err != nil {
		return nil, err
	}
	if tag.RowsAffected() == 0 {
		return nil, ErrNotActive // a concurrent turn-in won
	}
	out := &TurnInResult{Code: code, RewardXP: t.RewardXP, RewardGold: t.RewardGold}
	if t.RewardXP > 0 && s.xp != nil {
		if out.XP, err = s.xp.GrantXP(ctx, char, t.RewardXP); err != nil {
			s.log.ErrorContext(ctx, "quest: xp grant failed after turn-in (not retried)",
				slog.String("quest_id", id.String()), slog.String("err", err.Error()))
		}
	}
	return out, nil
}

// Record ticks open quests from game events (progression.Recorder). Like
// progression it never fails the caller's action: errors are logged.
func (s *Service) Record(ctx context.Context, char uuid.UUID, events ...progression.Event) {
	cat := s.cat.Load()
	today := DayKey(s.now())
	var completed []string
	err := pgx.BeginFunc(ctx, s.pool, func(tx pgx.Tx) error {
		open, err := openRows(ctx, tx, char, today, true)
		if err != nil {
			return err
		}
		for _, o := range open {
			t, ok := cat.Get(o.code)
			if !ok || o.state != "active" {
				continue
			}
			next, changed, done := Advance(t, o.progress, events)
			if !changed {
				continue
			}
			raw, _ := json.Marshal(next)
			state, completedAt := "active", (*time.Time)(nil)
			if done {
				now := s.now().UTC()
				state, completedAt = "completed", &now
				completed = append(completed, o.code)
			}
			if _, err := tx.Exec(ctx, `
				UPDATE progression.character_quest SET progress = $2, state = $3, completed_at = $4 WHERE id = $1`,
				o.id, raw, state, completedAt); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		s.log.WarnContext(ctx, "quest: record failed", "character_id", char.String(), "err", err.Error())
		return
	}
	if len(completed) == 0 || s.pub == nil {
		return
	}
	var account uuid.UUID
	if err := s.pool.QueryRow(ctx, `SELECT account_id FROM character.character WHERE id = $1`, char).Scan(&account); err != nil {
		return
	}
	for _, code := range completed {
		s.pub.Publish(account, "quest.completed", map[string]any{"character_id": char, "code": code})
	}
}

func isUnique(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "23505"
}
