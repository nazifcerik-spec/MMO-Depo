package item

// Rarity is the 6-tier item quality axis. Higher rarity → higher value + affix
// tier + drop rate lower. Client displays as color + border stroke.
type Rarity string

const (
	RarityCommon    Rarity = "common"
	RarityUncommon  Rarity = "uncommon"
	RarityRare      Rarity = "rare"
	RarityEpic      Rarity = "epic"
	RarityLegendary Rarity = "legendary"
	RarityMythic    Rarity = "mythic"
)

// ValueMultiplier scales the template's base value for rarity. Placeholder v1
// numbers — balance pass Phase 02.
func ValueMultiplier(r Rarity) float64 {
	switch r {
	case RarityCommon:
		return 1.0
	case RarityUncommon:
		return 1.5
	case RarityRare:
		return 3.0
	case RarityEpic:
		return 8.0
	case RarityLegendary:
		return 25.0
	case RarityMythic:
		return 100.0
	default:
		return 1.0
	}
}

// AffixRule describes whether a rarity gets prefix/suffix at drop time.
//
//	common/uncommon → no affix
//	rare            → 50% chance of ONE affix (prefix OR suffix)
//	epic            → guaranteed ONE affix; 50% chance both
//	legendary       → guaranteed BOTH
//	mythic          → guaranteed BOTH + reserved slot for Phase 02 unique effect
type AffixRule struct {
	HasPrefix      bool
	HasSuffix      bool
	OptionalPrefix bool // rolled probabilistically
	OptionalSuffix bool
}

func AffixRuleFor(r Rarity) AffixRule {
	switch r {
	case RarityRare:
		return AffixRule{OptionalPrefix: true, OptionalSuffix: true}
	case RarityEpic:
		return AffixRule{HasPrefix: true, OptionalSuffix: true}
	case RarityLegendary, RarityMythic:
		return AffixRule{HasPrefix: true, HasSuffix: true}
	default:
		return AffixRule{}
	}
}
