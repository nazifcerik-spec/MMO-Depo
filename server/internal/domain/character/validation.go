package character

import (
	"regexp"
	"strings"
)

// namePattern accepts 3-16 characters, ASCII letters + digits + underscore.
// Unicode + spaces are Phase 02 (locale + display-name policy needs legal review).
var namePattern = regexp.MustCompile(`^[A-Za-z][A-Za-z0-9_]{2,15}$`)

// V1 race + class allowlists — hardcoded until content pack ships.
// Faz 01 sprint 8+ replaces these with content.race / content.class reads.
var (
	knownRaces = map[string]bool{
		"human": true,
		"elf":   true,
		"dwarf": true,
		"orc":   true,
		// DD-267: six more races (looks only, like the first four).
		"troll":    true,
		"undead":   true,
		"dark_elf": true,
		"drakari":  true,
		"tauren":   true,
		"gnome":    true,
	}
	// The ten base classes (DD-252); promotions are content (progression.yaml).
	knownClasses = map[string]bool{
		"warrior": true,
		"rogue":   true,
		"ranger":  true,
		"mage":    true,
		"monk":    true,
		"cleric":  true,
		"paladin": true,
		"druid":   true,
		"bard":    true,
		"shaman":  true,
	}
)

func validateName(raw string) (string, string, error) {
	s := strings.TrimSpace(raw)
	if s == "" {
		return "", "", &ValidationError{Field: "name", Message: "required"}
	}
	if !namePattern.MatchString(s) {
		return "", "", &ValidationError{Field: "name", Message: "3-16 chars, ASCII letters/digits/_ (must start with a letter)"}
	}
	return s, strings.ToLower(s), nil
}

func validateRace(raw string) (string, error) {
	s := strings.ToLower(strings.TrimSpace(raw))
	if s == "" {
		return "", &ValidationError{Field: "race_code", Message: "required"}
	}
	if !knownRaces[s] {
		return "", &ValidationError{Field: "race_code", Message: "unknown race"}
	}
	return s, nil
}

func validateClass(raw string) (string, error) {
	s := strings.ToLower(strings.TrimSpace(raw))
	if s == "" {
		return "", &ValidationError{Field: "class_code", Message: "required"}
	}
	if !knownClasses[s] {
		return "", &ValidationError{Field: "class_code", Message: "unknown class"}
	}
	return s, nil
}

// KnownClass reports whether code is a playable class (content validation).
func KnownClass(code string) bool { return knownClasses[code] }
