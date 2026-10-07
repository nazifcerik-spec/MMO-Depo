//go:build integration

package integration

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/otter/server/internal/economy"
	"github.com/otter/server/internal/social"
)

func TestSocialFriendsAndBlocks(t *testing.T) {
	w := newWorld(t)
	svc := social.NewService(w.p, w.chars, w.txn, w.inv, nil)
	a, b, c := w.player("warrior", 1), w.player("mage", 1), w.player("rogue", 1)

	_, err := svc.RequestFriend(w.ctx, a.ID, a.Account, a.Name)
	require.ErrorIs(t, err, social.ErrSelfTarget)
	_, err = svc.RequestFriend(w.ctx, a.ID, a.Account, "Nobody"+randHex(3))
	require.ErrorIs(t, err, social.ErrTargetNotFound)
	_, err = svc.RequestFriend(w.ctx, a.ID, b.Account, b.Name)
	require.Error(t, err, "foreign account")

	accepted, err := svc.RequestFriend(w.ctx, a.ID, a.Account, b.Name)
	require.NoError(t, err)
	require.False(t, accepted)
	_, err = svc.RequestFriend(w.ctx, a.ID, a.Account, b.Name)
	require.ErrorIs(t, err, social.ErrAlreadyRequested)
	accepted, err = svc.RequestFriend(w.ctx, b.ID, b.Account, a.Name)
	require.NoError(t, err)
	require.True(t, accepted, "a mutual request auto-accepts")
	_, err = svc.RequestFriend(w.ctx, a.ID, a.Account, b.Name)
	require.ErrorIs(t, err, social.ErrAlreadyFriends)
	friends, err := svc.Friends(w.ctx, a.ID, a.Account)
	require.NoError(t, err)
	require.Equal(t, b.Name, friends[0].Name)

	// c → a: decline, then again and accept.
	_, err = svc.RequestFriend(w.ctx, c.ID, c.Account, a.Name)
	require.NoError(t, err)
	reqs, err := svc.FriendRequests(w.ctx, a.ID, a.Account)
	require.NoError(t, err)
	require.Len(t, reqs, 1)
	require.NoError(t, svc.DeclineFriend(w.ctx, a.ID, a.Account, reqs[0].ID))
	require.ErrorIs(t, svc.AcceptFriend(w.ctx, a.ID, a.Account, reqs[0].ID), social.ErrNotFound)
	_, err = svc.RequestFriend(w.ctx, c.ID, c.Account, a.Name)
	require.NoError(t, err)
	reqs, err = svc.FriendRequests(w.ctx, a.ID, a.Account)
	require.NoError(t, err)
	require.NoError(t, svc.AcceptFriend(w.ctx, a.ID, a.Account, reqs[0].ID))
	friends, err = svc.Friends(w.ctx, a.ID, a.Account)
	require.NoError(t, err)
	require.Len(t, friends, 2)
	require.NoError(t, svc.RemoveFriend(w.ctx, a.ID, a.Account, c.ID))
	require.ErrorIs(t, svc.RemoveFriend(w.ctx, a.ID, a.Account, c.ID), social.ErrNotFriends)

	// Blocking drops the friendship and stops requests/whispers both ways.
	blocked, err := svc.Block(w.ctx, a.ID, a.Account, b.Name)
	require.NoError(t, err)
	require.Equal(t, b.ID, blocked)
	friends, err = svc.Friends(w.ctx, a.ID, a.Account)
	require.NoError(t, err)
	require.Empty(t, friends)
	_, err = svc.RequestFriend(w.ctx, b.ID, b.Account, a.Name)
	require.ErrorIs(t, err, social.ErrBlocked)
	_, err = svc.SendChat(w.ctx, b.ID, b.Account, "whisper", a.Name, "hi")
	require.ErrorIs(t, err, social.ErrBlocked)
	list, err := svc.BlockList(w.ctx, a.ID, a.Account)
	require.NoError(t, err)
	require.Equal(t, b.Name, list[0].Name)
	require.NoError(t, svc.Unblock(w.ctx, a.ID, a.Account, b.ID))
	_, err = svc.RequestFriend(w.ctx, b.ID, b.Account, a.Name)
	require.NoError(t, err)
}

