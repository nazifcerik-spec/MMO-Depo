package contentpack

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

// Quest mirrors content.quest_template (DD-221).
type Quest struct {
	Code             string           `yaml:"code"`
	DisplayNameKey   string           `yaml:"display_name_key"`
	Category         string           `yaml:"category"`
	MinLevel         int16            `yaml:"min_level"`
	RequiresQuest    string           `yaml:"requires_quest"`
	Objectives       []QuestObjective `yaml:"objectives"`
	RewardXP         int32            `yaml:"reward_xp"`
	RewardGoldCopper int64            `yaml:"reward_gold_copper"`
	SortOrder        int32            `yaml:"sort_order"`
}

// QuestObjective counts one progression event kind.
type QuestObjective struct {
	Event string `yaml:"event" json:"event"`
	Count int64  `yaml:"count" json:"count"`
}

type questsFile struct {
	Quests []Quest `yaml:"quests"`
}

var (
	questCategories = map[string]bool{"main": true, "side": true, "daily": true}
	gatherKinds     = map[string]bool{"mining": true, "fishing": true, "woodcutting": true, "herbalism": true, "hunting": true}
)

func loadQuestsFile(path string) ([]Quest, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}
		return nil, fmt.Errorf("read %s: %w", path, err)
	}
	var f questsFile
	if err := yaml.Unmarshal(raw, &f); err != nil {
		return nil, fmt.Errorf("parse %s: %w", path, err)
	}
	return f.Quests, nil
}

// refs are the content codes quest objectives may point at.
type refs struct{ mobs, dungeons, items map[string]bool }

// validateQuests checks shape, the chain (requires_quest) and that every
// `<event>:<code>` objective names content that exists in this pack.
func validateQuests(quests []Quest, r refs) error {
	byCode := make(map[string]Quest, len(quests))
	for i, q := range quests {
		if q.Code == "" {
			return fmt.Errorf("quest[%d]: empty code", i)
		}
		if _, dup := byCode[q.Code]; dup {
			return fmt.Errorf("quest[%d] %s: duplicate code", i, q.Code)
		}
		byCode[q.Code] = q
		if q.DisplayNameKey == "" {
			return fmt.Errorf("quest %s: display_name_key required", q.Code)
		}
		if !questCategories[q.Category] {
			return fmt.Errorf("quest %s: category %q not in main|side|daily", q.Code, q.Category)
		}
		if q.MinLevel < 1 || q.MinLevel > 1000 { // the character cap (DD-278)
			return fmt.Errorf("quest %s: min_level must be 1..1000", q.Code)
		}
		if q.RewardXP < 0 || q.RewardGoldCopper < 0 {
			return fmt.Errorf("quest %s: rewards must be >= 0", q.Code)
		}
		if q.Category == "daily" && q.RequiresQuest != "" {
			return fmt.Errorf("quest %s: dailies cannot require another quest", q.Code)
		}
		if len(q.Objectives) == 0 || len(q.Objectives) > 5 {
			return fmt.Errorf("quest %s: 1..5 objectives required", q.Code)
		}
		seen := map[string]bool{}
		for j, o := range q.Objectives {
			if o.Count < 1 {
				return fmt.Errorf("quest %s: objectives[%d].count must be >= 1", q.Code, j)
			}
			if seen[o.Event] {
				return fmt.Errorf("quest %s: objectives[%d] duplicate event %q", q.Code, j, o.Event)
			}
			seen[o.Event] = true
			if err := validateQuestEvent(o.Event, r); err != nil {
				return fmt.Errorf("quest %s: objectives[%d]: %w", q.Code, j, err)
			}
		}
	}
	for _, q := range quests {
		if q.RequiresQuest == "" {
			continue
		}
		req, ok := byCode[q.RequiresQuest]
		if !ok {
			return fmt.Errorf("quest %s: requires unknown quest %q", q.Code, q.RequiresQuest)
		}
		if req.Category == "daily" {
			return fmt.Errorf("quest %s: cannot require daily %q", q.Code, q.RequiresQuest)
		}
		// Walk the chain to reject cycles.
		for cur, steps := q, 0; cur.RequiresQuest != ""; steps++ {
			if steps > len(quests) || cur.RequiresQuest == q.Code {
				return fmt.Errorf("quest %s: requires_quest cycle", q.Code)
			}
			cur = byCode[cur.RequiresQuest]
		}
	}
	return nil
}

func validateQuestEvent(ev string, r refs) error {
	kind, code, scoped := strings.Cut(ev, ":")
	check := func(set map[string]bool, what string) error {
		if !scoped {
			return nil
		}
		if !set[code] {
			return fmt.Errorf("event %q: unknown %s %q", ev, what, code)
		}
		return nil
	}
	switch kind {
	case "mob_killed":
		return check(r.mobs, "mob")
	case "dungeon_cleared":
		return check(r.dungeons, "dungeon")
	case "item_crafted":
		return check(r.items, "item")
	case "boss_killed":
		if scoped {
			return fmt.Errorf("event %q: boss_killed takes no code", ev)
		}
		return nil
	}
	if g, ok := strings.CutPrefix(kind, "gather_"); ok && gatherKinds[g] && !scoped {
		return nil
	}
	return fmt.Errorf("event %q: unknown kind", ev)
}

func codeSet[T any](list []T, code func(T) string) map[string]bool {
	out := make(map[string]bool, len(list))
	for _, v := range list {
		out[code(v)] = true
	}
	return out
}

func (l *Loader) upsertQuests(ctx context.Context, quests []Quest) error {
	tx, err := l.DB.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	const q = `
INSERT INTO content.quest_template (
  code, display_name_key, category, min_level, requires_quest_code,
  objectives, reward_xp, reward_gold_copper, sort_order
) VALUES ($1,$2,$3,$4,NULLIF($5,''),$6,$7,$8,$9)
ON CONFLICT (code) DO UPDATE SET
  display_name_key    = EXCLUDED.display_name_key,
  category            = EXCLUDED.category,
  min_level           = EXCLUDED.min_level,
  requires_quest_code = EXCLUDED.requires_quest_code,
  objectives          = EXCLUDED.objectives,
  reward_xp           = EXCLUDED.reward_xp,
  reward_gold_copper  = EXCLUDED.reward_gold_copper,
  sort_order          = EXCLUDED.sort_order,
  version             = content.quest_template.version + 1`
	for _, qu := range quests {
		obj, err := json.Marshal(qu.Objectives)
		if err != nil {
			return fmt.Errorf("quest %s: marshal objectives: %w", qu.Code, err)
		}
		if _, err := tx.Exec(ctx, q, qu.Code, qu.DisplayNameKey, qu.Category, qu.MinLevel, qu.RequiresQuest,
			obj, qu.RewardXP, qu.RewardGoldCopper, qu.SortOrder); err != nil {
			return fmt.Errorf("upsert quest %s: %w", qu.Code, err)
		}
	}
	return tx.Commit(ctx)
}
