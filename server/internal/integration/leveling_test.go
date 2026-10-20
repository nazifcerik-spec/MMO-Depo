//go:build integration

package integration

import (
	"context"
	"sync"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/otter/server/internal/domain/character"
	"github.com/otter/server/internal/inventory"
	"github.com/otter/server/internal/leveling"
)

// warrior creates a Warrior fixture at (spec, level) — deterministic test
// characters for DD-250 §39.
func (w *testEnv) warrior(spec string, level int16) player {
	w.t.Helper()
	p := w.player("warrior", level)
	if spec != "warrior" {
		_, err := w.p.Exec(w.ctx, `UPDATE character.character SET spec_code = $2 WHERE id = $1`, p.ID, spec)
		require.NoError(w.t, err)
	}
	return p
}

func (w *testEnv) grant(p player, amount int64, src string) *character.XPGrantResult {
	w.t.Helper()
	res, err := w.chars.GrantXP(w.ctx, p.ID, amount, character.XPSource{Kind: "test", ID: src})
	require.NoError(w.t, err)
	return res
}

// toLevel grants exactly the XP from the start of `from` to the start of `to`.
func toLevel(from, to int) int64 {
	cfg := leveling.Current()
	var n int64
	for l := from; l < to; l++ {
		n += cfg.XPToNext(l)
	}
	return n
}

func (w *testEnv) events(p player, kind string) []string {
	w.t.Helper()
	rows, err := w.p.Query(w.ctx, `SELECT ref FROM character.progression_event WHERE character_id = $1 AND kind = $2 ORDER BY id`, p.ID, kind)
	require.NoError(w.t, err)
	defer rows.Close()
	var out []string
	for rows.Next() {
		var ref string
		require.NoError(w.t, rows.Scan(&ref))
		out = append(out, ref)
	}
	return out
}

func (w *testEnv) titles(p player) map[string]bool {
	w.t.Helper()
	rows, err := w.p.Query(w.ctx, `SELECT title_code FROM progression.character_title WHERE character_id = $1`, p.ID)
	require.NoError(w.t, err)
	defer rows.Close()
	out := map[string]bool{}
	for rows.Next() {
		var c string
		require.NoError(w.t, rows.Scan(&c))
		out[c] = true
	}
	return out
}

func (w *testEnv) char(p player) *character.Character {
	w.t.Helper()
	c, err := w.chars.GetOwnedInternal(w.ctx, p.ID)
	require.NoError(w.t, err)
	return c
}

func TestLevelUpSingleMultipleOverflow(t *testing.T) {
	w := newWorld(t)
	cfg := leveling.Current()
	p := w.warrior("warrior", 1)
	res := w.grant(p, cfg.XPToNext(1)+5, "a")
	require.Equal(t, []int{2}, res.LevelsGained)
	require.Equal(t, int64(5), res.Character.XP)
	require.Equal(t, []string{"2"}, w.events(p, character.EventLevelUp))
	require.Equal(t, []string{"2"}, w.events(p, character.EventAttributePoints), "a point every two levels")

	q := w.warrior("warrior", 24)
	_, err := w.p.Exec(w.ctx, `UPDATE character.character SET xp = $2 WHERE id = $1`, q.ID, cfg.XPToNext(24)*9/10)
	require.NoError(t, err)
	res = w.grant(q, 40_000, "quest-1")
	require.GreaterOrEqual(t, res.Character.Level, int16(27), "one reward, several levels")
	require.Len(t, res.LevelsGained, int(res.Character.Level-24))
	require.Less(t, res.Character.XP, cfg.XPToNext(int(res.Character.Level)), "overflow carried, not discarded")
}

