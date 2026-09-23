package skill_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/otter/server/internal/domain/character"
	"github.com/otter/server/internal/skill"
)

type gear []skill.EquippedBonus

func (g gear) EquippedBonuses(context.Context, uuid.UUID) ([]skill.EquippedBonus, error) {
	return g, nil
}

func TestCharacterAdapter_AddsModifierSources(t *testing.T) {
	base := character.Sheet{Primary: character.PrimaryStats{Strength: 10, Agility: 10, Intelligence: 10,
		Vitality: 10, Charisma: 10, Luck: 10}}
	a := skill.NewCharacterAdapter(gear{{StatBonusesRaw: []byte(`[{"kind":"flat","target":"attack","amount":3}]`)}}, nil)

	plain, err := a.ComputeSheet(context.Background(), uuid.New(), base)
	require.NoError(t, err)
	require.Equal(t, 23, plain.Derived.Attack) // STR×2 + gear 3

	a.WithModifierSource(func(context.Context, uuid.UUID) ([]skill.Modifier, error) {
		return []skill.Modifier{
			{Kind: skill.KindFlat, Target: skill.TargetAttack, Amount: 2},
			{Kind: skill.KindFlat, Target: "gather_yield_mining", Amount: 1}, // not a sheet stat: ignored
		}, nil
	})
	got, err := a.ComputeSheet(context.Background(), uuid.New(), base)
	require.NoError(t, err)
	require.Equal(t, 25, got.Derived.Attack)
	require.Equal(t, plain.Primary, got.Primary)

	a.WithModifierSource(func(context.Context, uuid.UUID) ([]skill.Modifier, error) {
		return nil, errors.New("db down")
	})
	_, err = a.ComputeSheet(context.Background(), uuid.New(), base)
	require.Error(t, err, "a failing source must not silently drop bonuses")
}
