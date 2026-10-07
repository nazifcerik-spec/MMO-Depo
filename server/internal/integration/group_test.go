//go:build integration

package integration

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/otter/server/internal/economy"
	"github.com/otter/server/internal/guild"
	"github.com/otter/server/internal/party"
)

func TestPartyLifecycle(t *testing.T) {
	w := newWorld(t)
	ev := &events{}
	svc := party.NewService(w.p, w.chars, nil).WithPublisher(ev)
	lead, b, c := w.player("warrior", 3), w.player("mage", 2), w.player("rogue", 1)

	_, err := svc.Get(w.ctx, lead.ID, lead.Account)
	require.ErrorIs(t, err, party.ErrNotInParty)
	_, err = svc.Invite(w.ctx, lead.ID, lead.Account, b.Name)
	require.ErrorIs(t, err, party.ErrNotInParty, "invite needs a party")

	p, err := svc.Create(w.ctx, lead.ID, lead.Account)
	require.NoError(t, err)
	require.Equal(t, lead.ID, p.LeaderCharacterID)
	_, err = svc.Create(w.ctx, lead.ID, lead.Account)
	require.ErrorIs(t, err, party.ErrAlreadyInParty)
	_, err = svc.Invite(w.ctx, lead.ID, lead.Account, lead.Name)
	require.ErrorIs(t, err, party.ErrSelfTarget)
	_, err = svc.Invite(w.ctx, lead.ID, lead.Account, "Nobody"+randHex(3))
	require.ErrorIs(t, err, party.ErrTargetNotFound)

	inv, err := svc.Invite(w.ctx, lead.ID, lead.Account, b.Name)
	require.NoError(t, err)
	require.True(t, ev.has("party.invite"))
	_, err = svc.Invite(w.ctx, lead.ID, lead.Account, b.Name)
	require.ErrorIs(t, err, party.ErrAlreadyInvited)
	invites, err := svc.Invites(w.ctx, b.ID, b.Account)
	require.NoError(t, err)
	require.Equal(t, lead.Name, invites[0].InvitedByName)
	_, err = svc.Accept(w.ctx, c.ID, c.Account, inv.ID)
	require.ErrorIs(t, err, party.ErrNotFound, "someone else's invite does not exist for them")
	p, err = svc.Accept(w.ctx, b.ID, b.Account, inv.ID)
	require.NoError(t, err)
	require.Len(t, p.Members, 2)
	_, err = svc.Accept(w.ctx, b.ID, b.Account, inv.ID)
	require.Error(t, err, "an invite is used once")

	// Only the leader manages the party.
	_, err = svc.Invite(w.ctx, b.ID, b.Account, c.Name)
	require.ErrorIs(t, err, party.ErrNotLeader)
	require.ErrorIs(t, svc.Kick(w.ctx, b.ID, b.Account, lead.ID), party.ErrNotLeader)
	require.ErrorIs(t, svc.Promote(w.ctx, lead.ID, lead.Account, c.ID), party.ErrNotMember)

	inv2, err := svc.Invite(w.ctx, lead.ID, lead.Account, c.Name)
	require.NoError(t, err)
	require.NoError(t, svc.Decline(w.ctx, c.ID, c.Account, inv2.ID))
	_, err = svc.Accept(w.ctx, c.ID, c.Account, inv2.ID)
	require.ErrorIs(t, err, party.ErrInviteInvalid)
	inv3, err := svc.Invite(w.ctx, lead.ID, lead.Account, c.Name)
	require.NoError(t, err)
	_, err = svc.Accept(w.ctx, c.ID, c.Account, inv3.ID)
	require.NoError(t, err)

	require.NoError(t, svc.Kick(w.ctx, lead.ID, lead.Account, c.ID))
	_, err = svc.Get(w.ctx, c.ID, c.Account)
	require.ErrorIs(t, err, party.ErrNotInParty)

	require.NoError(t, svc.Promote(w.ctx, lead.ID, lead.Account, b.ID))
	p, err = svc.Get(w.ctx, lead.ID, lead.Account)
	require.NoError(t, err)
	require.Equal(t, b.ID, p.LeaderCharacterID)

	// The leader leaving hands over; the last one out disbands.
	require.NoError(t, svc.Leave(w.ctx, b.ID, b.Account))
	p, err = svc.Get(w.ctx, lead.ID, lead.Account)
	require.NoError(t, err)
	require.Equal(t, lead.ID, p.LeaderCharacterID)
	require.Len(t, p.Members, 1)
	require.NoError(t, svc.Leave(w.ctx, lead.ID, lead.Account))
	_, err = svc.Get(w.ctx, lead.ID, lead.Account)
	require.ErrorIs(t, err, party.ErrNotInParty)

	// Disband by the leader removes everyone.
	_, err = svc.Create(w.ctx, lead.ID, lead.Account)
	require.NoError(t, err)
	inv4, err := svc.Invite(w.ctx, lead.ID, lead.Account, c.Name)
	require.NoError(t, err)
	_, err = svc.Accept(w.ctx, c.ID, c.Account, inv4.ID)
	require.NoError(t, err)
	require.ErrorIs(t, svc.Disband(w.ctx, c.ID, c.Account), party.ErrNotLeader)
	require.NoError(t, svc.Disband(w.ctx, lead.ID, lead.Account))
	_, err = svc.Get(w.ctx, c.ID, c.Account)
	require.ErrorIs(t, err, party.ErrNotInParty)
}