// Level 100: the gate holds the extra XP; choosing the promotion pays it out.
func TestPromotionGateOverflowAndSelection(t *testing.T) {
	w := newWorld(t)
	p := w.warrior("warrior", 99)
	res := w.grant(p, toLevel(99, 104)+100, "afk-3h")
	require.Equal(t, int16(100), res.Character.Level, "stays at the gate")
	require.True(t, res.PromotionAvailable)
	require.Equal(t, toLevel(100, 104)+100, res.Character.OverflowXP, "the rest is stored")
	require.Equal(t, []string{"tier:warrior"}, w.events(p, character.EventPromotionAvail))
	require.True(t, w.titles(p)["veteran"])

	sum, err := w.chars.ProgressionSummary(w.ctx, p.ID, p.Account)
	require.NoError(t, err)
	require.True(t, sum.Promotion.Available)
	require.ElementsMatch(t, []string{"executioner", "guardian"}, []string{sum.Promotion.Choices[0].Code, sum.Promotion.Choices[1].Code})

	_, err = w.chars.SelectPromotion(w.ctx, p.ID, p.Account, "iron_fortress")
	require.ErrorIs(t, err, character.ErrInvalidPromotion, "not a tier-1 choice")
	_, err = w.chars.SelectPromotion(w.ctx, p.ID, p.Account, "mage")
	require.ErrorIs(t, err, character.ErrInvalidPromotion)
	_, err = w.chars.SelectPromotion(w.ctx, p.ID, uuid.New(), "guardian")
	require.Error(t, err, "someone else's character")

	g, err := w.chars.SelectPromotion(w.ctx, p.ID, p.Account, "guardian")
	require.NoError(t, err)
	require.Equal(t, "guardian", g.Character.Spec())
	require.Equal(t, int16(104), g.Character.Level, "overflow processed at once")
	require.Equal(t, int64(100), g.Character.XP)
	require.Zero(t, g.Character.OverflowXP)
	require.Contains(t, w.events(p, character.EventSkillUnlocked), "guardian_shield_slam")
	require.Contains(t, w.events(p, character.EventPassiveUnlocked), "guardian_unyielding")
	_, err = w.chars.SelectPromotion(w.ctx, p.ID, p.Account, "executioner")
	require.ErrorIs(t, err, character.ErrInvalidPromotion, "one choice per gate")

	// Level 300: the final specialization, same mechanism.
	q := w.warrior("guardian", 299)
	res = w.grant(q, toLevel(299, 300)+50, "b")
	require.Equal(t, int16(300), res.Character.Level)
	require.Equal(t, int64(50), res.Character.OverflowXP)
	_, err = w.chars.SelectPromotion(w.ctx, q.ID, q.Account, "berserker")
	require.ErrorIs(t, err, character.ErrInvalidPromotion, "berserker is an Executioner branch")
	g, err = w.chars.SelectPromotion(w.ctx, q.ID, q.Account, "iron_fortress")
	require.NoError(t, err)
	require.Equal(t, int64(50), g.Character.XP)
	require.Contains(t, w.events(q, character.EventPassiveUnlocked), "iron_fortress_bulwark")
	require.True(t, w.titles(q)["champion"])
}

// 600 awakening, 700 ultimate, 850 capstone, 1000 mastery + the cap.
func TestLateMilestonesAndCap(t *testing.T) {
	w := newWorld(t)
	p := w.warrior("iron_fortress", 599)
	w.grant(p, toLevel(599, 600), "a")
	require.Equal(t, []string{"iron_fortress_aegis_rebound"}, w.events(p, character.EventAwakening))
	require.True(t, w.titles(p)["mythic"])

	w.grant(p, toLevel(600, 700), "b")
	require.Equal(t, []string{"iron_fortress_unbreakable"}, w.events(p, character.EventUltimate))
	w.grant(p, toLevel(700, 850), "c")
	require.Equal(t, []string{"capstone"}, w.events(p, character.EventCapstone))

	res := w.grant(p, toLevel(850, 1000)+999_999, "d")
	require.Equal(t, int16(1000), res.Character.Level)
	require.Zero(t, res.Character.XP)
	require.Equal(t, int64(999_999), res.Character.PostCapXP, "post-cap XP is kept")
	require.Equal(t, []string{"iron_fortress"}, w.events(p, character.EventMastery))
	tt := w.titles(p)
	require.True(t, tt["the_unbreakable"] && tt["eternal"] && tt["immortal"])

	res = w.grant(p, 10_000_000, "e")
	require.Equal(t, int16(1000), res.Character.Level, "never 1001")
	require.Empty(t, res.LevelsGained)
	_, err := w.p.Exec(w.ctx, `UPDATE character.character SET level = 1001 WHERE id = $1`, p.ID)
	require.Error(t, err, "the column refuses 1001")

	sum, err := w.chars.ProgressionSummary(w.ctx, p.ID, p.Account)
	require.NoError(t, err)
	require.True(t, sum.Maxed)
	require.True(t, sum.Mastery.Unlocked)
	require.Equal(t, "the_unbreakable", sum.Mastery.Title)
	require.Equal(t, "eternal", sum.Title)
	require.Nil(t, sum.NextMilestone)
}

