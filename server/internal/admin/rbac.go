// Package admin is the operator backend (Sprint 26, docs/admin/*, DD-213):
// admin accounts with password + TOTP, role-based permissions, step-up for
// critical actions, full audit, and the /internal/ops moderation actions.
package admin

// Permissions (docs/admin/auth-rbac.md §4).
const (
	PermPlayerRead         = "player.read"
	PermPlayerPIIReveal    = "player.pii_reveal"
	PermPlayerMute24h      = "player.mute<=24h"
	PermPlayerBan7d        = "player.ban<=7d"
	PermPlayerBanPermanent = "player.ban_permanent"
	PermCharacterRead      = "character.read"
	PermCharacterFreeze    = "character.freeze"
	PermLedgerRead         = "ledger.read"
	PermLedgerReversal100k = "ledger.reversal<=100k_gold"
	PermLedgerReversalAny  = "ledger.reversal_any"
	PermMarketCancel       = "market.listing_cancel"
	PermGuildDisband       = "guild.disband_override"
	PermEventTrigger       = "content.event_trigger"
	PermAuditRead          = "audit.read"
	PermRoleGrant          = "admin.role_grant"
	PermSupportTickets     = "support.tickets"   // DD-234
	PermClientErrors       = "ops.client_errors" // DD-235
)

var readAll = []string{PermPlayerRead, PermCharacterRead, PermLedgerRead}

// RolePermissions is the v1 fixed matrix.
var RolePermissions = map[string][]string{
	"viewer":        readAll,
	"support_l1":    append(append([]string{}, readAll...), PermPlayerPIIReveal, PermPlayerMute24h, PermSupportTickets),
	"support_l2":    append(append([]string{}, readAll...), PermPlayerPIIReveal, PermPlayerMute24h, PermPlayerBan7d, PermCharacterFreeze, PermLedgerReversal100k, PermMarketCancel, PermAuditRead, PermSupportTickets, PermClientErrors),
	"game_designer": append(append([]string{}, readAll...), PermEventTrigger),
	"ops_engineer":  append(append([]string{}, readAll...), PermPlayerPIIReveal, PermCharacterFreeze, PermEventTrigger, PermAuditRead, PermClientErrors),
	"admin": {
		PermPlayerRead, PermPlayerPIIReveal, PermPlayerMute24h, PermPlayerBan7d, PermPlayerBanPermanent,
		PermCharacterRead, PermCharacterFreeze, PermLedgerRead, PermLedgerReversal100k, PermLedgerReversalAny,
		PermMarketCancel, PermGuildDisband, PermEventTrigger, PermAuditRead, PermRoleGrant, PermSupportTickets, PermClientErrors,
	},
}

// PermissionSet is the union of the given roles' permissions.
type PermissionSet map[string]bool

func Permissions(roles []string) PermissionSet {
	set := PermissionSet{}
	for _, r := range roles {
		for _, p := range RolePermissions[r] {
			set[p] = true
		}
	}
	return set
}

func (p PermissionSet) Has(perm string) bool { return p[perm] }

// ReversalLimitCopper is 100k gold (1 gold = 10 000 copper).
const ReversalLimitCopper = 100_000 * 10_000

// CanReverse reports whether a reversal moving `goldCopper` is allowed.
func (p PermissionSet) CanReverse(goldCopper int64) bool {
	if p.Has(PermLedgerReversalAny) {
		return true
	}
	return p.Has(PermLedgerReversal100k) && goldCopper <= ReversalLimitCopper
}

// CanBan reports whether a ban of `days` (0 = permanent) is allowed.
func (p PermissionSet) CanBan(days int) bool {
	if days <= 0 {
		return p.Has(PermPlayerBanPermanent)
	}
	return p.Has(PermPlayerBanPermanent) || (p.Has(PermPlayerBan7d) && days <= 7)
}

// CanMute reports whether a mute of `hours` (0 = permanent) is allowed.
// Mutes above 24h are reserved for the permanent-ban holders.
func (p PermissionSet) CanMute(hours int) bool {
	if hours > 0 && hours <= 24 && p.Has(PermPlayerMute24h) {
		return true
	}
	return p.Has(PermPlayerBanPermanent)
}

// ValidRole reports whether code is a known role.
func ValidRole(code string) bool { _, ok := RolePermissions[code]; return ok }
