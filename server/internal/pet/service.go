package pet

import (
	"context"
	"errors"
	"math/rand/v2"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/otter/server/internal/activity"
	"github.com/otter/server/internal/combat"
	"github.com/otter/server/internal/domain/character"
	"github.com/otter/server/internal/skill"
)

// MaxPets is how many pets a character keeps.
const MaxPets = 50

var (
	ErrNotFound   = errors.New("pet: not found")
	ErrNotEgg     = errors.New("pet: not a usable egg")
	ErrCap        = errors.New("pet: too many pets")
	ErrHatching   = errors.New("pet: still in its egg")
	ErrHunting    = errors.New("pet: already hunting")
	ErrIdle       = errors.New("pet: not hunting")
	ErrValidation = errors.New("pet: invalid hunt")
)

// Characters owns characters (character.Service).
type Characters interface {
	Owns(ctx context.Context, characterID, accountID uuid.UUID) error
}

// Mobs is the mob catalog (combat.Service).
type Mobs interface {
	Catalog() *combat.Catalog
}

type Service struct {
	pool  *pgxpool.Pool
	chars Characters
	mobs  Mobs
	now   func() time.Time
}

func NewService(pool *pgxpool.Pool, chars Characters, mobs Mobs) *Service {
	return &Service{pool: pool, chars: chars, mobs: mobs, now: time.Now}
}

// WithClock injects the clock (tests).
func (s *Service) WithClock(now func() time.Time) *Service { s.now = now; return s }

// Hunt is a pet's running hunt.
type Hunt struct {
	Mob       string    `json:"mob_code"`
	StartedAt time.Time `json:"started_at"`
	ReadyAt   time.Time `json:"ready_at"`
	seed      int64
}

// Pet is one of a character's pets as the API shows it.
type Pet struct {
	ID           uuid.UUID `json:"id"`
	Species      string    `json:"species"`
	Rank         int       `json:"rank"`
	Level        int       `json:"level"`
	XP           int64     `json:"xp"`
	XPToNext     int64     `json:"xp_to_next"`
	Power        int       `json:"power"` // the mob level it fights like
	HatchReadyAt time.Time `json:"hatch_ready_at"`
	Hatched      bool      `json:"hatched"`
	Active       bool      `json:"active"`
	Bonus        Modifier  `json:"bonus"` // at its level
	Hunt         *Hunt     `json:"hunt,omitempty"`
}

// Egg is a stack of eggs in the bag.
type Egg struct {
	ItemID   uuid.UUID `json:"item_id"`
	Template string    `json:"template_code"`
	Species  string    `json:"species"`
	Quantity int32     `json:"quantity"`
}

const petColumns = `id, species_code, level, xp, hatch_ready_at, active, hunt_mob, hunt_started_at, hunt_ready_at, hunt_seed`

func (s *Service) scan(row pgx.Row) (*Pet, error) {
	var p Pet
	var mob *string
	var started, ready *time.Time
	var seed *int64
	var level int16
	if err := row.Scan(&p.ID, &p.Species, &level, &p.XP, &p.HatchReadyAt, &p.Active, &mob, &started, &ready, &seed); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	p.Level = int(level)
	if mob != nil {
		p.Hunt = &Hunt{Mob: *mob, StartedAt: *started, ReadyAt: *ready, seed: *seed}
	}
	s.fill(&p)
	return &p, nil
}

// fill derives what the row does not store.
func (s *Service) fill(p *Pet) {
	p.XPToNext = XPToNext(p.Level)
	p.Power = Power(p.Level)
	p.Hatched = !s.now().Before(p.HatchReadyAt)
	if sp, ok := Current().Get(p.Species); ok {
		p.Rank = sp.Rank
		p.Bonus = sp.BonusAt(p.Level)
	}
}

