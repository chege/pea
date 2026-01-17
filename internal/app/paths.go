package app

import (
	"errors"
	"os"
	"path/filepath"
)

const (
	DefaultExt = ".md"
	LegacyExt  = ".txt"
)

func DefaultEntryPath(store, name string) string {
	return filepath.Join(store, name+DefaultExt)
}

func LegacyEntryPath(store, name string) string {
	return filepath.Join(store, name+LegacyExt)
}

// TargetEntryPath returns the path to use for creating or updating an entry,
// preferring an existing entry (either .md or .txt) and defaulting to .md.
func TargetEntryPath(store, name string) (string, string, error) {
	if p, ext, err := ExistingEntryPath(store, name); err == nil {
		return p, ext, nil
	} else if !errors.Is(err, os.ErrNotExist) {
		return "", "", err
	}
	return DefaultEntryPath(store, name), DefaultExt, nil
}

// ExistingEntryPath finds an existing entry path, preferring .md over legacy .txt.
func ExistingEntryPath(store, name string) (string, string, error) {
	p := DefaultEntryPath(store, name)
	if _, err := os.Stat(p); err == nil {
		return p, DefaultExt, nil
	} else if !errors.Is(err, os.ErrNotExist) {
		return "", "", err
	}

	lp := LegacyEntryPath(store, name)
	if _, err := os.Stat(lp); err == nil {
		return lp, LegacyExt, nil
	} else if !errors.Is(err, os.ErrNotExist) {
		return "", "", err
	}

	return "", "", os.ErrNotExist
}
