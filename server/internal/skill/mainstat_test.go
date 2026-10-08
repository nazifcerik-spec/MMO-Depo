package skill_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/otter/server/internal/skill"
)

// DD-229: attack follows the class main stat.
func TestComputeSheetFor_MainStat(t *testing.T) {
	base := skill.Primary{STR: 10, AGI: 20, INT: 30, VIT: 10, CHA: 10, Luck: 10}
	require.Equal(t, 20, skill.ComputeSheetFor(base, skill.TargetSTR, nil).Derived.Attack)
	require.Equal(t, 40, skill.ComputeSheetFor(base, skill.TargetAGI, nil).Derived.Attack)
	require.Equal(t, 60, skill.ComputeSheetFor(base, skill.TargetINT, nil).Derived.Attack)
}
