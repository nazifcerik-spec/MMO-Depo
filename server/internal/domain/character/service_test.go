package character_test

import (
	"context"
	"io"
	"log/slog"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/otter/server/internal/domain/character"
)

func newTestService(t *testing.T) (*character.Service, *character.MemoryRepo) {
	t.Helper()
	repo := character.NewMemoryRepo()
	svc := character.NewService(repo, character.ServiceConfig{
		SlotCap: 2,
		Now:     time.Now,
	}, slog.New(slog.NewTextHandler(io.Discard, nil)))
	return svc, repo
}

func mustCreate(t *testing.T, svc *character.Service, acct uuid.UUID, name string) *character.Character {
	t.Helper()
	c, err := svc.Create(context.Background(), character.CreateRequest{
		AccountID: acct,
		Name:      name,
		RaceCode:  "human",
		ClassCode: "warrior",
	})
	require.NoError(t, err)
	return c
}

func TestCreate_Success(t *testing.T) {
	svc, _ := newTestService(t)
	acct := uuid.New()

	c := mustCreate(t, svc, acct, "Grim")

	require.Equal(t, "Grim", c.Name)
	require.Equal(t, "grim", c.NameNormalized)
	require.Equal(t, character.StateActive, c.State)
	require.Equal(t, int16(1), c.Level)
	require.Equal(t, int64(0), c.XP)
}

func TestCreate_NormalizesRaceClass(t *testing.T) {
	svc, _ := newTestService(t)
	c, err := svc.Create(context.Background(), character.CreateRequest{
		AccountID: uuid.New(),
		Name:      "Aria",
		RaceCode:  "  ELF  ",
		ClassCode: "MAGE",
	})
	require.NoError(t, err)
	require.Equal(t, "elf", c.RaceCode)
	require.Equal(t, "mage", c.ClassCode)
}

// All ten races can be chosen (DD-267).
func TestCreate_AllRaces(t *testing.T) {
	svc, _ := newTestService(t)
	for i, race := range []string{"human", "elf", "dwarf", "orc", "troll", "undead", "dark_elf", "drakari", "tauren", "gnome"} {
		c, err := svc.Create(context.Background(), character.CreateRequest{
			AccountID: uuid.New(), Name: "Race" + string(rune('a'+i)) + "x", RaceCode: race, ClassCode: "warrior",
		})
		require.NoError(t, err, race)
		require.Equal(t, race, c.RaceCode)
	}
}

func TestCreate_SlotCap(t *testing.T) {
	svc, _ := newTestService(t)
	acct := uuid.New()

	_ = mustCreate(t, svc, acct, "First")
	_ = mustCreate(t, svc, acct, "Second")

	// 3rd exceeds default cap 2.
	_, err := svc.Create(context.Background(), character.CreateRequest{
		AccountID: acct, Name: "Third", RaceCode: "human", ClassCode: "warrior",
	})
	require.ErrorIs(t, err, character.ErrSlotCap)
}

func TestCreate_SlotCapPerAccountIsolated(t *testing.T) {
	svc, _ := newTestService(t)
	acctA := uuid.New()
	acctB := uuid.New()

	_ = mustCreate(t, svc, acctA, "One")
	_ = mustCreate(t, svc, acctA, "Two")

	// Different account not affected.
	c := mustCreate(t, svc, acctB, "Three")
	require.NotEqual(t, uuid.Nil, c.ID)
}

func TestCreate_DuplicateName(t *testing.T) {
	svc, _ := newTestService(t)
	acctA := uuid.New()
	acctB := uuid.New()

	_ = mustCreate(t, svc, acctA, "Grim")

	// Same-name across accounts also rejected (global name unique).
	_, err := svc.Create(context.Background(), character.CreateRequest{
		AccountID: acctB, Name: "GRIM", RaceCode: "elf", ClassCode: "mage",
	})
	require.ErrorIs(t, err, character.ErrNameTaken)
}

func TestCreate_ValidationErrors(t *testing.T) {
	svc, _ := newTestService(t)
	acct := uuid.New()

	cases := []struct {
		name  string
		req   character.CreateRequest
		field string
	}{
		{"bad name", character.CreateRequest{AccountID: acct, Name: "!!", RaceCode: "human", ClassCode: "warrior"}, "name"},
		{"empty name", character.CreateRequest{AccountID: acct, Name: "", RaceCode: "human", ClassCode: "warrior"}, "name"},
		{"unknown race", character.CreateRequest{AccountID: acct, Name: "Grim", RaceCode: "goblin", ClassCode: "warrior"}, "race_code"},
		{"unknown class", character.CreateRequest{AccountID: acct, Name: "Grim", RaceCode: "human", ClassCode: "ninja"}, "class_code"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := svc.Create(context.Background(), tc.req)
			var ve *character.ValidationError
			require.ErrorAs(t, err, &ve)
			require.Equal(t, tc.field, ve.Field)
		})
	}
}

