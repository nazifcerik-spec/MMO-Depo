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
)

// Strings are the content names per language (DD-230): "<kind>.<code>" →
// display name, from i18n/<lang>.yaml. Version changes whenever any file does.
type Strings struct {
	ByLang  map[string]map[string]string
	Version string
}

// DefaultLang answers languages without a file.
const DefaultLang = "en"

// Lookup returns a language's table, falling back to DefaultLang.
func (s *Strings) Lookup(lang string) map[string]string {
	if s == nil {
		return nil
	}
	if t, ok := s.ByLang[strings.ToLower(lang)]; ok {
		return t
	}
	return s.ByLang[DefaultLang]
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
		lang := strings.TrimSuffix(filepath.Base(f), ".yaml")
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
		for k := range base {
			if _, ok := m[k]; !ok {
				return fmt.Errorf("i18n %s: missing %s", lang, k)
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