func TestSocialChat(t *testing.T) {
	w := newWorld(t)
	ev := &events{}
	svc := social.NewService(w.p, w.chars, w.txn, w.inv, nil).WithPublisher(ev)
	a, b := w.player("warrior", 1), w.player("mage", 1)

	_, err := svc.SendChat(w.ctx, a.ID, a.Account, "trade", "", "hi")
	require.ErrorIs(t, err, social.ErrBadChannel)
	_, err = svc.SendChat(w.ctx, a.ID, a.Account, "global", "", "   ")
	require.ErrorIs(t, err, social.ErrEmpty)
	_, err = svc.SendChat(w.ctx, a.ID, a.Account, "guild", "", "hi")
	require.ErrorIs(t, err, social.ErrNotInChannel)
	_, err = svc.ChatHistory(w.ctx, a.ID, a.Account, "party", "")
	require.ErrorIs(t, err, social.ErrNotInChannel)

	body := "hello " + randHex(4)
	msg, err := svc.SendChat(w.ctx, a.ID, a.Account, "global", "", body)
	require.NoError(t, err)
	require.Equal(t, a.Name, msg.FromName)
	require.True(t, ev.has("chat.message"))
	hist, err := svc.ChatHistory(w.ctx, b.ID, b.Account, "global", "")
	require.NoError(t, err)
	require.Contains(t, bodies(hist), body)

	w1, err := svc.SendChat(w.ctx, a.ID, a.Account, "whisper", b.Name, "psst")
	require.NoError(t, err)
	require.Equal(t, b.ID, *w1.TargetID)
	_, err = svc.SendChat(w.ctx, a.ID, a.Account, "whisper", a.Name, "me")
	require.ErrorIs(t, err, social.ErrSelfTarget)
	whispers, err := svc.ChatHistory(w.ctx, b.ID, b.Account, "whisper", a.Name)
	require.NoError(t, err)
	require.Contains(t, bodies(whispers), "psst")

	// b blocks a: a's global lines disappear from b's history.
	_, err = svc.Block(w.ctx, b.ID, b.Account, a.Name)
	require.NoError(t, err)
	hist, err = svc.ChatHistory(w.ctx, b.ID, b.Account, "global", "")
	require.NoError(t, err)
	require.NotContains(t, bodies(hist), body)

	// A muted account cannot talk.
	var admin uuid.UUID
	require.NoError(t, w.p.QueryRow(w.ctx, `
		INSERT INTO ops.admin_user (email, password_hash, totp_secret) VALUES ($1, 'x', 'x') RETURNING id`,
		"it-"+randHex(4)+"@example.com").Scan(&admin))
	_, err = w.p.Exec(w.ctx, `
		INSERT INTO ops.account_sanction (account_id, kind, reason, created_by_admin_id) VALUES ($1, 'mute', 'it', $2)`,
		a.Account, admin)
	require.NoError(t, err)
	_, err = svc.SendChat(w.ctx, a.ID, a.Account, "global", "", "still here?")
	require.ErrorIs(t, err, social.ErrMuted)
}

func bodies(ms []social.ChatMessage) []string {
	out := make([]string, len(ms))
	for i, m := range ms {
		out[i] = m.Body
	}
	return out
}