// List is the character's pets (best and highest first) and its eggs.
func (s *Service) List(ctx context.Context, char, account uuid.UUID) ([]Pet, []Egg, error) {
	if err := s.chars.Owns(ctx, char, account); err != nil {
		return nil, nil, err
	}
	rows, err := s.pool.Query(ctx, `SELECT `+petColumns+` FROM pet.pet WHERE character_id = $1
		ORDER BY active DESC, level DESC, created_at`, char)
	if err != nil {
		return nil, nil, err
	}
	pets := []Pet{}
	for rows.Next() {
		p, err := s.scan(rows)
		if err != nil {
			rows.Close()
			return nil, nil, err
		}
		pets = append(pets, *p)
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return nil, nil, err
	}
	rows, err = s.pool.Query(ctx, `SELECT id, template_code, stack_quantity FROM inventory.item_instance
		WHERE owner_character_id = $1 AND template_code LIKE '`+EggPrefix+`%' AND NOT is_locked
		  AND equipment_slot IS NULL ORDER BY template_code`, char)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()
	eggs := []Egg{}
	for rows.Next() {
		var e Egg
		if err := rows.Scan(&e.ItemID, &e.Template, &e.Quantity); err != nil {
			return nil, nil, err
		}
		sp, ok := Current().ByEgg(e.Template)
		if !ok {
			continue
		}
		e.Species = sp.Code
		eggs = append(eggs, e)
	}
	return pets, eggs, rows.Err()
}

// Hatch puts one egg of the bag stack in the incubator: the egg is used and
// the pet comes out after its species' hatch time.
func (s *Service) Hatch(ctx context.Context, char, account, itemID uuid.UUID) (*Pet, error) {
	if err := s.chars.Owns(ctx, char, account); err != nil {
		return nil, err
	}
	var out *Pet
	err := pgx.BeginFunc(ctx, s.pool, func(tx pgx.Tx) error {
		// The character row orders hatches against the pet cap.
		if _, err := tx.Exec(ctx, `SELECT 1 FROM character.character WHERE id = $1 FOR UPDATE`, char); err != nil {
			return err
		}
		var tmpl string
		var stack int32
		err := tx.QueryRow(ctx, `SELECT template_code, stack_quantity FROM inventory.item_instance
			WHERE id = $1 AND owner_character_id = $2 AND NOT is_locked AND equipment_slot IS NULL
			FOR UPDATE`, itemID, char).Scan(&tmpl, &stack)
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrNotEgg
		}
		if err != nil {
			return err
		}
		sp, ok := Current().ByEgg(tmpl)
		if !ok {
			return ErrNotEgg
		}
		var n int
		if err := tx.QueryRow(ctx, `SELECT count(*) FROM pet.pet WHERE character_id = $1`, char).Scan(&n); err != nil {
			return err
		}
		if n >= MaxPets {
			return ErrCap
		}
		if stack > 1 {
			_, err = tx.Exec(ctx, `UPDATE inventory.item_instance SET stack_quantity = stack_quantity - 1,
				updated_at = now(), version = version + 1 WHERE id = $1`, itemID)
		} else {
			_, err = tx.Exec(ctx, `DELETE FROM inventory.item_instance WHERE id = $1`, itemID)
		}
		if err != nil {
			return err
		}
		now := s.now().UTC()
		out, err = s.scan(tx.QueryRow(ctx, `INSERT INTO pet.pet (id, character_id, species_code, hatch_ready_at, created_at)
			VALUES ($1, $2, $3, $4, $5) RETURNING `+petColumns,
			uuid.New(), char, sp.Code, now.Add(time.Duration(sp.HatchSeconds)*time.Second), now))
		return err
	})
	return out, err
}

func (s *Service) owned(ctx context.Context, char, account, id uuid.UUID) (*Pet, error) {
	if err := s.chars.Owns(ctx, char, account); err != nil {
		return nil, err
	}
	return s.scan(s.pool.QueryRow(ctx, `SELECT `+petColumns+` FROM pet.pet WHERE id = $1 AND character_id = $2`, id, char))
}

// SetActive makes the pet the character's active one (on) or none (off).
func (s *Service) SetActive(ctx context.Context, char, account, id uuid.UUID, on bool) error {
	p, err := s.owned(ctx, char, account, id)
	if err != nil {
		return err
	}
	if on && !p.Hatched {
		return ErrHatching
	}
	return pgx.BeginFunc(ctx, s.pool, func(tx pgx.Tx) error {
		if _, err := tx.Exec(ctx, `UPDATE pet.pet SET active = false WHERE character_id = $1 AND active`, char); err != nil {
			return err
		}
		if !on {
			return nil
		}
		_, err := tx.Exec(ctx, `UPDATE pet.pet SET active = true WHERE id = $1`, id)
		return err
	})
}

