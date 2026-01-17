package main

import (
	"fmt"
	"os"
)

// normalizeName converts input to snake_case and validates it is non-empty and not a reserved command.
func normalizeName(raw string) (string, error) {
	name := toSnake(raw)
	if name == "" {
		return "", fmt.Errorf("invalid name %q: use letters, numbers, or underscores", raw)
	}
	if isReserved(name) {
		return "", fmt.Errorf("invalid name %q: %q is a reserved command", raw, name)
	}
	return name, nil
}

func isReserved(name string) bool {
	switch name {
	case "add", "ls", "rm", "mv", "history", "search", "completion", "help":
		return true
	}
	return false
}

func fileExists(p string) bool {
	info, err := os.Stat(p)
	return err == nil && !info.IsDir()
}