func TestSocialMailEscrowClaimDelete(t *testing.T) {
	w := newWorld(t)
	ev := &events{}
	svc := social.NewService(w.p, w.chars, w.txn, w.inv, nil).WithPublisher(ev)
	a, b := w.player("warrior", 1), w.player("mage", 1)
	w.gold(a, 10_000)
	potion := w.give(a, "potion_hp_small", 3)

	send := func(req social.SendRequest) (*social.SendResult, error) {
		if req.IdempotencyKey == "" {
			req.IdempotencyKey = uuid.NewString()
		}
		return svc.SendMail(w.ctx, a.ID, a.Account, req)
	}
	_, err := svc.SendMail(w.ctx, a.ID, a.Account, social.SendRequest{RecipientName: b.Name, Subject: "x"})
	require.ErrorIs(t, err, social.ErrIdempotencyMissing)
	_, err = send(social.SendRequest{RecipientName: a.Name, Subject: "x"})
	require.ErrorIs(t, err, social.ErrSelfTarget)
	_, err = send(social.SendRequest{RecipientName: b.Name, Subject: "x", GoldCopper: -1})
	require.ErrorIs(t, err, social.ErrAttachmentInvalid)
	_, err = send(social.SendRequest{RecipientName: b.Name, Subject: "x", ItemIDs: []uuid.UUID{potion.ID, potion.ID}})
	require.ErrorIs(t, err, social.ErrAttachmentInvalid, "duplicate attachment")
	_, err = send(social.SendRequest{RecipientName: b.Name, Subject: "x", ItemIDs: []uuid.UUID{uuid.New()}})
	require.ErrorIs(t, err, social.ErrAttachmentInvalid, "not your item")
	_, err = send(social.SendRequest{RecipientName: b.Name, Subject: "x", GoldCopper: 1_000_000})
	require.ErrorIs(t, err, economy.ErrInsufficientBalance)

	key := uuid.NewString()
	res, err := send(social.SendRequest{RecipientName: b.Name, Subject: "gift", Body: "enjoy",
		GoldCopper: 2_000, ItemIDs: []uuid.UUID{potion.ID}, IdempotencyKey: key})
	require.NoError(t, err)
	postage := social.Postage(1, 2_000)
	require.Equal(t, postage, res.Postage)
	require.Equal(t, int64(10_000-2_000)-postage, w.balance(a), "gold + postage left the sender")
	owner, locked := w.owner(potion.ID)
	require.Equal(t, a.ID, owner)
	require.True(t, locked, "item escrowed in the mail")
	require.True(t, ev.has("mail.received"))

	again, err := send(social.SendRequest{RecipientName: b.Name, Subject: "gift", Body: "enjoy",
		GoldCopper: 2_000, ItemIDs: []uuid.UUID{potion.ID}, IdempotencyKey: key})
	require.NoError(t, err)
	require.Equal(t, res.MailID, again.MailID, "replay returns the same mail")
	require.Equal(t, int64(10_000-2_000)-postage, w.balance(a), "replay charges nothing")

	inbox, err := svc.Inbox(w.ctx, b.ID, b.Account)
	require.NoError(t, err)
	require.Equal(t, res.MailID, inbox[0].ID)
	require.ErrorIs(t, svc.DeleteMail(w.ctx, b.ID, b.Account, res.MailID), social.ErrHasAttachments)
	m, err := svc.ReadMail(w.ctx, b.ID, b.Account, res.MailID)
	require.NoError(t, err)
	require.Equal(t, "enjoy", m.Body)
	require.Equal(t, "potion_hp_small", m.Attachments[0].TemplateCode)
	_, err = svc.ReadMail(w.ctx, a.ID, a.Account, res.MailID)
	require.ErrorIs(t, err, social.ErrNotFound, "only the recipient reads it")

	_, err = svc.ClaimMail(w.ctx, b.ID, b.Account, res.MailID)
	require.NoError(t, err)
	require.Equal(t, int64(2_000), w.balance(b))
	owner, locked = w.owner(potion.ID)
	require.Equal(t, b.ID, owner)
	require.False(t, locked)
	require.True(t, w.slotted(potion.ID), "claimed item lands in a bag slot")
	_, err = svc.ClaimMail(w.ctx, b.ID, b.Account, res.MailID)
	require.ErrorIs(t, err, social.ErrNothingToClaim)
	require.Equal(t, int64(2_000), w.balance(b), "no double claim")
	require.NoError(t, svc.DeleteMail(w.ctx, b.ID, b.Account, res.MailID))
	require.ErrorIs(t, svc.DeleteMail(w.ctx, b.ID, b.Account, res.MailID), social.ErrNotFound)

	// A recipient who blocked the sender silently drops the mail.
	_, err = svc.Block(w.ctx, b.ID, b.Account, a.Name)
	require.NoError(t, err)
	before := w.balance(a)
	dropped, err := send(social.SendRequest{RecipientName: b.Name, Subject: "spam", GoldCopper: 100})
	require.NoError(t, err)
	require.Equal(t, uuid.Nil, dropped.MailID)
	require.Equal(t, before, w.balance(a), "dropped mail costs nothing")
}

