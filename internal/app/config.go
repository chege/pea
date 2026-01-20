package app

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

var ConfiguredRemote string

func DefaultPaths() (base string, store string) {
	home, err := os.UserHomeDir()
	if err != nil {
		return filepath.Join("~", ".pea"), filepath.Join("~", ".pea", "prompts")
	}
	base = filepath.Join(home, ".pea")
	store = filepath.Join(base, "prompts")
	return base, store
}

func EnsureStore() (string, error) {
	if v := os.Getenv("PEA_STORE"); v != "" {
		ConfiguredRemote = ""
		return prepareStore(v, "PEA_STORE", "")
	}

	base, defaultStore := DefaultPaths()

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

	store, remote, err := loadConfig(cfgPath, defaultStore)
	if err != nil {
		return "", err
	}
	ConfiguredRemote = remote

	return prepareStore(store, "config", remote)
}

type Config struct {
	StoreDir  string `toml:"store_dir,omitempty"`
	RemoteURL string `toml:"remote_url,omitempty"`
}

func loadConfig(cfgPath, defaultStore string) (string, string, error) {
	var conf Config

	if _, err := toml.DecodeFile(cfgPath, &conf); err != nil {
		return "", "", fmt.Errorf("invalid config %s: %w", cfgPath, err)
	}

	store := defaultStore
	if conf.StoreDir != "" {
		store = conf.StoreDir
	}

	if !filepath.IsAbs(store) {
		return "", "", fmt.Errorf("invalid config %s: store_dir must be an absolute path, got %q", cfgPath, store)
	}

	remote := conf.RemoteURL
	return store, remote, nil
}

func UpdateRemoteURL(url string) error {
	base, _ := DefaultPaths()
	cfgPath := filepath.Join(base, "config.toml")

	if err := os.MkdirAll(base, 0o755); err != nil {
		return fmt.Errorf("failed to create base dir %s: %w", base, err)
	}

	var conf Config
	if _, err := os.Stat(cfgPath); err == nil {
		if _, err := toml.DecodeFile(cfgPath, &conf); err != nil {
			return fmt.Errorf("invalid config %s: %w", cfgPath, err)
		}
	}

	conf.RemoteURL = url

	f, err := os.Create(cfgPath)
	if err != nil {
		return fmt.Errorf("failed to open config file %s: %w", cfgPath, err)
	}
	defer f.Close()

	if err := toml.NewEncoder(f).Encode(conf); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}
	return nil
}

func SetGitRemote(store, remote string) error {
	// Try to set-url first (if it exists)
	cmd := exec.Command("git", "remote", "set-url", "origin", remote)
	cmd.Dir = store
	if err := cmd.Run(); err == nil {
		return nil
	}

	// If set-url failed, try adding it
	add := exec.Command("git", "remote", "add", "origin", remote)
	add.Dir = store
	if out, err := add.CombinedOutput(); err != nil {
		return fmt.Errorf("git remote add/set-url failed: %v: %s", err, string(out))
	}
	return nil
}

func prepareStore(store, source, remote string) (string, error) {
	if !filepath.IsAbs(store) {
		return "", fmt.Errorf("%s must be an absolute path, got %q", source, store)
	}

	if err := os.MkdirAll(store, 0o755); err != nil {
		return "", fmt.Errorf("failed to create store dir %s: %w", store, err)
	}

	if err := ensureGitRepo(store, remote); err != nil {
		return "", err
	}

	return store, nil
}

func ensureGitRepo(store, remote string) error {
	if _, err := os.Stat(filepath.Join(store, ".git")); err == nil {
		// configure remote if provided and not set
		if remote != "" {
			if err := setRemoteIfMissing(store, remote); err != nil {
				return err
			}
		}
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

	if remote != "" {
		if err := setRemoteIfMissing(store, remote); err != nil {
			return err
		}
	}

	return nil
}

func setRemoteIfMissing(store, remote string) error {
	remoteCmd := exec.Command("git", "remote", "get-url", "origin")
	remoteCmd.Dir = store
	if err := remoteCmd.Run(); err == nil {
		return nil
	}
	add := exec.Command("git", "remote", "add", "origin", remote)
	add.Dir = store
	if out, err := add.CombinedOutput(); err != nil {
		return fmt.Errorf("git remote add failed: %v: %s", err, string(out))
	}
	return nil
}
