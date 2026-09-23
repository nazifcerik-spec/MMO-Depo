package economy

import "github.com/google/uuid"

// System account UUIDs are deterministic v5 UUIDs derived from a fixed
// namespace and the account name. Same across every environment so ledger
// dumps stay comparable, and no seed migration is needed.
//
// Namespace: uuid.NameSpaceOID (public constant, stable).
var systemNamespace = uuid.NameSpaceOID

func systemUUID(name string) uuid.UUID {
	return uuid.NewSHA1(systemNamespace, []byte("otter.system:"+name))
}

// Well-known system accounts — see docs/economy/transaction.md §Model.
var (
	SystemRewardSource  = systemUUID("reward_source")
	SystemLootSource    = systemUUID("loot_source")
	SystemMintSource    = systemUUID("mint_source")
	SystemVendorSink    = systemUUID("vendor_sink")
	SystemTaxSink       = systemUUID("tax_sink")
	SystemFeeSink       = systemUUID("fee_sink")
	SystemAuctionEscrow = systemUUID("auction_escrow")
	SystemMailEscrow    = systemUUID("mail_escrow")
)

// IsSystemAccount returns true when the id belongs to one of the well-known
// system accounts. Callers use this to skip cap / non-negative checks.
func IsSystemAccount(id uuid.UUID) bool {
	switch id {
	case SystemRewardSource,
		SystemLootSource,
		SystemMintSource,
		SystemVendorSink,
		SystemTaxSink,
		SystemFeeSink,
		SystemAuctionEscrow,
		SystemMailEscrow:
		return true
	}
	return false
}
