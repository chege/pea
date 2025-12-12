package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

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
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	base := filepath.Join(home, ".pea")
	store := filepath.Join(base, "prompts")
	if err := os.MkdirAll(store, 0o755); err != nil {
		return "", err
	}
	// Future: create default config at ~/.pea/config.toml if missing
	return store, nil
}

func listEntries(store string) ([]string, error) {
	files, err := os.ReadDir(store)
	if err != nil {
		return nil, err
	}
	var names []string
	for _, f := range files {
		if f.IsDir() { continue }
		name := f.Name()
		if strings.HasSuffix(name, ".txt") {
			name = strings.TrimSuffix(name, ".txt")
			names = append(names, name)
		}
	}
	sort.Strings(names)
	return names, nil
}