// The same source pays once however often (or concurrently) it is sent;
// different sources never lose an update.
func TestXPIdempotentAndConcurrent(t *testing.T) {
	w := newWorld(t)
	cfg := leveling.Current()
	p := w.warrior("warrior", 1)
	var wg sync.WaitGroup
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err := w.chars.GrantXP(context.Background(), p.ID, cfg.XPToNext(1), character.XPSource{Kind: "hunt", ID: "contract-1"})
			require.NoError(t, err)
		}()
	}
	wg.Wait()
	c := w.char(p)
	require.Equal(t, int16(2), c.Level, "eight identical reveals pay once")
	require.Equal(t, []string{"2"}, w.events(p, character.EventLevelUp))

	q := w.warrior("warrior", 1)
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			_, err := w.chars.GrantXP(context.Background(), q.ID, 1000, character.XPSource{Kind: "hunt", ID: uuid.NewString()})
			require.NoError(t, err)
		}(i)
	}
	wg.Wait()
	c = w.char(q)
	total := c.XP
	for l := 1; l < int(c.Level); l++ {
		total += cfg.XPToNext(l)
	}
	require.Equal(t, int64(8000), total, "no lost update under concurrency")

	// Replaying a whole level range adds no duplicate reward or title.
	r := w.warrior("warrior", 99)
	w.grant(r, toLevel(99, 100), "x")
	w.grant(r, toLevel(99, 100), "x")
	require.Len(t, w.events(r, character.EventTitleUnlocked), 1)
	require.Len(t, w.events(r, character.EventTalentPoint), 1)
}

