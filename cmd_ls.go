package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/spf13/cobra"
)

func addListCommand(root *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "ls",
		Short: "list stored entries",
		RunE: func(cmd *cobra.Command, args []string) error {
			store, err := ensureStore()
			if err != nil {
				return err
			}
			entries, err := listEntries(store)
			if err != nil {
				return err
			}
			for _, e := range entries {
				fmt.Fprintln(cmd.OutOrStdout(), e)
			}
			return nil
		},
	}
	root.AddCommand(cmd)
}

func ensureStore() (string, error) {
	// Env override
	if v := os.Getenv("PEA_STORE"); v != "" {
		if err := os.MkdirAll(v, 0o755); err != nil {
			return "", err
		}
		if _, err := os.Stat(filepath.Join(v, ".git")); os.IsNotExist(err) {
			_ = exec.Command("bash", "-c", "cd '"+v+"' && git init").Run()
			_ = exec.Command("bash", "-c", "cd '"+v+"' && git config user.name 'pea'").Run()
			_ = exec.Command("bash", "-c", "cd '"+v+"' && git config user.email 'pea@example.com'").Run()
		}
		return v, nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	base := filepath.Join(home, ".pea")
	store := filepath.Join(base, "prompts")
	// Ensure base and store
	if err := os.MkdirAll(store, 0o755); err != nil {
		return "", err
	}
	// Create default config if missing
	cfg := filepath.Join(base, "config.toml")
	if _, err := os.Stat(cfg); os.IsNotExist(err) {
		_ = os.WriteFile(cfg, []byte("# pea config\n# store_dir = \""+store+"\"\n"), 0o644)
	}
	// Read config.toml for store_dir using TOML parser if available
	var conf struct{ StoreDir string `toml:"store_dir"` }
	if _, err := toml.DecodeFile(cfg, &conf); err == nil {
		if conf.StoreDir != "" {
			store = conf.StoreDir
			_ = os.MkdirAll(store, 0o755)
		}
	}
	// Initialize git repo if missing (for versioning of entries)
	if _, err := os.Stat(filepath.Join(store, ".git")); os.IsNotExist(err) {
		_ = exec.Command("bash", "-c", "cd '"+store+"' && git init").Run()
		_ = exec.Command("bash", "-c", "cd '"+store+"' && git config user.name 'pea'").Run()
		_ = exec.Command("bash", "-c", "cd '"+store+"' && git config user.email 'pea@example.com'").Run()
	}
	return store, nil
}

func parseStoreDir(s string) string {
	for _, line := range strings.Split(s, "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "store_dir") {
			parts := strings.SplitN(line, "=", 2)
			if len(parts) == 2 {
				v := strings.TrimSpace(parts[1])
				v = strings.Trim(v, "\"'")
				if v != "" {
					return v
				}
			}
		}
	}
	return ""
}

func listEntries(store string) ([]string, error) {
	files, err := os.ReadDir(store)
	if err != nil {
		return nil, err
	}
	var names []string
	for _, f := range files {
		if f.IsDir() {
			continue
		}
		name := f.Name()
		if strings.HasSuffix(name, ".txt") {
			name = strings.TrimSuffix(name, ".txt")
			names = append(names, name)
		}
	}
	sort.Strings(names)
	return names, nil
}
