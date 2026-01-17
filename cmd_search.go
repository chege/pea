package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

func addSearchCommand(root *cobra.Command) {
	var tags []string
	cmd := &cobra.Command{
		Use:   "search <query>",
		Short: "search entries by name substring or tags",
		Args:  cobra.RangeArgs(0, 1),
		RunE: func(cmd *cobra.Command, args []string) error {
			store, err := ensureStore()
			if err != nil {
				return err
			}
			query := ""
			if len(args) == 1 {
				query = strings.ToLower(args[0])
			}

			entries, err := collectEntriesWithTags(store)
			if err != nil {
				return err
			}

			var out []string
			for _, e := range entries {
				match := false
				if query == "" {
					match = true
				} else {
					if strings.Contains(strings.ToLower(e.name), query) {
						match = true
					} else if strings.Contains(strings.ToLower(e.content), query) {
						match = true
					}
				}

				if match && len(tags) > 0 && !hasAllTags(e.tags, tags) {
					match = false
				}

				if match {
					out = append(out, e.name)
				}
			}

			for _, name := range out {
				if _, err := fmt.Fprintln(cmd.OutOrStdout(), name); err != nil {
					return err
				}
			}
			return nil
		},
	}
	cmd.Flags().StringArrayVar(&tags, "tag", nil, "filter by tag (repeatable)")
	root.AddCommand(cmd)
}

type entryWithTags struct {
	name    string
	tags    []string
	content string
}

func collectEntriesWithTags(store string) ([]entryWithTags, error) {
	files, err := os.ReadDir(store)
	if err != nil {
		return nil, err
	}
	var entries []entryWithTags
	seen := make(map[string]struct{})
	for _, f := range files {
		if f.IsDir() {
			continue
		}
		name := f.Name()
		var base string
		switch {
		case strings.HasSuffix(name, defaultExt):
			base = strings.TrimSuffix(name, defaultExt)
		case strings.HasSuffix(name, legacyExt):
			base = strings.TrimSuffix(name, legacyExt)
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
		entries = append(entries, entryWithTags{
			name:    base,
			tags:    parseTags(b),
			content: string(b),
		})
	}
	return entries, nil
}

func hasAllTags(entryTags, required []string) bool {
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
	tag = strings.Trim(tag, "\"'")
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
