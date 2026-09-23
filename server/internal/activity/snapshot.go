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

	// Mastery gather bonus frozen at start (DD-220): extra resources per
	// hour, then an additive percentage on the hourly rate. Absent on
	// contracts started before Sprint 40 → zero.
	YieldFlatPerHour float64 `json:"yield_flat_per_hour,omitempty"`
	YieldMult        float64 `json:"yield_mult,omitempty"`
}

// parseSnapshot reads a stored snapshot; a malformed one yields zero bonuses
// (reveal must stay possible).
func parseSnapshot(raw json.RawMessage) Snapshot {
	var s Snapshot
	_ = json.Unmarshal(raw, &s)
	return s
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
