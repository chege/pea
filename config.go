package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

func defaultPaths() (base string, store string) {
	home, err := os.UserHomeDir()
	if err != nil {
		return filepath.Join("~", ".pea"), filepath.Join("~", ".pea", "prompts")
	}
	base = filepath.Join(home, ".pea")
	store = filepath.Join(base, "prompts")
	return base, store
}

func ensureStore() (string, error) {
	if v := os.Getenv("PEA_STORE"); v != "" {
		return prepareStore(v, "PEA_STORE")
	}

	base, defaultStore := defaultPaths()

	if err := os.MkdirAll(base, 0o755); err != nil {
		return "", fmt.Errorf("failed to create base dir %s: %w", base, err)
	}

	cfgPath := filepath.Join(base, "config.toml")
	if _, err := os.Stat(cfgPath); errors.Is(err, os.ErrNotExist) {
		content := "# pea config\n# store_dir = \"" + defaultStore + "\"\n"
		if err := os.WriteFile(cfgPath, []byte(content), 0o644); err != nil {
			return "", fmt.Errorf("failed to write default config %s: %w", cfgPath, err)
		}
	}

	store, err := loadStoreDir(cfgPath, defaultStore)
	if err != nil {
		return "", err
	}

	return prepareStore(store, "config")
}

func loadStoreDir(cfgPath, defaultStore string) (string, error) {
	var conf struct {
		StoreDir string `toml:"store_dir"`
	}

	if _, err := toml.DecodeFile(cfgPath, &conf); err != nil {
		return "", fmt.Errorf("invalid config %s: %w", cfgPath, err)
	}

	store := defaultStore
	if conf.StoreDir != "" {
		store = conf.StoreDir
	}

	if !filepath.IsAbs(store) {
		return "", fmt.Errorf("invalid config %s: store_dir must be an absolute path, got %q", cfgPath, store)
	}

	return store, nil
}

func prepareStore(store, source string) (string, error) {
	if !filepath.IsAbs(store) {
		return "", fmt.Errorf("%s must be an absolute path, got %q", source, store)
	}

	if err := os.MkdirAll(store, 0o755); err != nil {
		return "", fmt.Errorf("failed to create store dir %s: %w", store, err)
	}

	if err := ensureGitRepo(store); err != nil {
		return "", err
	}

	return store, nil
}

func ensureGitRepo(store string) error {
	if _, err := os.Stat(filepath.Join(store, ".git")); err == nil {
		return nil
	} else if !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("failed to check git repo: %w", err)
	}

	cmds := [][]string{
		{"init"},
		{"config", "user.name", "pea"},
		{"config", "user.email", "pea@example.com"},
	}

	for _, args := range cmds {
		c := exec.Command("git", args...)
		c.Dir = store
		if out, err := c.CombinedOutput(); err != nil {
			return fmt.Errorf("git %v failed: %v: %s", args, err, string(out))
		}
	}

	return nil
}