// StartHunt sends a hatched pet hunting mob for duration.
func (s *Service) StartHunt(ctx context.Context, char, account, id uuid.UUID, mobCode string, duration int32) (*Pet, error) {
	p, err := s.owned(ctx, char, account, id)
	if err != nil {
		return nil, err
	}
	if !p.Hatched {
		return nil, ErrHatching
	}
	if p.Hunt != nil {
		return nil, ErrHunting
	}
	mob, ok := s.mobs.Catalog().Get(mobCode)
	if !ok || !CanHunt(p.Level, mob) || duration < activity.MinDurationSeconds || duration > activity.MaxDurationSeconds {
		return nil, ErrValidation
	}
	now := s.now().UTC()
	out, err := s.scan(s.pool.QueryRow(ctx, `UPDATE pet.pet SET hunt_mob = $2, hunt_started_at = $3, hunt_ready_at = $4,
		hunt_seed = $5 WHERE id = $1 AND hunt_mob IS NULL RETURNING `+petColumns,
		id, mobCode, now, now.Add(time.Duration(duration)*time.Second), rand.Int64()))
	if errors.Is(err, ErrNotFound) {
		return nil, ErrHunting // a parallel start won
	}
	return out, err
}

// HuntResult is what a collected pet hunt brought.
type HuntResult struct {
	Victories int   `json:"victories"`
	Died      bool  `json:"died"` // a defeat ended the hunt early
	XP        int64 `json:"xp"`
	LeveledUp bool  `json:"leveled_up"`
	Pet       *Pet  `json:"pet"`
}

// CollectHunt ends the pet's hunt any time: the fights of the time so far
// (all of it once ready) pay their XP.
func (s *Service) CollectHunt(ctx context.Context, char, account, id uuid.UUID) (*HuntResult, error) {
	if err := s.chars.Owns(ctx, char, account); err != nil {
		return nil, err
	}
	var out *HuntResult
	err := pgx.BeginFunc(ctx, s.pool, func(tx pgx.Tx) error {
		p, err := s.scan(tx.QueryRow(ctx, `SELECT `+petColumns+` FROM pet.pet WHERE id = $1 AND character_id = $2
			FOR UPDATE`, id, char))
		if err != nil {
			return err
		}
		if p.Hunt == nil {
			return ErrIdle
		}
		now := s.now().UTC()
		end := p.Hunt.ReadyAt
		if now.Before(end) {
			end = now
		}
		res := &HuntResult{}
		sp, okS := Current().Get(p.Species)
		mob, okM := s.mobs.Catalog().Get(p.Hunt.Mob)
		if budget := int32(end.Sub(p.Hunt.StartedAt).Seconds()); okS && okM && budget > 0 {
			hr := Roll(p.Hunt.seed, sp, p.Level, mob, budget, func(seed int64) combat.RNG { return activity.NewRNG(seed) })
			res.Victories, res.Died, res.XP = hr.Victories, hr.StoppedOnDeath, hr.TotalXP
		}
		level, xp := AddXP(p.Level, p.XP, res.XP)
		res.LeveledUp = level > p.Level
		res.Pet, err = s.scan(tx.QueryRow(ctx, `UPDATE pet.pet SET level = $2, xp = $3, hunt_mob = NULL,
			hunt_started_at = NULL, hunt_ready_at = NULL, hunt_seed = NULL WHERE id = $1 RETURNING `+petColumns,
			id, level, xp))
		out = res
		return err
	})
	return out, err
}

// Modifiers is the active, hatched pet's bonus on the character's sheet.
func (s *Service) Modifiers(ctx context.Context, id uuid.UUID, _ character.Sheet) ([]skill.Modifier, error) {
	var code string
	var level int16
	err := s.pool.QueryRow(ctx, `SELECT species_code, level FROM pet.pet
		WHERE character_id = $1 AND active AND hatch_ready_at <= $2`, id, s.now()).Scan(&code, &level)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	sp, ok := Current().Get(code)
	if !ok {
		return nil, nil
	}
	b := sp.BonusAt(int(level))
	return []skill.Modifier{{Kind: skill.ModifierKind(b.Kind), Target: skill.Target(b.Target), Amount: b.Amount}}, nil
}
