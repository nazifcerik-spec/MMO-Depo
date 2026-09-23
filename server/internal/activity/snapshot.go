package activity

import (
	"encoding/json"

	"github.com/otter/server/internal/domain/character"
)

// Snapshot captures the character state at contract start (DD-088). The reveal
// path must NOT read current character state — everything comes from this
// frozen record. Guarantees deterministic reveal even if the character was
// re-geared / respec'd during the AFK period.
//
// v1 minimal snapshot — expands with Sprint 11 (skill engine: gear score,
// mastery tier, race bonus, active buffs).
type Snapshot struct {
	CharacterLevel int    `json:"character_level"`
	ClassCode      string `json:"class_code"`
	RaceCode       string `json:"race_code"`

	// Future (Sprint 11+):
	// GearScore     int            `json:"gear_score,omitempty"`
	// MasteryTier   int            `json:"mastery_tier,omitempty"`
	// ActiveBuffs   []string       `json:"active_buffs,omitempty"`
	// GatherBonus   float64        `json:"gather_bonus,omitempty"`
}

// TakeSnapshot builds a snapshot from a live character.
func TakeSnapshot(c *character.Character) Snapshot {
	return Snapshot{
		CharacterLevel: int(c.Level),
		ClassCode:      c.ClassCode,
		RaceCode:       c.RaceCode,
	}
}

// MarshalSnapshot serializes for DB storage.
func MarshalSnapshot(s Snapshot) json.RawMessage {
	b, _ := json.Marshal(s)
	return b
}
