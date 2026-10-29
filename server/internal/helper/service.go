package helper

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/otter/server/internal/activity"
	"github.com/otter/server/internal/crafting"
	"github.com/otter/server/internal/economy"
	"github.com/otter/server/internal/inventory"
	"github.com/otter/server/internal/item"
)

// Gathering is the character's gathering skill (gathering.Service).
type Gathering interface {
	Nodes(kind activity.Kind) []activity.GatherNode
	Level(ctx context.Context, characterID uuid.UUID, kind activity.Kind) (int, error)
	AddXP(ctx context.Context, characterID uuid.UUID, kind activity.Kind, xp int) (activity.GatherProgress, error)
}

// Items grants gathered drops and crafted outputs (item.Service).
type Items interface {
	Grant(ctx context.Context, characterID, sourceRefID uuid.UUID, seed int64, d activity.Drops) ([]activity.GrantedItem, error)
	MintExact(ctx context.Context, characterID uuid.UUID, templateCode string, qty int32, sourceKind string,
		sourceRefID uuid.UUID) (*item.Instance, error)
}

// Bag takes a craft's materials (inventory.Service).
type Bag interface {
	ConsumeAll(ctx context.Context, characterID uuid.UUID, needs []inventory.Need) error
	PlaceOverflow(ctx context.Context, characterID uuid.UUID) error
}

// Crafting is the recipe book and the character's crafting skill.
type Crafting interface {
	Catalog() *crafting.Catalog
}

// Professions reads and raises crafting profession levels.
type Professions interface {
	Get(ctx context.Context, characterID uuid.UUID, professionCode string) (*crafting.CharacterProfession, error)
}

// CraftXP pays crafting profession XP (crafting.Service).
type CraftXP interface {
	GrantXP(ctx context.Context, characterID uuid.UUID, profession string, xp int32) error
}

// Characters owns and levels characters (character.Service).
type Characters interface {
	Owns(ctx context.Context, characterID, accountID uuid.UUID) error
	Level(ctx context.Context, characterID uuid.UUID) (int, error)
}

type Service struct {
	pool    *pgxpool.Pool
	txn     economy.TxnService
	chars   Characters
	gather  Gathering
	items   Items
	bag     Bag
	recipes Crafting
	profs   Professions
	craftXP CraftXP
	log     *slog.Logger
	now     func() time.Time
}

func NewService(pool *pgxpool.Pool, txn economy.TxnService, chars Characters, gather Gathering, items Items,
	bag Bag, recipes Crafting, profs Professions, craftXP CraftXP, logger *slog.Logger) *Service {
	if logger == nil {
		logger = slog.Default()
	}
	return &Service{pool: pool, txn: txn, chars: chars, gather: gather, items: items, bag: bag, recipes: recipes,
		profs: profs, craftXP: craftXP, log: logger, now: time.Now}
}

// WithClock is for tests.
func (s *Service) WithClock(now func() time.Time) *Service { s.now = now; return s }

// Spec is what a job does: a gathering run or a batch of crafts.
type Spec struct {
	Target   string `json:"target,omitempty"` // gathering resource
	Tier     int    `json:"tier,omitempty"`
	Duration int32  `json:"duration_seconds,omitempty"`
	Recipe   string `json:"recipe,omitempty"` // crafting
	Count    int    `json:"count,omitempty"`
	Gold     int64  `json:"gold_copper,omitempty"` // what the craft cost
}

// Job is a helper's current work.
type Job struct {
	ID        uuid.UUID `json:"id"`
	Spec      Spec      `json:"spec"`
	StartedAt time.Time `json:"started_at"`
	ReadyAt   time.Time `json:"ready_at"`
	seed      int64
}

// HelperView is one helper type for the character: hired, working.
type HelperView struct {
	Type
	Hired bool `json:"hired"`
	Job   *Job `json:"job,omitempty"`
}

