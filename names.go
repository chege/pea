package main

import (
	"fmt"
	"os"
)

// normalizeName converts input to snake_case and validates it is non-empty.
func normalizeName(raw string) (string, error) {
	name := toSnake(raw)
	if name == "" {
		return "", fmt.Errorf("invalid name %q: use letters, numbers, or underscores", raw)
	}
	return name, nil
}

func fileExists(p string) bool {
	info, err := os.Stat(p)
	return err == nil && !info.IsDir()
}
