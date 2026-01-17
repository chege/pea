package app

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

// NormalizeName converts input to snake_case and validates it is non-empty and not a reserved command.
func NormalizeName(raw string) (string, error) {
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

func FileExists(p string) bool {
	info, err := os.Stat(p)
	return err == nil && !info.IsDir()
}

var snakeRe = regexp.MustCompile(`[^a-z0-9_]+`)

func toSnake(s string) string {
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, " ", "_")
	s = snakeRe.ReplaceAllString(s, "")
	return s
}