func TestGetOwned_ForbiddenForOtherAccount(t *testing.T) {
	svc, _ := newTestService(t)
	owner := uuid.New()
	stranger := uuid.New()
	c := mustCreate(t, svc, owner, "Grim")

	_, err := svc.GetOwned(context.Background(), c.ID, stranger)
	require.ErrorIs(t, err, character.ErrForbidden)
}

func TestEnter_ReturnsSheet(t *testing.T) {
	svc, _ := newTestService(t)
	acct := uuid.New()
	c := mustCreate(t, svc, acct, "Grim")

	sheet, err := svc.Enter(context.Background(), c.ID, acct)
	require.NoError(t, err)
	require.Equal(t, "Grim", sheet.Name)
	require.Equal(t, "warrior", sheet.Class)
	require.Greater(t, sheet.Resources.HPMax, 0)
	require.Greater(t, sheet.Derived.Attack, 0)
}

func TestEnter_ForbiddenForOtherAccount(t *testing.T) {
	svc, _ := newTestService(t)
	c := mustCreate(t, svc, uuid.New(), "Grim")
	_, err := svc.Enter(context.Background(), c.ID, uuid.New())
	require.ErrorIs(t, err, character.ErrForbidden)
}

// pendingDelete puts a character in delete_grace, as an account deletion does.
func pendingDelete(t *testing.T, repo *character.MemoryRepo, id uuid.UUID) {
	t.Helper()
	require.NoError(t, repo.SetDeleteGrace(context.Background(), id, time.Now().Add(time.Hour)))
}

func TestDelete_RemovesAtOnce(t *testing.T) {
	svc, _ := newTestService(t)
	acct := uuid.New()
	c := mustCreate(t, svc, acct, "Grim")

	require.NoError(t, svc.Delete(context.Background(), c.ID, acct, "Grim"))
	_, err := svc.GetOwned(context.Background(), c.ID, acct)
	require.ErrorIs(t, err, character.ErrNotFound)
	list, err := svc.ListForAccount(context.Background(), acct)
	require.NoError(t, err)
	require.Empty(t, list)
	_ = mustCreate(t, svc, acct, "Grim") // the name is free again
}

func TestDelete_NeedsName(t *testing.T) {
	svc, _ := newTestService(t)
	acct := uuid.New()
	c := mustCreate(t, svc, acct, "Grim")
	for _, wrong := range []string{"", "Gri", "Grimm", "  "} {
		require.ErrorIs(t, svc.Delete(context.Background(), c.ID, acct, wrong), character.ErrNameMismatch, wrong)
	}
	require.NoError(t, svc.Delete(context.Background(), c.ID, acct, " grim "), "case and spacing do not matter")
}

func TestDelete_OthersAndPending(t *testing.T) {
	svc, repo := newTestService(t)
	acct := uuid.New()
	c := mustCreate(t, svc, acct, "Grim")
	require.ErrorIs(t, svc.Delete(context.Background(), c.ID, uuid.New(), "Grim"), character.ErrForbidden)

	pendingDelete(t, repo, c.ID)
	require.NoError(t, svc.Delete(context.Background(), c.ID, acct, "Grim"), "a pending deletion can be done now")
	require.ErrorIs(t, svc.Delete(context.Background(), c.ID, acct, "Grim"), character.ErrNotFound)
}

func TestCancelDelete_Success(t *testing.T) {
	svc, repo := newTestService(t)
	acct := uuid.New()
	c := mustCreate(t, svc, acct, "Grim")
	pendingDelete(t, repo, c.ID)

	restored, err := svc.CancelDelete(context.Background(), c.ID, acct)
	require.NoError(t, err)
	require.Equal(t, character.StateActive, restored.State)
	require.Nil(t, restored.DeleteGraceAt)
}

func TestCancelDelete_NotInGraceRejected(t *testing.T) {
	svc, _ := newTestService(t)
	acct := uuid.New()
	c := mustCreate(t, svc, acct, "Grim")

	_, err := svc.CancelDelete(context.Background(), c.ID, acct)
	require.ErrorIs(t, err, character.ErrNotInDeleteGrace)
}

func TestRename_Success(t *testing.T) {
	svc, _ := newTestService(t)
	acct := uuid.New()
	c := mustCreate(t, svc, acct, "Grim")

	renamed, err := svc.Rename(context.Background(), c.ID, acct, "Grimlock")
	require.NoError(t, err)
	require.Equal(t, "Grimlock", renamed.Name)
	require.Equal(t, "grimlock", renamed.NameNormalized)
}

func TestRename_NameTaken(t *testing.T) {
	svc, _ := newTestService(t)
	acct := uuid.New()
	_ = mustCreate(t, svc, acct, "Existing")
	c2 := mustCreate(t, svc, acct, "Grim")

	_, err := svc.Rename(context.Background(), c2.ID, acct, "existing")
	require.ErrorIs(t, err, character.ErrNameTaken)
}

func TestRename_NotPlayableRejected(t *testing.T) {
	svc, repo := newTestService(t)
	acct := uuid.New()
	c := mustCreate(t, svc, acct, "Grim")
	pendingDelete(t, repo, c.ID)

	_, err := svc.Rename(context.Background(), c.ID, acct, "NewName")
	require.ErrorIs(t, err, character.ErrNotPlayable)
}

