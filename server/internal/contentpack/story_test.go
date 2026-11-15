package contentpack

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/otter/server/internal/story"
)

// The shipped story (DD-367): every scene text exists in en + tr, every story
// quest is offered by one scene, and the reputation a choice grants names a
// seeded faction (migration 0021).
func TestShippedStory(t *testing.T) {
	p, err := Read("../../../content")
	require.NoError(t, err)
	require.NotEmpty(t, p.Story.Scenes)
	require.NoError(t, story.Validate(p.Story, storyQuestSet(p.Quests)))
	require.NoError(t, p.StoryText.Validate(p.Story))
	factions := map[string]bool{"ironhold_guard": true, "greenwood_wardens": true, "ember_covenant": true}
	for _, sc := range p.Story.Scenes {
		for _, o := range sc.Options {
			for f := range o.Reputation {
				require.True(t, factions[f], "%s.%s: faction %s", sc.Code, o.Code, f)
			}
		}
	}
	// Each story quest opens at its scene's level, so it can be taken there.
	level := map[string]int16{}
	for _, q := range p.Quests {
		level[q.Code] = q.MinLevel
	}
	for _, sc := range p.Story.Scenes {
		for _, q := range append(append([]string{}, sc.Quests...), sc.OptionalQuests...) {
			require.LessOrEqual(t, int(level[q]), sc.MinLevel, "%s: quest %s", sc.Code, q)
		}
	}
}

func TestStoryFileRejectsNonStoryQuests(t *testing.T) {
	dir := t.TempDir()
	path := dir + "/story.yaml"
	require.NoError(t, writeFile(path, "chapters: []\nquests:\n  - code: x\n    category: main\n"))
	_, _, err := loadStoryFile(path)
	require.ErrorContains(t, err, "category story")
	b, q, err := loadStoryFile(dir + "/none.yaml")
	require.NoError(t, err)
	require.Empty(t, b.Scenes)
	require.Empty(t, q)
}

func writeFile(path, s string) error { return os.WriteFile(path, []byte(s), 0o600) }
