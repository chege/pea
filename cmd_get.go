package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"pea/platform"
)

func readEntry(store, name string) ([]byte, error) {
	name = toSnake(name)
	path := filepath.Join(store, name+".txt")
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
