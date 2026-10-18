package contentpack

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/otter/server/internal/leveling"
)

// Strings are the content names per language (DD-230): "<kind>.<code>" →
// display name, from i18n/<lang>.yaml. Version changes whenever any file does.
type Strings struct {
	ByLang  map[string]map[string]string
	Version string
}

// DefaultLang answers languages without a file.
const DefaultLang = "en"

// CompleteLangs must name every content code; other languages (zh-cn, es,
// DD-250) may be partial — they must cover the progression keys and fall
// back to DefaultLang for the rest.
var CompleteLangs = map[string]bool{"en": true, "tr": true}

// Lookup returns a language's table ("zh" finds "zh-cn"), falling back to
// DefaultLang.
func (s *Strings) Lookup(lang string) map[string]string {
	if s == nil {
		return nil
	}
	lang = strings.ToLower(lang)
	if t, ok := s.ByLang[lang]; ok {
		return t
	}
	want, _, _ := strings.Cut(lang, "-") // "es-es" → "es", "zh" → "zh"
	if t, ok := s.ByLang[want]; ok {
		return t
	}
	for code, t := range s.ByLang {
		if base, _, _ := strings.Cut(code, "-"); base == want {
			return t
		}
	}
	return s.ByLang[DefaultLang]
}

// fillFallbacks completes partial languages with DefaultLang names.
func (s *Strings) fillFallbacks() {
	base := s.ByLang[DefaultLang]
	for lang, m := range s.ByLang {
		if CompleteLangs[lang] {
			continue
		}
		for k, v := range base {
			if _, ok := m[k]; !ok {
				m[k] = v
			}
		}
	}
}

func loadStrings(dir string) (*Strings, error) {
	files, err := filepath.Glob(filepath.Join(dir, "*.yaml"))
	if err != nil {
		return nil, err
	}
	sort.Strings(files)
	out := &Strings{ByLang: map[string]map[string]string{}}
	h := sha256.New()
	for _, f := range files {
		raw, err := os.ReadFile(f)
		if err != nil {
			return nil, fmt.Errorf("read %s: %w", f, err)
		}
		var m map[string]string
		if err := yaml.Unmarshal(raw, &m); err != nil {
			return nil, fmt.Errorf("parse %s: %w", f, err)
		}
		lang := strings.ToLower(strings.TrimSuffix(filepath.Base(f), ".yaml"))
		out.ByLang[lang] = m
		h.Write([]byte(lang))
		h.Write(raw)
	}
	out.Version = hex.EncodeToString(h.Sum(nil))[:16]
	return out, nil
}

// validateStrings: every language has the same keys as the default one, no
// empty names, and every content code the client shows has a name.
func validateStrings(s *Strings, p *Pack) error {
	if len(s.ByLang) == 0 {
		return nil
	}
	base, ok := s.ByLang[DefaultLang]
	if !ok {
		return fmt.Errorf("i18n: %s.yaml missing", DefaultLang)
	}
	for lang, m := range s.ByLang {
		for k, v := range m {
			if strings.TrimSpace(v) == "" {
				return fmt.Errorf("i18n %s: %s is empty", lang, k)
			}
			if _, ok := base[k]; !ok {
				return fmt.Errorf("i18n %s: %s not in %s.yaml", lang, k, DefaultLang)
			}
		}
		if CompleteLangs[lang] {
			for k := range base {
				if _, ok := m[k]; !ok {
					return fmt.Errorf("i18n %s: missing %s", lang, k)
				}
			}
		}
		for _, k := range progressionKeys(p) {
			if _, ok := m[k]; !ok {
				return fmt.Errorf("i18n %s: missing progression key %s", lang, k)
			}
		}
	}
	var keys []string
	add := func(kind string, codes ...string) {
		for _, c := range codes {
			keys = append(keys, kind+"."+c)
		}
	}
	for _, x := range p.Items {
		add("item", x.Code)
	}
	for _, x := range p.Mobs {
		add("mob", x.Code)
	}
	for _, x := range p.Dungeons {
		add("dungeon", x.Code)
	}
	for _, x := range p.Regions {
		add("region", x.Code)
	}
	for _, x := range p.Professions {
		add("profession", x.Code)
	}
	for _, x := range p.Vendors {
		add("vendor", x.Code)
	}
	for _, x := range p.Quests {
		add("quest", x.Code)
	}
	for _, x := range p.Skills {
		add("skill", x.Code)
	}
	var missing []string
	for _, k := range keys {
		if _, ok := base[k]; !ok {
			missing = append(missing, k)
		}
	}
	if len(missing) > 0 {
		return errors.New("i18n: no name for " + strings.Join(missing, ", "))
	}
	return nil
}

// progressionKeys are the player-facing progression strings every language
// must carry (DD-250): class nodes and roles, level and mastery titles,
// talents and trees, the skills of progression classes (name + description)
// and the attribute names.
func progressionKeys(p *Pack) []string {
	c := p.Progression
	if c == nil {
		return nil
	}
	var keys []string
	classes := map[string]bool{}
	for _, k := range c.Classes {
		classes[k.Code] = true
		keys = append(keys, "class."+k.Code, "class."+k.Code+".role")
		if k.Mastery != nil {
			keys = append(keys, "title."+k.Mastery.Title)
		}
	}
	for _, t := range c.Titles {
		keys = append(keys, "title."+t)
	}
	for _, t := range c.Trees {
		keys = append(keys, "talent_tree."+t.Code)
		for _, tl := range t.Talents {
			keys = append(keys, "talent."+tl.Code, "talent."+tl.Code+".desc")
		}
	}
	for _, s := range p.Skills {
		if classes[s.Class] {
			keys = append(keys, "skill."+s.Code, "skill."+s.Code+".desc")
		}
	}
	for kind := range leveling.ResourceKinds {
		keys = append(keys, "resource."+kind)
	}
	for _, k := range c.Classes {
		if k.Mastery != nil && k.Mastery.Achievement != "" {
			keys = append(keys, "achievement."+k.Mastery.Achievement)
		}
	}
	for _, a := range leveling.Attributes {
		keys = append(keys, "attribute."+a)
	}
	return keys
}