// List is every helper type with the character's hires and jobs.
func (s *Service) List(ctx context.Context, char, account uuid.UUID) ([]HelperView, error) {
	if err := s.chars.Owns(ctx, char, account); err != nil {
		return nil, err
	}
	hired := map[string]bool{}
	rows, err := s.pool.Query(ctx, `SELECT type_code FROM helper.hired WHERE character_id = $1`, char)
	if err != nil {
		return nil, err
	}
	codes, err := pgx.CollectRows(rows, pgx.RowTo[string])
	if err != nil {
		return nil, err
	}
	for _, c := range codes {
		hired[c] = true
	}
	jobs := map[string]*Job{}
	rows, err = s.pool.Query(ctx, `SELECT id, type_code, spec, started_at, ready_at, seed FROM helper.job
		WHERE character_id = $1 AND state = 'running'`, char)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var j Job
		var code string
		var spec []byte
		if err := rows.Scan(&j.ID, &code, &spec, &j.StartedAt, &j.ReadyAt, &j.seed); err != nil {
			rows.Close()
			return nil, err
		}
		_ = json.Unmarshal(spec, &j.Spec)
		jobs[code] = &j
	}
	rows.Close()
	out := make([]HelperView, 0, len(Types))
	for _, t := range Types {
		out = append(out, HelperView{Type: t, Hired: hired[t.Code], Job: jobs[t.Code]})
	}
	return out, rows.Err()
}

// Hire buys a helper for the character.
func (s *Service) Hire(ctx context.Context, char, account uuid.UUID, code string) error {
	t, ok := TypeOf(code)
	if !ok {
		return ErrUnknownType
	}
	if err := s.chars.Owns(ctx, char, account); err != nil {
		return err
	}
	var have bool
	if err := s.pool.QueryRow(ctx, `SELECT EXISTS (SELECT 1 FROM helper.hired WHERE character_id = $1 AND type_code = $2)`,
		char, code).Scan(&have); err != nil {
		return err
	}
	if have {
		return ErrOwned
	}
	if err := s.pay(ctx, fmt.Sprintf("helper-hire-%s-%s", char, code), char, t.PriceCopper, "helper_hire"); err != nil {
		return err
	}
	_, err := s.pool.Exec(ctx, `INSERT INTO helper.hired (character_id, type_code) VALUES ($1, $2) ON CONFLICT DO NOTHING`,
		char, code)
	return err
}

// StartGather sends a gathering helper out for duration on a resource the
// character's gathering level unlocks (the best one when target is empty).
func (s *Service) StartGather(ctx context.Context, char, account uuid.UUID, code, target string, duration int32) (*Job, error) {
	t, err := s.ready(ctx, char, account, code, KindGather)
	if err != nil {
		return nil, err
	}
	if duration < activity.MinDurationSeconds || duration > activity.MaxDurationSeconds {
		return nil, ErrValidation
	}
	kind := activity.Kind(t.Profession)
	level, err := s.gather.Level(ctx, char, kind)
	if err != nil {
		return nil, err
	}
	var node *activity.GatherNode
	for _, n := range s.gather.Nodes(kind) {
		if n.Level <= level && (target == "" || n.Code == target) {
			n := n
			node = &n
		}
	}
	if node == nil {
		return nil, ErrSkill
	}
	now := s.now().UTC()
	spec := Spec{Target: node.Code, Tier: node.Tier, Duration: duration}
	return s.insert(ctx, char, t, spec, now, now.Add(time.Duration(duration)*time.Second))
}