func TestPartyFull(t *testing.T) {
	w := newWorld(t)
	svc := party.NewService(w.p, w.chars, nil)
	lead := w.player("warrior", 1)
	_, err := svc.Create(w.ctx, lead.ID, lead.Account)
	require.NoError(t, err)
	for i := 1; i < party.SizeCap; i++ {
		m := w.player("mage", 1)
		inv, err := svc.Invite(w.ctx, lead.ID, lead.Account, m.Name)
		require.NoError(t, err)
		_, err = svc.Accept(w.ctx, m.ID, m.Account, inv.ID)
		require.NoError(t, err)
	}
	extra := w.player("rogue", 1)
	_, err = svc.Invite(w.ctx, lead.ID, lead.Account, extra.Name)
	require.ErrorIs(t, err, party.ErrPartyFull)
}

func TestGuildLifecycle(t *testing.T) {
	w := newWorld(t)
	ev := &events{}
	svc := guild.NewService(w.p, w.chars, w.txn, nil).WithPublisher(ev)
	lead, off, mem := w.player("warrior", 5), w.player("mage", 5), w.player("rogue", 5)
	name, tag := "Wolves "+randHex(3), "W"+randHex(2)[:3]

	create := func(p player, n, tg, key string) (*guild.Guild, error) {
		return svc.Create(w.ctx, p.ID, p.Account, guild.CreateRequest{Name: n, Tag: tg, IdempotencyKey: key})
	}
	_, err := create(lead, name, tag, "")
	require.ErrorIs(t, err, guild.ErrIdempotencyMissing)
	_, err = create(lead, "x", tag, uuid.NewString())
	require.ErrorIs(t, err, guild.ErrInvalidName)
	_, err = create(lead, name, "!", uuid.NewString())
	require.ErrorIs(t, err, guild.ErrInvalidTag)
	_, err = create(lead, name, tag, uuid.NewString())
	require.ErrorIs(t, err, economy.ErrInsufficientBalance, "10 gold fee")

	w.gold(lead, guild.CreationFeeCopper+500)
	key := uuid.NewString()
	g, err := create(lead, name, tag, key)
	require.NoError(t, err)
	require.Equal(t, int64(500), w.balance(lead))
	again, err := create(lead, name, tag, key)
	require.NoError(t, err)
	require.Equal(t, g.ID, again.ID, "replay returns the same guild")
	require.Equal(t, int64(500), w.balance(lead), "fee charged once")

	w.gold(off, guild.CreationFeeCopper)
	_, err = create(off, name, "Z"+randHex(2)[:3], uuid.NewString())
	require.ErrorIs(t, err, guild.ErrNameTaken)
	require.Equal(t, int64(guild.CreationFeeCopper), w.balance(off), "a lost name race refunds the fee")

	pub, err := svc.Public(w.ctx, g.ID)
	require.NoError(t, err)
	require.Equal(t, tag, pub.Tag)

	join := func(p player) {
		inv, err := svc.Invite(w.ctx, lead.ID, lead.Account, p.Name)
		require.NoError(t, err)
		invs, err := svc.Invites(w.ctx, p.ID, p.Account)
		require.NoError(t, err)
		require.Equal(t, name, invs[0].GuildName)
		_, err = svc.Accept(w.ctx, p.ID, p.Account, inv.ID)
		require.NoError(t, err)
	}
	join(off)
	join(mem)
	require.True(t, ev.has("guild.invite"))

	// Fresh members start as initiates and can't manage anything.
	require.ErrorIs(t, svc.SetMOTD(w.ctx, mem.ID, mem.Account, "hi"), guild.ErrForbidden)
	require.ErrorIs(t, svc.Kick(w.ctx, mem.ID, mem.Account, off.ID), guild.ErrForbidden)
	_, err = svc.Invite(w.ctx, mem.ID, mem.Account, "Nobody")
	require.Error(t, err)

	require.NoError(t, svc.SetRank(w.ctx, lead.ID, lead.Account, off.ID, guild.RankOfficer))
	require.ErrorIs(t, svc.SetRank(w.ctx, off.ID, off.Account, mem.ID, guild.RankOfficer), guild.ErrForbidden,
		"officers cannot create officers")
	require.NoError(t, svc.SetRank(w.ctx, off.ID, off.Account, mem.ID, guild.RankVeteran))
	require.NoError(t, svc.SetMOTD(w.ctx, off.ID, off.Account, "raid at 9"))
	mine, err := svc.Mine(w.ctx, mem.ID, mem.Account)
	require.NoError(t, err)
	require.Equal(t, "raid at 9", mine.MOTD)
	require.Equal(t, guild.RankVeteran, *mine.MyRank)

	require.ErrorIs(t, svc.Leave(w.ctx, lead.ID, lead.Account), guild.ErrLeaderCannotLeave)
	require.NoError(t, svc.TransferLeadership(w.ctx, lead.ID, lead.Account, off.ID))
	mine, err = svc.Mine(w.ctx, lead.ID, lead.Account)
	require.NoError(t, err)
	require.Equal(t, off.ID, mine.LeaderCharacterID)
	require.Equal(t, guild.RankOfficer, *mine.MyRank, "old leader becomes officer")

	require.NoError(t, svc.Kick(w.ctx, off.ID, off.Account, mem.ID))
	_, err = svc.Mine(w.ctx, mem.ID, mem.Account)
	require.ErrorIs(t, err, guild.ErrNotInGuild)
	require.NoError(t, svc.Leave(w.ctx, lead.ID, lead.Account))

	at, err := svc.ScheduleDisband(w.ctx, off.ID, off.Account)
	require.NoError(t, err)
	require.False(t, at.IsZero())
	_, err = svc.ScheduleDisband(w.ctx, off.ID, off.Account)
	require.ErrorIs(t, err, guild.ErrDisbandScheduled)
	require.NoError(t, svc.CancelDisband(w.ctx, off.ID, off.Account))
	require.ErrorIs(t, svc.CancelDisband(w.ctx, off.ID, off.Account), guild.ErrDisbandNotPending)
}
