package quest

import (
	"context"
	"math"
	"strings"

	"github.com/google/uuid"

	"github.com/otter/server/internal/activity"
	"github.com/otter/server/internal/crafting"
	"github.com/otter/server/internal/leveling"
	"github.com/otter/server/internal/progression"
)

// Level-scaled rewards (DD-318). Daily and side quests stay worth doing at
// any level: their XP is a share of the character's current level and
// their gold is that of a number of level-fitting kills (10 × L^1.5 copper
// each, as the daily reward counts). Each gathering objective also pays
// GatherXPPct of the gathering profession's current level. Main quests keep
// their content rewards (their chain is placed by level already).
const (
	DailyXPPct     = 10
	SideXPPct      = 25
	DailyGoldKills = 5
	SideGoldKills  = 15
	GatherXPPct    = 20
)

// ProfessionReward is profession XP a quest pays.
type ProfessionReward struct {
	Kind string `json:"kind"`
	XP   int64  `json:"xp"`
}

// Reward is what turning a quest in pays a character now.
type Reward struct {
	XP          int64              `json:"reward_xp"`
	GoldCopper  int64              `json:"reward_gold_copper"`
	Professions []ProfessionReward `json:"reward_professions,omitempty"`
}

func goldPerKill(level int16) float64 { return 10 * math.Pow(float64(max(level, 1)), 1.5) }

// RewardFor is t's reward for a character of level with gathering
// profession levels prof (kind → level; missing = 1).
func RewardFor(t *Template, level int16, prof map[string]int16) Reward {
	r := Reward{XP: t.RewardXP, GoldCopper: t.RewardGold}
	var xpPct, kills float64
	switch t.Category {
	case CategoryDaily:
		xpPct, kills = DailyXPPct, DailyGoldKills
	case CategorySide:
		xpPct, kills = SideXPPct, SideGoldKills
	}
	if xpPct > 0 {
		r.XP = int64(math.Round(float64(leveling.Current().XPToNext(int(max(level, 1)))) * xpPct / 100))
		r.GoldCopper = int64(math.Round(goldPerKill(level) * kills))
	}
	for _, o := range t.Objectives {
		kind, ok := strings.CutPrefix(o.Event, progression.EventGatherPrefix)
		if !ok || kind == "" {
			continue
		}
		pl := max(prof[kind], 1)
		r.Professions = append(r.Professions, ProfessionReward{Kind: kind,
			XP: int64(math.Round(float64(crafting.ProfessionXPToNext(pl)) * GatherXPPct / 100))})
	}
	return r
}

// GatherXP pays gathering profession XP (gathering.Service satisfies).
type GatherXP interface {
	AddXP(ctx context.Context, characterID uuid.UUID, kind activity.Kind, xp int) (activity.GatherProgress, error)
}

// WithGatherXP lets quests pay gathering profession XP (DD-318).
func (s *Service) WithGatherXP(g GatherXP) *Service { s.gather = g; return s }

// professionLevels is the character's gathering profession levels.
func (s *Service) professionLevels(ctx context.Context, char uuid.UUID) (map[string]int16, error) {
	rows, err := s.pool.Query(ctx, `SELECT profession_code, level FROM crafting.character_profession
		WHERE character_id = $1`, char)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := map[string]int16{}
	for rows.Next() {
		var code string
		var level int16
		if err := rows.Scan(&code, &level); err != nil {
			return nil, err
		}
		out[code] = level
	}
	return out, rows.Err()
}