func TestSocialMailExpiryReturnsValue(t *testing.T) {
	w := newWorld(t)
	svc := social.NewService(w.p, w.chars, w.txn, w.inv, nil)
	a, b := w.player("warrior", 1), w.player("mage", 1)
	w.gold(a, 5_000)
	sword := w.give(a, "sword_iron", 1)
	res, err := svc.SendMail(w.ctx, a.ID, a.Account, social.SendRequest{RecipientName: b.Name, Subject: "late",
		GoldCopper: 1_000, ItemIDs: []uuid.UUID{sword.ID}, IdempotencyKey: uuid.NewString()})
	require.NoError(t, err)
	empty, err := svc.SendMail(w.ctx, a.ID, a.Account, social.SendRequest{RecipientName: b.Name, Subject: "note",
		IdempotencyKey: uuid.NewString()})
	require.NoError(t, err)
	afterSend := w.balance(a)

	_, err = w.p.Exec(w.ctx, `UPDATE social.mail SET expires_at = now() - interval '1 minute' WHERE id = ANY($1)`,
		[]uuid.UUID{res.MailID, empty.MailID})
	require.NoError(t, err)
	n, err := social.ExpireMail(w.ctx, w.p, w.txn, w.inv, 500)
	require.NoError(t, err)
	require.GreaterOrEqual(t, n, int64(2))

	require.Equal(t, afterSend+1_000, w.balance(a), "escrowed gold returns to the sender")
	owner, locked := w.owner(sword.ID)
	require.Equal(t, a.ID, owner)
	require.False(t, locked, "escrowed item returns to the sender")
	require.True(t, w.slotted(sword.ID), "returned item is back in a bag slot (visible)")
	for _, id := range []uuid.UUID{res.MailID, empty.MailID} {
		var state string
		require.NoError(t, w.p.QueryRow(w.ctx, `SELECT state FROM social.mail WHERE id = $1`, id).Scan(&state))
		require.Contains(t, []string{"returned", "expired"}, state)
	}
	_, err = svc.ClaimMail(w.ctx, b.ID, b.Account, res.MailID)
	require.ErrorIs(t, err, social.ErrNothingToClaim)
}

// A claim whose ledger write committed but whose mail/item update rolled back
// must not strand the escrowed items when the mail later expires.
func TestSocialMailExpiryAfterHalfClaim(t *testing.T) {
	w := newWorld(t)
	svc := social.NewService(w.p, w.chars, w.txn, w.inv, nil)
	a, b := w.player("warrior", 1), w.player("mage", 1)
	w.gold(a, 5_000)
	sword := w.give(a, "sword_iron", 1)
	res, err := svc.SendMail(w.ctx, a.ID, a.Account, social.SendRequest{RecipientName: b.Name, Subject: "half",
		GoldCopper: 700, ItemIDs: []uuid.UUID{sword.ID}, IdempotencyKey: uuid.NewString()})
	require.NoError(t, err)

	// Simulate the half-done claim: the gold already moved in the ledger.
	_, err = w.txn.Write(w.ctx, economy.Intent{
		Kind: "mail_claim", RefKind: "mail", RefID: res.MailID, IdempotencyKey: "mail-claim-" + res.MailID.String(),
		Entries: []economy.Entry{
			{Currency: "gold", AccountKind: economy.AccountSystem, AccountRef: economy.SystemMailEscrow, Delta: -700},
			{Currency: "gold", AccountKind: economy.AccountCharacter, AccountRef: b.ID, Delta: 700},
		},
	})
	require.NoError(t, err)
	_, err = w.p.Exec(w.ctx, `UPDATE social.mail SET expires_at = now() - interval '1 minute' WHERE id = $1`, res.MailID)
	require.NoError(t, err)
	_, err = social.ExpireMail(w.ctx, w.p, w.txn, w.inv, 500)
	require.NoError(t, err)

	require.Equal(t, int64(700), w.balance(b), "gold stays with the recipient who claimed it")
	owner, locked := w.owner(sword.ID)
	require.Equal(t, b.ID, owner, "the claim is completed for the items too")
	require.False(t, locked, "items must never stay locked to a dead mail")
	require.True(t, w.slotted(sword.ID), "and visible in the recipient's bag")
}
