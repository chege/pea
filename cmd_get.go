package main

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"pea/platform"
)

func readEntry(store, name string) ([]byte, error) {
	name, err := normalizeName(name)
	if err != nil {
		return nil, err
	}
	path, _, err := existingEntryPath(store, name)
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
	return stripFrontMatter(b), nil
}

// stripFrontMatter removes simple YAML front matter delimited by lines starting with '---'.
func stripFrontMatter(b []byte) []byte {
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

func isTTY() bool {
	fi, err := os.Stdout.Stat()
	if err != nil {
		return false
	}

	return (fi.Mode() & os.ModeCharDevice) != 0
}

func copyToClipboard(s string) error {
	// Use platform clipboard abstraction
	if err := platform.ClipboardImpl.Init(); err != nil {
		return err
	}

	return platform.ClipboardImpl.WriteText(s)
}