// StartCraft sets a crafting helper on count units of a recipe the
// character could make: materials and gold are taken now.
func (s *Service) StartCraft(ctx context.Context, char, account uuid.UUID, code, recipeCode string, count int) (*Job, error) {
	t, err := s.ready(ctx, char, account, code, KindCraft)
	if err != nil {
		return nil, err
	}
	if count < 1 || count > crafting.MaxCraftCount {
		return nil, ErrValidation
	}
	r, ok := s.recipes.Catalog().Recipe(recipeCode)
	if !ok || r.ProfessionCode != t.Profession {
		return nil, ErrValidation
	}
	cp, err := s.profs.Get(ctx, char, r.ProfessionCode)
	if err != nil || cp.Level < r.MinSkill {
		return nil, ErrSkill
	}
	if level, err := s.chars.Level(ctx, char); err != nil {
		return nil, err
	} else if level < int(r.OutputMinLevel) {
		return nil, ErrSkill
	}
	id := uuid.New()
	gold := t.GoldCost(r.GoldCostCopper, count)
	if err := s.pay(ctx, "helper-craft-"+id.String(), char, gold, "helper_craft"); err != nil {
		return nil, err
	}
	needs := make([]inventory.Need, len(r.Inputs))
	for i, in := range r.Inputs {
		needs[i] = inventory.Need{TemplateCode: in.TemplateCode, Quantity: in.Quantity * int32(count)}
	}
	if err := s.bag.ConsumeAll(ctx, char, needs); err != nil {
		s.refund(ctx, "helper-craft-refund-"+id.String(), char, gold)
		if errors.Is(err, inventory.ErrInsufficient) {
			return nil, ErrMaterials
		}
		return nil, err
	}
	now := s.now().UTC()
	spec := Spec{Recipe: r.Code, Count: count, Gold: gold}
	job, err := s.insertID(ctx, id, char, t, spec, now, now.Add(time.Duration(t.CraftSeconds(r.CraftSeconds, count))*time.Second))
	if err != nil {
		s.giveBack(ctx, char, id, r, count, gold)
	}
	return job, err
}

// Result is what a collected job brought.
type Result struct {
	Items      []activity.GrantedItem `json:"items,omitempty"`
	GoldCopper int64                  `json:"gold_copper,omitempty"`
	Made       int                    `json:"made,omitempty"`    // crafted units (doubles included)
	Doubled    int                    `json:"doubled,omitempty"` // of them, the bonus ones
	XP         int                    `json:"xp"`
	Profession string                 `json:"profession"`
}

// Collect brings a helper's work in: a gathering run any time (what it has
// so far), a craft once it is done. XP goes to the character's profession.
func (s *Service) Collect(ctx context.Context, char, account uuid.UUID, code string) (*Result, error) {
	t, ok := TypeOf(code)
	if !ok {
		return nil, ErrUnknownType
	}
	if err := s.chars.Owns(ctx, char, account); err != nil {
		return nil, err
	}
	job, err := s.running(ctx, char, code)
	if err != nil {
		return nil, err
	}
	now := s.now().UTC()
	res := &Result{Profession: t.Profession}
	if t.Kind == KindGather {
		elapsed := min(int32(now.Sub(job.StartedAt).Seconds()), job.Spec.Duration)
		y := activity.ComputeYield(&activity.Contract{
			Kind: activity.Kind(t.Profession), DurationSeconds: job.Spec.Duration, Seed: job.seed,
			Snapshot: activity.MarshalSnapshot(activity.Snapshot{YieldMult: t.YieldMult(),
				Target: job.Spec.Target, TargetTier: job.Spec.Tier}),
		}, max(elapsed, 0))
		if res.Items, err = s.items.Grant(ctx, char, job.ID, job.seed, activity.Drops{Kind: activity.Kind(t.Profession),
			Target: job.Spec.Target, Common: y.CommonCount, Rare: y.RareCount, Legendary: y.LegendaryCount}); err != nil {
			return nil, err
		}
		if y.GoldCopper > 0 {
			if err := s.credit(ctx, "helper-gather-"+job.ID.String(), char, y.GoldCopper); err != nil {
				return nil, err
			}
		}
		res.GoldCopper = y.GoldCopper
		res.XP = t.XP((y.CommonCount + y.RareCount + y.LegendaryCount) * activity.TierXP(job.Spec.Tier))
	} else {
		if now.Before(job.ReadyAt) {
			return nil, ErrNotReady
		}
		r, ok := s.recipes.Catalog().Recipe(job.Spec.Recipe)
		if !ok {
			return nil, ErrValidation
		}
		res.Doubled = t.Doubles(job.Spec.Count, func(i int) int {
			return int(activity.NewRNG(job.seed+int64(i)).Uint64() % 100)
		})
		res.Made = job.Spec.Count + res.Doubled
		res.XP = t.XP(int(r.ProfessionXPReward) * job.Spec.Count)
	}
	// One collection: the job closes before anything is minted or XP paid.
	body, _ := json.Marshal(res)
	tag, err := s.pool.Exec(ctx, `UPDATE helper.job SET state = 'done', collected_at = $2, result = $3
		WHERE id = $1 AND state = 'running'`, job.ID, now, body)
	if err != nil {
		return nil, err
	}
	if tag.RowsAffected() == 0 {
		return nil, ErrIdle
	}
	if t.Kind == KindCraft {
		r, _ := s.recipes.Catalog().Recipe(job.Spec.Recipe)
		if err := s.mint(ctx, char, job.ID, r.OutputTemplateCode, int32(res.Made)*r.OutputStack); err != nil {
			s.log.ErrorContext(ctx, "helper: craft output mint failed", "job", job.ID.String(), "err", err.Error())
		}
		if res.XP > 0 {
			if err := s.craftXP.GrantXP(ctx, char, t.Profession, int32(res.XP)); err != nil {
				s.log.ErrorContext(ctx, "helper: craft xp failed", "job", job.ID.String(), "err", err.Error())
			}
		}
	} else if res.XP > 0 {
		if _, err := s.gather.AddXP(ctx, char, activity.Kind(t.Profession), res.XP); err != nil {
			s.log.ErrorContext(ctx, "helper: gather xp failed", "job", job.ID.String(), "err", err.Error())
		}
	}
	return res, nil
}

