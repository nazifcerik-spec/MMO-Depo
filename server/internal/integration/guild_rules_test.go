//go:build integration

package integration

import (
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/otter/server/internal/domain/character"
	"github.com/otter/server/internal/guild"
)

// Nobody acts through a character of another account, in any guild action.
func TestGuildRefusesForeignCharacter(t *testing.T) {
	w := newWorld(t)
	svc := guild.NewService(w.p, w.chars, w.txn, nil)
	owner, stranger := w.player("warrior", 5), w.player("mage", 5)
	other, acct := owner.ID, stranger.Account

	calls := map[string]func() error{
		"mine":    func() error { _, err := svc.Mine(w.ctx, other, acct); return err },
		"invites": func() error { _, err := svc.Invites(w.ctx, other, acct); return err },
		"create": func() error {
			_, err := svc.Create(w.ctx, other, acct, guild.CreateRequest{Name: "Nope " + randHex(3), Tag: randTag(), IdempotencyKey: uuid.NewString()})
			return err
		},
		"invite":   func() error { _, err := svc.Invite(w.ctx, other, acct, stranger.Name); return err },
		"accept":   func() error { _, err := svc.Accept(w.ctx, other, acct, uuid.New()); return err },
		"decline":  func() error { return svc.Decline(w.ctx, other, acct, uuid.New()) },
		"leave":    func() error { return svc.Leave(w.ctx, other, acct) },
		"kick":     func() error { return svc.Kick(w.ctx, other, acct, stranger.ID) },
		"rank":     func() error { return svc.SetRank(w.ctx, other, acct, stranger.ID, guild.RankOfficer) },
		"transfer": func() error { return svc.TransferLeadership(w.ctx, other, acct, stranger.ID) },
		"motd":     func() error { return svc.SetMOTD(w.ctx, other, acct, "hi") },
		"disband":  func() error { _, err := svc.ScheduleDisband(w.ctx, other, acct); return err },
		"vault":    func() error { _, err := svc.Vault(w.ctx, other, acct); return err },
		"deposit":  func() error { return svc.DepositGold(w.ctx, other, acct, 1, uuid.NewString()) },
		"withdraw": func() error { return svc.WithdrawGold(w.ctx, other, acct, 1, uuid.NewString()) },
		"item in":  func() error { return svc.DepositItem(w.ctx, other, acct, uuid.New(), 1) },
		"item out": func() error { return svc.WithdrawItem(w.ctx, other, acct, uuid.New(), 1) },
	}
	for name, call := range calls {
		require.ErrorIs(t, call(), character.ErrForbidden, name)
	}
}

