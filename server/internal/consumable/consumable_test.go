package consumable

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/otter/server/internal/domain/character"
	"github.com/otter/server/internal/skill"
)

func TestApplyHeal(t *testing.T) {
	hp, err := ApplyHeal(40, 200, 60)
	require.NoError(t, err)
	require.Equal(t, 100, hp)

	hp, err = ApplyHeal(180, 200, 60)
	require.NoError(t, err)
	require.Equal(t, 200, hp, "capped at max")

	_, err = ApplyHeal(200, 200, 60)
	require.ErrorIs(t, err, ErrNoEffect, "a potion at full HP must not be consumed")
}

func TestEffectiveHP(t *testing.T) {
	now := time.Date(2026, 9, 23, 12, 0, 0, 0, time.UTC)
	at := func(ago time.Duration) *time.Time { v := now.Add(-ago); return &v }
	hp := func(v int32) *int32 { return &v }

	require.Equal(t, 250, character.EffectiveHP(nil, nil, 250, 15, now), "never damaged = full")
	require.Equal(t, 10, character.EffectiveHP(hp(10), at(0), 250, 15, now), "no time passed")
	// VIT 15 → 2.5 HP per 2s step; 60s = 30 steps = 75 HP.
	require.Equal(t, 85, character.EffectiveHP(hp(10), at(time.Minute), 250, 15, now))
	require.Equal(t, 10, character.EffectiveHP(hp(10), at(time.Second), 250, 15, now), "partial step gives nothing")
	require.Equal(t, 250, character.EffectiveHP(hp(0), at(time.Hour), 250, 15, now), "capped at max")
	require.Equal(t, 5, character.EffectiveHP(hp(5), at(-time.Minute), 250, 15, now), "clock skew never drains HP")
	require.Equal(t, 120, character.EffectiveHP(hp(300), at(0), 120, 15, now), "max dropped (gear removed) → clamp")
}

// Food buffs become attack/defense percent modifiers (DD-269).
func TestBuffModifiers(t *testing.T) {
	mods := BuffModifiers(&Buff{AttackPct: 6, DefensePct: 4})
	require.Len(t, mods, 2)
	require.Equal(t, skill.Modifier{Kind: skill.KindMultAdd, Target: skill.TargetAttack, Amount: 0.06}, mods[0])
	require.Equal(t, skill.Modifier{Kind: skill.KindMultAdd, Target: skill.TargetDefense, Amount: 0.04}, mods[1])
	require.Empty(t, BuffModifiers(&Buff{}))
}
