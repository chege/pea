package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func readEntry(store, name string) ([]byte, error) {
	name = toSnake(name)
	path := filepath.Join(store, name+".txt")
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("not found: %s", name)
	}
	return b, nil
}
