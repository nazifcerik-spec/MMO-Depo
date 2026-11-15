package story

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/otter/server/internal/http/apierr"
)

// DefaultLang answers languages without a file; CompleteLangs must hold every
// scene key, other languages fall back to DefaultLang key by key.
const DefaultLang = "en"

var CompleteLangs = []string{"en", "tr"}

// Text is the scene text per language (content/story/<lang>.yaml), served
// apart from the content names so the long text is only fetched by readers.
type Text struct {
	ByLang  map[string]map[string]string
	Version string
}

// LoadText reads dir/*.yaml; a missing dir is no text.
func LoadText(dir string) (*Text, error) {
	files, err := filepath.Glob(filepath.Join(dir, "*.yaml"))
	if err != nil {
		return nil, err
	}
	sort.Strings(files)
	t := &Text{ByLang: map[string]map[string]string{}}
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
		t.ByLang[lang] = m
		h.Write([]byte(lang))
		h.Write(raw)
	}
	t.Version = hex.EncodeToString(h.Sum(nil))[:16]
	return t, nil
}

// Validate: complete languages hold every key of b, no language holds a key
// b does not use (a typo would never show) and no text is empty.
func (t *Text) Validate(b *Book) error {
	want := map[string]bool{}
	for _, k := range b.TextKeys() {
		want[k] = true
	}
	if t == nil || len(t.ByLang) == 0 {
		if len(want) > 0 {
			return errors.New("story text: no content/story/<lang>.yaml")
		}
		return nil
	}
	for _, lang := range CompleteLangs {
		m, ok := t.ByLang[lang]
		if !ok {
			return fmt.Errorf("story text: %s.yaml missing", lang)
		}
		var missing []string
		for k := range want {
			if _, ok := m[k]; !ok {
				missing = append(missing, k)
			}
		}
		if len(missing) > 0 {
			sort.Strings(missing)
			return fmt.Errorf("story text %s: missing %s", lang, strings.Join(missing, ", "))
		}
	}
	for lang, m := range t.ByLang {
		for k, v := range m {
			if !want[k] {
				return fmt.Errorf("story text %s: %s is not used by any scene", lang, k)
			}
			if strings.TrimSpace(v) == "" {
				return fmt.Errorf("story text %s: %s is empty", lang, k)
			}
		}
	}
	return nil
}

// lookup returns a language's table ("zh" finds "zh-cn"), DefaultLang keys
// filling what it lacks.
func (t *Text) lookup(lang string) map[string]string {
	if t == nil {
		return map[string]string{}
	}
	base := t.ByLang[DefaultLang]
	lang = strings.ToLower(lang)
	m, ok := t.ByLang[lang]
	if !ok {
		want, _, _ := strings.Cut(lang, "-")
		for code, tbl := range t.ByLang {
			if b, _, _ := strings.Cut(code, "-"); b == want {
				m, ok = tbl, true
				break
			}
		}
	}
	if !ok {
		return base
	}
	out := make(map[string]string, len(base))
	for k, v := range base {
		out[k] = v
	}
	for k, v := range m {
		out[k] = v
	}
	return out
}

// TextHandler is GET /v1/story/text?lang= (public, ETag + cache).
func TextHandler(w http.ResponseWriter, r *http.Request) {
	t := load().text
	lang := r.URL.Query().Get("lang")
	if lang == "" {
		lang = DefaultLang
	}
	version := ""
	if t != nil {
		version = t.Version
	}
	etag := `"` + version + "-" + strings.ToLower(lang) + `"`
	w.Header().Set("ETag", etag)
	w.Header().Set("Cache-Control", "public, max-age=300")
	if r.Header.Get("If-None-Match") == etag {
		w.WriteHeader(http.StatusNotModified)
		return
	}
	apierr.WriteJSON(w, http.StatusOK, map[string]any{"version": version, "text": t.lookup(lang)})
}
