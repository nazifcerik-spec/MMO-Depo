package skill

import (
	"context"
	"log/slog"

	"github.com/google/uuid"

	"github.com/otter/server/internal/domain/character"
	"github.com/otter/server/internal/leveling"
)

// EquippedFetcher exposes what the sheet computer needs from inventory:
// the character's currently equipped items (each carries its template's
// stat_bonuses JSONB).
type EquippedFetcher interface {
	EquippedBonuses(ctx context.Context, characterID uuid.UUID) ([]EquippedBonus, error)
}

// EquippedBonus is the projection of an equipped item's stat bonuses.
type EquippedBonus struct {
	StatBonusesRaw []byte
}

// CharacterAdapter satisfies character.SheetComputer by using an EquippedFetcher
// (implemented by inventory) + the skill package's pure ComputeSheet.
//
// character.Service.WithSheetComputer(charAdapter) wires this in main.go.
type CharacterAdapter struct {
	fetch EquippedFetcher
	extra []ModifierSource
	log   *slog.Logger
}

// ModifierSource contributes non-gear modifiers (mastery bonuses, DD-220).
// Targets the sheet does not know (e.g. gather_yield_*) are ignored by
// ComputeSheet.
// base is the pre-modifier sheet (class, level, ...).
type ModifierSource func(ctx context.Context, characterID uuid.UUID, base character.Sheet) ([]Modifier, error)

// WithModifierSource adds a source; its error fails the sheet like a gear
// fetch error would (a silently missing bonus would skew combat snapshots).
func (a *CharacterAdapter) WithModifierSource(src ModifierSource) *CharacterAdapter {
	a.extra = append(a.extra, src)
	return a
}

func NewCharacterAdapter(f EquippedFetcher, logger *slog.Logger) *CharacterAdapter {
	if logger == nil {
		logger = slog.Default()
	}
	return &CharacterAdapter{fetch: f, log: logger}
}

// ComputeSheet reads equipped bonuses and applies the modifier stack to the
// base sheet's primary stats. Derived + resources are then computed by
// skill.ComputeSheet's formulas — the placeholder derived values from the
// character package are overwritten.
func (a *CharacterAdapter) ComputeSheet(ctx context.Context, characterID uuid.UUID, base character.Sheet) (character.Sheet, error) {
	bonuses, err := a.fetch.EquippedBonuses(ctx, characterID)
	if err != nil {
		return base, err
	}
	mods := []Modifier{}
	for _, b := range bonuses {
		list, err := DecodeModifiers(b.StatBonusesRaw)
		if err != nil {
			a.log.WarnContext(ctx, "skill: bad stat_bonuses row — skipped",
				slog.String("err", err.Error()))
			continue
		}
		mods = append(mods, list...)
	}
	for _, src := range a.extra {
		list, err := src(ctx, characterID, base)
		if err != nil {
			return base, err
		}
		mods = append(mods, list...)
	}

	// Base primary from the placeholder sheet.
	basePrimary := Primary{
		STR:  base.Primary.Strength,
		AGI:  base.Primary.Agility,
		INT:  base.Primary.Intelligence,
		VIT:  base.Primary.Vitality,
		CHA:  base.Primary.Charisma,
		Luck: base.Primary.Luck,
	}
	var formula *Formula
	if k, ok := leveling.Current().Class(base.Class); ok && k.Derived != nil {
		formula = &Formula{Derived: *k.Derived, Level: base.Level}
	}
	computed := ComputeSheetWith(basePrimary, Target(character.MainStat(base.Class)), mods, formula)

	// Merge computed back into character.Sheet's shape.
	base.Primary = character.PrimaryStats{
		Strength:     computed.Primary.STR,
		Agility:      computed.Primary.AGI,
		Intelligence: computed.Primary.INT,
		Vitality:     computed.Primary.VIT,
		Charisma:     computed.Primary.CHA,
		Luck:         computed.Primary.Luck,
	}
	base.Derived = character.DerivedStats{
		Attack:  computed.Derived.Attack,
		Defense: computed.Derived.Defense,
		CritPct: computed.Derived.CritPct,
		HitPct:  computed.Derived.HitPct,

		EvasionPct:      computed.Derived.EvasionPct,
		BlockPct:        computed.Derived.BlockPct,
		CritDmgPct:      computed.Derived.CritDmgPct,
		ExecutePct:      computed.Derived.ExecutePct,
		MissingHPDmgPct: computed.Derived.MissingHPDmgPct,
		LifestealPct:    computed.Derived.LifestealPct,
		DmgReductionPct: computed.Derived.DmgReductionPct,
	}
	base.Resources = character.Resources{
		HPCurrent: computed.Resources.HPMax, HPMax: computed.Resources.HPMax,
		MPCurrent: computed.Resources.MPMax, MPMax: computed.Resources.MPMax,
		EnergyCurrent: computed.Resources.EnMax, EnergyMax: computed.Resources.EnMax,
		StaminaCurrent: computed.Resources.StaMax, StaminaMax: computed.Resources.StaMax,
	}
	return base, nil
}
