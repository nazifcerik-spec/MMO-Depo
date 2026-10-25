package hunt

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/otter/server/internal/combat"
	"github.com/otter/server/internal/domain/character"
)

type memRepo struct{ created []*Contract }

func (m *memRepo) Create(_ context.Context, c *Contract) error {
	m.created = append(m.created, c)
	return nil
}
func (m *memRepo) GetByID(context.Context, uuid.UUID) (*Contract, error) {
	return nil, errors.New("unused")
}
func (m *memRepo) Active(context.Context, uuid.UUID) (*Contract, error) {
	return nil, errors.New("unused")
}
func (m *memRepo) History(context.Context, uuid.UUID, int) ([]Contract, error) {
	return nil, errors.New("unused")
}
func (m *memRepo) MarkRevealed(context.Context, uuid.UUID, json.RawMessage, uuid.UUID, time.Time, []SupplyUse) error {
	return errors.New("unused")
}
func (m *memRepo) MarkCancelled(context.Context, uuid.UUID, time.Time) error {
	return errors.New("unused")
}

type mobs struct{ c *combat.Catalog }

func (m mobs) Catalog() *combat.Catalog { return m.c }

type sheets struct {
	sheet    *character.Sheet
	ownerErr error
}

func (s sheets) Sheet(context.Context, uuid.UUID, uuid.UUID) (*character.Sheet, error) {
	return s.sheet, nil
}
func (s sheets) Owns(context.Context, uuid.UUID, uuid.UUID) error { return s.ownerErr }

func TestStart_FreezesSkillsAndGuards(t *testing.T) {
	sheet := &character.Sheet{Class: "mage", Level: 6,
		Resources: character.Resources{HPCurrent: 90, HPMax: 100, MPMax: 125},
		Derived:   character.DerivedStats{Attack: 20, Defense: 5, CritPct: 5, HitPct: 90}}
	repo := &memRepo{}
	cat := combat.NewCatalog([]combat.Mob{{Code: "boar", HPMax: 40, Attack: 5, HitPct: 80, XPReward: 10}})
	src := sheets{sheet: sheet}
	svc := NewService(repo, mobs{cat}, src, nil, nil, nil, DefaultServiceConfig(), nil, newRNG).
		WithSkills(func(s *character.Sheet) []combat.Skill {
			require.Equal(t, "mage", s.Class)
			return []combat.Skill{{Code: "bolt", Effect: combat.EffectStrike, Power: 200, MPCost: 10, CooldownTurns: 2}}
		})
	req := StartRequest{CharacterID: uuid.New(), AccountID: uuid.New(), MobCode: "boar", DurationSeconds: 600}

	res, err := svc.Start(context.Background(), req)
	require.NoError(t, err)
	require.Len(t, repo.created, 1)
	var snap combat.Combatant
	require.NoError(t, json.Unmarshal(res.Contract.PlayerSnapshot, &snap))
	require.Equal(t, int32(90), snap.HPCurrent)
	require.Equal(t, int32(125), snap.MPMax)
	require.Equal(t, "bolt", snap.Skills[0].Code)

	bad := req
	bad.MobCode = ""
	_, err = svc.Start(context.Background(), bad)
	require.ErrorIs(t, err, ErrValidation)
	bad = req
	bad.DurationSeconds = 0
	_, err = svc.Start(context.Background(), bad)
	require.ErrorIs(t, err, ErrValidation)
	bad = req
	bad.MobCode = "dragon"
	_, err = svc.Start(context.Background(), bad)
	require.ErrorIs(t, err, ErrUnknownMob)
	bad = req
	bad.AccountID = uuid.Nil
	_, err = svc.Start(context.Background(), bad)
	require.ErrorIs(t, err, ErrForbidden)

	sheet.Resources.HPCurrent = 0
	_, err = svc.Start(context.Background(), req)
	require.ErrorIs(t, err, ErrTooWounded)
}

type gate []string

func (g gate) HuntRegions(context.Context, uuid.UUID) ([]string, error) { return g, nil }

// A gated hunt only starts on mobs in reach of where the character stands
// (DD-273); mobs without a region (old content) are never gated.
func TestStart_RegionGate(t *testing.T) {
	sheet := &character.Sheet{Class: "warrior", Level: 5,
		Resources: character.Resources{HPCurrent: 90, HPMax: 100},
		Derived:   character.DerivedStats{Attack: 20, Defense: 5, HitPct: 90}}
	cat := combat.NewCatalog([]combat.Mob{
		{Code: "boar", HPMax: 40, Attack: 5, HitPct: 80, RegionCode: "greenwood"},
		{Code: "wyrm", HPMax: 40, Attack: 5, HitPct: 80, RegionCode: "emberpeak"},
		{Code: "rat", HPMax: 40, Attack: 5, HitPct: 80},
	})
	svc := NewService(&memRepo{}, mobs{cat}, sheets{sheet: sheet}, nil, nil, nil, DefaultServiceConfig(), nil, newRNG).
		WithRegionGate(gate{"ironhold_city", "greenwood"})
	req := StartRequest{CharacterID: uuid.New(), AccountID: uuid.New(), DurationSeconds: 600}
	for mob, want := range map[string]error{"boar": nil, "rat": nil, "wyrm": ErrWrongRegion} {
		req.MobCode = mob
		_, err := svc.Start(context.Background(), req)
		if want == nil {
			require.NoError(t, err, mob)
		} else {
			require.ErrorIs(t, err, want, mob)
		}
	}
}
