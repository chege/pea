package app

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
)

// ListEntries returns a sorted list of entry names in the store.
func ListEntries(store string) ([]string, error) {
	files, err := os.ReadDir(store)
	if err != nil {
		return nil, err
	}
	nameSet := make(map[string]struct{})
	for _, f := range files {
		if f.IsDir() {
			continue
		}
		name := f.Name()
		switch {
		case strings.HasSuffix(name, DefaultExt):
			name = strings.TrimSuffix(name, DefaultExt)
			nameSet[name] = struct{}{}
		case strings.HasSuffix(name, LegacyExt):
			name = strings.TrimSuffix(name, LegacyExt)
			if _, exists := nameSet[name]; !exists {
				nameSet[name] = struct{}{}
			}
		}
	}
	var names []string
	for n := range nameSet {
		names = append(names, n)
	}
	sort.Strings(names)
	return names, nil
}

// ReadEntry reads the content of an entry, optionally at a specific git revision.
// It strips YAML front matter.
func ReadEntry(store, name, rev string) ([]byte, error) {
	name, err := NormalizeName(name)
	if err != nil {
		return nil, err
	}

	if rev == "" {
		path, _, err := ExistingEntryPath(store, name)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				return nil, fmt.Errorf("not found: %s", name)
			}
			return nil, err
		}
		b, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("not found: %s", name)
		}
		return StripFrontMatter(b), nil
	}

	if b, err := ShowAtRef(store, rev, name+DefaultExt); err == nil {
		return StripFrontMatter(b), nil
	}
	if b, err := ShowAtRef(store, rev, name+LegacyExt); err == nil {
		return StripFrontMatter(b), nil
	}
	return nil, fmt.Errorf("not found in ref %s: %s", rev, name)
}

func ShowAtRef(store, rev, path string) ([]byte, error) {
	c := exec.Command("git", "show", fmt.Sprintf("%s:%s", rev, path))
	c.Dir = store
	out, err := c.CombinedOutput()
	if err != nil {
		return nil, err
	}
	return out, nil
}

// StripFrontMatter removes simple YAML front matter delimited by lines starting with '---'.
func StripFrontMatter(b []byte) []byte {
	lines := strings.Split(string(b), "\n")
	if len(lines) > 0 && strings.TrimSpace(lines[0]) == "---" {
		// find closing '---'
		for i := 1; i < len(lines); i++ {
			if strings.TrimSpace(lines[i]) == "---" {
				return []byte(strings.Join(lines[i+1:], "\n"))
			}
		}
	}

	return b
}

type EntryWithTags struct {
	Name    string
	Tags    []string
	Content string
}

func CollectEntriesWithTags(store string) ([]EntryWithTags, error) {
	files, err := os.ReadDir(store)
	if err != nil {
		return nil, err
	}
	var entries []EntryWithTags
	seen := make(map[string]struct{})
	for _, f := range files {
		if f.IsDir() {
			continue
		}
		name := f.Name()
		var base string
		switch {
		case strings.HasSuffix(name, DefaultExt):
			base = strings.TrimSuffix(name, DefaultExt)
		case strings.HasSuffix(name, LegacyExt):
			base = strings.TrimSuffix(name, LegacyExt)
		default:
			continue
		}
		if _, ok := seen[base]; ok {
			continue
		}
		seen[base] = struct{}{}
		path := filepath.Join(store, name)
		b, err := os.ReadFile(path)
		if err != nil {
			continue
		}
		entries = append(entries, EntryWithTags{
			Name:    base,
			Tags:    parseTags(b),
			Content: string(b),
		})
	}
	return entries, nil
}

func HasAllTags(entryTags, required []string) bool {
	if len(required) == 0 {
		return true
	}
	set := make(map[string]struct{}, len(entryTags))
	for _, t := range entryTags {
		set[strings.ToLower(strings.TrimSpace(t))] = struct{}{}
	}
	for _, r := range required {
		if _, ok := set[strings.ToLower(strings.TrimSpace(r))]; !ok {
			return false
		}
	}
	return true
}

func parseTags(b []byte) []string {
	lines := strings.Split(string(b), "\n")
	inFrontMatter := false
	var tags []string
	for i := 0; i < len(lines); i++ {
		trimmed := strings.TrimSpace(lines[i])
		if trimmed == "---" {
			if !inFrontMatter {
				inFrontMatter = true
				continue
			}
			break
		}
		if !inFrontMatter {
			continue
		}
		if strings.HasPrefix(trimmed, "tags:") {
			value := strings.TrimSpace(strings.TrimPrefix(trimmed, "tags:"))
			if value != "" {
				tags = append(tags, splitInlineTagList(value)...)
			} else {
				for j := i + 1; j < len(lines); j++ {
					next := strings.TrimSpace(lines[j])
					if next == "" {
						continue
					}
					if strings.HasPrefix(next, "-") {
						tag := normalizeTag(strings.TrimSpace(strings.TrimPrefix(next, "-")))
						if tag != "" {
							tags = append(tags, tag)
						}
						continue
					}
					break
				}
			}
			break
		}
	}
	return uniqueTags(tags)
}

func splitInlineTagList(raw string) []string {
	if strings.HasPrefix(raw, "[") && strings.HasSuffix(raw, "]") {
		raw = strings.TrimSpace(raw[1 : len(raw)-1])
	}
	if raw == "" {
		return nil
	}
	parts := strings.Split(raw, ",")
	var out []string
	for _, part := range parts {
		if tag := normalizeTag(part); tag != "" {
			out = append(out, tag)
		}
	}
	return out
}

func normalizeTag(raw string) string {
	tag := strings.TrimSpace(raw)
	tag = strings.Trim(tag, "'\"")
	return strings.TrimSpace(tag)
}

func uniqueTags(tags []string) []string {
	seen := make(map[string]struct{}, len(tags))
	var out []string
	for _, t := range tags {
		if t == "" {
			continue
		}
		lower := strings.ToLower(t)
		if _, ok := seen[lower]; ok {
			continue
		}
		seen[lower] = struct{}{}
		out = append(out, t)
	}
	return out
}
