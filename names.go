package main

import "fmt"

// normalizeName converts input to snake_case and validates it is non-empty.
func normalizeName(raw string) (string, error) {
	name := toSnake(raw)
	if name == "" {
		return "", fmt.Errorf("invalid name %q: use letters, numbers, or underscores", raw)
	}
	return name, nil
}