// DD-344, DD-346: a new character is male with the first options unless it
// picks; the look can change later and a new body resets its faces and hair.
func TestAppearance(t *testing.T) {
	svc, _ := newTestService(t)
	ctx := context.Background()
	acct := uuid.New()
	c := mustCreate(t, svc, acct, "Grim")
	require.Equal(t, "male", c.Body)
	require.Equal(t, []int16{1, 1, 1, 1}, []int16{c.FaceID, c.HairStyleID, c.HairColorID, c.EyeColorID})

	got, err := svc.SetAppearance(ctx, c.ID, acct, character.Appearance{FaceID: 3, HairStyleID: 4, HairColorID: 2, EyeColorID: 5})
	require.NoError(t, err)
	require.Equal(t, []int16{3, 4, 2, 5}, []int16{got.FaceID, got.HairStyleID, got.HairColorID, got.EyeColorID})
	again, err := svc.GetOwned(ctx, c.ID, acct)
	require.NoError(t, err)
	require.Equal(t, int16(5), again.EyeColorID, "stored")

	got, err = svc.SetAppearance(ctx, c.ID, acct, character.Appearance{Body: "Female"})
	require.NoError(t, err)
	require.Equal(t, "female", got.Body)
	require.Equal(t, []int16{1, 1, 2, 5}, []int16{got.FaceID, got.HairStyleID, got.HairColorID, got.EyeColorID},
		"another body's face and hair do not carry over; colours do")
	got, err = svc.SetAppearance(ctx, c.ID, acct, character.Appearance{Body: "male", FaceID: 2})
	require.NoError(t, err)
	require.Equal(t, []int16{2, 1}, []int16{got.FaceID, got.HairStyleID}, "a chosen face stays")

	var verr *character.ValidationError
	for name, look := range map[string]character.Appearance{
		"body":  {Body: "dragon"},
		"face":  {FaceID: character.FaceOptions + 1},
		"style": {HairStyleID: -1},
		"hair":  {HairColorID: character.HairColorOptions + 1},
		"eyes":  {EyeColorID: character.EyeColorOptions + 1},
	} {
		_, err = svc.SetAppearance(ctx, c.ID, acct, look)
		require.ErrorAs(t, err, &verr, name)
	}
	_, err = svc.SetAppearance(ctx, c.ID, uuid.New(), character.Appearance{Body: "male"})
	require.Error(t, err, "only the owner")

	f, err := svc.Create(ctx, character.CreateRequest{AccountID: acct, Name: "Lira", RaceCode: "human", ClassCode: "warrior",
		Look: character.Appearance{Body: "female", FaceID: 4, HairStyleID: 2, HairColorID: 3, EyeColorID: 1}})
	require.NoError(t, err)
	require.Equal(t, "female", f.Body)
	require.Equal(t, []int16{4, 2, 3, 1}, []int16{f.FaceID, f.HairStyleID, f.HairColorID, f.EyeColorID})
	_, err = svc.Create(ctx, character.CreateRequest{AccountID: uuid.New(), Name: "Nope", RaceCode: "human", ClassCode: "warrior",
		Look: character.Appearance{FaceID: 99}})
	require.ErrorAs(t, err, &verr)
}

func TestListForAccount(t *testing.T) {
	svc, _ := newTestService(t)
	acct := uuid.New()
	_ = mustCreate(t, svc, acct, "One")
	_ = mustCreate(t, svc, acct, "Two")

	list, err := svc.ListForAccount(context.Background(), acct)
	require.NoError(t, err)
	require.Len(t, list, 2)
}

func TestOwns(t *testing.T) {
	svc, _ := newTestService(t)
	owner := uuid.New()
	stranger := uuid.New()
	c := mustCreate(t, svc, owner, "Grim")

	require.NoError(t, svc.Owns(context.Background(), c.ID, owner))
	require.ErrorIs(t, svc.Owns(context.Background(), c.ID, stranger), character.ErrForbidden)
	require.ErrorIs(t, svc.Owns(context.Background(), uuid.New(), owner), character.ErrNotFound)
}

func TestCreate_EmailVerificationGate(t *testing.T) {
	svc, _ := newTestService(t)
	verified := map[uuid.UUID]bool{}
	svc.WithEmailVerification(func(_ context.Context, id uuid.UUID) (bool, error) { return verified[id], nil })
	acct := uuid.New()
	req := character.CreateRequest{AccountID: acct, Name: "Gated", RaceCode: "human", ClassCode: "mage"}
	_, err := svc.Create(context.Background(), req)
	require.ErrorIs(t, err, character.ErrEmailNotVerified)
	verified[acct] = true
	_, err = svc.Create(context.Background(), req)
	require.NoError(t, err)
}

func TestValidationError_WrapsSentinel(t *testing.T) {
	err := error(&character.ValidationError{Field: "name", Message: "too short"})
	require.ErrorIs(t, err, character.ErrValidation)
	require.Contains(t, err.Error(), "name too short")
}