// Cancel calls a helper back: a gathering run brings nothing, a craft
// gives its materials and gold back.
func (s *Service) Cancel(ctx context.Context, char, account uuid.UUID, code string) error {
	t, ok := TypeOf(code)
	if !ok {
		return ErrUnknownType
	}
	if err := s.chars.Owns(ctx, char, account); err != nil {
		return err
	}
	job, err := s.running(ctx, char, code)
	if err != nil {
		return err
	}
	tag, err := s.pool.Exec(ctx, `UPDATE helper.job SET state = 'cancelled', collected_at = now()
		WHERE id = $1 AND state = 'running'`, job.ID)
	if err != nil || tag.RowsAffected() == 0 {
		return err
	}
	if t.Kind == KindCraft {
		if r, ok := s.recipes.Catalog().Recipe(job.Spec.Recipe); ok {
			s.giveBack(ctx, char, job.ID, r, job.Spec.Count, job.Spec.Gold)
		}
	}
	return nil
}

// --- internals ----------------------------------------------------------------

// ready checks the character owns an idle helper of the code and kind.
func (s *Service) ready(ctx context.Context, char, account uuid.UUID, code, kind string) (Type, error) {
	t, ok := TypeOf(code)
	if !ok || t.Kind != kind {
		return Type{}, ErrUnknownType
	}
	if err := s.chars.Owns(ctx, char, account); err != nil {
		return Type{}, err
	}
	var hired, busy bool
	if err := s.pool.QueryRow(ctx, `
		SELECT EXISTS (SELECT 1 FROM helper.hired WHERE character_id = $1 AND type_code = $2),
		       EXISTS (SELECT 1 FROM helper.job WHERE character_id = $1 AND type_code = $2 AND state = 'running')`,
		char, code).Scan(&hired, &busy); err != nil {
		return Type{}, err
	}
	if !hired {
		return Type{}, ErrNotOwned
	}
	if busy {
		return Type{}, ErrBusy
	}
	return t, nil
}

func (s *Service) running(ctx context.Context, char uuid.UUID, code string) (*Job, error) {
	var j Job
	var spec []byte
	err := s.pool.QueryRow(ctx, `SELECT id, spec, started_at, ready_at, seed FROM helper.job
		WHERE character_id = $1 AND type_code = $2 AND state = 'running'`, char, code).
		Scan(&j.ID, &spec, &j.StartedAt, &j.ReadyAt, &j.seed)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrIdle
	}
	if err != nil {
		return nil, err
	}
	_ = json.Unmarshal(spec, &j.Spec)
	return &j, nil
}

func (s *Service) insert(ctx context.Context, char uuid.UUID, t Type, spec Spec, start, ready time.Time) (*Job, error) {
	return s.insertID(ctx, uuid.New(), char, t, spec, start, ready)
}