// What a guild refuses by its own rules: who may do what, and to whom.
func TestGuildRules(t *testing.T) {
	w := newWorld(t)
	svc := guild.NewService(w.p, w.chars, w.txn, nil)
	lead, mem, outsider := w.player("warrior", 5), w.player("mage", 5), w.player("rogue", 5)

	create := func(p player, name, tag, desc string) (*guild.Guild, error) {
		return svc.Create(w.ctx, p.ID, p.Account, guild.CreateRequest{Name: name, Tag: tag, Description: desc, IdempotencyKey: uuid.NewString()})
	}
	_, err := create(lead, "Rules "+randHex(3), randTag(), strings.Repeat("x", guild.MaxDescription+1))
	require.ErrorIs(t, err, guild.ErrTooLong)

	w.gold(lead, guild.CreationFeeCopper)
	tag := randTag()
	g, err := create(lead, "Rules "+randHex(3), tag, "")
	require.NoError(t, err)
	_, err = create(lead, "Rules "+randHex(3), randTag(), "")
	require.ErrorIs(t, err, guild.ErrAlreadyInGuild)

	w.gold(outsider, guild.CreationFeeCopper)
	_, err = create(outsider, "Rules "+randHex(3), tag, "")
	require.ErrorIs(t, err, guild.ErrTagTaken)
	require.Equal(t, int64(guild.CreationFeeCopper), w.balance(outsider), "a taken tag costs nothing")

	inv, err := svc.Invite(w.ctx, lead.ID, lead.Account, mem.Name)
	require.NoError(t, err)
	_, err = svc.Invite(w.ctx, lead.ID, lead.Account, mem.Name)
	require.ErrorIs(t, err, guild.ErrAlreadyInvited)
	_, err = svc.Accept(w.ctx, outsider.ID, outsider.Account, inv.ID)
	require.ErrorIs(t, err, guild.ErrNotFound, "someone else's invite")
	_, err = svc.Accept(w.ctx, mem.ID, mem.Account, uuid.New())
	require.ErrorIs(t, err, guild.ErrNotFound)

	// A full guild takes nobody, at the invite or at the accept.
	_, err = w.p.Exec(w.ctx, `UPDATE social.guild SET member_cap = 1 WHERE id = $1`, g.ID)
	require.NoError(t, err)
	_, err = svc.Accept(w.ctx, mem.ID, mem.Account, inv.ID)
	require.ErrorIs(t, err, guild.ErrGuildFull)
	_, err = svc.Invite(w.ctx, lead.ID, lead.Account, outsider.Name)
	require.ErrorIs(t, err, guild.ErrGuildFull)
	_, err = w.p.Exec(w.ctx, `UPDATE social.guild SET member_cap = 50 WHERE id = $1`, g.ID)
	require.NoError(t, err)

	// An expired invite is dead.
	_, err = w.p.Exec(w.ctx, `UPDATE social.guild_invite SET expires_at = $2 WHERE id = $1`, inv.ID, time.Now().Add(-time.Minute))
	require.NoError(t, err)
	_, err = svc.Accept(w.ctx, mem.ID, mem.Account, inv.ID)
	require.ErrorIs(t, err, guild.ErrInviteInvalid)
	_, err = w.p.Exec(w.ctx, `UPDATE social.guild_invite SET expires_at = $2 WHERE id = $1`, inv.ID, time.Now().Add(time.Hour))
	require.NoError(t, err)
	_, err = svc.Accept(w.ctx, mem.ID, mem.Account, inv.ID)
	require.NoError(t, err)

	_, err = svc.Invite(w.ctx, lead.ID, lead.Account, mem.Name)
	require.ErrorIs(t, err, guild.ErrTargetInGuild)

	// Nobody targets themselves.
	require.ErrorIs(t, svc.Kick(w.ctx, lead.ID, lead.Account, lead.ID), guild.ErrSelfTarget)
	require.ErrorIs(t, svc.SetRank(w.ctx, lead.ID, lead.Account, lead.ID, guild.RankOfficer), guild.ErrSelfTarget)
	require.ErrorIs(t, svc.TransferLeadership(w.ctx, lead.ID, lead.Account, lead.ID), guild.ErrSelfTarget)

	// A new member has no leader's rights.
	_, err = svc.Invite(w.ctx, mem.ID, mem.Account, outsider.Name)
	require.ErrorIs(t, err, guild.ErrForbidden)
	require.ErrorIs(t, svc.TransferLeadership(w.ctx, mem.ID, mem.Account, lead.ID), guild.ErrForbidden)
	_, err = svc.ScheduleDisband(w.ctx, mem.ID, mem.Account)
	require.ErrorIs(t, err, guild.ErrForbidden)
	require.ErrorIs(t, svc.DepositItem(w.ctx, mem.ID, mem.Account, uuid.New(), 1), guild.ErrForbidden)
	require.ErrorIs(t, svc.WithdrawItem(w.ctx, mem.ID, mem.Account, uuid.New(), 1), guild.ErrForbidden)

	require.ErrorIs(t, svc.SetMOTD(w.ctx, lead.ID, lead.Account, strings.Repeat("x", guild.MaxMOTD+1)), guild.ErrTooLong)

	// The treasury wants a positive amount and an idempotency key.
	require.ErrorIs(t, svc.WithdrawGold(w.ctx, lead.ID, lead.Account, 0, uuid.NewString()), guild.ErrVaultAmount)
	require.ErrorIs(t, svc.WithdrawGold(w.ctx, lead.ID, lead.Account, 1, ""), guild.ErrIdempotencyMissing)
	require.ErrorIs(t, svc.DepositGold(w.ctx, lead.ID, lead.Account, 1, ""), guild.ErrIdempotencyMissing)

	// A guild on its way out takes nothing in.
	_, err = svc.ScheduleDisband(w.ctx, lead.ID, lead.Account)
	require.NoError(t, err)
	require.ErrorIs(t, svc.DepositGold(w.ctx, lead.ID, lead.Account, 1, uuid.NewString()), guild.ErrVaultClosing)
}
