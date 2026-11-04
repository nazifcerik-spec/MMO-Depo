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
		// DD-267: six more races. What each one brings is in progression.yaml (DD-349).
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

// Bodies are the portrait figures (DD-344); empty picks the first.
var Bodies = []string{"male", "female"}

func validateBody(raw string) (string, error) {
	s := strings.ToLower(strings.TrimSpace(raw))
	if s == "" {
		return Bodies[0], nil
	}
	for _, b := range Bodies {
		if s == b {
			return s, nil
		}
	}
	return "", &ValidationError{Field: "body", Message: "unknown body"}
}

// Customization option counts (DD-346); faces and hair styles are per body.
// New art only needs a higher count here and in the client.
const (
	FaceOptions      = 5
	HairStyleOptions = 5
	HairColorOptions = 5
	EyeColorOptions  = 5
)

// Appearance is a character's look (DD-346). A zero id means "unchanged"
// (the default 1 for a new character).
type Appearance struct {
	Body                                         string
	FaceID, HairStyleID, HairColorID, EyeColorID int16
}

func validateAppearanceID(field string, id int16, count int16) error {
	if id < 0 || id > count {
		return &ValidationError{Field: field, Message: "unknown option"}
	}
	return nil
}

// apply puts the requested look on c. Faces and hair styles belong to a body,
// so changing the body resets the ones the request does not choose.
func (a Appearance) apply(c *Character) error {
	body := c.Body
	if a.Body != "" {
		var err error
		if body, err = validateBody(a.Body); err != nil {
			return err
		}
	}
	for _, f := range []struct {
		field     string
		id, count int16
	}{
		{"face_id", a.FaceID, FaceOptions},
		{"hair_style_id", a.HairStyleID, HairStyleOptions},
		{"hair_color_id", a.HairColorID, HairColorOptions},
		{"eye_color_id", a.EyeColorID, EyeColorOptions},
	} {
		if err := validateAppearanceID(f.field, f.id, f.count); err != nil {
			return err
		}
	}
	if body != c.Body {
		c.FaceID, c.HairStyleID = 1, 1
	}
	c.Body = body
	set := func(dst *int16, id int16) {
		if id > 0 {
			*dst = id
		} else if *dst == 0 {
			*dst = 1
		}
	}
	set(&c.FaceID, a.FaceID)
	set(&c.HairStyleID, a.HairStyleID)
	set(&c.HairColorID, a.HairColorID)
	set(&c.EyeColorID, a.EyeColorID)
	return nil
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

// KnownRace reports whether a code is a playable race.
func KnownRace(code string) bool { return knownRaces[code] }