func (s *Service) insertID(ctx context.Context, id, char uuid.UUID, t Type, spec Spec, start, ready time.Time) (*Job, error) {
	body, _ := json.Marshal(spec)
	seed := int64(id.ID())<<32 | int64(char.ID())
	_, err := s.pool.Exec(ctx, `INSERT INTO helper.job (id, character_id, type_code, kind, spec, seed, started_at, ready_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`, id, char, t.Code, t.Kind, body, seed, start, ready)
	if isUnique(err) {
		return nil, ErrBusy
	}
	if err != nil {
		return nil, err
	}
	return &Job{ID: id, Spec: spec, StartedAt: start, ReadyAt: ready, seed: seed}, nil
}

// giveBack returns a craft's materials and gold.
func (s *Service) giveBack(ctx context.Context, char, job uuid.UUID, r crafting.Recipe, count int, gold int64) {
	for _, in := range r.Inputs {
		if err := s.mint(ctx, char, job, in.TemplateCode, in.Quantity*int32(count)); err != nil {
			s.log.ErrorContext(ctx, "helper: material refund failed", "job", job.String(), "err", err.Error())
		}
	}
	s.refund(ctx, "helper-craft-refund-"+job.String(), char, gold)
}

// mint puts qty of a template in the bag, in stacks of its max size.
func (s *Service) mint(ctx context.Context, char, ref uuid.UUID, code string, qty int32) error {
	var maxStack int32
	if err := s.pool.QueryRow(ctx, `SELECT max_stack FROM content.item_template WHERE code = $1`, code).Scan(&maxStack); err != nil {
		return err
	}
	maxStack = max(maxStack, 1)
	for qty > 0 {
		n := min(qty, maxStack)
		if _, err := s.items.MintExact(ctx, char, code, n, "helper", ref); err != nil {
			return err
		}
		qty -= n
	}
	return s.bag.PlaceOverflow(ctx, char)
}

func (s *Service) pay(ctx context.Context, key string, char uuid.UUID, copper int64, kind string) error {
	if copper <= 0 {
		return nil
	}
	_, err := s.txn.Write(ctx, economy.Intent{Kind: kind, RefKind: "helper", IdempotencyKey: key,
		Entries: []economy.Entry{
			{Currency: "gold", AccountKind: economy.AccountCharacter, AccountRef: char, Delta: -copper},
			{Currency: "gold", AccountKind: economy.AccountSystem, AccountRef: economy.SystemVendorSink, Delta: copper},
		}})
	if errors.Is(err, economy.ErrInsufficientBalance) {
		return ErrNoGold
	}
	return err
}

func (s *Service) refund(ctx context.Context, key string, char uuid.UUID, copper int64) {
	if copper <= 0 {
		return
	}
	if _, err := s.txn.Write(ctx, economy.Intent{Kind: "helper_refund", RefKind: "helper", IdempotencyKey: key,
		Entries: []economy.Entry{
			{Currency: "gold", AccountKind: economy.AccountSystem, AccountRef: economy.SystemVendorSink, Delta: -copper},
			{Currency: "gold", AccountKind: economy.AccountCharacter, AccountRef: char, Delta: copper},
		}}); err != nil {
		s.log.ErrorContext(ctx, "helper: gold refund failed", "key", key, "err", err.Error())
	}
}

func (s *Service) credit(ctx context.Context, key string, char uuid.UUID, copper int64) error {
	_, err := s.txn.Write(ctx, economy.Intent{Kind: "helper_gather", RefKind: "helper", IdempotencyKey: key,
		Entries: []economy.Entry{
			{Currency: "gold", AccountKind: economy.AccountSystem, AccountRef: economy.SystemRewardSource, Delta: -copper},
			{Currency: "gold", AccountKind: economy.AccountCharacter, AccountRef: char, Delta: copper},
		}})
	return err
}

func isUnique(err error) bool {
	var pg interface{ SQLState() string }
	return errors.As(err, &pg) && pg.SQLState() == "23505"
}
