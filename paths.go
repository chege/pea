package main

import (
	"errors"
	"os"
	"path/filepath"
)

const (
	defaultExt = ".md"
	legacyExt  = ".txt"
)

func defaultEntryPath(store, name string) string {
	return filepath.Join(store, name+defaultExt)
}

func legacyEntryPath(store, name string) string {
	return filepath.Join(store, name+legacyExt)
}

// targetEntryPath returns the path to use for creating or updating an entry,
// preferring an existing entry (either .md or .txt) and defaulting to .md.
func targetEntryPath(store, name string) (string, string, error) {
	if p, ext, err := existingEntryPath(store, name); err == nil {
		return p, ext, nil
	} else if !errors.Is(err, os.ErrNotExist) {
		return "", "", err
	}
	return defaultEntryPath(store, name), defaultExt, nil
}

// existingEntryPath finds an existing entry path, preferring .md over legacy .txt.
func existingEntryPath(store, name string) (string, string, error) {
	p := defaultEntryPath(store, name)
	if _, err := os.Stat(p); err == nil {
		return p, defaultExt, nil
	} else if !errors.Is(err, os.ErrNotExist) {
		return "", "", err
	}

	lp := legacyEntryPath(store, name)
	if _, err := os.Stat(lp); err == nil {
		return lp, legacyExt, nil
	} else if !errors.Is(err, os.ErrNotExist) {
		return "", "", err
	}

	return "", "", os.ErrNotExist
}