func TestAttributesAndTalents(t *testing.T) {
	w := newWorld(t)
	p := w.warrior("guardian", 100)
	cfg := leveling.Current()
	earned := cfg.AttributePointsEarned(100)
	require.Equal(t, 302, earned)

	_, err := w.chars.AllocateAttributes(w.ctx, p.ID, p.Account, map[string]int{"strength": earned + 1})
	require.ErrorIs(t, err, character.ErrNotEnoughPoints)
	_, err = w.chars.AllocateAttributes(w.ctx, p.ID, p.Account, map[string]int{"charm": 1})
	require.ErrorIs(t, err, character.ErrUnknownAttribute)
	c, err := w.chars.AllocateAttributes(w.ctx, p.ID, p.Account, map[string]int{"strength": 200, "vitality": 102})
	require.NoError(t, err)
	require.Equal(t, 200, c.Alloc.Strength)
	_, err = w.chars.AllocateAttributes(w.ctx, p.ID, p.Account, map[string]int{"luck": 1})
	require.ErrorIs(t, err, character.ErrNotEnoughPoints, "all spent")
	sheet, err := w.chars.Sheet(w.ctx, p.ID, p.Account)
	require.NoError(t, err)
	require.Equal(t, 210, sheet.Primary.Strength, "base 10 + 200 allocated")
	require.Equal(t, "rage", sheet.ResourceKind)

	// One talent point at 100.
	rank, err := w.chars.LearnTalent(w.ctx, p.ID, p.Account, "warrior_defense_toughness")
	require.NoError(t, err)
	require.Equal(t, 1, rank)
	_, err = w.chars.LearnTalent(w.ctx, p.ID, p.Account, "warrior_arms_accuracy")
	require.ErrorIs(t, err, character.ErrNotEnoughPoints)
	_, err = w.chars.LearnTalent(w.ctx, p.ID, p.Account, "mage_nothing")
	require.ErrorIs(t, err, character.ErrUnknownTalent)

	// Capstone: level 850 and 25 points in the tree.
	q := w.warrior("iron_fortress", 849)
	for _, code := range []string{"warrior_defense_toughness", "warrior_defense_plating", "warrior_defense_shieldwork",
		"warrior_defense_steadfast", "warrior_defense_riposte"} {
		_, err := w.p.Exec(w.ctx, `INSERT INTO character.talent (character_id, talent_code, rank) VALUES ($1, $2, 5)`, q.ID, code)
		require.NoError(t, err)
	}
	_, err = w.chars.LearnTalent(w.ctx, q.ID, q.Account, "warrior_defense_living_fortress")
	require.ErrorIs(t, err, character.ErrCapstoneLocked, "level 849")
	w.grant(q, toLevel(849, 850), "lv850")
	rank, err = w.chars.LearnTalent(w.ctx, q.ID, q.Account, "warrior_defense_living_fortress")
	require.NoError(t, err)
	require.Equal(t, 1, rank)
	_, err = w.chars.LearnTalent(w.ctx, q.ID, q.Account, "warrior_defense_living_fortress")
	require.ErrorIs(t, err, character.ErrTalentMaxed)
	r := w.warrior("iron_fortress", 900)
	_, err = w.chars.LearnTalent(w.ctx, r.ID, r.Account, "warrior_arms_perfect_technique")
	require.ErrorIs(t, err, character.ErrCapstoneLocked, "an empty tree")

	mods, err := w.chars.ProgressionModifiers(w.ctx, q.ID, character.Sheet{Class: "warrior", Spec: ptr("iron_fortress"), Level: 850})
	require.NoError(t, err)
	require.Len(t, mods, 25+3, "25 ranks + the capstone's three modifiers")
}

func ptr(s string) *string { return &s }

// Respec boundary: refused while busy; gear the new specialization cannot
// use goes back to the bag; level, XP and titles stay.
func TestRespecAndEquipmentRequirements(t *testing.T) {
	w := newWorld(t)
	p := w.warrior("iron_fortress", 600)
	_, err := w.p.Exec(w.ctx, `UPDATE character.character SET attr_str = 210, attr_vit = 290 WHERE id = $1`, p.ID)
	require.NoError(t, err)
	fetch := inventory.ProfileFetcher(func(ctx context.Context, id uuid.UUID) (inventory.Profile, error) {
		c, err := w.chars.GetOwnedInternal(ctx, id)
		if err != nil {
			return inventory.Profile{}, err
		}
		return inventory.Profile{Level: c.Level, Class: c.ClassCode, Path: leveling.Current().Path(c.Spec()),
			Raw: character.RawAttributes(c)}, nil
	})
	equip := inventory.NewEquipmentService(w.inv, w.items, fetch)
	shield := w.give(p, "titan_bulwark_shield", 1)
	slot, err := equip.Equip(w.ctx, p.ID, shield.ID)
	require.NoError(t, err, "guardian branch + raw STR 220 / VIT 300")
	require.Equal(t, "offhand", slot)

	weak := w.warrior("guardian", 600)
	s2 := w.give(weak, "titan_bulwark_shield", 1)
	_, err = equip.Equip(w.ctx, weak.ID, s2.ID)
	require.ErrorIs(t, err, inventory.ErrAttributeRequirement, "raw attributes too low")
	other := w.warrior("berserker", 600)
	_, err = w.p.Exec(w.ctx, `UPDATE character.character SET attr_str = 210, attr_vit = 290 WHERE id = $1`, other.ID)
	require.NoError(t, err)
	s3 := w.give(other, "titan_bulwark_shield", 1)
	_, err = equip.Equip(w.ctx, other.ID, s3.ID)
	require.ErrorIs(t, err, inventory.ErrEquipmentPermission, "berserkers do not use shields (DD-252)")
	unpromoted := w.warrior("warrior", 600)
	_, err = w.p.Exec(w.ctx, `UPDATE character.character SET attr_str = 210, attr_vit = 290 WHERE id = $1`, unpromoted.ID)
	require.NoError(t, err)
	s4 := w.give(unpromoted, "titan_bulwark_shield", 1)
	_, err = equip.Equip(w.ctx, unpromoted.ID, s4.ID)
	require.ErrorIs(t, err, inventory.ErrEquipmentPermission, "shields belong to the Guardian line (DD-272)")

	busy := fakeBusy(true)
	_, _, err = w.chars.ChangeSpecialization(w.ctx, p.ID, "berserker", busy, equip)
	require.ErrorIs(t, err, character.ErrBusy)
	_, _, err = w.chars.ChangeSpecialization(w.ctx, p.ID, "guardian", fakeBusy(false), equip)
	require.ErrorIs(t, err, character.ErrInvalidSpec, "respec keeps the tier")

	c, removed, err := w.chars.ChangeSpecialization(w.ctx, p.ID, "berserker", fakeBusy(false), equip)
	require.NoError(t, err)
	require.Equal(t, "berserker", c.Spec())
	require.Equal(t, 1, removed, "the Guardian shield came off")
	require.Equal(t, int16(600), c.Level)
	require.True(t, w.slotted(shield.ID), "back in the bag")
	require.Contains(t, w.events(p, character.EventPassiveUnlocked), "berserker_blood_fury")
}

type fakeBusy bool

func (b fakeBusy) Busy(context.Context, uuid.UUID) (bool, error) { return bool(b), nil }

// The DD-250 Warrior fixtures read back as expected (§39).
func TestWarriorFixtures(t *testing.T) {
	w := newWorld(t)
	for _, f := range []struct {
		spec      string
		level     int16
		promotion bool
		next      int
		title     string
		mastery   bool
	}{
		{"warrior", 1, false, 10, "", false},
		{"warrior", 99, false, 100, "", false},
		{"warrior", 100, true, 100, "veteran", false},
		{"guardian", 100, false, 150, "veteran", false},
		{"guardian", 299, false, 300, "elite", false},
		{"guardian", 300, true, 300, "champion", false},
		{"iron_fortress", 300, false, 350, "champion", false},
		{"iron_fortress", 599, false, 600, "legendary", false},
		{"iron_fortress", 600, false, 700, "mythic", false},
		{"iron_fortress", 699, false, 700, "mythic", false},
		{"iron_fortress", 700, false, 800, "ancient", false},
		{"iron_fortress", 849, false, 850, "ascendant", false},
		{"iron_fortress", 850, false, 900, "ascendant", false},
		{"iron_fortress", 999, false, 1000, "immortal", false},
		{"iron_fortress", 1000, false, 0, "eternal", true},
		{"warlord", 300, false, 350, "champion", false},
		{"berserker", 300, false, 350, "champion", false},
		{"weapon_master", 300, false, 350, "champion", false},
		{"warlord", 1000, false, 0, "eternal", true},
		{"berserker", 1000, false, 0, "eternal", true},
		{"weapon_master", 1000, false, 0, "eternal", true},
	} {
		p := w.warrior(f.spec, f.level)
		sum, err := w.chars.ProgressionSummary(w.ctx, p.ID, p.Account)
		require.NoError(t, err, "%s %d", f.spec, f.level)
		require.Equal(t, f.promotion, sum.Promotion.Available, "%s %d promotion", f.spec, f.level)
		require.Equal(t, f.title, sum.Title, "%s %d title", f.spec, f.level)
		require.Equal(t, f.mastery, sum.Mastery.Unlocked, "%s %d mastery", f.spec, f.level)
		if f.next == 0 {
			require.Nil(t, sum.NextMilestone, "%s %d", f.spec, f.level)
		} else {
			require.NotNil(t, sum.NextMilestone, "%s %d", f.spec, f.level)
			require.Equal(t, f.next, sum.NextMilestone.Level, "%s %d next milestone", f.spec, f.level)
		}
		require.Equal(t, "rage", sum.Resource.Kind)
		require.Equal(t, leveling.Current().AttributePointsEarned(int(f.level)), sum.Attributes.Available)
	}
}
